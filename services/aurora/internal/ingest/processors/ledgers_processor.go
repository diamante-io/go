package processors

import (
	"context"

	"github.com/diamcircle/go/ingest"
	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/xdr"
)

type LedgersProcessor struct {
	ledgersQ       history.QLedgers
	ledger         xdr.LedgerHeaderHistoryEntry
	ingestVersion  int
	successTxCount int
	failedTxCount  int
	opCount        int
	txSetOpCount   int
}

func NewLedgerProcessor(
	ledgerQ history.QLedgers,
	ledger xdr.LedgerHeaderHistoryEntry,
	ingestVersion int,
) *LedgersProcessor {
	return &LedgersProcessor{
		ledger:        ledger,
		ledgersQ:      ledgerQ,
		ingestVersion: ingestVersion,
	}
}

func (p *LedgersProcessor) ProcessTransaction(ctx context.Context, transaction ingest.LedgerTransaction) (err error) {
	opCount := len(transaction.Envelope.Operations())
	p.txSetOpCount += opCount
	if transaction.Result.Successful() {
		p.successTxCount++
		p.opCount += opCount
	} else {
		p.failedTxCount++
	}

	return nil
}

func (p *LedgersProcessor) Commit(ctx context.Context) error {
	rowsAffected, err := p.ledgersQ.InsertLedger(ctx,
		p.ledger,
		p.successTxCount,
		p.failedTxCount,
		p.opCount,
		p.txSetOpCount,
		p.ingestVersion,
	)

	if err != nil {
		return errors.Wrap(err, "Could not insert ledger")
	}

	sequence := uint32(p.ledger.Header.LedgerSeq)

	if rowsAffected != 1 {
		log.WithField("rowsAffected", rowsAffected).
			WithField("sequence", sequence).
			Error("Invalid number of rows affected when ingesting new ledger")
		return errors.Errorf(
			"0 rows affected when ingesting new ledger: %v",
			sequence,
		)
	}

	return nil
}
