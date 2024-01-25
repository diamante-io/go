package history

import (
	"context"
	"encoding/base64"

	sq "github.com/Masterminds/squirrel"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/xdr"
)

func (q *Q) CountAccountsData(ctx context.Context) (int, error) {
	sql := sq.Select("count(*)").From("accounts_data")

	var count int
	if err := q.Get(ctx, &count, sql); err != nil {
		return 0, errors.Wrap(err, "could not run select query")
	}

	return count, nil
}

// GetAccountDataByName loads account data for a given account ID and data name
func (q *Q) GetAccountDataByName(ctx context.Context, id, name string) (Data, error) {
	var data Data
	sql := selectAccountData.Where(sq.Eq{
		"account_id": id,
		"name":       name,
	}).Limit(1)
	err := q.Get(ctx, &data, sql)
	return data, err
}

// GetAccountDataByAccountID loads account data for a given account ID
func (q *Q) GetAccountDataByAccountID(ctx context.Context, id string) ([]Data, error) {
	var data []Data
	sql := selectAccountData.Where(sq.Eq{"account_id": id})
	err := q.Select(ctx, &data, sql)
	return data, err
}

// GetAccountDataByKeys loads a row from the `accounts_data` table, selected by multiple keys.
func (q *Q) GetAccountDataByKeys(ctx context.Context, keys []AccountDataKey) ([]Data, error) {
	var data []Data
	lkeys := make([]string, 0, len(keys))
	for _, key := range keys {
		lkey, err := accountDataKeyToString(key)
		if err != nil {
			return nil, errors.Wrap(err, "Error running accountDataKeyToString")
		}
		lkeys = append(lkeys, lkey)
	}
	sql := selectAccountData.Where(map[string]interface{}{"accounts_data.ledger_key": lkeys})
	err := q.Select(ctx, &data, sql)
	return data, err
}

func accountDataKeyToString(key AccountDataKey) (string, error) {
	var aid xdr.AccountId
	err := aid.SetAddress(key.AccountID)
	if err != nil {
		return "", err
	}
	var ledgerKey xdr.LedgerKey
	if err = ledgerKey.SetData(aid, key.DataName); err != nil {
		return "", errors.Wrap(err, "Error running ledgerKey.SetData")
	}
	lKey, err := ledgerKey.MarshalBinary()
	if err != nil {
		return "", errors.Wrap(err, "Error running MarshalBinaryCompress")
	}

	return base64.StdEncoding.EncodeToString(lKey), nil
}

// UpsertAccountData upserts a batch of data in the account_Data table.
func (q *Q) UpsertAccountData(ctx context.Context, data []Data) error {
	var ledgerKey, accountID, name, value, lastModifiedLedger, sponsor []interface{}

	for _, d := range data {
		key, err := accountDataKeyToString(AccountDataKey{
			AccountID: d.AccountID,
			DataName:  d.Name,
		})
		if err != nil {
			return err
		}
		ledgerKey = append(ledgerKey, key)
		accountID = append(accountID, d.AccountID)
		name = append(name, d.Name)
		value = append(value, d.Value)
		lastModifiedLedger = append(lastModifiedLedger, d.LastModifiedLedger)
		sponsor = append(sponsor, d.Sponsor)
	}

	upsertFields := []upsertField{
		{"ledger_key", "character varying(150)", ledgerKey},
		{"account_id", "character varying(56)", accountID},
		{"name", "character varying(64)", name},
		{"value", "character varying(90)", value},
		{"last_modified_ledger", "integer", lastModifiedLedger},
		{"sponsor", "text", sponsor},
	}

	return q.upsertRows(ctx, "accounts_data", "ledger_key", upsertFields)
}

// RemoveAccountData deletes a row in the accounts_data table.
// Returns number of rows affected and error.
func (q *Q) RemoveAccountData(ctx context.Context, keys []AccountDataKey) (int64, error) {
	lkeys := make([]string, 0, len(keys))
	for _, key := range keys {
		lkey, err := accountDataKeyToString(key)
		if err != nil {
			return 0, errors.Wrap(err, "Error running accountDataKeyToString")
		}
		lkeys = append(lkeys, lkey)
	}

	sql := sq.Delete("accounts_data").
		Where(sq.Eq{"ledger_key": lkeys})
	result, err := q.Exec(ctx, sql)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// GetAccountDataByAccountsID loads account data for a list of account ID
func (q *Q) GetAccountDataByAccountsID(ctx context.Context, id []string) ([]Data, error) {
	var data []Data
	sql := selectAccountData.Where(sq.Eq{"account_id": id})
	err := q.Select(ctx, &data, sql)
	return data, err
}

var selectAccountData = sq.Select(`
	account_id,
	name,
	value,
	last_modified_ledger,
	sponsor
`).From("accounts_data")
