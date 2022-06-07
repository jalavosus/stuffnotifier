package flightaware

import (
	"encoding/json"
	"strings"
	"time"

	"go.uber.org/zap"
)

const (
	apiTimeFormat string = "2006-01-02T15:04:05Z"
)

type ApiResponse struct {
	Links struct {
		Next string `json:"next"`
	} `json:"links"`
	NumPages int `json:"num_pages"`
}

type ApiResponseError struct {
	Title  string
	Reason string
	Detail string
	Status string
}

type FlightDataResponse struct {
	ApiResponse
	Flights []FlightData `json:"flights"`
}

type FlightData struct {
	GateArrivalTime     FlightTimestamp
	GateDepartureTime   FlightTimestamp
	RunwayDepartureTime FlightTimestamp
	RunwayArrivalTime   FlightTimestamp
	Origin              FlightOriginDestinationData
	Destination         FlightOriginDestinationData
	Operator            FlightOperator
	Identifiers         FlightIdentifiers
	InboundFlightId     string
	FlightNumber        string
	Route               string
	FlightId            string
	AircraftType        string
	AtcIdent            string
	Registration        string
	Status              string
	RouteWaypoints      []string
	CodesharesIata      []string
	Codeshares          []string
	ArrivalDelay        time.Duration
	FlightTime          time.Duration
	FiledAltitude       int64
	FlightProgress      int
	FiledAirspeed       int64
	DepartureDelay      time.Duration
	RouteDistance       int64
	PositionOnly        bool
	Cancelled           bool
	Diverted            bool
	Blocked             bool
	FlightType          FlightType
}

func (d *FlightData) UnmarshalJSON(data []byte) error {
	var raw rawFlightData

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*d = FlightData{
		Identifiers:         raw.FlightIdentifier(),
		Operator:            raw.FlightOperator(),
		FlightId:            raw.FaFlightId,
		FlightNumber:        raw.FlightNumber,
		Registration:        raw.Registration,
		InboundFlightId:     raw.InboundFaFlightId,
		Codeshares:          raw.Codeshares,
		CodesharesIata:      raw.CodesharesIata,
		Blocked:             raw.Blocked,
		Diverted:            raw.Diverted,
		Cancelled:           raw.Cancelled,
		PositionOnly:        raw.PositionOnly,
		Origin:              raw.FlightOrigin(),
		Destination:         raw.FlightDestination(),
		DepartureDelay:      time.Duration(raw.DepartureDelay) * time.Second,
		ArrivalDelay:        time.Duration(raw.ArrivalDelay) * time.Second,
		GateDepartureTime:   raw.GateDepartureTime(),
		RunwayDepartureTime: raw.RunwayDepartureTime(),
		GateArrivalTime:     raw.GateArrivalTime(),
		RunwayArrivalTime:   raw.RunwayArrivalTime(),
		FlightProgress:      raw.FlightProgress,
		Status:              raw.Status,
		AircraftType:        raw.AircraftType,
		Route:               raw.Route,
		FlightType:          flightTypeFromString(raw.Type),
	}

	if raw.AtcIdent != nil {
		d.AtcIdent = *raw.AtcIdent
	}
	if raw.RouteDistance != nil {
		d.RouteDistance = *raw.RouteDistance
	}
	if raw.FiledAirspeed != nil {
		d.FiledAirspeed = *raw.FiledAirspeed
	}
	if raw.FiledAltitude != nil {
		d.FiledAltitude = *raw.FiledAltitude
	}

	return nil
}

type FlightIdentifiers struct {
	// Either the operator code followed by the flight number for the flight (for commercial flights)
	// or the aircraft's registration (for general aviation)
	Identifier string
	// The ICAO operator code followed by the flight number for the flight (for commercial flights)
	ICAO string
	// The IATA operator code followed by the flight number for the flight (for commercial flights)
	IATA string
}

type FlightOperator struct {
	// ICAO code, if exists, of the operator of the flight, otherwise the IATA code
	Operator string
	// ICAO code of the operator of the flight
	ICAO string
	// IATA code of the operator of the flight
	IATA string
}

// FlightOriginDestinationData contains basic data about an airport a flight is flying to or from.
type FlightOriginDestinationData struct {
	Identifiers AirportIdentifiers
	// Aiport's website, if known
	InfoUrl string
	// Departure/Arrival gate, if known
	Gate string
	// Departure/Arrival terminal, if known
	Terminal string
}

// AirportIdentifiers contains various identifier codes for airports.
type AirportIdentifiers struct {
	Code string
	ICAO string
	IATA string
	LID  string
}

// FlightTimestamp contains the estimated, scheduled, and actual
// timestamps (stored in UTC) of a flight's departure and arrival.
type FlightTimestamp struct {
	// Scheduled arrival/departure time
	Scheduled time.Time
	// Estimated arrival/departure time
	Estimated time.Time
	// Actual arrival/departure time
	Actual time.Time
}

//go:generate stringer -type FlightType,IdentifierType,AirportType -linecomment -output models_string.go

