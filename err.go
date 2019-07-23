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
// and line. If err is found in exclude, then it is returned unchanged.
func Errorc(err error, exclude ...error) error {
	if err == nil {
		return nil
	}

	for _, e := range exclude {
		if err == e {
			return err
		}
	}
	file, line := logctx(1)
	return fmt.Errorf("%s:%d %s", trim(file), line, err)
}
