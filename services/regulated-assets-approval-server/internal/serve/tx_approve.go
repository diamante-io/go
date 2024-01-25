package serve

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"go/amount"
	"go/clients/auroraclient"
	"go/keypair"
	"go/services/regulated-assets-approval-server/internal/serve/httperror"
	"go/support/errors"
	"go/support/http/httpdecode"
	"go/support/log"
	"go/txnbuild"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type txApproveHandler struct {
	issuerKP          *keypair.Full
	assetCode         string
	auroraClient      auroraclient.ClientInterface
	networkPassphrase string
	db                *sqlx.DB
	kycThreshold      int64
	baseURL           string
}

type txApproveRequest struct {
	Tx string `json:"tx" form:"tx"`
}

// validate performs some validations on the provided handler data.
func (h txApproveHandler) validate() error {
	if h.issuerKP == nil {
		return errors.New("issuer keypair cannot be nil")
	}
	if h.assetCode == "" {
		return errors.New("asset code cannot be empty")
	}
	if h.auroraClient == nil {
		return errors.New("aurora client cannot be nil")
	}
	if h.networkPassphrase == "" {
		return errors.New("network passphrase cannot be empty")
	}
	if h.db == nil {
		return errors.New("database cannot be nil")
	}
	if h.kycThreshold <= 0 {
		return errors.New("kyc threshold cannot be less than or equal to zero")
	}
	if h.baseURL == "" {
		return errors.New("base url cannot be empty")
	}
	return nil
}

func (h txApproveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := h.validate()
	if err != nil {
		log.Ctx(ctx).Error(errors.Wrap(err, "validating txApproveHandler"))
		httperror.InternalServer.Render(w)
		return
	}

	in := txApproveRequest{}
	err = httpdecode.Decode(r, &in)
	if err != nil {
		log.Ctx(ctx).Error(errors.Wrap(err, "decoding txApproveRequest"))
		httperror.BadRequest.Render(w)
		return
	}

	txApproveResp, err := h.txApprove(ctx, in)
	if err != nil {
		log.Ctx(ctx).Error(errors.Wrap(err, "validating the input transaction for approval"))
		httperror.InternalServer.Render(w)
		return
	}

	txApproveResp.Render(w)
}

// validateInput validates if the input parameters contain a valid transaction
// and if the source account is not set in a way that would harm the issuer.
func (h txApproveHandler) validateInput(ctx context.Context, in txApproveRequest) (*txApprovalResponse, *txnbuild.Transaction) {
	if in.Tx == "" {
		log.Ctx(ctx).Error(`request is missing parameter "tx".`)
		return NewRejectedTxApprovalResponse(`Missing parameter "tx".`), nil
	}

	genericTx, err := txnbuild.TransactionFromXDR(in.Tx)
	if err != nil {
		log.Ctx(ctx).Error(errors.Wrap(err, "parsing transaction xdr"))
		return NewRejectedTxApprovalResponse(`Invalid parameter "tx".`), nil
	}

	tx, ok := genericTx.Transaction()
	if !ok {
		log.Ctx(ctx).Error(`invalid parameter "tx", generic transaction not given.`)
		return NewRejectedTxApprovalResponse(`Invalid parameter "tx".`), nil
	}

	if tx.SourceAccount().AccountID == h.issuerKP.Address() {
		log.Ctx(ctx).Errorf("transaction sourceAccount is the same as the server issuer account %s", h.issuerKP.Address())
		return NewRejectedTxApprovalResponse("Transaction source account is invalid."), nil
	}

	// only AllowTrust operations can have the issuer as their source account
	for _, op := range tx.Operations() {
		if _, ok := op.(*txnbuild.AllowTrust); ok {
			continue
		}

		if op.GetSourceAccount() == h.issuerKP.Address() {
			log.Ctx(ctx).Error("transaction contains one or more unauthorized operations where source account is the issuer account")
			return NewRejectedTxApprovalResponse("There are one or more unauthorized operations in the provided transaction."), nil
		}
	}

	return nil, tx
}

