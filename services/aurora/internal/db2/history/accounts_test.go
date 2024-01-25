package history

import (
	"testing"

	"github.com/guregu/null"
	"github.com/diamcircle/go/services/aurora/internal/db2"
	"github.com/diamcircle/go/services/aurora/internal/test"
	"github.com/diamcircle/go/xdr"
	"github.com/stretchr/testify/assert"
)

var (
	inflationDest = "GBUH7T6U36DAVEKECMKN5YEBQYZVRBPNSZAAKBCO6P5HBMDFSQMQL4Z4"
	sponsor       = "GCO26ZSBD63TKYX45H2C7D2WOFWOUSG5BMTNC3BG4QMXM3PAYI6WHKVZ"

	account1 = AccountEntry{
		LastModifiedLedger:   1234,
		AccountID:            "GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB",
		Balance:              20000,
		SequenceNumber:       223456789,
		NumSubEntries:        10,
		InflationDestination: inflationDest,
		Flags:                1,
		HomeDomain:           "diamcircle.org",
		MasterWeight:         1,
		ThresholdLow:         2,
		ThresholdMedium:      3,
		ThresholdHigh:        4,
		BuyingLiabilities:    3,
		SellingLiabilities:   4,
	}

	account2 = AccountEntry{
		LastModifiedLedger:   1235,
		AccountID:            "GCT2NQM5KJJEF55NPMY444C6M6CA7T33HRNCMA6ZFBIIXKNCRO6J25K7",
		Balance:              50000,
		SequenceNumber:       648736,
		NumSubEntries:        10,
		InflationDestination: inflationDest,
		Flags:                2,
		HomeDomain:           "meridian.diamcircle.org",
		MasterWeight:         5,
		ThresholdLow:         6,
		ThresholdMedium:      7,
		ThresholdHigh:        8,
		BuyingLiabilities:    30,
		SellingLiabilities:   40,
		NumSponsored:         1,
		NumSponsoring:        2,
		Sponsor:              null.StringFrom(sponsor),
	}

	account3 = AccountEntry{
		LastModifiedLedger:   1235,
		AccountID:            "GDPGOMFSP4IF7A4P7UBKA4UC4QTRLEHGBD6IMDIS3W3KBDNBFAQ7FXDY",
		Balance:              50000,
		SequenceNumber:       648736,
		NumSubEntries:        10,
		InflationDestination: inflationDest,
		Flags:                2,
		MasterWeight:         5,
		ThresholdLow:         6,
		ThresholdMedium:      7,
		ThresholdHigh:        8,
		BuyingLiabilities:    30,
		SellingLiabilities:   40,
	}
)

func TestInsertAccount(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}

	err := q.UpsertAccounts(tt.Ctx, []AccountEntry{account1, account2})
	assert.NoError(t, err)

	accounts, err := q.GetAccountsByIDs(tt.Ctx, []string{
		"GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB",
		"GCT2NQM5KJJEF55NPMY444C6M6CA7T33HRNCMA6ZFBIIXKNCRO6J25K7",
	})
	assert.NoError(t, err)
	assert.Len(t, accounts, 2)

	assert.Equal(t, "GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB", accounts[0].AccountID)
	assert.Equal(t, int64(20000), accounts[0].Balance)
	assert.Equal(t, int64(223456789), accounts[0].SequenceNumber)
	assert.Equal(t, uint32(10), accounts[0].NumSubEntries)
	assert.Equal(t, "GBUH7T6U36DAVEKECMKN5YEBQYZVRBPNSZAAKBCO6P5HBMDFSQMQL4Z4", accounts[0].InflationDestination)
	assert.Equal(t, uint32(1), accounts[0].Flags)
	assert.Equal(t, "diamcircle.org", accounts[0].HomeDomain)
	assert.Equal(t, byte(1), accounts[0].MasterWeight)
	assert.Equal(t, byte(2), accounts[0].ThresholdLow)
	assert.Equal(t, byte(3), accounts[0].ThresholdMedium)
	assert.Equal(t, byte(4), accounts[0].ThresholdHigh)
	assert.Equal(t, int64(3), accounts[0].BuyingLiabilities)
	assert.Equal(t, int64(4), accounts[0].SellingLiabilities)
	assert.Equal(t, uint32(0), accounts[0].NumSponsored)
	assert.Equal(t, uint32(0), accounts[0].NumSponsoring)
	assert.Equal(t, null.String{}, accounts[0].Sponsor)

	assert.Equal(t, uint32(1), accounts[1].NumSponsored)
	assert.Equal(t, uint32(2), accounts[1].NumSponsoring)
	assert.Equal(t, null.StringFrom(sponsor), accounts[1].Sponsor)
}

