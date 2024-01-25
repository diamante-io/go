package ledgerbackend

import (
	"context"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/diamcircle/go/historyarchive"
	"github.com/diamcircle/go/network"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/xdr"
)

// TODO: test frame decoding
// TODO: test from static base64-encoded data

type diamcircleCoreRunnerMock struct {
	mock.Mock
}

func (m *diamcircleCoreRunnerMock) context() context.Context {
	a := m.Called()
	return a.Get(0).(context.Context)
}

func (m *diamcircleCoreRunnerMock) catchup(from, to uint32) error {
	a := m.Called(from, to)
	return a.Error(0)
}

func (m *diamcircleCoreRunnerMock) runFrom(from uint32, hash string) error {
	a := m.Called(from, hash)
	return a.Error(0)
}

func (m *diamcircleCoreRunnerMock) getMetaPipe() <-chan metaResult {
	a := m.Called()
	return a.Get(0).(<-chan metaResult)
}

func (m *diamcircleCoreRunnerMock) getProcessExitError() (bool, error) {
	a := m.Called()
	return a.Bool(0), a.Error(1)
}

func (m *diamcircleCoreRunnerMock) close() error {
	a := m.Called()
	return a.Error(0)
}

func buildLedgerCloseMeta(header testLedgerHeader) xdr.LedgerCloseMeta {
	opResults := []xdr.OperationResult{}
	opMeta := []xdr.OperationMeta{}

	tmpHash, _ := hex.DecodeString("cde54da3901f5b9c0331d24fbb06ac9c5c5de76de9fb2d4a7b86c09e46f11d8c")
	var hash [32]byte
	copy(hash[:], tmpHash)

	var ledgerHash [32]byte
	if header.hash != "" {
		tmpHash, err := hex.DecodeString(header.hash)
		if err != nil {
			panic(err)
		}
		copy(ledgerHash[:], tmpHash)
	}

	var previousLedgerHash [32]byte
	if header.hash != "" {
		tmpHash, err := hex.DecodeString(header.previousLedgerHash)
		if err != nil {
			panic(err)
		}
		copy(previousLedgerHash[:], tmpHash)
	}

	source := xdr.MustAddress("GAEJJMDDCRYF752PKIJICUVL7MROJBNXDV2ZB455T7BAFHU2LCLSE2LW")
	return xdr.LedgerCloseMeta{
		V: 0,
		V0: &xdr.LedgerCloseMetaV0{
			LedgerHeader: xdr.LedgerHeaderHistoryEntry{
				Hash: ledgerHash,
				Header: xdr.LedgerHeader{
					LedgerSeq:          xdr.Uint32(header.sequence),
					PreviousLedgerHash: previousLedgerHash,
				},
			},
			TxSet: xdr.TransactionSet{
				Txs: []xdr.TransactionEnvelope{
					{
						Type: xdr.EnvelopeTypeEnvelopeTypeTx,
						V1: &xdr.TransactionV1Envelope{
							Tx: xdr.Transaction{
								SourceAccount: source.ToMuxedAccount(),
								Fee:           xdr.Uint32(header.sequence),
							},
						},
					},
				},
			},
			TxProcessing: []xdr.TransactionResultMeta{
				{
					Result: xdr.TransactionResultPair{
						TransactionHash: xdr.Hash(hash),
						Result: xdr.TransactionResult{
							FeeCharged: xdr.Int64(header.sequence),
							Result: xdr.TransactionResultResult{
								Code:    xdr.TransactionResultCodeTxSuccess,
								Results: &opResults,
							},
						},
					},
					TxApplyProcessing: xdr.TransactionMeta{
						Operations: &opMeta,
					},
				},
			},
		},
	}

}

type testLedgerHeader struct {
	sequence           uint32
	hash               string
	previousLedgerHash string
}

func TestCaptiveNew(t *testing.T) {
	executablePath := "/etc/diamcircle-core"
	networkPassphrase := network.PublicNetworkPassphrase
	historyURLs := []string{"http://history.diamcircle.org/prd/core-live/core_live_001"}

	captiveDiamcircleCore, err := NewCaptive(
		CaptiveCoreConfig{
			BinaryPath:         executablePath,
			NetworkPassphrase:  networkPassphrase,
			HistoryArchiveURLs: historyURLs,
		},
	)

	assert.NoError(t, err)
	assert.Equal(t, uint32(0), captiveDiamcircleCore.nextLedger)
	assert.NotNil(t, captiveDiamcircleCore.archive)
}

func TestCaptivePrepareRange(t *testing.T) {
	metaChan := make(chan metaResult, 100)

	// Core will actually start with the last checkpoint before the from ledger
	// and then rewind to the `from` ledger.
	for i := 64; i <= 100; i++ {
		meta := buildLedgerCloseMeta(testLedgerHeader{sequence: uint32(i)})
		metaChan <- metaResult{
			LedgerCloseMeta: &meta,
		}
	}

	ctx := context.Background()
	mockRunner := &diamcircleCoreRunnerMock{}
	mockRunner.On("catchup", uint32(100), uint32(200)).Return(nil).Once()
	mockRunner.On("getMetaPipe").Return((<-chan metaResult)(metaChan))
	mockRunner.On("context").Return(ctx)

	mockArchive := &historyarchive.MockArchive{}
	mockArchive.
		On("GetRootHAS").
		Return(historyarchive.HistoryArchiveState{
			CurrentLedger: uint32(200),
		}, nil)

	cancelCalled := false
	captiveBackend := CaptiveDiamcircleCore{
		archive: mockArchive,
		diamcircleCoreRunnerFactory: func(_ diamcircleCoreRunnerMode) (diamcircleCoreRunnerInterface, error) {
			return mockRunner, nil
		},
		checkpointManager: historyarchive.NewCheckpointManager(64),
		cancel: context.CancelFunc(func() {
			cancelCalled = true
		}),
	}

	err := captiveBackend.PrepareRange(ctx, BoundedRange(100, 200))
	assert.NoError(t, err)
	mockRunner.On("close").Return(nil).Once()
	err = captiveBackend.Close()
	assert.NoError(t, err)
	assert.True(t, cancelCalled)
	mockRunner.AssertExpectations(t)
	mockArchive.AssertExpectations(t)
}

