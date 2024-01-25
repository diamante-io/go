package txsub

import (
	"encoding/hex"
	"errors"
	"fmt"

	"go/services/aurora/internal/codes"
	"go/xdr"
)

var (
	ErrNoResults = errors.New("No result found")
	ErrCanceled  = errors.New("canceled")
	ErrTimeout   = errors.New("timeout")

	// ErrBadSequence is a canned error response for transactions whose sequence
	// number is wrong.
	ErrBadSequence = &FailedTransactionError{"AAAAAAAAAAD////7AAAAAA=="}
	// ErrNoAccount is returned when the source account for the transaction
	// cannot be found in the database
	ErrNoAccount = &FailedTransactionError{"AAAAAAAAAAD////4AAAAAA=="}
)

// FailedTransactionError represent an error that occurred because
// diamcircle-core rejected the transaction.  ResultXDR is a base64
// encoded TransactionResult struct
type FailedTransactionError struct {
	ResultXDR string
}

func (err *FailedTransactionError) Error() string {
	return fmt.Sprintf("tx failed: %s", err.ResultXDR)
}

func (fte *FailedTransactionError) Result() (result xdr.TransactionResult, err error) {
	err = xdr.SafeUnmarshalBase64(fte.ResultXDR, &result)
	return
}

// ResultCodes represents the result codes from a request attempting to submit a fee bump transaction.
type ResultCodes struct {
	Code      string
	InnerCode string
}

func (fte *FailedTransactionError) TransactionResultCodes(transactionHash string) (result ResultCodes, err error) {
	r, err := fte.Result()
	if err != nil {
		return
	}

	if innerResultPair, ok := r.Result.GetInnerResultPair(); ok {
		// This handles the case of a transaction which was fee bumped by another request.
		// The request submitting the inner transaction should have a view of the inner result,
		// instead of the fee bump transaction.
		if transactionHash == hex.EncodeToString(innerResultPair.TransactionHash[:]) {
			result.Code, err = codes.String(innerResultPair.Result.Result.Code)
			return
		}
		result.InnerCode, err = codes.String(innerResultPair.Result.Result.Code)
		if err != nil {
			return
		}
	}
	result.Code, err = codes.String(r.Result.Code)
	return
}

func (fte *FailedTransactionError) OperationResultCodes() (result []string, err error) {
	r, err := fte.Result()
	if err != nil {
		return
	}
	oprs, ok := r.OperationResults()

	if !ok {
		return
	}

	result = make([]string, len(oprs))

	for i, opr := range oprs {
		result[i], err = codes.ForOperationResult(opr)
		if err != nil {
			return
		}
	}

	return
}
