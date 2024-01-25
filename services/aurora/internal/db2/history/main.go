// Package history contains database record definitions useable for
// reading rows from a the history portion of aurora's database
package history

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"go/services/aurora/internal/db2"
	"go/support/db"
	"go/support/errors"
	strtime "go/support/time"
	"go/xdr"
)

const (
	// account effects

	// EffectAccountCreated effects occur when a new account is created
	EffectAccountCreated EffectType = 0 // from create_account

	// EffectAccountRemoved effects occur when one account is merged into another
	EffectAccountRemoved EffectType = 1 // from merge_account

	// EffectAccountCredited effects occur when an account receives some currency
	EffectAccountCredited EffectType = 2 // from create_account, payment, path_payment, merge_account

	// EffectAccountDebited effects occur when an account sends some currency
	EffectAccountDebited EffectType = 3 // from create_account, payment, path_payment, create_account

	// EffectAccountThresholdsUpdated effects occur when an account changes its
	// multisig thresholds.
	EffectAccountThresholdsUpdated EffectType = 4 // from set_options

	// EffectAccountHomeDomainUpdated effects occur when an account changes its
	// home domain.
	EffectAccountHomeDomainUpdated EffectType = 5 // from set_options

	// EffectAccountFlagsUpdated effects occur when an account changes its
	// account flags, either clearing or setting.
	EffectAccountFlagsUpdated EffectType = 6 // from set_options

	// EffectAccountInflationDestinationUpdated effects occur when an account changes its
	// inflation destination.
	EffectAccountInflationDestinationUpdated EffectType = 7 // from set_options

	// signer effects

	// EffectSignerCreated occurs when an account gains a signer
	EffectSignerCreated EffectType = 10 // from set_options

	// EffectSignerRemoved occurs when an account loses a signer
	EffectSignerRemoved EffectType = 11 // from set_options

	// EffectSignerUpdated occurs when an account changes the weight of one of its
	// signers.
	EffectSignerUpdated EffectType = 12 // from set_options

	// trustline effects

	// EffectTrustlineCreated occurs when an account trusts an anchor
	EffectTrustlineCreated EffectType = 20 // from change_trust

	// EffectTrustlineRemoved occurs when an account removes struct by setting the
	// limit of a trustline to 0
	EffectTrustlineRemoved EffectType = 21 // from change_trust

	// EffectTrustlineUpdated occurs when an account changes a trustline's limit
	EffectTrustlineUpdated EffectType = 22 // from change_trust, allow_trust

	// Deprecated: use EffectTrustlineFlagsUpdated instead.
	// EffectTrustlineAuthorized occurs when an anchor has AUTH_REQUIRED flag set
	// to true and it authorizes another account's trustline
	EffectTrustlineAuthorized EffectType = 23 // from allow_trust

	// Deprecated: use EffectTrustlineFlagsUpdated instead.
	// EffectTrustlineDeauthorized occurs when an anchor revokes access to a asset
	// it issues.
	EffectTrustlineDeauthorized EffectType = 24 // from allow_trust

	// Deprecated: use EffectTrustlineFlagsUpdated instead.
	// EffectTrustlineAuthorizedToMaintainLiabilities occurs when an anchor has AUTH_REQUIRED flag set
	// to true and it authorizes another account's trustline to maintain liabilities
	EffectTrustlineAuthorizedToMaintainLiabilities EffectType = 25 // from allow_trust

	// EffectTrustlineFlagsUpdated effects occur when a TrustLine changes its
	// flags, either clearing or setting.
	EffectTrustlineFlagsUpdated EffectType = 26 // from set_trust_line flags

	// trading effects

	// unused
	// EffectOfferCreated occurs when an account offers to trade an asset
	// EffectOfferCreated EffectType = 30 // from manage_offer, creat_passive_offer
	// EffectOfferRemoved occurs when an account removes an offer
	// EffectOfferRemoved EffectType = 31 // from manage_offer, create_passive_offer, path_payment
	// EffectOfferUpdated occurs when an offer is updated by the offering account.
	// EffectOfferUpdated EffectType = 32 // from manage_offer, creat_passive_offer, path_payment

	// EffectTrade occurs when a trade is initiated because of a path payment or
	// offer operation.
	EffectTrade EffectType = 33 // from manage_offer, creat_passive_offer, path_payment

	// data effects

	// EffectDataCreated occurs when an account gets a new data field
	EffectDataCreated EffectType = 40 // from manage_data

	// EffectDataRemoved occurs when an account removes a data field
	EffectDataRemoved EffectType = 41 // from manage_data

	// EffectDataUpdated occurs when an account changes a data field's value
	EffectDataUpdated EffectType = 42 // from manage_data

	// EffectSequenceBumped occurs when an account bumps their sequence number
	EffectSequenceBumped EffectType = 43 // from bump_sequence

	// claimable balance effects

	// EffectClaimableBalanceCreated occurs when a claimable balance is created
	EffectClaimableBalanceCreated EffectType = 50 // from create_claimable_balance

	// EffectClaimableBalanceClaimantCreated occurs when a claimable balance claimant is created
	EffectClaimableBalanceClaimantCreated EffectType = 51 // from create_claimable_balance

	// EffectClaimableBalanceClaimed occurs when a claimable balance is claimed
	EffectClaimableBalanceClaimed EffectType = 52 // from claim_claimable_balance

	// sponsorship effects

	// EffectAccountSponsorshipCreated occurs when an account ledger entry is sponsored
	EffectAccountSponsorshipCreated EffectType = 60 // from create_account

	// EffectAccountSponsorshipUpdated occurs when the sponsoring of an account ledger entry is updated
	EffectAccountSponsorshipUpdated EffectType = 61 // from revoke_sponsorship

	// EffectAccountSponsorshipRemoved occurs when the sponsorship of an account ledger entry is removed
	EffectAccountSponsorshipRemoved EffectType = 62 // from revoke_sponsorship

	// EffectTrustlineSponsorshipCreated occurs when a trustline ledger entry is sponsored
	EffectTrustlineSponsorshipCreated EffectType = 63 // from change_trust

	// EffectTrustlineSponsorshipUpdated occurs when the sponsoring of a trustline ledger entry is updated
	EffectTrustlineSponsorshipUpdated EffectType = 64 // from revoke_sponsorship

	// EffectTrustlineSponsorshipRemoved occurs when the sponsorship of a trustline ledger entry is removed
	EffectTrustlineSponsorshipRemoved EffectType = 65 // from revoke_sponsorship

	// EffectDataSponsorshipCreated occurs when a trustline ledger entry is sponsored
	EffectDataSponsorshipCreated EffectType = 66 // from manage_data

	// EffectDataSponsorshipUpdated occurs when the sponsoring of a trustline ledger entry is updated
	EffectDataSponsorshipUpdated EffectType = 67 // from revoke_sponsorship

	// EffectDataSponsorshipRemoved occurs when the sponsorship of a trustline ledger entry is removed
	EffectDataSponsorshipRemoved EffectType = 68 // from revoke_sponsorship

	// EffectClaimableBalanceSponsorshipCreated occurs when a claimable balance ledger entry is sponsored
	EffectClaimableBalanceSponsorshipCreated EffectType = 69 // from create_claimable_balance

	// EffectClaimableBalanceSponsorshipUpdated occurs when the sponsoring of a claimable balance ledger entry
	// is updated
	EffectClaimableBalanceSponsorshipUpdated EffectType = 70 // from revoke_sponsorship

	// EffectClaimableBalanceSponsorshipRemoved occurs when the sponsorship of a claimable balance ledger entry
	// is removed
	EffectClaimableBalanceSponsorshipRemoved EffectType = 71 // from claim_claimable_balance

	// EffectSignerSponsorshipCreated occurs when the sponsorship of a signer is created
	EffectSignerSponsorshipCreated EffectType = 72 // from set_options

	// EffectSignerSponsorshipUpdated occurs when the sponsorship of a signer is updated
	EffectSignerSponsorshipUpdated EffectType = 73 // from revoke_sponsorship

	// EffectSignerSponsorshipRemoved occurs when the sponsorship of a signer is removed
	EffectSignerSponsorshipRemoved EffectType = 74 // from revoke_sponsorship

	// EffectClaimableBalanceClawedBack occurs when a claimable balance is clawed back
	EffectClaimableBalanceClawedBack EffectType = 80 // from clawback_claimable_balance

	// EffectLiquidityPoolDeposited occurs when a liquidity pool incurs a deposit
	EffectLiquidityPoolDeposited EffectType = 90 // from liquidity_pool_deposit

	// EffectLiquidityPoolWithdrew occurs when a liquidity pool incurs a withdrawal
	EffectLiquidityPoolWithdrew EffectType = 91 // from liquidity_pool_withdraw

	// EffectLiquidityPoolTrade occurs when a trade happens in a liquidity pool
	EffectLiquidityPoolTrade EffectType = 92

	// EffectLiquidityPoolCreated occurs when a liquidity pool is created
	EffectLiquidityPoolCreated EffectType = 93 // from change_trust

	// EffectLiquidityPoolRemoved occurs when a liquidity pool is removed
	EffectLiquidityPoolRemoved EffectType = 94 // from change_trust

	// EffectLiquidityPoolRevoked occurs when a liquidity pool is revoked
	EffectLiquidityPoolRevoked EffectType = 95 // from change_trust_line_flags and allow_trust
)

