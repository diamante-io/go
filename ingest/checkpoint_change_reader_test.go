package ingest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"sync"
	"testing"
	"time"

	"github.com/diamcircle/go/historyarchive"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/xdr"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestSingleLedgerStateReaderTestSuite(t *testing.T) {
	suite.Run(t, new(SingleLedgerStateReaderTestSuite))
}

type SingleLedgerStateReaderTestSuite struct {
	suite.Suite
	mockArchive          *historyarchive.MockArchive
	reader               *CheckpointChangeReader
	has                  historyarchive.HistoryArchiveState
	mockBucketExistsCall *mock.Call
	mockBucketSizeCall   *mock.Call
}

func (s *SingleLedgerStateReaderTestSuite) SetupTest() {
	s.mockArchive = &historyarchive.MockArchive{}

	err := json.Unmarshal([]byte(hasExample), &s.has)
	s.Require().NoError(err)

	ledgerSeq := uint32(24123007)

	s.mockArchive.
		On("GetCheckpointHAS", ledgerSeq).
		Return(s.has, nil)

	// BucketExists should be called 21 times (11 levels, last without `snap`)
	s.mockBucketExistsCall = s.mockArchive.
		On("BucketExists", mock.AnythingOfType("historyarchive.Hash")).
		Return(true, nil).Times(21)

	// BucketSize should be called 21 times (11 levels, last without `snap`)
	s.mockBucketSizeCall = s.mockArchive.
		On("BucketSize", mock.AnythingOfType("historyarchive.Hash")).
		Return(int64(100), nil).Times(21)

	s.mockArchive.
		On("GetCheckpointManager").
		Return(historyarchive.NewCheckpointManager(
			historyarchive.DefaultCheckpointFrequency))

	s.reader, err = NewCheckpointChangeReader(
		context.Background(),
		s.mockArchive,
		ledgerSeq,
	)
	s.Require().NotNil(s.reader)
	s.Require().NoError(err)
	s.Assert().Equal(ledgerSeq, s.reader.sequence)

	// Disable hash validation. We trust historyarchive.XdrStream tests here.
	s.reader.disableBucketListHashValidation = true
}

func (s *SingleLedgerStateReaderTestSuite) TearDownTest() {
	s.mockArchive.AssertExpectations(s.T())
}

// TestSimple test reading buckets with a single live entry.
func (s *SingleLedgerStateReaderTestSuite) TestSimple() {
	curr1 := createXdrStream(
		metaEntry(11),
		entryAccount(xdr.BucketEntryTypeLiveentry, "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML", 1),
	)

	nextBucket := s.getNextBucketChannel()

	// Return curr1 stream for the first bucket...
	s.mockArchive.
		On("GetXdrStreamForHash", <-nextBucket).
		Return(curr1, nil).Once()

	// ...and empty streams for the rest of the buckets.
	for hash := range nextBucket {
		s.mockArchive.
			On("GetXdrStreamForHash", hash).
			Return(createXdrStream(), nil).Once()
	}

	var e Change
	var err error
	e, err = s.reader.Read()
	s.Require().NoError(err)

	id := e.Post.Data.MustAccount().AccountId
	s.Assert().Equal("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML", id.Address())

	_, err = s.reader.Read()
	s.Require().Equal(err, io.EOF)
}

// TestRemoved test reading buckets with a single live entry that was removed.
func (s *SingleLedgerStateReaderTestSuite) TestRemoved() {
	curr1 := createXdrStream(
		entryAccount(xdr.BucketEntryTypeDeadentry, "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML", 1),
	)

	snap1 := createXdrStream(
		entryAccount(xdr.BucketEntryTypeLiveentry, "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML", 1),
	)

	nextBucket := s.getNextBucketChannel()

	// Return curr1 and snap1 stream for the first two bucket...
	s.mockArchive.
		On("GetXdrStreamForHash", <-nextBucket).
		Return(curr1, nil).Once()

	s.mockArchive.
		On("GetXdrStreamForHash", <-nextBucket).
		Return(snap1, nil).Once()

	// ...and empty streams for the rest of the buckets.
	for hash := range nextBucket {
		s.mockArchive.
			On("GetXdrStreamForHash", hash).
			Return(createXdrStream(), nil).Once()
	}

	_, err := s.reader.Read()
	s.Require().Equal(err, io.EOF)
}

