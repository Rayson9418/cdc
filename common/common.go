package common

import (
	"time"
)

const (
	KTableNameFmt = "%s.%s" // database.table
)

func IsToday(timestamp int64) bool {
	_time := time.Unix(timestamp, 0)
	nowTime := time.Now()

	return _time.Year() == nowTime.Year() &&
		_time.Month() == nowTime.Month() &&
		_time.Day() == nowTime.Day()
}
