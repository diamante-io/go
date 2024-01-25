package processors

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"

	"github.com/guregu/null"
	"github.com/diamcircle/go/amount"
	"github.com/diamcircle/go/ingest"
	"github.com/diamcircle/go/keypair"
	"github.com/diamcircle/go/protocols/aurora/base"
	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/xdr"
)

// EffectProcessor process effects
type EffectProcessor struct {
	effects  []effect
	effectsQ history.QEffects
	sequence uint32
}

func NewEffectProcessor(effectsQ history.QEffects, sequence uint32) *EffectProcessor {
	return &EffectProcessor{
		effectsQ: effectsQ,
		sequence: sequence,
	}
}

func (p *EffectProcessor) loadAccountIDs(ctx context.Context, accountSet map[string]int64) error {
	addresses := make([]string, 0, len(accountSet))
	for address := range accountSet {
		addresses = append(addresses, address)
	}

	addressToID, err := p.effectsQ.CreateAccounts(ctx, addresses, maxBatchSize)
	if err != nil {
		return errors.Wrap(err, "Could not create account ids")
	}

	for _, address := range addresses {
		id, ok := addressToID[address]
		if !ok {
			return errors.Errorf("no id found for account address %s", address)
		}

		accountSet[address] = id
	}

	return nil
}

func operationsEffects(transaction ingest.LedgerTransaction, sequence uint32) ([]effect, error) {
	effects := []effect{}

	for opi, op := range transaction.Envelope.Operations() {
		operation := transactionOperationWrapper{
			index:          uint32(opi),
			transaction:    transaction,
			operation:      op,
			ledgerSequence: sequence,
		}

		p, err := operation.effects()
		if err != nil {
			return effects, errors.Wrapf(err, "reading operation %v effects", operation.ID())
		}
		effects = append(effects, p...)
	}

	return effects, nil
}

func (p *EffectProcessor) insertDBOperationsEffects(ctx context.Context, effects []effect, accountSet map[string]int64) error {
	batch := p.effectsQ.NewEffectBatchInsertBuilder(maxBatchSize)

	for _, effect := range effects {
		accountID, found := accountSet[effect.address]

		if !found {
			return errors.Errorf("Error finding history_account_id for address %v", effect.address)
		}

		var detailsJSON []byte
		detailsJSON, err := json.Marshal(effect.details)

		if err != nil {
			return errors.Wrapf(err, "Error marshaling details for operation effect %v", effect.operationID)
		}

		if err := batch.Add(ctx,
			accountID,
			effect.addressMuxed,
			effect.operationID,
			effect.order,
			effect.effectType,
			detailsJSON,
		); err != nil {
			return errors.Wrap(err, "could not insert operation effect in db")
		}
	}

	if err := batch.Exec(ctx); err != nil {
		return errors.Wrap(err, "could not flush operation effects to db")
	}
	return nil
}

func (p *EffectProcessor) ProcessTransaction(ctx context.Context, transaction ingest.LedgerTransaction) (err error) {
	// Failed transactions don't have operation effects
	if !transaction.Result.Successful() {
		return nil
	}

	var effectsForTx []effect
	effectsForTx, err = operationsEffects(transaction, p.sequence)
	if err != nil {
		return err
	}
	p.effects = append(p.effects, effectsForTx...)

	return nil
}

func (p *EffectProcessor) Commit(ctx context.Context) (err error) {
	if len(p.effects) > 0 {
		accountSet := map[string]int64{}

		for _, effect := range p.effects {
			accountSet[effect.address] = 0
		}

		if err = p.loadAccountIDs(ctx, accountSet); err != nil {
			return err
		}

		if err = p.insertDBOperationsEffects(ctx, p.effects, accountSet); err != nil {
			return err
		}
	}

	return err
}

type effect struct {
	address      string
	addressMuxed null.String
	operationID  int64
	details      map[string]interface{}
	effectType   history.EffectType
	order        uint32
}

// Effects returns the operation effects
func (operation *transactionOperationWrapper) effects() ([]effect, error) {
	if !operation.transaction.Result.Successful() {
		return []effect{}, nil
	}
	var (
		op  = operation.operation
		err error
	)

	changes, err := operation.transaction.GetOperationChanges(operation.index)
	if err != nil {
		return nil, err
	}

	wrapper := &effectsWrapper{
		effects:   []effect{},
		operation: operation,
	}

	switch operation.OperationType() {
	case xdr.OperationTypeCreateAccount:
		wrapper.addAccountCreatedEffects()
	case xdr.OperationTypePayment:
		wrapper.addPaymentEffects()
	case xdr.OperationTypePathPaymentStrictReceive:
		err = wrapper.pathPaymentStrictReceiveEffects()
	case xdr.OperationTypePathPaymentStrictSend:
		err = wrapper.addPathPaymentStrictSendEffects()
	case xdr.OperationTypeManageSellOffer:
		err = wrapper.addManageSellOfferEffects()
	case xdr.OperationTypeManageBuyOffer:
		err = wrapper.addManageBuyOfferEffects()
	case xdr.OperationTypeCreatePassiveSellOffer:
		err = wrapper.addCreatePassiveSellOfferEffect()
	case xdr.OperationTypeSetOptions:
		wrapper.addSetOptionsEffects()
	case xdr.OperationTypeChangeTrust:
		err = wrapper.addChangeTrustEffects()
	case xdr.OperationTypeAllowTrust:
		err = wrapper.addAllowTrustEffects()
	case xdr.OperationTypeAccountMerge:
		wrapper.addAccountMergeEffects()
	case xdr.OperationTypeInflation:
		wrapper.addInflationEffects()
	case xdr.OperationTypeManageData:
		err = wrapper.addManageDataEffects()
	case xdr.OperationTypeBumpSequence:
		err = wrapper.addBumpSequenceEffects()
	case xdr.OperationTypeCreateClaimableBalance:
		err = wrapper.addCreateClaimableBalanceEffects(changes)
	case xdr.OperationTypeClaimClaimableBalance:
		err = wrapper.addClaimClaimableBalanceEffects(changes)
	case xdr.OperationTypeBeginSponsoringFutureReserves, xdr.OperationTypeEndSponsoringFutureReserves, xdr.OperationTypeRevokeSponsorship:
	// The effects of these operations are obtained  indirectly from the ledger entries
	case xdr.OperationTypeClawback:
		err = wrapper.addClawbackEffects()
	case xdr.OperationTypeClawbackClaimableBalance:
		err = wrapper.addClawbackClaimableBalanceEffects(changes)
	case xdr.OperationTypeSetTrustLineFlags:
		err = wrapper.addSetTrustLineFlagsEffects()
	case xdr.OperationTypeLiquidityPoolDeposit:
		err = wrapper.addLiquidityPoolDepositEffect()
	case xdr.OperationTypeLiquidityPoolWithdraw:
		err = wrapper.addLiquidityPoolWithdrawEffect()
	default:
		return nil, fmt.Errorf("Unknown operation type: %s", op.Body.Type)
	}
	if err != nil {
		return nil, err
	}

	// Effects generated for multiple operations. Keep the effect categories
	// separated so they are "together" in case of different order or meta
	// changes generate by core (unordered_map).

	// Sponsorships
	for _, change := range changes {
		if err = wrapper.addLedgerEntrySponsorshipEffects(change); err != nil {
			return nil, err
		}
		wrapper.addSignerSponsorshipEffects(change)
	}

	// Liquidity pools
	for _, change := range changes {
		// Effects caused by ChangeTrust (creation), AllowTrust and SetTrustlineFlags (removal through revocation)
		wrapper.addLedgerEntryLiquidityPoolEffects(change)
	}

	return wrapper.effects, nil
}

