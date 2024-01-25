package resourceadapter

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/guregu/null"

	"go/amount"
	. "go/protocols/aurora"
	protocol "go/protocols/aurora"
	"go/services/aurora/internal/assets"
	"go/services/aurora/internal/db2/history"
	"go/support/test"
	"go/xdr"

	"github.com/stretchr/testify/assert"
)

var (
	accountID = xdr.MustAddress("GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB")

	data = []history.Data{
		{
			AccountID:          accountID.Address(),
			Name:               "test",
			Value:              []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			LastModifiedLedger: 1,
			Sponsor:            null.StringFrom("GBUH7T6U36DAVEKECMKN5YEBQYZVRBPNSZAAKBCO6P5HBMDFSQMQL4Z4"),
		},
		{
			AccountID:          accountID.Address(),
			Name:               "test2",
			Value:              []byte{10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
			LastModifiedLedger: 2,
		},
	}

	inflationDest = xdr.MustAddress("GBUH7T6U36DAVEKECMKN5YEBQYZVRBPNSZAAKBCO6P5HBMDFSQMQL4Z4")

	account = history.AccountEntry{
		AccountID:            accountID.Address(),
		Balance:              20000,
		SequenceNumber:       223456789,
		NumSubEntries:        10,
		InflationDestination: inflationDest.Address(),
		Flags:                0b1001, // required and clawback
		HomeDomain:           "diamcircle.org",
		ThresholdLow:         1,
		ThresholdMedium:      2,
		ThresholdHigh:        3,
		SellingLiabilities:   4,
		BuyingLiabilities:    3,
		LastModifiedLedger:   1000,
		NumSponsored:         12,
		NumSponsoring:        34,
		Sponsor:              null.StringFrom("GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H"),
	}

	ledgerWithCloseTime = &history.Ledger{
		ClosedAt: func() time.Time {
			t, err := time.Parse(time.RFC3339, "2019-03-05T13:23:50Z")
			if err != nil {
				panic(err)
			}
			return t
		}(),
	}

	trustLineIssuer = xdr.MustAddress("GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H")

	trustLines = []history.TrustLine{
		{
			AccountID:          accountID.Address(),
			AssetCode:          "EUR",
			AssetIssuer:        trustLineIssuer.Address(),
			AssetType:          xdr.AssetTypeAssetTypeCreditAlphanum4,
			Balance:            20000,
			Limit:              223456789,
			Flags:              1,
			SellingLiabilities: 3,
			BuyingLiabilities:  4,
			LastModifiedLedger: 900,
		},
		{
			AccountID:          accountID.Address(),
			AssetCode:          "USDDDDDDDDDD",
			AssetIssuer:        trustLineIssuer.Address(),
			AssetType:          xdr.AssetTypeAssetTypeCreditAlphanum12,
			Balance:            10000,
			Limit:              123456789,
			Flags:              0,
			SellingLiabilities: 2,
			BuyingLiabilities:  1,
			LastModifiedLedger: 900,
		},
		{
			AccountID:          accountID.Address(),
			AssetType:          xdr.AssetTypeAssetTypePoolShare,
			Balance:            10000,
			Limit:              123456789,
			Flags:              0,
			LastModifiedLedger: 900,
		},
	}

	signers = []history.AccountSigner{
		{
			Account: accountID.Address(),
			Signer:  accountID.Address(),
			Weight:  int32(3),
		},

		{
			Account: accountID.Address(),
			Signer:  "GCMQBJWOLTCSSMWNVDJAXL6E42SADH563IL5MN5B6RBBP4XP7TBRLJKE",
			Weight:  int32(1),
			Sponsor: null.StringFrom("GBXSGN5GX4PZOSBHB4JJF67CEGSGT56IN2N7LF3VGJ7WQ56BYWRVNNDX"),
		},
		{
			Account: accountID.Address(),
			Signer:  "GBXSGN5GX4PZOSBHB4JJF67CEGSGT56IN2N7LF3VGJ7WQ56BYWRVNNDX",
			Weight:  int32(2),
		},
		{
			Account: accountID.Address(),
			Signer:  "GBPXUGDRAOU5QUNUAXX6LYPBIOXYG45GLTKIRWKOCQ6HXP5QE5OCPFBY",
			Weight:  int32(3),
		},
	}
)

func TestPopulateAccountEntry(t *testing.T) {
	tt := assert.New(t)
	ctx, _ := test.ContextWithLogBuffer()
	hAccount := Account{}
	err := PopulateAccountEntry(ctx, &hAccount, account, data, signers, trustLines, ledgerWithCloseTime)
	tt.NoError(err)

	tt.Equal(account.AccountID, hAccount.ID)
	tt.Equal(account.AccountID, hAccount.AccountID)
	tt.Equal(account.AccountID, hAccount.PT)
	tt.Equal(strconv.FormatInt(account.SequenceNumber, 10), hAccount.Sequence)
	tt.Equal(int32(account.NumSubEntries), hAccount.SubentryCount)
	tt.Equal(account.InflationDestination, hAccount.InflationDestination)
	tt.Equal(account.HomeDomain, hAccount.HomeDomain)
	tt.Equal(account.LastModifiedLedger, hAccount.LastModifiedLedger)
	tt.NotNil(hAccount.LastModifiedTime)
	tt.Equal(ledgerWithCloseTime.ClosedAt, *hAccount.LastModifiedTime)
	tt.Equal(account.NumSponsoring, hAccount.NumSponsoring)
	tt.Equal(account.NumSponsored, hAccount.NumSponsored)
	tt.Equal(account.Sponsor.String, hAccount.Sponsor)

	wantAccountThresholds := AccountThresholds{
		LowThreshold:  account.ThresholdLow,
		MedThreshold:  account.ThresholdMedium,
		HighThreshold: account.ThresholdHigh,
	}
	tt.Equal(wantAccountThresholds, hAccount.Thresholds)

	wantFlags := AccountFlags{
		AuthRequired:        account.IsAuthRequired(),
		AuthRevocable:       account.IsAuthRevocable(),
		AuthImmutable:       account.IsAuthImmutable(),
		AuthClawbackEnabled: account.IsAuthClawbackEnabled(),
	}

	tt.Equal(wantFlags, hAccount.Flags)

	for _, d := range data {
		want, e := base64.StdEncoding.DecodeString(hAccount.Data[d.Name])
		tt.NoError(e)
		tt.Equal(d.Value, history.AccountDataValue(want))
	}

	tt.Len(hAccount.Balances, 4)

	for i, t := range trustLines {
		ht := hAccount.Balances[i]
		tt.Equal(t.AssetIssuer, ht.Issuer)
		tt.Equal(t.AssetCode, ht.Code)
		var wantType string
		if t.AssetType == xdr.AssetTypeAssetTypePoolShare {
			wantType = "liquidity_pool_shares"
		} else {
			wantType, err = assets.String(t.AssetType)
			tt.NoError(err)
		}
		tt.Equal(wantType, ht.Type)

		tt.Equal(amount.StringFromInt64(t.Balance), ht.Balance)

		wantBuy := ""
		if t.BuyingLiabilities != 0 {
			wantBuy = amount.StringFromInt64(t.BuyingLiabilities)
		}
		tt.Equal(wantBuy, ht.BuyingLiabilities)

		wantSell := ""
		if t.SellingLiabilities != 0 {
			wantSell = amount.StringFromInt64(t.SellingLiabilities)
		}
		tt.Equal(wantSell, ht.SellingLiabilities)

		tt.Equal(amount.StringFromInt64(t.Limit), ht.Limit)
		tt.Equal(t.LastModifiedLedger, ht.LastModifiedLedger)
		tt.Equal(t.IsAuthorized(), *ht.IsAuthorized)
	}

	native := hAccount.Balances[len(hAccount.Balances)-1]

	tt.Equal("native", native.Type)
	tt.Equal("0.0020000", native.Balance)
	tt.Equal("0.0000003", native.BuyingLiabilities)
	tt.Equal("0.0000004", native.SellingLiabilities)
	tt.Equal("", native.Limit)
	tt.Equal("", native.Issuer)
	tt.Equal("", native.Code)

	tt.Len(hAccount.Signers, 4)
	for i, s := range signers {
		hs := hAccount.Signers[i]
		tt.Equal(s.Signer, hs.Key)
		tt.Equal(s.Weight, hs.Weight)
		tt.Equal(protocol.MustKeyTypeFromAddress(s.Signer), hs.Type)
		var expectedSponsor string
		if s.Sponsor.Valid {
			expectedSponsor = s.Sponsor.String
		}
		tt.Equal(expectedSponsor, hs.Sponsor)
	}

	links, err := json.Marshal(hAccount.Links)
	want := `
	{
	  "data": {
		"href": "/accounts/GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB/data/{key}",
		"templated": true
	  },
	  "effects": {
		"href": "/accounts/GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB/effects{?cursor,limit,order}",
		"templated": true
	  },
	  "offers": {
		"href": "/accounts/GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB/offers{?cursor,limit,order}",
		"templated": true
	  },
	  "operations": {
		"href": "/accounts/GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB/operations{?cursor,limit,order}",
		"templated": true
	  },
	  "payments": {
		"href": "/accounts/GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB/payments{?cursor,limit,order}",
		"templated": true
	  },
	  "self": {
		"href": "/accounts/GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"
	  },
	  "trades": {
		"href": "/accounts/GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB/trades{?cursor,limit,order}",
		"templated": true
	  },
	  "transactions": {
		"href": "/accounts/GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB/transactions{?cursor,limit,order}",
		"templated": true
	  }
	}
	`
	tt.JSONEq(want, string(links))
}

func TestPopulateAccountEntryMasterMissingInSigners(t *testing.T) {
	tt := assert.New(t)
	ctx, _ := test.ContextWithLogBuffer()
	hAccount := Account{}

	account.MasterWeight = 0
	signers = []history.AccountSigner{
		{
			Account: accountID.Address(),
			Signer:  "GCMQBJWOLTCSSMWNVDJAXL6E42SADH563IL5MN5B6RBBP4XP7TBRLJKE",
			Weight:  int32(3),
		},
	}
	err := PopulateAccountEntry(ctx, &hAccount, account, data, signers, trustLines, nil)
	tt.NoError(err)

	tt.Len(hAccount.Signers, 2)

	signer := hAccount.Signers[1]
	tt.Equal(account.AccountID, signer.Key)
	tt.Equal(int32(account.MasterWeight), signer.Weight)
	tt.Equal(protocol.MustKeyTypeFromAddress(account.AccountID), signer.Type)
	tt.Nil(hAccount.LastModifiedTime)
}
