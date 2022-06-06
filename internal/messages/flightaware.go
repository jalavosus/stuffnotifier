package messages

import (
	"text/template"

	"go.uber.org/zap"

	"github.com/jalavosus/stuffnotifier/pkg/flightaware"
)

type FlightAwareAlert struct {
	baseMessage
	GateArrivalTime   flightaware.FlightTimestamp
	LandingTime       flightaware.FlightTimestamp
	TakeoffTime       flightaware.FlightTimestamp
	GateDepartureTime flightaware.FlightTimestamp
	plaintextTemplate *template.Template
	Origin            FlightAwareAirportInfo
	Destination       FlightAwareAirportInfo
	FlightNumber      string
	IsLanding         bool
	IsGateArrival     bool
	IsTakeoff         bool
	IsGateDeparture   bool
	UseLocalTimezone  bool
}

type FlightAwareAirportInfo struct {
	Airport  string
	Timezone string
	Gate     string
	Terminal string
}

func (a FlightAwareAlert) FormatPlaintext() string {
	msg, err := a.format(a.PlaintextTemplate(), a)
	if err != nil {
		logger.Panic("error formatting FlightAwareAlert plaintext template", zap.Error(err))
	}

	return msg
}

func (a FlightAwareAlert) FormatMarkdown() string {
	// msg, err := a.format(a.MarkdownTemplate(), a)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// return msg

	return ""
}

func (a FlightAwareAlert) PlaintextTemplate() *template.Template {
	return a.plaintextTemplate
}

func (a *FlightAwareAlert) SetPlaintextTemplate(tmpl *template.Template) *FlightAwareAlert {
	a.plaintextTemplate = tmpl
	return a
}

func (a FlightAwareAlert) MarkdownTemplate() *template.Template {
	return nil
}
