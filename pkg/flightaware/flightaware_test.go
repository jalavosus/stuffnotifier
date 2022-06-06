package flightaware_test

import (
	"context"
	"testing"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/assert"

	"github.com/stoicturtle/stuffnotifier/internal/authdata"
	"github.com/stoicturtle/stuffnotifier/pkg/flightaware"
)

func TestFlightInformation(t *testing.T) {
	// t.Skip()

	type args struct {
		flightIdentifier string
		identifierType   flightaware.IdentifierType
	}

	type want struct {
		originCodeIcao string
		originCodeIata string
		destCodeIcao   string
		destCodeIata   string
		flightType     flightaware.FlightType
	}

	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "test valid flight (UA2614)",
			args: args{
				flightIdentifier: "UA2614",
				identifierType:   flightaware.DesignatorIdent,
			},
			want: want{
				originCodeIcao: "KLAX",
				originCodeIata: "LAX",
				destCodeIcao:   "KEWR",
				destCodeIata:   "EWR",
				flightType:     flightaware.Airline,
			},
			wantErr: false,
		},
		{
			name: "test valid flight (UAL2614)",
			args: args{
				flightIdentifier: "UA2614",
				identifierType:   flightaware.DesignatorIdent,
			},
			want: want{
				originCodeIcao: "KLAX",
				originCodeIata: "LAX",
				destCodeIcao:   "KEWR",
				destCodeIata:   "EWR",
				flightType:     flightaware.Airline,
			},
			wantErr: false,
		},
		{
			name: "test valid flight (UAL2614-1653380561-fa-0007)",
			args: args{
				flightIdentifier: "UAL2614-1653380561-fa-0007",
				identifierType:   flightaware.FaFlightIdIdent,
			},
			want: want{
				originCodeIcao: "KLAX",
				originCodeIata: "LAX",
				destCodeIcao:   "KEWR",
				destCodeIata:   "EWR",
				flightType:     flightaware.Airline,
			},
			wantErr: false,
		},
	}

	rootCtx := context.Background()
	authData, _ := authdata.FlightAwareAPIAuth()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(rootCtx, 10*time.Second)
			defer cancel()

			got, err := flightaware.FlightInformation(
				ctx,
				authData,
				tt.args.flightIdentifier,
				tt.args.identifierType,
			)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NotEmpty(t, got.Flights)

			flightData := got.Flights[0]
			assert.Equal(t, tt.want.originCodeIcao, flightData.Origin.Identifiers.ICAO)
			assert.Equal(t, tt.want.originCodeIata, flightData.Origin.Identifiers.IATA)
			assert.Equal(t, tt.want.destCodeIcao, flightData.Destination.Identifiers.ICAO)
			assert.Equal(t, tt.want.destCodeIata, flightData.Destination.Identifiers.IATA)
			assert.Equal(t, tt.want.flightType, flightData.FlightType)
		})
	}
}

func TestAirportInformation(t *testing.T) {
	type args struct {
		identifier string
	}

	type want struct {
		codeIcao string
		codeIata string
		state    string
	}

	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name:    "LAX",
			args:    args{"LAX"},
			want:    want{"KLAX", "LAX", "CA"},
			wantErr: false,
		},
		{
			name:    "KLAX",
			args:    args{"KLAX"},
			want:    want{"KLAX", "LAX", "CA"},
			wantErr: false,
		},
		{
			name:    "YYZ",
			args:    args{"YYZ"},
			want:    want{"CYYZ", "YYZ", "Ontario"},
			wantErr: false,
		},
		{
			name:    "KYYZ",
			args:    args{"KYYZ"},
			want:    want{"", "", ""},
			wantErr: false,
		},
		{
			name:    "CYYZ",
			args:    args{"CYYZ"},
			want:    want{"CYYZ", "YYZ", "Ontario"},
			wantErr: false,
		},
	}

	rootCtx := context.Background()
	authData, _ := authdata.FlightAwareAPIAuth()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(rootCtx, 10*time.Second)
			defer cancel()

			got, err := flightaware.AirportInformation(
				ctx,
				authData,
				tt.args.identifier,
			)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.codeIcao, got.Identifiers.ICAO)
			assert.Equal(t, tt.want.codeIata, got.Identifiers.IATA)
			assert.Equal(t, tt.want.state, got.State)
		})
	}
}
