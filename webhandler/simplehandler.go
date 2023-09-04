package webhandler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/akfaew/syften/aelog"
)

type SimpleHandler func(http.ResponseWriter, *http.Request) *SimpleError

type SimpleError struct {
	Code  int   // HTTP response code
	Error error // for the logs
}

func NewSimpleError(code int, err error) *SimpleError {
	return &SimpleError{
		Code:  code,
		Error: err,
	}
}

func (fn SimpleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		ctx := r.Context()
		lg := aelog.New(ctx)

		m := meditation()
		lg.Errorf("Handler error %d. err=\"%v\", meditation=%s", e.Code, e.Error, m)
		http.Error(w, fmt.Sprintf("%d (%s). meditation: %s", e.Code, http.StatusText(e.Code), m), e.Code)
	}
}

func SimpleHandle(r *mux.Router, method string, path string, handler func(w http.ResponseWriter, r *http.Request) *SimpleError) {
	r.Methods(method).Path(path).Handler(SimpleHandler(handler))
}
