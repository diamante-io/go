package serve

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"go/exp/services/recoverysigner/internal/account"
	"go/exp/services/recoverysigner/internal/db/dbtest"
	"go/exp/services/recoverysigner/internal/serve/auth"
	"go/keypair"
	"go/network"
	supportlog "go/support/log"
	"go/txnbuild"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test that when the account does not exist it returns not found.
func TestAccountSign_signingAddressAuthenticatedButNotFound(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	h := accountSignHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningKeys: []*keypair.Full{
			keypair.MustParseFull("SBIB72S6JMTGJRC6LMKLC5XMHZ2IOHZSZH4SASTN47LECEEJ7QEB6EYK"), // GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H
			keypair.MustParseFull("SBJGZKZ7LU2FQNEFBUOBW4LHCA5BOZCABIJTR7BQIFWQ3P763ZW7MYDD"), // GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS
		},
		NetworkPassphrase: network.TestNetworkPassphrase,
	}

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4"})
	req := `{
	"transaction": ""
}`
	r := httptest.NewRequest("POST", "/GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4/sign/GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H", strings.NewReader(req))
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}/sign/{signing-address}", h.ServeHTTP)
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

// Test that when the account exists but the authenticated client does not have
// permission to access it returns not found.
func TestAccountSign_signingAddressAccountAuthenticatedButNotPermitted(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4",
	})
	s.Add(account.Account{
		Address: "GBLOP46WEVXWO5N75TDX7GXLYFQE3XLDT5NQ2VYIBEWWEMSZWR3AUISZ",
	})
	h := accountSignHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningKeys: []*keypair.Full{
			keypair.MustParseFull("SBIB72S6JMTGJRC6LMKLC5XMHZ2IOHZSZH4SASTN47LECEEJ7QEB6EYK"), // GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H
			keypair.MustParseFull("SBJGZKZ7LU2FQNEFBUOBW4LHCA5BOZCABIJTR7BQIFWQ3P763ZW7MYDD"), // GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS
		},
		NetworkPassphrase: network.TestNetworkPassphrase,
	}

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Address: "GBLOP46WEVXWO5N75TDX7GXLYFQE3XLDT5NQ2VYIBEWWEMSZWR3AUISZ"})
	req := `{
	"transaction": ""
}`
	r := httptest.NewRequest("POST", "/GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4/sign/GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H", strings.NewReader(req))
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}/sign/{signing-address}", h.ServeHTTP)
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

func TestAccountSign_signingAddressAccountAuthenticatedButInvalidAddress(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	h := accountSignHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningKeys: []*keypair.Full{
			keypair.MustParseFull("SBIB72S6JMTGJRC6LMKLC5XMHZ2IOHZSZH4SASTN47LECEEJ7QEB6EYK"), // GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H
			keypair.MustParseFull("SBJGZKZ7LU2FQNEFBUOBW4LHCA5BOZCABIJTR7BQIFWQ3P763ZW7MYDD"), // GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS
		},
		NetworkPassphrase: network.TestNetworkPassphrase,
	}

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4"})
	r := httptest.NewRequest("POST", "/ZA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4/sign/GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H", nil)
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}/sign/{signing-address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{"error": "The request was invalid in some way."}`
	assert.JSONEq(t, wantBody, string(body))
}

func TestAccountSign_signingAddressAccountAuthenticatedButEmptyAddress(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	h := accountSignHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningKeys: []*keypair.Full{
			keypair.MustParseFull("SBIB72S6JMTGJRC6LMKLC5XMHZ2IOHZSZH4SASTN47LECEEJ7QEB6EYK"), // GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H
			keypair.MustParseFull("SBJGZKZ7LU2FQNEFBUOBW4LHCA5BOZCABIJTR7BQIFWQ3P763ZW7MYDD"), // GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS
		},
		NetworkPassphrase: network.TestNetworkPassphrase,
	}

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4"})
	r := httptest.NewRequest("POST", "//sign/GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H", nil)
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}/sign/{signing-address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{"error": "The request was invalid in some way."}`
	assert.JSONEq(t, wantBody, string(body))
}

// Test that when the account exists but the authenticated client does not have
// permission to access it returns not found.
func TestAccountSign_signingAddressPhoneNumberAuthenticatedButNotPermitted(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4",
		Identities: []account.Identity{
			{
				Role: "sender",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypePhoneNumber, Value: "+10000000000"},
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
	s.Add(account.Account{
		Address: "GBLOP46WEVXWO5N75TDX7GXLYFQE3XLDT5NQ2VYIBEWWEMSZWR3AUISZ",
		Identities: []account.Identity{
			{
				Role: "sender",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypePhoneNumber, Value: "+20000000000"},
				},
			},
			{
				Role: "receiver",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypePhoneNumber, Value: "+20000000000"},
				},
			},
		},
	})
	h := accountSignHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningKeys: []*keypair.Full{
			keypair.MustParseFull("SBIB72S6JMTGJRC6LMKLC5XMHZ2IOHZSZH4SASTN47LECEEJ7QEB6EYK"), // GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H
			keypair.MustParseFull("SBJGZKZ7LU2FQNEFBUOBW4LHCA5BOZCABIJTR7BQIFWQ3P763ZW7MYDD"), // GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS
		},
		NetworkPassphrase: network.TestNetworkPassphrase,
	}

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{PhoneNumber: "+20000000000"})
	req := `{
	"transaction": ""
}`
	r := httptest.NewRequest("POST", "/GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4/sign/GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H", strings.NewReader(req))
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}/sign/{signing-address}", h.ServeHTTP)
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