// txApprove is called to validate the input transaction.
func (h txApproveHandler) txApprove(ctx context.Context, in txApproveRequest) (resp *txApprovalResponse, err error) {
	defer func() {
		log.Ctx(ctx).Debug("==== will log responses ====")
		log.Ctx(ctx).Debugf("req: %+v", in)
		log.Ctx(ctx).Debugf("resp: %+v", resp)
		log.Ctx(ctx).Debugf("err: %+v", err)
		log.Ctx(ctx).Debug("====  did log responses ====")
	}()

	rejectedResponse, tx := h.validateInput(ctx, in)
	if rejectedResponse != nil {
		return rejectedResponse, nil
	}

	txSuccessResp, err := h.handleSuccessResponseIfNeeded(ctx, tx)
	if err != nil {
		return nil, errors.Wrap(err, "checking if transaction in request was compliant")
	}
	if txSuccessResp != nil {
		return txSuccessResp, nil
	}

	// validate the revisable transaction has one operation.
	if len(tx.Operations()) != 1 {
		return NewRejectedTxApprovalResponse("Please submit a transaction with exactly one operation of type payment."), nil
	}

	paymentOp, ok := tx.Operations()[0].(*txnbuild.Payment)
	if !ok {
		log.Ctx(ctx).Error("transaction does not contain a payment operation")
		return NewRejectedTxApprovalResponse("There is one or more unauthorized operations in the provided transaction."), nil
	}
	paymentSource := paymentOp.SourceAccount
	if paymentSource == "" {
		paymentSource = tx.SourceAccount().AccountID
	}

	if paymentOp.Destination == h.issuerKP.Address() {
		return NewRejectedTxApprovalResponse("Can't transfer asset to its issuer."), nil
	}

	// validate payment asset is the one supported by the issuer
	issuerAddress := h.issuerKP.Address()
	if paymentOp.Asset.GetCode() != h.assetCode || paymentOp.Asset.GetIssuer() != issuerAddress {
		log.Ctx(ctx).Error(`the payment asset is not supported by this issuer`)
		return NewRejectedTxApprovalResponse("The payment asset is not supported by this issuer."), nil
	}

	acc, err := h.auroraClient.AccountDetail(auroraclient.AccountRequest{AccountID: paymentSource})
	if err != nil {
		return nil, errors.Wrapf(err, "getting detail for payment source account %s", paymentSource)
	}

	// validate the sequence number
	accountSequence, err := strconv.ParseInt(acc.Sequence, 10, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing account sequence number %q from string to int64", acc.Sequence)
	}
	if tx.SourceAccount().Sequence != accountSequence+1 {
		log.Ctx(ctx).Errorf(`invalid transaction sequence number tx.SourceAccount().Sequence: %d, accountSequence+1: %d`, tx.SourceAccount().Sequence, accountSequence+1)
		return NewRejectedTxApprovalResponse("Invalid transaction sequence number."), nil
	}

	actionRequiredResponse, err := h.handleActionRequiredResponseIfNeeded(ctx, paymentSource, paymentOp)
	if err != nil {
		return nil, errors.Wrap(err, "handling KYC required payment")
	}
	if actionRequiredResponse != nil {
		return actionRequiredResponse, nil
	}

	// build the transaction
	revisedOperations := []txnbuild.Operation{
		&txnbuild.AllowTrust{
			Trustor:       paymentSource,
			Type:          paymentOp.Asset,
			Authorize:     true,
			SourceAccount: issuerAddress,
		},
		&txnbuild.AllowTrust{
			Trustor:       paymentOp.Destination,
			Type:          paymentOp.Asset,
			Authorize:     true,
			SourceAccount: issuerAddress,
		},
		paymentOp,
		&txnbuild.AllowTrust{
			Trustor:       paymentOp.Destination,
			Type:          paymentOp.Asset,
			Authorize:     false,
			SourceAccount: issuerAddress,
		},
		&txnbuild.AllowTrust{
			Trustor:       paymentSource,
			Type:          paymentOp.Asset,
			Authorize:     false,
			SourceAccount: issuerAddress,
		},
	}
	revisedTx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount:        &acc,
		IncrementSequenceNum: true,
		Operations:           revisedOperations,
		BaseFee:              300,
		Timebounds:           txnbuild.NewTimeout(300),
	})
	if err != nil {
		return nil, errors.Wrap(err, "building transaction")
	}

	revisedTx, err = revisedTx.Sign(h.networkPassphrase, h.issuerKP)
	if err != nil {
		return nil, errors.Wrap(err, "signing transaction")
	}
	txe, err := revisedTx.Base64()
	if err != nil {
		return nil, errors.Wrap(err, "encoding revised transaction")
	}

	return NewRevisedTxApprovalResponse(txe), nil
}

