package http

import (
	"fmt"
	stdhttp "net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRun_setupDefault(t *testing.T) {

	// test that using no handler panics
	assert.Panics(t, func() {
		setup(Config{
			Handler: nil,
		})
	})

	// test defaults
	srv := setup(Config{
		Handler: stdhttp.NotFoundHandler(),
	})

	assert.Equal(t, defaultShutdownGracePeriod, srv.Timeout)
	assert.Equal(t, defaultReadTimeout, srv.ReadTimeout)
	assert.Equal(t, time.Duration(0), srv.WriteTimeout)
	assert.Equal(t, time.Duration(0), srv.IdleTimeout)
	assert.Equal(t, defaultListenAddr, srv.Server.Addr)
	assert.Equal(t, time.Duration(0), srv.TCPKeepAlive)
}

func TestRun_setupNonDefault(t *testing.T) {

	testHandler := stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {})

	onStarting := func() {
		fmt.Println("starting server")
	}
	onStopping := func() {
		fmt.Println("stopping server")
	}
	onStopped := func() {
		fmt.Println("stopped server")
	}

	srv := setup(Config{
		Handler:             testHandler,
		ListenAddr:          "1234",
		ShutdownGracePeriod: 25 * time.Second,
		ReadTimeout:         5 * time.Second,
		WriteTimeout:        35 * time.Second,
		IdleTimeout:         2 * time.Minute,
		TCPKeepAlive:        3 * time.Minute,
		OnStarting:          onStarting,
		OnStopping:          onStopping,
		OnStopped:           onStopped,
	})

	assert.Equal(t, "1234", srv.Addr)
	assert.Equal(t, 25*time.Second, srv.Timeout)
	assert.Equal(t, 5*time.Second, srv.ReadTimeout)
	assert.Equal(t, 35*time.Second, srv.WriteTimeout)
	assert.Equal(t, 2*time.Minute, srv.IdleTimeout)
	assert.Equal(t, 3*time.Minute, srv.TCPKeepAlive)
}