// Test that when the account exists but the authenticated client does not have
// permission to access it returns not found.
func TestAccountSign_signingAddressEmailAuthenticatedButNotPermitted(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4",
		Identities: []account.Identity{
			{
				Role: "sender",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypeEmail, Value: "user1@example.com"},
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
	s.Add(account.Account{
		Address: "GBLOP46WEVXWO5N75TDX7GXLYFQE3XLDT5NQ2VYIBEWWEMSZWR3AUISZ",
		Identities: []account.Identity{
			{
				Role: "sender",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypeEmail, Value: "user2@example.com"},
				},
			},
			{
				Role: "receiver",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypeEmail, Value: "user2@example.com"},
				},
			},
		},
	})
	h := accountSignHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningKeys: []*keypair.Full{
			keypair.MustParseFull("SBIB72S6JMTGJRC6LMKLC5XMHZ2IOHZSZH4SASTN47LECEEJ7QEB6EYK"), // GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H
			keypair.MustParseFull("SBJGZKZ7LU2FQNEFBUOBW4LHCA5BOZCABIJTR7BQIFWQ3P763ZW7MYDD"), // GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS
		},
		NetworkPassphrase: network.TestNetworkPassphrase,
	}

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{PhoneNumber: "user2@example.com"})
	req := `{
	"transaction": ""
}`
	r := httptest.NewRequest("POST", "/GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4/sign/GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H", strings.NewReader(req))
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}/sign/{signing-address}", h.ServeHTTP)
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

// Test that when the signing address is not valid.
func TestAccountSign_signingAddressAccountAuthenticatedButSigningAddressInvalid(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4",
	})
	h := accountSignHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningKeys: []*keypair.Full{
			keypair.MustParseFull("SBIB72S6JMTGJRC6LMKLC5XMHZ2IOHZSZH4SASTN47LECEEJ7QEB6EYK"), // GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H
			keypair.MustParseFull("SBJGZKZ7LU2FQNEFBUOBW4LHCA5BOZCABIJTR7BQIFWQ3P763ZW7MYDD"), // GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS
		},
		NetworkPassphrase: network.TestNetworkPassphrase,
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &txnbuild.SimpleAccount{AccountID: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4"},
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.SetOptions{
					Signer: &txnbuild.Signer{
						Address: "GD7CGJSJ5OBOU5KOP2UQDH3MPY75UTEY27HVV5XPSL2X6DJ2VGTOSXEU",
						Weight:  20,
					},
				},
			},
			BaseFee:    txnbuild.MinBaseFee,
			Timebounds: txnbuild.NewTimebounds(0, 1),
		},
	)
	require.NoError(t, err)
	txEnc, err := tx.Base64()
	require.NoError(t, err)
	t.Log("Tx:", txEnc)

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4"})
	req := `{
	"transaction": "` + txEnc + `"
}`
	r := httptest.NewRequest("POST", "/GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4/sign/GCF7VHXRMFMYODQTW7J4PFMHITSEP7XVYIM6X2AKSD7EKV3PMV2PCYHE", strings.NewReader(req))
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}/sign/{signing-address}", h.ServeHTTP)
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