// handleActionRequiredResponseIfNeeded validates and returns an action_required
// response if the payment requires KYC.
func (h txApproveHandler) handleActionRequiredResponseIfNeeded(ctx context.Context, diamcircleAddress string, paymentOp *txnbuild.Payment) (*txApprovalResponse, error) {
	paymentAmount, err := amount.ParseInt64(paymentOp.Amount)
	if err != nil {
		return nil, errors.Wrap(err, "parsing payment amount from string to Int64")
	}
	if paymentAmount <= h.kycThreshold {
		return nil, nil
	}

	intendedCallbackID := uuid.New().String()
	const q = `
		WITH new_row AS (
			INSERT INTO accounts_kyc_status (diamcircle_address, callback_id)
			VALUES ($1, $2)
			ON CONFLICT(diamcircle_address) DO NOTHING
			RETURNING *
		)
		SELECT callback_id, approved_at, rejected_at, pending_at FROM new_row
		UNION
		SELECT callback_id, approved_at, rejected_at, pending_at
		FROM accounts_kyc_status
		WHERE diamcircle_address = $1
	`
	var (
		callbackID                        string
		approvedAt, rejectedAt, pendingAt sql.NullTime
	)
	err = h.db.QueryRowContext(ctx, q, diamcircleAddress, intendedCallbackID).Scan(&callbackID, &approvedAt, &rejectedAt, &pendingAt)
	if err != nil {
		return nil, errors.Wrap(err, "inserting new row into accounts_kyc_status table")
	}

	if approvedAt.Valid {
		return nil, nil
	}

	kycThreshold, err := convertAmountToReadableString(h.kycThreshold)
	if err != nil {
		return nil, errors.Wrap(err, "converting kycThreshold to human readable string")
	}

	if rejectedAt.Valid {
		return NewRejectedTxApprovalResponse(fmt.Sprintf("Your KYC was rejected and you're not authorized for operations above %s %s.", kycThreshold, h.assetCode)), nil
	}

	if pendingAt.Valid {
		return NewPendingTxApprovalResponse(fmt.Sprintf("Your account could not be verified as approved nor rejected and was marked as pending. You will need staff authorization for operations above %s %s.", kycThreshold, h.assetCode)), nil
	}

	return NewActionRequiredTxApprovalResponse(
		fmt.Sprintf(`Payments exceeding %s %s require KYC approval. Please provide an email address.`, kycThreshold, h.assetCode),
		fmt.Sprintf("%s/kyc-status/%s", h.baseURL, callbackID),
		[]string{"email_address"},
	), nil
}

