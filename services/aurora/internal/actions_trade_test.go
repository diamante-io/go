//lint:file-ignore U1001 Ignore all unused code, thinks the code is unused because of the test skips
package aurora

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/guregu/null"
	"github.com/diamcircle/go/amount"
	"github.com/diamcircle/go/keypair"

	"github.com/diamcircle/go/protocols/aurora"
	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	diamcircleTime "github.com/diamcircle/go/support/time"
	"github.com/diamcircle/go/xdr"
)

func TestLiquidityPoolTrades(t *testing.T) {
	ht := StartHTTPTestWithoutScenario(t)
	defer ht.Finish()
	dbQ := &history.Q{ht.AuroraSession()}
	fixtures := history.TradeScenario(ht.T, dbQ)

	for _, liquidityPoolID := range fixtures.LiquidityPools {
		expected := fixtures.TradesByPool[liquidityPoolID]
		var records []aurora.Trade
		// All trades
		w := ht.Get("/liquidity_pools/" + liquidityPoolID + "/trades")
		if ht.Assert.Equal(200, w.Code) {
			ht.Assert.PageOf(len(expected), w.Body)
			ht.UnmarshalPage(w.Body, &records)
			for i, row := range expected {
				record := records[i]
				assertResponseTradeEqualsDBTrade(ht, row, record)
			}
		}
	}

	w := ht.Get("/liquidity_pools/" + fixtures.LiquidityPools[0] + "/trades?account_id=" + fixtures.Addresses[0])
	if ht.Assert.Equal(400, w.Code) {
		extras := ht.UnmarshalExtras(w.Body)
		ht.Assert.Equal("account_id,liquidity_pool_id,offer_id", extras["invalid_field"])
		ht.Assert.Equal("Use a single filter for trades, you can only use one of account_id, liquidity_pool_id, offer_id", extras["reason"])
	}

	w = ht.Get("/liquidity_pools/" + fixtures.LiquidityPools[0] + "/trades?offer_id=1")
	if ht.Assert.Equal(400, w.Code) {
		extras := ht.UnmarshalExtras(w.Body)
		ht.Assert.Equal("account_id,liquidity_pool_id,offer_id", extras["invalid_field"])
		ht.Assert.Equal("Use a single filter for trades, you can only use one of account_id, liquidity_pool_id, offer_id", extras["reason"])
	}

	w = ht.Get("/liquidity_pools/" + fixtures.LiquidityPools[0] + "/trades?trade_type=orderbook")
	if ht.Assert.Equal(400, w.Code) {
		extras := ht.UnmarshalExtras(w.Body)
		ht.Assert.Equal("trade_type", extras["invalid_field"])
		ht.Assert.Equal("trade_type orderbook cannot be used with the liquidity_pool_id filter", extras["reason"])
	}
}

func TestOrderbookTrades(t *testing.T) {
	ht := StartHTTPTestWithoutScenario(t)
	defer ht.Finish()
	dbQ := &history.Q{ht.AuroraSession()}
	fixtures := history.TradeScenario(ht.T, dbQ)

	for offerID, expected := range fixtures.TradesByOffer {
		var records []aurora.Trade
		// All trades
		w := ht.Get("/offers/" + strconv.FormatInt(offerID, 10) + "/trades")
		if ht.Assert.Equal(200, w.Code) {
			ht.Assert.PageOf(len(expected), w.Body)
			ht.UnmarshalPage(w.Body, &records)
			for i, row := range expected {
				record := records[i]
				assertResponseTradeEqualsDBTrade(ht, row, record)
			}
		}
	}

	w := ht.Get("/offers/1/trades?account_id=" + fixtures.Addresses[0])
	if ht.Assert.Equal(400, w.Code) {
		extras := ht.UnmarshalExtras(w.Body)
		ht.Assert.Equal("account_id,liquidity_pool_id,offer_id", extras["invalid_field"])
		ht.Assert.Equal("Use a single filter for trades, you can only use one of account_id, liquidity_pool_id, offer_id", extras["reason"])
	}

	w = ht.Get("/offers/1/trades?liquidity_pool_id=" + fixtures.LiquidityPools[0])
	if ht.Assert.Equal(400, w.Code) {
		extras := ht.UnmarshalExtras(w.Body)
		ht.Assert.Equal("account_id,liquidity_pool_id,offer_id", extras["invalid_field"])
		ht.Assert.Equal("Use a single filter for trades, you can only use one of account_id, liquidity_pool_id, offer_id", extras["reason"])
	}

	w = ht.Get("/offers/1/trades?trade_type=liquidity_pool")
	if ht.Assert.Equal(400, w.Code) {
		extras := ht.UnmarshalExtras(w.Body)
		ht.Assert.Equal("trade_type", extras["invalid_field"])
		ht.Assert.Equal("trade_type liquidity_pool cannot be used with the offer_id filter", extras["reason"])
	}
}

