// Copyright 2016 diamcircle Development Foundation and contributors. Licensed
// under the Apache License, Version 2.0. See the COPYING file at the root
// of this distribution or at http://www.apache.org/licenses/LICENSE-2.0

package historyarchive

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	"go/support/errors"
	"go/xdr"
)

const hexPrefixPat = "/[0-9a-f]{2}/[0-9a-f]{2}/[0-9a-f]{2}/"
const rootHASPath = ".well-known/diamcircle-history.json"

type CommandOptions struct {
	Concurrency int
	Range       Range
	DryRun      bool
	Force       bool
	Verify      bool
	Thorough    bool
}

type ConnectOptions struct {
	Context context.Context
	// NetworkPassphrase defines the expected network of history archive. It is
	// checked when getting HAS. If network passphrase does not match, error is
	// returned.
	NetworkPassphrase string
	S3Region          string
	S3Endpoint        string
	UnsignedRequests  bool
	// CheckpointFrequency is the number of ledgers between checkpoints
	// if unset, DefaultCheckpointFrequency will be used
	CheckpointFrequency uint32
}

type Ledger struct {
	Header            xdr.LedgerHeaderHistoryEntry
	Transaction       xdr.TransactionHistoryEntry
	TransactionResult xdr.TransactionHistoryResultEntry
}

type ArchiveBackend interface {
	Exists(path string) (bool, error)
	Size(path string) (int64, error)
	GetFile(path string) (io.ReadCloser, error)
	PutFile(path string, in io.ReadCloser) error
	ListFiles(path string) (chan string, chan error)
	CanListFiles() bool
}

type ArchiveInterface interface {
	GetPathHAS(path string) (HistoryArchiveState, error)
	PutPathHAS(path string, has HistoryArchiveState, opts *CommandOptions) error
	BucketExists(bucket Hash) (bool, error)
	BucketSize(bucket Hash) (int64, error)
	CategoryCheckpointExists(cat string, chk uint32) (bool, error)
	GetLedgerHeader(chk uint32) (xdr.LedgerHeaderHistoryEntry, error)
	GetRootHAS() (HistoryArchiveState, error)
	GetLedgers(start, end uint32) (map[uint32]*Ledger, error)
	GetCheckpointHAS(chk uint32) (HistoryArchiveState, error)
	PutCheckpointHAS(chk uint32, has HistoryArchiveState, opts *CommandOptions) error
	PutRootHAS(has HistoryArchiveState, opts *CommandOptions) error
	ListBucket(dp DirPrefix) (chan string, chan error)
	ListAllBuckets() (chan string, chan error)
	ListAllBucketHashes() (chan Hash, chan error)
	ListCategoryCheckpoints(cat string, pth string) (chan uint32, chan error)
	GetXdrStreamForHash(hash Hash) (*XdrStream, error)
	GetXdrStream(pth string) (*XdrStream, error)
	GetCheckpointManager() CheckpointManager
}

var _ ArchiveInterface = &Archive{}

type Archive struct {
	networkPassphrase string

	mutex             sync.Mutex
	checkpointFiles   map[string](map[uint32]bool)
	allBuckets        map[Hash]bool
	referencedBuckets map[Hash]bool

	expectLedgerHashes      map[uint32]Hash
	actualLedgerHashes      map[uint32]Hash
	expectTxSetHashes       map[uint32]Hash
	actualTxSetHashes       map[uint32]Hash
	expectTxResultSetHashes map[uint32]Hash
	actualTxResultSetHashes map[uint32]Hash

	invalidBuckets int

	invalidLedgers      int
	invalidTxSets       int
	invalidTxResultSets int

	checkpointManager CheckpointManager

	backend ArchiveBackend
}

func (arch *Archive) GetCheckpointManager() CheckpointManager {
	return arch.checkpointManager
}

func (a *Archive) GetPathHAS(path string) (HistoryArchiveState, error) {
	var has HistoryArchiveState
	rdr, err := a.backend.GetFile(path)
	if err != nil {
		return has, err
	}
	defer rdr.Close()
	dec := json.NewDecoder(rdr)
	err = dec.Decode(&has)
	if err != nil {
		return has, err
	}

	// Compare network passphrase only when non empty. The field was added in
	// diamcircle-Core 14.1.0.
	if has.NetworkPassphrase != "" && a.networkPassphrase != "" &&
		has.NetworkPassphrase != a.networkPassphrase {
		return has, errors.Errorf(
			"Network passphrase does not match! expected=%s actual=%s",
			a.networkPassphrase,
			has.NetworkPassphrase,
		)
	}

	return has, nil
}

func (a *Archive) PutPathHAS(path string, has HistoryArchiveState, opts *CommandOptions) error {
	exists, err := a.backend.Exists(path)
	if err != nil {
		return err
	}
	if exists && !opts.Force {
		log.Printf("skipping existing " + path)
		return nil
	}
	buf, err := json.MarshalIndent(has, "", "    ")
	if err != nil {
		return err
	}
	return a.backend.PutFile(path,
		ioutil.NopCloser(bytes.NewReader(buf)))
}

