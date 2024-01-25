//lint:file-ignore U1001 Ignore all unused code, staticcheck doesn't understand testify/suite

package processors

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/guregu/null"

	"github.com/diamcircle/go/ingest"
	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/services/aurora/internal/toid"
	"github.com/diamcircle/go/xdr"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TradeProcessorTestSuiteLedger struct {
	suite.Suite
	processor              *TradeProcessor
	mockQ                  *history.MockQTrades
	mockBatchInsertBuilder *history.MockTradeBatchInsertBuilder

	unmuxedSourceAccount       xdr.AccountId
	unmuxedOpSourceAccount     xdr.AccountId
	sourceAccount              xdr.MuxedAccount
	opSourceAccount            xdr.MuxedAccount
	strictReceiveTrade         xdr.ClaimAtom
	strictReceiveTradeLP       xdr.ClaimAtom
	strictSendTrade            xdr.ClaimAtom
	strictSendTradeLP          xdr.ClaimAtom
	buyOfferTrade              xdr.ClaimAtom
	sellOfferTrade             xdr.ClaimAtom
	passiveSellOfferTrade      xdr.ClaimAtom
	otherPassiveSellOfferTrade xdr.ClaimAtom
	allTrades                  []xdr.ClaimAtom
	sellPrices                 []xdr.Price

	assets []xdr.Asset

	lpToID             map[xdr.PoolId]int64
	unmuxedAccountToID map[string]int64
	assetToID          map[string]history.Asset

	txs []ingest.LedgerTransaction
}

func TestTradeProcessorTestSuiteLedger(t *testing.T) {
	suite.Run(t, new(TradeProcessorTestSuiteLedger))
}

