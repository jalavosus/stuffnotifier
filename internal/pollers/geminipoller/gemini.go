package geminipoller

import (
	"context"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"github.com/jalavosus/stuffnotifier/internal/messages"
	"github.com/jalavosus/stuffnotifier/internal/pollers/poller"
	"github.com/jalavosus/stuffnotifier/pkg/authdata"
	"github.com/jalavosus/stuffnotifier/pkg/gemini"
)

type Poller struct {
	*poller.BasePoller
	geminiConfig gemini.Config
	geminiClient *gemini.Client
}

func (p Poller) GeminiConfig() gemini.Config {
	return p.geminiConfig
}

func (p Poller) GeminiClient() *gemini.Client {
	return p.geminiClient
}

func (p Poller) initGeminiClient(authData authdata.AuthData) {
	p.geminiClient = gemini.NewClient(authData)
}

func (p *Poller) Start(ctx context.Context, baseCurrency, quoteCurrency string, baseAmt decimal.Decimal) error {
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

	cacheKey := "gemini:" + baseCurrency + "-" + quoteCurrency + "-" + baseAmt.String()
	concurrentParams := poller.NewConcurrentParams(authData, cacheKey)

	go p.pollSpotPrices(
		ctx,
		baseCurrency,
		quoteCurrency,
		baseAmt,
		concurrentParams,
	)

	return <-concurrentParams.ErrCh
}

func (p *Poller) pollSpotPrices(
	ctx context.Context,
	baseCurrency, quoteCurrency string,
	baseAmt decimal.Decimal,
	pollerParams *poller.ConcurrentParams,
) {

	symbol := strings.ToUpper(baseCurrency) + strings.ToUpper(quoteCurrency)
	// symbolData, symbolDataErr := gemini.SymbolDetails(ctx, authData, symbol)
	// if symbolDataErr != nil {
	// 	ch <- symbolDataErr
	// 	return
	// }

	ticker := time.NewTicker(p.PollInterval())
	cleanup := func(err error) {
		pollerParams.Cleanup(err, ticker)
	}

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

			msg := messages.SpotPriceAlert{
				EventTime:     t,
				SpotPrice:     spotPrice,
				BaseCurrency:  baseCurrency,
				QuoteCurrency: quoteCurrency,
			}

			if err := p.SendMessage(ctx, msg); err != nil {
				p.LogError("error sending notification", zap.Error(err))
			}
		}
	}
}

func (p Poller) fetchSpotPrice(ctx context.Context, symbol string) (decimal.Decimal, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	spotPrice, spotPriceErr := p.GeminiClient().Ticker(ctx, symbol)
	if spotPriceErr != nil {
		return decimal.Zero, spotPriceErr
	}

	return spotPrice.Last, nil
}
