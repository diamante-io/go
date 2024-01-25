package tickerdb

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go/support/db/dbtest"

	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/require"
)

func TestBulkInsertTrades(t *testing.T) {
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

	// Adding a seed issuer to be used later:
	tbl := session.GetTable("issuers")
	_, err = tbl.Insert(Issuer{
		PublicKey: "GCF3TQXKZJNFJK7HCMNE2O2CUNKCJH2Y2ROISTBPLC7C5EIA5NNG2XZB",
		Name:      "FOO BAR",
	}).IgnoreCols("id").Exec(ctx)
	require.NoError(t, err)
	var issuer Issuer
	err = session.GetRaw(ctx, &issuer, `
		SELECT *
		FROM issuers
		ORDER BY id DESC
		LIMIT 1`,
	)
	require.NoError(t, err)

	// Adding a seed asset to be used later:
	err = session.InsertOrUpdateAsset(ctx, &Asset{
		Code:     "XLM",
		IssuerID: issuer.ID,
	}, []string{"code", "issuer_id"})
	require.NoError(t, err)
	var asset1 Asset
	err = session.GetRaw(ctx, &asset1, `
		SELECT *
		FROM assets
		ORDER BY id DESC
		LIMIT 1`,
	)
	require.NoError(t, err)

	// Adding another asset to be used later:
	err = session.InsertOrUpdateAsset(ctx, &Asset{
		Code:     "BTC",
		IssuerID: issuer.ID,
	}, []string{"code", "issuer_id"})
	require.NoError(t, err)
	var asset2 Asset
	err = session.GetRaw(ctx, &asset2, `
		SELECT *
		FROM assets
		ORDER BY id DESC
		LIMIT 1`,
	)
	require.NoError(t, err)

	// Verify that we actually have two assets:
	assert.NotEqual(t, asset1.ID, asset2.ID)

	// Now let's create the trades:
	trades := []Trade{
		Trade{
			auroraID:        "hrzid1",
			BaseAssetID:     asset1.ID,
			CounterAssetID:  asset2.ID,
			LedgerCloseTime: time.Now(),
		},
		Trade{
			auroraID:        "hrzid2",
			BaseAssetID:     asset2.ID,
			CounterAssetID:  asset1.ID,
			LedgerCloseTime: time.Now(),
		},
	}
	err = session.BulkInsertTrades(ctx, trades)
	require.NoError(t, err)

	// Ensure only two were created:
	rows, err := session.QueryRaw(ctx, "SELECT * FROM trades")
	require.NoError(t, err)
	rowsCount := 0
	for rows.Next() {
		rowsCount++
	}
	assert.Equal(t, 2, rowsCount)

	// Re-insert the same trades and check if count remains = 2:
	err = session.BulkInsertTrades(ctx, trades)
	require.NoError(t, err)

	rows, err = session.QueryRaw(ctx, "SELECT * FROM trades")
	require.NoError(t, err)
	rowsCount2 := 0
	for rows.Next() {
		rowsCount2++
	}
	assert.Equal(t, 2, rowsCount2)
}

func TestGetLastTrade(t *testing.T) {
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

	// Sanity Check (there are no trades in the database)
	_, err = session.GetLastTrade(ctx)
	require.Error(t, err)

	// Adding a seed issuer to be used later:
	tbl := session.GetTable("issuers")
	_, err = tbl.Insert(Issuer{
		PublicKey: "GCF3TQXKZJNFJK7HCMNE2O2CUNKCJH2Y2ROISTBPLC7C5EIA5NNG2XZB",
		Name:      "FOO BAR",
	}).IgnoreCols("id").Exec(ctx)
	require.NoError(t, err)
	var issuer Issuer
	err = session.GetRaw(ctx, &issuer, `
		SELECT *
		FROM issuers
		ORDER BY id DESC
		LIMIT 1`,
	)
	require.NoError(t, err)

	// Adding a seed asset to be used later:
	err = session.InsertOrUpdateAsset(ctx, &Asset{
		Code:     "XLM",
		IssuerID: issuer.ID,
	}, []string{"code", "issuer_id"})
	require.NoError(t, err)
	var asset1 Asset
	err = session.GetRaw(ctx, &asset1, `
		SELECT *
		FROM assets
		ORDER BY id DESC
		LIMIT 1`,
	)
	require.NoError(t, err)

	// Adding another asset to be used later:
	err = session.InsertOrUpdateAsset(ctx, &Asset{
		Code:     "BTC",
		IssuerID: issuer.ID,
	}, []string{"code", "issuer_id"})
	require.NoError(t, err)
	var asset2 Asset
	err = session.GetRaw(ctx, &asset2, `
		SELECT *
		FROM assets
		ORDER BY id DESC
		LIMIT 1`,
	)
	require.NoError(t, err)

	// Verify that we actually have two assets:
	assert.NotEqual(t, asset1.ID, asset2.ID)

	now := time.Now()
	oneYearBefore := now.AddDate(-1, 0, 0)

	// Now let's create the trades:
	trades := []Trade{
		Trade{
			auroraID:        "hrzid2",
			BaseAssetID:     asset2.ID,
			CounterAssetID:  asset1.ID,
			LedgerCloseTime: oneYearBefore,
		},
		Trade{
			auroraID:        "hrzid1",
			BaseAssetID:     asset1.ID,
			CounterAssetID:  asset2.ID,
			LedgerCloseTime: now,
		},
		Trade{
			auroraID:        "hrzid2",
			BaseAssetID:     asset2.ID,
			CounterAssetID:  asset1.ID,
			LedgerCloseTime: oneYearBefore,
		},
	}

	// Re-insert the same trades and check if count remains = 2:
	err = session.BulkInsertTrades(ctx, trades)
	require.NoError(t, err)

	lastTrade, err := session.GetLastTrade(ctx)
	require.NoError(t, err)
	assert.WithinDuration(t, now.Local(), lastTrade.LedgerCloseTime.Local(), 10*time.Millisecond)
}

