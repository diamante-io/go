package auth

import (
	"errors"
	"net/http"
	"time"

	"go/keypair"
	"go/support/http/httpauthz"
	"go/support/log"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

// SEP10Middleware provides middleware for handling an authentication SEP-10 JWT.
func SEP10Middleware(issuer string, ks jose.JSONWebKeySet) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if address, k, ok := sep10ClaimsFromRequest(r, issuer, ks); ok {
				ctx := r.Context()
				auth, _ := FromContext(ctx)
				auth.Address = address

				log.Ctx(ctx).
					WithField("jwkkid", k.KeyID).
					WithField("address", address).
					Info("SEP-10 JWT verified.")

				ctx = NewContext(ctx, auth)
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		})
	}
}

type sep10JWTClaims struct {
	jwt.Claims
}

func (c sep10JWTClaims) Validate(issuer string) error {
	if c.Claims.IssuedAt == nil {
		return errors.New("validation failed, no issued at (iat) in token")
	}
	if c.Claims.Expiry == nil {
		return errors.New("validation failed, no expiry (exp) in token")
	}
	expectedClaims := jwt.Expected{
		Issuer: issuer,
		Time:   time.Now(),
	}
	return c.Claims.Validate(expectedClaims)
}

func sep10ClaimsFromRequest(r *http.Request, issuer string, ks jose.JSONWebKeySet) (address string, k jose.JSONWebKey, ok bool) {
	authHeader := r.Header.Get("Authorization")
	tokenEncoded := httpauthz.ParseBearerToken(authHeader)
	if tokenEncoded == "" {
		return "", jose.JSONWebKey{}, false
	}
	token, err := jwt.ParseSigned(tokenEncoded)
	if err != nil {
		return "", jose.JSONWebKey{}, false
	}
	tokenClaims := sep10JWTClaims{}
	verified := false
	verifiedWithKey := jose.JSONWebKey{}
	for _, k := range ks.Keys {
		err = token.Claims(k, &tokenClaims)
		if err == nil {
			verified = true
			verifiedWithKey = k
			break
		}
	}
	if !verified {
		return "", jose.JSONWebKey{}, false
	}
	err = tokenClaims.Validate(issuer)
	if err != nil {
		return "", jose.JSONWebKey{}, false
	}
	address = tokenClaims.Subject
	_, err = keypair.ParseAddress(address)
	if err != nil {
		return "", jose.JSONWebKey{}, false
	}
	return address, verifiedWithKey, true
}
