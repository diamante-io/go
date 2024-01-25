package xdr

import "fmt"

// LedgerKey implements the `Keyer` interface
func (entry *LedgerEntry) LedgerKey() LedgerKey {
	var body interface{}

	switch entry.Data.Type {
	case LedgerEntryTypeAccount:
		account := entry.Data.MustAccount()
		body = LedgerKeyAccount{
			AccountId: account.AccountId,
		}
	case LedgerEntryTypeData:
		data := entry.Data.MustData()
		body = LedgerKeyData{
			AccountId: data.AccountId,
			DataName:  data.DataName,
		}
	case LedgerEntryTypeOffer:
		offer := entry.Data.MustOffer()
		body = LedgerKeyOffer{
			SellerId: offer.SellerId,
			OfferId:  offer.OfferId,
		}
	case LedgerEntryTypeTrustline:
		tline := entry.Data.MustTrustLine()
		body = LedgerKeyTrustLine{
			AccountId: tline.AccountId,
			Asset:     tline.Asset,
		}
	case LedgerEntryTypeClaimableBalance:
		cBalance := entry.Data.MustClaimableBalance()
		body = LedgerKeyClaimableBalance{
			BalanceId: cBalance.BalanceId,
		}
	case LedgerEntryTypeLiquidityPool:
		lPool := entry.Data.MustLiquidityPool()
		body = LedgerKeyLiquidityPool{
			LiquidityPoolId: lPool.LiquidityPoolId,
		}
	default:
		panic(fmt.Errorf("Unknown entry type: %v", entry.Data.Type))
	}

	ret, err := NewLedgerKey(entry.Data.Type, body)
	if err != nil {
		panic(err)
	}

	return ret
}

// SponsoringID return SponsorshipDescriptor for a given ledger entry
func (entry *LedgerEntry) SponsoringID() SponsorshipDescriptor {
	var sponsor SponsorshipDescriptor
	if entry.Ext.V == 1 && entry.Ext.V1 != nil {
		sponsor = entry.Ext.V1.SponsoringId
	}
	return sponsor
}

// Normalize overwrites LedgerEntry with all the extensions set to default values
// (if extension is not present).
// This is helpful to compare two ledger entries that are the same but for one of
// them extensions are not set.
// Returns the same entry.
func (entry *LedgerEntry) Normalize() *LedgerEntry {
	// If ledgerEntry is V0, create ext=1 and add a nil sponsor
	if entry.Ext.V == 0 {
		entry.Ext = LedgerEntryExt{
			V: 1,
			V1: &LedgerEntryExtensionV1{
				SponsoringId: nil,
			},
		}
	}

	switch entry.Data.Type {
	case LedgerEntryTypeAccount:
		accountEntry := entry.Data.Account
		// Account can have ext=0. For those, create ext=1 with 0 liabilities.
		if accountEntry.Ext.V == 0 {
			accountEntry.Ext.V = 1
			accountEntry.Ext.V1 = &AccountEntryExtensionV1{
				Liabilities: Liabilities{
					Buying:  0,
					Selling: 0,
				},
			}
		}
		// if AccountEntryExtensionV1Ext is v=0, then create v2 with 0 values
		if accountEntry.Ext.V1.Ext.V == 0 {
			accountEntry.Ext.V1.Ext.V = 2
			accountEntry.Ext.V1.Ext.V2 = &AccountEntryExtensionV2{
				NumSponsored:        Uint32(0),
				NumSponsoring:       Uint32(0),
				SignerSponsoringIDs: make([]SponsorshipDescriptor, len(accountEntry.Signers)),
			}
		}

		signerSponsoringIDs := accountEntry.Ext.V1.Ext.V2.SignerSponsoringIDs

		// Map sponsors (signer => SponsorshipDescriptor)
		sponsorsMap := make(map[string]SponsorshipDescriptor)
		for i, signer := range accountEntry.Signers {
			sponsorsMap[signer.Key.Address()] = signerSponsoringIDs[i]
		}

		// Sort signers
		accountEntry.Signers = SortSignersByKey(accountEntry.Signers)

		// Recreate sponsors for sorted signers
		signerSponsoringIDs = make([]SponsorshipDescriptor, len(accountEntry.Signers))
		for i, signer := range accountEntry.Signers {
			signerSponsoringIDs[i] = sponsorsMap[signer.Key.Address()]
		}

		accountEntry.Ext.V1.Ext.V2.SignerSponsoringIDs = signerSponsoringIDs
	case LedgerEntryTypeTrustline:
		// Trust line can have ext=0. For those, create ext=1
		// with 0 liabilities.
		trustLineEntry := entry.Data.TrustLine
		if trustLineEntry.Ext.V == 0 {
			trustLineEntry.Ext.V = 1
			trustLineEntry.Ext.V1 = &TrustLineEntryV1{
				Liabilities: Liabilities{
					Buying:  0,
					Selling: 0,
				},
			}
		} else if trustLineEntry.Ext.V == 1 {
			// aurora doesn't make use of TrustLineEntryExtensionV2.liquidityPoolUseCount
			// so clear out those fields to make state verifier pass
			trustLineEntry.Ext.V1.Ext.V = 0
			trustLineEntry.Ext.V1.Ext.V2 = nil
		}
	case LedgerEntryTypeClaimableBalance:
		claimableBalanceEntry := entry.Data.ClaimableBalance
		claimableBalanceEntry.Claimants = SortClaimantsByDestination(claimableBalanceEntry.Claimants)

		if claimableBalanceEntry.Ext.V == 0 {
			claimableBalanceEntry.Ext.V = 1
			claimableBalanceEntry.Ext.V1 = &ClaimableBalanceEntryExtensionV1{
				Flags: 0,
			}
		}
	}

	return entry
}
