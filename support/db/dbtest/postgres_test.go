package dbtest

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgres(t *testing.T) {
	db := Postgres(t)
	t.Log("tempdb url", db.DSN)

	conn := db.Open()
	_, err := conn.Exec("SELECT 1")
	require.NoError(t, err)
	require.NoError(t, conn.Close())

	db.Close()

	conn = db.Open()
	_, err = conn.Exec("SELECT 1")
	require.EqualError(t, err, fmt.Sprintf("pq: database \"%s\" does not exist", db.dbName))
	require.Contains(t, err.Error(), "data")

	// ensure Close() can be called multiple times
	db.Close()
}

func TestPostgres_clientTimezone(t *testing.T) {
	db := Postgres(t)
	conn := db.Open()
	defer conn.Close()

	timestamp := time.Time{}
	err := conn.Get(&timestamp, "SELECT TO_TIMESTAMP('2020-03-19 16:56:00', 'YYYY-MM-DD HH24:MI:SS')")
	require.NoError(t, err)

	wantTimestamp := time.Date(2020, 3, 19, 16, 56, 0, 0, time.UTC)
	assert.Equal(t, wantTimestamp, timestamp)
}

func TestPostgres_Version(t *testing.T) {
	db := Postgres(t)
	majorVersion := db.Version()
	assert.GreaterOrEqual(t, majorVersion, 9)
}
