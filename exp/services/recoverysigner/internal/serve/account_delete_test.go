package serve

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/diamcircle/go/exp/services/recoverysigner/internal/account"
	"github.com/diamcircle/go/exp/services/recoverysigner/internal/db/dbtest"
	"github.com/diamcircle/go/exp/services/recoverysigner/internal/serve/auth"
	"github.com/diamcircle/go/keypair"
	supportlog "github.com/diamcircle/go/support/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test that when authenticated with an account, but deleting the wrong account,
// an error is returned.
func TestAccountDelete_authenticatedNotAuthorized(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
		Identities: []account.Identity{
			{
				Role: "sender",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypeAddress, Value: "GCGZ3CNBE47IWAA5YIKDZL2XYYLA2UKJPS55P5EJ4VOMLK523PF3G7EM"},
				},
			},
		},
	})
	h := accountDeleteHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningAddresses: []*keypair.FromAddress{
			keypair.MustParseAddress("GCAPXRXSU7P6D353YGXMP6ROJIC744HO5OZCIWTXZQK2X757YU5KCHUE"),
			keypair.MustParseAddress("GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS"),
		},
	}

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Address: "GCNPATZQVSFGGSAHR4T54WNELPHYEBTSKH4IIKUTC7CHPLG6EPPC4PJL"})
	r := httptest.NewRequest("DELETE", "/GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N", nil)
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Delete("/{address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
	"error": "The resource at the url requested was not found."
}`

	assert.JSONEq(t, wantBody, string(body))
}

func TestAccountDelete_notAuthenticated(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
		Identities: []account.Identity{
			{
				Role: "sender",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypeAddress, Value: "GCGZ3CNBE47IWAA5YIKDZL2XYYLA2UKJPS55P5EJ4VOMLK523PF3G7EM"},
				},
			},
		},
	})
	h := accountDeleteHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningAddresses: []*keypair.FromAddress{
			keypair.MustParseAddress("GCAPXRXSU7P6D353YGXMP6ROJIC744HO5OZCIWTXZQK2X757YU5KCHUE"),
			keypair.MustParseAddress("GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS"),
		},
	}

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{})
	r := httptest.NewRequest("DELETE", "/GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N", nil)
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Delete("/{address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
	"error": "The request could not be authenticated."
}`

	assert.JSONEq(t, wantBody, string(body))
}

func TestAccountDelete_notFound(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
		Identities: []account.Identity{
			{
				Role: "sender",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypeAddress, Value: "GCGZ3CNBE47IWAA5YIKDZL2XYYLA2UKJPS55P5EJ4VOMLK523PF3G7EM"},
				},
			},
		},
	})
	h := accountDeleteHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningAddresses: []*keypair.FromAddress{
			keypair.MustParseAddress("GCAPXRXSU7P6D353YGXMP6ROJIC744HO5OZCIWTXZQK2X757YU5KCHUE"),
			keypair.MustParseAddress("GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS"),
		},
	}

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Address: "GCGZ3CNBE47IWAA5YIKDZL2XYYLA2UKJPS55P5EJ4VOMLK523PF3G7EM"})
	r := httptest.NewRequest("DELETE", "/GCGZ3CNBE47IWAA5YIKDZL2XYYLA2UKJPS55P5EJ4VOMLK523PF3G7EM", nil)
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Delete("/{address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
	"error": "The resource at the url requested was not found."
}`

	assert.JSONEq(t, wantBody, string(body))
}

func TestAccountDelete_authenticatedByIdentityAddress(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
		Identities: []account.Identity{
			{
				Role: "sender",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypeAddress, Value: "GCGZ3CNBE47IWAA5YIKDZL2XYYLA2UKJPS55P5EJ4VOMLK523PF3G7EM"},
				},
			},
			{
				Role: "receiver",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypePhoneNumber, Value: "+10000000000"},
				},
			},
		},
	})
	h := accountDeleteHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningAddresses: []*keypair.FromAddress{
			keypair.MustParseAddress("GCAPXRXSU7P6D353YGXMP6ROJIC744HO5OZCIWTXZQK2X757YU5KCHUE"),
			keypair.MustParseAddress("GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS"),
		},
	}

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Address: "GCGZ3CNBE47IWAA5YIKDZL2XYYLA2UKJPS55P5EJ4VOMLK523PF3G7EM"})
	r := httptest.NewRequest("DELETE", "/GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N", nil)
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Delete("/{address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
	"address": "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
	"identities": [
		{
			"role": "sender",
			"authenticated": true
		},
		{
			"role": "receiver"
		}
	],
	"signers": [
		{
			"key": "GCAPXRXSU7P6D353YGXMP6ROJIC744HO5OZCIWTXZQK2X757YU5KCHUE",
			"added_at": "0001-01-01T00:00:00Z"
		},
		{
			"key": "GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS",
			"added_at": "0001-01-01T00:00:00Z"
		}
	]
}`

	assert.JSONEq(t, wantBody, string(body))

	_, err = s.Get("GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N")
	require.Error(t, err)
	assert.Equal(t, account.ErrNotFound, err)
}