func TestAccountSign_signingAddressAccountAuthenticatedOtherSignerSelected(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4",
	})
	h := accountSignHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningKeys: []*keypair.Full{
			keypair.MustParseFull("SBIB72S6JMTGJRC6LMKLC5XMHZ2IOHZSZH4SASTN47LECEEJ7QEB6EYK"), // GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H
			keypair.MustParseFull("SBJGZKZ7LU2FQNEFBUOBW4LHCA5BOZCABIJTR7BQIFWQ3P763ZW7MYDD"), // GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS
		},
		NetworkPassphrase: network.TestNetworkPassphrase,
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &txnbuild.SimpleAccount{AccountID: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4"},
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.SetOptions{
					Signer: &txnbuild.Signer{
						Address: "GD7CGJSJ5OBOU5KOP2UQDH3MPY75UTEY27HVV5XPSL2X6DJ2VGTOSXEU",
						Weight:  20,
					},
				},
			},
			BaseFee:    txnbuild.MinBaseFee,
			Timebounds: txnbuild.NewTimebounds(0, 1),
		},
	)
	require.NoError(t, err)
	txEnc, err := tx.Base64()
	require.NoError(t, err)
	t.Log("Tx:", txEnc)

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4"})
	req := `{
	"transaction": "` + txEnc + `"
}`
	r := httptest.NewRequest("POST", "/GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4/sign/GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS", strings.NewReader(req))
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}/sign/{signing-address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
	"signature": "8Tew6rTol9me8H9u7ezJXg6SEqzOr7cwSf1y9+Vri275XEDH9xWtZ2klTX2uUSPy56otUoIySPqsV3dUFs2kDA==",
	"network_passphrase": "Test SDF Network ; September 2015"
}`
	assert.JSONEq(t, wantBody, string(body))
}

// Test that when the source account of the transaction matches the account the
// request is for, that the transaction is signed and a signature is returned.
// The operation source account does not need to be set.
func TestAccountSign_signingAddressAccountAuthenticatedTxSourceAccountValid(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4",
	})
	h := accountSignHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningKeys: []*keypair.Full{
			keypair.MustParseFull("SBIB72S6JMTGJRC6LMKLC5XMHZ2IOHZSZH4SASTN47LECEEJ7QEB6EYK"), // GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H
			keypair.MustParseFull("SBJGZKZ7LU2FQNEFBUOBW4LHCA5BOZCABIJTR7BQIFWQ3P763ZW7MYDD"), // GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS
		},
		NetworkPassphrase: network.TestNetworkPassphrase,
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &txnbuild.SimpleAccount{AccountID: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4"},
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.SetOptions{
					Signer: &txnbuild.Signer{
						Address: "GD7CGJSJ5OBOU5KOP2UQDH3MPY75UTEY27HVV5XPSL2X6DJ2VGTOSXEU",
						Weight:  20,
					},
				},
			},
			BaseFee:    txnbuild.MinBaseFee,
			Timebounds: txnbuild.NewTimebounds(0, 1),
		},
	)
	require.NoError(t, err)
	txEnc, err := tx.Base64()
	require.NoError(t, err)
	t.Log("Tx:", txEnc)

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4"})
	req := `{
	"transaction": "` + txEnc + `"
}`
	r := httptest.NewRequest("POST", "/GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4/sign/GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H", strings.NewReader(req))
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}/sign/{signing-address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
	"signature": "okp0ISR/hjU6ItsfXie6ArlQ3YWkBBqEAM5TJrthALdawV5DzcpuwBKi0QE/iBgoU7eY0hY3RPdxm8mXGNYfCQ==",
	"network_passphrase": "Test SDF Network ; September 2015"
}`
	assert.JSONEq(t, wantBody, string(body))
}

// Test that when the source account of the transaction and operation are both
// set to values that match the account the request is for, that the
// transaction is signed and a signature is returned.
func TestAccountSign_signingAddressAccountAuthenticatedTxAndOpSourceAccountValid(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4",
	})
	h := accountSignHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningKeys: []*keypair.Full{
			keypair.MustParseFull("SBIB72S6JMTGJRC6LMKLC5XMHZ2IOHZSZH4SASTN47LECEEJ7QEB6EYK"), // GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H
			keypair.MustParseFull("SBJGZKZ7LU2FQNEFBUOBW4LHCA5BOZCABIJTR7BQIFWQ3P763ZW7MYDD"), // GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS
		},
		NetworkPassphrase: network.TestNetworkPassphrase,
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &txnbuild.SimpleAccount{AccountID: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4"},
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.SetOptions{
					SourceAccount: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4",
					Signer: &txnbuild.Signer{
						Address: "GD7CGJSJ5OBOU5KOP2UQDH3MPY75UTEY27HVV5XPSL2X6DJ2VGTOSXEU",
						Weight:  20,
					},
				},
			},
			BaseFee:    txnbuild.MinBaseFee,
			Timebounds: txnbuild.NewTimebounds(0, 1),
		},
	)
	require.NoError(t, err)
	txEnc, err := tx.Base64()
	require.NoError(t, err)
	t.Log("Tx:", txEnc)

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4"})
	req := `{
	"transaction": "` + txEnc + `"
}`
	r := httptest.NewRequest("POST", "/GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4/sign/GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H", strings.NewReader(req))
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}/sign/{signing-address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
	"signature": "MKAkl+R3VT5DJw6Qed8jO8ERD4RcQ4dJlN+UR2n7nT6AVBXnKBk0zqBZnDuB153zfTYmuA5kmsRiNr5terHVBg==",
	"network_passphrase": "Test SDF Network ; September 2015"
}`
	assert.JSONEq(t, wantBody, string(body))
}

// Test that when the source account of the transaction is not the account sign
// the request is calling sign on a bad request response is returned.
func TestAccountSign_signingAddressAccountAuthenticatedTxSourceAccountInvalid(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4",
	})
	h := accountSignHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningKeys: []*keypair.Full{
			keypair.MustParseFull("SBIB72S6JMTGJRC6LMKLC5XMHZ2IOHZSZH4SASTN47LECEEJ7QEB6EYK"), // GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H
			keypair.MustParseFull("SBJGZKZ7LU2FQNEFBUOBW4LHCA5BOZCABIJTR7BQIFWQ3P763ZW7MYDD"), // GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS
		},
		NetworkPassphrase: network.TestNetworkPassphrase,
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &txnbuild.SimpleAccount{AccountID: "GA47G3ZQBUR5NF2ZECG774O3QGKFMAW75XLXSCDICFDDV5GKGRFGFSOI"},
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.SetOptions{
					SourceAccount: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4",
					Signer: &txnbuild.Signer{
						Address: "GD7CGJSJ5OBOU5KOP2UQDH3MPY75UTEY27HVV5XPSL2X6DJ2VGTOSXEU",
						Weight:  20,
					},
				},
			},
			BaseFee:    txnbuild.MinBaseFee,
			Timebounds: txnbuild.NewTimebounds(0, 1),
		},
	)
	require.NoError(t, err)
	txEnc, err := tx.Base64()
	require.NoError(t, err)
	t.Log("Tx:", txEnc)

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4"})
	req := `{
	"transaction": "` + txEnc + `"
}`
	r := httptest.NewRequest("POST", "/GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4/sign/GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H", strings.NewReader(req))
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}/sign/{signing-address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{"error": "The request was invalid in some way."}`
	assert.JSONEq(t, wantBody, string(body))
}

