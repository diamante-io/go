package health

import (
	"net/http"

	"github.com/diamcircle/go/support/render/httpjson"
)

// PassHandler implements a simple handler that returns the most basic health
// response with a status of 'pass'.
type PassHandler struct{}

func (h PassHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Status: StatusPass,
	}
	httpjson.Render(w, response, httpjson.HEALTHJSON)
}
