package gql

import (
	"context"
	"errors"

	"github.com/graph-gophers/graphql-go"
	"github.com/diamcircle/go/services/ticker/internal/tickerdb"
	"github.com/diamcircle/go/services/ticker/internal/utils"
)

// Markets resolves the markets() GraphQL query.
func (r *resolver) Markets(ctx context.Context, args struct {
	BaseAssetCode      *string
	BaseAssetIssuer    *string
	CounterAssetCode   *string
	CounterAssetIssuer *string
	NumHoursAgo        *int32
}) (partialMarkets []*partialMarket, err error) {
	numHours, err := validateNumHoursAgo(args.NumHoursAgo)
	if err != nil {
		return
	}

	dbMarkets, err := r.db.RetrievePartialMarkets(ctx,
		args.BaseAssetCode,
		args.BaseAssetIssuer,
		args.CounterAssetCode,
		args.CounterAssetIssuer,
		numHours,
	)
	if err != nil {
		// obfuscating sql errors to avoid exposing underlying
		// implementation
		err = errors.New("could not retrieve the requested data")
		return
	}

	for _, dbMkt := range dbMarkets {
		partialMarkets = append(partialMarkets, dbMarketToPartialMarket(dbMkt))
	}
	return
}

// Ticker resolves the ticker() GraphQL query (TODO)
func (r *resolver) Ticker(ctx context.Context,
	args struct {
		Code        *string
		PairName    *string
		NumHoursAgo *int32
	},
) (partialMarkets []*partialMarket, err error) {
	numHours, err := validateNumHoursAgo(args.NumHoursAgo)
	if err != nil {
		return
	}

	dbMarkets, err := r.db.RetrievePartialAggMarkets(ctx, args.PairName, numHours)
	if err != nil {
		// obfuscating sql errors to avoid exposing underlying
		// implementation
		err = errors.New("could not retrieve the requested data")
		return
	}

	for _, dbMkt := range dbMarkets {
		partialMarkets = append(partialMarkets, dbMarketToPartialMarket(dbMkt))
	}
	return

}

// validateNumHoursAgo validates if the numHoursAgo parameter is within an acceptable
// time range (at most 168 hours ago = 7 days)
func validateNumHoursAgo(n *int32) (int, error) {
	if n == nil {
		return 24, nil // default numHours = 24
	}

	if *n <= 168 {
		return int(*n), nil
	}

	return 0, errors.New("numHoursAgo cannot be greater than 168 (7 days)")
}

// dbMarketToPartialMarket converts a tickerdb.PartialMarket to a *partialMarket
func dbMarketToPartialMarket(dbMarket tickerdb.PartialMarket) *partialMarket {
	spread, spreadMidPoint := utils.CalcSpread(dbMarket.HighestBid, dbMarket.LowestAsk)
	os := orderbookStats{
		BidCount:       BigInt(dbMarket.NumBids),
		BidVolume:      dbMarket.BidVolume,
		BidMax:         dbMarket.HighestBid,
		AskCount:       BigInt(dbMarket.NumAsks),
		AskVolume:      dbMarket.AskVolume,
		AskMin:         dbMarket.LowestAsk,
		Spread:         spread,
		SpreadMidPoint: spreadMidPoint,
	}

	return &partialMarket{
		TradePair:            dbMarket.TradePairName,
		BaseAssetCode:        dbMarket.BaseAssetCode,
		BaseAssetIssuer:      dbMarket.BaseAssetIssuer,
		CounterAssetCode:     dbMarket.CounterAssetCode,
		CounterAssetIssuer:   dbMarket.CounterAssetIssuer,
		BaseVolume:           dbMarket.BaseVolume,
		CounterVolume:        dbMarket.CounterVolume,
		TradeCount:           dbMarket.TradeCount,
		Open:                 dbMarket.Open,
		Low:                  dbMarket.Low,
		High:                 dbMarket.High,
		Change:               dbMarket.Change,
		Close:                dbMarket.Close,
		IntervalStart:        graphql.Time{Time: dbMarket.IntervalStart},
		FirstLedgerCloseTime: graphql.Time{Time: dbMarket.FirstLedgerCloseTime},
		LastLedgerCloseTime:  graphql.Time{Time: dbMarket.LastLedgerCloseTime},
		OrderbookStats:       os,
	}
}
