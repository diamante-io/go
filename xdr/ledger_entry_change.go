package xdr

import (
	"encoding/base64"
	"fmt"
)

// EntryType is a helper to get at the entry type for a change.
func (change *LedgerEntryChange) EntryType() LedgerEntryType {
	return change.LedgerKey().Type
}

// LedgerKey returns the key for the ledger entry that was changed
// in `change`.
func (change *LedgerEntryChange) LedgerKey() LedgerKey {
	switch change.Type {
	case LedgerEntryChangeTypeLedgerEntryCreated:
		change := change.MustCreated()
		return change.LedgerKey()
	case LedgerEntryChangeTypeLedgerEntryRemoved:
		return change.MustRemoved()
	case LedgerEntryChangeTypeLedgerEntryUpdated:
		change := change.MustUpdated()
		return change.LedgerKey()
	case LedgerEntryChangeTypeLedgerEntryState:
		change := change.MustState()
		return change.LedgerKey()
	default:
		panic(fmt.Errorf("Unknown change type: %v", change.Type))
	}
}

// MarshalBinaryBase64 marshals XDR into a binary form and then encodes it
// using base64.
func (change LedgerEntryChange) MarshalBinaryBase64() (string, error) {
	b, err := change.MarshalBinary()
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

// GetLedgerEntry returns the ledger entry that was changed in `change`, along
// with a boolean indicating whether the entry value was valid.
func (change *LedgerEntryChange) GetLedgerEntry() (LedgerEntry, bool) {
	switch change.Type {
	case LedgerEntryChangeTypeLedgerEntryCreated:
		return change.GetCreated()
	case LedgerEntryChangeTypeLedgerEntryState:
		return change.GetState()
	case LedgerEntryChangeTypeLedgerEntryUpdated:
		return change.GetUpdated()
	case LedgerEntryChangeTypeLedgerEntryRemoved:
		return LedgerEntry{}, false
	default:
		return LedgerEntry{}, false
	}
}
