package flightaware

import (
	"sort"
	"time"

	"github.com/jalavosus/stuffnotifier/internal/utils"
)

func SortFlights(flights []FlightData) (sorted []FlightData) {
	sorted = make([]FlightData, len(flights))

	sort.Slice(flights, func(i, j int) bool {
		flightAScheduled := scheduledGateDeparture(flights[i]).UTC()
		flightBScheduled := scheduledGateDeparture(flights[j]).UTC()

		return flightAScheduled.Before(flightBScheduled)
	})

	copy(sorted, flights)

	return
}

func FindFlightByDate(flights []FlightData, flightDate time.Time) (found *FlightData, ok bool) {
	flights = SortFlights(flights)

	flightDate = flightDate.UTC()

	for i := range flights {
		flightTime := scheduledGateDeparture(flights[i])

		if compareDates(flightDate, flightTime, true) {
			found = &flights[i]
			ok = true
			break
		}
	}

	return
}

func FindClosestFlight(flights []FlightData) (found *FlightData, ok bool) {
	flights = SortFlights(flights)

	now := utils.ToLocalTime(time.Now())

	for i := range flights {
		flightTime := utils.ToLocalTime(scheduledGateDeparture(flights[i]))

		if compareDates(now, flightTime, false) {
			found = &flights[i]
			ok = true
			break
		}
	}

	return
}

func FindFlightByApiId(flights []FlightData, apiId string) (found *FlightData, ok bool) {
	for i := range flights {
		if flights[i].FlightId == apiId {
			found = &flights[i]
			ok = true
			break
		}
	}

	return
}

func FindFlightFromParams(flights []FlightData, params FlightInformationParams) (found *FlightData, ok bool) {
	switch {
	case params.FetchByApiId():
		flightId, _ := utils.FromPointer(params.FlightId)
		found, ok = FindFlightByApiId(flights, flightId)
	case params.FetchByDate():
		desiredDate, _ := params.FlightDateParam()
		found, ok = FindFlightByDate(flights, desiredDate)
	default:
		found, ok = FindClosestFlight(flights)
	}

	return
}

func compareDates(t1 time.Time, t2 time.Time, wantEq bool) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()

	if wantEq {
		return y2 == y1 && m2 == m1 && d2 == d1
	}

	return y2 >= y1 && m2 >= m1 && d2 >= d1
}

func scheduledGateDeparture(flight FlightData) time.Time {
	return flight.GateDepartureTime.Scheduled
}
