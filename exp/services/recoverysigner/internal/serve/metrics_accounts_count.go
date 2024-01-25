package serve

import (
	"go/exp/services/recoverysigner/internal/account"
	supportlog "go/support/log"

	"github.com/prometheus/client_golang/prometheus"
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
