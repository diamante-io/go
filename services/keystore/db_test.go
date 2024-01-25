package keystore

import (
	"testing"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/diamcircle/go/support/db/dbtest"
)

//TODO: creating a DB for every single test is inefficient. Maybe we can
//improve our dbtest package so that we can just get a transaction.
func openKeystoreDB(t *testing.T) *dbtest.DB {
	db := dbtest.Postgres(t)
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}

	conn := db.Open()
	defer conn.Close()

	_, err := migrate.Exec(conn.DB, "postgres", migrations, migrate.Up)
	if err != nil {
		t.Fatal(err)
	}
	return db
}
