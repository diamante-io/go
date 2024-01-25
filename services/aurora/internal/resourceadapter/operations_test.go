package resourceadapter

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/guregu/null"
	"github.com/diamcircle/go/protocols/aurora/operations"
	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/support/test"
	"github.com/diamcircle/go/xdr"
	"github.com/stretchr/testify/assert"
)

func TestNewOperationAllTypesCovered(t *testing.T) {
	for typ, s := range xdr.OperationTypeToStringMap {
		row := history.Operation{
			Type: xdr.OperationType(typ),
		}
		op, err := NewOperation(context.Background(), row, "foo", &history.Transaction{}, history.Ledger{})
		assert.NoError(t, err, s)
		// if we got a base type, the operation is not covered
		if _, ok := op.(operations.Base); ok {
			assert.Fail(t, s)
		}
	}

	// make sure the check works for an unreasonable operation type
	row := history.Operation{
		Type: xdr.OperationType(200000),
	}
	op, err := NewOperation(context.Background(), row, "foo", &history.Transaction{}, history.Ledger{})
	assert.NoError(t, err)
	assert.IsType(t, op, operations.Base{})

}

// TestPopulateOperation_Successful tests operation object population.
func TestPopulateOperation_Successful(t *testing.T) {
	ctx, _ := test.ContextWithLogBuffer()

	var (
		dest   operations.Base
		row    history.Operation
		ledger = history.Ledger{}
	)

	dest = operations.Base{}
	row = history.Operation{TransactionSuccessful: true}

	assert.NoError(
		t,
		PopulateBaseOperation(ctx, &dest, row, "", nil, ledger),
	)
	assert.True(t, dest.TransactionSuccessful)
	assert.Nil(t, dest.Transaction)

	dest = operations.Base{}
	row = history.Operation{TransactionSuccessful: false}

	assert.NoError(
		t,
		PopulateBaseOperation(ctx, &dest, row, "", nil, ledger),
	)
	assert.False(t, dest.TransactionSuccessful)
	assert.Nil(t, dest.Transaction)
}

// TestPopulateOperation_WithTransaction tests PopulateBaseOperation when passing both an operation and a transaction.
func TestPopulateOperation_WithTransaction(t *testing.T) {
	ctx, _ := test.ContextWithLogBuffer()

	var (
		dest           operations.Base
		operationsRow  history.Operation
		ledger         = history.Ledger{}
		transactionRow history.Transaction
	)

	dest = operations.Base{}
	operationsRow = history.Operation{TransactionSuccessful: true}
	transactionRow = history.Transaction{
		TransactionWithoutLedger: history.TransactionWithoutLedger{
			Successful: true,
			MaxFee:     10000,
			FeeCharged: 100,
		},
	}

	assert.NoError(
		t,
		PopulateBaseOperation(ctx, &dest, operationsRow, transactionRow.TransactionHash, &transactionRow, ledger),
	)
	assert.True(t, dest.TransactionSuccessful)
	assert.True(t, dest.Transaction.Successful)
	assert.Equal(t, int64(100), dest.Transaction.FeeCharged)
	assert.Equal(t, int64(10000), dest.Transaction.MaxFee)
}

