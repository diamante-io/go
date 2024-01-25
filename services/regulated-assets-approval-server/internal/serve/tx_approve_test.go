package serve

import (
	"context"
	"net/http"
	"testing"

	"github.com/diamcircle/go/amount"
	"github.com/diamcircle/go/clients/auroraclient"
	"github.com/diamcircle/go/keypair"
	"github.com/diamcircle/go/network"
	"github.com/diamcircle/go/protocols/aurora"
	"github.com/diamcircle/go/services/regulated-assets-approval-server/internal/db/dbtest"
	"github.com/diamcircle/go/txnbuild"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTxApproveHandlerValidate(t *testing.T) {
	// empty issuer KP.
	h := txApproveHandler{}
	err := h.validate()
	require.EqualError(t, err, "issuer keypair cannot be nil")

	// empty asset code.
	issuerAccKeyPair := keypair.MustRandom()
	h = txApproveHandler{
		issuerKP: issuerAccKeyPair,
	}
	err = h.validate()
	require.EqualError(t, err, "asset code cannot be empty")

	// No Aurora client.
	h = txApproveHandler{
		issuerKP:  issuerAccKeyPair,
		assetCode: "FOOBAR",
	}
	err = h.validate()
	require.EqualError(t, err, "aurora client cannot be nil")

	// No network passphrase.
	auroraMock := auroraclient.MockClient{}
	h = txApproveHandler{
		issuerKP:      issuerAccKeyPair,
		assetCode:     "FOOBAR",
		auroraClient: &auroraMock,
	}
	err = h.validate()
	require.EqualError(t, err, "network passphrase cannot be empty")

	// No db.
	h = txApproveHandler{
		issuerKP:          issuerAccKeyPair,
		assetCode:         "FOOBAR",
		auroraClient:     &auroraMock,
		networkPassphrase: network.TestNetworkPassphrase,
	}
	err = h.validate()
	require.EqualError(t, err, "database cannot be nil")

	// Empty kycThreshold.
	db := dbtest.Open(t)
	defer db.Close()
	conn := db.Open()
	defer conn.Close()
	h = txApproveHandler{
		issuerKP:          issuerAccKeyPair,
		assetCode:         "FOOBAR",
		auroraClient:     &auroraMock,
		networkPassphrase: network.TestNetworkPassphrase,
		db:                conn,
	}
	err = h.validate()
	require.EqualError(t, err, "kyc threshold cannot be less than or equal to zero")

	// Negative kycThreshold.
	h = txApproveHandler{
		issuerKP:          issuerAccKeyPair,
		assetCode:         "FOOBAR",
		auroraClient:     &auroraMock,
		networkPassphrase: network.TestNetworkPassphrase,
		db:                conn,
		kycThreshold:      -1,
	}
	err = h.validate()
	require.EqualError(t, err, "kyc threshold cannot be less than or equal to zero")

	// no baseURL.
	h = txApproveHandler{
		issuerKP:          issuerAccKeyPair,
		assetCode:         "FOOBAR",
		auroraClient:     &auroraMock,
		networkPassphrase: network.TestNetworkPassphrase,
		db:                conn,
		kycThreshold:      1,
	}
	err = h.validate()
	require.EqualError(t, err, "base url cannot be empty")

	// Success.
	h = txApproveHandler{
		issuerKP:          issuerAccKeyPair,
		assetCode:         "FOOBAR",
		auroraClient:     &auroraMock,
		networkPassphrase: network.TestNetworkPassphrase,
		db:                conn,
		kycThreshold:      1,
		baseURL:           "https://example.com",
	}
	err = h.validate()
	require.NoError(t, err)
}

func TestTxApproveHandler_validateInput(t *testing.T) {
	h := txApproveHandler{}
	ctx := context.Background()

	// rejects if incoming tx is empty
	in := txApproveRequest{}
	txApprovalResp, gotTx := h.validateInput(ctx, in)
	require.Equal(t, NewRejectedTxApprovalResponse("Missing parameter \"tx\"."), txApprovalResp)
	require.Nil(t, gotTx)

	// rejects if incoming tx is invalid
	in = txApproveRequest{Tx: "foobar"}
	txApprovalResp, gotTx = h.validateInput(ctx, in)
	require.Equal(t, NewRejectedTxApprovalResponse("Invalid parameter \"tx\"."), txApprovalResp)
	require.Nil(t, gotTx)

	// rejects if incoming tx is a fee bump transaction
	in = txApproveRequest{Tx: "AAAABQAAAAAo/cVyQxyGh7F/Vsj0BzfDYuOJvrwgfHGyqYFpHB5RCAAAAAAAAADIAAAAAgAAAAAo/cVyQxyGh7F/Vsj0BzfDYuOJvrwgfHGyqYFpHB5RCAAAAGQAEfDJAAAAAQAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAQAAAAAo/cVyQxyGh7F/Vsj0BzfDYuOJvrwgfHGyqYFpHB5RCAAAAAAAAAAAAJiWgAAAAAAAAAAAAAAAAAAAAAA="}
	txApprovalResp, gotTx = h.validateInput(ctx, in)
	require.Equal(t, NewRejectedTxApprovalResponse("Invalid parameter \"tx\"."), txApprovalResp)
	require.Nil(t, gotTx)

	// rejects if tx source account is the issuer
	clientKP := keypair.MustRandom()
	h.issuerKP = keypair.MustRandom()

	tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount: &aurora.Account{
			AccountID: h.issuerKP.Address(),
			Sequence:  "1",
		},
		IncrementSequenceNum: true,
		Timebounds:           txnbuild.NewInfiniteTimeout(),
		BaseFee:              300,
		Operations: []txnbuild.Operation{
			&txnbuild.Payment{
				Destination: clientKP.Address(),
				Amount:      "1",
				Asset:       txnbuild.NativeAsset{},
			},
		},
	})
	require.NoError(t, err)
	txe, err := tx.Base64()
	require.NoError(t, err)

	in.Tx = txe
	txApprovalResp, gotTx = h.validateInput(ctx, in)
	require.Equal(t, NewRejectedTxApprovalResponse("Transaction source account is invalid."), txApprovalResp)
	require.Nil(t, gotTx)

	// rejects if there are any operations other than Allowtrust where the source account is the issuer
	tx, err = txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount: &aurora.Account{
			AccountID: clientKP.Address(),
			Sequence:  "1",
		},
		IncrementSequenceNum: true,
		Timebounds:           txnbuild.NewInfiniteTimeout(),
		BaseFee:              300,
		Operations: []txnbuild.Operation{
			&txnbuild.BumpSequence{},
			&txnbuild.Payment{
				Destination:   clientKP.Address(),
				Amount:        "1.0000000",
				Asset:         txnbuild.NativeAsset{},
				SourceAccount: h.issuerKP.Address(),
			},
		},
	})
	require.NoError(t, err)
	txe, err = tx.Base64()
	require.NoError(t, err)

	in.Tx = txe
	txApprovalResp, gotTx = h.validateInput(ctx, in)
	require.Equal(t, NewRejectedTxApprovalResponse("There are one or more unauthorized operations in the provided transaction."), txApprovalResp)
	require.Nil(t, gotTx)

	// validation success
	tx, err = txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount: &aurora.Account{
			AccountID: clientKP.Address(),
			Sequence:  "1",
		},
		IncrementSequenceNum: true,
		Timebounds:           txnbuild.NewInfiniteTimeout(),
		BaseFee:              300,
		Operations: []txnbuild.Operation{
			&txnbuild.Payment{
				Destination: clientKP.Address(),
				Amount:      "1.0000000",
				Asset:       txnbuild.NativeAsset{},
			},
		},
	})
	require.NoError(t, err)
	txe, err = tx.Base64()
	require.NoError(t, err)

	in.Tx = txe
	txApprovalResp, gotTx = h.validateInput(ctx, in)
	require.Nil(t, txApprovalResp)
	require.Equal(t, gotTx, tx)
}

