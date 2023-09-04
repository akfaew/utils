package ae

import (
	"errors"
	"net/http"
)

var (
	ErrNoCronHeader     = errors.New("Cron request does not have the X-Appengine-Cron header")
	ErrNoTasknameHeader = errors.New("Task request does not have the X-Appengine-Taskname header")
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