type (
	FlightType     uint8
	IdentifierType uint8
	AirportType    uint8
)

func flightTypeFromString(s string) FlightType {
	switch strings.ToLower(s) {
	case strings.ToLower(GeneralAviation.String()):
		return GeneralAviation
	case strings.ToLower(Airline.String()):
		return Airline
	default:
		return UnknownFlightType
	}
}

func airportTypeFromString(s string) AirportType {
	switch strings.ToLower(s) {
	case strings.ToLower(Airport.String()):
		return Airport
	case strings.ToLower(Heliport.String()):
		return Heliport
	case strings.ToLower(SeaplaneBase.String()):
		return SeaplaneBase
	case strings.ToLower(Ultralight.String()):
		return Ultralight
	case strings.ToLower(Stolport.String()):
		return Stolport
	case strings.ToLower(Gliderport.String()):
		return Gliderport
	case strings.ToLower(BalloonPort.String()):
		return BalloonPort
	default:
		return UnknownAirportType
	}
}

const (
	UnknownFlightType FlightType = iota // unknown
	GeneralAviation                     // General_Aviation
	Airline                             // Airline
)

const (
	DesignatorIdent   IdentifierType = iota // designator
	RegistrationIdent                       // registration
	FaFlightIdIdent                         // fa_flight_id
)

const (
	UnknownAirportType AirportType = iota // unknown
	Airport                               // Airport
	Heliport                              // Heliport
	SeaplaneBase                          // Seaplane Base
	Ultralight                            // Ultralight
	Stolport                              // Stolport
	Gliderport                            // Gliderport
	BalloonPort                           // BalloonPort
)

type rawFlightData struct {
	TerminalDestination *string                        `json:"terminal_destination,omitempty"`
	EstimatedIn         *string                        `json:"estimated_in,omitempty"`
	ScheduledIn         *string                        `json:"scheduled_in,omitempty"`
	ActualOn            *string                        `json:"actual_on,omitempty"`
	EstimatedOn         *string                        `json:"estimated_on,omitempty"`
	ActualIn            *string                        `json:"actual_in,omitempty"`
	ActualOff           *string                        `json:"actual_off,omitempty"`
	EstimatedOff        *string                        `json:"estimated_off,omitempty"`
	ScheduledOff        *string                        `json:"scheduled_off,omitempty"`
	AtcIdent            *string                        `json:"atc_ident,omitempty"`
	ActualOut           *string                        `json:"actual_out,omitempty"`
	ScheduledOn         *string                        `json:"scheduled_on,omitempty"`
	FiledAltitude       *int64                         `json:"filed_altitude,omitempty"`
	FiledAirspeed       *int64                         `json:"filed_airspeed,omitempty"`
	RouteDistance       *int64                         `json:"route_distance,omitempty"`
	GateDestination     *string                        `json:"gate_destination,omitempty"`
	ScheduledOut        *string                        `json:"scheduled_out,omitempty"`
	EstimatedOut        *string                        `json:"estimated_out,omitempty"`
	Destination         rawFlightOriginDestinationData `json:"destination"`
	Origin              rawFlightOriginDestinationData `json:"origin"`
	TerminalOrigin      string                         `json:"terminal_origin"`
	Type                string                         `json:"type"`
	Status              string                         `json:"status"`
	Route               string                         `json:"route"`
	InboundFaFlightId   string                         `json:"inbound_fa_flight_id"`
	Registration        string                         `json:"registration"`
	FlightNumber        string                         `json:"flight_number"`
	OperatorIata        string                         `json:"operator_iata"`
	OperatorIcao        string                         `json:"operator_icao"`
	Operator            string                         `json:"operator"`
	FaFlightId          string                         `json:"fa_flight_id"`
	IdentIata           string                         `json:"ident_iata"`
	IdentIcao           string                         `json:"ident_icao"`
	GateOrigin          string                         `json:"gate_origin"`
	Ident               string                         `json:"ident"`
	AircraftType        string                         `json:"aircraft_type"`
	CodesharesIata      []string                       `json:"codeshares_iata"`
	Codeshares          []string                       `json:"codeshares"`
	FlightEte           int64                          `json:"flight_ete"`
	ArrivalDelay        int64                          `json:"arrival_delay"`
	DepartureDelay      int64                          `json:"departure_delay"`
	FlightProgress      int                            `json:"flight_progress"`
	Cancelled           bool                           `json:"cancelled"`
	Diverted            bool                           `json:"diverted"`
	Blocked             bool                           `json:"blocked"`
	PositionOnly        bool                           `json:"position_only"`
}

func (d rawFlightData) FlightIdentifier() FlightIdentifiers {
	return FlightIdentifiers{
		Identifier: d.Ident,
		ICAO:       d.IdentIcao,
		IATA:       d.IdentIata,
	}
}

func (d rawFlightData) FlightOperator() FlightOperator {
	return FlightOperator{
		Operator: d.Operator,
		ICAO:     d.OperatorIcao,
		IATA:     d.OperatorIata,
	}
}

