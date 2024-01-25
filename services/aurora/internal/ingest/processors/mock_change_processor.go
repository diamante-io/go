package processors

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/diamcircle/go/ingest"
)

var _ ChangeProcessor = (*MockChangeProcessor)(nil)

type MockChangeProcessor struct {
	mock.Mock
}

func (m *MockChangeProcessor) ProcessChange(ctx context.Context, change ingest.Change) error {
	args := m.Called(ctx, change)
	return args.Error(0)
}
