package logmetrics

import (
	"bytes"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLogPackageMetrics(t *testing.T) {
	output := new(bytes.Buffer)
	l, m := New()
	l.DisableColors()
	l.SetLevel(logrus.DebugLevel)
	l.SetOutput(output)

	for _, meter := range *m {
		assert.Equal(t, float64(0), getMetricValue(meter).GetCounter().GetValue())
	}

	l.Debug("foo")
	l.Info("foo")
	l.Warn("foo")
	l.Error("foo")
	assert.Panics(t, func() {
		l.Panic("foo")
	})

	for _, meter := range *m {
		assert.Equal(t, float64(1), getMetricValue(meter).GetCounter().GetValue())
	}
}

func getMetricValue(metric prometheus.Metric) *dto.Metric {
	value := &dto.Metric{}
	err := metric.Write(value)
	if err != nil {
		panic(err)
	}
	return value
}
