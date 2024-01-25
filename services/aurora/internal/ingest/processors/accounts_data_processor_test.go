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

func TestAccountsDataProcessorTestSuiteState(t *testing.T) {
	suite.Run(t, new(AccountsDataProcessorTestSuiteState))
}

type AccountsDataProcessorTestSuiteState struct {
	suite.Suite
	ctx       context.Context
	processor *AccountDataProcessor
	mockQ     *history.MockQData
}

func (s *AccountsDataProcessorTestSuiteState) SetupTest() {
	s.ctx = context.Background()
	s.mockQ = &history.MockQData{}

	s.processor = NewAccountDataProcessor(s.mockQ)
}

func (s *AccountsDataProcessorTestSuiteState) TearDownTest() {
	s.Assert().NoError(s.processor.Commit(s.ctx))

	s.mockQ.AssertExpectations(s.T())
}

func (s *AccountsDataProcessorTestSuiteState) TestNoEntries() {
	// Nothing processed, assertions in TearDownTest.
}

func (s *AccountsDataProcessorTestSuiteState) TestCreatesAccounts() {
	data := xdr.DataEntry{
		AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
		DataName:  "test",
		DataValue: []byte{1, 1, 1, 1},
	}
	lastModifiedLedgerSeq := xdr.Uint32(123)
	entry := xdr.LedgerEntry{
		Data: xdr.LedgerEntryData{
			Type: xdr.LedgerEntryTypeData,
			Data: &data,
		},
		LastModifiedLedgerSeq: lastModifiedLedgerSeq,
	}
	historyData := history.Data{
		AccountID:          data.AccountId.Address(),
		Name:               string(data.DataName),
		Value:              history.AccountDataValue(data.DataValue),
		LastModifiedLedger: uint32(entry.LastModifiedLedgerSeq),
	}
	s.mockQ.On("UpsertAccountData", s.ctx, []history.Data{historyData}).Return(nil).Once()

	err := s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeData,
		Pre:  nil,
		Post: &entry,
	})
	s.Assert().NoError(err)
}

func TestAccountsDataProcessorTestSuiteLedger(t *testing.T) {
	suite.Run(t, new(AccountsDataProcessorTestSuiteLedger))
}

type AccountsDataProcessorTestSuiteLedger struct {
	suite.Suite
	ctx       context.Context
	processor *AccountDataProcessor
	mockQ     *history.MockQData
}

func (s *AccountsDataProcessorTestSuiteLedger) SetupTest() {
	s.ctx = context.Background()
	s.mockQ = &history.MockQData{}

	s.processor = NewAccountDataProcessor(s.mockQ)
}

func (s *AccountsDataProcessorTestSuiteLedger) TearDownTest() {
	s.Assert().NoError(s.processor.Commit(s.ctx))
	s.mockQ.AssertExpectations(s.T())
}

func (s *AccountsDataProcessorTestSuiteLedger) TestNoTransactions() {
	// Nothing processed, assertions in TearDownTest.
}

func (s *AccountsDataProcessorTestSuiteLedger) TestNewAccountData() {
	data := xdr.DataEntry{
		AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
		DataName:  "test",
		DataValue: []byte{1, 1, 1, 1},
	}
	lastModifiedLedgerSeq := xdr.Uint32(123)

	err := s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeData,
		Pre:  nil,
		Post: &xdr.LedgerEntry{
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeData,
				Data: &data,
			},
			LastModifiedLedgerSeq: lastModifiedLedgerSeq,
		},
	})
	s.Assert().NoError(err)

	updatedData := xdr.DataEntry{
		AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
		DataName:  "test",
		DataValue: []byte{2, 2, 2, 2},
	}

	updatedEntry := xdr.LedgerEntry{
		Data: xdr.LedgerEntryData{
			Type: xdr.LedgerEntryTypeData,
			Data: &updatedData,
		},
		LastModifiedLedgerSeq: lastModifiedLedgerSeq,
	}

	err = s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeData,
		Pre: &xdr.LedgerEntry{
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeData,
				Data: &data,
			},
			LastModifiedLedgerSeq: lastModifiedLedgerSeq - 1,
		},
		Post: &updatedEntry,
	})
	s.Assert().NoError(err)

	// We use LedgerEntryChangesCache so all changes are squashed
	historyData := history.Data{
		AccountID:          updatedData.AccountId.Address(),
		Name:               string(updatedData.DataName),
		Value:              history.AccountDataValue(updatedData.DataValue),
		LastModifiedLedger: uint32(updatedEntry.LastModifiedLedgerSeq),
	}
	s.mockQ.On("UpsertAccountData", s.ctx, []history.Data{historyData}).Return(nil).Once()
}

func (s *AccountsDataProcessorTestSuiteLedger) TestUpdateAccountData() {
	data := xdr.DataEntry{
		AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
		DataName:  "test",
		DataValue: []byte{1, 1, 1, 1},
	}
	lastModifiedLedgerSeq := xdr.Uint32(123)

	updatedData := xdr.DataEntry{
		AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
		DataName:  "test",
		DataValue: []byte{2, 2, 2, 2},
	}

	updatedEntry := xdr.LedgerEntry{
		Data: xdr.LedgerEntryData{
			Type: xdr.LedgerEntryTypeData,
			Data: &updatedData,
		},
		LastModifiedLedgerSeq: lastModifiedLedgerSeq,
	}

	err := s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeData,
		Pre: &xdr.LedgerEntry{
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeData,
				Data: &data,
			},
			LastModifiedLedgerSeq: lastModifiedLedgerSeq - 1,
		},
		Post: &updatedEntry,
	})
	s.Assert().NoError(err)

	historyData := history.Data{
		AccountID:          updatedData.AccountId.Address(),
		Name:               string(updatedData.DataName),
		Value:              history.AccountDataValue(updatedData.DataValue),
		LastModifiedLedger: uint32(updatedEntry.LastModifiedLedgerSeq),
	}
	s.mockQ.On("UpsertAccountData", s.ctx, []history.Data{historyData}).Return(nil).Once()
}

func (s *AccountsDataProcessorTestSuiteLedger) TestRemoveAccountData() {
	err := s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeData,
		Pre: &xdr.LedgerEntry{
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeData,
				Data: &xdr.DataEntry{
					AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					DataName:  "test",
					DataValue: []byte{1, 1, 1, 1},
				},
			},
		},
		Post: nil,
	})
	s.Assert().NoError(err)

	s.mockQ.On(
		"RemoveAccountData",
		s.ctx,
		[]history.AccountDataKey{
			{
				AccountID: "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML",
				DataName:  "test",
			},
		},
	).Return(int64(1), nil).Once()
}
