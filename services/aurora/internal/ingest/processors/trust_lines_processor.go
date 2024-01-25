package processors

import (
	"context"

	"go/ingest"
	"go/services/aurora/internal/db2/history"
	"go/support/errors"
	"go/xdr"
)

type TrustLinesProcessor struct {
	trustLinesQ history.QTrustLines

	cache *ingest.ChangeCompactor
}

func NewTrustLinesProcessor(trustLinesQ history.QTrustLines) *TrustLinesProcessor {
	p := &TrustLinesProcessor{trustLinesQ: trustLinesQ}
	p.reset()
	return p
}

func (p *TrustLinesProcessor) reset() {
	p.cache = ingest.NewChangeCompactor()
}

func (p *TrustLinesProcessor) ProcessChange(ctx context.Context, change ingest.Change) error {
	if change.Type != xdr.LedgerEntryTypeTrustline {
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

func trustLineLedgerKey(trustLineEntry xdr.TrustLineEntry) (string, error) {
	var ledgerKey xdr.LedgerKey
	var ledgerKeyString string

	err := ledgerKey.SetTrustline(trustLineEntry.AccountId, trustLineEntry.Asset)
	if err != nil {
		return "", errors.Wrap(err, "Error creating ledger key")
	}
	ledgerKeyString, err = ledgerKey.MarshalBinaryBase64()
	if err != nil {
		return "", errors.Wrap(err, "Error marshalling ledger key")
	}
	return ledgerKeyString, nil
}

func xdrToTrustline(ledgerEntry xdr.LedgerEntry) (history.TrustLine, error) {
	trustLineEntry := ledgerEntry.Data.MustTrustLine()
	ledgerKeyString, err := trustLineLedgerKey(trustLineEntry)
	if err != nil {
		return history.TrustLine{}, errors.Wrap(err, "Error extracting ledger key")
	}

	assetType := trustLineEntry.Asset.Type
	var assetCode, assetIssuer, poolID string
	if assetType == xdr.AssetTypeAssetTypePoolShare {
		poolID = PoolIDToString(trustLineEntry.Asset.MustLiquidityPoolId())
	} else {
		if err = trustLineEntry.Asset.ToAsset().Extract(&assetType, &assetCode, &assetIssuer); err != nil {
			return history.TrustLine{}, errors.Wrap(err, "Error extracting asset from trustline")
		}
	}

	liabilities := trustLineEntry.Liabilities()
	return history.TrustLine{
		AccountID:          trustLineEntry.AccountId.Address(),
		AssetType:          assetType,
		AssetIssuer:        assetIssuer,
		AssetCode:          assetCode,
		Balance:            int64(trustLineEntry.Balance),
		LedgerKey:          ledgerKeyString,
		Limit:              int64(trustLineEntry.Limit),
		LiquidityPoolID:    poolID,
		BuyingLiabilities:  int64(liabilities.Buying),
		SellingLiabilities: int64(liabilities.Selling),
		Flags:              uint32(trustLineEntry.Flags),
		LastModifiedLedger: uint32(ledgerEntry.LastModifiedLedgerSeq),
		Sponsor:            ledgerEntrySponsorToNullString(ledgerEntry),
	}, nil
}

func (p *TrustLinesProcessor) Commit(ctx context.Context) error {
	var batchUpsertTrustLines []history.TrustLine
	var batchRemoveTrustLineKeys []string

	changes := p.cache.GetChanges()
	for _, change := range changes {
		switch {
		case change.Post != nil:
			tl, err := xdrToTrustline(*change.Post)
			if err != nil {
				return errors.Wrap(err, "Error extracting trustline")
			}

			batchUpsertTrustLines = append(batchUpsertTrustLines, tl)
		case change.Pre != nil && change.Post == nil:
			// Removed
			trustLineEntry := change.Pre.Data.MustTrustLine()
			ledgerKeyString, err := trustLineLedgerKey(trustLineEntry)
			if err != nil {
				return errors.Wrap(err, "Error extracting ledger key")
			}
			batchRemoveTrustLineKeys = append(batchRemoveTrustLineKeys, ledgerKeyString)

		default:
			return errors.New("Invalid io.Change: change.Pre == nil && change.Post == nil")
		}
	}

	if len(batchUpsertTrustLines) > 0 {
		err := p.trustLinesQ.UpsertTrustLines(ctx, batchUpsertTrustLines)
		if err != nil {
			return errors.Wrap(err, "errors in UpsertTrustLines")
		}
	}

	if len(batchRemoveTrustLineKeys) > 0 {
		rowsAffected, err := p.trustLinesQ.RemoveTrustLines(ctx, batchRemoveTrustLineKeys)
		if err != nil {
			return err
		}

		if rowsAffected != int64(len(batchRemoveTrustLineKeys)) {
			return ingest.NewStateError(errors.Errorf(
				"%d rows affected when removing %d trust lines",
				rowsAffected,
				len(batchRemoveTrustLineKeys),
			))
		}
	}

	return nil
}
