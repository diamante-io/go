package integration

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"go/clients/auroraclient"
	"go/keypair"
	aurora "go/services/aurora/internal"
	"go/services/aurora/internal/db2/history"
	"go/xdr"

	"github.com/stretchr/testify/assert"

	"go/historyarchive"
	auroracmd "go/services/aurora/cmd"
	"go/services/aurora/internal/db2/schema"
	"go/services/aurora/internal/test/integration"
	"go/support/db"
	"go/support/db/dbtest"
	"go/txnbuild"
)

func generateLiquidityPoolOps(itest *integration.Test, tt *assert.Assertions) (lastLedger int32) {

	master := itest.Master()
	keys, accounts := itest.CreateAccounts(2, "1000")
	shareKeys, shareAccount := keys[0], accounts[0]
	tradeKeys, tradeAccount := keys[1], accounts[1]

	itest.MustSubmitMultiSigOperations(shareAccount, []*keypair.Full{shareKeys, master},
		&txnbuild.ChangeTrust{
			Line: txnbuild.ChangeTrustAssetWrapper{
				Asset: txnbuild.CreditAsset{
					Code:   "USD",
					Issuer: master.Address(),
				},
			},
			Limit: txnbuild.MaxTrustlineLimit,
		},
		&txnbuild.ChangeTrust{
			Line: txnbuild.LiquidityPoolShareChangeTrustAsset{
				LiquidityPoolParameters: txnbuild.LiquidityPoolParameters{
					AssetA: txnbuild.NativeAsset{},
					AssetB: txnbuild.CreditAsset{
						Code:   "USD",
						Issuer: master.Address(),
					},
					Fee: 30,
				},
			},
			Limit: txnbuild.MaxTrustlineLimit,
		},
		&txnbuild.Payment{
			SourceAccount: master.Address(),
			Destination:   shareAccount.GetAccountID(),
			Asset: txnbuild.CreditAsset{
				Code:   "USD",
				Issuer: master.Address(),
			},
			Amount: "1000",
		},
	)

	poolID, err := xdr.NewPoolId(
		xdr.MustNewNativeAsset(),
		xdr.MustNewCreditAsset("USD", master.Address()),
		30,
	)
	tt.NoError(err)
	poolIDHexString := xdr.Hash(poolID).HexString()

	itest.MustSubmitOperations(shareAccount, shareKeys,
		&txnbuild.LiquidityPoolDeposit{
			LiquidityPoolID: [32]byte(poolID),
			MaxAmountA:      "400",
			MaxAmountB:      "777",
			MinPrice:        "0.5",
			MaxPrice:        "2",
		},
	)

	itest.MustSubmitOperations(tradeAccount, tradeKeys,
		&txnbuild.ChangeTrust{
			Line: txnbuild.ChangeTrustAssetWrapper{
				Asset: txnbuild.CreditAsset{
					Code:   "USD",
					Issuer: master.Address(),
				},
			},
			Limit: txnbuild.MaxTrustlineLimit,
		},
		&txnbuild.PathPaymentStrictReceive{
			SendAsset: txnbuild.NativeAsset{},
			DestAsset: txnbuild.CreditAsset{
				Code:   "USD",
				Issuer: master.Address(),
			},
			SendMax:     "1000",
			DestAmount:  "2",
			Destination: tradeKeys.Address(),
		},
	)

	pool, err := itest.Client().LiquidityPoolDetail(auroraclient.LiquidityPoolRequest{
		LiquidityPoolID: poolIDHexString,
	})
	tt.NoError(err)

	txResp := itest.MustSubmitOperations(shareAccount, shareKeys,
		&txnbuild.LiquidityPoolWithdraw{
			LiquidityPoolID: [32]byte(poolID),
			Amount:          pool.TotalShares,
			MinAmountA:      "10",
			MinAmountB:      "20",
		},
	)

	return txResp.Ledger
}

func generatePaymentOps(itest *integration.Test, tt *assert.Assertions) (lastLedger int32) {
	txResp := itest.MustSubmitOperations(itest.MasterAccount(), itest.Master(),
		&txnbuild.Payment{
			Destination: itest.Master().Address(),
			Amount:      "10",
			Asset:       txnbuild.NativeAsset{},
		},
	)

	return txResp.Ledger
}

func initializeDBIntegrationTest(t *testing.T) (itest *integration.Test, reachedLedger int32) {
	config := integration.Config{ProtocolVersion: 18}
	itest = integration.NewTest(t, config)
	tt := assert.New(t)

	generatePaymentOps(itest, tt)
	reachedLedger = generateLiquidityPoolOps(itest, tt)

	root, err := itest.Client().Root()
	tt.NoError(err)
	tt.LessOrEqual(reachedLedger, root.auroraSequence)

	return
}

