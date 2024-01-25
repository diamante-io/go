package resourceadapter

import (
	"context"
	"encoding/json"
	"testing"

	"go/protocols/aurora/effects"
	"go/services/aurora/internal/db2/history"
	"go/support/render/hal"
	"go/support/test"

	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
)

func TestNewEffectAllEffectsCovered(t *testing.T) {
	for typ, s := range EffectTypeNames {
		if typ == history.EffectAccountRemoved || typ == history.EffectAccountInflationDestinationUpdated {
			// these effects use the base representation
			continue
		}
		e := history.Effect{
			Type: typ,
		}
		result, err := NewEffect(context.TODO(), e, history.Ledger{})
		assert.NoError(t, err, s)
		// it shouldn't be a base type
		_, ok := result.(effects.Base)
		assert.False(t, ok, s)
	}

	// verify that the check works for an unknown effect
	e := history.Effect{
		Type: 20000,
	}
	result, err := NewEffect(context.TODO(), e, history.Ledger{})
	assert.NoError(t, err)
	_, ok := result.(effects.Base)
	assert.True(t, ok)
}

func TestEffectTypeNamesAreConsistentWithAdapterTypeNames(t *testing.T) {
	for typ, s := range EffectTypeNames {
		s2, ok := effects.EffectTypeNames[effects.EffectType(typ)]
		assert.True(t, ok, s)
		assert.Equal(t, s, s2)
	}
	for typ, s := range effects.EffectTypeNames {
		s2, ok := EffectTypeNames[history.EffectType(typ)]
		assert.True(t, ok, s)
		assert.Equal(t, s, s2)
	}
}

func TestNewEffect_EffectTrustlineAuthorizedToMaintainLiabilities(t *testing.T) {
	tt := assert.New(t)
	ctx, _ := test.ContextWithLogBuffer()

	details := `{
		"asset_code":   "COP",
		"asset_issuer": "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
		"asset_type":   "credit_alphanum4",
		"trustor":      "GDQNY3PBOJOKYZSRMK2S7LHHGWZIUISD4QORETLMXEWXBI7KFZZMKTL3"
	}`

	hEffect := history.Effect{
		Account:            "GDQNY3PBOJOKYZSRMK2S7LHHGWZIUISD4QORETLMXEWXBI7KFZZMKTL3",
		HistoryOperationID: 1,
		Order:              1,
		Type:               history.EffectTrustlineAuthorizedToMaintainLiabilities,
		DetailsString:      null.StringFrom(details),
	}
	resource, err := NewEffect(ctx, hEffect, history.Ledger{})
	tt.NoError(err)

	var resourcePage hal.Page
	resourcePage.Add(resource)

	effect, ok := resource.(effects.TrustlineAuthorizedToMaintainLiabilities)
	tt.True(ok)
	tt.Equal("trustline_authorized_to_maintain_liabilities", effect.Type)

	binary, err := json.Marshal(resourcePage)
	tt.NoError(err)

	var page effects.EffectsPage
	tt.NoError(json.Unmarshal(binary, &page))
	tt.Len(page.Embedded.Records, 1)
	tt.Equal(effect, page.Embedded.Records[0].(effects.TrustlineAuthorizedToMaintainLiabilities))
}

func TestNewEffect_EffectTrade_Muxed(t *testing.T) {
	tt := assert.New(t)
	ctx, _ := test.ContextWithLogBuffer()

	details := `{
	"seller": "GAQAA5L65LSYH7CQ3VTJ7F3HHLGCL3DSLAR2Y47263D56MNNGHSQSTVY",
	"seller_muxed": "MAQAA5L65LSYH7CQ3VTJ7F3HHLGCL3DSLAR2Y47263D56MNNGHSQSAAAAAAAAAAE2LP26",
    "seller_muxed_id": 1234
	}`

	hEffect := history.Effect{
		Account:            "GAQAA5L65LSYH7CQ3VTJ7F3HHLGCL3DSLAR2Y47263D56MNNGHSQSTVY",
		AccountMuxed:       null.StringFrom("MAQAA5L65LSYH7CQ3VTJ7F3HHLGCL3DSLAR2Y47263D56MNNGHSQSAAAAAAAAAAE2LP26"),
		HistoryOperationID: 1,
		Order:              1,
		Type:               history.EffectTrade,
		DetailsString:      null.StringFrom(details),
	}
	resource, err := NewEffect(ctx, hEffect, history.Ledger{})
	tt.NoError(err)

	var resourcePage hal.Page
	resourcePage.Add(resource)

	effect, ok := resource.(effects.Trade)
	tt.True(ok)
	tt.Equal("trade", effect.Type)

	binary, err := json.Marshal(resourcePage)
	tt.NoError(err)

	var page effects.EffectsPage
	tt.NoError(json.Unmarshal(binary, &page))
	tt.Len(page.Embedded.Records, 1)
	tt.Equal(effect, page.Embedded.Records[0].(effects.Trade))
	tt.Equal("GAQAA5L65LSYH7CQ3VTJ7F3HHLGCL3DSLAR2Y47263D56MNNGHSQSTVY", effect.Account)
	tt.Equal("MAQAA5L65LSYH7CQ3VTJ7F3HHLGCL3DSLAR2Y47263D56MNNGHSQSAAAAAAAAAAE2LP26", effect.AccountMuxed)
	tt.Equal(uint64(1234), effect.AccountMuxedID)
	tt.Equal("GAQAA5L65LSYH7CQ3VTJ7F3HHLGCL3DSLAR2Y47263D56MNNGHSQSTVY", effect.Seller)
	tt.Equal("MAQAA5L65LSYH7CQ3VTJ7F3HHLGCL3DSLAR2Y47263D56MNNGHSQSAAAAAAAAAAE2LP26", effect.SellerMuxed)
	tt.Equal(uint64(1234), effect.SellerMuxedID)
}
