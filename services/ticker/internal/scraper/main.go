package scraper

import (
	"context"
	"time"

	auroraclient "github.com/diamcircle/go/clients/auroraclient"
	hProtocol "github.com/diamcircle/go/protocols/aurora"
	"github.com/diamcircle/go/services/ticker/internal/utils"
	hlog "github.com/diamcircle/go/support/log"
)

type ScraperConfig struct {
	Client *auroraclient.Client
	Logger *hlog.Entry
	Ctx    *context.Context
}

// TOMLDoc is the interface for storing TOML Issuer Documentation.
// See: https://github.com/diamcircle/diamcircle-protocol/blob/master/ecosystem/sep-0001.md#currency-documentation
type TOMLDoc struct {
	OrgName    string `toml:"ORG_NAME"`
	OrgURL     string `toml:"ORG_URL"`
	OrgTwitter string `toml:"ORG_TWITTER"`
}

// TOMLCurrency is the interface for storing TOML Currency Information.
// See: https://github.com/diamcircle/diamcircle-protocol/blob/master/ecosystem/sep-0001.md#currency-documentation
type TOMLCurrency struct {
	Code                        string   `toml:"code"`
	Issuer                      string   `toml:"issuer"`
	IsAssetAnchored             bool     `toml:"is_asset_anchored"`
	AnchorAsset                 string   `toml:"anchor_asset"`
	AnchorAssetType             string   `toml:"anchor_asset_type"`
	DisplayDecimals             int      `toml:"display_decimals"`
	Name                        string   `toml:"name"`
	Desc                        string   `toml:"desc"`
	Conditions                  string   `toml:"conditions"`
	FixedNumber                 int      `toml:"fixed_number"`
	MaxNumber                   int      `toml:"max_number"`
	IsUnlimited                 bool     `toml:"is_unlimited"`
	RedemptionInstructions      string   `toml:"redemption_instructions"`
	CollateralAddresses         []string `toml:"collateral_addresses"`
	CollateralAddressSignatures []string `toml:"collateral_address_signatures"`
	Status                      string   `toml:"status"`
}

// TOMLIssuer is the interface for storing TOML Issuer Information.
// See: https://github.com/diamcircle/diamcircle-protocol/blob/master/ecosystem/sep-0001.md#currency-documentation
type TOMLIssuer struct {
	FederationServer string         `toml:"FEDERATION_SERVER"`
	AuthServer       string         `toml:"AUTH_SERVER"`
	TransferServer   string         `toml:"TRANSFER_SERVER"`
	WebAuthEndpoint  string         `toml:"WEB_AUTH_ENDPOINT"`
	SigningKey       string         `toml:"SIGNING_KEY"`
	DepositServer    string         `toml:"DEPOSIT_SERVER"` // for legacy purposes
	Documentation    TOMLDoc        `toml:"DOCUMENTATION"`
	Currencies       []TOMLCurrency `toml:"CURRENCIES"`
	TOMLURL          string         `toml:"-"`
}

// FinalAsset is the interface to represent the aggregated Asset data.
type FinalAsset struct {
	Code                        string     `json:"code"`
	Issuer                      string     `json:"issuer"`
	Type                        string     `json:"type"`
	NumAccounts                 int32      `json:"num_accounts"`
	AuthRequired                bool       `json:"auth_required"`
	AuthRevocable               bool       `json:"auth_revocable"`
	Amount                      float64    `json:"amount"`
	IssuerDetails               TOMLIssuer `json:"-"`
	AssetControlledByDomain     bool       `json:"asset_controlled_by_domain"`
	Error                       string     `json:"-"`
	AnchorAsset                 string     `json:"anchor_asset"`
	AnchorAssetType             string     `json:"anchor_asset_type"`
	IsValid                     bool       `json:"-"`
	LastValid                   time.Time  `json:"-"`
	LastChecked                 time.Time  `json:"-"`
	IsTrash                     bool       `json:"-"`
	DisplayDecimals             int        `json:"display_decimals"`
	Name                        string     `json:"name"`
	Desc                        string     `json:"desc"`
	Conditions                  string     `json:"conditions"`
	IsAssetAnchored             bool       `json:"is_asset_anchored"`
	FixedNumber                 int        `json:"fixed_number"`
	MaxNumber                   int        `json:"max_number"`
	IsUnlimited                 bool       `json:"is_unlimited"`
	RedemptionInstructions      string     `json:"redemption_instructions"`
	CollateralAddresses         []string   `json:"collateral_addresses"`
	CollateralAddressSignatures []string   `json:"collateral_address_signatures"`
	Countries                   string     `json:"countries"`
	Status                      string     `json:"status"`
}