func TestTxApproveHandler_handleActionRequiredResponseIfNeeded(t *testing.T) {
	ctx := context.Background()
	db := dbtest.Open(t)
	defer db.Close()
	conn := db.Open()
	defer conn.Close()

	kycThreshold, err := amount.ParseInt64("500")
	require.NoError(t, err)
	h := txApproveHandler{
		assetCode:    "FOO",
		baseURL:      "https://example.com",
		kycThreshold: kycThreshold,
		db:           conn,
	}

	// payments up to the the threshold won't trigger "action_required"
	clientKP := keypair.MustRandom()
	paymentOp := &txnbuild.Payment{
		Amount: amount.StringFromInt64(kycThreshold),
	}
	txApprovalResp, err := h.handleActionRequiredResponseIfNeeded(ctx, clientKP.Address(), paymentOp)
	require.NoError(t, err)
	require.Nil(t, txApprovalResp)

	// payments greater than the threshold will trigger "action_required"
	paymentOp = &txnbuild.Payment{
		Amount: amount.StringFromInt64(kycThreshold + 1),
	}
	txApprovalResp, err = h.handleActionRequiredResponseIfNeeded(ctx, clientKP.Address(), paymentOp)
	require.NoError(t, err)

	var callbackID string
	q := `SELECT callback_id FROM accounts_kyc_status WHERE diamcircle_address = $1`
	err = conn.QueryRowContext(ctx, q, clientKP.Address()).Scan(&callbackID)
	require.NoError(t, err)

	wantResp := &txApprovalResponse{
		Status:       sep8StatusActionRequired,
		Message:      "Payments exceeding 500.00 FOO require KYC approval. Please provide an email address.",
		ActionMethod: "POST",
		StatusCode:   http.StatusOK,
		ActionURL:    "https://example.com/kyc-status/" + callbackID,
		ActionFields: []string{"email_address"},
	}
	require.Equal(t, wantResp, txApprovalResp)

	// if KYC was previously approved, handleActionRequiredResponseIfNeeded will return nil
	q = `
		UPDATE accounts_kyc_status
		SET 
			approved_at = NOW(),
			rejected_at = NULL,
			pending_at = NULL
		WHERE diamcircle_address = $1
	`
	_, err = conn.ExecContext(ctx, q, clientKP.Address())
	require.NoError(t, err)
	txApprovalResp, err = h.handleActionRequiredResponseIfNeeded(ctx, clientKP.Address(), paymentOp)
	require.NoError(t, err)
	require.Nil(t, txApprovalResp)

	// if KYC was previously rejected, handleActionRequiredResponseIfNeeded will return a "rejected" response
	q = `
		UPDATE accounts_kyc_status
		SET 
			approved_at = NULL,
			rejected_at = NOW(),
			pending_at = NULL
		WHERE diamcircle_address = $1
	`
	_, err = conn.ExecContext(ctx, q, clientKP.Address())
	require.NoError(t, err)
	txApprovalResp, err = h.handleActionRequiredResponseIfNeeded(ctx, clientKP.Address(), paymentOp)
	require.NoError(t, err)
	require.Equal(t, NewRejectedTxApprovalResponse("Your KYC was rejected and you're not authorized for operations above 500.00 FOO."), txApprovalResp)

	// if KYC was previously marked as pending, handleActionRequiredResponseIfNeeded will return a "pending" response
	q = `
		UPDATE accounts_kyc_status
		SET 
			approved_at = NULL,
			rejected_at = NULL,
			pending_at = NOW()
		WHERE diamcircle_address = $1
	`
	_, err = conn.ExecContext(ctx, q, clientKP.Address())
	require.NoError(t, err)
	txApprovalResp, err = h.handleActionRequiredResponseIfNeeded(ctx, clientKP.Address(), paymentOp)
	require.NoError(t, err)
	require.Equal(t, NewPendingTxApprovalResponse("Your account could not be verified as approved nor rejected and was marked as pending. You will need staff authorization for operations above 500.00 FOO."), txApprovalResp)
}

