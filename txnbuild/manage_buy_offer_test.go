package txnbuild

import (
	"github.com/diamcircle/go/xdr"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManageBuyOfferValidateSellingAsset(t *testing.T) {
	kp0 := newKeypair0()
	kp1 := newKeypair1()
	sourceAccount := NewSimpleAccount(kp1.Address(), int64(41137196761092))

	buyOffer := ManageBuyOffer{
		Selling: CreditAsset{"", kp0.Address()},
		Buying:  NativeAsset{},
		Amount:  "100",
		Price:   "0.01",
		OfferID: 0,
	}

	_, err := NewTransaction(
		TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: false,
			Operations:           []Operation{&buyOffer},
			BaseFee:              MinBaseFee,
			Timebounds:           NewInfiniteTimeout(),
		},
	)
	if assert.Error(t, err) {
		expected := "validation failed for *txnbuild.ManageBuyOffer operation: Field: Selling, Error: asset code length must be between 1 and 12 characters"
		assert.Contains(t, err.Error(), expected)
	}
}

func TestManageBuyOfferValidateBuyingAsset(t *testing.T) {
	kp0 := newKeypair0()
	kp1 := newKeypair1()
	sourceAccount := NewSimpleAccount(kp1.Address(), int64(41137196761092))

	buyOffer := ManageBuyOffer{
		Selling: CreditAsset{"ABC", kp0.Address()},
		Buying:  CreditAsset{"XYZ", ""},
		Amount:  "100",
		Price:   "0.01",
		OfferID: 0,
	}

	_, err := NewTransaction(
		TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: false,
			Operations:           []Operation{&buyOffer},
			BaseFee:              MinBaseFee,
			Timebounds:           NewInfiniteTimeout(),
		},
	)
	if assert.Error(t, err) {
		expected := "validation failed for *txnbuild.ManageBuyOffer operation: Field: Buying, Error: asset issuer: public key is undefined"
		assert.Contains(t, err.Error(), expected)
	}
}

func TestManageBuyOfferValidateAmount(t *testing.T) {
	kp0 := newKeypair0()
	kp1 := newKeypair1()
	sourceAccount := NewSimpleAccount(kp1.Address(), int64(41137196761092))

	buyOffer := ManageBuyOffer{
		Selling: CreditAsset{"ABCD", kp0.Address()},
		Buying:  NativeAsset{},
		Amount:  "",
		Price:   "0.01",
		OfferID: 0,
	}

	_, err := NewTransaction(
		TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: false,
			Operations:           []Operation{&buyOffer},
			BaseFee:              MinBaseFee,
			Timebounds:           NewInfiniteTimeout(),
		},
	)
	if assert.Error(t, err) {
		expected := "validation failed for *txnbuild.ManageBuyOffer operation: Field: Amount, Error: invalid amount format:"
		assert.Contains(t, err.Error(), expected)
	}
}

func TestManageBuyOfferValidatePrice(t *testing.T) {
	kp0 := newKeypair0()
	kp1 := newKeypair1()
	sourceAccount := NewSimpleAccount(kp1.Address(), int64(41137196761092))

	buyOffer := ManageBuyOffer{
		Selling: CreditAsset{"ABCD", kp0.Address()},
		Buying:  NativeAsset{},
		Amount:  "0",
		Price:   "-0.01",
		OfferID: 0,
	}

	_, err := NewTransaction(
		TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: false,
			Operations:           []Operation{&buyOffer},
			BaseFee:              MinBaseFee,
			Timebounds:           NewInfiniteTimeout(),
		},
	)
	if assert.Error(t, err) {
		expected := "validation failed for *txnbuild.ManageBuyOffer operation: Field: Price, Error: amount can not be negative"
		assert.Contains(t, err.Error(), expected)
	}
}

func TestManageBuyOfferValidateOfferID(t *testing.T) {
	kp0 := newKeypair0()
	kp1 := newKeypair1()
	sourceAccount := NewSimpleAccount(kp1.Address(), int64(41137196761092))

	buyOffer := ManageBuyOffer{
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
			Operations:           []Operation{&buyOffer},
			BaseFee:              MinBaseFee,
			Timebounds:           NewInfiniteTimeout(),
		},
	)
	if assert.Error(t, err) {
		expected := "validation failed for *txnbuild.ManageBuyOffer operation: Field: OfferID, Error: amount can not be negative"
		assert.Contains(t, err.Error(), expected)
	}
}

func TestManageBuyOfferPrice(t *testing.T) {
	kp0 := newKeypair0()

	mbo := ManageBuyOffer{
		Selling: CreditAsset{"ABCD", kp0.Address()},
		Buying:  NativeAsset{},
		Amount:  "1",
		Price:   "0.000000001",
		OfferID: 1,
	}

	xdrOp, err := mbo.BuildXDR(false)
	assert.NoError(t, err)
	expectedPrice := xdr.Price{N: 1, D: 1000000000}
	assert.Equal(t, expectedPrice, xdrOp.Body.ManageBuyOfferOp.Price)
	assert.Equal(t, mbo.Price, mbo.price.string())
	assert.Equal(t, expectedPrice, mbo.price.toXDR())

	parsed := ManageBuyOffer{}
	assert.NoError(t, parsed.FromXDR(xdrOp, false))
	assert.Equal(t, mbo.Price, parsed.Price)
	assert.Equal(t, mbo.price, parsed.price)
}

func TestManageBuyOfferRoundtrip(t *testing.T) {
	manageBuyOffer := ManageBuyOffer{
		SourceAccount: "GB7BDSZU2Y27LYNLALKKALB52WS2IZWYBDGY6EQBLEED3TJOCVMZRH7H",
		Selling:       CreditAsset{"USD", "GB7BDSZU2Y27LYNLALKKALB52WS2IZWYBDGY6EQBLEED3TJOCVMZRH7H"},
		Buying:        NativeAsset{},
		Amount:        "100.0000000",
		Price:         "0.01",
		OfferID:       0,
	}
	testOperationsMarshallingRoundtrip(t, []Operation{&manageBuyOffer}, false)

	// with muxed accounts
	manageBuyOffer = ManageBuyOffer{
		SourceAccount: "MA7QYNF7SOWQ3GLR2BGMZEHXAVIRZA4KVWLTJJFC7MGXUA74P7UJVAAAAAAAAAAAAAJLK",
		Selling:       CreditAsset{"USD", "GB7BDSZU2Y27LYNLALKKALB52WS2IZWYBDGY6EQBLEED3TJOCVMZRH7H"},
		Buying:        NativeAsset{},
		Amount:        "100.0000000",
		Price:         "0.01",
		OfferID:       0,
	}
	testOperationsMarshallingRoundtrip(t, []Operation{&manageBuyOffer}, true)
}
