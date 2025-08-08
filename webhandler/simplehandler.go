package webhandler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akfaew/utils/ae"
	"github.com/akfaew/utils/xctc"
	"github.com/gorilla/mux"
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
		lg := ae.NewLog(ctx)

		m := meditation()
		lg.Errorf("Handler error %d. err=\"%v\", meditation=%s", e.Code, e.Error, m)
		http.Error(w, fmt.Sprintf("%d (%s). meditation: %s", e.Code, http.StatusText(e.Code), m), e.Code)
	}
}

func SimpleHandle(r *mux.Router, method string, path string, handler func(w http.ResponseWriter, r *http.Request) *SimpleError) {
	wrappedHandler := func(w http.ResponseWriter, r *http.Request) *SimpleError {
		ctx := context.WithValue(r.Context(), xctc.XctcKey, r.Header.Get("Traceparent"))
		return handler(w, r.WithContext(ctx))
	}
	r.Methods(method).Path(path).Handler(SimpleHandler(wrappedHandler))
}
