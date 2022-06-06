package twilio

import (
	"context"
	"time"

	"github.com/pkg/errors"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"

	"github.com/stoicturtle/stuffnotifier/internal/messages"
	"github.com/stoicturtle/stuffnotifier/internal/utils"
)

func (c *Client) setHttpTimeout(ctx context.Context) {
	if timeout, ok := utils.TimeoutFromContext(ctx); ok {
		c.client.Client.SetTimeout(timeout)
	}
}

func (c *Client) SendMessage(ctx context.Context, msg messages.Message, recipient string) (SendSMSResponse, error) {
	var response = SendSMSResponse{}

	body := msg.FormatPlaintext()

	params := new(openapi.CreateMessageParams)
	params.SetFrom(c.sendNumber)
	params.SetTo(recipient)
	params.SetBody(body)

	c.setHttpTimeout(ctx)

	resp, err := c.client.Api.CreateMessage(params)
	if err != nil {
		response.Success = false
		err = sendMessageErr(err, c.sendNumber, recipient)

		return response, err
	}

	if resp.Status != nil {
		response.MessageStatus = parseStatus(resp.Status)
	}

	if resp.Sid != nil {
		response.MessageSid = *resp.Sid
	}

	if resp.DateSent != nil {
		response.TimeSent, _ = time.Parse(time.RFC822Z, *resp.DateSent)
	}

	return response, nil
}

func (c *Client) Message(ctx context.Context, msgSid string) (*openapi.ApiV2010Message, error) {
	params := new(openapi.FetchMessageParams)
	params.SetPathAccountSid(c.auth.Account())

	c.setHttpTimeout(ctx)

	resp, err := c.client.Api.FetchMessage(msgSid, params)
	if err != nil {
		err = errors.WithMessagef(err, "error fetching message with SID %[1]s", msgSid)
		return nil, err
	}

	return resp, nil
}

func (c *Client) MessageStatus(ctx context.Context, msgSid string) (MessageSendStatus, error) {
	resp, err := c.Message(ctx, msgSid)
	if err != nil {
		return "", err
	}

	if resp.Status == nil {
		return "", errors.Errorf("no status returned for message with SID %[1]s", msgSid)
	}

	status := parseStatus(resp.Status)
	if status == statusUnknown {
		return "", errors.Errorf("unknown message status %[1]s returned for message with SID %[2]s", *resp.Status, msgSid)
	}

	return status, nil
}

func sendMessageErr(err error, from, to string) error {
	return errors.WithMessagef(err, "error sending message from %[1]s to %[2]s", from, to)
}

func parseStatus(status *string) MessageSendStatus {
	switch *status {
	case string(StatusQueued):
		return StatusQueued
	case string(StatusFailed):
		return StatusFailed
	case string(StatusSent):
		return StatusSent
	case string(StatusDelivered):
		return StatusDelivered
	case string(StatusUndelivered):
		return StatusUndelivered
	default:
		return statusUnknown
	}
}
