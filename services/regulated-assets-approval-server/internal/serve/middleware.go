package serve

import (
	"net/http"

	"github.com/rs/cors"
)

func corsHandler(next http.Handler) http.Handler {
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"GET", "PUT", "POST", "PATCH", "DELETE", "HEAD", "OPTIONS"},
	})
	return cors.Handler(next)
}
