package txnbuild

import (
	"github.com/diamcircle/go/xdr"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePassiveSellOfferValidateBuyingAsset(t *testing.T) {
	kp1 := newKeypair1()
	sourceAccount := NewSimpleAccount(kp1.Address(), int64(41137196761100))

	createPassiveOffer := CreatePassiveSellOffer{
		Selling: NativeAsset{},
		Buying:  CreditAsset{"ABCD", ""},
		Amount:  "10",
		Price:   "1.0",
	}

	_, err := NewTransaction(
		TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: false,
			Operations:           []Operation{&createPassiveOffer},
			BaseFee:              MinBaseFee,
			Timebounds:           NewInfiniteTimeout(),
		},
	)
	if assert.Error(t, err) {
		expected := "validation failed for *txnbuild.CreatePassiveSellOffer operation: Field: Buying, Error: asset issuer: public key is undefined"
		assert.Contains(t, err.Error(), expected)
	}
}

func TestCreatePassiveSellOfferValidateSellingAsset(t *testing.T) {
	kp0 := newKeypair0()
	kp1 := newKeypair1()
	sourceAccount := NewSimpleAccount(kp1.Address(), int64(41137196761100))

	createPassiveOffer := CreatePassiveSellOffer{
		Selling: CreditAsset{"ABCD0123456789", kp0.Address()},
		Buying:  NativeAsset{},
		Amount:  "10",
		Price:   "1.0",
	}

	_, err := NewTransaction(
		TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: false,
			Operations:           []Operation{&createPassiveOffer},
			BaseFee:              MinBaseFee,
			Timebounds:           NewInfiniteTimeout(),
		},
	)
	if assert.Error(t, err) {
		expected := `validation failed for *txnbuild.CreatePassiveSellOffer operation: Field: Selling, Error: asset code length must be between 1 and 12 characters`
		assert.Contains(t, err.Error(), expected)
	}
}

func TestCreatePassiveSellOfferValidateAmount(t *testing.T) {
	kp0 := newKeypair0()
	kp1 := newKeypair1()
	sourceAccount := NewSimpleAccount(kp1.Address(), int64(41137196761100))

	createPassiveOffer := CreatePassiveSellOffer{
		Selling: CreditAsset{"ABCD", kp0.Address()},
		Buying:  NativeAsset{},
		Amount:  "-3",
		Price:   "1.0",
	}

	_, err := NewTransaction(
		TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: false,
			Operations:           []Operation{&createPassiveOffer},
			BaseFee:              MinBaseFee,
			Timebounds:           NewInfiniteTimeout(),
		},
	)
	if assert.Error(t, err) {
		expected := `validation failed for *txnbuild.CreatePassiveSellOffer operation: Field: Amount, Error: amount can not be negative`
		assert.Contains(t, err.Error(), expected)
	}
}

func TestCreatePassiveSellOfferValidatePrice(t *testing.T) {
	kp0 := newKeypair0()
	kp1 := newKeypair1()
	sourceAccount := NewSimpleAccount(kp1.Address(), int64(41137196761100))

	createPassiveOffer := CreatePassiveSellOffer{
		Selling: CreditAsset{"ABCD", kp0.Address()},
		Buying:  NativeAsset{},
		Amount:  "3",
		Price:   "-1.0",
	}

	_, err := NewTransaction(
		TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: false,
			Operations:           []Operation{&createPassiveOffer},
			BaseFee:              MinBaseFee,
			Timebounds:           NewInfiniteTimeout(),
		},
	)
	if assert.Error(t, err) {
		expected := `validation failed for *txnbuild.CreatePassiveSellOffer operation: Field: Price, Error: amount can not be negative`
		assert.Contains(t, err.Error(), expected)
	}
}

func TestCreatePassiveSellOfferPrice(t *testing.T) {
	kp0 := newKeypair0()

	offer := CreatePassiveSellOffer{
		Selling:       CreditAsset{"ABCD", kp0.Address()},
		Buying:        NativeAsset{},
		Amount:        "1",
		Price:         "0.000000001",
		SourceAccount: kp0.Address(),
	}

	xdrOp, err := offer.BuildXDR(false)
	assert.NoError(t, err)
	expectedPrice := xdr.Price{N: 1, D: 1000000000}
	assert.Equal(t, expectedPrice, xdrOp.Body.CreatePassiveSellOfferOp.Price)
	assert.Equal(t, offer.Price, offer.price.string())
	assert.Equal(t, expectedPrice, offer.price.toXDR())

	parsed := CreatePassiveSellOffer{}
	assert.NoError(t, parsed.FromXDR(xdrOp, false))
	assert.Equal(t, offer.Price, parsed.Price)
	assert.Equal(t, offer.price, parsed.price)
}
