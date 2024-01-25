package tickerdb

import (
	"context"
	"testing"

	"go/support/db/dbtest"

	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertOrUpdateIssuer(t *testing.T) {
	db := dbtest.Postgres(t)
	defer db.Close()

	var session TickerSession
	session.DB = db.Open()
	ctx := context.Background()
	defer session.DB.Close()

	// Run migrations to make sure the tests are run
	// on the most updated schema version
	migrations := &migrate.FileMigrationSource{
		Dir: "./migrations",
	}
	_, err := migrate.Exec(session.DB.DB, "postgres", migrations, migrate.Up)
	require.NoError(t, err)

	publicKey := "ASOKDASDKMAKSD19023ASDSAD0912309"
	name := "FOO BAR"

	// Adding a seed issuer to be used later:
	issuer := Issuer{
		PublicKey: publicKey,
		Name:      name,
	}
	id, err := session.InsertOrUpdateIssuer(ctx, &issuer, []string{"public_key"})

	require.NoError(t, err)
	var dbIssuer Issuer
	err = session.GetRaw(ctx, &dbIssuer, `
		SELECT *
		FROM issuers
		ORDER BY id DESC
		LIMIT 1`,
	)
	require.NoError(t, err)

	assert.Equal(t, publicKey, dbIssuer.PublicKey)
	assert.Equal(t, dbIssuer.ID, id)

	// Adding another issuer to validate we're correctly returning the ID
	issuer2 := Issuer{
		PublicKey: "ANOTHERKEY",
		Name:      "Hello from the other side",
	}
	id2, err := session.InsertOrUpdateIssuer(ctx, &issuer2, []string{"public_key"})

	require.NoError(t, err)
	var dbIssuer2 Issuer
	err = session.GetRaw(ctx, &dbIssuer2, `
		SELECT *
		FROM issuers
		ORDER BY id DESC
		LIMIT 1`,
	)
	require.NoError(t, err)

	assert.Equal(t, issuer2.Name, dbIssuer2.Name)
	assert.Equal(t, issuer2.PublicKey, dbIssuer2.PublicKey)
	assert.Equal(t, id2, dbIssuer2.ID)

	// Validate if it only changes the un-preserved fields
	name3 := "The Dark Side of the Moon"
	issuer3 := Issuer{
		PublicKey: publicKey,
		Name:      name3,
	}
	id, err = session.InsertOrUpdateIssuer(ctx, &issuer3, []string{"public_key"})
	require.NoError(t, err)

	var dbIssuer3 Issuer
	err = session.GetRaw(ctx,
		&dbIssuer3,
		"SELECT * FROM issuers WHERE id=?",
		id,
	)
	require.NoError(t, err)

	assert.Equal(t, dbIssuer.ID, dbIssuer3.ID)
	assert.Equal(t, dbIssuer.PublicKey, dbIssuer3.PublicKey)
	assert.Equal(t, name3, dbIssuer3.Name)
}
