package http

import (
	stdhttp "net/http"
	"strings"
)

// XFFMiddlewareConfig provides a configuration for XFFMiddleware.
type XFFMiddlewareConfig struct {
	BehindCloudflare      bool
	BehindAWSLoadBalancer bool
}

// XFFMiddleware is a middleware that replaces http.Request.RemoteAddr with a
// visitor value based on a given config:
//
//   * If BehindCloudflare is true CF-Connecting-IP header is used.
//   * If BehindAWSLoadBalancer is true the last value of X-Forwarded-For header
//     is used.
//   * If none of above is set the first value of X-Forwarded-For header is
//     used. Note: it's easy to spoof the real IP address if the application is
//     not behind a proxy that maintains a X-Forwarded-For header.
//
// Please note that the new RemoteAddr value may not contain the port part!
func XFFMiddleware(config XFFMiddlewareConfig) func(next stdhttp.Handler) stdhttp.Handler {
	if config.BehindCloudflare && config.BehindAWSLoadBalancer {
		panic("Only one of BehindCloudflare and BehindAWSLoadBalancer options can be selected")
	}

	return func(next stdhttp.Handler) stdhttp.Handler {
		fn := func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
			var newRemoteAddr string

			if config.BehindCloudflare {
				newRemoteAddr = r.Header.Get("CF-Connecting-IP")
			} else {
				ips := strings.Split(r.Header.Get("X-Forwarded-For"), ",")
				if len(ips) > 0 {
					if config.BehindAWSLoadBalancer {
						newRemoteAddr = ips[len(ips)-1]
					} else {
						newRemoteAddr = ips[0]
					}
				}
			}

			newRemoteAddr = strings.TrimSpace(newRemoteAddr)
			if newRemoteAddr != "" {
				r.RemoteAddr = newRemoteAddr
			}
			next.ServeHTTP(w, r)
		}
		return stdhttp.HandlerFunc(fn)
	}
}