type effectsWrapper struct {
	effects   []effect
	operation *transactionOperationWrapper
}

func (e *effectsWrapper) add(address string, addressMuxed null.String, effectType history.EffectType, details map[string]interface{}) {
	e.effects = append(e.effects, effect{
		address:      address,
		addressMuxed: addressMuxed,
		operationID:  e.operation.ID(),
		effectType:   effectType,
		order:        uint32(len(e.effects) + 1),
		details:      details,
	})
}

func (e *effectsWrapper) addUnmuxed(address *xdr.AccountId, effectType history.EffectType, details map[string]interface{}) {
	e.add(address.Address(), null.String{}, effectType, details)
}

func (e *effectsWrapper) addMuxed(address *xdr.MuxedAccount, effectType history.EffectType, details map[string]interface{}) {
	var addressMuxed null.String
	if address.Type == xdr.CryptoKeyTypeKeyTypeMuxedEd25519 {
		addressMuxed = null.StringFrom(address.Address())
	}
	accID := address.ToAccountId()
	e.add(accID.Address(), addressMuxed, effectType, details)
}

var sponsoringEffectsTable = map[xdr.LedgerEntryType]struct {
	created, updated, removed history.EffectType
}{
	xdr.LedgerEntryTypeAccount: {
		created: history.EffectAccountSponsorshipCreated,
		updated: history.EffectAccountSponsorshipUpdated,
		removed: history.EffectAccountSponsorshipRemoved,
	},
	xdr.LedgerEntryTypeTrustline: {
		created: history.EffectTrustlineSponsorshipCreated,
		updated: history.EffectTrustlineSponsorshipUpdated,
		removed: history.EffectTrustlineSponsorshipRemoved,
	},
	xdr.LedgerEntryTypeData: {
		created: history.EffectDataSponsorshipCreated,
		updated: history.EffectDataSponsorshipUpdated,
		removed: history.EffectDataSponsorshipRemoved,
	},
	xdr.LedgerEntryTypeClaimableBalance: {
		created: history.EffectClaimableBalanceSponsorshipCreated,
		updated: history.EffectClaimableBalanceSponsorshipUpdated,
		removed: history.EffectClaimableBalanceSponsorshipRemoved,
	},

	// We intentionally don't have Sponsoring effects for Offer
	// entries because we don't generate creation effects for them.
}

func (e *effectsWrapper) addSignerSponsorshipEffects(change ingest.Change) {
	if change.Type != xdr.LedgerEntryTypeAccount {
		return
	}

	preSigners := map[string]xdr.AccountId{}
	postSigners := map[string]xdr.AccountId{}
	if change.Pre != nil {
		account := change.Pre.Data.MustAccount()
		preSigners = account.SponsorPerSigner()
	}
	if change.Post != nil {
		account := change.Post.Data.MustAccount()
		postSigners = account.SponsorPerSigner()
	}

	var all []string
	for signer := range preSigners {
		all = append(all, signer)
	}
	for signer := range postSigners {
		if _, ok := preSigners[signer]; ok {
			continue
		}
		all = append(all, signer)
	}
	sort.Strings(all)

	for _, signer := range all {
		pre, foundPre := preSigners[signer]
		post, foundPost := postSigners[signer]
		details := map[string]interface{}{}

		switch {
		case !foundPre && !foundPost:
			continue
		case !foundPre && foundPost:
			details["sponsor"] = post.Address()
			details["signer"] = signer
			srcAccount := change.Post.Data.MustAccount().AccountId
			e.addUnmuxed(&srcAccount, history.EffectSignerSponsorshipCreated, details)
		case !foundPost && foundPre:
			details["former_sponsor"] = pre.Address()
			details["signer"] = signer
			srcAccount := change.Pre.Data.MustAccount().AccountId
			e.addUnmuxed(&srcAccount, history.EffectSignerSponsorshipRemoved, details)
		case foundPre && foundPost:
			formerSponsor := pre.Address()
			newSponsor := post.Address()
			if formerSponsor == newSponsor {
				continue
			}

			details["former_sponsor"] = formerSponsor
			details["new_sponsor"] = newSponsor
			details["signer"] = signer
			srcAccount := change.Post.Data.MustAccount().AccountId
			e.addUnmuxed(&srcAccount, history.EffectSignerSponsorshipUpdated, details)
		}
	}
}

