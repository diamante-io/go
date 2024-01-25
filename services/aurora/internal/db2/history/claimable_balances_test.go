package history

import (
	"testing"

	"github.com/guregu/null"
	"github.com/diamcircle/go/services/aurora/internal/db2"
	"github.com/diamcircle/go/services/aurora/internal/test"
	"github.com/diamcircle/go/xdr"
)

func TestRemoveClaimableBalance(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}

	accountID := "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"
	asset := xdr.MustNewCreditAsset("USD", accountID)
	balanceID := xdr.ClaimableBalanceId{
		Type: xdr.ClaimableBalanceIdTypeClaimableBalanceIdTypeV0,
		V0:   &xdr.Hash{1, 2, 3},
	}
	id, err := xdr.MarshalHex(balanceID)
	tt.Assert.NoError(err)
	cBalance := ClaimableBalance{
		BalanceID: id,
		Claimants: []Claimant{
			{
				Destination: accountID,
				Predicate: xdr.ClaimPredicate{
					Type: xdr.ClaimPredicateTypeClaimPredicateUnconditional,
				},
			},
		},
		Asset:              asset,
		LastModifiedLedger: 123,
		Amount:             10,
	}

	err = q.UpsertClaimableBalances(tt.Ctx, []ClaimableBalance{cBalance})
	tt.Assert.NoError(err)

	r, err := q.FindClaimableBalanceByID(tt.Ctx, id)
	tt.Assert.NoError(err)
	tt.Assert.NotNil(r)

	removed, err := q.RemoveClaimableBalances(tt.Ctx, []string{id})
	tt.Assert.NoError(err)
	tt.Assert.Equal(int64(1), removed)

	cbs := []ClaimableBalance{}
	err = q.Select(tt.Ctx, &cbs, selectClaimableBalances)

	if tt.Assert.NoError(err) {
		tt.Assert.Len(cbs, 0)
	}
}

func TestFindClaimableBalancesByDestination(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}

	dest1 := "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"
	dest2 := "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H"

	asset := xdr.MustNewCreditAsset("USD", dest1)
	balanceID := xdr.ClaimableBalanceId{
		Type: xdr.ClaimableBalanceIdTypeClaimableBalanceIdTypeV0,
		V0:   &xdr.Hash{1, 2, 3},
	}
	id, err := xdr.MarshalHex(balanceID)
	tt.Assert.NoError(err)
	cBalance := ClaimableBalance{
		BalanceID: id,
		Claimants: []Claimant{
			{
				Destination: dest1,
				Predicate: xdr.ClaimPredicate{
					Type: xdr.ClaimPredicateTypeClaimPredicateUnconditional,
				},
			},
		},
		Asset:              asset,
		LastModifiedLedger: 123,
		Amount:             10,
	}

	err = q.UpsertClaimableBalances(tt.Ctx, []ClaimableBalance{cBalance})
	tt.Assert.NoError(err)

	balanceID = xdr.ClaimableBalanceId{
		Type: xdr.ClaimableBalanceIdTypeClaimableBalanceIdTypeV0,
		V0:   &xdr.Hash{3, 2, 1},
	}
	id, err = xdr.MarshalHex(balanceID)
	tt.Assert.NoError(err)
	cBalance = ClaimableBalance{
		BalanceID: id,
		Claimants: []Claimant{
			{
				Destination: dest1,
				Predicate: xdr.ClaimPredicate{
					Type: xdr.ClaimPredicateTypeClaimPredicateUnconditional,
				},
			},
			{
				Destination: dest2,
				Predicate: xdr.ClaimPredicate{
					Type: xdr.ClaimPredicateTypeClaimPredicateUnconditional,
				},
			},
		},
		Asset:              asset,
		LastModifiedLedger: 123,
		Amount:             10,
	}

	err = q.UpsertClaimableBalances(tt.Ctx, []ClaimableBalance{cBalance})
	tt.Assert.NoError(err)

	query := ClaimableBalancesQuery{
		PageQuery: db2.MustPageQuery("", false, "", 10),
		Claimant:  xdr.MustAddressPtr(dest1),
	}

	cbs, err := q.GetClaimableBalances(tt.Ctx, query)
	tt.Assert.NoError(err)
	tt.Assert.Len(cbs, 2)

	for _, cb := range cbs {
		tt.Assert.Equal(dest1, cb.Claimants[0].Destination)
	}

	query.Claimant = xdr.MustAddressPtr(dest2)
	cbs, err = q.GetClaimableBalances(tt.Ctx, query)
	tt.Assert.NoError(err)
	tt.Assert.Len(cbs, 1)
	tt.Assert.Equal(dest2, cbs[0].Claimants[1].Destination)
}

