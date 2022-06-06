package utils_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stoicturtle/stuffnotifier/internal/utils"
)

func TestWait(t *testing.T) {
	t.Parallel()

	rootCtx := context.Background()

	tests := []struct {
		name string
		d    time.Duration
	}{
		{
			name: "1 second",
			d:    time.Second,
		},
		{
			name: "10 seconds",
			d:    10 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			now := time.Now()

			utils.Wait(rootCtx, tt.d)

			after := time.Now()

			diff := after.Sub(now)

			assert.LessOrEqual(t, tt.d.Seconds(), diff.Seconds())
		})
	}
}
