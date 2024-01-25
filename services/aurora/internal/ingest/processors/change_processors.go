package processors

import (
	"context"
	"io"

	"go/ingest"
	"go/support/errors"
)

type ChangeProcessor interface {
	ProcessChange(ctx context.Context, change ingest.Change) error
}

type LedgerTransactionProcessor interface {
	ProcessTransaction(ctx context.Context, transaction ingest.LedgerTransaction) error
}

func StreamLedgerTransactions(
	ctx context.Context,
	txProcessor LedgerTransactionProcessor,
	reader *ingest.LedgerTransactionReader,
) error {
	for {
		tx, err := reader.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return errors.Wrap(err, "could not read transaction")
		}
		if err = txProcessor.ProcessTransaction(ctx, tx); err != nil {
			return errors.Wrapf(
				err,
				"could not process transaction %v",
				tx.Index,
			)
		}
	}
}

func StreamChanges(
	ctx context.Context,
	changeProcessor ChangeProcessor,
	reader ingest.ChangeReader,
) error {
	for {
		change, err := reader.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return errors.Wrap(err, "could not read transaction")
		}

		if err = changeProcessor.ProcessChange(ctx, change); err != nil {
			return errors.Wrap(
				err,
				"could not process change",
			)
		}
	}
}
