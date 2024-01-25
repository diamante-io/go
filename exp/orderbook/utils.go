package orderbook

import (
	"github.com/diamcircle/go/xdr"
)

// getPoolAssets retrieves string representations of a pool's reserves
func getPoolAssets(pool xdr.LiquidityPoolEntry) (xdr.Asset, xdr.Asset) {
	params := pool.Body.MustConstantProduct().Params
	return params.AssetA, params.AssetB
}

func max(a, b xdr.Int64) xdr.Int64 {
	if a < b {
		return b
	}
	return a
}

// positiveMin returns the smallest positive value possible
func positiveMin(a, b xdr.Int64) xdr.Int64 {
	if b <= 0 {
		return a
	}

	if b < a || a <= 0 {
		return b
	}

	return a
}