func TestTxApproveHandler_txApprove_rejected(t *testing.T) {
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
			Sequence:  "2",
		}, nil)

	handler := txApproveHandler{
		issuerKP:          issuerKP,
		assetCode:         assetGOAT.GetCode(),
		auroraClient:     &auroraMock,
		networkPassphrase: network.TestNetworkPassphrase,
		db:                conn,
		kycThreshold:      kycThresholdAmount,
		baseURL:           "https://example.com",
	}

	// "rejected" if tx is empty
	rejectedResponse, err := handler.txApprove(ctx, txApproveRequest{})
	require.NoError(t, err)
	wantRejectedResponse := txApprovalResponse{
		Status:     "rejected",
		Error:      `Missing parameter "tx".`,
		StatusCode: http.StatusBadRequest,
	}
	assert.Equal(t, &wantRejectedResponse, rejectedResponse)

	// rejected if contains more than one operation
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount: &aurora.Account{
				AccountID: senderKP.Address(),
				Sequence:  "2",
			},
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.BumpSequence{},
				&txnbuild.Payment{
					Destination: receiverKP.Address(),
					Amount:      "1",
					Asset:       assetGOAT,
				},
			},
			BaseFee:    txnbuild.MinBaseFee,
			Timebounds: txnbuild.NewInfiniteTimeout(),
		},
	)
	require.NoError(t, err)
	txe, err := tx.Base64()
	require.NoError(t, err)

	txApprovalResp, err := handler.txApprove(ctx, txApproveRequest{Tx: txe})
	require.NoError(t, err)
	assert.Equal(t, NewRejectedTxApprovalResponse("Please submit a transaction with exactly one operation of type payment."), txApprovalResp)

	// rejected if the single operation is not a payment
	tx, err = txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount: &aurora.Account{
				AccountID: senderKP.Address(),
				Sequence:  "2",
			},
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.BumpSequence{},
			},
			BaseFee:    txnbuild.MinBaseFee,
			Timebounds: txnbuild.NewInfiniteTimeout(),
		},
	)
	require.NoError(t, err)
	txe, err = tx.Base64()
	require.NoError(t, err)

	txApprovalResp, err = handler.txApprove(ctx, txApproveRequest{Tx: txe})
	require.NoError(t, err)
	assert.Equal(t, NewRejectedTxApprovalResponse("There is one or more unauthorized operations in the provided transaction."), txApprovalResp)

	// rejected if attempting to transfer an asset to its own issuer
	tx, err = txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount: &aurora.Account{
				AccountID: senderKP.Address(),
				Sequence:  "2",
			},
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.Payment{
					Destination: issuerKP.Address(), // <--- this will trigger the rejection
					Amount:      "1",
					Asset: txnbuild.CreditAsset{
						Code:   "FOO",
						Issuer: keypair.MustRandom().Address(),
					},
				},
			},
			BaseFee:    txnbuild.MinBaseFee,
			Timebounds: txnbuild.NewInfiniteTimeout(),
		},
	)
	require.NoError(t, err)
	txe, err = tx.Base64()
	require.NoError(t, err)

	txApprovalResp, err = handler.txApprove(ctx, txApproveRequest{Tx: txe})
	require.NoError(t, err)
	assert.Equal(t, NewRejectedTxApprovalResponse("Can't transfer asset to its issuer."), txApprovalResp)

	// rejected if payment asset is not supported
	tx, err = txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount: &aurora.Account{
				AccountID: senderKP.Address(),
				Sequence:  "2",
			},
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.Payment{
					Destination: receiverKP.Address(),
					Amount:      "1",
					Asset: txnbuild.CreditAsset{
						Code:   "FOO",
						Issuer: keypair.MustRandom().Address(),
					},
				},
			},
			BaseFee:    txnbuild.MinBaseFee,
			Timebounds: txnbuild.NewInfiniteTimeout(),
		},
	)
	require.NoError(t, err)
	txe, err = tx.Base64()
	require.NoError(t, err)

	txApprovalResp, err = handler.txApprove(ctx, txApproveRequest{Tx: txe})
	require.NoError(t, err)
	assert.Equal(t, NewRejectedTxApprovalResponse("The payment asset is not supported by this issuer."), txApprovalResp)

	// rejected if sequence number is not incremental
	tx, err = txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount: &aurora.Account{
				AccountID: senderKP.Address(),
				Sequence:  "20",
			},
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.Payment{
					Destination: receiverKP.Address(),
					Amount:      "1",
					Asset:       assetGOAT,
				},
			},
			BaseFee:    txnbuild.MinBaseFee,
			Timebounds: txnbuild.NewInfiniteTimeout(),
		},
	)
	require.NoError(t, err)
	txe, err = tx.Base64()
	require.NoError(t, err)

	txApprovalResp, err = handler.txApprove(ctx, txApproveRequest{Tx: txe})
	require.NoError(t, err)
	assert.Equal(t, NewRejectedTxApprovalResponse("Invalid transaction sequence number."), txApprovalResp)
}

