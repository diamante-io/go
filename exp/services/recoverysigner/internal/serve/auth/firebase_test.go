package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	firebaseauth "firebase.google.com/go/auth"
	"github.com/stretchr/testify/assert"
)

// Test that if the token verifier says there is a Firebase token that contains
// a phone number claim, the claims stored in the context should contain it.
func TestFirebase_tokenWithPhoneNumber(t *testing.T) {
	tokenVerifier := FirebaseTokenVerifierFunc(func(_ *http.Request) (*firebaseauth.Token, bool) {
		token := &firebaseauth.Token{
			Claims: map[string]interface{}{
				"phone_number": "+10000000000",
			},
		}
		return token, true
	})

	claims := Auth{}
	claimsOK := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, claimsOK = FromContext(r.Context())
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	FirebaseMiddleware(tokenVerifier)(next).ServeHTTP(w, r)

	wantClaims := Auth{
		PhoneNumber: "+10000000000",
	}
	assert.Equal(t, wantClaims, claims)
	assert.Equal(t, true, claimsOK)
}

// Test that if the token verifier says there is a Firebase token that contains
// an email and email_verified=true claim, the claims stored in the context
// should contain it.
func TestFirebase_tokenWithEmailVerifiedTrue(t *testing.T) {
	tokenVerifier := FirebaseTokenVerifierFunc(func(_ *http.Request) (*firebaseauth.Token, bool) {
		token := &firebaseauth.Token{
			Claims: map[string]interface{}{
				"email":          "user@example.com",
				"email_verified": true,
			},
		}
		return token, true
	})

	claims := Auth{}
	claimsOK := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, claimsOK = FromContext(r.Context())
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	FirebaseMiddleware(tokenVerifier)(next).ServeHTTP(w, r)

	wantClaims := Auth{
		Email: "user@example.com",
	}
	assert.Equal(t, wantClaims, claims)
	assert.Equal(t, true, claimsOK)
}

// Test that if the token verifier says there is a Firebase token that contains
// an email claim and email_verified=false, the claims stored in the context
// should not contain the email.
func TestFirebase_tokenWithEmailVerifiedFalse(t *testing.T) {
	tokenVerifier := FirebaseTokenVerifierFunc(func(_ *http.Request) (*firebaseauth.Token, bool) {
		token := &firebaseauth.Token{
			Claims: map[string]interface{}{
				"email":          "user@example.com",
				"email_verified": false,
			},
		}
		return token, true
	})

	claims := Auth{}
	claimsOK := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, claimsOK = FromContext(r.Context())
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	FirebaseMiddleware(tokenVerifier)(next).ServeHTTP(w, r)

	wantClaims := Auth{}
	assert.Equal(t, wantClaims, claims)
	assert.Equal(t, true, claimsOK)
}

// Test that if the token verifier says there is a Firebase token that contains
// an email claim without email_verified, the claims stored in the context
// should not contain the email.
func TestFirebase_tokenWithEmail(t *testing.T) {
	tokenVerifier := FirebaseTokenVerifierFunc(func(_ *http.Request) (*firebaseauth.Token, bool) {
		token := &firebaseauth.Token{
			Claims: map[string]interface{}{
				"email": "user@example.com",
			},
		}
		return token, true
	})

	claims := Auth{}
	claimsOK := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, claimsOK = FromContext(r.Context())
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	FirebaseMiddleware(tokenVerifier)(next).ServeHTTP(w, r)

	wantClaims := Auth{}
	assert.Equal(t, wantClaims, claims)
	assert.Equal(t, true, claimsOK)
}

// Test that if the token verifier says there is a Firebase token that
// contains a phone number and an email claim, the claims stored in the
// context should contain both.
func TestFirebase_tokenWithPhoneNumberAndEmail(t *testing.T) {
	tokenVerifier := FirebaseTokenVerifierFunc(func(_ *http.Request) (*firebaseauth.Token, bool) {
		token := &firebaseauth.Token{
			Claims: map[string]interface{}{
				"phone_number":   "+10000000000",
				"email":          "user@example.com",
				"email_verified": true,
			},
		}
		return token, true
	})

	claims := Auth{}
	claimsOK := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, claimsOK = FromContext(r.Context())
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	FirebaseMiddleware(tokenVerifier)(next).ServeHTTP(w, r)

	wantClaims := Auth{
		PhoneNumber: "+10000000000",
		Email:       "user@example.com",
	}
	assert.Equal(t, wantClaims, claims)
	assert.Equal(t, true, claimsOK)
}

// Test that if the token verifier says there is a Firebase token that contains
// a phone number or an email claim, and there are other claims fields filled
// in, the claims stored in the context should contain the merging of both.
func TestFirebase_tokenWithPhoneNumberAndEmailAppendsToOtherClaims(t *testing.T) {
	tokenVerifier := FirebaseTokenVerifierFunc(func(_ *http.Request) (*firebaseauth.Token, bool) {
		token := &firebaseauth.Token{
			Claims: map[string]interface{}{
				"phone_number":   "+10000000000",
				"email":          "user@example.com",
				"email_verified": true,
			},
		}
		return token, true
	})

	claims := Auth{}
	claimsOK := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, claimsOK = FromContext(r.Context())
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	initialClaims := Auth{
		Address: "GCJYKECSRMIQX3KK62VPJ64NFNWV3EKJPUAXSXKKA7XSHN43VDHNKANO",
	}
	r = r.WithContext(NewContext(r.Context(), initialClaims))
	FirebaseMiddleware(tokenVerifier)(next).ServeHTTP(w, r)

	wantClaims := Auth{
		Address:     "GCJYKECSRMIQX3KK62VPJ64NFNWV3EKJPUAXSXKKA7XSHN43VDHNKANO",
		PhoneNumber: "+10000000000",
		Email:       "user@example.com",
	}
	assert.Equal(t, wantClaims, claims)
	assert.Equal(t, true, claimsOK)
}

// Test that if the token verifier says there is an empty Firebase token that
// does not have a phone number or email claim, the claims stored in the
// context should be empty.
func TestFirebase_tokenWithNone(t *testing.T) {
	tokenVerifier := FirebaseTokenVerifierFunc(func(_ *http.Request) (*firebaseauth.Token, bool) {
		return &firebaseauth.Token{}, true
	})

	claims := Auth{}
	claimsOK := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, claimsOK = FromContext(r.Context())
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	FirebaseMiddleware(tokenVerifier)(next).ServeHTTP(w, r)

	wantClaims := Auth{}
	assert.Equal(t, wantClaims, claims)
	assert.Equal(t, true, claimsOK)
}

// Test that if the token verifier says there is no Firebase token, the claims
// stored in the context should be empty.
func TestFirebase_noToken(t *testing.T) {
	tokenVerifier := FirebaseTokenVerifierFunc(func(_ *http.Request) (*firebaseauth.Token, bool) {
		return nil, false
	})

	claims := Auth{}
	claimsOK := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, claimsOK = FromContext(r.Context())
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	FirebaseMiddleware(tokenVerifier)(next).ServeHTTP(w, r)

	wantClaims := Auth{}
	assert.Equal(t, wantClaims, claims)
	assert.Equal(t, false, claimsOK)
}
