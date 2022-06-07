package flightawarepoller

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/jalavosus/stuffnotifier/internal/messages"
	"github.com/jalavosus/stuffnotifier/internal/pollers/poller"
	"github.com/jalavosus/stuffnotifier/internal/utils"
	"github.com/jalavosus/stuffnotifier/pkg/authdata"
	"github.com/jalavosus/stuffnotifier/pkg/flightaware"
)

func (p *Poller) Start(ctx context.Context, flightId string, flightIdType flightaware.IdentifierType) error {
	var authData authdata.AuthData

	if conf := p.FlightAwareConfig(); conf.Auth != nil {
		authData = conf.Auth
	} else {
		ad, authDataErr := authdata.FlightAwareAPIAuth()
		if authDataErr != nil {
			return authDataErr
		}

		authData = ad
	}

	p.initFlightAwareClient(authData)

	_, apiId, err := p.fetchFlightIdentifiers(ctx, flightId, flightIdType)
	if err != nil {
		return err
	}

	cacheKey := "flightdata:" + apiId
	concurrentParams := poller.NewConcurrentParams(authData, cacheKey)

	go p.pollFlightData(
		ctx,
		apiId,
		concurrentParams,
	)

	return <-concurrentParams.ErrCh
}

//nolint:gocognit
func (p *Poller) pollFlightData(
	ctx context.Context,
	flightId string,
	pollerParams *poller.ConcurrentParams,
) {

	var (
		flightData                  *flightaware.FlightData
		originInfo, destinationInfo *flightaware.AirportData
		gotCachedData               bool
		fetchAllErr                 error
	)

	var (
		notifsSent = new(SentNotifications)
		cacheKey   = pollerParams.CacheKey
	)

	ticker := time.NewTicker(p.PollInterval())
	cleanup := func(err error) {
		pollerParams.Cleanup(err, ticker)
	}

	p.LogDebug("checking for cached data...")
	if cached, ok, cacheErr := p.fetchCacheEntry(ctx, cacheKey); cacheErr != nil {
		p.LogError("error checking for cached data", zap.Error(cacheErr))
	} else if ok {
		flightData = cached.FlightData
		originInfo = cached.OriginData
		destinationInfo = cached.DestinationData
		notifsSent = cached.NotificationsSent

		gotCachedData = true
	}

	if !gotCachedData {
		flightData, originInfo, destinationInfo, fetchAllErr = p.fetchAll(ctx, flightId, flightaware.FaFlightIdIdent)
		if fetchAllErr != nil {
			cleanup(fetchAllErr)
			return
		}
	}

	notifsSent.SetDisabled(p.FlightAwareConfig().Notifications)

	p.LogInfo(
		"starting flight data poller",
		zap.String("flight_identifier", flightData.Identifiers.IATA),
		zap.String("flight_datetime", utils.FormatTimeWithZone(flightData.GateDepartureTime.Scheduled, originInfo.Timezone, true)),
		zap.String("check_interval", p.PollInterval().String()),
	)

	isInitial := true

	for {
		select {
		case <-ctx.Done():
			cleanup(nil)
			return
		case t := <-ticker.C:
			if notifsSent.SentAll() {
				cleanup(nil)
				return
			}

			if !isInitial {
				var flightDataErr error

				flightData, flightDataErr = p.fetchFlight(ctx, flightId, flightaware.FaFlightIdIdent)
				if flightDataErr != nil {
					p.LogError("error fetching flight data", zap.Error(flightDataErr))
					continue
				}
			} else {
				isInitial = false
			}

			setCacheErr := p.setCacheEntry(
				ctx,
				cacheKey,
				flightData,
				originInfo,
				destinationInfo,
				notifsSent,
			)

			if setCacheErr != nil {
				p.LogError("error setting flight data in cache", zap.Error(setCacheErr))
			}

			flightIdentifier := flightData.Identifiers.IATA

			var (
				flightEvents = newPastFlightEvents(flightData)
				notifType    = NoNotification
				msg          = newAlertMsg(
					flightIdentifier,
					flightData,
					originInfo,
					destinationInfo,
					p.FlightAwareConfig().Notifications.UseLocalTime,
				)
			)

		DetermineNotify:
			switch {
			case flightEvents.Pending():
				now := time.Now().UTC()
				departureOffset := flightData.GateDepartureTime.Scheduled.
					UTC().
					Sub(now)

				inPreDepartureWindow := departureOffset <= p.FlightAwareConfig().Notifications.PreDeparture.Offset
				if !notifsSent.PreDeparture && inPreDepartureWindow {
					msg.SetPlaintextTemplate(messages.PreDepartureAlertPlaintextTemplate)
					notifType = PreDepartureNotification

					break DetermineNotify
				}
			case flightEvents.InProgress():
				didSendGateDeparture := flightEvents.DepartedGate && notifsSent.GateDeparture
				didSendTakeoff := flightEvents.Takeoff && notifsSent.Takeoff

				arrivalOffset := flightData.RunwayArrivalTime.Estimated.
					UTC().
					Sub(t.UTC())

				inPreArrivalWindow := arrivalOffset <= p.FlightAwareConfig().Notifications.PreArrival.Offset

				var (
					isGateDeparture, isTakeoff bool
				)

				switch {
				case !notifsSent.PreArrival && inPreArrivalWindow:
					msg.SetPlaintextTemplate(messages.PreArrivalAlertPlaintextTemplate)
					notifType = PreArrivalNotification

					break DetermineNotify
				case didSendGateDeparture && didSendTakeoff:
					continue
				case flightEvents.DepartedGate && !flightEvents.Takeoff:
					notifType = GateDepartureNotification
					isGateDeparture = true
				case flightEvents.Takeoff:
					notifType = TakeoffNotification
					isTakeoff = true
				}

				msg = setMsgDepartureData(flightData, isGateDeparture, isTakeoff, msg)
			case flightEvents.ArrivedDestination():
				didSendLanding := flightEvents.Landed && notifsSent.Landing
				didSendGateArrival := flightEvents.ArrivedGate && notifsSent.GateArrival

				var (
					isGateArrival, isLanding bool
				)

				switch {
				case didSendGateArrival:
					break DetermineNotify
				case didSendLanding:
					continue
				case flightEvents.Landed && !flightEvents.ArrivedGate:
					notifType = LandingNotification
					isLanding = true
				case flightEvents.ArrivedGate:
					notifType = GateArrivalNotification
					isGateArrival = true
				}

				msg = setMsgArrivalData(flightData, isGateArrival, isLanding, msg)
			default:
				msg = nil
				notifType = NoNotification
			}

			if msg == nil || notifType == NoNotification {
				continue
			}

			if sendMsgErr := p.SendMessage(ctx, msg); sendMsgErr != nil {
				p.LogError("error sending notification", zap.Error(sendMsgErr))
				continue
			}

			notifsSent.SetSent(notifType)

			if notifsSent.SentAll() {
				cleanup(nil)
				return
			}
		}
	}
}