// TestConcurrentRead test concurrent reads for race conditions
func (s *SingleLedgerStateReaderTestSuite) TestConcurrentRead() {
	curr1 := createXdrStream(
		entryAccount(xdr.BucketEntryTypeDeadentry, "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML", 1),
	)

	snap1 := createXdrStream(
		entryAccount(xdr.BucketEntryTypeLiveentry, "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML", 1),
		entryAccount(xdr.BucketEntryTypeLiveentry, "GCMNSW2UZMSH3ZFRLWP6TW2TG4UX4HLSYO5HNIKUSFMLN2KFSF26JKWF", 1),
		entryAccount(xdr.BucketEntryTypeLiveentry, "GB6IPC7LIOSRY26MXHQ3QJ32MTELYAA6YFIRBXZVVGTU7AOI4KUFOQ54", 1),
		entryAccount(xdr.BucketEntryTypeLiveentry, "GCK45YKCFNIOICB4TWPCOPWLQYNUKCJVV7OMMHH55AB3DD67K4E54STO", 1),
	)

	nextBucket := s.getNextBucketChannel()

	// Return curr1 and snap1 stream for the first two bucket...
	s.mockArchive.
		On("GetXdrStreamForHash", <-nextBucket).
		Return(curr1, nil).Once()

	s.mockArchive.
		On("GetXdrStreamForHash", <-nextBucket).
		Return(snap1, nil).Once()

	// ...and empty streams for the rest of the buckets.
	for hash := range nextBucket {
		s.mockArchive.
			On("GetXdrStreamForHash", hash).
			Return(createXdrStream(), nil).Once()
	}

	// 3 live entries
	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			_, err := s.reader.Read()
			s.Assert().Nil(err)
			wg.Done()
		}()
	}

	wg.Wait()

	// Next call should return io.EOF
	_, err := s.reader.Read()
	s.Require().Equal(err, io.EOF)
}

// TestEnsureLatestLiveEntry tests if a live entry overrides an older initentry
func (s *SingleLedgerStateReaderTestSuite) TestEnsureLatestLiveEntry() {
	curr1 := createXdrStream(
		metaEntry(11),
		entryAccount(xdr.BucketEntryTypeLiveentry, "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML", 1),
		entryAccount(xdr.BucketEntryTypeInitentry, "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML", 2),
	)

	nextBucket := s.getNextBucketChannel()

	// Return curr1 stream, rest won't be read due to an error
	s.mockArchive.
		On("GetXdrStreamForHash", <-nextBucket).
		Return(curr1, nil).Once()

	// ...and empty streams for the rest of the buckets.
	for hash := range nextBucket {
		s.mockArchive.
			On("GetXdrStreamForHash", hash).
			Return(createXdrStream(), nil).Once()
	}

	entry, err := s.reader.Read()
	s.Require().Nil(err)
	// Latest entry balance is 1
	s.Assert().Equal(xdr.Int64(1), entry.Post.Data.Account.Balance)

	_, err = s.reader.Read()
	s.Require().Equal(err, io.EOF)
}

// TestMalformedProtocol11Bucket tests a buggy protocol 11 bucket (meta not the first entry)
func (s *SingleLedgerStateReaderTestSuite) TestMalformedProtocol11Bucket() {
	curr1 := createXdrStream(
		entryAccount(xdr.BucketEntryTypeLiveentry, "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML", 1),
		metaEntry(11),
	)

	nextBucket := s.getNextBucketChannel()

	// Return curr1 stream, rest won't be read due to an error
	s.mockArchive.
		On("GetXdrStreamForHash", <-nextBucket).
		Return(curr1, nil).Once()

	// Account entry
	_, err := s.reader.Read()
	s.Require().Nil(err)

	// Meta entry
	_, err = s.reader.Read()
	s.Require().NotNil(err)
	s.Assert().Equal("Error while reading from buckets: METAENTRY not the first entry (n=1) in the bucket hash '517bea4c6627a688a8ce501febd8c562e737e3d86b29689d9956217640f3c74b'", err.Error())
}

