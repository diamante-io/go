package processors

import (
	"context"
	"database/sql"

	"go/ingest"
	"go/services/aurora/internal/db2/history"
	"go/support/errors"
	"go/xdr"
)

type AssetStatsProcessor struct {
	assetStatsQ history.QAssetStats

	cache               *ingest.ChangeCompactor
	assetStatSet        AssetStatSet
	useLedgerEntryCache bool
}

// NewAssetStatsProcessor constructs a new AssetStatsProcessor instance.
// If useLedgerEntryCache is false we don't use ledger cache and we just
// add trust lines to assetStatSet, then we insert all the stats in one
// insert query. This is done to make history buckets processing faster
// (batch inserting).
func NewAssetStatsProcessor(
	assetStatsQ history.QAssetStats,
	useLedgerEntryCache bool,
) *AssetStatsProcessor {
	p := &AssetStatsProcessor{
		assetStatsQ:         assetStatsQ,
		useLedgerEntryCache: useLedgerEntryCache,
	}
	p.reset()
	return p
}

func (p *AssetStatsProcessor) reset() {
	p.cache = ingest.NewChangeCompactor()
	p.assetStatSet = AssetStatSet{}
}

func (p *AssetStatsProcessor) ProcessChange(ctx context.Context, change ingest.Change) error {
	if change.Type != xdr.LedgerEntryTypeLiquidityPool &&
		change.Type != xdr.LedgerEntryTypeClaimableBalance &&
		change.Type != xdr.LedgerEntryTypeTrustline {
		return nil
	}
	if p.useLedgerEntryCache {
		return p.addToCache(ctx, change)
	}
	if change.Pre != nil || change.Post == nil {
		return errors.New("AssetStatsProcessor is in insert only mode")
	}

	switch change.Type {
	case xdr.LedgerEntryTypeLiquidityPool:
		return p.assetStatSet.AddLiquidityPool(change)
	case xdr.LedgerEntryTypeClaimableBalance:
		return p.assetStatSet.AddClaimableBalance(change)
	case xdr.LedgerEntryTypeTrustline:
		return p.assetStatSet.AddTrustline(change)
	default:
		return nil
	}
}

func (p *AssetStatsProcessor) addToCache(ctx context.Context, change ingest.Change) error {
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

func (p *AssetStatsProcessor) Commit(ctx context.Context) error {
	if !p.useLedgerEntryCache {
		assetStatsDeltas := p.assetStatSet.All()
		if len(assetStatsDeltas) == 0 {
			return nil
		}
		return p.assetStatsQ.InsertAssetStats(ctx, assetStatsDeltas, maxBatchSize)
	}

	changes := p.cache.GetChanges()
	for _, change := range changes {
		var err error
		switch change.Type {
		case xdr.LedgerEntryTypeLiquidityPool:
			err = p.assetStatSet.AddLiquidityPool(change)
		case xdr.LedgerEntryTypeClaimableBalance:
			err = p.assetStatSet.AddClaimableBalance(change)
		case xdr.LedgerEntryTypeTrustline:
			err = p.assetStatSet.AddTrustline(change)
		default:
			return errors.Errorf("Change type %v is unexpected", change.Type)
		}

		if err != nil {
			return errors.Wrap(err, "Error adjusting asset stat")
		}
	}

	assetStatsDeltas := p.assetStatSet.All()
	for _, delta := range assetStatsDeltas {
		var rowsAffected int64
		var stat history.ExpAssetStat
		var err error

		stat, err = p.assetStatsQ.GetAssetStat(ctx,
			delta.AssetType,
			delta.AssetCode,
			delta.AssetIssuer,
		)
		assetStatNotFound := err == sql.ErrNoRows
		if !assetStatNotFound && err != nil {
			return errors.Wrap(err, "could not fetch asset stat from db")
		}

		if assetStatNotFound {
			// Safety checks
			if delta.Accounts.Authorized < 0 {
				return ingest.NewStateError(errors.Errorf(
					"Authorized accounts negative but DB entry does not exist for asset: %s %s %s",
					delta.AssetType,
					delta.AssetCode,
					delta.AssetIssuer,
				))
			} else if delta.Accounts.AuthorizedToMaintainLiabilities < 0 {
				return ingest.NewStateError(errors.Errorf(
					"AuthorizedToMaintainLiabilities accounts negative but DB entry does not exist for asset: %s %s %s",
					delta.AssetType,
					delta.AssetCode,
					delta.AssetIssuer,
				))
			} else if delta.Accounts.Unauthorized < 0 {
				return ingest.NewStateError(errors.Errorf(
					"Unauthorized accounts negative but DB entry does not exist for asset: %s %s %s",
					delta.AssetType,
					delta.AssetCode,
					delta.AssetIssuer,
				))
			} else if delta.Accounts.ClaimableBalances < 0 {
				return ingest.NewStateError(errors.Errorf(
					"Claimable balance negative but DB entry does not exist for asset: %s %s %s",
					delta.AssetType,
					delta.AssetCode,
					delta.AssetIssuer,
				))
			} else if delta.Accounts.LiquidityPools < 0 {
				return ingest.NewStateError(errors.Errorf(
					"Liquidity pools negative but DB entry does not exist for asset: %s %s %s",
					delta.AssetType,
					delta.AssetCode,
					delta.AssetIssuer,
				))
			}

			// Insert
			var errInsert error
			rowsAffected, errInsert = p.assetStatsQ.InsertAssetStat(ctx, delta)
			if errInsert != nil {
				return errors.Wrap(errInsert, "could not insert asset stat")
			}
		} else {
			var statBalances assetStatBalances
			if err = statBalances.Parse(&stat.Balances); err != nil {
				return errors.Wrap(err, "Error parsing balances")
			}

			var deltaBalances assetStatBalances
			if err = deltaBalances.Parse(&delta.Balances); err != nil {
				return errors.Wrap(err, "Error parsing balances")
			}

			statBalances = statBalances.Add(deltaBalances)
			statAccounts := stat.Accounts.Add(delta.Accounts)

			if statAccounts.IsZero() {
				// Remove stats
				if !statBalances.IsZero() {
					return ingest.NewStateError(errors.Errorf(
						"Removing asset stat by final amount non-zero for: %s %s %s",
						delta.AssetType,
						delta.AssetCode,
						delta.AssetIssuer,
					))
				}
				rowsAffected, err = p.assetStatsQ.RemoveAssetStat(ctx,
					delta.AssetType,
					delta.AssetCode,
					delta.AssetIssuer,
				)
				if err != nil {
					return errors.Wrap(err, "could not remove asset stat")
				}
			} else {
				// Update
				rowsAffected, err = p.assetStatsQ.UpdateAssetStat(ctx, history.ExpAssetStat{
					AssetType:   delta.AssetType,
					AssetCode:   delta.AssetCode,
					AssetIssuer: delta.AssetIssuer,
					Accounts:    statAccounts,
					Balances:    statBalances.ConvertToHistoryObject(),
					Amount:      statBalances.Authorized.String(),
					NumAccounts: statAccounts.Authorized,
				})
				if err != nil {
					return errors.Wrap(err, "could not update asset stat")
				}
			}
		}

		if rowsAffected != 1 {
			return ingest.NewStateError(errors.Errorf(
				"%d rows affected when adjusting asset stat for asset: %s %s %s",
				rowsAffected,
				delta.AssetType,
				delta.AssetCode,
				delta.AssetIssuer,
			))
		}
	}

	return nil
}
