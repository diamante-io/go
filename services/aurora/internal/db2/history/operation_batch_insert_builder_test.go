package history

import (
	"encoding/json"
	"testing"

	"github.com/guregu/null"
	"github.com/diamcircle/go/services/aurora/internal/test"
	"github.com/diamcircle/go/services/aurora/internal/toid"
	"github.com/diamcircle/go/xdr"
)

func TestAddOperation(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}

	txBatch := q.NewTransactionBatchInsertBuilder(0)

	builder := q.NewOperationBatchInsertBuilder(1)

	transactionHash := "2a805712c6d10f9e74bb0ccf54ae92a2b4b1e586451fe8133a2433816f6b567c"
	transactionResult := "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA="
	transaction := buildLedgerTransaction(
		t,
		testTransaction{
			index:         1,
			envelopeXDR:   "AAAAABpcjiETZ0uhwxJJhgBPYKWSVJy2TZ2LI87fqV1cUf/UAAAAZAAAADcAAAABAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAGlyOIRNnS6HDEkmGAE9gpZJUnLZNnYsjzt+pXVxR/9QAAAAAAAAAAAX14QAAAAAAAAAAAVxR/9QAAABAK6pcXYMzAEmH08CZ1LWmvtNDKauhx+OImtP/Lk4hVTMJRVBOebVs5WEPj9iSrgGT0EswuDCZ2i5AEzwgGof9Ag==",
			resultXDR:     transactionResult,
			metaXDR:       "AAAAAQAAAAIAAAADAAAAOAAAAAAAAAAAGlyOIRNnS6HDEkmGAE9gpZJUnLZNnYsjzt+pXVxR/9QAAAACVAvjnAAAADcAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAOAAAAAAAAAAAGlyOIRNnS6HDEkmGAE9gpZJUnLZNnYsjzt+pXVxR/9QAAAACVAvjnAAAADcAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAA==",
			feeChangesXDR: "AAAAAgAAAAMAAAA3AAAAAAAAAAAaXI4hE2dLocMSSYYAT2ClklSctk2diyPO36ldXFH/1AAAAAJUC+QAAAAANwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAA4AAAAAAAAAAAaXI4hE2dLocMSSYYAT2ClklSctk2diyPO36ldXFH/1AAAAAJUC+OcAAAANwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			hash:          transactionHash,
		},
	)

	sequence := int32(56)
	tt.Assert.NoError(txBatch.Add(tt.Ctx, transaction, uint32(sequence)))
	tt.Assert.NoError(txBatch.Exec(tt.Ctx))

	details, err := json.Marshal(map[string]string{
		"to":         "GANFZDRBCNTUXIODCJEYMACPMCSZEVE4WZGZ3CZDZ3P2SXK4KH75IK6Y",
		"from":       "GAQAA5L65LSYH7CQ3VTJ7F3HHLGCL3DSLAR2Y47263D56MNNGHSQSTVY",
		"amount":     "10.0000000",
		"asset_type": "native",
	})
	tt.Assert.NoError(err)

	sourceAccount := "GAQAA5L65LSYH7CQ3VTJ7F3HHLGCL3DSLAR2Y47263D56MNNGHSQSTVY"
	sourceAccountMuxed := "MAQAA5L65LSYH7CQ3VTJ7F3HHLGCL3DSLAR2Y47263D56MNNGHSQSAAAAAAAAAAE2LP26"
	err = builder.Add(tt.Ctx,
		toid.New(sequence, 1, 1).ToInt64(),
		toid.New(sequence, 1, 0).ToInt64(),
		1,
		xdr.OperationTypePayment,
		details,
		sourceAccount,
		null.StringFrom(sourceAccountMuxed),
	)
	tt.Assert.NoError(err)

	err = builder.Exec(tt.Ctx)
	tt.Assert.NoError(err)

	ops := []Operation{}
	err = q.Select(tt.Ctx, &ops, selectOperation)

	if tt.Assert.NoError(err) {
		tt.Assert.Len(ops, 1)

		op := ops[0]
		tt.Assert.Equal(int64(240518172673), op.ID)
		tt.Assert.Equal(int64(240518172672), op.TransactionID)
		tt.Assert.Equal(transactionHash, op.TransactionHash)
		tt.Assert.Equal(transactionResult, op.TxResult)
		tt.Assert.Equal(int32(1), op.ApplicationOrder)
		tt.Assert.Equal(xdr.OperationTypePayment, op.Type)
		tt.Assert.Equal(
			"{\"to\": \"GANFZDRBCNTUXIODCJEYMACPMCSZEVE4WZGZ3CZDZ3P2SXK4KH75IK6Y\", \"from\": \"GAQAA5L65LSYH7CQ3VTJ7F3HHLGCL3DSLAR2Y47263D56MNNGHSQSTVY\", \"amount\": \"10.0000000\", \"asset_type\": \"native\"}",
			op.DetailsString.String,
		)
		tt.Assert.Equal(sourceAccount, op.SourceAccount)
		tt.Assert.Equal(sourceAccountMuxed, op.SourceAccountMuxed.String)
		tt.Assert.Equal(true, op.TransactionSuccessful)
	}
}
