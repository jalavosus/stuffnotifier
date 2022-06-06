package twilio

import (
	"time"
)

type MessageSendStatus string

const (
	statusUnknown     MessageSendStatus = ""
	StatusQueued      MessageSendStatus = "queued"
	StatusFailed      MessageSendStatus = "failed"
	StatusSent        MessageSendStatus = "sent"
	StatusDelivered   MessageSendStatus = "delivered"
	StatusUndelivered MessageSendStatus = "undelivered"
)

type SendSMSResponse struct {
	Success       bool
	TimeSent      time.Time
	MessageSid    string
	MessageStatus MessageSendStatus
	MessageError  string
}
