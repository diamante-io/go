package httperror

import (
	"net/http"
	"testing"

	"github.com/diamcircle/go/clients/auroraclient"
	hProtocol "github.com/diamcircle/go/protocols/aurora"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/support/render/problem"
	"github.com/stretchr/testify/require"
)

func TestParseAuroraError(t *testing.T) {
	err := ParseAuroraError(nil)
	require.Nil(t, err)

	err = ParseAuroraError(errors.New("some error"))
	require.EqualError(t, err, "error submitting transaction: some error")

	auroraError := auroraclient.Error{
		Problem: problem.P{
			Type:   "bad_request",
			Title:  "Bad Request",
			Status: http.StatusBadRequest,
			Extras: map[string]interface{}{
				"result_codes": hProtocol.TransactionResultCodes{
					TransactionCode:      "tx_code_here",
					InnerTransactionCode: "",
					OperationCodes: []string{
						"op_success",
						"op_bad_auth",
					},
				},
			},
		},
	}
	err = ParseAuroraError(auroraError)
	require.EqualError(t, err, "error submitting transaction: problem: bad_request, &{TransactionCode:tx_code_here InnerTransactionCode: OperationCodes:[op_success op_bad_auth]}\n: aurora error: \"Bad Request\" (tx_code_here, op_success, op_bad_auth) - check aurora.Error.Problem for more information")
}
