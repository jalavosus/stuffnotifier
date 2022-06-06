package poller

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/stoicturtle/stuffnotifier/internal/messages"
	"github.com/stoicturtle/stuffnotifier/pkg/discord"
	"github.com/stoicturtle/stuffnotifier/pkg/twilio"
)

type BasePoller struct {
	pollInterval time.Duration
	config       Config
	logger       *zap.Logger
}

func NewBasePoller(conf Config) *BasePoller {
	p := &BasePoller{
		pollInterval: DefaultPollInterval,
		config:       conf,
		logger:       newLogger(),
	}

	return p
}

func (p BasePoller) PollInterval() time.Duration {
	return p.pollInterval
}

func (p *BasePoller) SetPollInterval(pollInterval time.Duration) *BasePoller {
	p.pollInterval = pollInterval

	return p
}

func (p BasePoller) LogStdout() bool {
	return p.config.LogStdout
}

func (p BasePoller) TwilioConfig() *twilio.Config {
	return p.config.Twilio
}

func (p BasePoller) DiscordConfig() *discord.Config {
	return p.config.Discord
}

func (p BasePoller) SendMessage(ctx context.Context, msg messages.Message) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if p.LogStdout() {
		fmt.Println(msg.FormatPlaintext())
	}

	if p.TwilioConfig() != nil {
		client, clientErr := twilio.ClientSingleton(p.TwilioConfig())
		if clientErr != nil {
			return clientErr
		}

		if _, err := client.Client().SendMessage(ctx, msg, p.TwilioConfig().RecipientNumber); err != nil {
			return err
		}
	}

	if p.DiscordConfig() != nil {
		// TODO
	}

	return nil
}
