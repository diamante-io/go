package xdr

func (account *AccountEntry) SignerSummary() map[string]int32 {
	ret := map[string]int32{}

	if account.MasterKeyWeight() > 0 {
		ret[account.AccountId.Address()] = int32(account.Thresholds[0])
	}
	for _, signer := range account.Signers {
		ret[signer.Key.Address()] = int32(signer.Weight)
	}

	return ret
}

func (account *AccountEntry) MasterKeyWeight() byte {
	return account.Thresholds.MasterKeyWeight()
}

func (account *AccountEntry) ThresholdLow() byte {
	return account.Thresholds.ThresholdLow()
}

func (account *AccountEntry) ThresholdMedium() byte {
	return account.Thresholds.ThresholdMedium()
}

func (account *AccountEntry) ThresholdHigh() byte {
	return account.Thresholds.ThresholdHigh()
}

// Liabilities returns AccountEntry's liabilities
func (account *AccountEntry) Liabilities() Liabilities {
	var liabilities Liabilities
	if account.Ext.V1 != nil {
		liabilities = account.Ext.V1.Liabilities
	}
	return liabilities
}

// NumSponsored returns NumSponsored value for account.
func (account *AccountEntry) NumSponsored() Uint32 {
	var numSponsored Uint32
	if account.Ext.V1 != nil &&
		account.Ext.V1.Ext.V2 != nil {
		numSponsored = account.Ext.V1.Ext.V2.NumSponsored
	}
	return numSponsored
}

// NumSponsoring returns NumSponsoring value for account.
func (account *AccountEntry) NumSponsoring() Uint32 {
	var numSponsoring Uint32
	if account.Ext.V1 != nil &&
		account.Ext.V1.Ext.V2 != nil {
		numSponsoring = account.Ext.V1.Ext.V2.NumSponsoring
	}
	return numSponsoring
}

// SignerSponsoringIDs returns SignerSponsoringIDs value for account.
// This will return a slice of nil values if V2 extension does not exist.
func (account *AccountEntry) SignerSponsoringIDs() []SponsorshipDescriptor {
	var ids []SponsorshipDescriptor
	if account.Ext.V1 != nil &&
		account.Ext.V1.Ext.V2 != nil {
		ids = account.Ext.V1.Ext.V2.SignerSponsoringIDs
	} else {
		ids = make([]SponsorshipDescriptor, len(account.Signers))
	}
	return ids
}

// SponsorPerSigner returns a mapping of signer to its sponsor
func (account *AccountEntry) SponsorPerSigner() map[string]AccountId {
	ids := account.SignerSponsoringIDs()

	signerToSponsor := map[string]AccountId{}

	for i, signerEntry := range account.Signers {
		if ids[i] != nil {
			signerToSponsor[signerEntry.Key.Address()] = *ids[i]
		}
	}

	return signerToSponsor
}
