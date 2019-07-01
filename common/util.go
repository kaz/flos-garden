package common

import (
	"time"
)

var (
	tz *time.Location
)

func Time(unixnano int64) time.Time {
	if tz == nil {
		tz, _ = time.LoadLocation("Asia/Tokyo")
	}
	return time.Unix(0, unixnano).In(tz)
}
