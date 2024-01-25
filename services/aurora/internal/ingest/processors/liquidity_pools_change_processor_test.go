//lint:file-ignore U1001 Ignore all unused code, staticcheck doesn't understand testify/suite

package processors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/diamcircle/go/ingest"
	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/xdr"
)

func TestLiquidityPoolsChangeProcessorTestSuiteState(t *testing.T) {
	suite.Run(t, new(LiquidityPoolsChangeProcessorTestSuiteState))
}

type LiquidityPoolsChangeProcessorTestSuiteState struct {
	suite.Suite
	ctx       context.Context
	processor *LiquidityPoolsChangeProcessor
	mockQ     *history.MockQLiquidityPools
	sequence  uint32
}

func (s *LiquidityPoolsChangeProcessorTestSuiteState) SetupTest() {
	s.ctx = context.Background()
	s.mockQ = &history.MockQLiquidityPools{}

	s.sequence = 456
	s.processor = NewLiquidityPoolsChangeProcessor(s.mockQ, s.sequence)
}

func (s *LiquidityPoolsChangeProcessorTestSuiteState) TearDownTest() {
	s.Assert().NoError(s.processor.Commit(s.ctx))
	s.mockQ.AssertExpectations(s.T())
}

func (s *LiquidityPoolsChangeProcessorTestSuiteState) TestNoEntries() {
	s.mockQ.On("CompactLiquidityPools", s.ctx, s.sequence-100).Return(int64(0), nil).Once()
}

func (s *LiquidityPoolsChangeProcessorTestSuiteState) TestNoEntriesWithSequenceLessThanWindow() {
	s.sequence = 50
	s.processor.sequence = s.sequence
	// Nothing processed, assertions in TearDownTest.
}

func (s *LiquidityPoolsChangeProcessorTestSuiteState) TestCreatesLiquidityPools() {
	lastModifiedLedgerSeq := xdr.Uint32(123)
	lpoolEntry := xdr.LiquidityPoolEntry{
		LiquidityPoolId: xdr.PoolId{0xca, 0xfe, 0xba, 0xbe, 0xde, 0xad, 0xbe, 0xef},
		Body: xdr.LiquidityPoolEntryBody{
			Type: xdr.LiquidityPoolTypeLiquidityPoolConstantProduct,
			ConstantProduct: &xdr.LiquidityPoolEntryConstantProduct{
				Params: xdr.LiquidityPoolConstantProductParameters{
					AssetA: xdr.MustNewCreditAsset("USD", "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					AssetB: xdr.MustNewNativeAsset(),
					Fee:    34,
				},
				ReserveA:                 450,
				ReserveB:                 500,
				TotalPoolShares:          412241,
				PoolSharesTrustLineCount: 52115,
			},
		},
	}
	lp := history.LiquidityPool{
		PoolID:         "cafebabedeadbeef000000000000000000000000000000000000000000000000",
		Type:           xdr.LiquidityPoolTypeLiquidityPoolConstantProduct,
		Fee:            34,
		TrustlineCount: 52115,
		ShareCount:     412241,
		AssetReserves: []history.LiquidityPoolAssetReserve{
			{
				xdr.MustNewCreditAsset("USD", "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
				450,
			},
			{
				xdr.MustNewNativeAsset(),
				500,
			},
		},
		LastModifiedLedger: 123,
	}
	s.mockQ.On("UpsertLiquidityPools", s.ctx, []history.LiquidityPool{lp}).Return(nil).Once()

	s.mockQ.On("CompactLiquidityPools", s.ctx, s.sequence-100).Return(int64(0), nil).Once()

	err := s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeLiquidityPool,
		Pre:  nil,
		Post: &xdr.LedgerEntry{
			Data: xdr.LedgerEntryData{
				Type:          xdr.LedgerEntryTypeLiquidityPool,
				LiquidityPool: &lpoolEntry,
			},
			LastModifiedLedgerSeq: lastModifiedLedgerSeq,
		},
	})
	s.Assert().NoError(err)
}

func TestLiquidityPoolsChangeProcessorTestSuiteLedger(t *testing.T) {
	suite.Run(t, new(LiquidityPoolsChangeProcessorTestSuiteLedger))
}

