package utils_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stoicturtle/stuffnotifier/internal/utils"
)

func TestTimeoutFromContext(t *testing.T) {
	rootCtx := context.Background()

	tests := []struct {
		name    string
		timeout time.Duration
		wantOk  bool
		check   assert.ValueAssertionFunc
	}{
		{
			name:    "ok=true",
			timeout: 3 * time.Second,
			wantOk:  true,
			check:   assert.NotZero,
		},
		{
			name:   "ok=false",
			wantOk: false,
			check:  assert.Zero,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				ctx    context.Context
				cancel context.CancelFunc
			)

			if tt.timeout != 0 {
				ctx, cancel = context.WithTimeout(rootCtx, tt.timeout)
			} else {
				ctx, cancel = context.WithCancel(rootCtx)
			}

			defer cancel()

			got, ok := utils.TimeoutFromContext(ctx)

			assert.Equal(t, tt.wantOk, ok)
			tt.check(t, got)
		})
	}
}
