package twilio

import (
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/stoicturtle/stuffnotifier/internal/utils"
	"github.com/stoicturtle/stuffnotifier/pkg/errs"
)

type Config struct {
	Auth            *AuthConfig `json:"auth,omitempty" yaml:"auth,omitempty" toml:"Auth,omitempty"`
	SenderNumber    string      `json:"sender_number" yaml:"sender_number" toml:"SenderNumber"`
	RecipientNumber string      `json:"recipient_number" yaml:"recipient_number" toml:"RecipientNumber"`
}

type AuthConfig struct {
	AuthToken *AuthTokenConfig `json:"auth_token,omitempty" yaml:"auth_token,omitempty" toml:"AuthToken,omitempty"`
	ApiKey    *ApiKeyConfig    `json:"api_key,omitempty" yaml:"api_key,omitempty" toml:"ApiKey,omitempty"`
}

type AuthTokenConfig struct {
	AccountSid string `json:"account_sid" yaml:"account_sid" toml:"AccountSid"`
	Token      string `json:"token" yaml:"token" toml:"Token"`
}

func (c AuthTokenConfig) Account() string {
	return c.AccountSid
}

func (c AuthTokenConfig) Key() string {
	return c.AccountSid
}

func (c AuthTokenConfig) Secret() string {
	return c.Token
}

type ApiKeyConfig struct {
	AccountSid string `json:"account_sid" yaml:"account_sid" toml:"AccountSid"`
	ApiKey     string `json:"key" yaml:"key" toml:"Key"`
	ApiSecret  string `json:"secret" yaml:"secret" toml:"Secret"`
}

func (c ApiKeyConfig) Account() string {
	return c.AccountSid
}

func (c ApiKeyConfig) Key() string {
	return c.ApiKey
}

func (c ApiKeyConfig) Secret() string {
	return c.ApiSecret
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