func TestPopulateOperation_AllowTrust(t *testing.T) {
	tt := assert.New(t)

	details := `{
		"asset_code":                        "COP",
		"asset_issuer":                      "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
		"asset_type":                        "credit_alphanum4",
		"authorize":                         false,
		"authorize_to_maintain_liabilities": true,
		"trustee":                           "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
		"trustor":                           "GDQNY3PBOJOKYZSRMK2S7LHHGWZIUISD4QORETLMXEWXBI7KFZZMKTL3"
	}`

	rsp, err := getJSONResponse(xdr.OperationTypeAllowTrust, details)
	tt.NoError(err)
	tt.Equal(false, rsp["authorize"])
	tt.Equal(true, rsp["authorize_to_maintain_liabilities"])

	details = `{
		"asset_code":                        "COP",
		"asset_issuer":                      "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
		"asset_type":                        "credit_alphanum4",
		"authorize":                         true,
		"authorize_to_maintain_liabilities": true,
		"trustee":                           "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
		"trustor":                           "GDQNY3PBOJOKYZSRMK2S7LHHGWZIUISD4QORETLMXEWXBI7KFZZMKTL3"
	}`

	rsp, err = getJSONResponse(xdr.OperationTypeAllowTrust, details)
	tt.NoError(err)
	tt.Equal(true, rsp["authorize"])
	tt.Equal(true, rsp["authorize_to_maintain_liabilities"])

	details = `{
		"asset_code":                        "COP",
		"asset_issuer":                      "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
		"asset_type":                        "credit_alphanum4",
		"authorize":                         false,
		"authorize_to_maintain_liabilities": false,
		"trustee":                           "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
		"trustor":                           "GDQNY3PBOJOKYZSRMK2S7LHHGWZIUISD4QORETLMXEWXBI7KFZZMKTL3"
	}`

	rsp, err = getJSONResponse(xdr.OperationTypeAllowTrust, details)
	tt.NoError(err)
	tt.Equal(false, rsp["authorize"])
	tt.Equal(false, rsp["authorize_to_maintain_liabilities"])
}

func TestPopulateOperation_CreateClaimableBalance(t *testing.T) {
	tt := assert.New(t)

	details := `{
		"asset":  "COP:GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
		"amount": "10.0000000",
		"claimants": [
			{
				"destination": "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
				"predicate": {
					"and": [
						{
							"or": [
								{"rel_before":"12"},
								{"abs_before": "2020-08-26T11:15:39Z"}
							]
						},
						{
							"not": {"unconditional": true}
						}
					]
				}
			}
		]
	}`

	resp, err := getJSONResponse(xdr.OperationTypeCreateClaimableBalance, details)
	tt.NoError(err)
	tt.Equal("COP:GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD", resp["asset"])
	tt.Equal("10.0000000", resp["amount"])
}

func TestPopulateOperation_ClaimClaimableBalance(t *testing.T) {
	tt := assert.New(t)

	details := `{
		"balance_id": "abc",
		"claimant": "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD"
	}`

	resp, err := getJSONResponse(xdr.OperationTypeClaimClaimableBalance, details)
	tt.NoError(err)
	tt.Equal("abc", resp["balance_id"])
	tt.Equal("GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD", resp["claimant"])
}

func TestPopulateOperation_ClaimClaimableBalance_Muxed(t *testing.T) {
	tt := assert.New(t)

	details := `{
		"claimant":                "GAQAA5L65LSYH7CQ3VTJ7F3HHLGCL3DSLAR2Y47263D56MNNGHSQSTVY",
		"claimant_muxed":          "MAQAA5L65LSYH7CQ3VTJ7F3HHLGCL3DSLAR2Y47263D56MNNGHSQSAAAAAAAAAAE2LP26",
		"claimant_muxed_id":       "1234",
		"balance_id":              "abc",
		"source_account_muxed":    "MAQAA5L65LSYH7CQ3VTJ7F3HHLGCL3DSLAR2Y47263D56MNNGHSQSAAAAAAAAAAE2LP26",
		"source_account_muxed_id": "1234"
	}`

	resp, err := getJSONResponse(xdr.OperationTypeClaimClaimableBalance, details)
	tt.NoError(err)
	tt.Equal("abc", resp["balance_id"])
	tt.Equal("GAQAA5L65LSYH7CQ3VTJ7F3HHLGCL3DSLAR2Y47263D56MNNGHSQSTVY", resp["claimant"])
	tt.Equal("MAQAA5L65LSYH7CQ3VTJ7F3HHLGCL3DSLAR2Y47263D56MNNGHSQSAAAAAAAAAAE2LP26", resp["claimant_muxed"])
	tt.Equal("1234", resp["claimant_muxed_id"])
	tt.Equal("MAQAA5L65LSYH7CQ3VTJ7F3HHLGCL3DSLAR2Y47263D56MNNGHSQSAAAAAAAAAAE2LP26", resp["source_account_muxed"])
	tt.Equal("1234", resp["source_account_muxed_id"])
}

