package resourceadapter

import (
	"encoding/json"
	"testing"

	"github.com/guregu/null"
	. "github.com/diamcircle/go/protocols/aurora"
	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/support/test"
	"github.com/diamcircle/go/xdr"
	"github.com/stretchr/testify/assert"
)

func TestPopulateClaimableBalance(t *testing.T) {
	tt := assert.New(t)
	ctx, _ := test.ContextWithLogBuffer()
	resource := ClaimableBalance{}

	balanceID := xdr.ClaimableBalanceId{
		Type: xdr.ClaimableBalanceIdTypeClaimableBalanceIdTypeV0,
		V0:   &xdr.Hash{1, 2, 3},
	}
	unconditional := &xdr.ClaimPredicate{
		Type: xdr.ClaimPredicateTypeClaimPredicateUnconditional,
	}
	relBefore := xdr.Int64(12)
	absBefore := xdr.Int64(1598440539)

	id, err := xdr.MarshalHex(&balanceID)
	tt.NoError(err)
	claimableBalance := history.ClaimableBalance{
		BalanceID: id,
		Asset:     xdr.MustNewNativeAsset(),
		Amount:    100000000,
		Sponsor:   null.StringFrom("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
		Claimants: history.Claimants{
			{
				Destination: "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML",
				Predicate: xdr.ClaimPredicate{
					Type: xdr.ClaimPredicateTypeClaimPredicateAnd,
					AndPredicates: &[]xdr.ClaimPredicate{
						{
							Type: xdr.ClaimPredicateTypeClaimPredicateOr,
							OrPredicates: &[]xdr.ClaimPredicate{
								{
									Type:      xdr.ClaimPredicateTypeClaimPredicateBeforeRelativeTime,
									RelBefore: &relBefore,
								},
								{
									Type:      xdr.ClaimPredicateTypeClaimPredicateBeforeAbsoluteTime,
									AbsBefore: &absBefore,
								},
							},
						},
						{
							Type:         xdr.ClaimPredicateTypeClaimPredicateNot,
							NotPredicate: &unconditional,
						},
					},
				},
			},
		},
		LastModifiedLedger: 123,
		Flags:              uint32(xdr.ClaimableBalanceFlagsClaimableBalanceClawbackEnabledFlag),
	}

	err = PopulateClaimableBalance(ctx, &resource, claimableBalance, nil)
	tt.NoError(err)

	tt.Equal("000000000102030000000000000000000000000000000000000000000000000000000000", resource.BalanceID)
	tt.Equal(claimableBalance.Asset.StringCanonical(), resource.Asset)
	tt.Equal("10.0000000", resource.Amount)
	tt.Equal("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML", resource.Sponsor)
	tt.Equal(uint32(123), resource.LastModifiedLedger)
	tt.Len(resource.Claimants, 1)
	tt.Equal("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML", resource.Claimants[0].Destination)
	tt.Equal("123-000000000102030000000000000000000000000000000000000000000000000000000000", resource.PagingToken())
	tt.True(resource.Flags.ClawbackEnabled)

	links, err := json.Marshal(resource.Links)
	tt.NoError(err)
	want := `
	{
	  "self": {
		"href": "/claimable_balances/000000000102030000000000000000000000000000000000000000000000000000000000"
	  },
	  "operations": {
		"href": "/claimable_balances/000000000102030000000000000000000000000000000000000000000000000000000000/operations{?cursor,limit,order}",
        "templated": true
	  },
	  "transactions": {
		"href": "/claimable_balances/000000000102030000000000000000000000000000000000000000000000000000000000/transactions{?cursor,limit,order}",
        "templated": true
	  }
	}
	`
	tt.JSONEq(want, string(links))

	predicate, err := json.Marshal(resource.Claimants[0].Predicate)
	tt.NoError(err)
	tt.JSONEq(`{"and":[{"or":[{"rel_before":"12"},{"abs_before":"2020-08-26T11:15:39Z","abs_before_epoch":"1598440539"}]},{"not":{"unconditional":true}}]}`, string(predicate))
}
