package flightaware

import (
	"fmt"
	"net/url"
	"time"

	"github.com/pkg/errors"

	"github.com/jalavosus/stuffnotifier/internal/utils"
)

const (
	FlightDateFormatLayout string = "2006-01-02"
)

var (
	ErrNoFlightsFound = errors.New("no flights found based on passed parameters")
	ErrNoFlightId     = errors.New("one of FlightId or FlightDesignator must be passed in FlightInformationParams")
)

type FlightInformationParams struct {
	// FlightAware AeroAPI ID for the flight.
	// Since this ID is tied to a single flight,
	// passing it almost certainly ensures that a
	// correct result will be returned.
	FlightId *string
	// "Human-readable" flight number, for example: UA1614 or DAL315.
	// If passed, a particular date should also be passed to ensure that
	// an accurate/desired result is returned.
	FlightDesignator *string
	// If querying flight info using the FlightDesignator param,
	// FlightDate should be passed as a time.Time in UTC.
	// Building a new time.Time using MakeFlightDateParam
	// or using the time.Date function is sufficient.
	FlightDate *time.Time
}

func (p FlightInformationParams) FlightIdentifier() string {
	if val, ok := utils.FromPointer(p.FlightId); ok {
		return val
	}

	if val, ok := utils.FromPointer(p.FlightDesignator); ok {
		return val
	}

	return ""
}

// UrlParams returns an API endpoint and constructed url query parameters
// based on the passed params.
func (p FlightInformationParams) UrlParams() (endpoint string, params url.Values, err error) {
	var (
		flightId string
		idType   IdentifierType
	)

	apiIdParam, hasApiId := utils.FromPointer(p.FlightId)
	designatorParam, hasDesignator := utils.FromPointer(p.FlightDesignator)

	switch {
	case hasApiId:
		flightId = apiIdParam
		idType = FaFlightIdIdent
	case hasDesignator:
		flightId = designatorParam
		idType = DesignatorIdent
	default:
		err = ErrNoFlightId
		return
	}

	params = url.Values{}
	params.Set(identTypeParam, idType.String())
	endpoint = fmt.Sprintf(flightsUri, flightId) + "?" + params.Encode()

	return
}

func (p FlightInformationParams) FlightDateParam() (time.Time, bool) {
	var (
		param    time.Time
		hasParam bool
	)

	if val, ok := utils.FromPointer(p.FlightDate); ok && !val.IsZero() {
		param = val
		hasParam = true
	}

	return param, hasParam
}

func (p FlightInformationParams) FetchByApiId() bool {
	_, ok := utils.FromPointer(p.FlightId)
	return ok
}

func (p FlightInformationParams) FetchByDate() bool {
	val, ok := utils.FromPointer(p.FlightDate)
	return ok && !val.IsZero()
}

func (p FlightInformationParams) FetchClosest() bool {
	if _, ok := p.FlightDateParam(); ok {
		return false
	}

	return true
}

// MakeFlightDateParam constructs a time.Time (in UTC)
// using the passed year, month, and day.
// `year` should be a "full" year, ex. 2006, 1997.
func MakeFlightDateParam(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
