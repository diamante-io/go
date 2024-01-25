package ticker

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"

	"go/services/ticker/internal/gql"
	"go/services/ticker/internal/tickerdb/tickerdbtest"
	hlog "go/support/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTicker_btcEth(t *testing.T) {
	session := tickerdbtest.SetupTickerTestSession(t, "./tickerdb/migrations")
	defer session.DB.Close()

	logger := hlog.New()
	resolver := gql.New(&session, logger)
	h := resolver.NewRelayHandler()
	m := chi.NewMux()
	m.Post("/graphql", h.ServeHTTP)

	req := `{
	"query": "query getTicker() { ticker(pairName: \"BTC_ETH\", numHoursAgo: 24) { tradePair, baseVolume, counterVolume, tradeCount, open, close, high, low, orderbookStats { bidCount, bidVolume, bidMax, askCount, askVolume, askMin, spread, spreadMidPoint, } } }",
	"operationName": "getTicker",
	"variables": {}
}`
	r := httptest.NewRequest("POST", "/graphql", strings.NewReader(req))
	w := httptest.NewRecorder()

	m.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	wantBody := `{
	"data": {
		"ticker": [
			{
				"tradePair": "BTC_ETH",
				"baseVolume": 174,
				"counterVolume": 86,
				"tradeCount": 3,
				"open": 1,
				"close": 0.92,
				"high": 1.0,
				"low": 0.1,
				"orderbookStats": {
					"bidCount": 16,
					"bidVolume": 0.25,
					"bidMax": 200,
					"askCount": 18,
					"askVolume": 45,
					"askMin": 0.1,
					"spread": -1999,
					"spreadMidPoint": -799.5
				}
			}
		]
	}
}`
	assert.JSONEq(t, wantBody, string(body))
}

func TestMarkets(t *testing.T) {
	session := tickerdbtest.SetupTickerTestSession(t, "./tickerdb/migrations")
	defer session.DB.Close()

	logger := hlog.New()
	resolver := gql.New(&session, logger)
	h := resolver.NewRelayHandler()
	m := chi.NewMux()
	m.Post("/graphql", h.ServeHTTP)

	issuerPK := "GCF3TQXKZJNFJK7HCMNE2O2CUNKCJH2Y2ROISTBPLC7C5EIA5NNG2XZB"

	query := fmt.Sprintf(`
	query getMarkets() {
		markets(
			baseAssetCode: \"BTC\",
			baseAssetIssuer: \"%s\",
			counterAssetCode: \"ETH\",
			counterAssetIssuer: \"%s\",
			numHoursAgo: 24
		)
		{
			tradePair,
			baseAssetCode,
			baseAssetIssuer,
			counterAssetCode,
			counterAssetIssuer,
			baseVolume,
			counterVolume,
			tradeCount,
			open,
			low,
			high,
			close,
			change,
			orderbookStats {
				bidCount,
				bidVolume,
				bidMax,
				askCount,
				askVolume,
				askMin
			}
		}
	}`, issuerPK, issuerPK)
	req := fmt.Sprintf(`{
		"query": "%s",
		"operationName": "getMarkets",
		"variables": {}
	}`, formatMultiline(query))
	t.Log(req)

	wantBody := fmt.Sprintf(`
	{"data":{"markets": [{
		"tradePair": "BTC:%s / ETH:%s",
		"baseAssetCode": "BTC",
		"baseAssetIssuer":"%s",
		"counterAssetCode": "ETH",
		"counterAssetIssuer":"%s",
		"baseVolume": 150,
		"counterVolume": 60,
		"tradeCount": 2,
		"open":1,
		"low":0.1,
		"high":1,
		"close":0.1,
		"change":-0.9,
		"orderbookStats": {
			"bidCount": 15,
			"bidVolume": 0.15,
			"bidMax": 200,
			"askCount": 17,
			"askVolume": 30,
			"askMin": 0.1
		}
	}]}}`, issuerPK, issuerPK, issuerPK, issuerPK)
	testRequest(t, m, req, wantBody)
}

func testRequest(t *testing.T, m *chi.Mux, req, wantBody string) {
	r := httptest.NewRequest("POST", "/graphql", strings.NewReader(req))
	w := httptest.NewRecorder()

	m.ServeHTTP(w, r)
	resp := w.Result()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.JSONEq(t, wantBody, string(body))
}

func formatMultiline(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	return strings.ReplaceAll(s, "\t", "")
}
