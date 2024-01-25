package tickerdb

import (
	"context"
	"testing"
	"time"

	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/diamcircle/go/support/db/dbtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertOrUpdateAsset(t *testing.T) {
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

	publicKey := "GCF3TQXKZJNFJK7HCMNE2O2CUNKCJH2Y2ROISTBPLC7C5EIA5NNG2XZB"
	issuerAccount := "AM2FQXKZJNFJK7HCMNE2O2CUNKCJH2Y2ROISTBPLC7C5EIA5NNG2XZB"
	name := "FOO BAR"
	code := "XLM"

	// Adding a seed issuer to be used later:
	issuer := Issuer{
		PublicKey: publicKey,
		Name:      name,
	}
	tbl := session.GetTable("issuers")
	_, err = tbl.Insert(issuer).IgnoreCols("id").Exec(ctx)
	require.NoError(t, err)
	var dbIssuer Issuer
	err = session.GetRaw(ctx, &dbIssuer, `
		SELECT *
		FROM issuers
		ORDER BY id DESC
		LIMIT 1`,
	)
	require.NoError(t, err)

	// Creating first asset:
	firstTime := time.Now()
	t.Log("firstTime:", firstTime)
	a := Asset{
		Code:          code,
		IssuerAccount: issuerAccount,
		IssuerID:      dbIssuer.ID,
		LastValid:     firstTime,
		LastChecked:   firstTime,
	}
	err = session.InsertOrUpdateAsset(ctx, &a, []string{"code", "issuer_account", "issuer_id"})
	require.NoError(t, err)

	var dbAsset1 Asset
	err = session.GetRaw(ctx, &dbAsset1, `
		SELECT *
		FROM assets
		ORDER BY id DESC
		LIMIT 1`,
	)
	require.NoError(t, err)

	assert.Equal(t, code, dbAsset1.Code)
	assert.Equal(t, issuerAccount, dbAsset1.IssuerAccount)
	assert.Equal(t, dbIssuer.ID, dbAsset1.IssuerID)
	assert.Equal(
		t,
		firstTime.Local().Round(time.Millisecond),
		dbAsset1.LastValid.Local().Round(time.Millisecond),
	)
	assert.Equal(
		t,
		firstTime.Local().Round(time.Millisecond),
		dbAsset1.LastChecked.Local().Round(time.Millisecond),
	)

	// Creating Seconde Asset:
	secondTime := time.Now()
	t.Log("secondTime:", secondTime)
	a.LastValid = secondTime
	a.LastChecked = secondTime
	err = session.InsertOrUpdateAsset(ctx, &a, []string{"code", "issuer_account", "issuer_id"})
	require.NoError(t, err)

	var dbAsset2 Asset
	err = session.GetRaw(ctx, &dbAsset2, `
		SELECT *
		FROM assets
		ORDER BY id DESC
		LIMIT 1`,
	)
	require.NoError(t, err)

	// Validating if changes match what was expected:
	assert.Equal(t, dbAsset1.ID, dbAsset2.ID)
	assert.Equal(t, code, dbAsset2.Code)
	assert.Equal(t, issuerAccount, dbAsset1.IssuerAccount)
	assert.Equal(t, dbIssuer.ID, dbAsset2.IssuerID)
	assert.True(t, dbAsset2.LastValid.After(firstTime))
	assert.True(t, dbAsset2.LastChecked.After(firstTime))
	assert.Equal(
		t,
		secondTime.Local().Round(time.Millisecond),
		dbAsset2.LastValid.Local().Round(time.Millisecond),
	)
	assert.Equal(
		t,
		secondTime.Local().Round(time.Millisecond),
		dbAsset2.LastChecked.Local().Round(time.Millisecond),
	)

	// Creating Third Asset:
	thirdTime := time.Now()
	t.Log("thirdTime:", thirdTime)
	a.LastValid = thirdTime
	a.LastChecked = thirdTime
	err = session.InsertOrUpdateAsset(ctx, &a, []string{"code", "issuer_id", "last_valid", "last_checked", "issuer_account"})
	require.NoError(t, err)
	var dbAsset3 Asset
	err = session.GetRaw(ctx, &dbAsset3, `
		SELECT *
		FROM assets
		ORDER BY id DESC
		LIMIT 1`,
	)
	require.NoError(t, err)

	// Validating if changes match what was expected:
	assert.Equal(t, dbAsset2.ID, dbAsset3.ID)
	assert.Equal(t, code, dbAsset3.Code)
	assert.Equal(t, issuerAccount, dbAsset3.IssuerAccount)
	assert.Equal(t, dbIssuer.ID, dbAsset3.IssuerID)
	assert.True(t, dbAsset3.LastValid.Before(thirdTime))
	assert.True(t, dbAsset3.LastChecked.Before(thirdTime))
	assert.Equal(
		t,
		dbAsset2.LastValid.Local().Round(time.Millisecond),
		dbAsset3.LastValid.Local().Round(time.Millisecond),
	)
	assert.Equal(
		t, dbAsset2.LastValid.Local().Round(time.Millisecond),
		dbAsset3.LastChecked.Local().Round(time.Millisecond),
	)
}

func TestGetAssetByCodeAndIssuerAccount(t *testing.T) {
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

	publicKey := "GCF3TQXKZJNFJK7HCMNE2O2CUNKCJH2Y2ROISTBPLC7C5EIA5NNG2XZB"
	name := "FOO BAR"
	code := "XLM"
	issuerAccount := "AM2FQXKZJNFJK7HCMNE2O2CUNKCJH2Y2ROISTBPLC7C5EIA5NNG2XZB"

	// Adding a seed issuer to be used later:
	issuer := Issuer{
		PublicKey: publicKey,
		Name:      name,
	}
	tbl := session.GetTable("issuers")
	_, err = tbl.Insert(issuer).IgnoreCols("id").Exec(ctx)
	require.NoError(t, err)
	var dbIssuer Issuer
	err = session.GetRaw(ctx, &dbIssuer, `
		SELECT *
		FROM issuers
		ORDER BY id DESC
		LIMIT 1`,
	)
	require.NoError(t, err)

	// Creating first asset:
	firstTime := time.Now()
	a := Asset{
		Code:          code,
		IssuerAccount: issuerAccount,
		IssuerID:      dbIssuer.ID,
		LastValid:     firstTime,
		LastChecked:   firstTime,
	}
	err = session.InsertOrUpdateAsset(ctx, &a, []string{"code", "issuer_account", "issuer_id"})
	require.NoError(t, err)

	var dbAsset Asset
	err = session.GetRaw(ctx, &dbAsset, `
		SELECT *
		FROM assets
		ORDER BY id DESC
		LIMIT 1`,
	)
	require.NoError(t, err)

	// Searching for an asset that exists:
	found, id, err := session.GetAssetByCodeAndIssuerAccount(ctx, code, issuerAccount)
	require.NoError(t, err)
	assert.Equal(t, dbAsset.ID, id)
	assert.True(t, found)

	// Now searching for an asset that does not exist:
	found, _, err = session.GetAssetByCodeAndIssuerAccount(ctx,
		"NONEXISTENT CODE",
		issuerAccount,
	)
	require.NoError(t, err)
	assert.False(t, found)
}
