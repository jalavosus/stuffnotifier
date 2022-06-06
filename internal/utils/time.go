package utils

import (
	"time"

	"go.uber.org/zap"

	"github.com/jalavosus/stuffnotifier/internal/logging"
)

var logger = logging.NewLogger()

const (
	timeFormatUTC        string = "2006-01-02 15:04:05 (UTC)"
	timeFormatNoOffset   string = "2006-01-02 15:04:05 MST"
	timeFormatWithOffset string = "2006-01-02 15:04:05 MST (-0700)"
)

func FormatTime(t time.Time, withUtcOffset bool) string {
	var timeFormat = timeFormatNoOffset
	if withUtcOffset {
		timeFormat = timeFormatWithOffset
	}

	return t.Format(timeFormat)
}

func FormatTimeWithZone(t time.Time, tz string, withUtcOffset bool) string {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		logger.Fatal("error parsing timezone name", zap.String("timezone", tz), zap.Error(err))
	}

	t = t.In(loc)

	return FormatTime(t, withUtcOffset)
}

func FormatTimeUTC(t time.Time) string {
	t = t.UTC()
	return t.Format(timeFormatUTC)
}

func ToLocalTime(t time.Time) time.Time {
	t = t.In(time.Local)
	return t
}
