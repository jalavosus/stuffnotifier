package authdata

import (
	"github.com/jalavosus/stuffnotifier/internal/env"
)

// GeminiAPIAuth returns an AuthData whose
// Key and Secret fields are populated with an
// api key and secret provided by Gemini.
func GeminiAPIAuth() (auth AuthData, err error) {
	apiKey, err := env.FromEnv(env.GeminiKey)
	if err != nil {
		return
	}

	apiSecret, err := env.FromEnv(env.GeminiSecret)
	if err != nil {
		return
	}

	auth = NewAuthData("", apiKey, apiSecret)

	return
}

// TwilioAPIAuth returns an AuthData whose
// Account field is populated with a Twilio Account SID,
// Key field populated with a Twilio API key,
// and Secret field populated with a Twilio API secret.
func TwilioAPIAuth() (auth AuthData, err error) {
	accountSid, err := env.FromEnv(env.TwilioSid)
	if err != nil {
		return
	}

	apiKey, err := env.FromEnv(env.TwilioKey)
	if err != nil {
		return
	}

	apiSecret, err := env.FromEnv(env.TwilioSecret)
	if err != nil {
		return
	}

	auth = NewAuthData(accountSid, apiKey, apiSecret)

	return
}

// TwilioAPITokenAuth returns an AuthData whose
// Account and Key fields are populated with a
// Twilio Account SID, and Secret field populated
// with a Twilio Auth Token.
func TwilioAPITokenAuth() (auth AuthData, err error) {
	accountSid, err := env.FromEnv(env.TwilioSid)
	if err != nil {
		return
	}

	authToken, err := env.FromEnv(env.TwilioToken)
	if err != nil {
		return
	}

	auth = NewAuthData(accountSid, accountSid, authToken)

	return
}

// DiscordAPIAuth returns an AuthData whose Secret
// field is populated with a Discord bot token.
func DiscordAPIAuth() (auth AuthData, err error) {
	apiToken, err := env.FromEnv(env.DiscordToken)
	if err != nil {
		return
	}

	auth = NewAuthData("", "", apiToken)

	return
}

// FlightAwareAPIAuth returns an AuthData whose Key
// field is populated with a FlightAware api key.
func FlightAwareAPIAuth() (auth AuthData, err error) {
	apiKey, err := env.FromEnv(env.FlightAwareKey)
	if err != nil {
		return
	}

	auth = NewAuthData("", apiKey, "")

	return
}

func SlackAPIAuth() (auth AuthData, err error) {
	token, err := env.FromEnv(env.SlackToken)
	if err != nil {
		return
	}

	auth = NewAuthData("", "", token)

	return
}

func RedisAuth() (auth ServiceAuthData, err error) {
	redisHost, _ := env.String(env.RedisHost, "localhost")
	redisPort, _ := env.Int(env.RedisPort, 6379)
	redisPassword, _ := env.String(env.RedisPassword, "")

	auth = NewServiceAuthData(redisHost, redisPort, "", redisPassword)

	return
}