func TestTxApproveHandler_txApprove_success(t *testing.T) {
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
			Sequence:  "2",
		}, nil)

	handler := txApproveHandler{
		issuerKP:          issuerKP,
		assetCode:         assetGOAT.GetCode(),
		auroraClient:     &auroraMock,
		networkPassphrase: network.TestNetworkPassphrase,
		db:                conn,
		kycThreshold:      kycThresholdAmount,
		baseURL:           "https://example.com",
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount: &aurora.Account{
				AccountID: senderKP.Address(),
				Sequence:  "2",
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
					SourceAccount: senderKP.Address(),
					Destination:   receiverKP.Address(),
					Amount:        "1",
					Asset:         assetGOAT,
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
			BaseFee:    txnbuild.MinBaseFee,
			Timebounds: txnbuild.NewInfiniteTimeout(),
		},
	)
	require.NoError(t, err)
	txe, err := tx.Base64()
	require.NoError(t, err)

	txApprovalResp, err := handler.txApprove(ctx, txApproveRequest{Tx: txe})
	require.NoError(t, err)
	require.Equal(t, NewSuccessTxApprovalResponse(txApprovalResp.Tx, "Transaction is compliant and signed by the issuer."), txApprovalResp)
}

func TestTxApproveHandler_txApprove_actionRequired(t *testing.T) {
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
			Sequence:  "2",
		}, nil)

	handler := txApproveHandler{
		issuerKP:          issuerKP,
		assetCode:         assetGOAT.GetCode(),
		auroraClient:     &auroraMock,
		networkPassphrase: network.TestNetworkPassphrase,
		db:                conn,
		kycThreshold:      kycThresholdAmount,
		baseURL:           "https://example.com",
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount: &aurora.Account{
				AccountID: senderKP.Address(),
				Sequence:  "2",
			},
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.Payment{
					Destination: receiverKP.Address(),
					Amount:      "501",
					Asset:       assetGOAT,
				},
			},
			BaseFee:    txnbuild.MinBaseFee,
			Timebounds: txnbuild.NewInfiniteTimeout(),
		},
	)
	require.NoError(t, err)
	txe, err := tx.Base64()
	require.NoError(t, err)

	txApprovalResp, err := handler.txApprove(ctx, txApproveRequest{Tx: txe})
	require.NoError(t, err)

	var callbackID string
	q := `SELECT callback_id FROM accounts_kyc_status WHERE diamcircle_address = $1`
	err = conn.QueryRowContext(ctx, q, senderKP.Address()).Scan(&callbackID)
	require.NoError(t, err)

	wantResp := &txApprovalResponse{
		Status:       sep8StatusActionRequired,
		Message:      "Payments exceeding 500.00 GOAT require KYC approval. Please provide an email address.",
		ActionMethod: "POST",
		StatusCode:   http.StatusOK,
		ActionURL:    "https://example.com/kyc-status/" + callbackID,
		ActionFields: []string{"email_address"},
	}
	require.Equal(t, wantResp, txApprovalResp)
}

