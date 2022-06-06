package flightawarepoller

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/stoicturtle/stuffnotifier/internal/authdata"
	"github.com/stoicturtle/stuffnotifier/internal/messages"
	"github.com/stoicturtle/stuffnotifier/internal/pollers/poller"
	"github.com/stoicturtle/stuffnotifier/internal/utils"
	"github.com/stoicturtle/stuffnotifier/pkg/flightaware"
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

	cacheKey := "flightdata:" + flightId + "_" + flightIdType.String()
	concurrentParams := poller.NewConcurrentParams(authData, cacheKey)

	go p.pollFlightData(
		ctx,
		flightId,
		flightIdType,
		concurrentParams,
	)

	return <-concurrentParams.ErrCh
}

func (p *Poller) pollFlightData(
	ctx context.Context,
	flightId string,
	flightIdType flightaware.IdentifierType,
	pollerParams *poller.ConcurrentParams,
) {

	fetchNearest := flightIdType == flightaware.DesignatorIdent

	var (
		flightData                  *flightaware.FlightData
		originInfo, destinationInfo *flightaware.AirportData
		gotCachedData               bool
	)

	var (
		notifsSent = new(SentNotifications)
		cacheKey   = pollerParams.CacheKey
	)

	ticker := time.NewTicker(p.PollInterval())
	cleanup := func(err error) {
		pollerParams.Cleanup(err, ticker)
	}

	if p.UseDatastore() {
		p.LogDebug("checking for cached data...")
		cached, ok, cacheErr := p.fetchCacheEntry(ctx, cacheKey)
		if cacheErr != nil {
			p.LogError("error checking for cached data", zap.Error(cacheErr))
		} else if ok {
			flightData = cached.FlightData
			originInfo = cached.OriginData
			destinationInfo = cached.DestinationData
			notifsSent = cached.NotificationsSent

			gotCachedData = true
		}
	}

	if !gotCachedData {
		var (
			flightDataErr, infoErr error
		)

		flightData, flightDataErr = p.fetchFlight(ctx, flightId, flightIdType, fetchNearest)
		if flightDataErr != nil {
			cleanup(flightDataErr)
			return
		}

		originInfo, infoErr = p.fetchAirport(ctx, flightData.Origin.Identifiers.ICAO)
		if infoErr != nil {
			cleanup(infoErr)
			return
		}

		destinationInfo, infoErr = p.fetchAirport(ctx, flightData.Destination.Identifiers.ICAO)
		if infoErr != nil {
			cleanup(infoErr)
			return
		}
	}

	notifsSent.SetDisabled(p.FlightAwareConfig().Notifications)

	p.LogInfo(
		"starting flight data poller",
		zap.String("flight_identifier", flightData.Identifier.IATA),
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

				flightData, flightDataErr = p.fetchFlight(ctx, flightId, flightIdType, fetchNearest)
				if flightDataErr != nil {
					p.LogError("error fetching flight data", zap.Error(flightDataErr))
					continue
				}
			} else {
				isInitial = false
			}

			if p.UseDatastore() {
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
			}

			flightIdentifier := flightData.Identifier.IATA

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

				if !notifsSent.PreArrival && inPreArrivalWindow {
					msg.SetPlaintextTemplate(messages.PreArrivalAlertPlaintextTemplate)
					notifType = PreArrivalNotification

					break DetermineNotify
				} else if didSendGateDeparture && didSendTakeoff {
					continue
				}

				var (
					isGateDeparture, isTakeoff bool
				)

				switch {
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

				if didSendGateArrival {
					break DetermineNotify
				} else if didSendLanding {
					continue
				}

				var (
					isGateArrival, isLanding bool
				)

				switch {
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
