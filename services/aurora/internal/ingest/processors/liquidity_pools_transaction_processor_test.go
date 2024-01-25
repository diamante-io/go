//lint:file-ignore U1001 Ignore all unused code, staticcheck doesn't understand testify/suite

package processors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"go/services/aurora/internal/db2/history"
	"go/services/aurora/internal/toid"
	"go/xdr"
)

type LiquidityPoolsTransactionProcessorTestSuiteLedger struct {
	suite.Suite
	ctx                               context.Context
	processor                         *LiquidityPoolsTransactionProcessor
	mockQ                             *history.MockQHistoryLiquidityPools
	mockTransactionBatchInsertBuilder *history.MockTransactionLiquidityPoolBatchInsertBuilder
	mockOperationBatchInsertBuilder   *history.MockOperationLiquidityPoolBatchInsertBuilder

	sequence uint32
}

func TestLiquidityPoolsTransactionProcessorTestSuiteLedger(t *testing.T) {
	suite.Run(t, new(LiquidityPoolsTransactionProcessorTestSuiteLedger))
}

func (s *LiquidityPoolsTransactionProcessorTestSuiteLedger) SetupTest() {
	s.ctx = context.Background()
	s.mockQ = &history.MockQHistoryLiquidityPools{}
	s.mockTransactionBatchInsertBuilder = &history.MockTransactionLiquidityPoolBatchInsertBuilder{}
	s.mockOperationBatchInsertBuilder = &history.MockOperationLiquidityPoolBatchInsertBuilder{}
	s.sequence = 20

	s.processor = NewLiquidityPoolsTransactionProcessor(
		s.mockQ,
		s.sequence,
	)
}

func (s *LiquidityPoolsTransactionProcessorTestSuiteLedger) TearDownTest() {
	s.mockQ.AssertExpectations(s.T())
	s.mockTransactionBatchInsertBuilder.AssertExpectations(s.T())
	s.mockOperationBatchInsertBuilder.AssertExpectations(s.T())
}

func (s *LiquidityPoolsTransactionProcessorTestSuiteLedger) mockTransactionBatchAdd(transactionID, internalID int64, err error) {
	s.mockTransactionBatchInsertBuilder.On("Add", s.ctx, transactionID, internalID).Return(err).Once()
}

func (s *LiquidityPoolsTransactionProcessorTestSuiteLedger) mockOperationBatchAdd(operationID, internalID int64, err error) {
	s.mockOperationBatchInsertBuilder.On("Add", s.ctx, operationID, internalID).Return(err).Once()
}

func (s *LiquidityPoolsTransactionProcessorTestSuiteLedger) TestEmptyLiquidityPools() {
	// What is this expecting? Doesn't seem to assert anything meaningful...
	err := s.processor.Commit(context.Background())
	s.Assert().NoError(err)
}

func (s *LiquidityPoolsTransactionProcessorTestSuiteLedger) testOperationInserts(poolID xdr.PoolId, body xdr.OperationBody, change xdr.LedgerEntryChange) {
	// Setup the transaction
	internalID := int64(1234)
	txn := createTransaction(true, 1)
	txn.Envelope.Operations()[0].Body = body
	txn.UnsafeMeta.V = 2
	txn.UnsafeMeta.V2.Operations = []xdr.OperationMeta{
		{Changes: xdr.LedgerEntryChanges{
			{
				Type: xdr.LedgerEntryChangeTypeLedgerEntryState,
				State: &xdr.LedgerEntry{
					Data: xdr.LedgerEntryData{
						Type: xdr.LedgerEntryTypeLiquidityPool,
						LiquidityPool: &xdr.LiquidityPoolEntry{
							LiquidityPoolId: poolID,
						},
					},
				},
			},
			change,
		}},
	}

	if body.Type == xdr.OperationTypeChangeTrust {
		// For insert test
		txn.Result.Result.Result.Results =
			&[]xdr.OperationResult{
				{
					Code: xdr.OperationResultCodeOpInner,
					Tr: &xdr.OperationResultTr{
						Type: xdr.OperationTypeChangeTrust,
						ChangeTrustResult: &xdr.ChangeTrustResult{
							Code: xdr.ChangeTrustResultCodeChangeTrustSuccess,
						},
					},
				},
			}
	}
	txnID := toid.New(int32(s.sequence), int32(txn.Index), 0).ToInt64()
	opID := (&transactionOperationWrapper{
		index:          uint32(0),
		transaction:    txn,
		operation:      txn.Envelope.Operations()[0],
		ledgerSequence: s.sequence,
	}).ID()

	hexID := PoolIDToString(poolID)

	// Setup a q
	s.mockQ.On("CreateHistoryLiquidityPools", s.ctx, mock.AnythingOfType("[]string"), maxBatchSize).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).([]string)
			s.Assert().ElementsMatch(
				[]string{hexID},
				arg,
			)
		}).Return(map[string]int64{
		hexID: internalID,
	}, nil).Once()

	// Prepare to process transactions successfully
	s.mockQ.On("NewTransactionLiquidityPoolBatchInsertBuilder", maxBatchSize).
		Return(s.mockTransactionBatchInsertBuilder).Once()
	s.mockTransactionBatchAdd(txnID, internalID, nil)
	s.mockTransactionBatchInsertBuilder.On("Exec", s.ctx).Return(nil).Once()

	// Prepare to process operations successfully
	s.mockQ.On("NewOperationLiquidityPoolBatchInsertBuilder", maxBatchSize).
		Return(s.mockOperationBatchInsertBuilder).Once()
	s.mockOperationBatchAdd(opID, internalID, nil)
	s.mockOperationBatchInsertBuilder.On("Exec", s.ctx).Return(nil).Once()

	// Process the transaction
	err := s.processor.ProcessTransaction(s.ctx, txn)
	s.Assert().NoError(err)
	err = s.processor.Commit(s.ctx)
	s.Assert().NoError(err)
}

