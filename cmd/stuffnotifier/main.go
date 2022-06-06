package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

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
		log.Println(err)
	}
}
