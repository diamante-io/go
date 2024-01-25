package txsub

// This file provides mock implementations for the txsub interfaces
// which are useful in a testing context.
//
// NOTE:  this file is not a test file so that other packages may import
// txsub and use these mocks in their own tests

import (
	"context"
	"database/sql"

	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/stretchr/testify/mock"
)

// MockSubmitter is a test helper that simplements the Submitter interface
type MockSubmitter struct {
	R              SubmissionResult
	WasSubmittedTo bool
}

// Submit implements `txsub.Submitter`
func (sub *MockSubmitter) Submit(ctx context.Context, env string) SubmissionResult {
	sub.WasSubmittedTo = true
	return sub.R
}

type mockDBQ struct {
	mock.Mock
}

func (m *mockDBQ) BeginTx(txOpts *sql.TxOptions) error {
	args := m.Called(txOpts)
	return args.Error(0)
}

func (m *mockDBQ) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockDBQ) NoRows(err error) bool {
	args := m.Called(err)
	return args.Bool(0)
}

func (m *mockDBQ) GetLatestHistoryLedger(ctx context.Context) (uint32, error) {
	args := m.Called()
	return args.Get(0).(uint32), args.Error(1)
}

func (m *mockDBQ) GetSequenceNumbers(ctx context.Context, addresses []string) (map[string]uint64, error) {
	args := m.Called(ctx, addresses)
	return args.Get(0).(map[string]uint64), args.Error(1)
}

func (m *mockDBQ) TransactionsByHashesSinceLedger(ctx context.Context, hashes []string, sinceLedgerSeq uint32) ([]history.Transaction, error) {
	args := m.Called(ctx, hashes, sinceLedgerSeq)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]history.Transaction), args.Error(1)
}

func (m *mockDBQ) TransactionByHash(ctx context.Context, dest interface{}, hash string) error {
	args := m.Called(ctx, dest, hash)
	return args.Error(0)
}
