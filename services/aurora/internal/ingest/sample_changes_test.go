//go:build update
// +build update

package ingest

import (
	"context"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"testing"

	"go/historyarchive"
	"go/ingest"
	"go/services/aurora/internal/ingest/processors"
	logpkg "go/support/log"
	"go/xdr"
)

const (
	sampleSize = 100
	sampleSeed = 37213
)

type sampleChangeReader struct {
	offers     []*xdr.LedgerEntry
	trustlines []*xdr.LedgerEntry
	accounts   []*xdr.LedgerEntry
	data       []*xdr.LedgerEntry
	size       int

	allAccounts map[string]*xdr.LedgerEntry
	allChanges  xdr.LedgerEntryChanges
	inner       ingest.ChangeReader
	random      *rand.Rand
	output      string
}

func newSampleChangeReader(output string, size int) (*sampleChangeReader, error) {
	archive, err := historyarchive.Connect(
		"http://history.diamcircle.org/prd/core-live/core_live_001",
		historyarchive.ConnectOptions{
			Context: context.Background(),
		},
	)
	if err != nil {
		return nil, err
	}

	historyAdapter := newHistoryArchiveAdapter(archive)
	checkpointLedger, err := historyAdapter.GetLatestLedgerSequence()
	if err != nil {
		return nil, err
	}

	inner, err := historyAdapter.GetState(context.Background(), checkpointLedger)
	if err != nil {
		return nil, err
	}
	inner = newloggingChangeReader(
		inner,
		"historyArchive",
		checkpointLedger,
		logFrequency,
		false,
	)

	r := &sampleChangeReader{
		offers:      []*xdr.LedgerEntry{},
		trustlines:  []*xdr.LedgerEntry{},
		accounts:    []*xdr.LedgerEntry{},
		data:        []*xdr.LedgerEntry{},
		inner:       inner,
		allAccounts: map[string]*xdr.LedgerEntry{},
		output:      output,
		random:      rand.New(rand.NewSource(sampleSeed)),
		size:        size,
		allChanges:  xdr.LedgerEntryChanges{},
	}
	return r, nil
}

func (r *sampleChangeReader) Read() (ingest.Change, error) {
	change, err := r.inner.Read()
	if err != nil {
		return change, err
	}

	switch change.Type {
	case xdr.LedgerEntryTypeAccount:
		r.allAccounts[change.Post.Data.Account.AccountId.Address()] = change.Post
	case xdr.LedgerEntryTypeData:
		if len(r.data) < r.size {
			r.data = append(r.data, change.Post)
		} else {
			r.data[r.random.Intn(r.size)] = change.Post
		}
	case xdr.LedgerEntryTypeOffer:
		if len(r.offers) < r.size {
			r.offers = append(r.offers, change.Post)
		} else {
			r.offers[r.random.Intn(r.size)] = change.Post
		}
	case xdr.LedgerEntryTypeTrustline:
		if len(r.trustlines) < r.size {
			r.trustlines = append(r.trustlines, change.Post)
		} else {
			r.trustlines[r.random.Intn(r.size)] = change.Post
		}
	}

	return change, nil
}

func getIssuer(asset xdr.Asset) string {
	if alphanum, ok := asset.GetAlphaNum12(); ok {
		return alphanum.Issuer.Address()
	}
	if alphanum, ok := asset.GetAlphaNum4(); ok {
		return alphanum.Issuer.Address()
	}
	return ""
}

func (r *sampleChangeReader) Close() error {
	if err := r.inner.Close(); err != nil {
		return err
	}

	for _, dataEntry := range r.data {
		address := dataEntry.Data.Data.AccountId.Address()
		if entry := r.allAccounts[address]; entry != nil {
			r.accounts = append(r.accounts, entry)
			delete(r.allAccounts, address)
		}
	}

	for _, trustlineEntry := range r.trustlines {
		address := trustlineEntry.Data.TrustLine.AccountId.Address()
		if entry := r.allAccounts[address]; entry != nil {
			r.accounts = append(r.accounts, entry)
			delete(r.allAccounts, address)
		}
	}

	for _, offerEntry := range r.offers {
		seller := offerEntry.Data.Offer.SellerId.Address()
		if entry := r.allAccounts[seller]; entry != nil {
			r.accounts = append(r.accounts, entry)
			delete(r.allAccounts, seller)
		}

		if issuer := getIssuer(offerEntry.Data.Offer.Buying); r.allAccounts[issuer] != nil {
			r.accounts = append(r.accounts, r.allAccounts[issuer])
			delete(r.allAccounts, issuer)
		}

		if issuer := getIssuer(offerEntry.Data.Offer.Selling); r.allAccounts[issuer] != nil {
			r.accounts = append(r.accounts, r.allAccounts[issuer])
			delete(r.allAccounts, issuer)
		}
	}

	extraAccounts := 0
	for _, entry := range r.allAccounts {
		if extraAccounts >= r.size {
			break
		}
		r.accounts = append(r.accounts, entry)
		extraAccounts++
	}

	for _, list := range [][]*xdr.LedgerEntry{
		r.accounts,
		r.data,
		r.offers,
		r.trustlines,
	} {
		for _, entry := range list {
			r.allChanges = append(r.allChanges, xdr.LedgerEntryChange{
				Type:  xdr.LedgerEntryChangeTypeLedgerEntryState,
				State: entry,
			})
		}
	}

	serialized, err := r.allChanges.MarshalBinary()
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(r.output, serialized, 0644); err != nil {
		return err
	}
	return nil
}

func TestUpdateSampleChanges(t *testing.T) {
	log.SetLevel(logpkg.InfoLevel)
	path := filepath.Join("testdata", "sample-changes.xdr")
	reader, err := newSampleChangeReader(path, sampleSize)
	if err != nil {
		t.Fatalf("could not create sample change reader: %v", err)
	}

	changeStats := &ingest.StatsChangeProcessor{}
	err = processors.StreamChanges(changeStats, reader)
	if err != nil {
		t.Fatalf("could not stream changes: %v", err)
	}
	err = reader.Close()
	if err != nil {
		t.Fatalf("could not close reader: %v", err)
	}

	results := changeStats.GetResults()
	log.WithFields(results.Map()).
		Info("Finished processing ledger entry changes")

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read sample file: %v", err)
	}
	entryChanges := xdr.LedgerEntryChanges{}
	if err := entryChanges.UnmarshalBinary(contents); err != nil {
		t.Fatalf("could not unmarshall sample file: %v", err)
	}

	for i, entry := range entryChanges {
		marshalledFileEntry, err := xdr.MarshalBase64(entry)
		if err != nil {
			t.Fatalf("could not marshall ledger entry change: %v", err)
		}

		marshalledSourceEntry, err := xdr.MarshalBase64(reader.allChanges[i])
		if err != nil {
			t.Fatalf("could not marshall ledger entry change: %v", err)
		}

		if marshalledFileEntry != marshalledSourceEntry {
			t.Fatalf(
				"ledger entry change from sample file '%s' does not match source '%s'",
				marshalledFileEntry,
				marshalledSourceEntry,
			)
		}
	}
}
