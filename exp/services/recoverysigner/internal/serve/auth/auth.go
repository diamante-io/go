package auth

import (
	"context"
)

type contextKey int

const (
	authContextKey contextKey = iota
)

// Auth holds a set of details that have been authenticated about a client.
type Auth struct {
	Address     string
	PhoneNumber string
	Email       string
}

// FromContext returns auth details that are stored in the context.
func FromContext(ctx context.Context) (Auth, bool) {
	if a, ok := ctx.Value(authContextKey).(Auth); ok {
		return a, true
	}
	return Auth{}, false
}

// NewContext returns a new context that is a copy of the given context with
// the auth details set within. An Auth can be retrieved from the context using
// FromContext.
func NewContext(ctx context.Context, a Auth) context.Context {
	return context.WithValue(ctx, authContextKey, a)
}
