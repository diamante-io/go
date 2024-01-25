//lint:file-ignore U1001 Ignore all unused code, staticcheck doesn't understand testify/suite

package processors

import (
	"context"
	"testing"

	"go/ingest"
	"go/services/aurora/internal/db2/history"
	"go/xdr"

	"github.com/stretchr/testify/suite"
)

func TestAccountsProcessorTestSuiteState(t *testing.T) {
	suite.Run(t, new(AccountsProcessorTestSuiteState))
}

type AccountsProcessorTestSuiteState struct {
	suite.Suite
	ctx       context.Context
	processor *AccountsProcessor
	mockQ     *history.MockQAccounts
}

func (s *AccountsProcessorTestSuiteState) SetupTest() {
	s.ctx = context.Background()
	s.mockQ = &history.MockQAccounts{}

	s.processor = NewAccountsProcessor(s.mockQ)
}

func (s *AccountsProcessorTestSuiteState) TearDownTest() {
	s.Assert().NoError(s.processor.Commit(s.ctx))
	s.mockQ.AssertExpectations(s.T())
}

func (s *AccountsProcessorTestSuiteState) TestNoEntries() {
	// Nothing processed, assertions in TearDownTest.
}

func (s *AccountsProcessorTestSuiteState) TestCreatesAccounts() {
	// We use LedgerEntryChangesCache so all changes are squashed
	s.mockQ.On(
		"UpsertAccounts", s.ctx,
		[]history.AccountEntry{
			{
				LastModifiedLedger: 123,
				AccountID:          "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML",
				MasterWeight:       1,
				ThresholdLow:       1,
				ThresholdMedium:    1,
				ThresholdHigh:      1,
			},
		},
	).Return(nil).Once()

	err := s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeAccount,
		Pre:  nil,
		Post: &xdr.LedgerEntry{
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId:  xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					Thresholds: [4]byte{1, 1, 1, 1},
				},
			},
			LastModifiedLedgerSeq: xdr.Uint32(123),
		},
	})
	s.Assert().NoError(err)
}

func TestAccountsProcessorTestSuiteLedger(t *testing.T) {
	suite.Run(t, new(AccountsProcessorTestSuiteLedger))
}

type AccountsProcessorTestSuiteLedger struct {
	suite.Suite
	ctx       context.Context
	processor *AccountsProcessor
	mockQ     *history.MockQAccounts
}

func (s *AccountsProcessorTestSuiteLedger) SetupTest() {
	s.ctx = context.Background()
	s.mockQ = &history.MockQAccounts{}

	s.processor = NewAccountsProcessor(s.mockQ)
}

func (s *AccountsProcessorTestSuiteLedger) TearDownTest() {
	s.Assert().NoError(s.processor.Commit(s.ctx))
	s.mockQ.AssertExpectations(s.T())
}

func (s *AccountsProcessorTestSuiteLedger) TestNoTransactions() {
	// Nothing processed, assertions in TearDownTest.
}

func (s *AccountsProcessorTestSuiteLedger) TestNewAccount() {
	account := xdr.AccountEntry{
		AccountId:  xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
		Thresholds: [4]byte{1, 1, 1, 1},
	}
	lastModifiedLedgerSeq := xdr.Uint32(123)

	err := s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeAccount,
		Pre:  nil,
		Post: &xdr.LedgerEntry{
			Data: xdr.LedgerEntryData{
				Type:    xdr.LedgerEntryTypeAccount,
				Account: &account,
			},
			LastModifiedLedgerSeq: lastModifiedLedgerSeq,
		},
	})
	s.Assert().NoError(err)

	updatedAccount := xdr.AccountEntry{
		AccountId:  xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
		Thresholds: [4]byte{0, 1, 2, 3},
		HomeDomain: "diamcircle.org",
	}

	err = s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeAccount,
		Pre: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: lastModifiedLedgerSeq - 1,
			Data: xdr.LedgerEntryData{
				Type:    xdr.LedgerEntryTypeAccount,
				Account: &account,
			},
		},
		Post: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: lastModifiedLedgerSeq,
			Data: xdr.LedgerEntryData{
				Type:    xdr.LedgerEntryTypeAccount,
				Account: &updatedAccount,
			},
		},
	})
	s.Assert().NoError(err)

	// We use LedgerEntryChangesCache so all changes are squashed
	s.mockQ.On(
		"UpsertAccounts",
		s.ctx,
		[]history.AccountEntry{
			{
				LastModifiedLedger: 123,
				AccountID:          "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML",
				MasterWeight:       0,
				ThresholdLow:       1,
				ThresholdMedium:    2,
				ThresholdHigh:      3,
				HomeDomain:         "diamcircle.org",
			},
		},
	).Return(nil).Once()
}

