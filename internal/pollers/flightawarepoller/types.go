package flightawarepoller

import (
	"time"

	"github.com/stoicturtle/stuffnotifier/internal/pollers/poller"
	"github.com/stoicturtle/stuffnotifier/pkg/flightaware"
)

type CacheEntry struct {
	FlightData        *flightaware.FlightData
	OriginData        *flightaware.AirportData
	DestinationData   *flightaware.AirportData
	NotificationsSent *SentNotifications
	InternalId        string
	FlightId          string
	RecipientConfig   poller.RecipientConfig
	Notifications     flightaware.NotificationsConfig
	PollInterval      time.Duration
}

type SentNotifications struct {
	GateDeparture bool
	Takeoff       bool
	Landing       bool
	GateArrival   bool
	PreArrival    bool
	PreDeparture  bool
	BaggageClaim  bool
}

func (s *SentNotifications) SetDisabled(notifsConfig flightaware.NotificationsConfig) {
	if !notifsConfig.GateDeparture {
		s.GateDeparture = true
	}
	if !notifsConfig.Takeoff {
		s.Takeoff = true
	}
	if !notifsConfig.Landing {
		s.Landing = true
	}
	if !notifsConfig.GateArrival {
		s.GateArrival = true
	}
	if !notifsConfig.PreArrival.Enabled {
		s.PreArrival = true
	}
	if !notifsConfig.PreDeparture.Enabled {
		s.PreDeparture = true
	}
}

func (s *SentNotifications) SetSent(notifType notificationType) {
	var allSentNotifs = []notificationType{notifType}

	if otherNotifs, ok := notifTypesMap[notifType]; ok {
		allSentNotifs = append(allSentNotifs, otherNotifs...)
	}

	s.setSent(allSentNotifs...)
}

func (s *SentNotifications) setSent(notifTypes ...notificationType) {
	for _, n := range notifTypes {
		switch n {
		case GateDepartureNotification:
			s.GateDeparture = true
		case TakeoffNotification:
			s.Takeoff = true
		case LandingNotification:
			s.Landing = true
		case GateArrivalNotification:
			s.GateArrival = true
		case PreDepartureNotification:
			s.PreDeparture = true
		case PreArrivalNotification:
			s.PreArrival = true
		case BaggageClaimNotification:
			s.BaggageClaim = true
		}
	}
}

func (s SentNotifications) SentAll() bool {
	return s.GateDeparture &&
		s.Takeoff &&
		s.Landing &&
		s.GateArrival &&
		s.PreArrival &&
		s.PreDeparture
}

var notifTypesMap = map[notificationType][]notificationType{
	GateDepartureNotification: {PreDepartureNotification},
	TakeoffNotification:       {PreDepartureNotification, GateDepartureNotification},
	LandingNotification:       {PreDepartureNotification, GateDepartureNotification, TakeoffNotification, PreArrivalNotification},
	GateArrivalNotification:   {PreDepartureNotification, GateDepartureNotification, TakeoffNotification, PreArrivalNotification, LandingNotification},
	PreArrivalNotification:    {PreDepartureNotification, PreArrivalNotification},
}

type pastFlightEvents struct {
	pending      bool
	DepartedGate bool
	Takeoff      bool
	Landed       bool
	ArrivedGate  bool
}

func newPastFlightEvents(flightData *flightaware.FlightData) pastFlightEvents {
	return pastFlightEvents{
		pending:      flightData.GateDepartureTime.Scheduled.UTC().After(time.Now().UTC()),
		DepartedGate: validTimestampActual(flightData.GateDepartureTime),
		Takeoff:      validTimestampActual(flightData.RunwayDepartureTime),
		Landed:       validTimestampActual(flightData.RunwayArrivalTime),
		ArrivedGate:  validTimestampActual(flightData.GateArrivalTime),
	}
}

func (e pastFlightEvents) Pending() bool {
	return e.pending
}

func (e pastFlightEvents) InProgress() bool {
	return (e.DepartedGate || e.Takeoff) && !(e.Landed || e.ArrivedGate)
}

func (e pastFlightEvents) ArrivedDestination() bool {
	return e.Landed || e.ArrivedGate
}
