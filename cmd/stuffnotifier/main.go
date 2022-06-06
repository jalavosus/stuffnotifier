package main

import (
	"os"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/stoicturtle/stuffnotifier/internal/logging"
)

var logger = logging.NewLogger()

func main() {
	app := &cli.App{
		Name: "stuffnotifier",
		Commands: []*cli.Command{
			&flightawareCmd,
			&geminiCmd,
		},
		Flags: []cli.Flag{
			&twilioConfigFlag,
			&discordConfigFlag,
			&twilioSidFlag,
			&twilioApiKeyFlag,
			&twilioApiSecretFlag,
			&twilioApiTokenFlag,
			&discordTokenFlag,
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal("", zap.Error(err))
	}
}