func (s *AccountsProcessorTestSuiteLedger) TestRemoveAccount() {
	s.mockQ.On(
		"RemoveAccounts",
		s.ctx,
		[]string{"GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"},
	).Return(int64(1), nil).Once()

	err := s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeAccount,
		Pre: &xdr.LedgerEntry{
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId:  xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					Thresholds: [4]byte{1, 1, 1, 1},
				},
			},
		},
		Post: nil,
	})
	s.Assert().NoError(err)
}

func (s *AccountsProcessorTestSuiteLedger) TestProcessUpgradeChange() {
	account := xdr.AccountEntry{
		AccountId:  xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
		Thresholds: [4]byte{1, 1, 1, 1},
	}
	lastModifiedLedgerSeq := xdr.Uint32(123)

	err := s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeAccount,
		Pre:  nil,
		Post: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: lastModifiedLedgerSeq,
			Data: xdr.LedgerEntryData{
				Type:    xdr.LedgerEntryTypeAccount,
				Account: &account,
			},
		},
	})
	s.Assert().NoError(err)

	updatedAccount := xdr.AccountEntry{
		AccountId:  xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
		Thresholds: [4]byte{0, 1, 2, 3},
		HomeDomain: "diamcircle.org",
	}

	err = s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeAccount,
		Pre: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: lastModifiedLedgerSeq,
			Data: xdr.LedgerEntryData{
				Type:    xdr.LedgerEntryTypeAccount,
				Account: &account,
			},
		},
		Post: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: lastModifiedLedgerSeq + 1,
			Data: xdr.LedgerEntryData{
				Type:    xdr.LedgerEntryTypeAccount,
				Account: &updatedAccount,
			},
		},
	})
	s.Assert().NoError(err)

	s.mockQ.On(
		"UpsertAccounts",
		s.ctx,
		[]history.AccountEntry{
			{
				LastModifiedLedger: uint32(lastModifiedLedgerSeq) + 1,
				AccountID:          "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML",
				MasterWeight:       0,
				ThresholdLow:       1,
				ThresholdMedium:    2,
				ThresholdHigh:      3,
				HomeDomain:         "diamcircle.org",
			},
		},
	).Return(nil).Once()
}

func (s *AccountsProcessorTestSuiteLedger) TestFeeProcessedBeforeEverythingElse() {
	err := s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeAccount,
		Pre: &xdr.LedgerEntry{
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GAHK7EEG2WWHVKDNT4CEQFZGKF2LGDSW2IVM4S5DP42RBW3K6BTODB4A"),
					Balance:   200,
				},
			},
		},
		Post: &xdr.LedgerEntry{
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GAHK7EEG2WWHVKDNT4CEQFZGKF2LGDSW2IVM4S5DP42RBW3K6BTODB4A"),
					Balance:   100,
				},
			},
		},
	})
	s.Assert().NoError(err)

	err = s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeAccount,
		Pre: &xdr.LedgerEntry{
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GAHK7EEG2WWHVKDNT4CEQFZGKF2LGDSW2IVM4S5DP42RBW3K6BTODB4A"),
					Balance:   100,
				},
			},
		},
		Post: &xdr.LedgerEntry{
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GAHK7EEG2WWHVKDNT4CEQFZGKF2LGDSW2IVM4S5DP42RBW3K6BTODB4A"),
					Balance:   300,
				},
			},
		},
	})
	s.Assert().NoError(err)

	s.mockQ.On(
		"UpsertAccounts",
		s.ctx,
		[]history.AccountEntry{
			{
				LastModifiedLedger: 0,
				AccountID:          "GAHK7EEG2WWHVKDNT4CEQFZGKF2LGDSW2IVM4S5DP42RBW3K6BTODB4A",
				Balance:            300,
			},
		},
	).Return(nil).Once()
}
