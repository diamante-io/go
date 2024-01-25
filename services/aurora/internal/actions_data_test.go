package aurora

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"

	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/services/aurora/internal/ingest"
	"github.com/diamcircle/go/services/aurora/internal/test"
	"github.com/diamcircle/go/xdr"
)

var (
	data1 = history.Data{
		LastModifiedLedger: 100,
		AccountID:          "GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB",
		Name:               "name1",
		// This also tests if base64 encoding is working as 0 is invalid UTF-8 byte
		Value: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},

		Sponsor: null.StringFrom("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
	}

	data2 = history.Data{
		LastModifiedLedger: 100,
		AccountID:          "GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB",
		Name:               "name ",
		Value:              []byte("it got spaces!"),
	}
)

func TestDataActions_Show(t *testing.T) {
	ht := StartHTTPTestWithoutScenario(t)
	defer ht.Finish()
	test.ResetAuroraDB(t, ht.AuroraDB)
	q := &history.Q{ht.AuroraSession()}

	// Makes StateMiddleware happy
	err := q.UpdateLastLedgerIngest(ht.Ctx, 100)
	ht.Assert.NoError(err)
	err = q.UpdateIngestVersion(ht.Ctx, ingest.CurrentVersion)
	ht.Assert.NoError(err)
	_, err = q.InsertLedger(ht.Ctx, xdr.LedgerHeaderHistoryEntry{
		Header: xdr.LedgerHeader{
			LedgerSeq: 100,
		},
	}, 0, 0, 0, 0, 0)
	ht.Assert.NoError(err)

	err = q.UpsertAccountData(ht.Ctx, []history.Data{data1, data2})
	assert.NoError(t, err)

	prefix := "/accounts/GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"
	result := map[string]string{}

	// json
	w := ht.Get(prefix + "/data/name1")
	if ht.Assert.Equal(200, w.Code) {
		err := json.Unmarshal(w.Body.Bytes(), &result)
		ht.Assert.NoError(err)
		decoded, err := base64.StdEncoding.DecodeString(result["value"])
		ht.Assert.NoError(err)
		ht.Assert.Equal([]byte(data1.Value), decoded)
		ht.Assert.Equal("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML", result["sponsor"])
	}

	// raw
	w = ht.Get(prefix+"/data/name1", test.RequestHelperRaw)
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.Equal([]byte(data1.Value), w.Body.Bytes())
	}

	result = map[string]string{}
	// regression: https://github.com/diamcircle/aurora/issues/325
	// names with special characters do not work
	w = ht.Get(prefix + "/data/name%20")
	if ht.Assert.Equal(200, w.Code) {
		err := json.Unmarshal(w.Body.Bytes(), &result)
		ht.Assert.NoError(err)

		decoded, err := base64.StdEncoding.DecodeString(result["value"])
		ht.Assert.NoError(err)
		ht.Assert.Equal([]byte(data2.Value), decoded)
		ht.Assert.Equal("", result["sponsor"])
	}

	w = ht.Get(prefix+"/data/name%20", test.RequestHelperRaw)
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.Equal("it got spaces!", w.Body.String())
	}

	// missing
	w = ht.Get(prefix + "/data/missing")
	ht.Assert.Equal(404, w.Code)

	w = ht.Get(prefix+"/data/missing", test.RequestHelperRaw)
	ht.Assert.Equal(404, w.Code)

	// Too long
	w = ht.Get(prefix+"/data/01234567890123456789012345678901234567890123456789012345678901234567890123456789", test.RequestHelperRaw)
	if ht.Assert.Equal(400, w.Code) {
		ht.Assert.Contains(w.Body.String(), "does not validate as length(1|64)")
	}

}
