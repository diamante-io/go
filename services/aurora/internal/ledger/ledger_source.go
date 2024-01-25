package ledger

import (
	"sync"
	"time"
)

// Source exposes two helpers methods to help you find out the current
// ledger and yield every time there is a new ledger. Call `Close` when
// source is no longer used.
type Source interface {
	CurrentLedger() uint32
	NextLedger(currentSequence uint32) chan uint32
	Close()
}

// HistoryDBSource utility struct to pass the SSE update frequency and a
// function to get the current ledger state.
type HistoryDBSource struct {
	updateFrequency time.Duration
	state           *State

	closedLock sync.Mutex
	closed     bool
}

// NewHistoryDBSource constructs a new instance of HistoryDBSource
func NewHistoryDBSource(updateFrequency time.Duration, state *State) *HistoryDBSource {
	return &HistoryDBSource{
		updateFrequency: updateFrequency,
		state:           state,
		closedLock:      sync.Mutex{},
	}
}

// CurrentLedger returns the current ledger.
func (source *HistoryDBSource) CurrentLedger() uint32 {
	return source.state.CurrentStatus().ExpHistoryLatest
}

// NextLedger returns a channel which yields every time there is a new ledger with a sequence number larger than currentSequence.
func (source *HistoryDBSource) NextLedger(currentSequence uint32) chan uint32 {
	// Make sure this is buffered channel of size 1. Otherwise, the go routine below
	// will never return if `newLedgers` channel is not read. From Effective Go:
	// > If the channel is unbuffered, the sender blocks until the receiver has received the value.
	newLedgers := make(chan uint32, 1)
	go func() {
		for {
			if source.updateFrequency > 0 {
				time.Sleep(source.updateFrequency)
			}

			source.closedLock.Lock()
			closed := source.closed
			source.closedLock.Unlock()
			if closed {
				return
			}

			currentLedgerState := source.state.CurrentStatus()
			if currentLedgerState.ExpHistoryLatest > currentSequence {
				newLedgers <- currentLedgerState.ExpHistoryLatest
				return
			}
		}
	}()

	return newLedgers
}

// Close closes the internal go routines.
func (source *HistoryDBSource) Close() {
	source.closedLock.Lock()
	defer source.closedLock.Unlock()
	source.closed = true
}

// TestingSource is helper struct which implements the LedgerSource
// interface.
type TestingSource struct {
	currentLedger uint32
	newLedgers    chan uint32
	lock          *sync.RWMutex
}

// NewTestingSource returns a TestingSource.
func NewTestingSource(currentLedger uint32) *TestingSource {
	return &TestingSource{
		currentLedger: currentLedger,
		newLedgers:    make(chan uint32),
		lock:          &sync.RWMutex{},
	}
}

// CurrentLedger returns the current ledger.
func (source *TestingSource) CurrentLedger() uint32 {
	source.lock.RLock()
	defer source.lock.RUnlock()
	return source.currentLedger
}

// AddLedger adds a new sequence to the newLedgers channel. AddLedger()
// will block until the new sequence is read
func (source *TestingSource) AddLedger(nextSequence uint32) {
	source.newLedgers <- nextSequence
}

// NextLedger returns a channel which yields every time there is a new ledger.
func (source *TestingSource) NextLedger(currentSequence uint32) chan uint32 {
	response := make(chan uint32, 1)

	go func() {
		for {
			nextLedger := <-source.newLedgers
			if nextLedger > source.currentLedger {
				source.lock.Lock()
				defer source.lock.Unlock()
				source.currentLedger = nextLedger
				response <- nextLedger
				return
			}
		}
	}()

	return response
}

func (source *TestingSource) Close() {}
