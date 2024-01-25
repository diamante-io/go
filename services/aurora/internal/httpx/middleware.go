package httpx

import (
	"context"
	"database/sql"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/diamcircle/go/services/aurora/internal/actions"
	auroraContext "github.com/diamcircle/go/services/aurora/internal/context"
	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/services/aurora/internal/errors"
	"github.com/diamcircle/go/services/aurora/internal/hchi"
	"github.com/diamcircle/go/services/aurora/internal/ingest"
	"github.com/diamcircle/go/services/aurora/internal/ledger"
	"github.com/diamcircle/go/services/aurora/internal/render"
	hProblem "github.com/diamcircle/go/services/aurora/internal/render/problem"
	"github.com/diamcircle/go/support/db"
	supportErrors "github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/support/log"
	"github.com/diamcircle/go/support/render/problem"
)

// requestCacheHeadersMiddleware adds caching headers to each response.
func requestCacheHeadersMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Before changing this read Stack Overflow answer about staled request
		// in older versions of Chrome:
		// https://stackoverflow.com/questions/27513994/chrome-stalls-when-making-multiple-requests-to-same-resource
		w.Header().Set("Cache-Control", "no-cache, no-store, max-age=0")
		h.ServeHTTP(w, r)
	})
}

func contextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = hchi.WithChiRequestID(ctx)
		ctx = auroraContext.RequestContext(ctx, w, r)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

const (
	clientNameHeader    = "X-Client-Name"
	clientVersionHeader = "X-Client-Version"
	appNameHeader       = "X-App-Name"
	appVersionHeader    = "X-App-Version"
)

func newWrapResponseWriter(w http.ResponseWriter, r *http.Request) middleware.WrapResponseWriter {
	mw, ok := w.(middleware.WrapResponseWriter)
	if !ok {
		mw = middleware.NewWrapResponseWriter(w, r.ProtoMajor)
	}

	return mw
}

// loggerMiddleware logs http requests and resposnes to the logging subsytem of aurora.
func loggerMiddleware(serverMetrics *ServerMetrics) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			mw := newWrapResponseWriter(w, r)

			logger := log.WithField("req", middleware.GetReqID(ctx))
			ctx = log.Set(ctx, logger)

			// Checking `Accept` header from user request because if the streaming connection
			// is reset before sending the first event no Content-Type header is sent in a response.
			acceptHeader := r.Header.Get("Accept")
			streaming := strings.Contains(acceptHeader, render.MimeEventStream)

			then := time.Now()
			next.ServeHTTP(mw, r.WithContext(ctx))
			duration := time.Since(then)
			logEndOfRequest(ctx, r, serverMetrics.RequestDurationSummary, duration, mw, streaming)
		})
	}
}

// timeoutMiddleware ensures the request is terminated after the given timeout
func timeoutMiddleware(timeout time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			mw := newWrapResponseWriter(w, r)
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer func() {
				cancel()
				if ctx.Err() == context.DeadlineExceeded {
					if mw.Status() == 0 {
						// only write the header if it hasn't been written yet
						mw.WriteHeader(http.StatusGatewayTimeout)
					}
				}
			}()

			// txsub has a custom timeout
			if r.Method != http.MethodPost {
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(mw, r)
		}
		return http.HandlerFunc(fn)
	}
}

// getClientData gets client data (name or version) from header or GET parameter
// (useful when not possible to set headers, like in EventStream).
func getClientData(r *http.Request, headerName string) string {
	value := r.Header.Get(headerName)
	if value != "" {
		return value
	}

	value = r.URL.Query().Get(headerName)
	if value == "" {
		value = "undefined"
	}

	return value
}

var routeRegexp = regexp.MustCompile("{([^:}]*):[^}]*}")

// https://prometheus.io/docs/instrumenting/exposition_formats/
// label_value can be any sequence of UTF-8 characters, but the backslash (\),
// double-quote ("), and line feed (\n) characters have to be escaped as \\,
// \", and \n, respectively.
func sanitizeMetricRoute(routePattern string) string {
	route := routeRegexp.ReplaceAllString(routePattern, "{$1}")
	route = strings.ReplaceAll(route, "\\", "\\\\")
	route = strings.ReplaceAll(route, "\"", "\\\"")
	route = strings.ReplaceAll(route, "\n", "\\n")
	if route == "" {
		// Can be empty when request did not reach the final route (ex. blocked by
		// a middleware). More info: https://github.com/go-chi/chi/issues/270
		return "undefined"
	}
	return route
}

