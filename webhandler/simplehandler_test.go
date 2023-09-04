package webhandler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func simpleSuccess(w http.ResponseWriter, r *http.Request) *SimpleError {
	return nil
}

func simpleFailure(w http.ResponseWriter, r *http.Request) *SimpleError {
	return NewSimpleError(http.StatusInternalServerError, fmt.Errorf("ooups"))
}

func apiRouter() *mux.Router {
	r := mux.NewRouter()
	r.StrictSlash(true)

	SimpleHandle(r, http.MethodGet, "/success", simpleSuccess)
	SimpleHandle(r, http.MethodGet, "/failure", simpleFailure)

	return r
}

func TestSimpleHandler(t *testing.T) {
	DisableMeditation = true
	t.Run("Success", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/success", nil)
		require.NoError(t, err)

		HTTPGetRouter(t, apiRouter(), req).Fixture()
	})

	t.Run("Failure", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/failure", nil)
		require.NoError(t, err)

		HTTPGetRouter(t, apiRouter(), req).Fixture()
	})
}
