package ledgerbackend

import (
	"context"

	"go/xdr"
)

// LedgerBackend represents the interface to a ledger data store.
type LedgerBackend interface {
	// GetLatestLedgerSequence returns the sequence of the latest ledger available
	// in the backend.
	GetLatestLedgerSequence(ctx context.Context) (sequence uint32, err error)
	// GetLedger will block until the ledger is available.
	GetLedger(ctx context.Context, sequence uint32) (xdr.LedgerCloseMeta, error)
	// PrepareRange prepares the given range (including from and to) to be loaded.
	// Some backends (like captive diamcircle-core) need to initalize data to be
	// able to stream ledgers. Blocks until the first ledger is available.
	PrepareRange(ctx context.Context, ledgerRange Range) error
	// IsPrepared returns true if a given ledgerRange is prepared.
	IsPrepared(ctx context.Context, ledgerRange Range) (bool, error)
	Close() error
}

// session is the interface needed to access a persistent database session.
// TODO can't use this until we add Close() to the existing db.Session object
type session interface {
	GetRaw(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectRaw(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Close() error
}

// ledgerHeaderHistory is a helper struct used to unmarshall header fields from a diamcircle-core DB.
type ledgerHeaderHistory struct {
	Hash   xdr.Hash         `db:"ledgerhash"`
	Header xdr.LedgerHeader `db:"data"`
}

// ledgerHeader holds a row of data from the diamcircle-core `ledgerheaders` table.
type ledgerHeader struct {
	LedgerHash     string           `db:"ledgerhash"`
	PrevHash       string           `db:"prevhash"`
	BucketListHash string           `db:"bucketlisthash"`
	CloseTime      int64            `db:"closetime"`
	LedgerSeq      uint32           `db:"ledgerseq"`
	Data           xdr.LedgerHeader `db:"data"`
}

// txHistory holds a row of data from the diamcircle-core `txhistory` table.
type txHistory struct {
	TXID      string                    `db:"txid"`
	LedgerSeq uint32                    `db:"ledgerseq"`
	TXIndex   uint32                    `db:"txindex"`
	TXBody    xdr.TransactionEnvelope   `db:"txbody"`
	TXResult  xdr.TransactionResultPair `db:"txresult"`
	TXMeta    xdr.TransactionMeta       `db:"txmeta"`
}

// txFeeHistory holds a row of data from the diamcircle-core `txfeehistory` table.
type txFeeHistory struct {
	TXID      string                 `db:"txid"`
	LedgerSeq uint32                 `db:"ledgerseq"`
	TXIndex   uint32                 `db:"txindex"`
	TXChanges xdr.LedgerEntryChanges `db:"txchanges"`
}

// scpHistory holds a row of data from the diamcircle-core `scphistory` table.
// type scpHistory struct {
// 	NodeID    string `db:"nodeid"`
// 	LedgerSeq uint32 `db:"ledgerseq"`
// 	Envelope  string `db:"envelope"`
// }

// upgradeHistory holds a row of data from the diamcircle-core `upgradehistory` table.
type upgradeHistory struct {
	LedgerSeq    uint32                 `db:"ledgerseq"`
	UpgradeIndex uint32                 `db:"upgradeindex"`
	Upgrade      xdr.LedgerUpgrade      `db:"upgrade"`
	Changes      xdr.LedgerEntryChanges `db:"changes"`
}
