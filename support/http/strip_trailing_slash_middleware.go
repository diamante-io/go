package http

import (
	stdhttp "net/http"
	"strings"
)

// StripTrailingSlashMiddleware strips trailing slash.
func StripTrailingSlashMiddleware(ignoredPrefixes ...string) func(next stdhttp.Handler) stdhttp.Handler {
	return func(next stdhttp.Handler) stdhttp.Handler {
		fn := func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
			path := r.URL.Path

			// Do not change ignored prefixes
			for _, prefix := range ignoredPrefixes {
				if strings.HasPrefix(path, prefix) {
					next.ServeHTTP(w, r)
					return
				}
			}

			l := len(path)

			// if the path is longer than 1 char (i.e., not '/')
			// and has a trailing slash, remove it.
			if l > 1 && path[l-1] == '/' {
				r.URL.Path = path[0 : l-1]
			}

			next.ServeHTTP(w, r)
		}
		return stdhttp.HandlerFunc(fn)
	}
}