func (s *TradeProcessorTestSuiteLedger) SetupTest() {
	s.mockQ = &history.MockQTrades{}
	s.mockBatchInsertBuilder = &history.MockTradeBatchInsertBuilder{}

	s.unmuxedSourceAccount = xdr.MustAddress("GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY")
	s.sourceAccount = xdr.MuxedAccount{
		Type: xdr.CryptoKeyTypeKeyTypeMuxedEd25519,
		Med25519: &xdr.MuxedAccountMed25519{
			Id:      0xdeadbeef,
			Ed25519: *s.unmuxedSourceAccount.Ed25519,
		},
	}
	s.unmuxedOpSourceAccount = xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML")
	s.opSourceAccount = xdr.MuxedAccount{
		Type: xdr.CryptoKeyTypeKeyTypeMuxedEd25519,
		Med25519: &xdr.MuxedAccountMed25519{
			Id:      0xcafebabe,
			Ed25519: *s.unmuxedOpSourceAccount.Ed25519,
		},
	}
	s.strictReceiveTrade = xdr.ClaimAtom{
		Type: xdr.ClaimAtomTypeClaimAtomTypeOrderBook,
		OrderBook: &xdr.ClaimOfferAtom{
			SellerId:     xdr.MustAddress("GA2YS6YBWIBUMUJCNYROC5TXYTTUA4TCZF7A4MJ2O4TTGT3LFNWIOMY4"),
			OfferId:      11,
			AssetSold:    xdr.MustNewNativeAsset(),
			AmountSold:   111,
			AmountBought: 211,
			AssetBought:  xdr.MustNewCreditAsset("HUF", s.unmuxedSourceAccount.Address()),
		},
	}
	s.strictReceiveTradeLP = xdr.ClaimAtom{
		Type: xdr.ClaimAtomTypeClaimAtomTypeLiquidityPool,
		LiquidityPool: &xdr.ClaimLiquidityAtom{
			LiquidityPoolId: xdr.PoolId{1, 2, 3},
			AssetSold:       xdr.MustNewCreditAsset("MAD", s.unmuxedSourceAccount.Address()),
			AmountSold:      20,
			AssetBought:     xdr.MustNewCreditAsset("GRE", s.unmuxedSourceAccount.Address()),
			AmountBought:    300,
		},
	}
	s.strictSendTrade = xdr.ClaimAtom{
		Type: xdr.ClaimAtomTypeClaimAtomTypeOrderBook,
		OrderBook: &xdr.ClaimOfferAtom{
			SellerId:     xdr.MustAddress("GALOBQKDZUSAEUDE7F4OYUIQTUZBL62G6TRCXU2ED6SA7TL72MBUQSYJ"),
			OfferId:      12,
			AssetSold:    xdr.MustNewCreditAsset("USD", s.unmuxedSourceAccount.Address()),
			AmountSold:   112,
			AmountBought: 212,
			AssetBought:  xdr.MustNewCreditAsset("RUB", s.unmuxedSourceAccount.Address()),
		},
	}
	s.strictSendTradeLP = xdr.ClaimAtom{
		Type: xdr.ClaimAtomTypeClaimAtomTypeLiquidityPool,
		LiquidityPool: &xdr.ClaimLiquidityAtom{
			LiquidityPoolId: xdr.PoolId{4, 5, 6},
			AssetSold:       xdr.MustNewCreditAsset("WER", s.unmuxedSourceAccount.Address()),
			AmountSold:      67,
			AssetBought:     xdr.MustNewCreditAsset("NIJ", s.unmuxedSourceAccount.Address()),
			AmountBought:    98,
		},
	}
	s.buyOfferTrade = xdr.ClaimAtom{
		Type: xdr.ClaimAtomTypeClaimAtomTypeOrderBook,
		OrderBook: &xdr.ClaimOfferAtom{
			SellerId:     xdr.MustAddress("GCWRLPH5X5A3GABFDLDILZ4RLY6O76AYOIIR5H2PAI6TNZZZNLZWBXSH"),
			OfferId:      13,
			AssetSold:    xdr.MustNewCreditAsset("EUR", s.unmuxedSourceAccount.Address()),
			AmountSold:   113,
			AmountBought: 213,
			AssetBought:  xdr.MustNewCreditAsset("NOK", s.unmuxedSourceAccount.Address()),
		},
	}
	s.sellOfferTrade = xdr.ClaimAtom{
		Type: xdr.ClaimAtomTypeClaimAtomTypeOrderBook,
		OrderBook: &xdr.ClaimOfferAtom{
			SellerId:     xdr.MustAddress("GAVOLNFXVVUJOELN4T3YVSH2FFA3VSP2XN4NJRYF2ZWVCHS77C5KXLHZ"),
			OfferId:      14,
			AssetSold:    xdr.MustNewCreditAsset("PLN", s.unmuxedSourceAccount.Address()),
			AmountSold:   114,
			AmountBought: 214,
			AssetBought:  xdr.MustNewCreditAsset("UAH", s.unmuxedSourceAccount.Address()),
		},
	}
	s.passiveSellOfferTrade = xdr.ClaimAtom{
		Type: xdr.ClaimAtomTypeClaimAtomTypeOrderBook,
		OrderBook: &xdr.ClaimOfferAtom{
			SellerId:     xdr.MustAddress("GDQWI6FKB72DPOJE4CGYCFQZKRPQQIOYXRMZ5KEVGXMG6UUTGJMBCASH"),
			OfferId:      15,
			AssetSold:    xdr.MustNewCreditAsset("SEK", s.unmuxedSourceAccount.Address()),
			AmountSold:   115,
			AmountBought: 215,
			AssetBought:  xdr.MustNewCreditAsset("GBP", s.unmuxedSourceAccount.Address()),
		},
	}
	s.otherPassiveSellOfferTrade = xdr.ClaimAtom{
		Type: xdr.ClaimAtomTypeClaimAtomTypeOrderBook,
		OrderBook: &xdr.ClaimOfferAtom{
			SellerId:     xdr.MustAddress("GCPZFOJON3PSSYUBNT7MCGEDSGP47UTSJSB4XGCVEWEJO4XQ6U4XN3N2"),
			OfferId:      16,
			AssetSold:    xdr.MustNewCreditAsset("CHF", s.unmuxedSourceAccount.Address()),
			AmountSold:   116,
			AmountBought: 216,
			AssetBought:  xdr.MustNewCreditAsset("JPY", s.unmuxedSourceAccount.Address()),
		},
	}

	s.unmuxedAccountToID = map[string]int64{
		s.unmuxedSourceAccount.Address():   1000,
		s.unmuxedOpSourceAccount.Address(): 1001,
	}
	s.assetToID = map[string]history.Asset{}
	s.allTrades = []xdr.ClaimAtom{
		s.strictReceiveTrade,
		s.strictSendTrade,
		s.buyOfferTrade,
		s.sellOfferTrade,
		s.passiveSellOfferTrade,
		s.otherPassiveSellOfferTrade,
		s.strictReceiveTradeLP,
		s.strictSendTradeLP,
	}

	s.assets = []xdr.Asset{}
	s.sellPrices = []xdr.Price{}
	s.lpToID = map[xdr.PoolId]int64{}
	for i, trade := range s.allTrades {
		if trade.Type == xdr.ClaimAtomTypeClaimAtomTypeOrderBook {
			s.unmuxedAccountToID[trade.SellerId().Address()] = int64(1002 + i)
			n := xdr.Int32(i + 1)
			s.sellPrices = append(s.sellPrices, xdr.Price{N: n, D: 100})
		} else {
			s.lpToID[trade.MustLiquidityPool().LiquidityPoolId] = int64(3000 + i)
			s.sellPrices = append(s.sellPrices, xdr.Price{N: xdr.Int32(trade.AmountBought()), D: xdr.Int32(trade.AmountSold())})
		}
		if i%2 == 0 {
			s.assetToID[trade.AssetSold().String()] = history.Asset{ID: int64(10000 + i)}
			s.assetToID[trade.AssetBought().String()] = history.Asset{ID: int64(100 + i)}
		} else {
			s.assetToID[trade.AssetSold().String()] = history.Asset{ID: int64(100 + i)}
			s.assetToID[trade.AssetBought().String()] = history.Asset{ID: int64(10000 + i)}
		}
		s.assets = append(s.assets, trade.AssetSold(), trade.AssetBought())
	}

	s.processor = NewTradeProcessor(
		s.mockQ,
		xdr.LedgerHeaderHistoryEntry{
			Header: xdr.LedgerHeader{
				LedgerSeq: 100,
			},
		},
	)
}

func (s *TradeProcessorTestSuiteLedger) TearDownTest() {
	s.mockQ.AssertExpectations(s.T())
	s.mockBatchInsertBuilder.AssertExpectations(s.T())
}

func (s *TradeProcessorTestSuiteLedger) TestIgnoreFailedTransactions() {
	ctx := context.Background()
	err := s.processor.ProcessTransaction(ctx, createTransaction(false, 1))
	s.Assert().NoError(err)

	err = s.processor.Commit(ctx)
	s.Assert().NoError(err)
}

