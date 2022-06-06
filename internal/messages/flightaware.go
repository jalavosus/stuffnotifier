package messages

import (
	"text/template"

	"github.com/stoicturtle/stuffnotifier/pkg/flightaware"
)

type FlightAwareAlert struct {
	baseMessage
	IsGateDeparture   bool
	IsTakeoff         bool
	IsGateArrival     bool
	IsLanding         bool
	UseLocalTimezone  bool
	FlightNumber      string
	Origin            FlightAwareAirportInfo
	Destination       FlightAwareAirportInfo
	GateDepartureTime flightaware.FlightTimestamp
	TakeoffTime       flightaware.FlightTimestamp
	LandingTime       flightaware.FlightTimestamp
	GateArrivalTime   flightaware.FlightTimestamp
	plaintextTemplate *template.Template
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
		panic(err)
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
