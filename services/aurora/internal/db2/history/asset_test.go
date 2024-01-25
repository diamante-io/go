package history

import (
	"sort"
	"testing"

	"go/services/aurora/internal/test"
	"go/xdr"
)

func TestCreateAssetsSortedOrder(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)

	q := &Q{tt.auroraSession()}
	assets := []xdr.Asset{
		usdAsset, nativeAsset, eurAsset,
		xdr.MustNewCreditAsset("CNY", issuer.Address()),
	}
	assetMap, err := q.CreateAssets(tt.Ctx,
		assets,
		2,
	)
	tt.Assert.NoError(err)

	idsToAsset := map[int64]string{}
	sortedIDs := []int64{}
	for assetString, asset := range assetMap {
		idsToAsset[asset.ID] = assetString
		sortedIDs = append(sortedIDs, asset.ID)
	}

	sort.Slice(assets, func(i, j int) bool {
		return assets[i].String() < assets[j].String()
	})
	sort.Slice(sortedIDs, func(i, j int) bool {
		return sortedIDs[i] < sortedIDs[j]
	})

	var assetStrings []string
	for _, asset := range assets {
		assetStrings = append(assetStrings, asset.String())
	}

	var values []string
	for _, id := range sortedIDs {
		values = append(values, idsToAsset[id])
	}
	tt.Assert.Equal(assetStrings, values)
}

func TestCreateAssets(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)

	q := &Q{tt.auroraSession()}

	// CreateAssets creates new rows
	assets := []xdr.Asset{
		nativeAsset, eurAsset,
	}
	assetMap, err := q.CreateAssets(tt.Ctx, assets, 1)
	tt.Assert.NoError(err)
	tt.Assert.Len(assetMap, len(assets))

	set := map[int64]bool{}
	for _, asset := range assets {
		row := assetMap[asset.String()]

		tt.Assert.False(set[row.ID])
		set[row.ID] = true

		var assetType, assetCode, assetIssuer string
		asset.MustExtract(&assetType, &assetCode, &assetIssuer)

		tt.Assert.Equal(row.Type, assetType)
		tt.Assert.Equal(row.Code, assetCode)
		tt.Assert.Equal(row.Issuer, assetIssuer)
	}

	// CreateAssets handles duplicates
	assetMap, err = q.CreateAssets(tt.Ctx,
		[]xdr.Asset{
			nativeAsset, nativeAsset, eurAsset, eurAsset,
			nativeAsset, nativeAsset, eurAsset, eurAsset,
		},
		2,
	)
	tt.Assert.NoError(err)
	tt.Assert.Len(assetMap, len(assets))

	for _, asset := range assets {
		row := assetMap[asset.String()]

		tt.Assert.True(set[row.ID])

		var assetType, assetCode, assetIssuer string
		asset.MustExtract(&assetType, &assetCode, &assetIssuer)

		tt.Assert.Equal(row.Type, assetType)
		tt.Assert.Equal(row.Code, assetCode)
		tt.Assert.Equal(row.Issuer, assetIssuer)
	}

	// CreateAssets handles duplicates and new rows
	assets = append(assets, usdAsset)
	assetMap, err = q.CreateAssets(tt.Ctx, assets, 2)
	tt.Assert.NoError(err)
	tt.Assert.Len(assetMap, len(assets))

	for _, asset := range assets {
		row := assetMap[asset.String()]

		inSet := !asset.Equals(usdAsset)
		tt.Assert.Equal(inSet, set[row.ID])

		var assetType, assetCode, assetIssuer string
		asset.MustExtract(&assetType, &assetCode, &assetIssuer)

		tt.Assert.Equal(row.Type, assetType)
		tt.Assert.Equal(row.Code, assetCode)
		tt.Assert.Equal(row.Issuer, assetIssuer)
	}
}