// TestMalformedProtocol11BucketNoMeta tests a buggy protocol 11 bucket (no meta entry)
func (s *SingleLedgerStateReaderTestSuite) TestMalformedProtocol11BucketNoMeta() {
	curr1 := createXdrStream(
		entryAccount(xdr.BucketEntryTypeInitentry, "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML", 1),
	)

	nextBucket := s.getNextBucketChannel()

	// Return curr1 stream, rest won't be read due to an error
	s.mockArchive.
		On("GetXdrStreamForHash", <-nextBucket).
		Return(curr1, nil).Once()

	// Init entry without meta
	_, err := s.reader.Read()
	s.Require().NotNil(err)
	s.Assert().Equal("Error while reading from buckets: Read INITENTRY from version <11 bucket: 0@517bea4c6627a688a8ce501febd8c562e737e3d86b29689d9956217640f3c74b", err.Error())
}

func TestBucketExistsTestSuite(t *testing.T) {
	suite.Run(t, new(BucketExistsTestSuite))
}

type BucketExistsTestSuite struct {
	suite.Suite
	mockArchive    *historyarchive.MockArchive
	reader         *CheckpointChangeReader
	cancel         context.CancelFunc
	expectedSleeps []time.Duration
}

func (s *BucketExistsTestSuite) SetupTest() {
	s.mockArchive = &historyarchive.MockArchive{}

	ledgerSeq := uint32(24123007)
	s.mockArchive.
		On("GetCheckpointHAS", ledgerSeq).
		Return(historyarchive.HistoryArchiveState{}, nil)

	s.mockArchive.
		On("GetCheckpointManager").
		Return(historyarchive.NewCheckpointManager(
			historyarchive.DefaultCheckpointFrequency))

	ctx, cancel := context.WithCancel(context.Background())
	var err error
	s.reader, err = NewCheckpointChangeReader(
		ctx,
		s.mockArchive,
		ledgerSeq,
	)
	s.cancel = cancel
	s.Require().NoError(err)
	s.reader.sleep = func(d time.Duration) {
		if len(s.expectedSleeps) == 0 {
			s.Assert().Fail("unexpected call to sleep()")
			return
		}
		s.Assert().Equal(s.expectedSleeps[0], d)
		s.expectedSleeps = s.expectedSleeps[1:]
	}
}

func (s *BucketExistsTestSuite) TearDownTest() {
	s.mockArchive.AssertExpectations(s.T())
}

func (s *BucketExistsTestSuite) testBucketExists(
	numErrors int, expectedSleeps []time.Duration,
) {
	for _, expected := range []bool{true, false} {
		hash := historyarchive.Hash{1, 2, 3}
		if numErrors > 0 {
			s.mockArchive.On("BucketExists", hash).
				Return(true, errors.New("transient error")).Times(numErrors)
		}
		s.mockArchive.On("BucketExists", hash).
			Return(expected, nil).Once()
		s.expectedSleeps = expectedSleeps
		exists, err := s.reader.bucketExists(hash)
		s.Assert().Equal(expected, exists)
		s.Assert().NoError(err)
		s.Assert().Empty(s.expectedSleeps)
	}
}

func (s *BucketExistsTestSuite) TestSucceedsFirstTime() {
	s.testBucketExists(0, []time.Duration{})
}

func (s *BucketExistsTestSuite) TestSucceedsSecondTime() {
	s.testBucketExists(1, []time.Duration{time.Second})
}

func (s *BucketExistsTestSuite) TestSucceedsThirdime() {
	s.testBucketExists(2, []time.Duration{time.Second, 2 * time.Second})
}

func (s *BucketExistsTestSuite) TestFailsAfterThirdTime() {
	hash := historyarchive.Hash{1, 2, 3}
	s.mockArchive.On("BucketExists", hash).
		Return(true, errors.New("transient error")).Times(4)
	s.expectedSleeps = []time.Duration{
		time.Second, 2 * time.Second, 4 * time.Second,
	}
	_, err := s.reader.bucketExists(hash)
	s.Assert().EqualError(err, "transient error")
	s.Assert().Empty(s.expectedSleeps)
}

func TestReadBucketEntryTestSuite(t *testing.T) {
	suite.Run(t, new(ReadBucketEntryTestSuite))
}

type ReadBucketEntryTestSuite struct {
	suite.Suite
	mockArchive *historyarchive.MockArchive
	reader      *CheckpointChangeReader
	cancel      context.CancelFunc
}

