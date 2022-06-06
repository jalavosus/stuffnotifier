package slack

import (
	"context"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"

	"github.com/jalavosus/stuffnotifier/internal/messages"
	"github.com/jalavosus/stuffnotifier/pkg/authdata"
)

// Client is a wrapper around the slack-go/slack
// client.
type Client struct {
	client   *slack.Client
	authData authdata.AuthData
	botId    string
	// clientConfig *slack.
}

func NewClient(conf *Config) (*Client, error) {
	var authData authdata.AuthData

	if conf != nil && conf.Auth != nil {
		authData = conf.Auth
	} else {
		ad, err := authdata.SlackAPIAuth()
		if err != nil {
			return nil, err
		}

		authData = ad
	}

	if authData == nil {
		return nil, errors.New("no slack token found in configuration or environment")
	}

	c := new(Client)

	c.client = slack.New(authData.Secret())
	c.authData = authData

	authResponse, err := c.client.AuthTest()
	if err != nil {
		return nil, err
	}

	c.botId = authResponse.BotID

	return c, nil
}

func (c Client) UserInfo(ctx context.Context, userId string) (*slack.User, error) {
	return c.client.GetUserInfoContext(ctx, userId)
}

func (c Client) SendMessage(ctx context.Context, msg messages.Message, channelId string) error {
	return c.sendMessage(ctx, msg, channelId)
}

func (c Client) sendMessage(ctx context.Context, msg messages.Message, channelId string) error {
	formattedMsg := msg.FormatMarkdown()

	attachment := slack.Attachment{
		Pretext: "StuffNotifier alert",
		Text:    formattedMsg,
	}

	_, _, err := c.client.PostMessageContext(ctx, channelId, slack.MsgOptionAttachments(attachment))

	return err
}