func TestUpsertAccount(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}

	ledgerEntries := []AccountEntry{account1, account2}
	err := q.UpsertAccounts(tt.Ctx, ledgerEntries)
	assert.NoError(t, err)

	modifiedAccount := AccountEntry{
		LastModifiedLedger:   1234,
		AccountID:            "GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB",
		Balance:              32847893,
		SequenceNumber:       223456789,
		NumSubEntries:        10,
		InflationDestination: inflationDest,
		Flags:                1,
		HomeDomain:           "diamcircle.org",
		MasterWeight:         1,
		ThresholdLow:         2,
		ThresholdMedium:      3,
		ThresholdHigh:        4,
		BuyingLiabilities:    3,
		SellingLiabilities:   4,
	}

	err = q.UpsertAccounts(tt.Ctx, []AccountEntry{modifiedAccount})
	assert.NoError(t, err)

	keys := []string{
		account1.AccountID,
		account2.AccountID,
	}
	accounts, err := q.GetAccountsByIDs(tt.Ctx, keys)
	assert.NoError(t, err)
	assert.Len(t, accounts, 2)

	assert.Equal(t, uint32(1), accounts[0].NumSponsored)
	assert.Equal(t, uint32(2), accounts[0].NumSponsoring)
	assert.Equal(t, null.StringFrom(sponsor), accounts[0].Sponsor)

	assert.Equal(t, uint32(0), accounts[1].NumSponsored)
	assert.Equal(t, uint32(0), accounts[1].NumSponsoring)
	assert.Equal(t, null.String{}, accounts[1].Sponsor)

	accounts, err = q.GetAccountsByIDs(tt.Ctx, []string{account1.AccountID})
	assert.NoError(t, err)
	assert.Len(t, accounts, 1)

	assert.Equal(t, modifiedAccount, accounts[0])
	assert.Equal(t, uint32(1234), accounts[0].LastModifiedLedger)

	accounts, err = q.GetAccountsByIDs(tt.Ctx, []string{account2.AccountID})
	assert.NoError(t, err)
	assert.Len(t, accounts, 1)

	assert.Equal(t, account2, accounts[0])
	assert.Equal(t, uint32(1235), accounts[0].LastModifiedLedger)
}

func TestRemoveAccount(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}

	err := q.UpsertAccounts(tt.Ctx, []AccountEntry{account1})
	assert.NoError(t, err)

	var rows int64
	rows, err = q.RemoveAccounts(tt.Ctx, []string{"GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rows)

	accounts, err := q.GetAccountsByIDs(tt.Ctx, []string{"GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"})
	assert.NoError(t, err)
	assert.Len(t, accounts, 0)

	// Doesn't exist anymore
	rows, err = q.RemoveAccounts(tt.Ctx, []string{"GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"})
	assert.NoError(t, err)
	assert.Equal(t, int64(0), rows)
}

func TestAccountsForAsset(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}

	eurTL := eurTrustLine
	usdTL := usdTrustLine
	psTL := poolShareTrustLine

	eurTL.AccountID = account1.AccountID
	usdTL.AccountID = account2.AccountID
	psTL.AccountID = account1.AccountID

	err := q.UpsertAccounts(tt.Ctx, []AccountEntry{account1, account2})
	assert.NoError(t, err)

	tt.Assert.NoError(q.UpsertTrustLines(tt.Ctx, []TrustLine{
		eurTL,
		usdTL,
		psTL,
	}))

	pq := db2.PageQuery{
		Order:  db2.OrderAscending,
		Limit:  db2.DefaultPageSize,
		Cursor: "",
	}

	accounts, err := q.AccountsForAsset(
		tt.Ctx,
		xdr.MustNewCreditAsset(eurTL.AssetCode, eurTL.AssetIssuer),
		pq,
	)
	assert.NoError(t, err)
	tt.Assert.Len(accounts, 1)
	tt.Assert.Equal(account1.AccountID, accounts[0].AccountID)

	accounts, err = q.AccountsForAsset(
		tt.Ctx,
		xdr.MustNewCreditAsset(usdTL.AssetCode, usdTL.AssetIssuer),
		pq,
	)
	assert.NoError(t, err)
	tt.Assert.Len(accounts, 1)
	tt.Assert.Equal(account2.AccountID, accounts[0].AccountID)

	pq.Cursor = account2.AccountID
	accounts, err = q.AccountsForAsset(
		tt.Ctx,
		xdr.MustNewCreditAsset(usdTL.AssetCode, usdTL.AssetIssuer),
		pq,
	)
	assert.NoError(t, err)
	tt.Assert.Len(accounts, 0)

	pq = db2.PageQuery{
		Order:  db2.OrderDescending,
		Limit:  db2.DefaultPageSize,
		Cursor: "",
	}

	accounts, err = q.AccountsForAsset(
		tt.Ctx,
		xdr.MustNewCreditAsset(usdTL.AssetCode, eurTL.AssetIssuer),
		pq,
	)
	assert.NoError(t, err)
	tt.Assert.Len(accounts, 1)
}

