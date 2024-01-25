package aurora

import (
	"encoding/json"
	"strings"
	"testing"

	"go/protocols/aurora"
	"go/services/aurora/internal/test"

	"github.com/stretchr/testify/assert"
)

func TestRootAction(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	server := test.NewStaticMockServer(`{
			"info": {
				"network": "test",
				"build": "test-core",
				"ledger": {
					"version": 3,
					"num": 64
				},
				"protocol_version": 4
			}
		}`)
	defer server.Close()

	ht.App.config.diamcircleCoreURL = server.URL
	ht.App.config.NetworkPassphrase = "test"
	assert.NoError(t, ht.App.UpdatediamcircleCoreInfo(ht.Ctx))
	ht.App.UpdateCoreLedgerState(ht.Ctx)
	ht.App.UpdateauroraLedgerState(ht.Ctx)

	w := ht.Get("/")

	if ht.Assert.Equal(200, w.Code) {
		var actual aurora.Root
		err := json.Unmarshal(w.Body.Bytes(), &actual)
		ht.Require.NoError(err)
		ht.Assert.Equal("devel", actual.auroraVersion)
		ht.Assert.Equal("test-core", actual.diamcircleCoreVersion)
		ht.Assert.Equal(int32(4), actual.CoreSupportedProtocolVersion)
		ht.Assert.Equal(int32(3), actual.CurrentProtocolVersion)
		ht.Assert.Equal(int32(64), actual.CoreSequence)

		err = json.Unmarshal(w.Body.Bytes(), &actual)
		ht.Require.NoError(err)
		ht.Assert.Equal(
			"http://localhost/accounts{?signer,sponsor,asset,liquidity_pool,cursor,limit,order}",
			actual.Links.Accounts.Href,
		)
		ht.Assert.Equal(
			"http://localhost/offers{?selling,buying,seller,sponsor,cursor,limit,order}",
			actual.Links.Offers.Href,
		)

		params := []string{
			"destination_account",
			"destination_assets",
			"source_asset_type",
			"source_asset_issuer",
			"source_asset_code",
			"source_amount",
		}

		ht.Assert.Equal(
			"http://localhost/paths/strict-send{?"+strings.Join(params, ",")+"}",
			actual.Links.StrictSendPaths.Href,
		)

		params = []string{
			"source_assets",
			"source_account",
			"destination_account",
			"destination_asset_type",
			"destination_asset_issuer",
			"destination_asset_code",
			"destination_amount",
		}

		ht.Assert.Equal(
			"http://localhost/paths/strict-receive{?"+strings.Join(params, ",")+"}",
			actual.Links.StrictReceivePaths.Href,
		)
	}
}

func TestRootCoreClientInfoErrored(t *testing.T) {
	ht := StartHTTPTestWithoutScenario(t)
	defer ht.Finish()

	// an empty payload causes the core client to err
	server := test.NewStaticMockServer(`{}`)
	defer server.Close()

	ht.App.config.diamcircleCoreURL = server.URL
	ht.App.UpdateCoreLedgerState(ht.Ctx)

	w := ht.Get("/")

	if ht.Assert.Equal(200, w.Code) {
		var actual aurora.Root
		err := json.Unmarshal(w.Body.Bytes(), &actual)
		ht.Require.NoError(err)
		ht.Assert.Equal(int32(18), actual.CurrentProtocolVersion)
	}
}
