package ticker

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	auroraclient "go/clients/auroraclient"
	"go/services/ticker/internal/scraper"
	"go/services/ticker/internal/tickerdb"
	"go/services/ticker/internal/utils"
	hlog "go/support/log"
)

// RefreshAssets scrapes the most recent asset list and ingests then into the db.
func RefreshAssets(ctx context.Context, s *tickerdb.TickerSession, c *auroraclient.Client, l *hlog.Entry) (err error) {
	sc := scraper.ScraperConfig{
		Client: c,
		Logger: l,
	}
	finalAssetList, err := sc.FetchAllAssets(0, 10)
	if err != nil {
		return
	}

	for _, finalAsset := range finalAssetList {
		dbIssuer := tomlIssuerToDBIssuer(finalAsset.IssuerDetails)
		if dbIssuer.PublicKey == "" {
			dbIssuer.PublicKey = finalAsset.Issuer
		}
		issuerID, err := s.InsertOrUpdateIssuer(ctx, &dbIssuer, []string{"public_key"})
		if err != nil {
			l.Error("Error inserting issuer:", dbIssuer, err)
			continue
		}

		dbAsset := finalAssetToDBAsset(finalAsset, issuerID)
		err = s.InsertOrUpdateAsset(ctx, &dbAsset, []string{"code", "issuer_account", "issuer_id"})
		if err != nil {
			l.Error("Error inserting asset:", dbAsset, err)
		}
	}

	return
}

// GenerateAssetsFile generates a file with the info about all valid scraped Assets
func GenerateAssetsFile(ctx context.Context, s *tickerdb.TickerSession, l *hlog.Entry, filename string) error {
	l.Info("Retrieving asset data from db...")
	var assets []Asset
	validAssets, err := s.GetAssetsWithNestedIssuer(ctx)
	if err != nil {
		return err
	}

	for _, dbAsset := range validAssets {
		asset := dbAssetToAsset(dbAsset)
		assets = append(assets, asset)
	}
	l.Info("Asset data successfully retrieved! Writing to: ", filename)
	now := time.Now()
	assetSummary := AssetSummary{
		GeneratedAt:        utils.TimeToUnixEpoch(now),
		GeneratedAtRFC3339: utils.TimeToRFC3339(now),
		Assets:             assets,
	}
	numBytes, err := writeAssetSummaryToFile(assetSummary, filename)
	if err != nil {
		return err
	}
	l.Infof("Wrote %d bytes to %s\n", numBytes, filename)
	return nil
}

// writeAssetSummaryToFile creates a list of assets exported in a JSON file.
func writeAssetSummaryToFile(assetSummary AssetSummary, filename string) (numBytes int, err error) {
	jsonAssets, err := json.MarshalIndent(assetSummary, "", "\t")
	if err != nil {
		return
	}

	numBytes, err = utils.WriteJSONToFile(jsonAssets, filename)
	if err != nil {
		return
	}
	return
}

// finalAssetToDBAsset converts a scraper.TOMLAsset to a tickerdb.Asset.
func finalAssetToDBAsset(asset scraper.FinalAsset, issuerID int32) tickerdb.Asset {
	return tickerdb.Asset{
		Code:                        asset.Code,
		IssuerID:                    issuerID,
		IssuerAccount:               asset.Issuer,
		Type:                        asset.Type,
		NumAccounts:                 asset.NumAccounts,
		AuthRequired:                asset.AuthRequired,
		AuthRevocable:               asset.AuthRevocable,
		Amount:                      asset.Amount,
		AssetControlledByDomain:     asset.AssetControlledByDomain,
		AnchorAssetCode:             asset.AnchorAsset,
		AnchorAssetType:             asset.AnchorAssetType,
		IsValid:                     asset.IsValid,
		ValidationError:             asset.Error,
		LastValid:                   asset.LastValid,
		LastChecked:                 asset.LastChecked,
		DisplayDecimals:             asset.DisplayDecimals,
		Name:                        asset.Name,
		Desc:                        asset.Desc,
		Conditions:                  asset.Conditions,
		IsAssetAnchored:             asset.IsAssetAnchored,
		FixedNumber:                 asset.FixedNumber,
		MaxNumber:                   asset.MaxNumber,
		IsUnlimited:                 asset.IsUnlimited,
		RedemptionInstructions:      asset.RedemptionInstructions,
		CollateralAddresses:         strings.Join(asset.CollateralAddresses, ","),
		CollateralAddressSignatures: strings.Join(asset.CollateralAddressSignatures, ","),
		Countries:                   asset.Countries,
		Status:                      asset.Status,
	}
}

// dbAssetToAsset converts a tickerdb.Asset to an Asset.
func dbAssetToAsset(dbAsset tickerdb.Asset) (a Asset) {
	collAddrs := strings.Split(dbAsset.CollateralAddresses, ",")
	if len(collAddrs) == 1 && collAddrs[0] == "" {
		collAddrs = []string{}
	}

	collAddrSigns := strings.Split(dbAsset.CollateralAddressSignatures, ",")
	if len(collAddrSigns) == 1 && collAddrSigns[0] == "" {
		collAddrSigns = []string{}
	}

	a.Code = dbAsset.Code
	a.Issuer = dbAsset.IssuerAccount
	a.Type = dbAsset.Type
	a.NumAccounts = dbAsset.NumAccounts
	a.AuthRequired = dbAsset.AuthRequired
	a.AuthRevocable = dbAsset.AuthRevocable
	a.Amount = dbAsset.Amount
	a.AssetControlledByDomain = dbAsset.AssetControlledByDomain
	a.AnchorAsset = dbAsset.AnchorAssetCode
	a.AnchorAssetType = dbAsset.AnchorAssetType
	a.LastValidTimestamp = utils.TimeToRFC3339(dbAsset.LastValid)
	a.DisplayDecimals = dbAsset.DisplayDecimals
	a.Name = dbAsset.Name
	a.Desc = dbAsset.Desc
	a.Conditions = dbAsset.Conditions
	a.IsAssetAnchored = dbAsset.IsAssetAnchored
	a.FixedNumber = dbAsset.FixedNumber
	a.MaxNumber = dbAsset.MaxNumber
	a.IsUnlimited = dbAsset.IsUnlimited
	a.RedemptionInstructions = dbAsset.RedemptionInstructions
	a.CollateralAddresses = collAddrs
	a.CollateralAddressSignatures = collAddrSigns
	a.Countries = dbAsset.Countries
	a.Status = dbAsset.Status

	i := Issuer{
		PublicKey:        dbAsset.Issuer.PublicKey,
		Name:             dbAsset.Issuer.Name,
		URL:              dbAsset.Issuer.URL,
		TOMLURL:          dbAsset.Issuer.TOMLURL,
		FederationServer: dbAsset.Issuer.FederationServer,
		AuthServer:       dbAsset.Issuer.AuthServer,
		TransferServer:   dbAsset.Issuer.TransferServer,
		WebAuthEndpoint:  dbAsset.Issuer.WebAuthEndpoint,
		DepositServer:    dbAsset.Issuer.DepositServer,
		OrgTwitter:       dbAsset.Issuer.OrgTwitter,
	}
	a.IssuerDetail = i

	return
}
