package flightaware

import (
	"encoding/json"
	"log"
	"strings"
	"time"
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
	Identifier FlightIdentifier
	Operator   FlightOperator
	// Unique identifier assigned by FlightAware for this specific flight.
	// If the flight is diverted, the new leg of the flight will have a different FlightId
	FlightId string
	// Bare flight number of the flight (ex. 2614)
	FlightNumber string
	// Aircraft registration (tail number) of the aircraft, when known
	Registration string
	// The identifier of the flight for Air Traffic Control purposes, when known and different than Identifier
	AtcIdent string
	// Unique identifier assigned by FlightAware for the previous flight of the aircraft serving this flight
	InboundFlightId string
	// List of any ICAO codeshares operating on this flight
	Codeshares []string
	// List of any IATA codeshares operating on this flight
	CodesharesIata []string
	// Indicates whether this flight is blocked from public viewing
	Blocked bool
	// Indicates whether this flight was diverted
	Diverted bool
	// Indicates that the flight is no longer being tracked by FlightAware.
	// There are a number of reasons this could happen including cancellation by the airline,
	// but that will not always be the case
	Cancelled bool
	// Indicates whether this flight has a flight plan, schedule, or other indication of intent available
	PositionOnly bool
	// Information for this flight's origin airport
	Origin FlightOriginDestinationData
	// Information for this flight's destination airport
	Destination FlightOriginDestinationData
	// Departure delay (in seconds) based on either actual or estimated gate departure time.
	// If gate time is unavailable then based on runway departure time.
	// A negative value indicates the flight has departed early
	DepartureDelay time.Duration
	// Arrival delay (in seconds) based on either actual or estimated gate arrival time.
	// If gate time is unavailable then based on runway arrival time.
	// A negative value indicates the flight has arrived early
	ArrivalDelay time.Duration
	// Runway-to-runway filed flight duration (in seconds)
	FlightTime time.Duration
	// Scheduled, estimated, and actual departure time from origin gate
	GateDepartureTime FlightTimestamp
	// Scheduled, estimated, and actual takeoff time from origin airport
	RunwayDepartureTime FlightTimestamp
	// Scheduled, estimated, and actual landing time at destination airport
	RunwayArrivalTime FlightTimestamp
	// Scheduled, estimated, and actual arrival time at destination airport gate
	GateArrivalTime FlightTimestamp
	// The percent completion of a flight, based on runway departure/arrival.
	// Null for en route position-only flights
	FlightProgress int
	// Human-readable summary of flight status
	Status string
	// Aircraft type will generally be ICAO code, but IATA code will be given when the ICAO code is not known
	AircraftType string
	// Planned flight distance (in statute miles) based on the filed route.
	// May vary from actual flown distance
	RouteDistance int64
	// Filed IFR airspeed (in knots)
	FiledAirspeed int64
	// Filed IFR altitude (in 100s of feet)
	FiledAltitude int64
	// The textual description of the flight's route
	Route string
	// The textual description of the flight's route, split by waypoint
	RouteWaypoints []string
	// Whether this is a commercial or general aviation flight
	FlightType FlightType
}

func (d *FlightData) UnmarshalJSON(data []byte) error {
	var raw rawFlightData

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*d = FlightData{
		Identifier:          raw.FlightIdentifier(),
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
		(*d).AtcIdent = *raw.AtcIdent
	}
	if raw.RouteDistance != nil {
		(*d).RouteDistance = *raw.RouteDistance
	}
	if raw.FiledAirspeed != nil {
		(*d).FiledAirspeed = *raw.FiledAirspeed
	}
	if raw.FiledAltitude != nil {
		(*d).FiledAltitude = *raw.FiledAltitude
	}

	return nil
}

