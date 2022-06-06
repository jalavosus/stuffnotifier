package env

import (
	"os"
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

func FromEnv(envKey string) (string, error) {
	val, ok := os.LookupEnv(envKey)
	if !ok || val == "" {
		return "", NewKeyNotSetError(envKey)
	}

	return val, nil
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
