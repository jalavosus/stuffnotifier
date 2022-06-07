package flightawarepoller

import (
	"context"
	"time"

	"github.com/jalavosus/stuffnotifier/internal/utils"
	"github.com/jalavosus/stuffnotifier/pkg/flightaware"
)

func (p Poller) buildCacheEntry(
	flightData *flightaware.FlightData,
	origin, destination *flightaware.AirportData,
	notificationsSent *SentNotifications,
) CacheEntry {

	return CacheEntry{
		PollerId:          p.PollerIdBytes(),
		FlightIdHash:      utils.SHA3(flightData.FlightId),
		InternalId:        flightData.FlightId,
		FlightId:          flightData.Identifiers.IATA,
		FlightData:        flightData,
		OriginData:        origin,
		DestinationData:   destination,
		PollInterval:      p.FlightAwareConfig().PollInterval,
		Notifications:     p.FlightAwareConfig().Notifications,
		RecipientConfig:   p.BuildRecipientConfig(),
		NotificationsSent: notificationsSent,
	}
}

func (p Poller) fetchFlightIdentifiers(
	ctx context.Context,
	flightId string,
	flightIdType flightaware.IdentifierType,
) (*flightaware.FlightIdentifiers, string, error) {

	params := buildFlightInformationParams(flightId, flightIdType)

	ctx, cancel := context.WithTimeout(ctx, fetchDataTimeout)
	defer cancel()

	identifiers, apiId, err := p.flightawareClient.FlightIdentifiers(ctx, params)
	if err != nil {
		return nil, "", err
	}

	return identifiers, apiId, nil
}

func (p Poller) fetchFlight(
	ctx context.Context,
	flightId string,
	flightIdType flightaware.IdentifierType,
) (*flightaware.FlightData, error) {

	params := buildFlightInformationParams(flightId, flightIdType)

	ctx, cancel := context.WithTimeout(ctx, fetchDataTimeout)
	defer cancel()

	data, err := p.flightawareClient.FlightInformation(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(data.Flights) == 0 {
		return nil, flightaware.ErrNoFlightsFound
	}

	flight, ok := flightaware.FindFlightFromParams(data.Flights, params)
	if ok {
		return flight, nil
	}

	return nil, flightaware.ErrNoFlightsFound
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

func (p Poller) fetchAll(
	ctx context.Context,
	flightId string,
	idType flightaware.IdentifierType,
) (flightData *flightaware.FlightData, originInfo, destinationInfo *flightaware.AirportData, err error) {

	flightData, err = p.fetchFlight(ctx, flightId, idType)
	if err != nil {
		return
	}

	originInfo, err = p.fetchAirport(ctx, flightData.Origin.Identifiers.ICAO)
	if err != nil {
		return
	}

	destinationInfo, err = p.fetchAirport(ctx, flightData.Destination.Identifiers.ICAO)
	if err != nil {
		return
	}

	return
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

func buildFlightInformationParams(flightId string, idType flightaware.IdentifierType) (params flightaware.FlightInformationParams) {
	params = flightaware.FlightInformationParams{}

	switch idType {
	case flightaware.FaFlightIdIdent:
		params.FlightId = utils.ToPointer(flightId)
	case flightaware.DesignatorIdent:
		params.FlightDesignator = utils.ToPointer(flightId)
	}

	return
}

func validTimestampActual(ts flightaware.FlightTimestamp) bool {
	return !ts.Actual.IsZero()
}