func TestAccountDelete_authenticatedByAccountAddress(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
		Identities: []account.Identity{
			{
				Role: "sender",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypeAddress, Value: "GCGZ3CNBE47IWAA5YIKDZL2XYYLA2UKJPS55P5EJ4VOMLK523PF3G7EM"},
				},
			},
			{
				Role: "receiver",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypePhoneNumber, Value: "+10000000000"},
				},
			},
		},
	})
	h := accountDeleteHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningAddresses: []*keypair.FromAddress{
			keypair.MustParseAddress("GCAPXRXSU7P6D353YGXMP6ROJIC744HO5OZCIWTXZQK2X757YU5KCHUE"),
			keypair.MustParseAddress("GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS"),
		},
	}

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N"})
	r := httptest.NewRequest("DELETE", "/GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N", nil)
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Delete("/{address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
	"address": "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
	"identities": [
		{
			"role": "sender"
		},
		{
			"role": "receiver"
		}
	],
	"signers": [
		{
			"key": "GCAPXRXSU7P6D353YGXMP6ROJIC744HO5OZCIWTXZQK2X757YU5KCHUE",
			"added_at": "0001-01-01T00:00:00Z"
		},
		{
			"key": "GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS",
			"added_at": "0001-01-01T00:00:00Z"
		}
	]
}`

	assert.JSONEq(t, wantBody, string(body))

	_, err = s.Get("GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N")
	require.Error(t, err)
	assert.Equal(t, account.ErrNotFound, err)
}

func TestAccountDelete_authenticatedByPhoneNumber(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
		Identities: []account.Identity{
			{
				Role: "sender",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypeAddress, Value: "GCGZ3CNBE47IWAA5YIKDZL2XYYLA2UKJPS55P5EJ4VOMLK523PF3G7EM"},
				},
			},
			{
				Role: "receiver",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypePhoneNumber, Value: "+10000000000"},
				},
			},
		},
	})
	h := accountDeleteHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningAddresses: []*keypair.FromAddress{
			keypair.MustParseAddress("GCAPXRXSU7P6D353YGXMP6ROJIC744HO5OZCIWTXZQK2X757YU5KCHUE"),
			keypair.MustParseAddress("GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS"),
		},
	}

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{PhoneNumber: "+10000000000"})
	r := httptest.NewRequest("DELETE", "/GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N", nil)
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Delete("/{address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
	"address": "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
	"identities": [
		{
			"role": "sender"
		},
		{
			"role": "receiver",
			"authenticated": true
		}
	],
	"signers": [
		{
			"key": "GCAPXRXSU7P6D353YGXMP6ROJIC744HO5OZCIWTXZQK2X757YU5KCHUE",
			"added_at": "0001-01-01T00:00:00Z"
		},
		{
			"key": "GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS",
			"added_at": "0001-01-01T00:00:00Z"
		}
	]
}`

	assert.JSONEq(t, wantBody, string(body))

	_, err = s.Get("GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N")
	require.Error(t, err)
	assert.Equal(t, account.ErrNotFound, err)
}

func TestAccountDelete_authenticatedByEmail(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
		Identities: []account.Identity{
			{
				Role: "sender",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypeAddress, Value: "GCGZ3CNBE47IWAA5YIKDZL2XYYLA2UKJPS55P5EJ4VOMLK523PF3G7EM"},
				},
			},
			{
				Role: "receiver",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypeEmail, Value: "user1@example.com"},
				},
			},
		},
	})
	h := accountDeleteHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningAddresses: []*keypair.FromAddress{
			keypair.MustParseAddress("GCAPXRXSU7P6D353YGXMP6ROJIC744HO5OZCIWTXZQK2X757YU5KCHUE"),
			keypair.MustParseAddress("GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS"),
		},
	}

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Email: "user1@example.com"})
	r := httptest.NewRequest("DELETE", "/GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N", nil)
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Delete("/{address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
	"address": "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
	"identities": [
		{
			"role": "sender"
		},
		{
			"role": "receiver",
			"authenticated": true
		}
	],
	"signers": [
		{
			"key": "GCAPXRXSU7P6D353YGXMP6ROJIC744HO5OZCIWTXZQK2X757YU5KCHUE",
			"added_at": "0001-01-01T00:00:00Z"
		},
		{
			"key": "GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS",
			"added_at": "0001-01-01T00:00:00Z"
		}
	]
}`

	assert.JSONEq(t, wantBody, string(body))

	_, err = s.Get("GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N")
	require.Error(t, err)
	assert.Equal(t, account.ErrNotFound, err)
}