func TestAccountTrades(t *testing.T) {
	ht := StartHTTPTestWithoutScenario(t)
	defer ht.Finish()
	dbQ := &history.Q{ht.AuroraSession()}
	fixtures := history.TradeScenario(ht.T, dbQ)

	for _, tradeType := range []string{"", history.AllTrades, history.OrderbookTrades, history.LiquidityPoolTrades} {
		for accountAddress, expected := range fixtures.TradesByAccount {
			var query string
			if tradeType != "" {
				expected = history.FilterTradesByType(expected, tradeType)
				query = "?trade_type=" + tradeType
			}
			var records []aurora.Trade
			// All trades
			w := ht.Get("/accounts/" + accountAddress + "/trades" + query)
			if ht.Assert.Equal(200, w.Code) {
				ht.Assert.PageOf(len(expected), w.Body)
				ht.UnmarshalPage(w.Body, &records)
				for i, row := range expected {
					record := records[i]
					assertResponseTradeEqualsDBTrade(ht, row, record)
				}
			}
		}
	}

	w := ht.Get("/accounts/" + fixtures.Addresses[0] + "/trades?offer_id=1")
	if ht.Assert.Equal(400, w.Code) {
		extras := ht.UnmarshalExtras(w.Body)
		ht.Assert.Equal("account_id,liquidity_pool_id,offer_id", extras["invalid_field"])
		ht.Assert.Equal("Use a single filter for trades, you can only use one of account_id, liquidity_pool_id, offer_id", extras["reason"])
	}

	w = ht.Get("/accounts/" + fixtures.Addresses[0] + "/trades?liquidity_pool_id=" + fixtures.LiquidityPools[0])
	if ht.Assert.Equal(400, w.Code) {
		extras := ht.UnmarshalExtras(w.Body)
		ht.Assert.Equal("account_id,liquidity_pool_id,offer_id", extras["invalid_field"])
		ht.Assert.Equal("Use a single filter for trades, you can only use one of account_id, liquidity_pool_id, offer_id", extras["reason"])
	}

	w = ht.Get("/accounts/" + fixtures.Addresses[0] + "/trades?trade_type=invalid")
	if ht.Assert.Equal(400, w.Code) {
		extras := ht.UnmarshalExtras(w.Body)
		ht.Assert.Equal("trade_type", extras["invalid_field"])
		ht.Assert.Equal("Trade type must be all, orderbook, or liquidity_pool", extras["reason"])
	}
}

func TestTrades(t *testing.T) {
	ht := StartHTTPTestWithoutScenario(t)
	defer ht.Finish()
	dbQ := &history.Q{ht.AuroraSession()}
	fixtures := history.TradeScenario(ht.T, dbQ)

	for _, tradeType := range []string{"", history.AllTrades, history.OrderbookTrades, history.LiquidityPoolTrades} {
		var query string
		expected := fixtures.Trades
		if tradeType != "" {
			expected = history.FilterTradesByType(expected, tradeType)
			query = "trade_type=" + tradeType
		}
		w := ht.Get("/trades?" + query)
		var records []aurora.Trade
		if ht.Assert.Equal(200, w.Code) {
			ht.Assert.PageOf(len(expected), w.Body)
			ht.UnmarshalPage(w.Body, &records)
			for i, row := range expected {
				record := records[i]
				assertResponseTradeEqualsDBTrade(ht, row, record)
			}
		}

		// reverseTrade order
		w = ht.Get("/trades?order=desc&" + query)
		if ht.Assert.Equal(200, w.Code) {
			ht.Assert.PageOf(len(records), w.Body)
			var reverseRecords []aurora.Trade
			ht.UnmarshalPage(w.Body, &reverseRecords)
			ht.Assert.Len(reverseRecords, len(records))

			// ensure that ordering is indeed reversed
			for i := 0; i < len(records); i++ {
				ht.Assert.Equal(records[i], reverseRecords[len(reverseRecords)-1-i])
			}
		}
	}

	w := ht.Get("/trades?trade_type=invalid")
	if ht.Assert.Equal(400, w.Code) {
		extras := ht.UnmarshalExtras(w.Body)
		ht.Assert.Equal("trade_type", extras["invalid_field"])
		ht.Assert.Equal("Trade type must be all, orderbook, or liquidity_pool", extras["reason"])
	}
}

