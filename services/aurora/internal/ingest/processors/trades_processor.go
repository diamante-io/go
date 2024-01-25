package processors

import (
	"context"
	"time"

	"github.com/guregu/null"

	"go/ingest"
	"go/services/aurora/internal/db2/history"
	"go/services/aurora/internal/toid"
	"go/support/errors"
	"go/xdr"
)

// TradeProcessor operations processor
type TradeProcessor struct {
	tradesQ history.QTrades
	ledger  xdr.LedgerHeaderHistoryEntry
	trades  []ingestTrade
	stats   TradeStats
}

func NewTradeProcessor(tradesQ history.QTrades, ledger xdr.LedgerHeaderHistoryEntry) *TradeProcessor {
	return &TradeProcessor{
		tradesQ: tradesQ,
		ledger:  ledger,
	}
}

type TradeStats struct {
	count int64
}

func (p *TradeProcessor) GetStats() TradeStats {
	return p.stats
}
func (stats *TradeStats) Map() map[string]interface{} {
	return map[string]interface{}{
		"stats_count": stats.count,
	}
}

// ProcessTransaction process the given transaction
func (p *TradeProcessor) ProcessTransaction(ctx context.Context, transaction ingest.LedgerTransaction) (err error) {
	if !transaction.Result.Successful() {
		return nil
	}

	trades, err := p.extractTrades(p.ledger, transaction)
	if err != nil {
		return err
	}

	p.trades = append(p.trades, trades...)
	p.stats.count += int64(len(trades))
	return nil
}

func (p *TradeProcessor) Commit(ctx context.Context) error {
	if len(p.trades) == 0 {
		return nil
	}

	batch := p.tradesQ.NewTradeBatchInsertBuilder(maxBatchSize)
	var poolIDs, accounts []string
	var assets []xdr.Asset
	for _, trade := range p.trades {
		if trade.buyerAccount != "" {
			accounts = append(accounts, trade.buyerAccount)
		}
		if trade.sellerAccount != "" {
			accounts = append(accounts, trade.sellerAccount)
		}
		if trade.liquidityPoolID != "" {
			poolIDs = append(poolIDs, trade.liquidityPoolID)
		}
		assets = append(assets, trade.boughtAsset)
		assets = append(assets, trade.soldAsset)
	}

	accountSet, err := p.tradesQ.CreateAccounts(ctx, accounts, maxBatchSize)
	if err != nil {
		return errors.Wrap(err, "Error creating account ids")
	}

	var assetMap map[string]history.Asset
	assetMap, err = p.tradesQ.CreateAssets(ctx, assets, maxBatchSize)
	if err != nil {
		return errors.Wrap(err, "Error creating asset ids")
	}

	var poolMap map[string]int64
	poolMap, err = p.tradesQ.CreateHistoryLiquidityPools(ctx, poolIDs, maxBatchSize)
	if err != nil {
		return errors.Wrap(err, "Error creating pool ids")
	}

	for _, trade := range p.trades {
		row := trade.row
		if id, ok := accountSet[trade.sellerAccount]; ok {
			row.BaseAccountID = null.IntFrom(id)
		} else if len(trade.sellerAccount) > 0 {
			return errors.Errorf("Could not find history account id for %s", trade.sellerAccount)
		}
		if id, ok := accountSet[trade.buyerAccount]; ok {
			row.CounterAccountID = null.IntFrom(id)
		} else if len(trade.buyerAccount) > 0 {
			return errors.Errorf("Could not find history account id for %s", trade.buyerAccount)
		}
		if id, ok := poolMap[trade.liquidityPoolID]; ok {
			row.BaseLiquidityPoolID = null.IntFrom(id)
		} else if len(trade.liquidityPoolID) > 0 {
			return errors.Errorf("Could not find history liquidity pool id for %s", trade.liquidityPoolID)
		}
		row.BaseAssetID = assetMap[trade.soldAsset.String()].ID
		row.CounterAssetID = assetMap[trade.boughtAsset.String()].ID

		if row.BaseAssetID > row.CounterAssetID {
			row.BaseIsSeller = false
			row.BaseAccountID, row.CounterAccountID = row.CounterAccountID, row.BaseAccountID
			row.BaseAssetID, row.CounterAssetID = row.CounterAssetID, row.BaseAssetID
			row.BaseAmount, row.CounterAmount = row.CounterAmount, row.BaseAmount
			row.BaseLiquidityPoolID, row.CounterLiquidityPoolID = row.CounterLiquidityPoolID, row.BaseLiquidityPoolID
			row.BaseOfferID, row.CounterOfferID = row.CounterOfferID, row.BaseOfferID
			row.PriceN, row.PriceD = row.PriceD, row.PriceN
		}

		if err = batch.Add(ctx, row); err != nil {
			return errors.Wrap(err, "Error adding trade to batch")
		}
	}

	if err = batch.Exec(ctx); err != nil {
		return errors.Wrap(err, "Error flushing operation batch")
	}
	return nil
}

