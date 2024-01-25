// Package httptest enhances the stdlib net/http/httptest package by integrating
// it with gopkg.in/gavv/httpexpect.v1, reducing the boilerplate needed for http
// tests.  In addition to the functions that make it easier to stand up a test
// server, this package also provides a set of tools to make it easier to mock
// out http client responses.
//
// Test Servers vs. Client mocking
//
// When building a testing fixture for HTTP, you can approach the problem in two
// ways:  Use a mocked server and make real http requests to the server, or use
// a mocked client and have _no_ server.  While this package provides facilities
// for both, we recommend that you follow the conventions for decided which to
// use in your tests.
//
// The test server system should be used when the object under test is our own
// server code; usually that means our own http.Handler implementations.  The
// mocked client system should be used when the object under test is making http
// requests.
package httptest

import (
	"net/http"
	stdtest "net/http/httptest"
	"testing"

	"github.com/jarcoal/httpmock"
	"gopkg.in/gavv/httpexpect.v1"
)

// Client represents an easier way to mock http client behavior.  It assumes
// that your packages use interfaces to store an http client instead of a
// concrete *http.Client.
type Client struct {
	*http.Client
	*httpmock.MockTransport
}

// ClientExpectation represents a in-process-of-being-built http client mocking
// operation.  The `On` method of `Client` returns an instance of this struct
// upon which you can call further methods to customize the response.
type ClientExpectation struct {
	Method string
	URL    string
	Client *Client
}

type Server struct {
	*httpexpect.Expect
	*stdtest.Server
}

// NewClient returns a new mocked http client.  A value being tested can be
// configured to use this http client allowing the tester to control the server
// responses without needing to run an actual server.
func NewClient() *Client {
	var result Client
	result.MockTransport = httpmock.NewMockTransport()
	result.Client = &http.Client{Transport: result.MockTransport}
	return &result
}

// NewServer returns a new test server instance using the provided handler.
func NewServer(t *testing.T, handler http.Handler) *Server {
	server := stdtest.NewServer(handler)
	return &Server{
		Server: server,
		Expect: httpexpect.New(t, server.URL),
	}
}
