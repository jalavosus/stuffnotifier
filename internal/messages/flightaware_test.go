package messages_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stoicturtle/stuffnotifier/internal/messages"
)

const timeFormat string = "2006-01-02 15:04:05 (-0700)"

func TestFlightAwareDepartureAlert_FormatPlaintext(t *testing.T) {
	type fields struct {
		IsGateDeparture      bool
		IsTakeoff            bool
		UseLocalTimezone     bool
		FlightNumber         string
		GateNumber           string
		OriginAirport        string
		DestinationAirport   string
		OriginTimezone       string
		DestinationTimezone  string
		GateDepartureTime    time.Time
		EstimatedTakeoffTime time.Time
		TakeoffTime          time.Time
		EstimatedLandingTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
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
			a := messages.FlightAwareDepartureAlert{
				IsGateDeparture:      tt.fields.IsGateDeparture,
				IsTakeoff:            tt.fields.IsTakeoff,
				UseLocalTimezone:     tt.fields.UseLocalTimezone,
				FlightNumber:         tt.fields.FlightNumber,
				GateNumber:           tt.fields.GateNumber,
				OriginAirport:        tt.fields.OriginAirport,
				DestinationAirport:   tt.fields.DestinationAirport,
				OriginTimezone:       tt.fields.OriginTimezone,
				DestinationTimezone:  tt.fields.DestinationTimezone,
				GateDepartureTime:    tt.fields.GateDepartureTime,
				EstimatedTakeoffTime: tt.fields.EstimatedTakeoffTime,
				TakeoffTime:          tt.fields.TakeoffTime,
				EstimatedLandingTime: tt.fields.EstimatedLandingTime,
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
