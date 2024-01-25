package keystore

import (
	"testing"

	"go/support/db/dbtest"

	migrate "github.com/rubenv/sql-migrate"
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
