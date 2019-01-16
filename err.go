package utils

import (
	"fmt"
)

func Err(format string, a ...interface{}) error {
	file, line := logctx(1)

	return fmt.Errorf("%s:%d %s", file, line, fmt.Sprintf(format, a...))
}

func ErrRich(err error) error {
	if err != nil {
		file, line := logctx(1)
		return fmt.Errorf("%s:%d %s", trim(file), line, err)
	}
	return nil
}
