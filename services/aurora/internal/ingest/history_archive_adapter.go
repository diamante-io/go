package ingest

import (
	"context"

	"go/historyarchive"
	"go/ingest"
	"go/support/errors"
	"go/xdr"
)

// historyArchiveAdapter is an adapter for the historyarchive package to read from history archives
type historyArchiveAdapter struct {
	archive historyarchive.ArchiveInterface
}

type historyArchiveAdapterInterface interface {
	GetLatestLedgerSequence() (uint32, error)
	BucketListHash(sequence uint32) (xdr.Hash, error)
	GetState(ctx context.Context, sequence uint32) (ingest.ChangeReader, error)
}

// newHistoryArchiveAdapter is a constructor to make a historyArchiveAdapter
func newHistoryArchiveAdapter(archive historyarchive.ArchiveInterface) historyArchiveAdapterInterface {
	return &historyArchiveAdapter{archive: archive}
}

// GetLatestLedgerSequence returns the latest ledger sequence or an error
func (haa *historyArchiveAdapter) GetLatestLedgerSequence() (uint32, error) {
	has, err := haa.archive.GetRootHAS()
	if err != nil {
		return 0, errors.Wrap(err, "could not get root HAS")
	}

	return has.CurrentLedger, nil
}

// BucketListHash returns the bucket list hash to compare with hash in the
// ledger header fetched from diamcircle-Core.
func (haa *historyArchiveAdapter) BucketListHash(sequence uint32) (xdr.Hash, error) {
	exists, err := haa.archive.CategoryCheckpointExists("history", sequence)
	if err != nil {
		return xdr.Hash{}, errors.Wrap(err, "error checking if category checkpoint exists")
	}
	if !exists {
		return xdr.Hash{}, errors.Errorf("history checkpoint does not exist for ledger %d", sequence)
	}

	has, err := haa.archive.GetCheckpointHAS(sequence)
	if err != nil {
		return xdr.Hash{}, errors.Wrapf(err, "unable to get checkpoint HAS at ledger sequence %d", sequence)
	}

	return has.BucketListHash()
}

// GetState returns a reader with the state of the ledger at the provided sequence number.
func (haa *historyArchiveAdapter) GetState(ctx context.Context, sequence uint32) (ingest.ChangeReader, error) {
	exists, err := haa.archive.CategoryCheckpointExists("history", sequence)
	if err != nil {
		return nil, errors.Wrap(err, "error checking if category checkpoint exists")
	}
	if !exists {
		return nil, errors.Errorf("history checkpoint does not exist for ledger %d", sequence)
	}

	sr, e := ingest.NewCheckpointChangeReader(ctx, haa.archive, sequence)
	if e != nil {
		return nil, errors.Wrap(e, "could not make memory state reader")
	}

	return sr, nil
}
