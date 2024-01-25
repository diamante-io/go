//lint:file-ignore U1001 Ignore all unused code, staticcheck doesn't understand testify/suite
package txnbuild

import (
	"go/amount"
	"go/support/errors"
	"go/xdr"
)

// LiquidityPoolDeposit represents the diamcircle liquidity pool deposit operation. See
// https://developers.diamcircle.org/docs/start/list-of-operations/
type LiquidityPoolDeposit struct {
	SourceAccount   string
	LiquidityPoolID LiquidityPoolId
	MaxAmountA      string
	MaxAmountB      string
	MinPrice        string
	MaxPrice        string
}

// NewLiquidityPoolDeposit creates a new LiquidityPoolDeposit operation,
// checking the ordering assets so we generate the correct pool id. minPrice,
// and maxPrice are in terms of a/b. Each AssetAmount is a pair of the asset
// with the maximum amount of that asset to deposit.
func NewLiquidityPoolDeposit(
	sourceAccount string,
	a, b AssetAmount,
	minPrice,
	maxPrice string,
) (LiquidityPoolDeposit, error) {
	if b.Asset.LessThan(a.Asset) {
		return LiquidityPoolDeposit{}, errors.New("AssetA must be <= AssetB")
	}

	poolId, err := NewLiquidityPoolId(a.Asset, b.Asset)
	if err != nil {
		return LiquidityPoolDeposit{}, err
	}

	return LiquidityPoolDeposit{
		SourceAccount:   sourceAccount,
		LiquidityPoolID: poolId,
		MaxAmountA:      a.Amount,
		MaxAmountB:      b.Amount,
		MinPrice:        minPrice,
		MaxPrice:        maxPrice,
	}, nil
}

// BuildXDR for LiquidityPoolDeposit returns a fully configured XDR Operation.
func (lpd *LiquidityPoolDeposit) BuildXDR(withMuxedAccounts bool) (xdr.Operation, error) {
	xdrLiquidityPoolId, err := lpd.LiquidityPoolID.ToXDR()
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "couldn't build liquidity pool ID XDR")
	}

	xdrMaxAmountA, err := amount.Parse(lpd.MaxAmountA)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to parse 'MaxAmountA'")
	}

	xdrMaxAmountB, err := amount.Parse(lpd.MaxAmountB)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to parse 'MaxAmountB'")
	}

	var minPrice, maxPrice price
	err = minPrice.parse(lpd.MinPrice)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to parse 'MinPrice'")
	}
	err = maxPrice.parse(lpd.MaxPrice)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to parse 'MaxPrice'")
	}

	xdrOp := xdr.LiquidityPoolDepositOp{
		LiquidityPoolId: xdrLiquidityPoolId,
		MaxAmountA:      xdrMaxAmountA,
		MaxAmountB:      xdrMaxAmountB,
		MinPrice:        minPrice.toXDR(),
		MaxPrice:        maxPrice.toXDR(),
	}

	opType := xdr.OperationTypeLiquidityPoolDeposit
	body, err := xdr.NewOperationBody(opType, xdrOp)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to build XDR OperationBody")
	}
	op := xdr.Operation{Body: body}
	if withMuxedAccounts {
		SetOpSourceMuxedAccount(&op, lpd.SourceAccount)
	} else {
		SetOpSourceAccount(&op, lpd.SourceAccount)
	}
	return op, nil
}

// FromXDR for LiquidityPoolDeposit initializes the txnbuild struct from the corresponding xdr Operation.
func (lpd *LiquidityPoolDeposit) FromXDR(xdrOp xdr.Operation, withMuxedAccounts bool) error {
	result, ok := xdrOp.Body.GetLiquidityPoolDepositOp()
	if !ok {
		return errors.New("error parsing liquidity_pool_deposit operation from xdr")
	}

	liquidityPoolID, err := liquidityPoolIdFromXDR(result.LiquidityPoolId)
	if err != nil {
		return errors.New("error parsing LiquidityPoolId in liquidity_pool_deposit operation from xdr")
	}
	lpd.LiquidityPoolID = liquidityPoolID

	lpd.SourceAccount = accountFromXDR(xdrOp.SourceAccount, withMuxedAccounts)
	lpd.MaxAmountA = amount.String(result.MaxAmountA)
	lpd.MaxAmountB = amount.String(result.MaxAmountB)
	if result.MinPrice != (xdr.Price{}) {
		lpd.MinPrice = priceFromXDR(result.MinPrice).string()
	}
	if result.MaxPrice != (xdr.Price{}) {
		lpd.MaxPrice = priceFromXDR(result.MaxPrice).string()
	}

	return nil
}

// Validate for LiquidityPoolDeposit validates the required struct fields. It returns an error if any of the fields are
// invalid. Otherwise, it returns nil.
func (lpd *LiquidityPoolDeposit) Validate(withMuxedAccounts bool) error {
	err := validateAmount(lpd.MaxAmountA)
	if err != nil {
		return NewValidationError("MaxAmountA", err.Error())
	}

	err = validateAmount(lpd.MaxAmountB)
	if err != nil {
		return NewValidationError("MaxAmountB", err.Error())
	}

	err = validateAmount(lpd.MinPrice)
	if err != nil {
		return NewValidationError("MinPrice", err.Error())
	}

	err = validateAmount(lpd.MaxPrice)
	if err != nil {
		return NewValidationError("MaxPrice", err.Error())
	}

	return nil
}

// GetSourceAccount returns the source account of the operation, or nil if not
// set.
func (lpd *LiquidityPoolDeposit) GetSourceAccount() string {
	return lpd.SourceAccount
}