func (s *TradeProcessorTestSuiteLedger) mockReadTradeTransactions(
	ledger xdr.LedgerHeaderHistoryEntry,
) []history.InsertTrade {
	closeTime := time.Unix(int64(ledger.Header.ScpValue.CloseTime), 0).UTC()
	inserts := []history.InsertTrade{
		{
			HistoryOperationID: toid.New(int32(ledger.Header.LedgerSeq), 1, 2).ToInt64(),
			Order:              1,
			LedgerCloseTime:    closeTime,
			BaseAmount:         int64(s.strictReceiveTrade.AmountBought()),
			BaseAccountID:      null.IntFrom(s.unmuxedAccountToID[s.unmuxedOpSourceAccount.Address()]),
			BaseAssetID:        s.assetToID[s.strictReceiveTrade.AssetBought().String()].ID,
			BaseOfferID:        null.IntFrom(EncodeOfferId(uint64(toid.New(int32(ledger.Header.LedgerSeq), 1, 2).ToInt64()), TOIDType)),
			CounterAmount:      int64(s.strictReceiveTrade.AmountSold()),
			CounterAccountID:   null.IntFrom(s.unmuxedAccountToID[s.strictReceiveTrade.SellerId().Address()]),
			CounterAssetID:     s.assetToID[s.strictReceiveTrade.AssetSold().String()].ID,
			CounterOfferID:     null.IntFrom(int64(s.strictReceiveTrade.OfferId())),
			BaseIsSeller:       false,
			PriceN:             int64(s.sellPrices[0].D),
			PriceD:             int64(s.sellPrices[0].N),
			Type:               history.OrderbookTradeType,
		},
		{
			HistoryOperationID: toid.New(int32(ledger.Header.LedgerSeq), 1, 3).ToInt64(),
			Order:              0,
			LedgerCloseTime:    closeTime,
			CounterAmount:      int64(s.strictSendTrade.AmountBought()),
			CounterAccountID:   null.IntFrom(s.unmuxedAccountToID[s.unmuxedOpSourceAccount.Address()]),
			CounterAssetID:     s.assetToID[s.strictSendTrade.AssetBought().String()].ID,
			CounterOfferID:     null.IntFrom(EncodeOfferId(uint64(toid.New(int32(ledger.Header.LedgerSeq), 1, 3).ToInt64()), TOIDType)),
			BaseAmount:         int64(s.strictSendTrade.AmountSold()),
			BaseAccountID:      null.IntFrom(s.unmuxedAccountToID[s.strictSendTrade.SellerId().Address()]),
			BaseAssetID:        s.assetToID[s.strictSendTrade.AssetSold().String()].ID,
			BaseIsSeller:       true,
			BaseOfferID:        null.IntFrom(int64(s.strictSendTrade.OfferId())),
			PriceN:             int64(s.sellPrices[1].N),
			PriceD:             int64(s.sellPrices[1].D),
			Type:               history.OrderbookTradeType,
		},
		{
			HistoryOperationID: toid.New(int32(ledger.Header.LedgerSeq), 1, 4).ToInt64(),
			Order:              1,
			LedgerCloseTime:    closeTime,
			BaseOfferID:        null.IntFrom(879136),
			BaseAmount:         int64(s.buyOfferTrade.AmountBought()),
			BaseAccountID:      null.IntFrom(s.unmuxedAccountToID[s.unmuxedOpSourceAccount.Address()]),
			BaseAssetID:        s.assetToID[s.buyOfferTrade.AssetBought().String()].ID,
			CounterAmount:      int64(s.buyOfferTrade.AmountSold()),
			CounterAccountID:   null.IntFrom(s.unmuxedAccountToID[s.buyOfferTrade.SellerId().Address()]),
			CounterAssetID:     s.assetToID[s.buyOfferTrade.AssetSold().String()].ID,
			BaseIsSeller:       false,
			CounterOfferID:     null.IntFrom(int64(s.buyOfferTrade.OfferId())),
			PriceN:             int64(s.sellPrices[2].D),
			PriceD:             int64(s.sellPrices[2].N),
			Type:               history.OrderbookTradeType,
		},
		{
			HistoryOperationID: toid.New(int32(ledger.Header.LedgerSeq), 1, 5).ToInt64(),
			Order:              2,
			LedgerCloseTime:    closeTime,
			CounterAmount:      int64(s.sellOfferTrade.AmountBought()),
			CounterAssetID:     s.assetToID[s.sellOfferTrade.AssetBought().String()].ID,
			CounterAccountID:   null.IntFrom(s.unmuxedAccountToID[s.unmuxedOpSourceAccount.Address()]),
			CounterOfferID:     null.IntFrom(EncodeOfferId(uint64(toid.New(int32(ledger.Header.LedgerSeq), 1, 5).ToInt64()), TOIDType)),
			BaseAmount:         int64(s.sellOfferTrade.AmountSold()),
			BaseAccountID:      null.IntFrom(s.unmuxedAccountToID[s.sellOfferTrade.SellerId().Address()]),
			BaseAssetID:        s.assetToID[s.sellOfferTrade.AssetSold().String()].ID,
			BaseIsSeller:       true,
			BaseOfferID:        null.IntFrom(int64(s.sellOfferTrade.OfferId())),
			PriceN:             int64(s.sellPrices[3].N),
			PriceD:             int64(s.sellPrices[3].D),
			Type:               history.OrderbookTradeType,
		},
		{
			HistoryOperationID: toid.New(int32(ledger.Header.LedgerSeq), 1, 6).ToInt64(),
			Order:              0,
			LedgerCloseTime:    closeTime,
			BaseAmount:         int64(s.passiveSellOfferTrade.AmountBought()),
			BaseAssetID:        s.assetToID[s.passiveSellOfferTrade.AssetBought().String()].ID,
			BaseAccountID:      null.IntFrom(s.unmuxedAccountToID[s.unmuxedSourceAccount.Address()]),
			BaseOfferID:        null.IntFrom(EncodeOfferId(uint64(toid.New(int32(ledger.Header.LedgerSeq), 1, 6).ToInt64()), TOIDType)),
			CounterAmount:      int64(s.passiveSellOfferTrade.AmountSold()),
			CounterAccountID:   null.IntFrom(s.unmuxedAccountToID[s.passiveSellOfferTrade.SellerId().Address()]),
			CounterAssetID:     s.assetToID[s.passiveSellOfferTrade.AssetSold().String()].ID,
			BaseIsSeller:       false,
			CounterOfferID:     null.IntFrom(int64(s.passiveSellOfferTrade.OfferId())),
			PriceN:             int64(s.sellPrices[4].D),
			PriceD:             int64(s.sellPrices[4].N),
			Type:               history.OrderbookTradeType,
		},

		{
			HistoryOperationID: toid.New(int32(ledger.Header.LedgerSeq), 1, 7).ToInt64(),
			Order:              0,
			LedgerCloseTime:    closeTime,

			CounterAmount:    int64(s.otherPassiveSellOfferTrade.AmountBought()),
			CounterAssetID:   s.assetToID[s.otherPassiveSellOfferTrade.AssetBought().String()].ID,
			CounterAccountID: null.IntFrom(s.unmuxedAccountToID[s.unmuxedOpSourceAccount.Address()]),
			CounterOfferID:   null.IntFrom(EncodeOfferId(uint64(toid.New(int32(ledger.Header.LedgerSeq), 1, 7).ToInt64()), TOIDType)),
			BaseAmount:       int64(s.otherPassiveSellOfferTrade.AmountSold()),
			BaseAccountID:    null.IntFrom(s.unmuxedAccountToID[s.otherPassiveSellOfferTrade.SellerId().Address()]),
			BaseAssetID:      s.assetToID[s.otherPassiveSellOfferTrade.AssetSold().String()].ID,
			BaseIsSeller:     true,
			BaseOfferID:      null.IntFrom(int64(s.otherPassiveSellOfferTrade.OfferId())),
			PriceN:           int64(s.sellPrices[5].N),
			PriceD:           int64(s.sellPrices[5].D),
			Type:             history.OrderbookTradeType,
		},
		{
			HistoryOperationID:     toid.New(int32(ledger.Header.LedgerSeq), 1, 8).ToInt64(),
			Order:                  1,
			LedgerCloseTime:        closeTime,
			BaseAmount:             int64(s.strictReceiveTradeLP.AmountBought()),
			BaseAssetID:            s.assetToID[s.strictReceiveTradeLP.AssetBought().String()].ID,
			BaseAccountID:          null.IntFrom(s.unmuxedAccountToID[s.unmuxedOpSourceAccount.Address()]),
			BaseOfferID:            null.IntFrom(EncodeOfferId(uint64(toid.New(int32(ledger.Header.LedgerSeq), 1, 8).ToInt64()), TOIDType)),
			CounterAmount:          int64(s.strictReceiveTradeLP.AmountSold()),
			CounterLiquidityPoolID: null.IntFrom(s.lpToID[s.strictReceiveTradeLP.MustLiquidityPool().LiquidityPoolId]),
			CounterAssetID:         s.assetToID[s.strictReceiveTradeLP.AssetSold().String()].ID,
			BaseIsSeller:           false,
			LiquidityPoolFee:       null.IntFrom(int64(xdr.LiquidityPoolFeeV18)),
			PriceN:                 int64(s.sellPrices[6].D),
			PriceD:                 int64(s.sellPrices[6].N),
			Type:                   history.LiquidityPoolTradeType,
		},
		{
			HistoryOperationID:  toid.New(int32(ledger.Header.LedgerSeq), 1, 9).ToInt64(),
			Order:               0,
			LedgerCloseTime:     closeTime,
			CounterAmount:       int64(s.strictSendTradeLP.AmountBought()),
			CounterAssetID:      s.assetToID[s.strictSendTradeLP.AssetBought().String()].ID,
			CounterAccountID:    null.IntFrom(s.unmuxedAccountToID[s.unmuxedOpSourceAccount.Address()]),
			CounterOfferID:      null.IntFrom(EncodeOfferId(uint64(toid.New(int32(ledger.Header.LedgerSeq), 1, 9).ToInt64()), TOIDType)),
			BaseAmount:          int64(s.strictSendTradeLP.AmountSold()),
			BaseLiquidityPoolID: null.IntFrom(s.lpToID[s.strictSendTradeLP.MustLiquidityPool().LiquidityPoolId]),
			BaseAssetID:         s.assetToID[s.strictSendTradeLP.AssetSold().String()].ID,
			BaseIsSeller:        true,
			LiquidityPoolFee:    null.IntFrom(int64(xdr.LiquidityPoolFeeV18)),
			PriceN:              int64(s.sellPrices[7].N),
			PriceD:              int64(s.sellPrices[7].D),
			Type:                history.LiquidityPoolTradeType,
		},
	}

	emptyTrade := xdr.ClaimAtom{
		Type: xdr.ClaimAtomTypeClaimAtomTypeOrderBook,
		OrderBook: &xdr.ClaimOfferAtom{
			SellerId:     s.sourceAccount.ToAccountId(),
			OfferId:      123,
			AssetSold:    xdr.MustNewNativeAsset(),
			AmountSold:   0,
			AssetBought:  xdr.MustNewCreditAsset("EUR", s.unmuxedSourceAccount.Address()),
			AmountBought: 0,
		},
	}

	operationResults := []xdr.OperationResult{
		{
			Tr: &xdr.OperationResultTr{
				Type: xdr.OperationTypeBumpSequence,
				BumpSeqResult: &xdr.BumpSequenceResult{
					Code: xdr.BumpSequenceResultCodeBumpSequenceSuccess,
				},
			},
		},
		{
			Tr: &xdr.OperationResultTr{
				Type: xdr.OperationTypePathPaymentStrictReceive,
				PathPaymentStrictReceiveResult: &xdr.PathPaymentStrictReceiveResult{
					Code: xdr.PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveSuccess,
					Success: &xdr.PathPaymentStrictReceiveResultSuccess{
						Offers: []xdr.ClaimAtom{
							emptyTrade,
							s.strictReceiveTrade,
						},
					},
				},
			},
		},
		{
			Tr: &xdr.OperationResultTr{
				Type: xdr.OperationTypePathPaymentStrictSend,
				PathPaymentStrictSendResult: &xdr.PathPaymentStrictSendResult{
					Code: xdr.PathPaymentStrictSendResultCodePathPaymentStrictSendSuccess,
					Success: &xdr.PathPaymentStrictSendResultSuccess{
						Offers: []xdr.ClaimAtom{
							s.strictSendTrade,
							emptyTrade,
						},
					},
				},
			},
		},
		{
			Tr: &xdr.OperationResultTr{
				Type: xdr.OperationTypeManageBuyOffer,
				ManageBuyOfferResult: &xdr.ManageBuyOfferResult{
					Code: xdr.ManageBuyOfferResultCodeManageBuyOfferSuccess,
					Success: &xdr.ManageOfferSuccessResult{
						OffersClaimed: []xdr.ClaimAtom{
							emptyTrade,
							s.buyOfferTrade,
							emptyTrade,
						},
						Offer: xdr.ManageOfferSuccessResultOffer{
							Offer: &xdr.OfferEntry{
								OfferId: 879136,
							},
						},
					},
				},
			},
		},
		{
			Tr: &xdr.OperationResultTr{
				Type: xdr.OperationTypeManageSellOffer,
				ManageSellOfferResult: &xdr.ManageSellOfferResult{
					Code: xdr.ManageSellOfferResultCodeManageSellOfferSuccess,
					Success: &xdr.ManageOfferSuccessResult{
						OffersClaimed: []xdr.ClaimAtom{
							emptyTrade,
							emptyTrade,
							s.sellOfferTrade,
						},
						Offer: xdr.ManageOfferSuccessResultOffer{
							Effect: xdr.ManageOfferEffectManageOfferDeleted,
						},
					},
				},
			},
		},
		{
			Tr: &xdr.OperationResultTr{
				Type: xdr.OperationTypeManageSellOffer,
				ManageSellOfferResult: &xdr.ManageSellOfferResult{
					Code: xdr.ManageSellOfferResultCodeManageSellOfferSuccess,
					Success: &xdr.ManageOfferSuccessResult{
						OffersClaimed: []xdr.ClaimAtom{
							s.passiveSellOfferTrade,
							emptyTrade,
							emptyTrade,
						},
						Offer: xdr.ManageOfferSuccessResultOffer{
							Effect: xdr.ManageOfferEffectManageOfferDeleted,
						},
					},
				},
			},
		},
		{
			Tr: &xdr.OperationResultTr{
				Type: xdr.OperationTypeCreatePassiveSellOffer,
				CreatePassiveSellOfferResult: &xdr.ManageSellOfferResult{
					Code: xdr.ManageSellOfferResultCodeManageSellOfferSuccess,
					Success: &xdr.ManageOfferSuccessResult{
						OffersClaimed: []xdr.ClaimAtom{
							s.otherPassiveSellOfferTrade,
						},
						Offer: xdr.ManageOfferSuccessResultOffer{
							Effect: xdr.ManageOfferEffectManageOfferDeleted,
						},
					},
				},
			},
		},
		{
			Tr: &xdr.OperationResultTr{
				Type: xdr.OperationTypePathPaymentStrictReceive,
				PathPaymentStrictReceiveResult: &xdr.PathPaymentStrictReceiveResult{
					Code: xdr.PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveSuccess,
					Success: &xdr.PathPaymentStrictReceiveResultSuccess{
						Offers: []xdr.ClaimAtom{
							emptyTrade,
							s.strictReceiveTradeLP,
						},
					},
				},
			},
		},
		{
			Tr: &xdr.OperationResultTr{
				Type: xdr.OperationTypePathPaymentStrictSend,
				PathPaymentStrictSendResult: &xdr.PathPaymentStrictSendResult{
					Code: xdr.PathPaymentStrictSendResultCodePathPaymentStrictSendSuccess,
					Success: &xdr.PathPaymentStrictSendResultSuccess{
						Offers: []xdr.ClaimAtom{
							s.strictSendTradeLP,
						},
					},
				},
			},
		},
	}

	operations := []xdr.Operation{
		{
			Body: xdr.OperationBody{
				Type:           xdr.OperationTypeBumpSequence,
				BumpSequenceOp: &xdr.BumpSequenceOp{BumpTo: 30000},
			},
		},
		{
			Body: xdr.OperationBody{
				Type:                       xdr.OperationTypePathPaymentStrictReceive,
				PathPaymentStrictReceiveOp: &xdr.PathPaymentStrictReceiveOp{},
			},
			SourceAccount: &s.opSourceAccount,
		},
		{
			Body: xdr.OperationBody{
				Type:                    xdr.OperationTypePathPaymentStrictSend,
				PathPaymentStrictSendOp: &xdr.PathPaymentStrictSendOp{},
			},
			SourceAccount: &s.opSourceAccount,
		},
		{
			Body: xdr.OperationBody{
				Type:             xdr.OperationTypeManageBuyOffer,
				ManageBuyOfferOp: &xdr.ManageBuyOfferOp{},
			},
			SourceAccount: &s.opSourceAccount,
		},
		{
			Body: xdr.OperationBody{
				Type:              xdr.OperationTypeManageSellOffer,
				ManageSellOfferOp: &xdr.ManageSellOfferOp{},
			},
			SourceAccount: &s.opSourceAccount,
		},
		{
			Body: xdr.OperationBody{
				Type:                     xdr.OperationTypeCreatePassiveSellOffer,
				CreatePassiveSellOfferOp: &xdr.CreatePassiveSellOfferOp{},
			},
		},
		{
			Body: xdr.OperationBody{
				Type:                     xdr.OperationTypeCreatePassiveSellOffer,
				CreatePassiveSellOfferOp: &xdr.CreatePassiveSellOfferOp{},
			},
			SourceAccount: &s.opSourceAccount,
		},
		{
			Body: xdr.OperationBody{
				Type:                       xdr.OperationTypePathPaymentStrictReceive,
				PathPaymentStrictReceiveOp: &xdr.PathPaymentStrictReceiveOp{},
			},
			SourceAccount: &s.opSourceAccount,
		},
		{
			Body: xdr.OperationBody{
				Type:                    xdr.OperationTypePathPaymentStrictSend,
				PathPaymentStrictSendOp: &xdr.PathPaymentStrictSendOp{},
			},
			SourceAccount: &s.opSourceAccount,
		},
	}

	tx := ingest.LedgerTransaction{
		Result: xdr.TransactionResultPair{
			Result: xdr.TransactionResult{
				Result: xdr.TransactionResultResult{
					Code:    xdr.TransactionResultCodeTxSuccess,
					Results: &operationResults,
				},
			},
		},
		Envelope: xdr.TransactionEnvelope{
			Type: xdr.EnvelopeTypeEnvelopeTypeTx,
			V1: &xdr.TransactionV1Envelope{
				Tx: xdr.Transaction{
					SourceAccount: s.sourceAccount,
					Operations:    operations,
				},
			},
		},
		Index:      1,
		FeeChanges: []xdr.LedgerEntryChange{},
		UnsafeMeta: xdr.TransactionMeta{
			V: 2,
			V2: &xdr.TransactionMetaV2{
				Operations: []xdr.OperationMeta{
					{
						Changes: xdr.LedgerEntryChanges{},
					},
				},
			},
		},
	}

	for i, trade := range s.allTrades {
		var changes xdr.LedgerEntryChanges
		if trade.Type == xdr.ClaimAtomTypeClaimAtomTypeOrderBook {
			changes = xdr.LedgerEntryChanges{
				xdr.LedgerEntryChange{
					Type: xdr.LedgerEntryChangeTypeLedgerEntryState,
					State: &xdr.LedgerEntry{
						Data: xdr.LedgerEntryData{
							Type: xdr.LedgerEntryTypeOffer,
							Offer: &xdr.OfferEntry{
								Price:    s.sellPrices[i],
								SellerId: trade.SellerId(),
								OfferId:  trade.OfferId(),
							},
						},
					},
				},
				xdr.LedgerEntryChange{
					Type: xdr.LedgerEntryChangeTypeLedgerEntryRemoved,
					Removed: &xdr.LedgerKey{
						Type: xdr.LedgerEntryTypeOffer,
						Offer: &xdr.LedgerKeyOffer{
							SellerId: trade.SellerId(),
							OfferId:  trade.OfferId(),
						},
					},
				},
			}
		} else {
			changes = xdr.LedgerEntryChanges{
				xdr.LedgerEntryChange{
					Type: xdr.LedgerEntryChangeTypeLedgerEntryState,
					State: &xdr.LedgerEntry{
						Data: xdr.LedgerEntryData{
							Type: xdr.LedgerEntryTypeLiquidityPool,
							LiquidityPool: &xdr.LiquidityPoolEntry{
								LiquidityPoolId: trade.MustLiquidityPool().LiquidityPoolId,
								Body: xdr.LiquidityPoolEntryBody{
									Type: xdr.LiquidityPoolTypeLiquidityPoolConstantProduct,
									ConstantProduct: &xdr.LiquidityPoolEntryConstantProduct{
										Params: xdr.LiquidityPoolConstantProductParameters{
											AssetA: trade.AssetSold(),
											AssetB: trade.AssetSold(),
											Fee:    xdr.LiquidityPoolFeeV18,
										},
										ReserveA:                 100,
										ReserveB:                 200,
										TotalPoolShares:          40,
										PoolSharesTrustLineCount: 50,
									},
								},
							},
						},
					},
				},
				xdr.LedgerEntryChange{
					Type: xdr.LedgerEntryChangeTypeLedgerEntryUpdated,
					Updated: &xdr.LedgerEntry{
						Data: xdr.LedgerEntryData{
							Type: xdr.LedgerEntryTypeLiquidityPool,
							LiquidityPool: &xdr.LiquidityPoolEntry{
								LiquidityPoolId: trade.MustLiquidityPool().LiquidityPoolId,
								Body: xdr.LiquidityPoolEntryBody{
									Type: xdr.LiquidityPoolTypeLiquidityPoolConstantProduct,
									ConstantProduct: &xdr.LiquidityPoolEntryConstantProduct{
										Params: xdr.LiquidityPoolConstantProductParameters{
											AssetA: trade.AssetSold(),
											AssetB: trade.AssetSold(),
											Fee:    xdr.LiquidityPoolFeeV18,
										},
										ReserveA:                 100,
										ReserveB:                 200,
										TotalPoolShares:          40,
										PoolSharesTrustLineCount: 50,
									},
								},
							},
						},
					},
				},
			}
		}
		tx.UnsafeMeta.V2.Operations = append(tx.UnsafeMeta.V2.Operations, xdr.OperationMeta{
			Changes: changes,
		})
	}

	s.txs = []ingest.LedgerTransaction{
		tx,
	}

	s.mockQ.On("NewTradeBatchInsertBuilder", maxBatchSize).
		Return(s.mockBatchInsertBuilder).Once()

	return inserts
}

