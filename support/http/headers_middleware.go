package http

import (
	stdhttp "net/http"
	"strings"
)

// HeadersMiddleware sends headers
func HeadersMiddleware(headers stdhttp.Header, ignoredPrefixes ...string) func(next stdhttp.Handler) stdhttp.Handler {
	return func(next stdhttp.Handler) stdhttp.Handler {
		fn := func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
			// Do not change ignored prefixes
			for _, prefix := range ignoredPrefixes {
				if strings.HasPrefix(r.URL.Path, prefix) {
					next.ServeHTTP(w, r)
					return
				}
			}

			// headers.Write(w)
			for key := range headers {
				w.Header().Set(key, headers.Get(key))
			}
			next.ServeHTTP(w, r)
		}
		return stdhttp.HandlerFunc(fn)
	}
}
