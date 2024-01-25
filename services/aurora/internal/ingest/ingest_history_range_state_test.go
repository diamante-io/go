//lint:file-ignore U1001 Ignore all unused code, staticcheck doesn't understand testify/suite

package ingest

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"

	"github.com/diamcircle/go/ingest/ledgerbackend"
	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/services/aurora/internal/ingest/processors"
	"github.com/diamcircle/go/services/aurora/internal/toid"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/xdr"
)

func TestIngestHistoryRangeStateTestSuite(t *testing.T) {
	suite.Run(t, new(IngestHistoryRangeStateTestSuite))
}

type IngestHistoryRangeStateTestSuite struct {
	suite.Suite
	ctx            context.Context
	historyQ       *mockDBQ
	historyAdapter *mockHistoryArchiveAdapter
	ledgerBackend  *ledgerbackend.MockDatabaseBackend
	runner         *mockProcessorsRunner
	system         *system
}

func (s *IngestHistoryRangeStateTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.historyQ = &mockDBQ{}
	s.ledgerBackend = &ledgerbackend.MockDatabaseBackend{}
	s.historyAdapter = &mockHistoryArchiveAdapter{}
	s.runner = &mockProcessorsRunner{}
	s.system = &system{
		ctx:            s.ctx,
		historyQ:       s.historyQ,
		historyAdapter: s.historyAdapter,
		ledgerBackend:  s.ledgerBackend,
		runner:         s.runner,
	}
	s.system.initMetrics()

	s.historyQ.On("Rollback").Return(nil).Once()

	s.ledgerBackend.On("IsPrepared", s.ctx, ledgerbackend.UnboundedRange(100)).Return(false, nil).Once()
	s.ledgerBackend.On("PrepareRange", s.ctx, ledgerbackend.UnboundedRange(100)).Return(nil).Once()
}

func (s *IngestHistoryRangeStateTestSuite) TearDownTest() {
	t := s.T()
	s.historyQ.AssertExpectations(t)
	s.ledgerBackend.AssertExpectations(t)
	s.historyAdapter.AssertExpectations(t)
	s.runner.AssertExpectations(t)
}

func (s *IngestHistoryRangeStateTestSuite) TestInvalidRange() {
	// Recreate mock in this single test to remove Rollback assertion.
	*s.historyQ = mockDBQ{}
	*s.ledgerBackend = ledgerbackend.MockDatabaseBackend{}

	next, err := historyRangeState{fromLedger: 0, toLedger: 0}.run(s.system)
	s.Assert().Error(err)
	s.Assert().EqualError(err, "invalid range: [0, 0]")
	s.Assert().Equal(transition{node: startState{}, sleepDuration: defaultSleep}, next)

	next, err = historyRangeState{fromLedger: 0, toLedger: 100}.run(s.system)
	s.Assert().Error(err)
	s.Assert().EqualError(err, "invalid range: [0, 100]")
	s.Assert().Equal(transition{node: startState{}, sleepDuration: defaultSleep}, next)

	next, err = historyRangeState{fromLedger: 100, toLedger: 0}.run(s.system)
	s.Assert().Error(err)
	s.Assert().EqualError(err, "invalid range: [100, 0]")
	s.Assert().Equal(transition{node: startState{}, sleepDuration: defaultSleep}, next)

	next, err = historyRangeState{fromLedger: 100, toLedger: 99}.run(s.system)
	s.Assert().Error(err)
	s.Assert().EqualError(err, "invalid range: [100, 99]")
	s.Assert().Equal(transition{node: startState{}, sleepDuration: defaultSleep}, next)
}

