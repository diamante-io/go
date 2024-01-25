//lint:file-ignore U1001 Ignore all unused code, staticcheck doesn't understand testify/suite

package ingest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/diamcircle/go/support/errors"
)

func TestStressTestStateTestSuite(t *testing.T) {
	suite.Run(t, new(StressTestStateTestSuite))
}

type StressTestStateTestSuite struct {
	suite.Suite
	ctx            context.Context
	historyQ       *mockDBQ
	historyAdapter *mockHistoryArchiveAdapter
	runner         *mockProcessorsRunner
	system         *system
}

func (s *StressTestStateTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.historyQ = &mockDBQ{}
	s.historyAdapter = &mockHistoryArchiveAdapter{}
	s.runner = &mockProcessorsRunner{}
	s.system = &system{
		ctx:            s.ctx,
		historyQ:       s.historyQ,
		historyAdapter: s.historyAdapter,
		runner:         s.runner,
	}
	s.system.initMetrics()

	s.historyQ.On("GetTx").Return(nil).Once()
	s.historyQ.On("Rollback").Return(nil).Once()
	s.runner.On("EnableMemoryStatsLogging").Return()
}

func (s *StressTestStateTestSuite) TearDownTest() {
	t := s.T()
	s.historyQ.AssertExpectations(t)
	s.historyAdapter.AssertExpectations(t)
	s.runner.AssertExpectations(t)
}

func (s *StressTestStateTestSuite) TestBounds() {
	*s.historyQ = mockDBQ{}
	*s.runner = mockProcessorsRunner{}

	err := s.system.StressTest(-1, 4)
	s.Assert().EqualError(err, "transactions must be positive")

	err = s.system.StressTest(0, 4)
	s.Assert().EqualError(err, "transactions must be positive")

	err = s.system.StressTest(100, -2)
	s.Assert().EqualError(err, "changes per transaction must be positive")
}

func (s *StressTestStateTestSuite) TestBeginReturnsError() {
	*s.historyQ = mockDBQ{}
	s.historyQ.On("GetTx").Return(nil).Once()
	s.historyQ.On("Begin").Return(errors.New("my error")).Once()

	err := s.system.StressTest(10, 4)
	s.Assert().EqualError(err, "Error starting a transaction: my error")
}

func (s *StressTestStateTestSuite) TestGetLastLedgerIngestReturnsError() {
	s.historyQ.On("Begin").Return(nil).Once()
	s.historyQ.On("GetLastLedgerIngest", s.ctx).Return(uint32(0), errors.New("my error")).Once()

	err := s.system.StressTest(10, 4)
	s.Assert().EqualError(err, "Error getting last ingested ledger: my error")
}

func (s *StressTestStateTestSuite) TestGetLastLedgerIngestNonEmpty() {
	s.historyQ.On("Begin").Return(nil).Once()
	s.historyQ.On("GetLastLedgerIngest", s.ctx).Return(uint32(100), nil).Once()

	err := s.system.StressTest(10, 4)
	s.Assert().EqualError(err, "Database not empty")
}

func (s *StressTestStateTestSuite) TestRunAllProcessorsOnLedgerReturnsError() {
	s.historyQ.On("Begin").Return(nil).Once()
	s.historyQ.On("GetLastLedgerIngest", s.ctx).Return(uint32(0), nil).Once()

	s.runner.On("RunAllProcessorsOnLedger", mock.AnythingOfType("xdr.LedgerCloseMeta")).Return(
		ledgerStats{},
		errors.New("my error"),
	).Once()

	err := s.system.StressTest(10, 4)
	s.Assert().EqualError(err, "Error running processors on ledger: my error")
}

func (s *StressTestStateTestSuite) TestUpdateLastLedgerIngestReturnsError() {
	s.historyQ.On("Begin").Return(nil).Once()
	s.historyQ.On("GetLastLedgerIngest", s.ctx).Return(uint32(0), nil).Once()
	s.runner.On("RunAllProcessorsOnLedger", mock.AnythingOfType("xdr.LedgerCloseMeta")).Return(
		ledgerStats{},
		nil,
	).Once()
	s.historyQ.On("UpdateLastLedgerIngest", s.ctx, uint32(1)).Return(errors.New("my error")).Once()

	err := s.system.StressTest(10, 4)
	s.Assert().EqualError(err, "Error updating last ingested ledger: my error")
}

func (s *StressTestStateTestSuite) TestCommitReturnsError() {
	s.historyQ.On("Begin").Return(nil).Once()
	s.historyQ.On("GetLastLedgerIngest", s.ctx).Return(uint32(0), nil).Once()
	s.runner.On("RunAllProcessorsOnLedger", mock.AnythingOfType("xdr.LedgerCloseMeta")).Return(
		ledgerStats{},
		nil,
	).Once()
	s.historyQ.On("UpdateLastLedgerIngest", s.ctx, uint32(1)).Return(nil).Once()
	s.historyQ.On("Commit").Return(errors.New("my error")).Once()

	err := s.system.StressTest(10, 4)
	s.Assert().EqualError(err, "Error committing db transaction: my error")
}

func (s *StressTestStateTestSuite) TestSucceeds() {
	s.historyQ.On("Begin").Return(nil).Once()
	s.historyQ.On("GetLastLedgerIngest", s.ctx).Return(uint32(0), nil).Once()
	s.runner.On("RunAllProcessorsOnLedger", mock.AnythingOfType("xdr.LedgerCloseMeta")).Return(
		ledgerStats{},
		nil,
	).Once()
	s.historyQ.On("UpdateLastLedgerIngest", s.ctx, uint32(1)).Return(nil).Once()
	s.historyQ.On("Commit").Return(nil).Once()

	err := s.system.StressTest(10, 4)
	s.Assert().NoError(err)
}
