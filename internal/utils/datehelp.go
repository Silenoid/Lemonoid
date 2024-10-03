package utils

import (
	"time"
)

func ToReadableDate(nowTime time.Time) string {
	return nowTime.Format("2006-01-02 15:04:05")
}

func ToReadableSince(nowTime time.Time, pastTime time.Time) string {
	return nowTime.Sub(pastTime).Round(time.Minute).String()
}

func ToReadableHowLongTo(nowTime time.Time, pastTime time.Time, timePeriod time.Duration) string {
	return pastTime.Add(timePeriod).Sub(nowTime).Round(time.Minute).String()
}
