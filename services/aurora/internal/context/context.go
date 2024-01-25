package context

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"go/services/aurora/internal/db2/history"
	"go/support/db"
)

type CtxKey string

var RequestContextKey = CtxKey("request")
var SessionContextKey = CtxKey("session")

func RequestFromContext(ctx context.Context) *http.Request {
	found, _ := ctx.Value(&RequestContextKey).(*http.Request)
	return found
}

// requestContext returns a context representing the provided http action.
func RequestContext(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	if r == nil {
		panic("Cannot bind nil *http.Request to context tree")
	}

	return context.WithValue(ctx, &RequestContextKey, r)
}

// BaseURL returns the "base" url for this request, defined as a url containing
// the Host and Scheme portions of the request uri.
func BaseURL(ctx context.Context) *url.URL {
	r := RequestFromContext(ctx)
	if r == nil {
		return nil
	}

	var scheme string
	switch {
	case r.Header.Get("X-Forwarded-Proto") != "":
		scheme = r.Header.Get("X-Forwarded-Proto")
	case r.TLS != nil:
		scheme = "https"
	default:
		scheme = "http"
	}

	return &url.URL{
		Scheme: scheme,
		Host:   r.Host,
	}
}

func HistoryQFromRequest(request *http.Request) (*history.Q, error) {
	ctx := request.Context()
	session, ok := ctx.Value(&SessionContextKey).(db.SessionInterface)
	if !ok {
		return nil, errors.New("missing session in request context")
	}
	return &history.Q{session}, nil
}
