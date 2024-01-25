//lint:file-ignore U1001 Ignore all unused code, staticcheck doesn't understand testify/suite

package processors

import (
	"context"
	"testing"

	"github.com/diamcircle/go/ingest"
	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/xdr"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestOffersProcessorTestSuiteState(t *testing.T) {
	suite.Run(t, new(OffersProcessorTestSuiteState))
}

type OffersProcessorTestSuiteState struct {
	suite.Suite
	ctx       context.Context
	processor *OffersProcessor
	mockQ     *history.MockQOffers
	sequence  uint32
}

func (s *OffersProcessorTestSuiteState) SetupTest() {
	s.ctx = context.Background()
	s.mockQ = &history.MockQOffers{}

	s.sequence = 456
	s.processor = NewOffersProcessor(s.mockQ, s.sequence)
}

func (s *OffersProcessorTestSuiteState) TearDownTest() {
	s.mockQ.On("CompactOffers", s.ctx, s.sequence-100).Return(int64(0), nil).Once()
	s.Assert().NoError(s.processor.Commit(s.ctx))

	s.mockQ.AssertExpectations(s.T())
}

func (s *OffersProcessorTestSuiteState) TestCreateOffer() {
	offer := xdr.OfferEntry{
		SellerId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
		OfferId:  xdr.Int64(1),
		Price:    xdr.Price{1, 2},
	}
	lastModifiedLedgerSeq := xdr.Uint32(123)
	entry := xdr.LedgerEntry{
		Data: xdr.LedgerEntryData{
			Type:  xdr.LedgerEntryTypeOffer,
			Offer: &offer,
		},
		LastModifiedLedgerSeq: lastModifiedLedgerSeq,
	}

	s.mockQ.On("UpsertOffers", s.ctx, []history.Offer{
		{
			SellerID:           "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML",
			OfferID:            1,
			Pricen:             int32(1),
			Priced:             int32(2),
			Price:              float64(0.5),
			LastModifiedLedger: uint32(lastModifiedLedgerSeq),
		},
	}).Return(nil).Once()

	err := s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeOffer,
		Pre:  nil,
		Post: &entry,
	})
	s.Assert().NoError(err)
}

func TestOffersProcessorTestSuiteLedger(t *testing.T) {
	suite.Run(t, new(OffersProcessorTestSuiteLedger))
}

type OffersProcessorTestSuiteLedger struct {
	suite.Suite
	ctx       context.Context
	processor *OffersProcessor
	mockQ     *history.MockQOffers
	sequence  uint32
}

func (s *OffersProcessorTestSuiteLedger) SetupTest() {
	s.ctx = context.Background()
	s.mockQ = &history.MockQOffers{}

	s.sequence = 456
	s.processor = NewOffersProcessor(s.mockQ, s.sequence)
}

func (s *OffersProcessorTestSuiteLedger) TearDownTest() {
	s.mockQ.AssertExpectations(s.T())
}

func (s *OffersProcessorTestSuiteLedger) setupInsertOffer() {
	// should be ignored because it's not an offer type
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
		},
	})
	s.Assert().NoError(err)

	// add offer
	offer := xdr.OfferEntry{
		SellerId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
		OfferId:  xdr.Int64(2),
		Price:    xdr.Price{1, 2},
	}
	lastModifiedLedgerSeq := xdr.Uint32(1234)

	err = s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeOffer,
		Pre:  nil,
		Post: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: lastModifiedLedgerSeq,
			Data: xdr.LedgerEntryData{
				Type:  xdr.LedgerEntryTypeOffer,
				Offer: &offer,
			},
		},
	})
	s.Assert().NoError(err)

	updatedOffer := xdr.OfferEntry{
		SellerId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
		OfferId:  xdr.Int64(2),
		Price:    xdr.Price{1, 6},
	}

	updatedEntry := xdr.LedgerEntry{
		LastModifiedLedgerSeq: lastModifiedLedgerSeq,
		Data: xdr.LedgerEntryData{
			Type:  xdr.LedgerEntryTypeOffer,
			Offer: &updatedOffer,
		},
	}

	err = s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeOffer,
		Pre: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: lastModifiedLedgerSeq - 1,
			Data: xdr.LedgerEntryData{
				Type:  xdr.LedgerEntryTypeOffer,
				Offer: &offer,
			},
		},
		Post: &updatedEntry,
	})
	s.Assert().NoError(err)

	// We use LedgerEntryChangesCache so all changes are squashed
	s.mockQ.On("UpsertOffers", s.ctx, []history.Offer{
		{
			SellerID:           "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML",
			OfferID:            2,
			Pricen:             int32(1),
			Priced:             int32(6),
			Price:              float64(1) / float64(6),
			LastModifiedLedger: uint32(lastModifiedLedgerSeq),
		},
	}).Return(nil).Once()
}