func (s *IngestHistoryRangeStateTestSuite) TestRangeNotPreparedFailPrepare() {
	// Recreate mock in this single test to remove assertions.
	*s.historyQ = mockDBQ{}
	*s.ledgerBackend = ledgerbackend.MockDatabaseBackend{}

	s.ledgerBackend.On("IsPrepared", s.ctx, ledgerbackend.UnboundedRange(100)).Return(false, nil).Once()
	s.ledgerBackend.On("PrepareRange", s.ctx, ledgerbackend.UnboundedRange(100)).Return(errors.New("my error")).Once()

	next, err := historyRangeState{fromLedger: 100, toLedger: 200}.run(s.system)
	s.Assert().Error(err)
	s.Assert().EqualError(err, "error preparing range: my error")
	s.Assert().Equal(transition{node: startState{}, sleepDuration: defaultSleep}, next)
}

func (s *IngestHistoryRangeStateTestSuite) TestBeginReturnsError() {
	// Recreate mock in this single test to remove Rollback assertion.
	*s.historyQ = mockDBQ{}

	s.historyQ.On("Begin").Return(errors.New("my error")).Once()

	next, err := historyRangeState{fromLedger: 100, toLedger: 200}.run(s.system)
	s.Assert().Error(err)
	s.Assert().EqualError(err, "Error starting a transaction: my error")
	s.Assert().Equal(transition{node: startState{}, sleepDuration: defaultSleep}, next)
}

func (s *IngestHistoryRangeStateTestSuite) TestGetLastLedgerIngestReturnsError() {
	s.historyQ.On("Begin").Return(nil).Once()
	s.historyQ.On("GetLastLedgerIngest", s.ctx).Return(uint32(0), errors.New("my error")).Once()

	next, err := historyRangeState{fromLedger: 100, toLedger: 200}.run(s.system)
	s.Assert().Error(err)
	s.Assert().EqualError(err, "Error getting last ingested ledger: my error")
	s.Assert().Equal(transition{node: startState{}, sleepDuration: defaultSleep}, next)
}

func (s *IngestHistoryRangeStateTestSuite) TestGetLatestLedgerReturnsError() {
	s.historyQ.On("Begin").Return(nil).Once()
	s.historyQ.On("GetLastLedgerIngest", s.ctx).Return(uint32(0), nil).Once()
	s.historyQ.On("GetLatestHistoryLedger", s.ctx).Return(uint32(0), errors.New("my error")).Once()

	next, err := historyRangeState{fromLedger: 100, toLedger: 200}.run(s.system)
	s.Assert().Error(err)
	s.Assert().EqualError(err, "could not get latest history ledger: my error")
	s.Assert().Equal(transition{node: startState{}, sleepDuration: defaultSleep}, next)
}

// TestAnotherNodeIngested tests the case when another node has ingested the range.
// In such case we go back to `init` state without processing.
func (s *IngestHistoryRangeStateTestSuite) TestAnotherNodeIngested() {
	s.historyQ.On("Begin").Return(nil).Once()
	s.historyQ.On("GetLastLedgerIngest", s.ctx).Return(uint32(0), nil).Once()
	s.historyQ.On("GetLatestHistoryLedger", s.ctx).Return(uint32(200), nil).Once()

	next, err := historyRangeState{fromLedger: 100, toLedger: 200}.run(s.system)
	s.Assert().NoError(err)
	s.Assert().Equal(transition{node: startState{}, sleepDuration: defaultSleep}, next)
}

func (s *IngestHistoryRangeStateTestSuite) TestRunTransactionProcessorsOnLedgerReturnsError() {
	s.historyQ.On("Begin").Return(nil).Once()
	s.historyQ.On("GetLastLedgerIngest", s.ctx).Return(uint32(0), nil).Once()
	s.historyQ.On("GetLatestHistoryLedger", s.ctx).Return(uint32(99), nil).Once()

	meta := xdr.LedgerCloseMeta{
		V0: &xdr.LedgerCloseMetaV0{
			LedgerHeader: xdr.LedgerHeaderHistoryEntry{
				Header: xdr.LedgerHeader{
					LedgerSeq: 100,
				},
			},
		},
	}
	s.ledgerBackend.On("GetLedger", s.ctx, uint32(100)).Return(meta, nil).Once()

	s.runner.On("RunTransactionProcessorsOnLedger", meta).Return(
		processors.StatsLedgerTransactionProcessorResults{},
		processorsRunDurations{},
		processors.TradeStats{},
		errors.New("my error"),
	).Once()

	next, err := historyRangeState{fromLedger: 100, toLedger: 200}.run(s.system)
	s.Assert().Error(err)
	s.Assert().EqualError(err, "error processing ledger sequence=100: my error")
	s.Assert().Equal(transition{node: startState{}, sleepDuration: defaultSleep}, next)
}

