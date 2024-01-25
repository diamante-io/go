package db2

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPageQuery(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var p PageQuery
	var err error

	p, err = NewPageQuery("10", true, "desc", 15)
	require.NoError(err)
	assert.Equal("10", p.Cursor)
	assert.Equal("desc", p.Order)
	assert.Equal(uint64(15), p.Limit)

	// Defaults
	p, err = NewPageQuery("", true, "", 1)
	require.NoError(err)
	assert.Equal("asc", p.Order)
	c, err := p.CursorInt64()
	require.NoError(err)
	assert.Equal(int64(0), c)
	assert.Equal(uint64(1), p.Limit)
	p, err = NewPageQuery("", true, "desc", 1)
	require.NoError(err)
	c, err = p.CursorInt64()
	require.NoError(err)
	assert.Equal(int64(9223372036854775807), c)

	// Max
	p, err = NewPageQuery("", true, "", 200)
	require.NoError(err)

	// Error states
	_, err = NewPageQuery("", true, "foo", 1)
	assert.Error(err)
	_, err = NewPageQuery("", true, "", 0)
	assert.Error(err)
	_, err = NewPageQuery("", true, "", 201)
	assert.Error(err)

}

func TestPageQuery_CursorInt64(t *testing.T) {
	assertInstance := assert.New(t)
	requireInstance := require.New(t)

	var p PageQuery
	var err error

	p = MustPageQuery("1231-4456", false, "asc", 1)
	l, r, err := p.CursorInt64Pair("-")
	requireInstance.NoError(err)
	assertInstance.Equal(int64(1231), l)
	assertInstance.Equal(int64(4456), r)

	// Defaults
	p = MustPageQuery("", false, "asc", 1)
	l, r, err = p.CursorInt64Pair("-")
	requireInstance.NoError(err)
	assertInstance.Equal(int64(0), l)
	assertInstance.Equal(int64(0), r)
	p = MustPageQuery("", false, "desc", 1)
	l, r, err = p.CursorInt64Pair("-")
	requireInstance.NoError(err)
	assertInstance.Equal(int64(math.MaxInt64), l)
	assertInstance.Equal(int64(math.MaxInt64), r)
	p = MustPageQuery("0", false, "", 1)
	_, r, err = p.CursorInt64Pair("-")
	requireInstance.NoError(err)
	assertInstance.Equal(int64(math.MaxInt64), r)

	// Errors
	p = MustPageQuery("123-foo", false, "", 1)
	_, _, err = p.CursorInt64Pair("-")
	assertInstance.Error(err)
	p = MustPageQuery("foo-123", false, "", 1)
	_, _, err = p.CursorInt64Pair("-")
	assertInstance.Error(err)
	p = MustPageQuery("-1:123", false, "", 1)
	_, _, err = p.CursorInt64Pair("-")
	assertInstance.Error(err)
	p = MustPageQuery("111:-123", false, "", 1)
	_, _, err = p.CursorInt64Pair("-")
	assertInstance.Error(err)

	// Regression: -23667108046966785
	p = MustPageQuery("-23667108046966785", false, "asc", 100)
	_, err = p.CursorInt64()
	assertInstance.Error(err)
}
