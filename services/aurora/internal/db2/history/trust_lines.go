package history

import (
	"context"

	"go/support/errors"
	"go/xdr"

	sq "github.com/Masterminds/squirrel"
)

// IsAuthorized returns true if issuer has authorized account to perform
// transactions with its credit
func (trustLine TrustLine) IsAuthorized() bool {
	return xdr.TrustLineFlags(trustLine.Flags).IsAuthorized()
}

// IsAuthorizedToMaintainLiabilities returns true if issuer has authorized the account to maintain
// liabilities with its credit
func (trustLine TrustLine) IsAuthorizedToMaintainLiabilities() bool {
	return xdr.TrustLineFlags(trustLine.Flags).IsAuthorizedToMaintainLiabilitiesFlag()
}

// IsClawbackEnabled returns true if issuer has authorized the account to claw
// assets back
func (trustLine TrustLine) IsClawbackEnabled() bool {
	return xdr.TrustLineFlags(trustLine.Flags).IsClawbackEnabledFlag()
}

func (q *Q) CountTrustLines(ctx context.Context) (int, error) {
	sql := sq.Select("count(*)").From("trust_lines")

	var count int
	if err := q.Get(ctx, &count, sql); err != nil {
		return 0, errors.Wrap(err, "could not run select query")
	}

	return count, nil
}

func (q *Q) GetSortedTrustLinesByAccountID(ctx context.Context, id string) ([]TrustLine, error) {
	return q.GetSortedTrustLinesByAccountIDs(ctx, []string{id})
}

// GetTrustLinesByKeys loads a row from the `trust_lines` table, selected by multiple keys.
func (q *Q) GetTrustLinesByKeys(ctx context.Context, ledgerKeys []string) ([]TrustLine, error) {
	var trustLines []TrustLine
	sql := selectTrustLines.Where(map[string]interface{}{"trust_lines.ledger_key": ledgerKeys})
	err := q.Select(ctx, &trustLines, sql)
	return trustLines, err
}

// UpsertTrustLines upserts a batch of trust lines in the trust lines table.
// There's currently no limit of the number of trust lines this method can
// accept other than 2GB limit of the query string length what should be enough
// for each ledger with the current limits.
func (q *Q) UpsertTrustLines(ctx context.Context, trustLines []TrustLine) error {

	var accountID, assetType, assetIssuer, assetCode, balance, ledgerKey, limit, liquidityPoolID, buyingLiabilities,
		sellingLiabilites, flags, lastModifiedLedger, sponsor []interface{}

	for _, trustLine := range trustLines {
		accountID = append(accountID, trustLine.AccountID)
		assetCode = append(assetCode, trustLine.AssetCode)
		assetIssuer = append(assetIssuer, trustLine.AssetIssuer)
		assetType = append(assetType, trustLine.AssetType)
		balance = append(balance, trustLine.Balance)
		ledgerKey = append(ledgerKey, trustLine.LedgerKey)
		limit = append(limit, trustLine.Limit)
		liquidityPoolID = append(liquidityPoolID, trustLine.LiquidityPoolID)
		buyingLiabilities = append(buyingLiabilities, trustLine.BuyingLiabilities)
		sellingLiabilites = append(sellingLiabilites, trustLine.SellingLiabilities)
		flags = append(flags, trustLine.Flags)
		lastModifiedLedger = append(lastModifiedLedger, trustLine.LastModifiedLedger)
		sponsor = append(sponsor, trustLine.Sponsor)
	}

	upsertFields := []upsertField{
		{"account_id", "character varying(56)", accountID},
		{"asset_code", "character varying(12)", assetCode},
		{"asset_issuer", "character varying(56)", assetIssuer},
		{"asset_type", "int", assetType},
		{"balance", "bigint", balance},
		{"ledger_key", "character varying(150)", ledgerKey},
		{"trust_line_limit", "bigint", limit},
		{"liquidity_pool_id", "text", liquidityPoolID},
		{"buying_liabilities", "bigint", buyingLiabilities},
		{"selling_liabilities", "bigint", sellingLiabilites},
		{"flags", "int", flags},
		{"last_modified_ledger", "int", lastModifiedLedger},
		{"sponsor", "text", sponsor},
	}

	return q.upsertRows(ctx, "trust_lines", "ledger_key", upsertFields)
}

// RemoveTrustLine deletes a row in the trust lines table.
// Returns number of rows affected and error.
func (q *Q) RemoveTrustLines(ctx context.Context, ledgerKeys []string) (int64, error) {
	sql := sq.Delete("trust_lines").
		Where(map[string]interface{}{"ledger_key": ledgerKeys})
	result, err := q.Exec(ctx, sql)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// GetSortedTrustLinesByAccountIDs loads trust lines for a list of accounts ID, ordered by asset and issuer
func (q *Q) GetSortedTrustLinesByAccountIDs(ctx context.Context, id []string) ([]TrustLine, error) {
	var data []TrustLine
	sql := selectTrustLines.Where(sq.Eq{"account_id": id}).
		OrderBy("asset_code", "asset_issuer", "liquidity_pool_id")
	err := q.Select(ctx, &data, sql)
	return data, err
}

var selectTrustLines = sq.Select(`
	ledger_key,
	account_id,
	asset_type,
	asset_issuer,
	asset_code,
	COALESCE(liquidity_pool_id, '') as liquidity_pool_id,
	balance,
	trust_line_limit,
	buying_liabilities,
	selling_liabilities,
	flags,
	last_modified_ledger,
	sponsor
`).From("trust_lines")
