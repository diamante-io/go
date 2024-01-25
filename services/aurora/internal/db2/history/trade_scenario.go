package history

import (
	"encoding/hex"
	"strings"
	"time"

	"github.com/guregu/null"

	"go/services/aurora/internal/test"
	"go/services/aurora/internal/toid"
	"go/xdr"
)

func createInsertTrades(
	accountIDs, assetIDs, poolIDs []int64, ledger int32,
) []InsertTrade {
	first := InsertTrade{
		HistoryOperationID: toid.New(ledger, 1, 1).ToInt64(),
		Order:              1,
		LedgerCloseTime:    time.Unix(10000000, 0).UTC(),
		CounterOfferID:     null.IntFrom(32145),
		BaseAccountID:      null.IntFrom(accountIDs[0]),
		CounterAccountID:   null.IntFrom(accountIDs[1]),
		BaseAssetID:        assetIDs[0],
		CounterAssetID:     assetIDs[1],
		BaseOfferID:        null.IntFrom(214515),
		BaseIsSeller:       true,
		BaseAmount:         7986,
		CounterAmount:      896,
		PriceN:             1,
		PriceD:             3,
		Type:               OrderbookTradeType,
	}

	second := first
	second.CounterOfferID = null.Int{}
	second.Order = 2

	third := InsertTrade{
		HistoryOperationID: toid.New(ledger, 2, 1).ToInt64(),
		Order:              1,
		LedgerCloseTime:    time.Unix(10000001, 0).UTC(),
		CounterOfferID:     null.IntFrom(2),
		BaseAccountID:      null.IntFrom(accountIDs[0]),
		CounterAccountID:   null.IntFrom(accountIDs[1]),
		BaseAssetID:        assetIDs[1],
		CounterAssetID:     assetIDs[2],
		BaseOfferID:        null.IntFrom(7),
		BaseIsSeller:       false,
		BaseAmount:         123,
		CounterAmount:      6,
		PriceN:             1156,
		PriceD:             3,
		Type:               OrderbookTradeType,
	}

	fourth := InsertTrade{
		HistoryOperationID:  toid.New(ledger, 2, 2).ToInt64(),
		Order:               3,
		LedgerCloseTime:     time.Unix(10000001, 0).UTC(),
		CounterAssetID:      assetIDs[4],
		CounterAmount:       675,
		CounterAccountID:    null.IntFrom(accountIDs[0]),
		LiquidityPoolFee:    null.IntFrom(xdr.LiquidityPoolFeeV18),
		BaseAssetID:         assetIDs[3],
		BaseAmount:          981,
		BaseLiquidityPoolID: null.IntFrom(poolIDs[0]),
		BaseIsSeller:        true,
		PriceN:              675,
		PriceD:              981,
		Type:                LiquidityPoolTradeType,
	}

	fifth := InsertTrade{
		HistoryOperationID:  toid.New(ledger, 3, 1).ToInt64(),
		Order:               1,
		LedgerCloseTime:     time.Unix(10002000, 0).UTC(),
		CounterAssetID:      assetIDs[1],
		CounterAmount:       300,
		CounterAccountID:    null.IntFrom(accountIDs[1]),
		LiquidityPoolFee:    null.IntFrom(xdr.LiquidityPoolFeeV18),
		BaseAssetID:         assetIDs[0],
		BaseAmount:          200,
		BaseLiquidityPoolID: null.IntFrom(poolIDs[1]),
		BaseIsSeller:        true,
		PriceN:              43,
		PriceD:              56,
		Type:                LiquidityPoolTradeType,
	}

	return []InsertTrade{
		first,
		second,
		third,
		fourth,
		fifth,
	}
}

