package gemini

import (
	"time"
)

const (
	defaultHttpTimeout = 10 * time.Second
)

const (
	v1Endpoint string = "/v1"
	v2Endpoint string = "/v2"
)

const (
	apiUri     string = "api.gemini.com"
	sandboxUri string = "api.sandbox.gemini.com"
)

const (
	symbolsEndpoint  = v1Endpoint + "/symbols"
	symbolDetailsUri = symbolsEndpoint + "/details"
)

const (
	tickerUri   = v1Endpoint + "/pubticker"
	tickerV2Uri = v2Endpoint + "/ticker"
)

const (
	priceFeedEndpoint = v1Endpoint + "/pricefeed"
)

const (
	marketDataUri = v1Endpoint + "/marketdata"
)

const (
	headerPrefix    string = "X-GEMINI-"
	apiKeyHeader           = headerPrefix + "APIKEY"
	payloadHeader          = headerPrefix + "PAYLOAD"
	signatureHeader        = headerPrefix + "SIGNATURE"
)

func endpointWithSymbol(endpoint, symbol string) string {
	return endpoint + "/" + symbol
}

func symbolDetailsEndpoint(symbol string) string {
	return endpointWithSymbol(symbolDetailsUri, symbol)
}

func tickerEndpoint(symbol string) string {
	return endpointWithSymbol(tickerUri, symbol)
}

func tickerV2Endpoint(symbol string) string {
	return endpointWithSymbol(tickerV2Uri, symbol)
}

func marketDataEndpoint(symbol string) string {
	return endpointWithSymbol(marketDataUri, symbol)
}
