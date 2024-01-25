package httpx

import (
	"net/http"
	"strings"
	"time"

	"github.com/diamcircle/throttled"

	"go/services/aurora/internal/ledger"
	hProblem "go/services/aurora/internal/render/problem"
	"go/support/render/problem"
)

const lruCacheSize = 50000

type historyLedgerSourceFactory struct {
	updateFrequency time.Duration
	ledgerState     *ledger.State
}

func (f historyLedgerSourceFactory) Get() ledger.Source {
	return ledger.NewHistoryDBSource(f.updateFrequency, f.ledgerState)
}

func remoteAddrIP(r *http.Request) string {
	// To support IPv6
	lastSemicolon := strings.LastIndex(r.RemoteAddr, ":")
	if lastSemicolon == -1 {
		return r.RemoteAddr
	} else {
		return r.RemoteAddr[0:lastSemicolon]
	}
}

type VaryByRemoteIP struct{}

func (v VaryByRemoteIP) Key(r *http.Request) string {
	return remoteAddrIP(r)
}

func newRateLimiter(rateQuota *throttled.RateQuota) (*throttled.HTTPRateLimiter, error) {
	rateLimiter, err := throttled.NewGCRARateLimiter(lruCacheSize, *rateQuota)
	if err != nil {
		return nil, err
	}

	result := &throttled.HTTPRateLimiter{
		RateLimiter: rateLimiter,
		DeniedHandler: http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			problem.Render(request.Context(), w, hProblem.RateLimitExceeded)
		}),
		VaryBy: VaryByRemoteIP{},
	}
	return result, nil
}
