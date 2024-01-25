package resourceadapter

import (
	"context"
	"fmt"
	"strconv"

	protocol "go/protocols/aurora"
	auroraContext "go/services/aurora/internal/context"
	"go/services/aurora/internal/db2/history"
	"go/support/errors"
	"go/support/render/hal"
	"go/xdr"
)

// PopulateAccountEntry fills out the resource's fields
func PopulateAccountEntry(
	ctx context.Context,
	dest *protocol.Account,
	account history.AccountEntry,
	accountData []history.Data,
	accountSigners []history.AccountSigner,
	trustLines []history.TrustLine,
	ledger *history.Ledger,
) error {
	dest.ID = account.AccountID
	dest.PT = account.AccountID
	dest.AccountID = account.AccountID
	dest.Sequence = strconv.FormatInt(account.SequenceNumber, 10)
	dest.SubentryCount = int32(account.NumSubEntries)
	dest.InflationDestination = account.InflationDestination
	dest.HomeDomain = account.HomeDomain
	dest.LastModifiedLedger = account.LastModifiedLedger
	if ledger != nil {
		dest.LastModifiedTime = &ledger.ClosedAt
	}

	dest.Flags.AuthRequired = account.IsAuthRequired()
	dest.Flags.AuthRevocable = account.IsAuthRevocable()
	dest.Flags.AuthImmutable = account.IsAuthImmutable()
	dest.Flags.AuthClawbackEnabled = account.IsAuthClawbackEnabled()

	dest.Thresholds.LowThreshold = account.ThresholdLow
	dest.Thresholds.MedThreshold = account.ThresholdMedium
	dest.Thresholds.HighThreshold = account.ThresholdHigh

	// populate balances
	dest.Balances = make([]protocol.Balance, len(trustLines)+1)
	for i, tl := range trustLines {
		var err error

		switch tl.AssetType {
		case xdr.AssetTypeAssetTypePoolShare:
			err = PopulatePoolShareBalance(&dest.Balances[i], tl)
		default:
			err = PopulateAssetBalance(&dest.Balances[i], tl)
		}

		if err != nil {
			return errors.Wrap(err, "populating balance")
		}
	}

	// add native balance
	err := PopulateNativeBalance(
		&dest.Balances[len(dest.Balances)-1],
		xdr.Int64(account.Balance),
		xdr.Int64(account.BuyingLiabilities),
		xdr.Int64(account.SellingLiabilities),
	)
	if err != nil {
		return errors.Wrap(err, "populating native balance")
	}

	// populate data
	dest.Data = make(map[string]string)
	for _, d := range accountData {
		dest.Data[d.Name] = d.Value.Base64()
	}

	masterKeyIncluded := false

	// populate signers
	dest.Signers = make([]protocol.Signer, len(accountSigners))
	for i, signer := range accountSigners {
		dest.Signers[i].Weight = signer.Weight
		dest.Signers[i].Key = signer.Signer
		dest.Signers[i].Type = protocol.MustKeyTypeFromAddress(signer.Signer)
		if signer.Sponsor.Valid {
			dest.Signers[i].Sponsor = signer.Sponsor.String
		}

		if account.AccountID == signer.Signer {
			masterKeyIncluded = true
		}
	}

	if !masterKeyIncluded {
		dest.Signers = append(dest.Signers, protocol.Signer{
			Weight: int32(account.MasterWeight),
			Key:    account.AccountID,
			Type:   protocol.MustKeyTypeFromAddress(account.AccountID),
		})
	}

	dest.NumSponsoring = account.NumSponsoring
	dest.NumSponsored = account.NumSponsored
	if account.Sponsor.Valid {
		dest.Sponsor = account.Sponsor.String
	}

	lb := hal.LinkBuilder{auroraContext.BaseURL(ctx)}
	self := fmt.Sprintf("/accounts/%s", account.AccountID)
	dest.Links.Self = lb.Link(self)
	dest.Links.Transactions = lb.PagedLink(self, "transactions")
	dest.Links.Operations = lb.PagedLink(self, "operations")
	dest.Links.Payments = lb.PagedLink(self, "payments")
	dest.Links.Effects = lb.PagedLink(self, "effects")
	dest.Links.Offers = lb.PagedLink(self, "offers")
	dest.Links.Trades = lb.PagedLink(self, "trades")
	dest.Links.Data = lb.Link(self, "data/{key}")
	dest.Links.Data.PopulateTemplated()
	return nil
}
