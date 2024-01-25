// Package jwtkey provides utility functions for generating, serializing and
// deserializing JWT ECDSA keys.
//
// TODO: Replace EC function usages with PKCS8 functions for supporting ECDSA
// and RSA keys instead of only supporting ECDSA. The fact this package only
// supports ECDSA is unnecessary.
package jwtkey

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

	"github.com/diamcircle/go/support/errors"
)

// GenerateKey is a convenience function for generating an ECDSA key for use as
// a JWT key. It uses the P256 curve. To use other curves use the crypto/ecdsa
// package directly.
func GenerateKey() (*ecdsa.PrivateKey, error) {
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, errors.Wrap(err, "generating ECDSA key")
	}
	return k, nil
}
