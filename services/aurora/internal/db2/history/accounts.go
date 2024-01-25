package history

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/diamcircle/go/services/aurora/internal/db2"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/xdr"
)

// IsAuthRequired returns true if the account has the "AUTH_REQUIRED" option
// turned on.
func (account AccountEntry) IsAuthRequired() bool {
	return xdr.AccountFlags(account.Flags).IsAuthRequired()
}

// IsAuthRevocable returns true if the account has the "AUTH_REVOCABLE" option
// turned on.
func (account AccountEntry) IsAuthRevocable() bool {
	return xdr.AccountFlags(account.Flags).IsAuthRevocable()
}

// IsAuthImmutable returns true if the account has the "AUTH_IMMUTABLE" option
// turned on.
func (account AccountEntry) IsAuthImmutable() bool {
	return xdr.AccountFlags(account.Flags).IsAuthImmutable()
}

// IsAuthClawbackEnabled returns true if the account has the "AUTH_CLAWBACK_ENABLED" option
// turned on.
func (account AccountEntry) IsAuthClawbackEnabled() bool {
	return xdr.AccountFlags(account.Flags).IsAuthClawbackEnabled()
}

func (q *Q) CountAccounts(ctx context.Context) (int, error) {
	sql := sq.Select("count(*)").From("accounts")

	var count int
	if err := q.Get(ctx, &count, sql); err != nil {
		return 0, errors.Wrap(err, "could not run select query")
	}

	return count, nil
}

func (q *Q) GetAccountByID(ctx context.Context, id string) (AccountEntry, error) {
	var account AccountEntry
	sql := selectAccounts.Where(sq.Eq{"account_id": id})
	err := q.Get(ctx, &account, sql)
	return account, err
}

func (q *Q) GetAccountsByIDs(ctx context.Context, ids []string) ([]AccountEntry, error) {
	var accounts []AccountEntry
	sql := selectAccounts.Where(map[string]interface{}{"accounts.account_id": ids})
	err := q.Select(ctx, &accounts, sql)
	return accounts, err
}

// UpsertAccounts upserts a batch of accounts in the accounts table.
// There's currently no limit of the number of accounts this method can
// accept other than 2GB limit of the query string length what should be enough
// for each ledger with the current limits.
func (q *Q) UpsertAccounts(ctx context.Context, accounts []AccountEntry) error {
	var accountID, inflationDestination, homeDomain, balance, buyingLiabilities,
		sellingLiabilities, sequenceNumber, numSubEntries, flags, lastModifiedLedger,
		numSponsored, numSponsoring, masterWeight, thresholdLow, thresholdMedium,
		thresholdHigh, sponsor []interface{}

	for _, account := range accounts {
		accountID = append(accountID, account.AccountID)
		balance = append(balance, account.Balance)
		buyingLiabilities = append(buyingLiabilities, account.BuyingLiabilities)
		sellingLiabilities = append(sellingLiabilities, account.SellingLiabilities)
		sequenceNumber = append(sequenceNumber, account.SequenceNumber)
		numSubEntries = append(numSubEntries, account.NumSubEntries)
		inflationDestination = append(inflationDestination, account.InflationDestination)
		homeDomain = append(homeDomain, account.HomeDomain)
		flags = append(flags, account.Flags)
		masterWeight = append(masterWeight, account.MasterWeight)
		thresholdLow = append(thresholdLow, account.ThresholdLow)
		thresholdMedium = append(thresholdMedium, account.ThresholdMedium)
		thresholdHigh = append(thresholdHigh, account.ThresholdHigh)
		lastModifiedLedger = append(lastModifiedLedger, account.LastModifiedLedger)
		sponsor = append(sponsor, account.Sponsor)
		numSponsored = append(numSponsored, account.NumSponsored)
		numSponsoring = append(numSponsoring, account.NumSponsoring)
	}

	upsertFields := []upsertField{
		{"account_id", "text", accountID},
		{"balance", "bigint", balance},
		{"buying_liabilities", "bigint", buyingLiabilities},
		{"selling_liabilities", "bigint", sellingLiabilities},
		{"sequence_number", "bigint", sequenceNumber},
		{"num_subentries", "int", numSubEntries},
		{"inflation_destination", "text", inflationDestination},
		{"flags", "int", flags},
		{"home_domain", "text", homeDomain},
		{"master_weight", "int", masterWeight},
		{"threshold_low", "int", thresholdLow},
		{"threshold_medium", "int", thresholdMedium},
		{"threshold_high", "int", thresholdHigh},
		{"last_modified_ledger", "int", lastModifiedLedger},
		{"sponsor", "text", sponsor},
		{"num_sponsored", "int", numSponsored},
		{"num_sponsoring", "int", numSponsoring},
	}

	return q.upsertRows(ctx, "accounts", "account_id", upsertFields)
}

