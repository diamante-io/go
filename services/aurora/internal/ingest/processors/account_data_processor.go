package processors

import (
	"context"

	"go/ingest"
	"go/services/aurora/internal/db2/history"
	"go/support/errors"
	"go/xdr"
)

type AccountDataProcessor struct {
	dataQ history.QData

	cache *ingest.ChangeCompactor
}

func NewAccountDataProcessor(dataQ history.QData) *AccountDataProcessor {
	p := &AccountDataProcessor{dataQ: dataQ}
	p.reset()
	return p
}

func (p *AccountDataProcessor) reset() {
	p.cache = ingest.NewChangeCompactor()
}

func (p *AccountDataProcessor) ProcessChange(ctx context.Context, change ingest.Change) error {
	// We're interested in data only
	if change.Type != xdr.LedgerEntryTypeData {
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

func (p *AccountDataProcessor) Commit(ctx context.Context) error {
	var (
		datasToUpsert []history.Data
		datasToDelete []history.AccountDataKey
	)
	changes := p.cache.GetChanges()
	for _, change := range changes {
		switch {
		case change.Pre == nil && change.Post != nil:
			// Created
			datasToUpsert = append(datasToUpsert, p.ledgerEntryToRow(change.Post))
		case change.Pre != nil && change.Post == nil:
			// Removed
			data := change.Pre.Data.MustData()
			key := history.AccountDataKey{
				AccountID: data.AccountId.Address(),
				DataName:  string(data.DataName),
			}
			datasToDelete = append(datasToDelete, key)
		default:
			// Updated
			datasToUpsert = append(datasToUpsert, p.ledgerEntryToRow(change.Post))
		}
	}

	if len(datasToUpsert) > 0 {
		if err := p.dataQ.UpsertAccountData(ctx, datasToUpsert); err != nil {
			return errors.Wrap(err, "error executing upsert")
		}
	}

	if len(datasToDelete) > 0 {
		count, err := p.dataQ.RemoveAccountData(ctx, datasToDelete)
		if err != nil {
			return errors.Wrap(err, "error executing removal")
		}
		if count != int64(len(datasToDelete)) {
			return ingest.NewStateError(errors.Errorf(
				"%d rows affected when deleting %d account data",
				count,
				len(datasToDelete),
			))
		}
	}

	return nil
}

func (p *AccountDataProcessor) ledgerEntryToRow(entry *xdr.LedgerEntry) history.Data {
	data := entry.Data.MustData()
	return history.Data{
		AccountID:          data.AccountId.Address(),
		Name:               string(data.DataName),
		Value:              history.AccountDataValue(data.DataValue),
		LastModifiedLedger: uint32(entry.LastModifiedLedgerSeq),
		Sponsor:            ledgerEntrySponsorToNullString(*entry),
	}
}
