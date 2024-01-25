package history

import (
	"testing"

	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"

	"github.com/diamcircle/go/services/aurora/internal/test"
	"github.com/diamcircle/go/xdr"
)

var (
	trustLineIssuer = "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H"

	eurTrustLine = TrustLine{
		AccountID:          account1.AccountID,
		AssetType:          xdr.AssetTypeAssetTypeCreditAlphanum4,
		AssetIssuer:        trustLineIssuer,
		AssetCode:          "EUR",
		Balance:            30000,
		LedgerKey:          "abcdef",
		Limit:              223456789,
		LiquidityPoolID:    "",
		BuyingLiabilities:  3,
		SellingLiabilities: 4,
		Flags:              1,
		LastModifiedLedger: 1234,
		Sponsor:            null.StringFrom(sponsor),
	}

	usdTrustLine = TrustLine{
		AccountID:          "GCYVFGI3SEQJGBNQQG7YCMFWEYOHK3XPVOVPA6C566PXWN4SN7LILZSM",
		AssetType:          xdr.AssetTypeAssetTypeCreditAlphanum4,
		AssetIssuer:        trustLineIssuer,
		AssetCode:          "USD",
		Balance:            10000,
		LedgerKey:          "jhkli",
		Limit:              123456789,
		LiquidityPoolID:    "",
		BuyingLiabilities:  1,
		SellingLiabilities: 2,
		Flags:              0,
		LastModifiedLedger: 1235,
		Sponsor:            null.String{},
	}

	usdTrustLine2 = TrustLine{
		AccountID:          "GBYSBDAJZMHL5AMD7QXQ3JEP3Q4GLKADWIJURAAHQALNAWD6Z5XF2RAC",
		AssetType:          xdr.AssetTypeAssetTypeCreditAlphanum4,
		AssetIssuer:        trustLineIssuer,
		AssetCode:          "USD",
		Balance:            10000,
		LedgerKey:          "lkjpoi",
		Limit:              123456789,
		LiquidityPoolID:    "",
		BuyingLiabilities:  1,
		SellingLiabilities: 2,
		Flags:              0,
		LastModifiedLedger: 1234,
		Sponsor:            null.String{},
	}

	poolShareTrustLine = TrustLine{
		AccountID:          "GBB4JST32UWKOLGYYSCEYBHBCOFL2TGBHDVOMZP462ET4ZRD4ULA7S2L",
		AssetType:          xdr.AssetTypeAssetTypePoolShare,
		Balance:            976,
		LedgerKey:          "mlmn908",
		Limit:              87654,
		LiquidityPoolID:    "mpolnbv",
		Flags:              1,
		LastModifiedLedger: 1235,
		Sponsor:            null.String{},
	}
)

func TestIsAuthorized(t *testing.T) {
	tt := assert.New(t)
	tl := TrustLine{
		Flags: 1,
	}
	tt.True(tl.IsAuthorized())

	tl = TrustLine{
		Flags: 0,
	}
	tt.False(tl.IsAuthorized())
}
func TestInsertTrustLine(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}

	tt.Assert.NoError(q.UpsertTrustLines(tt.Ctx, []TrustLine{eurTrustLine, usdTrustLine}))

	lines, err := q.GetTrustLinesByKeys(tt.Ctx, []string{eurTrustLine.LedgerKey, usdTrustLine.LedgerKey})
	tt.Assert.NoError(err)
	tt.Assert.Len(lines, 2)

	tt.Assert.Equal(null.StringFrom(sponsor), lines[0].Sponsor)
	tt.Assert.Equal([]TrustLine{eurTrustLine, usdTrustLine}, lines)
}

func TestUpdateTrustLine(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}

	tt.Assert.NoError(q.UpsertTrustLines(tt.Ctx, []TrustLine{eurTrustLine}))

	lines, err := q.GetTrustLinesByKeys(tt.Ctx, []string{eurTrustLine.LedgerKey})
	assert.NoError(t, err)
	assert.Len(t, lines, 1)
	assert.Equal(t, eurTrustLine, lines[0])

	modifiedTrustLine := eurTrustLine
	modifiedTrustLine.Balance = 30000
	modifiedTrustLine.Sponsor = null.String{}

	tt.Assert.NoError(q.UpsertTrustLines(tt.Ctx, []TrustLine{modifiedTrustLine}))
	lines, err = q.GetTrustLinesByKeys(tt.Ctx, []string{eurTrustLine.LedgerKey})
	assert.NoError(t, err)
	assert.Len(t, lines, 1)
	assert.Equal(t, modifiedTrustLine, lines[0])
}

