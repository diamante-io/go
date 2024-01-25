package history

import (
	"testing"
	"time"

	"go/services/aurora/internal/test"
)

func TestLatestLedger(t *testing.T) {
	tt := test.Start(t)
	tt.Scenario("base")
	defer tt.Finish()
	q := &Q{tt.auroraSession()}

	var seq int
	err := q.LatestLedger(tt.Ctx, &seq)

	if tt.Assert.NoError(err) {
		tt.Assert.Equal(3, seq)
	}
}

func TestLatestLedgerSequenceClosedAt(t *testing.T) {
	tt := test.Start(t)
	tt.Scenario("base")
	defer tt.Finish()
	q := &Q{tt.auroraSession()}

	sequence, closedAt, err := q.LatestLedgerSequenceClosedAt(tt.Ctx)
	if tt.Assert.NoError(err) {
		tt.Assert.Equal(int32(3), sequence)
		tt.Assert.Equal("2019-10-31T13:19:46Z", closedAt.Format(time.RFC3339))
	}

	test.ResetauroraDB(t, tt.auroraDB)

	sequence, closedAt, err = q.LatestLedgerSequenceClosedAt(tt.Ctx)
	if tt.Assert.NoError(err) {
		tt.Assert.Equal(int32(0), sequence)
		tt.Assert.Equal("0001-01-01T00:00:00Z", closedAt.Format(time.RFC3339))
	}
}

func TestGetLatestHistoryLedgerEmptyDB(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &Q{tt.auroraSession()}

	value, err := q.GetLatestHistoryLedger(tt.Ctx)
	tt.Assert.NoError(err)
	tt.Assert.Equal(uint32(0), value)
}

func TestElderLedger(t *testing.T) {
	tt := test.Start(t)
	tt.Scenario("base")
	defer tt.Finish()
	q := &Q{tt.auroraSession()}

	var seq int
	err := q.ElderLedger(tt.Ctx, &seq)

	if tt.Assert.NoError(err) {
		tt.Assert.Equal(1, seq)
	}
}
