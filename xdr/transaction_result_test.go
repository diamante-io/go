package xdr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTxResult(code TransactionResultCode) TransactionResult {
	return TransactionResult{
		FeeCharged: 123,
		Result: TransactionResultResult{
			Code:    code,
			Results: &[]OperationResult{},
		},
	}
}

func TestSuccessful(t *testing.T) {
	for _, testCase := range []struct {
		code     TransactionResultCode
		expected bool
	}{
		{TransactionResultCodeTxSuccess, true},
		{TransactionResultCodeTxFeeBumpInnerSuccess, true},
		{TransactionResultCodeTxFailed, false},
		{TransactionResultCodeTxFeeBumpInnerFailed, false},
		{TransactionResultCodeTxBadSeq, false},
	} {

		result := createTxResult(testCase.code)
		assert.Equal(t, testCase.expected, result.Successful())
		resultPair := TransactionResultPair{
			Result: result,
		}
		assert.Equal(t, testCase.expected, resultPair.Successful())
	}
}

func TestOperationResults(t *testing.T) {
	successfulEmptyTx := createTxResult(TransactionResultCodeTxSuccess)
	successfulEmptyTx.Result.Results = &[]OperationResult{}

	successfulEmptyFeeBumpTx := createTxResult(TransactionResultCodeTxFeeBumpInnerSuccess)
	successfulEmptyFeeBumpTx.Result.InnerResultPair = &InnerTransactionResultPair{
		Result: InnerTransactionResult{
			Result: InnerTransactionResultResult{
				Code:    TransactionResultCodeTxSuccess,
				Results: &[]OperationResult{},
			},
		},
	}

	failedEmptyTx := createTxResult(TransactionResultCodeTxFailed)
	failedEmptyTx.Result.Results = &[]OperationResult{}

	failedEmptyFeeBumpTx := createTxResult(TransactionResultCodeTxFeeBumpInnerFailed)
	failedEmptyFeeBumpTx.Result.InnerResultPair = &InnerTransactionResultPair{
		Result: InnerTransactionResult{
			Result: InnerTransactionResultResult{
				Code:    TransactionResultCodeTxFailed,
				Results: &[]OperationResult{},
			},
		},
	}

	bumpSeqOp := OperationResult{
		Tr: &OperationResultTr{
			Type: OperationTypeBumpSequence,
			BumpSeqResult: &BumpSequenceResult{
				Code: BumpSequenceResultCodeBumpSequenceSuccess,
			},
		},
	}
	inflationOp := OperationResult{
		Tr: &OperationResultTr{
			Type: OperationTypeInflation,
			InflationResult: &InflationResult{
				Code: InflationResultCodeInflationNotTime,
			},
		},
	}

	successfulTx := createTxResult(TransactionResultCodeTxSuccess)
	successfulTx.Result.Results = &[]OperationResult{bumpSeqOp}

	successfulFeeBumpTx := createTxResult(TransactionResultCodeTxFeeBumpInnerSuccess)
	successfulFeeBumpTx.Result.InnerResultPair = &InnerTransactionResultPair{
		Result: InnerTransactionResult{
			Result: InnerTransactionResultResult{
				Code:    TransactionResultCodeTxSuccess,
				Results: &[]OperationResult{inflationOp},
			},
		},
	}

	failedBumpSeqOp := OperationResult{
		Tr: &OperationResultTr{
			Type: OperationTypeBumpSequence,
			BumpSeqResult: &BumpSequenceResult{
				Code: BumpSequenceResultCodeBumpSequenceBadSeq,
			},
		},
	}
	failedPaymentOp := OperationResult{
		Tr: &OperationResultTr{
			Type: OperationTypePayment,
			PaymentResult: &PaymentResult{
				Code: PaymentResultCodePaymentMalformed,
			},
		},
	}

	failedTx := createTxResult(TransactionResultCodeTxFailed)
	failedTx.Result.Results = &[]OperationResult{failedBumpSeqOp}

	failedFeeBumpTx := createTxResult(TransactionResultCodeTxFeeBumpInnerFailed)
	failedFeeBumpTx.Result.InnerResultPair = &InnerTransactionResultPair{
		Result: InnerTransactionResult{
			Result: InnerTransactionResultResult{
				Code:    TransactionResultCodeTxFailed,
				Results: &[]OperationResult{failedPaymentOp},
			},
		},
	}

	for _, testCase := range []struct {
		result             TransactionResult
		expectedOk         bool
		expectedOperations []OperationResult
	}{
		{
			createTxResult(TransactionResultCodeTxBadSeq),
			false,
			nil,
		},
		{
			successfulEmptyTx,
			true,
			[]OperationResult{},
		},
		{
			successfulEmptyFeeBumpTx,
			true,
			[]OperationResult{},
		},
		{
			failedEmptyTx,
			true,
			[]OperationResult{},
		},
		{
			failedEmptyFeeBumpTx,
			true,
			[]OperationResult{},
		},
		{
			successfulTx,
			true,
			[]OperationResult{bumpSeqOp},
		},
		{
			successfulFeeBumpTx,
			true,
			[]OperationResult{inflationOp},
		},
		{
			failedTx,
			true,
			[]OperationResult{failedBumpSeqOp},
		},
		{
			failedFeeBumpTx,
			true,
			[]OperationResult{failedPaymentOp},
		},
	} {
		opResults, ok := testCase.result.OperationResults()
		assert.Equal(t, testCase.expectedOk, ok)
		assert.Equal(t, testCase.expectedOperations, opResults)

		resultPair := TransactionResultPair{
			Result: testCase.result,
		}
		opResults, ok = resultPair.OperationResults()
		assert.Equal(t, testCase.expectedOk, ok)
		assert.Equal(t, testCase.expectedOperations, opResults)
	}
}

func TestInnerHash(t *testing.T) {
	tx := TransactionResultPair{
		TransactionHash: Hash{1, 1, 1},
		Result:          createTxResult(TransactionResultCodeTxSuccess),
	}
	assert.Panics(t, func() {
		tx.InnerHash()
	})

	feeBumpTx := TransactionResultPair{
		TransactionHash: Hash{1, 1, 1},
		Result:          createTxResult(TransactionResultCodeTxFeeBumpInnerSuccess),
	}
	feeBumpTx.Result.Result.InnerResultPair = &InnerTransactionResultPair{
		TransactionHash: Hash{1, 2, 3},
	}
	assert.Equal(t, Hash{1, 2, 3}, feeBumpTx.InnerHash())
}
