package cmp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeRoute(t *testing.T) {
	var testCases = []struct {
		Input  string
		Output string
	}{
		{"/accounts/*/payments", "^/accounts/[^?/]+/payments[?/]?[^/]*$"},
		{"/accounts/*", "^/accounts/[^?/]+[?/]?[^/]*$"},
		{"/accounts", "^/accounts[?/]?[^/]*$"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s -> %s", tc.Input, tc.Output), func(t *testing.T) {
			assert.Equal(t, tc.Output, MakeRoute(tc.Input).regexp.String())
		})
	}
}
