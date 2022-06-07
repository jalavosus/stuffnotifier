package gemini

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/jalavosus/stuffnotifier/internal/utils"
	"github.com/jalavosus/stuffnotifier/pkg/errs"
)

const (
	DefaultPollInterval = 30 * time.Second
)

type Config struct {
	PollInterval  time.Duration       `json:"poll_interval" yaml:"poll_interval" toml:"PollInterval"`
	Auth          *AuthConfig         `json:"auth,omitempty" yaml:"auth,omitempty" toml:"Auth"`
	Notifications NotificationsConfig `json:"notifications" yaml:"notifications" toml:"Notifications"`
}

type AuthConfig struct {
	ApiKey    string `json:"api_key" yaml:"api_key"`
	ApiSecret string `json:"api_secret" yaml:"api_secret"`
}

func (c AuthConfig) Account() string {
	return ""
}

func (c AuthConfig) Key() string {
	return c.ApiKey
}

func (c AuthConfig) Secret() string {
	return c.ApiSecret
}

type NotificationsConfig struct {
	SpotPrice []SpotPriceNotificationsConfig `json:"spot_price,omitempty" yaml:"spot_price,omitempty" toml:"SpotPrice,omitempty"`
}

type SpotPriceNotificationsConfig struct {
	BaseCurrency  string           `json:"base_currency" yaml:"base_currency" toml:"BaseCurrency"`
	QuoteCurrency string           `json:"quote_currency" yaml:"quote_currency" toml:"QuoteCurrency"`
	Symbol        *string          `json:"symbol,omitempty" yaml:"symbol,omitempty" toml:"Symbol,omitempty"`
	BaseAmount    *decimal.Decimal `json:"base_amount,omitempty" yaml:"base_amount,omitempty" toml:"BaseAmount,omitempty"`
}

func (c SpotPriceNotificationsConfig) CurrencySymbol() string {
	if symbol, ok := utils.FromPointer(c.Symbol); ok {
		return symbol
	}

	return strings.ToUpper(c.BaseCurrency) + strings.ToUpper(c.QuoteCurrency)
}

func (c SpotPriceNotificationsConfig) BaseAmt() decimal.Decimal {
	if baseAmt, ok := utils.FromPointer(c.BaseAmount); ok && !baseAmt.IsZero() {
		return baseAmt
	}

	return decimal.NewFromInt(1)
}

func DefaultConfig() *Config {
	spotPriceConf := []SpotPriceNotificationsConfig{
		{Symbol: utils.ToPointer("ETHUSD")},
	}

	notifsConf := NotificationsConfig{SpotPrice: spotPriceConf}

	return &Config{
		PollInterval:  DefaultPollInterval,
		Notifications: notifsConf,
	}
}

// LoadConfig reads configuration data from the file at the passed path
// and returns it as a fully loaded Config.
// `confPath` is expected to be an absolute file path.
// Supported file types: json, yaml, toml
func LoadConfig(confPath string) (conf Config, err error) {
	confBytes, readErr := ioutil.ReadFile(confPath)
	if readErr != nil {
		err = errs.ReadFileError(readErr, confPath)
		return
	}

	confType := utils.ConfigFileTypeFromExtension(filepath.Ext(confPath))

	if unmarshalErr := utils.UnmarshalConfig(confBytes, confType, &conf); unmarshalErr != nil {
		err = errors.WithMessagef(unmarshalErr, "error unmarshalling data from file %[1]s", confPath)
		return
	}

	return
}
