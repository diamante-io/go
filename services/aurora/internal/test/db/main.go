// Package db provides helpers to connect to test databases.  It has no
// internal dependencies on aurora and so should be able to be imported by
// any aurora package.
package db

import (
	"fmt"
	"log"
	"testing"

	"github.com/jmoiron/sqlx"
	// pq enables postgres support
	db "go/support/db/dbtest"

	_ "github.com/lib/pq"
)

var (
	coreDB    *sqlx.DB
	coreUrl   *string
	auroraDB  *sqlx.DB
	auroraUrl *string
)

// aurora returns a connection to the aurora test database
func aurora(t *testing.T) *sqlx.DB {
	if auroraDB != nil {
		return auroraDB
	}
	postgres := db.Postgres(t)
	auroraUrl = &postgres.DSN
	auroraDB = postgres.Open()

	return auroraDB
}

// auroraURL returns the database connection the url any test
// use when connecting to the history/aurora database
func auroraURL() string {
	if auroraUrl == nil {
		log.Panic(fmt.Errorf("aurora not initialized"))
	}
	return *auroraUrl
}

// diamcircleCore returns a connection to the diamcircle core test database
func diamcircleCore(t *testing.T) *sqlx.DB {
	if coreDB != nil {
		return coreDB
	}
	postgres := db.Postgres(t)
	coreUrl = &postgres.DSN
	coreDB = postgres.Open()
	return coreDB
}

// diamcircleCoreURL returns the database connection the url any test
// use when connecting to the diamcircle-core database
func diamcircleCoreURL() string {
	if coreUrl == nil {
		log.Panic(fmt.Errorf("diamcircleCore not initialized"))
	}
	return *coreUrl
}
