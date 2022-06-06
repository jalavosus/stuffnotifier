package poller

import (
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/pkg/errors"

	"github.com/stoicturtle/stuffnotifier/internal/datastore"
	"github.com/stoicturtle/stuffnotifier/internal/utils"
	"github.com/stoicturtle/stuffnotifier/pkg/discord"
	"github.com/stoicturtle/stuffnotifier/pkg/errs"
	"github.com/stoicturtle/stuffnotifier/pkg/flightaware"
	"github.com/stoicturtle/stuffnotifier/pkg/gemini"
	"github.com/stoicturtle/stuffnotifier/pkg/slack"
	"github.com/stoicturtle/stuffnotifier/pkg/twilio"
)

type Config struct {
	Slack        *slack.Config       `json:"slack,omitempty" yaml:"slack,omitempty" toml:"Slack,omitempty"`
	Cache        *datastore.Config   `json:"cache,omitempty" yaml:"cache,omitempty" toml:"Cache,omitempty"`
	Gemini       *gemini.Config      `json:"gemini,omitempty" yaml:"gemini,omitempty" toml:"Gemini,omitempty"`
	FlightAware  *flightaware.Config `json:"flightaware,omitempty" yaml:"flightaware,omitempty" toml:"FlightAware,omitempty"`
	Twilio       *twilio.Config      `json:"twilio,omitempty" yaml:"twilio,omitempty" toml:"Twilio,omitempty"`
	Discord      *discord.Config     `json:"discord,omitempty" yaml:"discord,omitempty" toml:"Discord,omitempty"`
	PollInterval time.Duration       `json:"poll_interval" yaml:"poll_interval" toml:"PollInterval"`
	LogStdout    bool                `json:"log_stdout" yaml:"log_stdout" toml:"LogStdout"`
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
