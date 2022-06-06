package flightawarepoller

import (
	"time"

	"github.com/stoicturtle/stuffnotifier/internal/datastore"
	"github.com/stoicturtle/stuffnotifier/internal/pollers/poller"
	"github.com/stoicturtle/stuffnotifier/pkg/authdata"
	"github.com/stoicturtle/stuffnotifier/pkg/flightaware"
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

func NewFlightAwarePoller(conf poller.Config, flightawareConf *flightaware.Config) (*Poller, error) {
	var (
		faConf       = flightaware.DefaultConfig()
		datastoreErr error
	)

	if flightawareConf != nil {
		if conf.FlightAware != nil {
			faConf = *conf.FlightAware
		} else {
			faConf = *flightawareConf
		}
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

func (p Poller) FlightAwareConfig() flightaware.Config {
	return p.flightawareConfig
}

func (p Poller) UseDatastore() bool {
	return p.datastore != nil
}

func (p Poller) Datastore() datastore.Datastore[CacheEntry] {
	return p.datastore
}

func (p Poller) FlightAwareClient() *flightaware.Client {
	return p.flightawareClient
}

func (p *Poller) initFlightAwareClient(authData authdata.AuthData) {
	p.flightawareClient = flightaware.NewClient(authData)
}

// func prettyPrint(data any) {
// 	d, _ := yaml.Marshal(data)
// 	log.Println(string(d))
// }
