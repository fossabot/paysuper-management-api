package common

import "time"

// GetTimeRangeFrom
func GetTimeRangeFrom(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// GetTimeRangeTo
func GetTimeRangeTo(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 99, t.Location())
}
