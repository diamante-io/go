package scraper

import (
	"time"

	"go/services/ticker/internal/utils"

	auroraclient "go/clients/auroraclient"
	hProtocol "go/protocols/aurora"
)

// checkRecords check if a list of records contains entries older than minTime. If it does,
// it will return a filtered page with only the passing records and lastPage = true.
func (c *ScraperConfig) checkRecords(trades []hProtocol.Trade, minTime time.Time) (lastPage bool, cleanTrades []hProtocol.Trade) {
	lastPage = false
	for _, t := range trades {
		if t.LedgerCloseTime.After(minTime) {
			NormalizeTradeAssets(&t)
			cleanTrades = append(cleanTrades, t)
		} else {
			c.Logger.Debug("Reached entries older than the acceptable time range:", t.LedgerCloseTime)
			lastPage = true
			return
		}
	}
	return
}

// retrieveTrades retrieves trades from the aurora API for the last timeDelta period.
// If limit = 0, will fetch all trades within that period.
func (c *ScraperConfig) retrieveTrades(since time.Time, limit int) (trades []hProtocol.Trade, err error) {
	r := auroraclient.TradeRequest{Limit: 200, Order: auroraclient.OrderDesc}

	tradesPage, err := c.Client.Trades(r)
	if err != nil {
		return
	}

	for tradesPage.Links.Next.Href != tradesPage.Links.Self.Href {
		// Enforcing time boundaries:
		last, cleanTrades := c.checkRecords(tradesPage.Embedded.Records, since)
		trades = append(trades, cleanTrades...)
		if last {
			break
		}

		// Enforcing limit of results:
		if limit != 0 {
			numTrades := len(trades)
			if numTrades >= limit {
				diff := numTrades - limit
				trades = trades[0 : numTrades-diff]
				break
			}
		}

		// Finding next page's params:
		nextURL := tradesPage.Links.Next.Href
		n, err := nextCursor(nextURL)
		if err != nil {
			return trades, err
		}
		c.Logger.Debug("Cursor currently at:", n)
		r.Cursor = n

		err = utils.Retry(5, 5*time.Second, c.Logger, func() error {
			tradesPage, err = c.Client.Trades(r)
			if err != nil {
				c.Logger.Info("aurora rate limit reached!")
			}
			return err
		})
		if err != nil {
			return trades, err
		}
	}

	return
}

// streamTrades streams trades directly from aurora and calls the handler function
// whenever a new trade appears.
func (c *ScraperConfig) streamTrades(h auroraclient.TradeHandler, cursor string) error {
	if cursor == "" {
		cursor = "now"
	}

	r := auroraclient.TradeRequest{
		Limit:  200,
		Cursor: cursor,
	}

	return r.StreamTrades(*c.Ctx, c.Client, h)
}

// addNativeData adds additional fields when one of the assets is native.
func addNativeData(trade *hProtocol.Trade) {
	if trade.BaseAssetType == "native" {
		trade.BaseAssetCode = "XLM"
		trade.BaseAssetIssuer = "native"
	}

	if trade.CounterAssetType == "native" {
		trade.CounterAssetCode = "XLM"
		trade.CounterAssetIssuer = "native"
	}
}

// reverseAssets swaps out the base and counter assets of a trade.
func reverseAssets(trade *hProtocol.Trade) {
	trade.BaseAmount, trade.CounterAmount = trade.CounterAmount, trade.BaseAmount
	trade.BaseAccount, trade.CounterAccount = trade.CounterAccount, trade.BaseAccount
	trade.BaseAssetCode, trade.CounterAssetCode = trade.CounterAssetCode, trade.BaseAssetCode
	trade.BaseAssetType, trade.CounterAssetType = trade.CounterAssetType, trade.BaseAssetType
	trade.BaseAssetIssuer, trade.CounterAssetIssuer = trade.CounterAssetIssuer, trade.BaseAssetIssuer

	trade.BaseIsSeller = !trade.BaseIsSeller
	trade.Price.N, trade.Price.D = trade.Price.D, trade.Price.N
}
