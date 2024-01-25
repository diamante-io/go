package auth

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/square/go-jose.v2"
)

func TestSEP10_addsAddressToClaimIfJWTValid(t *testing.T) {
	issuer := "https://webauth.example.com"

	k1, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	jwks := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{Key: &k1.PublicKey},
		},
	}

	ctx := context.Context(nil)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
	})
	middleware := SEP10Middleware(issuer, jwks)
	handler := middleware(next)

	r := httptest.NewRequest("GET", "/", nil)
	jwtClaims := jwt.MapClaims{
		"iss": "https://webauth.example.com",
		"sub": "GDKABHI4LTLG7UCE6O7Y4D6REHJVS4DLXTVVXTE3BPRRLXPASHSOKG2D",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodES256, jwtClaims).SignedString(k1)
	require.NoError(t, err)
	r.Header.Set("Authorization", "Bearer "+jwtToken)
	handler.ServeHTTP(nil, r)

	assert.NotNil(t, ctx)
	claims, ok := FromContext(ctx)
	assert.Equal(t, true, ok)

	wantClaims := Auth{
		Address: "GDKABHI4LTLG7UCE6O7Y4D6REHJVS4DLXTVVXTE3BPRRLXPASHSOKG2D",
	}
	assert.Equal(t, wantClaims, claims)
}

func TestSEP10_addsAddressToClaimIfJWTValidMultipleJWKS(t *testing.T) {
	issuer := "https://webauth.example.com"

	k1, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	k2, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	keys := []*ecdsa.PrivateKey{k1, k2}
	jwks := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{Key: &k1.PublicKey},
			{Key: &k2.PublicKey},
		},
	}

	ctx := context.Context(nil)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
	})
	middleware := SEP10Middleware(issuer, jwks)
	handler := middleware(next)

	for i, k := range keys {
		t.Run(fmt.Sprintf("known key %d", i), func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", nil)
			jwtClaims := jwt.MapClaims{
				"iss": "https://webauth.example.com",
				"sub": "GDKABHI4LTLG7UCE6O7Y4D6REHJVS4DLXTVVXTE3BPRRLXPASHSOKG2D",
				"iat": time.Now().Unix(),
				"exp": time.Now().Add(time.Hour).Unix(),
			}
			jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodES256, jwtClaims).SignedString(k)
			require.NoError(t, err)
			r.Header.Set("Authorization", "Bearer "+jwtToken)
			handler.ServeHTTP(nil, r)

			assert.NotNil(t, ctx)
			claims, ok := FromContext(ctx)
			assert.Equal(t, true, ok)

			wantClaims := Auth{
				Address: "GDKABHI4LTLG7UCE6O7Y4D6REHJVS4DLXTVVXTE3BPRRLXPASHSOKG2D",
			}
			assert.Equal(t, wantClaims, claims)
		})
	}
	t.Run("unknown key", func(t *testing.T) {
		k3, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		require.NoError(t, err)

		r := httptest.NewRequest("GET", "/", nil)
		jwtClaims := jwt.MapClaims{
			"iss": "https://webauth.example.com",
			"sub": "GDKABHI4LTLG7UCE6O7Y4D6REHJVS4DLXTVVXTE3BPRRLXPASHSOKG2D",
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(time.Hour).Unix(),
		}
		jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodES256, jwtClaims).SignedString(k3)
		require.NoError(t, err)
		r.Header.Set("Authorization", "Bearer "+jwtToken)
		handler.ServeHTTP(nil, r)

		assert.NotNil(t, ctx)
		claims, ok := FromContext(ctx)
		assert.Equal(t, false, ok)

		wantClaims := Auth{}
		assert.Equal(t, wantClaims, claims)
	})
}

func TestSEP10_doesNotAddAddressToClaimIfJWTNotPresent(t *testing.T) {
	issuer := "https://webauth.example.com"

	k1, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	jwks := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{Key: &k1.PublicKey},
		},
	}

	ctx := context.Context(nil)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
	})
	middleware := SEP10Middleware(issuer, jwks)
	handler := middleware(next)

	r := httptest.NewRequest("GET", "/", nil)
	handler.ServeHTTP(nil, r)

	assert.NotNil(t, ctx)
	claims, ok := FromContext(ctx)
	assert.Equal(t, false, ok)

	wantClaims := Auth{}
	assert.Equal(t, wantClaims, claims)
}

