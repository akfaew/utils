package errs

import (
	"fmt"
	"runtime"
)

// Error wraps an error with with the current file name and line number. Don't
// create it manually.
type Error struct {
	File    string
	Line    int
	Message string // Can be empty, in which case it's skipped
	Err     error  // The underlying (wrapped) error
}

func logctx(skip int) (file string, line int) {
	_, file, line, _ = runtime.Caller(skip + 1)

	return
}

func (e *Error) Error() string {
	if len(e.Message) > 0 {
		return fmt.Sprintf("%s:%d %s: %s", e.File, e.Line, e.Message, e.Err)
	} else {
		return fmt.Sprintf("%s:%d: %s", e.File, e.Line, e.Err)
	}
}

func (e *Error) Unwrap() error {
	return e.Err
}

// Errorf creates a new error and wraps it with the current file name and line number.
func Errorf(format string, a ...interface{}) error {
	file, line := logctx(1)
	return &Error{
		File: file,
		Line: line,
		Err:  fmt.Errorf(format, a...),
	}
}

// Wrapf wraps the error with the current file name, line number, and an additional message.
func Wrapf(err error, format string, a ...interface{}) error {
	if err == nil {
		return nil
	}

	file, line := logctx(1)
	return &Error{
		File:    file,
		Line:    line,
		Message: fmt.Sprintf(format, a...),
		Err:     err,
	}
}

// Wrap wraps the error with the current file name and line number.
func Wrap(err error) error {
	if err == nil {
		return nil
	}

	file, line := logctx(1)
	return &Error{
		File: file,
		Line: line,
		Err:  err,
	}
}
