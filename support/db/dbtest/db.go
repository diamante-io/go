package dbtest

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/diamcircle/go/support/db/sqlutils"
	"github.com/diamcircle/go/support/errors"
	"github.com/stretchr/testify/require"
)

// DB represents an ephemeral database that starts blank and can be used
// to run tests against.
type DB struct {
	Dialect string
	DSN     string
	dbName  string
	t       testing.TB
	closer  func()
	closed  bool
}

// randomName returns a new psuedo-random name that is sufficient for naming a
// test database.  In the event that reading from the source of randomness
// fails, a panic will occur.
func randomName() string {
	raw := make([]byte, 6)

	_, err := rand.Read(raw)
	if err != nil {
		err = errors.Wrap(err, "read from rand failed")
		panic(err)
	}

	enc := hex.EncodeToString(raw)

	return fmt.Sprintf("test_%s", enc)
}

// Close closes and deletes the database represented by `db`
func (db *DB) Close() {
	if db.closed {
		return
	}

	db.closer()
	db.closed = true
}

// Load executes all of the statements in the provided sql script against the
// test database, panicking if any fail.  The receiver is returned allowing for
// chain-style calling within your test functions.
func (db *DB) Load(sql string) *DB {
	conn := db.Open()
	defer conn.Close()

	tx, err := conn.Begin()
	require.NoError(db.t, err)

	defer tx.Rollback()

	for i, cmd := range sqlutils.AllStatements(sql) {
		_, err = tx.Exec(cmd)
		require.NoError(db.t, err, "failed execing statement: %d", i)
	}

	err = tx.Commit()
	require.NoError(db.t, err)

	return db
}

// Open opens a sqlx connection to the db.
func (db *DB) Open() *sqlx.DB {
	conn, err := sqlx.Open(db.Dialect, db.DSN)
	require.NoError(db.t, err)

	return conn
}

func (db *DB) Version() (major int) {
	conn := db.Open()
	defer conn.Close()

	versionFull := ""
	err := conn.Get(&versionFull, "SHOW server_version")
	require.NoError(db.t, err)

	version := strings.Fields(versionFull)
	parts := strings.Split(version[0], ".")
	major, err = strconv.Atoi(parts[0])
	require.NoError(db.t, err)

	return major
}

func execStatement(t testing.TB, pguser, query string) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("postgres://%s@localhost/?sslmode=disable", pguser))
	require.NoError(t, err)
	_, err = db.Exec(query)
	require.NoError(t, err)
	require.NoError(t, db.Close())
}

// Postgres provisions a new, blank database with a random name on the localhost
// of the running process.  It assumes that you have postgres running on the
// default port, have the command line postgres tools installed, and that the
// current user has access to the server.  It panics on the event of a failure.
func Postgres(t testing.TB) *DB {
	var result DB
	result.dbName = randomName()
	result.Dialect = "postgres"
	result.t = t

	t.Log("Test Database:", result.dbName)

	pgUser := os.Getenv("PGUSER")
	if len(pgUser) == 0 {
		pgUser = "postgres"
	}

	// create the db
	execStatement(t, pgUser, "CREATE DATABASE "+pq.QuoteIdentifier(result.dbName))

	result.DSN = fmt.Sprintf("postgres://%s@localhost/%s?sslmode=disable&timezone=UTC", pgUser, result.dbName)

	result.closer = func() {
		execStatement(t, pgUser, "DROP DATABASE "+pq.QuoteIdentifier(result.dbName))
	}

	return &result
}