func TestAccountsForLiquidityPool(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}

	eurTL := eurTrustLine
	psTL := poolShareTrustLine

	eurTL.AccountID = account1.AccountID
	psTL.AccountID = account1.AccountID

	err := q.UpsertAccounts(tt.Ctx, []AccountEntry{account1})
	assert.NoError(t, err)

	tt.Assert.NoError(q.UpsertTrustLines(tt.Ctx, []TrustLine{
		eurTL,
		psTL,
	}))

	pq := db2.PageQuery{
		Order:  db2.OrderAscending,
		Limit:  db2.DefaultPageSize,
		Cursor: "",
	}

	accounts, err := q.AccountsForLiquidityPool(
		tt.Ctx,
		psTL.LiquidityPoolID,
		pq,
	)
	assert.NoError(t, err)
	tt.Assert.Len(accounts, 1)
	tt.Assert.Equal(account1.AccountID, accounts[0].AccountID)

	pq.Cursor = account1.AccountID
	accounts, err = q.AccountsForLiquidityPool(
		tt.Ctx,
		psTL.LiquidityPoolID,
		pq,
	)
	assert.NoError(t, err)
	tt.Assert.Len(accounts, 0)
}

func TestAccountsForSponsor(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}

	eurTL := eurTrustLine
	usdTL := usdTrustLine

	eurTL.AccountID = account1.AccountID
	usdTL.AccountID = account2.AccountID

	err := q.UpsertAccounts(tt.Ctx, []AccountEntry{account1, account2, account3})
	assert.NoError(t, err)

	tt.Assert.NoError(q.UpsertTrustLines(tt.Ctx, []TrustLine{eurTL, usdTL}))

	_, err = q.CreateAccountSigner(tt.Ctx, account1.AccountID, account1.AccountID, 1, nil)
	tt.Assert.NoError(err)
	_, err = q.CreateAccountSigner(tt.Ctx, account2.AccountID, account2.AccountID, 1, nil)
	tt.Assert.NoError(err)
	_, err = q.CreateAccountSigner(tt.Ctx, account3.AccountID, account3.AccountID, 1, nil)
	tt.Assert.NoError(err)
	_, err = q.CreateAccountSigner(tt.Ctx, account1.AccountID, account3.AccountID, 1, nil)
	tt.Assert.NoError(err)
	_, err = q.CreateAccountSigner(tt.Ctx, account2.AccountID, account3.AccountID, 1, nil)
	tt.Assert.NoError(err)

	pq := db2.PageQuery{
		Order:  db2.OrderAscending,
		Limit:  db2.DefaultPageSize,
		Cursor: "",
	}

	accounts, err := q.AccountsForSponsor(tt.Ctx, sponsor, pq)
	assert.NoError(t, err)
	tt.Assert.Len(accounts, 2)
	tt.Assert.Equal(account1.AccountID, accounts[0].AccountID)
	tt.Assert.Equal(account2.AccountID, accounts[1].AccountID)
}