func TestSEP10_doesNotAddAddressToClaimIfJWTNoSignature(t *testing.T) {
	issuer := "https://webauth.example.com"

	k1, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	jwks := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{Key: &k1.PublicKey},
		},
	}

	ctx := context.Context(nil)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
	})
	middleware := SEP10Middleware(issuer, jwks)
	handler := middleware(next)

	r := httptest.NewRequest("GET", "/", nil)
	jwtClaims := jwt.MapClaims{
		"iss": "https://webauth.example.com",
		"sub": "GDKABHI4LTLG7UCE6O7Y4D6REHJVS4DLXTVVXTE3BPRRLXPASHSOKG2D",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodES256, jwtClaims).SigningString()
	require.NoError(t, err)
	r.Header.Set("Authorization", "Bearer "+jwtToken)
	handler.ServeHTTP(nil, r)

	assert.NotNil(t, ctx)
	claims, ok := FromContext(ctx)
	assert.Equal(t, false, ok)

	wantClaims := Auth{}
	assert.Equal(t, wantClaims, claims)
}

func TestSEP10_doesNotAddAddressToClaimIfJWTWrongAlg(t *testing.T) {
	issuer := "https://webauth.example.com"

	k1, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	jwks := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{Key: &k1.PublicKey},
		},
	}

	ctx := context.Context(nil)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
	})
	middleware := SEP10Middleware(issuer, jwks)
	handler := middleware(next)

	r := httptest.NewRequest("GET", "/", nil)
	jwtClaims := jwt.MapClaims{
		"iss": "https://webauth.example.com",
		"sub": "GDKABHI4LTLG7UCE6O7Y4D6REHJVS4DLXTVVXTE3BPRRLXPASHSOKG2D",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodNone, jwtClaims).SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)
	r.Header.Set("Authorization", "Bearer "+jwtToken)
	handler.ServeHTTP(nil, r)

	assert.NotNil(t, ctx)
	claims, ok := FromContext(ctx)
	assert.Equal(t, false, ok)

	wantClaims := Auth{}
	assert.Equal(t, wantClaims, claims)
}

func TestSEP10_doesNotAddAddressToClaimIfJWTInvalidSignature(t *testing.T) {
	issuer := "https://webauth.example.com"

	k1, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	jwks := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{Key: &k1.PublicKey},
		},
	}

	k2, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	ctx := context.Context(nil)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
	})
	middleware := SEP10Middleware(issuer, jwks)
	handler := middleware(next)

	r := httptest.NewRequest("GET", "/", nil)
	jwtClaims := jwt.MapClaims{
		"iss": "https://webauth.example.com",
		"sub": "GDKABHI4LTLG7UCE6O7Y4D6REHJVS4DLXTVVXTE3BPRRLXPASHSOKG2D",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodES256, jwtClaims).SignedString(k2)
	require.NoError(t, err)
	r.Header.Set("Authorization", "Bearer "+jwtToken)
	handler.ServeHTTP(nil, r)

	assert.NotNil(t, ctx)
	claims, ok := FromContext(ctx)
	assert.Equal(t, false, ok)

	wantClaims := Auth{}
	assert.Equal(t, wantClaims, claims)
}

func TestSEP10_doesNotAddAddressToClaimIfJWTExpired(t *testing.T) {
	issuer := "https://webauth.example.com"

	k1, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	jwks := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{Key: &k1.PublicKey},
		},
	}

	ctx := context.Context(nil)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
	})
	middleware := SEP10Middleware(issuer, jwks)
	handler := middleware(next)

	r := httptest.NewRequest("GET", "/", nil)
	jwtClaims := jwt.MapClaims{
		"iss": "https://webauth.example.com",
		"sub": "GDKABHI4LTLG7UCE6O7Y4D6REHJVS4DLXTVVXTE3BPRRLXPASHSOKG2D",
		"iat": 1,
		"exp": 1,
	}
	jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodES256, jwtClaims).SignedString(k1)
	require.NoError(t, err)
	r.Header.Set("Authorization", "Bearer "+jwtToken)
	handler.ServeHTTP(nil, r)

	assert.NotNil(t, ctx)
	claims, ok := FromContext(ctx)
	assert.Equal(t, false, ok)

	wantClaims := Auth{}
	assert.Equal(t, wantClaims, claims)
}

