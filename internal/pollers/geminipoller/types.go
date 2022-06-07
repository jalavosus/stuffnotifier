package geminipoller

import (
	"time"

	"github.com/jalavosus/stuffnotifier/internal/pollers/poller"
	"github.com/jalavosus/stuffnotifier/pkg/gemini"
)

type CacheEntry struct {
	PollerId        []byte
	SymbolHash      []byte
	RecipientConfig poller.RecipientConfig
	Notifications   gemini.NotificationsConfig
	PollInterval    time.Duration
}
