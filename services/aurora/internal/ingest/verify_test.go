package ingest

import (
	"io"
	"math/rand"
	"regexp"
	"testing"

	"github.com/diamcircle/go/gxdr"
	"github.com/diamcircle/go/historyarchive"
	"github.com/diamcircle/go/ingest"
	"github.com/diamcircle/go/randxdr"
	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/services/aurora/internal/test"
	"github.com/diamcircle/go/support/db"
	"github.com/diamcircle/go/xdr"
)

func genAccount(tt *test.T, gen randxdr.Generator) xdr.LedgerEntryChange {
	change := xdr.LedgerEntryChange{}
	shape := &gxdr.LedgerEntryChange{}
	numSigners := uint32(rand.Int31n(xdr.MaxSigners))
	gen.Next(
		shape,
		[]randxdr.Preset{
			{randxdr.FieldEquals("type"), randxdr.SetU32(gxdr.LEDGER_ENTRY_CREATED.GetU32())},
			{randxdr.FieldEquals("created.data.type"), randxdr.SetU32(gxdr.ACCOUNT.GetU32())},
			{randxdr.FieldEquals("created.lastModifiedLedgerSeq"), randxdr.SetPositiveNum32},
			{randxdr.FieldEquals("created.data.account.seqNum"), randxdr.SetPositiveNum64},
			{randxdr.FieldEquals("created.data.account.numSubEntries"), randxdr.SetPositiveNum32},
			{randxdr.FieldEquals("created.data.account.balance"), randxdr.SetPositiveNum64},
			{randxdr.FieldEquals("created.data.account.homeDomain"), randxdr.SetPrintableASCII},
			{randxdr.FieldEquals("created.data.account.flags"), randxdr.SetPositiveNum32},
			{randxdr.FieldEquals("created.data.account.signers"), randxdr.SetVecLen(numSigners)},
			{randxdr.FieldMatches(regexp.MustCompile("created\\.data\\.account\\.signers.*weight")), randxdr.SetPositiveNum32},
			{randxdr.FieldEquals("created.data.account.ext.v1.liabilities.selling"), randxdr.SetPositiveNum64},
			{randxdr.FieldEquals("created.data.account.ext.v1.liabilities.buying"), randxdr.SetPositiveNum64},
			{randxdr.FieldEquals("created.data.account.ext.v1.ext.v2.numSponsoring"), randxdr.SetPositiveNum32},
			{randxdr.FieldEquals("created.data.account.ext.v1.ext.v2.numSponsored"), randxdr.SetPositiveNum32},
			{randxdr.FieldEquals("created.data.account.ext.v1.ext.v2.signerSponsoringIDs"), randxdr.SetVecLen(numSigners)},
		},
	)

	tt.Assert.NoError(gxdr.Convert(shape, &change))
	return change
}

func genAccountData(tt *test.T, gen randxdr.Generator) xdr.LedgerEntryChange {
	change := xdr.LedgerEntryChange{}
	shape := &gxdr.LedgerEntryChange{}
	gen.Next(
		shape,
		[]randxdr.Preset{
			{randxdr.FieldEquals("type"), randxdr.SetU32(gxdr.LEDGER_ENTRY_CREATED.GetU32())},
			{randxdr.FieldEquals("created.data.type"), randxdr.SetU32(gxdr.DATA.GetU32())},
			{randxdr.FieldEquals("created.lastModifiedLedgerSeq"), randxdr.SetPositiveNum32},
			{randxdr.FieldEquals("created.data.data.dataName"), randxdr.SetPrintableASCII},
		},
	)

	tt.Assert.NoError(gxdr.Convert(shape, &change))
	return change
}

