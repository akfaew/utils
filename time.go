package utils

import "time"

type TimeReal struct{}
type TimeTest struct{}

type MockTime interface {
	Now() time.Time
	Since(t time.Time) time.Duration
}

var Time MockTime = TimeReal{}

func (_ TimeReal) Now() time.Time {
	return time.Now()
}

func (_ TimeReal) Since(t time.Time) time.Duration {
	return time.Since(t)
}

func (_ TimeTest) Now() time.Time {
	now, _ := time.Parse("2006-01-02 15:04:05", "2018-10-29 14:34:45")
	return now
}

func (_ TimeTest) Since(t time.Time) time.Duration {
	now, _ := time.Parse("2006-01-02 15:04:05", "2018-10-29 14:34:45")
	return now.Sub(t)
}