func TestUpsertTrustLines(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}

	// Upserting nothing is no op
	err := q.UpsertTrustLines(tt.Ctx, []TrustLine{})
	assert.NoError(t, err)

	err = q.UpsertTrustLines(tt.Ctx, []TrustLine{eurTrustLine, usdTrustLine})
	assert.NoError(t, err)

	lines, err := q.GetTrustLinesByKeys(tt.Ctx, []string{eurTrustLine.LedgerKey})
	assert.NoError(t, err)
	assert.Len(t, lines, 1)

	lines, err = q.GetTrustLinesByKeys(tt.Ctx, []string{usdTrustLine.LedgerKey})
	assert.NoError(t, err)
	assert.Len(t, lines, 1)

	modifiedTrustLine := eurTrustLine
	modifiedTrustLine.Balance = 30000

	err = q.UpsertTrustLines(tt.Ctx, []TrustLine{modifiedTrustLine})
	assert.NoError(t, err)

	lines, err = q.GetTrustLinesByKeys(tt.Ctx, []string{eurTrustLine.LedgerKey})
	assert.NoError(t, err)
	assert.Len(t, lines, 1)

	assert.Equal(t, modifiedTrustLine, lines[0])
	assert.Equal(t, uint32(1234), lines[0].LastModifiedLedger)

	lines, err = q.GetTrustLinesByKeys(tt.Ctx, []string{usdTrustLine.LedgerKey})
	assert.NoError(t, err)
	assert.Len(t, lines, 1)

	assert.Equal(t, usdTrustLine, lines[0])
	assert.Equal(t, uint32(1235), lines[0].LastModifiedLedger)
}

func TestRemoveTrustLine(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}

	err := q.UpsertTrustLines(tt.Ctx, []TrustLine{eurTrustLine})
	assert.NoError(t, err)

	rows, err := q.RemoveTrustLines(tt.Ctx, []string{eurTrustLine.LedgerKey})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rows)

	lines, err := q.GetTrustLinesByKeys(tt.Ctx, []string{eurTrustLine.LedgerKey})
	assert.NoError(t, err)
	assert.Len(t, lines, 0)

	// Doesn't exist anymore
	rows, err = q.RemoveTrustLines(tt.Ctx, []string{eurTrustLine.LedgerKey})
	assert.NoError(t, err)
	assert.Equal(t, int64(0), rows)
}

func TestGetSortedTrustLinesByAccountsID(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}

	err := q.UpsertTrustLines(tt.Ctx, []TrustLine{eurTrustLine, usdTrustLine, usdTrustLine2, poolShareTrustLine})
	assert.NoError(t, err)

	ids := []string{
		eurTrustLine.AccountID,
		usdTrustLine.AccountID,
		poolShareTrustLine.AccountID,
	}

	records, err := q.GetSortedTrustLinesByAccountIDs(tt.Ctx, ids)
	tt.Assert.NoError(err)
	tt.Assert.Len(records, 3)

	m := map[string]TrustLine{
		eurTrustLine.AccountID:       eurTrustLine,
		usdTrustLine.AccountID:       usdTrustLine,
		poolShareTrustLine.AccountID: poolShareTrustLine,
	}

	tt.Assert.Equal(poolShareTrustLine, records[0])
	delete(m, poolShareTrustLine.AccountID)

	lastAssetCode := ""
	lastIssuer := records[1].AssetIssuer
	for _, record := range records[1:] {
		tt.Assert.LessOrEqual(lastAssetCode, record.AssetCode)
		lastAssetCode = record.AssetCode
		tt.Assert.LessOrEqual(lastIssuer, record.AssetIssuer)
		lastIssuer = record.AssetIssuer
		xtl, ok := m[record.AccountID]
		tt.Assert.True(ok)
		tt.Assert.Equal(record, xtl)
		delete(m, record.AccountID)
	}

	tt.Assert.Len(m, 0)
}

func TestGetTrustLinesByAccountID(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}

	records, err := q.GetSortedTrustLinesByAccountID(tt.Ctx, eurTrustLine.AccountID)
	tt.Assert.NoError(err)
	tt.Assert.Empty(records)

	err = q.UpsertTrustLines(tt.Ctx, []TrustLine{eurTrustLine})
	tt.Assert.NoError(err)

	records, err = q.GetSortedTrustLinesByAccountID(tt.Ctx, eurTrustLine.AccountID)
	tt.Assert.NoError(err)

	tt.Assert.Equal(eurTrustLine, records[0])

}

func TestAssetsForAddress(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}

	ledgerEntries := []AccountEntry{account1}

	err := q.UpsertAccounts(tt.Ctx, ledgerEntries)
	assert.NoError(t, err)

	err = q.UpsertTrustLines(tt.Ctx, []TrustLine{eurTrustLine})
	tt.Assert.NoError(err)

	records, err := q.GetSortedTrustLinesByAccountID(tt.Ctx, eurTrustLine.AccountID)
	tt.Assert.NoError(err)
	tt.Assert.Equal(eurTrustLine, records[0])
}
