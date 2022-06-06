package flightaware

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/jalavosus/stuffnotifier/internal/utils"
	"github.com/jalavosus/stuffnotifier/pkg/authdata"
	"github.com/jalavosus/stuffnotifier/pkg/errs"
)

// Client is a superficial wrapper around
// API functions in this package,
// allowing users to call those functions without
// passing an authdata.AuthData object for every call.
type Client struct {
	authData   authdata.AuthData
	httpClient *http.Client
}

func NewClient(authData authdata.AuthData) *Client {
	return &Client{
		authData:   authData,
		httpClient: utils.HttpClientWithTimeout(defaultHttpTimeout),
	}
}

// SetHttpTimeout sets the timeout of the Client's http.Client instance.
// By default, this is set to 10 seconds.
func (c *Client) SetHttpTimeout(timeout time.Duration) *Client {
	c.httpClient.Timeout = timeout
	return c
}

func (c Client) FlightInformation(ctx context.Context, flightIdentifier string, flightIdentifierType IdentifierType) (*FlightDataResponse, error) {
	var response *FlightDataResponse

	params := url.Values{}
	params.Set(identTypeParam, flightIdentifierType.String())

	endpoint := fmt.Sprintf(flightsUri, flightIdentifier) + "?" + params.Encode()

	resp, err := c.httpRequest(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(resp, &response); err != nil {
		return nil, errs.HttpUnmarshalResponseBodyError(err)
	}

	return response, nil
}

func (c Client) AirportInformation(ctx context.Context, airportIdentifier string) (*AirportData, error) {
	var response *AirportData

	endpoint := fmt.Sprintf(airportsUri, airportIdentifier)

	resp, err := c.httpRequest(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(resp, &response); err != nil {
		return nil, errs.HttpUnmarshalResponseBodyError(err)
	}

	return response, nil
}

func (c Client) httpRequest(ctx context.Context, endpoint string) ([]byte, error) {
	uri := utils.BuildRequestEndpoint(baseApiEndpoint, endpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, errs.HttpBuildRequestError(err)
	}

	req.Header.Set(apiKeyHeader, c.authData.Key())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errs.HttpResponseError(err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errs.HttpReadBodyError(err)
	}

	return respBody, nil
}
