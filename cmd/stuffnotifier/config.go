package main

import (
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/jalavosus/stuffnotifier/internal/pollers/poller"
	"github.com/jalavosus/stuffnotifier/pkg/discord"
	"github.com/jalavosus/stuffnotifier/pkg/flightaware"
	"github.com/jalavosus/stuffnotifier/pkg/gemini"
	"github.com/jalavosus/stuffnotifier/pkg/slack"
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

func loadSlackConfig(c *cli.Context) (*slack.Config, error) {
	return loadConfigFromFlag[slack.Config](c, slackConfigFlag, slack.LoadConfig)
}

func loadFlightAwareConfig(c *cli.Context) (*flightaware.Config, error) {
	return loadConfigFromFlag[flightaware.Config](c, flightawareConfigFlag, flightaware.LoadConfig)
}

func loadGeminiConfig(c *cli.Context) (*gemini.Config, error) {
	return loadConfigFromFlag[gemini.Config](c, geminiConfigFlag, gemini.LoadConfig)
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

func loadBuildPollerConfig(c *cli.Context) (config poller.Config) {
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

	if config.Slack == nil {
		slackConf, confErr := loadSlackConfig(c)
		if confErr != nil {
			logger.Warn("error loading Slack config", zap.Error(confErr))
		} else if slackConf != nil {
			config.Slack = slackConf
		}
	}

	return
}
