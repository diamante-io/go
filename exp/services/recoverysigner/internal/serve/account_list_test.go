package serve

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"go/exp/services/recoverysigner/internal/account"
	"go/exp/services/recoverysigner/internal/db/dbtest"
	"go/exp/services/recoverysigner/internal/serve/auth"
	"go/keypair"
	supportlog "go/support/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test that when authenticated with an account, but no matching accounts,
// empty list is returned.
func TestAccountList_authenticatedButNonePermitted(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
	})
	s.Add(account.Account{
		Address: "GDU2CH4V3QYQB2BLMX45XQLVBEKSIN2EZLP37I6MZZ7NAR5U3TLZDQEY",
		Identities: []account.Identity{
			{
				Role: "sender",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypeAddress, Value: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N"},
				},
			},
		},
	})
	s.Add(account.Account{
		Address: "GCS4CVAAX7MVUNHP24655TNHIJ4YFN7GW5V3RFDC2BXVVMVDTB3GYH5U",
		Identities: []account.Identity{
			{
				Role: "receiver",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypeAddress, Value: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N"},
				},
			},
		},
	})
	h := accountListHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningAddresses: []*keypair.FromAddress{
			keypair.MustParseAddress("GCAPXRXSU7P6D353YGXMP6ROJIC744HO5OZCIWTXZQK2X757YU5KCHUE"),
			keypair.MustParseAddress("GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS"),
		},
	}

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Address: "GCNPATZQVSFGGSAHR4T54WNELPHYEBTSKH4IIKUTC7CHPLG6EPPC4PJL"})
	r := httptest.NewRequest("", "/", nil)
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
	"accounts": []
}`
	assert.JSONEq(t, wantBody, string(body))
}

func TestAccountList_authenticatedByPhoneNumber(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
		Identities: []account.Identity{
			{
				Role: "sender",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypePhoneNumber, Value: "+10000000000"},
				},
			},
		},
	})
	s.Add(account.Account{
		Address: "GDU2CH4V3QYQB2BLMX45XQLVBEKSIN2EZLP37I6MZZ7NAR5U3TLZDQEY",
		Identities: []account.Identity{
			{
				Role: "receiver",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypePhoneNumber, Value: "+10000000000"},
				},
			},
		},
	})
	s.Add(account.Account{
		Address: "GCS4CVAAX7MVUNHP24655TNHIJ4YFN7GW5V3RFDC2BXVVMVDTB3GYH5U",
		Identities: []account.Identity{
			{
				Role: "receiver",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypePhoneNumber, Value: "+20000000000"},
				},
			},
		},
	})
	h := accountListHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningAddresses: []*keypair.FromAddress{
			keypair.MustParseAddress("GCAPXRXSU7P6D353YGXMP6ROJIC744HO5OZCIWTXZQK2X757YU5KCHUE"),
			keypair.MustParseAddress("GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS"),
		},
	}

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{PhoneNumber: "+10000000000"})
	r := httptest.NewRequest("", "/", nil)
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
	"accounts": [
		{
			"address": "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
			"identities": [
				{ "role": "sender", "authenticated": true }
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
		},
		{
			"address": "GDU2CH4V3QYQB2BLMX45XQLVBEKSIN2EZLP37I6MZZ7NAR5U3TLZDQEY",
			"identities": [
				{ "role": "receiver", "authenticated": true }
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
		}
	]
}`
	assert.JSONEq(t, wantBody, string(body))
}

func TestAccountList_authenticatedByEmail(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
		Identities: []account.Identity{
			{
				Role: "sender",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypeEmail, Value: "user1@example.com"},
				},
			},
		},
	})
	s.Add(account.Account{
		Address: "GDU2CH4V3QYQB2BLMX45XQLVBEKSIN2EZLP37I6MZZ7NAR5U3TLZDQEY",
		Identities: []account.Identity{
			{
				Role: "receiver",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypeEmail, Value: "user1@example.com"},
				},
			},
		},
	})
	s.Add(account.Account{
		Address: "GCS4CVAAX7MVUNHP24655TNHIJ4YFN7GW5V3RFDC2BXVVMVDTB3GYH5U",
		Identities: []account.Identity{
			{
				Role: "receiver",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypeEmail, Value: "user2@example.com"},
				},
			},
		},
	})
	h := accountListHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningAddresses: []*keypair.FromAddress{
			keypair.MustParseAddress("GCAPXRXSU7P6D353YGXMP6ROJIC744HO5OZCIWTXZQK2X757YU5KCHUE"),
			keypair.MustParseAddress("GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS"),
		},
	}

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Email: "user1@example.com"})
	r := httptest.NewRequest("", "/", nil)
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
	"accounts": [
		{
			"address": "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
			"identities": [
				{ "role": "sender", "authenticated": true }
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
		},
		{
			"address": "GDU2CH4V3QYQB2BLMX45XQLVBEKSIN2EZLP37I6MZZ7NAR5U3TLZDQEY",
			"identities": [
				{ "role": "receiver", "authenticated": true }
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
		}
	]
}`
	assert.JSONEq(t, wantBody, string(body))
}

func TestAccountList_notAuthenticated(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GDIXCQJ2W2N6TAS6AYW4LW2EBV7XNRUCLNHQB37FARDEWBQXRWP47Q6N",
		Identities: []account.Identity{
			{
				Role: "sender",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypeEmail, Value: "user1@example.com"},
				},
			},
		},
	})
	s.Add(account.Account{
		Address: "GDU2CH4V3QYQB2BLMX45XQLVBEKSIN2EZLP37I6MZZ7NAR5U3TLZDQEY",
		Identities: []account.Identity{
			{
				Role: "receiver",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypeEmail, Value: "user1@example.com"},
				},
			},
		},
	})
	s.Add(account.Account{
		Address: "GCS4CVAAX7MVUNHP24655TNHIJ4YFN7GW5V3RFDC2BXVVMVDTB3GYH5U",
		Identities: []account.Identity{
			{
				Role: "receiver",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypeEmail, Value: "user2@example.com"},
				},
			},
		},
	})
	h := accountListHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningAddresses: []*keypair.FromAddress{
			keypair.MustParseAddress("GCAPXRXSU7P6D353YGXMP6ROJIC744HO5OZCIWTXZQK2X757YU5KCHUE"),
			keypair.MustParseAddress("GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS"),
		},
	}

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{})
	r := httptest.NewRequest("", "/", nil)
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
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
