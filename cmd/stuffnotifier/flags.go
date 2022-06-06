package main

import (
	"github.com/urfave/cli/v2"

	"github.com/jalavosus/stuffnotifier/internal/env"
)

const (
	categoryConfig      string = "Config"
	categoryAuth        string = "Auth"
	categoryTwilioAuth         = categoryAuth + " - Twilio"
	categoryDiscordAuth        = categoryAuth + " - Discord"
)

const (
	configFlagName      string = "config"
	twilioFlagName      string = "twilio"
	discordFlagName     string = "discord"
	geminiFlagName      string = "gemini"
	flightAwareFlagName string = "flightaware"
)

const (
	apiKeyFlag    string = "apiKey"
	apiSecretFlag string = "apiSecret"
	authTokenFlag string = "authToken"
)

func makeSuffixedFlagName(flagName, suffix string) string {
	return flagName + "." + suffix
}

func makeConfigFlag(flagName string) cli.PathFlag {
	var (
		suffixedFlagName string
		usage            string
	)

	if flagName == "" {
		suffixedFlagName = configFlagName
		usage = "`path` to a config file"
	} else {
		suffixedFlagName = makeSuffixedFlagName(configFlagName, flagName)
		usage = "`path` to a " + flagName + " config file"
	}

	return cli.PathFlag{
		Name:     suffixedFlagName,
		Usage:    usage,
		Category: categoryConfig,
		Required: false,
	}
}

var (
	pollerConfigFlag      = makeConfigFlag("")
	geminiConfigFlag      = makeConfigFlag(geminiFlagName)
	flightawareConfigFlag = makeConfigFlag(flightAwareFlagName)
	twilioConfigFlag      = makeConfigFlag(twilioFlagName)
	discordConfigFlag     = makeConfigFlag(discordFlagName)
)

var (
	twilioSidFlag = cli.StringFlag{
		Name:     makeSuffixedFlagName(twilioFlagName, "accountSid"),
		Usage:    "Twilio account `sid`. Required for any form of Twilio API authentication",
		Category: categoryTwilioAuth,
		Required: false,
		EnvVars:  []string{env.TwilioSid},
	}
	twilioApiKeyFlag = cli.StringFlag{
		Name:     makeSuffixedFlagName(twilioFlagName, apiKeyFlag),
		Usage:    "Twilio API `key` (for API key based authentication)",
		Category: categoryTwilioAuth,
		Required: false,
		EnvVars:  []string{env.TwilioKey},
	}
	twilioApiSecretFlag = cli.StringFlag{
		Name:     makeSuffixedFlagName(twilioFlagName, apiSecretFlag),
		Usage:    "Twilio API `secret` (for API key based authentication). Required if twilioApiKey is set",
		Category: categoryTwilioAuth,
		Required: false,
		EnvVars:  []string{env.TwilioSecret},
	}
	twilioApiTokenFlag = cli.StringFlag{
		Name:     makeSuffixedFlagName(twilioFlagName, authTokenFlag),
		Usage:    "Twilio API auth `token` (for auth token based authentication)",
		Category: categoryTwilioAuth,
		Required: false,
		EnvVars:  []string{env.TwilioToken},
	}
)

var (
	discordTokenFlag = cli.StringFlag{
		Name:     makeSuffixedFlagName(discordFlagName, authTokenFlag),
		Usage:    "Discord Bot `token`. Required for any/all interactions with Discord",
		Category: categoryDiscordAuth,
		Required: false,
		EnvVars:  []string{env.DiscordToken},
	}
)

var (
	geminiApiKeyFlag = cli.StringFlag{
		Name:     makeSuffixedFlagName(geminiFlagName, apiKeyFlag),
		Usage:    "Gemini API `key`. Required for any/all interactions with Gemini",
		Category: categoryAuth,
		Required: false,
		EnvVars:  []string{env.GeminiKey},
	}
	geminiApiSecret = cli.StringFlag{
		Name:     makeSuffixedFlagName(geminiFlagName, apiSecretFlag),
		Usage:    "Gemini API `secret`. Required for any/all interactions with Gemini",
		Category: categoryAuth,
		Required: false,
		EnvVars:  []string{env.GeminiSecret},
	}
)

var (
	flightAwareApiKeyFlag = cli.StringFlag{
		Name:     makeSuffixedFlagName(flightAwareFlagName, apiKeyFlag),
		Usage:    "FlightAware API `key`. Required for any/all interfactions with FlightAware",
		Category: categoryAuth,
		Required: false,
		EnvVars:  []string{env.FlightAwareKey},
	}
)
