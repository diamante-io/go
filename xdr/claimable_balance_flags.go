package xdr

// IsClawbackEnabled returns true if the claimable balance has the "CLAWBACK_ENABLED" option
// turned on.
func (cbFlags ClaimableBalanceFlags) IsClawbackEnabled() bool {
	return (cbFlags & ClaimableBalanceFlagsClaimableBalanceClawbackEnabledFlag) != 0
}