// Account is a row of data from the `history_accounts` table
type Account struct {
	ID      int64  `db:"id"`
	Address string `db:"address"`
}

// AccountEntry is a row of data from the `account` table
type AccountEntry struct {
	AccountID            string      `db:"account_id"`
	Balance              int64       `db:"balance"`
	BuyingLiabilities    int64       `db:"buying_liabilities"`
	SellingLiabilities   int64       `db:"selling_liabilities"`
	SequenceNumber       int64       `db:"sequence_number"`
	NumSubEntries        uint32      `db:"num_subentries"`
	InflationDestination string      `db:"inflation_destination"`
	HomeDomain           string      `db:"home_domain"`
	Flags                uint32      `db:"flags"`
	MasterWeight         byte        `db:"master_weight"`
	ThresholdLow         byte        `db:"threshold_low"`
	ThresholdMedium      byte        `db:"threshold_medium"`
	ThresholdHigh        byte        `db:"threshold_high"`
	LastModifiedLedger   uint32      `db:"last_modified_ledger"`
	Sponsor              null.String `db:"sponsor"`
	NumSponsored         uint32      `db:"num_sponsored"`
	NumSponsoring        uint32      `db:"num_sponsoring"`
}

type IngestionQ interface {
	QAccounts
	QAssetStats
	QClaimableBalances
	QHistoryClaimableBalances
	QData
	QEffects
	QLedgers
	QLiquidityPools
	QHistoryLiquidityPools
	QOffers
	QOperations
	// QParticipants
	// Copy the small interfaces with shared methods directly, otherwise error:
	// duplicate method CreateAccounts
	NewTransactionParticipantsBatchInsertBuilder(maxBatchSize int) TransactionParticipantsBatchInsertBuilder
	NewOperationParticipantBatchInsertBuilder(maxBatchSize int) OperationParticipantBatchInsertBuilder
	QSigners
	//QTrades
	NewTradeBatchInsertBuilder(maxBatchSize int) TradeBatchInsertBuilder
	RebuildTradeAggregationTimes(ctx context.Context, from, to strtime.Millis) error
	RebuildTradeAggregationBuckets(ctx context.Context, fromLedger, toLedger uint32) error
	CreateAssets(ctx context.Context, assets []xdr.Asset, batchSize int) (map[string]Asset, error)
	QTransactions
	QTrustLines

	Begin() error
	BeginTx(*sql.TxOptions) error
	Commit() error
	CloneIngestionQ() IngestionQ
	Rollback() error
	GetTx() *sqlx.Tx
	GetIngestVersion(context.Context) (int, error)
	UpdateExpStateInvalid(context.Context, bool) error
	UpdateIngestVersion(context.Context, int) error
	GetExpStateInvalid(context.Context) (bool, error)
	GetLatestHistoryLedger(context.Context) (uint32, error)
	GetOfferCompactionSequence(context.Context) (uint32, error)
	GetLiquidityPoolCompactionSequence(context.Context) (uint32, error)
	TruncateIngestStateTables(context.Context) error
	DeleteRangeAll(ctx context.Context, start, end int64) error
}

