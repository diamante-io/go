package history

import (
	"context"
	"database/sql"
	"math/big"

	"github.com/diamcircle/go/amount"
	"github.com/diamcircle/go/support/db"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/xdr"
)

type priceLevel struct {
	Type   string  `db:"type"`
	Pricen int32   `db:"pricen"`
	Priced int32   `db:"priced"`
	Price  float64 `db:"price"`
}

type offerSummary struct {
	Type   string  `db:"type"`
	Amount string  `db:"amount"`
	Price  float64 `db:"price"`
}

// PriceLevel represents an aggregation of offers to trade at a certain
// price.
type PriceLevel struct {
	Pricen int32
	Priced int32
	Pricef string
	Amount string
}

// OrderBookSummary is a summary of a set of offers for a given base and
// counter currency
type OrderBookSummary struct {
	Asks []PriceLevel
	Bids []PriceLevel
}

// GetOrderBookSummary returns an OrderBookSummary for a given trading pair.
// GetOrderBookSummary should only be called in a repeatable read transaction.
func (q *Q) GetOrderBookSummary(ctx context.Context, sellingAsset, buyingAsset xdr.Asset, maxPriceLevels int) (OrderBookSummary, error) {
	var result OrderBookSummary

	if tx := q.GetTx(); tx == nil {
		return result, errors.New("cannot be called outside of a transaction")
	}
	if opts := q.GetTxOptions(); opts == nil || !opts.ReadOnly || opts.Isolation != sql.LevelRepeatableRead {
		return result, errors.New("should only be called in a repeatable read transaction")
	}

	selling, err := xdr.MarshalBase64(sellingAsset)
	if err != nil {
		return result, errors.Wrap(err, "cannot marshal selling asset")
	}
	buying, err := xdr.MarshalBase64(buyingAsset)
	if err != nil {
		return result, errors.Wrap(err, "cannot marshal Buying asset")
	}

	var levels []priceLevel
	// First, obtain the price fractions for each price level.
	// In the next query, we'll sum the amounts for each price level.
	// Finally, we will combine the results to produce a OrderBookSummary.
	selectPriceLevels := `
		(SELECT DISTINCT ON (price)
			'ask' as type, pricen, priced, price
		FROM offers
		WHERE selling_asset = $1 AND buying_asset = $2 AND deleted = false
		ORDER BY price ASC LIMIT $3)
		UNION ALL
		(SELECT DISTINCT ON (price)
			'bid' as type, pricen, priced, price
		FROM offers
		WHERE selling_asset = $2 AND buying_asset = $1 AND deleted = false
		ORDER BY price ASC LIMIT $3)
	`

	var offers []offerSummary
	// The SUM() value in postgres has type decimal which means it will
	// handle values that exceed max int64 so we don't need to worry about
	// overflows.
	selectOfferSummaries := `
		(
			SELECT
				'ask' as type, co.price, SUM(co.amount) as amount
			FROM  offers co
			WHERE selling_asset = $1 AND buying_asset = $2 AND deleted = false
			GROUP BY co.price
			ORDER BY co.price ASC
			LIMIT $3
		) UNION ALL (
			SELECT
				'bid'  as type, co.price, SUM(co.amount) as amount
			FROM offers co
			WHERE selling_asset = $2 AND buying_asset = $1 AND deleted = false
			GROUP BY co.price
			ORDER BY co.price ASC
			LIMIT $3
		)
	`
	// Add explicit query type for prometheus metrics, since we use raw sql.
	ctx = context.WithValue(ctx, &db.QueryTypeContextKey, db.SelectQueryType)
	err = q.SelectRaw(ctx, &levels, selectPriceLevels, selling, buying, maxPriceLevels)
	if err != nil {
		return result, errors.Wrap(err, "cannot select price levels")
	}

	err = q.SelectRaw(ctx, &offers, selectOfferSummaries, selling, buying, maxPriceLevels)
	if err != nil {
		return result, errors.Wrap(err, "cannot select offer summaries")
	}

	// we don't expect there to be any inconsistency between levels and offers because
	// this function should only be invoked in a repeatable read transaction
	if len(levels) != len(offers) {
		return result, errors.Wrap(err, "price levels length does not match summaries length")
	}
	for i, level := range levels {
		sum := offers[i]
		if level.Type != sum.Type {
			return result, errors.Wrap(err, "price level type does not match offer summary type")
		}
		if level.Price != sum.Price {
			return result, errors.Wrap(err, "price level price does not match offer summary price")
		}
		// use big.Rat to get reduced fractions
		priceFraction := big.NewRat(int64(level.Pricen), int64(level.Priced))
		if sum.Type == "bid" {
			// only invert bids
			if level.Pricen == 0 {
				return result, errors.Wrap(err, "bid has price denominator equal to 0")
			}
			priceFraction = priceFraction.Inv(priceFraction)
		}

		entry := PriceLevel{
			Pricef: priceFraction.FloatString(7),
			Pricen: int32(priceFraction.Num().Int64()),
			Priced: int32(priceFraction.Denom().Int64()),
		}
		entry.Amount, err = amount.IntStringToAmount(sum.Amount)
		if err != nil {
			return result, errors.Wrap(err, "could not determine summary amount")
		}
		if sum.Type == "ask" {
			result.Asks = append(result.Asks, entry)
		} else if sum.Type == "bid" {
			result.Bids = append(result.Bids, entry)
		} else {
			return result, errors.Wrap(err, "invalid offer type")
		}
	}

	return result, nil
}
