//lint:file-ignore U1001 Ignore all unused code, staticcheck doesn't understand testify/suite

package ingest

import (
	"context"
	"testing"

	"github.com/diamcircle/go/ingest"
	"github.com/diamcircle/go/ingest/ledgerbackend"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/xdr"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestBuildStateTestSuite(t *testing.T) {
	suite.Run(t, new(BuildStateTestSuite))
}

type BuildStateTestSuite struct {
	suite.Suite
	ctx               context.Context
	historyQ          *mockDBQ
	historyAdapter    *mockHistoryArchiveAdapter
	ledgerBackend     *ledgerbackend.MockDatabaseBackend
	system            *system
	runner            *mockProcessorsRunner
	diamcircleCoreClient *mockDiamcircleCoreClient
	checkpointLedger  uint32
	lastLedger        uint32
}

func (s *BuildStateTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.historyQ = &mockDBQ{}
	s.runner = &mockProcessorsRunner{}
	s.historyAdapter = &mockHistoryArchiveAdapter{}
	s.ledgerBackend = &ledgerbackend.MockDatabaseBackend{}
	s.diamcircleCoreClient = &mockDiamcircleCoreClient{}
	s.checkpointLedger = uint32(63)
	s.lastLedger = 0
	s.system = &system{
		ctx:               s.ctx,
		historyQ:          s.historyQ,
		historyAdapter:    s.historyAdapter,
		ledgerBackend:     s.ledgerBackend,
		runner:            s.runner,
		diamcircleCoreClient: s.diamcircleCoreClient,
	}
	s.system.initMetrics()

	s.historyQ.On("Begin").Return(nil).Once()
	s.historyQ.On("Rollback").Return(nil).Once()

	s.ledgerBackend.On("IsPrepared", s.ctx, ledgerbackend.UnboundedRange(63)).Return(false, nil).Once()
	s.ledgerBackend.On("PrepareRange", s.ctx, ledgerbackend.UnboundedRange(63)).Return(nil).Once()
	s.ledgerBackend.On("GetLedger", s.ctx, uint32(63)).Return(xdr.LedgerCloseMeta{
		V0: &xdr.LedgerCloseMetaV0{
			LedgerHeader: xdr.LedgerHeaderHistoryEntry{
				Header: xdr.LedgerHeader{
					LedgerSeq:      63,
					LedgerVersion:  xdr.Uint32(MaxSupportedProtocolVersion),
					BucketListHash: xdr.Hash{1, 2, 3},
				},
			},
		},
	}, nil).Once()
}

func (s *BuildStateTestSuite) TearDownTest() {
	t := s.T()
	s.historyQ.AssertExpectations(t)
	s.historyAdapter.AssertExpectations(t)
	s.runner.AssertExpectations(t)
	s.diamcircleCoreClient.AssertExpectations(t)
	s.ledgerBackend.AssertExpectations(t)
}

func (s *BuildStateTestSuite) mockCommonHistoryQ() {
	s.historyQ.On("GetLastLedgerIngest", s.ctx).Return(s.lastLedger, nil).Once()
	s.historyQ.On("GetIngestVersion", s.ctx).Return(CurrentVersion, nil).Once()
	s.historyQ.On("UpdateLastLedgerIngest", s.ctx, s.lastLedger).Return(nil).Once()
	s.historyQ.On("UpdateExpStateInvalid", s.ctx, false).Return(nil).Once()
	s.historyQ.On("TruncateIngestStateTables", s.ctx).Return(nil).Once()
	s.diamcircleCoreClient.On(
		"SetCursor",
		mock.AnythingOfType("*context.timerCtx"),
		defaultCoreCursorName,
		int32(62),
	).Return(nil).Once()
}

