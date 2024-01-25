package actions

import (
	"context"
	"encoding/hex"
	"mime"
	"net/http"

	"github.com/diamcircle/go/network"
	"github.com/diamcircle/go/protocols/aurora"
	hProblem "github.com/diamcircle/go/services/aurora/internal/render/problem"
	"github.com/diamcircle/go/services/aurora/internal/resourceadapter"
	"github.com/diamcircle/go/services/aurora/internal/txsub"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/support/render/hal"
	"github.com/diamcircle/go/support/render/problem"
	"github.com/diamcircle/go/xdr"
)

type NetworkSubmitter interface {
	Submit(
		ctx context.Context,
		rawTx string,
		envelope xdr.TransactionEnvelope,
		hash string) <-chan txsub.Result
}

type SubmitTransactionHandler struct {
	Submitter         NetworkSubmitter
	NetworkPassphrase string
	CoreStateGetter
}

type envelopeInfo struct {
	hash   string
	raw    string
	parsed xdr.TransactionEnvelope
}

func extractEnvelopeInfo(raw string, passphrase string) (envelopeInfo, error) {
	result := envelopeInfo{raw: raw}
	err := xdr.SafeUnmarshalBase64(raw, &result.parsed)
	if err != nil {
		return result, err
	}

	var hash [32]byte
	hash, err = network.HashTransactionInEnvelope(result.parsed, passphrase)
	if err != nil {
		return result, err
	}
	result.hash = hex.EncodeToString(hash[:])
	return result, nil
}

func (handler SubmitTransactionHandler) validateBodyType(r *http.Request) error {
	c := r.Header.Get("Content-Type")
	if c == "" {
		return nil
	}

	mt, _, err := mime.ParseMediaType(c)
	if err != nil {
		return errors.Wrap(err, "Could not determine mime type")
	}

	if mt != "application/x-www-form-urlencoded" && mt != "multipart/form-data" {
		return &hProblem.UnsupportedMediaType
	}
	return nil
}

func (handler SubmitTransactionHandler) response(r *http.Request, info envelopeInfo, result txsub.Result) (hal.Pageable, error) {
	if result.Err == nil {
		var resource aurora.Transaction
		err := resourceadapter.PopulateTransaction(
			r.Context(),
			info.hash,
			&resource,
			result.Transaction,
		)
		return resource, err
	}

	if result.Err == txsub.ErrTimeout {
		return nil, &hProblem.Timeout
	}

	if result.Err == txsub.ErrCanceled {
		return nil, &hProblem.ClientDisconnected
	}

	switch err := result.Err.(type) {
	case *txsub.FailedTransactionError:
		rcr := aurora.TransactionResultCodes{}
		resourceadapter.PopulateTransactionResultCodes(
			r.Context(),
			info.hash,
			&rcr,
			err,
		)

		return nil, &problem.P{
			Type:   "transaction_failed",
			Title:  "Transaction Failed",
			Status: http.StatusBadRequest,
			Detail: "The transaction failed when submitted to the diamcircle network. " +
				"The `extras.result_codes` field on this response contains further " +
				"details.  Descriptions of each code can be found at: " +
				"https://developers.diamcircle.org/api/errors/http-status-codes/aurora-specific/transaction-failed/",
			Extras: map[string]interface{}{
				"envelope_xdr": info.raw,
				"result_xdr":   err.ResultXDR,
				"result_codes": rcr,
			},
		}
	}

	return nil, result.Err
}

func (handler SubmitTransactionHandler) GetResource(w HeaderWriter, r *http.Request) (interface{}, error) {
	if err := handler.validateBodyType(r); err != nil {
		return nil, err
	}

	raw, err := getString(r, "tx")
	if err != nil {
		return nil, err
	}

	info, err := extractEnvelopeInfo(raw, handler.NetworkPassphrase)
	if err != nil {
		return nil, &problem.P{
			Type:   "transaction_malformed",
			Title:  "Transaction Malformed",
			Status: http.StatusBadRequest,
			Detail: "Aurora could not decode the transaction envelope in this " +
				"request. A transaction should be an XDR TransactionEnvelope struct " +
				"encoded using base64.  The envelope read from this request is " +
				"echoed in the `extras.envelope_xdr` field of this response for your " +
				"convenience.",
			Extras: map[string]interface{}{
				"envelope_xdr": raw,
			},
		}
	}

	coreState := handler.GetCoreState()
	if !coreState.Synced {
		return nil, hProblem.StaleHistory
	}

	submission := handler.Submitter.Submit(
		r.Context(),
		info.raw,
		info.parsed,
		info.hash,
	)

	select {
	case result := <-submission:
		return handler.response(r, info, result)
	case <-r.Context().Done():
		if r.Context().Err() == context.Canceled {
			return nil, hProblem.ClientDisconnected
		}
		return nil, hProblem.Timeout
	}
}