func mapKeysToList(set map[string]int64) []string {
	keys := make([]string, 0, len(set))
	for key := range set {
		keys = append(keys, key)
	}
	return keys
}

func uniq(list []string) []string {
	var deduped []string
	set := map[string]bool{}
	for _, s := range list {
		if set[s] {
			continue
		}
		deduped = append(deduped, s)
		set[s] = true
	}
	return deduped
}

func (s *TradeProcessorTestSuiteLedger) TestIngestTradesSucceeds() {
	ctx := context.Background()
	inserts := s.mockReadTradeTransactions(s.processor.ledger)

	s.mockCreateAccounts(ctx)

	s.mockCreateAssets(ctx)

	s.mockCreateHistoryLiquidityPools(ctx)

	for _, insert := range inserts {
		s.mockBatchInsertBuilder.On("Add", ctx, []history.InsertTrade{
			insert,
		}).Return(nil).Once()
	}

	s.mockBatchInsertBuilder.On("Exec", ctx).Return(nil).Once()

	for _, tx := range s.txs {
		err := s.processor.ProcessTransaction(ctx, tx)
		s.Assert().NoError(err)
	}

	err := s.processor.Commit(ctx)
	s.Assert().NoError(err)
}

func (s *TradeProcessorTestSuiteLedger) mockCreateHistoryLiquidityPools(ctx context.Context) {
	lpIDs, lpStrToID := s.extractLpIDs()
	s.mockQ.On("CreateHistoryLiquidityPools", ctx, mock.AnythingOfType("[]string"), maxBatchSize).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).([]string)
			s.Assert().ElementsMatch(
				lpIDs,
				arg,
			)
		}).Return(lpStrToID, nil).Once()
}

