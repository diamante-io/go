package errors

import (
	"fmt"
	"net/http"

	"github.com/getsentry/raven-go"
	"github.com/go-errors/errors"
)

// FromPanic extracts the err from the result of a recover() call.
func FromPanic(rec interface{}) error {
	err, ok := rec.(error)
	if !ok {
		err = fmt.Errorf("%s", rec)
	}

	return errors.Wrap(err, 4)
}

// ReportToSentry reports err to the configured sentry server.  Optionally,
// specifying a non-nil `r` will include information in the report about the
// current http request.
func ReportToSentry(err error, r *http.Request) {
	st := raven.NewStacktrace(4, 3, []string{"github.org/diamcircle"})
	exc := raven.NewException(err, st)

	var packet *raven.Packet
	if r != nil {
		h := raven.NewHttp(r)
		packet = raven.NewPacket(err.Error(), exc, h)
	} else {
		packet = raven.NewPacket(err.Error(), exc)
	}

	raven.Capture(packet, nil)
}
