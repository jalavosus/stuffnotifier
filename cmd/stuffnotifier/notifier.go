package main

import (
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/urfave/cli/v2"

	fapoller "github.com/jalavosus/stuffnotifier/internal/pollers/flightawarepoller"
	"github.com/jalavosus/stuffnotifier/internal/pollers/geminipoller"
	"github.com/jalavosus/stuffnotifier/pkg/flightaware"
	"github.com/jalavosus/stuffnotifier/pkg/gemini"
)

var (
	geminiCmd = cli.Command{
		Name:        "gemini",
		Description: "Monitor crypto prices via Gemini",
		Usage:       "stuffnotifier monitor gemini [stuff]",
		Action:      geminiCmdAction,
		Flags: []cli.Flag{
			&pollerConfigFlag,
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

func geminiCmdAction(c *cli.Context) error {
	var (
		geminiConfig *gemini.Config
		pollInterval = gemini.DefaultPollInterval
	)

	config := loadBuildPollerConfig(c)

	if config.Gemini == nil {
		geminiConf, confErr := loadGeminiConfig(c)
		if confErr != nil {
			logger.Warn("error loading Gemini config", zap.Error(confErr))
		} else if geminiConf != nil {
			config.Gemini = geminiConf
			geminiConfig = geminiConf
			pollInterval = geminiConfig.PollInterval
		}
	} else {
		geminiConfig = config.Gemini
		pollInterval = geminiConfig.PollInterval
	}

	config.PollInterval = pollInterval

	poller, err := geminipoller.NewPoller(config, geminiConfig)
	if err != nil {
		return err
	}

	return poller.Start(c.Context)
}

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
		faConfig     *flightaware.Config
	)

	config := loadBuildPollerConfig(c)

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

	faPoller, err := fapoller.NewPoller(config, faConfig)
	if err != nil {
		return err
	}

	return faPoller.Start(c.Context, flightId, flightIdType)
}
