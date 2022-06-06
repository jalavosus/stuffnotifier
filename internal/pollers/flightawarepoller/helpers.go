package flightawarepoller

import (
	"context"
	"sort"
	"time"

	"github.com/pkg/errors"

	"github.com/stoicturtle/stuffnotifier/internal/utils"
	"github.com/stoicturtle/stuffnotifier/pkg/flightaware"
)

func (p Poller) buildCacheEntry(
	flightData *flightaware.FlightData,
	origin, destination *flightaware.AirportData,
	notificationsSent *SentNotifications,
) CacheEntry {

	return CacheEntry{
		InternalId:        flightData.FlightId,
		FlightId:          flightData.Identifier.IATA,
		FlightData:        flightData,
		OriginData:        origin,
		DestinationData:   destination,
		PollInterval:      p.FlightAwareConfig().PollInterval,
		Notifications:     p.FlightAwareConfig().Notifications,
		RecipientConfig:   p.BuildRecipientConfig(),
		NotificationsSent: notificationsSent,
	}
}

func (p Poller) fetchFlight(ctx context.Context, flightId string, flightIdType flightaware.IdentifierType, nearest bool) (*flightaware.FlightData, error) {
	var flights []flightaware.FlightData

	ctx, cancel := context.WithTimeout(ctx, fetchDataTimeout)
	defer cancel()

	data, err := p.flightawareClient.FlightInformation(ctx, flightId, flightIdType)
	if err != nil {
		return nil, err
	}

	if len(data.Flights) == 0 {
		return nil, errors.New("no flights found")
	}

	switch flightIdType {
	case flightaware.FaFlightIdIdent:
		flights = data.Flights
	case flightaware.DesignatorIdent:
		now := utils.ToLocalTime(time.Now())
		y1, m1, d1 := now.Date()

		for _, fl := range data.Flights {
			var appendFlight = true

			if nearest {
				departure := utils.ToLocalTime(fl.GateDepartureTime.Scheduled)
				y2, m2, d2 := departure.Date()

				appendFlight = y2 >= y1 && m2 >= m1 && d2 >= d1
			}

			if appendFlight {
				flights = append(flights, fl)
			}
		}
	}

	sort.Slice(flights, func(i, j int) bool {
		return flights[i].GateDepartureTime.Scheduled.Before(flights[j].GateDepartureTime.Scheduled)
	})

	flight := utils.SliceFirst(flights)

	return &flight, nil
}

func (p Poller) fetchAirport(ctx context.Context, airportId string) (*flightaware.AirportData, error) {
	ctx, cancel := context.WithTimeout(ctx, fetchDataTimeout)
	defer cancel()

	data, err := p.flightawareClient.AirportInformation(ctx, airportId)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (p Poller) fetchCacheEntry(ctx context.Context, cacheKey string) (cacheData *CacheEntry, ok bool, err error) {
	ctx, cancel := context.WithTimeout(ctx, fetchDataTimeout)
	defer cancel()

	cacheData, ok, err = p.Datastore().Get(ctx, cacheKey)
	if err != nil {
		ok = false
	}

	return
}

func (p *Poller) setCacheEntry(
	ctx context.Context,
	cacheKey string,
	flightData *flightaware.FlightData,
	origin, dest *flightaware.AirportData,
	notificationsSent *SentNotifications,
) error {

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	cacheData := p.buildCacheEntry(flightData, origin, dest, notificationsSent)

	return p.Datastore().Insert(ctx, cacheKey, cacheData)
}

func validTimestampActual(ts flightaware.FlightTimestamp) bool {
	return !ts.Actual.IsZero()
}