func genOffer(tt *test.T, gen randxdr.Generator) xdr.LedgerEntryChange {
	change := xdr.LedgerEntryChange{}
	shape := &gxdr.LedgerEntryChange{}
	gen.Next(
		shape,
		[]randxdr.Preset{
			{randxdr.FieldEquals("type"), randxdr.SetU32(gxdr.LEDGER_ENTRY_CREATED.GetU32())},
			{randxdr.FieldEquals("created.data.type"), randxdr.SetU32(gxdr.OFFER.GetU32())},
			{randxdr.FieldEquals("created.lastModifiedLedgerSeq"), randxdr.SetPositiveNum32},
			{randxdr.FieldEquals("created.data.offer.amount"), randxdr.SetPositiveNum64},
			{randxdr.FieldEquals("created.data.offer.price.n"), randxdr.SetPositiveNum32},
			{randxdr.FieldEquals("created.data.offer.price.d"), randxdr.SetPositiveNum32},
		},
	)
	tt.Assert.NoError(gxdr.Convert(shape, &change))
	return change
}

func genLiquidityPool(tt *test.T, gen randxdr.Generator) xdr.LedgerEntryChange {
	change := xdr.LedgerEntryChange{}
	shape := &gxdr.LedgerEntryChange{}
	gen.Next(
		shape,
		[]randxdr.Preset{
			{randxdr.FieldEquals("type"), randxdr.SetU32(gxdr.LEDGER_ENTRY_CREATED.GetU32())},
			// liquidity pools cannot be sponsored
			{randxdr.FieldEquals("created.ext.v1.sponsoringID"), randxdr.SetPtr(false)},
			{randxdr.FieldEquals("created.lastModifiedLedgerSeq"), randxdr.SetPositiveNum32},
			{randxdr.FieldEquals("created.data.type"), randxdr.SetU32(gxdr.LIQUIDITY_POOL.GetU32())},
			{randxdr.FieldEquals("created.data.liquidityPool.body.constantProduct.params.fee"), randxdr.SetPositiveNum32},
			{randxdr.FieldEquals("created.data.liquidityPool.body.constantProduct.params.assetA.alphaNum4.assetCode"), randxdr.SetAssetCode},
			{randxdr.FieldEquals("created.data.liquidityPool.body.constantProduct.params.assetA.alphaNum12.assetCode"), randxdr.SetAssetCode},
			{randxdr.FieldEquals("created.data.liquidityPool.body.constantProduct.params.assetB.alphaNum4.assetCode"), randxdr.SetAssetCode},
			{randxdr.FieldEquals("created.data.liquidityPool.body.constantProduct.params.assetB.alphaNum12.assetCode"), randxdr.SetAssetCode},
			{randxdr.FieldEquals("created.data.liquidityPool.body.constantProduct.reserveA"), randxdr.SetPositiveNum64},
			{randxdr.FieldEquals("created.data.liquidityPool.body.constantProduct.reserveB"), randxdr.SetPositiveNum64},
			{randxdr.FieldEquals("created.data.liquidityPool.body.constantProduct.totalPoolShares"), randxdr.SetPositiveNum64},
			{randxdr.FieldEquals("created.data.liquidityPool.body.constantProduct.poolSharesTrustLineCount"), randxdr.SetPositiveNum64},
		},
	)
	tt.Assert.NoError(gxdr.Convert(shape, &change))
	return change
}

func genTrustLine(tt *test.T, gen randxdr.Generator) xdr.LedgerEntryChange {
	change := xdr.LedgerEntryChange{}
	shape := &gxdr.LedgerEntryChange{}
	gen.Next(
		shape,
		[]randxdr.Preset{
			{randxdr.FieldEquals("type"), randxdr.SetU32(gxdr.LEDGER_ENTRY_CREATED.GetU32())},
			{randxdr.FieldEquals("created.lastModifiedLedgerSeq"), randxdr.SetPositiveNum32},
			{randxdr.FieldEquals("created.data.type"), randxdr.SetU32(gxdr.TRUSTLINE.GetU32())},
			{randxdr.FieldEquals("created.data.trustLine.flags"), randxdr.SetPositiveNum32},
			{randxdr.FieldEquals("created.data.trustLine.asset.alphaNum4.assetCode"), randxdr.SetAssetCode},
			{randxdr.FieldEquals("created.data.trustLine.asset.alphaNum12.assetCode"), randxdr.SetAssetCode},
			{randxdr.FieldEquals("created.data.trustLine.balance"), randxdr.SetPositiveNum64},
			{randxdr.FieldEquals("created.data.trustLine.limit"), randxdr.SetPositiveNum64},
			{randxdr.FieldEquals("created.data.trustLine.ext.v1.liabilities.selling"), randxdr.SetPositiveNum64},
			{randxdr.FieldEquals("created.data.trustLine.ext.v1.liabilities.buying"), randxdr.SetPositiveNum64},
		},
	)
	tt.Assert.NoError(gxdr.Convert(shape, &change))
	return change
}