// RemoveAccounts deletes a row in the accounts table.
// Returns number of rows affected and error.
func (q *Q) RemoveAccounts(ctx context.Context, accountIDs []string) (int64, error) {
	sql := sq.Delete("accounts").Where(sq.Eq{"account_id": accountIDs})
	result, err := q.Exec(ctx, sql)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// AccountsForAsset returns a list of `AccountEntry` rows who are trustee to an
// asset
func (q *Q) AccountsForAsset(ctx context.Context, asset xdr.Asset, page db2.PageQuery) ([]AccountEntry, error) {
	var assetType, code, issuer string
	asset.MustExtract(&assetType, &code, &issuer)

	sql := sq.
		Select("accounts.*").
		From("accounts").
		Join("trust_lines ON accounts.account_id = trust_lines.account_id").
		Where(map[string]interface{}{
			"trust_lines.asset_type":   int32(asset.Type),
			"trust_lines.asset_issuer": issuer,
			"trust_lines.asset_code":   code,
		})

	sql, err := page.ApplyToUsingCursor(sql, "trust_lines.account_id", page.Cursor)
	if err != nil {
		return nil, errors.Wrap(err, "could not apply query to page")
	}

	var results []AccountEntry
	if err := q.Select(ctx, &results, sql); err != nil {
		return nil, errors.Wrap(err, "could not run select query")
	}

	return results, nil
}

// AccountsForLiquidityPool returns a list of `AccountEntry` rows who are trustee to a
// liquidity pool share asset
func (q *Q) AccountsForLiquidityPool(ctx context.Context, poolID string, page db2.PageQuery) ([]AccountEntry, error) {
	sql := sq.
		Select("accounts.*").
		From("accounts").
		Join("trust_lines ON accounts.account_id = trust_lines.account_id").
		Where(map[string]interface{}{
			"trust_lines.liquidity_pool_id": poolID,
		})

	sql, err := page.ApplyToUsingCursor(sql, "trust_lines.account_id", page.Cursor)
	if err != nil {
		return nil, errors.Wrap(err, "could not apply query to page")
	}

	var results []AccountEntry
	if err := q.Select(ctx, &results, sql); err != nil {
		return nil, errors.Wrap(err, "could not run select query")
	}

	return results, nil
}

func selectBySponsor(table, sponsor string, page db2.PageQuery) (sq.SelectBuilder, error) {
	sql := sq.
		Select("account_id").
		From(table).
		Where(map[string]interface{}{
			"sponsor": sponsor,
		})

	sql, err := page.ApplyToUsingCursor(sql, "account_id", page.Cursor)
	if err != nil {
		return sql, errors.Wrap(err, "could not apply query to page")
	}
	return sql, err
}

func selectUnionBySponsor(tables []string, sponsor string, page db2.PageQuery) (sq.SelectBuilder, error) {
	var selectIDs sq.SelectBuilder
	for i, table := range tables {
		sql, err := selectBySponsor(table, sponsor, page)
		if err != nil {
			return sql, errors.Wrap(err, "could not construct account id query")
		}
		sql = sql.Prefix("(").Suffix(")")

		if i == 0 {
			selectIDs = sql
			continue
		}

		sqlStr, args, err := sql.ToSql()
		if err != nil {
			return sql, errors.Wrap(err, "could not construct account id query")
		}
		selectIDs = selectIDs.Suffix("UNION "+sqlStr, args...)
	}

	return sq.
		Select("accounts.*").
		FromSelect(selectIDs, "accountSet").
		Join("accounts ON accounts.account_id = accountSet.account_id").
		OrderBy("accounts.account_id " + page.Order).
		Limit(page.Limit), nil
}

// AccountsForSponsor return all the accounts where `sponsor`` is sponsoring the account entry or
// any of its subentries (trust lines, signers, data, or account entry)
func (q *Q) AccountsForSponsor(ctx context.Context, sponsor string, page db2.PageQuery) ([]AccountEntry, error) {
	sql, err := selectUnionBySponsor(
		[]string{"accounts", "accounts_data", "accounts_signers", "trust_lines"},
		sponsor,
		page,
	)
	if err != nil {
		return nil, errors.Wrap(err, "could not construct accounts query")
	}

	var results []AccountEntry
	if err := q.Select(ctx, &results, sql); err != nil {
		return nil, errors.Wrap(err, "could not run select query")
	}

	return results, nil
}

// AccountEntriesForSigner returns a list of `AccountEntry` rows for a given signer
func (q *Q) AccountEntriesForSigner(ctx context.Context, signer string, page db2.PageQuery) ([]AccountEntry, error) {
	sql := sq.
		Select("accounts.*").
		From("accounts").
		Join("accounts_signers ON accounts.account_id = accounts_signers.account_id").
		Where(map[string]interface{}{
			"accounts_signers.signer": signer,
		})

	sql, err := page.ApplyToUsingCursor(sql, "accounts_signers.account_id", page.Cursor)
	if err != nil {
		return nil, errors.Wrap(err, "could not apply query to page")
	}

	var results []AccountEntry
	if err := q.Select(ctx, &results, sql); err != nil {
		return nil, errors.Wrap(err, "could not run select query")
	}

	return results, nil
}

var selectAccounts = sq.Select(`
	account_id,
	balance,
	buying_liabilities,
	selling_liabilities,
	sequence_number,
	num_subentries,
	inflation_destination,
	flags,
	home_domain,
	master_weight,
	threshold_low,
	threshold_medium,
	threshold_high,
	last_modified_ledger,
	sponsor,
	num_sponsored,
	num_sponsoring
`).From("accounts")
