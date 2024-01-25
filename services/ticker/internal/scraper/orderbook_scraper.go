package scraper

import (
	"math"
	"strconv"
	"time"

	auroraclient "go/clients/auroraclient"
	hProtocol "go/protocols/aurora"
	"go/services/ticker/internal/utils"

	"github.com/pkg/errors"
)

// fetchOrderbook fetches the orderbook stats for the base and counter assets provided in the parameters
func (c *ScraperConfig) fetchOrderbook(bType, bCode, bIssuer, cType, cCode, cIssuer string) (OrderbookStats, error) {
	var (
		err     error
		summary hProtocol.OrderBookSummary
	)

	obStats := OrderbookStats{
		BaseAssetCode:      bType,
		BaseAssetType:      bCode,
		BaseAssetIssuer:    bIssuer,
		CounterAssetCode:   cType,
		CounterAssetType:   cCode,
		CounterAssetIssuer: cIssuer,
		HighestBid:         math.Inf(-1), // start with -Inf to make sure we catch the correct max bid
		LowestAsk:          math.Inf(1),  // start with +Inf to make sure we catch the correct min ask
	}
	r := createOrderbookRequest(bType, bCode, bIssuer, cType, cCode, cIssuer)

	err = utils.Retry(5, 5*time.Second, c.Logger, func() error {
		summary, err = c.Client.OrderBook(r)
		if err != nil {
			c.Logger.Info("aurora rate limit reached!")
		}
		return err
	})
	if err != nil {
		return obStats, errors.Wrap(err, "could not fetch orderbook summary")
	}

	err = calcOrderbookStats(&obStats, summary)
	if err != nil {
		return obStats, errors.Wrap(err, "could not calculate orderbook stats")
	}
	return obStats, nil
}

// calcOrderbookStats calculates the NumBids, BidVolume, BidMax, NumAsks, AskVolume and AskMin
// statistics for a given OrdebookStats instance
func calcOrderbookStats(obStats *OrderbookStats, summary hProtocol.OrderBookSummary) error {
	// Calculate Bid Data:
	obStats.NumBids = len(summary.Bids)
	if obStats.NumBids == 0 {
		obStats.HighestBid = 0
	}
	for _, bid := range summary.Bids {
		pricef := float64(bid.PriceR.N) / float64(bid.PriceR.D)
		if pricef > obStats.HighestBid {
			obStats.HighestBid = pricef
		}

		amountf, err := strconv.ParseFloat(bid.Amount, 64)
		if err != nil {
			return errors.Wrap(err, "invalid bid amount")
		}
		obStats.BidVolume += amountf
	}

	// Calculate Ask Data:
	obStats.NumAsks = len(summary.Asks)
	if obStats.NumAsks == 0 {
		obStats.LowestAsk = 0
	}
	for _, ask := range summary.Asks {
		pricef := float64(ask.PriceR.N) / float64(ask.PriceR.D)
		amountf, err := strconv.ParseFloat(ask.Amount, 64)
		if err != nil {
			return errors.Wrap(err, "invalid ask amount")
		}

		// On aurora, Ask prices are in units of counter, but
		// amount is in units of base. Therefore, real amount = amount * price
		// See: https://go/issues/612
		obStats.AskVolume += pricef * amountf
		if pricef < obStats.LowestAsk {
			obStats.LowestAsk = pricef
		}
	}

	obStats.Spread, obStats.SpreadMidPoint = utils.CalcSpread(obStats.HighestBid, obStats.LowestAsk)

	// Clean up remaining infinity values:
	if math.IsInf(obStats.LowestAsk, 0) {
		obStats.LowestAsk = 0
	}

	if math.IsInf(obStats.HighestBid, 0) {
		obStats.HighestBid = 0
	}

	return nil
}

// createOrderbookRequest generates a auroraclient.OrderBookRequest based on the base
// and counter asset parameters provided
func createOrderbookRequest(bType, bCode, bIssuer, cType, cCode, cIssuer string) auroraclient.OrderBookRequest {
	r := auroraclient.OrderBookRequest{
		SellingAssetType: auroraclient.AssetType(bType),
		BuyingAssetType:  auroraclient.AssetType(cType),
		// NOTE (Alex C, 2019-05-02):
		// Orderbook requests are currently not paginated on aurora.
		// This limit has been added to ensure we capture at least 200
		// orderbook entries once pagination is added.
		Limit: 200,
	}

	// The aurora API requires *AssetCode and *AssetIssuer fields to be empty
	// when an Asset is native. As we store "XLM" as the asset code for native,
	// we should only add Code and Issuer info in case we're dealing with
	// non-native assets.
	// See: https://developers.diamcircle.org/api/aggregations/order-books/single/
	if bType != string(auroraclient.AssetTypeNative) {
		r.SellingAssetCode = bCode
		r.SellingAssetIssuer = bIssuer
	}
	if cType != string(auroraclient.AssetTypeNative) {
		r.BuyingAssetCode = cCode
		r.BuyingAssetIssuer = cIssuer
	}

	return r
}
