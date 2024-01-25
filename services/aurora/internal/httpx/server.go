package httpx

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/diamcircle/go/services/aurora/internal/db2"
	"github.com/diamcircle/go/services/aurora/internal/ledger"
	hProblem "github.com/diamcircle/go/services/aurora/internal/render/problem"
	"github.com/diamcircle/go/services/aurora/internal/render/sse"
	"github.com/diamcircle/go/services/aurora/internal/txsub/sequence"
	"github.com/diamcircle/go/support/db"
	"github.com/diamcircle/go/support/log"
	"github.com/diamcircle/go/support/render/problem"
)

type ServerMetrics struct {
	RequestDurationSummary  *prometheus.SummaryVec
	ReplicaLagErrorsCounter prometheus.Counter
}

type TLSConfig struct {
	CertPath, KeyPath string
}
type ServerConfig struct {
	Port      uint16
	TLSConfig *TLSConfig
	AdminPort uint16
}

// Server contains the http server related fields for aurora: the Router,
// rate limiter, etc.
type Server struct {
	Router         *Router
	Metrics        *ServerMetrics
	server         *http.Server
	config         ServerConfig
	internalServer *http.Server
}

func init() {
	// register problems
	problem.SetLogFilter(problem.LogUnknownErrors)
	problem.RegisterError(sql.ErrNoRows, problem.NotFound)
	problem.RegisterError(sequence.ErrNoMoreRoom, hProblem.ServerOverCapacity)
	problem.RegisterError(db2.ErrInvalidCursor, problem.BadRequest)
	problem.RegisterError(db2.ErrInvalidLimit, problem.BadRequest)
	problem.RegisterError(db2.ErrInvalidOrder, problem.BadRequest)
	problem.RegisterError(sse.ErrRateLimited, hProblem.RateLimitExceeded)
	problem.RegisterError(context.DeadlineExceeded, hProblem.Timeout)
	problem.RegisterError(context.Canceled, hProblem.ClientDisconnected)
	problem.RegisterError(db.ErrCancelled, hProblem.ClientDisconnected)
	problem.RegisterError(db.ErrTimeout, hProblem.ServiceUnavailable)
	problem.RegisterError(db.ErrConflictWithRecovery, hProblem.ServiceUnavailable)
	problem.RegisterError(db.ErrBadConnection, hProblem.ServiceUnavailable)
}

func NewServer(serverConfig ServerConfig, routerConfig RouterConfig, ledgerState *ledger.State) (*Server, error) {
	sm := &ServerMetrics{
		RequestDurationSummary: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace: "aurora", Subsystem: "http", Name: "requests_duration_seconds",
				Help: "HTTP requests durations, sliding window = 10m",
			},
			[]string{"status", "route", "streaming", "method"},
		),
		ReplicaLagErrorsCounter: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: "aurora", Subsystem: "http", Name: "replica_lag_errors_count",
				Help: "Count of HTTP errors returned due to replica lag",
			},
		),
	}
	router, err := NewRouter(&routerConfig, sm, ledgerState)
	if err != nil {
		return nil, err
	}
	addr := fmt.Sprintf(":%d", serverConfig.Port)
	result := &Server{
		Router:  router,
		Metrics: sm,
		config:  serverConfig,
		server: &http.Server{
			Addr:        addr,
			Handler:     router,
			ReadTimeout: 5 * time.Second,
		},
	}

	if serverConfig.AdminPort != 0 {
		adminAddr := fmt.Sprintf(":%d", serverConfig.AdminPort)
		result.internalServer = &http.Server{
			Addr:        adminAddr,
			Handler:     result.Router.Internal,
			ReadTimeout: 5 * time.Second,
		}
	}
	return result, nil
}

// RegisterMetrics registers the prometheus metrics
func (s *Server) RegisterMetrics(registry *prometheus.Registry) {
	registry.MustRegister(s.Metrics.RequestDurationSummary)
	registry.MustRegister(s.Metrics.ReplicaLagErrorsCounter)
}

func (s *Server) Serve() error {
	if s.internalServer != nil {
		go func() {
			log.Infof("Starting internal server on %s", s.internalServer.Addr)
			if err := s.internalServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Warn(errors.Wrap(err, "error in internalServer.ListenAndServe()"))
			}
		}()
	}

	var err error
	if s.config.TLSConfig != nil {
		err = s.server.ListenAndServeTLS(s.config.TLSConfig.CertPath, s.config.TLSConfig.KeyPath)
	} else {
		err = s.server.ListenAndServe()
	}
	return err
}

func (s *Server) Shutdown(ctx context.Context) error {
	var wg sync.WaitGroup
	defer wg.Wait()
	if s.internalServer != nil {
		wg.Add(1)
		go func() {
			err := s.internalServer.Shutdown(ctx)
			if err != nil {
				log.Warn(errors.Wrap(err, "error in internalServer.Shutdown()"))
			}
			wg.Done()
		}()
	}
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}
