package history

import (
	"github.com/diamcircle/go/xdr"
	"testing"

	"github.com/diamcircle/go/services/aurora/internal/db2"
	"github.com/diamcircle/go/services/aurora/internal/test"
)

var (
	ascPQ  = db2.MustPageQuery("", false, "asc", 100)
	descPQ = db2.MustPageQuery("", false, "desc", 100)
)

func assertTradesAreEqual(tt *test.T, expected, rows []Trade) {
	tt.Assert.Len(rows, len(expected))
	for i := 0; i < len(rows); i++ {
		tt.Assert.Equal(expected[i].LedgerCloseTime.Unix(), rows[i].LedgerCloseTime.Unix())
		rows[i].LedgerCloseTime = expected[i].LedgerCloseTime
		tt.Assert.Equal(
			expected[i],
			rows[i],
		)
	}
}

const allAccounts = ""

func filterByAccount(trades []Trade, account string) []Trade {
	var result []Trade
	for _, trade := range trades {
		if account == allAccounts ||
			(trade.BaseAccount.Valid && trade.BaseAccount.String == account) ||
			(trade.CounterAccount.Valid && trade.CounterAccount.String == account) {
			result = append(result, trade)
		}
	}
	return result
}

func TestSelectTrades(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}
	fixtures := TradeScenario(tt, q)

	for _, account := range append([]string{allAccounts}, fixtures.Addresses...) {
		for _, tradeType := range []string{AllTrades, OrderbookTrades, LiquidityPoolTrades} {
			expected := filterByAccount(FilterTradesByType(fixtures.Trades, tradeType), account)
			rows, err := q.GetTrades(tt.Ctx, ascPQ, account, tradeType)
			tt.Assert.NoError(err)

			assertTradesAreEqual(tt, expected, rows)

			rows, err = q.GetTrades(tt.Ctx, descPQ, account, tradeType)
			tt.Assert.NoError(err)
			start, end := 0, len(rows)-1
			for start < end {
				rows[start], rows[end] = rows[end], rows[start]
				start++
				end--
			}

			assertTradesAreEqual(tt, expected, rows)
		}
	}
}

func TestSelectTradesCursor(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}
	fixtures := TradeScenario(tt, q)

	for _, account := range append([]string{allAccounts}, fixtures.Addresses...) {
		for _, tradeType := range []string{AllTrades, OrderbookTrades, LiquidityPoolTrades} {
			expected := filterByAccount(FilterTradesByType(fixtures.Trades, tradeType), account)
			if len(expected) == 0 {
				continue
			}

			rows, err := q.GetTrades(
				tt.Ctx,
				db2.MustPageQuery(expected[0].PagingToken(), false, "asc", 100),
				account,
				tradeType,
			)
			tt.Assert.NoError(err)
			assertTradesAreEqual(tt, rows, expected[1:])

			if len(expected) == 1 {
				continue
			}

			rows, err = q.GetTrades(
				tt.Ctx,
				db2.MustPageQuery(expected[1].PagingToken(), false, "asc", 100),
				account,
				tradeType,
			)
			tt.Assert.NoError(err)
			assertTradesAreEqual(tt, rows, expected[2:])
		}
	}
}

func TestTradesQueryForOffer(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}
	fixtures := TradeScenario(tt, q)
	tt.Assert.NotEmpty(fixtures.TradesByOffer)

	for offer, expected := range fixtures.TradesByOffer {
		trades, err := q.GetTradesForOffer(tt.Ctx, ascPQ, offer)
		tt.Assert.NoError(err)
		assertTradesAreEqual(tt, expected, trades)

		trades, err = q.GetTradesForOffer(
			tt.Ctx,
			db2.MustPageQuery(expected[0].PagingToken(), false, "asc", 100),
			offer,
		)
		tt.Assert.NoError(err)
		assertTradesAreEqual(tt, expected[1:], trades)
	}
}

