package render

import (
	"net/http"

	"go/support/log"

	"github.com/adjust/goautoneg"
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
