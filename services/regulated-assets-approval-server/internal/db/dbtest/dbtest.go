package dbtest

import (
	"path"
	"runtime"
	"testing"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/diamcircle/go/support/db/dbtest"
)

func OpenWithoutMigrations(t *testing.T) *dbtest.DB {
	db := dbtest.Postgres(t)
	return db
}

func Open(t *testing.T) *dbtest.DB {
	db := OpenWithoutMigrations(t)

	// Get the folder holding the migrations relative to this file. We cannot
	// hardcode "../dbmigrate/migrations" because Open is called from tests in
	// multiple packages and tests are executed with the current working
	// directory set to the package the test lives in.
	_, filename, _, _ := runtime.Caller(0)
	migrationsDir := path.Join(path.Dir(filename), "..", "dbmigrate", "migrations")

	migrations := &migrate.FileMigrationSource{
		Dir: migrationsDir,
	}

	conn := db.Open()
	defer conn.Close()

	_, err := migrate.Exec(conn.DB, "postgres", migrations, migrate.Up)
	if err != nil {
		t.Fatal(err)
	}
	return db
}