func (e *effectsWrapper) addLedgerEntrySponsorshipEffects(change ingest.Change) error {
	effectsForEntryType, found := sponsoringEffectsTable[change.Type]
	if !found {
		return nil
	}

	details := map[string]interface{}{}
	var effectType history.EffectType

	switch {
	case (change.Pre == nil || change.Pre.SponsoringID() == nil) &&
		(change.Post != nil && change.Post.SponsoringID() != nil):
		effectType = effectsForEntryType.created
		details["sponsor"] = (*change.Post.SponsoringID()).Address()
	case (change.Pre != nil && change.Pre.SponsoringID() != nil) &&
		(change.Post == nil || change.Post.SponsoringID() == nil):
		effectType = effectsForEntryType.removed
		details["former_sponsor"] = (*change.Pre.SponsoringID()).Address()
	case (change.Pre != nil && change.Pre.SponsoringID() != nil) &&
		(change.Post != nil && change.Post.SponsoringID() != nil):
		preSponsor := (*change.Pre.SponsoringID()).Address()
		postSponsor := (*change.Post.SponsoringID()).Address()
		if preSponsor == postSponsor {
			return nil
		}
		effectType = effectsForEntryType.updated
		details["new_sponsor"] = postSponsor
		details["former_sponsor"] = preSponsor
	default:
		return nil
	}

	var (
		accountID    *xdr.AccountId
		muxedAccount *xdr.MuxedAccount
	)

	var data xdr.LedgerEntryData
	if change.Post != nil {
		data = change.Post.Data
	} else {
		data = change.Pre.Data
	}

	switch change.Type {
	case xdr.LedgerEntryTypeAccount:
		a := data.MustAccount().AccountId
		accountID = &a
	case xdr.LedgerEntryTypeTrustline:
		tl := data.MustTrustLine()
		accountID = &tl.AccountId
		if tl.Asset.Type == xdr.AssetTypeAssetTypePoolShare {
			details["asset_type"] = "liquidity_pool"
			details["liquidity_pool_id"] = PoolIDToString(*tl.Asset.LiquidityPoolId)
		} else {
			details["asset"] = tl.Asset.ToAsset().StringCanonical()
		}
	case xdr.LedgerEntryTypeData:
		muxedAccount = e.operation.SourceAccount()
		details["data_name"] = data.MustData().DataName
	case xdr.LedgerEntryTypeClaimableBalance:
		muxedAccount = e.operation.SourceAccount()
		var err error
		details["balance_id"], err = xdr.MarshalHex(data.MustClaimableBalance().BalanceId)
		if err != nil {
			return errors.Wrapf(err, "Invalid balanceId in change from op %d", e.operation.index)
		}
	case xdr.LedgerEntryTypeLiquidityPool:
		// liquidity pools cannot be sponsored
		fallthrough
	default:
		return errors.Errorf("invalid sponsorship ledger entry type %v", change.Type.String())
	}

	if accountID != nil {
		e.addUnmuxed(accountID, effectType, details)
	} else {
		e.addMuxed(muxedAccount, effectType, details)
	}

	return nil
}

func (e *effectsWrapper) addLedgerEntryLiquidityPoolEffects(change ingest.Change) error {
	if change.Type != xdr.LedgerEntryTypeLiquidityPool {
		return nil
	}
	var effectType history.EffectType

	var details map[string]interface{}
	switch {
	case change.Pre == nil && change.Post != nil:
		effectType = history.EffectLiquidityPoolCreated
		details = map[string]interface{}{
			"liquidity_pool": liquidityPoolDetails(change.Post.Data.LiquidityPool),
		}
	case change.Pre != nil && change.Post == nil:
		effectType = history.EffectLiquidityPoolRemoved
		poolID := change.Pre.Data.LiquidityPool.LiquidityPoolId
		details = map[string]interface{}{
			"liquidity_pool_id": PoolIDToString(poolID),
		}
	default:
		return nil
	}
	e.addMuxed(
		e.operation.SourceAccount(),
		effectType,
		details,
	)

	return nil
}

func (e *effectsWrapper) addAccountCreatedEffects() {
	op := e.operation.operation.Body.MustCreateAccountOp()

	e.addUnmuxed(
		&op.Destination,
		history.EffectAccountCreated,
		map[string]interface{}{
			"starting_balance": amount.String(op.StartingBalance),
		},
	)
	e.addMuxed(
		e.operation.SourceAccount(),
		history.EffectAccountDebited,
		map[string]interface{}{
			"asset_type": "native",
			"amount":     amount.String(op.StartingBalance),
		},
	)
	e.addUnmuxed(
		&op.Destination,
		history.EffectSignerCreated,
		map[string]interface{}{
			"public_key": op.Destination.Address(),
			"weight":     keypair.DefaultSignerWeight,
		},
	)
}

func (e *effectsWrapper) addPaymentEffects() {
	op := e.operation.operation.Body.MustPaymentOp()

	details := map[string]interface{}{"amount": amount.String(op.Amount)}
	addAssetDetails(details, op.Asset, "")

	e.addMuxed(
		&op.Destination,
		history.EffectAccountCredited,
		details,
	)
	e.addMuxed(
		e.operation.SourceAccount(),
		history.EffectAccountDebited,
		details,
	)
}

