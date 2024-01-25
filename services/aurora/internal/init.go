package aurora

import (
	"context"
	"net/http"
	"runtime"

	"go/exp/orderbook"
	"go/services/aurora/internal/db2/history"
	"go/services/aurora/internal/ingest"
	"go/services/aurora/internal/simplepath"
	"go/services/aurora/internal/txsub"
	"go/services/aurora/internal/txsub/sequence"
	"go/support/db"
	"go/support/log"

	"github.com/getsentry/raven-go"
	"github.com/prometheus/client_golang/prometheus"
)

func mustNewDBSession(subservice db.Subservice, databaseURL string, maxIdle, maxOpen int, registry *prometheus.Registry) db.SessionInterface {
	session, err := db.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("cannot open %v DB: %v", subservice, err)
	}

	session.DB.SetMaxIdleConns(maxIdle)
	session.DB.SetMaxOpenConns(maxOpen)
	return db.RegisterMetrics(session, "aurora", subservice, registry)
}

func mustInitauroraDB(app *App) {
	maxIdle := app.config.auroraDBMaxIdleConnections
	maxOpen := app.config.auroraDBMaxOpenConnections
	if app.config.Ingest {
		maxIdle -= ingest.MaxDBConnections
		maxOpen -= ingest.MaxDBConnections
		if maxIdle <= 0 {
			log.Fatalf("max idle connections to aurora db must be greater than %d", ingest.MaxDBConnections)
		}
		if maxOpen <= 0 {
			log.Fatalf("max open connections to aurora db must be greater than %d", ingest.MaxDBConnections)
		}
	}

	if app.config.RoDatabaseURL == "" {
		app.historyQ = &history.Q{mustNewDBSession(
			db.HistorySubservice,
			app.config.DatabaseURL,
			maxIdle,
			maxOpen,
			app.prometheusRegistry,
		)}
	} else {
		// If RO set, use it for all DB queries
		app.historyQ = &history.Q{mustNewDBSession(
			db.HistorySubservice,
			app.config.RoDatabaseURL,
			maxIdle,
			maxOpen,
			app.prometheusRegistry,
		)}

		app.primaryHistoryQ = &history.Q{mustNewDBSession(
			db.HistoryPrimarySubservice,
			app.config.DatabaseURL,
			maxIdle,
			maxOpen,
			app.prometheusRegistry,
		)}
	}
}

func initIngester(app *App) {
	var err error
	var coreSession db.SessionInterface
	if !app.config.EnableCaptiveCoreIngestion {
		coreSession = mustNewDBSession(
			db.CoreSubservice, app.config.diamcircleCoreDatabaseURL, ingest.MaxDBConnections, ingest.MaxDBConnections, app.prometheusRegistry)
	}
	app.ingester, err = ingest.NewSystem(ingest.Config{
		CoreSession: coreSession,
		HistorySession: mustNewDBSession(
			db.IngestSubservice, app.config.DatabaseURL, ingest.MaxDBConnections, ingest.MaxDBConnections, app.prometheusRegistry,
		),
		NetworkPassphrase: app.config.NetworkPassphrase,
		// TODO:
		// Use the first archive for now. We don't have a mechanism to
		// use multiple archives at the same time currently.
		HistoryArchiveURL:            app.config.HistoryArchiveURLs[0],
		CheckpointFrequency:          app.config.CheckpointFrequency,
		diamcircleCoreURL:            app.config.diamcircleCoreURL,
		diamcircleCoreCursor:         app.config.CursorName,
		CaptiveCoreBinaryPath:        app.config.CaptiveCoreBinaryPath,
		CaptiveCoreStoragePath:       app.config.CaptiveCoreStoragePath,
		CaptiveCoreToml:              app.config.CaptiveCoreToml,
		RemoteCaptiveCoreURL:         app.config.RemoteCaptiveCoreURL,
		EnableCaptiveCore:            app.config.EnableCaptiveCoreIngestion,
		DisableStateVerification:     app.config.IngestDisableStateVerification,
		EnableExtendedLogLedgerStats: app.config.IngestEnableExtendedLogLedgerStats,
	})

	if err != nil {
		log.Fatal(err)
	}
}