type FlightIdentifier struct {
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
	Ident               string                         `json:"ident"`
	IdentIcao           string                         `json:"ident_icao"`
	IdentIata           string                         `json:"ident_iata"`
	FaFlightId          string                         `json:"fa_flight_id"`
	Operator            string                         `json:"operator"`
	OperatorIcao        string                         `json:"operator_icao"`
	OperatorIata        string                         `json:"operator_iata"`
	FlightNumber        string                         `json:"flight_number"`
	Registration        string                         `json:"registration"`
	AtcIdent            *string                        `json:"atc_ident,omitempty"`
	InboundFaFlightId   string                         `json:"inbound_fa_flight_id"`
	Codeshares          []string                       `json:"codeshares"`
	CodesharesIata      []string                       `json:"codeshares_iata"`
	Blocked             bool                           `json:"blocked"`
	Diverted            bool                           `json:"diverted"`
	Cancelled           bool                           `json:"cancelled"`
	PositionOnly        bool                           `json:"position_only"`
	Origin              rawFlightOriginDestinationData `json:"origin"`
	Destination         rawFlightOriginDestinationData `json:"destination"`
	DepartureDelay      int64                          `json:"departure_delay"`
	ArrivalDelay        int64                          `json:"arrival_delay"`
	FlightEte           int64                          `json:"flight_ete"`
	ScheduledOut        *string                        `json:"scheduled_out,omitempty"`
	EstimatedOut        *string                        `json:"estimated_out,omitempty"`
	ActualOut           *string                        `json:"actual_out,omitempty"`
	ScheduledOff        *string                        `json:"scheduled_off,omitempty"`
	EstimatedOff        *string                        `json:"estimated_off,omitempty"`
	ActualOff           *string                        `json:"actual_off,omitempty"`
	ScheduledOn         *string                        `json:"scheduled_on,omitempty"`
	EstimatedOn         *string                        `json:"estimated_on,omitempty"`
	ActualOn            *string                        `json:"actual_on,omitempty"`
	ScheduledIn         *string                        `json:"scheduled_in,omitempty"`
	EstimatedIn         *string                        `json:"estimated_in,omitempty"`
	ActualIn            *string                        `json:"actual_in,omitempty"`
	GateOrigin          string                         `json:"gate_origin"`
	GateDestination     *string                        `json:"gate_destination,omitempty"`
	TerminalOrigin      string                         `json:"terminal_origin"`
	TerminalDestination *string                        `json:"terminal_destination,omitempty"`
	FlightProgress      int                            `json:"flight_progress"`
	Status              string                         `json:"status"`
	AircraftType        string                         `json:"aircraft_type"`
	RouteDistance       *int64                         `json:"route_distance,omitempty"`
	FiledAirspeed       *int64                         `json:"filed_airspeed,omitempty"`
	FiledAltitude       *int64                         `json:"filed_altitude,omitempty"`
	Route               string                         `json:"route"`
	Type                string                         `json:"type"`
}

func (d rawFlightData) FlightIdentifier() FlightIdentifier {
	return FlightIdentifier{
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
	// IATA, ICAO, and LID identifier codes for the airport
	Identifiers AirportIdentifiers
	// Common name for the airport
	Name string
	// Type of airport
	AirportType AirportType
	// Height above Mean Sea Level (MSL) (in feet)
	Elevation int64
	// Closest city to this airport
	City string
	// State/province where the airport resides if applicable.
	// For US states this will be their 2-letter code;
	// for provinces or other entities, it will be the full name
	State string
	// Airport's latitude, generally the center point of the airport
	Latitude float64
	// Airport's longitude, generally the center point of the airport
	Longitude float64
	// Applicable timezone for the airport, in the TZ database format
	Timezone string
	// 2-letter code of country where the airport resides (ISO 3166-1 alpha-2)
	CountryCode string
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
	AirportCode string  `json:"airport_code"`
	CodeIcao    string  `json:"code_icao"`
	CodeIata    string  `json:"code_iata"`
	CodeLid     string  `json:"code_lid"`
	Name        string  `json:"name"`
	Type        string  `json:"type,omitempty"`
	Elevation   int64   `json:"elevation"`
	City        string  `json:"city"`
	State       string  `json:"state"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Timezone    string  `json:"timezone"`
	CountryCode string  `json:"country_code"`
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
			log.Panicln(parseErr)
		}
	}

	if checkValidString(estimated) {
		parsedEstimated, parseErr = time.Parse(apiTimeFormat, *estimated)
		if parseErr != nil {
			log.Panicln(parseErr)
		}
	}

	if checkValidString(actual) {
		parsedActual, parseErr = time.Parse(apiTimeFormat, *actual)
		if parseErr != nil {
			log.Panicln(parseErr)
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

// {
// "atc_ident": null,
// "baggage_claim": "6",
// "seats_cabin_business": 21,
// "seats_cabin_coach": 253,
// "seats_cabin_first": 44,
// },
