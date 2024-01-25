package httpx

import (
	"compress/flate"
	"fmt"
	"net/http"
	"net/http/pprof"
	"net/url"
	"time"

	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"github.com/diamcircle/throttled"

	"github.com/diamcircle/go/services/aurora/internal/actions"
	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/services/aurora/internal/ledger"
	"github.com/diamcircle/go/services/aurora/internal/paths"
	"github.com/diamcircle/go/services/aurora/internal/render/sse"
	"github.com/diamcircle/go/services/aurora/internal/txsub"
	"github.com/diamcircle/go/support/db"
	supporthttp "github.com/diamcircle/go/support/http"
	"github.com/diamcircle/go/support/render/problem"
)

type RouterConfig struct {
	DBSession        db.SessionInterface
	PrimaryDBSession db.SessionInterface
	TxSubmitter      *txsub.System
	RateQuota        *throttled.RateQuota

	BehindCloudflare        bool
	BehindAWSLoadBalancer   bool
	SSEUpdateFrequency      time.Duration
	StaleThreshold          uint
	ConnectionTimeout       time.Duration
	NetworkPassphrase       string
	MaxPathLength           uint
	MaxAssetsPerPathRequest int
	PathFinder              paths.Finder
	PrometheusRegistry      *prometheus.Registry
	CoreGetter              actions.CoreStateGetter
	AuroraVersion          string
	FriendbotURL            *url.URL
	HealthCheck             http.Handler
}

type Router struct {
	*chi.Mux
	Internal *chi.Mux
}

func NewRouter(config *RouterConfig, serverMetrics *ServerMetrics, ledgerState *ledger.State) (*Router, error) {
	result := Router{
		Mux:      chi.NewMux(),
		Internal: chi.NewMux(),
	}
	var rateLimiter *throttled.HTTPRateLimiter
	if config.RateQuota != nil {
		var err error
		rateLimiter, err = newRateLimiter(config.RateQuota)
		if err != nil {
			return nil, fmt.Errorf("unable to create RateLimiter: %v", err)
		}
	}
	result.addMiddleware(config, rateLimiter, serverMetrics)
	result.addRoutes(config, rateLimiter, ledgerState)
	return &result, nil
}

func (r *Router) addMiddleware(config *RouterConfig,
	rateLimitter *throttled.HTTPRateLimiter,
	serverMetrics *ServerMetrics) {

	r.Use(chimiddleware.StripSlashes)

	r.Use(requestCacheHeadersMiddleware)
	r.Use(chimiddleware.RequestID)
	r.Use(contextMiddleware)
	r.Use(supporthttp.XFFMiddleware(supporthttp.XFFMiddlewareConfig{
		BehindCloudflare:      config.BehindCloudflare,
		BehindAWSLoadBalancer: config.BehindAWSLoadBalancer,
	}))
	r.Use(loggerMiddleware(serverMetrics))
	r.Use(timeoutMiddleware(config.ConnectionTimeout))
	r.Use(recoverMiddleware)
	r.Use(chimiddleware.Compress(flate.DefaultCompression, "application/hal+json"))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Date", "Latest-Ledger"},
	})
	r.Use(c.Handler)

	if rateLimitter != nil {
		r.Use(rateLimitter.RateLimit)
	}

	if config.PrimaryDBSession != nil {
		replicaSyncMiddleware := ReplicaSyncCheckMiddleware{
			PrimaryHistoryQ: &history.Q{config.PrimaryDBSession},
			ReplicaHistoryQ: &history.Q{config.DBSession},
			ServerMetrics:   serverMetrics,
		}
		r.Use(replicaSyncMiddleware.Wrap)
	}

	// Internal middlewares
	r.Internal.Use(chimiddleware.StripSlashes)
	r.Internal.Use(chimiddleware.RequestID)
	r.Internal.Use(loggerMiddleware(serverMetrics))
}