func genClaimableBalance(tt *test.T, gen randxdr.Generator) xdr.LedgerEntryChange {
	change := xdr.LedgerEntryChange{}
	shape := &gxdr.LedgerEntryChange{}
	gen.Next(
		shape,
		[]randxdr.Preset{
			{randxdr.FieldEquals("type"), randxdr.SetU32(gxdr.LEDGER_ENTRY_CREATED.GetU32())},
			{randxdr.FieldEquals("created.lastModifiedLedgerSeq"), randxdr.SetPositiveNum32},
			{randxdr.FieldEquals("created.data.type"), randxdr.SetU32(gxdr.CLAIMABLE_BALANCE.GetU32())},
			{randxdr.FieldEquals("created.data.claimableBalance.ext.v1.flags"), randxdr.SetPositiveNum32},
			{
				randxdr.And(
					randxdr.IsPtr,
					randxdr.FieldMatches(regexp.MustCompile("created\\.data\\.claimableBalance\\.claimants.*notPredicate")),
				),
				randxdr.SetPtr(true),
			},
			{randxdr.FieldEquals("created.data.claimableBalance.amount"), randxdr.SetPositiveNum64},
			{randxdr.FieldEquals("created.data.claimableBalance.asset.alphaNum4.assetCode"), randxdr.SetAssetCode},
			{randxdr.FieldEquals("created.data.claimableBalance.asset.alphaNum12.assetCode"), randxdr.SetAssetCode},
		},
	)
	tt.Assert.NoError(gxdr.Convert(shape, &change))
	return change
}

func TestStateVerifier(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &history.Q{&db.Session{DB: tt.AuroraDB}}

	checkpointLedger := uint32(63)
	changeProcessor := buildChangeProcessor(q, &ingest.StatsChangeProcessor{}, ledgerSource, checkpointLedger)
	mockChangeReader := &ingest.MockChangeReader{}

	gen := randxdr.NewGenerator()
	var changes []xdr.LedgerEntryChange
	for i := 0; i < 100; i++ {
		changes = append(changes,
			genLiquidityPool(tt, gen),
			genClaimableBalance(tt, gen),
			genOffer(tt, gen),
			genTrustLine(tt, gen),
			genAccount(tt, gen),
			genAccountData(tt, gen),
		)
	}
	for _, change := range ingest.GetChangesFromLedgerEntryChanges(changes) {
		mockChangeReader.On("Read").Return(change, nil).Once()
		tt.Assert.NoError(changeProcessor.ProcessChange(tt.Ctx, change))
	}
	tt.Assert.NoError(changeProcessor.Commit(tt.Ctx))

	q.UpdateLastLedgerIngest(tt.Ctx, checkpointLedger)

	mockChangeReader.On("Read").Return(ingest.Change{}, io.EOF).Twice()
	mockChangeReader.On("Close").Return(nil).Once()

	mockHistoryAdapter := &mockHistoryArchiveAdapter{}
	mockHistoryAdapter.On("GetState", tt.Ctx, uint32(checkpointLedger)).Return(mockChangeReader, nil).Once()

	sys := &system{
		ctx:               tt.Ctx,
		historyQ:          q,
		historyAdapter:    mockHistoryAdapter,
		checkpointManager: historyarchive.NewCheckpointManager(64),
	}
	sys.initMetrics()

	tt.Assert.NoError(sys.verifyState(false))
	mockChangeReader.AssertExpectations(t)
	mockHistoryAdapter.AssertExpectations(t)
}
