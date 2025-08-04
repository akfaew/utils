package webhandler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

var simple = ParseTemplate("base.html", "simple.html")
var failure = ParseTemplate("base.html", "simple.html")
var nolog = ParseTemplate("base.html", "simple.html")
var redirect = ParseTemplate("base.html", "simple.html")
var nobase = ParseTemplate("nobase.html")

func webContext(w http.ResponseWriter, r *http.Request, tmpl *WebTemplate) (any, *WebError) {
	switch tmpl {
	case failure:
		return nil, WebErrorf(http.StatusInternalServerError, fmt.Errorf("ooups"), "User error")
	case nolog:
		return nil, WebErrorf(http.StatusInternalServerError, nil, "")
	case redirect:
		http.Redirect(w, r, "/redirect_target", http.StatusFound)
	}

	return struct {
		Title string
	}{
		Title: "Template Title",
	}, nil
}

func errorContext(message, meditation string) any {
	return struct {
		Message string
	}{
		Message: message,
	}
}

func webRouter() *mux.Router {
	r := mux.NewRouter()
	r.StrictSlash(true)

	ErrorContextFunc = errorContext

	WebHandle(r, http.MethodGet, "/success", simple.Executor(webContext))
	WebHandle(r, http.MethodGet, "/failure", failure.Executor(webContext))
	WebHandle(r, http.MethodGet, "/nolog", nolog.Executor(webContext))
	WebHandle(r, http.MethodGet, "/nobase", nobase.Executor(webContext))
	WebHandle(r, http.MethodGet, "/redirect", redirect.Executor(webContext))

	return r
}

func TestWebHandler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/success", nil)
		require.NoError(t, err)

		HTTPGetRouter(t, webRouter(), req).Fixture()
	})

	t.Run("Failure", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/failure", nil)
		require.NoError(t, err)

		HTTPGetRouter(t, webRouter(), req).Fixture()
		HTTPGetRouter(t, webRouter(), req).Status(http.StatusInternalServerError)
	})

	t.Run("No Log", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/nolog", nil)
		require.NoError(t, err)

		HTTPGetRouter(t, webRouter(), req).Fixture()
	})

	t.Run("No Base", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/nobase", nil)
		require.NoError(t, err)

		HTTPGetRouter(t, webRouter(), req).Fixture()
	})

	t.Run("Redirect", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/redirect", nil)
		require.NoError(t, err)

		HTTPGetRouter(t, webRouter(), req).Fixture()
	})
}
