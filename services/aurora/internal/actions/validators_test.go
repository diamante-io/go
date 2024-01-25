package actions

import (
	"fmt"
	"testing"

	"github.com/asaskevich/govalidator"
	"github.com/stretchr/testify/assert"
)

func TestAssetTypeValidator(t *testing.T) {
	type Query struct {
		AssetType string `valid:"assetType,optional"`
	}

	for _, testCase := range []struct {
		assetType string
		valid     bool
	}{
		{
			"native",
			true,
		},
		{
			"credit_alphanum4",
			true,
		},
		{
			"credit_alphanum12",
			true,
		},
		{
			"",
			true,
		},
		{
			"diamcircle_asset_type",
			false,
		},
	} {
		t.Run(testCase.assetType, func(t *testing.T) {
			tt := assert.New(t)

			q := Query{
				AssetType: testCase.assetType,
			}

			result, err := govalidator.ValidateStruct(q)
			if testCase.valid {
				tt.NoError(err)
				tt.True(result)
			} else {
				tt.Equal("AssetType: diamcircle_asset_type does not validate as assetType", err.Error())
			}
		})
	}
}

func TestAccountIDValidator(t *testing.T) {
	type Query struct {
		Account string `valid:"accountID,optional"`
	}

	for _, testCase := range []struct {
		name          string
		value         string
		expectedError string
	}{
		{
			"invalid diamcircle address",
			"FON4WOTCFSASG3J6SGLLQZURDDUVNBQANAHEQJ3PBNDZ74X63UZWQPZW",
			"Account: FON4WOTCFSASG3J6SGLLQZURDDUVNBQANAHEQJ3PBNDZ74X63UZWQPZW does not validate as accountID",
		},
		{
			"valid diamcircle address",
			"GAN4WOTCFSASG3J6SGLLQZURDDUVNBQANAHEQJ3PBNDZ74X63UZWQPZW",
			"",
		},
		{
			"empty diamcircle address should not be validated",
			"",
			"",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			tt := assert.New(t)

			q := Query{
				Account: testCase.value,
			}

			result, err := govalidator.ValidateStruct(q)
			if testCase.expectedError == "" {
				tt.NoError(err)
				tt.True(result)
			} else {
				tt.Equal(testCase.expectedError, err.Error())
			}
		})
	}
}

func TestAssetValidator(t *testing.T) {
	type Query struct {
		Asset string `valid:"asset"`
	}

	for _, testCase := range []struct {
		desc  string
		asset string
		valid bool
	}{
		{
			"native",
			"native",
			true,
		},
		{
			"credit_alphanum4",
			"USD:GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
			true,
		},
		{
			"credit_alphanum12",
			"SDFUSD:GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
			true,
		},
		{
			"invalid credit_alphanum12",
			"SDFUSDSDFUSDSDFUSD:GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
			false,
		},
		{
			"invalid no issuer",
			"FOO",
			false,
		},
		{
			"invalid issuer",
			"FOO:BAR",
			false,
		},
		{
			"empty colon",
			":",
			false,
		},
	} {
		t.Run(testCase.desc, func(t *testing.T) {
			tt := assert.New(t)

			q := Query{
				Asset: testCase.asset,
			}

			result, err := govalidator.ValidateStruct(q)
			if testCase.valid {
				tt.NoError(err)
				tt.True(result)
			} else {
				tt.Error(err)
			}
		})
	}
}

func TestAmountValidator(t *testing.T) {
	type Query struct {
		Amount string `valid:"amount"`
	}

	for _, testCase := range []struct {
		name          string
		value         string
		expectedError string
	}{
		{
			"valid",
			"10",
			"",
		},
		{
			"zero",
			"0",
			"Amount: 0 does not validate as amount",
		},
		{
			"negative",
			"-1",
			"Amount: -1 does not validate as amount",
		},
		{
			"non-number",
			"one",
			"Amount: one does not validate as amount",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			tt := assert.New(t)

			q := Query{
				Amount: testCase.value,
			}

			result, err := govalidator.ValidateStruct(q)
			if testCase.expectedError == "" {
				tt.NoError(err)
				tt.True(result)
			} else {
				tt.Equal(testCase.expectedError, err.Error())
			}
		})
	}
}

func TestTransactionHashValidator(t *testing.T) {
	type Query struct {
		TransactionHash string `valid:"transactionHash,optional"`
	}

	for _, testCase := range []struct {
		name  string
		value string
		valid bool
	}{
		{
			"length 63",
			"1d2a4be72470658f68db50eef29ea0af3f985ce18b5c218f03461d40c47dc29",
			false,
		},
		{
			"length 66",
			"1d2a4be72470658f68db50eef29ea0af3f985ce18b5c218f03461d40c47dc29222",
			false,
		},
		{
			"uppercase hash",
			"2374E99349B9EF7DBA9A5DB3339B78FDA8F34777B1AF33BA468AD5C0DF946D4D",
			false,
		},
		{
			"badly formated tx hash",
			"%00%1E4%5E%EF%BF%BD%EF%BF%BD%EF%BF%BDpVP%EF%BF%BDI&R%0BK%EF%BF%BD%1D%EF%BF%BD%EF%BF%BD=%EF%BF%BD%3F%23%EF%BF%BD%EF%BF%BDl%EF%BF%BD%1El%EF%BF%BD%EF%BF%BD",
			false,
		},
		{
			"valid tx hash",
			"2374e99349b9ef7dba9a5db3339b78fda8f34777b1af33ba468ad5c0df946d4d",
			true,
		},
		{
			"empty transaction hash should not be validated",
			"",
			true,
		},
		{
			"0x prefixed hash",
			"0x2374e99349b9ef7dba9a5db3339b78fda8f34777b1af33ba468ad5c0df946d4d",
			false,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			tt := assert.New(t)

			q := Query{
				TransactionHash: testCase.value,
			}

			result, err := govalidator.ValidateStruct(q)
			if testCase.valid {
				tt.NoError(err)
				tt.True(result)
			} else {
				expected := fmt.Sprintf("TransactionHash: %s does not validate as transactionHash", testCase.value)
				tt.Equal(expected, err.Error())
			}
		})
	}
}