func TestTxApproveHandler_txApprove_revised(t *testing.T) {
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
			Sequence:  "2",
		}, nil)

	handler := txApproveHandler{
		issuerKP:          issuerKP,
		assetCode:         assetGOAT.GetCode(),
		auroraClient:     &auroraMock,
		networkPassphrase: network.TestNetworkPassphrase,
		db:                conn,
		kycThreshold:      kycThresholdAmount,
		baseURL:           "https://example.com",
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount: &aurora.Account{
				AccountID: senderKP.Address(),
				Sequence:  "2",
			},
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.Payment{
					Destination: receiverKP.Address(),
					Amount:      "500",
					Asset:       assetGOAT,
				},
			},
			BaseFee:    txnbuild.MinBaseFee,
			Timebounds: txnbuild.NewInfiniteTimeout(),
		},
	)
	require.NoError(t, err)
	txe, err := tx.Base64()
	require.NoError(t, err)

	txApprovalResp, err := handler.txApprove(ctx, txApproveRequest{Tx: txe})
	require.NoError(t, err)
	require.Equal(t, sep8StatusRevised, txApprovalResp.Status)
	require.Equal(t, http.StatusOK, txApprovalResp.StatusCode)
	require.Equal(t, "Authorization and deauthorization operations were added.", txApprovalResp.Message)

	gotGenericTx, err := txnbuild.TransactionFromXDR(txApprovalResp.Tx)
	require.NoError(t, err)
	gotTx, ok := gotGenericTx.Transaction()
	require.True(t, ok)
	require.Equal(t, senderKP.Address(), gotTx.SourceAccount().AccountID)
	require.Equal(t, int64(3), gotTx.SourceAccount().Sequence)

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

