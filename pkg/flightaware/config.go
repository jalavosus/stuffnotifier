package flightaware

import (
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/pkg/errors"

	"github.com/jalavosus/stuffnotifier/internal/utils"
	"github.com/jalavosus/stuffnotifier/pkg/errs"
)

const (
	DefaultPollInterval     = time.Minute
	DefaultPreArrivalOffset = time.Hour
)

const (
	defaultHttpTimeout = 10 * time.Second
)

// Config contains the configuration for a FlightAware poller.
type Config struct {
	Auth *AuthConfig `json:"auth,omitempty" yaml:"auth,omitempty" toml:"Auth,omitempty"`
	// Interval between flight data checks.
	// Default: 1 minute
	PollInterval time.Duration `json:"poll_interval" yaml:"poll_interval" toml:"PollInterval"`
	// Configuration for various notifications.
	Notifications NotificationsConfig `json:"notifications" yaml:"notifications" toml:"Notifications"`
}

func DefaultConfig() Config {
	return Config{
		Auth:          nil,
		PollInterval:  DefaultPollInterval,
		Notifications: DefaultNotificationsConfig(),
	}
}

// NotificationsConfig contains configuration data respective to
// which notifications will be sent for certain flight events.
type NotificationsConfig struct {
	PreDeparture  PreEventConfig `json:"pre_departure" yaml:"pre_departure" toml:"PreDeparture"`
	PreArrival    PreEventConfig `json:"pre_arrival" yaml:"pre_arrival" toml:"PreArrival"`
	Takeoff       bool           `json:"takeoff" yaml:"takeoff" toml:"Takeoff"`
	Landing       bool           `json:"landing" yaml:"landing" toml:"Landing"`
	GateArrival   bool           `json:"gate_arrival" yaml:"gate_arrival" toml:"GateArrival"`
	BaggageClaim  bool           `json:"baggage_claim" yaml:"baggage_claim" toml:"BaggageClaim"`
	GateDeparture bool           `json:"gate_departure" yaml:"gate_departure" toml:"GateDeparture"`
	UseLocalTime  bool           `json:"use_local_time" yaml:"use_local_time" toml:"UseLocalTime"`
}

// PreEventConfig contains configuration details
// for pre-arrival notifications.
type PreEventConfig struct {
	// If true, enables notifications pre-arrival.
	// Amount of time pre-arrival to send notifications can be set with
	// Offset.
	// Default: true
	Enabled bool `json:"enabled" yaml:"enabled" toml:"Enabled"`
	// If true, a notification will be sent at Offset time
	// before the estimated arrival time of a flight.
	// Default: true
	Estimated bool `json:"estimated" yaml:"estimated" toml:"Estimated"`
	// If true, a notification will be sent at Offset time
	// before the scheduled arrival time of a flight.
	// Default: false
	Scheduled bool `json:"scheduled" yaml:"scheduled" toml:"Scheduled"`
	// If Enabled is true, Offset is the amount of time
	// before both/either of estimated and scheduled arrival times
	// at which to send notifications.
	// Default: 1 hour
	Offset time.Duration `json:"offset" yaml:"offset" toml:"Offset"`
}

// DefaultNotificationsConfig returns a NotificationsConfig
// set with default values.
func DefaultNotificationsConfig() NotificationsConfig {
	return NotificationsConfig{
		GateDeparture: true,
		Takeoff:       true,
		Landing:       true,
		GateArrival:   true,
		BaggageClaim:  false,
		PreArrival: PreEventConfig{
			Enabled:   true,
			Estimated: true,
			Scheduled: false,
			Offset:    DefaultPreArrivalOffset,
		},
		PreDeparture: PreEventConfig{
			Enabled:   false,
			Estimated: false,
			Scheduled: true,
			Offset:    DefaultPreArrivalOffset,
		},
		UseLocalTime: true,
	}
}

type AuthConfig struct {
	ApiKey string `json:"api_key" yaml:"api_key" toml:"ApiKey"`
}

func (c AuthConfig) Account() string {
	return ""
}

func (c AuthConfig) Key() string {
	return c.ApiKey
}

func (c AuthConfig) Secret() string {
	return ""
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