func (e *effectsWrapper) pathPaymentStrictReceiveEffects() error {
	op := e.operation.operation.Body.MustPathPaymentStrictReceiveOp()
	resultSuccess := e.operation.OperationResult().MustPathPaymentStrictReceiveResult().MustSuccess()
	source := e.operation.SourceAccount()

	details := map[string]interface{}{"amount": amount.String(op.DestAmount)}
	addAssetDetails(details, op.DestAsset, "")

	e.addMuxed(
		&op.Destination,
		history.EffectAccountCredited,
		details,
	)

	result := e.operation.OperationResult().MustPathPaymentStrictReceiveResult()
	details = map[string]interface{}{"amount": amount.String(result.SendAmount())}
	addAssetDetails(details, op.SendAsset, "")

	e.addMuxed(
		source,
		history.EffectAccountDebited,
		details,
	)

	return e.addIngestTradeEffects(*source, resultSuccess.Offers)
}

func (e *effectsWrapper) addPathPaymentStrictSendEffects() error {
	source := e.operation.SourceAccount()
	op := e.operation.operation.Body.MustPathPaymentStrictSendOp()
	resultSuccess := e.operation.OperationResult().MustPathPaymentStrictSendResult().MustSuccess()
	result := e.operation.OperationResult().MustPathPaymentStrictSendResult()

	details := map[string]interface{}{"amount": amount.String(result.DestAmount())}
	addAssetDetails(details, op.DestAsset, "")
	e.addMuxed(&op.Destination, history.EffectAccountCredited, details)

	details = map[string]interface{}{"amount": amount.String(op.SendAmount)}
	addAssetDetails(details, op.SendAsset, "")
	e.addMuxed(source, history.EffectAccountDebited, details)

	return e.addIngestTradeEffects(*source, resultSuccess.Offers)
}

func (e *effectsWrapper) addManageSellOfferEffects() error {
	source := e.operation.SourceAccount()
	result := e.operation.OperationResult().MustManageSellOfferResult().MustSuccess()
	return e.addIngestTradeEffects(*source, result.OffersClaimed)
}

func (e *effectsWrapper) addManageBuyOfferEffects() error {
	source := e.operation.SourceAccount()
	result := e.operation.OperationResult().MustManageBuyOfferResult().MustSuccess()
	return e.addIngestTradeEffects(*source, result.OffersClaimed)
}

func (e *effectsWrapper) addCreatePassiveSellOfferEffect() error {
	result := e.operation.OperationResult()
	source := e.operation.SourceAccount()

	var claims []xdr.ClaimAtom

	// KNOWN ISSUE:  diamcircle-core creates results for CreatePassiveOffer operations
	// with the wrong result arm set.
	if result.Type == xdr.OperationTypeManageSellOffer {
		claims = result.MustManageSellOfferResult().MustSuccess().OffersClaimed
	} else {
		claims = result.MustCreatePassiveSellOfferResult().MustSuccess().OffersClaimed
	}

	return e.addIngestTradeEffects(*source, claims)
}

func (e *effectsWrapper) addSetOptionsEffects() error {
	source := e.operation.SourceAccount()
	op := e.operation.operation.Body.MustSetOptionsOp()

	if op.HomeDomain != nil {
		e.addMuxed(source, history.EffectAccountHomeDomainUpdated,
			map[string]interface{}{
				"home_domain": string(*op.HomeDomain),
			},
		)
	}

	thresholdDetails := map[string]interface{}{}

	if op.LowThreshold != nil {
		thresholdDetails["low_threshold"] = *op.LowThreshold
	}

	if op.MedThreshold != nil {
		thresholdDetails["med_threshold"] = *op.MedThreshold
	}

	if op.HighThreshold != nil {
		thresholdDetails["high_threshold"] = *op.HighThreshold
	}

	if len(thresholdDetails) > 0 {
		e.addMuxed(source, history.EffectAccountThresholdsUpdated, thresholdDetails)
	}

	flagDetails := map[string]interface{}{}
	if op.SetFlags != nil {
		setAuthFlagDetails(flagDetails, xdr.AccountFlags(*op.SetFlags), true)
	}
	if op.ClearFlags != nil {
		setAuthFlagDetails(flagDetails, xdr.AccountFlags(*op.ClearFlags), false)
	}

	if len(flagDetails) > 0 {
		e.addMuxed(source, history.EffectAccountFlagsUpdated, flagDetails)
	}

	if op.InflationDest != nil {
		e.addMuxed(source, history.EffectAccountInflationDestinationUpdated,
			map[string]interface{}{
				"inflation_destination": op.InflationDest.Address(),
			},
		)
	}
	changes, err := e.operation.transaction.GetOperationChanges(e.operation.index)
	if err != nil {
		return err
	}

	for _, change := range changes {
		if change.Type != xdr.LedgerEntryTypeAccount {
			continue
		}

		beforeAccount := change.Pre.Data.MustAccount()
		afterAccount := change.Post.Data.MustAccount()

		before := beforeAccount.SignerSummary()
		after := afterAccount.SignerSummary()

		// if before and after are the same, the signers have not changed
		if reflect.DeepEqual(before, after) {
			continue
		}

		beforeSortedSigners := []string{}
		for signer := range before {
			beforeSortedSigners = append(beforeSortedSigners, signer)
		}
		sort.Strings(beforeSortedSigners)

		for _, addy := range beforeSortedSigners {
			weight, ok := after[addy]
			if !ok {
				e.addMuxed(source, history.EffectSignerRemoved, map[string]interface{}{
					"public_key": addy,
				})
				continue
			}

			if weight != before[addy] {
				e.addMuxed(source, history.EffectSignerUpdated, map[string]interface{}{
					"public_key": addy,
					"weight":     weight,
				})
			}
		}

		afterSortedSigners := []string{}
		for signer := range after {
			afterSortedSigners = append(afterSortedSigners, signer)
		}
		sort.Strings(afterSortedSigners)

		// Add the "created" effects
		for _, addy := range afterSortedSigners {
			weight := after[addy]
			// if `addy` is in before, the previous for loop should have recorded
			// the update, so skip this key
			if _, ok := before[addy]; ok {
				continue
			}

			e.addMuxed(source, history.EffectSignerCreated, map[string]interface{}{
				"public_key": addy,
				"weight":     weight,
			})
		}
	}
	return nil
}