func (s *TradeProcessorTestSuiteLedger) extractLpIDs() ([]string, map[string]int64) {
	var lpIDs []string
	lpStrToID := map[string]int64{}
	for lpID, id := range s.lpToID {
		lpIDStr := PoolIDToString(lpID)
		lpIDs = append(lpIDs, lpIDStr)
		lpStrToID[lpIDStr] = id
	}
	return lpIDs, lpStrToID
}

func (s *TradeProcessorTestSuiteLedger) TestCreateAccountsError() {
	ctx := context.Background()
	s.mockReadTradeTransactions(s.processor.ledger)

	s.mockQ.On("CreateAccounts", ctx, mock.AnythingOfType("[]string"), maxBatchSize).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).([]string)
			s.Assert().ElementsMatch(
				mapKeysToList(s.unmuxedAccountToID),
				uniq(arg),
			)
		}).Return(map[string]int64{}, fmt.Errorf("create accounts error")).Once()

	for _, tx := range s.txs {
		err := s.processor.ProcessTransaction(ctx, tx)
		s.Assert().NoError(err)
	}

	err := s.processor.Commit(ctx)

	s.Assert().EqualError(err, "Error creating account ids: create accounts error")
}

func (s *TradeProcessorTestSuiteLedger) TestCreateAssetsError() {
	ctx := context.Background()
	s.mockReadTradeTransactions(s.processor.ledger)

	s.mockCreateAccounts(ctx)

	s.mockQ.On("CreateAssets", ctx, mock.AnythingOfType("[]xdr.Asset"), maxBatchSize).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).([]xdr.Asset)
			s.Assert().ElementsMatch(
				s.assets,
				arg,
			)
		}).Return(s.assetToID, fmt.Errorf("create assets error")).Once()

	for _, tx := range s.txs {
		err := s.processor.ProcessTransaction(ctx, tx)
		s.Assert().NoError(err)
	}

	err := s.processor.Commit(ctx)
	s.Assert().EqualError(err, "Error creating asset ids: create assets error")
}

