package gemini

import (
	"context"

	"github.com/jalavosus/stuffnotifier/pkg/authdata"
)

func Symbols(ctx context.Context, authData authdata.AuthData) (*SymbolsResponse, error) {
	return NewClient(authData).Symbols(ctx)
}

func SymbolDetails(ctx context.Context, authData authdata.AuthData, symbol string) (*SymbolDetailsResponse, error) {
	return NewClient(authData).SymbolDetails(ctx, symbol)
}

func Ticker(ctx context.Context, authData authdata.AuthData, symbol string) (*TickerResponse, error) {
	return NewClient(authData).Ticker(ctx, symbol)
}

func TickerV2(ctx context.Context, authData authdata.AuthData, symbol string) (*TickerV2Response, error) {
	return NewClient(authData).TickerV2(ctx, symbol)
}

func PriceFeed(ctx context.Context, authData authdata.AuthData) (*PriceFeedResponse, error) {
	return NewClient(authData).PriceFeed(ctx)
}

// func SubscribeTrades(ctx context.Context, symbol string) (<-chan *SubscribeTradesResponse, func(), error) {
// 	ch := make(chan *SubscribeTradesResponse)
//
// 	uri := marketDataEndpoint(symbol)
//
// 	stop := func() {
//
// 	}
// }
