package orderbook

import (
	"math"

	"github.com/holiman/uint256"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/xdr"
)

// There are two different exchanges that can be simulated:
//
// 1. You know how much you can *give* to the pool, and are curious about the
// resulting payout. We call this a "deposit", and you should pass
// tradeTypeDeposit.
//
// 2. You know how much you'd like to *receive* from the pool, and want to know
// how much to deposit to achieve this. We call this an "expectation", and you
// should pass tradeTypeExpectation.
const (
	tradeTypeDeposit     = iota // deposit into pool, what's the payout?
	tradeTypeExpectation = iota // expect payout, what to deposit?
)

var (
	errPoolOverflows = errors.New("Liquidity pool overflows from this exchange")
	errBadPoolType   = errors.New("Unsupported liquidity pool: must be ConstantProduct")
	errBadTradeType  = errors.New("Unknown pool exchange type requested")
	errBadAmount     = errors.New("Exchange amount must be positive")
)

// makeTrade simulates execution of an exchange with a liquidity pool.
//
// In (1), this returns the amount that would be paid out by the pool (in terms
// of the *other* asset) for depositing `amount` of `asset`.
//
// In (2), this returns the amount of `asset` you'd need to deposit to get
// `amount` of the *other* asset in return.
//
// Refer to https://github.com/diamcircle/diamcircle-protocol/blob/master/core/cap-0038.md#pathpaymentstrictsendop-and-pathpaymentstrictreceiveop
// and the calculation functions (below) for details on the exchange algorithm.
//
// Warning: If you pass an asset that is NOT one of the pool reserves, the
// behavior of this function is undefined (for performance).
func makeTrade(
	pool liquidityPool,
	asset int32,
	tradeType int,
	amount xdr.Int64,
) (xdr.Int64, error) {
	details, ok := pool.Body.GetConstantProduct()
	if !ok {
		return 0, errBadPoolType
	}

	if amount <= 0 {
		return 0, errBadAmount
	}

	// determine which asset `amount` corresponds to
	X, Y := details.ReserveA, details.ReserveB
	if pool.assetA != asset {
		X, Y = Y, X
	}

	ok = false
	var result xdr.Int64
	switch tradeType {
	case tradeTypeDeposit:
		result, ok = calculatePoolPayout(X, Y, amount, details.Params.Fee)

	case tradeTypeExpectation:
		result, ok = calculatePoolExpectation(X, Y, amount, details.Params.Fee)

	default:
		return 0, errBadTradeType
	}

	if !ok {
		// the error isn't strictly accurate (e.g. it could be div-by-0), but
		// from the caller's perspective it's true enough
		return 0, errPoolOverflows
	}
	return result, nil
}

// calculatePoolPayout calculates the amount of `reserveB` disbursed from the
// pool for a `received` amount of `reserveA` . From CAP-38:
//
//      y = floor[(1 - F) Yx / (X + x - Fx)]
//
// It returns false if the calculation overflows.
func calculatePoolPayout(reserveA, reserveB, received xdr.Int64, feeBips xdr.Int32) (xdr.Int64, bool) {
	X, Y := uint256.NewInt(uint64(reserveA)), uint256.NewInt(uint64(reserveB))
	F, x := uint256.NewInt(uint64(feeBips)), uint256.NewInt(uint64(received))

	// would this deposit overflow the reserve?
	if received > math.MaxInt64-reserveA {
		return 0, false
	}

	// We do all of the math in bips, so it's all upscaled by this value.
	maxBips := uint256.NewInt(10000)
	f := new(uint256.Int).Sub(maxBips, F) // upscaled 1 - F

	// right half: X + (1 - F)x
	denom := X.Mul(X, maxBips).Add(X, new(uint256.Int).Mul(x, f))
	if denom.IsZero() { // avoid div-by-zero panic
		return 0, false
	}

	// left half, a: (1 - F) Yx
	numer := Y.Mul(Y, x).Mul(Y, f)

	// divide & check overflow
	result := numer.Div(numer, denom)

	val := xdr.Int64(result.Uint64())
	return val, result.IsUint64() && val >= 0
}

// calculatePoolExpectation determines how much of `reserveA` you would need to
// put into a pool to get the `disbursed` amount of `reserveB`.
//
//      x = ceil[Xy / ((Y - y)(1 - F))]
//
// It returns false if the calculation overflows.
func calculatePoolExpectation(
	reserveA, reserveB, disbursed xdr.Int64, feeBips xdr.Int32,
) (xdr.Int64, bool) {
	X, Y := uint256.NewInt(uint64(reserveA)), uint256.NewInt(uint64(reserveB))
	F, y := uint256.NewInt(uint64(feeBips)), uint256.NewInt(uint64(disbursed))

	// sanity check: disbursing shouldn't underflow the reserve
	if disbursed >= reserveB {
		return 0, false
	}

	// We do all of the math in bips, so it's all upscaled by this value.
	maxBips := uint256.NewInt(10000)
	f := new(uint256.Int).Sub(maxBips, F) // upscaled 1 - F

	denom := Y.Sub(Y, y).Mul(Y, f) // right half: (Y - y)(1 - F)
	if denom.IsZero() {            // avoid div-by-zero panic
		return 0, false
	}

	numer := X.Mul(X, y).Mul(X, maxBips) // left half: Xy

	result, rem := new(uint256.Int), new(uint256.Int)
	result.Div(numer, denom)
	rem.Mod(numer, denom)

	// hacky way to ceil(): if there's a remainder, add 1
	if !rem.IsZero() {
		result.AddUint64(result, 1)
	}

	val := xdr.Int64(result.Uint64())
	return val, result.IsUint64() && val >= 0
}

// getOtherAsset returns the other asset in the liquidity pool. Note that
// doesn't check to make sure the passed in `asset` is actually part of the
// pool; behavior in that case is undefined.
func getOtherAsset(asset int32, pool liquidityPool) int32 {
	if pool.assetA == asset {
		return pool.assetB
	}
	return pool.assetA
}
