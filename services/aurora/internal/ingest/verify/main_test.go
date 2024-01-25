//lint:file-ignore U1001 Ignore all unused code, staticcheck doesn't understand testify/suite

package verify

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"go/ingest"
	"go/support/errors"
	"go/xdr"
)

func assertStateError(t *testing.T, err error, expectStateError bool) {
	_, ok := err.(ingest.StateError)
	if expectStateError {
		assert.True(t, ok, "err should be StateError")
	} else {
		assert.False(t, ok, "err should not be StateError")
	}
}

func TestStateVerifierTestSuite(t *testing.T) {
	suite.Run(t, new(StateVerifierTestSuite))
}

type StateVerifierTestSuite struct {
	suite.Suite
	verifier        *StateVerifier
	mockStateReader *ingest.MockChangeReader
}

func (s *StateVerifierTestSuite) SetupTest() {
	s.mockStateReader = &ingest.MockChangeReader{}
	s.verifier = NewStateVerifier(s.mockStateReader, nil)
}

func (s *StateVerifierTestSuite) TearDownTest() {
	s.mockStateReader.AssertExpectations(s.T())
}

func (s *StateVerifierTestSuite) TestNoEntries() {
	s.mockStateReader.On("Read").Return(ingest.Change{}, io.EOF).Once()

	keys, err := s.verifier.GetLedgerKeys(10)
	s.Assert().NoError(err)
	s.Assert().Len(keys, 0)
}

func (s *StateVerifierTestSuite) TestReturnErrorOnStateReaderError() {
	s.mockStateReader.On("Read").Return(ingest.Change{}, errors.New("Read error")).Once()

	_, err := s.verifier.GetLedgerKeys(10)
	s.Assert().EqualError(err, "Read error")
}

func (s *StateVerifierTestSuite) TestCurrentEntriesNotEmpty() {
	entry := makeAccountLedgerEntry()
	entryBase64, err := xdr.MarshalBase64(entry)
	s.Assert().NoError(err)

	ledgerKey := entry.LedgerKey()
	ledgerKeyBase64, err := xdr.MarshalBase64(ledgerKey)
	s.Assert().NoError(err)

	s.verifier.currentEntries = map[string]xdr.LedgerEntry{
		ledgerKeyBase64: entry,
	}

	_, err = s.verifier.GetLedgerKeys(10)
	s.Assert().Error(err)
	assertStateError(s.T(), err, true)
	s.Assert().EqualError(err, "Entries (1) not found locally, example: "+entryBase64)

	err = s.verifier.Verify(10)
	s.Assert().Error(err)
	assertStateError(s.T(), err, true)
	s.Assert().EqualError(err, "Entries (1) not found locally, example: "+entryBase64)
}

func (s *StateVerifierTestSuite) TestTransformFunction() {
	accountEntry := makeAccountLedgerEntry()
	s.mockStateReader.
		On("Read").
		Return(ingest.Change{
			Type: xdr.LedgerEntryTypeAccount,
			Post: &accountEntry,
		}, nil).Once()

	offerEntry := makeOfferLedgerEntry()
	s.mockStateReader.
		On("Read").
		Return(ingest.Change{
			Type: xdr.LedgerEntryTypeOffer,
			Post: &offerEntry,
		}, nil).Once()

	s.mockStateReader.On("Read").Return(ingest.Change{}, io.EOF).Once()

	s.verifier.transformFunction =
		func(entry xdr.LedgerEntry) (ignore bool, newEntry xdr.LedgerEntry) {
			// Leave Account ID only for accounts, ignore the rest
			switch entry.Data.Type {
			case xdr.LedgerEntryTypeAccount:
				accountEntry := entry.Data.Account

				return false, xdr.LedgerEntry{
					Data: xdr.LedgerEntryData{
						Type: xdr.LedgerEntryTypeAccount,
						Account: &xdr.AccountEntry{
							AccountId: accountEntry.AccountId,
						},
					},
				}
			default:
				return true, xdr.LedgerEntry{}
			}
		}

	_, err := s.verifier.GetLedgerKeys(10)
	s.Assert().NoError(err)

	// Check currentEntries
	ledgerKey, err := accountEntry.LedgerKey().MarshalBinary()
	s.Assert().NoError(err)

	// Account entry transformed and offer entry ignored
	s.Assert().Len(s.verifier.currentEntries, 1)
	s.Assert().Equal(accountEntry, s.verifier.currentEntries[string(ledgerKey)])
}

func (s *StateVerifierTestSuite) TestOnlyRequestedNumberOfKeysReturned() {
	accountEntry := makeAccountLedgerEntry()
	s.mockStateReader.
		On("Read").
		Return(ingest.Change{
			Type: xdr.LedgerEntryTypeAccount,
			Post: &accountEntry,
		}, nil).Once()

	// We don't mock Read() -> (io.Change{}, stdio.EOF) call here
	// because this would execute `stdio.EOF` code path.

	keys, err := s.verifier.GetLedgerKeys(1)
	s.Assert().NoError(err)
	s.Assert().Len(keys, 1)

	// In such case Verify() should notice that not all entries read from buckets
	err = s.verifier.Write(accountEntry)
	s.Assert().NoError(err)

	err = s.verifier.Verify(1)
	s.Assert().Error(err)
	assertStateError(s.T(), err, false)
	s.Assert().EqualError(err, "There are unread entries in state reader. Process all entries before calling Verify.")
}

