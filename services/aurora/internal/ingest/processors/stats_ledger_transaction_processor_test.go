package processors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/diamcircle/go/ingest"
	"github.com/diamcircle/go/xdr"
)

func TestStatsLedgerTransactionProcessoAllOpTypesCovered(t *testing.T) {
	txTemplate := ingest.LedgerTransaction{
		Envelope: xdr.TransactionEnvelope{
			Type: xdr.EnvelopeTypeEnvelopeTypeTx,
			V1: &xdr.TransactionV1Envelope{
				Tx: xdr.Transaction{
					Operations: []xdr.Operation{
						{Body: xdr.OperationBody{Type: 0}},
					},
				},
			},
		},
	}
	for typ, s := range xdr.OperationTypeToStringMap {
		tx := txTemplate
		txTemplate.Envelope.V1.Tx.Operations[0].Body.Type = xdr.OperationType(typ)
		f := func() {
			var p StatsLedgerTransactionProcessor
			p.ProcessTransaction(context.Background(), tx)
		}
		assert.NotPanics(t, f, s)
	}

	// make sure the check works for an unreasonable operation type
	tx := txTemplate
	txTemplate.Envelope.V1.Tx.Operations[0].Body.Type = 20000
	f := func() {
		var p StatsLedgerTransactionProcessor
		p.ProcessTransaction(context.Background(), tx)
	}
	assert.Panics(t, f)
}

func TestStatsLedgerTransactionProcessor(t *testing.T) {
	processor := &StatsLedgerTransactionProcessor{}

	// Successful
	assert.NoError(t, processor.ProcessTransaction(context.Background(), ingest.LedgerTransaction{
		Result: xdr.TransactionResultPair{
			Result: xdr.TransactionResult{
				Result: xdr.TransactionResultResult{
					Code: xdr.TransactionResultCodeTxSuccess,
				},
			},
		},
		Envelope: xdr.TransactionEnvelope{
			Type: xdr.EnvelopeTypeEnvelopeTypeTx,
			V1: &xdr.TransactionV1Envelope{
				Tx: xdr.Transaction{
					Operations: []xdr.Operation{
						{Body: xdr.OperationBody{Type: xdr.OperationTypeCreateAccount}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypePayment}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypePathPaymentStrictReceive}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeManageSellOffer}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeCreatePassiveSellOffer}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeSetOptions}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeChangeTrust}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeAllowTrust}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeAccountMerge}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeInflation}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeManageData}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeBumpSequence}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeManageBuyOffer}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypePathPaymentStrictSend}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeCreateClaimableBalance}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeClaimClaimableBalance}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeBeginSponsoringFutureReserves}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeEndSponsoringFutureReserves}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeRevokeSponsorship}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeClawback}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeClawbackClaimableBalance}},
					},
				},
			},
		},
	}))

	// Failed
	assert.NoError(t, processor.ProcessTransaction(context.Background(), ingest.LedgerTransaction{
		Result: xdr.TransactionResultPair{
			Result: xdr.TransactionResult{
				Result: xdr.TransactionResultResult{
					Code: xdr.TransactionResultCodeTxFailed,
				},
			},
		},
		Envelope: xdr.TransactionEnvelope{
			Type: xdr.EnvelopeTypeEnvelopeTypeTx,
			V1: &xdr.TransactionV1Envelope{
				Tx: xdr.Transaction{
					Operations: []xdr.Operation{
						{Body: xdr.OperationBody{Type: xdr.OperationTypeCreateAccount}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypePayment}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypePathPaymentStrictReceive}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeManageSellOffer}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeCreatePassiveSellOffer}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeSetOptions}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeChangeTrust}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeAllowTrust}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeAccountMerge}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeInflation}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeManageData}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeBumpSequence}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeManageBuyOffer}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypePathPaymentStrictSend}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeCreateClaimableBalance}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeClaimClaimableBalance}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeBeginSponsoringFutureReserves}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeEndSponsoringFutureReserves}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeRevokeSponsorship}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeClawback}},
						{Body: xdr.OperationBody{Type: xdr.OperationTypeClawbackClaimableBalance}},
					},
				},
			},
		},
	}))

	results := processor.GetResults()

	assert.Equal(t, int64(2), results.Transactions)
	assert.Equal(t, int64(1), results.TransactionsSuccessful)
	assert.Equal(t, int64(1), results.TransactionsFailed)

	assert.Equal(t, int64(21*2), results.Operations)
	assert.Equal(t, int64(21), results.OperationsInSuccessful)
	assert.Equal(t, int64(21), results.OperationsInFailed)

	assert.Equal(t, int64(2), results.OperationsCreateAccount)
	assert.Equal(t, int64(2), results.OperationsPayment)
	assert.Equal(t, int64(2), results.OperationsPathPaymentStrictReceive)
	assert.Equal(t, int64(2), results.OperationsManageSellOffer)
	assert.Equal(t, int64(2), results.OperationsCreatePassiveSellOffer)
	assert.Equal(t, int64(2), results.OperationsSetOptions)
	assert.Equal(t, int64(2), results.OperationsChangeTrust)
	assert.Equal(t, int64(2), results.OperationsAllowTrust)
	assert.Equal(t, int64(2), results.OperationsAccountMerge)
	assert.Equal(t, int64(2), results.OperationsInflation)
	assert.Equal(t, int64(2), results.OperationsManageData)
	assert.Equal(t, int64(2), results.OperationsBumpSequence)
	assert.Equal(t, int64(2), results.OperationsManageBuyOffer)
	assert.Equal(t, int64(2), results.OperationsPathPaymentStrictSend)
	assert.Equal(t, int64(2), results.OperationsCreateClaimableBalance)
	assert.Equal(t, int64(2), results.OperationsClaimClaimableBalance)
	assert.Equal(t, int64(2), results.OperationsBeginSponsoringFutureReserves)
	assert.Equal(t, int64(2), results.OperationsEndSponsoringFutureReserves)
	assert.Equal(t, int64(2), results.OperationsRevokeSponsorship)
	assert.Equal(t, int64(2), results.OperationsClawback)
	assert.Equal(t, int64(2), results.OperationsClawbackClaimableBalance)
}
