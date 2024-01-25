package toid

import (
	"errors"
	"fmt"
)

//
// ID represents the total order of Ledgers, Transactions and
// Operations.
//
// Operations within the diamcircle network have a total order, expressed by three
// pieces of information:  the ledger sequence the operation was validated in,
// the order which the operation's containing transaction was applied in
// that ledger, and the index of the operation within that parent transaction.
//
// We express this order by packing those three pieces of information into a
// single signed 64-bit number (we used a signed number for SQL compatibility).
//
// The follow diagram shows this format:
//
//    0                   1                   2                   3
//    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |                    Ledger Sequence Number                     |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |     Transaction Application Order     |       Op Index        |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//
// By component:
//
// Ledger Sequence: 32-bits
//
//   A complete ledger sequence number in which the operation was validated.
//
//   Expressed in network byte order.
//
// Transaction Application Order: 20-bits
//
//   The order that the transaction was applied within the ledger it was
//   validated.  Accommodates up to 1,048,575 transactions in a single ledger.
//
//   Expressed in network byte order.
//
// Operation Index: 12-bits
//
//   The index of the operation within its parent transaction. Accommodates up
//   to 4095 operations per transaction.
//
//   Expressed in network byte order.
//
//
// Note: API Clients should not be interpreting this value.  We will use it
// as an opaque paging token that clients can parrot back to us after having read
// it within a resource to page from the represented position in time.
//
// Note: This does not uniquely identify an object.  Given a ledger, it will
// share its id with its first transaction and the first operation of that
// transaction as well.  Given that this ID is only meant for ordering within a
// single type of object, the sharing of ids across object types seems
// acceptable.
//
type ID struct {
	LedgerSequence   int32
	TransactionOrder int32
	OperationOrder   int32
}

const (
	// LedgerMask is the bitmask to mask out ledger sequences in a
	// TotalOrderID
	LedgerMask = (1 << 32) - 1
	// TransactionMask is the bitmask to mask out transaction indexes
	TransactionMask = (1 << 20) - 1
	// OperationMask is the bitmask to mask out operation indexes
	OperationMask = (1 << 12) - 1

	// LedgerShift is the number of bits to shift an int64 to target the
	// ledger component
	LedgerShift = 32
	// TransactionShift is the number of bits to shift an int64 to
	// target the transaction component
	TransactionShift = 12
	// OperationShift is the number of bits to shift an int64 to target
	// the operation component
	OperationShift = 0
)

// AfterLedger returns a new toid that represents the ledger time _after_ any
// contents (e.g. transactions, operations) that occur within the specified
// ledger.
func AfterLedger(seq int32) *ID {
	return New(seq, TransactionMask, OperationMask)
}

// LedgerRangeInclusive returns inclusive range representation between two
// ledgers inclusive. The second value points at the to+1 ledger so when using
// this value make sure < order is used.
func LedgerRangeInclusive(from, to int32) (int64, int64, error) {
	if from > to {
		return 0, 0, errors.New("Invalid range: from > to")
	}

	if from <= 0 || to <= 0 {
		return 0, 0, errors.New("Invalid range: from or to negative")
	}

	var toidFrom, toidTo int64
	if from == 1 {
		toidFrom = 0
	} else {
		toidFrom = New(from, 0, 0).ToInt64()
	}

	toidTo = New(to+1, 0, 0).ToInt64()

	return toidFrom, toidTo, nil
}

// IncOperationOrder increments the operation order, rolling over to the next
// ledger if overflow occurs.  This allows queries to easily advance a cursor to
// the next operation.
func (id *ID) IncOperationOrder() {
	id.OperationOrder++

	if id.OperationOrder > OperationMask {
		id.OperationOrder = 0
		id.LedgerSequence++
	}
}

// New creates a new total order ID
func New(ledger int32, tx int32, op int32) *ID {
	return &ID{
		LedgerSequence:   ledger,
		TransactionOrder: tx,
		OperationOrder:   op,
	}
}

// ToInt64 converts this struct back into an int64
func (id ID) ToInt64() (result int64) {

	if id.LedgerSequence < 0 {
		panic("invalid ledger sequence")
	}

	if id.TransactionOrder > TransactionMask {
		panic("transaction order overflow")
	}

	if id.OperationOrder > OperationMask {
		panic("operation order overflow")
	}

	result = result | ((int64(id.LedgerSequence) & LedgerMask) << LedgerShift)
	result = result | ((int64(id.TransactionOrder) & TransactionMask) << TransactionShift)
	result = result | ((int64(id.OperationOrder) & OperationMask) << OperationShift)
	return
}

// String returns a string representation of this id
func (id ID) String() string {
	return fmt.Sprintf("%d", id.ToInt64())
}

// Parse parses an int64 into a TotalOrderID struct
func Parse(id int64) (result ID) {
	result.LedgerSequence = int32((id >> LedgerShift) & LedgerMask)
	result.TransactionOrder = int32((id >> TransactionShift) & TransactionMask)
	result.OperationOrder = int32((id >> OperationShift) & OperationMask)

	return
}
