package flightaware

import (
	"context"

	"github.com/jalavosus/stuffnotifier/internal/logging"
	"github.com/jalavosus/stuffnotifier/pkg/authdata"
)

var logger = logging.NewLogger()

const (
	apiKeyHeader    string = "x-apikey"
	baseApiEndpoint string = "aeroapi.flightaware.com/aeroapi"
)

const (
	identTypeParam string = "ident_type"
)

const (
	flightsUri  string = "flights/%s"
	airportsUri string = "airports/%s"
)

func GetFlightInformation(ctx context.Context, authData authdata.AuthData, params FlightInformationParams) (*FlightDataResponse, error) {
	return NewClient(authData).FlightInformation(ctx, params)
}

func GetAirportInformation(ctx context.Context, authData authdata.AuthData, airportIdentifier string) (*AirportData, error) {
	return NewClient(authData).AirportInformation(ctx, airportIdentifier)
}

func GetFlightIdentifiers(ctx context.Context, authData authdata.AuthData, params FlightInformationParams) (*FlightIdentifiers, string, error) {
	return NewClient(authData).FlightIdentifiers(ctx, params)
}
