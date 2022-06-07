package flightawarepoller

import (
	"time"

	"github.com/jalavosus/stuffnotifier/internal/datastore"
	"github.com/jalavosus/stuffnotifier/internal/pollers/poller"
	"github.com/jalavosus/stuffnotifier/pkg/authdata"
	"github.com/jalavosus/stuffnotifier/pkg/flightaware"
)

const (
	fetchDataTimeout = time.Minute
)

type Poller struct {
	datastore datastore.Datastore[CacheEntry]
	*poller.BasePoller
	flightawareClient *flightaware.Client
	flightawareConfig flightaware.Config
}

func NewPoller(conf poller.Config, flightawareConf *flightaware.Config) (*Poller, error) {
	var (
		faConf       = flightaware.DefaultConfig()
		datastoreErr error
	)

	switch {
	case conf.FlightAware != nil:
		faConf = *conf.FlightAware
	case flightawareConf != nil:
		faConf = *flightawareConf
	}

	p := &Poller{
		BasePoller:        poller.NewBasePoller(conf),
		flightawareConfig: faConf,
	}

	p.datastore, datastoreErr = datastore.NewDatastore[CacheEntry](conf.Cache)
	if datastoreErr != nil {
		return nil, datastoreErr
	}

	p.SetPollInterval(faConf.PollInterval)

	return p, nil
}

func (p *Poller) FlightAwareConfig() flightaware.Config {
	return p.flightawareConfig
}

func (p *Poller) Datastore() datastore.Datastore[CacheEntry] {
	return p.datastore
}

func (p *Poller) FlightAwareClient() *flightaware.Client {
	return p.flightawareClient
}

func (p *Poller) initFlightAwareClient(authData authdata.AuthData) {
	p.flightawareClient = flightaware.NewClient(authData)
}