func TestSEP10_doesNotAddAddressToClaimIfJWTMissingIAT(t *testing.T) {
	issuer := "https://webauth.example.com"

	k1, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	jwks := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{Key: &k1.PublicKey},
		},
	}

	ctx := context.Context(nil)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
	})
	middleware := SEP10Middleware(issuer, jwks)
	handler := middleware(next)

	r := httptest.NewRequest("GET", "/", nil)
	jwtClaims := jwt.MapClaims{
		"iss": "https://webauth.example.com",
		"sub": "GDKABHI4LTLG7UCE6O7Y4D6REHJVS4DLXTVVXTE3BPRRLXPASHSOKG2D",
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodES256, jwtClaims).SignedString(k1)
	require.NoError(t, err)
	r.Header.Set("Authorization", "Bearer "+jwtToken)
	handler.ServeHTTP(nil, r)

	assert.NotNil(t, ctx)
	claims, ok := FromContext(ctx)
	assert.Equal(t, false, ok)

	wantClaims := Auth{}
	assert.Equal(t, wantClaims, claims)
}

func TestSEP10_doesNotAddAddressToClaimIfJWTMissingEXP(t *testing.T) {
	issuer := "https://webauth.example.com"

	k1, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	jwks := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{Key: &k1.PublicKey},
		},
	}

	ctx := context.Context(nil)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
	})
	middleware := SEP10Middleware(issuer, jwks)
	handler := middleware(next)

	r := httptest.NewRequest("GET", "/", nil)
	jwtClaims := jwt.MapClaims{
		"iss": "https://webauth.example.com",
		"sub": "GDKABHI4LTLG7UCE6O7Y4D6REHJVS4DLXTVVXTE3BPRRLXPASHSOKG2D",
		"iat": time.Now().Unix(),
	}
	jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodES256, jwtClaims).SignedString(k1)
	require.NoError(t, err)
	r.Header.Set("Authorization", "Bearer "+jwtToken)
	handler.ServeHTTP(nil, r)

	assert.NotNil(t, ctx)
	claims, ok := FromContext(ctx)
	assert.Equal(t, false, ok)

	wantClaims := Auth{}
	assert.Equal(t, wantClaims, claims)
}

func TestSEP10_doesNotAddAddressToClaimIfJWTMissingSUB(t *testing.T) {
	issuer := "https://webauth.example.com"

	k1, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	jwks := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{Key: &k1.PublicKey},
		},
	}

	ctx := context.Context(nil)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
	})
	middleware := SEP10Middleware(issuer, jwks)
	handler := middleware(next)

	r := httptest.NewRequest("GET", "/", nil)
	jwtClaims := jwt.MapClaims{
		"iss": "https://webauth.example.com",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodES256, jwtClaims).SignedString(k1)
	require.NoError(t, err)
	r.Header.Set("Authorization", "Bearer "+jwtToken)
	handler.ServeHTTP(nil, r)

	assert.NotNil(t, ctx)
	claims, ok := FromContext(ctx)
	assert.Equal(t, false, ok)

	wantClaims := Auth{}
	assert.Equal(t, wantClaims, claims)
}

func TestSEP10_doesNotAddAddressToClaimIfJWTHasSUBNotContainingGStrkey(t *testing.T) {
	issuer := "https://webauth.example.com"

	k1, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	jwks := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{Key: &k1.PublicKey},
		},
	}

	ctx := context.Context(nil)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
	})
	middleware := SEP10Middleware(issuer, jwks)
	handler := middleware(next)

	r := httptest.NewRequest("GET", "/", nil)
	jwtClaims := jwt.MapClaims{
		"iss": "https://webauth.example.com",
		"sub": "SBAZWVXOQ5LWT5PJSVOA62PVIYZIV3T3HQ3GFC2RUZ6K43QFNF5BLLDE",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodES256, jwtClaims).SignedString(k1)
	require.NoError(t, err)
	r.Header.Set("Authorization", "Bearer "+jwtToken)
	handler.ServeHTTP(nil, r)

	assert.NotNil(t, ctx)
	claims, ok := FromContext(ctx)
	assert.Equal(t, false, ok)

	wantClaims := Auth{}
	assert.Equal(t, wantClaims, claims)
}

