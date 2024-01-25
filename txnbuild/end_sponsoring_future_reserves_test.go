package txnbuild

import "testing"

func TestEndSponsoringFutureReservesRoundTrip(t *testing.T) {
	withoutMuxedAccounts := &EndSponsoringFutureReserves{SourceAccount: "GB7BDSZU2Y27LYNLALKKALB52WS2IZWYBDGY6EQBLEED3TJOCVMZRH7H"}
	testOperationsMarshallingRoundtrip(t, []Operation{withoutMuxedAccounts}, false)
	withMuxedAccounts := &EndSponsoringFutureReserves{SourceAccount: "MA7QYNF7SOWQ3GLR2BGMZEHXAVIRZA4KVWLTJJFC7MGXUA74P7UJVAAAAAAAAAAAAAJLK"}
	testOperationsMarshallingRoundtrip(t, []Operation{withMuxedAccounts}, true)
}
