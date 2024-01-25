package httpauthz

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBearerToken(t *testing.T) {
	testCases := []struct {
		AuthorizationHeaderValue string
		WantToken                string
	}{
		// An empty header results in no token.
		{"", ""},

		// Other schemes are ignored.
		{"Basic dXNlcm5hbWU6cGFzc3dvcmQ=", ""},

		// Bearer <token> is the classic example.
		{"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjEyMzM3NDk1Mzh9.eyJpYXQiOjEyMzM3NDk1Mzh9", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjEyMzM3NDk1Mzh9.eyJpYXQiOjEyMzM3NDk1Mzh9"},

		// Bearer may be any case.
		{"bEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjEyMzM3NDk1Mzh9.eyJpYXQiOjEyMzM3NDk1Mzh9", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjEyMzM3NDk1Mzh9.eyJpYXQiOjEyMzM3NDk1Mzh9"},

		// Bearer is required.
		{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjEyMzM3NDk1Mzh9.eyJpYXQiOjEyMzM3NDk1Mzh9", ""},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			token := ParseBearerToken(tc.AuthorizationHeaderValue)
			assert.Equal(t, tc.WantToken, token)
		})
	}
}