// Test that when the source account of the operation is not the account sign
// the request is calling sign on a bad request response is returned.
func TestAccountSign_signingAddressAccountAuthenticatedOpSourceAccountInvalid(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4",
	})
	h := accountSignHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningKeys: []*keypair.Full{
			keypair.MustParseFull("SBIB72S6JMTGJRC6LMKLC5XMHZ2IOHZSZH4SASTN47LECEEJ7QEB6EYK"), // GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H
			keypair.MustParseFull("SBJGZKZ7LU2FQNEFBUOBW4LHCA5BOZCABIJTR7BQIFWQ3P763ZW7MYDD"), // GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS
		},
		NetworkPassphrase: network.TestNetworkPassphrase,
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &txnbuild.SimpleAccount{AccountID: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4"},
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.SetOptions{
					SourceAccount: "GA47G3ZQBUR5NF2ZECG774O3QGKFMAW75XLXSCDICFDDV5GKGRFGFSOI",
					Signer: &txnbuild.Signer{
						Address: "GD7CGJSJ5OBOU5KOP2UQDH3MPY75UTEY27HVV5XPSL2X6DJ2VGTOSXEU",
						Weight:  20,
					},
				},
			},
			BaseFee:    txnbuild.MinBaseFee,
			Timebounds: txnbuild.NewTimebounds(0, 1),
		},
	)
	require.NoError(t, err)
	txEnc, err := tx.Base64()
	require.NoError(t, err)
	t.Log("Tx:", txEnc)

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4"})
	req := `{
	"transaction": "` + txEnc + `"
}`
	r := httptest.NewRequest("POST", "/GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4/sign/GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H", strings.NewReader(req))
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}/sign/{signing-address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{"error": "The request was invalid in some way."}`
	assert.JSONEq(t, wantBody, string(body))
}

// Test that when the source account of the operation and transaction is not
// the account sign the request is calling sign on a bad request response is
// returned.
func TestAccountSign_signingAddressAccountAuthenticatedTxAndOpSourceAccountInvalid(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4",
	})
	h := accountSignHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningKeys: []*keypair.Full{
			keypair.MustParseFull("SBIB72S6JMTGJRC6LMKLC5XMHZ2IOHZSZH4SASTN47LECEEJ7QEB6EYK"), // GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H
			keypair.MustParseFull("SBJGZKZ7LU2FQNEFBUOBW4LHCA5BOZCABIJTR7BQIFWQ3P763ZW7MYDD"), // GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS
		},
		NetworkPassphrase: network.TestNetworkPassphrase,
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &txnbuild.SimpleAccount{AccountID: "GA47G3ZQBUR5NF2ZECG774O3QGKFMAW75XLXSCDICFDDV5GKGRFGFSOI"},
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.SetOptions{
					SourceAccount: "GA47G3ZQBUR5NF2ZECG774O3QGKFMAW75XLXSCDICFDDV5GKGRFGFSOI",
					Signer: &txnbuild.Signer{
						Address: "GD7CGJSJ5OBOU5KOP2UQDH3MPY75UTEY27HVV5XPSL2X6DJ2VGTOSXEU",
						Weight:  20,
					},
				},
			},
			BaseFee:    txnbuild.MinBaseFee,
			Timebounds: txnbuild.NewTimebounds(0, 1),
		},
	)
	require.NoError(t, err)
	txEnc, err := tx.Base64()
	require.NoError(t, err)
	t.Log("Tx:", txEnc)

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4"})
	req := `{
	"transaction": "` + txEnc + `"
}`
	r := httptest.NewRequest("POST", "/GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4/sign/GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H", strings.NewReader(req))
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}/sign/{signing-address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{"error": "The request was invalid in some way."}`
	assert.JSONEq(t, wantBody, string(body))
}