func (s *ReadBucketEntryTestSuite) SetupTest() {
	s.mockArchive = &historyarchive.MockArchive{}

	ledgerSeq := uint32(24123007)
	s.mockArchive.
		On("GetCheckpointHAS", ledgerSeq).
		Return(historyarchive.HistoryArchiveState{}, nil)

	s.mockArchive.
		On("GetCheckpointManager").
		Return(historyarchive.NewCheckpointManager(
			historyarchive.DefaultCheckpointFrequency))

	ctx, cancel := context.WithCancel(context.Background())
	var err error
	s.reader, err = NewCheckpointChangeReader(
		ctx,
		s.mockArchive,
		ledgerSeq,
	)
	s.cancel = cancel
	s.Require().NoError(err)
}

func (s *ReadBucketEntryTestSuite) TearDownTest() {
	s.mockArchive.AssertExpectations(s.T())
}

func (s *ReadBucketEntryTestSuite) TestNewXDRStream() {
	emptyHash := historyarchive.EmptyXdrArrayHash()
	expectedStream := createXdrStream(metaEntry(1), metaEntry(2))

	hash, ok := expectedStream.ExpectedHash()
	s.Require().NotEqual(historyarchive.Hash(hash), emptyHash)
	s.Require().False(ok)

	s.mockArchive.
		On("GetXdrStreamForHash", emptyHash).
		Return(expectedStream, nil).Once()

	stream, err := s.reader.newXDRStream(emptyHash)
	s.Require().NoError(err)
	s.Require().True(stream == expectedStream)

	hash, ok = stream.ExpectedHash()
	s.Require().Equal(historyarchive.Hash(hash), emptyHash)
	s.Require().True(ok)
}

func (s *ReadBucketEntryTestSuite) TestReadAllEntries() {
	emptyHash := historyarchive.EmptyXdrArrayHash()
	firstEntry := metaEntry(1)
	secondEntry := metaEntry(2)
	s.mockArchive.
		On("GetXdrStreamForHash", emptyHash).
		Return(createXdrStream(firstEntry, secondEntry), nil).Once()

	stream, err := s.reader.newXDRStream(emptyHash)
	s.Require().NoError(err)

	entry, err := s.reader.readBucketEntry(stream, emptyHash)
	s.Require().NoError(err)
	s.Require().Equal(entry, firstEntry)

	entry, err = s.reader.readBucketEntry(stream, emptyHash)
	s.Require().NoError(err)
	s.Require().Equal(entry, secondEntry)

	_, err = s.reader.readBucketEntry(stream, emptyHash)
	s.Require().Equal(io.EOF, err)
}

func (s *ReadBucketEntryTestSuite) TestFirstReadFailsWithContextError() {
	emptyHash := historyarchive.EmptyXdrArrayHash()
	firstEntry := metaEntry(1)
	secondEntry := metaEntry(2)
	s.mockArchive.
		On("GetXdrStreamForHash", emptyHash).
		Return(createXdrStream(firstEntry, secondEntry), nil).Once()

	stream, err := s.reader.newXDRStream(emptyHash)
	s.Require().NoError(err)
	s.cancel()

	_, err = s.reader.readBucketEntry(stream, emptyHash)
	s.Require().Equal(context.Canceled, err)
}

func (s *ReadBucketEntryTestSuite) TestSecondReadFailsWithContextError() {
	emptyHash := historyarchive.EmptyXdrArrayHash()
	firstEntry := metaEntry(1)
	secondEntry := metaEntry(2)
	s.mockArchive.
		On("GetXdrStreamForHash", emptyHash).
		Return(createXdrStream(firstEntry, secondEntry), nil).Once()

	stream, err := s.reader.newXDRStream(emptyHash)
	s.Require().NoError(err)

	entry, err := s.reader.readBucketEntry(stream, emptyHash)
	s.Require().NoError(err)
	s.Require().Equal(entry, firstEntry)
	s.cancel()

	_, err = s.reader.readBucketEntry(stream, emptyHash)
	s.Require().Equal(context.Canceled, err)
}

