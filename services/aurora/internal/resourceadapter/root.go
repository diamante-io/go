package resourceadapter

import (
	"context"
	"net/url"

	"go/protocols/aurora"
	auroraContext "go/services/aurora/internal/context"
	"go/services/aurora/internal/ledger"
	"go/support/render/hal"
)

// Populate fills in the details
func PopulateRoot(
	ctx context.Context,
	dest *aurora.Root,
	ledgerState ledger.Status,
	hVersion, cVersion string,
	passphrase string,
	currentProtocolVersion int32,
	coreSupportedProtocolVersion int32,
	friendBotURL *url.URL,
	templates map[string]string,
) {
	dest.IngestSequence = ledgerState.ExpHistoryLatest
	dest.auroraSequence = ledgerState.HistoryLatest
	dest.auroraLatestClosedAt = ledgerState.HistoryLatestClosedAt
	dest.HistoryElderSequence = ledgerState.HistoryElder
	dest.CoreSequence = ledgerState.CoreLatest
	dest.auroraVersion = hVersion
	dest.diamcircleCoreVersion = cVersion
	dest.NetworkPassphrase = passphrase
	dest.CurrentProtocolVersion = currentProtocolVersion
	dest.CoreSupportedProtocolVersion = coreSupportedProtocolVersion

	lb := hal.LinkBuilder{Base: auroraContext.BaseURL(ctx)}
	if friendBotURL != nil {
		friendbotLinkBuild := hal.LinkBuilder{Base: friendBotURL}
		l := friendbotLinkBuild.Link("{?addr}")
		dest.Links.Friendbot = &l
	}

	dest.Links.Account = lb.Link("/accounts/{account_id}")
	dest.Links.AccountTransactions = lb.PagedLink("/accounts/{account_id}/transactions")
	dest.Links.Assets = lb.Link("/assets{?asset_code,asset_issuer,cursor,limit,order}")
	dest.Links.Effects = lb.Link("/effects{?cursor,limit,order}")
	dest.Links.Ledger = lb.Link("/ledger/{sequence}")
	dest.Links.Ledgers = lb.Link("/ledgers{?cursor,limit,order}")
	dest.Links.FeeStats = lb.Link("/fee_stats")
	dest.Links.Operation = lb.Link("/operations/{id}")
	dest.Links.Operations = lb.Link("/operations{?cursor,limit,order,include_failed}")
	dest.Links.Payments = lb.Link("/payments{?cursor,limit,order,include_failed}")
	dest.Links.TradeAggregations = lb.Link("/trade_aggregations?base_asset_type={base_asset_type}&base_asset_code={base_asset_code}&base_asset_issuer={base_asset_issuer}&counter_asset_type={counter_asset_type}&counter_asset_code={counter_asset_code}&counter_asset_issuer={counter_asset_issuer}")
	dest.Links.Trades = lb.Link("/trades?base_asset_type={base_asset_type}&base_asset_code={base_asset_code}&base_asset_issuer={base_asset_issuer}&counter_asset_type={counter_asset_type}&counter_asset_code={counter_asset_code}&counter_asset_issuer={counter_asset_issuer}")

	accountsLink := lb.Link(templates["accounts"])
	claimableBalancesLink := lb.Link(templates["claimableBalances"])
	liquidityPoolsLink := lb.Link(templates["liquidityPools"])
	offerLink := lb.Link("/offers/{offer_id}")
	offersLink := lb.Link(templates["offers"])
	strictReceivePaths := lb.Link(templates["strictReceivePaths"])
	strictSendPaths := lb.Link(templates["strictSendPaths"])
	dest.Links.Accounts = &accountsLink
	dest.Links.ClaimableBalances = &claimableBalancesLink
	dest.Links.LiquidityPools = &liquidityPoolsLink
	dest.Links.Offer = &offerLink
	dest.Links.Offers = &offersLink
	dest.Links.StrictReceivePaths = &strictReceivePaths
	dest.Links.StrictSendPaths = &strictSendPaths

	dest.Links.OrderBook = lb.Link("/order_book{?selling_asset_type,selling_asset_code,selling_asset_issuer,buying_asset_type,buying_asset_code,buying_asset_issuer,limit}")
	dest.Links.Self = lb.Link("/")
	dest.Links.Transaction = lb.Link("/transactions/{hash}")
	dest.Links.Transactions = lb.PagedLink("/transactions")
}