func TestTradesForAssetPair(t *testing.T) {
	ht := StartHTTPTestWithoutScenario(t)
	defer ht.Finish()
	dbQ := &history.Q{ht.AuroraSession()}
	fixtures := history.TradeScenario(ht.T, dbQ)

	q := make(url.Values)
	q.Add("base_asset_type", fixtures.Trades[0].BaseAssetType)
	q.Add("base_asset_code", fixtures.Trades[0].BaseAssetCode)
	q.Add("base_asset_issuer", fixtures.Trades[0].BaseAssetIssuer)
	q.Add("counter_asset_type", fixtures.Trades[0].CounterAssetType)
	q.Add("counter_asset_code", fixtures.Trades[0].CounterAssetCode)
	q.Add("counter_asset_issuer", fixtures.Trades[0].CounterAssetIssuer)

	reverseQ := make(url.Values)
	reverseQ.Add("counter_asset_type", fixtures.Trades[0].BaseAssetType)
	reverseQ.Add("counter_asset_code", fixtures.Trades[0].BaseAssetCode)
	reverseQ.Add("counter_asset_issuer", fixtures.Trades[0].BaseAssetIssuer)
	reverseQ.Add("base_asset_type", fixtures.Trades[0].CounterAssetType)
	reverseQ.Add("base_asset_code", fixtures.Trades[0].CounterAssetCode)
	reverseQ.Add("base_asset_issuer", fixtures.Trades[0].CounterAssetIssuer)

	baseAsset, err := xdr.BuildAsset(
		fixtures.Trades[0].BaseAssetType, fixtures.Trades[0].BaseAssetIssuer, fixtures.Trades[0].BaseAssetCode,
	)
	ht.Assert.NoError(err)
	counterAsset, err := xdr.BuildAsset(
		fixtures.Trades[0].CounterAssetType, fixtures.Trades[0].CounterAssetIssuer, fixtures.Trades[0].CounterAssetCode,
	)
	ht.Assert.NoError(err)

	rows := fixtures.TradesByAssetPair(baseAsset, counterAsset)

	for _, tradeType := range []string{"", history.AllTrades, history.OrderbookTrades, history.LiquidityPoolTrades} {
		expected := rows
		if tradeType != "" {
			expected = history.FilterTradesByType(expected, tradeType)
			q.Set("trade_type", tradeType)
			reverseQ.Set("trade_type", tradeType)
		}

		w := ht.GetWithParams("/trades", q)
		var tradesForPair []aurora.Trade
		if ht.Assert.Equal(200, w.Code) {
			ht.UnmarshalPage(w.Body, &tradesForPair)

			ht.Assert.Equal(len(expected), len(tradesForPair))
			for i, row := range expected {
				assertResponseTradeEqualsDBTrade(ht, row, tradesForPair[i])
			}
		}

		w = ht.GetWithParams("/trades", reverseQ)
		if ht.Assert.Equal(200, w.Code) {
			var trades []aurora.Trade
			ht.UnmarshalPage(w.Body, &trades)
			ht.Assert.Equal(len(tradesForPair), len(trades))

			for i, expected := range tradesForPair {
				ht.Assert.Equal(reverseTrade(expected), trades[i])
			}
		}
	}
}

func reverseTrade(expected aurora.Trade) aurora.Trade {
	expected.Links.Base, expected.Links.Counter = expected.Links.Counter, expected.Links.Base
	expected.BaseIsSeller = !expected.BaseIsSeller
	expected.BaseAssetCode, expected.CounterAssetCode = expected.CounterAssetCode, expected.BaseAssetCode
	expected.BaseAssetIssuer, expected.CounterAssetIssuer = expected.CounterAssetIssuer, expected.BaseAssetIssuer
	expected.BaseOfferID, expected.CounterOfferID = expected.CounterOfferID, expected.BaseOfferID
	expected.BaseLiquidityPoolID, expected.CounterLiquidityPoolID = expected.CounterLiquidityPoolID, expected.BaseLiquidityPoolID
	expected.BaseAssetType, expected.CounterAssetType = expected.CounterAssetType, expected.BaseAssetType
	expected.BaseAccount, expected.CounterAccount = expected.CounterAccount, expected.BaseAccount
	expected.BaseAmount, expected.CounterAmount = expected.CounterAmount, expected.BaseAmount
	expected.Price.N, expected.Price.D = expected.Price.D, expected.Price.N
	return expected
}