func (s *ReadBucketEntryTestSuite) TestReadEntryAllRetriesFail() {
	emptyHash := historyarchive.EmptyXdrArrayHash()

	for i := 0; i < 4; i++ {
		s.mockArchive.
			On("GetXdrStreamForHash", emptyHash).
			Return(createInvalidXdrStream(nil), nil).Once()
	}

	stream, err := s.reader.newXDRStream(emptyHash)
	s.Require().NoError(err)

	_, err = s.reader.readBucketEntry(stream, emptyHash)
	s.Require().EqualError(err, "Read wrong number of bytes from XDR")
}

func (s *ReadBucketEntryTestSuite) TestReadEntryRetryIgnoresProtocolCloseError() {
	emptyHash := historyarchive.EmptyXdrArrayHash()

	s.mockArchive.
		On("GetXdrStreamForHash", emptyHash).
		Return(
			createInvalidXdrStream(errors.New("stream error: stream ID 75; PROTOCOL_ERROR")),
			nil,
		).Once()

	expectedEntry := metaEntry(1)
	s.mockArchive.
		On("GetXdrStreamForHash", emptyHash).
		Return(createXdrStream(expectedEntry), nil).Once()

	stream, err := s.reader.newXDRStream(emptyHash)
	s.Require().NoError(err)

	entry, err := s.reader.readBucketEntry(stream, emptyHash)
	s.Require().NoError(err)
	s.Require().Equal(entry, expectedEntry)

	hash, ok := stream.ExpectedHash()
	s.Require().Equal(historyarchive.Hash(hash), emptyHash)
	s.Require().True(ok)

	_, err = s.reader.readBucketEntry(stream, emptyHash)
	s.Require().Equal(err, io.EOF)
}

func (s *ReadBucketEntryTestSuite) TestReadEntryRetryFailsToCreateNewStream() {
	emptyHash := historyarchive.EmptyXdrArrayHash()

	s.mockArchive.
		On("GetXdrStreamForHash", emptyHash).
		Return(createInvalidXdrStream(nil), nil).Once()

	var nilStream *historyarchive.XdrStream
	s.mockArchive.
		On("GetXdrStreamForHash", emptyHash).
		Return(nilStream, errors.New("cannot create new stream")).Times(3)

	stream, err := s.reader.newXDRStream(emptyHash)
	s.Require().NoError(err)

	_, err = s.reader.readBucketEntry(stream, emptyHash)
	s.Require().EqualError(err, "Error creating new xdr stream: cannot create new stream")
}

func (s *ReadBucketEntryTestSuite) TestReadEntryRetrySucceedsAfterFailsToCreateNewStream() {
	emptyHash := historyarchive.EmptyXdrArrayHash()

	s.mockArchive.
		On("GetXdrStreamForHash", emptyHash).
		Return(createInvalidXdrStream(nil), nil).Once()

	var nilStream *historyarchive.XdrStream
	s.mockArchive.
		On("GetXdrStreamForHash", emptyHash).
		Return(nilStream, errors.New("cannot create new stream")).Once()

	firstEntry := metaEntry(1)

	s.mockArchive.
		On("GetXdrStreamForHash", emptyHash).
		Return(createXdrStream(firstEntry), nil).Once()

	stream, err := s.reader.newXDRStream(emptyHash)
	s.Require().NoError(err)

	entry, err := s.reader.readBucketEntry(stream, emptyHash)
	s.Require().NoError(err)
	s.Require().Equal(entry, firstEntry)

	_, err = s.reader.readBucketEntry(stream, emptyHash)
	s.Require().Equal(io.EOF, err)
}

func (s *ReadBucketEntryTestSuite) TestReadEntryRetrySucceeds() {
	emptyHash := historyarchive.EmptyXdrArrayHash()

	s.mockArchive.
		On("GetXdrStreamForHash", emptyHash).
		Return(createInvalidXdrStream(nil), nil).Once()

	s.mockArchive.
		On("GetXdrStreamForHash", emptyHash).
		Return(createInvalidXdrStream(nil), nil).Once()

	expectedEntry := metaEntry(1)
	s.mockArchive.
		On("GetXdrStreamForHash", emptyHash).
		Return(createXdrStream(expectedEntry), nil).Once()

	stream, err := s.reader.newXDRStream(emptyHash)
	s.Require().NoError(err)

	entry, err := s.reader.readBucketEntry(stream, emptyHash)
	s.Require().NoError(err)
	s.Require().Equal(entry, expectedEntry)

	hash, ok := stream.ExpectedHash()
	s.Require().Equal(historyarchive.Hash(hash), emptyHash)
	s.Require().True(ok)

	_, err = s.reader.readBucketEntry(stream, emptyHash)
	s.Require().Equal(err, io.EOF)
}

