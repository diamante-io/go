package tickerdb

import (
	"context"
)

// InsertOrUpdateOrderbookStats inserts an OrdebookStats entry on the database (if new),
// or updates an existing one
func (s *TickerSession) InsertOrUpdateOrderbookStats(ctx context.Context, o *OrderbookStats, preserveFields []string) (err error) {
	return s.performUpsertQuery(ctx, *o, "orderbook_stats", "orderbook_stats_base_counter_asset_key", preserveFields)
}
