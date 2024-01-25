//lint:file-ignore U1001 Ignore all unused code, staticcheck doesn't understand testify/suite

package processors

import (
	"context"
	"testing"

	"github.com/guregu/null"
	"github.com/stretchr/testify/suite"

	"github.com/diamcircle/go/ingest"
	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/xdr"
)

func TestClaimableBalancesChangeProcessorTestSuiteState(t *testing.T) {
	suite.Run(t, new(ClaimableBalancesChangeProcessorTestSuiteState))
}

type ClaimableBalancesChangeProcessorTestSuiteState struct {
	suite.Suite
	ctx       context.Context
	processor *ClaimableBalancesChangeProcessor
	mockQ     *history.MockQClaimableBalances
}

func (s *ClaimableBalancesChangeProcessorTestSuiteState) SetupTest() {
	s.ctx = context.Background()
	s.mockQ = &history.MockQClaimableBalances{}

	s.processor = NewClaimableBalancesChangeProcessor(s.mockQ)
}

func (s *ClaimableBalancesChangeProcessorTestSuiteState) TearDownTest() {
	s.Assert().NoError(s.processor.Commit(s.ctx))
	s.mockQ.AssertExpectations(s.T())
}

func (s *ClaimableBalancesChangeProcessorTestSuiteState) TestNoEntries() {
	// Nothing processed, assertions in TearDownTest.
}

func (s *ClaimableBalancesChangeProcessorTestSuiteState) TestCreatesClaimableBalances() {
	lastModifiedLedgerSeq := xdr.Uint32(123)
	balanceID := xdr.ClaimableBalanceId{
		Type: xdr.ClaimableBalanceIdTypeClaimableBalanceIdTypeV0,
		V0:   &xdr.Hash{1, 2, 3},
	}

	cBalance := xdr.ClaimableBalanceEntry{
		BalanceId: balanceID,
		Claimants: []xdr.Claimant{},
		Asset:     xdr.MustNewCreditAsset("USD", "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
		Amount:    10,
	}
	id, err := xdr.MarshalHex(balanceID)
	s.Assert().NoError(err)
	s.mockQ.On("UpsertClaimableBalances", s.ctx, []history.ClaimableBalance{
		{
			BalanceID:          id,
			Claimants:          []history.Claimant{},
			Asset:              cBalance.Asset,
			Amount:             cBalance.Amount,
			LastModifiedLedger: uint32(lastModifiedLedgerSeq),
		},
	}).Return(nil).Once()

	err = s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeClaimableBalance,
		Pre:  nil,
		Post: &xdr.LedgerEntry{
			Data: xdr.LedgerEntryData{
				Type:             xdr.LedgerEntryTypeClaimableBalance,
				ClaimableBalance: &cBalance,
			},
			LastModifiedLedgerSeq: lastModifiedLedgerSeq,
		},
	})
	s.Assert().NoError(err)
}

func TestClaimableBalancesChangeProcessorTestSuiteLedger(t *testing.T) {
	suite.Run(t, new(ClaimableBalancesChangeProcessorTestSuiteLedger))
}

type ClaimableBalancesChangeProcessorTestSuiteLedger struct {
	suite.Suite
	ctx       context.Context
	processor *ClaimableBalancesChangeProcessor
	mockQ     *history.MockQClaimableBalances
}

func (s *ClaimableBalancesChangeProcessorTestSuiteLedger) SetupTest() {
	s.ctx = context.Background()
	s.mockQ = &history.MockQClaimableBalances{}

	s.processor = NewClaimableBalancesChangeProcessor(s.mockQ)
}

func (s *ClaimableBalancesChangeProcessorTestSuiteLedger) TearDownTest() {
	s.Assert().NoError(s.processor.Commit(s.ctx))
	s.mockQ.AssertExpectations(s.T())
}

func (s *ClaimableBalancesChangeProcessorTestSuiteLedger) TestNoTransactions() {
	// Nothing processed, assertions in TearDownTest.
}

