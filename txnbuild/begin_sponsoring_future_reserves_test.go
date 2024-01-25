package txnbuild

import (
	"testing"
)

func TestBeginSponsoringFutureReservesRoundTrip(t *testing.T) {
	beginSponsoring := &BeginSponsoringFutureReserves{
		SponsoredID: newKeypair1().Address(),
	}

	testOperationsMarshallingRoundtrip(t, []Operation{beginSponsoring}, false)
}