// OrderbookStats represents the Orderbook stats for a given asset
type OrderbookStats struct {
	BaseAssetCode      string
	BaseAssetType      string
	BaseAssetIssuer    string
	CounterAssetCode   string
	CounterAssetType   string
	CounterAssetIssuer string
	NumBids            int
	BidVolume          float64
	HighestBid         float64
	NumAsks            int
	AskVolume          float64
	LowestAsk          float64
	Spread             float64
	SpreadMidPoint     float64
}

// FetchAllAssets fetches assets from the Aurora public net. If limit = 0, will fetch all assets.
func (c *ScraperConfig) FetchAllAssets(limit int, parallelism int) (assets []FinalAsset, err error) {
	dirtyAssets, err := c.retrieveAssets(limit)
	if err != nil {
		return
	}

	assets, numTrash := c.parallelProcessAssets(dirtyAssets, parallelism)

	c.Logger.Infof(
		"Scanned %d entries. Trash: %d. Non-trash: %d\n",
		len(dirtyAssets),
		numTrash,
		len(assets),
	)
	return
}

// FetchAllTrades fetches all trades for a given period, respecting the limit. If limit = 0,
// will fetch all trades for that given period.
func (c *ScraperConfig) FetchAllTrades(since time.Time, limit int) (trades []hProtocol.Trade, err error) {
	c.Logger.Info("Fetching trades from Aurora")

	trades, err = c.retrieveTrades(since, limit)

	c.Logger.Info("Last close time ingested:", trades[len(trades)-1].LedgerCloseTime)
	c.Logger.Infof("Fetched: %d trades\n", len(trades))
	return
}

// StreamNewTrades streams trades directly from aurora and calls the handler function
// whenever a new trade appears.
func (c *ScraperConfig) StreamNewTrades(cursor string, h auroraclient.TradeHandler) error {
	c.Logger.Info("Starting to stream trades with cursor at:", cursor)
	return c.streamTrades(h, cursor)
}

// FetchOrderbookForAssets fetches the orderbook stats for the base and counter assets provided in the parameters
func (c *ScraperConfig) FetchOrderbookForAssets(bType, bCode, bIssuer, cType, cCode, cIssuer string) (OrderbookStats, error) {
	c.Logger.Infof("Fetching orderbook info for %s:%s / %s:%s\n", bCode, bIssuer, cCode, cIssuer)
	return c.fetchOrderbook(bType, bCode, bIssuer, cType, cCode, cIssuer)
}

// NormalizeTradeAssets enforces the following rules:
// 1. native asset type refers to a "XLM" code and a "native" issuer
// 2. native is always the base asset (and if not, base and counter are swapped)
// 3. when trades are between two non-native, the base is the asset whose string
// comes first alphabetically.
func NormalizeTradeAssets(t *hProtocol.Trade) {
	addNativeData(t)
	if t.BaseAssetType == "native" {
		return
	}
	if t.CounterAssetType == "native" {
		reverseAssets(t)
		return
	}
	bAssetString := utils.GetAssetString(t.BaseAssetType, t.BaseAssetCode, t.BaseAssetIssuer)
	cAssetString := utils.GetAssetString(t.CounterAssetType, t.CounterAssetCode, t.CounterAssetIssuer)
	if bAssetString > cAssetString {
		reverseAssets(t)
	}
}
