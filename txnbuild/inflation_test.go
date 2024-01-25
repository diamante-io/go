package txnbuild

import "testing"

func TestInflationRoundtrip(t *testing.T) {
	inflation := Inflation{
		SourceAccount: "GB7BDSZU2Y27LYNLALKKALB52WS2IZWYBDGY6EQBLEED3TJOCVMZRH7H",
	}
	testOperationsMarshallingRoundtrip(t, []Operation{&inflation}, false)

	// with muxed accounts
	inflation = Inflation{
		SourceAccount: "MA7QYNF7SOWQ3GLR2BGMZEHXAVIRZA4KVWLTJJFC7MGXUA74P7UJVAAAAAAAAAAAAAJLK",
	}
	testOperationsMarshallingRoundtrip(t, []Operation{&inflation}, true)
}
