package tickerdb

import (
	"time"

	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
	bdata "github.com/diamcircle/go/services/ticker/internal/tickerdb/migrations"
	"github.com/diamcircle/go/support/db"
)

//go:generate go run github.com/kevinburke/go-bindata/go-bindata@v3.18.0+incompatible -nometadata -ignore .+\.go$ -pkg bdata -o migrations/bindata.go ./...

// TickerSession provides helper methods for making queries against `DB`.
type TickerSession struct {
	db.Session
}

// Asset represents an entry on the assets table
type Asset struct {
	ID                          int32     `db:"id"`
	Code                        string    `db:"code"`
	IssuerAccount               string    `db:"issuer_account"`
	Type                        string    `db:"type"`
	NumAccounts                 int32     `db:"num_accounts"`
	AuthRequired                bool      `db:"auth_required"`
	AuthRevocable               bool      `db:"auth_revocable"`
	Amount                      float64   `db:"amount"`
	AssetControlledByDomain     bool      `db:"asset_controlled_by_domain"`
	AnchorAssetCode             string    `db:"anchor_asset_code"`
	AnchorAssetType             string    `db:"anchor_asset_type"`
	IsValid                     bool      `db:"is_valid"`
	ValidationError             string    `db:"validation_error"`
	LastValid                   time.Time `db:"last_valid"`
	LastChecked                 time.Time `db:"last_checked"`
	DisplayDecimals             int       `db:"display_decimals"`
	Name                        string    `db:"name"`
	Desc                        string    `db:"description"`
	Conditions                  string    `db:"conditions"`
	IsAssetAnchored             bool      `db:"is_asset_anchored"`
	FixedNumber                 int       `db:"fixed_number"`
	MaxNumber                   int       `db:"max_number"`
	IsUnlimited                 bool      `db:"is_unlimited"`
	RedemptionInstructions      string    `db:"redemption_instructions"`
	CollateralAddresses         string    `db:"collateral_addresses"`
	CollateralAddressSignatures string    `db:"collateral_address_signatures"`
	Countries                   string    `db:"countries"`
	Status                      string    `db:"status"`
	IssuerID                    int32     `db:"issuer_id"`
	Issuer                      Issuer    `db:"-"`
}

// Issuer represents an entry on the issuers table
type Issuer struct {
	ID               int32  `db:"id"`
	PublicKey        string `db:"public_key"`
	Name             string `db:"name"`
	URL              string `db:"url"`
	TOMLURL          string `db:"toml_url"`
	FederationServer string `db:"federation_server"`
	AuthServer       string `db:"auth_server"`
	TransferServer   string `db:"transfer_server"`
	WebAuthEndpoint  string `db:"web_auth_endpoint"`
	DepositServer    string `db:"deposit_server"`
	OrgTwitter       string `db:"org_twitter"`
}

// Trade represents an entry on the trades table
type Trade struct {
	ID              int32     `db:"id"`
	AuroraID       string    `db:"aurora_id"`
	LedgerCloseTime time.Time `db:"ledger_close_time"`
	OfferID         string    `db:"offer_id"`
	BaseOfferID     string    `db:"base_offer_id"`
	BaseAccount     string    `db:"base_account"`
	BaseAmount      float64   `db:"base_amount"`
	BaseAssetID     int32     `db:"base_asset_id"`
	CounterOfferID  string    `db:"counter_offer_id"`
	CounterAccount  string    `db:"counter_account"`
	CounterAmount   float64   `db:"counter_amount"`
	CounterAssetID  int32     `db:"counter_asset_id"`
	BaseIsSeller    bool      `db:"base_is_seller"`
	Price           float64   `db:"price"`
}

// OrderbookStats represents an entry on the orderbook_stats table
type OrderbookStats struct {
	ID             int32     `db:"id"`
	BaseAssetID    int32     `db:"base_asset_id"`
	CounterAssetID int32     `db:"counter_asset_id"`
	NumBids        int       `db:"num_bids"`
	BidVolume      float64   `db:"bid_volume"`
	HighestBid     float64   `db:"highest_bid"`
	NumAsks        int       `db:"num_asks"`
	AskVolume      float64   `db:"ask_volume"`
	LowestAsk      float64   `db:"lowest_ask"`
	Spread         float64   `db:"spread"`
	SpreadMidPoint float64   `db:"spread_mid_point"`
	UpdatedAt      time.Time `db:"updated_at"`
}

