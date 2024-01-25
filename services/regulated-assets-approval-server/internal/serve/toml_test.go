package serve

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"go/network"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTomlHandler_validate(t *testing.T) {
	// empty network passphrase
	h := diamcircleTOMLHandler{}
	err := h.validate()
	require.EqualError(t, err, "network passphrase cannot be empty")

	// empty asset code
	h = diamcircleTOMLHandler{
		networkPassphrase: network.TestNetworkPassphrase,
	}
	err = h.validate()
	require.EqualError(t, err, "asset code cannot be empty")

	// empty asset issuer address
	h = diamcircleTOMLHandler{
		networkPassphrase: network.TestNetworkPassphrase,
		assetCode:         "FOOBAR",
	}
	err = h.validate()
	require.EqualError(t, err, "asset issuer address cannot be empty")

	// invalid asset issuer address
	h = diamcircleTOMLHandler{
		networkPassphrase: network.TestNetworkPassphrase,
		assetCode:         "FOOBAR",
		issuerAddress:     "foobar",
	}
	err = h.validate()
	require.EqualError(t, err, "asset issuer address is not a valid public key")

	// empty approval server
	h = diamcircleTOMLHandler{
		networkPassphrase: network.TestNetworkPassphrase,
		assetCode:         "FOOBAR",
		issuerAddress:     "GCVDOU4YHHXGM3QYVSDHPQIFMZKXTFSIYO4HJOJZOTR7GURVQO6IQ5HM",
	}
	err = h.validate()
	require.EqualError(t, err, "approval server cannot be empty")

	// empty kyc threshold
	h = diamcircleTOMLHandler{
		networkPassphrase: network.TestNetworkPassphrase,
		assetCode:         "FOOBAR",
		issuerAddress:     "GCVDOU4YHHXGM3QYVSDHPQIFMZKXTFSIYO4HJOJZOTR7GURVQO6IQ5HM",
		approvalServer:    "localhost:8000/tx-approve",
	}
	err = h.validate()
	require.EqualError(t, err, "kyc threshold cannot be less than or equal to zero")

	// negative kyc threshold
	h = diamcircleTOMLHandler{
		networkPassphrase: network.TestNetworkPassphrase,
		assetCode:         "FOOBAR",
		issuerAddress:     "GCVDOU4YHHXGM3QYVSDHPQIFMZKXTFSIYO4HJOJZOTR7GURVQO6IQ5HM",
		approvalServer:    "localhost:8000/tx-approve",
		kycThreshold:      -500,
	}
	err = h.validate()
	require.EqualError(t, err, "kyc threshold cannot be less than or equal to zero")

	// success
	h = diamcircleTOMLHandler{
		networkPassphrase: network.TestNetworkPassphrase,
		assetCode:         "FOOBAR",
		issuerAddress:     "GCVDOU4YHHXGM3QYVSDHPQIFMZKXTFSIYO4HJOJZOTR7GURVQO6IQ5HM",
		approvalServer:    "localhost:8000/tx-approve",
		kycThreshold:      500,
	}
	err = h.validate()
	require.NoError(t, err)
}

func TestTomlHandler_ServeHTTP(t *testing.T) {
	mux := chi.NewMux()
	mux.Get("/.well-known/diamcircle.toml", diamcircleTOMLHandler{
		networkPassphrase: network.TestNetworkPassphrase,
		assetCode:         "FOO",
		issuerAddress:     "GCVDOU4YHHXGM3QYVSDHPQIFMZKXTFSIYO4HJOJZOTR7GURVQO6IQ5HM",
		approvalServer:    "localhost:8000/tx-approve",
		kycThreshold:      5000000000,
	}.ServeHTTP)

	ctx := context.Background()
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/.well-known/diamcircle.toml", nil)
	r = r.WithContext(ctx)
	mux.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "text/plain; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `NETWORK_PASSPHRASE="` + network.TestNetworkPassphrase + `"
[[CURRENCIES]]
code="FOO"
issuer="GCVDOU4YHHXGM3QYVSDHPQIFMZKXTFSIYO4HJOJZOTR7GURVQO6IQ5HM"
regulated=true
approval_server="localhost:8000/tx-approve"
approval_criteria="The approval server currently only accepts payments. The transaction must have exactly one operation of type payment. If the payment amount exceeds 500.00 FOO it will need KYC approval if the account hasn’t been previously approved."`
	require.Equal(t, wantBody, string(body))
}
