package processors

import (
	"github.com/guregu/null"
	logpkg "github.com/diamcircle/go/support/log"
	"github.com/diamcircle/go/xdr"
)

var log = logpkg.DefaultLogger.WithField("service", "ingest")

const maxBatchSize = 100000

func ledgerEntrySponsorToNullString(entry xdr.LedgerEntry) null.String {
	sponsoringID := entry.SponsoringID()

	var sponsor null.String
	if sponsoringID != nil {
		sponsor.SetValid((*sponsoringID).Address())
	}

	return sponsor
}
