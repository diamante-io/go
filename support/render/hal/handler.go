package hal

import (
	"context"
	"net/http"

	"github.com/diamcircle/go/support/render/httpjson"
)

func Handler(fn, param interface{}) (http.Handler, error) {
	return httpjson.Handler(fn, param, httpjson.HALJSON)
}

func ExecuteFunc(ctx context.Context, fn, param interface{}) (interface{}, bool, error) {
	return httpjson.ExecuteFunc(ctx, fn, param, httpjson.HALJSON)
}
