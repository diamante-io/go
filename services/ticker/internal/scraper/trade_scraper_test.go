package scraper

import (
	"fmt"
	"testing"

	hProtocol "go/protocols/aurora"

	"github.com/stretchr/testify/assert"
)

func TestReverseAssets(t *testing.T) {
	baseAmount := "10.0"
	baseAccount := "BASEACCOUNT"
	baseAssetCode := "BASECODE"
	baseAssetType := "BASEASSETTYPE"
	baseAssetIssuer := "BASEASSETISSUER"

	counterAmount := "5.0"
	counterAccount := "COUNTERACCOUNT"
	counterAssetCode := "COUNTERASSETCODE"
	counterAssetType := "COUNTERASSETTYPE"
	counterAssetIssuer := "COUNTERASSETISSUER"

	baseIsSeller := true

	n := int64(10)
	d := int64(50)

	price := hProtocol.TradePrice{
		N: n,
		D: d,
	}

	trade1 := hProtocol.Trade{
		BaseAmount:         baseAmount,
		BaseAccount:        baseAccount,
		BaseAssetCode:      baseAssetCode,
		BaseAssetType:      baseAssetType,
		BaseAssetIssuer:    baseAssetIssuer,
		CounterAmount:      counterAmount,
		CounterAccount:     counterAccount,
		CounterAssetCode:   counterAssetCode,
		CounterAssetType:   counterAssetType,
		CounterAssetIssuer: counterAssetIssuer,
		BaseIsSeller:       baseIsSeller,
		Price:              price,
	}

	fmt.Println(trade1)

	reverseAssets(&trade1)

	assert.Equal(t, counterAmount, trade1.BaseAmount)
	assert.Equal(t, counterAccount, trade1.BaseAccount)
	assert.Equal(t, counterAssetCode, trade1.BaseAssetCode)
	assert.Equal(t, counterAssetType, trade1.BaseAssetType)
	assert.Equal(t, counterAssetIssuer, trade1.BaseAssetIssuer)

	assert.Equal(t, baseAmount, trade1.CounterAmount)
	assert.Equal(t, baseAccount, trade1.CounterAccount)
	assert.Equal(t, baseAssetCode, trade1.CounterAssetCode)
	assert.Equal(t, baseAssetType, trade1.CounterAssetType)
	assert.Equal(t, baseAssetIssuer, trade1.CounterAssetIssuer)

	assert.Equal(t, !baseIsSeller, trade1.BaseIsSeller)

	assert.Equal(t, d, trade1.Price.N)
	assert.Equal(t, n, trade1.Price.D)
}

func TestAddNativeData(t *testing.T) {
	trade1 := hProtocol.Trade{
		BaseAssetType: "native",
	}

	addNativeData(&trade1)
	assert.Equal(t, "XLM", trade1.BaseAssetCode)
	assert.Equal(t, "native", trade1.BaseAssetIssuer)

	trade2 := hProtocol.Trade{
		CounterAssetType: "native",
	}
	addNativeData(&trade2)
	assert.Equal(t, "XLM", trade2.CounterAssetCode)
	assert.Equal(t, "native", trade2.CounterAssetIssuer)
}

func TestNormalizeTradeAssets(t *testing.T) {
	baseAmount := "10.0"
	baseAccount := "BASEACCOUNT"
	baseAssetCode := "BASECODE"
	baseAssetType := "BASEASSETTYPE"
	baseAssetIssuer := "BASEASSETISSUER"

	n := int64(10)
	d := int64(50)

	price := hProtocol.TradePrice{
		N: n,
		D: d,
	}

	trade1 := hProtocol.Trade{
		BaseAmount:       baseAmount,
		BaseAccount:      baseAccount,
		BaseAssetCode:    baseAssetCode,
		BaseAssetType:    baseAssetType,
		BaseAssetIssuer:  baseAssetIssuer,
		CounterAssetType: "native",
		Price:            price,
	}

	NormalizeTradeAssets(&trade1)
	assert.Equal(t, baseAmount, trade1.CounterAmount)
	assert.Equal(t, baseAccount, trade1.CounterAccount)
	assert.Equal(t, baseAssetCode, trade1.CounterAssetCode)
	assert.Equal(t, baseAssetType, trade1.CounterAssetType)
	assert.Equal(t, baseAssetIssuer, trade1.CounterAssetIssuer)
	assert.Equal(t, "native", trade1.BaseAssetType)

	counterAmount := "5.0"
	counterAccount := "COUNTERACCOUNT"
	counterAssetType := "COUNTERASSETTYPE"
	counterAssetIssuer := "COUNTERASSETISSUER"

	trade2 := hProtocol.Trade{
		BaseAmount:         baseAmount,
		BaseAccount:        baseAccount,
		BaseAssetCode:      "BBB",
		BaseAssetType:      baseAssetType,
		BaseAssetIssuer:    baseAssetIssuer,
		CounterAmount:      counterAmount,
		CounterAccount:     counterAccount,
		CounterAssetCode:   "AAA",
		CounterAssetType:   counterAssetType,
		CounterAssetIssuer: counterAssetIssuer,
		Price:              price,
	}
	NormalizeTradeAssets(&trade2)
	assert.Equal(t, baseAmount, trade2.CounterAmount)
	assert.Equal(t, baseAccount, trade2.CounterAccount)
	assert.Equal(t, "BBB", trade2.CounterAssetCode)
	assert.Equal(t, baseAssetType, trade2.CounterAssetType)
	assert.Equal(t, baseAssetIssuer, trade2.CounterAssetIssuer)

	assert.Equal(t, counterAmount, trade2.BaseAmount)
	assert.Equal(t, counterAccount, trade2.BaseAccount)
	assert.Equal(t, "AAA", trade2.BaseAssetCode)
	assert.Equal(t, counterAssetType, trade2.BaseAssetType)
	assert.Equal(t, counterAssetIssuer, trade2.BaseAssetIssuer)
}
