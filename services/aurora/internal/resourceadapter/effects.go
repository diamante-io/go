package resourceadapter

import (
	"context"

	"go/protocols/aurora/base"
	"go/protocols/aurora/effects"
	auroraContext "go/services/aurora/internal/context"
	"go/services/aurora/internal/db2/history"
	"go/support/render/hal"
	"go/xdr"
)

var EffectTypeNames = map[history.EffectType]string{
	history.EffectAccountCreated:                           "account_created",
	history.EffectAccountRemoved:                           "account_removed",
	history.EffectAccountCredited:                          "account_credited",
	history.EffectAccountDebited:                           "account_debited",
	history.EffectAccountThresholdsUpdated:                 "account_thresholds_updated",
	history.EffectAccountHomeDomainUpdated:                 "account_home_domain_updated",
	history.EffectAccountFlagsUpdated:                      "account_flags_updated",
	history.EffectAccountInflationDestinationUpdated:       "account_inflation_destination_updated",
	history.EffectSignerCreated:                            "signer_created",
	history.EffectSignerRemoved:                            "signer_removed",
	history.EffectSignerUpdated:                            "signer_updated",
	history.EffectTrustlineCreated:                         "trustline_created",
	history.EffectTrustlineRemoved:                         "trustline_removed",
	history.EffectTrustlineUpdated:                         "trustline_updated",
	history.EffectTrustlineAuthorized:                      "trustline_authorized",
	history.EffectTrustlineAuthorizedToMaintainLiabilities: "trustline_authorized_to_maintain_liabilities",
	history.EffectTrustlineDeauthorized:                    "trustline_deauthorized",
	history.EffectTrustlineFlagsUpdated:                    "trustline_flags_updated",
	// unused
	// history.EffectOfferCreated:                             "offer_created",
	// history.EffectOfferRemoved:                             "offer_removed",
	// history.EffectOfferUpdated:                             "offer_updated",
	history.EffectTrade:                              "trade",
	history.EffectDataCreated:                        "data_created",
	history.EffectDataRemoved:                        "data_removed",
	history.EffectDataUpdated:                        "data_updated",
	history.EffectSequenceBumped:                     "sequence_bumped",
	history.EffectClaimableBalanceCreated:            "claimable_balance_created",
	history.EffectClaimableBalanceClaimantCreated:    "claimable_balance_claimant_created",
	history.EffectClaimableBalanceClaimed:            "claimable_balance_claimed",
	history.EffectAccountSponsorshipCreated:          "account_sponsorship_created",
	history.EffectAccountSponsorshipUpdated:          "account_sponsorship_updated",
	history.EffectAccountSponsorshipRemoved:          "account_sponsorship_removed",
	history.EffectTrustlineSponsorshipCreated:        "trustline_sponsorship_created",
	history.EffectTrustlineSponsorshipUpdated:        "trustline_sponsorship_updated",
	history.EffectTrustlineSponsorshipRemoved:        "trustline_sponsorship_removed",
	history.EffectDataSponsorshipCreated:             "data_sponsorship_created",
	history.EffectDataSponsorshipUpdated:             "data_sponsorship_updated",
	history.EffectDataSponsorshipRemoved:             "data_sponsorship_removed",
	history.EffectClaimableBalanceSponsorshipCreated: "claimable_balance_sponsorship_created",
	history.EffectClaimableBalanceSponsorshipUpdated: "claimable_balance_sponsorship_updated",
	history.EffectClaimableBalanceSponsorshipRemoved: "claimable_balance_sponsorship_removed",
	history.EffectSignerSponsorshipCreated:           "signer_sponsorship_created",
	history.EffectSignerSponsorshipUpdated:           "signer_sponsorship_updated",
	history.EffectSignerSponsorshipRemoved:           "signer_sponsorship_removed",
	history.EffectClaimableBalanceClawedBack:         "claimable_balance_clawed_back",
	history.EffectLiquidityPoolDeposited:             "liquidity_pool_deposited",
	history.EffectLiquidityPoolWithdrew:              "liquidity_pool_withdrew",
	history.EffectLiquidityPoolTrade:                 "liquidity_pool_trade",
	history.EffectLiquidityPoolCreated:               "liquidity_pool_created",
	history.EffectLiquidityPoolRemoved:               "liquidity_pool_removed",
	history.EffectLiquidityPoolRevoked:               "liquidity_pool_revoked",
}