func (s *ReadBucketEntryTestSuite) TestReadEntryRetrySucceedsWithDiscard() {
	emptyHash := historyarchive.EmptyXdrArrayHash()

	firstEntry := metaEntry(1)
	secondEntry := metaEntry(2)

	b := &bytes.Buffer{}
	s.Require().NoError(xdr.MarshalFramed(b, firstEntry))
	writeInvalidFrame(b)

	s.mockArchive.
		On("GetXdrStreamForHash", emptyHash).
		Return(xdrStreamFromBuffer(b), nil).Once()

	s.mockArchive.
		On("GetXdrStreamForHash", emptyHash).
		Return(createXdrStream(firstEntry, secondEntry), nil).Once()

	stream, err := s.reader.newXDRStream(emptyHash)
	s.Require().NoError(err)

	entry, err := s.reader.readBucketEntry(stream, emptyHash)
	s.Require().NoError(err)
	s.Require().Equal(entry, firstEntry)

	entry, err = s.reader.readBucketEntry(stream, emptyHash)
	s.Require().NoError(err)
	s.Require().Equal(entry, secondEntry)

	hash, ok := stream.ExpectedHash()
	s.Require().Equal(historyarchive.Hash(hash), emptyHash)
	s.Require().True(ok)

	_, err = s.reader.readBucketEntry(stream, emptyHash)
	s.Require().Equal(err, io.EOF)
}

func (s *ReadBucketEntryTestSuite) TestReadEntryRetryFailsWithDiscardError() {
	emptyHash := historyarchive.EmptyXdrArrayHash()

	firstEntry := metaEntry(1)

	b := &bytes.Buffer{}
	s.Require().NoError(xdr.MarshalFramed(b, firstEntry))
	writeInvalidFrame(b)

	s.mockArchive.
		On("GetXdrStreamForHash", emptyHash).
		Return(xdrStreamFromBuffer(b), nil).Times(4)

	b = &bytes.Buffer{}
	b.WriteString("a")

	stream, err := s.reader.newXDRStream(emptyHash)
	s.Require().NoError(err)

	entry, err := s.reader.readBucketEntry(stream, emptyHash)
	s.Require().NoError(err)
	s.Require().Equal(entry, firstEntry)

	_, err = s.reader.readBucketEntry(stream, emptyHash)
	s.Require().EqualError(err, "Error discarding from xdr stream: EOF")
}

func (s *ReadBucketEntryTestSuite) TestReadEntryRetrySucceedsAfterDiscardError() {
	emptyHash := historyarchive.EmptyXdrArrayHash()

	firstEntry := metaEntry(1)
	secondEntry := metaEntry(2)

	b := &bytes.Buffer{}
	s.Require().NoError(xdr.MarshalFramed(b, firstEntry))
	writeInvalidFrame(b)

	s.mockArchive.
		On("GetXdrStreamForHash", emptyHash).
		Return(xdrStreamFromBuffer(b), nil).Once()

	b = &bytes.Buffer{}
	b.WriteString("a")

	s.mockArchive.
		On("GetXdrStreamForHash", emptyHash).
		Return(xdrStreamFromBuffer(b), nil).Once()

	s.mockArchive.
		On("GetXdrStreamForHash", emptyHash).
		Return(createXdrStream(firstEntry, secondEntry), nil).Once()

	stream, err := s.reader.newXDRStream(emptyHash)
	s.Require().NoError(err)

	entry, err := s.reader.readBucketEntry(stream, emptyHash)
	s.Require().NoError(err)
	s.Require().Equal(entry, firstEntry)

	entry, err = s.reader.readBucketEntry(stream, emptyHash)
	s.Require().NoError(err)
	s.Require().Equal(entry, secondEntry)

	_, err = s.reader.readBucketEntry(stream, emptyHash)
	s.Require().Equal(io.EOF, err)
}

func TestCheckpointLedgersTestSuite(t *testing.T) {
	suite.Run(t, new(CheckpointLedgersTestSuite))
}

type CheckpointLedgersTestSuite struct {
	suite.Suite
}

