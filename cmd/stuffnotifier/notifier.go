package main

import (
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/urfave/cli/v2"

	fapoller "github.com/jalavosus/stuffnotifier/internal/pollers/flightawarepoller"
	"github.com/jalavosus/stuffnotifier/internal/pollers/poller"
	"github.com/jalavosus/stuffnotifier/pkg/flightaware"
)

var (
	geminiCmd = cli.Command{
		Name:        "gemini",
		Description: "Monitor crypto prices via Gemini",
		Usage:       "stuffnotifier monitor gemini [stuff]",
		ArgsUsage:   "",
		Action:      nil,
		Flags: []cli.Flag{
			&geminiConfigFlag,
			&geminiApiKeyFlag,
			&geminiApiSecret,
		},
	}
	flightawareCmd = cli.Command{
		Name:        "flightaware",
		Description: "Monitor flightaware flights",
		ArgsUsage: "FLIGHT_IDENTIFIER: (Flight number OR FlightAware flight ID for a flight) " +
			"IDENTIFIER_TYPE: ('flight_number' if passing a flight number, 'flightaware_id' if passing a FlightAware ID number)",
		Action: flightawareCmdAction,
		Flags: []cli.Flag{
			&pollerConfigFlag,
			&flightawareConfigFlag,
			&flightAwareApiKeyFlag,
		},
	}
)

func flightawareCmdAction(c *cli.Context) error {
	if c.Args().Len() < 2 {
		return errors.New("not enough arguments passed")
	}

	flightId := c.Args().First()
	flightIdTypeStr := c.Args().Get(1)

	var flightIdType flightaware.IdentifierType

	switch strings.ToLower(flightIdTypeStr) {
	case "flight_number":
		flightIdType = flightaware.DesignatorIdent
	case "flightaware_id":
		flightIdType = flightaware.FaFlightIdIdent
	default:
		return errors.Errorf("unknown identifier type %[1]s. Allowed values: 'flight_number', 'flightaware_id'", flightIdType)
	}

	var (
		pollInterval time.Duration
		config       poller.Config
		faConfig     *flightaware.Config
	)

	pollerConfig, confErr := loadPollerConfig(c)
	if confErr != nil {
		logger.Warn("error loading Poller config", zap.Error(confErr))
	} else {
		if pollerConfig != nil {
			config = *pollerConfig
		} else {
			config = poller.Config{}
		}
	}

	if config.Twilio == nil {
		twilioConf, confErr := loadTwilioConfig(c)
		if confErr != nil {
			logger.Warn("error loading Twilio config", zap.Error(confErr))
		} else if twilioConf != nil {
			config.Twilio = twilioConf
		}
	}

	if config.Discord == nil {
		discordConf, confErr := loadDiscordConfig(c)
		if confErr != nil {
			logger.Warn("error loading Discord config", zap.Error(confErr))
		} else if discordConf != nil {
			config.Discord = discordConf
		}
	}

	if config.FlightAware == nil {
		faConf, confErr := loadFlightAwareConfig(c)
		if confErr != nil {
			logger.Warn("error loading FlightAware config", zap.Error(confErr))
			pollInterval = flightaware.DefaultPollInterval
		} else if faConf != nil {
			config.FlightAware = faConf
			faConfig = faConf
			pollInterval = faConfig.PollInterval
		}
	} else {
		faConfig = config.FlightAware
	}

	config.PollInterval = pollInterval

	faPoller, err := fapoller.NewFlightAwarePoller(config, faConfig)
	if err != nil {
		return err
	}

	return faPoller.Start(c.Context, flightId, flightIdType)
}