// NewEffect creates a new effect resource from the provided database representation
// of the effect.
func NewEffect(
	ctx context.Context,
	row history.Effect,
	ledger history.Ledger,
) (result hal.Pageable, err error) {

	basev := effects.Base{}
	PopulateBaseEffect(ctx, &basev, row, ledger)

	switch row.Type {
	case history.EffectAccountCreated:
		e := effects.AccountCreated{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectAccountCredited:
		e := effects.AccountCredited{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectAccountDebited:
		e := effects.AccountDebited{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectAccountThresholdsUpdated:
		e := effects.AccountThresholdsUpdated{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectAccountHomeDomainUpdated:
		e := effects.AccountHomeDomainUpdated{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectAccountFlagsUpdated:
		e := effects.AccountFlagsUpdated{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectSignerCreated:
		e := effects.SignerCreated{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectSignerUpdated:
		e := effects.SignerUpdated{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectSignerRemoved:
		e := effects.SignerRemoved{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectTrustlineCreated:
		e := effects.TrustlineCreated{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectTrustlineUpdated:
		e := effects.TrustlineUpdated{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectTrustlineRemoved:
		e := effects.TrustlineRemoved{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectTrustlineAuthorized:
		e := effects.TrustlineAuthorized{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectTrustlineAuthorizedToMaintainLiabilities:
		e := effects.TrustlineAuthorizedToMaintainLiabilities{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectTrustlineDeauthorized:
		e := effects.TrustlineDeauthorized{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectTrustlineFlagsUpdated:
		e := effects.TrustlineFlagsUpdated{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectTrade:
		e := effects.Trade{Base: basev}
		tradeDetails := history.TradeEffectDetails{}
		err = row.UnmarshalDetails(&tradeDetails)
		if err == nil {
			e.Seller = tradeDetails.Seller
			e.SellerMuxed = tradeDetails.SellerMuxed
			e.SellerMuxedID = tradeDetails.SellerMuxedID
			e.OfferID = tradeDetails.OfferID
			e.SoldAmount = tradeDetails.SoldAmount
			e.SoldAssetType = tradeDetails.SoldAssetType
			e.SoldAssetCode = tradeDetails.SoldAssetCode
			e.SoldAssetIssuer = tradeDetails.SoldAssetIssuer
			e.BoughtAmount = tradeDetails.BoughtAmount
			e.BoughtAssetType = tradeDetails.BoughtAssetType
			e.BoughtAssetCode = tradeDetails.BoughtAssetCode
			e.BoughtAssetIssuer = tradeDetails.BoughtAssetIssuer
		}
		result = e
	case history.EffectDataCreated:
		e := effects.DataCreated{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectDataUpdated:
		e := effects.DataUpdated{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectDataRemoved:
		e := effects.DataRemoved{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectSequenceBumped:
		e := effects.SequenceBumped{Base: basev}
		hsb := history.SequenceBumped{}
		err = row.UnmarshalDetails(&hsb)
		if err == nil {
			e.NewSeq = hsb.NewSeq
		}
		result = e
	case history.EffectClaimableBalanceCreated:
		e := effects.ClaimableBalanceCreated{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectClaimableBalanceClaimed:
		e := effects.ClaimableBalanceClaimed{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectClaimableBalanceClaimantCreated:
		e := effects.ClaimableBalanceClaimantCreated{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectAccountSponsorshipCreated:
		e := effects.AccountSponsorshipCreated{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectAccountSponsorshipUpdated:
		e := effects.AccountSponsorshipUpdated{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectAccountSponsorshipRemoved:
		e := effects.AccountSponsorshipRemoved{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectTrustlineSponsorshipCreated:
		e := effects.TrustlineSponsorshipCreated{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectTrustlineSponsorshipUpdated:
		e := effects.TrustlineSponsorshipUpdated{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectTrustlineSponsorshipRemoved:
		e := effects.TrustlineSponsorshipRemoved{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectDataSponsorshipCreated:
		e := effects.DataSponsorshipCreated{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectDataSponsorshipUpdated:
		e := effects.DataSponsorshipUpdated{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectDataSponsorshipRemoved:
		e := effects.DataSponsorshipRemoved{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectClaimableBalanceSponsorshipCreated:
		e := effects.ClaimableBalanceSponsorshipCreated{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectClaimableBalanceSponsorshipUpdated:
		e := effects.ClaimableBalanceSponsorshipUpdated{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectClaimableBalanceSponsorshipRemoved:
		e := effects.ClaimableBalanceSponsorshipRemoved{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectSignerSponsorshipCreated:
		e := effects.SignerSponsorshipCreated{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectSignerSponsorshipUpdated:
		e := effects.SignerSponsorshipUpdated{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectSignerSponsorshipRemoved:
		e := effects.SignerSponsorshipRemoved{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectClaimableBalanceClawedBack:
		e := effects.ClaimableBalanceClawedBack{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectLiquidityPoolDeposited:
		e := effects.LiquidityPoolDeposited{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectLiquidityPoolWithdrew:
		e := effects.LiquidityPoolWithdrew{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectLiquidityPoolTrade:
		e := effects.LiquidityPoolTrade{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectLiquidityPoolCreated:
		e := effects.LiquidityPoolCreated{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectLiquidityPoolRemoved:
		e := effects.LiquidityPoolRemoved{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectLiquidityPoolRevoked:
		e := effects.LiquidityPoolRevoked{Base: basev}
		err = row.UnmarshalDetails(&e)
		result = e
	case history.EffectAccountRemoved:
		// there is no explicit data structure for account removed
		fallthrough
	default:
		result = basev
	}

	if err != nil {
		return
	}

	rh, ok := result.(base.Rehydratable)

	if ok {
		err = rh.Rehydrate()
	}

	return
}

// Populate loads this resource from `row`
func PopulateBaseEffect(ctx context.Context, this *effects.Base, row history.Effect, ledger history.Ledger) {
	this.ID = row.ID()
	this.PT = row.PagingToken()
	this.Account = row.Account
	if row.AccountMuxed.Valid {
		this.AccountMuxed = row.AccountMuxed.String
		muxedAccount := xdr.MustMuxedAddress(row.AccountMuxed.String)
		this.AccountMuxedID = uint64(muxedAccount.Med25519.Id)
	}
	populateEffectType(this, row)
	this.LedgerCloseTime = ledger.ClosedAt

	lb := hal.LinkBuilder{auroraContext.BaseURL(ctx)}
	this.Links.Operation = lb.Linkf("/operations/%d", row.HistoryOperationID)
	this.Links.Succeeds = lb.Linkf("/effects?order=desc&cursor=%s", this.PT)
	this.Links.Precedes = lb.Linkf("/effects?order=asc&cursor=%s", this.PT)
}

func populateEffectType(this *effects.Base, row history.Effect) {
	var ok bool
	this.TypeI = int32(row.Type)
	this.Type, ok = EffectTypeNames[row.Type]

	if !ok {
		this.Type = "unknown"
	}
}