func (p *TradeProcessor) findTradeSellPrice(
	transaction ingest.LedgerTransaction,
	opidx int,
	trade xdr.ClaimAtom,
) (int64, int64, error) {
	if trade.Type == xdr.ClaimAtomTypeClaimAtomTypeLiquidityPool {
		return int64(trade.AmountBought()), int64(trade.AmountSold()), nil
	}

	key := xdr.LedgerKey{}
	if err := key.SetOffer(trade.SellerId(), uint64(trade.OfferId())); err != nil {
		return 0, 0, errors.Wrap(err, "Could not create offer ledger key")
	}

	change, err := p.findOperationChange(transaction, opidx, key)
	if err != nil {
		return 0, 0, errors.Wrap(err, "could not find change for trade offer")
	}

	return int64(change.Pre.Data.Offer.Price.N), int64(change.Pre.Data.Offer.Price.D), nil
}

func (p *TradeProcessor) findOperationChange(tx ingest.LedgerTransaction, opidx int, key xdr.LedgerKey) (ingest.Change, error) {
	changes, err := tx.GetOperationChanges(uint32(opidx))
	if err != nil {
		return ingest.Change{}, errors.Wrap(err, "could not determine changes for operation")
	}

	var change ingest.Change
	for i := len(changes) - 1; i >= 0; i-- {
		change = changes[i]
		if change.Pre != nil && key.Equals(change.Pre.LedgerKey()) {
			return change, nil
		}
	}
	return ingest.Change{}, errors.Errorf("could not find operation for key %v", key)
}

func (p *TradeProcessor) findPoolFee(
	transaction ingest.LedgerTransaction,
	opidx int,
	poolID xdr.PoolId,
) (uint32, error) {
	key := xdr.LedgerKey{}
	if err := key.SetLiquidityPool(poolID); err != nil {
		return 0, errors.Wrap(err, "Could not create liquidity pool ledger key")

	}

	change, err := p.findOperationChange(transaction, opidx, key)
	if err != nil {
		return 0, errors.Wrap(err, "could not find change for liquidity pool")
	}

	return uint32(change.Pre.Data.MustLiquidityPool().Body.MustConstantProduct().Params.Fee), nil
}

type ingestTrade struct {
	row             history.InsertTrade
	sellerAccount   string
	liquidityPoolID string
	buyerAccount    string
	boughtAsset     xdr.Asset
	soldAsset       xdr.Asset
}

