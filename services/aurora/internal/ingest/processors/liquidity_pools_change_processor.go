package processors

import (
	"context"

	"github.com/diamcircle/go/ingest"
	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/xdr"
)

type LiquidityPoolsChangeProcessor struct {
	qLiquidityPools history.QLiquidityPools
	cache           *ingest.ChangeCompactor
	sequence        uint32
}

func NewLiquidityPoolsChangeProcessor(Q history.QLiquidityPools, sequence uint32) *LiquidityPoolsChangeProcessor {
	p := &LiquidityPoolsChangeProcessor{
		qLiquidityPools: Q,
		sequence:        sequence,
	}
	p.reset()
	return p
}

func (p *LiquidityPoolsChangeProcessor) reset() {
	p.cache = ingest.NewChangeCompactor()
}

func (p *LiquidityPoolsChangeProcessor) ProcessChange(ctx context.Context, change ingest.Change) error {
	if change.Type != xdr.LedgerEntryTypeLiquidityPool {
		return nil
	}

	err := p.cache.AddChange(change)
	if err != nil {
		return errors.Wrap(err, "error adding to ledgerCache")
	}

	if p.cache.Size() > maxBatchSize {
		err = p.Commit(ctx)
		if err != nil {
			return errors.Wrap(err, "error in Commit")
		}
		p.reset()
	}

	return nil
}

func (p *LiquidityPoolsChangeProcessor) Commit(ctx context.Context) error {

	changes := p.cache.GetChanges()
	var lps []history.LiquidityPool
	for _, change := range changes {
		switch {
		case change.Pre == nil && change.Post != nil:
			// Created
			lps = append(lps, p.ledgerEntryToRow(change.Post))
		case change.Pre != nil && change.Post == nil:
			// Removed
			lp := p.ledgerEntryToRow(change.Pre)
			lp.Deleted = true
			lp.LastModifiedLedger = p.sequence
			lps = append(lps, lp)
		default:
			// Updated
			lps = append(lps, p.ledgerEntryToRow(change.Post))
		}
	}

	if len(lps) > 0 {
		if err := p.qLiquidityPools.UpsertLiquidityPools(ctx, lps); err != nil {
			return errors.Wrap(err, "error upserting liquidity pools")
		}
	}

	if p.sequence > compactionWindow {
		// trim liquidity pools table by removing liquidity pools which were deleted before the cutoff ledger
		if removed, err := p.qLiquidityPools.CompactLiquidityPools(ctx, p.sequence-compactionWindow); err != nil {
			return errors.Wrap(err, "could not compact liquidity pools")
		} else {
			log.WithField("liquidity_pool_rows_removed", removed).Info("Trimmed liquidity pools table")
		}
	}

	return nil
}

func (p *LiquidityPoolsChangeProcessor) ledgerEntryToRow(entry *xdr.LedgerEntry) history.LiquidityPool {
	lPool := entry.Data.MustLiquidityPool()
	cp := lPool.Body.MustConstantProduct()
	ar := history.LiquidityPoolAssetReserves{
		{
			Asset:   cp.Params.AssetA,
			Reserve: uint64(cp.ReserveA),
		},
		{
			Asset:   cp.Params.AssetB,
			Reserve: uint64(cp.ReserveB),
		},
	}
	return history.LiquidityPool{
		PoolID:             PoolIDToString(lPool.LiquidityPoolId),
		Type:               lPool.Body.Type,
		Fee:                uint32(cp.Params.Fee),
		TrustlineCount:     uint64(cp.PoolSharesTrustLineCount),
		ShareCount:         uint64(cp.TotalPoolShares),
		AssetReserves:      ar,
		LastModifiedLedger: uint32(entry.LastModifiedLedgerSeq),
	}
}

// PoolIDToString encodes a liquidity pool id xdr value to its string form
func PoolIDToString(id xdr.PoolId) string {
	return xdr.Hash(id).HexString()
}
