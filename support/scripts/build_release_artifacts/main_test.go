package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractFromTag(t *testing.T) {
	cases := []struct {
		Name            string
		Tag             string
		ExpectedBinary  string
		ExpectedVersion string
	}{
		// Successful cases
		{"No hyphen in binary", "federation-v1.0.0", "federation", "v1.0.0"},
		{"hyphen in binary", "diamcircle-sign-v1.0.0", "diamcircle-sign", "v1.0.0"},
		{"non-semver", "federation-master", "federation", "master"},
		// Faileds cases
		{"capitalized", "Federation-v1.0.0", "", ""},
	}

	for _, kase := range cases {
		b, v := extractFromTag(kase.Tag)
		assert.Equal(t, kase.ExpectedBinary, b,
			fmt.Sprintf("Case \"%s\" failed the binary assertion", kase.Name))
		assert.Equal(t, kase.ExpectedVersion, v,
			fmt.Sprintf("Case \"%s\" failed the version assertion", kase.Name))
	}
}
