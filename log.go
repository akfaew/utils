package utils

import (
	"fmt"
	"log"
	"os"
)

var Debug = false

func LogcDebugf(format string, a ...any) {
	if Debug {
		log.Printf(format, a...)
	}
}

func LogcDebugfd(format string, a ...any) {
	if Debug {
		file, line := logctx(1)

		log.Printf("%s:%d %s", file, line, fmt.Sprintf(format, a...))
	}
}

func LogcInfof(format string, a ...any) {
	log.Printf(format, a...)
}

func LogcInfofd(format string, a ...any) {
	file, line := logctx(1)

	log.Printf("%s:%d %s", file, line, fmt.Sprintf(format, a...))
}

func LogcErrorf(format string, a ...any) {
	log.Printf("ERROR: "+format, a...)
}

func LogcErrorfd(format string, a ...any) {
	file, line := logctx(1)

	log.Printf("ERROR: %s:%d %s", file, line, fmt.Sprintf(format, a...))
}

func LogcFatalf(format string, a ...any) {
	log.Printf("FATAL: "+format, a...)
	os.Exit(1)
}

func LogcFatalfd(format string, a ...any) {
	file, line := logctx(1)

	log.Printf("FATAL: %s:%d %s", file, line, fmt.Sprintf(format, a...))
	os.Exit(1)
}
