package aurora

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"go/protocols/diamcirclecore"
	"go/support/clock"
	"go/support/db"
	"go/support/log"
)

const (
	dbPingTimeout      = 5 * time.Second
	infoRequestTimeout = 5 * time.Second
	healthCacheTTL     = 500 * time.Millisecond
)

var healthLogger = log.WithField("service", "healthCheck")

type diamcircleCoreClient interface {
	Info(ctx context.Context) (*diamcirclecore.InfoResponse, error)
}

type healthCache struct {
	response   healthResponse
	lastUpdate time.Time
	ttl        time.Duration
	clock      clock.Clock
	lock       sync.Mutex
}

func (h *healthCache) get(runCheck func() healthResponse) healthResponse {
	h.lock.Lock()
	defer h.lock.Unlock()

	if h.clock.Now().Sub(h.lastUpdate) > h.ttl {
		h.response = runCheck()
		h.lastUpdate = h.clock.Now()
	}

	return h.response
}

func newHealthCache(ttl time.Duration) *healthCache {
	return &healthCache{ttl: ttl}
}

type healthCheck struct {
	session db.SessionInterface
	ctx     context.Context
	core    diamcircleCoreClient
	cache   *healthCache
}

type healthResponse struct {
	DatabaseConnected bool `json:"database_connected"`
	CoreUp            bool `json:"core_up"`
	CoreSynced        bool `json:"core_synced"`
}

func (h healthCheck) runCheck() healthResponse {
	response := healthResponse{
		DatabaseConnected: true,
		CoreUp:            true,
		CoreSynced:        true,
	}
	if err := h.session.Ping(h.ctx, dbPingTimeout); err != nil {
		healthLogger.Warnf("could not ping db: %s", err)
		response.DatabaseConnected = false
	}
	if resp, err := h.core.Info(h.ctx); err != nil {
		healthLogger.Warnf("request to diamcircle core failed: %s", err)
		response.CoreUp = false
		response.CoreSynced = false
	} else {
		response.CoreSynced = resp.IsSynced()
	}

	return response
}

func (h healthCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := h.cache.get(h.runCheck)

	if !response.DatabaseConnected || !response.CoreSynced || !response.CoreUp {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		healthLogger.Warnf("could not write response: %s", err)
	}
}