// TestNonCheckpointLedger ensures that the reader errors on a non-checkpoint ledger
func (s *CheckpointLedgersTestSuite) TestNonCheckpointLedger() {
	mockArchive := &historyarchive.MockArchive{}
	ledgerSeq := uint32(123456)

	for _, freq := range []uint32{historyarchive.DefaultCheckpointFrequency, 5} {
		mockArchive.
			On("GetCheckpointManager").
			Return(historyarchive.NewCheckpointManager(freq))

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		_, err := NewCheckpointChangeReader(ctx, mockArchive, ledgerSeq)
		s.Require().Error(err)
	}
}

func metaEntry(version uint32) xdr.BucketEntry {
	return xdr.BucketEntry{
		Type: xdr.BucketEntryTypeMetaentry,
		MetaEntry: &xdr.BucketMetadata{
			LedgerVersion: xdr.Uint32(version),
		},
	}
}

func entryAccount(t xdr.BucketEntryType, id string, balance uint32) xdr.BucketEntry {
	switch t {
	case xdr.BucketEntryTypeLiveentry, xdr.BucketEntryTypeInitentry:
		return xdr.BucketEntry{
			Type: t,
			LiveEntry: &xdr.LedgerEntry{
				Data: xdr.LedgerEntryData{
					Type: xdr.LedgerEntryTypeAccount,
					Account: &xdr.AccountEntry{
						AccountId: xdr.MustAddress(id),
						Balance:   xdr.Int64(balance),
					},
				},
			},
		}
	case xdr.BucketEntryTypeDeadentry:
		return xdr.BucketEntry{
			Type: xdr.BucketEntryTypeDeadentry,
			DeadEntry: &xdr.LedgerKey{
				Type:    xdr.LedgerEntryTypeAccount,
				Account: &xdr.LedgerKeyAccount{xdr.MustAddress(id)},
			},
		}
	default:
		panic("Unkown entry type")
	}
}

type errCloser struct {
	io.Reader
	err error
}

func (e errCloser) Close() error { return e.err }

func createInvalidXdrStream(closeError error) *historyarchive.XdrStream {
	b := &bytes.Buffer{}
	writeInvalidFrame(b)

	return historyarchive.NewXdrStream(errCloser{b, closeError})
}

func writeInvalidFrame(b *bytes.Buffer) {
	bufferSize := b.Len()
	err := xdr.MarshalFramed(b, metaEntry(1))
	if err != nil {
		panic(err)
	}
	frameSize := b.Len() - bufferSize
	b.Truncate(bufferSize + frameSize/2)
}

func createXdrStream(entries ...xdr.BucketEntry) *historyarchive.XdrStream {
	b := &bytes.Buffer{}
	for _, e := range entries {
		err := xdr.MarshalFramed(b, e)
		if err != nil {
			panic(err)
		}
	}

	return xdrStreamFromBuffer(b)
}

func xdrStreamFromBuffer(b *bytes.Buffer) *historyarchive.XdrStream {
	return historyarchive.NewXdrStream(ioutil.NopCloser(b))
}

// getNextBucket is a helper that returns next bucket hash in the order of processing.
// This allows to write simpler test code that ensures that mocked calls are in a
// correct order.
func (s *SingleLedgerStateReaderTestSuite) getNextBucketChannel() <-chan (historyarchive.Hash) {
	// 11 levels with 2 buckets each = buffer of 22
	c := make(chan (historyarchive.Hash), 22)

	for i := 0; i < len(s.has.CurrentBuckets); i++ {
		b := s.has.CurrentBuckets[i]

		curr := historyarchive.MustDecodeHash(b.Curr)
		if !curr.IsZero() {
			c <- curr
		}

		snap := historyarchive.MustDecodeHash(b.Snap)
		if !snap.IsZero() {
			c <- snap
		}
	}

	close(c)
	return c
}

