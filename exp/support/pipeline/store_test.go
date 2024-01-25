package pipeline

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	var s Store

	s.Lock()
	s.Put("value", 0)
	s.Unlock()

	s.Lock()
	v := s.Get("value")
	s.Put("value", v.(int)+1)
	s.Unlock()

	assert.Equal(t, 1, s.Get("value"))
}
