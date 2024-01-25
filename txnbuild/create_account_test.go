package txnbuild

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccountValidateDestination(t *testing.T) {
	kp0 := newKeypair0()
	sourceAccount := NewSimpleAccount(kp0.Address(), int64(9605939170639897))

	createAccount := CreateAccount{
		Destination: "",
		Amount:      "43",
	}

	_, err := NewTransaction(
		TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: false,
			Operations:           []Operation{&createAccount},
			BaseFee:              MinBaseFee,
			Timebounds:           NewInfiniteTimeout(),
		},
	)
	if assert.Error(t, err) {
		expected := "validation failed for *txnbuild.CreateAccount operation: Field: Destination, Error: public key is undefined"
		assert.Contains(t, err.Error(), expected)
	}
}

func TestCreateAccountValidateAmount(t *testing.T) {
	kp0 := newKeypair0()
	sourceAccount := NewSimpleAccount(kp0.Address(), int64(9605939170639897))

	createAccount := CreateAccount{
		Destination: "GDYNXQFHU6W5RBW2CCCDDAAU3TMTSU2RMGIBM6HGHAR4NJJKY3IJETHT",
		Amount:      "",
	}

	_, err := NewTransaction(
		TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: false,
			Operations:           []Operation{&createAccount},
			BaseFee:              MinBaseFee,
			Timebounds:           NewInfiniteTimeout(),
		},
	)

	if assert.Error(t, err) {
		expected := "validation failed for *txnbuild.CreateAccount operation: Field: Amount, Error: invalid amount format"
		assert.Contains(t, err.Error(), expected)
	}
}

func TestCreateAccountRoundtrip(t *testing.T) {
	createAccount := CreateAccount{
		SourceAccount: "GDYNXQFHU6W5RBW2CCCDDAAU3TMTSU2RMGIBM6HGHAR4NJJKY3IJETHT",
		Destination:   "GDYNXQFHU6W5RBW2CCCDDAAU3TMTSU2RMGIBM6HGHAR4NJJKY3IJETHT",
		Amount:        "1.0000000",
	}
	testOperationsMarshallingRoundtrip(t, []Operation{&createAccount}, false)

	// with muxed accounts
	createAccount = CreateAccount{
		SourceAccount: "MA7QYNF7SOWQ3GLR2BGMZEHXAVIRZA4KVWLTJJFC7MGXUA74P7UJVAAAAAAAAAAAAAJLK",
		Destination:   "GDYNXQFHU6W5RBW2CCCDDAAU3TMTSU2RMGIBM6HGHAR4NJJKY3IJETHT",
		Amount:        "1.0000000",
	}
	testOperationsMarshallingRoundtrip(t, []Operation{&createAccount}, true)

}
