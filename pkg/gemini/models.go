package gemini

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type OrderBookStatus string

const (
	Open       OrderBookStatus = "open"
	Closed     OrderBookStatus = "closed"
	CancelOnly OrderBookStatus = "cancel_only"
	PostOnly   OrderBookStatus = "post_only"
	LimitOnly  OrderBookStatus = "limit_only"
)

type SymbolsResponse []string

type SymbolDetailsResponse struct {
	Symbol         string          `json:"symbol"`
	BaseCurrency   string          `json:"base_currency"`
	QuoteCurrency  string          `json:"quote_currency"`
	TickSize       decimal.Decimal `json:"tick_size"`
	QuoteIncrement decimal.Decimal `json:"quote_increment"`
	MinOrderSize   decimal.Decimal `json:"min_order_size"`
	Status         OrderBookStatus `json:"status"`
	WrapEnabled    bool            `json:"wrap_enabled"`
}

type TickerResponse struct {
	Bid    decimal.Decimal `json:"bid"`
	Ask    decimal.Decimal `json:"ask"`
	Last   decimal.Decimal `json:"last"`
	Volume TickerVolume    `json:"volume"`
}

type TickerVolume struct {
	Timestamp time.Time `json:"timestamp"`
	amounts   map[string]decimal.Decimal
}

func (tv *TickerVolume) UnmarshalJSON(data []byte) error {
	raw := make(map[string]any)
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	ts := raw["timestamp"].(float64)
	tv.Timestamp = time.UnixMilli(int64(ts))

	tv.amounts = make(map[string]decimal.Decimal)
	for k, v := range raw {
		if k == "timestamp" {
			continue
		}

		if _, ok := tv.amounts[k]; ok {
			continue
		}

		amtStr := v.(string)
		amt, err := decimal.NewFromString(amtStr)
		if err != nil {
			err = errors.WithMessagef(err, "error parsing string %[1]s into decimal.Decimal", amtStr)
			return err
		}

		tv.amounts[k] = amt
	}

	return nil
}

// Volume returns the trading volume denominated in the passed symbol.
func (tv TickerVolume) Volume(symbol string) decimal.Decimal {
	val, ok := tv.amounts[symbol]
	if ok {
		return val
	}

	return decimal.Zero
}

type TickerV2Response struct {
	Symbol  string            `json:"symbol"`
	Open    decimal.Decimal   `json:"open"`
	High    decimal.Decimal   `json:"high"`
	Low     decimal.Decimal   `json:"low"`
	Close   decimal.Decimal   `json:"close"`
	Changes []decimal.Decimal `json:"changes"`
	Bid     decimal.Decimal   `json:"bid"`
	Ask     decimal.Decimal   `json:"ask"`
}

type PriceFeedResponse []PriceFeedData

type PriceFeedData struct {
	Pair              string          `json:"pair"`
	Price             decimal.Decimal `json:"price"`
	PercentChange24Hr decimal.Decimal `json:"percentChange24h"`
}

type WebsocketResponse struct {
	Type           string    `json:"type"`
	EventId        int64     `json:"eventId"`
	Timestamp      time.Time `json:"timestamp"`
	TimestampMs    time.Time `json:"timestampms"`
	SocketSequence int64     `json:"socket_sequence"`
}

type SubscribeTradesResponse struct {
	WebsocketResponse
	Events []TradeEvent `json:"events"`
}

type TradeEvent struct {
	Type      string          `json:"type"`
	Tid       int64           `json:"tid"`
	Price     decimal.Decimal `json:"price"`
	Amount    decimal.Decimal `json:"amount"`
	MakerSide string          `json:"makerSide"`
}
