package xdr

import (
	"sort"
)

// SortClaimantsByDestination returns a new []Claimant array sorted by destination.
func SortClaimantsByDestination(claimants []Claimant) []Claimant {
	keys := make([]string, 0, len(claimants))
	keysMap := make(map[string]Claimant)
	newClaimants := make([]Claimant, 0, len(claimants))

	for _, claimant := range claimants {
		v0 := claimant.MustV0()
		key := v0.Destination.Address()
		keys = append(keys, key)
		keysMap[key] = claimant
	}

	sort.Strings(keys)

	for _, key := range keys {
		newClaimants = append(newClaimants, keysMap[key])
	}

	return newClaimants
}
