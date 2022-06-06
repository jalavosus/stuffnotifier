package messages_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stoicturtle/stuffnotifier/internal/messages"
	"github.com/stoicturtle/stuffnotifier/pkg/flightaware"
)

const timeFormat string = "2006-01-02 15:04:05 (-0700)"

func TestFlightAwareDepartureAlert_FormatPlaintext(t *testing.T) {
	type fields struct {
		EstimatedLandingTime time.Time
		TakeoffTime          time.Time
		EstimatedTakeoffTime time.Time
		GateDepartureTime    time.Time
		DestinationAirport   string
		OriginAirport        string
		GateNumber           string
		OriginTimezone       string
		DestinationTimezone  string
		FlightNumber         string
		UseLocalTimezone     bool
		IsTakeoff            bool
		IsGateDeparture      bool
	}
	tests := []struct {
		fields fields
		name   string
	}{
		{
			name: "test takeoff/local timezone",
			fields: fields{
				IsGateDeparture:      false,
				IsTakeoff:            true,
				UseLocalTimezone:     true,
				FlightNumber:         "UA2614",
				GateNumber:           "37",
				OriginAirport:        "LAX",
				DestinationAirport:   "EWR",
				OriginTimezone:       "America/Los_Angeles",
				DestinationTimezone:  "America/New_York",
				GateDepartureTime:    mustParseTime(t, "2022-05-31 12:35:00 (-0700)"),
				EstimatedTakeoffTime: mustParseTime(t, "2022-05-31 12:45:00 (-0700)"),
				TakeoffTime:          mustParseTime(t, "2022-05-31 12:40:00 (-0700)"),
				EstimatedLandingTime: mustParseTime(t, "2022-05-31 20:53:00 (-0400)"),
			},
		},
		{
			name: "test takeoff/utc",
			fields: fields{
				IsGateDeparture:      false,
				IsTakeoff:            true,
				UseLocalTimezone:     false,
				FlightNumber:         "UA2614",
				GateNumber:           "37",
				OriginAirport:        "LAX",
				DestinationAirport:   "EWR",
				OriginTimezone:       "America/Los_Angeles",
				DestinationTimezone:  "America/New_York",
				GateDepartureTime:    mustParseTime(t, "2022-05-31 12:35:00 (-0700)"),
				EstimatedTakeoffTime: mustParseTime(t, "2022-05-31 12:45:00 (-0700)"),
				TakeoffTime:          mustParseTime(t, "2022-05-31 12:40:00 (-0700)"),
				EstimatedLandingTime: mustParseTime(t, "2022-05-31 20:53:00 (-0400)"),
			},
		},
		{
			name: "test gate departure/local timezone",
			fields: fields{
				IsGateDeparture:      true,
				IsTakeoff:            false,
				UseLocalTimezone:     true,
				FlightNumber:         "UA2614",
				GateNumber:           "37",
				OriginAirport:        "LAX",
				DestinationAirport:   "EWR",
				OriginTimezone:       "America/Los_Angeles",
				DestinationTimezone:  "America/New_York",
				GateDepartureTime:    mustParseTime(t, "2022-05-31 12:35:00 (-0700)"),
				EstimatedTakeoffTime: mustParseTime(t, "2022-05-31 12:45:00 (-0700)"),
				TakeoffTime:          mustParseTime(t, "2022-05-31 12:40:00 (-0700)"),
				EstimatedLandingTime: mustParseTime(t, "2022-05-31 20:53:00 (-0400)"),
			},
		},
		{
			name: "test gate departure/utc",
			fields: fields{
				IsGateDeparture:      true,
				IsTakeoff:            false,
				UseLocalTimezone:     false,
				FlightNumber:         "UA2614",
				GateNumber:           "37",
				OriginAirport:        "LAX",
				DestinationAirport:   "EWR",
				OriginTimezone:       "America/Los_Angeles",
				DestinationTimezone:  "America/New_York",
				GateDepartureTime:    mustParseTime(t, "2022-05-31 12:35:00 (-0700)"),
				EstimatedTakeoffTime: mustParseTime(t, "2022-05-31 12:45:00 (-0700)"),
				TakeoffTime:          mustParseTime(t, "2022-05-31 12:40:00 (-0700)"),
				EstimatedLandingTime: mustParseTime(t, "2022-05-31 20:53:00 (-0400)"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := messages.FlightAwareAlert{
				IsGateDeparture:  tt.fields.IsGateDeparture,
				IsTakeoff:        tt.fields.IsTakeoff,
				UseLocalTimezone: tt.fields.UseLocalTimezone,
				FlightNumber:     tt.fields.FlightNumber,
				Origin: messages.FlightAwareAirportInfo{
					Airport:  tt.fields.OriginAirport,
					Gate:     tt.fields.GateNumber,
					Timezone: tt.fields.OriginTimezone,
				},
				Destination: messages.FlightAwareAirportInfo{
					Airport:  tt.fields.DestinationAirport,
					Timezone: tt.fields.DestinationTimezone,
				},
				GateDepartureTime: flightaware.FlightTimestamp{Actual: tt.fields.GateDepartureTime},
				TakeoffTime:       flightaware.FlightTimestamp{Estimated: tt.fields.EstimatedTakeoffTime, Actual: tt.fields.TakeoffTime},
				LandingTime:       flightaware.FlightTimestamp{Estimated: tt.fields.EstimatedLandingTime},
			}

			_ = a.FormatPlaintext()
		})
	}
}

func mustParseTime(t *testing.T, s string) time.Time {
	t.Helper()

	parsed, err := time.Parse(timeFormat, s)
	assert.NoErrorf(t, err, "error parsing %[1]s into time.Time", s)

	return parsed.UTC()
}
