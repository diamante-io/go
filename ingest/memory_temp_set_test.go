package ingest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryTempSet(t *testing.T) {
	s := memoryTempSet{}
	assert.Nil(t, s.m)
	err := s.Open()
	assert.NoError(t, err)
	assert.NotNil(t, s.m)

	err = s.Add("a")
	assert.NoError(t, err)

	err = s.Add("b")
	assert.NoError(t, err)

	v, err := s.Exist("a")
	assert.NoError(t, err)
	assert.True(t, v)

	v, err = s.Exist("b")
	assert.NoError(t, err)
	assert.True(t, v)

	// Get for not-set key should return false
	v, err = s.Exist("c")
	assert.NoError(t, err)
	assert.False(t, v)

	err = s.Close()
	assert.NoError(t, err)
	assert.Nil(t, s.m)
}
