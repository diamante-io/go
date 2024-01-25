package serve

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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

func TestAccountPost_newWithRoleOwnerContentTypeJSON(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	h := accountPostHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningAddresses: []*keypair.FromAddress{
			keypair.MustParseAddress("GCAPXRXSU7P6D353YGXMP6ROJIC744HO5OZCIWTXZQK2X757YU5KCHUE"),
			keypair.MustParseAddress("GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS"),
		},
	}

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N"})
	req := `{
	"identities": [
		{
			"role": "owner",
			"auth_methods": [
				{ "type": "diamcircle_address", "value": "GBF3XFXGBGNQDN3HOSZ7NVRF6TJ2JOD5U6ELIWJOOEI6T5WKMQT2YSXQ" },
				{ "type": "phone_number", "value": "+10000000000" },
				{ "type": "email", "value": "user1@example.com" }
			]
		}
	]
}`
	r := httptest.NewRequest("POST", "/GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N", strings.NewReader(req))
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
	"address": "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
	"identities": [
		{ "role": "owner" }
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

	acc, err := s.Get("GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N")
	require.NoError(t, err)
	wantAcc := account.Account{
		Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
		Identities: []account.Identity{
			{
				Role: "owner",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypeAddress, Value: "GBF3XFXGBGNQDN3HOSZ7NVRF6TJ2JOD5U6ELIWJOOEI6T5WKMQT2YSXQ"},
					{Type: account.AuthMethodTypePhoneNumber, Value: "+10000000000"},
					{Type: account.AuthMethodTypeEmail, Value: "user1@example.com"},
				},
			},
		},
	}
	assert.Equal(t, wantAcc, acc)
}

func TestAccountPost_newWithRoleOwnerContentTypeForm(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	h := accountPostHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningAddresses: []*keypair.FromAddress{
			keypair.MustParseAddress("GCAPXRXSU7P6D353YGXMP6ROJIC744HO5OZCIWTXZQK2X757YU5KCHUE"),
			keypair.MustParseAddress("GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS"),
		},
	}

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N"})
	reqValues := url.Values{}
	reqValues.Set("identities.0.role", "owner")
	reqValues.Set("identities.0.auth_methods.0.type", "diamcircle_address")
	reqValues.Set("identities.0.auth_methods.0.value", "GBF3XFXGBGNQDN3HOSZ7NVRF6TJ2JOD5U6ELIWJOOEI6T5WKMQT2YSXQ")
	reqValues.Set("identities.0.auth_methods.1.type", "phone_number")
	reqValues.Set("identities.0.auth_methods.1.value", "+10000000000")
	reqValues.Set("identities.0.auth_methods.2.type", "email")
	reqValues.Set("identities.0.auth_methods.2.value", "user1@example.com")
	req := reqValues.Encode()
	t.Log("Request Body:", req)
	r := httptest.NewRequest("POST", "/GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N", strings.NewReader(req))
	r = r.WithContext(ctx)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
	"address": "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
	"identities": [
		{ "role": "owner" }
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

	acc, err := s.Get("GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N")
	require.NoError(t, err)
	wantAcc := account.Account{
		Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
		Identities: []account.Identity{
			{
				Role: "owner",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypeAddress, Value: "GBF3XFXGBGNQDN3HOSZ7NVRF6TJ2JOD5U6ELIWJOOEI6T5WKMQT2YSXQ"},
					{Type: account.AuthMethodTypePhoneNumber, Value: "+10000000000"},
					{Type: account.AuthMethodTypeEmail, Value: "user1@example.com"},
				},
			},
		},
	}
	assert.Equal(t, wantAcc, acc)
}

func TestAccountPost_newWithRolesSenderReceiverContentTypeJSON(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	h := accountPostHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningAddresses: []*keypair.FromAddress{
			keypair.MustParseAddress("GCAPXRXSU7P6D353YGXMP6ROJIC744HO5OZCIWTXZQK2X757YU5KCHUE"),
			keypair.MustParseAddress("GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS"),
		},
	}

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N"})
	req := `{
	"identities": [
		{
			"role": "sender",
			"auth_methods": [
				{ "type": "diamcircle_address", "value": "GBF3XFXGBGNQDN3HOSZ7NVRF6TJ2JOD5U6ELIWJOOEI6T5WKMQT2YSXQ" },
				{ "type": "phone_number", "value": "+10000000000" },
				{ "type": "email", "value": "user1@example.com" }
			]
		},
		{
			"role": "receiver",
			"auth_methods": [
				{ "type": "diamcircle_address", "value": "GB5VOTKJ3IPGIYQBJ6GVJMUVEAIYGXZUJE4WYLPBJSHOTKLZTETBYOBI" },
				{ "type": "phone_number", "value": "+20000000000" },
				{ "type": "email", "value": "user2@example.com" }
			]
		}
	]
}`
	r := httptest.NewRequest("POST", "/GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N", strings.NewReader(req))
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
	"address": "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
	"identities": [
		{ "role": "sender" },
		{ "role": "receiver" }
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

	acc, err := s.Get("GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N")
	require.NoError(t, err)
	wantAcc := account.Account{
		Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
		Identities: []account.Identity{
			{
				Role: "sender",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypeAddress, Value: "GBF3XFXGBGNQDN3HOSZ7NVRF6TJ2JOD5U6ELIWJOOEI6T5WKMQT2YSXQ"},
					{Type: account.AuthMethodTypePhoneNumber, Value: "+10000000000"},
					{Type: account.AuthMethodTypeEmail, Value: "user1@example.com"},
				},
			},
			{
				Role: "receiver",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypeAddress, Value: "GB5VOTKJ3IPGIYQBJ6GVJMUVEAIYGXZUJE4WYLPBJSHOTKLZTETBYOBI"},
					{Type: account.AuthMethodTypePhoneNumber, Value: "+20000000000"},
					{Type: account.AuthMethodTypeEmail, Value: "user2@example.com"},
				},
			},
		},
	}
	assert.Equal(t, wantAcc, acc)
}