func assertResponseTradeEqualsDBTrade(ht *HTTPT, row history.Trade, record aurora.Trade) {
	ht.Assert.Equal(row.BaseAssetCode, record.BaseAssetCode)
	ht.Assert.Equal(row.BaseAssetType, record.BaseAssetType)
	ht.Assert.Equal(row.BaseAssetIssuer, record.BaseAssetIssuer)
	if row.BaseOfferID.Valid {
		ht.Assert.Equal(strconv.FormatInt(row.BaseOfferID.Int64, 10), record.BaseOfferID)
	} else {
		ht.Assert.Equal("", record.BaseOfferID)
	}
	ht.Assert.Equal(row.BaseAmount, int64(amount.MustParse(record.BaseAmount)))
	ht.Assert.Equal(row.BaseLiquidityPoolID.String, record.BaseLiquidityPoolID)
	ht.Assert.Equal(row.BaseAccount.String, record.BaseAccount)
	ht.Assert.Equal(row.BaseLiquidityPoolID.String, record.BaseLiquidityPoolID)
	ht.Assert.Equal(row.BaseIsSeller, record.BaseIsSeller)

	ht.Assert.Equal(row.CounterAssetCode, record.CounterAssetCode)
	ht.Assert.Equal(row.CounterAssetType, record.CounterAssetType)
	ht.Assert.Equal(row.CounterAssetIssuer, record.CounterAssetIssuer)
	if row.CounterOfferID.Valid {
		ht.Assert.Equal(strconv.FormatInt(row.CounterOfferID.Int64, 10), record.CounterOfferID)
	} else {
		ht.Assert.Equal("", record.CounterOfferID)
	}
	ht.Assert.Equal(row.CounterAmount, int64(amount.MustParse(record.CounterAmount)))
	ht.Assert.Equal(row.CounterLiquidityPoolID.String, record.CounterLiquidityPoolID)
	ht.Assert.Equal(row.CounterAccount.String, record.CounterAccount)
	ht.Assert.Equal(row.CounterLiquidityPoolID.String, record.CounterLiquidityPoolID)

	ht.Assert.Equal(uint32(row.LiquidityPoolFee.Int64), record.LiquidityPoolFeeBP)
	ht.Assert.Equal(row.PagingToken(), record.PagingToken())
	ht.Assert.Equal(row.LedgerCloseTime.Unix(), record.LedgerCloseTime.Unix())
	ht.Assert.Equal(row.PriceN.Int64, record.Price.N)
	ht.Assert.Equal(row.PriceD.Int64, record.Price.D)

	switch row.Type {
	case history.OrderbookTradeType:
		ht.Assert.Equal(history.OrderbookTrades, record.TradeType)
	case history.LiquidityPoolTradeType:
		ht.Assert.Equal(history.LiquidityPoolTrades, record.TradeType)
	default:
		ht.Assert.Fail("invalid trade type %v", row.Type)
	}
}

// setAssetQuery adds an asset filter with a given prefix to a query
func setAssetQuery(q *url.Values, prefix string, asset xdr.Asset) {
	var assetType, assetCode, assetFilter string
	asset.Extract(&assetType, &assetCode, &assetFilter)
	q.Add(prefix+"asset_type", assetType)
	q.Add(prefix+"asset_code", assetCode)
	q.Add(prefix+"asset_issuer", assetFilter)
}

// unsetAssetQuery removes an asset filter with a given prefix from a query
func unsetAssetQuery(q *url.Values, prefix string) {
	q.Del(prefix + "asset_type")
	q.Del(prefix + "asset_code")
	q.Del(prefix + "asset_issuer")
}

//testPrice ensures that the price float string is equal to the rational price
func testPrice(t *HTTPT, priceStr string, priceR aurora.TradePrice) {
	price, err := strconv.ParseFloat(priceStr, 64)
	if t.Assert.NoError(err) {
		t.Assert.Equal(price, float64(priceR.N)/float64(priceR.D))
	}
}

func testTradeAggregationPrices(t *HTTPT, record aurora.TradeAggregation) {
	testPrice(t, record.High, record.HighR)
	testPrice(t, record.Low, record.LowR)
	testPrice(t, record.Open, record.OpenR)
	testPrice(t, record.Close, record.CloseR)
}

const minute = int64(time.Minute / time.Millisecond)
const hour = int64(time.Hour / time.Millisecond)
const day = int64(24 * time.Hour / time.Millisecond)
const week = int64(7 * 24 * time.Hour / time.Millisecond)
const aggregationPath = "/trade_aggregations"