func (s *ClaimableBalancesChangeProcessorTestSuiteLedger) TestNewClaimableBalance() {
	lastModifiedLedgerSeq := xdr.Uint32(123)
	balanceID := xdr.ClaimableBalanceId{
		Type: xdr.ClaimableBalanceIdTypeClaimableBalanceIdTypeV0,
		V0:   &xdr.Hash{1, 2, 3},
	}
	cBalance := xdr.ClaimableBalanceEntry{
		BalanceId: balanceID,
		Claimants: []xdr.Claimant{},
		Asset:     xdr.MustNewCreditAsset("USD", "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
		Amount:    10,
	}
	entry := xdr.LedgerEntry{
		Data: xdr.LedgerEntryData{
			Type:             xdr.LedgerEntryTypeClaimableBalance,
			ClaimableBalance: &cBalance,
		},
		LastModifiedLedgerSeq: lastModifiedLedgerSeq,
		Ext: xdr.LedgerEntryExt{
			V: 1,
			V1: &xdr.LedgerEntryExtensionV1{
				SponsoringId: nil,
			},
		},
	}
	err := s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeClaimableBalance,
		Pre:  nil,
		Post: &entry,
	})
	s.Assert().NoError(err)

	// add sponsor
	updated := xdr.LedgerEntry{
		Data: xdr.LedgerEntryData{
			Type:             xdr.LedgerEntryTypeClaimableBalance,
			ClaimableBalance: &cBalance,
		},
		LastModifiedLedgerSeq: lastModifiedLedgerSeq,
		Ext: xdr.LedgerEntryExt{
			V: 1,
			V1: &xdr.LedgerEntryExtensionV1{
				SponsoringId: xdr.MustAddressPtr("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
			},
		},
	}

	entry.LastModifiedLedgerSeq = entry.LastModifiedLedgerSeq - 1
	err = s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeClaimableBalance,
		Pre:  &entry,
		Post: &updated,
	})
	s.Assert().NoError(err)

	id, err := xdr.MarshalHex(balanceID)
	s.Assert().NoError(err)
	// We use LedgerEntryChangesCache so all changes are squashed
	s.mockQ.On(
		"UpsertClaimableBalances",
		s.ctx,
		[]history.ClaimableBalance{
			{
				BalanceID:          id,
				Claimants:          []history.Claimant{},
				Asset:              cBalance.Asset,
				Amount:             cBalance.Amount,
				LastModifiedLedger: uint32(lastModifiedLedgerSeq),
				Sponsor:            null.StringFrom("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
			},
		},
	).Return(nil).Once()
}

func (s *ClaimableBalancesChangeProcessorTestSuiteLedger) TestUpdateClaimableBalance() {
	balanceID := xdr.ClaimableBalanceId{
		Type: xdr.ClaimableBalanceIdTypeClaimableBalanceIdTypeV0,
		V0:   &xdr.Hash{1, 2, 3},
	}
	cBalance := xdr.ClaimableBalanceEntry{
		BalanceId: balanceID,
		Claimants: []xdr.Claimant{},
		Asset:     xdr.MustNewCreditAsset("USD", "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
		Amount:    10,
	}
	lastModifiedLedgerSeq := xdr.Uint32(123)

	pre := xdr.LedgerEntry{
		Data: xdr.LedgerEntryData{
			Type:             xdr.LedgerEntryTypeClaimableBalance,
			ClaimableBalance: &cBalance,
		},
		LastModifiedLedgerSeq: lastModifiedLedgerSeq - 1,
		Ext: xdr.LedgerEntryExt{
			V: 1,
			V1: &xdr.LedgerEntryExtensionV1{
				SponsoringId: nil,
			},
		},
	}

	// add sponsor
	updated := xdr.LedgerEntry{
		Data: xdr.LedgerEntryData{
			Type:             xdr.LedgerEntryTypeClaimableBalance,
			ClaimableBalance: &cBalance,
		},
		LastModifiedLedgerSeq: lastModifiedLedgerSeq,
		Ext: xdr.LedgerEntryExt{
			V: 1,
			V1: &xdr.LedgerEntryExtensionV1{
				SponsoringId: xdr.MustAddressPtr("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
			},
		},
	}

	err := s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeClaimableBalance,
		Pre:  &pre,
		Post: &updated,
	})
	s.Assert().NoError(err)

	id, err := xdr.MarshalHex(balanceID)
	s.Assert().NoError(err)
	s.mockQ.On(
		"UpsertClaimableBalances",
		s.ctx,
		[]history.ClaimableBalance{
			{
				BalanceID:          id,
				Claimants:          []history.Claimant{},
				Asset:              cBalance.Asset,
				Amount:             cBalance.Amount,
				LastModifiedLedger: uint32(lastModifiedLedgerSeq),
				Sponsor:            null.StringFrom("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
			},
		},
	).Return(nil).Once()
}

func (s *ClaimableBalancesChangeProcessorTestSuiteLedger) TestRemoveClaimableBalance() {
	balanceID := xdr.ClaimableBalanceId{
		Type: xdr.ClaimableBalanceIdTypeClaimableBalanceIdTypeV0,
		V0:   &xdr.Hash{1, 2, 3},
	}
	cBalance := xdr.ClaimableBalanceEntry{
		BalanceId: balanceID,
		Claimants: []xdr.Claimant{},
		Asset:     xdr.MustNewCreditAsset("USD", "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
		Amount:    10,
	}
	lastModifiedLedgerSeq := xdr.Uint32(123)
	pre := xdr.LedgerEntry{
		Data: xdr.LedgerEntryData{
			Type:             xdr.LedgerEntryTypeClaimableBalance,
			ClaimableBalance: &cBalance,
		},
		LastModifiedLedgerSeq: lastModifiedLedgerSeq - 1,
		Ext: xdr.LedgerEntryExt{
			V: 1,
			V1: &xdr.LedgerEntryExtensionV1{
				SponsoringId: nil,
			},
		},
	}
	err := s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeClaimableBalance,
		Pre:  &pre,
		Post: nil,
	})
	s.Assert().NoError(err)

	id, err := xdr.MarshalHex(balanceID)
	s.Assert().NoError(err)
	s.mockQ.On(
		"RemoveClaimableBalances",
		s.ctx,
		[]string{id},
	).Return(int64(1), nil).Once()
}