func TestTradesQueryForLiquidityPool(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}
	fixtures := TradeScenario(tt, q)
	tt.Assert.NotEmpty(fixtures.TradesByOffer)

	for poolID, expected := range fixtures.TradesByPool {
		trades, err := q.GetTradesForLiquidityPool(tt.Ctx, ascPQ, poolID)
		tt.Assert.NoError(err)
		assertTradesAreEqual(tt, expected, trades)

		trades, err = q.GetTradesForLiquidityPool(
			tt.Ctx,
			db2.MustPageQuery(expected[0].PagingToken(), false, "asc", 100),
			poolID,
		)
		tt.Assert.NoError(err)
		assertTradesAreEqual(tt, expected[1:], trades)
	}
}

func TestTradesForAssetPair(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}
	fixtures := TradeScenario(tt, q)
	eurAsset := xdr.MustNewCreditAsset("EUR", issuer.Address())
	chfAsset := xdr.MustNewCreditAsset("CHF", "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU")
	allTrades := fixtures.TradesByAssetPair(eurAsset, chfAsset)

	for _, account := range append([]string{allAccounts}, fixtures.Addresses...) {
		for _, tradeType := range []string{AllTrades, OrderbookTrades, LiquidityPoolTrades} {
			expected := filterByAccount(FilterTradesByType(allTrades, tradeType), account)

			trades, err := q.GetTradesForAssets(tt.Ctx, ascPQ, account, tradeType, chfAsset, eurAsset)
			tt.Assert.NoError(err)
			assertTradesAreEqual(tt, expected, trades)

			if len(expected) == 0 {
				continue
			}

			trades, err = q.GetTradesForAssets(
				tt.Ctx,
				db2.MustPageQuery(expected[0].PagingToken(), false, "asc", 100),
				account,
				tradeType,
				chfAsset,
				eurAsset,
			)
			tt.Assert.NoError(err)
			assertTradesAreEqual(tt, expected[1:], trades)
		}
	}
}

func reverseTrade(expected Trade) Trade {
	expected.BaseIsSeller = !expected.BaseIsSeller
	expected.BaseAssetCode, expected.CounterAssetCode = expected.CounterAssetCode, expected.BaseAssetCode
	expected.BaseAssetIssuer, expected.CounterAssetIssuer = expected.CounterAssetIssuer, expected.BaseAssetIssuer
	expected.BaseOfferID, expected.CounterOfferID = expected.CounterOfferID, expected.BaseOfferID
	expected.BaseLiquidityPoolID, expected.CounterLiquidityPoolID = expected.CounterLiquidityPoolID, expected.BaseLiquidityPoolID
	expected.BaseAssetType, expected.CounterAssetType = expected.CounterAssetType, expected.BaseAssetType
	expected.BaseAccount, expected.CounterAccount = expected.CounterAccount, expected.BaseAccount
	expected.BaseAmount, expected.CounterAmount = expected.CounterAmount, expected.BaseAmount
	expected.PriceN, expected.PriceD = expected.PriceD, expected.PriceN
	return expected
}

func TestTradesForReverseAssetPair(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}
	fixtures := TradeScenario(tt, q)
	eurAsset := xdr.MustNewCreditAsset("EUR", issuer.Address())
	chfAsset := xdr.MustNewCreditAsset("CHF", "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU")
	allTrades := fixtures.TradesByAssetPair(eurAsset, chfAsset)

	for _, account := range append([]string{allAccounts}, fixtures.Addresses...) {
		for _, tradeType := range []string{AllTrades, OrderbookTrades, LiquidityPoolTrades} {
			expected := filterByAccount(FilterTradesByType(allTrades, tradeType), account)
			for i := range expected {
				expected[i] = reverseTrade(expected[i])
			}

			trades, err := q.GetTradesForAssets(tt.Ctx, ascPQ, account, tradeType, eurAsset, chfAsset)
			tt.Assert.NoError(err)
			assertTradesAreEqual(tt, expected, trades)

			if len(expected) == 0 {
				continue
			}

			trades, err = q.GetTradesForAssets(
				tt.Ctx,
				db2.MustPageQuery(expected[0].PagingToken(), false, "asc", 100),
				account,
				tradeType,
				eurAsset,
				chfAsset,
			)
			tt.Assert.NoError(err)
			assertTradesAreEqual(tt, expected[1:], trades)
		}
	}
}
