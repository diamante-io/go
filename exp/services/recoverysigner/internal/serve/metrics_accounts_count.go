package serve

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/diamcircle/go/exp/services/recoverysigner/internal/account"
	supportlog "github.com/diamcircle/go/support/log"
)

type metricAccountsCount struct {
	Logger       *supportlog.Entry
	AccountStore account.Store
}

func (m metricAccountsCount) NewCollector() prometheus.Collector {
	opts := prometheus.GaugeOpts{
		Name: "accounts_count",
		Help: "Number of active accounts.",
	}
	return prometheus.NewGaugeFunc(opts, m.gauge)
}

func (m metricAccountsCount) gauge() float64 {
	count, err := m.AccountStore.Count()
	if err != nil {
		m.Logger.Warnf("Error getting count from account store: %v", err)
		return 0
	}
	return float64(count)
}
