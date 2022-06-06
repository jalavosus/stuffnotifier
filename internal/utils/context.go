package utils

import (
	"context"
	"time"
)

func TimeoutFromContext(ctx context.Context) (time.Duration, bool) {
	var timeout time.Duration
	
	deadline, ok := ctx.Deadline()
	if !ok {
		return timeout, ok
	}

	return time.Until(deadline), ok
}