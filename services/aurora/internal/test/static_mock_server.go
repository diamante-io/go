package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// NewStaticMockServer creates a new mock server that always responds with
// `response`
func NewStaticMockServer(response string) *StaticMockServer {
	result := &StaticMockServer{}
	result.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result.LastRequest = r
		fmt.Fprintln(w, response)
	}))

	return result
}
