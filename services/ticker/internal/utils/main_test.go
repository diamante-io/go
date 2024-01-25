package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceDiff(t *testing.T) {
	slice1 := []string{"a", "b", "c"}
	slice2 := []string{"a", "b"}

	diff := SliceDiff(slice1, slice2)
	assert.Contains(t, diff, "c")
	assert.NotContains(t, diff, "a")
	assert.NotContains(t, diff, "b")
	assert.Equal(t, 1, len(diff))
}