func (r *Router) addRoutes(config *RouterConfig, rateLimiter *throttled.HTTPRateLimiter, ledgerState *ledger.State) {
	stateMiddleware := StateMiddleware{
		AuroraSession: config.DBSession,
	}

	r.Method(http.MethodGet, "/health", config.HealthCheck)

	r.Method(http.MethodGet, "/", ObjectActionHandler{Action: actions.GetRootHandler{
		LedgerState:       ledgerState,
		CoreStateGetter:   config.CoreGetter,
		NetworkPassphrase: config.NetworkPassphrase,
		FriendbotURL:      config.FriendbotURL,
		AuroraVersion:    config.AuroraVersion,
	}})

	streamHandler := sse.StreamHandler{
		RateLimiter:         rateLimiter,
		LedgerSourceFactory: historyLedgerSourceFactory{ledgerState: ledgerState, updateFrequency: config.SSEUpdateFrequency},
	}

	historyMiddleware := NewHistoryMiddleware(ledgerState, int32(config.StaleThreshold), config.DBSession)
	// State endpoints behind stateMiddleware
	r.Group(func(r chi.Router) {
		r.Route("/accounts", func(r chi.Router) {
			r.With(stateMiddleware.Wrap).Method(http.MethodGet, "/", restPageHandler(ledgerState, actions.GetAccountsHandler{LedgerState: ledgerState}))
			r.Route("/{account_id}", func(r chi.Router) {
				r.With(stateMiddleware.Wrap).Method(
					http.MethodGet,
					"/",
					streamableObjectActionHandler{
						streamHandler: streamHandler,
						action:        actions.GetAccountByIDHandler{},
					},
				)
				accountData := actions.GetAccountDataHandler{}
				r.With(stateMiddleware.Wrap).Method(http.MethodGet, "/data/{key}", WrapRaw(
					streamableObjectActionHandler{streamHandler: streamHandler, action: accountData},
					accountData,
				))
				r.With(stateMiddleware.Wrap).Method(http.MethodGet, "/offers", streamableStatePageHandler(ledgerState, actions.GetAccountOffersHandler{LedgerState: ledgerState}, streamHandler))
			})
		})

		r.Route("/claimable_balances", func(r chi.Router) {
			r.With(stateMiddleware.Wrap).Method(http.MethodGet, "/", restPageHandler(ledgerState, actions.GetClaimableBalancesHandler{LedgerState: ledgerState}))
			r.With(stateMiddleware.Wrap).Method(http.MethodGet, "/{id}", ObjectActionHandler{actions.GetClaimableBalanceByIDHandler{}})
		})

		r.Route("/liquidity_pools", func(r chi.Router) {
			r.With(stateMiddleware.Wrap).Method(http.MethodGet, "/", restPageHandler(ledgerState, actions.GetLiquidityPoolsHandler{LedgerState: ledgerState}))
			r.Route("/{liquidity_pool_id:\\w+}", func(r chi.Router) {
				r.With(stateMiddleware.Wrap).Method(http.MethodGet, "/", ObjectActionHandler{actions.GetLiquidityPoolByIDHandler{}})
				r.With(historyMiddleware).Method(http.MethodGet, "/operations", streamableHistoryPageHandler(ledgerState, actions.GetOperationsHandler{
					LedgerState:  ledgerState,
					OnlyPayments: false,
				}, streamHandler))
				r.With(historyMiddleware).Method(http.MethodGet, "/transactions", streamableHistoryPageHandler(ledgerState, actions.GetTransactionsHandler{LedgerState: ledgerState}, streamHandler))
				r.With(historyMiddleware).Method(http.MethodGet, "/effects", streamableHistoryPageHandler(ledgerState, actions.GetEffectsHandler{LedgerState: ledgerState}, streamHandler))
				r.With(historyMiddleware).Method(http.MethodGet, "/trades", streamableHistoryPageHandler(ledgerState, actions.GetTradesHandler{LedgerState: ledgerState, CoreStateGetter: config.CoreGetter}, streamHandler))
			})
		})

		r.Route("/offers", func(r chi.Router) {
			r.With(stateMiddleware.Wrap).Method(http.MethodGet, "/", restPageHandler(ledgerState, actions.GetOffersHandler{LedgerState: ledgerState}))
			r.With(stateMiddleware.Wrap).Method(http.MethodGet, "/{offer_id}", ObjectActionHandler{actions.GetOfferByID{}})
		})

		r.With(stateMiddleware.Wrap).Method(http.MethodGet, "/assets", restPageHandler(ledgerState, actions.AssetStatsHandler{LedgerState: ledgerState}))

		findPaths := ObjectActionHandler{actions.FindPathsHandler{
			StaleThreshold:       config.StaleThreshold,
			SetLastLedgerHeader:  true,
			MaxPathLength:        config.MaxPathLength,
			MaxAssetsParamLength: config.MaxAssetsPerPathRequest,
			PathFinder:           config.PathFinder,
		}}
		findFixedPaths := ObjectActionHandler{actions.FindFixedPathsHandler{
			MaxPathLength:        config.MaxPathLength,
			SetLastLedgerHeader:  true,
			MaxAssetsParamLength: config.MaxAssetsPerPathRequest,
			PathFinder:           config.PathFinder,
		}}
		r.With(stateMiddleware.Wrap).Method(http.MethodGet, "/paths", findPaths)
		r.With(stateMiddleware.Wrap).Method(http.MethodGet, "/paths/strict-receive", findPaths)
		r.With(stateMiddleware.Wrap).Method(http.MethodGet, "/paths/strict-send", findFixedPaths)
		r.With(stateMiddleware.Wrap).Method(
			http.MethodGet,
			"/order_book",
			streamableObjectActionHandler{
				streamHandler: streamHandler,
				action:        actions.GetOrderbookHandler{},
			},
		)
	})

	// account actions - /accounts/{account_id} has been created above so we
	// need to use absolute routes here. Make sure we use regexp check here for
	// emptiness. Without it, requesting `/accounts//payments` return all payments!
	r.Group(func(r chi.Router) {
		r.With(historyMiddleware).Method(http.MethodGet, "/accounts/{account_id:\\w+}/effects", streamableHistoryPageHandler(ledgerState, actions.GetEffectsHandler{LedgerState: ledgerState}, streamHandler))
		r.With(historyMiddleware).Method(http.MethodGet, "/accounts/{account_id:\\w+}/operations", streamableHistoryPageHandler(ledgerState, actions.GetOperationsHandler{
			LedgerState:  ledgerState,
			OnlyPayments: false,
		}, streamHandler))
		r.With(historyMiddleware).Method(http.MethodGet, "/accounts/{account_id:\\w+}/payments", streamableHistoryPageHandler(ledgerState, actions.GetOperationsHandler{
			LedgerState:  ledgerState,
			OnlyPayments: true,
		}, streamHandler))
		r.With(historyMiddleware).Method(http.MethodGet, "/accounts/{account_id:\\w+}/trades", streamableHistoryPageHandler(ledgerState, actions.GetTradesHandler{LedgerState: ledgerState, CoreStateGetter: config.CoreGetter}, streamHandler))
		r.With(historyMiddleware).Method(http.MethodGet, "/accounts/{account_id:\\w+}/transactions", streamableHistoryPageHandler(ledgerState, actions.GetTransactionsHandler{LedgerState: ledgerState}, streamHandler))
	})
	// ledger actions
	r.Route("/ledgers", func(r chi.Router) {
		r.With(historyMiddleware).Method(http.MethodGet, "/", streamableHistoryPageHandler(ledgerState, actions.GetLedgersHandler{LedgerState: ledgerState}, streamHandler))
		r.Route("/{ledger_id}", func(r chi.Router) {
			r.With(historyMiddleware).Method(http.MethodGet, "/", ObjectActionHandler{actions.GetLedgerByIDHandler{LedgerState: ledgerState}})
			r.With(historyMiddleware).Method(http.MethodGet, "/transactions", streamableHistoryPageHandler(ledgerState, actions.GetTransactionsHandler{LedgerState: ledgerState}, streamHandler))
			r.Group(func(r chi.Router) {
				r.With(historyMiddleware).Method(http.MethodGet, "/effects", streamableHistoryPageHandler(ledgerState, actions.GetEffectsHandler{LedgerState: ledgerState}, streamHandler))
				r.With(historyMiddleware).Method(http.MethodGet, "/operations", streamableHistoryPageHandler(ledgerState, actions.GetOperationsHandler{
					LedgerState:  ledgerState,
					OnlyPayments: false,
				}, streamHandler))
				r.With(historyMiddleware).Method(http.MethodGet, "/payments", streamableHistoryPageHandler(ledgerState, actions.GetOperationsHandler{
					LedgerState:  ledgerState,
					OnlyPayments: true,
				}, streamHandler))
			})
		})
	})

	// claimable balance actions
	r.Group(func(r chi.Router) {
		r.With(historyMiddleware).Method(http.MethodGet, "/claimable_balances/{claimable_balance_id:\\w+}/operations", streamableHistoryPageHandler(ledgerState, actions.GetOperationsHandler{
			LedgerState:  ledgerState,
			OnlyPayments: false,
		}, streamHandler))
		r.With(historyMiddleware).Method(http.MethodGet, "/claimable_balances/{claimable_balance_id:\\w+}/transactions", streamableHistoryPageHandler(ledgerState, actions.GetTransactionsHandler{LedgerState: ledgerState}, streamHandler))
	})

	// transaction history actions
	r.Route("/transactions", func(r chi.Router) {
		r.With(historyMiddleware).Method(http.MethodGet, "/", streamableHistoryPageHandler(ledgerState, actions.GetTransactionsHandler{LedgerState: ledgerState}, streamHandler))
		r.Route("/{tx_id}", func(r chi.Router) {
			r.With(historyMiddleware).Method(http.MethodGet, "/", ObjectActionHandler{actions.GetTransactionByHashHandler{}})
			r.With(historyMiddleware).Method(http.MethodGet, "/effects", streamableHistoryPageHandler(ledgerState, actions.GetEffectsHandler{LedgerState: ledgerState}, streamHandler))
			r.With(historyMiddleware).Method(http.MethodGet, "/operations", streamableHistoryPageHandler(ledgerState, actions.GetOperationsHandler{
				LedgerState:  ledgerState,
				OnlyPayments: false,
			}, streamHandler))
			r.With(historyMiddleware).Method(http.MethodGet, "/payments", streamableHistoryPageHandler(ledgerState, actions.GetOperationsHandler{
				LedgerState:  ledgerState,
				OnlyPayments: true,
			}, streamHandler))
		})
	})

	// operation actions
	r.Route("/operations", func(r chi.Router) {
		r.With(historyMiddleware).Method(http.MethodGet, "/", streamableHistoryPageHandler(ledgerState, actions.GetOperationsHandler{
			LedgerState:  ledgerState,
			OnlyPayments: false,
		}, streamHandler))
		r.With(historyMiddleware).Method(http.MethodGet, "/{id}", ObjectActionHandler{actions.GetOperationByIDHandler{LedgerState: ledgerState}})
		r.With(historyMiddleware).Method(http.MethodGet, "/{op_id}/effects", streamableHistoryPageHandler(ledgerState, actions.GetEffectsHandler{LedgerState: ledgerState}, streamHandler))
	})

	r.Group(func(r chi.Router) {
		// payment actions
		r.With(historyMiddleware).Method(http.MethodGet, "/payments", streamableHistoryPageHandler(ledgerState, actions.GetOperationsHandler{
			LedgerState:  ledgerState,
			OnlyPayments: true,
		}, streamHandler))

		// effect actions
		r.With(historyMiddleware).Method(http.MethodGet, "/effects", streamableHistoryPageHandler(ledgerState, actions.GetEffectsHandler{LedgerState: ledgerState}, streamHandler))

		// trading related endpoints
		r.With(historyMiddleware).Method(http.MethodGet, "/trades", streamableHistoryPageHandler(ledgerState, actions.GetTradesHandler{LedgerState: ledgerState, CoreStateGetter: config.CoreGetter}, streamHandler))
		r.With(historyMiddleware).Method(http.MethodGet, "/trade_aggregations", ObjectActionHandler{actions.GetTradeAggregationsHandler{LedgerState: ledgerState, CoreStateGetter: config.CoreGetter}})
		// /offers/{offer_id} has been created above so we need to use absolute
		// routes here.
		r.With(historyMiddleware).Method(http.MethodGet, "/offers/{offer_id}/trades", streamableHistoryPageHandler(ledgerState, actions.GetTradesHandler{LedgerState: ledgerState, CoreStateGetter: config.CoreGetter}, streamHandler))
	})

	// Transaction submission API
	r.Method(http.MethodPost, "/transactions", ObjectActionHandler{actions.SubmitTransactionHandler{
		Submitter:         config.TxSubmitter,
		NetworkPassphrase: config.NetworkPassphrase,
		CoreStateGetter:   config.CoreGetter,
	}})

	// Network state related endpoints
	r.Method(http.MethodGet, "/fee_stats", ObjectActionHandler{actions.FeeStatsHandler{}})

	// friendbot
	if config.FriendbotURL != nil {
		redirectFriendbot := func(w http.ResponseWriter, r *http.Request) {
			redirectURL := config.FriendbotURL.String() + "?" + r.URL.RawQuery
			http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
		}
		r.Post("/friendbot", redirectFriendbot)
		r.Get("/friendbot", redirectFriendbot)
	}

	r.NotFound(func(w http.ResponseWriter, request *http.Request) {
		problem.Render(request.Context(), w, problem.NotFound)
	})

	// internal
	r.Internal.Get("/metrics", promhttp.HandlerFor(config.PrometheusRegistry, promhttp.HandlerOpts{}).ServeHTTP)
	r.Internal.Get("/debug/pprof/heap", pprof.Index)
	r.Internal.Get("/debug/pprof/profile", pprof.Profile)
}
