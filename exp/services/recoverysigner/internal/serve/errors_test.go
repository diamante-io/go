package serve

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorResponseRender(t *testing.T) {
	w := httptest.NewRecorder()
	serverError.Render(w)
	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.JSONEq(t, `{"error":"An error occurred while processing this request."}`, string(body))
}

func TestErrorHandler(t *testing.T) {
	r := httptest.NewRequest("GET", "/404", nil)
	w := httptest.NewRecorder()
	handler := errorHandler{Error: notFound}
	handler.ServeHTTP(w, r)
	resp := w.Result()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.JSONEq(t, `{"error":"The resource at the url requested was not found."}`, string(body))
}
