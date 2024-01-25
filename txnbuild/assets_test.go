package txnbuild

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssetsSorting(t *testing.T) {
	// Native is always first
	a := NativeAsset{}

	// Type is Alphanum4
	b := CreditAsset{Code: "BCDE", Issuer: "GB7TAYRUZGE6TVT7NHP5SMIZRNQA6PLM423EYISAOAP3MKYIQMVYP2JO"}

	// Type is Alphanum12
	c := CreditAsset{Code: "ABCD1", Issuer: "GB7TAYRUZGE6TVT7NHP5SMIZRNQA6PLM423EYISAOAP3MKYIQMVYP2JO"}

	// Code is >
	d := CreditAsset{Code: "ABCD2", Issuer: "GB7TAYRUZGE6TVT7NHP5SMIZRNQA6PLM423EYISAOAP3MKYIQMVYP2JO"}

	// Issuer is >
	e := CreditAsset{Code: "ABCD2", Issuer: "GCEZWKCA5VLDNRLN3RPRJMRZOX3Z6G5CHCGSNFHEYVXM3XOJMDS674JZ"}

	expected := Assets([]Asset{a, b, c, d, e})

	t.Run("basic check it doesn't change stuff", func(t *testing.T) {
		assets := Assets([]Asset{a, b, c, d, e})
		sort.Sort(assets)
		assert.Equal(t, expected, assets)
	})

	// Reverse it and check it still sorts to the same
	t.Run("reverse it and check it sorts the same", func(t *testing.T) {
		assets := Assets([]Asset{e, d, c, b, a})
		sort.Sort(assets)
		assert.Equal(t, expected, assets)
	})
}