// QAccounts defines account related queries.
type QAccounts interface {
	GetAccountsByIDs(ctx context.Context, ids []string) ([]AccountEntry, error)
	UpsertAccounts(ctx context.Context, accounts []AccountEntry) error
	RemoveAccounts(ctx context.Context, accountIDs []string) (int64, error)
}

// AccountSigner is a row of data from the `accounts_signers` table
type AccountSigner struct {
	Account string      `db:"account_id"`
	Signer  string      `db:"signer"`
	Weight  int32       `db:"weight"`
	Sponsor null.String `db:"sponsor"`
}

type AccountSignersBatchInsertBuilder interface {
	Add(ctx context.Context, signer AccountSigner) error
	Exec(ctx context.Context) error
}

// accountSignersBatchInsertBuilder is a simple wrapper around db.BatchInsertBuilder
type accountSignersBatchInsertBuilder struct {
	builder db.BatchInsertBuilder
}

// Data is a row of data from the `account_data` table
type Data struct {
	AccountID          string           `db:"account_id"`
	Name               string           `db:"name"`
	Value              AccountDataValue `db:"value"`
	LastModifiedLedger uint32           `db:"last_modified_ledger"`
	Sponsor            null.String      `db:"sponsor"`
}

type AccountDataValue []byte

type AccountDataKey struct {
	AccountID string
	DataName  string
}