// Market represent the aggregated market data retrieved from the database.
// Note: this struct does *not* directly map to a db entity.
type Market struct {
	TradePair          string    `db:"trade_pair_name"`
	BaseVolume24h      float64   `db:"base_volume_24h"`
	CounterVolume24h   float64   `db:"counter_volume_24h"`
	TradeCount24h      int64     `db:"trade_count_24h"`
	OpenPrice24h       float64   `db:"open_price_24h"`
	LowestPrice24h     float64   `db:"lowest_price_24h"`
	HighestPrice24h    float64   `db:"highest_price_24h"`
	PriceChange24h     float64   `db:"price_change_24h"`
	BaseVolume7d       float64   `db:"base_volume_7d"`
	CounterVolume7d    float64   `db:"counter_volume_7d"`
	TradeCount7d       int64     `db:"trade_count_7d"`
	OpenPrice7d        float64   `db:"open_price_7d"`
	LowestPrice7d      float64   `db:"lowest_price_7d"`
	HighestPrice7d     float64   `db:"highest_price_7d"`
	PriceChange7d      float64   `db:"price_change_7d"`
	LastPriceCloseTime time.Time `db:"close_time"`
	LastPrice          float64   `db:"last_price"`
	NumBids            int       `db:"num_bids"`
	BidVolume          float64   `db:"bid_volume"`
	HighestBid         float64   `db:"highest_bid"`
	NumAsks            int       `db:"num_asks"`
	AskVolume          float64   `db:"ask_volume"`
	LowestAsk          float64   `db:"lowest_ask"`
}

// PartialMarket represents the aggregated market data for a
// specific pair of assets (or asset codes) during an arbitrary
// time range.
// Note: this struct does *not* directly map to a db entity.
type PartialMarket struct {
	TradePairName        string    `db:"trade_pair_name"`
	BaseAssetID          int32     `db:"base_asset_id"`
	BaseAssetCode        string    `db:"base_asset_code"`
	BaseAssetIssuer      string    `db:"base_asset_issuer"`
	BaseAssetType        string    `db:"base_asset_type"`
	CounterAssetID       int32     `db:"counter_asset_id"`
	CounterAssetCode     string    `db:"counter_asset_code"`
	CounterAssetIssuer   string    `db:"counter_asset_issuer"`
	CounterAssetType     string    `db:"counter_asset_type"`
	BaseVolume           float64   `db:"base_volume"`
	CounterVolume        float64   `db:"counter_volume"`
	TradeCount           int32     `db:"trade_count"`
	Open                 float64   `db:"open_price"`
	Low                  float64   `db:"lowest_price"`
	High                 float64   `db:"highest_price"`
	Change               float64   `db:"price_change"`
	Close                float64   `db:"last_price"`
	NumBids              int       `db:"num_bids"`
	BidVolume            float64   `db:"bid_volume"`
	HighestBid           float64   `db:"highest_bid"`
	NumAsks              int       `db:"num_asks"`
	AskVolume            float64   `db:"ask_volume"`
	LowestAsk            float64   `db:"lowest_ask"`
	IntervalStart        time.Time `db:"interval_start"`
	FirstLedgerCloseTime time.Time `db:"first_ledger_close_time"`
	LastLedgerCloseTime  time.Time `db:"last_ledger_close_time"`
}

// CreateSession returns a new TickerSession that connects to the given db settings
func CreateSession(driverName, dataSourceName string) (session TickerSession, err error) {
	dbconn, err := sqlx.Connect(driverName, dataSourceName)
	if err != nil {
		return
	}

	session.DB = dbconn
	return
}

func MigrateDB(s *TickerSession) (int, error) {
	migrations := &migrate.AssetMigrationSource{
		Asset:    bdata.Asset,
		AssetDir: bdata.AssetDir,
		Dir:      "migrations",
	}
	migrate.SetTable("migrations")
	return migrate.Exec(s.DB.DB, "postgres", migrations, migrate.Up)
}
