package history

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type MockAccountSignersBatchInsertBuilder struct {
	mock.Mock
}

func (m *MockAccountSignersBatchInsertBuilder) Add(ctx context.Context, signer AccountSigner) error {
	a := m.Called(ctx, signer)
	return a.Error(0)
}

func (m *MockAccountSignersBatchInsertBuilder) Exec(ctx context.Context) error {
	a := m.Called(ctx)
	return a.Error(0)
}