func TestTradeActions_Aggregation(t *testing.T) {
	ht := StartHTTPTestWithoutScenario(t)
	defer ht.Finish()

	const numOfTrades = 10

	//a realistic millis (since epoch) value to start the test from
	//it represents a round hour and is bigger than a max int32
	const start = int64(1510693200000)

	dbQ := &history.Q{ht.AuroraSession()}
	ass1, ass2, err := PopulateTestTrades(dbQ, start, numOfTrades, minute, 0)
	ht.Require.NoError(err)

	//add other trades as noise, to ensure asset filtering is working
	_, _, err = PopulateTestTrades(dbQ, start, numOfTrades, minute, numOfTrades)
	ht.Require.NoError(err)

	var records []aurora.TradeAggregation
	var record aurora.TradeAggregation
	var nextLink string

	q := make(url.Values)
	setAssetQuery(&q, "base_", ass1)
	setAssetQuery(&q, "counter_", ass2)

	q.Add("start_time", strconv.FormatInt(start, 10))
	q.Add("end_time", strconv.FormatInt(start+hour, 10))
	q.Add("order", "asc")

	//test no resolution provided
	w := ht.GetWithParams(aggregationPath, q)
	ht.Assert.Equal(400, w.Code)

	//test illegal resolution
	if history.StrictResolutionFiltering {
		q.Add("resolution", strconv.FormatInt(hour/2, 10))
		w = ht.GetWithParams(aggregationPath, q)
		ht.Assert.Equal(400, w.Code)
	}

	//test one bucket for all trades
	q.Set("resolution", strconv.FormatInt(hour, 10))
	w = ht.GetWithParams(aggregationPath, q)
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(1, w.Body)
		ht.UnmarshalPage(w.Body, &records)
		record = records[0] //Save the single aggregation record for next test
		testTradeAggregationPrices(ht, record)
		ht.Assert.Equal("0.0005500", records[0].BaseVolume)
	}

	//test reverseTrade one bucket - make sure values don't change
	q.Set("order", "desc")
	w = ht.GetWithParams(aggregationPath, q)
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(1, w.Body)
		ht.UnmarshalPage(w.Body, &records)
		ht.Assert.Equal(record, records[0])
	}

	//Test bucket per trade
	q.Set("order", "asc")
	q.Set("resolution", strconv.FormatInt(minute, 10))
	w = ht.GetWithParams(aggregationPath, q)
	if ht.Assert.Equal(200, w.Code) {
		if ht.Assert.PageOf(numOfTrades, w.Body) {
			//test that asset filters work
			ht.UnmarshalPage(w.Body, &records)
			ht.Assert.Equal(int64(1), records[0].TradeCount)
			ht.Assert.Equal("0.0000100", records[0].BaseVolume)
			ht.Assert.Equal("1.0000000", records[0].Average)
		}
	}

	//test partial range by modifying endTime to be one minute above half range.
	//half of the results are expected
	endTime := start + (numOfTrades/2)*minute
	q.Set("end_time", strconv.FormatInt(endTime, 10))
	w = ht.GetWithParams(aggregationPath, q)
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(numOfTrades/2, w.Body)
	}

	//test that page limit works
	limit := 3
	q.Add("limit", strconv.Itoa(limit))
	w = ht.GetWithParams(aggregationPath, q)
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(limit, w.Body)
	}

	//test that next page delivers the correct amount of records
	w = ht.GetWithParams(aggregationPath, q)
	nextLink = ht.UnmarshalNext(w.Body)
	//make sure the next link is a full url and not just a path
	ht.Assert.Equal(true, strings.HasPrefix(nextLink, "http://localhost"))
	w = ht.Get(nextLink)
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(numOfTrades/2-limit, w.Body)
		ht.UnmarshalPage(w.Body, &records)
		//test for expected value on timestamp of first record on next page
		ht.Assert.Equal(start+int64(limit)*minute, records[0].Timestamp)
	}

	//test direction (desc)
	q.Set("order", "desc")
	w = ht.GetWithParams(aggregationPath, q)
	if ht.Assert.Equal(200, w.Code) {
		if ht.Assert.PageOf(limit, w.Body) {
			ht.UnmarshalPage(w.Body, &records)
			ht.Assert.Equal(int64(start+(numOfTrades/2-1)*minute), records[0].Timestamp)
		}
	}

	//test next link desc
	w = ht.GetWithParams(aggregationPath, q)
	nextLink = ht.UnmarshalNext(w.Body)
	w = ht.Get(nextLink)
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(numOfTrades/2-limit, w.Body)
	}

	//test next next link empty
	w = ht.GetWithParams(aggregationPath, q)
	nextLink = ht.UnmarshalNext(w.Body)
	w = ht.Get(nextLink)
	nextLink = ht.UnmarshalNext(w.Body)
	w = ht.Get(nextLink)
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(0, w.Body)
	}

	//test non-existent base asset
	foo := GetTestAsset("FOO")

	unsetAssetQuery(&q, "base_")
	setAssetQuery(&q, "base_", foo)

	w = ht.GetWithParams(aggregationPath, q)
	ht.Assert.Equal(404, w.Code)

	jsonErr := map[string]interface{}{}
	err = json.Unmarshal(w.Body.Bytes(), &jsonErr)
	ht.Assert.NoError(err)
	ht.Assert.Equal(float64(404), jsonErr["status"])
	ht.Assert.Equal(
		map[string]interface{}{
			"invalid_field": "base_asset",
			"reason":        "not found",
		},
		jsonErr["extras"],
	)

	unsetAssetQuery(&q, "base_")
	setAssetQuery(&q, "base_", ass1)

	//test non-existent counter asset
	unsetAssetQuery(&q, "counter_")
	setAssetQuery(&q, "counter_", foo)

	w = ht.GetWithParams(aggregationPath, q)
	ht.Assert.Equal(404, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &jsonErr)
	ht.Assert.NoError(err)
	ht.Assert.Equal(float64(404), jsonErr["status"])
	ht.Assert.Equal(
		map[string]interface{}{
			"invalid_field": "counter_asset",
			"reason":        "not found",
		},
		jsonErr["extras"],
	)
}

