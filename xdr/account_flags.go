package xdr

// IsAuthRequired returns true if the account has the "AUTH_REQUIRED" option
// turned on.
func (accountFlags AccountFlags) IsAuthRequired() bool {
	return (accountFlags & AccountFlagsAuthRequiredFlag) != 0
}

// IsAuthRevocable returns true if the account has the "AUTH_REVOCABLE" option
// turned on.
func (accountFlags AccountFlags) IsAuthRevocable() bool {
	return (accountFlags & AccountFlagsAuthRevocableFlag) != 0
}

// IsAuthImmutable returns true if the account has the "AUTH_IMMUTABLE" option
// turned on.
func (accountFlags AccountFlags) IsAuthImmutable() bool {
	return (accountFlags & AccountFlagsAuthImmutableFlag) != 0
}

// IsAuthClawbackEnabled returns true if the account has the "AUTH_CLAWBACK_ENABLED" option
// turned on.
func (accountFlags AccountFlags) IsAuthClawbackEnabled() bool {
	return (accountFlags & AccountFlagsAuthClawbackEnabledFlag) != 0
}
