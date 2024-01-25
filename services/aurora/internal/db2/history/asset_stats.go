package history

import (
	"context"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/diamcircle/go/services/aurora/internal/db2"
	"github.com/diamcircle/go/support/db"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/xdr"
)

func assetStatToMap(assetStat ExpAssetStat) map[string]interface{} {
	return map[string]interface{}{
		"asset_type":   assetStat.AssetType,
		"asset_code":   assetStat.AssetCode,
		"asset_issuer": assetStat.AssetIssuer,
		"accounts":     assetStat.Accounts,
		"balances":     assetStat.Balances,
		"amount":       assetStat.Amount,
		"num_accounts": assetStat.NumAccounts,
	}
}

func assetStatToPrimaryKeyMap(assetStat ExpAssetStat) map[string]interface{} {
	return map[string]interface{}{
		"asset_type":   assetStat.AssetType,
		"asset_code":   assetStat.AssetCode,
		"asset_issuer": assetStat.AssetIssuer,
	}
}

// InsertAssetStats a set of asset stats into the exp_asset_stats
func (q *Q) InsertAssetStats(ctx context.Context, assetStats []ExpAssetStat, batchSize int) error {
	builder := &db.BatchInsertBuilder{
		Table:        q.GetTable("exp_asset_stats"),
		MaxBatchSize: batchSize,
	}

	for _, assetStat := range assetStats {
		if err := builder.Row(ctx, assetStatToMap(assetStat)); err != nil {
			return errors.Wrap(err, "could not insert asset assetStat row")
		}
	}

	if err := builder.Exec(ctx); err != nil {
		return errors.Wrap(err, "could not exec asset assetStats insert builder")
	}

	return nil
}

// InsertAssetStat a single asset assetStat row into the exp_asset_stats
// Returns number of rows affected and error.
func (q *Q) InsertAssetStat(ctx context.Context, assetStat ExpAssetStat) (int64, error) {
	sql := sq.Insert("exp_asset_stats").SetMap(assetStatToMap(assetStat))
	result, err := q.Exec(ctx, sql)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// UpdateAssetStat updates a row in the exp_asset_stats table.
// Returns number of rows affected and error.
func (q *Q) UpdateAssetStat(ctx context.Context, assetStat ExpAssetStat) (int64, error) {
	sql := sq.Update("exp_asset_stats").
		SetMap(assetStatToMap(assetStat)).
		Where(assetStatToPrimaryKeyMap(assetStat))
	result, err := q.Exec(ctx, sql)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// RemoveAssetStat removes a row in the exp_asset_stats table.
func (q *Q) RemoveAssetStat(ctx context.Context, assetType xdr.AssetType, assetCode, assetIssuer string) (int64, error) {
	sql := sq.Delete("exp_asset_stats").
		Where(map[string]interface{}{
			"asset_type":   assetType,
			"asset_code":   assetCode,
			"asset_issuer": assetIssuer,
		})
	result, err := q.Exec(ctx, sql)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// GetAssetStat returns a row in the exp_asset_stats table.
func (q *Q) GetAssetStat(ctx context.Context, assetType xdr.AssetType, assetCode, assetIssuer string) (ExpAssetStat, error) {
	sql := selectAssetStats.Where(map[string]interface{}{
		"asset_type":   assetType,
		"asset_code":   assetCode,
		"asset_issuer": assetIssuer,
	})
	var assetStat ExpAssetStat
	err := q.Get(ctx, &assetStat, sql)
	return assetStat, err
}

func parseAssetStatsCursor(cursor string) (string, string, error) {
	parts := strings.SplitN(cursor, "_", 3)
	if len(parts) != 3 {
		return "", "", fmt.Errorf("invalid asset stats cursor: %v", cursor)
	}

	code, issuer, assetType := parts[0], parts[1], parts[2]
	var issuerAccount xdr.AccountId
	var asset xdr.Asset

	if err := issuerAccount.SetAddress(issuer); err != nil {
		return "", "", errors.Wrap(
			err,
			fmt.Sprintf("invalid issuer in asset stats cursor: %v", cursor),
		)
	}

	if err := asset.SetCredit(code, issuerAccount); err != nil {
		return "", "", errors.Wrap(
			err,
			fmt.Sprintf("invalid asset stats cursor: %v", cursor),
		)
	}

	if _, ok := xdr.StringToAssetType[assetType]; !ok {
		return "", "", errors.Errorf("invalid asset type in asset stats cursor: %v", cursor)
	}

	return code, issuer, nil
}

// GetAssetStats returns a page of exp_asset_stats rows.
func (q *Q) GetAssetStats(ctx context.Context, assetCode, assetIssuer string, page db2.PageQuery) ([]ExpAssetStat, error) {
	sql := selectAssetStats
	filters := map[string]interface{}{}
	if assetCode != "" {
		filters["asset_code"] = assetCode
	}
	if assetIssuer != "" {
		filters["asset_issuer"] = assetIssuer
	}

	if len(filters) > 0 {
		sql = sql.Where(filters)
	}

	var cursorComparison, orderBy string
	switch page.Order {
	case "asc":
		cursorComparison, orderBy = ">", "asc"
	case "desc":
		cursorComparison, orderBy = "<", "desc"
	default:
		return nil, fmt.Errorf("invalid page order %s", page.Order)
	}

	if page.Cursor != "" {
		cursorCode, cursorIssuer, err := parseAssetStatsCursor(page.Cursor)
		if err != nil {
			return nil, err
		}

		sql = sql.Where("((asset_code, asset_issuer) "+cursorComparison+" (?,?))", cursorCode, cursorIssuer)
	}

	sql = sql.OrderBy("(asset_code, asset_issuer) " + orderBy).Limit(page.Limit)

	var results []ExpAssetStat
	if err := q.Select(ctx, &results, sql); err != nil {
		return nil, errors.Wrap(err, "could not run select query")
	}

	return results, nil
}

var selectAssetStats = sq.Select("exp_asset_stats.*").From("exp_asset_stats")