func TestPopulateOperation_BeginSponsoringFutureReserves(t *testing.T) {
	tt := assert.New(t)

	details := `{
		"sponsored_id": "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD"
	}`

	resp, err := getJSONResponse(xdr.OperationTypeBeginSponsoringFutureReserves, details)
	tt.NoError(err)
	tt.Equal("GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD", resp["sponsored_id"])
}

func TestPopulateOperation_EndSponsoringFutureReserves(t *testing.T) {
	tt := assert.New(t)

	details := `{
		"begin_sponsor": "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD"
	}`

	resp, err := getJSONResponse(xdr.OperationTypeEndSponsoringFutureReserves, details)
	tt.NoError(err)
	tt.Equal("GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD", resp["begin_sponsor"])
}

func TestPopulateOperation_OperationTypeRevokeSponsorship_Account(t *testing.T) {
	tt := assert.New(t)

	details := `{
		"account_id": "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD"
	}`

	resp, err := getJSONResponse(xdr.OperationTypeRevokeSponsorship, details)
	tt.NoError(err)
	tt.Equal("GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD", resp["account_id"])
}

func TestPopulateOperation_OperationTypeRevokeSponsorship_Data(t *testing.T) {
	tt := assert.New(t)

	details := `{
		"data_account_id": "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
		"data_name": "name"
	}`

	resp, err := getJSONResponse(xdr.OperationTypeRevokeSponsorship, details)
	tt.NoError(err)
	tt.Equal("GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD", resp["data_account_id"])
	tt.Equal("name", resp["data_name"])
}

func TestPopulateOperation_OperationTypeRevokeSponsorship_Offer(t *testing.T) {
	tt := assert.New(t)

	details := `{
		"offer_id": "1000"
	}`

	resp, err := getJSONResponse(xdr.OperationTypeRevokeSponsorship, details)
	tt.NoError(err)
	tt.Equal("1000", resp["offer_id"])
}

func TestPopulateOperation_OperationTypeRevokeSponsorship_Trustline(t *testing.T) {
	tt := assert.New(t)

	details := `{
		"trustline_account_id": "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
		"trustline_asset": "COP:GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD"
	}`

	resp, err := getJSONResponse(xdr.OperationTypeRevokeSponsorship, details)
	tt.NoError(err)
	tt.Equal("GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD", resp["trustline_account_id"])
	tt.Equal("COP:GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD", resp["trustline_asset"])
}

func getJSONResponse(typ xdr.OperationType, details string) (rsp map[string]interface{}, err error) {
	ctx, _ := test.ContextWithLogBuffer()
	transactionRow := history.Transaction{
		TransactionWithoutLedger: history.TransactionWithoutLedger{
			Successful: true,
			MaxFee:     10000,
			FeeCharged: 100,
		},
	}
	operationsRow := history.Operation{
		TransactionSuccessful: true,
		Type:                  typ,
		DetailsString:         null.StringFrom(details),
	}
	resource, err := NewOperation(ctx, operationsRow, "", &transactionRow, history.Ledger{})
	if err != nil {
		return
	}

	data, err := json.Marshal(resource)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &rsp)
	return
}

