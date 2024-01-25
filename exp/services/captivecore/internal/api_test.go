package internal

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"go/ingest/ledgerbackend"
	"go/support/log"
	"go/xdr"
)

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

type APITestSuite struct {
	suite.Suite
	ctx           context.Context
	ledgerBackend *ledgerbackend.MockDatabaseBackend
	api           CaptiveCoreAPI
}

func (s *APITestSuite) SetupTest() {
	s.ctx = context.Background()
	s.ledgerBackend = &ledgerbackend.MockDatabaseBackend{}
	s.api = NewCaptiveCoreAPI(s.ledgerBackend, log.New())
}

func (s *APITestSuite) TearDownTest() {
	s.ledgerBackend.AssertExpectations(s.T())
}

func (s *APITestSuite) TestLatestSeqActiveRequestInvalid() {
	_, err := s.api.GetLatestLedgerSequence(s.ctx)
	s.Assert().Equal(err, ErrMissingPrepareRange)
}

func (s *APITestSuite) TestGetLedgerActiveRequestInvalid() {
	_, err := s.api.GetLedger(s.ctx, 64)
	s.Assert().Equal(err, ErrMissingPrepareRange)
}

func (s *APITestSuite) runBeforeReady(prepareRangeErr error, f func()) {
	waitChan := make(chan time.Time)
	ledgerRange := ledgerbackend.UnboundedRange(63)
	s.ledgerBackend.On("PrepareRange", mock.Anything, ledgerRange).
		WaitUntil(waitChan).
		Return(prepareRangeErr).Once()

	response, err := s.api.PrepareRange(s.ctx, ledgerRange)
	s.Assert().NoError(err)
	s.Assert().False(response.Ready)
	s.Assert().Equal(response.LedgerRange, ledgerRange)

	f()

	close(waitChan)
	s.api.wg.Wait()
}

func (s *APITestSuite) TestLatestSeqActiveRequestNotReady() {
	s.runBeforeReady(nil, func() {
		_, err := s.api.GetLatestLedgerSequence(s.ctx)
		s.Assert().Equal(err, ErrPrepareRangeNotReady)
	})
}

func (s *APITestSuite) TestGetLedgerNotReady() {
	s.runBeforeReady(nil, func() {
		_, err := s.api.GetLedger(s.ctx, 64)
		s.Assert().Equal(err, ErrPrepareRangeNotReady)
	})
}

func (s *APITestSuite) waitUntilReady(ledgerRange ledgerbackend.Range) {
	s.ledgerBackend.On("PrepareRange", mock.Anything, ledgerRange).
		Return(nil).Once()

	response, err := s.api.PrepareRange(s.ctx, ledgerRange)
	s.Assert().NoError(err)
	s.Assert().False(response.Ready)
	s.Assert().Equal(response.LedgerRange, ledgerRange)

	s.api.wg.Wait()
}

func (s *APITestSuite) TestLatestSeqError() {
	s.waitUntilReady(ledgerbackend.UnboundedRange(63))

	expectedErr := fmt.Errorf("test error")
	s.ledgerBackend.On("GetLatestLedgerSequence", s.ctx).Return(uint32(0), expectedErr).Once()

	_, err := s.api.GetLatestLedgerSequence(s.ctx)
	s.Assert().Equal(err, expectedErr)
}

func (s *APITestSuite) TestGetLedgerError() {
	s.waitUntilReady(ledgerbackend.UnboundedRange(63))

	expectedErr := fmt.Errorf("test error")
	s.ledgerBackend.On("GetLedger", s.ctx, uint32(64)).
		Return(xdr.LedgerCloseMeta{}, expectedErr).Once()

	_, err := s.api.GetLedger(s.ctx, 64)
	s.Assert().Equal(err, expectedErr)
}

func (s *APITestSuite) TestLatestSeqSucceeds() {
	s.waitUntilReady(ledgerbackend.UnboundedRange(63))

	expectedSeq := uint32(100)
	s.ledgerBackend.On("GetLatestLedgerSequence", s.ctx).Return(expectedSeq, nil).Once()
	seq, err := s.api.GetLatestLedgerSequence(s.ctx)
	s.Assert().NoError(err)
	s.Assert().Equal(seq, ledgerbackend.LatestLedgerSequenceResponse{Sequence: expectedSeq})
}

func (s *APITestSuite) TestGetLedgerSucceeds() {
	s.waitUntilReady(ledgerbackend.UnboundedRange(63))

	expectedLedger := xdr.LedgerCloseMeta{
		V0: &xdr.LedgerCloseMetaV0{
			LedgerHeader: xdr.LedgerHeaderHistoryEntry{
				Header: xdr.LedgerHeader{
					LedgerSeq: 64,
				},
			},
		},
	}
	s.ledgerBackend.On("GetLedger", s.ctx, uint32(64)).
		Return(expectedLedger, nil).Once()
	seq, err := s.api.GetLedger(s.ctx, 64)

	s.Assert().NoError(err)
	s.Assert().Equal(seq, ledgerbackend.LedgerResponse{
		Ledger: ledgerbackend.Base64Ledger(expectedLedger),
	})
}

