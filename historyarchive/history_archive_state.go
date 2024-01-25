// Copyright 2016 diamcircle Development Foundation and contributors. Licensed
// under the Apache License, Version 2.0. See the COPYING file at the root
// of this distribution or at http://www.apache.org/licenses/LICENSE-2.0

package historyarchive

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"go/support/errors"
	"go/xdr"
)

const NumLevels = 11

type HistoryArchiveState struct {
	Version       int    `json:"version"`
	Server        string `json:"server"`
	CurrentLedger uint32 `json:"currentLedger"`
	// NetworkPassphrase was added in diamcircle-Core v14.1.0. Can be missing
	// in HAS created by previous versions.
	NetworkPassphrase string `json:"networkPassphrase"`
	CurrentBuckets    [NumLevels]struct {
		Curr string `json:"curr"`
		Snap string `json:"snap"`
		Next struct {
			State  uint32 `json:"state"`
			Output string `json:"output,omitempty"`
		} `json:"next"`
	} `json:"currentBuckets"`
}

func (h *HistoryArchiveState) LevelSummary() (string, int, error) {
	summ := ""
	nz := 0
	for _, b := range h.CurrentBuckets {
		state := '_'
		for _, bs := range []string{
			b.Curr, b.Snap, b.Next.Output,
		} {
			// Ignore empty values
			if bs == "" {
				continue
			}

			h, err := DecodeHash(bs)
			if err != nil {
				return summ, nz, err
			}

			if !h.IsZero() {
				state = '#'
			}
		}
		if state != '_' {
			nz += 1
		}
		summ += string(state)
	}
	return summ, nz, nil
}

func (h *HistoryArchiveState) Buckets() ([]Hash, error) {
	r := []Hash{}
	for _, b := range h.CurrentBuckets {
		for _, bs := range []string{
			b.Curr, b.Snap, b.Next.Output,
		} {
			// Ignore empty values
			if bs == "" {
				continue
			}

			h, err := DecodeHash(bs)
			if err != nil {
				return r, err
			}
			if !h.IsZero() {
				r = append(r, h)
			}
		}
	}
	return r, nil
}

// BucketListHash calculates the hash of bucket list in the HistoryArchiveState.
// This can be later compared with LedgerHeader.BucketListHash of the checkpoint
// ledger to ensure data in history archive has not been changed by a malicious
// actor.
// Warning: Ledger header should be fetched from a trusted (!) diamcircle-core
// instead of ex. history archives!
func (h *HistoryArchiveState) BucketListHash() (xdr.Hash, error) {
	total := []byte{}

	for i, b := range h.CurrentBuckets {
		curr, err := hex.DecodeString(b.Curr)
		if err != nil {
			return xdr.Hash{}, errors.Wrap(err, fmt.Sprintf("Error decoding hex of %d.curr", i))
		}
		snap, err := hex.DecodeString(b.Snap)
		if err != nil {
			return xdr.Hash{}, errors.Wrap(err, fmt.Sprintf("Error decoding hex of %d.snap", i))
		}
		both := append(curr, snap...)
		bothHash := sha256.Sum256(both)
		total = append(total, bothHash[:]...)
	}

	return sha256.Sum256(total), nil
}

func (h *HistoryArchiveState) Range() Range {
	return Range{Low: 63, High: h.CurrentLedger}
}