func (s *BuildStateTestSuite) TestCheckPointLedgerIsZero() {
	// Recreate mock in this single test to remove assertions.
	*s.historyQ = mockDBQ{}
	*s.ledgerBackend = ledgerbackend.MockDatabaseBackend{}

	next, err := buildState{checkpointLedger: 0}.run(s.system)
	s.Assert().Error(err)
	s.Assert().EqualError(err, "unexpected checkpointLedger value")
	s.Assert().Equal(transition{node: startState{}, sleepDuration: defaultSleep}, next)
}

func (s *BuildStateTestSuite) TestRangeNotPreparedFailPrepare() {
	// Recreate mock in this single test to remove assertions.
	*s.historyQ = mockDBQ{}
	*s.ledgerBackend = ledgerbackend.MockDatabaseBackend{}

	s.ledgerBackend.On("IsPrepared", s.ctx, ledgerbackend.UnboundedRange(63)).Return(false, nil).Once()
	s.ledgerBackend.On("PrepareRange", s.ctx, ledgerbackend.UnboundedRange(63)).Return(errors.New("my error")).Once()

	next, err := buildState{checkpointLedger: s.checkpointLedger}.run(s.system)

	s.Assert().Error(err)
	s.Assert().EqualError(err, "error preparing range: my error")
	s.Assert().Equal(transition{node: startState{}, sleepDuration: defaultSleep}, next)
}

func (s *BuildStateTestSuite) TestRangeNotPreparedSuccessPrepareGetLedgerFail() {
	// Recreate mock in this single test to remove assertions.
	*s.historyQ = mockDBQ{}
	*s.ledgerBackend = ledgerbackend.MockDatabaseBackend{}

	s.ledgerBackend.On("IsPrepared", s.ctx, ledgerbackend.UnboundedRange(63)).Return(false, nil).Once()
	s.ledgerBackend.On("PrepareRange", s.ctx, ledgerbackend.UnboundedRange(63)).Return(nil).Once()
	s.ledgerBackend.On("GetLedger", s.ctx, uint32(63)).Return(xdr.LedgerCloseMeta{}, errors.New("my error")).Once()

	next, err := buildState{checkpointLedger: s.checkpointLedger}.run(s.system)

	s.Assert().Error(err)
	s.Assert().EqualError(err, "error getting ledger blocking: my error")
	s.Assert().Equal(transition{node: startState{}, sleepDuration: defaultSleep}, next)
}

func (s *BuildStateTestSuite) TestBeginReturnsError() {
	// Recreate mock in this single test to remove assertions.
	*s.historyQ = mockDBQ{}
	s.historyQ.On("Begin").Return(errors.New("my error")).Once()

	next, err := buildState{checkpointLedger: s.checkpointLedger}.run(s.system)
	s.Assert().Error(err)
	s.Assert().EqualError(err, "Error starting a transaction: my error")
	s.Assert().Equal(transition{node: startState{}, sleepDuration: defaultSleep}, next)
}

func (s *BuildStateTestSuite) TestGetLastLedgerIngestReturnsError() {
	s.historyQ.On("GetLastLedgerIngest", s.ctx).Return(s.lastLedger, errors.New("my error")).Once()

	next, err := buildState{checkpointLedger: s.checkpointLedger}.run(s.system)
	s.Assert().Error(err)
	s.Assert().EqualError(err, "Error getting last ingested ledger: my error")
	s.Assert().Equal(transition{node: startState{}, sleepDuration: defaultSleep}, next)
}
func (s *BuildStateTestSuite) TestGetIngestVersionReturnsError() {
	s.historyQ.On("GetLastLedgerIngest", s.ctx).Return(s.lastLedger, nil).Once()
	s.historyQ.On("GetIngestVersion", s.ctx).Return(CurrentVersion, errors.New("my error")).Once()

	next, err := buildState{checkpointLedger: s.checkpointLedger}.run(s.system)
	s.Assert().Error(err)
	s.Assert().EqualError(err, "Error getting ingestion version: my error")
	s.Assert().Equal(transition{node: startState{}, sleepDuration: defaultSleep}, next)
}