func TestValidateTransactionOperationsForSuccess(t *testing.T) {
	ctx := context.Background()
	senderKP := keypair.MustRandom()
	receiverKP := keypair.MustRandom()
	issuerKP := keypair.MustRandom()
	assetGOAT := txnbuild.CreditAsset{
		Code:   "GOAT",
		Issuer: issuerKP.Address(),
	}

	// rejected if number of operations is unsupported
	tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
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
		BaseFee:    300,
		Timebounds: txnbuild.NewTimeout(300),
	})
	require.NoError(t, err)

	txApprovalResp, paymentOp, paymentSource := validateTransactionOperationsForSuccess(ctx, tx, issuerKP.Address())
	assert.Equal(t, NewRejectedTxApprovalResponse("Unsupported number of operations."), txApprovalResp)
	assert.Nil(t, paymentOp)
	assert.Empty(t, paymentSource)

	// rejected if operation at index "2" is not a payment
	tx, err = txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount: &aurora.Account{
			AccountID: senderKP.Address(),
			Sequence:  "5",
		},
		IncrementSequenceNum: true,
		Operations: []txnbuild.Operation{
			&txnbuild.BumpSequence{},
			&txnbuild.BumpSequence{},
			&txnbuild.BumpSequence{},
			&txnbuild.BumpSequence{},
			&txnbuild.BumpSequence{},
		},
		BaseFee:    300,
		Timebounds: txnbuild.NewTimeout(300),
	})
	require.NoError(t, err)

	txApprovalResp, paymentOp, paymentSource = validateTransactionOperationsForSuccess(ctx, tx, issuerKP.Address())
	assert.Equal(t, NewRejectedTxApprovalResponse("There are one or more unexpected operations in the provided transaction."), txApprovalResp)
	assert.Nil(t, paymentOp)
	assert.Empty(t, paymentSource)

	// rejected if the operations list don't match the expected format [AllowTrust, AllowTrust, Payment, AllowTrust, AllowTrust]
	tx, err = txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount: &aurora.Account{
			AccountID: senderKP.Address(),
			Sequence:  "5",
		},
		IncrementSequenceNum: true,
		Operations: []txnbuild.Operation{
			&txnbuild.BumpSequence{},
			&txnbuild.BumpSequence{},
			&txnbuild.Payment{
				SourceAccount: senderKP.Address(),
				Destination:   receiverKP.Address(),
				Amount:        "1",
				Asset:         assetGOAT,
			},
			&txnbuild.BumpSequence{},
			&txnbuild.BumpSequence{},
		},
		BaseFee:    300,
		Timebounds: txnbuild.NewTimeout(300),
	})
	require.NoError(t, err)

	txApprovalResp, paymentOp, paymentSource = validateTransactionOperationsForSuccess(ctx, tx, issuerKP.Address())
	assert.Equal(t, NewRejectedTxApprovalResponse("There are one or more unexpected operations in the provided transaction."), txApprovalResp)
	assert.Nil(t, paymentOp)
	assert.Empty(t, paymentSource)

	// rejected if the values inside the operations list don't match the expected format
	tx, err = txnbuild.NewTransaction(txnbuild.TransactionParams{
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
				SourceAccount: senderKP.Address(),
				Destination:   receiverKP.Address(),
				Amount:        "1",
				Asset:         assetGOAT,
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
				Authorize:     true, // <--- this flag is the only wrong value in this transaction
				SourceAccount: issuerKP.Address(),
			},
		},
		BaseFee:    300,
		Timebounds: txnbuild.NewTimeout(300),
	})
	require.NoError(t, err)

	txApprovalResp, paymentOp, paymentSource = validateTransactionOperationsForSuccess(ctx, tx, issuerKP.Address())
	assert.Equal(t, NewRejectedTxApprovalResponse("There are one or more unexpected operations in the provided transaction."), txApprovalResp)
	assert.Nil(t, paymentOp)
	assert.Empty(t, paymentSource)

	// success
	tx, err = txnbuild.NewTransaction(txnbuild.TransactionParams{
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
				SourceAccount: senderKP.Address(),
				Destination:   receiverKP.Address(),
				Amount:        "1",
				Asset:         assetGOAT,
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

	txApprovalResp, paymentOp, paymentSource = validateTransactionOperationsForSuccess(ctx, tx, issuerKP.Address())
	assert.Nil(t, txApprovalResp)
	assert.Equal(t, senderKP.Address(), paymentSource)
	wantPaymentOp := &txnbuild.Payment{
		SourceAccount: senderKP.Address(),
		Destination:   receiverKP.Address(),
		Amount:        "1",
		Asset:         assetGOAT,
	}
	assert.Equal(t, wantPaymentOp, paymentOp)
}

func TestTxApproveHandler_handleSuccessResponseIfNeeded_revisable(t *testing.T) {
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
			Sequence:  "2",
		}, nil)

	handler := txApproveHandler{
		issuerKP:          issuerKP,
		assetCode:         assetGOAT.GetCode(),
		auroraClient:     &auroraMock,
		networkPassphrase: network.TestNetworkPassphrase,
		db:                conn,
		kycThreshold:      kycThresholdAmount,
		baseURL:           "https://example.com",
	}

	revisableTx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount: &aurora.Account{
			AccountID: senderKP.Address(),
			Sequence:  "2",
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
		BaseFee:    300,
		Timebounds: txnbuild.NewTimeout(300),
	})
	require.NoError(t, err)

	txSuccessResponse, err := handler.handleSuccessResponseIfNeeded(ctx, revisableTx)
	require.NoError(t, err)
	assert.Nil(t, txSuccessResponse)
}