func TestUpdateClaimableBalance(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}

	accountID := "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"
	lastModifiedLedgerSeq := xdr.Uint32(123)
	asset := xdr.MustNewCreditAsset("USD", accountID)
	balanceID := xdr.ClaimableBalanceId{
		Type: xdr.ClaimableBalanceIdTypeClaimableBalanceIdTypeV0,
		V0:   &xdr.Hash{1, 2, 3},
	}
	id, err := xdr.MarshalHex(balanceID)
	tt.Assert.NoError(err)
	cBalance := ClaimableBalance{
		BalanceID: id,
		Claimants: []Claimant{
			{
				Destination: accountID,
				Predicate: xdr.ClaimPredicate{
					Type: xdr.ClaimPredicateTypeClaimPredicateUnconditional,
				},
			},
		},
		Asset:              asset,
		LastModifiedLedger: 123,
		Amount:             10,
	}

	err = q.UpsertClaimableBalances(tt.Ctx, []ClaimableBalance{cBalance})
	tt.Assert.NoError(err)

	// add sponsor
	cBalance2 := ClaimableBalance{
		BalanceID: id,
		Claimants: []Claimant{
			{
				Destination: accountID,
				Predicate: xdr.ClaimPredicate{
					Type: xdr.ClaimPredicateTypeClaimPredicateUnconditional,
				},
			},
		},
		Asset:              asset,
		LastModifiedLedger: 123 + 1,
		Amount:             10,
		Sponsor:            null.StringFrom("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
	}

	err = q.UpsertClaimableBalances(tt.Ctx, []ClaimableBalance{cBalance2})
	tt.Assert.NoError(err)

	cbs := []ClaimableBalance{}
	err = q.Select(tt.Ctx, &cbs, selectClaimableBalances)
	tt.Assert.NoError(err)
	tt.Assert.Len(cbs, 1)
	tt.Assert.Equal("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML", cbs[0].Sponsor.String)
	tt.Assert.Equal(uint32(lastModifiedLedgerSeq+1), cbs[0].LastModifiedLedger)
}

func TestFindClaimableBalance(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}

	accountID := "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"
	asset := xdr.MustNewCreditAsset("USD", accountID)
	balanceID := xdr.ClaimableBalanceId{
		Type: xdr.ClaimableBalanceIdTypeClaimableBalanceIdTypeV0,
		V0:   &xdr.Hash{1, 2, 3},
	}
	id, err := xdr.MarshalHex(balanceID)
	tt.Assert.NoError(err)
	cBalance := ClaimableBalance{
		BalanceID: id,
		Claimants: []Claimant{
			{
				Destination: accountID,
				Predicate: xdr.ClaimPredicate{
					Type: xdr.ClaimPredicateTypeClaimPredicateUnconditional,
				},
			},
		},
		Asset:              asset,
		LastModifiedLedger: 123,
		Amount:             10,
	}

	err = q.UpsertClaimableBalances(tt.Ctx, []ClaimableBalance{cBalance})
	tt.Assert.NoError(err)

	cb, err := q.FindClaimableBalanceByID(tt.Ctx, id)
	tt.Assert.NoError(err)

	tt.Assert.Equal(cBalance.BalanceID, cb.BalanceID)
	tt.Assert.Equal(cBalance.Asset, cb.Asset)
	tt.Assert.Equal(cBalance.Amount, cb.Amount)

	for i, hClaimant := range cb.Claimants {
		tt.Assert.Equal(cBalance.Claimants[i].Destination, hClaimant.Destination)
		tt.Assert.Equal(cBalance.Claimants[i].Predicate, hClaimant.Predicate)
	}
}
func TestGetClaimableBalancesByID(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}

	accountID := "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"
	asset := xdr.MustNewCreditAsset("USD", accountID)
	balanceID := xdr.ClaimableBalanceId{
		Type: xdr.ClaimableBalanceIdTypeClaimableBalanceIdTypeV0,
		V0:   &xdr.Hash{1, 2, 3},
	}
	id, err := xdr.MarshalHex(balanceID)
	tt.Assert.NoError(err)
	cBalance := ClaimableBalance{
		BalanceID: id,
		Claimants: []Claimant{
			{
				Destination: accountID,
				Predicate: xdr.ClaimPredicate{
					Type: xdr.ClaimPredicateTypeClaimPredicateUnconditional,
				},
			},
		},
		Asset:              asset,
		LastModifiedLedger: 123,
		Amount:             10,
	}

	err = q.UpsertClaimableBalances(tt.Ctx, []ClaimableBalance{cBalance})
	tt.Assert.NoError(err)

	r, err := q.GetClaimableBalancesByID(tt.Ctx, []string{id})
	tt.Assert.NoError(err)
	tt.Assert.Len(r, 1)

	removed, err := q.RemoveClaimableBalances(tt.Ctx, []string{id})
	tt.Assert.NoError(err)
	tt.Assert.Equal(int64(1), removed)

	r, err = q.GetClaimableBalancesByID(tt.Ctx, []string{id})
	tt.Assert.NoError(err)
	tt.Assert.Len(r, 0)
}
