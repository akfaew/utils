package utils

import "log"

var Debug = false

func LogcDebugf(format string, a ...interface{}) {
	if Debug {
		log.Printf(format, a...)
	}
}