func (a *Archive) BucketExists(bucket Hash) (bool, error) {
	return a.backend.Exists(BucketPath(bucket))
}

func (a *Archive) BucketSize(bucket Hash) (int64, error) {
	return a.backend.Size(BucketPath(bucket))
}

func (a *Archive) CategoryCheckpointExists(cat string, chk uint32) (bool, error) {
	return a.backend.Exists(CategoryCheckpointPath(cat, chk))
}

func (a *Archive) GetLedgerHeader(ledger uint32) (xdr.LedgerHeaderHistoryEntry, error) {
	checkpoint := ledger
	if !a.checkpointManager.IsCheckpoint(checkpoint) {
		checkpoint = a.checkpointManager.NextCheckpoint(ledger)
	}
	path := CategoryCheckpointPath("ledger", checkpoint)
	xdrStream, err := a.GetXdrStream(path)
	if err != nil {
		return xdr.LedgerHeaderHistoryEntry{}, errors.Wrap(err, "error opening ledger stream")
	}
	defer xdrStream.Close()

	for {
		var ledgerHeader xdr.LedgerHeaderHistoryEntry
		err = xdrStream.ReadOne(&ledgerHeader)
		if err != nil {
			if err == io.EOF {
				break
			}
			return ledgerHeader, errors.Wrap(err, "error reading from ledger stream")
		}

		if uint32(ledgerHeader.Header.LedgerSeq) == ledger {
			return ledgerHeader, nil
		}
	}

	return xdr.LedgerHeaderHistoryEntry{}, errors.New("ledger header not found in checkpoint")
}

func (a *Archive) GetRootHAS() (HistoryArchiveState, error) {
	return a.GetPathHAS(rootHASPath)
}

func (a *Archive) GetLedgers(start, end uint32) (map[uint32]*Ledger, error) {
	if start > end {
		return nil, errors.Errorf("range is invalid, start: %d end: %d", start, end)
	}
	startCheckpoint := a.GetCheckpointManager().GetCheckpoint(start)
	endCheckpoint := a.GetCheckpointManager().GetCheckpoint(end)
	cache := map[uint32]*Ledger{}
	for cur := startCheckpoint; cur <= endCheckpoint; cur += a.GetCheckpointManager().GetCheckpointFrequency() {
		for _, category := range []string{"ledger", "transactions", "results"} {
			if exists, err := a.CategoryCheckpointExists(category, cur); err != nil {
				return nil, errors.Wrap(err, "could not check if category checkpoint exists")
			} else if !exists {
				return nil, errors.Errorf("checkpoint %d is not published", cur)
			}

			if err := a.fetchCategory(cache, category, cur); err != nil {
				return nil, errors.Wrap(err, "could not fetch category checkpoint")
			}
		}
	}

	return cache, nil
}

func (a *Archive) fetchCategory(cache map[uint32]*Ledger, category string, checkpointSequence uint32) error {
	checkpointPath := CategoryCheckpointPath(category, checkpointSequence)
	xdrStream, err := a.GetXdrStream(checkpointPath)
	if err != nil {
		return errors.Wrapf(err, "error opening %s stream", category)
	}
	defer xdrStream.Close()

	for {
		switch category {
		case "ledger":
			var object xdr.LedgerHeaderHistoryEntry
			if err = xdrStream.ReadOne(&object); err == nil {
				entry := cache[uint32(object.Header.LedgerSeq)]
				if entry == nil {
					entry = &Ledger{}
				}
				entry.Header = object
				cache[uint32(object.Header.LedgerSeq)] = entry
			}
		case "transactions":
			var object xdr.TransactionHistoryEntry
			if err = xdrStream.ReadOne(&object); err == nil {
				entry := cache[uint32(object.LedgerSeq)]
				if entry == nil {
					entry = &Ledger{}
				}
				entry.Transaction = object
				cache[uint32(object.LedgerSeq)] = entry
			}
		case "results":
			var object xdr.TransactionHistoryResultEntry
			if err = xdrStream.ReadOne(&object); err == nil {
				entry := cache[uint32(object.LedgerSeq)]
				if entry == nil {
					entry = &Ledger{}
				}
				entry.TransactionResult = object
				cache[uint32(object.LedgerSeq)] = entry
			}
		default:
			panic("unknown category")
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return errors.Wrapf(err, "error reading from %s stream", category)
		}
	}

	return nil
}

func (a *Archive) GetCheckpointHAS(chk uint32) (HistoryArchiveState, error) {
	return a.GetPathHAS(CategoryCheckpointPath("history", chk))
}

func (a *Archive) PutCheckpointHAS(chk uint32, has HistoryArchiveState, opts *CommandOptions) error {
	return a.PutPathHAS(CategoryCheckpointPath("history", chk), has, opts)
}

