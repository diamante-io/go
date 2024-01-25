package httpjson

import (
	"io/ioutil"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	cases := []struct {
		data            interface{}
		contentType     contentType
		wantContentType string
		wantBody        string
	}{
		{
			data:            map[string]interface{}{"key": "value"},
			contentType:     JSON,
			wantContentType: "application/json; charset=utf-8",
			wantBody:        `{"key":"value"}`,
		},
		{
			data:            map[string]interface{}{"key": "value"},
			contentType:     HALJSON,
			wantContentType: "application/hal+json; charset=utf-8",
			wantBody:        `{"key":"value"}`,
		},
		{
			data:            map[string]interface{}{"key": "value"},
			contentType:     HEALTHJSON,
			wantContentType: "application/health+json; charset=utf-8",
			wantBody:        `{"key":"value"}`,
		},
	}

	for i, tc := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			w := httptest.NewRecorder()
			Render(w, tc.data, tc.contentType)
			resp := w.Result()

			assert.Equal(t, 200, resp.StatusCode)
			assert.Equal(t, tc.wantContentType, resp.Header.Get("Content-Type"))

			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.JSONEq(t, tc.wantBody, string(body))
		})
	}
}

func TestRenderStatus(t *testing.T) {
	cases := []struct {
		data            interface{}
		status          int
		contentType     contentType
		wantContentType string
		wantBody        string
	}{
		{
			data:            map[string]interface{}{"key": "value"},
			status:          200,
			contentType:     JSON,
			wantContentType: "application/json; charset=utf-8",
			wantBody:        `{"key":"value"}`,
		},
		{
			data:            map[string]interface{}{"key": "value"},
			status:          400,
			contentType:     HALJSON,
			wantContentType: "application/hal+json; charset=utf-8",
			wantBody:        `{"key":"value"}`,
		},
		{
			data:            map[string]interface{}{"key": "value"},
			status:          400,
			contentType:     HEALTHJSON,
			wantContentType: "application/health+json; charset=utf-8",
			wantBody:        `{"key":"value"}`,
		},
	}

	for i, tc := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			w := httptest.NewRecorder()
			RenderStatus(w, tc.status, tc.data, tc.contentType)
			resp := w.Result()

			assert.Equal(t, tc.status, resp.StatusCode)
			assert.Equal(t, tc.wantContentType, resp.Header.Get("Content-Type"))

			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.JSONEq(t, tc.wantBody, string(body))
		})
	}
}
