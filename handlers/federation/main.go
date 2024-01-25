// Package federation provides a pluggable handler that satisfies the diamcircle
// federation protocol.  Add an instance of `Handler` onto your router to allow
// a server to satisfy the protocol.
//
// The central type in this package is the "Driver" interfaces.  Implementing
// these interfaces allows a developer to plug in their own back end, whether it
// be a RDBMS, a KV store, or even just an in memory data structure.
//
// A pre-baked implementation of `Driver` and `ReverseDriver` that provides
// simple access to SQL systems is included. See `SQLDriver` for more details.
package federation

import (
	"context"
	"database/sql"
	"net/url"
	"sync"

	"go/support/db"
)

// Driver represents a data source against which federation queries can be
// executed.
type Driver interface {
	// LookupRecord is called when a handler receives a so-called "name"
	// federation request to lookup a `Record` using the provided diamcircle address.
	// An implementation should return a nil `*Record` value if the lookup
	// successfully executed but no result was found.
	LookupRecord(ctx context.Context, name, domain string) (*Record, error)
}

// ErrorResponse represents the JSON response sent to a client when the request
// triggered an error. FederationDriver methods can return this as an error and
// it will be passed to end user.
type ErrorResponse struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (response ErrorResponse) Error() string {
	return response.Message
}

// Handler represents an http handler that can service http requests that
// conform to the diamcircle federation protocol.  This handler should be added to
// your chosen mux at the path `/federation` (and for good measure
// `/federation/` if your middleware doesn't normalize trailing slashes).
type Handler struct {
	// Driver is the backend against which queries will be evaluated.
	Driver Driver
}

// Record represents the result from the database when performing a
// federation request.
type Record struct {
	AccountID string `db:"id"`
	MemoType  string `db:"memo_type"`
	Memo      string `db:"memo"`
}

// ReverseDriver represents a data source against which federation queries can
// be executed.
type ReverseDriver interface {
	// LookupReverseRecord is called when a handler receives a "reverse"
	// federation request to lookup a `ReverseRecord` using the provided strkey
	// encoded accountID. An implementation should return a nil `*ReverseRecord`
	// value if the lookup successfully executed but no result was found.
	LookupReverseRecord(ctx context.Context, accountID string) (*ReverseRecord, error)
}

// ForwardDriver represents a data source against which forward queries can
// be executed.
type ForwardDriver interface {
	// Forward is called when a handler receives a so-called "forward"
	// federation request to lookup a `Record` using the provided data (ex. bank
	// account number).
	// An implementation should return a nil `*Record` value if the lookup
	// successfully executed but no result was found.
	LookupForwardingRecord(query url.Values) (*Record, error)
}

// ReverseRecord represents the result from performing a "Reverse federation"
// lookup, in which an Account ID is used to lookup an associated address.
type ReverseRecord struct {
	Name   string `db:"name"`
	Domain string `db:"domain"`
}

// ReverseSQLDriver provides a `ReverseDriver` implementation based upon a SQL
// Server.  See `SQLDriver`, the forward only version, for more details.
type ReverseSQLDriver struct {
	SQLDriver

	// LookupReverseRecordQuery is a SQL query used for performing "reverse"
	// federation queries.  This query should accomodate a single parameter, using
	// "?" as the placeholder.  This provided parameter will be a strkey-encoded
	// diamcircle account id to lookup, such as
	// "GDOP3VI4UA5LS7AMLJI66RJUXEQ4HX46WUXTRTJGI5IKDLNWUBOW3FUK".
	LookupReverseRecordQuery string
}

// SQLDriver represents an implementation of `Driver` that
// provides a simple way to incorporate a SQL-backed federation handler into an
// application.  Note: this type is not designed for dynamic configuration
// changes.  Once a method is called on the struct the public fields of this
// struct should be considered frozen.
type SQLDriver struct {
	// DB is the target database that federation queries will be executed against
	DB *sql.DB

	// Dialect is the type of database peer field `DB` is communicating with.  It
	// is equivalent to the `driverName` params used in a call to `sql.Open` from
	// the standard library.
	Dialect string

	// LookupRecordQuery is a SQL query used for performing "forward" federation
	// queries.  This query should accomodate one or two parameters, using "?" as
	// the placeholder.  This provided parameters will be a name and domain
	LookupRecordQuery string

	init sync.Once
	db   *db.Session
}
