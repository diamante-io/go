package hchi

import (
	"context"
	"testing"

	"github.com/go-chi/chi/middleware"
	"github.com/stretchr/testify/assert"
)

func TestRequestContext(t *testing.T) {
	ctx := WithRequestID(context.Background(), "2")
	assert.Equal(t, "2", ctx.Value(reqidKey))

	ctx2 := WithRequestID(ctx, "3")
	assert.Equal(t, "3", ctx2.Value(reqidKey))
	assert.Equal(t, "2", ctx.Value(reqidKey))
}

func TestRequestContextFromCHI(t *testing.T) {
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "foobar")
	ctx2 := WithChiRequestID(ctx)
	assert.Equal(t, "foobar", RequestID(ctx2))
}

func TestRequestFromContext(t *testing.T) {
	ctx := WithRequestID(context.Background(), "2")
	ctx2 := WithRequestID(ctx, "3")
	assert.Equal(t, "", RequestID(context.Background()))
	assert.Equal(t, "2", RequestID(ctx))
	assert.Equal(t, "3", RequestID(ctx2))
}