func initPathFinder(app *App) {
	orderBookGraph := orderbook.NewOrderBookGraph()
	app.orderBookStream = ingest.NewOrderBookStream(
		&history.Q{app.auroraSession()},
		orderBookGraph,
	)

	app.paths = simplepath.NewInMemoryFinder(orderBookGraph, !app.config.DisablePoolPathFinding)
}

// initSentry initialized the default sentry client with the configured DSN
func initSentry(app *App) {
	if app.config.SentryDSN == "" {
		return
	}

	log.WithField("dsn", app.config.SentryDSN).Info("Initializing sentry")
	err := raven.SetDSN(app.config.SentryDSN)
	if err != nil {
		log.Fatal(err)
	}
}

// initLogglyLog attaches a loggly hook to our logging system.
func initLogglyLog(app *App) {
	if app.config.LogglyToken == "" {
		return
	}

	log.WithFields(log.F{
		"token": app.config.LogglyToken,
		"tag":   app.config.LogglyTag,
	}).Info("Initializing loggly hook")

	hook := log.NewLogglyHook(app.config.LogglyToken, app.config.LogglyTag)
	log.DefaultLogger.AddHook(hook)

	go func() {
		<-app.ctx.Done()
		hook.Flush()
	}()
}

func initDbMetrics(app *App) {
	app.buildInfoGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{Namespace: "aurora", Subsystem: "build", Name: "info"},
		[]string{"version", "goversion"},
	)
	app.prometheusRegistry.MustRegister(app.buildInfoGauge)
	app.buildInfoGauge.With(prometheus.Labels{
		"version":   app.auroraVersion,
		"goversion": runtime.Version(),
	}).Inc()

	app.ingestingGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{Namespace: "aurora", Subsystem: "ingest", Name: "enabled"},
	)
	app.prometheusRegistry.MustRegister(app.ingestingGauge)

	app.ledgerState.RegisterMetrics(app.prometheusRegistry)

	app.coreState.RegisterMetrics(app.prometheusRegistry)

	app.prometheusRegistry.MustRegister(app.orderBookStream.LatestLedgerGauge)
}

// initGoMetrics registers the Go collector provided by prometheus package which
// includes Go-related metrics.
func initGoMetrics(app *App) {
	app.prometheusRegistry.MustRegister(prometheus.NewGoCollector())
}

// initProcessMetrics registers the process collector provided by prometheus
// package. This is only available on operating systems with a Linux-style proc
// filesystem and on Microsoft Windows.
func initProcessMetrics(app *App) {
	app.prometheusRegistry.MustRegister(
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
	)
}

// initIngestMetrics registers the metrics for the ingestion into the provided
// app's metrics registry.
func initIngestMetrics(app *App) {
	if app.ingester == nil {
		return
	}

	app.ingestingGauge.Inc()
	app.ingester.RegisterMetrics(app.prometheusRegistry)
}

func initTxSubMetrics(app *App) {
	app.submitter.Init()
	app.submitter.RegisterMetrics(app.prometheusRegistry)
}

func initWebMetrics(app *App) {
	app.webServer.RegisterMetrics(app.prometheusRegistry)
}

func initSubmissionSystem(app *App) {
	app.submitter = &txsub.System{
		Pending:         txsub.NewDefaultSubmissionList(),
		Submitter:       txsub.NewDefaultSubmitter(http.DefaultClient, app.config.diamcircleCoreURL),
		SubmissionQueue: sequence.NewManager(),
		DB: func(ctx context.Context) txsub.auroraDB {
			return &history.Q{SessionInterface: app.auroraSession()}
		},
	}
}
