package time

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMillisRoundUp(t *testing.T) {
	assert.Equal(t, MillisFromInt64(20), MillisFromInt64(13).RoundUp(10))
	assert.Equal(t, MillisFromInt64(10), MillisFromInt64(10).RoundUp(10))
	assert.Equal(t, MillisFromInt64(15), MillisFromInt64(15).RoundUp(0))
}

func TestMillisRoundDown(t *testing.T) {
	assert.Equal(t, MillisFromInt64(10), MillisFromInt64(13).RoundDown(10))
	assert.Equal(t, MillisFromInt64(10), MillisFromInt64(10).RoundDown(10))
}

func TestMillisParsing(t *testing.T) {
	expected := Now()
	actual, err := MillisFromString(expected.String())
	assert.NoError(t, err)
	assert.Equal(t, actual, expected)
}

func TestMillisToTime(t *testing.T) {
	assert.Equal(t, int64(1510831636149000000), Millis(1510831636149).ToTime().UnixNano())
}

func TestMillisFromSeconds(t *testing.T) {
	assert.Equal(t, MillisFromInt64(10000), MillisFromSeconds(10))
}
