package ingest

import (
	"context"
	"testing"

	"go/xdr"

	"github.com/stretchr/testify/assert"
)

func TestStatsChangeProcessor(t *testing.T) {
	ctx := context.Background()
	processor := &StatsChangeProcessor{}

	// Created
	assert.NoError(t, processor.ProcessChange(ctx, Change{
		Type: xdr.LedgerEntryTypeAccount,
		Pre:  nil,
		Post: &xdr.LedgerEntry{},
	}))

	assert.NoError(t, processor.ProcessChange(ctx, Change{
		Type: xdr.LedgerEntryTypeClaimableBalance,
		Pre:  nil,
		Post: &xdr.LedgerEntry{},
	}))

	assert.NoError(t, processor.ProcessChange(ctx, Change{
		Type: xdr.LedgerEntryTypeData,
		Pre:  nil,
		Post: &xdr.LedgerEntry{},
	}))

	assert.NoError(t, processor.ProcessChange(ctx, Change{
		Type: xdr.LedgerEntryTypeOffer,
		Pre:  nil,
		Post: &xdr.LedgerEntry{},
	}))

	assert.NoError(t, processor.ProcessChange(ctx, Change{
		Type: xdr.LedgerEntryTypeTrustline,
		Pre:  nil,
		Post: &xdr.LedgerEntry{},
	}))

	assert.NoError(t, processor.ProcessChange(ctx, Change{
		Type: xdr.LedgerEntryTypeLiquidityPool,
		Pre:  nil,
		Post: &xdr.LedgerEntry{},
	}))

	// Updated
	assert.NoError(t, processor.ProcessChange(ctx, Change{
		Type: xdr.LedgerEntryTypeAccount,
		Pre:  &xdr.LedgerEntry{},
		Post: &xdr.LedgerEntry{},
	}))

	assert.NoError(t, processor.ProcessChange(ctx, Change{
		Type: xdr.LedgerEntryTypeClaimableBalance,
		Pre:  &xdr.LedgerEntry{},
		Post: &xdr.LedgerEntry{},
	}))

	assert.NoError(t, processor.ProcessChange(ctx, Change{
		Type: xdr.LedgerEntryTypeData,
		Pre:  &xdr.LedgerEntry{},
		Post: &xdr.LedgerEntry{},
	}))

	assert.NoError(t, processor.ProcessChange(ctx, Change{
		Type: xdr.LedgerEntryTypeOffer,
		Pre:  &xdr.LedgerEntry{},
		Post: &xdr.LedgerEntry{},
	}))

	assert.NoError(t, processor.ProcessChange(ctx, Change{
		Type: xdr.LedgerEntryTypeTrustline,
		Pre:  &xdr.LedgerEntry{},
		Post: &xdr.LedgerEntry{},
	}))

	assert.NoError(t, processor.ProcessChange(ctx, Change{
		Type: xdr.LedgerEntryTypeLiquidityPool,
		Pre:  &xdr.LedgerEntry{},
		Post: &xdr.LedgerEntry{},
	}))

	// Removed
	assert.NoError(t, processor.ProcessChange(ctx, Change{
		Type: xdr.LedgerEntryTypeAccount,
		Pre:  &xdr.LedgerEntry{},
		Post: nil,
	}))

	assert.NoError(t, processor.ProcessChange(ctx, Change{
		Type: xdr.LedgerEntryTypeClaimableBalance,
		Pre:  &xdr.LedgerEntry{},
		Post: nil,
	}))

	assert.NoError(t, processor.ProcessChange(ctx, Change{
		Type: xdr.LedgerEntryTypeData,
		Pre:  &xdr.LedgerEntry{},
		Post: nil,
	}))

	assert.NoError(t, processor.ProcessChange(ctx, Change{
		Type: xdr.LedgerEntryTypeOffer,
		Pre:  &xdr.LedgerEntry{},
		Post: nil,
	}))

	assert.NoError(t, processor.ProcessChange(ctx, Change{
		Type: xdr.LedgerEntryTypeTrustline,
		Pre:  &xdr.LedgerEntry{},
		Post: nil,
	}))

	assert.NoError(t, processor.ProcessChange(ctx, Change{
		Type: xdr.LedgerEntryTypeLiquidityPool,
		Pre:  &xdr.LedgerEntry{},
		Post: nil,
	}))

	results := processor.GetResults()

	assert.Equal(t, int64(1), results.AccountsCreated)
	assert.Equal(t, int64(1), results.ClaimableBalancesCreated)
	assert.Equal(t, int64(1), results.DataCreated)
	assert.Equal(t, int64(1), results.OffersCreated)
	assert.Equal(t, int64(1), results.TrustLinesCreated)
	assert.Equal(t, int64(1), results.LiquidityPoolsCreated)

	assert.Equal(t, int64(1), results.AccountsUpdated)
	assert.Equal(t, int64(1), results.ClaimableBalancesUpdated)
	assert.Equal(t, int64(1), results.DataUpdated)
	assert.Equal(t, int64(1), results.OffersUpdated)
	assert.Equal(t, int64(1), results.TrustLinesUpdated)
	assert.Equal(t, int64(1), results.LiquidityPoolsUpdated)

	assert.Equal(t, int64(1), results.AccountsRemoved)
	assert.Equal(t, int64(1), results.ClaimableBalancesRemoved)
	assert.Equal(t, int64(1), results.DataRemoved)
	assert.Equal(t, int64(1), results.OffersRemoved)
	assert.Equal(t, int64(1), results.TrustLinesRemoved)
	assert.Equal(t, int64(1), results.LiquidityPoolsRemoved)
}