func (s *TradeProcessorTestSuiteLedger) TestCreateHistoryLiquidityPoolsError() {
	ctx := context.Background()
	s.mockReadTradeTransactions(s.processor.ledger)

	s.mockCreateAccounts(ctx)

	s.mockCreateAssets(ctx)

	lpIDs, lpStrToID := s.extractLpIDs()
	s.mockQ.On("CreateHistoryLiquidityPools", ctx, mock.AnythingOfType("[]string"), maxBatchSize).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).([]string)
			s.Assert().ElementsMatch(
				lpIDs,
				arg,
			)
		}).Return(lpStrToID, fmt.Errorf("create liqudity pool id error")).Once()

	for _, tx := range s.txs {
		err := s.processor.ProcessTransaction(ctx, tx)
		s.Assert().NoError(err)
	}

	err := s.processor.Commit(ctx)
	s.Assert().EqualError(err, "Error creating pool ids: create liqudity pool id error")
}

func (s *TradeProcessorTestSuiteLedger) mockCreateAssets(ctx context.Context) {
	s.mockQ.On("CreateAssets", ctx, mock.AnythingOfType("[]xdr.Asset"), maxBatchSize).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).([]xdr.Asset)
			s.Assert().ElementsMatch(
				s.assets,
				arg,
			)
		}).Return(s.assetToID, nil).Once()
}