func TestFeeBumpOperation(t *testing.T) {
	ctx, _ := test.ContextWithLogBuffer()
	dest := operations.Base{}
	operationsRow := history.Operation{TransactionSuccessful: true}
	transactionRow := history.Transaction{
		TransactionWithoutLedger: history.TransactionWithoutLedger{
			Successful:           true,
			MaxFee:               123,
			FeeCharged:           100,
			TransactionHash:      "cebb875a00ff6e1383aef0fd251a76f22c1f9ab2a2dffcb077855736ade2659a",
			FeeAccount:           null.StringFrom("GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU"),
			Account:              "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
			NewMaxFee:            null.IntFrom(10000),
			InnerTransactionHash: null.StringFrom("2374e99349b9ef7dba9a5db3339b78fda8f34777b1af33ba468ad5c0df946d4d"),
			Signatures:           []string{"a", "b", "c"},
			InnerSignatures:      []string{"d", "e", "f"},
		},
	}

	assert.NoError(
		t,
		PopulateBaseOperation(ctx, &dest, operationsRow, transactionRow.TransactionHash, nil, history.Ledger{}),
	)
	assert.Equal(t, transactionRow.TransactionHash, dest.TransactionHash)

	assert.NoError(
		t,
		PopulateBaseOperation(ctx, &dest, operationsRow, transactionRow.InnerTransactionHash.String, nil, history.Ledger{}),
	)
	assert.Equal(t, transactionRow.InnerTransactionHash.String, dest.TransactionHash)

	assert.NoError(
		t,
		PopulateBaseOperation(ctx, &dest, operationsRow, transactionRow.TransactionHash, &transactionRow, history.Ledger{}),
	)

	assert.Equal(t, transactionRow.TransactionHash, dest.TransactionHash)
	assert.Equal(t, transactionRow.TransactionHash, dest.Transaction.Hash)
	assert.Equal(t, transactionRow.TransactionHash, dest.Transaction.ID)
	assert.Equal(t, transactionRow.FeeAccount.String, dest.Transaction.FeeAccount)
	assert.Equal(t, transactionRow.Account, dest.Transaction.Account)
	assert.Equal(t, transactionRow.FeeCharged, dest.Transaction.FeeCharged)
	assert.Equal(t, transactionRow.NewMaxFee.Int64, dest.Transaction.MaxFee)
	assert.Equal(t, []string{"a", "b", "c"}, dest.Transaction.Signatures)
	assert.Equal(t, transactionRow.InnerTransactionHash.String, dest.Transaction.InnerTransaction.Hash)
	assert.Equal(t, transactionRow.MaxFee, dest.Transaction.InnerTransaction.MaxFee)
	assert.Equal(t, []string{"d", "e", "f"}, dest.Transaction.InnerTransaction.Signatures)
	assert.Equal(t, transactionRow.TransactionHash, dest.Transaction.FeeBumpTransaction.Hash)
	assert.Equal(t, []string{"a", "b", "c"}, dest.Transaction.FeeBumpTransaction.Signatures)

	assert.NoError(
		t,
		PopulateBaseOperation(ctx, &dest, operationsRow, transactionRow.InnerTransactionHash.String, &transactionRow, history.Ledger{}),
	)
	assert.Equal(t, transactionRow.InnerTransactionHash.String, dest.TransactionHash)
	assert.Equal(t, transactionRow.InnerTransactionHash.String, dest.Transaction.Hash)
	assert.Equal(t, transactionRow.InnerTransactionHash.String, dest.Transaction.ID)
	assert.Equal(t, transactionRow.FeeAccount.String, dest.Transaction.FeeAccount)
	assert.Equal(t, transactionRow.Account, dest.Transaction.Account)
	assert.Equal(t, transactionRow.FeeCharged, dest.Transaction.FeeCharged)
	assert.Equal(t, transactionRow.NewMaxFee.Int64, dest.Transaction.MaxFee)
	assert.Equal(t, []string{"d", "e", "f"}, dest.Transaction.Signatures)
	assert.Equal(t, transactionRow.InnerTransactionHash.String, dest.Transaction.InnerTransaction.Hash)
	assert.Equal(t, transactionRow.MaxFee, dest.Transaction.InnerTransaction.MaxFee)
	assert.Equal(t, []string{"d", "e", "f"}, dest.Transaction.InnerTransaction.Signatures)
	assert.Equal(t, transactionRow.TransactionHash, dest.Transaction.FeeBumpTransaction.Hash)
	assert.Equal(t, []string{"a", "b", "c"}, dest.Transaction.FeeBumpTransaction.Signatures)
}

func TestPopulateOperation_OperationTypeManageSellOffer(t *testing.T) {
	tt := assert.New(t)

	details := `{
		"offer_id": 1000
	}`

	resp, err := getJSONResponse(xdr.OperationTypeManageSellOffer, details)
	tt.NoError(err)
	tt.Equal("1000", resp["offer_id"])
}

func TestPopulateOperation_OperationTypeManageBuyOffer(t *testing.T) {
	tt := assert.New(t)

	details := `{
		"offer_id": 1000
	}`

	resp, err := getJSONResponse(xdr.OperationTypeManageBuyOffer, details)
	tt.NoError(err)
	tt.Equal("1000", resp["offer_id"])
}
