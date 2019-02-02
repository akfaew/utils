package utils

import (
	"fmt"
	"log"
	"os"
)

var Debug = false

func LogcDebugf(format string, a ...interface{}) {
	if Debug {
		log.Printf(format, a...)
	}
}

func LogcDebugfd(format string, a ...interface{}) {
	if Debug {
		file, line := logctx(1)

		log.Printf("%s:%d %s", file, line, fmt.Sprintf(format, a...))
	}
}

func LogcErrorf(format string, a ...interface{}) {
	log.Printf("ERROR: "+format, a...)
}

func LogcErrorfd(format string, a ...interface{}) {
	file, line := logctx(1)

	log.Printf("ERROR: %s:%d %s", file, line, fmt.Sprintf(format, a...))
}

func LogcFatalf(format string, a ...interface{}) {
	log.Printf("FATAL: "+format, a...)
	os.Exit(1)
}

func LogcFatalfd(format string, a ...interface{}) {
	file, line := logctx(1)

	log.Printf("FATAL: %s:%d %s", file, line, fmt.Sprintf(format, a...))
	os.Exit(1)
}