// Test that when authenticated with a phone number signing is possible.
func TestAccountSign_signingAddressPhoneNumberOwnerAuthenticated(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4",
		Identities: []account.Identity{
			{
				Role: "sender",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypePhoneNumber, Value: "+10000000000"},
				},
			},
		},
	})
	h := accountSignHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningKeys: []*keypair.Full{
			keypair.MustParseFull("SBIB72S6JMTGJRC6LMKLC5XMHZ2IOHZSZH4SASTN47LECEEJ7QEB6EYK"), // GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H
			keypair.MustParseFull("SBJGZKZ7LU2FQNEFBUOBW4LHCA5BOZCABIJTR7BQIFWQ3P763ZW7MYDD"), // GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS
		},
		NetworkPassphrase: network.TestNetworkPassphrase,
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &txnbuild.SimpleAccount{AccountID: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4"},
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.SetOptions{
					Signer: &txnbuild.Signer{
						Address: "GD7CGJSJ5OBOU5KOP2UQDH3MPY75UTEY27HVV5XPSL2X6DJ2VGTOSXEU",
						Weight:  20,
					},
				},
			},
			BaseFee:    txnbuild.MinBaseFee,
			Timebounds: txnbuild.NewTimebounds(0, 1),
		},
	)
	require.NoError(t, err)
	txEnc, err := tx.Base64()
	require.NoError(t, err)
	t.Log("Tx:", txEnc)

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{PhoneNumber: "+10000000000"})
	req := `{
	"transaction": "` + txEnc + `"
}`
	r := httptest.NewRequest("POST", "/GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4/sign/GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H", strings.NewReader(req))
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}/sign/{signing-address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
	"signature": "okp0ISR/hjU6ItsfXie6ArlQ3YWkBBqEAM5TJrthALdawV5DzcpuwBKi0QE/iBgoU7eY0hY3RPdxm8mXGNYfCQ==",
	"network_passphrase": "Test SDF Network ; September 2015"
}`
	assert.JSONEq(t, wantBody, string(body))
}

// Test that when authenticated with a phone number signing is possible.
func TestAccountSign_signingAddressPhoneNumberOtherAuthenticated(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4",
		Identities: []account.Identity{
			{
				Role: "receiver",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypePhoneNumber, Value: "+10000000000"},
				},
			},
		},
	})
	h := accountSignHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningKeys: []*keypair.Full{
			keypair.MustParseFull("SBIB72S6JMTGJRC6LMKLC5XMHZ2IOHZSZH4SASTN47LECEEJ7QEB6EYK"), // GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H
			keypair.MustParseFull("SBJGZKZ7LU2FQNEFBUOBW4LHCA5BOZCABIJTR7BQIFWQ3P763ZW7MYDD"), // GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS
		},
		NetworkPassphrase: network.TestNetworkPassphrase,
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &txnbuild.SimpleAccount{AccountID: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4"},
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.SetOptions{
					Signer: &txnbuild.Signer{
						Address: "GD7CGJSJ5OBOU5KOP2UQDH3MPY75UTEY27HVV5XPSL2X6DJ2VGTOSXEU",
						Weight:  20,
					},
				},
			},
			BaseFee:    txnbuild.MinBaseFee,
			Timebounds: txnbuild.NewTimebounds(0, 1),
		},
	)
	require.NoError(t, err)
	txEnc, err := tx.Base64()
	require.NoError(t, err)
	t.Log("Tx:", txEnc)

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{PhoneNumber: "+10000000000"})
	req := `{
	"transaction": "` + txEnc + `"
}`
	r := httptest.NewRequest("POST", "/GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4/sign/GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H", strings.NewReader(req))
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}/sign/{signing-address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
	"signature": "okp0ISR/hjU6ItsfXie6ArlQ3YWkBBqEAM5TJrthALdawV5DzcpuwBKi0QE/iBgoU7eY0hY3RPdxm8mXGNYfCQ==",
	"network_passphrase": "Test SDF Network ; September 2015"
}`
	assert.JSONEq(t, wantBody, string(body))
}

