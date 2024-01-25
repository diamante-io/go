package history

import (
	"context"
)

// TruncateIngestStateTables clears out ingestion state tables.
// Ingestion state tables are aurora database tables populated by
// the ingestion system using history archive snapshots.
// Any aurora database tables which cannot be populated using
// history archive snapshots will not be truncated.
func (q *Q) TruncateIngestStateTables(ctx context.Context) error {
	return q.TruncateTables(ctx, []string{
		"accounts",
		"accounts_data",
		"accounts_signers",
		"claimable_balances",
		"exp_asset_stats",
		"liquidity_pools",
		"offers",
		"trust_lines",
	})
}
