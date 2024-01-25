package txnbuild

import (
	"github.com/diamcircle/go/xdr"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManageSellOfferValidateSellingAsset(t *testing.T) {
	kp0 := newKeypair0()
	kp1 := newKeypair1()
	sourceAccount := NewSimpleAccount(kp1.Address(), int64(41137196761092))

	selling := CreditAsset{"", kp0.Address()}
	buying := NativeAsset{}
	sellAmount := "100"
	price := "0.01"
	createOffer, err := CreateOfferOp(selling, buying, sellAmount, price)
	check(err)

	_, err = NewTransaction(
		TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: false,
			Operations:           []Operation{&createOffer},
			BaseFee:              MinBaseFee,
			Timebounds:           NewInfiniteTimeout(),
		},
	)
	if assert.Error(t, err) {
		expected := "validation failed for *txnbuild.ManageSellOffer operation: Field: Selling, Error: asset code length must be between 1 and 12 characters"
		assert.Contains(t, err.Error(), expected)
	}
}

func TestManageSellOfferValidateBuyingAsset(t *testing.T) {
	kp0 := newKeypair0()
	kp1 := newKeypair1()
	sourceAccount := NewSimpleAccount(kp1.Address(), int64(41137196761092))

	selling := NativeAsset{}
	buying := CreditAsset{"", kp0.Address()}
	sellAmount := "100"
	price := "0.01"
	createOffer, err := CreateOfferOp(selling, buying, sellAmount, price)
	check(err)

	_, err = NewTransaction(
		TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: false,
			Operations:           []Operation{&createOffer},
			BaseFee:              MinBaseFee,
			Timebounds:           NewInfiniteTimeout(),
		},
	)
	if assert.Error(t, err) {
		expected := "validation failed for *txnbuild.ManageSellOffer operation: Field: Buying, Error: asset code length must be between 1 and 12 characters"
		assert.Contains(t, err.Error(), expected)
	}
}

func TestManageSellOfferValidateAmount(t *testing.T) {
	kp0 := newKeypair0()
	kp1 := newKeypair1()
	sourceAccount := NewSimpleAccount(kp1.Address(), int64(41137196761092))

	selling := NativeAsset{}
	buying := CreditAsset{"ABCD", kp0.Address()}
	sellAmount := "-1"
	price := "0.01"
	createOffer, err := CreateOfferOp(selling, buying, sellAmount, price)
	check(err)

	_, err = NewTransaction(
		TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: false,
			Operations:           []Operation{&createOffer},
			BaseFee:              MinBaseFee,
			Timebounds:           NewInfiniteTimeout(),
		},
	)
	if assert.Error(t, err) {
		expected := "validation failed for *txnbuild.ManageSellOffer operation: Field: Amount, Error: amount can not be negative"
		assert.Contains(t, err.Error(), expected)
	}
}

func TestManageSellOfferValidatePrice(t *testing.T) {
	kp0 := newKeypair0()
	kp1 := newKeypair1()
	sourceAccount := NewSimpleAccount(kp1.Address(), int64(41137196761092))

	selling := NativeAsset{}
	buying := CreditAsset{"ABCD", kp0.Address()}
	sellAmount := "0"
	price := "-0.01"
	createOffer, err := CreateOfferOp(selling, buying, sellAmount, price)
	check(err)

	_, err = NewTransaction(
		TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: false,
			Operations:           []Operation{&createOffer},
			BaseFee:              MinBaseFee,
			Timebounds:           NewInfiniteTimeout(),
		},
	)
	if assert.Error(t, err) {
		expected := "validation failed for *txnbuild.ManageSellOffer operation: Field: Price, Error: amount can not be negative"
		assert.Contains(t, err.Error(), expected)
	}
}

func TestManageSellOfferValidateOfferID(t *testing.T) {
	kp0 := newKeypair0()
	kp1 := newKeypair1()
	sourceAccount := NewSimpleAccount(kp1.Address(), int64(41137196761092))

	mso := ManageSellOffer{
		Selling: CreditAsset{"ABCD", kp0.Address()},
		Buying:  NativeAsset{},
		Amount:  "0",
		Price:   "0.01",
		OfferID: -1,
	}

	_, err := NewTransaction(
		TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: false,
			Operations:           []Operation{&mso},
			BaseFee:              MinBaseFee,
			Timebounds:           NewInfiniteTimeout(),
		},
	)
	if assert.Error(t, err) {
		expected := "validation failed for *txnbuild.ManageSellOffer operation: Field: OfferID, Error: amount can not be negative"
		assert.Contains(t, err.Error(), expected)
	}
}

func TestManageSellOfferPrice(t *testing.T) {
	kp0 := newKeypair0()

	mso := ManageSellOffer{
		Selling: CreditAsset{"ABCD", kp0.Address()},
		Buying:  NativeAsset{},
		Amount:  "1",
		Price:   "0.000000001",
		OfferID: 1,
	}

	xdrOp, err := mso.BuildXDR(false)
	assert.NoError(t, err)
	expectedPrice := xdr.Price{N: 1, D: 1000000000}
	assert.Equal(t, expectedPrice, xdrOp.Body.ManageSellOfferOp.Price)
	assert.Equal(t, mso.Price, mso.price.string())
	assert.Equal(t, expectedPrice, mso.price.toXDR())

	parsed := ManageSellOffer{}
	assert.NoError(t, parsed.FromXDR(xdrOp, false))
	assert.Equal(t, mso.Price, parsed.Price)
	assert.Equal(t, mso.price, parsed.price)
}

func TestManageSellOfferRoundtrip(t *testing.T) {
	manageSellOffer := ManageSellOffer{
		SourceAccount: "GB7BDSZU2Y27LYNLALKKALB52WS2IZWYBDGY6EQBLEED3TJOCVMZRH7H",
		Selling:       CreditAsset{"USD", "GB7BDSZU2Y27LYNLALKKALB52WS2IZWYBDGY6EQBLEED3TJOCVMZRH7H"},
		Buying:        NativeAsset{},
		Amount:        "100.0000000",
		Price:         "0.01",
		OfferID:       0,
	}
	testOperationsMarshallingRoundtrip(t, []Operation{&manageSellOffer}, false)

	// with muxed accounts
	manageSellOffer = ManageSellOffer{
		SourceAccount: "MA7QYNF7SOWQ3GLR2BGMZEHXAVIRZA4KVWLTJJFC7MGXUA74P7UJVAAAAAAAAAAAAAJLK",
		Selling:       CreditAsset{"USD", "GB7BDSZU2Y27LYNLALKKALB52WS2IZWYBDGY6EQBLEED3TJOCVMZRH7H"},
		Buying:        NativeAsset{},
		Amount:        "100.0000000",
		Price:         "0.01",
		OfferID:       0,
	}
	testOperationsMarshallingRoundtrip(t, []Operation{&manageSellOffer}, true)
}
