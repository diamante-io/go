package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go/support/render/problem"
	"go/xdr"
)

var (
	native = xdr.MustNewNativeAsset()
)

func TestSellingBuyingAssetQueryParams(t *testing.T) {
	testCases := []struct {
		desc                 string
		urlParams            map[string]string
		expectedInvalidField string
		expectedErr          string
	}{
		{
			desc: "Invalid selling_asset_type",
			urlParams: map[string]string{
				"selling_asset_type": "invalid",
			},
			expectedInvalidField: "selling_asset_type",
			expectedErr:          "Asset type must be native, credit_alphanum4 or credit_alphanum12",
		},
		{
			desc: "Invalid buying_asset_type",
			urlParams: map[string]string{
				"buying_asset_type": "invalid",
			},
			expectedInvalidField: "buying_asset_type",
			expectedErr:          "Asset type must be native, credit_alphanum4 or credit_alphanum12",
		}, {
			desc: "Invalid selling_asset_code for credit_alphanum4",
			urlParams: map[string]string{
				"selling_asset_type":   "credit_alphanum4",
				"selling_asset_code":   "invalid",
				"selling_asset_issuer": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
			},
			expectedInvalidField: "selling_asset_code",
			expectedErr:          "Asset code must be 1-12 alphanumeric characters",
		}, {
			desc: "Invalid buying_asset_code for credit_alphanum4",
			urlParams: map[string]string{
				"buying_asset_type": "credit_alphanum4",
				"buying_asset_code": "invalid",
			},
			expectedInvalidField: "buying_asset_code",
			expectedErr:          "Asset code must be 1-12 alphanumeric characters",
		}, {
			desc: "Empty selling_asset_code for credit_alphanum4",
			urlParams: map[string]string{
				"selling_asset_type": "credit_alphanum4",
				"selling_asset_code": "",
			},
			expectedInvalidField: "selling_asset_code",
			expectedErr:          "Asset code must be 1-12 alphanumeric characters",
		}, {
			desc: "Empty buying_asset_code for credit_alphanum4",
			urlParams: map[string]string{
				"buying_asset_type": "credit_alphanum4",
				"buying_asset_code": "",
			},
			expectedInvalidField: "buying_asset_code",
			expectedErr:          "Asset code must be 1-12 alphanumeric characters",
		}, {
			desc: "Empty selling_asset_code for credit_alphanum12",
			urlParams: map[string]string{
				"selling_asset_type": "credit_alphanum12",
				"selling_asset_code": "",
			},
			expectedInvalidField: "selling_asset_code",
			expectedErr:          "Asset code must be 1-12 alphanumeric characters",
		}, {
			desc: "Empty buying_asset_code for credit_alphanum12",
			urlParams: map[string]string{
				"buying_asset_type": "credit_alphanum12",
				"buying_asset_code": "",
			},
			expectedInvalidField: "buying_asset_code",
			expectedErr:          "Asset code must be 1-12 alphanumeric characters",
		}, {
			desc: "Invalid selling_asset_code for credit_alphanum12",
			urlParams: map[string]string{
				"selling_asset_type": "credit_alphanum12",
				"selling_asset_code": "OHLOOOOOOOOOONG",
			},
			expectedInvalidField: "selling_asset_code",
			expectedErr:          "Asset code must be 1-12 alphanumeric characters",
		}, {
			desc: "Invalid buying_asset_code for credit_alphanum12",
			urlParams: map[string]string{
				"buying_asset_type": "credit_alphanum12",
				"buying_asset_code": "OHLOOOOOOOOOONG",
			},
			expectedInvalidField: "buying_asset_code",
			expectedErr:          "Asset code must be 1-12 alphanumeric characters",
		}, {
			desc: "Invalid selling_asset_issuer",
			urlParams: map[string]string{
				"selling_asset_issuer": "GFOOO",
			},
			expectedInvalidField: "selling_asset_issuer",
			expectedErr:          "Account ID must start with `G` and contain 56 alphanum characters",
		}, {
			desc: "Invalid buying_asset_issuer",
			urlParams: map[string]string{
				"buying_asset_issuer": "GFOOO",
			},
			expectedInvalidField: "buying_asset_issuer",
			expectedErr:          "Account ID must start with `G` and contain 56 alphanum characters",
		}, {
			desc: "Missing selling_asset_type",
			urlParams: map[string]string{
				"selling_asset_code":   "OHLOOOOOOOOOONG",
				"selling_asset_issuer": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
			},
			expectedInvalidField: "selling_asset_type",
			expectedErr:          "Missing parameter",
		}, {
			desc: "Missing buying_asset_type",
			urlParams: map[string]string{
				"buying_asset_code":   "OHLOOOOOOOOOONG",
				"buying_asset_issuer": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
			},
			expectedInvalidField: "buying_asset_type",
			expectedErr:          "Missing parameter",
		}, {
			desc: "Missing selling_asset_issuer",
			urlParams: map[string]string{
				"selling_asset_type": "credit_alphanum4",
				"selling_asset_code": "USD",
			},
			expectedInvalidField: "selling_asset_issuer",
			expectedErr:          "Missing parameter",
		}, {
			desc: "Missing buying_asset_issuer",
			urlParams: map[string]string{
				"buying_asset_type": "credit_alphanum4",
				"buying_asset_code": "USD",
			},
			expectedInvalidField: "buying_asset_issuer",
			expectedErr:          "Missing parameter",
		}, {
			desc: "Native with issued asset info: buying_asset_issuer",
			urlParams: map[string]string{
				"buying_asset_type":   "native",
				"buying_asset_issuer": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
			},
			expectedInvalidField: "buying_asset_issuer",
			expectedErr:          "native asset does not have an issuer",
		}, {
			desc: "Native with issued asset info: buying_asset_code",
			urlParams: map[string]string{
				"buying_asset_type": "native",
				"buying_asset_code": "USD",
			},
			expectedInvalidField: "buying_asset_code",
			expectedErr:          "native asset does not have a code",
		}, {
			desc: "Native with issued asset info: selling_asset_issuer",
			urlParams: map[string]string{
				"selling_asset_type":   "native",
				"selling_asset_issuer": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
			},
			expectedInvalidField: "selling_asset_issuer",
			expectedErr:          "native asset does not have an issuer",
		}, {
			desc: "Native with issued asset info: selling_asset_code",
			urlParams: map[string]string{
				"selling_asset_type": "native",
				"selling_asset_code": "USD",
			},
			expectedInvalidField: "selling_asset_code",
			expectedErr:          "native asset does not have a code",
		}, {
			desc: "Valid parameters",
			urlParams: map[string]string{
				"buying_asset_type":    "credit_alphanum4",
				"buying_asset_code":    "USD",
				"buying_asset_issuer":  "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
				"selling_asset_type":   "credit_alphanum4",
				"selling_asset_code":   "EUR",
				"selling_asset_issuer": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
			},
		},
		{
			desc: "Valid parameters with canonical representation",
			urlParams: map[string]string{
				"buying":  "USD:GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
				"selling": "EUR:GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tt := assert.New(t)
			r := makeTestActionRequest("/", tc.urlParams)
			qp := SellingBuyingAssetQueryParams{}
			err := getParams(&qp, r)

			if len(tc.expectedInvalidField) == 0 {
				tt.NoError(err)
			} else {
				if tt.IsType(&problem.P{}, err) {
					p := err.(*problem.P)
					tt.Equal("bad_request", p.Type)
					tt.Equal(tc.expectedInvalidField, p.Extras["invalid_field"])
					tt.Equal(
						tc.expectedErr,
						p.Extras["reason"],
					)
				}
			}

		})
	}
}

