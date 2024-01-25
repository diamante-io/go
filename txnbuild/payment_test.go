package txnbuild

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaymentValidateDestination(t *testing.T) {
	kp0 := newKeypair0()
	sourceAccount := NewSimpleAccount(kp0.Address(), int64(9605939170639898))

	payment := Payment{
		Destination: "",
		Amount:      "10",
		Asset:       NativeAsset{},
	}

	_, err := NewTransaction(
		TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: false,
			Operations:           []Operation{&payment},
			BaseFee:              MinBaseFee,
			Timebounds:           NewInfiniteTimeout(),
		},
	)
	if assert.Error(t, err) {
		expected := "validation failed for *txnbuild.Payment operation: Field: Destination"
		assert.Contains(t, err.Error(), expected)
	}
}

func TestPaymentValidateAmount(t *testing.T) {
	kp0 := newKeypair0()
	sourceAccount := NewSimpleAccount(kp0.Address(), int64(9605939170639898))

	payment := Payment{
		Destination: "GB7BDSZU2Y27LYNLALKKALB52WS2IZWYBDGY6EQBLEED3TJOCVMZRH7H",
		Amount:      "ten",
		Asset:       NativeAsset{},
	}

	_, err := NewTransaction(
		TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: false,
			Operations:           []Operation{&payment},
			BaseFee:              MinBaseFee,
			Timebounds:           NewInfiniteTimeout(),
		},
	)
	if assert.Error(t, err) {
		expected := "validation failed for *txnbuild.Payment operation: Field: Amount, Error: invalid amount format: ten"
		assert.Contains(t, err.Error(), expected)
	}
}

func TestPaymentValidateAsset(t *testing.T) {
	kp0 := newKeypair0()
	sourceAccount := NewSimpleAccount(kp0.Address(), int64(9605939170639898))

	payment := Payment{
		Destination: "GB7BDSZU2Y27LYNLALKKALB52WS2IZWYBDGY6EQBLEED3TJOCVMZRH7H",
		Amount:      "10",
		Asset:       CreditAsset{},
	}

	_, err := NewTransaction(
		TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: false,
			Operations:           []Operation{&payment},
			BaseFee:              MinBaseFee,
			Timebounds:           NewInfiniteTimeout(),
		},
	)
	if assert.Error(t, err) {
		expected := "validation failed for *txnbuild.Payment operation: Field: Asset, Error: asset code length must be between 1 and 12 characters"
		assert.Contains(t, err.Error(), expected)
	}
}

func TestPaymentRoundtrip(t *testing.T) {
	payment := Payment{
		SourceAccount: "GB7BDSZU2Y27LYNLALKKALB52WS2IZWYBDGY6EQBLEED3TJOCVMZRH7H",
		Destination:   "GB7BDSZU2Y27LYNLALKKALB52WS2IZWYBDGY6EQBLEED3TJOCVMZRH7H",
		Amount:        "10.0000000",
		Asset:         NativeAsset{},
	}
	testOperationsMarshallingRoundtrip(t, []Operation{&payment}, false)

	// with muxed accounts
	payment = Payment{
		SourceAccount: "MA7QYNF7SOWQ3GLR2BGMZEHXAVIRZA4KVWLTJJFC7MGXUA74P7UJVAAAAAAAAAAAAAJLK",
		Destination:   "MA7QYNF7SOWQ3GLR2BGMZEHXAVIRZA4KVWLTJJFC7MGXUA74P7UJVAAAAAAAAAAAAAJLK",
		Amount:        "10.0000000",
		Asset:         NativeAsset{},
	}
	testOperationsMarshallingRoundtrip(t, []Operation{&payment}, true)
}