type LiquidityPoolsChangeProcessorTestSuiteLedger struct {
	suite.Suite
	ctx       context.Context
	processor *LiquidityPoolsChangeProcessor
	mockQ     *history.MockQLiquidityPools
	sequence  uint32
}

func (s *LiquidityPoolsChangeProcessorTestSuiteLedger) SetupTest() {
	s.ctx = context.Background()
	s.mockQ = &history.MockQLiquidityPools{}

	s.sequence = 456
	s.processor = NewLiquidityPoolsChangeProcessor(s.mockQ, s.sequence)
}

func (s *LiquidityPoolsChangeProcessorTestSuiteLedger) TearDownTest() {
	s.Assert().NoError(s.processor.Commit(s.ctx))
	s.mockQ.AssertExpectations(s.T())
}

func (s *LiquidityPoolsChangeProcessorTestSuiteLedger) TestNoTransactions() {
	s.mockQ.On("CompactLiquidityPools", s.ctx, s.sequence-100).Return(int64(0), nil).Once()
}

func (s *LiquidityPoolsChangeProcessorTestSuiteLedger) TestNoEntriesWithSequenceLessThanWindow() {
	s.sequence = 50
	s.processor.sequence = s.sequence
	// Nothing processed, assertions in TearDownTest.
}

func (s *LiquidityPoolsChangeProcessorTestSuiteLedger) TestNewLiquidityPool() {
	lastModifiedLedgerSeq := xdr.Uint32(123)
	lpEntry := xdr.LiquidityPoolEntry{
		LiquidityPoolId: xdr.PoolId{0xca, 0xfe, 0xba, 0xbe, 0xde, 0xad, 0xbe, 0xef},
		Body: xdr.LiquidityPoolEntryBody{
			Type: xdr.LiquidityPoolTypeLiquidityPoolConstantProduct,
			ConstantProduct: &xdr.LiquidityPoolEntryConstantProduct{
				Params: xdr.LiquidityPoolConstantProductParameters{
					AssetA: xdr.MustNewCreditAsset("USD", "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					AssetB: xdr.MustNewNativeAsset(),
					Fee:    34,
				},
				ReserveA:                 450,
				ReserveB:                 500,
				TotalPoolShares:          412241,
				PoolSharesTrustLineCount: 52115,
			},
		},
	}
	pre := xdr.LedgerEntry{
		Data: xdr.LedgerEntryData{
			Type:          xdr.LedgerEntryTypeLiquidityPool,
			LiquidityPool: &lpEntry,
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
		Type: xdr.LedgerEntryTypeLiquidityPool,
		Pre:  nil,
		Post: &pre,
	})
	s.Assert().NoError(err)

	// add sponsor
	post := xdr.LedgerEntry{
		Data: xdr.LedgerEntryData{
			Type:          xdr.LedgerEntryTypeLiquidityPool,
			LiquidityPool: &lpEntry,
		},
		LastModifiedLedgerSeq: lastModifiedLedgerSeq,
		Ext: xdr.LedgerEntryExt{
			V: 1,
			V1: &xdr.LedgerEntryExtensionV1{
				SponsoringId: xdr.MustAddressPtr("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
			},
		},
	}

	pre.LastModifiedLedgerSeq = pre.LastModifiedLedgerSeq - 1
	err = s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeLiquidityPool,
		Pre:  &pre,
		Post: &post,
	})
	s.Assert().NoError(err)

	postLP := history.LiquidityPool{
		PoolID:         "cafebabedeadbeef000000000000000000000000000000000000000000000000",
		Type:           xdr.LiquidityPoolTypeLiquidityPoolConstantProduct,
		Fee:            34,
		TrustlineCount: 52115,
		ShareCount:     412241,
		AssetReserves: []history.LiquidityPoolAssetReserve{
			{
				xdr.MustNewCreditAsset("USD", "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
				450,
			},
			{
				xdr.MustNewNativeAsset(),
				500,
			},
		},
		LastModifiedLedger: 123,
	}
	s.mockQ.On("UpsertLiquidityPools", s.ctx, []history.LiquidityPool{postLP}).Return(nil).Once()
	s.mockQ.On("CompactLiquidityPools", s.ctx, s.sequence-100).Return(int64(0), nil).Once()
}

func (s *LiquidityPoolsChangeProcessorTestSuiteLedger) TestUpdateLiquidityPool() {
	lpEntry := xdr.LiquidityPoolEntry{
		LiquidityPoolId: xdr.PoolId{0xca, 0xfe, 0xba, 0xbe, 0xde, 0xad, 0xbe, 0xef},
		Body: xdr.LiquidityPoolEntryBody{
			Type: xdr.LiquidityPoolTypeLiquidityPoolConstantProduct,
			ConstantProduct: &xdr.LiquidityPoolEntryConstantProduct{
				Params: xdr.LiquidityPoolConstantProductParameters{
					AssetA: xdr.MustNewCreditAsset("USD", "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					AssetB: xdr.MustNewNativeAsset(),
					Fee:    34,
				},
				ReserveA:                 450,
				ReserveB:                 500,
				TotalPoolShares:          412241,
				PoolSharesTrustLineCount: 52115,
			},
		},
	}
	lastModifiedLedgerSeq := xdr.Uint32(123)

	pre := xdr.LedgerEntry{
		Data: xdr.LedgerEntryData{
			Type:          xdr.LedgerEntryTypeLiquidityPool,
			LiquidityPool: &lpEntry,
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
	post := xdr.LedgerEntry{
		Data: xdr.LedgerEntryData{
			Type:          xdr.LedgerEntryTypeLiquidityPool,
			LiquidityPool: &lpEntry,
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
		Type: xdr.LedgerEntryTypeLiquidityPool,
		Pre:  &pre,
		Post: &post,
	})
	s.Assert().NoError(err)

	postLP := history.LiquidityPool{
		PoolID:         "cafebabedeadbeef000000000000000000000000000000000000000000000000",
		Type:           xdr.LiquidityPoolTypeLiquidityPoolConstantProduct,
		Fee:            34,
		TrustlineCount: 52115,
		ShareCount:     412241,
		AssetReserves: []history.LiquidityPoolAssetReserve{
			{
				xdr.MustNewCreditAsset("USD", "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
				450,
			},
			{
				xdr.MustNewNativeAsset(),
				500,
			},
		},
		LastModifiedLedger: 123,
	}

	s.mockQ.On("UpsertLiquidityPools", s.ctx, []history.LiquidityPool{postLP}).Return(nil).Once()
	s.mockQ.On("CompactLiquidityPools", s.ctx, s.sequence-100).Return(int64(0), nil).Once()
}

func (s *LiquidityPoolsChangeProcessorTestSuiteLedger) TestRemoveLiquidityPool() {
	lpEntry := xdr.LiquidityPoolEntry{
		LiquidityPoolId: xdr.PoolId{0xca, 0xfe, 0xba, 0xbe, 0xde, 0xad, 0xbe, 0xef},
		Body: xdr.LiquidityPoolEntryBody{
			Type: xdr.LiquidityPoolTypeLiquidityPoolConstantProduct,
			ConstantProduct: &xdr.LiquidityPoolEntryConstantProduct{
				Params: xdr.LiquidityPoolConstantProductParameters{
					AssetA: xdr.MustNewCreditAsset("USD", "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					AssetB: xdr.MustNewNativeAsset(),
					Fee:    34,
				},
				ReserveA:                 450,
				ReserveB:                 123,
				TotalPoolShares:          412241,
				PoolSharesTrustLineCount: 52115,
			},
		},
	}
	lastModifiedLedgerSeq := xdr.Uint32(123)
	pre := xdr.LedgerEntry{
		Data: xdr.LedgerEntryData{
			Type:          xdr.LedgerEntryTypeLiquidityPool,
			LiquidityPool: &lpEntry,
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
		Type: xdr.LedgerEntryTypeLiquidityPool,
		Pre:  &pre,
		Post: nil,
	})
	s.Assert().NoError(err)

	deleted := s.processor.ledgerEntryToRow(&pre)
	deleted.Deleted = true
	deleted.LastModifiedLedger = s.processor.sequence
	s.mockQ.On("UpsertLiquidityPools", s.ctx, []history.LiquidityPool{deleted}).Return(nil).Once()
	s.mockQ.On("CompactLiquidityPools", s.ctx, s.sequence-100).Return(int64(0), nil).Once()
}
