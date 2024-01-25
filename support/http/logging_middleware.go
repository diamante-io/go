package http

import (
	stdhttp "net/http"
	"strings"
	"time"

	"go/support/http/mutil"
	"go/support/log"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Options allow the middleware logger to accept additional information.
type Options struct {
	ExtraHeaders []string
}

// SetLogger is a middleware that sets a logger on the context.
func SetLoggerMiddleware(l *log.Entry) func(stdhttp.Handler) stdhttp.Handler {
	return func(next stdhttp.Handler) stdhttp.Handler {
		return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
			ctx := r.Context()
			ctx = log.Set(ctx, l)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// LoggingMiddleware is a middleware that logs requests to the logger.
func LoggingMiddleware(next stdhttp.Handler) stdhttp.Handler {
	return LoggingMiddlewareWithOptions(Options{})(next)
}

// LoggingMiddlewareWithOptions is a middleware that logs requests to the logger.
// Requires an Options struct to accept additional information.
func LoggingMiddlewareWithOptions(options Options) func(stdhttp.Handler) stdhttp.Handler {
	return func(next stdhttp.Handler) stdhttp.Handler {
		return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
			mw := mutil.WrapWriter(w)
			ctx := log.PushContext(r.Context(), func(l *log.Entry) *log.Entry {
				return l.WithFields(log.F{
					"req": middleware.GetReqID(r.Context()),
				})
			})
			r = r.WithContext(ctx)

			logStartOfRequest(r, options.ExtraHeaders)

			then := time.Now()
			next.ServeHTTP(mw, r)
			duration := time.Since(then)

			logEndOfRequest(r, duration, mw)
		})
	}
}

// logStartOfRequest emits the logline that reports that an http request is
// beginning processing.
func logStartOfRequest(
	r *stdhttp.Request,
	extraHeaders []string,
) {
	fields := log.F{}
	for _, header := range extraHeaders {
		// Strips "-" characters and lowercases new logrus.Fields keys to be uniform with the other keys in the logger.
		// Simplifies querying extended fields.
		var headerkey = strings.ToLower(strings.ReplaceAll(header, "-", ""))
		fields[headerkey] = r.Header.Get(header)
	}
	fields["subsys"] = "http"
	fields["path"] = r.URL.String()
	fields["method"] = r.Method
	fields["ip"] = r.RemoteAddr
	fields["host"] = r.Host
	fields["useragent"] = r.Header.Get("User-Agent")
	l := log.Ctx(r.Context()).WithFields(fields)

	l.Info("starting request")
}

// logEndOfRequest emits the logline for the end of the request
func logEndOfRequest(
	r *stdhttp.Request,
	duration time.Duration,
	mw mutil.WriterProxy,
) {
	l := log.Ctx(r.Context()).WithFields(log.F{
		"subsys":   "http",
		"path":     r.URL.String(),
		"method":   r.Method,
		"status":   mw.Status(),
		"bytes":    mw.BytesWritten(),
		"duration": duration,
	})
	if routeContext := chi.RouteContext(r.Context()); routeContext != nil {
		l = l.WithField("route", routeContext.RoutePattern())
	}
	l.Info("finished request")
}
