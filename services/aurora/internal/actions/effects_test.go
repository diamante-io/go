package actions

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"go/support/http/httptest"
	"go/support/render/problem"
)

func TestEffectsQuery_BadOperationID(t *testing.T) {
	called := false
	s := httptest.NewServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		qp := EffectsQuery{}
		err := getParams(&qp, r)
		assert.Error(t, err)
		p, ok := err.(*problem.P)
		if assert.True(t, ok) {
			assert.Equal(t, 400, p.Status)
			assert.NotNil(t, p.Extras)
			assert.Equal(t, "op_id", p.Extras["invalid_field"])
			assert.Equal(t, "Operation ID must be an integer higher than 0", p.Extras["reason"])
		}
		called = true
	}))
	defer s.Close()

	_, err := http.Get(s.URL + "/?op_id=-1")
	assert.NoError(t, err)
	assert.True(t, called)

	called = false
	_, err = http.Get(s.URL + "/?op_id=foobar")
	assert.NoError(t, err)
	assert.True(t, called)
}
