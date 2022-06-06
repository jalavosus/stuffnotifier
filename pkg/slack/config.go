package slack

import (
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/jalavosus/stuffnotifier/internal/utils"
	"github.com/jalavosus/stuffnotifier/pkg/errs"
)

type Config struct {
	Auth *AuthConfig `json:"auth,omitempty" yaml:"auth,omitempty" toml:"Auth,omitempty"`
	// Which user IDs to send notifications to (as private/direct messages)
	Users []string `json:"users" yaml:"users" toml:"Users"`
	// Which channel IDs to send notifications to
	Channels []string `json:"channels" yaml:"channels" toml:"Channels"`
}

type AuthConfig struct {
	Token string `json:"token" yaml:"token" toml:"Token"`
}

func (c AuthConfig) Account() string {
	return ""
}

func (c AuthConfig) Key() string {
	return ""
}

func (c AuthConfig) Secret() string {
	return c.Token
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