func TestAccountEntriesForSigner(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}

	eurTL := eurTrustLine
	usdTL := usdTrustLine

	eurTL.AccountID = account1.AccountID
	usdTL.AccountID = account2.AccountID

	err := q.UpsertAccounts(tt.Ctx, []AccountEntry{account1, account2, account3})
	assert.NoError(t, err)

	tt.Assert.NoError(q.UpsertTrustLines(tt.Ctx, []TrustLine{eurTL, usdTL}))

	_, err = q.CreateAccountSigner(tt.Ctx, account1.AccountID, account1.AccountID, 1, nil)
	tt.Assert.NoError(err)
	_, err = q.CreateAccountSigner(tt.Ctx, account2.AccountID, account2.AccountID, 1, nil)
	tt.Assert.NoError(err)
	_, err = q.CreateAccountSigner(tt.Ctx, account3.AccountID, account3.AccountID, 1, nil)
	tt.Assert.NoError(err)
	_, err = q.CreateAccountSigner(tt.Ctx, account1.AccountID, account3.AccountID, 1, nil)
	tt.Assert.NoError(err)
	_, err = q.CreateAccountSigner(tt.Ctx, account2.AccountID, account3.AccountID, 1, nil)
	tt.Assert.NoError(err)

	pq := db2.PageQuery{
		Order:  db2.OrderAscending,
		Limit:  db2.DefaultPageSize,
		Cursor: "",
	}

	accounts, err := q.AccountEntriesForSigner(tt.Ctx, account1.AccountID, pq)
	assert.NoError(t, err)
	tt.Assert.Len(accounts, 1)
	tt.Assert.Equal(account1.AccountID, accounts[0].AccountID)

	accounts, err = q.AccountEntriesForSigner(tt.Ctx, account2.AccountID, pq)
	assert.NoError(t, err)
	tt.Assert.Len(accounts, 1)
	tt.Assert.Equal(account2.AccountID, accounts[0].AccountID)

	want := map[string]bool{
		account1.AccountID: true,
		account2.AccountID: true,
		account3.AccountID: true,
	}

	accounts, err = q.AccountEntriesForSigner(tt.Ctx, account3.AccountID, pq)
	assert.NoError(t, err)
	tt.Assert.Len(accounts, 3)

	for _, account := range accounts {
		tt.Assert.True(want[account.AccountID])
		delete(want, account.AccountID)
	}

	tt.Assert.Len(want, 0)

	pq.Cursor = accounts[len(accounts)-1].AccountID
	accounts, err = q.AccountEntriesForSigner(tt.Ctx, account3.AccountID, pq)
	assert.NoError(t, err)
	tt.Assert.Len(accounts, 0)

	pq.Order = "desc"
	accounts, err = q.AccountEntriesForSigner(tt.Ctx, account3.AccountID, pq)
	assert.NoError(t, err)
	tt.Assert.Len(accounts, 2)

	pq = db2.PageQuery{
		Order:  db2.OrderDescending,
		Limit:  db2.DefaultPageSize,
		Cursor: "",
	}

	accounts, err = q.AccountEntriesForSigner(tt.Ctx, account1.AccountID, pq)
	assert.NoError(t, err)
	tt.Assert.Len(accounts, 1)
}

func TestGetAccountByID(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}

	err := q.UpsertAccounts(tt.Ctx, []AccountEntry{account1})
	assert.NoError(t, err)

	resultAccount, err := q.GetAccountByID(tt.Ctx, account1.AccountID)
	assert.NoError(t, err)

	assert.Equal(t, "GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB", resultAccount.AccountID)
	assert.Equal(t, int64(20000), resultAccount.Balance)
	assert.Equal(t, int64(223456789), resultAccount.SequenceNumber)
	assert.Equal(t, uint32(10), resultAccount.NumSubEntries)
	assert.Equal(t, "GBUH7T6U36DAVEKECMKN5YEBQYZVRBPNSZAAKBCO6P5HBMDFSQMQL4Z4", resultAccount.InflationDestination)
	assert.Equal(t, uint32(1), resultAccount.Flags)
	assert.Equal(t, "diamcircle.org", resultAccount.HomeDomain)
	assert.Equal(t, byte(1), resultAccount.MasterWeight)
	assert.Equal(t, byte(2), resultAccount.ThresholdLow)
	assert.Equal(t, byte(3), resultAccount.ThresholdMedium)
	assert.Equal(t, byte(4), resultAccount.ThresholdHigh)
	assert.Equal(t, int64(3), resultAccount.BuyingLiabilities)
	assert.Equal(t, int64(4), resultAccount.SellingLiabilities)
}