func (s *IngestHistoryRangeStateTestSuite) TestSuccess() {
	s.historyQ.On("Begin").Return(nil).Once()
	s.historyQ.On("GetLastLedgerIngest", s.ctx).Return(uint32(0), nil).Once()
	s.historyQ.On("GetLatestHistoryLedger", s.ctx).Return(uint32(99), nil).Once()

	for i := 100; i <= 200; i++ {
		meta := xdr.LedgerCloseMeta{
			V0: &xdr.LedgerCloseMetaV0{
				LedgerHeader: xdr.LedgerHeaderHistoryEntry{
					Header: xdr.LedgerHeader{
						LedgerSeq: xdr.Uint32(i),
					},
				},
			},
		}
		s.ledgerBackend.On("GetLedger", s.ctx, uint32(i)).Return(meta, nil).Once()

		s.runner.On("RunTransactionProcessorsOnLedger", meta).Return(
			processors.StatsLedgerTransactionProcessorResults{},
			processorsRunDurations{},
			processors.TradeStats{},
			nil,
		).Once()
	}

	s.historyQ.On("Commit").Return(nil).Once()

	next, err := historyRangeState{fromLedger: 100, toLedger: 200}.run(s.system)
	s.Assert().NoError(err)
	s.Assert().Equal(transition{node: startState{}, sleepDuration: defaultSleep}, next)
}

func (s *IngestHistoryRangeStateTestSuite) TestSuccessOneLedger() {
	s.historyQ.On("Begin").Return(nil).Once()
	s.historyQ.On("GetLastLedgerIngest", s.ctx).Return(uint32(0), nil).Once()
	s.historyQ.On("GetLatestHistoryLedger", s.ctx).Return(uint32(99), nil).Once()

	meta := xdr.LedgerCloseMeta{
		V0: &xdr.LedgerCloseMetaV0{
			LedgerHeader: xdr.LedgerHeaderHistoryEntry{
				Header: xdr.LedgerHeader{
					LedgerSeq: xdr.Uint32(100),
				},
			},
		},
	}
	s.ledgerBackend.On("GetLedger", s.ctx, uint32(100)).Return(meta, nil).Once()

	s.runner.On("RunTransactionProcessorsOnLedger", meta).Return(
		processors.StatsLedgerTransactionProcessorResults{},
		processorsRunDurations{},
		processors.TradeStats{},
		nil,
	).Once()

	s.historyQ.On("Commit").Return(nil).Once()

	next, err := historyRangeState{fromLedger: 100, toLedger: 100}.run(s.system)
	s.Assert().NoError(err)
	s.Assert().Equal(transition{node: startState{}, sleepDuration: defaultSleep}, next)
}