func (p *TradeProcessor) extractTrades(
	ledger xdr.LedgerHeaderHistoryEntry,
	transaction ingest.LedgerTransaction,
) ([]ingestTrade, error) {
	var result []ingestTrade

	closeTime := time.Unix(int64(ledger.Header.ScpValue.CloseTime), 0).UTC()

	opResults, ok := transaction.Result.OperationResults()
	if !ok {
		return result, errors.New("transaction has no operation results")
	}
	for opidx, op := range transaction.Envelope.Operations() {
		var trades []xdr.ClaimAtom
		var buyOfferExists bool
		var buyOffer xdr.OfferEntry

		switch op.Body.Type {
		case xdr.OperationTypePathPaymentStrictReceive:
			trades = opResults[opidx].MustTr().MustPathPaymentStrictReceiveResult().
				MustSuccess().
				Offers

		case xdr.OperationTypePathPaymentStrictSend:
			trades = opResults[opidx].MustTr().
				MustPathPaymentStrictSendResult().
				MustSuccess().
				Offers

		case xdr.OperationTypeManageBuyOffer:
			manageOfferResult := opResults[opidx].MustTr().MustManageBuyOfferResult().
				MustSuccess()
			trades = manageOfferResult.OffersClaimed
			buyOffer, buyOfferExists = manageOfferResult.Offer.GetOffer()

		case xdr.OperationTypeManageSellOffer:
			manageOfferResult := opResults[opidx].MustTr().MustManageSellOfferResult().
				MustSuccess()
			trades = manageOfferResult.OffersClaimed
			buyOffer, buyOfferExists = manageOfferResult.Offer.GetOffer()

		case xdr.OperationTypeCreatePassiveSellOffer:
			result := opResults[opidx].MustTr()

			// KNOWN ISSUE:  diamcircle-core creates results for CreatePassiveOffer operations
			// with the wrong result arm set.
			if result.Type == xdr.OperationTypeManageSellOffer {
				manageOfferResult := result.MustManageSellOfferResult().MustSuccess()
				trades = manageOfferResult.OffersClaimed
				buyOffer, buyOfferExists = manageOfferResult.Offer.GetOffer()
			} else {
				passiveOfferResult := result.MustCreatePassiveSellOfferResult().MustSuccess()
				trades = passiveOfferResult.OffersClaimed
				buyOffer, buyOfferExists = passiveOfferResult.Offer.GetOffer()
			}
		}

		opID := toid.New(
			int32(ledger.Header.LedgerSeq), int32(transaction.Index), int32(opidx+1),
		).ToInt64()
		for order, trade := range trades {
			// diamcircle-core will opportunistically garbage collect invalid offers (in the
			// event that a trader spends down their balance).  These garbage collected
			// offers get emitted in the result with the amount values set to zero.
			//
			// These zeroed ClaimOfferAtom values do not represent trades, and so we
			// skip them.
			if trade.AmountBought() == 0 && trade.AmountSold() == 0 {
				continue
			}

			sellPriceN, sellPriceD, err := p.findTradeSellPrice(transaction, opidx, trade)
			if err != nil {
				return result, err
			}

			row := history.InsertTrade{
				HistoryOperationID: opID,
				Order:              int32(order),
				LedgerCloseTime:    closeTime,
				CounterAmount:      int64(trade.AmountBought()),
				BaseAmount:         int64(trade.AmountSold()),
				BaseIsSeller:       true,
				PriceN:             sellPriceN,
				PriceD:             sellPriceD,
			}

			var sellerAccount, liquidityPoolID string
			if trade.Type == xdr.ClaimAtomTypeClaimAtomTypeLiquidityPool {
				id := trade.MustLiquidityPool().LiquidityPoolId
				liquidityPoolID = PoolIDToString(id)
				var fee uint32
				if fee, err = p.findPoolFee(transaction, opidx, id); err != nil {
					return nil, err
				}
				row.LiquidityPoolFee = null.IntFrom(int64(fee))
				row.Type = history.LiquidityPoolTradeType
			} else {
				row.BaseOfferID = null.IntFrom(int64(trade.OfferId()))
				sellerAccount = trade.SellerId().Address()
				row.Type = history.OrderbookTradeType
			}

			if buyOfferExists {
				row.CounterOfferID = null.IntFrom(EncodeOfferId(uint64(buyOffer.OfferId), CoreOfferIDType))
			} else {
				row.CounterOfferID = null.IntFrom(EncodeOfferId(uint64(opID), TOIDType))
			}

			var buyerAddress string
			if buyer := op.SourceAccount; buyer != nil {
				accid := buyer.ToAccountId()
				buyerAddress = accid.Address()
			} else {
				sa := transaction.Envelope.SourceAccount().ToAccountId()
				buyerAddress = sa.Address()
			}

			result = append(result, ingestTrade{
				row:             row,
				sellerAccount:   sellerAccount,
				liquidityPoolID: liquidityPoolID,
				buyerAccount:    buyerAddress,
				boughtAsset:     trade.AssetBought(),
				soldAsset:       trade.AssetSold(),
			})
		}
	}

	return result, nil
}