func (s *OffersProcessorTestSuiteLedger) TestInsertOffer() {
	s.setupInsertOffer()
	s.mockQ.On("CompactOffers", s.ctx, s.sequence-100).Return(int64(0), nil).Once()
	s.Assert().NoError(s.processor.Commit(s.ctx))
}

func (s *OffersProcessorTestSuiteLedger) TestSkipCompactionIfSequenceEqualsWindow() {
	s.processor.sequence = compactionWindow
	s.setupInsertOffer()
	s.Assert().NoError(s.processor.Commit(s.ctx))
}

func (s *OffersProcessorTestSuiteLedger) TestSkipCompactionIfSequenceLessThanWindow() {
	s.processor.sequence = compactionWindow - 1
	s.setupInsertOffer()
	s.Assert().NoError(s.processor.Commit(s.ctx))
}

func (s *OffersProcessorTestSuiteLedger) TestCompactionError() {
	s.setupInsertOffer()
	s.mockQ.On("CompactOffers", s.ctx, s.sequence-100).
		Return(int64(0), errors.New("compaction error")).Once()
	s.Assert().EqualError(s.processor.Commit(s.ctx), "could not compact offers: compaction error")
}

func (s *OffersProcessorTestSuiteLedger) TestUpsertManyOffers() {
	lastModifiedLedgerSeq := xdr.Uint32(1234)

	offer := xdr.OfferEntry{
		SellerId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
		OfferId:  xdr.Int64(2),
		Price:    xdr.Price{1, 2},
	}
	updatedOffer := xdr.OfferEntry{
		SellerId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
		OfferId:  xdr.Int64(2),
		Price:    xdr.Price{1, 6},
	}

	anotherOffer := xdr.OfferEntry{
		SellerId: xdr.MustAddress("GDMUVYVYPYZYBDXNJWKFT3X2GCZCICTL3GSVP6AWBGB4ZZG7ZRDA746P"),
		OfferId:  xdr.Int64(3),
		Price:    xdr.Price{2, 3},
	}

	updatedEntry := xdr.LedgerEntry{
		LastModifiedLedgerSeq: lastModifiedLedgerSeq,
		Data: xdr.LedgerEntryData{
			Type:  xdr.LedgerEntryTypeOffer,
			Offer: &updatedOffer,
		},
	}

	err := s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeOffer,
		Pre: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: lastModifiedLedgerSeq - 1,
			Data: xdr.LedgerEntryData{
				Type:  xdr.LedgerEntryTypeOffer,
				Offer: &offer,
			},
		},
		Post: &updatedEntry,
	})
	s.Assert().NoError(err)

	err = s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeOffer,
		Pre:  nil,
		Post: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: lastModifiedLedgerSeq,
			Data: xdr.LedgerEntryData{
				Type:  xdr.LedgerEntryTypeOffer,
				Offer: &anotherOffer,
			},
		},
	})
	s.Assert().NoError(err)

	s.mockQ.On("UpsertOffers", s.ctx, mock.Anything).Run(func(args mock.Arguments) {
		// To fix order issue due to using ChangeCompactor
		offers := args.Get(1).([]history.Offer)
		s.Assert().ElementsMatch(
			offers,
			[]history.Offer{
				{
					SellerID:           "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML",
					OfferID:            2,
					Pricen:             int32(1),
					Priced:             int32(6),
					Price:              float64(1) / float64(6),
					LastModifiedLedger: uint32(lastModifiedLedgerSeq),
				},
				{
					SellerID:           "GDMUVYVYPYZYBDXNJWKFT3X2GCZCICTL3GSVP6AWBGB4ZZG7ZRDA746P",
					OfferID:            3,
					Pricen:             int32(2),
					Priced:             int32(3),
					Price:              float64(2) / float64(3),
					LastModifiedLedger: uint32(lastModifiedLedgerSeq),
				},
			},
		)
	}).Return(nil).Once()
	s.mockQ.On("CompactOffers", s.ctx, s.sequence-100).Return(int64(0), nil).Once()
	s.Assert().NoError(s.processor.Commit(s.ctx))
}

