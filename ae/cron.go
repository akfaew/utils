package ae

import (
	"errors"
	"net/http"
)

var (
	ErrNoCronHeader     = errors.New("cron request does not have the X-Appengine-Cron header")
	ErrNoTasknameHeader = errors.New("task request does not have the X-Appengine-Taskname header")
)

func ValidateCron(r *http.Request) error {
	if r.Header.Get("X-Appengine-Cron") == "" {
		return ErrNoCronHeader
	}

	return nil
}

func ValidateTask(r *http.Request) error {
	if r.Header.Get("X-Appengine-Taskname") == "" {
		return ErrNoTasknameHeader
	}

	return nil
}
