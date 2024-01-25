package txnbuild

import (
	"github.com/diamcircle/go/amount"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/xdr"
)

// ManageBuyOffer represents the Diamcircle manage buy offer operation. See
// https://developers.diamcircle.org/docs/start/list-of-operations/
type ManageBuyOffer struct {
	Selling       Asset
	Buying        Asset
	Amount        string
	Price         string
	price         price
	OfferID       int64
	SourceAccount string
}

// BuildXDR for ManageBuyOffer returns a fully configured XDR Operation.
func (mo *ManageBuyOffer) BuildXDR(withMuxedAccounts bool) (xdr.Operation, error) {
	xdrSelling, err := mo.Selling.ToXDR()
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to set XDR 'Selling' field")
	}

	xdrBuying, err := mo.Buying.ToXDR()
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to set XDR 'Buying' field")
	}

	xdrAmount, err := amount.Parse(mo.Amount)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to parse 'Amount'")
	}

	if err = mo.price.parse(mo.Price); err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to parse 'Price'")
	}

	opType := xdr.OperationTypeManageBuyOffer
	xdrOp := xdr.ManageBuyOfferOp{
		Selling:   xdrSelling,
		Buying:    xdrBuying,
		BuyAmount: xdrAmount,
		Price:     mo.price.toXDR(),
		OfferId:   xdr.Int64(mo.OfferID),
	}
	body, err := xdr.NewOperationBody(opType, xdrOp)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to build XDR OperationBody")
	}

	op := xdr.Operation{Body: body}
	if withMuxedAccounts {
		SetOpSourceMuxedAccount(&op, mo.SourceAccount)
	} else {
		SetOpSourceAccount(&op, mo.SourceAccount)
	}
	return op, nil
}

// FromXDR for ManageBuyOffer initialises the txnbuild struct from the corresponding xdr Operation.
func (mo *ManageBuyOffer) FromXDR(xdrOp xdr.Operation, withMuxedAccounts bool) error {
	result, ok := xdrOp.Body.GetManageBuyOfferOp()
	if !ok {
		return errors.New("error parsing manage_buy_offer operation from xdr")
	}

	mo.SourceAccount = accountFromXDR(xdrOp.SourceAccount, withMuxedAccounts)
	mo.OfferID = int64(result.OfferId)
	mo.Amount = amount.String(result.BuyAmount)
	if result.Price != (xdr.Price{}) {
		mo.price.fromXDR(result.Price)
		mo.Price = mo.price.string()
	}
	buyingAsset, err := assetFromXDR(result.Buying)
	if err != nil {
		return errors.Wrap(err, "error parsing buying_asset in manage_buy_offer operation")
	}
	mo.Buying = buyingAsset

	sellingAsset, err := assetFromXDR(result.Selling)
	if err != nil {
		return errors.Wrap(err, "error parsing selling_asset in manage_buy_offer operation")
	}
	mo.Selling = sellingAsset
	return nil
}

// Validate for ManageBuyOffer validates the required struct fields. It returns an error if any
// of the fields are invalid. Otherwise, it returns nil.
func (mo *ManageBuyOffer) Validate(withMuxedAccounts bool) error {
	return validateOffer(mo.Buying, mo.Selling, mo.Amount, mo.Price, mo.OfferID)
}

// GetSourceAccount returns the source account of the operation, or the empty string if not
// set.
func (mo *ManageBuyOffer) GetSourceAccount() string {
	return mo.SourceAccount
}