func (s *BuildStateTestSuite) TestAnotherInstanceHasCompletedBuildState() {
	s.historyQ.On("GetLastLedgerIngest", s.ctx).Return(s.checkpointLedger, nil).Once()
	s.historyQ.On("GetIngestVersion", s.ctx).Return(CurrentVersion, nil).Once()

	next, err := buildState{checkpointLedger: s.checkpointLedger}.run(s.system)
	s.Assert().NoError(err)
	s.Assert().Equal(transition{node: startState{}, sleepDuration: defaultSleep}, next)
}

func (s *BuildStateTestSuite) TestUpdateLastLedgerIngestReturnsError() {
	s.historyQ.On("GetLastLedgerIngest", s.ctx).Return(s.lastLedger, nil).Once()
	s.historyQ.On("GetIngestVersion", s.ctx).Return(CurrentVersion, nil).Once()
	s.historyQ.On("UpdateLastLedgerIngest", s.ctx, s.lastLedger).Return(errors.New("my error")).Once()
	s.diamcircleCoreClient.On(
		"SetCursor",
		mock.AnythingOfType("*context.timerCtx"),
		defaultCoreCursorName,
		int32(62),
	).Return(nil).Once()

	next, err := buildState{checkpointLedger: s.checkpointLedger}.run(s.system)

	s.Assert().Error(err)
	s.Assert().EqualError(err, "Error updating last ingested ledger: my error")
	s.Assert().Equal(transition{node: startState{}, sleepDuration: defaultSleep}, next)
}

func (s *BuildStateTestSuite) TestUpdateExpStateInvalidReturnsError() {
	s.historyQ.On("GetLastLedgerIngest", s.ctx).Return(s.lastLedger, nil).Once()
	s.historyQ.On("GetIngestVersion", s.ctx).Return(CurrentVersion, nil).Once()
	s.historyQ.On("UpdateLastLedgerIngest", s.ctx, s.lastLedger).Return(nil).Once()
	s.historyQ.On("UpdateExpStateInvalid", s.ctx, false).Return(errors.New("my error")).Once()
	s.diamcircleCoreClient.On(
		"SetCursor",
		mock.AnythingOfType("*context.timerCtx"),
		defaultCoreCursorName,
		int32(62),
	).Return(nil).Once()

	next, err := buildState{checkpointLedger: s.checkpointLedger}.run(s.system)

	s.Assert().Error(err)
	s.Assert().EqualError(err, "Error updating state invalid value: my error")
	s.Assert().Equal(transition{node: startState{}, sleepDuration: defaultSleep}, next)
}

func (s *BuildStateTestSuite) TestTruncateIngestStateTablesReturnsError() {
	s.historyQ.On("GetLastLedgerIngest", s.ctx).Return(s.lastLedger, nil).Once()
	s.historyQ.On("GetIngestVersion", s.ctx).Return(CurrentVersion, nil).Once()
	s.historyQ.On("UpdateLastLedgerIngest", s.ctx, s.lastLedger).Return(nil).Once()
	s.historyQ.On("UpdateExpStateInvalid", s.ctx, false).Return(nil).Once()
	s.historyQ.On("TruncateIngestStateTables", s.ctx).Return(errors.New("my error")).Once()

	s.diamcircleCoreClient.On(
		"SetCursor",
		mock.AnythingOfType("*context.timerCtx"),
		defaultCoreCursorName,
		int32(62),
	).Return(nil).Once()

	next, err := buildState{checkpointLedger: s.checkpointLedger}.run(s.system)

	s.Assert().Error(err)
	s.Assert().EqualError(err, "Error clearing ingest tables: my error")
	s.Assert().Equal(transition{node: startState{}, sleepDuration: defaultSleep}, next)
}

