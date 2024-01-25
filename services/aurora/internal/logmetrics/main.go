package logmetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/diamcircle/go/support/log"
)

// Metrics is a logrus hook-compliant struct that records metrics about logging
// when added to a logrus.Logger
type Metrics map[logrus.Level]prometheus.Counter

var DefaultMetrics = NewMetrics()

func init() {
	_, DefaultMetrics = New()
}

// New creates a new logger according to aurora specifications.
func New() (l *log.Entry, m *Metrics) {
	m = NewMetrics()
	l = log.New()
	l.SetLevel(logrus.WarnLevel)
	l.AddHook(m)
	return
}

// NewMetrics creates a new hook for recording metrics.
func NewMetrics() *Metrics {
	return &Metrics{
		logrus.DebugLevel: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "aurora", Subsystem: "log", Name: "debug_total",
		}),
		logrus.InfoLevel: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "aurora", Subsystem: "log", Name: "info_total",
		}),
		logrus.WarnLevel: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "aurora", Subsystem: "log", Name: "warn_total",
		}),
		logrus.ErrorLevel: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "aurora", Subsystem: "log", Name: "error_total",
		}),
		logrus.PanicLevel: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "aurora", Subsystem: "log", Name: "panic_total",
		}),
	}
}

// Fire is triggered by logrus, in response to a logging event
func (m *Metrics) Fire(e *logrus.Entry) error {
	(*m)[e.Level].Inc()
	return nil
}

// Levels returns the logging levels that will trigger this hook to run.  In
// this case, all of them.
func (m *Metrics) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.PanicLevel,
	}
}
