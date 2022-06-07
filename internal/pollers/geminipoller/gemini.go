package geminipoller

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"github.com/jalavosus/stuffnotifier/internal/datastore"
	"github.com/jalavosus/stuffnotifier/internal/messages"
	"github.com/jalavosus/stuffnotifier/internal/pollers/poller"
	"github.com/jalavosus/stuffnotifier/pkg/authdata"
	"github.com/jalavosus/stuffnotifier/pkg/gemini"
)

type Poller struct {
	datastore datastore.Datastore[CacheEntry]
	*poller.BasePoller
	geminiClient *gemini.Client
	geminiConfig gemini.Config
}

func NewPoller(conf poller.Config, geminiConfig *gemini.Config) (*Poller, error) {
	var (
		geminiConf   = gemini.DefaultConfig()
		datastoreErr error
	)

	switch {
	case conf.Gemini != nil:
		geminiConf = conf.Gemini
	case geminiConfig != nil:
		geminiConf = geminiConfig
	}

	p := &Poller{
		BasePoller:   poller.NewBasePoller(conf),
		geminiConfig: *geminiConf,
	}

	p.SetPollInterval(geminiConf.PollInterval)

	p.datastore, datastoreErr = datastore.NewDatastore[CacheEntry](conf.Cache)
	if datastoreErr != nil {
		return nil, datastoreErr
	}

	return p, nil
}

func (p *Poller) GeminiConfig() gemini.Config {
	return p.geminiConfig
}

func (p *Poller) Datastore() datastore.Datastore[CacheEntry] {
	return p.datastore
}

func (p *Poller) GeminiClient() *gemini.Client {
	return p.geminiClient
}

func (p *Poller) initGeminiClient(authData authdata.AuthData) {
	p.geminiClient = gemini.NewClient(authData)
}

func (p *Poller) Start(ctx context.Context) error {
	var authData authdata.AuthData
	if conf := p.GeminiConfig(); conf.Auth != nil {
		authData = conf.Auth
	} else {
		ad, authDataErr := authdata.GeminiAPIAuth()
		if authDataErr != nil {
			return authDataErr
		}

		authData = ad
	}

	p.initGeminiClient(authData)

	errCh := make(chan error, 1)

	for _, pollerConf := range p.GeminiConfig().Notifications.SpotPrice {
		cacheKey := "gemini:" + pollerConf.CurrencySymbol()
		concurrentParams := poller.NewConcurrentParamsWithChannel(authData, cacheKey, errCh)

		go p.pollSpotPrices(
			ctx,
			pollerConf,
			concurrentParams,
		)
	}

	return <-errCh
}

func (p *Poller) pollSpotPrices(
	ctx context.Context,
	pollerConf gemini.SpotPriceNotificationsConfig,
	pollerParams *poller.ConcurrentParams,
) {

	symbol := pollerConf.CurrencySymbol()

	baseCurrency := pollerConf.BaseCurrency
	quoteCurrency := pollerConf.QuoteCurrency
	baseAmt := pollerConf.BaseAmt()

	ticker := time.NewTicker(p.PollInterval())
	cleanup := func(err error) {
		pollerParams.Cleanup(err, ticker)
	}

	p.LogInfo(
		"starting gemini spot price poller",
		zap.String("symbol", symbol),
		zap.String("check_interval", p.PollInterval().String()),
	)

	for {
		select {
		case <-ctx.Done():
			cleanup(nil)
			return
		case t := <-ticker.C:
			spotPrice, spotPriceErr := p.fetchSpotPrice(ctx, symbol)
			if spotPriceErr != nil {
				p.LogError("error fetching spot price", zap.String("symbol", symbol), zap.Error(spotPriceErr))
				continue
			}

			spotPrice = spotPrice.Mul(baseAmt)

			msg := messages.SpotPriceAlert{
				EventTime:     t,
				SpotPrice:     spotPrice,
				BaseAmount:    baseAmt,
				BaseCurrency:  baseCurrency,
				QuoteCurrency: quoteCurrency,
			}

			if err := p.SendMessage(ctx, msg); err != nil {
				p.LogError("error sending notification", zap.Error(err))
			}
		}
	}
}

func (p *Poller) fetchSpotPrice(ctx context.Context, symbol string) (decimal.Decimal, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	spotPrice, spotPriceErr := p.GeminiClient().Ticker(ctx, symbol)
	if spotPriceErr != nil {
		return decimal.Zero, spotPriceErr
	}

	return spotPrice.Last, nil
}
