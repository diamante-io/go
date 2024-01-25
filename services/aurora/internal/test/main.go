// Package test contains simple test helpers that should not
// have any dependencies on aurora's packages.  think constants,
// custom matchers, generic helpers etc.
package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"go/support/log"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// StaticMockServer is a test helper that records it's last request
type StaticMockServer struct {
	*httptest.Server
	LastRequest *http.Request
}

// T provides a common set of functionality for each test in aurora
type T struct {
	T          *testing.T
	Assert     *assert.Assertions
	Require    *require.Assertions
	Ctx        context.Context
	auroraDB   *sqlx.DB
	CoreDB     *sqlx.DB
	EndLogTest func() []logrus.Entry
}

// Context provides a context suitable for testing in tests that do not create
// a full App instance (in which case your tests should be using the app's
// context).  This context has a logger bound to it suitable for testing.
func Context() context.Context {
	return log.Set(context.Background(), testLogger)
}

// Database returns a connection to the aurora test database
//
// DEPRECATED:  use `aurora()` from test/db package
func Database(t *testing.T) *sqlx.DB {
	return tdb.aurora(t)
}

// DatabaseURL returns the database connection the url any test
// use when connecting to the history/aurora database
//
// DEPRECATED:  use `auroraURL()` from test/db package
func DatabaseURL() string {
	return tdb.auroraURL()
}

// Start initializes a new test helper object, a new instance of log,
// and conceptually "starts" a new test
func Start(t *testing.T) *T {
	result := &T{}
	result.T = t
	logger := log.New()

	result.Ctx = log.Set(context.Background(), logger)
	result.auroraDB = Database(t)
	result.CoreDB = diamcircleCoreDatabase(t)
	result.Assert = assert.New(t)
	result.Require = require.New(t)
	result.EndLogTest = logger.StartTest(log.DebugLevel)

	return result
}

// diamcircleCoreDatabase returns a connection to the diamcircle core test database
//
// DEPRECATED:  use `diamcircleCore()` from test/db package
func diamcircleCoreDatabase(t *testing.T) *sqlx.DB {
	return tdb.diamcircleCore(t)
}

// diamcircleCoreDatabaseURL returns the database connection the url any test
// use when connecting to the diamcircle-core database
//
// DEPRECATED:  use `diamcircleCoreURL()` from test/db package
func diamcircleCoreDatabaseURL() string {
	return tdb.diamcircleCoreURL()
}
