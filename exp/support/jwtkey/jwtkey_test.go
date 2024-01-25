package jwtkey

import (
	"crypto/elliptic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	key, err := GenerateKey()
	require.NoError(t, err)
	assert.Equal(t, elliptic.P256(), key.Curve)
}