// QData defines account data related queries.
type QData interface {
	CountAccountsData(ctx context.Context) (int, error)
	GetAccountDataByKeys(ctx context.Context, keys []AccountDataKey) ([]Data, error)
	UpsertAccountData(ctx context.Context, data []Data) error
	RemoveAccountData(ctx context.Context, keys []AccountDataKey) (int64, error)
}

// Asset is a row of data from the `history_assets` table
type Asset struct {
	ID     int64  `db:"id"`
	Type   string `db:"asset_type"`
	Code   string `db:"asset_code"`
	Issuer string `db:"asset_issuer"`
}

// ExpAssetStat is a row in the exp_asset_stats table representing the stats per Asset
type ExpAssetStat struct {
	AssetType   xdr.AssetType        `db:"asset_type"`
	AssetCode   string               `db:"asset_code"`
	AssetIssuer string               `db:"asset_issuer"`
	Accounts    ExpAssetStatAccounts `db:"accounts"`
	Balances    ExpAssetStatBalances `db:"balances"`
	Amount      string               `db:"amount"`
	NumAccounts int32                `db:"num_accounts"`
}

// PagingToken returns a cursor for this asset stat
func (e ExpAssetStat) PagingToken() string {
	return fmt.Sprintf(
		"%s_%s_%s",
		e.AssetCode,
		e.AssetIssuer,
		xdr.AssetTypeToString[e.AssetType],
	)
}

// ExpAssetStatAccounts represents the summarized acount numbers for a single Asset
type ExpAssetStatAccounts struct {
	Authorized                      int32 `json:"authorized"`
	AuthorizedToMaintainLiabilities int32 `json:"authorized_to_maintain_liabilities"`
	ClaimableBalances               int32 `json:"claimable_balances"`
	LiquidityPools                  int32 `json:"liquidity_pools"`
	Unauthorized                    int32 `json:"unauthorized"`
}

func (e ExpAssetStatAccounts) Value() (driver.Value, error) {
	return json.Marshal(e)
}

func (e *ExpAssetStatAccounts) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	return json.Unmarshal(source, &e)
}

func (a ExpAssetStatAccounts) Add(b ExpAssetStatAccounts) ExpAssetStatAccounts {
	return ExpAssetStatAccounts{
		Authorized:                      a.Authorized + b.Authorized,
		AuthorizedToMaintainLiabilities: a.AuthorizedToMaintainLiabilities + b.AuthorizedToMaintainLiabilities,
		ClaimableBalances:               a.ClaimableBalances + b.ClaimableBalances,
		LiquidityPools:                  a.LiquidityPools + b.LiquidityPools,
		Unauthorized:                    a.Unauthorized + b.Unauthorized,
	}
}

func (a ExpAssetStatAccounts) IsZero() bool {
	return a == ExpAssetStatAccounts{}
}

// ExpAssetStatBalances represents the summarized balances for a single Asset
// Note: the string representation is in stroops!
type ExpAssetStatBalances struct {
	Authorized                      string `json:"authorized"`
	AuthorizedToMaintainLiabilities string `json:"authorized_to_maintain_liabilities"`
	ClaimableBalances               string `json:"claimable_balances"`
	LiquidityPools                  string `json:"liquidity_pools"`
	Unauthorized                    string `json:"unauthorized"`
}

func (e ExpAssetStatBalances) Value() (driver.Value, error) {
	return json.Marshal(e)
}

func (e *ExpAssetStatBalances) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	err := json.Unmarshal(source, &e)
	if err != nil {
		return err
	}

	// Sets zero values for empty balances
	if e.Authorized == "" {
		e.Authorized = "0"
	}
	if e.AuthorizedToMaintainLiabilities == "" {
		e.AuthorizedToMaintainLiabilities = "0"
	}
	if e.ClaimableBalances == "" {
		e.ClaimableBalances = "0"
	}
	if e.LiquidityPools == "" {
		e.LiquidityPools = "0"
	}
	if e.Unauthorized == "" {
		e.Unauthorized = "0"
	}

	return nil
}

// QAssetStats defines exp_asset_stats related queries.
type QAssetStats interface {
	InsertAssetStats(ctx context.Context, stats []ExpAssetStat, batchSize int) error
	InsertAssetStat(ctx context.Context, stat ExpAssetStat) (int64, error)
	UpdateAssetStat(ctx context.Context, stat ExpAssetStat) (int64, error)
	GetAssetStat(ctx context.Context, assetType xdr.AssetType, assetCode, assetIssuer string) (ExpAssetStat, error)
	RemoveAssetStat(ctx context.Context, assetType xdr.AssetType, assetCode, assetIssuer string) (int64, error)
	GetAssetStats(ctx context.Context, assetCode, assetIssuer string, page db2.PageQuery) ([]ExpAssetStat, error)
	CountTrustLines(ctx context.Context) (int, error)
}

