package tickerdb

import (
	"context"
)

// InsertOrUpdateAsset inserts an Asset on the database (if new),
// or updates an existing one
func (s *TickerSession) InsertOrUpdateAsset(ctx context.Context, a *Asset, preserveFields []string) (err error) {
	return s.performUpsertQuery(ctx, *a, "assets", "assets_code_issuer_account", preserveFields)
}

// GetAssetByCodeAndIssuerAccount searches for an Asset with the given code
// and public key, and returns its ID in case it is found.
func (s *TickerSession) GetAssetByCodeAndIssuerAccount(ctx context.Context,
	code string,
	issuerAccount string,
) (found bool, id int32, err error) {
	var assets []Asset
	tbl := s.GetTable("assets")

	err = tbl.Select(
		&assets,
		"assets.code = ? AND assets.issuer_account = ?",
		code,
		issuerAccount,
	).Exec(ctx)
	if err != nil {
		return
	}

	if len(assets) > 0 {
		id = assets[0].ID
		found = true
	}
	return
}

// GetAllValidAssets returns a slice with all assets in the database
// with is_valid = true
func (s *TickerSession) GetAllValidAssets(ctx context.Context) (assets []Asset, err error) {
	tbl := s.GetTable("assets")

	err = tbl.Select(
		&assets,
		"assets.is_valid = TRUE",
	).Exec(ctx)

	return
}

// GetAssetsWithNestedIssuer returns a slice with all assets in the database
// with is_valid = true, also adding the nested Issuer attribute
func (s *TickerSession) GetAssetsWithNestedIssuer(ctx context.Context) (assets []Asset, err error) {
	const q = `
		SELECT
			a.code, a.issuer_account, a.type, a.num_accounts, a.auth_required, a.auth_revocable,
			a.amount, a.asset_controlled_by_domain, a.anchor_asset_code, a.anchor_asset_type,
			a.is_valid, a.validation_error, a.last_valid, a.last_checked, a.display_decimals,
			a.name, a.description, a.conditions, a.is_asset_anchored, a.fixed_number, a.max_number,
			a.is_unlimited, a.redemption_instructions, a.collateral_addresses, a.collateral_address_signatures,
			a.countries, a.status, a.issuer_id, i.public_key, i.name, i.url, i.toml_url, i.federation_server,
			i.auth_server, i.transfer_server, i.web_auth_endpoint, i.deposit_server, i.org_twitter
		FROM assets AS a
		INNER JOIN issuers AS i ON a.issuer_id = i.id
		WHERE a.is_valid = TRUE
	`

	rows, err := s.DB.QueryContext(ctx, q)
	if err != nil {
		return
	}

	for rows.Next() {
		var (
			a Asset
			i Issuer
		)

		err = rows.Scan(
			&a.Code, &a.IssuerAccount, &a.Type, &a.NumAccounts, &a.AuthRequired, &a.AuthRevocable,
			&a.Amount, &a.AssetControlledByDomain, &a.AnchorAssetCode, &a.AnchorAssetType,
			&a.IsValid, &a.ValidationError, &a.LastValid, &a.LastChecked, &a.DisplayDecimals,
			&a.Name, &a.Desc, &a.Conditions, &a.IsAssetAnchored, &a.FixedNumber, &a.MaxNumber,
			&a.IsUnlimited, &a.RedemptionInstructions, &a.CollateralAddresses, &a.CollateralAddressSignatures,
			&a.Countries, &a.Status, &a.IssuerID, &i.PublicKey, &i.Name, &i.URL, &i.TOMLURL, &i.FederationServer,
			&i.AuthServer, &i.TransferServer, &i.WebAuthEndpoint, &i.DepositServer, &i.OrgTwitter,
		)
		if err != nil {
			return
		}

		a.Issuer = i
		assets = append(assets, a)
	}

	return
}
