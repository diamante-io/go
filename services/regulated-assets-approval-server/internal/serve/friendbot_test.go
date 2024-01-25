package serve

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/diamcircle/go/clients/auroraclient"
	"github.com/diamcircle/go/network"
	"github.com/diamcircle/go/protocols/aurora"
	"github.com/diamcircle/go/protocols/aurora/base"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/support/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFriendbotHandler_validate(t *testing.T) {
	// missing secret seed
	fh := friendbotHandler{}
	err := fh.validate()
	require.EqualError(t, err, "issuer secret cannot be empty")

	// invalid secret seed
	fh = friendbotHandler{
		issuerAccountSecret: "foo bar",
	}
	err = fh.validate()
	require.EqualError(t, err, "the provided string \"foo bar\" is not a valid Diamcircle account seed")

	// missing asset code
	fh = friendbotHandler{
		issuerAccountSecret: "SB6SFUY6ZJ2ETQHTY456GDAQ547R6NDAU74DTI2CKVVI4JERTUXKB2R4",
	}
	err = fh.validate()
	require.EqualError(t, err, "asset code cannot be empty")

	// missing aurora client
	fh = friendbotHandler{
		issuerAccountSecret: "SB6SFUY6ZJ2ETQHTY456GDAQ547R6NDAU74DTI2CKVVI4JERTUXKB2R4",
		assetCode:           "FOO",
	}
	err = fh.validate()
	require.EqualError(t, err, "aurora client cannot be nil")

	// missing aurora URL
	fh = friendbotHandler{
		issuerAccountSecret: "SB6SFUY6ZJ2ETQHTY456GDAQ547R6NDAU74DTI2CKVVI4JERTUXKB2R4",
		assetCode:           "FOO",
		auroraClient:       auroraclient.DefaultTestNetClient,
	}
	err = fh.validate()
	require.EqualError(t, err, "aurora url cannot be empty")

	// missing network passphrase
	fh = friendbotHandler{
		issuerAccountSecret: "SB6SFUY6ZJ2ETQHTY456GDAQ547R6NDAU74DTI2CKVVI4JERTUXKB2R4",
		assetCode:           "FOO",
		auroraClient:       auroraclient.DefaultTestNetClient,
		auroraURL:          "https://aurora-testnet.diamcircle.org/",
	}
	err = fh.validate()
	require.EqualError(t, err, "network passphrase cannot be empty")

	// missing payment amount
	fh = friendbotHandler{
		issuerAccountSecret: "SB6SFUY6ZJ2ETQHTY456GDAQ547R6NDAU74DTI2CKVVI4JERTUXKB2R4",
		assetCode:           "FOO",
		auroraClient:       auroraclient.DefaultTestNetClient,
		auroraURL:          "https://aurora-testnet.diamcircle.org/",
		networkPassphrase:   network.TestNetworkPassphrase,
	}
	err = fh.validate()
	require.EqualError(t, err, "payment amount must be greater than zero")

	// success!
	fh = friendbotHandler{
		issuerAccountSecret: "SB6SFUY6ZJ2ETQHTY456GDAQ547R6NDAU74DTI2CKVVI4JERTUXKB2R4",
		assetCode:           "FOO",
		auroraClient:       auroraclient.DefaultTestNetClient,
		auroraURL:          "https://aurora-testnet.diamcircle.org/",
		networkPassphrase:   network.TestNetworkPassphrase,
		paymentAmount:       1,
	}
	err = fh.validate()
	require.NoError(t, err)
}

