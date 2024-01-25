package tickerdb

import (
	"context"
	"testing"
	"time"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/diamcircle/go/support/db/dbtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertOrUpdateOrderbokStats(t *testing.T) {
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
	assert.WithinDuration(t, firstTime.Local(), dbAsset1.LastValid.Local(), 10*time.Millisecond)
	assert.WithinDuration(t, firstTime.Local(), dbAsset1.LastChecked.Local(), 10*time.Millisecond)

	// Creating Seconde Asset:
	secondTime := time.Now()
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

	// Creating an orderbook_stats entry:
	obTime := time.Now()
	orderbookStats := OrderbookStats{
		BaseAssetID:    dbAsset1.ID,
		CounterAssetID: dbAsset2.ID,
		NumBids:        15,
		BidVolume:      0.15,
		HighestBid:     200.0,
		NumAsks:        17,
		AskVolume:      30.0,
		LowestAsk:      0.1,
		Spread:         0.93,
		SpreadMidPoint: 0.35,
		UpdatedAt:      obTime,
	}
	err = session.InsertOrUpdateOrderbookStats(ctx,
		&orderbookStats,
		[]string{"base_asset_id", "counter_asset_id"},
	)
	require.NoError(t, err)

	var dbOS OrderbookStats
	err = session.GetRaw(ctx, &dbOS, `
		SELECT *
		FROM orderbook_stats
		ORDER BY id DESC
		LIMIT 1`,
	)
	require.NoError(t, err)

	assert.Equal(t, dbAsset1.ID, dbOS.BaseAssetID)
	assert.Equal(t, dbAsset2.ID, dbOS.CounterAssetID)
	assert.Equal(t, 15, dbOS.NumBids)
	assert.Equal(t, 0.15, dbOS.BidVolume)
	assert.Equal(t, 200.0, dbOS.HighestBid)
	assert.Equal(t, 17, dbOS.NumAsks)
	assert.Equal(t, 30.0, dbOS.AskVolume)
	assert.Equal(t, 0.1, dbOS.LowestAsk)
	assert.Equal(t, 0.93, dbOS.Spread)
	assert.Equal(t, 0.35, dbOS.SpreadMidPoint)
	assert.WithinDuration(t, obTime.Local(), dbOS.UpdatedAt.Local(), 10*time.Millisecond)

	// Making sure we're upserting:
	obTime2 := time.Now()
	orderbookStats2 := OrderbookStats{
		BaseAssetID:    dbAsset1.ID,
		CounterAssetID: dbAsset2.ID,
		NumBids:        30,
		BidVolume:      0.3,
		HighestBid:     400.0,
		NumAsks:        34,
		AskVolume:      60.0,
		LowestAsk:      0.2,
		Spread:         1.86,
		SpreadMidPoint: 0.7,
		UpdatedAt:      obTime2,
	}
	err = session.InsertOrUpdateOrderbookStats(ctx,
		&orderbookStats2,
		[]string{"base_asset_id", "counter_asset_id", "lowest_ask"},
	)
	require.NoError(t, err)

	var dbOS2 OrderbookStats
	err = session.GetRaw(ctx, &dbOS2, `
		SELECT *
		FROM orderbook_stats
		ORDER BY id DESC
		LIMIT 1`,
	)
	require.NoError(t, err)

	assert.Equal(t, dbOS2.ID, dbOS.ID) // shouldn't create another instance

	assert.Equal(t, dbAsset1.ID, dbOS2.BaseAssetID)
	assert.Equal(t, dbAsset2.ID, dbOS2.CounterAssetID)
	assert.Equal(t, 30, dbOS2.NumBids)
	assert.Equal(t, 0.3, dbOS2.BidVolume)
	assert.Equal(t, 400.0, dbOS2.HighestBid)
	assert.Equal(t, 34, dbOS2.NumAsks)
	assert.Equal(t, 60.0, dbOS2.AskVolume)
	assert.Equal(t, 0.1, dbOS2.LowestAsk) // should keep the old value, since on preserveFields
	assert.Equal(t, 1.86, dbOS2.Spread)
	assert.Equal(t, 0.7, dbOS2.SpreadMidPoint)
	assert.WithinDuration(t, obTime2.Local(), dbOS2.UpdatedAt.Local(), 10*time.Millisecond)
}