func newAlertMsg(
	flightIdentifier string,
	flightData *flightaware.FlightData,
	originInfo, destinationInfo *flightaware.AirportData,
	useLocalTime bool,
) *messages.FlightAwareAlert {

	return &messages.FlightAwareAlert{
		UseLocalTimezone: useLocalTime,
		FlightNumber:     flightIdentifier,
		Origin: messages.FlightAwareAirportInfo{
			Airport:  originInfo.Name,
			Timezone: originInfo.Timezone,
			Gate:     flightData.Origin.Gate,
			Terminal: flightData.Origin.Terminal,
		},
		Destination: messages.FlightAwareAirportInfo{
			Airport:  destinationInfo.Name,
			Timezone: destinationInfo.Timezone,
			Gate:     flightData.Destination.Gate,
			Terminal: flightData.Destination.Terminal,
		},
		GateDepartureTime: flightData.GateDepartureTime,
		TakeoffTime:       flightData.RunwayDepartureTime,
		LandingTime:       flightData.RunwayArrivalTime,
		GateArrivalTime:   flightData.GateArrivalTime,
	}
}

func setMsgDepartureData(
	flightData *flightaware.FlightData,
	isGateDeparture bool,
	isTakeoff bool,
	msg *messages.FlightAwareAlert,
) *messages.FlightAwareAlert {

	msg.IsGateDeparture = isGateDeparture
	msg.IsTakeoff = isTakeoff
	msg.GateDepartureTime.Actual = flightData.GateDepartureTime.Actual
	msg.TakeoffTime.Actual = flightData.RunwayDepartureTime.Actual
	msg.SetPlaintextTemplate(messages.DepartureAlertPlaintextTemplate)

	return msg
}

func setMsgArrivalData(
	flightData *flightaware.FlightData,
	isGateArrival bool,
	isLanding bool,
	msg *messages.FlightAwareAlert,
) *messages.FlightAwareAlert {

	msg.IsGateArrival = isGateArrival
	msg.IsLanding = isLanding
	msg.GateArrivalTime.Actual = flightData.GateArrivalTime.Actual
	msg.LandingTime.Actual = flightData.RunwayArrivalTime.Actual
	msg.SetPlaintextTemplate(messages.ArrivalAlertPlaintextTemplate)

	return msg
}
