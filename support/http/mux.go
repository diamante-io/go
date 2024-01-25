package http

import (
	"go/support/log"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)

// NewMux returns a new server mux configured with the common defaults used across all
// diamcircle services.
func NewMux(l *log.Entry) *chi.Mux {
	mux := chi.NewMux()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.Recoverer)
	mux.Use(SetLoggerMiddleware(l))
	mux.Use(LoggingMiddleware)

	return mux
}

// NewAPIMux returns a new server mux configured with the common defaults used for a web API in
// diamcircle.
func NewAPIMux(l *log.Entry) *chi.Mux {
	mux := NewMux(l)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"GET", "PUT", "POST", "PATCH", "DELETE", "HEAD", "OPTIONS"},
	})

	mux.Use(c.Handler)
	return mux
}