type QCreateAccountsHistory interface {
	CreateAccounts(ctx context.Context, addresses []string, maxBatchSize int) (map[string]int64, error)
}

// Effect is a row of data from the `history_effects` table
type Effect struct {
	HistoryAccountID   int64       `db:"history_account_id"`
	Account            string      `db:"address"`
	AccountMuxed       null.String `db:"address_muxed"`
	HistoryOperationID int64       `db:"history_operation_id"`
	Order              int32       `db:"order"`
	Type               EffectType  `db:"type"`
	DetailsString      null.String `db:"details"`
}

// TradeEffectDetails is a struct of data from `effects.DetailsString`
// when the effect type is trade
type TradeEffectDetails struct {
	Seller            string `json:"seller"`
	SellerMuxed       string `json:"seller_muxed,omitempty"`
	SellerMuxedID     uint64 `json:"seller_muxed_id,omitempty"`
	OfferID           int64  `json:"offer_id"`
	SoldAmount        string `json:"sold_amount"`
	SoldAssetType     string `json:"sold_asset_type"`
	SoldAssetCode     string `json:"sold_asset_code,omitempty"`
	SoldAssetIssuer   string `json:"sold_asset_issuer,omitempty"`
	BoughtAmount      string `json:"bought_amount"`
	BoughtAssetType   string `json:"bought_asset_type"`
	BoughtAssetCode   string `json:"bought_asset_code,omitempty"`
	BoughtAssetIssuer string `json:"bought_asset_issuer,omitempty"`
}

// SequenceBumped is a struct of data from `effects.DetailsString`
// when the effect type is sequence bumped.
type SequenceBumped struct {
	NewSeq int64 `json:"new_seq"`
}

// EffectsQ is a helper struct to aid in configuring queries that loads
// slices of Ledger structs.
type EffectsQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

// EffectType is the numeric type for an effect, used as the `type` field in the
// `history_effects` table.
type EffectType int

// FeeStats is a row of data from the min, mode, percentile aggregate functions over the
// `history_transactions` table.
type FeeStats struct {
	FeeChargedMax  null.Int `db:"fee_charged_max"`
	FeeChargedMin  null.Int `db:"fee_charged_min"`
	FeeChargedMode null.Int `db:"fee_charged_mode"`
	FeeChargedP10  null.Int `db:"fee_charged_p10"`
	FeeChargedP20  null.Int `db:"fee_charged_p20"`
	FeeChargedP30  null.Int `db:"fee_charged_p30"`
	FeeChargedP40  null.Int `db:"fee_charged_p40"`
	FeeChargedP50  null.Int `db:"fee_charged_p50"`
	FeeChargedP60  null.Int `db:"fee_charged_p60"`
	FeeChargedP70  null.Int `db:"fee_charged_p70"`
	FeeChargedP80  null.Int `db:"fee_charged_p80"`
	FeeChargedP90  null.Int `db:"fee_charged_p90"`
	FeeChargedP95  null.Int `db:"fee_charged_p95"`
	FeeChargedP99  null.Int `db:"fee_charged_p99"`
	MaxFeeMax      null.Int `db:"max_fee_max"`
	MaxFeeMin      null.Int `db:"max_fee_min"`
	MaxFeeMode     null.Int `db:"max_fee_mode"`
	MaxFeeP10      null.Int `db:"max_fee_p10"`
	MaxFeeP20      null.Int `db:"max_fee_p20"`
	MaxFeeP30      null.Int `db:"max_fee_p30"`
	MaxFeeP40      null.Int `db:"max_fee_p40"`
	MaxFeeP50      null.Int `db:"max_fee_p50"`
	MaxFeeP60      null.Int `db:"max_fee_p60"`
	MaxFeeP70      null.Int `db:"max_fee_p70"`
	MaxFeeP80      null.Int `db:"max_fee_p80"`
	MaxFeeP90      null.Int `db:"max_fee_p90"`
	MaxFeeP95      null.Int `db:"max_fee_p95"`
	MaxFeeP99      null.Int `db:"max_fee_p99"`
}

// LatestLedger represents a response from the raw LatestLedgerBaseFeeAndSequence
// query.
type LatestLedger struct {
	BaseFee  int32 `db:"base_fee"`
	Sequence int32 `db:"sequence"`
}