func (e *effectsWrapper) addChangeTrustEffects() error {
	source := e.operation.SourceAccount()

	op := e.operation.operation.Body.MustChangeTrustOp()
	changes, err := e.operation.transaction.GetOperationChanges(e.operation.index)
	if err != nil {
		return err
	}

	// NOTE:  when an account trusts itself, the transaction is successful but
	// no ledger entries are actually modified.
	for _, change := range changes {
		if change.Type != xdr.LedgerEntryTypeTrustline {
			continue
		}

		var (
			effect    history.EffectType
			trustLine xdr.TrustLineEntry
		)

		switch {
		case change.Pre == nil && change.Post != nil:
			effect = history.EffectTrustlineCreated
			trustLine = *change.Post.Data.TrustLine
		case change.Pre != nil && change.Post == nil:
			effect = history.EffectTrustlineRemoved
			trustLine = *change.Pre.Data.TrustLine
		case change.Pre != nil && change.Post != nil:
			effect = history.EffectTrustlineUpdated
			trustLine = *change.Post.Data.TrustLine
		default:
			panic("Invalid change")
		}

		// We want to add a single effect for change_trust op. If it's modifying
		// credit_asset search for credit_asset trustline, otherwise search for
		// liquidity_pool.
		if op.Line.Type != trustLine.Asset.Type {
			continue
		}

		details := map[string]interface{}{"limit": amount.String(op.Limit)}
		if trustLine.Asset.Type == xdr.AssetTypeAssetTypePoolShare {
			// The only change_trust ops that can modify LP are those with
			// asset=liquidity_pool so *op.Line.LiquidityPool below is available.
			if err := addLiquidityPoolAssetDetails(details, *op.Line.LiquidityPool); err != nil {
				return err
			}
		} else {
			addAssetDetails(details, op.Line.ToAsset(), "")
		}

		e.addMuxed(source, effect, details)
		break
	}

	return nil
}

func (e *effectsWrapper) addAllowTrustEffects() error {
	source := e.operation.SourceAccount()
	op := e.operation.operation.Body.MustAllowTrustOp()
	asset := op.Asset.ToAsset(source.ToAccountId())
	details := map[string]interface{}{
		"trustor": op.Trustor.Address(),
	}
	addAssetDetails(details, asset, "")

	switch {
	case xdr.TrustLineFlags(op.Authorize).IsAuthorized():
		e.addMuxed(source, history.EffectTrustlineAuthorized, details)
		// Forward compatibility
		setFlags := xdr.Uint32(xdr.TrustLineFlagsAuthorizedFlag)
		e.addTrustLineFlagsEffect(source, &op.Trustor, asset, &setFlags, nil)
	case xdr.TrustLineFlags(op.Authorize).IsAuthorizedToMaintainLiabilitiesFlag():
		e.addMuxed(
			source,
			history.EffectTrustlineAuthorizedToMaintainLiabilities,
			details,
		)
		// Forward compatibility
		setFlags := xdr.Uint32(xdr.TrustLineFlagsAuthorizedToMaintainLiabilitiesFlag)
		e.addTrustLineFlagsEffect(source, &op.Trustor, asset, &setFlags, nil)
	default:
		e.addMuxed(source, history.EffectTrustlineDeauthorized, details)
		// Forward compatibility, show both as cleared
		clearFlags := xdr.Uint32(xdr.TrustLineFlagsAuthorizedFlag | xdr.TrustLineFlagsAuthorizedToMaintainLiabilitiesFlag)
		e.addTrustLineFlagsEffect(source, &op.Trustor, asset, nil, &clearFlags)
	}
	return e.addLiquidityPoolRevokedEffect()
}

func (e *effectsWrapper) addAccountMergeEffects() {
	source := e.operation.SourceAccount()

	dest := e.operation.operation.Body.MustDestination()
	result := e.operation.OperationResult().MustAccountMergeResult()
	details := map[string]interface{}{
		"amount":     amount.String(result.MustSourceAccountBalance()),
		"asset_type": "native",
	}

	e.addMuxed(source, history.EffectAccountDebited, details)
	e.addMuxed(&dest, history.EffectAccountCredited, details)
	e.addMuxed(source, history.EffectAccountRemoved, map[string]interface{}{})
}

func (e *effectsWrapper) addInflationEffects() {
	payouts := e.operation.OperationResult().MustInflationResult().MustPayouts()
	for _, payout := range payouts {
		e.addUnmuxed(&payout.Destination, history.EffectAccountCredited,
			map[string]interface{}{
				"amount":     amount.String(payout.Amount),
				"asset_type": "native",
			},
		)
	}
}

func (e *effectsWrapper) addManageDataEffects() error {
	source := e.operation.SourceAccount()
	op := e.operation.operation.Body.MustManageDataOp()
	details := map[string]interface{}{"name": op.DataName}
	effect := history.EffectType(0)
	changes, err := e.operation.transaction.GetOperationChanges(e.operation.index)
	if err != nil {
		return err
	}

	for _, change := range changes {
		if change.Type != xdr.LedgerEntryTypeData {
			continue
		}

		before := change.Pre
		after := change.Post

		if after != nil {
			raw := after.Data.MustData().DataValue
			details["value"] = base64.StdEncoding.EncodeToString(raw)
		}

		switch {
		case before == nil && after != nil:
			effect = history.EffectDataCreated
		case before != nil && after == nil:
			effect = history.EffectDataRemoved
		case before != nil && after != nil:
			effect = history.EffectDataUpdated
		default:
			panic("Invalid before-and-after state")
		}

		break
	}

	e.addMuxed(source, effect, details)
	return nil
}

