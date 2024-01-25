package strutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKebabToConstantCase(t *testing.T) {
	assert.Equal(t, "ENABLE_ASSET_STATS", KebabToConstantCase("enable-asset-stats"), "ordinary use")
	assert.Equal(t, "ABC_DEF", KebabToConstantCase("ABC_DEF"), "ignores uppercase and underscores")
}
