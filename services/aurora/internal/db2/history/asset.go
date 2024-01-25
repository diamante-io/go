package history

import (
	"context"
	"sort"

	sq "github.com/Masterminds/squirrel"
	"github.com/diamcircle/go/support/db"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/xdr"
)

// GetAssetID fetches the id for an Asset
func (q *Q) GetAssetID(ctx context.Context, asset xdr.Asset) (id int64, err error) {
	var (
		assetType   string
		assetCode   string
		assetIssuer string
	)

	err = asset.Extract(&assetType, &assetCode, &assetIssuer)
	if err != nil {
		return
	}

	sql := sq.Select("id").From("history_assets").Limit(1).Where(sq.Eq{
		"asset_type":   assetType,
		"asset_code":   assetCode,
		"asset_issuer": assetIssuer})

	err = q.Get(ctx, &id, sql)
	return
}

// CreateAssets creates rows in the history_assets table for a given list of assets.
func (q *Q) CreateAssets(ctx context.Context, assets []xdr.Asset, batchSize int) (map[string]Asset, error) {
	searchStrings := make([]string, 0, len(assets))
	assetToKey := map[[3]string]string{}

	builder := &db.BatchInsertBuilder{
		Table:        q.GetTable("history_assets"),
		MaxBatchSize: batchSize,
		Suffix:       "ON CONFLICT (asset_code, asset_type, asset_issuer) DO NOTHING",
	}

	// sort assets before inserting rows into history_assets to prevent deadlocks on acquiring a ShareLock
	// https://github.com/diamcircle/go/issues/2370
	sort.Slice(assets, func(i, j int) bool {
		return assets[i].String() < assets[j].String()
	})
	for _, asset := range assets {
		var assetType, assetCode, assetIssuer string
		err := asset.Extract(&assetType, &assetCode, &assetIssuer)
		if err != nil {
			return nil, errors.Wrap(err, "could not extract asset details")
		}

		assetTuple := [3]string{
			assetType,
			assetCode,
			assetIssuer,
		}
		if _, contains := assetToKey[assetTuple]; !contains {
			searchStrings = append(searchStrings, assetType+"/"+assetCode+"/"+assetIssuer)
			assetToKey[assetTuple] = asset.String()

			err = builder.Row(ctx, map[string]interface{}{
				"asset_type":   assetType,
				"asset_code":   assetCode,
				"asset_issuer": assetIssuer,
			})
			if err != nil {
				return nil, errors.Wrap(err, "could not insert history_assets row")
			}
		}
	}

	err := builder.Exec(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "could not exec asset insert builder")
	}
	assetMap := map[string]Asset{}

	const selectBatchSize = 1000
	var rows []Asset
	for i := 0; i < len(searchStrings); i += selectBatchSize {
		end := i + selectBatchSize
		if end > len(searchStrings) {
			end = len(searchStrings)
		}
		subset := searchStrings[i:end]

		err = q.Select(ctx, &rows, sq.Select("*").From("history_assets").Where(sq.Eq{
			"concat(asset_type, '/', asset_code, '/', asset_issuer)": subset,
		}))
		if err != nil {
			return nil, errors.Wrap(err, "could not select assets")
		}

		for _, row := range rows {
			key := assetToKey[[3]string{
				row.Type,
				row.Code,
				row.Issuer,
			}]
			assetMap[key] = row
		}
	}

	return assetMap, nil
}
