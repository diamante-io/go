package txsub

import (
	"context"
	"database/sql"

	"go/services/aurora/internal/db2/history"
	"go/support/errors"
	"go/xdr"
)

func txResultByHash(ctx context.Context, db auroraDB, hash string) (history.Transaction, error) {
	// query history database
	var hr history.Transaction
	err := db.TransactionByHash(ctx, &hr, hash)
	if err == nil {
		return txResultFromHistory(hr)
	}

	if !db.NoRows(err) {
		return hr, errors.Wrap(err, "could not lookup transaction by hash")
	}

	// if no result was found in either db, return ErrNoResults
	return hr, ErrNoResults
}

func txResultFromHistory(tx history.Transaction) (history.Transaction, error) {
	var txResult xdr.TransactionResult
	err := xdr.SafeUnmarshalBase64(tx.TxResult, &txResult)
	if err == nil {
		if !txResult.Successful() {
			err = &FailedTransactionError{
				ResultXDR: tx.TxResult,
			}
		}
	} else {
		err = errors.Wrap(err, "could not unmarshall transaction result")
	}

	return tx, err
}

// checkTxAlreadyExists uses a repeatable read transaction to look up both transaction results
// and sequence numbers. Without the repeatable read transaction it is possible that the two database
// queries execute on different ledgers. In this case, txsub can mistakenly respond with a bad_seq error
// because the first query occurs when the tx is not yet ingested and the second query occurs when the tx
// is ingested.
func checkTxAlreadyExists(ctx context.Context, db auroraDB, hash, sourceAddress string) (history.Transaction, uint64, error) {
	err := db.BeginTx(&sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	})
	if err != nil {
		return history.Transaction{}, 0, errors.Wrap(err, "cannot start repeatable read tx")
	}
	defer db.Rollback()

	tx, err := txResultByHash(ctx, db, hash)
	if err == ErrNoResults {
		var sequenceNumbers map[string]uint64
		sequenceNumbers, err = db.GetSequenceNumbers(ctx, []string{sourceAddress})
		if err != nil {
			return tx, 0, errors.Wrapf(err, "cannot fetch sequence number for %v", sourceAddress)
		}

		num, ok := sequenceNumbers[sourceAddress]
		if !ok {
			return tx, 0, ErrNoAccount
		}
		return tx, num, ErrNoResults
	}
	return tx, 0, err
}
