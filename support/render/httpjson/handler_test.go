package httpjson

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	cases := []struct {
		input   interface{}
		output  string
		f       interface{}
		cType   contentType
		wantErr bool
	}{
		{`foo`, `"foo"`, func(ctx context.Context, s string) (string, error) { return s, nil }, JSON, false},
		{struct{ Foo int }{1}, `1`, func(ctx context.Context, param struct{ Foo int }) (int, error) { return param.Foo, nil }, HALJSON, false},
		{``, ``, func(ctx context.Context) (int, error) { return 0, errors.New("test") }, JSON, true},
	}

	for _, tc := range cases {
		h, err := Handler(tc.f, tc.input, tc.cType)
		if err != nil {
			t.Errorf("Handler(%v) got err %v", tc.f, err)
			continue
		}

		resp := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		h.ServeHTTP(resp, req.WithContext(context.Background()))
		if tc.wantErr {
			if resp.Code != 500 {
				t.Errorf("%T response code = %d want 200", tc.f, resp.Code)
			}
			continue
		}

		if resp.Code != 200 {
			t.Errorf("%T response code = %d want 200", tc.f, resp.Code)
		}

		if tc.cType == JSON {
			want := "application/json; charset=utf-8"
			if ct := resp.Header().Get("Content-Type"); ct != want {
				t.Errorf(`Content-Type = %s, want %s`, ct, want)
			}
		} else if tc.cType == HALJSON {
			want := "application/hal+json; charset=utf-8"
			if ct := resp.Header().Get("Content-Type"); ct != want {
				t.Errorf(`Content-Type = %s, want %s`, ct, want)
			}
		}

		got := resp.Body.String()
		if got != tc.output {
			t.Errorf("%T response body = %#q want %#q", tc.f, got, tc.output)
		}
	}
}

func TestReqBodyHandler(t *testing.T) {
	cases := []struct {
		input   string
		output  string
		f       interface{}
		cType   contentType
		wantErr bool
	}{
		{`{"Foo":1}`, `1`, func(ctx context.Context, param struct{ Foo int }) (int, error) { return param.Foo, nil }, JSON, false},
		{``, ``, func(ctx context.Context) (int, error) { return 0, errors.New("test") }, JSON, true},
	}

	for _, tc := range cases {
		h, err := ReqBodyHandler(tc.f, tc.cType)
		if err != nil {
			t.Errorf("Handler(%v) got err %v", tc.f, err)
			continue
		}

		resp := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/", strings.NewReader(tc.input))
		h.ServeHTTP(resp, req.WithContext(context.Background()))
		if tc.wantErr {
			if resp.Code != 500 {
				t.Errorf("%T response code = %d want 200", tc.f, resp.Code)
			}
			continue
		}

		if resp.Code != 200 {
			t.Errorf("%T response code = %d want 200", tc.f, resp.Code)
		}

		if tc.cType == JSON {
			want := "application/json; charset=utf-8"
			if ct := resp.Header().Get("Content-Type"); ct != want {
				t.Errorf(`Content-Type = %s, want %s`, ct, want)
			}
		}

		got := resp.Body.String()
		if got != tc.output {
			t.Errorf("%T response body = %#q want %#q", tc.f, got, tc.output)
		}
	}
}

func TestFuncParamTypeError(t *testing.T) {
	cases := []interface{}{
		0,                                        // not a function
		"a string",                               // not a function
		func() (int, error) { return 0, nil },    // no inputs
		func(int) (int, error) { return 0, nil }, // first input is not context
		func(context.Context, int, int) (int, error) { return 0, nil }, // too many inputs
		func(context.Context, int) (int, int) { return 0, 0 },          // second return value is not an error
		func() (int, int, error) { return 0, 0, nil },                  // too many return values
	}

	for _, tc := range cases {
		_, err := funcParamType(reflect.ValueOf(tc))
		if err == nil {
			t.Errorf("funcParamType(%T) wants error", tc)
		}
	}
}

func TestFuncParamTypeNoError(t *testing.T) {
	cases := []interface{}{
		func(context.Context) {},                                  // no return value
		func(context.Context) int { return 0 },                    // one non-error type return value
		func(context.Context) error { return nil },                // one error type return value
		func(context.Context, int) (int, error) { return 0, nil }, // two return values
	}

	for _, tc := range cases {
		_, err := funcParamType(reflect.ValueOf(tc))
		if err != nil {
			t.Errorf("funcParamType(%T) got error %v", tc, err)
		}
	}
}
