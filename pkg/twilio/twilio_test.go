package twilio_test

import (
	"context"
	"testing"
	"text/template"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/jalavosus/stuffnotifier/pkg/twilio"
)

type testMessage struct {
	body string
}

func (t testMessage) FormatPlaintext() string {
	return t.body
}

func (t testMessage) FormatMarkdown() string {
	return t.body
}

func (t testMessage) PlaintextTemplate() *template.Template {
	return nil
}

func (t testMessage) MarkdownTemplate() *template.Template {
	return nil
}

func TestSendMessage(t *testing.T) {
	t.Skip()

	var testMsg = testMessage{"Hello, world!"}

	rootCtx := context.Background()

	type testCase struct {
		name        string
		msg         testMessage
		recipient   string
		wantStatus  twilio.MessageSendStatus
		wantSuccess bool
		wantError   bool
	}

	testCases := []testCase{
		{
			name:        "invalid phone number",
			msg:         testMsg,
			recipient:   "+15005550001",
			wantStatus:  twilio.StatusUndelivered,
			wantSuccess: false,
		},
		{
			name:        "with country code",
			msg:         testMsg,
			recipient:   "+16037065541",
			wantStatus:  twilio.StatusDelivered,
			wantSuccess: true,
		},
		{
			name:        "with country code/+ dashes",
			msg:         testMsg,
			recipient:   "+1-603-706-5541",
			wantSuccess: true,
			wantStatus:  twilio.StatusDelivered,
		},
		{
			name:        "no country code",
			msg:         testMsg,
			recipient:   "6037065541",
			wantSuccess: true,
			wantStatus:  twilio.StatusDelivered,
		},
		{
			name:        "no country code/+ dashes",
			msg:         testMsg,
			recipient:   "603-706-5541",
			wantSuccess: true,
			wantStatus:  twilio.StatusDelivered,
		},
		{
			name:        "with country code/+ dashes/no country code dash",
			msg:         testMsg,
			recipient:   "+1 603-706-5541",
			wantSuccess: true,
			wantStatus:  twilio.StatusDelivered,
		},
		{
			name:        "with country code/+ dashes/no country code plus",
			msg:         testMsg,
			recipient:   "1-603-706-5541",
			wantSuccess: true,
			wantStatus:  twilio.StatusDelivered,
		},
		{
			name:        "with country code/+ dashes/no country code dash or plus",
			msg:         testMsg,
			recipient:   "1 603-706-5541",
			wantSuccess: true,
			wantStatus:  twilio.StatusDelivered,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(rootCtx, 10*time.Second)
			defer cancel()

			client, err := twilio.NewClient()
			assert.NoError(t, err)

			resp, err := client.SendMessage(ctx, tc.msg, tc.recipient)
			if tc.wantError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equalf(t, tc.wantSuccess, resp.Success, "expected Success to be %[1]t, got %[2]t", tc.wantSuccess, resp.Success)
			assert.Equal(t, tc.wantStatus, resp.MessageStatus)
		})
	}
}