// handleSuccessResponseIfNeeded inspects the incoming transaction and returns a
// "success" response if it's already compliant with the SEP-8 authorization spec.
func (h txApproveHandler) handleSuccessResponseIfNeeded(ctx context.Context, tx *txnbuild.Transaction) (*txApprovalResponse, error) {
	if len(tx.Operations()) != 5 {
		return nil, nil
	}

	rejectedResp, paymentOp, paymentSource := validateTransactionOperationsForSuccess(ctx, tx, h.issuerKP.Address())
	if rejectedResp != nil {
		return rejectedResp, nil
	}

	if paymentOp.Destination == h.issuerKP.Address() {
		return NewRejectedTxApprovalResponse("Can't transfer asset to its issuer."), nil
	}

	// pull current account details from the network then validate the tx sequence number
	acc, err := h.auroraClient.AccountDetail(auroraclient.AccountRequest{AccountID: paymentSource})
	if err != nil {
		return nil, errors.Wrapf(err, "getting detail for payment source account %s", paymentSource)
	}
	accountSequence, err := strconv.ParseInt(acc.Sequence, 10, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing account sequence number %q from string to int64", acc.Sequence)
	}
	if tx.SourceAccount().Sequence != accountSequence+1 {
		log.Ctx(ctx).Errorf(`invalid transaction sequence number tx.SourceAccount().Sequence: %d, accountSequence+1: %d`, tx.SourceAccount().Sequence, accountSequence+1)
		return NewRejectedTxApprovalResponse("Invalid transaction sequence number."), nil
	}

	kycRequiredResponse, err := h.handleActionRequiredResponseIfNeeded(ctx, paymentSource, paymentOp)
	if err != nil {
		return nil, errors.Wrap(err, "handling KYC required payment")
	}
	if kycRequiredResponse != nil {
		return kycRequiredResponse, nil
	}

	// sign transaction with issuer's signature and encode it
	tx, err = tx.Sign(h.networkPassphrase, h.issuerKP)
	if err != nil {
		return nil, errors.Wrap(err, "signing transaction")
	}
	txe, err := tx.Base64()
	if err != nil {
		return nil, errors.Wrap(err, "encoding revised transaction")
	}

	return NewSuccessTxApprovalResponse(txe, "Transaction is compliant and signed by the issuer."), nil
}

// validateTransactionOperationsForSuccess checks if the incoming transaction
// operations are compliant with the anchor's SEP-8 policy.
func validateTransactionOperationsForSuccess(ctx context.Context, tx *txnbuild.Transaction, issuerAddress string) (resp *txApprovalResponse, paymentOp *txnbuild.Payment, paymentSource string) {
	if len(tx.Operations()) != 5 {
		return NewRejectedTxApprovalResponse("Unsupported number of operations."), nil, ""
	}

	// extract the payment operation and payment source account.
	paymentOp, ok := tx.Operations()[2].(*txnbuild.Payment)
	if !ok {
		log.Ctx(ctx).Error(`third operation is not of type payment`)
		return NewRejectedTxApprovalResponse("There are one or more unexpected operations in the provided transaction."), nil, ""
	}
	paymentSource = paymentOp.SourceAccount
	if paymentSource == "" {
		paymentSource = tx.SourceAccount().AccountID
	}

	assetCode := paymentOp.Asset.GetCode()

	operationsValid := func() bool {
		op0, ok := tx.Operations()[0].(*txnbuild.AllowTrust)
		if !ok ||
			op0.Trustor != paymentSource ||
			op0.Type.GetCode() != assetCode ||
			!op0.Authorize ||
			op0.SourceAccount != issuerAddress {
			return false
		}

		op1, ok := tx.Operations()[1].(*txnbuild.AllowTrust)
		if !ok ||
			op1.Trustor != paymentOp.Destination ||
			op1.Type.GetCode() != assetCode ||
			!op1.Authorize ||
			op1.SourceAccount != issuerAddress {
			return false
		}

		op2, ok := tx.Operations()[2].(*txnbuild.Payment)
		if !ok || op2 != paymentOp {
			return false
		}

		op3, ok := tx.Operations()[3].(*txnbuild.AllowTrust)
		if !ok ||
			op3.Trustor != paymentOp.Destination ||
			op3.Type.GetCode() != assetCode ||
			op3.Authorize ||
			op3.SourceAccount != issuerAddress {
			return false
		}

		op4, ok := tx.Operations()[4].(*txnbuild.AllowTrust)
		if !ok ||
			op4.Trustor != paymentSource ||
			op4.Type.GetCode() != assetCode ||
			op4.Authorize ||
			op4.SourceAccount != issuerAddress {
			return false
		}

		return true
	}()
	if !operationsValid {
		return NewRejectedTxApprovalResponse("There are one or more unexpected operations in the provided transaction."), nil, ""
	}

	return nil, paymentOp, paymentSource
}

func convertAmountToReadableString(threshold int64) (string, error) {
	amountStr := amount.StringFromInt64(threshold)
	amountFloat, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return "", errors.Wrap(err, "converting threshold amount from string to float")
	}
	return fmt.Sprintf("%.2f", amountFloat), nil
}