// Test that when authenticated with a email signing is possible.
func TestAccountSign_signingAddressEmailOwnerAuthenticated(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4",
		Identities: []account.Identity{
			{
				Role: "sender",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypeEmail, Value: "user1@example.com"},
				},
			},
		},
	})
	h := accountSignHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningKeys: []*keypair.Full{
			keypair.MustParseFull("SBIB72S6JMTGJRC6LMKLC5XMHZ2IOHZSZH4SASTN47LECEEJ7QEB6EYK"), // GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H
			keypair.MustParseFull("SBJGZKZ7LU2FQNEFBUOBW4LHCA5BOZCABIJTR7BQIFWQ3P763ZW7MYDD"), // GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS
		},
		NetworkPassphrase: network.TestNetworkPassphrase,
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &txnbuild.SimpleAccount{AccountID: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4"},
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.SetOptions{
					Signer: &txnbuild.Signer{
						Address: "GD7CGJSJ5OBOU5KOP2UQDH3MPY75UTEY27HVV5XPSL2X6DJ2VGTOSXEU",
						Weight:  20,
					},
				},
			},
			BaseFee:    txnbuild.MinBaseFee,
			Timebounds: txnbuild.NewTimebounds(0, 1),
		},
	)
	require.NoError(t, err)
	txEnc, err := tx.Base64()
	require.NoError(t, err)
	t.Log("Tx:", txEnc)

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Email: "user1@example.com"})
	req := `{
	"transaction": "` + txEnc + `"
}`
	r := httptest.NewRequest("POST", "/GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4/sign/GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H", strings.NewReader(req))
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}/sign/{signing-address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
	"signature": "okp0ISR/hjU6ItsfXie6ArlQ3YWkBBqEAM5TJrthALdawV5DzcpuwBKi0QE/iBgoU7eY0hY3RPdxm8mXGNYfCQ==",
	"network_passphrase": "Test SDF Network ; September 2015"
}`
	assert.JSONEq(t, wantBody, string(body))
}

// Test that when the source account of the operation is listed in the allowed
// source accounts a successful response is returned.
func TestAccountSign_signingAddressEmailOwnerAuthenticatedOpSourceAccountIsAllowedSourceAccount(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4",
		Identities: []account.Identity{
			{
				Role: "sender",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypeEmail, Value: "user1@example.com"},
				},
			},
		},
	})
	h := accountSignHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningKeys: []*keypair.Full{
			keypair.MustParseFull("SBIB72S6JMTGJRC6LMKLC5XMHZ2IOHZSZH4SASTN47LECEEJ7QEB6EYK"), // GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H
			keypair.MustParseFull("SBJGZKZ7LU2FQNEFBUOBW4LHCA5BOZCABIJTR7BQIFWQ3P763ZW7MYDD"), // GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS
		},
		NetworkPassphrase: network.TestNetworkPassphrase,
		AllowedSourceAccounts: []*keypair.FromAddress{
			keypair.MustParseAddress("GDR3RJVOHYR5A4RSLZ7D3GOSTPBGD2FY7KJD7ZB7363ROOQHWYDVVULS"),
		},
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &txnbuild.SimpleAccount{AccountID: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4"},
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.BeginSponsoringFutureReserves{
					SponsoredID:   "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4",
					SourceAccount: "GDR3RJVOHYR5A4RSLZ7D3GOSTPBGD2FY7KJD7ZB7363ROOQHWYDVVULS",
				},
				&txnbuild.SetOptions{
					Signer: &txnbuild.Signer{
						Address: "GD7CGJSJ5OBOU5KOP2UQDH3MPY75UTEY27HVV5XPSL2X6DJ2VGTOSXEU",
						Weight:  20,
					},
				},
				&txnbuild.EndSponsoringFutureReserves{},
			},
			BaseFee:    txnbuild.MinBaseFee,
			Timebounds: txnbuild.NewTimebounds(0, 1),
		},
	)
	require.NoError(t, err)
	txEnc, err := tx.Base64()
	require.NoError(t, err)
	t.Log("Tx:", txEnc)

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Email: "user1@example.com"})
	req := fmt.Sprintf(`{"transaction": "%s"}`, txEnc)
	r := httptest.NewRequest("POST", "/GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4/sign/GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H", strings.NewReader(req))
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}/sign/{signing-address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
		"signature": "Tpl/yZoKkahakaX4fSrdIeBLL2oi4uKegs5bLXFj5fG6Rcfe2D4EeSHcjJmmO2ZscuY8pX8+YPo70AvCtfw9Ag==",
		"network_passphrase": "Test SDF Network ; September 2015"
	}`
	assert.JSONEq(t, wantBody, string(body))
}

