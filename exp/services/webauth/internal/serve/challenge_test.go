package serve

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go/keypair"
	"go/network"
	supportlog "go/support/log"
	"go/xdr"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChallenge(t *testing.T) {
	serverKey := keypair.MustRandom()
	account := keypair.MustRandom()

	h := challengeHandler{
		Logger:             supportlog.DefaultLogger,
		NetworkPassphrase:  network.TestNetworkPassphrase,
		SigningKey:         serverKey,
		ChallengeExpiresIn: time.Minute,
		Domain:             "webauthdomain",
		HomeDomains:        []string{"testdomain"},
	}

	r := httptest.NewRequest("GET", "/?account="+account.Address(), nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	res := struct {
		Transaction       string `json:"transaction"`
		NetworkPassphrase string `json:"network_passphrase"`
	}{}
	err := json.NewDecoder(resp.Body).Decode(&res)
	require.NoError(t, err)

	var tx xdr.TransactionEnvelope
	err = xdr.SafeUnmarshalBase64(res.Transaction, &tx)
	require.NoError(t, err)

	assert.Len(t, tx.Signatures(), 1)
	sourceAccount := tx.SourceAccount().ToAccountId()
	assert.Equal(t, serverKey.Address(), sourceAccount.Address())
	assert.Equal(t, tx.SeqNum(), int64(0))
	assert.Equal(t, time.Unix(int64(tx.TimeBounds().MaxTime), 0).Sub(time.Unix(int64(tx.TimeBounds().MinTime), 0)), time.Minute)
	assert.Len(t, tx.Operations(), 2)
	op0SourceAccount := tx.Operations()[0].SourceAccount.ToAccountId()
	assert.Equal(t, account.Address(), op0SourceAccount.Address())
	assert.Equal(t, xdr.OperationTypeManageData, tx.Operations()[0].Body.Type)
	assert.Regexp(t, "^testdomain auth", tx.Operations()[0].Body.ManageDataOp.DataName)
	op1SourceAccount := tx.Operations()[1].SourceAccount.ToAccountId()
	assert.Equal(t, sourceAccount.Address(), op1SourceAccount.Address())
	assert.Equal(t, xdr.OperationTypeManageData, tx.Operations()[1].Body.Type)
	assert.Equal(t, "web_auth_domain", string(tx.Operations()[1].Body.ManageDataOp.DataName))
	assert.Equal(t, "webauthdomain", string(*tx.Operations()[1].Body.ManageDataOp.DataValue))

	hash, err := network.HashTransactionInEnvelope(tx, res.NetworkPassphrase)
	require.NoError(t, err)
	assert.NoError(t, serverKey.FromAddress().Verify(hash[:], tx.Signatures()[0].Signature))

	assert.Equal(t, network.TestNetworkPassphrase, res.NetworkPassphrase)
}

func TestChallenge_anotherHomeDomain(t *testing.T) {
	serverKey := keypair.MustRandom()
	account := keypair.MustRandom()
	anotherDomain := "anotherdomain"

	h := challengeHandler{
		Logger:             supportlog.DefaultLogger,
		NetworkPassphrase:  network.TestNetworkPassphrase,
		SigningKey:         serverKey,
		ChallengeExpiresIn: time.Minute,
		Domain:             "webauthdomain",
		HomeDomains:        []string{"testdomain", anotherDomain},
	}

	r := httptest.NewRequest("GET", fmt.Sprintf("/?account=%s&home_domain=%s", account.Address(), anotherDomain), nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	res := struct {
		Transaction       string `json:"transaction"`
		NetworkPassphrase string `json:"network_passphrase"`
	}{}
	err := json.NewDecoder(resp.Body).Decode(&res)
	require.NoError(t, err)

	var tx xdr.TransactionEnvelope
	err = xdr.SafeUnmarshalBase64(res.Transaction, &tx)
	require.NoError(t, err)

	assert.Len(t, tx.Signatures(), 1)
	sourceAccount := tx.SourceAccount().ToAccountId()
	assert.Equal(t, serverKey.Address(), sourceAccount.Address())
	assert.Equal(t, tx.SeqNum(), int64(0))
	assert.Equal(t, time.Unix(int64(tx.TimeBounds().MaxTime), 0).Sub(time.Unix(int64(tx.TimeBounds().MinTime), 0)), time.Minute)
	assert.Len(t, tx.Operations(), 2)
	op0SourceAccount := tx.Operations()[0].SourceAccount.ToAccountId()
	assert.Equal(t, account.Address(), op0SourceAccount.Address())
	assert.Equal(t, xdr.OperationTypeManageData, tx.Operations()[0].Body.Type)
	assert.Regexp(t, "^anotherdomain auth", tx.Operations()[0].Body.ManageDataOp.DataName)
	op1SourceAccount := tx.Operations()[1].SourceAccount.ToAccountId()
	assert.Equal(t, sourceAccount.Address(), op1SourceAccount.Address())
	assert.Equal(t, xdr.OperationTypeManageData, tx.Operations()[1].Body.Type)
	assert.Equal(t, "web_auth_domain", string(tx.Operations()[1].Body.ManageDataOp.DataName))
	assert.Equal(t, "webauthdomain", string(*tx.Operations()[1].Body.ManageDataOp.DataValue))

	hash, err := network.HashTransactionInEnvelope(tx, res.NetworkPassphrase)
	require.NoError(t, err)
	assert.NoError(t, serverKey.FromAddress().Verify(hash[:], tx.Signatures()[0].Signature))

	assert.Equal(t, network.TestNetworkPassphrase, res.NetworkPassphrase)
}

func TestChallenge_noAccount(t *testing.T) {
	h := challengeHandler{
		SigningKey: keypair.MustRandom(),
	}

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.JSONEq(t, `{"error":"The request was invalid in some way."}`, string(body))
}

func TestChallenge_invalidAccount(t *testing.T) {
	h := challengeHandler{
		SigningKey: keypair.MustRandom(),
	}

	r := httptest.NewRequest("GET", "/?account=GREATACCOUNT", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.JSONEq(t, `{"error":"The request was invalid in some way."}`, string(body))
}

func TestChallenge_invalidHomeDomain(t *testing.T) {
	account := keypair.MustRandom()
	anotherDomain := "anotherdomain"

	h := challengeHandler{
		SigningKey:  keypair.MustRandom(),
		HomeDomains: []string{"testdomain"},
	}

	r := httptest.NewRequest("GET", fmt.Sprintf("/?account=%s&home_domain=%s", account.Address(), anotherDomain), nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.JSONEq(t, `{"error":"The request was invalid in some way."}`, string(body))
}
