package poller

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/stoicturtle/stuffnotifier/internal/messages"
	"github.com/stoicturtle/stuffnotifier/internal/utils"
	"github.com/stoicturtle/stuffnotifier/pkg/discord"
	"github.com/stoicturtle/stuffnotifier/pkg/slack"
	"github.com/stoicturtle/stuffnotifier/pkg/twilio"
)

type BasePoller struct {
	logger       *zap.Logger
	config       Config
	pollInterval time.Duration
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

func (p BasePoller) SlackConfig() *slack.Config {
	return p.config.Slack
}

func (p BasePoller) SendMessage(ctx context.Context, msg messages.Message) error {
	if p.LogStdout() {
		fmt.Println(msg.FormatPlaintext())
	}

	if p.TwilioConfig() != nil {
		if err := p.sendTwilio(ctx, msg); err != nil {
			return err
		}
	}

	if p.DiscordConfig() != nil {
		// TODO
	}

	if p.SlackConfig() != nil {
		if err := p.sendSlack(ctx, msg); err != nil {
			return err
		}
	}

	return nil
}

func (p BasePoller) sendTwilio(ctx context.Context, msg messages.Message) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, clientErr := twilio.ClientSingleton(p.TwilioConfig())
	if clientErr != nil {
		return clientErr
	}

	if _, err := client.Client().SendMessage(ctx, msg, p.TwilioConfig().RecipientNumber); err != nil {
		return err
	}

	return nil
}

func (p BasePoller) sendSlack(ctx context.Context, msg messages.Message) error {
	slackConf := p.SlackConfig()

	client, clientErr := slack.NewClient(slackConf)
	if clientErr != nil {
		return clientErr
	}

	recipientIds := utils.AppendSlices(slackConf.Channels, slackConf.Users)

	for _, id := range recipientIds {
		if err := p.sendSlackMsg(ctx, msg, id, client); err != nil {
			return err
		}
	}

	return nil
}

func (p BasePoller) sendSlackMsg(ctx context.Context, msg messages.Message, recipientId string, client *slack.Client) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return client.SendMessage(ctx, msg, recipientId)
}