func (s *IngestHistoryRangeStateTestSuite) TestCommitsWorkOnLedgerBackendFailure() {
	s.historyQ.On("Begin").Return(nil).Once()
	s.historyQ.On("GetLastLedgerIngest", s.ctx).Return(uint32(0), nil).Once()
	s.historyQ.On("GetLatestHistoryLedger", s.ctx).Return(uint32(99), nil).Once()

	meta := xdr.LedgerCloseMeta{
		V0: &xdr.LedgerCloseMetaV0{
			LedgerHeader: xdr.LedgerHeaderHistoryEntry{
				Header: xdr.LedgerHeader{
					LedgerSeq: xdr.Uint32(100),
				},
			},
		},
	}
	s.ledgerBackend.On("GetLedger", s.ctx, uint32(100)).Return(meta, nil).Once()
	s.ledgerBackend.On("GetLedger", s.ctx, uint32(101)).
		Return(xdr.LedgerCloseMeta{}, errors.New("my error")).Once()

	s.runner.On("RunTransactionProcessorsOnLedger", meta).Return(
		processors.StatsLedgerTransactionProcessorResults{},
		processorsRunDurations{},
		processors.TradeStats{},
		nil,
	).Once()

	s.historyQ.On("Commit").Return(nil).Once()

	next, err := historyRangeState{fromLedger: 100, toLedger: 102}.run(s.system)
	s.Assert().Error(err)
	s.Assert().EqualError(err, "error getting ledger: my error")
	s.Assert().Equal(transition{node: startState{}, sleepDuration: defaultSleep}, next)
}

func TestReingestHistoryRangeStateTestSuite(t *testing.T) {
	suite.Run(t, new(ReingestHistoryRangeStateTestSuite))
}

type ReingestHistoryRangeStateTestSuite struct {
	suite.Suite
	ctx            context.Context
	historyQ       *mockDBQ
	historyAdapter *mockHistoryArchiveAdapter
	ledgerBackend  *mockLedgerBackend
	runner         *mockProcessorsRunner
	system         *system
}

func (s *ReingestHistoryRangeStateTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.historyQ = &mockDBQ{}
	s.historyAdapter = &mockHistoryArchiveAdapter{}
	s.ledgerBackend = &mockLedgerBackend{}
	s.runner = &mockProcessorsRunner{}
	s.system = &system{
		ctx:            s.ctx,
		historyQ:       s.historyQ,
		historyAdapter: s.historyAdapter,
		ledgerBackend:  s.ledgerBackend,
		runner:         s.runner,
	}

	s.historyQ.On("GetTx").Return(nil).Once()
	s.historyQ.On("Rollback").Return(nil).Once()
	s.historyQ.On("Begin").Return(nil).Once()

	s.ledgerBackend.On("PrepareRange", s.ctx, ledgerbackend.BoundedRange(100, 200)).Return(nil).Once()
}

func (s *ReingestHistoryRangeStateTestSuite) TearDownTest() {
	t := s.T()
	s.historyQ.AssertExpectations(t)
	s.historyAdapter.AssertExpectations(t)
	s.runner.AssertExpectations(t)
}

func (s *ReingestHistoryRangeStateTestSuite) TestInvalidRange() {
	// Recreate mock in this single test to remove Rollback assertion.
	*s.historyQ = mockDBQ{}

	err := s.system.ReingestRange([]history.LedgerRange{{0, 0}}, false)
	s.Assert().EqualError(err, "Invalid range: {0 0} genesis ledger starts at 1")

	err = s.system.ReingestRange([]history.LedgerRange{{0, 100}}, false)
	s.Assert().EqualError(err, "Invalid range: {0 100} genesis ledger starts at 1")

	err = s.system.ReingestRange([]history.LedgerRange{{100, 0}}, false)
	s.Assert().EqualError(err, "Invalid range: {100 0} from > to")

	err = s.system.ReingestRange([]history.LedgerRange{{100, 99}}, false)
	s.Assert().EqualError(err, "Invalid range: {100 99} from > to")
}

func (s *ReingestHistoryRangeStateTestSuite) TestBeginReturnsError() {
	// Recreate mock in this single test to remove Rollback assertion.
	*s.historyQ = mockDBQ{}
	s.historyQ.On("GetTx").Return(nil)
	s.historyQ.On("GetLastLedgerIngestNonBlocking", s.ctx).Return(uint32(0), nil).Once()

	s.historyQ.On("Begin").Return(errors.New("my error")).Once()

	err := s.system.ReingestRange([]history.LedgerRange{{100, 200}}, false)
	s.Assert().EqualError(err, "Error starting a transaction: my error")
}

