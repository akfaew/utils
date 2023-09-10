package webhandler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"

	"github.com/akfaew/utils/fixture"
)

type WebTestResponse struct {
	t                *testing.T
	ResponseRecorder *httptest.ResponseRecorder
}

// Use like this:
//
//	func HTTPGet(t *testing.T, req *http.Request) *WebTestResponse {
//		return webhandler.HTTPGetRouter(t, Router(), req)
//	}
func HTTPGetRouter(t *testing.T, router *mux.Router, req *http.Request) *WebTestResponse {
	t.Helper()

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	return &WebTestResponse{
		t:                t,
		ResponseRecorder: res,
	}
}

func (r *WebTestResponse) Body() string {
	return r.ResponseRecorder.Body.String()
}

func (r *WebTestResponse) Location() string {
	r.t.Helper()

	require.NotEmpty(r.t, r.ResponseRecorder.Result().Header["Location"])

	return r.ResponseRecorder.Result().Header["Location"][0]
}

func (r *WebTestResponse) Fixture() {
	r.t.Helper()

	r.FixtureExtra("")
}

func (r *WebTestResponse) FixtureExtra(extra string) {
	r.t.Helper()

	code := fmt.Sprintf("Code: %d\n\n", r.ResponseRecorder.Code)
	fixture.FixtureExtra(r.t, extra, code+r.Body())
}

func (r *WebTestResponse) Status(want int) {
	r.t.Helper()

	if r.ResponseRecorder.Code != want {
		r.t.Fatalf("Status Code == %d (expected %d)", r.ResponseRecorder.Code, want)
	}
}