func TestTxApproveHandler_handleSuccessResponseIfNeeded_rejected(t *testing.T) {
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
			Sequence:  "2",
		}, nil)

	handler := txApproveHandler{
		issuerKP:          issuerKP,
		assetCode:         assetGOAT.GetCode(),
		auroraClient:     &auroraMock,
		networkPassphrase: network.TestNetworkPassphrase,
		db:                conn,
		kycThreshold:      kycThresholdAmount,
		baseURL:           "https://example.com",
	}

	// rejected if operations don't match the expected format
	tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount: &aurora.Account{
			AccountID: senderKP.Address(),
			Sequence:  "2",
		},
		IncrementSequenceNum: true,
		Operations: []txnbuild.Operation{
			&txnbuild.BumpSequence{},
			&txnbuild.BumpSequence{},
			&txnbuild.Payment{
				SourceAccount: senderKP.Address(),
				Destination:   receiverKP.Address(),
				Amount:        "1",
				Asset:         assetGOAT,
			},
			&txnbuild.BumpSequence{},
			&txnbuild.BumpSequence{},
		},
		BaseFee:    300,
		Timebounds: txnbuild.NewTimeout(300),
	})
	require.NoError(t, err)

	txApprovalResp, err := handler.handleSuccessResponseIfNeeded(ctx, tx)
	require.NoError(t, err)
	assert.Equal(t, NewRejectedTxApprovalResponse("There are one or more unexpected operations in the provided transaction."), txApprovalResp)

	// rejected if attempting to transfer an asset to its own issuer
	tx, err = txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount: &aurora.Account{
			AccountID: senderKP.Address(),
			Sequence:  "2",
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
				Trustor:       issuerKP.Address(),
				Type:          assetGOAT,
				Authorize:     true,
				SourceAccount: issuerKP.Address(),
			},
			&txnbuild.Payment{
				SourceAccount: senderKP.Address(),
				Destination:   issuerKP.Address(), // <--- this will trigger the rejection
				Amount:        "1",
				Asset:         assetGOAT,
			},
			&txnbuild.AllowTrust{
				Trustor:       issuerKP.Address(),
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

	txApprovalResp, err = handler.handleSuccessResponseIfNeeded(ctx, tx)
	require.NoError(t, err)
	assert.Equal(t, NewRejectedTxApprovalResponse("Can't transfer asset to its issuer."), txApprovalResp)

	// rejected if sequence number is not incremental
	compliantOps := []txnbuild.Operation{
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
			SourceAccount: senderKP.Address(),
			Destination:   receiverKP.Address(),
			Amount:        "1",
			Asset:         assetGOAT,
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
	}
	tx, err = txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount: &aurora.Account{
			AccountID: senderKP.Address(),
			Sequence:  "3",
		},
		IncrementSequenceNum: true,
		Operations:           compliantOps,
		BaseFee:              300,
		Timebounds:           txnbuild.NewTimeout(300),
	})
	require.NoError(t, err)

	txApprovalResp, err = handler.handleSuccessResponseIfNeeded(ctx, tx)
	require.NoError(t, err)
	assert.Equal(t, NewRejectedTxApprovalResponse("Invalid transaction sequence number."), txApprovalResp)
}

