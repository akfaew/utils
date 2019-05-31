package utils

import (
	"fmt"
)

// Errorfc creates a new error and wraps it with context information using Errorc().
func Errorfc(format string, a ...interface{}) error {
	file, line := logctx(1)
	return fmt.Errorf("%s:%d %s", trim(file), line, fmt.Errorf(format, a...))
}

// Errorc wraps the error, if not nil, with context information, such as the file name
// and line.
func Errorc(err error) error {
	if err != nil {
		file, line := logctx(1)
		return fmt.Errorf("%s:%d %s", trim(file), line, err)
	}
	return nil
}
