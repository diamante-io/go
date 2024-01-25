package problem

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/diamcircle/go/support/render/problem"
	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()
var testRender = func(ctx context.Context, err error) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	problem.Render(ctx, w, err)
	return w
}

func TestCommonProblems(t *testing.T) {
	testCases := []struct {
		testName     string
		p            problem.P
		expectedCode int
	}{
		{"NotFound", problem.NotFound, 404},
		{"RateLimitExceeded", RateLimitExceeded, 429},
		{"ClientDisconneted", ClientDisconnected, 499},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			w := testRender(ctx, tc.p)
			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}

func TestMakeProblemWithInvalidField(t *testing.T) {
	tt := assert.New(t)

	p := problem.NewProblemWithInvalidField(
		problem.NotFound,
		"key",
		errors.New("not found"),
	)

	expectedErr := map[string]interface{}{
		"invalid_field": "key",
		"reason":        "not found",
	}

	tt.Equal(expectedErr, p.Extras)
	tt.Equal(p.Type, "not_found")

	// it doesn't add keys to source problem
	tt.Len(problem.NotFound.Extras, 0)
}

func TestMakeInvalidFieldProblem(t *testing.T) {
	tt := assert.New(t)

	p := problem.MakeInvalidFieldProblem(
		"key",
		errors.New("not found"),
	)

	expectedErr := map[string]interface{}{
		"invalid_field": "key",
		"reason":        "not found",
	}

	tt.Equal(expectedErr, p.Extras)
	tt.Equal(p.Type, "bad_request")

	// it doesn't add keys to source problem
	tt.Len(problem.BadRequest.Extras, 0)
}