func TestTradeActions_AmountsExceedInt64(t *testing.T) {
	ht := StartHTTPTestWithoutScenario(t)
	defer ht.Finish()
	dbQ := &history.Q{ht.AuroraSession()}

	const start = int64(1510693200000)

	acc1 := GetTestAccount()
	acc2 := GetTestAccount()
	ass1 := GetTestAsset("euro")
	ass2 := GetTestAsset("usd")
	for i := 1; i <= 3; i++ {
		timestamp := diamcircleTime.MillisFromInt64(start + (minute * int64(i-1)))
		err := IngestTestTrade(
			dbQ, ass1, ass2, acc1, acc2, int64(9131689504000000000), int64(9131689504000000000), timestamp, int64(i))
		ht.Require.NoError(err)
	}

	var records []aurora.TradeAggregation

	q := make(url.Values)
	setAssetQuery(&q, "base_", ass1)
	setAssetQuery(&q, "counter_", ass2)

	q.Add("start_time", strconv.FormatInt(start, 10))
	q.Add("end_time", strconv.FormatInt(start+hour, 10))
	q.Add("order", "asc")
	q.Set("resolution", strconv.FormatInt(hour, 10))

	w := ht.GetWithParams(aggregationPath, q)
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(1, w.Body)
		ht.UnmarshalPage(w.Body, &records)
		ht.Assert.Equal("2739506851200.0000000", records[0].BaseVolume)
		ht.Assert.Equal("2739506851200.0000000", records[0].CounterVolume)
	}
}

func TestTradeActions_IndexRegressions(t *testing.T) {
	t.Run("Assets Dont Exist trades - 404", func(t *testing.T) {
		ht := StartHTTPTestWithoutScenario(t)
		defer ht.Finish()

		var q = make(url.Values)
		q.Add("base_asset_type", "credit_alphanum4")
		q.Add("base_asset_code", "EUR")
		q.Add("base_asset_issuer", "GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG")
		q.Add("counter_asset_type", "native")

		w := ht.Get("/trades?" + q.Encode())

		ht.Assert.Equal(404, w.Code) //This used to be 200 with length 0
	})

	t.Run("Regression for nil prices: https://github.com/diamcircle/go/issues/357", func(t *testing.T) {
		ht := StartHTTPTestWithoutScenario(t)
		dbQ := &history.Q{ht.AuroraSession()}
		history.TradeScenario(ht.T, dbQ)
		defer ht.Finish()

		w := ht.Get("/trades")
		ht.Require.Equal(200, w.Code)

		_ = ht.AuroraDB.MustExec("UPDATE history_trades SET price_n = NULL, price_d = NULL")
		w = ht.Get("/trades")
		ht.Assert.Equal(200, w.Code, "nil-price trades failed")
	})
}

// TestTradeActions_AggregationOrdering checks that open/close aggregation
// fields are correct for multiple trades that occur in the same ledger
// https://github.com/diamcircle/go/issues/215
func TestTradeActions_AggregationOrdering(t *testing.T) {
	ht := StartHTTPTestWithoutScenario(t)
	defer ht.Finish()

	seller := GetTestAccount()
	buyer := GetTestAccount()
	ass1 := GetTestAsset("euro")
	ass2 := GetTestAsset("usd")

	dbQ := &history.Q{ht.AuroraSession()}
	IngestTestTrade(dbQ, ass1, ass2, seller, buyer, 1, 3, 0, 3)
	IngestTestTrade(dbQ, ass1, ass2, seller, buyer, 1, 1, 0, 1)
	IngestTestTrade(dbQ, ass1, ass2, seller, buyer, 1, 2, 0, 2)

	q := make(url.Values)
	setAssetQuery(&q, "base_", ass1)
	setAssetQuery(&q, "counter_", ass2)

	q.Add("start_time", "0")
	q.Add("end_time", "60000")
	q.Add("order", "asc")
	q.Add("resolution", "60000")

	var records []aurora.TradeAggregation
	w := ht.GetWithParams("/trade_aggregations", q)
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(1, w.Body)
		ht.UnmarshalPage(w.Body, &records)
		ht.Assert.Equal("1.0000000", records[0].Open)
		ht.Assert.Equal("3.0000000", records[0].Close)
	}
}

func TestTradeActions_AssetValidation(t *testing.T) {
	ht := StartHTTPTestWithoutScenario(t)
	defer ht.Finish()

	var q = make(url.Values)
	q.Add("base_asset_type", "native")

	w := ht.GetWithParams("/trades", q)
	ht.Assert.Equal(400, w.Code)

	extras := ht.UnmarshalExtras(w.Body)
	ht.Assert.Equal("base_asset_type,counter_asset_type", extras["invalid_field"])
	ht.Assert.Equal("this endpoint supports asset pairs but only one asset supplied", extras["reason"])
}

