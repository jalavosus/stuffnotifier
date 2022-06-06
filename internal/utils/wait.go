package utils

import (
	"context"
	"time"
)

func Wait(ctx context.Context, d time.Duration) {
	t := time.NewTicker(d)

	cleanup := func() {
		t.Stop()
	}

	for {
		select {
		case <-ctx.Done():
			cleanup()
			return
		case <-t.C:
			cleanup()
			return
		}
	}
}