func TestReingestDB(t *testing.T) {
	itest, reachedLedger := initializeDBIntegrationTest(t)
	tt := assert.New(t)

	// Create a fresh aurora database
	newDB := dbtest.Postgres(t)
	// TODO: Unfortunately aurora's ingestion System leaves open sessions behind,leading to
	//       a "database  is being accessed by other users" error when trying to drop it
	// defer newDB.Close()
	freshauroraPostgresURL := newDB.DSN
	auroraConfig := itest.GetauroraConfig()
	auroraConfig.DatabaseURL = freshauroraPostgresURL
	// Initialize the DB schema
	dbConn, err := db.Open("postgres", freshauroraPostgresURL)
	tt.NoError(err)
	defer dbConn.Close()
	_, err = schema.Migrate(dbConn.DB.DB, schema.MigrateUp, 0)
	tt.NoError(err)

	t.Run("validate parallel range", func(t *testing.T) {
		auroracmd.RootCmd.SetArgs(command(auroraConfig,
			"db",
			"reingest",
			"range",
			"--parallel-workers=2",
			"10",
			"2",
		))

		assert.EqualError(t, auroracmd.RootCmd.Execute(), "Invalid range: {10 2} from > to")
	})

	// cap reachedLedger to the nearest checkpoint ledger because reingest range cannot ingest past the most
	// recent checkpoint ledger when using captive core
	toLedger := uint32(reachedLedger)
	archive, err := historyarchive.Connect(auroraConfig.HistoryArchiveURLs[0], historyarchive.ConnectOptions{
		NetworkPassphrase:   auroraConfig.NetworkPassphrase,
		CheckpointFrequency: auroraConfig.CheckpointFrequency,
	})
	tt.NoError(err)

	// make sure a full checkpoint has elapsed otherwise there will be nothing to reingest
	var latestCheckpoint uint32
	publishedFirstCheckpoint := func() bool {
		has, requestErr := archive.GetRootHAS()
		tt.NoError(requestErr)
		latestCheckpoint = has.CurrentLedger
		return latestCheckpoint > 1
	}
	tt.Eventually(publishedFirstCheckpoint, 10*time.Second, time.Second)

	if toLedger > latestCheckpoint {
		toLedger = latestCheckpoint
	}

	// We just want to test reingestion, so there's no reason for a background
	// aurora to run. Keeping it running will actually cause the Captive Core
	// subprocesses to conflict.
	itest.Stopaurora()

	auroraConfig.CaptiveCoreConfigPath = filepath.Join(
		filepath.Dir(auroraConfig.CaptiveCoreConfigPath),
		"captive-core-reingest-range-integration-tests.cfg",
	)

	auroracmd.RootCmd.SetArgs(command(auroraConfig, "db",
		"reingest",
		"range",
		"--parallel-workers=1",
		"1",
		fmt.Sprintf("%d", toLedger),
	))

	tt.NoError(auroracmd.RootCmd.Execute())
	tt.NoError(auroracmd.RootCmd.Execute(), "Repeat the same reingest range against db, should not have errors.")
}

func command(auroraConfig aurora.Config, args ...string) []string {
	return append([]string{
		"--diamcircle-core-url",
		auroraConfig.diamcircleCoreURL,
		"--history-archive-urls",
		auroraConfig.HistoryArchiveURLs[0],
		"--db-url",
		auroraConfig.DatabaseURL,
		"--diamcircle-core-db-url",
		auroraConfig.diamcircleCoreDatabaseURL,
		"--diamcircle-core-binary-path",
		auroraConfig.CaptiveCoreBinaryPath,
		"--captive-core-config-path",
		auroraConfig.CaptiveCoreConfigPath,
		"--enable-captive-core-ingestion=" + strconv.FormatBool(auroraConfig.EnableCaptiveCoreIngestion),
		"--network-passphrase",
		auroraConfig.NetworkPassphrase,
		// due to ARTIFICIALLY_ACCELERATE_TIME_FOR_TESTING
		"--checkpoint-frequency",
		"8",
	}, args...)
}

