package gql

import (
	"net/http"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/diamcircle/go/services/ticker/internal/gql/static"
	"github.com/diamcircle/go/services/ticker/internal/tickerdb"
	hlog "github.com/diamcircle/go/support/log"
)

// asset represents a Diamcircle asset, with some type
// adaptations to match the GraphQL type system
type asset struct {
	Code                        string
	IssuerAccount               string
	Type                        string
	NumAccounts                 int32
	AuthRequired                bool
	AuthRevocable               bool
	Amount                      float64
	AssetControlledByDomain     bool
	AnchorAssetCode             string
	AnchorAssetType             string
	IsValid                     bool
	DisplayDecimals             BigInt
	Name                        string
	Desc                        string
	Conditions                  string
	IsAssetAnchored             bool
	FixedNumber                 BigInt
	MaxNumber                   BigInt
	IsUnlimited                 bool
	RedemptionInstructions      string
	CollateralAddresses         string
	CollateralAddressSignatures string
	Countries                   string
	Status                      string
	IssuerID                    int32
	OrderbookStats              orderbookStats
}

// partialMarket represents the aggregated market data for a
// specific pair of assets since <Since>
type partialMarket struct {
	TradePair            string
	BaseAssetCode        string
	BaseAssetIssuer      string
	CounterAssetCode     string
	CounterAssetIssuer   string
	BaseVolume           float64
	CounterVolume        float64
	TradeCount           int32
	Open                 float64
	Low                  float64
	High                 float64
	Change               float64
	Close                float64
	IntervalStart        graphql.Time
	FirstLedgerCloseTime graphql.Time
	LastLedgerCloseTime  graphql.Time
	OrderbookStats       orderbookStats
}

// orderbookStats represents the orderbook stats for a
// specific pair of assets (aggregated or not)
type orderbookStats struct {
	BidCount       BigInt
	BidVolume      float64
	BidMax         float64
	AskCount       BigInt
	AskVolume      float64
	AskMin         float64
	Spread         float64
	SpreadMidPoint float64
}

type resolver struct {
	db     *tickerdb.TickerSession
	logger *hlog.Entry
}

// New creates a new GraphQL resolver
func New(s *tickerdb.TickerSession, l *hlog.Entry) *resolver {
	if s == nil {
		panic("A valid database session must be provided for the GraphQL server")
	}
	return &resolver{db: s, logger: l}
}

// Serve creates a GraphQL interface on <address>/graphql and a GraphiQL explorer on /graphiql
func (r *resolver) Serve(address string) {
	relayHandler := r.NewRelayHandler()
	mux := http.NewServeMux()
	mux.Handle("/graphql", http.HandlerFunc(func(wr http.ResponseWriter, re *http.Request) {
		r.logger.Infof("%s %s %s\n", re.RemoteAddr, re.Method, re.URL)
		relayHandler.ServeHTTP(wr, re)
	}))
	mux.Handle("/graphiql", GraphiQL{})

	server := &http.Server{
		Addr:        address,
		Handler:     mux,
		ReadTimeout: 5 * time.Second,
	}
	r.logger.Infof("Starting to serve on address %s\n", address)

	if err := server.ListenAndServe(); err != nil {
		r.logger.Error("server.ListenAndServe:", err)
	}
}

// NewRelayHandler sets up the response handler.
func (r *resolver) NewRelayHandler() relay.Handler {
	opts := []graphql.SchemaOpt{graphql.UseFieldResolvers()}
	r.logger.Info("Validating GraphQL schema")
	s := graphql.MustParseSchema(static.Schema(), r, opts...)
	r.logger.Infof("Schema Validated!")
	relayHandler := relay.Handler{Schema: s}
	return relayHandler
}