func (s *APITestSuite) TestShutDownBeforePrepareRange() {
	s.ledgerBackend.On("Close").Return(nil).Once()
	s.api.Shutdown()
	_, err := s.api.PrepareRange(s.ctx, ledgerbackend.UnboundedRange(63))
	s.Assert().EqualError(err, "Cannot prepare range when shut down")
}

func (s *APITestSuite) TestShutDownDuringPrepareRange() {
	s.runBeforeReady(nil, func() {
		s.api.cancel()
	})

	s.Assert().False(s.api.activeRequest.ready)
}

func (s *APITestSuite) TestPrepareRangeInvalidActiveRequest() {
	s.runBeforeReady(nil, func() {
		s.Assert().True(s.api.activeRequest.valid)
		s.api.activeRequest.valid = false
	})
	s.Assert().False(s.api.activeRequest.ready)

	s.api.activeRequest = &rangeRequest{}

	s.runBeforeReady(fmt.Errorf("with error"), func() {
		s.Assert().True(s.api.activeRequest.valid)
		s.api.activeRequest.valid = false
	})
	s.Assert().False(s.api.activeRequest.ready)
}

func (s *APITestSuite) TestPrepareRangeDoesNotMatchActiveRequestRange() {
	s.runBeforeReady(nil, func() {
		s.Assert().Equal(ledgerbackend.UnboundedRange(63), s.api.activeRequest.ledgerRange)
		s.api.activeRequest.ledgerRange = ledgerbackend.UnboundedRange(1000)
	})
	s.Assert().False(s.api.activeRequest.ready)
	s.Assert().Equal(ledgerbackend.UnboundedRange(1000), s.api.activeRequest.ledgerRange)

	s.api.activeRequest = &rangeRequest{}

	s.runBeforeReady(fmt.Errorf("with error"), func() {
		s.Assert().Equal(ledgerbackend.UnboundedRange(63), s.api.activeRequest.ledgerRange)
		s.api.activeRequest.ledgerRange = ledgerbackend.UnboundedRange(10)
	})
	s.Assert().False(s.api.activeRequest.ready)
	s.Assert().Equal(ledgerbackend.UnboundedRange(10), s.api.activeRequest.ledgerRange)
}

func (s *APITestSuite) TestPrepareRangeActiveRequestReady() {
	s.runBeforeReady(nil, func() {
		s.api.activeRequest.ready = true
	})
	s.Assert().True(s.api.activeRequest.ready)
	s.Assert().True(s.api.activeRequest.valid)
	s.Assert().Equal(0, s.api.activeRequest.readyDuration)

	s.api.activeRequest = &rangeRequest{}

	s.runBeforeReady(fmt.Errorf("with error"), func() {
		s.api.activeRequest.ready = true
	})
	s.Assert().True(s.api.activeRequest.ready)
	s.Assert().True(s.api.activeRequest.valid)
	s.Assert().Equal(0, s.api.activeRequest.readyDuration)
}

func (s *APITestSuite) TestPrepareRangeError() {
	s.runBeforeReady(fmt.Errorf("with error"), func() {
		s.Assert().False(s.api.activeRequest.ready)
		s.Assert().True(s.api.activeRequest.valid)
	})
	s.Assert().False(s.api.activeRequest.ready)
	s.Assert().False(s.api.activeRequest.valid)

	s.api.activeRequest = &rangeRequest{}
}

func (s *APITestSuite) TestRangeAlreadyPrepared() {
	superSetRange := ledgerbackend.UnboundedRange(63)
	s.waitUntilReady(superSetRange)

	for _, ledgerRange := range []ledgerbackend.Range{
		superSetRange,
		ledgerbackend.UnboundedRange(100),
		ledgerbackend.BoundedRange(63, 70),
	} {
		response, err := s.api.PrepareRange(s.ctx, ledgerRange)
		s.Assert().NoError(err)
		s.Assert().True(response.Ready)
		s.Assert().Equal(superSetRange, response.LedgerRange)
	}
}

func (s *APITestSuite) TestNewPrepareRange() {
	s.waitUntilReady(ledgerbackend.UnboundedRange(63))
	s.waitUntilReady(ledgerbackend.UnboundedRange(50))
	s.waitUntilReady(ledgerbackend.BoundedRange(45, 50))
	s.waitUntilReady(ledgerbackend.UnboundedRange(46))
}