func TestDeleteOldTrades(t *testing.T) {
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

	// Adding a seed issuer to be used later:
	tbl := session.GetTable("issuers")
	_, err = tbl.Insert(Issuer{
		PublicKey: "GCF3TQXKZJNFJK7HCMNE2O2CUNKCJH2Y2ROISTBPLC7C5EIA5NNG2XZB",
		Name:      "FOO BAR",
	}).IgnoreCols("id").Exec(ctx)
	require.NoError(t, err)
	var issuer Issuer
	err = session.GetRaw(ctx, &issuer, `
		SELECT *
		FROM issuers
		ORDER BY id DESC
		LIMIT 1`,
	)
	require.NoError(t, err)

	// Adding a seed asset to be used later:
	err = session.InsertOrUpdateAsset(ctx, &Asset{
		Code:     "XLM",
		IssuerID: issuer.ID,
	}, []string{"code", "issuer_id"})
	require.NoError(t, err)
	var asset1 Asset
	err = session.GetRaw(ctx, &asset1, `
		SELECT *
		FROM assets
		ORDER BY id DESC
		LIMIT 1`,
	)
	require.NoError(t, err)

	// Adding another asset to be used later:
	err = session.InsertOrUpdateAsset(ctx, &Asset{
		Code:     "BTC",
		IssuerID: issuer.ID,
	}, []string{"code", "issuer_id"})
	require.NoError(t, err)
	var asset2 Asset
	err = session.GetRaw(ctx, &asset2, `
		SELECT *
		FROM assets
		ORDER BY id DESC
		LIMIT 1`,
	)
	require.NoError(t, err)

	// Verify that we actually have two assets:
	assert.NotEqual(t, asset1.ID, asset2.ID)

	// Setting up some times for testing
	now := time.Now()
	oneDayAgo := now.AddDate(0, 0, -1)
	oneMonthAgo := now.AddDate(0, -1, 0)
	oneYearAgo := now.AddDate(-1, 0, 0)

	// Now let's create the trades:
	trades := []Trade{
		Trade{
			auroraID:        "hrzid1",
			BaseAssetID:     asset1.ID,
			CounterAssetID:  asset2.ID,
			LedgerCloseTime: now,
		},
		Trade{
			auroraID:        "hrzid2",
			BaseAssetID:     asset2.ID,
			CounterAssetID:  asset1.ID,
			LedgerCloseTime: oneDayAgo,
		},
		Trade{
			auroraID:        "hrzid3",
			BaseAssetID:     asset2.ID,
			CounterAssetID:  asset1.ID,
			LedgerCloseTime: oneMonthAgo,
		},
		Trade{
			auroraID:        "hrzid4",
			BaseAssetID:     asset2.ID,
			CounterAssetID:  asset1.ID,
			LedgerCloseTime: oneYearAgo,
		},
	}
	err = session.BulkInsertTrades(ctx, trades)
	require.NoError(t, err)

	// Deleting trades older than 1 day ago:
	err = session.DeleteOldTrades(ctx, oneDayAgo)
	require.NoError(t, err)

	var dbTrades []Trade
	var trade1, trade2 Trade
	err = session.SelectRaw(ctx, &dbTrades, "SELECT * FROM trades")
	require.NoError(t, err)
	assert.Equal(t, 2, len(dbTrades))

	// Make sure we're actually deleting the entries we wanted:
	for i, trade := range dbTrades {
		if trade.auroraID == "hrzid1" {
			trade1 = dbTrades[i]
		}

		if trade.auroraID == "hrzid2" {
			trade2 = dbTrades[i]
		}
	}

	assert.NotEqual(t, trade1.auroraID, "")
	assert.NotEqual(t, trade2.auroraID, "")
	assert.WithinDuration(t, now.Local(), trade1.LedgerCloseTime.Local(), 10*time.Millisecond)
	assert.WithinDuration(t, oneDayAgo.Local(), trade2.LedgerCloseTime.Local(), 10*time.Millisecond)
}
