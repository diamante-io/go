package history

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/guregu/null"
	"github.com/diamcircle/go/protocols/aurora/effects"
	"github.com/diamcircle/go/services/aurora/internal/db2"
	"github.com/diamcircle/go/services/aurora/internal/test"
	"github.com/diamcircle/go/services/aurora/internal/toid"
)

func TestEffectsForLiquidityPool(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}

	// Insert Effect
	address := "GAQAA5L65LSYH7CQ3VTJ7F3HHLGCL3DSLAR2Y47263D56MNNGHSQSTVY"
	muxedAddres := "MAQAA5L65LSYH7CQ3VTJ7F3HHLGCL3DSLAR2Y47263D56MNNGHSQSAAAAAAAAAAE2LP26"
	accountIDs, err := q.CreateAccounts(tt.Ctx, []string{address}, 1)
	tt.Assert.NoError(err)

	builder := q.NewEffectBatchInsertBuilder(2)
	sequence := int32(56)
	details, err := json.Marshal(map[string]string{
		"amount":     "1000.0000000",
		"asset_type": "native",
	})
	opID := toid.New(sequence, 1, 1).ToInt64()
	err = builder.Add(tt.Ctx,
		accountIDs[address],
		null.StringFrom(muxedAddres),
		opID,
		1,
		3,
		details,
	)
	tt.Assert.NoError(err)

	err = builder.Exec(tt.Ctx)
	tt.Assert.NoError(err)

	// Insert Liquidity Pool history
	liquidityPoolID := "abcde"
	toInternalID, err := q.CreateHistoryLiquidityPools(tt.Ctx, []string{liquidityPoolID}, 2)
	tt.Assert.NoError(err)
	operationBuilder := q.NewOperationLiquidityPoolBatchInsertBuilder(2)
	tt.Assert.NoError(err)
	internalID, ok := toInternalID[liquidityPoolID]
	tt.Assert.True(ok)
	err = operationBuilder.Add(tt.Ctx, opID, internalID)
	tt.Assert.NoError(err)
	err = operationBuilder.Exec(tt.Ctx)
	tt.Assert.NoError(err)

	var result []Effect
	err = q.Effects().ForLiquidityPool(tt.Ctx, db2.PageQuery{
		Cursor: "0-0",
		Order:  "asc",
		Limit:  10,
	}, liquidityPoolID).Select(tt.Ctx, &result)
	tt.Assert.NoError(err)

	tt.Assert.Len(result, 1)
	tt.Assert.Equal(result[0].Account, address)

}

func TestEffectsForTrustlinesSponsorshipEmptyAssetType(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}

	address := "GAQAA5L65LSYH7CQ3VTJ7F3HHLGCL3DSLAR2Y47263D56MNNGHSQSTVY"
	muxedAddres := "MAQAA5L65LSYH7CQ3VTJ7F3HHLGCL3DSLAR2Y47263D56MNNGHSQSAAAAAAAAAAE2LP26"
	accountIDs, err := q.CreateAccounts(tt.Ctx, []string{address}, 1)
	tt.Assert.NoError(err)

	builder := q.NewEffectBatchInsertBuilder(1)
	sequence := int32(56)
	tests := []struct {
		effectType        EffectType
		details           map[string]string
		expectedAssetType string
	}{
		{
			EffectTrustlineSponsorshipCreated,
			map[string]string{
				"asset":   "USD:GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
				"sponsor": "GDMQUXK7ZUCWM5472ZU3YLDP4BMJLQQ76DEMNYDEY2ODEEGGRKLEWGW2",
			},
			"credit_alphanum4",
		},
		{
			EffectTrustlineSponsorshipCreated,
			map[string]string{
				"asset":   "USDCE:GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
				"sponsor": "GDMQUXK7ZUCWM5472ZU3YLDP4BMJLQQ76DEMNYDEY2ODEEGGRKLEWGW2",
			},
			"credit_alphanum12",
		},
		{
			EffectTrustlineSponsorshipUpdated,
			map[string]string{
				"asset":   "USD:GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
				"sponsor": "GDMQUXK7ZUCWM5472ZU3YLDP4BMJLQQ76DEMNYDEY2ODEEGGRKLEWGW2",
			},
			"credit_alphanum4",
		},
		{
			EffectTrustlineSponsorshipUpdated,
			map[string]string{
				"asset":   "USDCE:GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
				"sponsor": "GDMQUXK7ZUCWM5472ZU3YLDP4BMJLQQ76DEMNYDEY2ODEEGGRKLEWGW2",
			},
			"credit_alphanum12",
		},
		{
			EffectTrustlineSponsorshipRemoved,
			map[string]string{
				"asset":   "USD:GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
				"sponsor": "GDMQUXK7ZUCWM5472ZU3YLDP4BMJLQQ76DEMNYDEY2ODEEGGRKLEWGW2",
			},
			"credit_alphanum4",
		},
		{
			EffectTrustlineSponsorshipRemoved,
			map[string]string{
				"asset":   "USDCE:GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
				"sponsor": "GDMQUXK7ZUCWM5472ZU3YLDP4BMJLQQ76DEMNYDEY2ODEEGGRKLEWGW2",
			},
			"credit_alphanum12",
		},
	}
	opID := toid.New(sequence, 1, 1).ToInt64()

	for i, test := range tests {
		var bytes []byte
		bytes, err = json.Marshal(test.details)
		tt.Require.NoError(err)

		err = builder.Add(tt.Ctx,
			accountIDs[address],
			null.StringFrom(muxedAddres),
			opID,
			uint32(i),
			test.effectType,
			bytes,
		)
		tt.Require.NoError(err)
	}

	err = builder.Exec(tt.Ctx)
	tt.Require.NoError(err)

	var results []Effect
	err = q.Effects().Select(tt.Ctx, &results)
	tt.Require.NoError(err)
	tt.Require.Len(results, len(tests))

	for i, test := range tests {
		switch test.effectType {
		case EffectTrustlineSponsorshipCreated:
			var eff effects.TrustlineSponsorshipCreated
			err := results[i].UnmarshalDetails(&eff)
			tt.Require.NoError(err)
			tt.Assert.Equal(test.expectedAssetType, eff.Type)
		case EffectTrustlineSponsorshipUpdated:
			var eff effects.TrustlineSponsorshipUpdated
			err := results[i].UnmarshalDetails(&eff)
			tt.Require.NoError(err)
			tt.Assert.Equal(test.expectedAssetType, eff.Type)
		case EffectTrustlineSponsorshipRemoved:
			var eff effects.TrustlineSponsorshipRemoved
			err := results[i].UnmarshalDetails(&eff)
			tt.Require.NoError(err)
			tt.Assert.Equal(test.expectedAssetType, eff.Type)
		default:
			panic(fmt.Sprintf("Unknown type %v", test.effectType))
		}
	}
}