func (s *StateVerifierTestSuite) TestWriteEntryNotExist() {
	entry := makeAccountLedgerEntry()
	entryBase64, err := xdr.MarshalBase64(entry)
	s.Assert().NoError(err)

	ledgerKey := entry.LedgerKey()
	ledgerKeyBase64, err := xdr.MarshalBase64(ledgerKey)
	s.Assert().NoError(err)

	err = s.verifier.Write(entry)
	s.Assert().Error(err)
	assertStateError(s.T(), err, true)
	errorMsg := fmt.Sprintf(
		"Cannot find entry in currentEntries map: %s (key = %s)",
		entryBase64,
		ledgerKeyBase64,
	)
	s.Assert().EqualError(err, errorMsg)
}

func (s *StateVerifierTestSuite) TestTransformFunctionBuggyIgnore() {
	accountEntry := makeAccountLedgerEntry()
	s.mockStateReader.
		On("Read").
		Return(ingest.Change{
			Type: xdr.LedgerEntryTypeAccount,
			Post: &accountEntry,
		}, nil).Once()

	s.verifier.transformFunction =
		func(entry xdr.LedgerEntry) (ignore bool, newEntry xdr.LedgerEntry) {
			return false, xdr.LedgerEntry{}
		}

	keys, err := s.verifier.GetLedgerKeys(1)
	s.Assert().NoError(err)
	s.Assert().Len(keys, 1)

	// Check the behaviour of transformFunction to code path to test.
	s.verifier.transformFunction =
		func(entry xdr.LedgerEntry) (ignore bool, newEntry xdr.LedgerEntry) {
			return true, xdr.LedgerEntry{}
		}

	entryBase64, err := xdr.MarshalBase64(accountEntry)
	s.Assert().NoError(err)
	errorMsg := fmt.Sprintf(
		"Entry ignored in GetEntries but not ignored in Write: %s. Possibly transformFunction is buggy.",
		entryBase64,
	)
	err = s.verifier.Write(accountEntry)
	s.Assert().EqualError(err, errorMsg)
}

func (s *StateVerifierTestSuite) TestActualExpectedEntryNotEqualWrite() {
	expectedEntry := makeAccountLedgerEntry()
	s.mockStateReader.
		On("Read").
		Return(ingest.Change{
			Type: xdr.LedgerEntryTypeAccount,
			Post: &expectedEntry,
		}, nil).Once()

	keys, err := s.verifier.GetLedgerKeys(1)
	s.Assert().NoError(err)
	s.Assert().Len(keys, 1)

	actualEntry := makeAccountLedgerEntry()
	actualEntry.Data.Account.Thresholds = [4]byte{1, 1, 1, 0}
	actualEntry.Normalize()

	expectedEntryBase64, err := xdr.MarshalBase64(expectedEntry)
	s.Assert().NoError(err)
	actualEntryBase64, err := xdr.MarshalBase64(actualEntry)
	s.Assert().NoError(err)

	errorMsg := fmt.Sprintf(
		"Entry does not match the fetched entry. Expected: %s (pretransform = %s), actual: %s",
		expectedEntryBase64,
		expectedEntryBase64,
		actualEntryBase64,
	)
	err = s.verifier.Write(actualEntry)
	s.Assert().Error(err)
	assertStateError(s.T(), err, true)
	s.Assert().EqualError(err, errorMsg)
}

func (s *StateVerifierTestSuite) TestVerifyCountersMatch() {
	accountEntry := makeAccountLedgerEntry()
	s.mockStateReader.
		On("Read").
		Return(ingest.Change{
			Type: xdr.LedgerEntryTypeAccount,
			Post: &accountEntry,
		}, nil).Once()

	s.mockStateReader.On("Read").Return(ingest.Change{}, io.EOF).Once()

	keys, err := s.verifier.GetLedgerKeys(2)
	s.Assert().NoError(err)
	s.Assert().Len(keys, 1)

	err = s.verifier.Write(accountEntry)
	s.Assert().NoError(err)

	err = s.verifier.Verify(10)
	s.Assert().Error(err)
	assertStateError(s.T(), err, true)
	errorMsg := fmt.Sprintf(
		"Number of entries read using GetEntries (%d) does not match number of entries in your storage (%d).",
		1,
		10,
	)
	s.Assert().EqualError(err, errorMsg)

	err = s.verifier.Verify(1)
	s.Assert().NoError(err)
}

func makeAccountLedgerEntry() xdr.LedgerEntry {
	entry := xdr.LedgerEntry{
		Data: xdr.LedgerEntryData{
			Type: xdr.LedgerEntryTypeAccount,
			Account: &xdr.AccountEntry{
				AccountId:  xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
				Thresholds: [4]byte{1, 1, 1, 1},
			},
		},
	}
	entry.Normalize()
	return entry
}

func makeOfferLedgerEntry() xdr.LedgerEntry {
	entry := xdr.LedgerEntry{
		Data: xdr.LedgerEntryData{
			Type: xdr.LedgerEntryTypeOffer,
			Offer: &xdr.OfferEntry{
				SellerId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
			},
		},
	}
	entry.Normalize()
	return entry
}
