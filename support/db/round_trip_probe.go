package db

import (
	"context"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type roundTripProbe struct {
	session              SessionInterface
	roundTripTimeSummary prometheus.Summary

	closeChan chan struct{}
	closeOnce sync.Once
}

func (p *roundTripProbe) start() {
	p.closeChan = make(chan struct{})
	// session must be cloned because will be used concurrently in a
	// separate go routine in roundTripProbe
	p.session = p.session.Clone()
	ticker := time.NewTicker(time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				startTime := time.Now()
				_, err := p.session.ExecRaw(ctx, "select 1")
				duration := time.Since(startTime).Seconds()
				if err != nil {
					duration = 1
				}
				p.roundTripTimeSummary.Observe(duration)
				cancel()
			case <-p.closeChan:
				ticker.Stop()
				return
			}
		}
	}()
}

func (p *roundTripProbe) close() {
	p.closeOnce.Do(func() {
		close(p.closeChan)
	})
}