// Ledger is a row of data from the `history_ledgers` table
type Ledger struct {
	TotalOrderID
	Sequence                   int32       `db:"sequence"`
	ImporterVersion            int32       `db:"importer_version"`
	LedgerHash                 string      `db:"ledger_hash"`
	PreviousLedgerHash         null.String `db:"previous_ledger_hash"`
	TransactionCount           int32       `db:"transaction_count"`
	SuccessfulTransactionCount *int32      `db:"successful_transaction_count"`
	FailedTransactionCount     *int32      `db:"failed_transaction_count"`
	OperationCount             int32       `db:"operation_count"`
	TxSetOperationCount        *int32      `db:"tx_set_operation_count"`
	ClosedAt                   time.Time   `db:"closed_at"`
	CreatedAt                  time.Time   `db:"created_at"`
	UpdatedAt                  time.Time   `db:"updated_at"`
	TotalCoins                 int64       `db:"total_coins"`
	FeePool                    int64       `db:"fee_pool"`
	BaseFee                    int32       `db:"base_fee"`
	BaseReserve                int32       `db:"base_reserve"`
	MaxTxSetSize               int32       `db:"max_tx_set_size"`
	ProtocolVersion            int32       `db:"protocol_version"`
	LedgerHeaderXDR            null.String `db:"ledger_header"`
}

// LedgerCapacityUsageStats contains ledgers fullness stats.
type LedgerCapacityUsageStats struct {
	CapacityUsage null.String `db:"ledger_capacity_usage"`
}

// LedgerCache is a helper struct to load ledger data related to a batch of
// sequences.
type LedgerCache struct {
	Records map[int32]Ledger

	lock   sync.Mutex
	queued map[int32]struct{}
}

type LedgerRange struct {
	StartSequence uint32 `db:"start"`
	EndSequence   uint32 `db:"end"`
}

// LedgersQ is a helper struct to aid in configuring queries that loads
// slices of Ledger structs.
type LedgersQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

// Operation is a row of data from the `history_operations` table
type Operation struct {
	TotalOrderID
	TransactionID         int64             `db:"transaction_id"`
	TransactionHash       string            `db:"transaction_hash"`
	TxResult              string            `db:"tx_result"`
	ApplicationOrder      int32             `db:"application_order"`
	Type                  xdr.OperationType `db:"type"`
	DetailsString         null.String       `db:"details"`
	SourceAccount         string            `db:"source_account"`
	SourceAccountMuxed    null.String       `db:"source_account_muxed"`
	TransactionSuccessful bool              `db:"transaction_successful"`
}

// ManageOffer is a struct of data from `operations.DetailsString`
// when the operation type is manage sell offer or manage buy offer
type ManageOffer struct {
	OfferID int64 `json:"offer_id"`
}

// upsertField is used in upsertRows function generating upsert query for
// different tables.
type upsertField struct {
	name    string
	dbType  string
	objects []interface{}
}

// Offer is row of data from the `offers` table from aurora DB
type Offer struct {
	SellerID string `db:"seller_id"`
	OfferID  int64  `db:"offer_id"`

	SellingAsset xdr.Asset `db:"selling_asset"`
	BuyingAsset  xdr.Asset `db:"buying_asset"`

	Amount             int64       `db:"amount"`
	Pricen             int32       `db:"pricen"`
	Priced             int32       `db:"priced"`
	Price              float64     `db:"price"`
	Flags              int32       `db:"flags"`
	Deleted            bool        `db:"deleted"`
	LastModifiedLedger uint32      `db:"last_modified_ledger"`
	Sponsor            null.String `db:"sponsor"`
}

// OperationsQ is a helper struct to aid in configuring queries that loads
// slices of Operation structs.
type OperationsQ struct {
	Err                 error
	parent              *Q
	sql                 sq.SelectBuilder
	opIdCol             string
	includeFailed       bool
	includeTransactions bool
}

// Q is a helper struct on which to hang common_trades queries against a history
// portion of the aurora database.
type Q struct {
	db.SessionInterface
}

// QSigners defines signer related queries.
type QSigners interface {
	GetLastLedgerIngestNonBlocking(ctx context.Context) (uint32, error)
	GetLastLedgerIngest(ctx context.Context) (uint32, error)
	UpdateLastLedgerIngest(ctx context.Context, ledgerSequence uint32) error
	AccountsForSigner(ctx context.Context, signer string, page db2.PageQuery) ([]AccountSigner, error)
	NewAccountSignersBatchInsertBuilder(maxBatchSize int) AccountSignersBatchInsertBuilder
	CreateAccountSigner(ctx context.Context, account, signer string, weight int32, sponsor *string) (int64, error)
	RemoveAccountSigner(ctx context.Context, account, signer string) (int64, error)
	SignersForAccounts(ctx context.Context, accounts []string) ([]AccountSigner, error)
	CountAccounts(ctx context.Context) (int, error)
}