func TestAccountPost_newWithRolesSenderReceiverContentTypeForm(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	h := accountPostHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningAddresses: []*keypair.FromAddress{
			keypair.MustParseAddress("GCAPXRXSU7P6D353YGXMP6ROJIC744HO5OZCIWTXZQK2X757YU5KCHUE"),
			keypair.MustParseAddress("GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS"),
		},
	}

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N"})
	reqValues := url.Values{}
	reqValues.Set("identities.0.role", "sender")
	reqValues.Set("identities.0.auth_methods.0.type", "diamcircle_address")
	reqValues.Set("identities.0.auth_methods.0.value", "GBF3XFXGBGNQDN3HOSZ7NVRF6TJ2JOD5U6ELIWJOOEI6T5WKMQT2YSXQ")
	reqValues.Set("identities.0.auth_methods.1.type", "phone_number")
	reqValues.Set("identities.0.auth_methods.1.value", "+10000000000")
	reqValues.Set("identities.0.auth_methods.2.type", "email")
	reqValues.Set("identities.0.auth_methods.2.value", "user1@example.com")
	reqValues.Set("identities.1.role", "receiver")
	reqValues.Set("identities.1.auth_methods.0.type", "diamcircle_address")
	reqValues.Set("identities.1.auth_methods.0.value", "GB5VOTKJ3IPGIYQBJ6GVJMUVEAIYGXZUJE4WYLPBJSHOTKLZTETBYOBI")
	reqValues.Set("identities.1.auth_methods.1.type", "phone_number")
	reqValues.Set("identities.1.auth_methods.1.value", "+20000000000")
	reqValues.Set("identities.1.auth_methods.2.type", "email")
	reqValues.Set("identities.1.auth_methods.2.value", "user2@example.com")
	req := reqValues.Encode()
	t.Log("Request Body:", req)
	r := httptest.NewRequest("POST", "/GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N", strings.NewReader(req))
	r = r.WithContext(ctx)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
	"address": "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
	"identities": [
		{ "role": "sender" },
		{ "role": "receiver" }
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

	acc, err := s.Get("GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N")
	require.NoError(t, err)
	wantAcc := account.Account{
		Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
		Identities: []account.Identity{
			{
				Role: "sender",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypeAddress, Value: "GBF3XFXGBGNQDN3HOSZ7NVRF6TJ2JOD5U6ELIWJOOEI6T5WKMQT2YSXQ"},
					{Type: account.AuthMethodTypePhoneNumber, Value: "+10000000000"},
					{Type: account.AuthMethodTypeEmail, Value: "user1@example.com"},
				},
			},
			{
				Role: "receiver",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypeAddress, Value: "GB5VOTKJ3IPGIYQBJ6GVJMUVEAIYGXZUJE4WYLPBJSHOTKLZTETBYOBI"},
					{Type: account.AuthMethodTypePhoneNumber, Value: "+20000000000"},
					{Type: account.AuthMethodTypeEmail, Value: "user2@example.com"},
				},
			},
		},
	}
	assert.Equal(t, wantAcc, acc)
}

func TestAccountPost_accountAddressInvalid(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
	})
	h := accountPostHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningAddresses: []*keypair.FromAddress{
			keypair.MustParseAddress("GCAPXRXSU7P6D353YGXMP6ROJIC744HO5OZCIWTXZQK2X757YU5KCHUE"),
			keypair.MustParseAddress("GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS"),
		},
	}

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N"})
	req := `{}`
	r := httptest.NewRequest("POST", "/ZDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N", strings.NewReader(req))
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
	"error": "The request was invalid in some way."
}`
	assert.JSONEq(t, wantBody, string(body))

	_, err = s.Get("ZDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N")
	assert.Equal(t, account.ErrNotFound, err)
}

func TestAccountPost_accountAlreadyExists(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
	})
	h := accountPostHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningAddresses: []*keypair.FromAddress{
			keypair.MustParseAddress("GCAPXRXSU7P6D353YGXMP6ROJIC744HO5OZCIWTXZQK2X757YU5KCHUE"),
			keypair.MustParseAddress("GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS"),
		},
	}

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N"})
	req := `{
	"identities": [
		{
			"role": "owner",
			"auth_methods": [
				{ "type": "phone_number", "value": "+10000000000" }
			]
		}
	]
}`
	r := httptest.NewRequest("POST", "/GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N", strings.NewReader(req))
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusConflict, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
	"error": "The request could not be completed because the resource already exists."
}`
	assert.JSONEq(t, wantBody, string(body))

	acc, err := s.Get("GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N")
	require.NoError(t, err)
	wantAcc := account.Account{
		Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
	}
	assert.Equal(t, wantAcc, acc)
}

