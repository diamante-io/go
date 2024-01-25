package toid

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ledger = int64(4294967296) // ledger sequence 1
var tx = int64(4096)           // tx index 1
var op = int64(1)              // op index 1

func TestID_ToInt64(t *testing.T) {
	testCases := []struct {
		id          *ID
		expected    int64
		shouldPanic bool
	}{
		// accomodates 12-bits of precision for the operation field
		{
			id:       &ID{0, 0, 1},
			expected: 1,
		},
		{
			id:       &ID{0, 0, 4095},
			expected: 4095,
		},
		{
			id:          &ID{0, 0, 4096},
			shouldPanic: true,
		},
		// accomodates 20-bits of precision for the transaction field
		{
			id:       &ID{0, 1, 0},
			expected: 4096,
		},
		{
			id:       &ID{0, 1048575, 0},
			expected: 4294963200,
		},
		{
			id:          &ID{0, 1048576, 0},
			shouldPanic: true,
		},
		// accomodates 32-bits of precision for the ledger field
		{
			id:       &ID{1, 0, 0},
			expected: 4294967296,
		},
		{
			id:       &ID{math.MaxInt32, 0, 0},
			expected: 9223372032559808512,
		},
		{
			id:          &ID{-1, 0, 0},
			shouldPanic: true,
		},
		{
			id:          &ID{math.MinInt32, 0, 0},
			shouldPanic: true,
		},
		// works as expected
		{
			id:       &ID{1, 1, 1},
			expected: ledger + tx + op,
		},
		{
			id:       &ID{1, 1, 0},
			expected: ledger + tx,
		},
		{
			id:       &ID{1, 0, 1},
			expected: ledger + op,
		},
		{
			id:       &ID{1, 0, 0},
			expected: ledger,
		},
		{
			id:       &ID{0, 1, 0},
			expected: tx,
		},
		{
			id:       &ID{0, 0, 1},
			expected: op,
		},
		{
			id:       &ID{0, 0, 0},
			expected: 0,
		},
	}
	for _, tc := range testCases {
		t.Run("Testing ToInt64", func(t *testing.T) {
			if tc.shouldPanic {
				assert.Panics(t, func() {
					tc.id.ToInt64()
				})
				return
			}
			assert.Equal(t, tc.expected, tc.id.ToInt64())
		})
	}
}

func TestParse(t *testing.T) {
	testCases := []struct {
		parsed   ID
		expected ID
	}{
		{Parse(ledger + tx + op), ID{1, 1, 1}},
		{Parse(ledger + tx), ID{1, 1, 0}},
		{Parse(ledger + op), ID{1, 0, 1}},
		{Parse(ledger), ID{1, 0, 0}},
		{Parse(tx), ID{0, 1, 0}},
		{Parse(op), ID{0, 0, 1}},
	}
	for _, tc := range testCases {
		t.Run("Testing Parse", func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.parsed)
		})
	}
}

// Test InOperationOrder to make sure it rolls over to the next ledger sequence if overflow occurs.
func TestID_IncOperationOrder(t *testing.T) {
	tid := ID{0, 0, 0}
	tid.IncOperationOrder()
	assert.Equal(t, int32(1), tid.OperationOrder)
	tid.OperationOrder = OperationMask
	tid.IncOperationOrder()
	assert.Equal(t, int32(0), tid.OperationOrder)
	assert.Equal(t, int32(1), tid.LedgerSequence)
}

func ExampleParse() {
	toid := Parse(12884910080)
	fmt.Printf("ledger:%d, tx:%d, op:%d", toid.LedgerSequence, toid.TransactionOrder, toid.OperationOrder)
	// Output: ledger:3, tx:2, op:0
}

func TestLedgerRangeInclusive(t *testing.T) {
	testCases := []struct {
		from int32
		to   int32

		fromLedger int32
		toLedger   int32
	}{
		{1, 1, 0, 2},
		{1, 2, 0, 3},
		{2, 2, 2, 3},
		{2, 3, 2, 4},
	}
	for _, tc := range testCases {
		t.Run("Testing TestLedgerRangeInclusive", func(t *testing.T) {
			toidFrom, toidTo, err := LedgerRangeInclusive(tc.from, tc.to)
			assert.NoError(t, err)

			id := Parse(toidFrom)
			assert.Equal(t, tc.fromLedger, id.LedgerSequence)
			assert.Equal(t, int32(0), id.TransactionOrder)
			assert.Equal(t, int32(0), id.OperationOrder)

			id = Parse(toidTo)
			assert.Equal(t, tc.toLedger, id.LedgerSequence)
			assert.Equal(t, int32(0), id.TransactionOrder)
			assert.Equal(t, int32(0), id.OperationOrder)
		})
	}

	_, _, err := LedgerRangeInclusive(2, 1)
	assert.Error(t, err)

	_, _, err = LedgerRangeInclusive(-1, 1)
	assert.Error(t, err)

	_, _, err = LedgerRangeInclusive(-3, -5)
	assert.Error(t, err)
}