func createHistoryIDs(
	tt *test.T, q *Q, accounts []string, assets []xdr.Asset, pools []string,
) ([]int64, []int64, []int64) {
	addressToAccounts, err := q.CreateAccounts(tt.Ctx, accounts, len(accounts))
	tt.Assert.NoError(err)

	accountIDs := []int64{}
	for _, account := range accounts {
		accountIDs = append(accountIDs, addressToAccounts[account])
	}

	assetMap, err := q.CreateAssets(tt.Ctx, assets, len(assets))
	tt.Assert.NoError(err)

	assetIDs := []int64{}
	for _, asset := range assets {
		assetIDs = append(assetIDs, assetMap[asset.String()].ID)
	}

	poolsMap, err := q.CreateHistoryLiquidityPools(tt.Ctx, pools, len(pools))
	tt.Assert.NoError(err)
	poolIDs := []int64{}
	for _, pool := range pools {
		poolIDs = append(poolIDs, poolsMap[pool])
	}

	return accountIDs, assetIDs, poolIDs
}

func buildIDtoAccountMapping(addresses []string, ids []int64) map[int64]xdr.AccountId {
	idToAccount := map[int64]xdr.AccountId{}
	for i, id := range ids {
		account := xdr.MustAddress(addresses[i])
		idToAccount[id] = account
	}

	return idToAccount
}

func buildIDtoAssetMapping(assets []xdr.Asset, ids []int64) map[int64]xdr.Asset {
	idToAsset := map[int64]xdr.Asset{}
	for i, id := range ids {
		idToAsset[id] = assets[i]
	}

	return idToAsset
}

// TradeFixtures contains the data inserted into the database
// when running TradeScenario
type TradeFixtures struct {
	Addresses       []string
	Assets          []xdr.Asset
	Trades          []Trade
	LiquidityPools  []string
	TradesByAccount map[string][]Trade
	TradesByAsset   map[string][]Trade
	TradesByPool    map[string][]Trade
	TradesByOffer   map[int64][]Trade
}

// TradesByAssetPair returns the trades which match a given trading pair
func (f TradeFixtures) TradesByAssetPair(a, b xdr.Asset) []Trade {
	set := map[string]bool{}
	var intersection []Trade
	for _, trade := range f.TradesByAsset[a.String()] {
		set[trade.PagingToken()] = true
	}

	for _, trade := range f.TradesByAsset[b.String()] {
		if set[trade.PagingToken()] {
			intersection = append(intersection, trade)
		}
	}
	return intersection
}

// FilterTradesByType filters the given trades by type
func FilterTradesByType(trades []Trade, tradeType string) []Trade {
	var result []Trade
	for _, trade := range trades {
		switch tradeType {
		case AllTrades:
			result = append(result, trade)
		case OrderbookTrades:
			if trade.BaseOfferID.Valid || trade.CounterOfferID.Valid {
				result = append(result, trade)
			}
		case LiquidityPoolTrades:
			if trade.BaseLiquidityPoolID.Valid || trade.CounterLiquidityPoolID.Valid {
				result = append(result, trade)
			}
		}
	}
	return result
}