func (s *LiquidityPoolsTransactionProcessorTestSuiteLedger) TestIngestLiquidityPoolsRemoval() {
	poolID := xdr.PoolId{0xca, 0xfe, 0xba, 0xbe}
	s.testOperationInserts(poolID,
		xdr.OperationBody{
			Type: xdr.OperationTypeLiquidityPoolDeposit,
			ChangeTrustOp: &xdr.ChangeTrustOp{
				Line: xdr.ChangeTrustAsset{
					Type: xdr.AssetTypeAssetTypePoolShare,
					LiquidityPool: &xdr.LiquidityPoolParameters{
						Type: xdr.LiquidityPoolTypeLiquidityPoolConstantProduct,
						ConstantProduct: &xdr.LiquidityPoolConstantProductParameters{
							AssetA: xdr.MustNewCreditAsset("EUR", "GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"),
							AssetB: xdr.MustNewNativeAsset(),
							Fee:    30,
						},
					},
				},
				Limit: 0,
			},
		},
		xdr.LedgerEntryChange{
			Type: xdr.LedgerEntryChangeTypeLedgerEntryRemoved,
			Removed: &xdr.LedgerKey{
				Type: xdr.LedgerEntryTypeLiquidityPool,
				LiquidityPool: &xdr.LedgerKeyLiquidityPool{
					LiquidityPoolId: poolID,
				},
			},
		},
	)
}

func (s *LiquidityPoolsTransactionProcessorTestSuiteLedger) TestIngestLiquidityPoolsUpdate() {
	poolID := xdr.PoolId{0xca, 0xfe, 0xba, 0xbe}
	s.testOperationInserts(poolID,
		xdr.OperationBody{
			Type: xdr.OperationTypeLiquidityPoolDeposit,
			ChangeTrustOp: &xdr.ChangeTrustOp{
				Line: xdr.ChangeTrustAsset{
					Type: xdr.AssetTypeAssetTypePoolShare,
					LiquidityPool: &xdr.LiquidityPoolParameters{
						Type: xdr.LiquidityPoolTypeLiquidityPoolConstantProduct,
						ConstantProduct: &xdr.LiquidityPoolConstantProductParameters{
							AssetA: xdr.MustNewCreditAsset("EUR", "GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"),
							AssetB: xdr.MustNewNativeAsset(),
							Fee:    30,
						},
					},
				},
				Limit: 10,
			},
		},
		xdr.LedgerEntryChange{
			Type: xdr.LedgerEntryChangeTypeLedgerEntryUpdated,
			Updated: &xdr.LedgerEntry{
				Data: xdr.LedgerEntryData{
					Type: xdr.LedgerEntryTypeLiquidityPool,
					LiquidityPool: &xdr.LiquidityPoolEntry{
						LiquidityPoolId: poolID,
						Body:            xdr.LiquidityPoolEntryBody{},
					},
				},
			},
		},
	)
}

func (s *LiquidityPoolsTransactionProcessorTestSuiteLedger) TestIngestLiquidityPoolsCreate() {
	poolID := xdr.PoolId{0xca, 0xfe, 0xba, 0xbe}
	s.testOperationInserts(poolID,
		xdr.OperationBody{
			Type: xdr.OperationTypeLiquidityPoolDeposit,
			ChangeTrustOp: &xdr.ChangeTrustOp{
				Line: xdr.ChangeTrustAsset{
					Type: xdr.AssetTypeAssetTypePoolShare,
					LiquidityPool: &xdr.LiquidityPoolParameters{
						Type: xdr.LiquidityPoolTypeLiquidityPoolConstantProduct,
						ConstantProduct: &xdr.LiquidityPoolConstantProductParameters{
							AssetA: xdr.MustNewCreditAsset("EUR", "GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"),
							AssetB: xdr.MustNewNativeAsset(),
							Fee:    30,
						},
					},
				},
				Limit: 10,
			},
		},
		xdr.LedgerEntryChange{
			Type: xdr.LedgerEntryChangeTypeLedgerEntryCreated,
			Created: &xdr.LedgerEntry{
				Data: xdr.LedgerEntryData{
					Type: xdr.LedgerEntryTypeLiquidityPool,
					LiquidityPool: &xdr.LiquidityPoolEntry{
						LiquidityPoolId: poolID,
						Body:            xdr.LiquidityPoolEntryBody{},
					},
				},
			},
		},
	)
}
