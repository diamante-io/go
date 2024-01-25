package txnbuild

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBumpSequenceValidate(t *testing.T) {
	kp1 := newKeypair1()
	sourceAccount := NewSimpleAccount(kp1.Address(), int64(9606132444168199))

	bumpSequence := BumpSequence{
		BumpTo: -10,
	}

	_, err := NewTransaction(
		TransactionParams{
			SourceAccount: &sourceAccount,
			Operations:    []Operation{&bumpSequence},
			Timebounds:    NewInfiniteTimeout(),
			BaseFee:       MinBaseFee,
		},
	)
	if assert.Error(t, err) {
		expected := "validation failed for *txnbuild.BumpSequence operation: Field: BumpTo, Error: amount can not be negative"
		assert.Contains(t, err.Error(), expected)
	}
}

func TestBumpSequenceRountrip(t *testing.T) {
	bumpSequence := BumpSequence{
		SourceAccount: "GB7BDSZU2Y27LYNLALKKALB52WS2IZWYBDGY6EQBLEED3TJOCVMZRH7H",
		BumpTo:        10,
	}
	testOperationsMarshallingRoundtrip(t, []Operation{&bumpSequence}, false)

	bumpSequence = BumpSequence{
		SourceAccount: "MA7QYNF7SOWQ3GLR2BGMZEHXAVIRZA4KVWLTJJFC7MGXUA74P7UJVAAAAAAAAAAAAAJLK",
		BumpTo:        10,
	}
	testOperationsMarshallingRoundtrip(t, []Operation{&bumpSequence}, true)
}