func TestAccountPost_identitiesNotProvided(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	h := accountPostHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningAddresses: []*keypair.FromAddress{
			keypair.MustParseAddress("GCAPXRXSU7P6D353YGXMP6ROJIC744HO5OZCIWTXZQK2X757YU5KCHUE"),
			keypair.MustParseAddress("GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS"),
		},
	}

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N"})
	req := `{}`
	r := httptest.NewRequest("POST", "/GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N", strings.NewReader(req))
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
	"error": "The request was invalid in some way."
}`
	assert.JSONEq(t, wantBody, string(body))

	_, err = s.Get("GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N")
	assert.Equal(t, account.ErrNotFound, err)
}

func TestAccountPost_roleNotProvided(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	h := accountPostHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningAddresses: []*keypair.FromAddress{
			keypair.MustParseAddress("GCAPXRXSU7P6D353YGXMP6ROJIC744HO5OZCIWTXZQK2X757YU5KCHUE"),
			keypair.MustParseAddress("GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS"),
		},
	}

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N"})
	req := `{
	"identities": [
		{
			"auth_methods": [
				{ "type": "diamcircle_address", "value": "GBF3XFXGBGNQDN3HOSZ7NVRF6TJ2JOD5U6ELIWJOOEI6T5WKMQT2YSXQ" },
				{ "type": "phone_number", "value": "+10000000000" },
				{ "type": "email", "value": "user1@example.com" }
			]
		}
	]
}`
	r := httptest.NewRequest("POST", "/GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N", strings.NewReader(req))
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
	"error": "The request was invalid in some way."
}`
	assert.JSONEq(t, wantBody, string(body))

	_, err = s.Get("GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N")
	assert.Equal(t, account.ErrNotFound, err)
}

func TestAccountPost_authMethodsNotProvided(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	h := accountPostHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningAddresses: []*keypair.FromAddress{
			keypair.MustParseAddress("GCAPXRXSU7P6D353YGXMP6ROJIC744HO5OZCIWTXZQK2X757YU5KCHUE"),
			keypair.MustParseAddress("GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS"),
		},
	}

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N"})
	req := `{
	"identities": [
		{
			"role": "owner"
		}
	]
}`
	r := httptest.NewRequest("POST", "/GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N", strings.NewReader(req))
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
	"error": "The request was invalid in some way."
}`
	assert.JSONEq(t, wantBody, string(body))

	_, err = s.Get("GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N")
	assert.Equal(t, account.ErrNotFound, err)
}

func TestAccountPost_authMethodTypeUnrecognized(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
	})
	h := accountPostHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningAddresses: []*keypair.FromAddress{
			keypair.MustParseAddress("GCAPXRXSU7P6D353YGXMP6ROJIC744HO5OZCIWTXZQK2X757YU5KCHUE"),
			keypair.MustParseAddress("GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS"),
		},
	}

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N"})
	req := `{
	"identities": [
		{
			"role": "owner",
			"auth_methods": [
				{ "type": "diamcircle_address", "value": "GBF3XFXGBGNQDN3HOSZ7NVRF6TJ2JOD5U6ELIWJOOEI6T5WKMQT2YSXQ" },
				{ "type": "phone_number", "value": "+10000000000" },
				{ "type": "wormhole_technology", "value": "galaxy5.earth3.asdfuaiosufd" },
				{ "type": "email", "value": "user1@example.com" }
			]
		}
	]
}`
	r := httptest.NewRequest("POST", "/GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N", strings.NewReader(req))
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
	"error": "The request was invalid in some way."
}`
	assert.JSONEq(t, wantBody, string(body))

	acc, err := s.Get("GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N")
	require.NoError(t, err)
	wantAcc := account.Account{
		Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
	}
	assert.Equal(t, wantAcc, acc)
}

func TestAccountPost_notAuthenticatedForAccount(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	h := accountPostHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningAddresses: []*keypair.FromAddress{
			keypair.MustParseAddress("GCAPXRXSU7P6D353YGXMP6ROJIC744HO5OZCIWTXZQK2X757YU5KCHUE"),
			keypair.MustParseAddress("GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS"),
		},
	}

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N"})
	req := `{}`
	r := httptest.NewRequest("POST", "/GDUKTYDY3RDNTNOUFJ2GPL5PIZTMTRD5P2CT274SYH67Q5J3NYI7XKYB", strings.NewReader(req))
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
	"error": "The request could not be authenticated."
}`
	assert.JSONEq(t, wantBody, string(body))

	_, err = s.Get("GDUKTYDY3RDNTNOUFJ2GPL5PIZTMTRD5P2CT274SYH67Q5J3NYI7XKYB")
	assert.Equal(t, account.ErrNotFound, err)
}
