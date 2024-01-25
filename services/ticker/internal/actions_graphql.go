package ticker

import (
	"go/services/ticker/internal/gql"
	"go/services/ticker/internal/tickerdb"
	hlog "go/support/log"
)

func StartGraphQLServer(s *tickerdb.TickerSession, l *hlog.Entry, port string) {
	graphql := gql.New(s, l)

	graphql.Serve(port)
}