func (e *effectsWrapper) addBumpSequenceEffects() error {
	source := e.operation.SourceAccount()
	changes, err := e.operation.transaction.GetOperationChanges(e.operation.index)
	if err != nil {
		return err
	}

	for _, change := range changes {
		if change.Type != xdr.LedgerEntryTypeAccount {
			continue
		}

		before := change.Pre
		after := change.Post

		beforeAccount := before.Data.MustAccount()
		afterAccount := after.Data.MustAccount()

		if beforeAccount.SeqNum != afterAccount.SeqNum {
			details := map[string]interface{}{"new_seq": afterAccount.SeqNum}
			e.addMuxed(source, history.EffectSequenceBumped, details)
		}
		break
	}

	return nil
}

func setClaimableBalanceFlagDetails(details map[string]interface{}, flags xdr.ClaimableBalanceFlags) {
	if flags.IsClawbackEnabled() {
		details["claimable_balance_clawback_enabled_flag"] = true
		return
	}
}

func (e *effectsWrapper) addCreateClaimableBalanceEffects(changes []ingest.Change) error {
	source := e.operation.SourceAccount()
	var cb *xdr.ClaimableBalanceEntry
	for _, change := range changes {
		if change.Type != xdr.LedgerEntryTypeClaimableBalance || change.Post == nil {
			continue
		}
		cb = change.Post.Data.ClaimableBalance
		e.addClaimableBalanceEntryCreatedEffects(source, cb)
		break
	}
	if cb == nil {
		return errors.New("claimable balance entry not found")
	}

	details := map[string]interface{}{
		"amount": amount.String(cb.Amount),
	}
	addAssetDetails(details, cb.Asset, "")
	e.addMuxed(
		source,
		history.EffectAccountDebited,
		details,
	)

	return nil
}

func (e *effectsWrapper) addClaimableBalanceEntryCreatedEffects(source *xdr.MuxedAccount, cb *xdr.ClaimableBalanceEntry) error {
	id, err := xdr.MarshalHex(cb.BalanceId)
	if err != nil {
		return err
	}
	details := map[string]interface{}{
		"balance_id": id,
		"amount":     amount.String(cb.Amount),
		"asset":      cb.Asset.StringCanonical(),
	}
	setClaimableBalanceFlagDetails(details, cb.Flags())
	e.addMuxed(
		source,
		history.EffectClaimableBalanceCreated,
		details,
	)
	// EffectClaimableBalanceClaimantCreated can be generated by
	// `create_claimable_balance` operation but also by `liquidity_pool_withdraw`
	// operation causing a revocation.
	// In case of `create_claimable_balance` we use `op.Claimants` to make
	// effects backward compatible. The reason for this is that Diamcircle-Core
	// changes all `rel_before` predicated to `abs_before` when tx is included
	// in the ledger.
	var claimants []xdr.Claimant
	if op, ok := e.operation.operation.Body.GetCreateClaimableBalanceOp(); ok {
		claimants = op.Claimants
	} else {
		claimants = cb.Claimants
	}
	for _, c := range claimants {
		cv0 := c.MustV0()
		e.addUnmuxed(
			&cv0.Destination,
			history.EffectClaimableBalanceClaimantCreated,
			map[string]interface{}{
				"balance_id": id,
				"amount":     amount.String(cb.Amount),
				"predicate":  cv0.Predicate,
				"asset":      cb.Asset.StringCanonical(),
			},
		)
	}
	return err
}

func (e *effectsWrapper) addClaimClaimableBalanceEffects(changes []ingest.Change) error {
	op := e.operation.operation.Body.MustClaimClaimableBalanceOp()

	balanceID, err := xdr.MarshalHex(op.BalanceId)
	if err != nil {
		return fmt.Errorf("Invalid balanceId in op: %d", e.operation.index)
	}

	var cBalance xdr.ClaimableBalanceEntry
	found := false
	for _, change := range changes {
		if change.Type != xdr.LedgerEntryTypeClaimableBalance {
			continue
		}

		if change.Pre != nil && change.Post == nil {
			cBalance = change.Pre.Data.MustClaimableBalance()
			preBalanceID, err := xdr.MarshalHex(cBalance.BalanceId)
			if err != nil {
				return fmt.Errorf("Invalid balanceId in meta changes for op: %d", e.operation.index)
			}

			if preBalanceID == balanceID {
				found = true
				break
			}
		}
	}

	if !found {
		return fmt.Errorf("Change not found for balanceId : %s", balanceID)
	}

	details := map[string]interface{}{
		"amount":     amount.String(cBalance.Amount),
		"balance_id": balanceID,
		"asset":      cBalance.Asset.StringCanonical(),
	}
	setClaimableBalanceFlagDetails(details, cBalance.Flags())
	source := e.operation.SourceAccount()
	e.addMuxed(
		source,
		history.EffectClaimableBalanceClaimed,
		details,
	)

	details = map[string]interface{}{
		"amount": amount.String(cBalance.Amount),
	}
	addAssetDetails(details, cBalance.Asset, "")
	e.addMuxed(
		source,
		history.EffectAccountCredited,
		details,
	)

	return nil
}

func (e *effectsWrapper) addIngestTradeEffects(buyer xdr.MuxedAccount, claims []xdr.ClaimAtom) error {
	for _, claim := range claims {
		if claim.AmountSold() == 0 && claim.AmountBought() == 0 {
			continue
		}
		switch claim.Type {
		case xdr.ClaimAtomTypeClaimAtomTypeLiquidityPool:
			if err := e.addClaimLiquidityPoolTradeEffect(claim); err != nil {
				return err
			}
		default:
			e.addClaimTradeEffects(buyer, claim)
		}
	}
	return nil
}

func (e *effectsWrapper) addClaimTradeEffects(buyer xdr.MuxedAccount, claim xdr.ClaimAtom) {
	seller := claim.SellerId()
	bd, sd := tradeDetails(buyer, seller, claim)

	e.addMuxed(
		&buyer,
		history.EffectTrade,
		bd,
	)

	e.addUnmuxed(
		&seller,
		history.EffectTrade,
		sd,
	)
}

