package diamcircletoml

import (
	"strings"
	"testing"

	"net/http"

	"github.com/diamcircle/go/support/http/httptest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientURL(t *testing.T) {
	//HACK:  we're testing an internal method rather than setting up a http client
	//mock.

	c := &Client{UseHTTP: false}
	assert.Equal(t, "https://diamcircle.org/.well-known/diamcircle.toml", c.url("diamcircle.org"))

	c = &Client{UseHTTP: true}
	assert.Equal(t, "http://diamcircle.org/.well-known/diamcircle.toml", c.url("diamcircle.org"))
}

func TestClient(t *testing.T) {
	h := httptest.NewClient()
	c := &Client{HTTP: h}

	// happy path
	h.
		On("GET", "https://diamcircle.org/.well-known/diamcircle.toml").
		ReturnString(http.StatusOK,
			`FEDERATION_SERVER="https://localhost/federation"`,
		)
	stoml, err := c.GetDiamcircleToml("diamcircle.org")
	require.NoError(t, err)
	assert.Equal(t, "https://localhost/federation", stoml.FederationServer)

	// diamcircle.toml exceeds limit
	h.
		On("GET", "https://toobig.org/.well-known/diamcircle.toml").
		ReturnString(http.StatusOK,
			`FEDERATION_SERVER="https://localhost/federation`+strings.Repeat("0", DiamcircleTomlMaxSize)+`"`,
		)
	stoml, err = c.GetDiamcircleToml("toobig.org")
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "diamcircle.toml response exceeds")
	}

	// not found
	h.
		On("GET", "https://missing.org/.well-known/diamcircle.toml").
		ReturnNotFound()
	stoml, err = c.GetDiamcircleToml("missing.org")
	assert.EqualError(t, err, "http request failed with non-200 status code")

	// invalid toml
	h.
		On("GET", "https://json.org/.well-known/diamcircle.toml").
		ReturnJSON(http.StatusOK, map[string]string{"hello": "world"})
	stoml, err = c.GetDiamcircleToml("json.org")

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "toml decode failed")
	}
}
