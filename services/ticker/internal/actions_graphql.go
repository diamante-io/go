package ticker

import (
	"github.com/diamcircle/go/services/ticker/internal/gql"
	"github.com/diamcircle/go/services/ticker/internal/tickerdb"
	hlog "github.com/diamcircle/go/support/log"
)

func StartGraphQLServer(s *tickerdb.TickerSession, l *hlog.Entry, port string) {
	graphql := gql.New(s, l)

	graphql.Serve(port)
}