func (e *effectsWrapper) addClaimLiquidityPoolTradeEffect(claim xdr.ClaimAtom) error {
	lp, _, err := e.operation.getLiquidityPoolAndProductDelta(&claim.LiquidityPool.LiquidityPoolId)
	if err != nil {
		return err
	}
	details := map[string]interface{}{
		"liquidity_pool": liquidityPoolDetails(lp),
		"sold": map[string]string{
			"asset":  claim.LiquidityPool.AssetSold.StringCanonical(),
			"amount": amount.String(claim.LiquidityPool.AmountSold),
		},
		"bought": map[string]string{
			"asset":  claim.LiquidityPool.AssetBought.StringCanonical(),
			"amount": amount.String(claim.LiquidityPool.AmountBought),
		},
	}
	e.addMuxed(e.operation.SourceAccount(), history.EffectLiquidityPoolTrade, details)
	return nil
}

func (e *effectsWrapper) addClawbackEffects() error {
	op := e.operation.operation.Body.MustClawbackOp()
	details := map[string]interface{}{
		"amount": amount.String(op.Amount),
	}
	source := e.operation.SourceAccount()
	addAssetDetails(details, op.Asset, "")

	// The funds will be burned, but even with that, we generated an account credited effect
	e.addMuxed(
		source,
		history.EffectAccountCredited,
		details,
	)

	e.addMuxed(
		&op.From,
		history.EffectAccountDebited,
		details,
	)

	return nil
}

func (e *effectsWrapper) addClawbackClaimableBalanceEffects(changes []ingest.Change) error {
	op := e.operation.operation.Body.MustClawbackClaimableBalanceOp()
	balanceId, err := xdr.MarshalHex(op.BalanceId)
	if err != nil {
		return errors.Wrapf(err, "Invalid balanceId in op %d", e.operation.index)
	}
	details := map[string]interface{}{
		"balance_id": balanceId,
	}
	source := e.operation.SourceAccount()
	e.addMuxed(
		source,
		history.EffectClaimableBalanceClawedBack,
		details,
	)

	// Generate the account credited effect (although the funds will be burned) for the asset issuer
	for _, c := range changes {
		if c.Type == xdr.LedgerEntryTypeClaimableBalance && c.Post == nil && c.Pre != nil {
			cb := c.Pre.Data.ClaimableBalance
			details = map[string]interface{}{"amount": amount.String(cb.Amount)}
			addAssetDetails(details, cb.Asset, "")
			e.addMuxed(
				source,
				history.EffectAccountCredited,
				details,
			)
			break
		}
	}

	return nil
}

func (e *effectsWrapper) addSetTrustLineFlagsEffects() error {
	source := e.operation.SourceAccount()
	op := e.operation.operation.Body.MustSetTrustLineFlagsOp()
	e.addTrustLineFlagsEffect(source, &op.Trustor, op.Asset, &op.SetFlags, &op.ClearFlags)
	return e.addLiquidityPoolRevokedEffect()
}

func (e *effectsWrapper) addTrustLineFlagsEffect(
	account *xdr.MuxedAccount,
	trustor *xdr.AccountId,
	asset xdr.Asset,
	setFlags *xdr.Uint32,
	clearFlags *xdr.Uint32) {
	details := map[string]interface{}{
		"trustor": trustor.Address(),
	}
	addAssetDetails(details, asset, "")

	var flagDetailsAdded bool
	if setFlags != nil {
		setTrustLineFlagDetails(details, xdr.TrustLineFlags(*setFlags), true)
		flagDetailsAdded = true
	}
	if clearFlags != nil {
		setTrustLineFlagDetails(details, xdr.TrustLineFlags(*clearFlags), false)
		flagDetailsAdded = true
	}

	if flagDetailsAdded {
		e.addMuxed(account, history.EffectTrustlineFlagsUpdated, details)
	}
}

func setTrustLineFlagDetails(flagDetails map[string]interface{}, flags xdr.TrustLineFlags, setValue bool) {
	if flags.IsAuthorized() {
		flagDetails["authorized_flag"] = setValue
	}
	if flags.IsAuthorizedToMaintainLiabilitiesFlag() {
		flagDetails["authorized_to_maintain_liabilites"] = setValue
	}
	if flags.IsClawbackEnabledFlag() {
		flagDetails["clawback_enabled_flag"] = setValue
	}
}

type sortableClaimableBalanceEntries []*xdr.ClaimableBalanceEntry

