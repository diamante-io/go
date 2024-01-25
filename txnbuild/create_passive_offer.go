package txnbuild

import (
	"go/amount"
	"go/support/errors"
	"go/xdr"
)

// CreatePassiveSellOffer represents the diamcircle create passive offer operation. See
// https://developers.diamcircle.org/docs/start/list-of-operations/
type CreatePassiveSellOffer struct {
	Selling       Asset
	Buying        Asset
	Amount        string
	Price         string
	price         price
	SourceAccount string
}

// BuildXDR for CreatePassiveSellOffer returns a fully configured XDR Operation.
func (cpo *CreatePassiveSellOffer) BuildXDR(withMuxedAccounts bool) (xdr.Operation, error) {
	xdrSelling, err := cpo.Selling.ToXDR()
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to set XDR 'Selling' field")
	}

	xdrBuying, err := cpo.Buying.ToXDR()
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to set XDR 'Buying' field")
	}

	xdrAmount, err := amount.Parse(cpo.Amount)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to parse 'Amount'")
	}

	if err = cpo.price.parse(cpo.Price); err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to parse 'Price'")
	}

	xdrOp := xdr.CreatePassiveSellOfferOp{
		Selling: xdrSelling,
		Buying:  xdrBuying,
		Amount:  xdrAmount,
		Price:   cpo.price.toXDR(),
	}

	opType := xdr.OperationTypeCreatePassiveSellOffer
	body, err := xdr.NewOperationBody(opType, xdrOp)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to build XDR OperationBody")
	}
	op := xdr.Operation{Body: body}
	if withMuxedAccounts {
		SetOpSourceMuxedAccount(&op, cpo.SourceAccount)
	} else {
		SetOpSourceAccount(&op, cpo.SourceAccount)
	}
	return op, nil
}

// FromXDR for CreatePassiveSellOffer initialises the txnbuild struct from the corresponding xdr Operation.
func (cpo *CreatePassiveSellOffer) FromXDR(xdrOp xdr.Operation, withMuxedAccounts bool) error {
	result, ok := xdrOp.Body.GetCreatePassiveSellOfferOp()
	if !ok {
		return errors.New("error parsing create_passive_sell_offer operation from xdr")
	}

	cpo.SourceAccount = accountFromXDR(xdrOp.SourceAccount, withMuxedAccounts)
	cpo.Amount = amount.String(result.Amount)
	if result.Price != (xdr.Price{}) {
		cpo.price.fromXDR(result.Price)
		cpo.Price = cpo.price.string()
	}
	buyingAsset, err := assetFromXDR(result.Buying)
	if err != nil {
		return errors.Wrap(err, "error parsing buying_asset in create_passive_sell_offer operation")
	}
	cpo.Buying = buyingAsset

	sellingAsset, err := assetFromXDR(result.Selling)
	if err != nil {
		return errors.Wrap(err, "error parsing selling_asset in create_passive_sell_offer operation")
	}
	cpo.Selling = sellingAsset
	return nil
}

// Validate for CreatePassiveSellOffer validates the required struct fields. It returns an error if any
// of the fields are invalid. Otherwise, it returns nil.
func (cpo *CreatePassiveSellOffer) Validate(withMuxedAccounts bool) error {
	return validatePassiveOffer(cpo.Buying, cpo.Selling, cpo.Amount, cpo.Price)
}

// GetSourceAccount returns the source account of the operation, or the empty string if not
// set.
func (cpo *CreatePassiveSellOffer) GetSourceAccount() string {
	return cpo.SourceAccount
}
