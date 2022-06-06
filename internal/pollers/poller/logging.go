package poller

import (
	"go.uber.org/zap"
)

func newLogger() *zap.Logger {
	l, _ := zap.NewProduction(zap.AddCallerSkip(1))
	defer func() { _ = l.Sync() }()

	return l
}

func (p BasePoller) sugarLogger() *zap.SugaredLogger {
	return p.logger.Sugar()
}

func (p BasePoller) LogDebug(msg string, fields ...zap.Field) {
	p.logger.Debug(msg, fields...)
}

func (p BasePoller) LogInfo(msg string, fields ...zap.Field) {
	p.logger.Info(msg, fields...)
}

func (p BasePoller) LogWarning(msg string, fields ...zap.Field) {
	p.logger.Warn(msg, fields...)
}

func (p BasePoller) LogError(msg string, fields ...zap.Field) {
	p.logger.Error(msg, fields...)
}

func (p BasePoller) LogFatal(msg string, fields ...zap.Field) {
	p.logger.Fatal(msg, fields...)
}

func (p BasePoller) LogPanic(msg string, fields ...zap.Field) {
	p.logger.Panic(msg, fields...)
}

func (p BasePoller) LogDebugf(template string, args ...any) {
	p.sugarLogger().Debugf(template, args...)
}

func (p BasePoller) LogInfof(template string, args ...any) {
	p.sugarLogger().Infof(template, args...)
}

func (p BasePoller) LogWarningf(template string, args ...any) {
	p.sugarLogger().Warnf(template, args...)
}

func (p BasePoller) LogErrorf(template string, args ...any) {
	p.sugarLogger().Errorf(template, args...)
}

func (p BasePoller) LogFatalf(template string, args ...any) {
	p.sugarLogger().Fatalf(template, args...)
}

func (p BasePoller) LogPanicf(template string, args ...any) {
	p.sugarLogger().Panicf(template, args...)
}