func TestSellingBuyingAssetQueryParamsWithCanonicalRepresenation(t *testing.T) {

	testCases := []struct {
		desc                 string
		urlParams            map[string]string
		expectedSelling      *xdr.Asset
		expectedBuying       *xdr.Asset
		expectedInvalidField string
		expectedErr          string
	}{
		{
			desc: "selling native and buying issued asset",
			urlParams: map[string]string{
				"buying":  "native",
				"selling": "EUR:GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
			},
			expectedBuying:  &native,
			expectedSelling: &euro,
		},
		{
			desc: "selling issued and buying native asset",
			urlParams: map[string]string{
				"buying":  "USD:GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
				"selling": "native",
			},
			expectedBuying:  &usd,
			expectedSelling: &native,
		},
		{
			desc: "selling and buying issued assets",
			urlParams: map[string]string{
				"buying":  "USD:GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
				"selling": "EUR:GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
			},
			expectedBuying:  &usd,
			expectedSelling: &euro,
		},
		{
			desc: "new and old format for buying",
			urlParams: map[string]string{
				"buying":              "USD:GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
				"buying_asset_type":   "credit_alphanum4",
				"buying_asset_code":   "USD",
				"buying_asset_issuer": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
			},
			expectedInvalidField: "buying_asset_type",
			expectedErr:          "Ambiguous parameter, you can't include both `buying` and `buying_asset_type`. Remove all parameters of the form `buying_`",
		},
		{
			desc: "new and old format for selling",
			urlParams: map[string]string{
				"selling":              "USD:GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
				"selling_asset_type":   "credit_alphanum4",
				"selling_asset_code":   "USD",
				"selling_asset_issuer": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
			},
			expectedInvalidField: "selling_asset_type",
			expectedErr:          "Ambiguous parameter, you can't include both `selling` and `selling_asset_type`. Remove all parameters of the form `selling_`",
		},
		{
			desc: "invalid selling asset",
			urlParams: map[string]string{
				"selling": "LOLUSD",
			},
			expectedInvalidField: "selling",
			expectedErr:          "Asset must be the string \"native\" or a string of the form \"Code:IssuerAccountID\" for issued assets.",
		},
		{
			desc: "invalid buying asset",
			urlParams: map[string]string{
				"buying": "LOLEUR:",
			},
			expectedInvalidField: "buying",
			expectedErr:          "Asset must be the string \"native\" or a string of the form \"Code:IssuerAccountID\" for issued assets.",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tt := assert.New(t)
			r := makeTestActionRequest("/", tc.urlParams)
			qp := SellingBuyingAssetQueryParams{}
			err := getParams(&qp, r)

			if len(tc.expectedInvalidField) == 0 {
				tt.NoError(err)
				selling, sellingErr := qp.Selling()
				tt.NoError(sellingErr)
				buying, buyingErr := qp.Buying()
				tt.NoError(buyingErr)
				tt.Equal(tc.expectedBuying, buying)
				tt.Equal(tc.expectedSelling, selling)
			} else {
				if tt.IsType(&problem.P{}, err) {
					p := err.(*problem.P)
					tt.Equal("bad_request", p.Type)
					tt.Equal(tc.expectedInvalidField, p.Extras["invalid_field"])
					tt.Equal(
						tc.expectedErr,
						p.Extras["reason"],
					)
				}
			}

		})
	}
}
