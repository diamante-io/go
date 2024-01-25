package aurora

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/diamcircle/go/services/aurora/internal/test"
)

type HTTPT struct {
	Assert     *test.Assertions
	App        *App
	RH         test.RequestHelper
	coreServer *test.StaticMockServer
	*test.T
}

func startHTTPTest(t *testing.T, scenario string) *HTTPT {
	ret := &HTTPT{T: test.Start(t)}
	if scenario == "" {
		test.ResetAuroraDB(t, ret.AuroraDB)
	} else {
		ret.Scenario(scenario)
	}
	ret.App = NewTestApp()
	ret.RH = test.NewRequestHelper(ret.App.webServer.Router.Mux)
	ret.Assert = &test.Assertions{ret.T.Assert}

	ret.coreServer = test.NewStaticMockServer(`{
		"info": {
			"network": "test",
			"build": "test-core",
			"ledger": {
				"version": 18,
				"num": 64
			},
			"protocol_version": 18,
			"network": "Test SDF Network ; September 2015"
		}
	}`)

	ret.App.config.DiamcircleCoreURL = ret.coreServer.URL
	ret.App.UpdateCoreLedgerState(context.Background())
	ret.App.UpdateDiamcircleCoreInfo(context.Background())
	ret.App.UpdateAuroraLedgerState(context.Background())

	return ret
}

// StartHTTPTest is a helper function to setup a new test that will make http
// requests. Pair it with a deferred call to FinishHTTPTest.
func StartHTTPTest(t *testing.T, scenario string) *HTTPT {
	if scenario == "" {
		t.Fatal("scenario cannot be empty string")
	}
	return startHTTPTest(t, scenario)
}

// StartHTTPTestWithoutScenario is like StartHTTPTest except it does not use
// a sql scenario
func StartHTTPTestWithoutScenario(t *testing.T) *HTTPT {
	return startHTTPTest(t, "")
}

// Get delegates to the test's request helper
func (ht *HTTPT) Get(
	path string,
	fn ...func(*http.Request),
) *httptest.ResponseRecorder {
	return ht.RH.Get(path, fn...)
}

// GetWithParams delegates to the test's request helper and encodes along with the query params
func (ht *HTTPT) GetWithParams(
	path string,
	queryParams url.Values,
	fn ...func(*http.Request),
) *httptest.ResponseRecorder {
	return ht.RH.Get(path+"?"+queryParams.Encode(), fn...)
}

// Finish closes the test app and finishes the test
func (ht *HTTPT) Finish() {
	ht.T.Finish()
	ht.App.Close()
	ht.coreServer.Close()
}

// Post delegates to the test's request helper
func (ht *HTTPT) Post(
	path string,
	form url.Values,
	mods ...func(*http.Request),
) *httptest.ResponseRecorder {
	return ht.RH.Post(path, form, mods...)
}

// ReapHistory causes the test server to run `DeleteUnretainedHistory`, after
// setting the retention count to the provided number.
func (ht *HTTPT) ReapHistory(retention uint) {
	ht.App.reaper.RetentionCount = retention
	err := ht.App.DeleteUnretainedHistory(context.Background())
	ht.Require.NoError(err)
	ht.App.UpdateCoreLedgerState(context.Background())
	ht.App.UpdateAuroraLedgerState(context.Background())
}