// Test that when the source account of the operation is not the account sign
// the request is calling sign on a bad request response is returned.
func TestAccountSign_signingAddressEmailOwnerAuthenticatedOpSourceAccountInvalid(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4",
		Identities: []account.Identity{
			{
				Role: "sender",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypeEmail, Value: "user1@example.com"},
				},
			},
		},
	})
	h := accountSignHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningKeys: []*keypair.Full{
			keypair.MustParseFull("SBIB72S6JMTGJRC6LMKLC5XMHZ2IOHZSZH4SASTN47LECEEJ7QEB6EYK"), // GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H
			keypair.MustParseFull("SBJGZKZ7LU2FQNEFBUOBW4LHCA5BOZCABIJTR7BQIFWQ3P763ZW7MYDD"), // GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS
		},
		NetworkPassphrase:     network.TestNetworkPassphrase,
		AllowedSourceAccounts: nil,
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &txnbuild.SimpleAccount{AccountID: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4"},
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.BeginSponsoringFutureReserves{
					SponsoredID:   "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4",
					SourceAccount: "GDR3RJVOHYR5A4RSLZ7D3GOSTPBGD2FY7KJD7ZB7363ROOQHWYDVVULS",
				},
				&txnbuild.SetOptions{
					Signer: &txnbuild.Signer{
						Address: "GD7CGJSJ5OBOU5KOP2UQDH3MPY75UTEY27HVV5XPSL2X6DJ2VGTOSXEU",
						Weight:  20,
					},
				},
				&txnbuild.EndSponsoringFutureReserves{},
			},
			BaseFee:    txnbuild.MinBaseFee,
			Timebounds: txnbuild.NewTimebounds(0, 1),
		},
	)
	require.NoError(t, err)
	txEnc, err := tx.Base64()
	require.NoError(t, err)
	t.Log("Tx:", txEnc)

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Email: "user1@example.com"})
	req := fmt.Sprintf(`{"transaction": "%s"}`, txEnc)
	r := httptest.NewRequest("POST", "/GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4/sign/GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H", strings.NewReader(req))
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}/sign/{signing-address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{"error": "The request was invalid in some way."}`
	assert.JSONEq(t, wantBody, string(body))
}

// Test that when authenticated with a email signing is possible.
func TestAccountSign_signingAddressEmailOtherAuthenticated(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4",
		Identities: []account.Identity{
			{
				Role: "receiver",
				AuthMethods: []account.AuthMethod{
					{Type: account.AuthMethodTypeEmail, Value: "user1@example.com"},
				},
			},
		},
	})
	h := accountSignHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningKeys: []*keypair.Full{
			keypair.MustParseFull("SBIB72S6JMTGJRC6LMKLC5XMHZ2IOHZSZH4SASTN47LECEEJ7QEB6EYK"), // GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H
			keypair.MustParseFull("SBJGZKZ7LU2FQNEFBUOBW4LHCA5BOZCABIJTR7BQIFWQ3P763ZW7MYDD"), // GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS
		},
		NetworkPassphrase: network.TestNetworkPassphrase,
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &txnbuild.SimpleAccount{AccountID: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4"},
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.SetOptions{
					Signer: &txnbuild.Signer{
						Address: "GD7CGJSJ5OBOU5KOP2UQDH3MPY75UTEY27HVV5XPSL2X6DJ2VGTOSXEU",
						Weight:  20,
					},
				},
			},
			BaseFee:    txnbuild.MinBaseFee,
			Timebounds: txnbuild.NewTimebounds(0, 1),
		},
	)
	require.NoError(t, err)
	txEnc, err := tx.Base64()
	require.NoError(t, err)
	t.Log("Tx:", txEnc)

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Email: "user1@example.com"})
	req := `{
	"transaction": "` + txEnc + `"
}`
	r := httptest.NewRequest("POST", "/GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4/sign/GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H", strings.NewReader(req))
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}/sign/{signing-address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
	"signature": "okp0ISR/hjU6ItsfXie6ArlQ3YWkBBqEAM5TJrthALdawV5DzcpuwBKi0QE/iBgoU7eY0hY3RPdxm8mXGNYfCQ==",
	"network_passphrase": "Test SDF Network ; September 2015"
}`
	assert.JSONEq(t, wantBody, string(body))
}

