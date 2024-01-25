package ticker

import (
	"fmt"
	"testing"

	hProtocol "github.com/diamcircle/go/protocols/aurora"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProtocolTradeToDBTrade_priceMapping(t *testing.T) {
	testCases := []struct {
		N         int64
		D         int64
		WantPrice float64
	}{
		{100, 200, 2},
		{1, 2, 2},
		{1187492342, 283724929, 0.23892779680763618},
	}
	for _, tc := range testCases {
		name := fmt.Sprintf("%d/%d=%f", tc.N, tc.D, tc.WantPrice)
		t.Run(name, func(t *testing.T) {
			hpt := hProtocol.Trade{
				BaseAmount:    "0",
				CounterAmount: "0",
				Price: hProtocol.TradePrice{
					N: tc.N,
					D: tc.D,
				},
			}
			dbTrade, err := hProtocolTradeToDBTrade(hpt, 0, 0)
			require.NoError(t, err)
			assert.Equal(t, tc.WantPrice, dbTrade.Price)
		})
	}
}
