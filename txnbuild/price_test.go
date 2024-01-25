package txnbuild

import (
	"math"
	"testing"

	"go/xdr"

	"github.com/stretchr/testify/assert"
)

func TestPriceFromXDR(t *testing.T) {
	for _, testCase := range []struct {
		name     string
		input    xdr.Price
		expected price
	}{
		{
			"1/2",
			xdr.Price{N: 1, D: 2},
			price{n: 1, d: 2, s: "0.5"},
		},
		{
			"1",
			xdr.Price{N: 1, D: 1},
			price{n: 1, d: 1, s: "1"},
		},
		{
			"1 / 1000000000",
			xdr.Price{N: 1, D: 1000000000},
			price{n: 1, d: 1000000000, s: "0.000000001"},
		},
		{
			"max int 32",
			xdr.Price{N: math.MaxInt32, D: 1},
			price{n: math.MaxInt32, d: 1, s: "2147483600"},
		},
		{
			"1/3",
			xdr.Price{N: 1, D: 3},
			price{n: 1, d: 3, s: "0.33333334"},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			var p price
			p.fromXDR(testCase.input)

			assert.Equal(t, testCase.expected, p)
			assert.Equal(t, testCase.input, p.toXDR())

			assert.NoError(t, p.parse(p.string()))
			assert.Equal(t, testCase.expected, p)
		})
	}
}

func TestPriceParse(t *testing.T) {
	var p price
	assert.NoError(t, p.parse("0.5"))
	assert.Equal(t, price{n: 1, d: 2, s: "0.5"}, p)

	assert.NoError(t, p.parse("00.5"))
	assert.Equal(t, price{n: 1, d: 2, s: "0.5"}, p)

	assert.NoError(t, p.parse("0.50"))
	assert.Equal(t, price{n: 1, d: 2, s: "0.5"}, p)

	assert.EqualError(t, p.parse(""), "cannot parse price from empty string")
	assert.Equal(t, price{n: 1, d: 2, s: "0.5"}, p)

	assert.EqualError(t, p.parse("abc"), "failed to parse price from string: invalid price format: abc")
	assert.Equal(t, price{n: 1, d: 2, s: "0.5"}, p)

	assert.NoError(t, p.parse("0.33333334"))
	assert.Equal(t, price{n: 16666667, d: 50000000, s: "0.33333334"}, p)

	p.fromXDR(xdr.Price{N: 1, D: 3})
	assert.Equal(t, price{n: 1, d: 3, s: "0.33333334"}, p)
	assert.NoError(t, p.parse("00.33333334"))
	assert.Equal(t, price{n: 1, d: 3, s: "0.33333334"}, p)
	assert.NoError(t, p.parse("0.333333340"))
	assert.Equal(t, price{n: 1, d: 3, s: "0.33333334"}, p)
}
