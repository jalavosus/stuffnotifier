package flightaware

import (
	"context"

	"github.com/stoicturtle/stuffnotifier/internal/logging"
	"github.com/stoicturtle/stuffnotifier/pkg/authdata"
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

func FlightInformation(ctx context.Context, authData authdata.AuthData, flightIdentifier string, flightIdentifierType IdentifierType) (*FlightDataResponse, error) {
	return NewClient(authData).FlightInformation(ctx, flightIdentifier, flightIdentifierType)
}

func AirportInformation(ctx context.Context, authData authdata.AuthData, airportIdentifier string) (*AirportData, error) {
	return NewClient(authData).AirportInformation(ctx, airportIdentifier)
}