func TestTradeActions_AggregationInvalidOffset(t *testing.T) {
	ht := StartHTTPTestWithoutScenario(t)
	defer ht.Finish()

	dbQ := &history.Q{ht.AuroraSession()}
	ass1, ass2, err := PopulateTestTrades(dbQ, 0, 100, hour, 1)
	ht.Require.NoError(err)

	q := make(url.Values)
	setAssetQuery(&q, "base_", ass1)
	setAssetQuery(&q, "counter_", ass2)
	q.Add("order", "asc")

	testCases := []struct {
		offset     int64
		resolution int64
		startTime  int64
		endTime    int64
	}{
		{offset: minute, resolution: hour},                                            // Test invalid offset value that's not hour aligned
		{offset: 25 * hour, resolution: week},                                         // Test invalid offset value that's greater than 24 hours
		{offset: 3 * hour, resolution: hour},                                          // Test invalid offset value that's greater than the resolution
		{offset: 3 * hour, startTime: 28 * hour, endTime: 26 * hour, resolution: day}, // Test invalid end time that's less than the start time
		{offset: 3 * hour, startTime: 6 * hour, endTime: 26 * hour, resolution: day},  // Test invalid end time that's less than the offset-adjusted start time
		{offset: 1 * hour, startTime: 5 * hour, endTime: 3 * hour, resolution: day},   // Test invalid end time that's less than the offset-adjusted start time
		{offset: 3 * hour, endTime: 1 * hour, resolution: day},                        // Test invalid end time that's less than the offset
		{startTime: 3 * minute, endTime: 1 * minute, resolution: minute},              // Test invalid end time that's less than the start time (no offset)
	}

	for _, tc := range testCases {
		t.Run("Testing invalid offset parameters", func(t *testing.T) {
			q.Add("offset", strconv.FormatInt(tc.offset, 10))
			q.Add("resolution", strconv.FormatInt(tc.resolution, 10))
			q.Add("start_time", strconv.FormatInt(tc.startTime, 10))
			if tc.endTime != 0 {
				q.Add("end_time", strconv.FormatInt(tc.endTime, 10))
			}
			w := ht.GetWithParams(aggregationPath, q)
			ht.Assert.Equal(400, w.Code)
		})
	}
}

func TestTradeActions_AggregationOffset(t *testing.T) {
	ht := StartHTTPTestWithoutScenario(t)
	defer ht.Finish()

	dbQ := &history.Q{ht.AuroraSession()}
	// One trade every hour
	ass1, ass2, err := PopulateTestTrades(dbQ, 0, 100, hour, 1)
	ht.Require.NoError(err)

	q := make(url.Values)
	setAssetQuery(&q, "base_", ass1)
	setAssetQuery(&q, "counter_", ass2)
	q.Add("order", "asc")

	q.Set("resolution", strconv.FormatInt(day, 10))
	testCases := []struct {
		offset             int64
		startTime          int64
		endTime            int64
		expectedTimestamps []int64
	}{
		{offset: 2 * hour, expectedTimestamps: []int64{2 * hour, 26 * hour, 50 * hour, 74 * hour, 98 * hour}}, //Test with no start time
		{offset: 1 * hour, startTime: 25 * hour, expectedTimestamps: []int64{25 * hour, 49 * hour, 73 * hour, 97 * hour}},
		{offset: 3 * hour, startTime: 10 * hour, expectedTimestamps: []int64{27 * hour, 51 * hour, 75 * hour, 99 * hour}},
		{offset: 6 * hour, startTime: 1 * hour, expectedTimestamps: []int64{6 * hour, 30 * hour, 54 * hour, 78 * hour}},
		{offset: 18 * hour, startTime: 30 * hour, expectedTimestamps: []int64{42 * hour, 66 * hour, 90 * hour}},
		{offset: 10 * hour, startTime: 35 * hour, expectedTimestamps: []int64{58 * hour, 82 * hour}},
		{offset: 18 * hour, startTime: 96 * hour, expectedTimestamps: []int64{}}, // No results since last timestamp is at 100
		{offset: 1 * hour, startTime: 5 * hour, endTime: 95 * hour, expectedTimestamps: []int64{25 * hour, 49 * hour}},
		{offset: 1 * hour, startTime: 5 * hour, endTime: 26 * hour, expectedTimestamps: []int64{}}, // end time and start time should both be at 25 hours
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Testing trade aggregations bucket with offset %d (hour) start time %d (hour)",
			tc.offset/hour, tc.startTime/hour), func(t *testing.T) {
			q.Set("offset", strconv.FormatInt(tc.offset, 10))
			if tc.startTime != 0 {
				q.Set("start_time", strconv.FormatInt(tc.startTime, 10))
			}
			if tc.endTime != 0 {
				q.Set("end_time", strconv.FormatInt(tc.endTime, 10))
			}
			w := ht.GetWithParams(aggregationPath, q)
			if ht.Assert.Equal(200, w.Code) {
				ht.Assert.PageOf(len(tc.expectedTimestamps), w.Body)
				var records []aurora.TradeAggregation
				ht.UnmarshalPage(w.Body, &records)
				if len(records) > 0 {
					for i, record := range records {
						ht.Assert.Equal(tc.expectedTimestamps[i], record.Timestamp)
					}
				}
			}
		})
	}
}