func TestTxApproveHandler_handleSuccessResponseIfNeeded_actionRequired(t *testing.T) {
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
			Sequence:  "2",
		}, nil)

	handler := txApproveHandler{
		issuerKP:          issuerKP,
		assetCode:         assetGOAT.GetCode(),
		auroraClient:     &auroraMock,
		networkPassphrase: network.TestNetworkPassphrase,
		db:                conn,
		kycThreshold:      kycThresholdAmount,
		baseURL:           "https://example.com",
	}

	// compliant operations with a payment above threshold will return "action_required" if the user hasn't gone through KYC yet
	tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount: &aurora.Account{
			AccountID: senderKP.Address(),
			Sequence:  "2",
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
				SourceAccount: senderKP.Address(),
				Destination:   receiverKP.Address(),
				Amount:        "501",
				Asset:         assetGOAT,
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

	txApprovalResponse, err := handler.handleSuccessResponseIfNeeded(ctx, tx)
	require.NoError(t, err)

	var callbackID string
	q := `SELECT callback_id FROM accounts_kyc_status WHERE diamcircle_address = $1`
	err = conn.QueryRowContext(ctx, q, senderKP.Address()).Scan(&callbackID)
	require.NoError(t, err)

	wantTxApprovalResponse := NewActionRequiredTxApprovalResponse(
		"Payments exceeding 500.00 GOAT require KYC approval. Please provide an email address.",
		"https://example.com/kyc-status/"+callbackID,
		[]string{"email_address"},
	)
	assert.Equal(t, wantTxApprovalResponse, txApprovalResponse)

	// compliant operations with a payment above threshold will return "rejected" if the user's KYC was rejected
	query := `
		UPDATE accounts_kyc_status
		SET
			approved_at = NULL,
			rejected_at = NOW(),
			pending_at = NULL
		WHERE diamcircle_address = $1
	`
	_, err = handler.db.ExecContext(ctx, query, senderKP.Address())
	require.NoError(t, err)
	txApprovalResponse, err = handler.handleSuccessResponseIfNeeded(ctx, tx)
	require.NoError(t, err)
	assert.Equal(t, NewRejectedTxApprovalResponse("Your KYC was rejected and you're not authorized for operations above 500.00 GOAT."), txApprovalResponse)

	// compliant operations with a payment above threshold will return "pending" if the user's KYC was marked as pending
	query = `
		UPDATE accounts_kyc_status
		SET
			approved_at = NULL,
			rejected_at = NULL,
			pending_at = NOW()
		WHERE diamcircle_address = $1
	`
	_, err = handler.db.ExecContext(ctx, query, senderKP.Address())
	require.NoError(t, err)
	txApprovalResponse, err = handler.handleSuccessResponseIfNeeded(ctx, tx)
	require.NoError(t, err)
	assert.Equal(t, NewPendingTxApprovalResponse("Your account could not be verified as approved nor rejected and was marked as pending. You will need staff authorization for operations above 500.00 GOAT."), txApprovalResponse)

	// compliant operations with a payment above threshold will return "success" if the user's KYC was approved
	query = `
		UPDATE accounts_kyc_status
		SET
			approved_at = NOW(),
			rejected_at = NULL,
			pending_at = NULL
		WHERE diamcircle_address = $1
	`
	_, err = handler.db.ExecContext(ctx, query, senderKP.Address())
	require.NoError(t, err)
	txApprovalResponse, err = handler.handleSuccessResponseIfNeeded(ctx, tx)
	require.NoError(t, err)
	assert.Equal(t, NewSuccessTxApprovalResponse(txApprovalResponse.Tx, "Transaction is compliant and signed by the issuer."), txApprovalResponse)
}

func TestTxApproveHandler_handleSuccessResponseIfNeeded_success(t *testing.T) {
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
			Sequence:  "2",
		}, nil)

	handler := txApproveHandler{
		issuerKP:          issuerKP,
		assetCode:         assetGOAT.GetCode(),
		auroraClient:     &auroraMock,
		networkPassphrase: network.TestNetworkPassphrase,
		db:                conn,
		kycThreshold:      kycThresholdAmount,
		baseURL:           "https://example.com",
	}

	tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount: &aurora.Account{
			AccountID: senderKP.Address(),
			Sequence:  "2",
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
				SourceAccount: senderKP.Address(),
				Destination:   receiverKP.Address(),
				Amount:        "1",
				Asset:         assetGOAT,
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

	txApprovalResponse, err := handler.handleSuccessResponseIfNeeded(ctx, tx)
	require.NoError(t, err)
	require.Equal(t, NewSuccessTxApprovalResponse(txApprovalResponse.Tx, "Transaction is compliant and signed by the issuer."), txApprovalResponse)

	gotGenericTx, err := txnbuild.TransactionFromXDR(txApprovalResponse.Tx)
	require.NoError(t, err)
	gotTx, ok := gotGenericTx.Transaction()
	require.True(t, ok)

	// test transaction params
	assert.Equal(t, tx.SourceAccount(), gotTx.SourceAccount())
	assert.Equal(t, tx.BaseFee(), gotTx.BaseFee())
	assert.Equal(t, tx.Timebounds(), gotTx.Timebounds())
	assert.Equal(t, tx.SequenceNumber(), gotTx.SequenceNumber())

	// test if the operations are as expected
	resp, _, _ := validateTransactionOperationsForSuccess(ctx, gotTx, issuerKP.Address())
	assert.Nil(t, resp)

	// check if the transaction contains the issuer's signature
	gotTxHash, err := gotTx.Hash(handler.networkPassphrase)
	require.NoError(t, err)
	err = handler.issuerKP.Verify(gotTxHash[:], gotTx.Signatures()[0].Signature)
	require.NoError(t, err)
}

func TestConvertAmountToReadableString(t *testing.T) {
	parsedAmount, err := amount.ParseInt64("500")
	require.NoError(t, err)
	assert.Equal(t, int64(5000000000), parsedAmount)

	readableAmount, err := convertAmountToReadableString(parsedAmount)
	require.NoError(t, err)
	assert.Equal(t, "500.00", readableAmount)
}