// Author: https://github.com/rliebz
// From: https://github.com/go-chi/chi/issues/270#issuecomment-479184559
// https://github.com/go-chi/chi/blob/master/LICENSE
func getRoutePattern(r *http.Request) string {
	rctx := chi.RouteContext(r.Context())
	if pattern := rctx.RoutePattern(); pattern != "" {
		// Pattern is already available
		return pattern
	}

	routePath := r.URL.Path
	if r.URL.RawPath != "" {
		routePath = r.URL.RawPath
	}

	tctx := chi.NewRouteContext()
	if !rctx.Routes.Match(tctx, r.Method, routePath) {
		return ""
	}

	// tctx has the updated pattern, since Match mutates it
	return tctx.RoutePattern()
}

func logEndOfRequest(ctx context.Context, r *http.Request, requestDurationSummary *prometheus.SummaryVec, duration time.Duration, mw middleware.WrapResponseWriter, streaming bool) {
	route := sanitizeMetricRoute(getRoutePattern(r))

	referer := r.Referer()
	if referer == "" {
		referer = r.Header.Get("Origin")
	}
	if referer == "" {
		referer = "undefined"
	}

	log.Ctx(ctx).WithFields(log.F{
		"bytes":           mw.BytesWritten(),
		"client_name":     getClientData(r, clientNameHeader),
		"client_version":  getClientData(r, clientVersionHeader),
		"app_name":        getClientData(r, appNameHeader),
		"app_version":     getClientData(r, appVersionHeader),
		"duration":        duration.Seconds(),
		"x_forwarder_for": r.Header.Get("X-Forwarded-For"),
		"host":            r.Host,
		"ip":              remoteAddrIP(r),
		"ip_port":         r.RemoteAddr,
		"method":          r.Method,
		"path":            r.URL.String(),
		"route":           route,
		"status":          mw.Status(),
		"streaming":       streaming,
		"referer":         referer,
	}).Info("Finished request")

	requestDurationSummary.With(prometheus.Labels{
		"status":    strconv.FormatInt(int64(mw.Status()), 10),
		"route":     route,
		"streaming": strconv.FormatBool(streaming),
		"method":    r.Method,
	}).Observe(float64(duration.Seconds()))
}

// recoverMiddleware helps the server recover from panics. It ensures that
// no request can fully bring down the aurora server, and it also logs the
// panics to the logging subsystem.
func recoverMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		defer func() {
			if rec := recover(); rec != nil {
				err := errors.FromPanic(rec)
				errors.ReportToSentry(err, r)
				problem.Render(ctx, w, err)
			}
		}()

		h.ServeHTTP(w, r)
	})
}

// NewHistoryMiddleware adds session to the request context and ensures Aurora
// is not in a stale state, which is when the difference between latest core
// ledger and latest history ledger is higher than the given threshold
func NewHistoryMiddleware(ledgerState *ledger.State, staleThreshold int32, session db.SessionInterface) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			chiRoute := chi.RouteContext(ctx)
			if chiRoute != nil {
				ctx = context.WithValue(ctx, &db.RouteContextKey, sanitizeMetricRoute(chiRoute.RoutePattern()))
			}
			if staleThreshold > 0 {
				ls := ledgerState.CurrentStatus()
				isStale := (ls.CoreLatest - ls.HistoryLatest) > int32(staleThreshold)
				if isStale {
					err := hProblem.StaleHistory
					err.Extras = map[string]interface{}{
						"history_latest_ledger": ls.HistoryLatest,
						"core_latest_ledger":    ls.CoreLatest,
					}
					problem.Render(ctx, w, err)
					return
				}
			}

			requestSession := session.Clone()
			h.ServeHTTP(w, r.WithContext(
				context.WithValue(
					ctx,
					&auroraContext.SessionContextKey,
					requestSession,
				),
			))
		})
	}
}

// StateMiddleware is a middleware which enables a state handler if the state
// has been initialized.
// Unless NoStateVerification is set, it ensures that the state (ledger entries)
// has been verified and is correct (Otherwise returns `500 Internal Server Error` to prevent
// returning invalid data to the user)
type StateMiddleware struct {
	AuroraSession      db.SessionInterface
	NoStateVerification bool
}

