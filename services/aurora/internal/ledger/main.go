// Package ledger provides useful utilities concerning ledgers within diamcircle,
// specifically as a central location to store a cached snapshot of the state of
// both aurora's and diamcircle-core's views of the ledger.  This package is
// intended to be at the lowest levels of aurora's dependency tree, please keep
// it free of dependencies to other aurora packages.
package ledger

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Status represents a snapshot of both aurora's and diamcircle-core's view of the
// ledger.
type Status struct {
	CoreStatus
	AuroraStatus
}

type CoreStatus struct {
	CoreLatest int32 `db:"core_latest"`
}

type AuroraStatus struct {
	HistoryLatest         int32     `db:"history_latest"`
	HistoryLatestClosedAt time.Time `db:"history_latest_closed_at"`
	HistoryElder          int32     `db:"history_elder"`
	ExpHistoryLatest      uint32    `db:"exp_history_latest"`
}

// State is an in-memory data structure which holds a snapshot of both
// aurora's and diamcircle-core's view of the the network
type State struct {
	sync.RWMutex
	current Status

	Metrics struct {
		HistoryLatestLedgerCounter        prometheus.CounterFunc
		HistoryLatestLedgerClosedAgoGauge prometheus.GaugeFunc
		HistoryElderLedgerCounter         prometheus.CounterFunc
		CoreLatestLedgerCounter           prometheus.CounterFunc
	}
}

// CurrentStatus returns the cached snapshot of ledger state
func (c *State) CurrentStatus() Status {
	c.RLock()
	defer c.RUnlock()
	ret := c.current
	return ret
}

// SetStatus updates the cached snapshot of the ledger state
func (c *State) SetStatus(next Status) {
	c.Lock()
	defer c.Unlock()
	c.current = next
}

// SetCoreStatus updates the cached snapshot of the ledger state of Diamcircle-Core
func (c *State) SetCoreStatus(next CoreStatus) {
	c.Lock()
	defer c.Unlock()
	c.current.CoreStatus = next
}

// SetAuroraStatus updates the cached snapshot of the ledger state of Aurora
func (c *State) SetAuroraStatus(next AuroraStatus) {
	c.Lock()
	defer c.Unlock()
	c.current.AuroraStatus = next
}

func (c *State) RegisterMetrics(registry *prometheus.Registry) {
	c.Metrics.HistoryLatestLedgerCounter = prometheus.NewCounterFunc(
		prometheus.CounterOpts{Namespace: "aurora", Subsystem: "history", Name: "latest_ledger"},
		func() float64 {
			ls := c.CurrentStatus()
			return float64(ls.HistoryLatest)
		},
	)
	registry.MustRegister(c.Metrics.HistoryLatestLedgerCounter)

	c.Metrics.HistoryLatestLedgerClosedAgoGauge = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: "aurora", Subsystem: "history", Name: "latest_ledger_closed_ago_seconds",
			Help: "seconds since the close of the last ingested ledger",
		},
		func() float64 {
			ls := c.CurrentStatus()
			return time.Since(ls.HistoryLatestClosedAt).Seconds()
		},
	)
	registry.MustRegister(c.Metrics.HistoryLatestLedgerClosedAgoGauge)

	c.Metrics.HistoryElderLedgerCounter = prometheus.NewCounterFunc(
		prometheus.CounterOpts{Namespace: "aurora", Subsystem: "history", Name: "elder_ledger"},
		func() float64 {
			ls := c.CurrentStatus()
			return float64(ls.HistoryElder)
		},
	)
	registry.MustRegister(c.Metrics.HistoryElderLedgerCounter)

	c.Metrics.CoreLatestLedgerCounter = prometheus.NewCounterFunc(
		prometheus.CounterOpts{Namespace: "aurora", Subsystem: "diamcircle_core", Name: "latest_ledger"},
		func() float64 {
			ls := c.CurrentStatus()
			return float64(ls.CoreLatest)
		},
	)
	registry.MustRegister(c.Metrics.CoreLatestLedgerCounter)
}
