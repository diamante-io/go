package ticker

import (
	"context"
	"encoding/json"
	"time"

	"go/services/ticker/internal/tickerdb"
	"go/services/ticker/internal/utils"
	hlog "go/support/log"
)

// GenerateMarketSummaryFile generates a MarketSummary with the statistics for all
// valid markets within the database and outputs it to <filename>.
func GenerateMarketSummaryFile(s *tickerdb.TickerSession, l *hlog.Entry, filename string) error {
	l.Info("Generating market data...")
	marketSummary, err := GenerateMarketSummary(s)
	if err != nil {
		return err
	}
	l.Info("Market data successfully generated!")

	jsonMkt, err := json.MarshalIndent(marketSummary, "", "    ")
	if err != nil {
		return err
	}

	l.Info("Writing market data to: ", filename)
	numBytes, err := utils.WriteJSONToFile(jsonMkt, filename)
	if err != nil {
		return err
	}
	l.Infof("Wrote %d bytes to %s\n", numBytes, filename)
	return nil
}

// GenerateMarketSummary outputs a MarketSummary with the statistics for all
// valid markets within the database.
func GenerateMarketSummary(s *tickerdb.TickerSession) (ms MarketSummary, err error) {
	var marketStatsSlice []MarketStats
	now := time.Now()
	nowMillis := utils.TimeToUnixEpoch(now)
	nowRFC339 := utils.TimeToRFC3339(now)
	ctx := context.Background()

	dbMarkets, err := s.RetrieveMarketData(ctx)
	if err != nil {
		return
	}

	for _, dbMarket := range dbMarkets {
		marketStats := dbMarketToMarketStats(dbMarket)
		marketStatsSlice = append(marketStatsSlice, marketStats)
	}

	ms = MarketSummary{
		GeneratedAt:        nowMillis,
		GeneratedAtRFC3339: nowRFC339,
		Pairs:              marketStatsSlice,
	}
	return
}

func dbMarketToMarketStats(m tickerdb.Market) MarketStats {
	closeTime := utils.TimeToRFC3339(m.LastPriceCloseTime)

	spread, spreadMidPoint := utils.CalcSpread(m.HighestBid, m.LowestAsk)
	return MarketStats{
		TradePairName:    m.TradePair,
		BaseVolume24h:    m.BaseVolume24h,
		CounterVolume24h: m.CounterVolume24h,
		TradeCount24h:    m.TradeCount24h,
		Open24h:          m.OpenPrice24h,
		Low24h:           m.LowestPrice24h,
		High24h:          m.HighestPrice24h,
		Change24h:        m.PriceChange24h,
		BaseVolume7d:     m.BaseVolume7d,
		CounterVolume7d:  m.CounterVolume7d,
		TradeCount7d:     m.TradeCount7d,
		Open7d:           m.OpenPrice7d,
		Low7d:            m.LowestPrice7d,
		High7d:           m.HighestPrice7d,
		Change7d:         m.PriceChange7d,
		Price:            m.LastPrice,
		Close:            m.LastPrice,
		BidCount:         m.NumBids,
		BidVolume:        m.BidVolume,
		BidMax:           m.HighestBid,
		AskCount:         m.NumAsks,
		AskVolume:        m.AskVolume,
		AskMin:           m.LowestAsk,
		Spread:           spread,
		SpreadMidPoint:   spreadMidPoint,
		CloseTime:        closeTime,
	}
}