var hasExample = `{
    "version": 1,
    "server": "v11.1.0",
    "currentLedger": 24123007,
    "currentBuckets": [
        {
            "curr": "517bea4c6627a688a8ce501febd8c562e737e3d86b29689d9956217640f3c74b",
            "next": {
                "state": 0
            },
            "snap": "75c8c5540a825da61e05ae23d0b0be9d29f2bdb8fdfa550a3f3496f030f62ffd"
        },
        {
            "curr": "5bca6165dbf6832ff4550e67d0e564eca56494acfc9b7fd46c740f4d66c74609",
            "next": {
                "state": 1,
                "output": "75c8c5540a825da61e05ae23d0b0be9d29f2bdb8fdfa550a3f3496f030f62ffd"
            },
            "snap": "b6bad6183a3394087aae3d05ed393c5dcb80e35ed557e2c8935cba855f20dfa5"
        },
        {
            "curr": "56b70bb56fcb27dd05759b00b07bc3c9dc7cc6dbfc9d409cfec2a41d9fd4a1e8",
            "next": {
                "state": 1,
                "output": "cfa973ce4ba1fbdf3b5767e398a5b7b86e30461855d24b7b50dc499eb313b4d0"
            },
            "snap": "974a089a6980bf25d8ad1a6a71370bac2663e9bb14702ba90b9db657464c0b3a"
        },
        {
            "curr": "16742c8e61a4dde3b35179bedbdd7c56e67d03a5faf8973a6094c57e430322df",
            "next": {
                "state": 1,
                "output": "ef39804657a928139750e801c63d1d911532d7d126c80f151ba362f49147972e"
            },
            "snap": "b415a283c5e33d8c425cbb003a86c780f73e8d2016fb5dcc6ba1477e551a2c66"
        },
        {
            "curr": "b081e1c075c9114a6c74cf87a0767ee877f02e88e18a8bf97b8f268ff120ad0d",
            "next": {
                "state": 1,
                "output": "162b859558c7c51c6416f659dbd8d70236c75540196e5d7a5dee2a66744ebbf5"
            },
            "snap": "66f8fb3f36bbe328bbbe14151951891d455ad0fba1d19d05531226c0909a84c7"
        },
        {
            "curr": "822b766e755e83d4ad08a38e86466f47452a2d7c4702295ebd3235332db76a05",
            "next": {
                "state": 1,
                "output": "1c04dc66c3410efc535044f4250c02490627b549f99a8873e4857b2cec4d51c8"
            },
            "snap": "163a49fa560761217710f6bbbf85179514aa7714d373337dde7f200f8d6c623a"
        },
        {
            "curr": "75b77814875529876258760ed6b6f37d81b5a39183812c684b9e3014bb6b8cf6",
            "next": {
                "state": 1,
                "output": "444088f447eb7ea3d397e7098d57c4f63b66912d24c4a26a29bf1dde7a4fdecc"
            },
            "snap": "35472156c463eaf62867c9b229b92e8192e2fe40cf86e269cab65fd0045c996f"
        },
        {
            "curr": "b331675d693bdb4456f409083a1b8cbadbcef977df765ba7d4ddd787800bdc84",
            "next": {
                "state": 1,
                "output": "3d9627fa5ef81486688dc584f52445560a55496d3b961a7664b0e536655179bb"
            },
            "snap": "5a7996730755a90ea5cbd2d726a982f3f3703c3db8bc2a2217bd496b9c0cf3d1"
        },
        {
            "curr": "11f8c2f8e1cb0d47576f74d9e2fa838f5f3a37180907a24a85d0ad8b647862e4",
            "next": {
                "state": 1,
                "output": "6c0133dfd0411f9975c74d792911bb80fc1555830a943249cea6c2a80e5064d1"
            },
            "snap": "48f435285dd96511d0822f7ae1a20e28c6c28019e385313713655fc76fe3bc03"
        },
        {
            "curr": "5f351041761b45f3e725f98bb8b6713873e30ab6c8aee56ba0823d357c7ebd0d",
            "next": {
                "state": 1,
                "output": "264d3a93bc5fff47a968cc53f0f2f50297e5f9015300bbc357cfb8dec30899c6"
            },
            "snap": "4100ad3b1085bd14d1c808ece3b38db97171532d0d11ed5edd57aff0e416e06a"
        },
        {
            "curr": "a4811c9ba9505e421f0015e5fcfd9f5d204ae85b584766759e844ef85db10d47",
            "next": {
                "state": 1,
                "output": "be4ecc289998a40319be24662c88f161f5e78d4be846b083923614573aa17336"
            },
            "snap": "0000000000000000000000000000000000000000000000000000000000000000"
        }
    ]
}`
