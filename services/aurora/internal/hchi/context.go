// Package hchi provides functions to support embedded and retrieving
// a request id from a go context tree
package hchi

import (
	"context"

	"github.com/go-chi/chi/middleware"
)

// key is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type key int

const (
	reqidKey key = iota
)

// WithRequestID sets the reqid in a new context and returns that context.
func WithRequestID(ctx context.Context, reqid string) context.Context {
	return context.WithValue(ctx, reqidKey, reqid)
}

// WithChiRequestID gets the request id from the chi middleware, sets in a new
// context and returns the context.
func WithChiRequestID(ctx context.Context) context.Context {
	reqid := middleware.GetReqID(ctx)
	return WithRequestID(ctx, reqid)
}

// RequestID returns the request id carried in the context, if any. It returns
// "" if no request id has been set or the context is nil.
func RequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	val, _ := ctx.Value(reqidKey).(string)
	return val
}