func (s *TradeProcessorTestSuiteLedger) mockCreateAccounts(ctx context.Context) {
	s.mockQ.On("CreateAccounts", ctx, mock.AnythingOfType("[]string"), maxBatchSize).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).([]string)
			s.Assert().ElementsMatch(
				mapKeysToList(s.unmuxedAccountToID),
				uniq(arg),
			)
		}).Return(s.unmuxedAccountToID, nil).Once()
}

func (s *TradeProcessorTestSuiteLedger) TestBatchAddError() {
	ctx := context.Background()
	s.mockReadTradeTransactions(s.processor.ledger)

	s.mockCreateAccounts(ctx)

	s.mockCreateAssets(ctx)

	s.mockCreateHistoryLiquidityPools(ctx)

	s.mockBatchInsertBuilder.On("Add", ctx, mock.AnythingOfType("[]history.InsertTrade")).
		Return(fmt.Errorf("batch add error")).Once()

	for _, tx := range s.txs {
		err := s.processor.ProcessTransaction(ctx, tx)
		s.Assert().NoError(err)
	}

	err := s.processor.Commit(ctx)
	s.Assert().EqualError(err, "Error adding trade to batch: batch add error")
}

func (s *TradeProcessorTestSuiteLedger) TestBatchExecError() {
	ctx := context.Background()
	insert := s.mockReadTradeTransactions(s.processor.ledger)

	s.mockCreateAccounts(ctx)

	s.mockCreateAssets(ctx)

	s.mockCreateHistoryLiquidityPools(ctx)

	s.mockBatchInsertBuilder.On("Add", ctx, mock.AnythingOfType("[]history.InsertTrade")).
		Return(nil).Times(len(insert))
	s.mockBatchInsertBuilder.On("Exec", ctx).Return(fmt.Errorf("exec error")).Once()
	for _, tx := range s.txs {
		err := s.processor.ProcessTransaction(ctx, tx)
		s.Assert().NoError(err)
	}

	err := s.processor.Commit(ctx)
	s.Assert().EqualError(err, "Error flushing operation batch: exec error")
}

