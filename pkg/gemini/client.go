package gemini

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/stoicturtle/stuffnotifier/internal/authdata"
	"github.com/stoicturtle/stuffnotifier/internal/authdata/env"
	"github.com/stoicturtle/stuffnotifier/internal/utils"
	"github.com/stoicturtle/stuffnotifier/pkg/errs"
)

// Client is a superficial wrapper around
// API functions in this package,
// allowing users to call those functions without
// passing an authdata.AuthData object for every call.
type Client struct {
	authData   authdata.AuthData
	httpClient *http.Client
	baseApiUri string
}

// NewClient returns a Client instance.
func NewClient(authData authdata.AuthData) *Client {
	var baseApiUri = apiUri
	if env.GeminiUseSandbox() {
		baseApiUri = sandboxUri
	}

	return &Client{
		authData:   authData,
		httpClient: utils.HttpClientWithTimeout(defaultHttpTimeout),
		baseApiUri: baseApiUri,
	}
}

// SetHttpTimeout sets the timeout of the Client's http.Client instance.
// By default, this is set to 10 seconds.
func (c *Client) SetHttpTimeout(timeout time.Duration) *Client {
	c.httpClient.Timeout = timeout
	return c
}

func (c Client) Symbols(ctx context.Context) (*SymbolsResponse, error) {
	var response *SymbolsResponse

	resp, err := c.httpRequest(ctx, symbolsEndpoint)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(resp, &response); err != nil {
		return nil, errs.HttpUnmarshalResponseBodyError(err)
	}

	return response, nil
}

func (c Client) SymbolDetails(ctx context.Context, symbol string) (*SymbolDetailsResponse, error) {
	var response *SymbolDetailsResponse

	resp, err := c.httpRequest(ctx, symbolDetailsEndpoint(symbol))
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(resp, &response); err != nil {
		return nil, errs.HttpUnmarshalResponseBodyError(err)
	}

	return response, nil
}

func (c Client) Ticker(ctx context.Context, symbol string) (*TickerResponse, error) {
	var response *TickerResponse

	resp, err := c.httpRequest(ctx, tickerEndpoint(symbol))
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(resp, &response); err != nil {
		return nil, errs.HttpUnmarshalResponseBodyError(err)
	}

	return response, nil
}

func (c Client) TickerV2(ctx context.Context, symbol string) (*TickerV2Response, error) {
	var response *TickerV2Response

	resp, err := c.httpRequest(ctx, tickerV2Endpoint(symbol))
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(resp, &response); err != nil {
		return nil, errs.HttpUnmarshalResponseBodyError(err)
	}

	return response, nil
}

func (c Client) PriceFeed(ctx context.Context) (*PriceFeedResponse, error) {
	var response *PriceFeedResponse

	resp, err := c.httpRequest(ctx, priceFeedEndpoint)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(resp, &response); err != nil {
		return nil, errs.HttpUnmarshalResponseBodyError(err)
	}

	return response, nil
}

func (c Client) httpRequest(ctx context.Context, endpoint string) ([]byte, error) {
	uri := utils.BuildRequestEndpoint(c.baseApiUri, endpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, errs.HttpBuildRequestError(err)
	}

	sig, payload := c.buildNoncePayload(endpoint)
	req.Header.Set(apiKeyHeader, c.authData.Key())
	req.Header.Set(payloadHeader, payload)
	req.Header.Set(signatureHeader, sig)

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

func (c Client) buildNoncePayload(endpoint string) (sig, payload string) {
	return BuildNonceWithPayload(c.authData, endpoint)
}