// TradeScenario inserts trade rows into the aurora DB
func TradeScenario(tt *test.T, q *Q) TradeFixtures {
	builder := q.NewTradeBatchInsertBuilder(0)

	addresses := []string{
		"GB2QIYT2IAUFMRXKLSLLPRECC6OCOGJMADSPTRK7TGNT2SFR2YGWDARD",
		"GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU",
	}
	issuer := xdr.MustAddress("GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H")
	nativeAsset := xdr.MustNewNativeAsset()
	eurAsset := xdr.MustNewCreditAsset("EUR", issuer.Address())
	usdAsset := xdr.MustNewCreditAsset("USD", issuer.Address())

	assets := []xdr.Asset{
		eurAsset,
		usdAsset,
		nativeAsset,
		xdr.MustNewCreditAsset("JPY", addresses[0]),
		xdr.MustNewCreditAsset("CHF", addresses[1]),
	}
	hash := [32]byte{1, 2, 3, 4, 5}
	otherHash := [32]byte{6, 7, 8, 9, 10}
	pools := []string{hex.EncodeToString(hash[:]), hex.EncodeToString(otherHash[:])}
	accountIDs, assetIDs, poolIDs := createHistoryIDs(
		tt, q,
		addresses,
		assets,
		pools,
	)

	inserts := createInsertTrades(accountIDs, assetIDs, poolIDs, 3)

	tt.Assert.NoError(
		builder.Add(tt.Ctx, inserts...),
	)
	tt.Assert.NoError(builder.Exec(tt.Ctx))

	idToAccount := buildIDtoAccountMapping(addresses, accountIDs)
	idToAsset := buildIDtoAssetMapping(assets, assetIDs)

	firstSellerAccount := idToAccount[inserts[0].BaseAccountID.Int64]
	firstBuyerAccount := idToAccount[inserts[0].CounterAccountID.Int64]
	var firstSoldAssetType, firstSoldAssetCode, firstSoldAssetIssuer string
	idToAsset[inserts[0].BaseAssetID].MustExtract(
		&firstSoldAssetType, &firstSoldAssetCode, &firstSoldAssetIssuer,
	)
	var firstBoughtAssetType, firstBoughtAssetCode, firstBoughtAssetIssuer string
	idToAsset[inserts[0].CounterAssetID].MustExtract(
		&firstBoughtAssetType, &firstBoughtAssetCode, &firstBoughtAssetIssuer,
	)

	secondSellerAccount := idToAccount[inserts[1].BaseAccountID.Int64]
	secondBuyerAccount := idToAccount[inserts[1].CounterAccountID.Int64]
	var secondSoldAssetType, secondSoldAssetCode, secondSoldAssetIssuer string
	idToAsset[inserts[1].BaseAssetID].MustExtract(
		&secondSoldAssetType, &secondSoldAssetCode, &secondSoldAssetIssuer,
	)
	var secondBoughtAssetType, secondBoughtAssetCode, secondBoughtAssetIssuer string
	idToAsset[inserts[1].CounterAssetID].MustExtract(
		&secondBoughtAssetType, &secondBoughtAssetCode, &secondBoughtAssetIssuer,
	)

	thirdSellerAccount := idToAccount[inserts[2].BaseAccountID.Int64]
	thirdBuyerAccount := idToAccount[inserts[2].CounterAccountID.Int64]
	var thirdSoldAssetType, thirdSoldAssetCode, thirdSoldAssetIssuer string
	idToAsset[inserts[2].BaseAssetID].MustExtract(
		&thirdSoldAssetType, &thirdSoldAssetCode, &thirdSoldAssetIssuer,
	)
	var thirdBoughtAssetType, thirdBoughtAssetCode, thirdBoughtAssetIssuer string
	idToAsset[inserts[2].CounterAssetID].MustExtract(
		&thirdBoughtAssetType, &thirdBoughtAssetCode, &thirdBoughtAssetIssuer,
	)

	var fourthSoldAssetType, fourthSoldAssetCode, fourthSoldAssetIssuer string
	idToAsset[inserts[3].BaseAssetID].MustExtract(
		&fourthSoldAssetType, &fourthSoldAssetCode, &fourthSoldAssetIssuer,
	)
	var fourthBoughtAssetType, fourthBoughtAssetCode, fourthBoughtAssetIssuer string
	idToAsset[inserts[3].CounterAssetID].MustExtract(
		&fourthBoughtAssetType, &fourthBoughtAssetCode, &fourthBoughtAssetIssuer,
	)

	trades := []Trade{
		{
			HistoryOperationID: inserts[0].HistoryOperationID,
			Order:              inserts[0].Order,
			LedgerCloseTime:    inserts[0].LedgerCloseTime,
			BaseOfferID:        inserts[0].BaseOfferID,
			BaseAccount:        null.StringFrom(firstSellerAccount.Address()),
			BaseAssetType:      firstSoldAssetType,
			BaseAssetIssuer:    firstSoldAssetIssuer,
			BaseAssetCode:      firstSoldAssetCode,
			BaseAmount:         inserts[0].BaseAmount,
			CounterOfferID:     inserts[0].CounterOfferID,
			CounterAccount:     null.StringFrom(firstBuyerAccount.Address()),
			CounterAssetType:   firstBoughtAssetType,
			CounterAssetIssuer: firstBoughtAssetIssuer,
			CounterAssetCode:   firstBoughtAssetCode,
			CounterAmount:      inserts[0].CounterAmount,
			BaseIsSeller:       true,
			PriceN:             null.IntFrom(inserts[0].PriceN),
			PriceD:             null.IntFrom(inserts[0].PriceD),
			Type:               OrderbookTradeType,
		},
		{
			HistoryOperationID: inserts[1].HistoryOperationID,
			Order:              inserts[1].Order,
			LedgerCloseTime:    inserts[1].LedgerCloseTime,
			BaseOfferID:        inserts[1].BaseOfferID,
			BaseAccount:        null.StringFrom(secondSellerAccount.Address()),
			BaseAssetType:      secondSoldAssetType,
			BaseAssetIssuer:    secondSoldAssetIssuer,
			BaseAssetCode:      secondSoldAssetCode,
			BaseAmount:         inserts[1].BaseAmount,
			CounterOfferID:     null.Int{},
			CounterAccount:     null.StringFrom(secondBuyerAccount.Address()),
			CounterAssetType:   secondBoughtAssetType,
			CounterAssetCode:   secondBoughtAssetCode,
			CounterAssetIssuer: secondBoughtAssetIssuer,
			CounterAmount:      inserts[1].CounterAmount,
			BaseIsSeller:       true,
			PriceN:             null.IntFrom(inserts[1].PriceN),
			PriceD:             null.IntFrom(inserts[1].PriceD),
			Type:               OrderbookTradeType,
		},
		{
			HistoryOperationID: inserts[2].HistoryOperationID,
			Order:              inserts[2].Order,
			LedgerCloseTime:    inserts[2].LedgerCloseTime,
			BaseOfferID:        inserts[2].BaseOfferID,
			BaseAccount:        null.StringFrom(thirdSellerAccount.Address()),
			BaseAssetType:      thirdSoldAssetType,
			BaseAssetCode:      thirdSoldAssetCode,
			BaseAssetIssuer:    thirdSoldAssetIssuer,
			BaseAmount:         inserts[2].BaseAmount,
			CounterOfferID:     inserts[2].CounterOfferID,
			CounterAccount:     null.StringFrom(thirdBuyerAccount.Address()),
			CounterAssetType:   thirdBoughtAssetType,
			CounterAssetCode:   thirdBoughtAssetCode,
			CounterAssetIssuer: thirdBoughtAssetIssuer,
			CounterAmount:      inserts[2].CounterAmount,
			BaseIsSeller:       false,
			PriceN:             null.IntFrom(inserts[2].PriceN),
			PriceD:             null.IntFrom(inserts[2].PriceD),
			Type:               OrderbookTradeType,
		},
		{
			HistoryOperationID:  inserts[3].HistoryOperationID,
			Order:               inserts[3].Order,
			LedgerCloseTime:     inserts[3].LedgerCloseTime,
			BaseOfferID:         inserts[3].BaseOfferID,
			BaseAssetType:       fourthSoldAssetType,
			BaseAssetCode:       fourthSoldAssetCode,
			BaseAssetIssuer:     fourthSoldAssetIssuer,
			BaseLiquidityPoolID: null.StringFrom(pools[0]),
			BaseAmount:          inserts[3].BaseAmount,
			CounterOfferID:      null.Int{},
			CounterAccount:      null.StringFrom(thirdSellerAccount.Address()),
			CounterAssetType:    fourthBoughtAssetType,
			CounterAssetCode:    fourthBoughtAssetCode,
			CounterAssetIssuer:  fourthBoughtAssetIssuer,
			CounterAmount:       inserts[3].CounterAmount,
			BaseIsSeller:        inserts[3].BaseIsSeller,
			LiquidityPoolFee:    inserts[3].LiquidityPoolFee,
			PriceN:              null.IntFrom(inserts[3].PriceN),
			PriceD:              null.IntFrom(inserts[3].PriceD),
			Type:                LiquidityPoolTradeType,
		},
		{
			HistoryOperationID:  inserts[4].HistoryOperationID,
			Order:               inserts[4].Order,
			LedgerCloseTime:     inserts[4].LedgerCloseTime,
			BaseOfferID:         inserts[4].BaseOfferID,
			BaseAssetType:       firstSoldAssetType,
			BaseAssetIssuer:     firstSoldAssetIssuer,
			BaseAssetCode:       firstSoldAssetCode,
			BaseLiquidityPoolID: null.StringFrom(pools[1]),
			BaseAmount:          inserts[4].BaseAmount,
			CounterOfferID:      null.Int{},
			CounterAccount:      null.StringFrom(thirdBuyerAccount.Address()),
			CounterAssetType:    firstBoughtAssetType,
			CounterAssetIssuer:  firstBoughtAssetIssuer,
			CounterAssetCode:    firstBoughtAssetCode,
			CounterAmount:       inserts[4].CounterAmount,
			BaseIsSeller:        inserts[4].BaseIsSeller,
			LiquidityPoolFee:    inserts[4].LiquidityPoolFee,
			PriceN:              null.IntFrom(inserts[4].PriceN),
			PriceD:              null.IntFrom(inserts[4].PriceD),
			Type:                LiquidityPoolTradeType,
		},
	}

	fixtures := TradeFixtures{
		Addresses:       addresses,
		Assets:          assets,
		Trades:          trades,
		LiquidityPools:  pools,
		TradesByAccount: map[string][]Trade{},
		TradesByAsset:   map[string][]Trade{},
		TradesByPool:    map[string][]Trade{},
		TradesByOffer:   map[int64][]Trade{},
	}

	for _, trade := range trades {
		if trade.BaseAccount.Valid {
			fixtures.TradesByAccount[trade.BaseAccount.String] = append(fixtures.TradesByAccount[trade.BaseAccount.String], trade)
		}
		if trade.CounterAccount.Valid {
			fixtures.TradesByAccount[trade.CounterAccount.String] = append(fixtures.TradesByAccount[trade.CounterAccount.String], trade)
		}
		baseAsset := strings.Join([]string{trade.BaseAssetType, trade.BaseAssetCode, trade.BaseAssetIssuer}, "/")
		fixtures.TradesByAsset[baseAsset] = append(fixtures.TradesByAsset[baseAsset], trade)

		counterAsset := strings.Join([]string{trade.CounterAssetType, trade.CounterAssetCode, trade.CounterAssetIssuer}, "/")
		fixtures.TradesByAsset[counterAsset] = append(fixtures.TradesByAsset[counterAsset], trade)

		if trade.BaseLiquidityPoolID.Valid {
			fixtures.TradesByPool[trade.BaseLiquidityPoolID.String] = append(fixtures.TradesByPool[trade.BaseLiquidityPoolID.String], trade)
		}
		if trade.CounterLiquidityPoolID.Valid {
			fixtures.TradesByPool[trade.CounterLiquidityPoolID.String] = append(fixtures.TradesByPool[trade.CounterLiquidityPoolID.String], trade)
		}
		if trade.BaseOfferID.Valid {
			fixtures.TradesByOffer[trade.BaseOfferID.Int64] = append(fixtures.TradesByOffer[trade.BaseOfferID.Int64], trade)
		}
		if trade.CounterOfferID.Valid {
			fixtures.TradesByOffer[trade.CounterOfferID.Int64] = append(fixtures.TradesByOffer[trade.CounterOfferID.Int64], trade)
		}
	}

	return fixtures
}
