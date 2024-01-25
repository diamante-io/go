package gql

import (
	"context"
	"errors"

	"go/services/ticker/internal/tickerdb"
)

// Assets resolves the assets() GraphQL query.
func (r *resolver) Assets(ctx context.Context) (assets []*asset, err error) {
	dbAssets, err := r.db.GetAllValidAssets(ctx)
	if err != nil {
		// obfuscating sql errors to avoid exposing underlying
		// implementation
		err = errors.New("could not retrieve the requested data")
		return
	}

	for _, dbAsset := range dbAssets {
		assets = append(assets, dbAssetToAsset(dbAsset))
	}
	return
}

// dbAssetToAsset converts a tickerdb.Asset to an *asset
func dbAssetToAsset(dbAsset tickerdb.Asset) *asset {
	return &asset{
		Code:                        dbAsset.Code,
		IssuerAccount:               dbAsset.IssuerAccount,
		Type:                        dbAsset.Type,
		NumAccounts:                 dbAsset.NumAccounts,
		AuthRequired:                dbAsset.AuthRequired,
		AuthRevocable:               dbAsset.AuthRevocable,
		Amount:                      dbAsset.Amount,
		AssetControlledByDomain:     dbAsset.AssetControlledByDomain,
		AnchorAssetCode:             dbAsset.AnchorAssetCode,
		AnchorAssetType:             dbAsset.AnchorAssetType,
		IsValid:                     dbAsset.IsValid,
		DisplayDecimals:             BigInt(dbAsset.DisplayDecimals),
		Name:                        dbAsset.Name,
		Desc:                        dbAsset.Desc,
		Conditions:                  dbAsset.Conditions,
		IsAssetAnchored:             dbAsset.IsAssetAnchored,
		FixedNumber:                 BigInt(dbAsset.FixedNumber),
		MaxNumber:                   BigInt(dbAsset.MaxNumber),
		IsUnlimited:                 dbAsset.IsUnlimited,
		RedemptionInstructions:      dbAsset.RedemptionInstructions,
		CollateralAddresses:         dbAsset.CollateralAddresses,
		CollateralAddressSignatures: dbAsset.CollateralAddressSignatures,
		Countries:                   dbAsset.Countries,
		Status:                      dbAsset.Status,
		IssuerID:                    dbAsset.IssuerID,
	}
}
