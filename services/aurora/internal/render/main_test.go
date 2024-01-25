package render

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNegotiate(t *testing.T) {
	r, err := http.NewRequest("GET", "/ledgers", nil)
	assert.Nil(t, err)
	r.WithContext(context.Background())

	testCases := []struct {
		Header               string
		ExpectedResponseType string
	}{
		// Obeys the Accept header's prioritization
		{"application/hal+json", MimeHal},
		{"text/event-stream,application/hal+json", MimeEventStream},
		// Defaults to HAL
		{"text/event-stream;q=0.5,application/hal+json", MimeHal},
		{"", MimeHal},
		// Returns empty string for invalid type
		{"text/plain", ""},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			r.Header.Set("Accept", tc.Header)
			assert.Equal(t, tc.ExpectedResponseType, Negotiate(r))
		})
	}

	// Defaults to MimeHal even with no Accept key set
	r.Header.Del("Accept")
	assert.Equal(t, MimeHal, Negotiate(r))
}