func (s *BuildStateTestSuite) TestRunHistoryArchiveIngestionReturnsError() {
	s.mockCommonHistoryQ()
	s.runner.
		On("RunHistoryArchiveIngestion", s.checkpointLedger, MaxSupportedProtocolVersion, xdr.Hash{1, 2, 3}).
		Return(ingest.StatsChangeProcessorResults{}, errors.New("my error")).
		Once()
	next, err := buildState{checkpointLedger: s.checkpointLedger}.run(s.system)

	s.Assert().Error(err)
	s.Assert().EqualError(err, "Error ingesting history archive: my error")
	s.Assert().Equal(transition{node: startState{}, sleepDuration: defaultSleep}, next)
}

func (s *BuildStateTestSuite) TestRunHistoryArchiveIngestionGenesisReturnsError() {
	// Recreate mock in this single test to remove assertions.
	*s.ledgerBackend = ledgerbackend.MockDatabaseBackend{}

	s.historyQ.On("GetLastLedgerIngest", s.ctx).Return(uint32(0), nil).Once()
	s.historyQ.On("GetIngestVersion", s.ctx).Return(CurrentVersion, nil).Once()
	s.historyQ.On("UpdateLastLedgerIngest", s.ctx, uint32(0)).Return(nil).Once()
	s.historyQ.On("UpdateExpStateInvalid", s.ctx, false).Return(nil).Once()
	s.historyQ.On("TruncateIngestStateTables", s.ctx).Return(nil).Once()
	s.diamcircleCoreClient.On(
		"SetCursor",
		mock.AnythingOfType("*context.timerCtx"),
		defaultCoreCursorName,
		int32(0),
	).Return(nil).Once()

	s.runner.
		On("RunGenesisStateIngestion").
		Return(ingest.StatsChangeProcessorResults{}, errors.New("my error")).
		Once()
	next, err := buildState{checkpointLedger: 1}.run(s.system)

	s.Assert().Error(err)
	s.Assert().EqualError(err, "Error ingesting history archive: my error")
	s.Assert().Equal(transition{node: startState{}, sleepDuration: defaultSleep}, next)
}

func (s *BuildStateTestSuite) TestUpdateLastLedgerIngestAfterIngestReturnsError() {
	s.mockCommonHistoryQ()
	s.runner.
		On("RunHistoryArchiveIngestion", s.checkpointLedger, MaxSupportedProtocolVersion, xdr.Hash{1, 2, 3}).
		Return(ingest.StatsChangeProcessorResults{}, nil).
		Once()
	s.historyQ.On("UpdateIngestVersion", s.ctx, CurrentVersion).
		Return(nil).
		Once()
	s.historyQ.On("UpdateLastLedgerIngest", s.ctx, s.checkpointLedger).
		Return(errors.New("my error")).
		Once()
	next, err := buildState{checkpointLedger: s.checkpointLedger}.run(s.system)

	s.Assert().Error(err)
	s.Assert().EqualError(err, "Error updating last ingested ledger: my error")
	s.Assert().Equal(transition{node: startState{}, sleepDuration: defaultSleep}, next)
}

func (s *BuildStateTestSuite) TestUpdateIngestVersionIngestReturnsError() {
	s.mockCommonHistoryQ()
	s.runner.
		On("RunHistoryArchiveIngestion", s.checkpointLedger, MaxSupportedProtocolVersion, xdr.Hash{1, 2, 3}).
		Return(ingest.StatsChangeProcessorResults{}, nil).
		Once()
	s.historyQ.On("UpdateIngestVersion", s.ctx, CurrentVersion).
		Return(errors.New("my error")).
		Once()
	next, err := buildState{checkpointLedger: s.checkpointLedger}.run(s.system)

	s.Assert().Error(err)
	s.Assert().EqualError(err, "Error updating ingestion version: my error")
	s.Assert().Equal(transition{node: startState{}, sleepDuration: defaultSleep}, next)
}

