package poller

import (
	"sync"
	"time"

	"github.com/jalavosus/stuffnotifier/internal/utils"
	"github.com/jalavosus/stuffnotifier/pkg/authdata"
)

type RecipientConfig struct {
	TwilioRecipientNumber     string
	DiscordRecipientUserId    string
	DiscordRecipientChannelId string
	UseTwilio                 bool
	UseDiscord                bool
}

func (p BasePoller) BuildRecipientConfig() RecipientConfig {
	d := RecipientConfig{}

	if p.TwilioConfig() != nil {
		d.UseTwilio = true
		d.TwilioRecipientNumber = p.TwilioConfig().RecipientNumber
	}

	if p.DiscordConfig() != nil {
		if userId, ok := utils.FromPointer(p.DiscordConfig().UserId); ok && userId != "" {
			d.UseDiscord = true
			d.DiscordRecipientUserId = userId
		}

		if channelId, ok := utils.FromPointer(p.DiscordConfig().ChannelId); ok && channelId != "" {
			d.UseDiscord = true
			d.DiscordRecipientChannelId = channelId
		}
	}

	return d
}

type ConcurrentParams struct {
	AuthData authdata.AuthData
	ErrCh    chan error
	CacheKey string
	once     sync.Once // for cleanup
}

func NewConcurrentParams(authData authdata.AuthData, cacheKey string) *ConcurrentParams {
	p := new(ConcurrentParams)
	p.AuthData = authData
	p.CacheKey = cacheKey
	p.ErrCh = make(chan error, 1)

	return p
}

func (p *ConcurrentParams) Cleanup(err error, ticker *time.Ticker) {
	ticker.Stop()
	p.ErrCh <- err
}