func ingestionStatus(ctx context.Context, q *history.Q) (uint32, bool, error) {
	version, err := q.GetIngestVersion(ctx)
	if err != nil {
		return 0, false, supportErrors.Wrap(
			err, "Error running GetIngestVersion",
		)
	}

	lastIngestedLedger, err := q.GetLastLedgerIngestNonBlocking(ctx)
	if err != nil {
		return 0, false, supportErrors.Wrap(
			err, "Error running GetLastLedgerIngestNonBlocking",
		)
	}

	var lastHistoryLedger uint32
	err = q.LatestLedger(ctx, &lastHistoryLedger)
	if err != nil {
		return 0, false, supportErrors.Wrap(err, "Error running LatestLedger")
	}

	ready := version == ingest.CurrentVersion &&
		lastIngestedLedger > 0 &&
		lastIngestedLedger == lastHistoryLedger

	return lastIngestedLedger, ready, nil
}

// WrapFunc executes the middleware on a given HTTP handler function
func (m *StateMiddleware) WrapFunc(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		chiRoute := chi.RouteContext(ctx)
		if chiRoute != nil {
			ctx = context.WithValue(ctx, &db.RouteContextKey, sanitizeMetricRoute(chiRoute.RoutePattern()))
		}
		session := m.AuroraSession.Clone()
		q := &history.Q{session}
		sseRequest := render.Negotiate(r) == render.MimeEventStream

		// We want to start a repeatable read session to ensure that the data we
		// fetch from the db belong to the same ledger.
		// Otherwise, because the ingestion system is running concurrently with this request,
		// it is possible to have one read fetch data from ledger N and another read
		// fetch data from ledger N+1 .
		err := session.BeginTx(&sql.TxOptions{
			Isolation: sql.LevelRepeatableRead,
			ReadOnly:  true,
		})
		if err != nil {
			err = supportErrors.Wrap(err, "Error starting ingestion read transaction")
			problem.Render(ctx, w, err)
			return
		}
		defer session.Rollback()

		if !m.NoStateVerification {
			stateInvalid, invalidErr := q.GetExpStateInvalid(ctx)
			if invalidErr != nil {
				invalidErr = supportErrors.Wrap(invalidErr, "Error running GetExpStateInvalid")
				problem.Render(ctx, w, invalidErr)
				return
			}
			if stateInvalid {
				problem.Render(ctx, w, problem.ServerError)
				return
			}
		}

		lastIngestedLedger, ready, err := ingestionStatus(ctx, q)
		if err != nil {
			problem.Render(ctx, w, err)
			return
		}
		if !m.NoStateVerification && !ready {
			problem.Render(ctx, w, hProblem.StillIngesting)
			return
		}

		// for SSE requests we need to discard the repeatable read transaction
		// otherwise, the stream will not pick up updates occurring in future
		// ledgers
		if sseRequest {
			if err = session.Rollback(); err != nil {
				problem.Render(
					ctx,
					w,
					supportErrors.Wrap(
						err,
						"Could not roll back repeatable read session for SSE request",
					),
				)
				return
			}
		} else {
			actions.SetLastLedgerHeader(w, lastIngestedLedger)
		}

		h.ServeHTTP(w, r.WithContext(
			context.WithValue(ctx, &auroraContext.SessionContextKey, session),
		))
	}
}

// WrapFunc executes the middleware on a given HTTP handler function
func (m *StateMiddleware) Wrap(h http.Handler) http.Handler {
	return m.WrapFunc(h.ServeHTTP)
}

type ReplicaSyncCheckMiddleware struct {
	PrimaryHistoryQ *history.Q
	ReplicaHistoryQ *history.Q
	ServerMetrics   *ServerMetrics
}

// WrapFunc executes the middleware on a given HTTP handler function
func (m *ReplicaSyncCheckMiddleware) WrapFunc(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for attempt := 1; attempt <= 4; attempt++ {
			primaryIngestLedger, err := m.PrimaryHistoryQ.GetLastLedgerIngestNonBlocking(r.Context())
			if err != nil {
				problem.Render(r.Context(), w, err)
				return
			}

			replicaIngestLedger, err := m.ReplicaHistoryQ.GetLastLedgerIngestNonBlocking(r.Context())
			if err != nil {
				problem.Render(r.Context(), w, err)
				return
			}

			if replicaIngestLedger >= primaryIngestLedger {
				break
			}

			switch attempt {
			case 1:
				time.Sleep(20 * time.Millisecond)
			case 2:
				time.Sleep(40 * time.Millisecond)
			case 3:
				time.Sleep(80 * time.Millisecond)
			case 4:
				problem.Render(r.Context(), w, hProblem.StaleHistory)
				m.ServerMetrics.ReplicaLagErrorsCounter.Inc()
				return
			}
		}

		h.ServeHTTP(w, r)
	}
}

func (m *ReplicaSyncCheckMiddleware) Wrap(h http.Handler) http.Handler {
	return m.WrapFunc(h.ServeHTTP)
}
