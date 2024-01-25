package actions

import (
	"context"
	"net/http"
	"testing"

	auroraContext "go/services/aurora/internal/context"
	"go/services/aurora/internal/db2/history"
	"go/services/aurora/internal/test"

	"github.com/stretchr/testify/assert"
)

func TestAssetsForAddressRequiresTransaction(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &history.Q{tt.auroraSession()}

	r := &http.Request{}
	ctx := context.WithValue(
		r.Context(),
		&auroraContext.SessionContextKey,
		q,
	)

	_, _, err := assetsForAddress(r.WithContext(ctx), "GCATOZ7YJV2FANQQLX47TIV6P7VMPJCEEJGQGR6X7TONPKBN3UCLKEIS")
	assert.EqualError(t, err, "cannot be called outside of a transaction")

	assert.NoError(t, q.Begin())
	defer q.Rollback()

	_, _, err = assetsForAddress(r.WithContext(ctx), "GCATOZ7YJV2FANQQLX47TIV6P7VMPJCEEJGQGR6X7TONPKBN3UCLKEIS")
	assert.EqualError(t, err, "should only be called in a repeatable read transaction")
}