func (a *Archive) PutRootHAS(has HistoryArchiveState, opts *CommandOptions) error {
	force := opts.Force
	opts.Force = true
	e := a.PutPathHAS(rootHASPath, has, opts)
	opts.Force = force
	return e
}

func (a *Archive) ListBucket(dp DirPrefix) (chan string, chan error) {
	return a.backend.ListFiles(path.Join("bucket", dp.Path()))
}

func (a *Archive) ListAllBuckets() (chan string, chan error) {
	return a.backend.ListFiles("bucket")
}

func (a *Archive) ListAllBucketHashes() (chan Hash, chan error) {
	sch, errs := a.backend.ListFiles("bucket")
	ch := make(chan Hash)
	rx := regexp.MustCompile("bucket" + hexPrefixPat + "bucket-([0-9a-f]{64})\\.xdr\\.gz$")
	errs = makeErrorPump(errs)
	go func() {
		for s := range sch {
			m := rx.FindStringSubmatch(s)
			if m != nil {
				ch <- MustDecodeHash(m[1])
			}
		}
		close(ch)
	}()
	return ch, errs
}

func (a *Archive) ListCategoryCheckpoints(cat string, pth string) (chan uint32, chan error) {
	ext := categoryExt(cat)
	rx := regexp.MustCompile(cat + hexPrefixPat + cat +
		"-([0-9a-f]{8})\\." + regexp.QuoteMeta(ext) + "$")
	sch, errs := a.backend.ListFiles(path.Join(cat, pth))
	ch := make(chan uint32)
	errs = makeErrorPump(errs)

	go func() {
		for s := range sch {
			m := rx.FindStringSubmatch(s)
			if m != nil {
				i, e := strconv.ParseUint(m[1], 16, 32)
				if e == nil {
					ch <- uint32(i)
				} else {
					errs <- errors.New("decoding checkpoint number in filename " + s)
				}
			}
		}
		close(ch)
	}()
	return ch, errs
}

func (a *Archive) GetBucketPathForHash(hash Hash) string {
	return fmt.Sprintf(
		"bucket/%s/bucket-%s.xdr.gz",
		HashPrefix(hash).Path(),
		hash.String(),
	)
}

func (a *Archive) GetXdrStreamForHash(hash Hash) (*XdrStream, error) {
	return a.GetXdrStream(a.GetBucketPathForHash(hash))
}

func (a *Archive) GetXdrStream(pth string) (*XdrStream, error) {
	if !strings.HasSuffix(pth, ".xdr.gz") {
		return nil, errors.New("File has non-.xdr.gz suffix: " + pth)
	}
	rdr, err := a.backend.GetFile(pth)
	if err != nil {
		return nil, err
	}
	return NewXdrGzStream(rdr)
}

func Connect(u string, opts ConnectOptions) (*Archive, error) {
	arch := Archive{
		networkPassphrase:       opts.NetworkPassphrase,
		checkpointFiles:         make(map[string](map[uint32]bool)),
		allBuckets:              make(map[Hash]bool),
		referencedBuckets:       make(map[Hash]bool),
		expectLedgerHashes:      make(map[uint32]Hash),
		actualLedgerHashes:      make(map[uint32]Hash),
		expectTxSetHashes:       make(map[uint32]Hash),
		actualTxSetHashes:       make(map[uint32]Hash),
		expectTxResultSetHashes: make(map[uint32]Hash),
		actualTxResultSetHashes: make(map[uint32]Hash),
		checkpointManager:       NewCheckpointManager(opts.CheckpointFrequency),
	}
	for _, cat := range Categories() {
		arch.checkpointFiles[cat] = make(map[uint32]bool)
	}

	if u == "" {
		return &arch, errors.New("URL is empty")
	}

	parsed, err := url.Parse(u)
	if err != nil {
		return &arch, err
	}

	if opts.Context == nil {
		opts.Context = context.Background()
	}

	pth := parsed.Path
	if parsed.Scheme == "s3" {
		// Inside s3, all paths start _without_ the leading /
		if len(pth) > 0 && pth[0] == '/' {
			pth = pth[1:]
		}
		arch.backend, err = makeS3Backend(parsed.Host, pth, opts)
	} else if parsed.Scheme == "file" {
		pth = path.Join(parsed.Host, pth)
		arch.backend = makeFsBackend(pth, opts)
	} else if parsed.Scheme == "http" || parsed.Scheme == "https" {
		arch.backend = makeHttpBackend(parsed, opts)
	} else if parsed.Scheme == "mock" {
		arch.backend = makeMockBackend(opts)
	} else {
		err = errors.New("unknown URL scheme: '" + parsed.Scheme + "'")
	}
	return &arch, err
}

func MustConnect(u string, opts ConnectOptions) *Archive {
	arch, err := Connect(u, opts)
	if err != nil {
		log.Fatal(err)
	}
	return arch
}
