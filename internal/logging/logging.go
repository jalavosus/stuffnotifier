package logging

import (
	"go.uber.org/zap"
)

func NewLogger() *zap.Logger {
	return newLogger(false)
}

func NewWrappedLogger() *zap.Logger {
	return newLogger(true)
}

func newLogger(withCallSkip bool) *zap.Logger {
	var opts []zap.Option
	if withCallSkip {
		opts = append(opts, zap.AddCallerSkip(1))
	}

	l, _ := zap.NewProduction(opts...)
	defer func() { _ = l.Sync() }()

	return l
}
