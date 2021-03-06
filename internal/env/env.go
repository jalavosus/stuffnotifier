package env

import (
	"os"
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

const (
	GeminiKey     string = "GEMINI_API_KEY"
	GeminiSecret  string = "GEMINI_API_SECRET"
	GeminiSandbox string = "GEMINI_SANDBOX"
)

const (
	DiscordToken string = "DISCORD_TOKEN"
)

const (
	TwilioSid    string = "TWILIO_ACCOUNT_SID"
	TwilioKey    string = "TWILIO_API_KEY"
	TwilioSecret string = "TWILIO_API_SECRET"
	TwilioToken  string = "TWILIO_API_TOKEN"
	TwilioSender string = "TWILIO_SENDER"
)

const (
	FlightAwareKey string = "FLIGHTAWARE_API_KEY"
)

const (
	SlackToken string = "SLACK_TOKEN"
)

const (
	RedisHost     string = "REDIS_HOST"
	RedisPort     string = "REDIS_PORT"
	RedisPassword string = "REDIS_PASSWORD"
)

func FromEnv(envKey string) (string, error) {
	val, ok := os.LookupEnv(envKey)
	if !ok || val == "" {
		return "", NewKeyNotSetError(envKey)
	}

	return val, nil
}

func String(envKey string, fallback string) (string, bool) {
	val, ok := os.LookupEnv(envKey)
	if ok {
		return val, true
	}

	return fallback, false
}

func Int(envKey string, fallback int) (int, bool) {
	val, ok := os.LookupEnv(envKey)
	if ok {
		n := parseInt64(val)
		return int(n), true
	}

	return fallback, false
}

func Int64(envKey string, fallback int64) (int64, bool) {
	val, ok := os.LookupEnv(envKey)
	if ok {
		return parseInt64(val), true
	}

	return fallback, false
}

func parseInt64(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}

func GeminiUseSandbox() bool {
	val, _ := FromEnv(GeminiSandbox)

	switch strings.ToLower(val) {
	case "false", "no":
		return false
	case "true", "yes":
		return true
	}

	return false
}

func TwilioFromNumber() (string, error) {
	return FromEnv(TwilioSender)
}
