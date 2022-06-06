package main

import (
	"github.com/urfave/cli/v2"

	"github.com/jalavosus/stuffnotifier/internal/pollers/poller"
	"github.com/jalavosus/stuffnotifier/pkg/discord"
	"github.com/jalavosus/stuffnotifier/pkg/flightaware"
	"github.com/jalavosus/stuffnotifier/pkg/twilio"
)

func loadPollerConfig(c *cli.Context) (*poller.Config, error) {
	return loadConfigFromFlag[poller.Config](c, pollerConfigFlag, poller.LoadConfig)
}

func loadTwilioConfig(c *cli.Context) (*twilio.Config, error) {
	return loadConfigFromFlag[twilio.Config](c, twilioConfigFlag, twilio.LoadConfig)
}

func loadDiscordConfig(c *cli.Context) (*discord.Config, error) {
	return loadConfigFromFlag[discord.Config](c, discordConfigFlag, discord.LoadConfig)
}

func loadFlightAwareConfig(c *cli.Context) (*flightaware.Config, error) {
	return loadConfigFromFlag[flightaware.Config](c, flightawareConfigFlag, flightaware.LoadConfig)
}

func loadConfigFromFlag[T any](c *cli.Context, flag cli.PathFlag, loadConf func(string) (T, error)) (*T, error) {
	var conf *T

	if confPath := flag.Get(c); confPath != "" {
		tryConf, confErr := loadConf(confPath)
		if confErr != nil {
			return nil, confErr
		} else {
			conf = &tryConf
		}
	}

	return conf, nil
}
