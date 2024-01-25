package history

import (
	"context"

	"github.com/stretchr/testify/mock"

	"go/xdr"
)

type MockQLedgers struct {
	mock.Mock
}

func (m *MockQLedgers) InsertLedger(ctx context.Context,
	ledger xdr.LedgerHeaderHistoryEntry,
	successTxsCount int,
	failedTxsCount int,
	opCount int,
	txSetOpCount int,
	ingestVersion int,
) (int64, error) {
	a := m.Called(ctx, ledger, successTxsCount, failedTxsCount, opCount, txSetOpCount, ingestVersion)
	return a.Get(0).(int64), a.Error(1)
}