func (s *ReingestHistoryRangeStateTestSuite) TestGetLastLedgerIngestNonBlockingError() {
	*s.historyQ = mockDBQ{}
	s.historyQ.On("GetTx").Return(nil).Once()

	s.historyQ.On("GetLastLedgerIngestNonBlocking", s.ctx).Return(uint32(0), errors.New("my error")).Once()

	err := s.system.ReingestRange([]history.LedgerRange{{100, 200}}, false)
	s.Assert().EqualError(err, "Error getting last ingested ledger: my error")
}

func (s *ReingestHistoryRangeStateTestSuite) TestReingestRangeOverlaps() {
	*s.historyQ = mockDBQ{}
	s.historyQ.On("GetTx").Return(nil).Once()

	s.historyQ.On("GetLastLedgerIngestNonBlocking", s.ctx).Return(uint32(190), nil).Once()

	err := s.system.ReingestRange([]history.LedgerRange{{100, 200}}, false)
	s.Assert().Equal(ErrReingestRangeConflict{190}, err)
}

func (s *ReingestHistoryRangeStateTestSuite) TestReingestRangeOverlapsAtEnd() {
	*s.historyQ = mockDBQ{}
	s.historyQ.On("GetTx").Return(nil).Once()

	s.historyQ.On("GetLastLedgerIngestNonBlocking", s.ctx).Return(uint32(200), nil).Once()

	err := s.system.ReingestRange([]history.LedgerRange{{100, 200}}, false)
	s.Assert().Equal(ErrReingestRangeConflict{200}, err)
}

func (s *ReingestHistoryRangeStateTestSuite) TestClearHistoryFails() {
	*s.historyQ = mockDBQ{}
	s.historyQ.On("GetTx").Return(nil).Once()
	s.historyQ.On("GetLastLedgerIngestNonBlocking", s.ctx).Return(uint32(0), nil).Once()

	s.historyQ.On("Begin").Return(nil).Once()
	s.historyQ.On("GetTx").Return(&sqlx.Tx{}).Once()
	toidFrom := toid.New(100, 0, 0)
	toidTo := toid.New(101, 0, 0)
	s.historyQ.On(
		"DeleteRangeAll", s.ctx, toidFrom.ToInt64(), toidTo.ToInt64(),
	).Return(errors.New("my error")).Once()

	s.historyQ.On("Rollback").Return(nil).Once()

	err := s.system.ReingestRange([]history.LedgerRange{{100, 200}}, false)
	s.Assert().EqualError(err, "error in DeleteRangeAll: my error")
}

func (s *ReingestHistoryRangeStateTestSuite) TestRunTransactionProcessorsOnLedgerReturnsError() {
	*s.historyQ = mockDBQ{}
	s.historyQ.On("GetTx").Return(nil).Once()
	s.historyQ.On("GetLastLedgerIngestNonBlocking", s.ctx).Return(uint32(0), nil).Once()

	s.historyQ.On("Begin").Return(nil).Once()
	s.historyQ.On("GetTx").Return(&sqlx.Tx{}).Once()
	toidFrom := toid.New(100, 0, 0)
	toidTo := toid.New(101, 0, 0)
	s.historyQ.On(
		"DeleteRangeAll", s.ctx, toidFrom.ToInt64(), toidTo.ToInt64(),
	).Return(nil).Once()

	meta := xdr.LedgerCloseMeta{
		V0: &xdr.LedgerCloseMetaV0{
			LedgerHeader: xdr.LedgerHeaderHistoryEntry{
				Header: xdr.LedgerHeader{
					LedgerSeq: xdr.Uint32(100),
				},
			},
		},
	}
	s.ledgerBackend.On("GetLedger", s.ctx, uint32(100)).Return(meta, nil).Once()

	s.runner.On("RunTransactionProcessorsOnLedger", meta).
		Return(
			processors.StatsLedgerTransactionProcessorResults{},
			processorsRunDurations{},
			processors.TradeStats{},
			errors.New("my error"),
		).Once()
	s.historyQ.On("Rollback").Return(nil).Once()

	err := s.system.ReingestRange([]history.LedgerRange{{100, 200}}, false)
	s.Assert().EqualError(err, "error processing ledger sequence=100: my error")
}