func (s *TradeProcessorTestSuiteLedger) TestIgnoreCheckIfSmallLedger() {
	ctx := context.Background()
	insert := s.mockReadTradeTransactions(s.processor.ledger)

	s.mockCreateAccounts(ctx)

	s.mockCreateAssets(ctx)

	s.mockCreateHistoryLiquidityPools(ctx)
	s.mockBatchInsertBuilder.On("Add", ctx, mock.AnythingOfType("[]history.InsertTrade")).
		Return(nil).Times(len(insert))
	s.mockBatchInsertBuilder.On("Exec", ctx).Return(nil).Once()

	for _, tx := range s.txs {
		err := s.processor.ProcessTransaction(ctx, tx)
		s.Assert().NoError(err)
	}

	err := s.processor.Commit(ctx)
	s.Assert().NoError(err)
}

func TestTradeProcessor_ProcessTransaction_MuxedAccount(t *testing.T) {
	unmuxed := xdr.MustAddress("GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2")
	muxed := xdr.MuxedAccount{
		Type: xdr.CryptoKeyTypeKeyTypeMuxedEd25519,
		Med25519: &xdr.MuxedAccountMed25519{
			Id:      0xdeadbeefdeadbeef,
			Ed25519: *unmuxed.Ed25519,
		},
	}
	tx := createTransaction(true, 1)
	tx.Index = 1
	tx.Envelope.Operations()[0].Body = xdr.OperationBody{
		Type: xdr.OperationTypePayment,
		PaymentOp: &xdr.PaymentOp{
			Destination: muxed,
			Asset:       xdr.Asset{Type: xdr.AssetTypeAssetTypeNative},
			Amount:      100,
		},
	}
}
