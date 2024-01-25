package ingest

import (
	"go/support/errors"
	"go/xdr"
)

// LedgerTransaction represents the data for a single transaction within a ledger.
type LedgerTransaction struct {
	Index    uint32
	Envelope xdr.TransactionEnvelope
	Result   xdr.TransactionResultPair
	// FeeChanges and UnsafeMeta are low level values, do not use them directly unless
	// you know what you are doing.
	// Use LedgerTransaction.GetChanges() for higher level access to ledger
	// entry changes.
	FeeChanges xdr.LedgerEntryChanges
	UnsafeMeta xdr.TransactionMeta
}

func (t *LedgerTransaction) txInternalError() bool {
	return t.Result.Result.Result.Code == xdr.TransactionResultCodeTxInternalError
}

// GetFeeChanges returns a developer friendly representation of LedgerEntryChanges
// connected to fees.
func (t *LedgerTransaction) GetFeeChanges() []Change {
	return GetChangesFromLedgerEntryChanges(t.FeeChanges)
}

// GetChanges returns a developer friendly representation of LedgerEntryChanges.
// It contains transaction changes and operation changes in that order. If the
// transaction failed with TxInternalError, operations and txChangesAfter are
// omitted. It doesn't support legacy TransactionMeta.V=0.
func (t *LedgerTransaction) GetChanges() ([]Change, error) {
	var changes []Change

	// Transaction meta
	switch t.UnsafeMeta.V {
	case 0:
		return changes, errors.New("TransactionMeta.V=0 not supported")
	case 1:
		v1Meta := t.UnsafeMeta.MustV1()
		txChanges := GetChangesFromLedgerEntryChanges(v1Meta.TxChanges)
		changes = append(changes, txChanges...)

		// Ignore operations meta if txInternalError https://go/issues/2111
		if t.txInternalError() {
			return changes, nil
		}

		for _, operationMeta := range v1Meta.Operations {
			opChanges := GetChangesFromLedgerEntryChanges(
				operationMeta.Changes,
			)
			changes = append(changes, opChanges...)
		}

	case 2:
		v2Meta := t.UnsafeMeta.MustV2()
		txChangesBefore := GetChangesFromLedgerEntryChanges(v2Meta.TxChangesBefore)
		changes = append(changes, txChangesBefore...)

		// Ignore operations meta and txChangesAfter if txInternalError
		// https://go/issues/2111
		if t.txInternalError() {
			return changes, nil
		}

		for _, operationMeta := range v2Meta.Operations {
			opChanges := GetChangesFromLedgerEntryChanges(
				operationMeta.Changes,
			)
			changes = append(changes, opChanges...)
		}

		txChangesAfter := GetChangesFromLedgerEntryChanges(v2Meta.TxChangesAfter)
		changes = append(changes, txChangesAfter...)
	default:
		return changes, errors.New("Unsupported TransactionMeta version")
	}

	return changes, nil
}

// GetOperationChanges returns a developer friendly representation of LedgerEntryChanges.
// It contains only operation changes.
func (t *LedgerTransaction) GetOperationChanges(operationIndex uint32) ([]Change, error) {
	changes := []Change{}

	// Transaction meta
	switch t.UnsafeMeta.V {
	case 0:
		return changes, errors.New("TransactionMeta.V=0 not supported")
	case 1:
		// Ignore operations meta if txInternalError https://go/issues/2111
		if t.txInternalError() {
			return changes, nil
		}

		v1Meta := t.UnsafeMeta.MustV1()
		changes = operationChanges(v1Meta.Operations, operationIndex)
	case 2:
		// Ignore operations meta if txInternalError https://go/issues/2111
		if t.txInternalError() {
			return changes, nil
		}

		v2Meta := t.UnsafeMeta.MustV2()
		changes = operationChanges(v2Meta.Operations, operationIndex)
	default:
		return changes, errors.New("Unsupported TransactionMeta version")
	}

	return changes, nil
}

func operationChanges(ops []xdr.OperationMeta, index uint32) []Change {
	if len(ops) == 0 || int(index) >= len(ops) {
		return []Change{}
	}

	operationMeta := ops[index]
	return GetChangesFromLedgerEntryChanges(
		operationMeta.Changes,
	)
}
