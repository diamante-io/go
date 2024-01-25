package ledgerbackend

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/mock"

	"go/support/db"
)

// TrustedLedgerHashStore is used to query ledger data from a trusted source.
// The store should contain ledgers verified by diamcircle-Core, do not use untrusted
// source like history archives.
type TrustedLedgerHashStore interface {
	// GetLedgerHash returns the ledger hash for the given sequence number
	GetLedgerHash(ctx context.Context, seq uint32) (string, bool, error)
	Close() error
}

// auroraDBLedgerHashStore is a TrustedLedgerHashStore which uses aurora's db to look up ledger hashes
type auroraDBLedgerHashStore struct {
	session db.SessionInterface
}

// NewauroraDBLedgerHashStore constructs a new TrustedLedgerHashStore backed by the aurora db
func NewauroraDBLedgerHashStore(session db.SessionInterface) TrustedLedgerHashStore {
	return auroraDBLedgerHashStore{session: session}
}

// GetLedgerHash returns the ledger hash for the given sequence number
func (h auroraDBLedgerHashStore) GetLedgerHash(ctx context.Context, seq uint32) (string, bool, error) {
	sql := sq.Select("hl.ledger_hash").From("history_ledgers hl").
		Limit(1).Where("sequence = ?", seq)

	var hash string
	err := h.session.Get(ctx, &hash, sql)
	if h.session.NoRows(err) {
		return hash, false, nil
	}
	return hash, true, err
}

func (h auroraDBLedgerHashStore) Close() error {
	return h.session.Close()
}

// MockLedgerHashStore is a mock implementation of TrustedLedgerHashStore
type MockLedgerHashStore struct {
	mock.Mock
}

// GetLedgerHash returns the ledger hash for the given sequence number
func (m *MockLedgerHashStore) GetLedgerHash(ctx context.Context, seq uint32) (string, bool, error) {
	args := m.Called(ctx, seq)
	return args.Get(0).(string), args.Get(1).(bool), args.Error(2)
}

func (m *MockLedgerHashStore) Close() error {
	args := m.Called()
	return args.Error(0)
}
