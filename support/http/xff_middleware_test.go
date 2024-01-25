package http

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func TestXFFMiddlewareWrongConfig(t *testing.T) {
	assert.Panics(t, func() {
		XFFMiddleware(XFFMiddlewareConfig{
			BehindCloudflare:      true,
			BehindAWSLoadBalancer: true,
		})
	})
}

func TestXFFMiddlewareCloudFlare(t *testing.T) {
	xff := XFFMiddleware(XFFMiddlewareConfig{
		BehindCloudflare: true,
	})

	mockHandler := &MockHandler{}
	mockHandler.On("ServeHTTP", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			r := args.Get(1).(*http.Request)
			assert.Equal(t, "2.2.2.2", r.RemoteAddr)
		}).Once()
	handler := xff(mockHandler)
	handler.ServeHTTP(nil, &http.Request{
		RemoteAddr: "127.0.0.1",
		Header:     xffHeaders("CF-Connecting-IP", "2.2.2.2"),
	})
}

func TestXFFMiddlewareAWS(t *testing.T) {
	xff := XFFMiddleware(XFFMiddlewareConfig{
		BehindAWSLoadBalancer: true,
	})

	mockHandler := &MockHandler{}
	mockHandler.On("ServeHTTP", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			r := args.Get(1).(*http.Request)
			assert.Equal(t, "2.2.2.2", r.RemoteAddr)
		}).Once()
	handler := xff(mockHandler)
	handler.ServeHTTP(nil, &http.Request{
		RemoteAddr: "127.0.0.1",
		Header:     xffHeaders("X-Forwarded-For", "1.1.1.1,2.2.2.2"),
	})
}

func TestXFFMiddlewareNormalProxy(t *testing.T) {
	xff := XFFMiddleware(XFFMiddlewareConfig{})

	mockHandler := &MockHandler{}
	mockHandler.On("ServeHTTP", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			r := args.Get(1).(*http.Request)
			assert.Equal(t, "1.1.1.1", r.RemoteAddr)
		}).Once()
	handler := xff(mockHandler)
	handler.ServeHTTP(nil, &http.Request{
		RemoteAddr: "127.0.0.1",
		Header:     xffHeaders("X-Forwarded-For", "1.1.1.1,2.2.2.2"),
	})
}

func xffHeaders(name, value string) http.Header {
	headers := http.Header{}
	headers.Add(name, value)
	return headers
}
