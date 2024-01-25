// Package operationfeestats provides useful utilities concerning operation fee
// stats within diamcircle,specifically as a central location to store a cached snapshot
// of the state of network per operation fees and surge pricing. This package is
// intended to be at the lowest levels of aurora's dependency tree, please keep
// it free of dependencies to other aurora packages.
package operationfeestats

import (
	"sync"
)

// State represents a snapshot of aurora's view of the state of operation fee's
// on the network.
type State struct {
	FeeChargedMax  int64
	FeeChargedMin  int64
	FeeChargedMode int64
	FeeChargedP10  int64
	FeeChargedP20  int64
	FeeChargedP30  int64
	FeeChargedP40  int64
	FeeChargedP50  int64
	FeeChargedP60  int64
	FeeChargedP70  int64
	FeeChargedP80  int64
	FeeChargedP90  int64
	FeeChargedP95  int64
	FeeChargedP99  int64

	// MaxFee
	MaxFeeMax  int64
	MaxFeeMin  int64
	MaxFeeMode int64
	MaxFeeP10  int64
	MaxFeeP20  int64
	MaxFeeP30  int64
	MaxFeeP40  int64
	MaxFeeP50  int64
	MaxFeeP60  int64
	MaxFeeP70  int64
	MaxFeeP80  int64
	MaxFeeP90  int64
	MaxFeeP95  int64
	MaxFeeP99  int64

	LastBaseFee         int64
	LastLedger          uint32
	LedgerCapacityUsage string
}

// CurrentState returns the cached snapshot of operation fee state and a boolean indicating
// if the cache has been populated
func CurrentState() (State, bool) {
	lock.RLock()
	ret := current
	ok := present
	lock.RUnlock()
	return ret, ok
}

// SetState updates the cached snapshot of the operation fee state
func SetState(next State) {
	lock.Lock()
	// in case of one query taking longer than another, this makes
	// sure we don't overwrite the latest fee stats with old stats
	if current.LastLedger < next.LastLedger {
		current = next
	}
	present = true
	lock.Unlock()
}

// ResetState is used only for testing purposes
func ResetState() {
	current = State{}
	present = false
}

var current State
var present bool
var lock sync.RWMutex
