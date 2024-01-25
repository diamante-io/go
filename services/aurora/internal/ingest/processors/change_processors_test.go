package processors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"go/ingest"
	"go/support/errors"
)

func TestStreamReaderError(t *testing.T) {
	tt := assert.New(t)
	ctx := context.Background()

	mockChangeReader := &ingest.MockChangeReader{}
	mockChangeReader.
		On("Read").
		Return(ingest.Change{}, errors.New("transient error")).Once()
	mockChangeProcessor := &MockChangeProcessor{}

	err := StreamChanges(ctx, mockChangeProcessor, mockChangeReader)
	tt.EqualError(err, "could not read transaction: transient error")
}

func TestStreamChangeProcessorError(t *testing.T) {
	tt := assert.New(t)
	ctx := context.Background()

	change := ingest.Change{}
	mockChangeReader := &ingest.MockChangeReader{}
	mockChangeReader.
		On("Read").
		Return(change, nil).Once()

	mockChangeProcessor := &MockChangeProcessor{}
	mockChangeProcessor.
		On(
			"ProcessChange", ctx,
			change,
		).
		Return(errors.New("transient error")).Once()

	err := StreamChanges(ctx, mockChangeProcessor, mockChangeReader)
	tt.EqualError(err, "could not process change: transient error")
}