//GetTestAsset generates an issuer on the fly and creates a CreditAlphanum4 Asset with given code
func GetTestAsset(code string) xdr.Asset {
	var codeBytes [4]byte
	copy(codeBytes[:], []byte(code))
	ca4 := xdr.AlphaNum4{Issuer: GetTestAccount(), AssetCode: codeBytes}
	return xdr.Asset{Type: xdr.AssetTypeAssetTypeCreditAlphanum4, AlphaNum4: &ca4, AlphaNum12: nil}
}

//Get generates and returns an account on the fly
func GetTestAccount() xdr.AccountId {
	var key xdr.Uint256
	kp, _ := keypair.Random()
	copy(key[:], kp.Address())
	acc, _ := xdr.NewAccountId(xdr.PublicKeyTypePublicKeyTypeEd25519, key)
	return acc
}

func abs(a xdr.Int32) xdr.Int32 {
	if a < 0 {
		return -a
	}
	return a
}

//IngestTestTrade mock ingests a trade
func IngestTestTrade(
	q *history.Q,
	assetSold xdr.Asset,
	assetBought xdr.Asset,
	seller xdr.AccountId,
	buyer xdr.AccountId,
	amountSold int64,
	amountBought int64,
	timestamp diamcircleTime.Millis,
	opCounter int64) error {

	trade := xdr.ClaimAtom{
		Type: xdr.ClaimAtomTypeClaimAtomTypeOrderBook,
		OrderBook: &xdr.ClaimOfferAtom{
			AmountBought: xdr.Int64(amountBought),
			SellerId:     seller,
			AmountSold:   xdr.Int64(amountSold),
			AssetBought:  assetBought,
			AssetSold:    assetSold,
			OfferId:      100,
		},
	}

	price := xdr.Price{
		N: abs(xdr.Int32(amountBought)),
		D: abs(xdr.Int32(amountSold)),
	}

	ctx := context.Background()
	accounts, err := q.CreateAccounts(ctx, []string{seller.Address(), buyer.Address()}, 2)
	if err != nil {
		return err
	}
	assets, err := q.CreateAssets(ctx, []xdr.Asset{assetBought, assetSold}, 2)
	if err != nil {
		return err
	}

	batch := q.NewTradeBatchInsertBuilder(0)
	batch.Add(ctx, history.InsertTrade{
		HistoryOperationID: opCounter,
		Order:              0,
		CounterAssetID:     assets[assetBought.String()].ID,
		CounterAccountID:   null.IntFrom(accounts[buyer.Address()]),
		CounterAmount:      amountBought,

		BaseAssetID:     assets[assetSold.String()].ID,
		BaseAccountID:   null.IntFrom(accounts[seller.Address()]),
		BaseAmount:      amountSold,
		BaseOfferID:     null.IntFrom(int64(trade.OfferId())),
		BaseIsSeller:    true,
		PriceN:          int64(price.N),
		PriceD:          int64(price.D),
		LedgerCloseTime: timestamp.ToTime(),

		Type: history.OrderbookTradeType,
	})
	err = batch.Exec(ctx)
	if err != nil {
		return err
	}

	err = q.RebuildTradeAggregationTimes(context.Background(), timestamp, timestamp)
	if err != nil {
		return err
	}

	return nil
}

//PopulateTestTrades generates and ingests trades between two assets according to given parameters
func PopulateTestTrades(
	q *history.Q,
	startTs int64,
	numOfTrades int,
	delta int64,
	opStart int64) (ass1 xdr.Asset, ass2 xdr.Asset, err error) {

	acc1 := GetTestAccount()
	acc2 := GetTestAccount()
	ass1 = GetTestAsset("euro")
	ass2 = GetTestAsset("usd")
	for i := 1; i <= numOfTrades; i++ {
		timestamp := diamcircleTime.MillisFromInt64(startTs + (delta * int64(i-1)))
		err = IngestTestTrade(
			q, ass1, ass2, acc1, acc2, int64(i*100), int64(i*100)*int64(i), timestamp, opStart+int64(i))
		//tt.Assert.NoError(err)
		if err != nil {
			return
		}
	}
	return
}
