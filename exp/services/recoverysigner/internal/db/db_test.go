package db

import (
	"fmt"
	"net"
	"regexp"
	"testing"

	"github.com/diamcircle/go/support/db/dbtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpen_openAndPingSucceeds(t *testing.T) {
	db := dbtest.Postgres(t)

	sqlxDB, err := Open(db.DSN)
	require.NoError(t, err)
	assert.Equal(t, "postgres", sqlxDB.DriverName())

	err = sqlxDB.Ping()
	require.NoError(t, err)
}

func TestOpen_openAndPingFails(t *testing.T) {
	// Find an empty port
	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	port := listener.Addr().(*net.TCPAddr).Port
	require.NoError(t, listener.Close())
	// Slight race here with other stuff on the system, which could claim this port.

	sqlxDB, err := Open(fmt.Sprintf("postgres://localhost:%d", port))
	require.NoError(t, err)
	assert.Equal(t, "postgres", sqlxDB.DriverName())

	err = sqlxDB.Ping()
	require.Error(t, err)
	require.Regexp(
		t,
		regexp.MustCompile(
			// regex to support both ipv4 and ipv6, on the port we found.
			fmt.Sprintf("dial tcp (127\\.0\\.0\\.1|\\[::1\\]):%d: connect: connection refused", port),
		),
		err.Error(),
	)
}
