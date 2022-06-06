package discord

import (
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/jalavosus/stuffnotifier/internal/utils"
	"github.com/jalavosus/stuffnotifier/pkg/errs"
)

type Config struct {
	UserId    *string     `json:"user_id,omitempty" yaml:"user_id,omitempty" toml:"UserId,omitempty"`
	ChannelId *string     `json:"channel_id,omitempty" yaml:"channel_id,omitempty" toml:"ChannelId,omitempty"`
	Auth      *AuthConfig `json:"auth,omitempty" yaml:"auth,omitempty" toml:"Auth,omitempty"`
}

type AuthConfig struct {
	Token string `json:"token" yaml:"token" toml:"Token"`
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
