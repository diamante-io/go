package ingest

import (
	"context"
	"fmt"
	stdio "io"
	"testing"

	"go/historyarchive"
	"go/ingest"
	"go/xdr"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockHistoryArchiveAdapter struct {
	mock.Mock
}

func (m *mockHistoryArchiveAdapter) GetLatestLedgerSequence() (uint32, error) {
	args := m.Called()
	return args.Get(0).(uint32), args.Error(1)
}

func (m *mockHistoryArchiveAdapter) BucketListHash(sequence uint32) (xdr.Hash, error) {
	args := m.Called(sequence)
	return args.Get(0).(xdr.Hash), args.Error(1)
}

func (m *mockHistoryArchiveAdapter) GetState(ctx context.Context, sequence uint32) (ingest.ChangeReader, error) {
	args := m.Called(ctx, sequence)
	return args.Get(0).(ingest.ChangeReader), args.Error(1)
}

func TestGetState_Read(t *testing.T) {
	archive, e := getTestArchive()
	if !assert.NoError(t, e) {
		return
	}
	haa := newHistoryArchiveAdapter(archive)

	sr, e := haa.GetState(context.Background(), 21686847)
	if !assert.NoError(t, e) {
		return
	}

	lec, e := sr.Read()
	if !assert.NoError(t, e) {
		return
	}
	assert.NotEqual(t, e, stdio.EOF)

	if !assert.NotNil(t, lec) {
		return
	}
	assert.Equal(t, "GAFBQT4VRORLEVEECUYDQGWNVQ563ZN76LGRJR7T7KDL32EES54UOQST", lec.Post.Data.Account.AccountId.Address())
}

func getTestArchive() (*historyarchive.Archive, error) {
	return historyarchive.Connect(
		fmt.Sprintf("s3://history.diamcircle.org/prd/core-live/core_live_001/"),
		historyarchive.ConnectOptions{
			S3Region:            "eu-west-1",
			UnsignedRequests:    true,
			CheckpointFrequency: 64,
		},
	)
}
