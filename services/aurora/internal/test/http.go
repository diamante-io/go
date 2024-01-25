package test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/go-chi/chi"
)

type RequestHelper interface {
	Get(string, ...func(*http.Request)) *httptest.ResponseRecorder
	Post(string, url.Values, ...func(*http.Request)) *httptest.ResponseRecorder
}

type requestHelper struct {
	router *chi.Mux
}

func RequestHelperRaw(r *http.Request) {
	r.Header.Set("Accept", "application/octet-stream")
}

func RequestHelperStreaming(r *http.Request) {
	r.Header.Set("Accept", "text/event-stream")
}

func NewRequestHelper(router *chi.Mux) RequestHelper {
	return &requestHelper{router}
}

func (rh *requestHelper) Get(
	path string,
	mods ...func(*http.Request),
) *httptest.ResponseRecorder {

	req, _ := http.NewRequest("GET", path, nil)
	return rh.Execute(req, mods)
}

func (rh *requestHelper) Post(
	path string,
	form url.Values,
	mods ...func(*http.Request),
) *httptest.ResponseRecorder {

	body := strings.NewReader(form.Encode())
	req, _ := http.NewRequest("POST", path, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return rh.Execute(req, mods)
}

func (rh *requestHelper) Execute(
	req *http.Request,
	requestModFns []func(*http.Request),
) *httptest.ResponseRecorder {

	req.RemoteAddr = "127.0.0.1"
	req.Host = "localhost"
	for _, fn := range requestModFns {
		fn(req)
	}

	w := httptest.NewRecorder()

	rh.router.ServeHTTP(w, req)
	return w

}