// OffersQuery is a helper struct to configure queries to offers
type OffersQuery struct {
	PageQuery db2.PageQuery
	SellerID  string
	Sponsor   string
	Selling   *xdr.Asset
	Buying    *xdr.Asset
}

// TotalOrderID represents the ID portion of rows that are identified by the
// "TotalOrderID".  See total_order_id.go in the `db` package for details.
type TotalOrderID struct {
	ID int64 `db:"id"`
}

// Trade represents a trade from the trades table, joined with asset information from the assets table and account
// addresses from the accounts table
type Trade struct {
	HistoryOperationID     int64       `db:"history_operation_id"`
	Order                  int32       `db:"order"`
	LedgerCloseTime        time.Time   `db:"ledger_closed_at"`
	BaseOfferID            null.Int    `db:"base_offer_id"`
	BaseAccount            null.String `db:"base_account"`
	BaseAssetType          string      `db:"base_asset_type"`
	BaseAssetCode          string      `db:"base_asset_code"`
	BaseAssetIssuer        string      `db:"base_asset_issuer"`
	BaseAmount             int64       `db:"base_amount"`
	BaseLiquidityPoolID    null.String `db:"base_liquidity_pool_id"`
	CounterOfferID         null.Int    `db:"counter_offer_id"`
	CounterAccount         null.String `db:"counter_account"`
	CounterAssetType       string      `db:"counter_asset_type"`
	CounterAssetCode       string      `db:"counter_asset_code"`
	CounterAssetIssuer     string      `db:"counter_asset_issuer"`
	CounterAmount          int64       `db:"counter_amount"`
	CounterLiquidityPoolID null.String `db:"counter_liquidity_pool_id"`
	LiquidityPoolFee       null.Int    `db:"liquidity_pool_fee"`
	BaseIsSeller           bool        `db:"base_is_seller"`
	PriceN                 null.Int    `db:"price_n"`
	PriceD                 null.Int    `db:"price_d"`
	Type                   TradeType   `db:"trade_type"`
}

// Transaction is a row of data from the `history_transactions` table
type Transaction struct {
	LedgerCloseTime time.Time `db:"ledger_close_time"`
	TransactionWithoutLedger
}

// TransactionsQ is a helper struct to aid in configuring queries that loads
// slices of transaction structs.
type TransactionsQ struct {
	Err           error
	parent        *Q
	sql           sq.SelectBuilder
	includeFailed bool
}

// TrustLine is row of data from the `trust_lines` table from aurora DB
type TrustLine struct {
	AccountID          string        `db:"account_id"`
	AssetType          xdr.AssetType `db:"asset_type"`
	AssetIssuer        string        `db:"asset_issuer"`
	AssetCode          string        `db:"asset_code"`
	Balance            int64         `db:"balance"`
	LedgerKey          string        `db:"ledger_key"`
	Limit              int64         `db:"trust_line_limit"`
	LiquidityPoolID    string        `db:"liquidity_pool_id"`
	BuyingLiabilities  int64         `db:"buying_liabilities"`
	SellingLiabilities int64         `db:"selling_liabilities"`
	Flags              uint32        `db:"flags"`
	LastModifiedLedger uint32        `db:"last_modified_ledger"`
	Sponsor            null.String   `db:"sponsor"`
}

// QTrustLines defines trust lines related queries.
type QTrustLines interface {
	GetTrustLinesByKeys(ctx context.Context, ledgerKeys []string) ([]TrustLine, error)
	UpsertTrustLines(ctx context.Context, trustlines []TrustLine) error
	RemoveTrustLines(ctx context.Context, ledgerKeys []string) (int64, error)
}

func (q *Q) NewAccountSignersBatchInsertBuilder(maxBatchSize int) AccountSignersBatchInsertBuilder {
	return &accountSignersBatchInsertBuilder{
		builder: db.BatchInsertBuilder{
			Table:        q.GetTable("accounts_signers"),
			MaxBatchSize: maxBatchSize,
		},
	}
}

// ElderLedger loads the oldest ledger known to the history database
func (q *Q) ElderLedger(ctx context.Context, dest interface{}) error {
	return q.GetRaw(ctx, dest, `SELECT COALESCE(MIN(sequence), 0) FROM history_ledgers`)
}

// GetLatestHistoryLedger loads the latest known ledger. Returns 0 if no ledgers in
// `history_ledgers` table.
func (q *Q) GetLatestHistoryLedger(ctx context.Context) (uint32, error) {
	var value uint32
	err := q.LatestLedger(ctx, &value)
	return value, err
}

