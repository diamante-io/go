package aurora

import (
	"context"
	"net/http"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/diamcircle/go/services/aurora/internal/test"
)

func TestGenericHTTPFeatures(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	// CORS
	w := ht.Get("/")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.Empty(w.HeaderMap.Get("Access-Control-Allow-Origin"))
	}

	w = ht.Get("/", func(r *http.Request) {
		r.Header.Set("Origin", "somewhere.com")
	})

	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.Equal(
			"somewhere.com",
			w.HeaderMap.Get("Access-Control-Allow-Origin"),
		)
	}

	// Trailing slash is stripped
	w = ht.Get("/ledgers")
	ht.Assert.Equal(200, w.Code)
	w = ht.Get("/ledgers/")
	ht.Assert.Equal(200, w.Code)
}

func TestMetrics(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	adminRouterRH := test.NewRequestHelper(ht.App.webServer.Router.Internal)
	w := adminRouterRH.Get("/metrics")
	ht.Assert.Equal(200, w.Code)

	hl := ht.App.ledgerState.Metrics.HistoryLatestLedgerCounter
	hlc := ht.App.ledgerState.Metrics.HistoryLatestLedgerClosedAgoGauge
	he := ht.App.ledgerState.Metrics.HistoryElderLedgerCounter
	cl := ht.App.ledgerState.Metrics.CoreLatestLedgerCounter

	ht.Require.EqualValues(3, getMetricValue(hl).GetCounter().GetValue())
	ht.Require.Less(float64(1000), getMetricValue(hlc).GetGauge().GetValue())
	ht.Require.EqualValues(1, getMetricValue(he).GetCounter().GetValue())
	ht.Require.EqualValues(64, getMetricValue(cl).GetCounter().GetValue())
}

func TestTick(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	// Just sanity-check that we return the context error...
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := ht.App.Tick(ctx)
	ht.Assert.EqualError(err, context.Canceled.Error())
}

func getMetricValue(metric prometheus.Metric) *dto.Metric {
	value := &dto.Metric{}
	err := metric.Write(value)
	if err != nil {
		panic(err)
	}
	return value
}
