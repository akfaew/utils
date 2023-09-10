package utils

import (
	"fmt"
)

type UserError struct {
	Err         error
	UserMessage string
}

func (e UserError) Error() string {
	return e.Err.Error()
}

func (e UserError) Unwrap() error {
	return e.Err
}

func (e UserError) Message() string {
	return e.UserMessage
}

func UserErrorfc(err error, format string, a ...interface{}) error {
	file, line := logctx(1)

	return UserError{
		Err:         fmt.Errorf("%s:%d %w", trim(file), line, err),
		UserMessage: fmt.Sprintf(format, a...),
	}
}

func Errorfc(format string, a ...interface{}) error {
	file, line := logctx(1)

	return fmt.Errorf("%s:%d %w", trim(file), line, fmt.Errorf(format, a...))
}

func Errorc(err error) error {
	if err == nil {
		return nil
	}
	file, line := logctx(1)

	return fmt.Errorf("%s:%d %w", trim(file), line, err)
}
