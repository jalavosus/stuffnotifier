package gemini_test

import (
	"context"
	"testing"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/assert"

	"github.com/stoicturtle/stuffnotifier/internal/utils"
	"github.com/stoicturtle/stuffnotifier/pkg/authdata"
	"github.com/stoicturtle/stuffnotifier/pkg/gemini"
)

type TestCase struct {
	name    string
	symbol  string
	wantErr bool
	valid   bool
}

var (
	ethBtcTestCase = TestCase{
		name:    "ETHBTC (valid)",
		symbol:  "ETHBTC",
		wantErr: false,
		valid:   true,
	}
	ustUsdTestCase = TestCase{
		name:    "USTUSD (valid)",
		symbol:  "USTUSD",
		wantErr: false,
		valid:   true,
	}
	ustFraxTestCase = TestCase{
		name:    "USTFRAX (invalid)",
		symbol:  "USTFRAX",
		wantErr: false,
		valid:   false,
	}
)

func TestSymbols(t *testing.T) {
	t.Skip()

	rootCtx := context.Background()
	ctx, cancel := context.WithTimeout(rootCtx, 10*time.Second)
	defer cancel()

	auth, _ := authdata.GeminiAPIAuth()

	got, err := gemini.Symbols(ctx, auth)

	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.NotEmpty(t, got)
}

func TestSymbolDetails(t *testing.T) {
	type detailsWant struct {
		symbol string
		base   string
		quote  string
		status gemini.OrderBookStatus
	}

	tests := []struct {
		TestCase
		want detailsWant
	}{
		{
			ethBtcTestCase,
			detailsWant{
				symbol: "ETHBTC",
				base:   "ETH",
				quote:  "BTC",
				status: gemini.Open,
			},
		},
		{
			TestCase{
				name:    "ETHBTC (valid-badstatus)",
				symbol:  "ETHBTC",
				wantErr: false,
				valid:   true,
			},
			detailsWant{
				symbol: "ETHBTC",
				base:   "ETH",
				quote:  "BTC",
				status: gemini.PostOnly,
			},
		},
		// {
		// 	ustUsdTestCase,
		// 	detailsWant{
		// 		symbol: "USTUSD",
		// 		base:   "UST",
		// 		quote:  "USD",
		// 		status: gemini.LimitOnly,
		// 	},
		// },
		{
			ustFraxTestCase,
			detailsWant{
				symbol: "USTFRAX",
				base:   "UST",
				quote:  "FRAX",
			},
		},
	}

	rootCtx := context.Background()
	auth, _ := authdata.GeminiAPIAuth()

	var limitedStatuses = []gemini.OrderBookStatus{
		gemini.LimitOnly,
		gemini.CancelOnly,
		gemini.PostOnly,
	}

	checkLimitedStatus := func(want, got gemini.OrderBookStatus) func() (success bool) {
		return func() (success bool) {
			gotExpected := got == want
			if utils.SliceIncludes(limitedStatuses, want) {
				return gotExpected || got == gemini.Open
			}

			return gotExpected
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(rootCtx, 10*time.Second)
			defer cancel()

			got, err := gemini.SymbolDetails(ctx, auth, tt.symbol)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			if tt.valid {
				assert.Equal(t, tt.want.symbol, got.Symbol)
				assert.Equal(t, tt.want.base, got.BaseCurrency)
				assert.Equal(t, tt.want.quote, got.QuoteCurrency)
				assert.Condition(t, checkLimitedStatus(tt.want.status, got.Status))
			} else {
				assert.Equal(t, "", got.Symbol)
				assert.Equal(t, "", got.BaseCurrency)
				assert.Equal(t, "", got.QuoteCurrency)
			}
		})
	}
}

func TestTicker(t *testing.T) {
	type args struct {
		symbol  string
		symbolA string
		symbolB string
	}

	tests := []struct {
		TestCase
		args args
	}{
		{
			ethBtcTestCase,
			args{"ETHBTC", "ETH", "BTC"},
		},
		// {
		// 	ustUsdTestCase,
		// 	args{"USTUSD", "UST", "USD"},
		// },
		{
			ustFraxTestCase,
			args{"USTFRAX", "UST", "FRAX"},
		},
	}

	rootCtx := context.Background()

	auth, _ := authdata.GeminiAPIAuth()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(rootCtx, 10*time.Second)
			defer cancel()

			got, err := gemini.Ticker(ctx, auth, tt.args.symbol)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			if tt.valid {
				assert.False(t, got.Last.IsZero(), "expected got.Last to not be zero")
				assert.Falsef(t, got.Volume.Volume(tt.args.symbolA).IsZero(), "expected got.Volume(%[1]s) to not be zero", tt.args.symbolA)
				assert.Falsef(t, got.Volume.Volume(tt.args.symbolB).IsZero(), "expected got.Volume(%[1]s) to not be zero", tt.args.symbolB)
				assert.Truef(t, got.Volume.Volume(tt.args.symbol).IsZero(), "expected got.Volume(%[1]s) to be zero", tt.args.symbolA)
			} else {
				assert.True(t, got.Last.IsZero(), "expected got.Last to be zero")
				assert.Truef(t, got.Volume.Volume(tt.args.symbolA).IsZero(), "expected got.Volume(%[1]s) to be zero", tt.args.symbolA)
				assert.Truef(t, got.Volume.Volume(tt.args.symbolB).IsZero(), "expected got.Volume(%[1]s) to be zero", tt.args.symbolB)
			}
		})
	}
}

func TestTickerV2(t *testing.T) {
	tests := []struct {
		name    string
		symbol  string
		wantErr bool
		valid   bool
	}{
		ethBtcTestCase,
		ustUsdTestCase,
		ustFraxTestCase,
	}

	rootCtx := context.Background()
	auth, _ := authdata.GeminiAPIAuth()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(rootCtx, 10*time.Second)
			defer cancel()

			got, err := gemini.TickerV2(ctx, auth, tt.symbol)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			if tt.valid {
				assert.False(t, got.Open.IsZero(), "expected got.Open to not be zero")
				assert.False(t, got.Close.IsZero(), "expected got.Close to not be zero")
				assert.NotEmpty(t, got.Changes)
			} else {
				assert.True(t, got.Open.IsZero(), "expected got.Open to be zero")
				assert.True(t, got.Close.IsZero(), "expected got.Close to be zero")
				assert.Empty(t, got.Changes)
			}
		})
	}
}

func TestPriceFeed(t *testing.T) {
	// t.Skip()

	rootCtx := context.Background()
	ctx, cancel := context.WithTimeout(rootCtx, 10*time.Second)
	defer cancel()

	auth, _ := authdata.GeminiAPIAuth()

	got, err := gemini.PriceFeed(ctx, auth)

	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.NotEmpty(t, got)
}
