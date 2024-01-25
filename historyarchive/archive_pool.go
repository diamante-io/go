// Copyright 2021 Diamcircle Development Foundation and contributors. Licensed
// under the Apache License, Version 2.0. See the COPYING file at the root
// of this distribution or at http://www.apache.org/licenses/LICENSE-2.0

package historyarchive

import (
	"math/rand"

	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/xdr"
)

// A ArchivePool is just a collection of `ArchiveInterface`s so that we can
// distribute requests fairly throughout the pool.
type ArchivePool []ArchiveInterface

// NewArchivePool tries connecting to each of the provided history archive URLs,
// returning a pool of valid archives.
//
// If none of the archives work, this returns the error message of the last
// failed archive. Note that the errors for each individual archive are hard to
// track if there's success overall.
func NewArchivePool(archiveURLs []string, config ConnectOptions) (ArchivePool, error) {
	if len(archiveURLs) <= 0 {
		return nil, errors.New("No history archives provided")
	}

	var lastErr error = nil

	// Try connecting to all of the listed archives, but only store valid ones.
	var validArchives ArchivePool
	for _, url := range archiveURLs {
		archive, err := Connect(
			url,
			ConnectOptions{
				NetworkPassphrase:   config.NetworkPassphrase,
				CheckpointFrequency: config.CheckpointFrequency,
				Context:             config.Context,
			},
		)

		if err != nil {
			lastErr = errors.Wrapf(err, "Error connecting to history archive (%s)", url)
			continue
		}

		validArchives = append(validArchives, archive)
	}

	if len(validArchives) == 0 {
		return nil, lastErr
	}

	return validArchives, nil
}

// Ensure the pool conforms to the ArchiveInterface
var _ ArchiveInterface = ArchivePool{}

// Below are the ArchiveInterface method implementations.

func (pa ArchivePool) GetAnyArchive() ArchiveInterface {
	return pa[rand.Intn(len(pa))]
}

func (pa ArchivePool) GetPathHAS(path string) (HistoryArchiveState, error) {
	return pa.GetAnyArchive().GetPathHAS(path)
}

func (pa ArchivePool) PutPathHAS(path string, has HistoryArchiveState, opts *CommandOptions) error {
	return pa.GetAnyArchive().PutPathHAS(path, has, opts)
}

func (pa ArchivePool) BucketExists(bucket Hash) (bool, error) {
	return pa.GetAnyArchive().BucketExists(bucket)
}

func (pa ArchivePool) BucketSize(bucket Hash) (int64, error) {
	return pa.GetAnyArchive().BucketSize(bucket)
}

func (pa ArchivePool) CategoryCheckpointExists(cat string, chk uint32) (bool, error) {
	return pa.GetAnyArchive().CategoryCheckpointExists(cat, chk)
}

func (pa ArchivePool) GetLedgerHeader(chk uint32) (xdr.LedgerHeaderHistoryEntry, error) {
	return pa.GetAnyArchive().GetLedgerHeader(chk)
}

func (pa ArchivePool) GetRootHAS() (HistoryArchiveState, error) {
	return pa.GetAnyArchive().GetRootHAS()
}

func (pa ArchivePool) GetLedgers(start, end uint32) (map[uint32]*Ledger, error) {
	return pa.GetAnyArchive().GetLedgers(start, end)
}

func (pa ArchivePool) GetCheckpointHAS(chk uint32) (HistoryArchiveState, error) {
	return pa.GetAnyArchive().GetCheckpointHAS(chk)
}

func (pa ArchivePool) PutCheckpointHAS(chk uint32, has HistoryArchiveState, opts *CommandOptions) error {
	return pa.GetAnyArchive().PutCheckpointHAS(chk, has, opts)
}

func (pa ArchivePool) PutRootHAS(has HistoryArchiveState, opts *CommandOptions) error {
	return pa.GetAnyArchive().PutRootHAS(has, opts)
}

func (pa ArchivePool) ListBucket(dp DirPrefix) (chan string, chan error) {
	return pa.GetAnyArchive().ListBucket(dp)
}

func (pa ArchivePool) ListAllBuckets() (chan string, chan error) {
	return pa.GetAnyArchive().ListAllBuckets()
}

func (pa ArchivePool) ListAllBucketHashes() (chan Hash, chan error) {
	return pa.GetAnyArchive().ListAllBucketHashes()
}

func (pa ArchivePool) ListCategoryCheckpoints(cat string, pth string) (chan uint32, chan error) {
	return pa.GetAnyArchive().ListCategoryCheckpoints(cat, pth)
}

func (pa ArchivePool) GetXdrStreamForHash(hash Hash) (*XdrStream, error) {
	return pa.GetAnyArchive().GetXdrStreamForHash(hash)
}

func (pa ArchivePool) GetXdrStream(pth string) (*XdrStream, error) {
	return pa.GetAnyArchive().GetXdrStream(pth)
}

func (pa ArchivePool) GetCheckpointManager() CheckpointManager {
	return pa.GetAnyArchive().GetCheckpointManager()
}