// Test that when the transaction cannot be parsed it errors.
func TestAccountSign_signingAddressCannotParseTransaction(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4",
	})
	h := accountSignHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningKeys: []*keypair.Full{
			keypair.MustParseFull("SBIB72S6JMTGJRC6LMKLC5XMHZ2IOHZSZH4SASTN47LECEEJ7QEB6EYK"), // GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H
			keypair.MustParseFull("SBJGZKZ7LU2FQNEFBUOBW4LHCA5BOZCABIJTR7BQIFWQ3P763ZW7MYDD"), // GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS
		},
		NetworkPassphrase: network.TestNetworkPassphrase,
	}

	txEnc := "AAAAADx2k+7TdIuhwctzxD5y0/w5ASkvTD68az9Nh2fWY+ShAAAAZAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAPHaT7tN0i6HBy3PEPnLT/DkBKS9MPrxrP02HZ9Zj5KEAAAAAAJiWgAAAAAAAAAA"

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4"})
	req := `{
	"transaction": "` + txEnc + `"
}`
	r := httptest.NewRequest("POST", "/GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4/sign/GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H", strings.NewReader(req))
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}/sign/{signing-address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{"error": "The request was invalid in some way."}`
	assert.JSONEq(t, wantBody, string(body))
}

// Test that the sign endpoint responds with an error when given a fee bump transaction.
func TestAccountSign_signingAddressRejectsFeeBumpTx(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4",
	})
	h := accountSignHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningKeys: []*keypair.Full{
			keypair.MustParseFull("SBIB72S6JMTGJRC6LMKLC5XMHZ2IOHZSZH4SASTN47LECEEJ7QEB6EYK"), // GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H
			keypair.MustParseFull("SBJGZKZ7LU2FQNEFBUOBW4LHCA5BOZCABIJTR7BQIFWQ3P763ZW7MYDD"), // GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS
		},
		NetworkPassphrase: network.TestNetworkPassphrase,
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &txnbuild.SimpleAccount{AccountID: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4"},
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.SetOptions{
					Signer: &txnbuild.Signer{
						Address: "GD7CGJSJ5OBOU5KOP2UQDH3MPY75UTEY27HVV5XPSL2X6DJ2VGTOSXEU",
						Weight:  20,
					},
				},
			},
			BaseFee:    txnbuild.MinBaseFee,
			Timebounds: txnbuild.NewTimebounds(0, 1),
		},
	)
	require.NoError(t, err)

	feeBumpTx, err := txnbuild.NewFeeBumpTransaction(
		txnbuild.FeeBumpTransactionParams{
			Inner:      tx,
			FeeAccount: "GD7CGJSJ5OBOU5KOP2UQDH3MPY75UTEY27HVV5XPSL2X6DJ2VGTOSXEU",
			BaseFee:    2 * txnbuild.MinBaseFee,
		},
	)

	feeBumpTxB64, err := feeBumpTx.Base64()
	require.NoError(t, err)
	t.Log("Tx:", feeBumpTxB64)

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4"})
	req := `{
	"transaction": "` + feeBumpTxB64 + `"
}`
	r := httptest.NewRequest("POST", "/GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4/sign/GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H", strings.NewReader(req))
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}/sign/{signing-address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{"error": "The request was invalid in some way."}`
	assert.JSONEq(t, wantBody, string(body))
}

// Test that the request can be made as content-text form instead of JSON.
func TestAccountSign_signingAddressValidContentTypeForm(t *testing.T) {
	s := &account.DBStore{DB: dbtest.Open(t).Open()}
	s.Add(account.Account{
		Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4",
	})
	h := accountSignHandler{
		Logger:       supportlog.DefaultLogger,
		AccountStore: s,
		SigningKeys: []*keypair.Full{
			keypair.MustParseFull("SBIB72S6JMTGJRC6LMKLC5XMHZ2IOHZSZH4SASTN47LECEEJ7QEB6EYK"), // GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H
			keypair.MustParseFull("SBJGZKZ7LU2FQNEFBUOBW4LHCA5BOZCABIJTR7BQIFWQ3P763ZW7MYDD"), // GAPE22DOMALCH42VOR4S3HN6KIZZ643G7D3GNTYF4YOWWXP6UVRAF5JS
		},
		NetworkPassphrase: network.TestNetworkPassphrase,
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &txnbuild.SimpleAccount{AccountID: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4"},
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.SetOptions{
					Signer: &txnbuild.Signer{
						Address: "GD7CGJSJ5OBOU5KOP2UQDH3MPY75UTEY27HVV5XPSL2X6DJ2VGTOSXEU",
						Weight:  20,
					},
				},
			},
			BaseFee:    txnbuild.MinBaseFee,
			Timebounds: txnbuild.NewTimebounds(0, 1),
		},
	)
	require.NoError(t, err)
	txEnc, err := tx.Base64()
	require.NoError(t, err)
	t.Log("Tx:", txEnc)

	ctx := context.Background()
	ctx = auth.NewContext(ctx, auth.Auth{Address: "GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4"})
	reqValues := url.Values{}
	reqValues.Set("transaction", txEnc)
	req := reqValues.Encode()
	t.Log("Request Body:", req)
	r := httptest.NewRequest("POST", "/GA6HNE7O2N2IXIOBZNZ4IPTS2P6DSAJJF5GD5PDLH5GYOZ6WMPSKCXD4/sign/GBOG4KF66M4AFRBUHOTJQJRO7BGGFCSGIICTI5BHXHKXCWV2C67QRN5H", strings.NewReader(req))
	r = r.WithContext(ctx)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	m := chi.NewMux()
	m.Post("/{address}/sign/{signing-address}", h.ServeHTTP)
	m.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
	"signature": "okp0ISR/hjU6ItsfXie6ArlQ3YWkBBqEAM5TJrthALdawV5DzcpuwBKi0QE/iBgoU7eY0hY3RPdxm8mXGNYfCQ==",
	"network_passphrase": "Test SDF Network ; September 2015"
}`
	assert.JSONEq(t, wantBody, string(body))
}