func (s *ReingestHistoryRangeStateTestSuite) TestCommitFails() {
	*s.historyQ = mockDBQ{}
	s.historyQ.On("GetTx").Return(nil).Once()
	s.historyQ.On("GetLastLedgerIngestNonBlocking", s.ctx).Return(uint32(0), nil).Once()

	s.historyQ.On("Begin").Return(nil).Once()
	s.historyQ.On("GetTx").Return(&sqlx.Tx{}).Once()
	toidFrom := toid.New(100, 0, 0)
	toidTo := toid.New(101, 0, 0)
	s.historyQ.On(
		"DeleteRangeAll", s.ctx, toidFrom.ToInt64(), toidTo.ToInt64(),
	).Return(nil).Once()

	meta := xdr.LedgerCloseMeta{
		V0: &xdr.LedgerCloseMetaV0{
			LedgerHeader: xdr.LedgerHeaderHistoryEntry{
				Header: xdr.LedgerHeader{
					LedgerSeq: xdr.Uint32(100),
				},
			},
		},
	}
	s.ledgerBackend.On("GetLedger", s.ctx, uint32(100)).Return(meta, nil).Once()

	s.runner.On("RunTransactionProcessorsOnLedger", meta).Return(
		processors.StatsLedgerTransactionProcessorResults{},
		processorsRunDurations{},
		processors.TradeStats{},
		nil,
	).Once()

	s.historyQ.On("Commit").Return(errors.New("my error")).Once()
	s.historyQ.On("Rollback").Return(nil).Once()

	err := s.system.ReingestRange([]history.LedgerRange{{100, 200}}, false)
	s.Assert().EqualError(err, "Error committing db transaction: my error")
}

func (s *ReingestHistoryRangeStateTestSuite) TestSuccess() {
	*s.historyQ = mockDBQ{}
	s.historyQ.On("GetTx").Return(nil).Once()
	s.historyQ.On("GetLastLedgerIngestNonBlocking", s.ctx).Return(uint32(0), nil).Once()

	for i := uint32(100); i <= uint32(200); i++ {
		s.historyQ.On("Begin").Return(nil).Once()
		s.historyQ.On("GetTx").Return(&sqlx.Tx{}).Once()

		toidFrom := toid.New(int32(i), 0, 0)
		toidTo := toid.New(int32(i+1), 0, 0)
		s.historyQ.On(
			"DeleteRangeAll", s.ctx, toidFrom.ToInt64(), toidTo.ToInt64(),
		).Return(nil).Once()

		meta := xdr.LedgerCloseMeta{
			V0: &xdr.LedgerCloseMetaV0{
				LedgerHeader: xdr.LedgerHeaderHistoryEntry{
					Header: xdr.LedgerHeader{
						LedgerSeq: xdr.Uint32(i),
					},
				},
			},
		}
		s.ledgerBackend.On("GetLedger", s.ctx, uint32(i)).Return(meta, nil).Once()

		s.runner.On("RunTransactionProcessorsOnLedger", meta).Return(
			processors.StatsLedgerTransactionProcessorResults{},
			processorsRunDurations{},
			processors.TradeStats{},
			nil,
		).Once()

		s.historyQ.On("Commit").Return(nil).Once()
		s.historyQ.On("Rollback").Return(nil).Once()
	}
	s.historyQ.On("RebuildTradeAggregationBuckets", s.ctx, uint32(100), uint32(200)).Return(nil).Once()

	err := s.system.ReingestRange([]history.LedgerRange{{100, 200}}, false)
	s.Assert().NoError(err)
}

