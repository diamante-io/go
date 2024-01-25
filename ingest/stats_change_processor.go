package ingest

import (
	"context"

	"go/xdr"
)

// StatsChangeProcessor is a state processors that counts number of changes types
// and entry types.
type StatsChangeProcessor struct {
	results StatsChangeProcessorResults
}

// StatsChangeProcessorResults contains results after running StatsChangeProcessor.
type StatsChangeProcessorResults struct {
	AccountsCreated int64
	AccountsUpdated int64
	AccountsRemoved int64

	ClaimableBalancesCreated int64
	ClaimableBalancesUpdated int64
	ClaimableBalancesRemoved int64

	DataCreated int64
	DataUpdated int64
	DataRemoved int64

	OffersCreated int64
	OffersUpdated int64
	OffersRemoved int64

	TrustLinesCreated int64
	TrustLinesUpdated int64
	TrustLinesRemoved int64

	LiquidityPoolsCreated int64
	LiquidityPoolsUpdated int64
	LiquidityPoolsRemoved int64
}

func (p *StatsChangeProcessor) ProcessChange(ctx context.Context, change Change) error {
	switch change.Type {
	case xdr.LedgerEntryTypeAccount:
		switch change.LedgerEntryChangeType() {
		case xdr.LedgerEntryChangeTypeLedgerEntryCreated:
			p.results.AccountsCreated++
		case xdr.LedgerEntryChangeTypeLedgerEntryUpdated:
			p.results.AccountsUpdated++
		case xdr.LedgerEntryChangeTypeLedgerEntryRemoved:
			p.results.AccountsRemoved++
		}
	case xdr.LedgerEntryTypeClaimableBalance:
		switch change.LedgerEntryChangeType() {
		case xdr.LedgerEntryChangeTypeLedgerEntryCreated:
			p.results.ClaimableBalancesCreated++
		case xdr.LedgerEntryChangeTypeLedgerEntryUpdated:
			p.results.ClaimableBalancesUpdated++
		case xdr.LedgerEntryChangeTypeLedgerEntryRemoved:
			p.results.ClaimableBalancesRemoved++
		}
	case xdr.LedgerEntryTypeData:
		switch change.LedgerEntryChangeType() {
		case xdr.LedgerEntryChangeTypeLedgerEntryCreated:
			p.results.DataCreated++
		case xdr.LedgerEntryChangeTypeLedgerEntryUpdated:
			p.results.DataUpdated++
		case xdr.LedgerEntryChangeTypeLedgerEntryRemoved:
			p.results.DataRemoved++
		}
	case xdr.LedgerEntryTypeOffer:
		switch change.LedgerEntryChangeType() {
		case xdr.LedgerEntryChangeTypeLedgerEntryCreated:
			p.results.OffersCreated++
		case xdr.LedgerEntryChangeTypeLedgerEntryUpdated:
			p.results.OffersUpdated++
		case xdr.LedgerEntryChangeTypeLedgerEntryRemoved:
			p.results.OffersRemoved++
		}
	case xdr.LedgerEntryTypeTrustline:
		switch change.LedgerEntryChangeType() {
		case xdr.LedgerEntryChangeTypeLedgerEntryCreated:
			p.results.TrustLinesCreated++
		case xdr.LedgerEntryChangeTypeLedgerEntryUpdated:
			p.results.TrustLinesUpdated++
		case xdr.LedgerEntryChangeTypeLedgerEntryRemoved:
			p.results.TrustLinesRemoved++
		}
	case xdr.LedgerEntryTypeLiquidityPool:
		switch change.LedgerEntryChangeType() {
		case xdr.LedgerEntryChangeTypeLedgerEntryCreated:
			p.results.LiquidityPoolsCreated++
		case xdr.LedgerEntryChangeTypeLedgerEntryUpdated:
			p.results.LiquidityPoolsUpdated++
		case xdr.LedgerEntryChangeTypeLedgerEntryRemoved:
			p.results.LiquidityPoolsRemoved++
		}
	}

	return nil
}

func (p *StatsChangeProcessor) GetResults() StatsChangeProcessorResults {
	return p.results
}

func (stats *StatsChangeProcessorResults) Map() map[string]interface{} {
	return map[string]interface{}{
		"stats_accounts_created": stats.AccountsCreated,
		"stats_accounts_updated": stats.AccountsUpdated,
		"stats_accounts_removed": stats.AccountsRemoved,

		"stats_claimable_balances_created": stats.ClaimableBalancesCreated,
		"stats_claimable_balances_updated": stats.ClaimableBalancesUpdated,
		"stats_claimable_balances_removed": stats.ClaimableBalancesRemoved,

		"stats_data_created": stats.DataCreated,
		"stats_data_updated": stats.DataUpdated,
		"stats_data_removed": stats.DataRemoved,

		"stats_offers_created": stats.OffersCreated,
		"stats_offers_updated": stats.OffersUpdated,
		"stats_offers_removed": stats.OffersRemoved,

		"stats_trust_lines_created": stats.TrustLinesCreated,
		"stats_trust_lines_updated": stats.TrustLinesUpdated,
		"stats_trust_lines_removed": stats.TrustLinesRemoved,

		"stats_liquidity_pools_created": stats.LiquidityPoolsCreated,
		"stats_liquidity_pools_updated": stats.LiquidityPoolsUpdated,
		"stats_liquidity_pools_removed": stats.LiquidityPoolsRemoved,
	}
}