func TestFillGaps(t *testing.T) {
	itest, reachedLedger := initializeDBIntegrationTest(t)
	tt := assert.New(t)

	// Create a fresh aurora database
	newDB := dbtest.Postgres(t)
	// TODO: Unfortunately aurora's ingestion System leaves open sessions behind,leading to
	//       a "database  is being accessed by other users" error when trying to drop it
	// defer newDB.Close()
	freshauroraPostgresURL := newDB.DSN
	auroraConfig := itest.GetauroraConfig()
	auroraConfig.DatabaseURL = freshauroraPostgresURL
	// Initialize the DB schema
	dbConn, err := db.Open("postgres", freshauroraPostgresURL)
	defer dbConn.Close()
	_, err = schema.Migrate(dbConn.DB.DB, schema.MigrateUp, 0)
	tt.NoError(err)

	// cap reachedLedger to the nearest checkpoint ledger because reingest range cannot ingest past the most
	// recent checkpoint ledger when using captive core
	toLedger := uint32(reachedLedger)
	archive, err := historyarchive.Connect(auroraConfig.HistoryArchiveURLs[0], historyarchive.ConnectOptions{
		NetworkPassphrase:   auroraConfig.NetworkPassphrase,
		CheckpointFrequency: auroraConfig.CheckpointFrequency,
	})
	tt.NoError(err)

	t.Run("validate parallel range", func(t *testing.T) {
		auroracmd.RootCmd.SetArgs(command(auroraConfig,
			"db",
			"fill-gaps",
			"--parallel-workers=2",
			"10",
			"2",
		))

		assert.EqualError(t, auroracmd.RootCmd.Execute(), "Invalid range: {10 2} from > to")
	})

	// make sure a full checkpoint has elapsed otherwise there will be nothing to reingest
	var latestCheckpoint uint32
	publishedFirstCheckpoint := func() bool {
		has, requestErr := archive.GetRootHAS()
		tt.NoError(requestErr)
		latestCheckpoint = has.CurrentLedger
		return latestCheckpoint > 1
	}
	tt.Eventually(publishedFirstCheckpoint, 10*time.Second, time.Second)

	if toLedger > latestCheckpoint {
		toLedger = latestCheckpoint
	}

	// We just want to test reingestion, so there's no reason for a background
	// aurora to run. Keeping it running will actually cause the Captive Core
	// subprocesses to conflict.
	itest.Stopaurora()

	historyQ := history.Q{dbConn}
	var oldestLedger, latestLedger int64
	tt.NoError(historyQ.ElderLedger(context.Background(), &oldestLedger))
	tt.NoError(historyQ.LatestLedger(context.Background(), &latestLedger))
	tt.NoError(historyQ.DeleteRangeAll(context.Background(), oldestLedger, latestLedger))

	auroraConfig.CaptiveCoreConfigPath = filepath.Join(
		filepath.Dir(auroraConfig.CaptiveCoreConfigPath),
		"captive-core-reingest-range-integration-tests.cfg",
	)
	auroracmd.RootCmd.SetArgs(command(auroraConfig, "db", "fill-gaps", "--parallel-workers=1"))
	tt.NoError(auroracmd.RootCmd.Execute())

	tt.NoError(historyQ.LatestLedger(context.Background(), &latestLedger))
	tt.Equal(int64(0), latestLedger)

	auroracmd.RootCmd.SetArgs(command(auroraConfig, "db", "fill-gaps", "3", "4"))
	tt.NoError(auroracmd.RootCmd.Execute())
	tt.NoError(historyQ.LatestLedger(context.Background(), &latestLedger))
	tt.NoError(historyQ.ElderLedger(context.Background(), &oldestLedger))
	tt.Equal(int64(3), oldestLedger)
	tt.Equal(int64(4), latestLedger)

	auroracmd.RootCmd.SetArgs(command(auroraConfig, "db", "fill-gaps", "6", "7"))
	tt.NoError(auroracmd.RootCmd.Execute())
	tt.NoError(historyQ.LatestLedger(context.Background(), &latestLedger))
	tt.NoError(historyQ.ElderLedger(context.Background(), &oldestLedger))
	tt.Equal(int64(3), oldestLedger)
	tt.Equal(int64(7), latestLedger)
	var gaps []history.LedgerRange
	gaps, err = historyQ.GetLedgerGaps(context.Background())
	tt.NoError(err)
	tt.Equal([]history.LedgerRange{{StartSequence: 5, EndSequence: 5}}, gaps)

	auroracmd.RootCmd.SetArgs(command(auroraConfig, "db", "fill-gaps"))
	tt.NoError(auroracmd.RootCmd.Execute())
	tt.NoError(historyQ.LatestLedger(context.Background(), &latestLedger))
	tt.NoError(historyQ.ElderLedger(context.Background(), &oldestLedger))
	tt.Equal(int64(3), oldestLedger)
	tt.Equal(int64(7), latestLedger)
	gaps, err = historyQ.GetLedgerGaps(context.Background())
	tt.NoError(err)
	tt.Empty(gaps)

	auroracmd.RootCmd.SetArgs(command(auroraConfig, "db", "fill-gaps", "2", "8"))
	tt.NoError(auroracmd.RootCmd.Execute())
	tt.NoError(historyQ.LatestLedger(context.Background(), &latestLedger))
	tt.NoError(historyQ.ElderLedger(context.Background(), &oldestLedger))
	tt.Equal(int64(2), oldestLedger)
	tt.Equal(int64(8), latestLedger)
	gaps, err = historyQ.GetLedgerGaps(context.Background())
	tt.NoError(err)
	tt.Empty(gaps)
}

func TestResumeFromInitializedDB(t *testing.T) {
	itest, reachedLedger := initializeDBIntegrationTest(t)
	tt := assert.New(t)

	// Stop the integration test, and restart it with the same database
	oldDBURL := itest.GetauroraConfig().DatabaseURL
	itestConfig := protocol15Config
	itestConfig.PostgresURL = oldDBURL

	err := itest.Restartaurora()
	tt.NoError(err)

	successfullyResumed := func() bool {
		root, err := itest.Client().Root()
		tt.NoError(err)
		// It must be able to reach the ledger and surpass it
		const ledgersPastStopPoint = 4
		return root.auroraSequence > (reachedLedger + ledgersPastStopPoint)
	}

	tt.Eventually(successfullyResumed, 1*time.Minute, 1*time.Second)
}
