package tickerdbtest

import (
	"context"
	"testing"
	"time"

	"go/services/ticker/internal/tickerdb"
	"go/support/db/dbtest"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SetupTickerTestSession sets up the database for testing the GraphQL endpoints
// and associated query logic.
func SetupTickerTestSession(t *testing.T, migrationsDir string) (session tickerdb.TickerSession) {
	db := dbtest.Postgres(t)
	session.DB = db.Open()
	ctx := context.Background()

	// Run migrations to make sure the tests are run
	// on the most updated schema version
	migrations := &migrate.FileMigrationSource{
		Dir: migrationsDir,
	}
	_, err := migrate.Exec(session.DB.DB, "postgres", migrations, migrate.Up)
	require.NoError(t, err)

	// Adding a seed issuer to be used later:
	issuer1PK := "GCF3TQXKZJNFJK7HCMNE2O2CUNKCJH2Y2ROISTBPLC7C5EIA5NNG2XZB"
	tbl := session.GetTable("issuers")
	_, err = tbl.Insert(tickerdb.Issuer{
		PublicKey: issuer1PK,
		Name:      "FOO BAR",
	}).IgnoreCols("id").Exec(ctx)
	require.NoError(t, err)

	var issuer1 tickerdb.Issuer
	err = session.GetRaw(ctx, &issuer1, `
		SELECT *
		FROM issuers
		WHERE public_key = ?`,
		issuer1PK,
	)
	require.NoError(t, err)

	// Adding another issuer to be used later:
	issuer2PK := "ABF3TQXKZJNFJK7HCMNE2O2CUNKCJH2Y2ROISTBPLC7C5EIA5NNG2XZB"
	_, err = tbl.Insert(tickerdb.Issuer{
		PublicKey: issuer2PK,
		Name:      "FOO BAR",
	}).IgnoreCols("id").Exec(ctx)
	require.NoError(t, err)

	var issuer2 tickerdb.Issuer
	err = session.GetRaw(ctx, &issuer2, `
		SELECT *
		FROM issuers
		WHERE public_key = ?`,
		issuer2PK,
	)
	require.NoError(t, err)

	// Adding a seed asset to be used later:
	err = session.InsertOrUpdateAsset(ctx, &tickerdb.Asset{
		Code:          "ETH",
		IssuerAccount: issuer1PK,
		IssuerID:      issuer1.ID,
		IsValid:       true,
	}, []string{"code", "issuer_id"})
	require.NoError(t, err)
	var ethAsset1 tickerdb.Asset
	err = session.GetRaw(ctx, &ethAsset1, `
		SELECT *
		FROM assets
		WHERE code = ?
		AND issuer_account = ?`,
		"ETH",
		issuer1PK,
	)
	require.NoError(t, err)

	// Adding a seed asset to be used later:
	err = session.InsertOrUpdateAsset(ctx, &tickerdb.Asset{
		Code:          "ETH",
		IssuerAccount: issuer2PK,
		IssuerID:      issuer2.ID,
		IsValid:       true,
	}, []string{"code", "issuer_id"})
	require.NoError(t, err)

	var ethAsset2 tickerdb.Asset
	err = session.GetRaw(ctx, &ethAsset2, `
		SELECT *
		FROM assets
		WHERE code = ?
		AND issuer_account = ?`,
		"ETH",
		issuer2PK,
	)
	require.NoError(t, err)

	// Adding another asset to be used later:
	err = session.InsertOrUpdateAsset(ctx, &tickerdb.Asset{
		Code:          "BTC",
		IssuerAccount: issuer1PK,
		IssuerID:      issuer1.ID,
		IsValid:       true,
	}, []string{"code", "issuer_id"})
	require.NoError(t, err)

	var btcAsset tickerdb.Asset
	err = session.GetRaw(ctx, &btcAsset, `
		SELECT *
		FROM assets
		WHERE code = ?
		AND issuer_account = ?`,
		"BTC",
		issuer1PK,
	)
	require.NoError(t, err)

	// A few times to be used:
	now := time.Now()
	tenMinutesAgo := now.Add(-10 * time.Minute)
	oneHourAgo := now.Add(-1 * time.Hour)
	threeDaysAgo := now.AddDate(0, 0, -3)

	// Now let's create the trades:
	trades := []tickerdb.Trade{
		tickerdb.Trade{ // BTC_ETH  trade (ETH is from issuer 1)
			auroraID:        "hrzid1",
			BaseAssetID:     btcAsset.ID,
			BaseAmount:      100.0,
			CounterAssetID:  ethAsset1.ID,
			CounterAmount:   10.0,
			Price:           0.1,
			LedgerCloseTime: tenMinutesAgo,
		},
		tickerdb.Trade{ // BTC_ETH trade (ETH is from issuer 2)
			auroraID:        "hrzid3",
			BaseAssetID:     btcAsset.ID,
			BaseAmount:      24.0,
			CounterAssetID:  ethAsset2.ID,
			CounterAmount:   26.0,
			Price:           0.92,
			LedgerCloseTime: now,
		},
		tickerdb.Trade{ // BTC_ETH  trade (ETH is from issuer 1)
			auroraID:        "hrzid2",
			BaseAssetID:     btcAsset.ID,
			BaseAmount:      50.0,
			CounterAssetID:  ethAsset1.ID,
			CounterAmount:   50.0,
			Price:           1.0,
			LedgerCloseTime: oneHourAgo,
		},
		tickerdb.Trade{ // BTC_ETH  trade (ETH is from issuer 1)
			auroraID:        "hrzid4",
			BaseAssetID:     btcAsset.ID,
			BaseAmount:      50.0,
			CounterAssetID:  ethAsset1.ID,
			CounterAmount:   6.0,
			Price:           0.12,
			LedgerCloseTime: threeDaysAgo,
		},
	}
	err = session.BulkInsertTrades(ctx, trades)
	require.NoError(t, err)

	// Adding some orderbook stats:
	obTime := time.Now()
	orderbookStats := tickerdb.OrderbookStats{
		BaseAssetID:    btcAsset.ID,
		CounterAssetID: ethAsset1.ID,
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

	var obBTCETH1 tickerdb.OrderbookStats
	err = session.GetRaw(ctx, &obBTCETH1, `
		SELECT *
		FROM orderbook_stats
		ORDER BY id DESC
		LIMIT 1`,
	)
	require.NoError(t, err)

	orderbookStats = tickerdb.OrderbookStats{
		BaseAssetID:    ethAsset1.ID,
		CounterAssetID: btcAsset.ID,
		NumBids:        10,
		BidVolume:      0.90,
		HighestBid:     100.0,
		NumAsks:        12,
		AskVolume:      25.0,
		LowestAsk:      0.2,
		Spread:         0.55,
		SpreadMidPoint: 0.85,
	}
	err = session.InsertOrUpdateOrderbookStats(ctx,
		&orderbookStats,
		[]string{"base_asset_id", "counter_asset_id"},
	)
	require.NoError(t, err)

	var obETH1BTC tickerdb.OrderbookStats
	err = session.GetRaw(ctx, &obETH1BTC, `
		SELECT *
		FROM orderbook_stats
		ORDER BY id DESC
		LIMIT 1`,
	)
	require.NoError(t, err)

	orderbookStats = tickerdb.OrderbookStats{
		BaseAssetID:    btcAsset.ID,
		CounterAssetID: ethAsset2.ID,
		NumBids:        1,
		BidVolume:      0.1,
		HighestBid:     20.0,
		NumAsks:        1,
		AskVolume:      15.0,
		LowestAsk:      0.2,
		Spread:         0.96,
		SpreadMidPoint: 0.36,
		UpdatedAt:      obTime,
	}
	err = session.InsertOrUpdateOrderbookStats(ctx,
		&orderbookStats,
		[]string{"base_asset_id", "counter_asset_id"},
	)
	require.NoError(t, err)

	var obBTCETH2 tickerdb.OrderbookStats
	err = session.GetRaw(ctx, &obBTCETH2, `
		SELECT *
		FROM orderbook_stats
		ORDER BY id DESC
		LIMIT 1`,
	)
	require.NoError(t, err)
	assert.NotEqual(t, obBTCETH1.ID, obBTCETH2.ID)

	orderbookStats = tickerdb.OrderbookStats{
		BaseAssetID:    ethAsset2.ID,
		CounterAssetID: btcAsset.ID,
		NumBids:        20,
		BidVolume:      0.60,
		HighestBid:     300.0,
		NumAsks:        20,
		AskVolume:      256.0,
		LowestAsk:      0.70,
		Spread:         150.0,
		SpreadMidPoint: 200.0,
	}
	err = session.InsertOrUpdateOrderbookStats(ctx,
		&orderbookStats,
		[]string{"base_asset_id", "counter_asset_id"},
	)
	require.NoError(t, err)

	var obETH2BTC tickerdb.OrderbookStats
	err = session.GetRaw(ctx, &obETH2BTC, `
		SELECT *
		FROM orderbook_stats
		ORDER BY id DESC
		LIMIT 1`,
	)
	require.NoError(t, err)

	// Add an XLM asset.
	err = session.InsertOrUpdateAsset(ctx, &tickerdb.Asset{
		Code:          "XLM",
		IssuerAccount: issuer1PK,
		IssuerID:      issuer1.ID,
		IsValid:       true,
	}, []string{"code", "issuer_id"})
	require.NoError(t, err)

	var xlmAsset tickerdb.Asset
	err = session.GetRaw(ctx, &xlmAsset, `
		SELECT *
		FROM assets
		WHERE code = ?
		AND issuer_account = ?`,
		"XLM",
		issuer1PK,
	)
	require.NoError(t, err)

	// Add XLM/BTC trades.
	trades = []tickerdb.Trade{
		tickerdb.Trade{
			auroraID:        "hrzid5",
			BaseAssetID:     xlmAsset.ID,
			BaseAmount:      10.0,
			CounterAssetID:  btcAsset.ID,
			CounterAmount:   10.0,
			Price:           0.5, // close price & lowest price
			LedgerCloseTime: tenMinutesAgo,
		},
		tickerdb.Trade{
			auroraID:        "hrzid6",
			BaseAssetID:     xlmAsset.ID,
			BaseAmount:      10.0,
			CounterAssetID:  btcAsset.ID,
			CounterAmount:   10.0,
			Price:           1.0, // open price & highest price
			LedgerCloseTime: now,
		},
	}
	err = session.BulkInsertTrades(ctx, trades)
	require.NoError(t, err)

	return
}