func (s *BuildStateTestSuite) TestUpdateCommitReturnsError() {
	s.mockCommonHistoryQ()
	s.runner.
		On("RunHistoryArchiveIngestion", s.checkpointLedger, MaxSupportedProtocolVersion, xdr.Hash{1, 2, 3}).
		Return(ingest.StatsChangeProcessorResults{}, nil).
		Once()
	s.historyQ.On("UpdateLastLedgerIngest", s.ctx, s.checkpointLedger).
		Return(nil).
		Once()
	s.historyQ.On("UpdateIngestVersion", s.ctx, CurrentVersion).
		Return(nil).
		Once()
	s.historyQ.On("Commit").
		Return(errors.New("my error")).
		Once()
	next, err := buildState{checkpointLedger: s.checkpointLedger}.run(s.system)

	s.Assert().Error(err)
	s.Assert().EqualError(err, "Error committing db transaction: my error")
	s.Assert().Equal(transition{node: startState{}, sleepDuration: defaultSleep}, next)
}

func (s *BuildStateTestSuite) TestBuildStateSucceeds() {
	s.mockCommonHistoryQ()
	s.runner.
		On("RunHistoryArchiveIngestion", s.checkpointLedger, MaxSupportedProtocolVersion, xdr.Hash{1, 2, 3}).
		Return(ingest.StatsChangeProcessorResults{}, nil).
		Once()
	s.historyQ.On("UpdateLastLedgerIngest", s.ctx, s.checkpointLedger).
		Return(nil).
		Once()
	s.historyQ.On("UpdateIngestVersion", s.ctx, CurrentVersion).
		Return(nil).
		Once()
	s.historyQ.On("Commit").
		Return(nil).
		Once()

	next, err := buildState{checkpointLedger: s.checkpointLedger}.run(s.system)

	s.Assert().NoError(err)
	s.Assert().Equal(
		transition{
			node:          resumeState{latestSuccessfullyProcessedLedger: s.checkpointLedger},
			sleepDuration: defaultSleep,
		},
		next,
	)
}

func (s *BuildStateTestSuite) TestUpdateCommitReturnsErrorStop() {
	s.mockCommonHistoryQ()
	s.runner.
		On("RunHistoryArchiveIngestion", s.checkpointLedger, MaxSupportedProtocolVersion, xdr.Hash{1, 2, 3}).
		Return(ingest.StatsChangeProcessorResults{}, nil).
		Once()
	s.historyQ.On("UpdateLastLedgerIngest", s.ctx, s.checkpointLedger).
		Return(nil).
		Once()
	s.historyQ.On("UpdateIngestVersion", s.ctx, CurrentVersion).
		Return(nil).
		Once()
	s.historyQ.On("Commit").
		Return(errors.New("my error")).
		Once()
	next, err := buildState{checkpointLedger: s.checkpointLedger, stop: true}.run(s.system)

	s.Assert().Error(err)
	s.Assert().EqualError(err, "Error committing db transaction: my error")
	s.Assert().Equal(transition{node: stopState{}, sleepDuration: 0}, next)
}

func (s *BuildStateTestSuite) TestBuildStateSucceedStop() {
	s.mockCommonHistoryQ()
	s.runner.
		On("RunHistoryArchiveIngestion", s.checkpointLedger, MaxSupportedProtocolVersion, xdr.Hash{1, 2, 3}).
		Return(ingest.StatsChangeProcessorResults{}, nil).
		Once()
	s.historyQ.On("UpdateLastLedgerIngest", s.ctx, s.checkpointLedger).
		Return(nil).
		Once()
	s.historyQ.On("UpdateIngestVersion", s.ctx, CurrentVersion).
		Return(nil).
		Once()
	s.historyQ.On("Commit").
		Return(nil).
		Once()

	next, err := buildState{checkpointLedger: s.checkpointLedger, stop: true}.run(s.system)

	s.Assert().NoError(err)
	s.Assert().Equal(
		transition{
			node:          stopState{},
			sleepDuration: 0,
		},
		next,
	)
}
