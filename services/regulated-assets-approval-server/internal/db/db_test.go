package db

import (
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
	sqlxDB, err := Open("postgres://127.0.0.1:0")
	require.NoError(t, err)
	assert.Equal(t, "postgres", sqlxDB.DriverName())

	err = sqlxDB.Ping()
	require.Error(t, err)
	require.Contains(t, err.Error(), "dial tcp 127.0.0.1:0: connect:")
}