func TestCaptivePrepareRangeCrash(t *testing.T) {
	metaChan := make(chan metaResult)
	close(metaChan)
	ctx := context.Background()
	mockRunner := &diamcircleCoreRunnerMock{}
	mockRunner.On("catchup", uint32(100), uint32(200)).Return(nil).Once()
	mockRunner.On("getProcessExitError").Return(true, errors.New("exit code -1"))
	mockRunner.On("getMetaPipe").Return((<-chan metaResult)(metaChan))
	mockRunner.On("close").Return(nil).Once()
	mockRunner.On("context").Return(ctx)

	mockArchive := &historyarchive.MockArchive{}
	mockArchive.
		On("GetRootHAS").
		Return(historyarchive.HistoryArchiveState{
			CurrentLedger: uint32(200),
		}, nil)

	captiveBackend := CaptiveDiamcircleCore{
		archive: mockArchive,
		diamcircleCoreRunnerFactory: func(_ diamcircleCoreRunnerMode) (diamcircleCoreRunnerInterface, error) {
			return mockRunner, nil
		},
		checkpointManager: historyarchive.NewCheckpointManager(64),
	}

	err := captiveBackend.PrepareRange(ctx, BoundedRange(100, 200))
	assert.EqualError(t, err, "Error fast-forwarding to 100: diamcircle core exited unexpectedly: exit code -1")
	mockRunner.AssertExpectations(t)
	mockArchive.AssertExpectations(t)
}

func TestCaptivePrepareRangeTerminated(t *testing.T) {
	metaChan := make(chan metaResult, 100)

	// Core will actually start with the last checkpoint before the from ledger
	// and then rewind to the `from` ledger.
	for i := 64; i <= 100; i++ {
		meta := buildLedgerCloseMeta(testLedgerHeader{sequence: uint32(i)})
		metaChan <- metaResult{
			LedgerCloseMeta: &meta,
		}
	}
	close(metaChan)
	ctx := context.Background()
	mockRunner := &diamcircleCoreRunnerMock{}
	mockRunner.On("catchup", uint32(100), uint32(200)).Return(nil).Once()
	mockRunner.On("getMetaPipe").Return((<-chan metaResult)(metaChan))
	mockRunner.On("context").Return(ctx)

	mockArchive := &historyarchive.MockArchive{}
	mockArchive.
		On("GetRootHAS").
		Return(historyarchive.HistoryArchiveState{
			CurrentLedger: uint32(200),
		}, nil)

	captiveBackend := CaptiveDiamcircleCore{
		archive: mockArchive,
		diamcircleCoreRunnerFactory: func(_ diamcircleCoreRunnerMode) (diamcircleCoreRunnerInterface, error) {
			return mockRunner, nil
		},
		checkpointManager: historyarchive.NewCheckpointManager(64),
	}

	err := captiveBackend.PrepareRange(ctx, BoundedRange(100, 200))
	assert.NoError(t, err)
	mockRunner.AssertExpectations(t)
	mockArchive.AssertExpectations(t)
}