func (s *OffersProcessorTestSuiteLedger) TestRemoveOffer() {
	err := s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeOffer,
		Pre: &xdr.LedgerEntry{
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeOffer,
				Offer: &xdr.OfferEntry{
					SellerId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					OfferId:  xdr.Int64(3),
					Price:    xdr.Price{3, 1},
				},
			},
		},
		Post: nil,
	})
	s.Assert().NoError(err)

	s.mockQ.On("UpsertOffers", s.ctx, []history.Offer{
		{
			SellerID:           "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML",
			OfferID:            3,
			Pricen:             3,
			Priced:             1,
			Price:              3,
			Deleted:            true,
			LastModifiedLedger: 456,
		},
	}).Return(nil).Once()
	s.mockQ.On("CompactOffers", s.ctx, s.sequence-100).Return(int64(0), nil).Once()
	s.Assert().NoError(s.processor.Commit(s.ctx))
}

func (s *OffersProcessorTestSuiteLedger) TestProcessUpgradeChange() {
	// add offer
	offer := xdr.OfferEntry{
		SellerId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
		OfferId:  xdr.Int64(2),
		Price:    xdr.Price{1, 2},
	}
	lastModifiedLedgerSeq := xdr.Uint32(1234)

	err := s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeOffer,
		Post: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: lastModifiedLedgerSeq,
			Data: xdr.LedgerEntryData{
				Type:  xdr.LedgerEntryTypeOffer,
				Offer: &offer,
			},
		},
	})
	s.Assert().NoError(err)

	updatedOffer := xdr.OfferEntry{
		SellerId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
		OfferId:  xdr.Int64(2),
		Price:    xdr.Price{1, 6},
	}

	updatedEntry := xdr.LedgerEntry{
		LastModifiedLedgerSeq: lastModifiedLedgerSeq,
		Data: xdr.LedgerEntryData{
			Type:  xdr.LedgerEntryTypeOffer,
			Offer: &updatedOffer,
		},
	}

	err = s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeOffer,
		Pre: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: lastModifiedLedgerSeq,
			Data: xdr.LedgerEntryData{
				Type:  xdr.LedgerEntryTypeOffer,
				Offer: &offer,
			},
		},
		Post: &updatedEntry,
	})
	s.Assert().NoError(err)

	// We use LedgerEntryChangesCache so all changes are squashed
	s.mockQ.On("UpsertOffers", s.ctx, []history.Offer{
		{
			SellerID:           "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML",
			OfferID:            2,
			Pricen:             1,
			Priced:             6,
			Price:              float64(1) / float64(6),
			LastModifiedLedger: uint32(lastModifiedLedgerSeq),
		},
	}).Return(nil).Once()

	s.mockQ.On("CompactOffers", s.ctx, s.sequence-100).Return(int64(0), nil).Once()
	s.Assert().NoError(s.processor.Commit(s.ctx))
}

func (s *OffersProcessorTestSuiteLedger) TestRemoveMultipleOffers() {
	err := s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeOffer,
		Pre: &xdr.LedgerEntry{
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeOffer,
				Offer: &xdr.OfferEntry{
					SellerId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					OfferId:  xdr.Int64(3),
					Price:    xdr.Price{3, 1},
				},
			},
		},
		Post: nil,
	})
	s.Assert().NoError(err)

	err = s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeOffer,
		Pre: &xdr.LedgerEntry{
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeOffer,
				Offer: &xdr.OfferEntry{
					SellerId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					OfferId:  xdr.Int64(4),
					Price:    xdr.Price{3, 1},
				},
			},
		},
		Post: nil,
	})
	s.Assert().NoError(err)

	s.mockQ.On("CompactOffers", s.ctx, s.sequence-100).Return(int64(0), nil).Once()
	s.mockQ.On("UpsertOffers", s.ctx, mock.Anything).Run(func(args mock.Arguments) {
		// To fix order issue due to using ChangeCompactor
		offers := args.Get(1).([]history.Offer)
		s.Assert().ElementsMatch(
			offers,
			[]history.Offer{
				{
					SellerID:           "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML",
					OfferID:            3,
					Pricen:             3,
					Priced:             1,
					Price:              3,
					LastModifiedLedger: 456,
					Deleted:            true,
				},
				{
					SellerID:           "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML",
					OfferID:            4,
					Pricen:             3,
					Priced:             1,
					Price:              3,
					LastModifiedLedger: 456,
					Deleted:            true,
				},
			},
		)
	}).Return(nil).Once()

	err = s.processor.Commit(s.ctx)
	s.Assert().NoError(err)
}
