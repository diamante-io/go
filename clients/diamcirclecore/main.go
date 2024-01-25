// Package diamcirclecore is a client library for communicating with an
// instance of diamcircle-core using through the server's HTTP port.
package diamcirclecore

import "net/http"

// SetCursorDone is the success message returned by diamcircle-core when a cursor
// update succeeds.
const SetCursorDone = "Done"

// HTTP represents the http client that a diamcirclecore client uses to make http
// requests.
type HTTP interface {
	Do(req *http.Request) (*http.Response, error)
}

// confirm interface conformity
var _ HTTP = http.DefaultClient