func TestCaptivePrepareRangeCloseNotFullyTerminated(t *testing.T) {
	metaChan := make(chan metaResult, 100)
	for i := 64; i <= 100; i++ {
		meta := buildLedgerCloseMeta(testLedgerHeader{sequence: uint32(i)})
		metaChan <- metaResult{
			LedgerCloseMeta: &meta,
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	mockRunner := &diamcircleCoreRunnerMock{}
	mockRunner.On("catchup", uint32(100), uint32(200)).Return(nil).Once()
	mockRunner.On("getMetaPipe").Return((<-chan metaResult)(metaChan))
	mockRunner.On("context").Return(ctx)

	mockArchive := &historyarchive.MockArchive{}
	mockArchive.
		On("GetRootHAS").
		Return(historyarchive.HistoryArchiveState{
			CurrentLedger: uint32(200),
		}, nil)

	captiveBackend := CaptiveDiamcircleCore{
		archive: mockArchive,
		diamcircleCoreRunnerFactory: func(_ diamcircleCoreRunnerMode) (diamcircleCoreRunnerInterface, error) {
			return mockRunner, nil
		},
		checkpointManager: historyarchive.NewCheckpointManager(64),
	}

	err := captiveBackend.PrepareRange(ctx, BoundedRange(100, 200))
	assert.NoError(t, err)

	// Simulates a long (but graceful) shutdown...
	cancel()

	err = captiveBackend.PrepareRange(ctx, BoundedRange(100, 200))
	assert.NoError(t, err)

	mockRunner.AssertExpectations(t)
	mockArchive.AssertExpectations(t)
}

func TestCaptivePrepareRange_ErrClosingSession(t *testing.T) {
	ctx := context.Background()
	mockRunner := &diamcircleCoreRunnerMock{}
	mockRunner.On("close").Return(fmt.Errorf("transient error"))

	captiveBackend := CaptiveDiamcircleCore{
		nextLedger:        300,
		diamcircleCoreRunner: mockRunner,
	}

	err := captiveBackend.PrepareRange(ctx, BoundedRange(100, 200))
	assert.EqualError(t, err, "error starting prepare range: error closing existing session: transient error")

	err = captiveBackend.PrepareRange(ctx, UnboundedRange(64))
	assert.EqualError(t, err, "error starting prepare range: error closing existing session: transient error")

	mockRunner.AssertExpectations(t)
}

func TestCaptivePrepareRange_ErrGettingRootHAS(t *testing.T) {
	ctx := context.Background()
	mockArchive := &historyarchive.MockArchive{}
	mockArchive.
		On("GetRootHAS").
		Return(historyarchive.HistoryArchiveState{}, errors.New("transient error"))

	captiveBackend := CaptiveDiamcircleCore{
		archive: mockArchive,
	}

	err := captiveBackend.PrepareRange(ctx, BoundedRange(100, 200))
	assert.EqualError(t, err, "error starting prepare range: opening subprocess: error getting latest checkpoint sequence: error getting root HAS: transient error")

	err = captiveBackend.PrepareRange(ctx, UnboundedRange(100))
	assert.EqualError(t, err, "error starting prepare range: opening subprocess: error getting latest checkpoint sequence: error getting root HAS: transient error")

	mockArchive.AssertExpectations(t)
}

func TestCaptivePrepareRange_FromIsAheadOfRootHAS(t *testing.T) {
	ctx := context.Background()
	mockArchive := &historyarchive.MockArchive{}
	mockArchive.
		On("GetRootHAS").
		Return(historyarchive.HistoryArchiveState{
			CurrentLedger: uint32(64),
		}, nil)

	captiveBackend := CaptiveDiamcircleCore{
		archive: mockArchive,
	}

	err := captiveBackend.PrepareRange(ctx, BoundedRange(100, 200))
	assert.EqualError(t, err, "error starting prepare range: opening subprocess: from sequence: 100 is greater than max available in history archives: 64")

	err = captiveBackend.PrepareRange(ctx, UnboundedRange(100))
	assert.EqualError(t, err, "error starting prepare range: opening subprocess: trying to start online mode too far (latest checkpoint=64), only two checkpoints in the future allowed")

	mockArchive.AssertExpectations(t)
}

func TestCaptivePrepareRange_ToIsAheadOfRootHAS(t *testing.T) {
	mockRunner := &diamcircleCoreRunnerMock{}
	mockArchive := &historyarchive.MockArchive{}
	mockArchive.
		On("GetRootHAS").
		Return(historyarchive.HistoryArchiveState{
			CurrentLedger: uint32(192),
		}, nil)

	captiveBackend := CaptiveDiamcircleCore{
		archive: mockArchive,
		diamcircleCoreRunnerFactory: func(_ diamcircleCoreRunnerMode) (diamcircleCoreRunnerInterface, error) {
			return mockRunner, nil
		},
		checkpointManager: historyarchive.NewCheckpointManager(64),
	}

	err := captiveBackend.PrepareRange(context.Background(), BoundedRange(100, 200))
	assert.EqualError(t, err, "error starting prepare range: opening subprocess: to sequence: 200 is greater than max available in history archives: 192")

	mockArchive.AssertExpectations(t)
	mockRunner.AssertExpectations(t)
}

func TestCaptivePrepareRange_ErrCatchup(t *testing.T) {
	mockRunner := &diamcircleCoreRunnerMock{}
	mockRunner.On("catchup", uint32(100), uint32(192)).Return(errors.New("transient error")).Once()
	mockRunner.On("close").Return(nil).Once()

	mockArchive := &historyarchive.MockArchive{}
	mockArchive.
		On("GetRootHAS").
		Return(historyarchive.HistoryArchiveState{
			CurrentLedger: uint32(192),
		}, nil)

	ctx := context.Background()
	cancelCalled := false
	captiveBackend := CaptiveDiamcircleCore{
		archive: mockArchive,
		diamcircleCoreRunnerFactory: func(_ diamcircleCoreRunnerMode) (diamcircleCoreRunnerInterface, error) {
			return mockRunner, nil
		},
		cancel: context.CancelFunc(func() {
			cancelCalled = true
		}),
	}

	err := captiveBackend.PrepareRange(ctx, BoundedRange(100, 192))
	assert.EqualError(t, err, "error starting prepare range: opening subprocess: error running diamcircle-core: transient error")

	// make sure we can Close without errors
	assert.NoError(t, captiveBackend.Close())
	assert.True(t, cancelCalled)

	mockArchive.AssertExpectations(t)
	mockRunner.AssertExpectations(t)
}

func TestCaptivePrepareRangeUnboundedRange_ErrRunFrom(t *testing.T) {
	mockRunner := &diamcircleCoreRunnerMock{}
	mockRunner.On("runFrom", uint32(127), "0000000000000000000000000000000000000000000000000000000000000000").Return(errors.New("transient error")).Once()
	mockRunner.On("close").Return(nil).Once()

	mockArchive := &historyarchive.MockArchive{}
	mockArchive.
		On("GetRootHAS").
		Return(historyarchive.HistoryArchiveState{
			CurrentLedger: uint32(127),
		}, nil)

	mockArchive.
		On("GetLedgerHeader", uint32(128)).
		Return(xdr.LedgerHeaderHistoryEntry{}, nil)

	ctx := context.Background()
	cancelCalled := false
	captiveBackend := CaptiveDiamcircleCore{
		archive: mockArchive,
		diamcircleCoreRunnerFactory: func(_ diamcircleCoreRunnerMode) (diamcircleCoreRunnerInterface, error) {
			return mockRunner, nil
		},
		checkpointManager: historyarchive.NewCheckpointManager(64),
		cancel: context.CancelFunc(func() {
			cancelCalled = true
		}),
	}

	err := captiveBackend.PrepareRange(ctx, UnboundedRange(128))
	assert.EqualError(t, err, "error starting prepare range: opening subprocess: error running diamcircle-core: transient error")

	// make sure we can Close without errors
	assert.NoError(t, captiveBackend.Close())
	assert.True(t, cancelCalled)

	mockArchive.AssertExpectations(t)
	mockRunner.AssertExpectations(t)
}

func TestCaptivePrepareRangeUnboundedRange_ReuseSession(t *testing.T) {
	metaChan := make(chan metaResult, 100)

	// Core will actually start with the last checkpoint before the from ledger
	// and then rewind to the `from` ledger.
	for i := 2; i <= 65; i++ {
		meta := buildLedgerCloseMeta(testLedgerHeader{sequence: uint32(i)})
		metaChan <- metaResult{
			LedgerCloseMeta: &meta,
		}
	}

	ctx := context.Background()
	mockRunner := &diamcircleCoreRunnerMock{}
	mockRunner.On("runFrom", uint32(64), "0000000000000000000000000000000000000000000000000000000000000000").Return(nil).Once()
	mockRunner.On("getMetaPipe").Return((<-chan metaResult)(metaChan))
	mockRunner.On("context").Return(ctx)

	mockArchive := &historyarchive.MockArchive{}
	mockArchive.
		On("GetRootHAS").
		Return(historyarchive.HistoryArchiveState{
			CurrentLedger: uint32(129),
		}, nil)
	mockArchive.
		On("GetLedgerHeader", uint32(65)).
		Return(xdr.LedgerHeaderHistoryEntry{}, nil)

	captiveBackend := CaptiveDiamcircleCore{
		archive: mockArchive,
		diamcircleCoreRunnerFactory: func(_ diamcircleCoreRunnerMode) (diamcircleCoreRunnerInterface, error) {
			return mockRunner, nil
		},
		checkpointManager: historyarchive.NewCheckpointManager(64),
	}

	err := captiveBackend.PrepareRange(ctx, UnboundedRange(65))
	assert.NoError(t, err)

	captiveBackend.nextLedger = 64
	err = captiveBackend.PrepareRange(ctx, UnboundedRange(65))
	assert.NoError(t, err)

	mockArchive.AssertExpectations(t)
	mockRunner.AssertExpectations(t)
}

func TestGetLatestLedgerSequence(t *testing.T) {
	metaChan := make(chan metaResult, 300)

	// Core will actually start with the last checkpoint before the `from` ledger
	// and then rewind to the `from` ledger.
	for i := 2; i <= 200; i++ {
		meta := buildLedgerCloseMeta(testLedgerHeader{sequence: uint32(i)})
		metaChan <- metaResult{
			LedgerCloseMeta: &meta,
		}
	}

	ctx := context.Background()
	mockRunner := &diamcircleCoreRunnerMock{}
	mockRunner.On("runFrom", uint32(63), "0000000000000000000000000000000000000000000000000000000000000000").Return(nil).Once()
	mockRunner.On("getMetaPipe").Return((<-chan metaResult)(metaChan))
	mockRunner.On("context").Return(ctx)

	mockArchive := &historyarchive.MockArchive{}
	mockArchive.
		On("GetRootHAS").
		Return(historyarchive.HistoryArchiveState{
			CurrentLedger: uint32(200),
		}, nil)

	mockArchive.
		On("GetLedgerHeader", uint32(64)).
		Return(xdr.LedgerHeaderHistoryEntry{}, nil)

	captiveBackend := CaptiveDiamcircleCore{
		archive: mockArchive,
		diamcircleCoreRunnerFactory: func(_ diamcircleCoreRunnerMode) (diamcircleCoreRunnerInterface, error) {
			return mockRunner, nil
		},
		checkpointManager: historyarchive.NewCheckpointManager(64),
	}

	err := captiveBackend.PrepareRange(ctx, UnboundedRange(64))
	assert.NoError(t, err)

	latest, err := captiveBackend.GetLatestLedgerSequence(ctx)
	assert.NoError(t, err)
	assert.Equal(t, uint32(200), latest)

	mockArchive.AssertExpectations(t)
	mockRunner.AssertExpectations(t)
}

func TestCaptiveGetLedger(t *testing.T) {
	tt := assert.New(t)
	metaChan := make(chan metaResult, 300)

	for i := 64; i <= 66; i++ {
		meta := buildLedgerCloseMeta(testLedgerHeader{sequence: uint32(i)})
		metaChan <- metaResult{
			LedgerCloseMeta: &meta,
		}
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	mockRunner := &diamcircleCoreRunnerMock{}
	mockRunner.On("catchup", uint32(65), uint32(66)).Return(nil)
	mockRunner.On("getMetaPipe").Return((<-chan metaResult)(metaChan))
	mockRunner.On("context").Return(ctx)

	mockArchive := &historyarchive.MockArchive{}
	mockArchive.
		On("GetRootHAS").
		Return(historyarchive.HistoryArchiveState{
			CurrentLedger: uint32(200),
		}, nil)

	captiveBackend := CaptiveDiamcircleCore{
		archive: mockArchive,
		diamcircleCoreRunnerFactory: func(_ diamcircleCoreRunnerMode) (diamcircleCoreRunnerInterface, error) {
			return mockRunner, nil
		},
		checkpointManager: historyarchive.NewCheckpointManager(64),
	}

	// requires PrepareRange
	_, err := captiveBackend.GetLedger(ctx, 64)
	tt.EqualError(err, "session is not prepared, call PrepareRange first")

	ledgerRange := BoundedRange(65, 66)
	tt.False(captiveBackend.isPrepared(ledgerRange), "core is not prepared until explicitly prepared")
	tt.False(captiveBackend.isClosed())
	err = captiveBackend.PrepareRange(ctx, ledgerRange)
	assert.NoError(t, err)

	tt.True(captiveBackend.isPrepared(ledgerRange))
	tt.False(captiveBackend.isClosed())

	_, err = captiveBackend.GetLedger(ctx, 64)
	tt.Error(err, "requested ledger 64 is behind the captive core stream (expected=66)")

	// reads value from buffer
	meta, err := captiveBackend.GetLedger(ctx, 65)
	tt.NoError(err)
	tt.Equal(xdr.Uint32(65), meta.V0.LedgerHeader.Header.LedgerSeq)

	// reads value from cachedMeta
	cachedMeta, err := captiveBackend.GetLedger(ctx, 65)
	tt.NoError(err)
	tt.Equal(meta, cachedMeta)

	// next sequence number didn't get consumed
	tt.Equal(uint32(66), captiveBackend.nextLedger)

	mockRunner.On("close").Return(nil).Run(func(args mock.Arguments) {
		cancel()
	}).Once()

	_, err = captiveBackend.GetLedger(ctx, 66)
	tt.NoError(err)

	tt.False(captiveBackend.isPrepared(ledgerRange))
	tt.False(captiveBackend.isClosed())
	_, err = captiveBackend.GetLedger(ctx, 66)
	tt.NoError(err)

	// core is not closed unless it's explicitly closed
	tt.False(captiveBackend.isClosed())

	mockArchive.AssertExpectations(t)
	mockRunner.AssertExpectations(t)
}

// TestCaptiveGetLedgerCacheLatestLedger test the following case:
// 1. Prepare Unbounded range.
// 2. GetLedger that is still not in the buffer.
// 3. Get latest ledger in the buffer using GetLedger.
//
// Before 3d97762 this test failed because cachedMeta was only updated when
// the ledger with a requested sequence was reached while streaming meta.
//
// TODO: Not sure this test is really valid or worth it anymore, now that GetLedger is always blocking.
func TestCaptiveGetLedgerCacheLatestLedger(t *testing.T) {
	tt := assert.New(t)
	metaChan := make(chan metaResult, 300)

	for i := 2; i <= 67; i++ {
		meta := buildLedgerCloseMeta(testLedgerHeader{sequence: uint32(i)})
		metaChan <- metaResult{
			LedgerCloseMeta: &meta,
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mockRunner := &diamcircleCoreRunnerMock{}
	mockRunner.On("runFrom", uint32(65), "0101010100000000000000000000000000000000000000000000000000000000").Return(nil).Once()
	mockRunner.On("getMetaPipe").Return((<-chan metaResult)(metaChan))
	mockRunner.On("context").Return(ctx)

	mockArchive := &historyarchive.MockArchive{}
	mockArchive.
		On("GetRootHAS").
		Return(historyarchive.HistoryArchiveState{
			CurrentLedger: uint32(200),
		}, nil)

	mockArchive.
		On("GetLedgerHeader", uint32(66)).
		Return(xdr.LedgerHeaderHistoryEntry{
			Header: xdr.LedgerHeader{
				PreviousLedgerHash: xdr.Hash{1, 1, 1, 1},
			},
		}, nil).Once()

	captiveBackend := CaptiveDiamcircleCore{
		archive: mockArchive,
		diamcircleCoreRunnerFactory: func(_ diamcircleCoreRunnerMode) (diamcircleCoreRunnerInterface, error) {
			return mockRunner, nil
		},
		checkpointManager: historyarchive.NewCheckpointManager(64),
	}

	err := captiveBackend.PrepareRange(ctx, UnboundedRange(66))
	assert.NoError(t, err)

	// found, _, err := captiveBackend.GetLedger(ctx, 68)
	// tt.NoError(err)
	// tt.False(found)
	// tt.Equal(uint32(67), captiveBackend.cachedMeta.LedgerSequence())
	// tt.Equal(uint32(68), captiveBackend.nextLedger)

	meta, err := captiveBackend.GetLedger(ctx, 67)
	tt.NoError(err)
	tt.Equal(uint32(67), meta.LedgerSequence())

	mockArchive.AssertExpectations(t)
	mockRunner.AssertExpectations(t)
}

func TestCaptiveGetLedger_NextLedgerIsDifferentToLedgerFromBuffer(t *testing.T) {
	metaChan := make(chan metaResult, 100)

	for i := 64; i <= 65; i++ {
		meta := buildLedgerCloseMeta(testLedgerHeader{sequence: uint32(i)})
		metaChan <- metaResult{
			LedgerCloseMeta: &meta,
		}
	}
	{
		meta := buildLedgerCloseMeta(testLedgerHeader{sequence: uint32(68)})
		metaChan <- metaResult{
			LedgerCloseMeta: &meta,
		}
	}

	ctx := context.Background()
	mockRunner := &diamcircleCoreRunnerMock{}
	mockRunner.On("catchup", uint32(65), uint32(66)).Return(nil)
	mockRunner.On("getMetaPipe").Return((<-chan metaResult)(metaChan))
	mockRunner.On("context").Return(ctx)
	mockRunner.On("close").Return(nil)

	mockArchive := &historyarchive.MockArchive{}
	mockArchive.
		On("GetRootHAS").
		Return(historyarchive.HistoryArchiveState{
			CurrentLedger: uint32(200),
		}, nil)

	captiveBackend := CaptiveDiamcircleCore{
		archive: mockArchive,
		diamcircleCoreRunnerFactory: func(_ diamcircleCoreRunnerMode) (diamcircleCoreRunnerInterface, error) {
			return mockRunner, nil
		},
		checkpointManager: historyarchive.NewCheckpointManager(64),
	}

	err := captiveBackend.PrepareRange(ctx, BoundedRange(65, 66))
	assert.NoError(t, err)

	_, err = captiveBackend.GetLedger(ctx, 66)
	assert.EqualError(t, err, "unexpected ledger sequence (expected=66 actual=68)")

	// TODO assertions should work - to be fixed in a separate PR.
	// _, err = captiveBackend.GetLedger(ctx, 66)
	// assert.EqualError(t, err, "session is closed, call PrepareRange first")

	mockArchive.AssertExpectations(t)
	mockRunner.AssertExpectations(t)
}

func TestCaptiveGetLedger_NextLedger0RangeFromIsSmallerThanLedgerFromBuffer(t *testing.T) {
	metaChan := make(chan metaResult, 100)

	for i := 66; i <= 66; i++ {
		meta := buildLedgerCloseMeta(testLedgerHeader{sequence: uint32(i)})
		metaChan <- metaResult{
			LedgerCloseMeta: &meta,
		}
	}

	ctx := context.Background()
	mockRunner := &diamcircleCoreRunnerMock{}
	mockRunner.On("runFrom", uint32(64), mock.Anything).Return(nil)
	mockRunner.On("getMetaPipe").Return((<-chan metaResult)(metaChan))
	mockRunner.On("context").Return(ctx)
	mockRunner.On("close").Return(nil)

	mockArchive := &historyarchive.MockArchive{}
	mockArchive.
		On("GetRootHAS").
		Return(historyarchive.HistoryArchiveState{
			CurrentLedger: uint32(200),
		}, nil)

	mockArchive.
		On("GetLedgerHeader", uint32(65)).
		Return(xdr.LedgerHeaderHistoryEntry{}, nil)

	captiveBackend := CaptiveDiamcircleCore{
		archive: mockArchive,
		diamcircleCoreRunnerFactory: func(_ diamcircleCoreRunnerMode) (diamcircleCoreRunnerInterface, error) {
			return mockRunner, nil
		},
		checkpointManager: historyarchive.NewCheckpointManager(64),
	}

	err := captiveBackend.PrepareRange(ctx, UnboundedRange(65))
	assert.EqualError(t, err, "Error fast-forwarding to 65: unexpected ledger sequence (expected=<=65 actual=66)")

	// TODO assertions should work - to be fixed in a separate PR.
	// prepared, err := captiveBackend.IsPrepared(ctx, UnboundedRange(65))
	// assert.NoError(t, err)
	// assert.False(t, prepared)

	mockArchive.AssertExpectations(t)
	mockRunner.AssertExpectations(t)
}

func TestCaptiveDiamcircleCore_PrepareRangeAfterClose(t *testing.T) {
	ctx := context.Background()
	executablePath := "/etc/diamcircle-core"
	networkPassphrase := network.PublicNetworkPassphrase
	historyURLs := []string{"http://localhost"}

	captiveCoreToml, err := NewCaptiveCoreToml(CaptiveCoreTomlParams{})
	assert.NoError(t, err)

	captiveDiamcircleCore, err := NewCaptive(
		CaptiveCoreConfig{
			BinaryPath:         executablePath,
			NetworkPassphrase:  networkPassphrase,
			HistoryArchiveURLs: historyURLs,
			Toml:               captiveCoreToml,
		},
	)
	assert.NoError(t, err)

	assert.NoError(t, captiveDiamcircleCore.Close())

	assert.EqualError(
		t,
		captiveDiamcircleCore.PrepareRange(ctx, BoundedRange(65, 66)),
		"error starting prepare range: opening subprocess: error getting latest checkpoint sequence: "+
			"error getting root HAS: Get \"http://localhost/.well-known/diamcircle-history.json\": context canceled",
	)

	// even if the request to fetch the latest checkpoint succeeds, we should fail at creating the subprocess
	mockArchive := &historyarchive.MockArchive{}
	mockArchive.
		On("GetRootHAS").
		Return(historyarchive.HistoryArchiveState{
			CurrentLedger: uint32(200),
		}, nil)
	captiveDiamcircleCore.archive = mockArchive
	assert.EqualError(
		t,
		captiveDiamcircleCore.PrepareRange(ctx, BoundedRange(65, 66)),
		"error starting prepare range: opening subprocess: error running diamcircle-core: context canceled",
	)
	mockArchive.AssertExpectations(t)
}

func TestCaptiveGetLedger_ErrReadingMetaResult(t *testing.T) {
	tt := assert.New(t)
	metaChan := make(chan metaResult, 100)

	for i := 64; i <= 65; i++ {
		meta := buildLedgerCloseMeta(testLedgerHeader{sequence: uint32(i)})
		metaChan <- metaResult{
			LedgerCloseMeta: &meta,
		}
	}
	metaChan <- metaResult{
		err: fmt.Errorf("unmarshalling error"),
	}

	ctx := context.Background()
	mockRunner := &diamcircleCoreRunnerMock{}
	mockRunner.On("catchup", uint32(65), uint32(66)).Return(nil)
	mockRunner.On("getMetaPipe").Return((<-chan metaResult)(metaChan))
	ctx, cancel := context.WithCancel(ctx)
	mockRunner.On("context").Return(ctx)
	mockRunner.On("close").Return(nil).Run(func(args mock.Arguments) {
		cancel()
	}).Once()

	// even if the request to fetch the latest checkpoint succeeds, we should fail at creating the subprocess
	mockArchive := &historyarchive.MockArchive{}
	mockArchive.
		On("GetRootHAS").
		Return(historyarchive.HistoryArchiveState{
			CurrentLedger: uint32(200),
		}, nil)

	captiveBackend := CaptiveDiamcircleCore{
		archive: mockArchive,
		diamcircleCoreRunnerFactory: func(_ diamcircleCoreRunnerMode) (diamcircleCoreRunnerInterface, error) {
			return mockRunner, nil
		},
		checkpointManager: historyarchive.NewCheckpointManager(64),
	}

	err := captiveBackend.PrepareRange(ctx, BoundedRange(65, 66))
	assert.NoError(t, err)

	meta, err := captiveBackend.GetLedger(ctx, 65)
	tt.NoError(err)
	tt.Equal(xdr.Uint32(65), meta.V0.LedgerHeader.Header.LedgerSeq)

	tt.False(captiveBackend.isClosed())

	// try reading from an empty buffer
	_, err = captiveBackend.GetLedger(ctx, 66)
	tt.EqualError(err, "unmarshalling error")

	// not closed even if there is an error getting ledger
	tt.False(captiveBackend.isClosed())

	mockArchive.AssertExpectations(t)
	mockRunner.AssertExpectations(t)
}

func TestCaptiveGetLedger_ErrClosingAfterLastLedger(t *testing.T) {
	tt := assert.New(t)
	metaChan := make(chan metaResult, 100)

	for i := 64; i <= 66; i++ {
		meta := buildLedgerCloseMeta(testLedgerHeader{sequence: uint32(i)})
		metaChan <- metaResult{
			LedgerCloseMeta: &meta,
		}
	}

	ctx := context.Background()
	mockRunner := &diamcircleCoreRunnerMock{}
	mockRunner.On("catchup", uint32(65), uint32(66)).Return(nil)
	mockRunner.On("getMetaPipe").Return((<-chan metaResult)(metaChan))
	mockRunner.On("context").Return(ctx)
	mockRunner.On("close").Return(fmt.Errorf("transient error")).Once()

	mockArchive := &historyarchive.MockArchive{}
	mockArchive.
		On("GetRootHAS").
		Return(historyarchive.HistoryArchiveState{
			CurrentLedger: uint32(200),
		}, nil)

	captiveBackend := CaptiveDiamcircleCore{
		archive: mockArchive,
		diamcircleCoreRunnerFactory: func(_ diamcircleCoreRunnerMode) (diamcircleCoreRunnerInterface, error) {
			return mockRunner, nil
		},
		checkpointManager: historyarchive.NewCheckpointManager(64),
	}

	err := captiveBackend.PrepareRange(ctx, BoundedRange(65, 66))
	assert.NoError(t, err)

	_, err = captiveBackend.GetLedger(ctx, 66)
	tt.EqualError(err, "error closing session: transient error")

	mockArchive.AssertExpectations(t)
	mockRunner.AssertExpectations(t)
}

func TestCaptiveAfterClose(t *testing.T) {
	metaChan := make(chan metaResult, 100)

	for i := 64; i <= 66; i++ {
		meta := buildLedgerCloseMeta(testLedgerHeader{sequence: uint32(i)})
		metaChan <- metaResult{
			LedgerCloseMeta: &meta,
		}
	}

	mockRunner := &diamcircleCoreRunnerMock{}
	ctx, cancel := context.WithCancel(context.Background())
	mockRunner.On("catchup", uint32(65), uint32(66)).Return(nil)
	mockRunner.On("getMetaPipe").Return((<-chan metaResult)(metaChan))
	mockRunner.On("context").Return(ctx)
	mockRunner.On("close").Return(nil).Once()

	mockArchive := &historyarchive.MockArchive{}
	mockArchive.
		On("GetRootHAS").
		Return(historyarchive.HistoryArchiveState{
			CurrentLedger: uint32(200),
		}, nil)

	captiveBackend := CaptiveDiamcircleCore{
		archive: mockArchive,
		diamcircleCoreRunnerFactory: func(_ diamcircleCoreRunnerMode) (diamcircleCoreRunnerInterface, error) {
			return mockRunner, nil
		},
		checkpointManager: historyarchive.NewCheckpointManager(64),
		cancel:            cancel,
	}

	boundedRange := BoundedRange(65, 66)
	err := captiveBackend.PrepareRange(ctx, boundedRange)
	assert.NoError(t, err)

	assert.NoError(t, captiveBackend.Close())
	assert.True(t, captiveBackend.isClosed())

	_, err = captiveBackend.GetLedger(ctx, boundedRange.to)
	assert.EqualError(t, err, "diamcircle-core is no longer usable")

	var prepared bool
	prepared, err = captiveBackend.IsPrepared(ctx, boundedRange)
	assert.False(t, prepared)
	assert.NoError(t, err)

	_, err = captiveBackend.GetLatestLedgerSequence(ctx)
	assert.EqualError(t, err, "diamcircle-core is no longer usable")

	mockArchive.AssertExpectations(t)
	mockRunner.AssertExpectations(t)
}

func TestGetLedgerBoundsCheck(t *testing.T) {
	metaChan := make(chan metaResult, 100)

	for i := 128; i <= 130; i++ {
		meta := buildLedgerCloseMeta(testLedgerHeader{sequence: uint32(i)})
		metaChan <- metaResult{
			LedgerCloseMeta: &meta,
		}
	}

	ctx := context.Background()
	mockRunner := &diamcircleCoreRunnerMock{}
	mockRunner.On("catchup", uint32(128), uint32(130)).Return(nil).Once()
	mockRunner.On("getMetaPipe").Return((<-chan metaResult)(metaChan))
	mockRunner.On("context").Return(ctx)

	mockArchive := &historyarchive.MockArchive{}
	mockArchive.
		On("GetRootHAS").
		Return(historyarchive.HistoryArchiveState{
			CurrentLedger: uint32(200),
		}, nil)

	captiveBackend := CaptiveDiamcircleCore{
		archive: mockArchive,
		diamcircleCoreRunnerFactory: func(_ diamcircleCoreRunnerMode) (diamcircleCoreRunnerInterface, error) {
			return mockRunner, nil
		},
		checkpointManager: historyarchive.NewCheckpointManager(64),
	}

	err := captiveBackend.PrepareRange(ctx, BoundedRange(128, 130))
	assert.NoError(t, err)

	meta, err := captiveBackend.GetLedger(ctx, 128)
	assert.NoError(t, err)
	assert.Equal(t, uint32(128), meta.LedgerSequence())

	prev := meta
	meta, err = captiveBackend.GetLedger(ctx, 128)
	assert.NoError(t, err)
	assert.Equal(t, prev, meta)

	_, err = captiveBackend.GetLedger(ctx, 64)
	assert.EqualError(t, err, "requested ledger 64 is behind the captive core stream (expected=129)")

	mockArchive.AssertExpectations(t)
	mockRunner.AssertExpectations(t)
}

func TestCaptiveGetLedgerTerminatedUnexpectedly(t *testing.T) {
	ledger64 := buildLedgerCloseMeta(testLedgerHeader{sequence: uint32(64)})

	for _, testCase := range []struct {
		name               string
		ctx                context.Context
		ledgers            []metaResult
		processExited      bool
		processExitedError error
		expectedError      string
	}{
		{
			"diamcircle core exited unexpectedly without error",
			context.Background(),
			[]metaResult{{LedgerCloseMeta: &ledger64}},
			true,
			nil,
			"diamcircle core exited unexpectedly",
		},
		{
			"diamcircle core exited unexpectedly with an error",
			context.Background(),
			[]metaResult{{LedgerCloseMeta: &ledger64}},
			true,
			fmt.Errorf("signal kill"),
			"diamcircle core exited unexpectedly: signal kill",
		},
		{
			"diamcircle core exited unexpectedly without error and closed channel",
			context.Background(),
			[]metaResult{{LedgerCloseMeta: &ledger64}},
			true,
			nil,
			"diamcircle core exited unexpectedly",
		},
		{
			"diamcircle core exited unexpectedly with an error and closed channel",
			context.Background(),
			[]metaResult{{LedgerCloseMeta: &ledger64}},
			true,
			fmt.Errorf("signal kill"),
			"diamcircle core exited unexpectedly: signal kill",
		},
		{
			"meta pipe closed unexpectedly",
			context.Background(),
			[]metaResult{{LedgerCloseMeta: &ledger64}},
			false,
			nil,
			"meta pipe closed unexpectedly",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			metaChan := make(chan metaResult, 100)

			for _, result := range testCase.ledgers {
				metaChan <- result
			}
			close(metaChan)

			ctx := testCase.ctx
			mockRunner := &diamcircleCoreRunnerMock{}
			mockRunner.On("catchup", uint32(64), uint32(100)).Return(nil).Once()
			mockRunner.On("getMetaPipe").Return((<-chan metaResult)(metaChan))
			mockRunner.On("context").Return(ctx)
			mockRunner.On("getProcessExitError").Return(testCase.processExited, testCase.processExitedError)
			mockRunner.On("close").Return(nil).Once()

			mockArchive := &historyarchive.MockArchive{}
			mockArchive.
				On("GetRootHAS").
				Return(historyarchive.HistoryArchiveState{
					CurrentLedger: uint32(200),
				}, nil)

			captiveBackend := CaptiveDiamcircleCore{
				archive: mockArchive,
				diamcircleCoreRunnerFactory: func(_ diamcircleCoreRunnerMode) (diamcircleCoreRunnerInterface, error) {
					return mockRunner, nil
				},
				checkpointManager: historyarchive.NewCheckpointManager(64),
			}

			err := captiveBackend.PrepareRange(ctx, BoundedRange(64, 100))
			assert.NoError(t, err)

			meta, err := captiveBackend.GetLedger(ctx, 64)
			assert.NoError(t, err)
			assert.Equal(t, uint32(64), meta.LedgerSequence())

			_, err = captiveBackend.GetLedger(ctx, 65)
			assert.EqualError(t, err, testCase.expectedError)

			mockArchive.AssertExpectations(t)
			mockRunner.AssertExpectations(t)
		})
	}
}

func TestCaptiveUseOfLedgerHashStore(t *testing.T) {
	ctx := context.Background()
	mockArchive := &historyarchive.MockArchive{}
	mockArchive.
		On("GetLedgerHeader", uint32(300)).
		Return(xdr.LedgerHeaderHistoryEntry{
			Header: xdr.LedgerHeader{
				PreviousLedgerHash: xdr.Hash{1, 1, 1, 1},
			},
		}, nil)

	mockLedgerHashStore := &MockLedgerHashStore{}
	mockLedgerHashStore.On("GetLedgerHash", ctx, uint32(1049)).
		Return("", false, fmt.Errorf("transient error")).Once()
	mockLedgerHashStore.On("GetLedgerHash", ctx, uint32(299)).
		Return("", false, nil).Once()
	mockLedgerHashStore.On("GetLedgerHash", ctx, uint32(85)).
		Return("cde", true, nil).Once()
	mockLedgerHashStore.On("GetLedgerHash", ctx, uint32(127)).
		Return("ghi", true, nil).Once()
	mockLedgerHashStore.On("GetLedgerHash", ctx, uint32(2)).
		Return("mnb", true, nil).Once()

	cancelCalled := false
	captiveBackend := CaptiveDiamcircleCore{
		archive:           mockArchive,
		ledgerHashStore:   mockLedgerHashStore,
		checkpointManager: historyarchive.NewCheckpointManager(64),
		cancel: context.CancelFunc(func() {
			cancelCalled = true
		}),
	}

	runFrom, ledgerHash, err := captiveBackend.runFromParams(ctx, 24)
	assert.NoError(t, err)
	assert.Equal(t, uint32(2), runFrom)
	assert.Equal(t, "mnb", ledgerHash)

	runFrom, ledgerHash, err = captiveBackend.runFromParams(ctx, 86)
	assert.NoError(t, err)
	assert.Equal(t, uint32(85), runFrom)
	assert.Equal(t, "cde", ledgerHash)

	runFrom, ledgerHash, err = captiveBackend.runFromParams(ctx, 128)
	assert.NoError(t, err)
	assert.Equal(t, uint32(127), runFrom)
	assert.Equal(t, "ghi", ledgerHash)

	_, _, err = captiveBackend.runFromParams(ctx, 1050)
	assert.EqualError(t, err, "error trying to read ledger hash 1049: transient error")

	runFrom, ledgerHash, err = captiveBackend.runFromParams(ctx, 300)
	assert.NoError(t, err)
	assert.Equal(t, uint32(299), runFrom, "runFrom")
	assert.Equal(t, "0101010100000000000000000000000000000000000000000000000000000000", ledgerHash)

	mockLedgerHashStore.On("Close").Return(nil).Once()
	err = captiveBackend.Close()
	assert.NoError(t, err)
	assert.True(t, cancelCalled)
	mockLedgerHashStore.AssertExpectations(t)
	mockArchive.AssertExpectations(t)
}

func TestCaptiveRunFromParams(t *testing.T) {
	var tests = []struct {
		from           uint32
		runFrom        uint32
		ledgerArchives uint32
	}{
		// Before and including 1st checkpoint:
		{2, 2, 3},
		{3, 2, 3},
		{3, 2, 3},
		{4, 2, 3},
		{62, 2, 3},
		{63, 2, 3},

		// Starting from 64 we go normal path: between 1st and 2nd checkpoint:
		{64, 63, 64},
		{65, 64, 65},
		{66, 65, 66},
		{126, 125, 126},

		// between 2nd and 3rd checkpoint... and so on.
		{127, 126, 127},
		{128, 127, 128},
		{129, 128, 129},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("from_%d", tc.from), func(t *testing.T) {
			tt := assert.New(t)
			mockArchive := &historyarchive.MockArchive{}
			mockArchive.
				On("GetLedgerHeader", uint32(tc.ledgerArchives)).
				Return(xdr.LedgerHeaderHistoryEntry{
					Header: xdr.LedgerHeader{
						PreviousLedgerHash: xdr.Hash{1, 1, 1, 1},
					},
				}, nil)

			captiveBackend := CaptiveDiamcircleCore{
				archive:           mockArchive,
				checkpointManager: historyarchive.NewCheckpointManager(64),
			}

			ctx := context.Background()
			runFrom, ledgerHash, err := captiveBackend.runFromParams(ctx, tc.from)
			tt.NoError(err)
			tt.Equal(tc.runFrom, runFrom, "runFrom")
			tt.Equal("0101010100000000000000000000000000000000000000000000000000000000", ledgerHash)

			mockArchive.AssertExpectations(t)
		})
	}
}

func TestCaptiveIsPrepared(t *testing.T) {
	mockRunner := &diamcircleCoreRunnerMock{}
	mockRunner.On("context").Return(context.Background()).Maybe()

	// c.prepared == nil
	captiveBackend := CaptiveDiamcircleCore{
		nextLedger: 0,
	}

	result := captiveBackend.isPrepared(UnboundedRange(100))
	assert.False(t, result)

	// c.prepared != nil:
	var tests = []struct {
		nextLedger    uint32
		lastLedger    uint32
		cachedLedger  uint32
		preparedRange Range
		ledgerRange   Range
		result        bool
	}{
		// If nextLedger == 0, prepared range is checked
		{0, 0, 0, UnboundedRange(100), UnboundedRange(100), true},
		{0, 0, 0, UnboundedRange(100), UnboundedRange(99), false},
		{0, 0, 0, UnboundedRange(100), BoundedRange(100, 200), true},

		{100, 0, 0, UnboundedRange(99), UnboundedRange(101), true},
		{101, 0, 100, UnboundedRange(99), UnboundedRange(100), true},
		{100, 200, 0, BoundedRange(99, 200), UnboundedRange(100), false},

		{100, 200, 0, BoundedRange(99, 200), BoundedRange(100, 200), true},
		{100, 200, 0, BoundedRange(99, 200), BoundedRange(100, 201), false},
		{100, 201, 0, BoundedRange(99, 201), BoundedRange(100, 200), true},
		{101, 200, 100, BoundedRange(99, 200), BoundedRange(100, 200), true},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("next_%d_last_%d_cached_%d_range_%v", tc.nextLedger, tc.lastLedger, tc.cachedLedger, tc.ledgerRange), func(t *testing.T) {
			captiveBackend := CaptiveDiamcircleCore{
				diamcircleCoreRunner: mockRunner,
				nextLedger:        tc.nextLedger,
				prepared:          &tc.preparedRange,
			}
			if tc.lastLedger > 0 {
				captiveBackend.lastLedger = &tc.lastLedger
			}
			if tc.cachedLedger > 0 {
				meta := buildLedgerCloseMeta(testLedgerHeader{
					sequence: tc.cachedLedger,
				})
				captiveBackend.cachedMeta = &meta
			}

			result := captiveBackend.isPrepared(tc.ledgerRange)
			assert.Equal(t, tc.result, result)
		})
	}
}

// TestCaptivePreviousLedgerCheck checks if previousLedgerHash is set in PrepareRange
// and then checked and updated in GetLedger.
func TestCaptivePreviousLedgerCheck(t *testing.T) {
	metaChan := make(chan metaResult, 200)

	h := 3
	for i := 192; i <= 300; i++ {
		meta := buildLedgerCloseMeta(testLedgerHeader{
			sequence:           uint32(i),
			hash:               fmt.Sprintf("%02x00000000000000000000000000000000000000000000000000000000000000", h),
			previousLedgerHash: fmt.Sprintf("%02x00000000000000000000000000000000000000000000000000000000000000", h-1),
		})
		metaChan <- metaResult{
			LedgerCloseMeta: &meta,
		}
		h++
	}

	{
		// Write invalid hash
		meta := buildLedgerCloseMeta(testLedgerHeader{
			sequence:           301,
			hash:               "0000000000000000000000000000000000000000000000000000000000000000",
			previousLedgerHash: "0000000000000000000000000000000000000000000000000000000000000000",
		})
		metaChan <- metaResult{
			LedgerCloseMeta: &meta,
		}

	}

	ctx := context.Background()
	mockRunner := &diamcircleCoreRunnerMock{}
	mockRunner.On("runFrom", uint32(299), "0101010100000000000000000000000000000000000000000000000000000000").Return(nil).Once()
	mockRunner.On("getMetaPipe").Return((<-chan metaResult)(metaChan))
	mockRunner.On("context").Return(ctx)
	mockRunner.On("close").Return(nil).Once()

	mockArchive := &historyarchive.MockArchive{}
	mockArchive.
		On("GetRootHAS").
		Return(historyarchive.HistoryArchiveState{
			CurrentLedger: uint32(255),
		}, nil)
	mockArchive.
		On("GetLedgerHeader", uint32(300)).
		Return(xdr.LedgerHeaderHistoryEntry{
			Header: xdr.LedgerHeader{
				PreviousLedgerHash: xdr.Hash{1, 1, 1, 1},
			},
		}, nil).Once()

	mockLedgerHashStore := &MockLedgerHashStore{}
	mockLedgerHashStore.On("GetLedgerHash", ctx, uint32(299)).
		Return("", false, nil).Once()

	captiveBackend := CaptiveDiamcircleCore{
		archive: mockArchive,
		diamcircleCoreRunnerFactory: func(_ diamcircleCoreRunnerMode) (diamcircleCoreRunnerInterface, error) {
			return mockRunner, nil
		},
		ledgerHashStore:   mockLedgerHashStore,
		checkpointManager: historyarchive.NewCheckpointManager(64),
	}

	err := captiveBackend.PrepareRange(ctx, UnboundedRange(300))
	assert.NoError(t, err)

	meta, err := captiveBackend.GetLedger(ctx, 300)
	assert.NoError(t, err)
	assert.NotNil(t, captiveBackend.previousLedgerHash)
	assert.Equal(t, uint32(301), captiveBackend.nextLedger)
	assert.Equal(t, meta.LedgerHash().HexString(), *captiveBackend.previousLedgerHash)

	_, err = captiveBackend.GetLedger(ctx, 301)
	assert.EqualError(t, err, "unexpected previous ledger hash for ledger 301 (expected=6f00000000000000000000000000000000000000000000000000000000000000 actual=0000000000000000000000000000000000000000000000000000000000000000)")

	mockRunner.AssertExpectations(t)
	mockArchive.AssertExpectations(t)
	mockLedgerHashStore.AssertExpectations(t)
}
