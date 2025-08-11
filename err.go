package utils

import (
	"errors"
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

func UserErrorfc(err error, format string, a ...any) error {
	file, line := logctx(1)

	return UserError{
		Err:         fmt.Errorf("%s:%d %w", file, line, err),
		UserMessage: fmt.Sprintf(format, a...),
	}
}

func Errorfc(format string, a ...any) error {
	file, line := logctx(1)

	return fmt.Errorf("%s:%d %w", file, line, fmt.Errorf(format, a...))
}

func Errorc(err error) error {
	if err == nil {
		return nil
	}
	file, line := logctx(1)

	return fmt.Errorf("%s:%d %w", file, line, err)
}

func RootCause(err error) error {
	for {
		unwrapped := errors.Unwrap(err)
		if unwrapped == nil {
			return err
		}
		err = unwrapped
	}
}
