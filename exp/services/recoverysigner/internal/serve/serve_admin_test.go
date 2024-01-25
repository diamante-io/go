package serve

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"go/exp/services/recoverysigner/internal/account"
	"go/exp/services/recoverysigner/internal/db/dbtest"
	supportlog "go/support/log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminHandler_metricsAccountsCountNone(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}

	mr := prometheus.NewRegistry()
	mr.MustRegister(metricAccountsCount{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
	}.NewCollector())

	deps := adminDeps{
		Logger:          supportlog.DefaultLogger,
		MetricsGatherer: mr,
	}
	h := adminHandler(deps)

	r := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `accounts_count 0`
	assert.Contains(t, string(body), wantBody)
}

func TestAdminHandler_metricsAccountsCountSome(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address:    "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
		Identities: []account.Identity{{Role: "owner", AuthMethods: []account.AuthMethod{{Type: account.AuthMethodTypeAddress, Value: "GCGZ3CNBE47IWAA5YIKDZL2XYYLA2UKJPS55P5EJ4VOMLK523PF3G7EM"}}}},
	})
	s.Add(account.Account{
		Address:    "GCGZ3CNBE47IWAA5YIKDZL2XYYLA2UKJPS55P5EJ4VOMLK523PF3G7EM",
		Identities: []account.Identity{{Role: "owner", AuthMethods: []account.AuthMethod{{Type: account.AuthMethodTypeAddress, Value: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N"}}}},
	})

	mr := prometheus.NewRegistry()
	mr.MustRegister(metricAccountsCount{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
	}.NewCollector())

	deps := adminDeps{
		Logger:          supportlog.DefaultLogger,
		MetricsGatherer: mr,
	}

	h := adminHandler(deps)

	r := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `accounts_count 2`
	assert.Contains(t, string(body), wantBody)
}

func TestAdminHandler_metricsAccountsCountSomeDeleted(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address:    "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
		Identities: []account.Identity{{Role: "owner", AuthMethods: []account.AuthMethod{{Type: account.AuthMethodTypeAddress, Value: "GCGZ3CNBE47IWAA5YIKDZL2XYYLA2UKJPS55P5EJ4VOMLK523PF3G7EM"}}}},
	})
	s.Add(account.Account{
		Address:    "GCGZ3CNBE47IWAA5YIKDZL2XYYLA2UKJPS55P5EJ4VOMLK523PF3G7EM",
		Identities: []account.Identity{{Role: "owner", AuthMethods: []account.AuthMethod{{Type: account.AuthMethodTypeAddress, Value: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N"}}}},
	})
	s.Delete("GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N")

	mr := prometheus.NewRegistry()
	mr.MustRegister(metricAccountsCount{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
	}.NewCollector())

	deps := adminDeps{
		Logger:          supportlog.DefaultLogger,
		MetricsGatherer: mr,
	}

	h := adminHandler(deps)

	r := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `accounts_count 1`
	assert.Contains(t, string(body), wantBody)
}
