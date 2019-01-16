package utils

import "time"

type TimeReal struct{}
type TimeTest struct{}

var Time interface {
	Since(t time.Time) time.Duration
} = TimeReal{}

func (_ TimeReal) Since(t time.Time) time.Duration {
	return time.Since(t)
}

func (_ TimeTest) Since(t time.Time) time.Duration {
	now, _ := time.Parse("2006-01-02 15:04:05", "2018-10-29 14:34:45")
	return now.Sub(t)
}