func TestFriendbotHandler_serveHTTP_missingAddress(t *testing.T) {
	ctx := context.Background()

	handler := friendbotHandler{
		issuerAccountSecret: "SB6SFUY6ZJ2ETQHTY456GDAQ547R6NDAU74DTI2CKVVI4JERTUXKB2R4",
		assetCode:           "FOO",
		auroraClient:       auroraclient.DefaultTestNetClient,
		auroraURL:          "https://aurora-testnet.diamcircle.org/",
		networkPassphrase:   network.TestNetworkPassphrase,
		paymentAmount:       10000,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/friendbot", nil)
	r = r.WithContext(ctx)
	m := chi.NewMux()
	m.Get("/friendbot", handler.ServeHTTP)
	m.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	wantBody := `{
		"error":"Missing query paramater \"addr\"."
	}`
	require.JSONEq(t, wantBody, string(body))
}

func TestFriendbotHandler_serveHTTP_invalidAddress(t *testing.T) {
	ctx := context.Background()

	handler := friendbotHandler{
		issuerAccountSecret: "SB6SFUY6ZJ2ETQHTY456GDAQ547R6NDAU74DTI2CKVVI4JERTUXKB2R4",
		assetCode:           "FOO",
		auroraClient:       auroraclient.DefaultTestNetClient,
		auroraURL:          "https://aurora-testnet.diamcircle.org/",
		networkPassphrase:   network.TestNetworkPassphrase,
		paymentAmount:       10000,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/friendbot?addr=123", nil)
	r = r.WithContext(ctx)
	m := chi.NewMux()
	m.Get("/friendbot", handler.ServeHTTP)
	m.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	wantBody := `{
		"error":"\"addr\" is not a valid Diamcircle address."
	}`
	require.JSONEq(t, wantBody, string(body))
}

func TestFriendbotHandler_serveHTTP_accountDoesntExist(t *testing.T) {
	ctx := context.Background()

	auroraMock := auroraclient.MockClient{}
	auroraMock.
		On("AccountDetail", auroraclient.AccountRequest{AccountID: "GA2ILZPZAQ4R5PRKZ2X2AFAZK3ND6AGA4VFBQGR66BH36PV3VKMWLLZP"}).
		Return(aurora.Account{}, errors.New("something went wrong")) // account doesn't exist on ledger

	handler := friendbotHandler{
		issuerAccountSecret: "SB6SFUY6ZJ2ETQHTY456GDAQ547R6NDAU74DTI2CKVVI4JERTUXKB2R4",
		assetCode:           "FOO",
		auroraClient:       &auroraMock,
		auroraURL:          "https://aurora-testnet.diamcircle.org/",
		networkPassphrase:   network.TestNetworkPassphrase,
		paymentAmount:       10000,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/friendbot?addr=GA2ILZPZAQ4R5PRKZ2X2AFAZK3ND6AGA4VFBQGR66BH36PV3VKMWLLZP", nil)
	r = r.WithContext(ctx)
	m := chi.NewMux()
	m.Get("/friendbot", handler.ServeHTTP)
	m.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	wantBody := `{
		"error":"Please make sure the provided account address already exists in the network."
	}`
	require.JSONEq(t, wantBody, string(body))
}

func TestFriendbotHandler_serveHTTP_missingTrustline(t *testing.T) {
	ctx := context.Background()

	auroraMock := auroraclient.MockClient{}
	auroraMock.
		On("AccountDetail", auroraclient.AccountRequest{AccountID: "GA2ILZPZAQ4R5PRKZ2X2AFAZK3ND6AGA4VFBQGR66BH36PV3VKMWLLZP"}).
		Return(aurora.Account{}, nil)

	handler := friendbotHandler{
		issuerAccountSecret: "SB6SFUY6ZJ2ETQHTY456GDAQ547R6NDAU74DTI2CKVVI4JERTUXKB2R4",
		assetCode:           "FOO",
		auroraClient:       &auroraMock,
		auroraURL:          "https://aurora-testnet.diamcircle.org/",
		networkPassphrase:   network.TestNetworkPassphrase,
		paymentAmount:       10000,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/friendbot?addr=GA2ILZPZAQ4R5PRKZ2X2AFAZK3ND6AGA4VFBQGR66BH36PV3VKMWLLZP", nil)
	r = r.WithContext(ctx)
	m := chi.NewMux()
	m.Get("/friendbot", handler.ServeHTTP)
	m.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	wantBody := `{
		"error":"Account with address GA2ILZPZAQ4R5PRKZ2X2AFAZK3ND6AGA4VFBQGR66BH36PV3VKMWLLZP doesn't have a trustline for FOO:GDCRZMSHZGQYSRXPWDMIUNUQW36SV2NIC3C7R6WWT6XDO267WCI2TTBR"
	}`
	require.JSONEq(t, wantBody, string(body))
}

func TestFriendbotHandler_serveHTTP_issuerAccountDoesntExist(t *testing.T) {
	ctx := context.Background()

	// declare a logging buffer to validate output logs
	buf := new(strings.Builder)
	log.DefaultLogger.SetOutput(buf)
	log.DefaultLogger.SetLevel(log.InfoLevel)

	auroraMock := auroraclient.MockClient{}
	auroraMock.
		On("AccountDetail", auroraclient.AccountRequest{AccountID: "GA2ILZPZAQ4R5PRKZ2X2AFAZK3ND6AGA4VFBQGR66BH36PV3VKMWLLZP"}).
		Return(aurora.Account{
			Balances: []aurora.Balance{
				{
					Asset:   base.Asset{Code: "FOO", Issuer: "GDDIO6SFRD4SJEQFJOSKPIDYTDM7LM4METFBKN4NFGVR5DTGB7H75N5S"},
					Balance: "0",
				},
			},
		}, nil)
	auroraMock.
		On("AccountDetail", auroraclient.AccountRequest{AccountID: "GDDIO6SFRD4SJEQFJOSKPIDYTDM7LM4METFBKN4NFGVR5DTGB7H75N5S"}).
		Return(aurora.Account{}, errors.New("account doesn't exist")) // issuer account doesn't exist on ledger

	handler := friendbotHandler{
		issuerAccountSecret: "SDVFEIZ3WH5F6GHGK56QITTC5IO6QJ2UIQDWCHE72DAFZFSXYPIHQ6EV", // GDDIO6SFRD4SJEQFJOSKPIDYTDM7LM4METFBKN4NFGVR5DTGB7H75N5S
		assetCode:           "FOO",
		auroraClient:       &auroraMock,
		auroraURL:          "https://aurora-testnet.diamcircle.org/",
		networkPassphrase:   network.TestNetworkPassphrase,
		paymentAmount:       10000,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/friendbot?addr=GA2ILZPZAQ4R5PRKZ2X2AFAZK3ND6AGA4VFBQGR66BH36PV3VKMWLLZP", nil)
	r = r.WithContext(ctx)
	m := chi.NewMux()
	m.Get("/friendbot", handler.ServeHTTP)
	m.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	wantBody := `{
		"error":"An error occurred while processing this request."
	}`
	require.JSONEq(t, wantBody, string(body))
	require.Contains(t, buf.String(), "getting detail for issuer account")
}

func TestFriendbotHandler_serveHTTP(t *testing.T) {
	ctx := context.Background()

	auroraMock := auroraclient.MockClient{}
	auroraMock.
		On("AccountDetail", auroraclient.AccountRequest{AccountID: "GA2ILZPZAQ4R5PRKZ2X2AFAZK3ND6AGA4VFBQGR66BH36PV3VKMWLLZP"}).
		Return(aurora.Account{
			Balances: []aurora.Balance{
				{
					Asset:   base.Asset{Code: "FOO", Issuer: "GDDIO6SFRD4SJEQFJOSKPIDYTDM7LM4METFBKN4NFGVR5DTGB7H75N5S"},
					Balance: "0",
				},
			},
		}, nil)
	auroraMock.
		On("AccountDetail", auroraclient.AccountRequest{AccountID: "GDDIO6SFRD4SJEQFJOSKPIDYTDM7LM4METFBKN4NFGVR5DTGB7H75N5S"}).
		Return(aurora.Account{
			AccountID: "GDDIO6SFRD4SJEQFJOSKPIDYTDM7LM4METFBKN4NFGVR5DTGB7H75N5S",
			Sequence:  "1",
		}, nil)
	auroraMock.
		On("SubmitTransaction", mock.AnythingOfType("*txnbuild.Transaction")).
		Return(aurora.Transaction{}, nil)

	handler := friendbotHandler{
		issuerAccountSecret: "SDVFEIZ3WH5F6GHGK56QITTC5IO6QJ2UIQDWCHE72DAFZFSXYPIHQ6EV", // GDDIO6SFRD4SJEQFJOSKPIDYTDM7LM4METFBKN4NFGVR5DTGB7H75N5S
		assetCode:           "FOO",
		auroraClient:       &auroraMock,
		auroraURL:          "https://aurora-testnet.diamcircle.org/",
		networkPassphrase:   network.TestNetworkPassphrase,
		paymentAmount:       10000,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/friendbot?addr=GA2ILZPZAQ4R5PRKZ2X2AFAZK3ND6AGA4VFBQGR66BH36PV3VKMWLLZP", nil)
	r = r.WithContext(ctx)
	m := chi.NewMux()
	m.Get("/friendbot", handler.ServeHTTP)
	m.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	wantBody := `{
		"message":"ok"
	}`
	require.JSONEq(t, wantBody, string(body))
}