func (s *ReingestHistoryRangeStateTestSuite) TestSuccessOneLedger() {
	s.historyQ.On("GetLastLedgerIngestNonBlocking", s.ctx).Return(uint32(0), nil).Once()
	s.historyQ.On("GetTx").Return(&sqlx.Tx{}).Once()

	toidFrom := toid.New(100, 0, 0)
	toidTo := toid.New(101, 0, 0)
	s.historyQ.On(
		"DeleteRangeAll", s.ctx, toidFrom.ToInt64(), toidTo.ToInt64(),
	).Return(nil).Once()

	meta := xdr.LedgerCloseMeta{
		V0: &xdr.LedgerCloseMetaV0{
			LedgerHeader: xdr.LedgerHeaderHistoryEntry{
				Header: xdr.LedgerHeader{
					LedgerSeq: xdr.Uint32(100),
				},
			},
		},
	}

	s.runner.On("RunTransactionProcessorsOnLedger", meta).Return(
		processors.StatsLedgerTransactionProcessorResults{},
		processorsRunDurations{},
		processors.TradeStats{},
		nil,
	).Once()
	s.historyQ.On("Commit").Return(nil).Once()
	s.historyQ.On("RebuildTradeAggregationBuckets", s.ctx, uint32(100), uint32(100)).Return(nil).Once()

	// Recreate mock in this single test to remove previous assertion.
	*s.ledgerBackend = mockLedgerBackend{}
	s.ledgerBackend.On("PrepareRange", s.ctx, ledgerbackend.BoundedRange(100, 100)).Return(nil).Once()
	s.ledgerBackend.On("GetLedger", s.ctx, uint32(100)).Return(meta, nil).Once()

	err := s.system.ReingestRange([]history.LedgerRange{{100, 100}}, false)
	s.Assert().NoError(err)
}

func (s *ReingestHistoryRangeStateTestSuite) TestGetLastLedgerIngestError() {
	s.historyQ.On("GetLastLedgerIngest", s.ctx).Return(uint32(0), errors.New("my error")).Once()

	err := s.system.ReingestRange([]history.LedgerRange{{100, 200}}, true)
	s.Assert().EqualError(err, "Error getting last ingested ledger: my error")
}

func (s *ReingestHistoryRangeStateTestSuite) TestReingestRangeForce() {
	s.historyQ.On("GetLastLedgerIngest", s.ctx).Return(uint32(190), nil).Once()

	s.historyQ.On("GetTx").Return(&sqlx.Tx{}).Once()

	toidFrom := toid.New(100, 0, 0)
	toidTo := toid.New(201, 0, 0)
	s.historyQ.On(
		"DeleteRangeAll", s.ctx, toidFrom.ToInt64(), toidTo.ToInt64(),
	).Return(nil).Once()

	for i := 100; i <= 200; i++ {
		meta := xdr.LedgerCloseMeta{
			V0: &xdr.LedgerCloseMetaV0{
				LedgerHeader: xdr.LedgerHeaderHistoryEntry{
					Header: xdr.LedgerHeader{
						LedgerSeq: xdr.Uint32(i),
					},
				},
			},
		}
		s.ledgerBackend.On("GetLedger", s.ctx, uint32(i)).Return(meta, nil).Once()

		s.runner.On("RunTransactionProcessorsOnLedger", meta).Return(
			processors.StatsLedgerTransactionProcessorResults{},
			processorsRunDurations{},
			processors.TradeStats{},
			nil,
		).Once()
	}

	s.historyQ.On("Commit").Return(nil).Once()

	s.historyQ.On("RebuildTradeAggregationBuckets", s.ctx, uint32(100), uint32(200)).Return(nil).Once()

	err := s.system.ReingestRange([]history.LedgerRange{{100, 200}}, true)
	s.Assert().NoError(err)
}
