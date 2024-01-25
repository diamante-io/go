package history

import (
	"encoding/json"
	"testing"

	"go/services/aurora/internal/test"
	"go/services/aurora/internal/toid"

	"github.com/guregu/null"
)

func TestAddEffect(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &Q{tt.auroraSession()}

	address := "GAQAA5L65LSYH7CQ3VTJ7F3HHLGCL3DSLAR2Y47263D56MNNGHSQSTVY"
	muxedAddres := "MAQAA5L65LSYH7CQ3VTJ7F3HHLGCL3DSLAR2Y47263D56MNNGHSQSAAAAAAAAAAE2LP26"
	accounIDs, err := q.CreateAccounts(tt.Ctx, []string{address}, 1)
	tt.Assert.NoError(err)

	builder := q.NewEffectBatchInsertBuilder(2)
	sequence := int32(56)
	details, err := json.Marshal(map[string]string{
		"amount":     "1000.0000000",
		"asset_type": "native",
	})

	err = builder.Add(tt.Ctx,
		accounIDs[address],
		null.StringFrom(muxedAddres),
		toid.New(sequence, 1, 1).ToInt64(),
		1,
		3,
		details,
	)
	tt.Assert.NoError(err)

	err = builder.Exec(tt.Ctx)
	tt.Assert.NoError(err)

	effects := []Effect{}
	tt.Assert.NoError(q.Effects().Select(tt.Ctx, &effects))
	tt.Assert.Len(effects, 1)

	effect := effects[0]
	tt.Assert.Equal(address, effect.Account)
	tt.Assert.Equal(muxedAddres, effect.AccountMuxed.String)
	tt.Assert.Equal(int64(240518172673), effect.HistoryOperationID)
	tt.Assert.Equal(int32(1), effect.Order)
	tt.Assert.Equal(EffectType(3), effect.Type)
	tt.Assert.Equal("{\"amount\": \"1000.0000000\", \"asset_type\": \"native\"}", effect.DetailsString.String)
}
