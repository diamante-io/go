package gql

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/diamcircle/go/services/ticker/internal/gql/static"
)

// GraphiQL is an in-browser IDE for exploring GraphiQL APIs.
// This handler returns GraphiQL when requested.
//
// For more information, see https://github.com/graphql/graphiql.
type GraphiQL struct{}

func (h GraphiQL) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		buf := bytes.Buffer{}
		fmt.Fprint(&buf, `{"error": "Only GET is allowed"}`)
		_, _ = w.Write(buf.Bytes())
		return
	}

	graphiql, _ := static.Asset("graphiql.html")
	_, _ = w.Write(graphiql)
}