func TestSEP10_doesNotAddAddressToClaimIfJWTMissingISSButRequired(t *testing.T) {
	issuer := "https://webauth.example.com"

	k1, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	jwks := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{Key: &k1.PublicKey},
		},
	}

	ctx := context.Context(nil)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
	})
	middleware := SEP10Middleware(issuer, jwks)
	handler := middleware(next)

	r := httptest.NewRequest("GET", "/", nil)
	jwtClaims := jwt.MapClaims{
		"sub": "GDKABHI4LTLG7UCE6O7Y4D6REHJVS4DLXTVVXTE3BPRRLXPASHSOKG2D",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodES256, jwtClaims).SignedString(k1)
	require.NoError(t, err)
	r.Header.Set("Authorization", "Bearer "+jwtToken)
	handler.ServeHTTP(nil, r)

	assert.NotNil(t, ctx)
	claims, ok := FromContext(ctx)
	assert.Equal(t, false, ok)

	wantClaims := Auth{}
	assert.Equal(t, wantClaims, claims)
}

func TestSEP10_doesNotAddAddressToClaimIfJWTHasISSButUnexpectedValue(t *testing.T) {
	issuer := "https://webauth.example.com"

	k1, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	jwks := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{Key: &k1.PublicKey},
		},
	}

	ctx := context.Context(nil)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
	})
	middleware := SEP10Middleware(issuer, jwks)
	handler := middleware(next)

	r := httptest.NewRequest("GET", "/", nil)
	jwtClaims := jwt.MapClaims{
		"iss": "otherissuer",
		"sub": "GDKABHI4LTLG7UCE6O7Y4D6REHJVS4DLXTVVXTE3BPRRLXPASHSOKG2D",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodES256, jwtClaims).SignedString(k1)
	require.NoError(t, err)
	r.Header.Set("Authorization", "Bearer "+jwtToken)
	handler.ServeHTTP(nil, r)

	assert.NotNil(t, ctx)
	claims, ok := FromContext(ctx)
	assert.Equal(t, false, ok)

	wantClaims := Auth{}
	assert.Equal(t, wantClaims, claims)
}

func TestSEP10_addAddressToClaimIfJWTMissingISSButNotRequired(t *testing.T) {
	issuer := ""

	k1, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	jwks := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{Key: &k1.PublicKey},
		},
	}

	ctx := context.Context(nil)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
	})
	middleware := SEP10Middleware(issuer, jwks)
	handler := middleware(next)

	r := httptest.NewRequest("GET", "/", nil)
	jwtClaims := jwt.MapClaims{
		"sub": "GDKABHI4LTLG7UCE6O7Y4D6REHJVS4DLXTVVXTE3BPRRLXPASHSOKG2D",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodES256, jwtClaims).SignedString(k1)
	require.NoError(t, err)
	r.Header.Set("Authorization", "Bearer "+jwtToken)
	handler.ServeHTTP(nil, r)

	assert.NotNil(t, ctx)
	claims, ok := FromContext(ctx)
	assert.Equal(t, true, ok)

	wantClaims := Auth{
		Address: "GDKABHI4LTLG7UCE6O7Y4D6REHJVS4DLXTVVXTE3BPRRLXPASHSOKG2D",
	}
	assert.Equal(t, wantClaims, claims)
}

func TestSEP10_addAddressToClaimIfJWTHasISSButNotRequired(t *testing.T) {
	issuer := ""

	k1, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	jwks := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{Key: &k1.PublicKey},
		},
	}

	ctx := context.Context(nil)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
	})
	middleware := SEP10Middleware(issuer, jwks)
	handler := middleware(next)

	r := httptest.NewRequest("GET", "/", nil)
	jwtClaims := jwt.MapClaims{
		"iss": "otherservice",
		"sub": "GDKABHI4LTLG7UCE6O7Y4D6REHJVS4DLXTVVXTE3BPRRLXPASHSOKG2D",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodES256, jwtClaims).SignedString(k1)
	require.NoError(t, err)
	r.Header.Set("Authorization", "Bearer "+jwtToken)
	handler.ServeHTTP(nil, r)

	assert.NotNil(t, ctx)
	claims, ok := FromContext(ctx)
	assert.Equal(t, true, ok)

	wantClaims := Auth{
		Address: "GDKABHI4LTLG7UCE6O7Y4D6REHJVS4DLXTVVXTE3BPRRLXPASHSOKG2D",
	}
	assert.Equal(t, wantClaims, claims)
}