func (s sortableClaimableBalanceEntries) Len() int           { return len(s) }
func (s sortableClaimableBalanceEntries) Less(i, j int) bool { return s[i].Asset.LessThan(s[j].Asset) }
func (s sortableClaimableBalanceEntries) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (e *effectsWrapper) addLiquidityPoolRevokedEffect() error {
	source := e.operation.SourceAccount()
	lp, delta, err := e.operation.getLiquidityPoolAndProductDelta(nil)
	if err != nil {
		if err == errLiquidityPoolChangeNotFound {
			// no revocation happened
			return nil
		}
		return err
	}
	changes, err := e.operation.transaction.GetOperationChanges(e.operation.index)
	if err != nil {
		return err
	}
	assetToCBID := map[string]string{}
	var cbs sortableClaimableBalanceEntries
	for _, change := range changes {
		if change.Type == xdr.LedgerEntryTypeClaimableBalance && change.Pre == nil && change.Post != nil {
			cb := change.Post.Data.ClaimableBalance
			id, err := xdr.MarshalHex(cb.BalanceId)
			if err != nil {
				return err
			}
			assetToCBID[cb.Asset.StringCanonical()] = id
			cbs = append(cbs, cb)
		}
	}
	if len(assetToCBID) == 0 {
		// no claimable balances were created, and thus, no revocation happened
		return nil
	}
	// Core's claimable balance metadata isn't ordered, so we order it ourselves
	// so that effects are ordered consistently
	sort.Sort(cbs)
	for _, cb := range cbs {
		if err := e.addClaimableBalanceEntryCreatedEffects(source, cb); err != nil {
			return err
		}
	}

	reservesRevoked := make([]map[string]string, 0, 2)
	for _, aa := range []base.AssetAmount{
		{
			Asset:  lp.Body.ConstantProduct.Params.AssetA.StringCanonical(),
			Amount: amount.String(-delta.ReserveA),
		},
		{
			Asset:  lp.Body.ConstantProduct.Params.AssetB.StringCanonical(),
			Amount: amount.String(-delta.ReserveB),
		},
	} {
		if cbID, ok := assetToCBID[aa.Asset]; ok {
			assetAmountDetail := map[string]string{
				"asset":                aa.Asset,
				"amount":               aa.Amount,
				"claimable_balance_id": cbID,
			}
			reservesRevoked = append(reservesRevoked, assetAmountDetail)
		}
	}
	details := map[string]interface{}{
		"liquidity_pool":   liquidityPoolDetails(lp),
		"reserves_revoked": reservesRevoked,
		"shares_revoked":   amount.String(-delta.TotalPoolShares),
	}
	e.addMuxed(source, history.EffectLiquidityPoolRevoked, details)
	return nil
}

func setAuthFlagDetails(flagDetails map[string]interface{}, flags xdr.AccountFlags, setValue bool) {
	if flags.IsAuthRequired() {
		flagDetails["auth_required_flag"] = setValue
	}
	if flags.IsAuthRevocable() {
		flagDetails["auth_revocable_flag"] = setValue
	}
	if flags.IsAuthImmutable() {
		flagDetails["auth_immutable_flag"] = setValue
	}
	if flags.IsAuthClawbackEnabled() {
		flagDetails["auth_clawback_enabled_flag"] = setValue
	}
}

func tradeDetails(buyer xdr.MuxedAccount, seller xdr.AccountId, claim xdr.ClaimAtom) (bd map[string]interface{}, sd map[string]interface{}) {
	bd = map[string]interface{}{
		"offer_id":      claim.OfferId(),
		"seller":        seller.Address(),
		"bought_amount": amount.String(claim.AmountSold()),
		"sold_amount":   amount.String(claim.AmountBought()),
	}
	addAssetDetails(bd, claim.AssetSold(), "bought_")
	addAssetDetails(bd, claim.AssetBought(), "sold_")

	sd = map[string]interface{}{
		"offer_id":      claim.OfferId(),
		"bought_amount": amount.String(claim.AmountBought()),
		"sold_amount":   amount.String(claim.AmountSold()),
	}
	addAccountAndMuxedAccountDetails(sd, buyer, "seller")
	addAssetDetails(sd, claim.AssetBought(), "bought_")
	addAssetDetails(sd, claim.AssetSold(), "sold_")

	return
}

func liquidityPoolDetails(lp *xdr.LiquidityPoolEntry) map[string]interface{} {
	return map[string]interface{}{
		"id":               PoolIDToString(lp.LiquidityPoolId),
		"fee_bp":           uint32(lp.Body.ConstantProduct.Params.Fee),
		"type":             "constant_product",
		"total_trustlines": strconv.FormatInt(int64(lp.Body.ConstantProduct.PoolSharesTrustLineCount), 10),
		"total_shares":     amount.String(lp.Body.ConstantProduct.TotalPoolShares),
		"reserves": []base.AssetAmount{
			{
				Asset:  lp.Body.ConstantProduct.Params.AssetA.StringCanonical(),
				Amount: amount.String(lp.Body.ConstantProduct.ReserveA),
			},
			{
				Asset:  lp.Body.ConstantProduct.Params.AssetB.StringCanonical(),
				Amount: amount.String(lp.Body.ConstantProduct.ReserveB),
			},
		},
	}
}

func (e *effectsWrapper) addLiquidityPoolDepositEffect() error {
	op := e.operation.operation.Body.MustLiquidityPoolDepositOp()
	lp, delta, err := e.operation.getLiquidityPoolAndProductDelta(&op.LiquidityPoolId)
	if err != nil {
		return err
	}
	details := map[string]interface{}{
		"liquidity_pool": liquidityPoolDetails(lp),
		"reserves_deposited": []base.AssetAmount{
			{
				Asset:  lp.Body.ConstantProduct.Params.AssetA.StringCanonical(),
				Amount: amount.String(delta.ReserveA),
			},
			{
				Asset:  lp.Body.ConstantProduct.Params.AssetB.StringCanonical(),
				Amount: amount.String(delta.ReserveB),
			},
		},
		"shares_received": amount.String(delta.TotalPoolShares),
	}
	e.addMuxed(e.operation.SourceAccount(), history.EffectLiquidityPoolDeposited, details)
	return nil
}

func (e *effectsWrapper) addLiquidityPoolWithdrawEffect() error {
	op := e.operation.operation.Body.MustLiquidityPoolWithdrawOp()
	lp, delta, err := e.operation.getLiquidityPoolAndProductDelta(&op.LiquidityPoolId)
	if err != nil {
		return err
	}
	details := map[string]interface{}{
		"liquidity_pool": liquidityPoolDetails(lp),
		"reserves_received": []base.AssetAmount{
			{
				Asset:  lp.Body.ConstantProduct.Params.AssetA.StringCanonical(),
				Amount: amount.String(-delta.ReserveA),
			},
			{
				Asset:  lp.Body.ConstantProduct.Params.AssetB.StringCanonical(),
				Amount: amount.String(-delta.ReserveB),
			},
		},
		"shares_redeemed": amount.String(-delta.TotalPoolShares),
	}
	e.addMuxed(e.operation.SourceAccount(), history.EffectLiquidityPoolWithdrew, details)
	return nil
}
