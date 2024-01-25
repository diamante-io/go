package hal

import (
	"net/http"

	"github.com/diamcircle/go/support/render/httpjson"
)

// Render write data to w, after marshalling to json
func Render(w http.ResponseWriter, data interface{}) {
	httpjson.Render(w, data, httpjson.HALJSON)
}
