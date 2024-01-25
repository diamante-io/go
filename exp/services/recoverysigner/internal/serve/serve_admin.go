package serve

import (
	"fmt"
	"net/http"

	supporthttp "go/support/http"
	supportlog "go/support/log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func serveAdmin(opts Options, deps adminDeps) {
	adminHandler := adminHandler(deps)

	addr := fmt.Sprintf(":%d", opts.AdminPort)
	supporthttp.Run(supporthttp.Config{
		ListenAddr: addr,
		Handler:    adminHandler,
		OnStarting: func() {
			deps.Logger.Infof("Starting admin port server on %s", addr)
		},
	})
}

type adminDeps struct {
	Logger          *supportlog.Entry
	MetricsGatherer prometheus.Gatherer
}

func adminHandler(deps adminDeps) http.Handler {
	mux := supporthttp.NewMux(deps.Logger)
	mux.Handle("/metrics", promhttp.HandlerFor(deps.MetricsGatherer, promhttp.HandlerOpts{}))
	return mux
}
