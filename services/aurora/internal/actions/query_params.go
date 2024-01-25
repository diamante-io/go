package actions

import (
	"fmt"
	"strings"

	"go/support/errors"
	"go/support/render/problem"
	"go/xdr"
)

// SellingBuyingAssetQueryParams query struct for end-points requiring a selling
// and buying asset
type SellingBuyingAssetQueryParams struct {
	SellingAssetType   string `schema:"selling_asset_type" valid:"assetType,optional"`
	SellingAssetIssuer string `schema:"selling_asset_issuer" valid:"accountID,optional"`
	SellingAssetCode   string `schema:"selling_asset_code" valid:"-"`
	BuyingAssetType    string `schema:"buying_asset_type" valid:"assetType,optional"`
	BuyingAssetIssuer  string `schema:"buying_asset_issuer" valid:"accountID,optional"`
	BuyingAssetCode    string `schema:"buying_asset_code" valid:"-"`

	// allow selling and buying using an asset's canonical representation. We
	// are keeping the former selling_* and buying_* for backwards compatibility
	// but it should not be documented.
	SellingAsset string `schema:"selling" valid:"asset,optional"`
	BuyingAsset  string `schema:"buying" valid:"asset,optional"`
}

// Validate runs custom validations buying and selling
func (q SellingBuyingAssetQueryParams) Validate() error {
	ambiguousErr := "Ambiguous parameter, you can't include both `%[1]s` and `%[1]s_asset_type`. Remove all parameters of the form `%[1]s_`"
	if len(q.SellingAssetType) > 0 && len(q.SellingAsset) > 0 {
		return problem.MakeInvalidFieldProblem(
			"selling_asset_type",
			errors.New(fmt.Sprintf(ambiguousErr, "selling")),
		)
	}
	err := validateAssetParams(q.SellingAssetType, q.SellingAssetCode, q.SellingAssetIssuer, "selling_")
	if err != nil {
		return err
	}

	if len(q.BuyingAssetType) > 0 && len(q.BuyingAsset) > 0 {
		return problem.MakeInvalidFieldProblem(
			"buying_asset_type",
			errors.New(fmt.Sprintf(ambiguousErr, "buying")),
		)
	}
	err = validateAssetParams(q.BuyingAssetType, q.BuyingAssetCode, q.BuyingAssetIssuer, "buying_")
	if err != nil {
		return err
	}
	return nil
}

// Selling returns an xdr.Asset representing the selling side of the offer.
func (q SellingBuyingAssetQueryParams) Selling() (*xdr.Asset, error) {
	if len(q.SellingAsset) > 0 {
		switch q.SellingAsset {
		case "native":
			asset := xdr.MustNewNativeAsset()
			return &asset, nil
		default:
			parts := strings.Split(q.SellingAsset, ":")
			if len(parts) != 2 {
				return nil, problem.MakeInvalidFieldProblem(
					"selling",
					errors.New("missing colon"),
				)
			}
			asset, err := xdr.NewCreditAsset(parts[0], parts[1])
			if err != nil {
				return nil, problem.MakeInvalidFieldProblem(
					"selling",
					err,
				)
			}
			return &asset, err
		}
	}

	if len(q.SellingAssetType) == 0 {
		return nil, nil
	}

	selling, err := xdr.BuildAsset(
		q.SellingAssetType,
		q.SellingAssetIssuer,
		q.SellingAssetCode,
	)

	if err != nil {
		p := problem.BadRequest
		p.Extras = map[string]interface{}{"reason": fmt.Sprintf("bad selling asset: %s", err.Error())}
		return nil, p
	}

	return &selling, nil
}

// Buying returns an *xdr.Asset representing the buying side of the offer.
func (q SellingBuyingAssetQueryParams) Buying() (*xdr.Asset, error) {
	if len(q.BuyingAsset) > 0 {
		switch q.BuyingAsset {
		case "native":
			asset := xdr.MustNewNativeAsset()
			return &asset, nil
		default:
			parts := strings.Split(q.BuyingAsset, ":")
			if len(parts) != 2 {
				return nil, problem.MakeInvalidFieldProblem(
					"buying",
					errors.New("missing colon"),
				)
			}
			asset, err := xdr.NewCreditAsset(parts[0], parts[1])
			if err != nil {
				return nil, problem.MakeInvalidFieldProblem(
					"buying",
					err,
				)
			}
			return &asset, err
		}
	}

	if len(q.BuyingAssetType) == 0 {
		return nil, nil
	}

	buying, err := xdr.BuildAsset(
		q.BuyingAssetType,
		q.BuyingAssetIssuer,
		q.BuyingAssetCode,
	)

	if err != nil {
		p := problem.BadRequest
		p.Extras = map[string]interface{}{"reason": fmt.Sprintf("bad buying asset: %s", err.Error())}
		return nil, p
	}

	return &buying, nil
}
