package processors

import (
	"context"

	"go/ingest"
	"go/services/aurora/internal/db2/history"
	"go/support/errors"
	"go/xdr"
)

type AccountsProcessor struct {
	accountsQ history.QAccounts

	cache *ingest.ChangeCompactor
}

func NewAccountsProcessor(accountsQ history.QAccounts) *AccountsProcessor {
	p := &AccountsProcessor{accountsQ: accountsQ}
	p.reset()
	return p
}

func (p *AccountsProcessor) reset() {
	p.cache = ingest.NewChangeCompactor()
}

func (p *AccountsProcessor) ProcessChange(ctx context.Context, change ingest.Change) error {
	if change.Type != xdr.LedgerEntryTypeAccount {
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

func (p *AccountsProcessor) Commit(ctx context.Context) error {
	batchUpsertAccounts := []history.AccountEntry{}
	removeBatch := []string{}

	changes := p.cache.GetChanges()
	for _, change := range changes {
		changed, err := change.AccountChangedExceptSigners()
		if err != nil {
			return errors.Wrap(err, "Error running change.AccountChangedExceptSigners")
		}

		if !changed {
			continue
		}

		switch {
		case change.Post != nil:
			// Created and updated
			row := p.ledgerEntryToRow(*change.Post)
			batchUpsertAccounts = append(batchUpsertAccounts, row)
		case change.Pre != nil && change.Post == nil:
			// Removed
			account := change.Pre.Data.MustAccount()
			accountID := account.AccountId.Address()
			removeBatch = append(removeBatch, accountID)
		default:
			return errors.New("Invalid io.Change: change.Pre == nil && change.Post == nil")
		}
	}

	// Upsert accounts
	if len(batchUpsertAccounts) > 0 {
		err := p.accountsQ.UpsertAccounts(ctx, batchUpsertAccounts)
		if err != nil {
			return errors.Wrap(err, "errors in UpsertAccounts")
		}
	}

	if len(removeBatch) > 0 {
		rowsAffected, err := p.accountsQ.RemoveAccounts(ctx, removeBatch)
		if err != nil {
			return errors.Wrap(err, "error in RemoveAccounts")
		}

		if rowsAffected != int64(len(removeBatch)) {
			return ingest.NewStateError(errors.Errorf(
				"%d rows affected when removing %d accounts",
				rowsAffected,
				len(removeBatch),
			))
		}
	}

	return nil
}

func (p *AccountsProcessor) ledgerEntryToRow(entry xdr.LedgerEntry) history.AccountEntry {
	account := entry.Data.MustAccount()
	liabilities := account.Liabilities()

	var inflationDestination = ""
	if account.InflationDest != nil {
		inflationDestination = account.InflationDest.Address()
	}

	return history.AccountEntry{
		AccountID:            account.AccountId.Address(),
		Balance:              int64(account.Balance),
		BuyingLiabilities:    int64(liabilities.Buying),
		SellingLiabilities:   int64(liabilities.Selling),
		SequenceNumber:       int64(account.SeqNum),
		NumSubEntries:        uint32(account.NumSubEntries),
		InflationDestination: inflationDestination,
		Flags:                uint32(account.Flags),
		HomeDomain:           string(account.HomeDomain),
		MasterWeight:         account.MasterKeyWeight(),
		ThresholdLow:         account.ThresholdLow(),
		ThresholdMedium:      account.ThresholdMedium(),
		ThresholdHigh:        account.ThresholdHigh(),
		LastModifiedLedger:   uint32(entry.LastModifiedLedgerSeq),
		Sponsor:              ledgerEntrySponsorToNullString(entry),
		NumSponsored:         uint32(account.NumSponsored()),
		NumSponsoring:        uint32(account.NumSponsoring()),
	}
}
