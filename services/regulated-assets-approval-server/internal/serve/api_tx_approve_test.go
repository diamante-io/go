package serve

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go/amount"
	"go/clients/auroraclient"
	"go/keypair"
	"go/network"
	"go/protocols/aurora"
	"go/services/regulated-assets-approval-server/internal/db/dbtest"
	"go/services/regulated-assets-approval-server/internal/serve/kycstatus"
	"go/txnbuild"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPI_txApprove_rejected(t *testing.T) {
	ctx := context.Background()
	db := dbtest.Open(t)
	defer db.Close()
	conn := db.Open()
	defer conn.Close()

	issuerKP := keypair.MustRandom()
	kycThresholdAmount, err := amount.ParseInt64("500")
	require.NoError(t, err)

	handler := txApproveHandler{
		issuerKP:          issuerKP,
		assetCode:         "FOO",
		auroraClient:      &auroraclient.MockClient{},
		networkPassphrase: network.TestNetworkPassphrase,
		db:                conn,
		kycThreshold:      kycThresholdAmount,
		baseURL:           "https://example.com",
	}

	// rejected if no transaction "tx" is submitted
	m := chi.NewMux()
	m.Post("/tx-approve", handler.ServeHTTP)
	r := httptest.NewRequest("POST", "/tx-approve", nil)
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
		"status": "rejected",
		"error": "Missing parameter \"tx\"."
	}`
	require.JSONEq(t, wantBody, string(body))
}

func TestAPI_txApprove_revised(t *testing.T) {
	ctx := context.Background()
	db := dbtest.Open(t)
	defer db.Close()
	conn := db.Open()
	defer conn.Close()

	senderKP := keypair.MustRandom()
	receiverKP := keypair.MustRandom()
	issuerKP := keypair.MustRandom()
	assetGOAT := txnbuild.CreditAsset{
		Code:   "GOAT",
		Issuer: issuerKP.Address(),
	}
	kycThresholdAmount, err := amount.ParseInt64("500")
	require.NoError(t, err)

	auroraMock := auroraclient.MockClient{}
	auroraMock.
		On("AccountDetail", auroraclient.AccountRequest{AccountID: senderKP.Address()}).
		Return(aurora.Account{
			AccountID: senderKP.Address(),
			Sequence:  "5",
		}, nil)

	handler := txApproveHandler{
		issuerKP:          issuerKP,
		assetCode:         assetGOAT.GetCode(),
		auroraClient:      &auroraMock,
		networkPassphrase: network.TestNetworkPassphrase,
		db:                conn,
		kycThreshold:      kycThresholdAmount,
		baseURL:           "https://example.com",
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount: &aurora.Account{
				AccountID: senderKP.Address(),
				Sequence:  "5",
			},
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.Payment{
					SourceAccount: senderKP.Address(),
					Destination:   receiverKP.Address(),
					Amount:        "1",
					Asset:         assetGOAT,
				},
			},
			BaseFee:    txnbuild.MinBaseFee,
			Timebounds: txnbuild.NewInfiniteTimeout(),
		},
	)
	require.NoError(t, err)
	txe, err := tx.Base64()
	require.NoError(t, err)

	m := chi.NewMux()
	m.Post("/tx-approve", handler.ServeHTTP)
	r := httptest.NewRequest("POST", "/tx-approve", strings.NewReader(`{"tx": "`+txe+`"}`))
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var gotResponse txApprovalResponse
	err = json.Unmarshal(body, &gotResponse)
	require.NoError(t, err)
	require.Equal(t, sep8StatusRevised, gotResponse.Status)
	require.Equal(t, "Authorization and deauthorization operations were added.", gotResponse.Message)

	gotGenericTx, err := txnbuild.TransactionFromXDR(gotResponse.Tx)
	require.NoError(t, err)
	gotTx, ok := gotGenericTx.Transaction()
	require.True(t, ok)

	require.Len(t, gotTx.Operations(), 5)
	// AllowTrust op where issuer fully authorizes sender, asset GOAT
	op0, ok := gotTx.Operations()[0].(*txnbuild.AllowTrust)
	require.True(t, ok)
	assert.Equal(t, op0.Trustor, senderKP.Address())
	assert.Equal(t, op0.Type.GetCode(), assetGOAT.GetCode())
	require.True(t, op0.Authorize)
	// AllowTrust op where issuer fully authorizes receiver, asset GOAT
	op1, ok := gotTx.Operations()[1].(*txnbuild.AllowTrust)
	require.True(t, ok)
	assert.Equal(t, op1.Trustor, receiverKP.Address())
	assert.Equal(t, op1.Type.GetCode(), assetGOAT.GetCode())
	require.True(t, op1.Authorize)
	// Payment from sender to receiver
	op2, ok := gotTx.Operations()[2].(*txnbuild.Payment)
	require.True(t, ok)
	assert.Equal(t, op2.Destination, receiverKP.Address())
	assert.Equal(t, op2.Asset, assetGOAT)
	// AllowTrust op where issuer fully deauthorizes receiver, asset GOAT
	op3, ok := gotTx.Operations()[3].(*txnbuild.AllowTrust)
	require.True(t, ok)
	assert.Equal(t, op3.Trustor, receiverKP.Address())
	assert.Equal(t, op3.Type.GetCode(), assetGOAT.GetCode())
	require.False(t, op3.Authorize)
	// AllowTrust op where issuer fully deauthorizes sender, asset GOAT
	op4, ok := gotTx.Operations()[4].(*txnbuild.AllowTrust)
	require.True(t, ok)
	assert.Equal(t, op4.Trustor, senderKP.Address())
	assert.Equal(t, op4.Type.GetCode(), assetGOAT.GetCode())
	require.False(t, op4.Authorize)
}

func TestAPI_txAprove_actionRequired(t *testing.T) {
	ctx := context.Background()
	db := dbtest.Open(t)
	defer db.Close()
	conn := db.Open()
	defer conn.Close()

	// prepare handler dependencies
	senderKP := keypair.MustRandom()
	receiverKP := keypair.MustRandom()
	issuerKP := keypair.MustRandom()
	assetGOAT := txnbuild.CreditAsset{
		Code:   "GOAT",
		Issuer: issuerKP.Address(),
	}
	kycThresholdAmount, err := amount.ParseInt64("500")
	require.NoError(t, err)

	auroraMock := auroraclient.MockClient{}
	auroraMock.
		On("AccountDetail", auroraclient.AccountRequest{AccountID: senderKP.Address()}).
		Return(aurora.Account{
			AccountID: senderKP.Address(),
			Sequence:  "1",
		}, nil)

	handler := txApproveHandler{
		issuerKP:          issuerKP,
		assetCode:         assetGOAT.GetCode(),
		auroraClient:      &auroraMock,
		networkPassphrase: network.TestNetworkPassphrase,
		db:                conn,
		kycThreshold:      kycThresholdAmount,
		baseURL:           "https://example.com",
	}

	// setup route handlers
	m := chi.NewMux()
	m.Post("/tx-approve", handler.ServeHTTP)
	m.Post("/kyc-status/{callback_id}", kycstatus.PostHandler{DB: conn}.ServeHTTP)

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount: &aurora.Account{
				AccountID: senderKP.Address(),
				Sequence:  "1",
			},
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.Payment{
					SourceAccount: senderKP.Address(),
					Destination:   receiverKP.Address(),
					Amount:        "501",
					Asset:         assetGOAT,
				},
			},
			BaseFee:    txnbuild.MinBaseFee,
			Timebounds: txnbuild.NewInfiniteTimeout(),
		},
	)
	require.NoError(t, err)
	txe, err := tx.Base64()
	require.NoError(t, err)

	r := httptest.NewRequest("POST", "/tx-approve", strings.NewReader(`{"tx": "`+txe+`"}`))
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, r)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var callbackID string
	q := `SELECT callback_id FROM accounts_kyc_status WHERE diamcircle_address = $1`
	err = conn.QueryRowContext(ctx, q, senderKP.Address()).Scan(&callbackID)
	require.NoError(t, err)

	var gotTxApprovalResponse txApprovalResponse
	err = json.Unmarshal(body, &gotTxApprovalResponse)
	require.NoError(t, err)
	wantTxApprovalResponse := txApprovalResponse{
		Status:       sep8Status("action_required"),
		Message:      "Payments exceeding 500.00 GOAT require KYC approval. Please provide an email address.",
		ActionURL:    "https://example.com/kyc-status/" + callbackID,
		ActionMethod: "POST",
		ActionFields: []string{"email_address"},
	}
	assert.Equal(t, wantTxApprovalResponse, gotTxApprovalResponse)
}

func TestAPI_txAprove_actionRequiredFlow(t *testing.T) {
	ctx := context.Background()
	db := dbtest.Open(t)
	defer db.Close()
	conn := db.Open()
	defer conn.Close()

	// prepare handler dependencies
	senderKP := keypair.MustRandom()
	receiverKP := keypair.MustRandom()
	issuerKP := keypair.MustRandom()
	assetGOAT := txnbuild.CreditAsset{
		Code:   "GOAT",
		Issuer: issuerKP.Address(),
	}
	kycThresholdAmount, err := amount.ParseInt64("500")
	require.NoError(t, err)

	auroraMock := auroraclient.MockClient{}
	auroraMock.
		On("AccountDetail", auroraclient.AccountRequest{AccountID: senderKP.Address()}).
		Return(aurora.Account{
			AccountID: senderKP.Address(),
			Sequence:  "1",
		}, nil)

	handler := txApproveHandler{
		issuerKP:          issuerKP,
		assetCode:         assetGOAT.GetCode(),
		auroraClient:      &auroraMock,
		networkPassphrase: network.TestNetworkPassphrase,
		db:                conn,
		kycThreshold:      kycThresholdAmount,
		baseURL:           "https://example.com",
	}

	// setup route handlers
	m := chi.NewMux()
	m.Post("/tx-approve", handler.ServeHTTP)
	m.Post("/kyc-status/{callback_id}", kycstatus.PostHandler{DB: conn}.ServeHTTP)

	// Step 1: client sends payment with 500+ GOAT
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount: &aurora.Account{
				AccountID: senderKP.Address(),
				Sequence:  "1",
			},
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.Payment{
					SourceAccount: senderKP.Address(),
					Destination:   receiverKP.Address(),
					Amount:        "501",
					Asset:         assetGOAT,
				},
			},
			BaseFee:    txnbuild.MinBaseFee,
			Timebounds: txnbuild.NewInfiniteTimeout(),
		},
	)
	require.NoError(t, err)
	txe, err := tx.Base64()
	require.NoError(t, err)

	r := httptest.NewRequest("POST", "/tx-approve", strings.NewReader(`{"tx": "`+txe+`"}`))
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, r)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var callbackID string
	q := `SELECT callback_id FROM accounts_kyc_status WHERE diamcircle_address = $1`
	err = conn.QueryRowContext(ctx, q, senderKP.Address()).Scan(&callbackID)
	require.NoError(t, err)

	var gotTxApprovalResponse txApprovalResponse
	err = json.Unmarshal(body, &gotTxApprovalResponse)
	require.NoError(t, err)
	wantTxApprovalResponse := txApprovalResponse{
		Status:       sep8Status("action_required"),
		Message:      "Payments exceeding 500.00 GOAT require KYC approval. Please provide an email address.",
		ActionURL:    "https://example.com/kyc-status/" + callbackID,
		ActionMethod: "POST",
		ActionFields: []string{"email_address"},
	}
	assert.Equal(t, wantTxApprovalResponse, gotTxApprovalResponse)

	// Step 2: client follows up with action required. KYC should get approved for emails not starting with "x" nor "y"
	actionMethod := gotTxApprovalResponse.ActionMethod
	actionURL := gotTxApprovalResponse.ActionURL
	actionFields := strings.NewReader(`{"email_address": "test@email.com"}`)
	r = httptest.NewRequest(actionMethod, actionURL, actionFields)
	r = r.WithContext(ctx)
	w = httptest.NewRecorder()
	m.ServeHTTP(w, r)
	resp = w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err = ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	wantBody := `{"result": "no_further_action_required"}`
	require.JSONEq(t, wantBody, string(body))

	// Step 3: verify transactions with 500+ GOAT can now be revised
	r = httptest.NewRequest("POST", "/tx-approve", strings.NewReader(`{"tx": "`+txe+`"}`))
	r = r.WithContext(ctx)
	w = httptest.NewRecorder()
	m.ServeHTTP(w, r)
	resp = w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
	body, err = ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	gotTxApprovalResponse = txApprovalResponse{}
	err = json.Unmarshal(body, &gotTxApprovalResponse)
	require.NoError(t, err)
	assert.Equal(t, sep8StatusRevised, gotTxApprovalResponse.Status)
	assert.Equal(t, "Authorization and deauthorization operations were added.", gotTxApprovalResponse.Message)
	require.NotEmpty(t, gotTxApprovalResponse.Tx)

	// Step 4: client follows up with action required again. This time KYC will get rejected as the email starts with "x"
	actionFields = strings.NewReader(`{"email_address": "xtest@email.com"}`)
	r = httptest.NewRequest(actionMethod, actionURL, actionFields)
	r = r.WithContext(ctx)
	w = httptest.NewRecorder()
	m.ServeHTTP(w, r)
	resp = w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err = ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	wantBody = `{"result": "no_further_action_required"}`
	require.JSONEq(t, wantBody, string(body))

	// Step 5: verify transactions with 500+ GOAT are rejected
	r = httptest.NewRequest("POST", "/tx-approve", strings.NewReader(`{"tx": "`+txe+`"}`))
	r = r.WithContext(ctx)
	w = httptest.NewRecorder()
	m.ServeHTTP(w, r)
	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err = ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	wantBody = `{
		"status": "rejected",
		"error": "Your KYC was rejected and you're not authorized for operations above 500.00 GOAT."
	}`
	require.JSONEq(t, wantBody, string(body))

	// Step 6: client follows up with action required again. This time KYC will be marked as pending as the email starts with "y"
	actionFields = strings.NewReader(`{"email_address": "ytest@email.com"}`)
	r = httptest.NewRequest(actionMethod, actionURL, actionFields)
	r = r.WithContext(ctx)
	w = httptest.NewRecorder()
	m.ServeHTTP(w, r)
	resp = w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err = ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	wantBody = `{"result": "no_further_action_required"}`
	require.JSONEq(t, wantBody, string(body))

	// Step 7: verify transactions with 500+ GOAT are pending
	r = httptest.NewRequest("POST", "/tx-approve", strings.NewReader(`{"tx": "`+txe+`"}`))
	r = r.WithContext(ctx)
	w = httptest.NewRecorder()
	m.ServeHTTP(w, r)
	resp = w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err = ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	wantBody = `{
		"status": "pending",
		"message": "Your account could not be verified as approved nor rejected and was marked as pending. You will need staff authorization for operations above 500.00 GOAT.",
		"timeout": 0
	}`
	require.JSONEq(t, wantBody, string(body))
}

func TestAPI_txApprove_success(t *testing.T) {
	ctx := context.Background()
	db := dbtest.Open(t)
	defer db.Close()
	conn := db.Open()
	defer conn.Close()

	senderKP := keypair.MustRandom()
	receiverKP := keypair.MustRandom()
	issuerKP := keypair.MustRandom()
	assetGOAT := txnbuild.CreditAsset{
		Code:   "GOAT",
		Issuer: issuerKP.Address(),
	}
	kycThresholdAmount, err := amount.ParseInt64("500")
	require.NoError(t, err)

	auroraMock := auroraclient.MockClient{}
	auroraMock.
		On("AccountDetail", auroraclient.AccountRequest{AccountID: senderKP.Address()}).
		Return(aurora.Account{
			AccountID: senderKP.Address(),
			Sequence:  "5",
		}, nil)

	handler := txApproveHandler{
		issuerKP:          issuerKP,
		assetCode:         assetGOAT.GetCode(),
		auroraClient:      &auroraMock,
		networkPassphrase: network.TestNetworkPassphrase,
		db:                conn,
		kycThreshold:      kycThresholdAmount,
		baseURL:           "https://example.com",
	}
	m := chi.NewMux()
	m.Post("/tx-approve", handler.ServeHTTP)

	// prepare SEP-8 compliant transaction
	tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount: &aurora.Account{
			AccountID: senderKP.Address(),
			Sequence:  "5",
		},
		IncrementSequenceNum: true,
		Operations: []txnbuild.Operation{
			&txnbuild.AllowTrust{
				Trustor:       senderKP.Address(),
				Type:          assetGOAT,
				Authorize:     true,
				SourceAccount: issuerKP.Address(),
			},
			&txnbuild.AllowTrust{
				Trustor:       receiverKP.Address(),
				Type:          assetGOAT,
				Authorize:     true,
				SourceAccount: issuerKP.Address(),
			},
			&txnbuild.Payment{
				Destination: receiverKP.Address(),
				Amount:      "1",
				Asset:       assetGOAT,
			},
			&txnbuild.AllowTrust{
				Trustor:       receiverKP.Address(),
				Type:          assetGOAT,
				Authorize:     false,
				SourceAccount: issuerKP.Address(),
			},
			&txnbuild.AllowTrust{
				Trustor:       senderKP.Address(),
				Type:          assetGOAT,
				Authorize:     false,
				SourceAccount: issuerKP.Address(),
			},
		},
		BaseFee:    300,
		Timebounds: txnbuild.NewTimeout(300),
	})
	require.NoError(t, err)
	txe, err := tx.Base64()
	require.NoError(t, err)

	r := httptest.NewRequest("POST", "/tx-approve", strings.NewReader(`{"tx": "`+txe+`"}`))
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var gotSuccessResponse txApprovalResponse
	err = json.Unmarshal(body, &gotSuccessResponse)
	require.NoError(t, err)
	wantSuccessResponse := txApprovalResponse{
		Status:  sep8Status("success"),
		Tx:      gotSuccessResponse.Tx,
		Message: "Transaction is compliant and signed by the issuer.",
	}
	assert.Equal(t, wantSuccessResponse, gotSuccessResponse)

	genericTx, err := txnbuild.TransactionFromXDR(gotSuccessResponse.Tx)
	require.NoError(t, err)
	tx, ok := genericTx.Transaction()
	require.True(t, ok)
	require.Equal(t, senderKP.Address(), tx.SourceAccount().AccountID)

	require.Len(t, tx.Operations(), 5)
	// AllowTrust op where issuer fully authorizes sender, asset GOAT
	op0, ok := tx.Operations()[0].(*txnbuild.AllowTrust)
	require.True(t, ok)
	assert.Equal(t, op0.Trustor, senderKP.Address())
	assert.Equal(t, op0.Type.GetCode(), assetGOAT.GetCode())
	require.True(t, op0.Authorize)
	// AllowTrust op where issuer fully authorizes receiver, asset GOAT
	op1, ok := tx.Operations()[1].(*txnbuild.AllowTrust)
	require.True(t, ok)
	assert.Equal(t, op1.Trustor, receiverKP.Address())
	assert.Equal(t, op1.Type.GetCode(), assetGOAT.GetCode())
	require.True(t, op1.Authorize)
	// Payment from sender to receiver
	op2, ok := tx.Operations()[2].(*txnbuild.Payment)
	require.True(t, ok)
	assert.Equal(t, op2.Destination, receiverKP.Address())
	assert.Equal(t, op2.Asset, assetGOAT)
	// AllowTrust op where issuer fully deauthorizes receiver, asset GOAT
	op3, ok := tx.Operations()[3].(*txnbuild.AllowTrust)
	require.True(t, ok)
	assert.Equal(t, op3.Trustor, receiverKP.Address())
	assert.Equal(t, op3.Type.GetCode(), assetGOAT.GetCode())
	require.False(t, op3.Authorize)
	// AllowTrust op where issuer fully deauthorizes sender, asset GOAT
	op4, ok := tx.Operations()[4].(*txnbuild.AllowTrust)
	require.True(t, ok)
	assert.Equal(t, op4.Trustor, senderKP.Address())
	assert.Equal(t, op4.Type.GetCode(), assetGOAT.GetCode())
	require.False(t, op4.Authorize)

	// check if the transaction contains the issuer's signature
	txHash, err := tx.Hash(handler.networkPassphrase)
	require.NoError(t, err)
	err = handler.issuerKP.Verify(txHash[:], tx.Signatures()[0].Signature)
	require.NoError(t, err)
}