// LatestLedger loads the latest known ledger
func (q *Q) LatestLedger(ctx context.Context, dest interface{}) error {
	return q.GetRaw(ctx, dest, `SELECT COALESCE(MAX(sequence), 0) FROM history_ledgers`)
}

// LatestLedgerSequenceClosedAt loads the latest known ledger sequence and close time,
// returns empty values if no ledgers in a DB.
func (q *Q) LatestLedgerSequenceClosedAt(ctx context.Context) (int32, time.Time, error) {
	ledger := struct {
		Sequence int32     `db:"sequence"`
		ClosedAt time.Time `db:"closed_at"`
	}{}
	err := q.GetRaw(ctx, &ledger, `SELECT sequence, closed_at FROM history_ledgers ORDER BY sequence DESC LIMIT 1`)
	if err == sql.ErrNoRows {
		// Will return empty values
		return ledger.Sequence, ledger.ClosedAt, nil
	}
	return ledger.Sequence, ledger.ClosedAt, err
}

// LatestLedgerBaseFeeAndSequence loads the latest known ledger's base fee and
// sequence number.
func (q *Q) LatestLedgerBaseFeeAndSequence(ctx context.Context, dest interface{}) error {
	return q.GetRaw(ctx, dest, `
		SELECT base_fee, sequence
		FROM history_ledgers
		WHERE sequence = (SELECT COALESCE(MAX(sequence), 0) FROM history_ledgers)
	`)
}

// CloneIngestionQ clones underlying db.Session and returns IngestionQ
func (q *Q) CloneIngestionQ() IngestionQ {
	return &Q{q.Clone()}
}

// DeleteRangeAll deletes a range of rows from all history tables between
// `start` and `end` (exclusive).
func (q *Q) DeleteRangeAll(ctx context.Context, start, end int64) error {
	for table, column := range map[string]string{
		"history_effects":                        "history_operation_id",
		"history_ledgers":                        "id",
		"history_operation_claimable_balances":   "history_operation_id",
		"history_operation_participants":         "history_operation_id",
		"history_operation_liquidity_pools":      "history_operation_id",
		"history_operations":                     "id",
		"history_trades":                         "history_operation_id",
		"history_trades_60000":                   "open_ledger_toid",
		"history_transaction_claimable_balances": "history_transaction_id",
		"history_transaction_participants":       "history_transaction_id",
		"history_transaction_liquidity_pools":    "history_transaction_id",
		"history_transactions":                   "id",
	} {
		err := q.DeleteRange(ctx, start, end, table, column)
		if err != nil {
			return errors.Wrapf(err, "Error clearing %s", table)
		}
	}
	return nil
}

// upsertRows builds and executes an upsert query that allows very fast upserts
// to a given table. The final query is of form:
//
// WITH r AS
// 		(SELECT
//			/* unnestPart */
// 			unnest(?::type1[]), /* field1 */
// 			unnest(?::type2[]), /* field2 */
//			...
// 		)
// 	INSERT INTO table (
//		/* insertFieldsPart */
// 		field1,
// 		field2,
//		...
// 	)
// 	SELECT * from r
// 	ON CONFLICT (conflictField) DO UPDATE SET
//		/* onConflictPart */
// 		field1 = excluded.field1,
// 		field2 = excluded.field2,
//		...
func (q *Q) upsertRows(ctx context.Context, table string, conflictField string, fields []upsertField) error {
	unnestPart := make([]string, 0, len(fields))
	insertFieldsPart := make([]string, 0, len(fields))
	onConflictPart := make([]string, 0, len(fields))
	pqArrays := make([]interface{}, 0, len(fields))

	for _, field := range fields {
		unnestPart = append(
			unnestPart,
			fmt.Sprintf("unnest(?::%s[]) /* %s */", field.dbType, field.name),
		)
		insertFieldsPart = append(
			insertFieldsPart,
			field.name,
		)
		onConflictPart = append(
			onConflictPart,
			fmt.Sprintf("%s = excluded.%s", field.name, field.name),
		)
		pqArrays = append(
			pqArrays,
			pq.Array(field.objects),
		)
	}

	sql := `
	WITH r AS
		(SELECT ` + strings.Join(unnestPart, ",") + `)
	INSERT INTO ` + table + `
		(` + strings.Join(insertFieldsPart, ",") + `)
	SELECT * from r
	ON CONFLICT (` + conflictField + `) DO UPDATE SET
		` + strings.Join(onConflictPart, ",")

	_, err := q.ExecRaw(
		context.WithValue(ctx, &db.QueryTypeContextKey, db.UpsertQueryType),
		sql,
		pqArrays...,
	)
	return err
}