func (d rawFlightData) GateDepartureTime() FlightTimestamp {
	return parseFlightTimestamp(d.ScheduledOut, d.EstimatedOut, d.ActualOut)
}

func (d rawFlightData) RunwayDepartureTime() FlightTimestamp {
	return parseFlightTimestamp(d.ScheduledOff, d.EstimatedOff, d.ActualOff)
}

func (d rawFlightData) GateArrivalTime() FlightTimestamp {
	return parseFlightTimestamp(d.ScheduledIn, d.EstimatedIn, d.ActualIn)
}

func (d rawFlightData) RunwayArrivalTime() FlightTimestamp {
	return parseFlightTimestamp(d.ScheduledOn, d.EstimatedOn, d.ActualOn)
}

func (d rawFlightData) FlightOrigin() FlightOriginDestinationData {
	data := d.Origin.FlightOriginDestinationData()
	data.Gate = d.GateOrigin
	data.Terminal = d.TerminalOrigin

	return data
}

func (d rawFlightData) FlightDestination() FlightOriginDestinationData {
	data := d.Destination.FlightOriginDestinationData()

	if d.GateDestination != nil {
		data.Gate = *d.GateDestination
	}

	if d.TerminalDestination != nil {
		data.Terminal = *d.TerminalDestination
	}

	return data
}

type rawFlightOriginDestinationData struct {
	Code           string `json:"code"`
	CodeIcao       string `json:"code_icao"`
	CodeIata       string `json:"code_iata"`
	CodeLid        string `json:"code_lid"`
	AirportInfoUrl string `json:"airport_info_url"`
}

func (d rawFlightOriginDestinationData) FlightOriginDestinationData() FlightOriginDestinationData {
	return FlightOriginDestinationData{
		Identifiers: AirportIdentifiers{
			Code: d.Code,
			IATA: d.CodeIata,
			ICAO: d.CodeIcao,
			LID:  d.CodeLid,
		},
		InfoUrl: d.AirportInfoUrl,
	}
}

type AirportData struct {
	Identifiers AirportIdentifiers
	Name        string
	Timezone    string
	CountryCode string
	City        string
	State       string
	Latitude    float64
	Longitude   float64
	Elevation   int64
	AirportType AirportType
}

func (d *AirportData) UnmarshalJSON(data []byte) error {
	var raw rawAirportData
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*d = AirportData{
		Identifiers: raw.AirportIdentifiers(),
		Name:        raw.Name,
		AirportType: airportTypeFromString(raw.Type),
		Elevation:   raw.Elevation,
		City:        raw.City,
		State:       raw.State,
		Latitude:    raw.Latitude,
		Longitude:   raw.Longitude,
		Timezone:    raw.Timezone,
		CountryCode: raw.CountryCode,
	}

	return nil
}

type rawAirportData struct {
	CountryCode string  `json:"country_code"`
	CodeIcao    string  `json:"code_icao"`
	CodeIata    string  `json:"code_iata"`
	CodeLid     string  `json:"code_lid"`
	Name        string  `json:"name"`
	Type        string  `json:"type,omitempty"`
	AirportCode string  `json:"airport_code"`
	City        string  `json:"city"`
	State       string  `json:"state"`
	Timezone    string  `json:"timezone"`
	Longitude   float64 `json:"longitude"`
	Latitude    float64 `json:"latitude"`
	Elevation   int64   `json:"elevation"`
}

func (d rawAirportData) AirportIdentifiers() AirportIdentifiers {
	return AirportIdentifiers{
		Code: d.AirportCode,
		IATA: d.CodeIata,
		ICAO: d.CodeIcao,
		LID:  d.CodeLid,
	}
}

func parseFlightTimestamp(scheduled, estimated, actual *string) FlightTimestamp {
	var (
		parsedScheduled, parsedEstimated, parsedActual time.Time
		parseErr                                       error
	)

	if checkValidString(scheduled) {
		parsedScheduled, parseErr = time.Parse(apiTimeFormat, *scheduled)
		if parseErr != nil {
			logger.Panic("error parsing time string", zap.String("input", *scheduled), zap.Error(parseErr))
		}
	}

	if checkValidString(estimated) {
		parsedEstimated, parseErr = time.Parse(apiTimeFormat, *estimated)
		if parseErr != nil {
			logger.Panic("error parsing time string", zap.String("input", *estimated), zap.Error(parseErr))
		}
	}

	if checkValidString(actual) {
		parsedActual, parseErr = time.Parse(apiTimeFormat, *actual)
		if parseErr != nil {
			logger.Panic("error parsing time string", zap.String("input", *actual), zap.Error(parseErr))
		}
	}

	return FlightTimestamp{
		Scheduled: parsedScheduled,
		Estimated: parsedEstimated,
		Actual:    parsedActual,
	}
}

func checkValidString(s *string) bool {
	return s != nil && *s != ""
}
