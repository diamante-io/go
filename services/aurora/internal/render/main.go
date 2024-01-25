package render

import (
	"net/http"

	"github.com/adjust/goautoneg"
	"github.com/diamcircle/go/support/log"
)

// Negotiate inspects the Accept header of the provided request and determines
// what the most appropriate response type should be.  Defaults to HAL.
func Negotiate(r *http.Request) string {
	ctx := r.Context()
	alternatives := []string{MimeHal, MimeJSON, MimeEventStream, MimeRaw}
	accept := r.Header.Get("Accept")

	if accept == "" {
		return MimeHal
	}

	result := goautoneg.Negotiate(r.Header.Get("Accept"), alternatives)

	log.Ctx(ctx).WithFields(log.F{
		"content_type": result,
		"accept":       accept,
	}).Debug("Negotiated content type")

	return result
}
