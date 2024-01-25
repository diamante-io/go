package ingest

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/diamcircle/go/ingest/ledgerbackend"
	"github.com/diamcircle/go/network"
	"github.com/diamcircle/go/xdr"
)

const (
	feeAddress     = "GAHK7EEG2WWHVKDNT4CEQFZGKF2LGDSW2IVM4S5DP42RBW3K6BTODB4A"
	metaAddress    = "GACMZD5VJXTRLKVET72CETCYKELPNCOTTBDC6DHFEUPLG5DHEK534JQX"
	upgradeAddress = "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2"
)

func TestNewLedgerChangeReaderFails(t *testing.T) {
	ctx := context.Background()
	mock := &ledgerbackend.MockDatabaseBackend{}
	seq := uint32(123)
	mock.On("GetLedger", ctx, seq).Return(
		xdr.LedgerCloseMeta{},
		fmt.Errorf("ledger error"),
	).Once()
	_, err := NewLedgerChangeReader(ctx, mock, network.TestNetworkPassphrase, seq)
	assert.EqualError(
		t,
		err,
		"error getting ledger from the backend: ledger error",
	)
}

func TestNewLedgerChangeReaderSucceeds(t *testing.T) {
	ctx := context.Background()
	mock := &ledgerbackend.MockDatabaseBackend{}
	seq := uint32(123)

	header := xdr.LedgerHeaderHistoryEntry{
		Hash: xdr.Hash{1, 2, 3},
		Header: xdr.LedgerHeader{
			LedgerVersion: 7,
		},
	}

	mock.On("GetLedger", ctx, seq).Return(
		xdr.LedgerCloseMeta{
			V0: &xdr.LedgerCloseMetaV0{
				LedgerHeader: header,
			},
		},
		nil,
	).Once()

	reader, err := NewLedgerChangeReader(ctx, mock, network.TestNetworkPassphrase, seq)
	assert.NoError(t, err)

	assert.Equal(t, reader.GetHeader(), header)
}

func buildChange(account string, balance int64) xdr.LedgerEntryChange {
	return xdr.LedgerEntryChange{
		Type: xdr.LedgerEntryChangeTypeLedgerEntryCreated,
		Created: &xdr.LedgerEntry{
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress(account),
					Balance:   xdr.Int64(balance),
				},
			},
		},
	}
}

type balanceEntry struct {
	address string
	balance int64
}

func parseChange(change Change) balanceEntry {
	address := change.Post.Data.Account.AccountId.Address()
	balance := int64(change.Post.Data.Account.Balance)

	return balanceEntry{address, balance}
}

func assertChangesEqual(
	t *testing.T,
	ctx context.Context,
	sequence uint32,
	backend ledgerbackend.LedgerBackend,
	expected []balanceEntry,
) {
	reader, err := NewLedgerChangeReader(ctx, backend, network.TestNetworkPassphrase, sequence)
	assert.NoError(t, err)

	changes := []balanceEntry{}
	for {
		change, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		changes = append(changes, parseChange(change))
	}

	assert.Equal(t, expected, changes)
}

func TestLedgerChangeReaderOrder(t *testing.T) {
	ctx := context.Background()
	mock := &ledgerbackend.MockDatabaseBackend{}
	seq := uint32(123)

	src := xdr.MustAddress("GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON")
	firstTx := xdr.TransactionEnvelope{
		Type: xdr.EnvelopeTypeEnvelopeTypeTx,
		V1: &xdr.TransactionV1Envelope{
			Tx: xdr.Transaction{
				Fee:           1,
				SourceAccount: src.ToMuxedAccount(),
			},
		},
	}
	firstTxHash, err := network.HashTransactionInEnvelope(firstTx, network.TestNetworkPassphrase)
	assert.NoError(t, err)

	src = xdr.MustAddress("GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU")
	secondTx := xdr.TransactionEnvelope{
		Type: xdr.EnvelopeTypeEnvelopeTypeTx,
		V1: &xdr.TransactionV1Envelope{
			Tx: xdr.Transaction{
				Fee:           2,
				SourceAccount: src.ToMuxedAccount(),
			},
		},
	}
	secondTxHash, err := network.HashTransactionInEnvelope(secondTx, network.TestNetworkPassphrase)
	assert.NoError(t, err)

	ledger := xdr.LedgerCloseMeta{
		V0: &xdr.LedgerCloseMetaV0{
			LedgerHeader: xdr.LedgerHeaderHistoryEntry{Header: xdr.LedgerHeader{LedgerVersion: 10}},
			TxSet: xdr.TransactionSet{
				Txs: []xdr.TransactionEnvelope{
					secondTx,
					firstTx,
				},
			},
			TxProcessing: []xdr.TransactionResultMeta{
				{
					Result: xdr.TransactionResultPair{TransactionHash: firstTxHash},
					FeeProcessing: xdr.LedgerEntryChanges{
						buildChange(feeAddress, 100),
						buildChange(feeAddress, 200),
					},
					TxApplyProcessing: xdr.TransactionMeta{
						V: 1,
						V1: &xdr.TransactionMetaV1{
							Operations: []xdr.OperationMeta{
								{
									Changes: xdr.LedgerEntryChanges{
										buildChange(
											metaAddress,
											300,
										),
										buildChange(
											metaAddress,
											400,
										),
									},
								},
							},
						},
					},
				},
				{
					Result: xdr.TransactionResultPair{TransactionHash: secondTxHash},
					FeeProcessing: xdr.LedgerEntryChanges{
						buildChange(feeAddress, 300),
					},
					TxApplyProcessing: xdr.TransactionMeta{
						V: 2,
						V2: &xdr.TransactionMetaV2{
							TxChangesBefore: xdr.LedgerEntryChanges{
								buildChange(metaAddress, 600),
							},
							Operations: []xdr.OperationMeta{
								{
									Changes: xdr.LedgerEntryChanges{
										buildChange(metaAddress, 700),
									},
								},
							},
							TxChangesAfter: xdr.LedgerEntryChanges{
								buildChange(metaAddress, 800),
								buildChange(metaAddress, 900),
							},
						},
					},
				},
			},
			UpgradesProcessing: []xdr.UpgradeEntryMeta{
				{
					Changes: xdr.LedgerEntryChanges{
						buildChange(upgradeAddress, 2),
					},
				},
				{
					Changes: xdr.LedgerEntryChanges{
						buildChange(upgradeAddress, 3),
					},
				},
			},
		},
	}
	mock.On("GetLedger", ctx, seq).Return(ledger, nil).Once()

	assertChangesEqual(t, ctx, seq, mock, []balanceEntry{
		{feeAddress, 100},
		{feeAddress, 200},
		{feeAddress, 300},
		{metaAddress, 300},
		{metaAddress, 400},
		{metaAddress, 600},
		{metaAddress, 700},
		{metaAddress, 800},
		{metaAddress, 900},
		{upgradeAddress, 2},
		{upgradeAddress, 3},
	})
	mock.AssertExpectations(t)

	ledger.V0.LedgerHeader.Header.LedgerVersion = 8
	mock.On("GetLedger", ctx, seq).Return(ledger, nil).Once()
	_, err = NewLedgerChangeReader(ctx, mock, network.TestNetworkPassphrase, seq)
	assert.EqualError(
		t,
		err,
		"error extracting transactions from ledger close meta: TransactionMeta.V=2 is required in protocol"+
			" version older than version 10. Please process ledgers again using the latest diamcircle-core version.",
	)
	mock.AssertExpectations(t)

	ledger.V0.LedgerHeader.Header.LedgerVersion = 9
	ledger.V0.TxProcessing[0].FeeProcessing = xdr.LedgerEntryChanges{}
	ledger.V0.TxProcessing[1].FeeProcessing = xdr.LedgerEntryChanges{}
	mock.On("GetLedger", ctx, seq).Return(ledger, nil).Once()

	assertChangesEqual(t, ctx, seq, mock, []balanceEntry{
		{metaAddress, 300},
		{metaAddress, 400},
		{metaAddress, 600},
		{metaAddress, 700},
		{metaAddress, 800},
		{metaAddress, 900},
		{upgradeAddress, 2},
		{upgradeAddress, 3},
	})
	mock.AssertExpectations(t)

	ledger.V0.LedgerHeader.Header.LedgerVersion = 10
	ledger.V0.TxProcessing[0].FeeProcessing = xdr.LedgerEntryChanges{}
	ledger.V0.TxProcessing[1].FeeProcessing = xdr.LedgerEntryChanges{}
	mock.On("GetLedger", ctx, seq).Return(ledger, nil).Once()

	assertChangesEqual(t, ctx, seq, mock, []balanceEntry{
		{metaAddress, 300},
		{metaAddress, 400},
		{metaAddress, 600},
		{metaAddress, 700},
		{metaAddress, 800},
		{metaAddress, 900},
		{upgradeAddress, 2},
		{upgradeAddress, 3},
	})
	mock.AssertExpectations(t)

	ledger.V0.UpgradesProcessing = []xdr.UpgradeEntryMeta{
		{
			Changes: xdr.LedgerEntryChanges{},
		},
		{
			Changes: xdr.LedgerEntryChanges{},
		},
	}
	mock.On("GetLedger", ctx, seq).Return(ledger, nil).Once()

	assertChangesEqual(t, ctx, seq, mock, []balanceEntry{
		{metaAddress, 300},
		{metaAddress, 400},
		{metaAddress, 600},
		{metaAddress, 700},
		{metaAddress, 800},
		{metaAddress, 900},
	})
	mock.AssertExpectations(t)

	ledger.V0.TxProcessing[0].TxApplyProcessing = xdr.TransactionMeta{
		V: 1,
		V1: &xdr.TransactionMetaV1{
			Operations: []xdr.OperationMeta{},
		},
	}
	ledger.V0.TxProcessing[1].TxApplyProcessing = xdr.TransactionMeta{
		V: 1,
		V1: &xdr.TransactionMetaV1{
			Operations: []xdr.OperationMeta{},
		},
	}
	mock.On("GetLedger", ctx, seq).Return(ledger, nil).Once()

	assertChangesEqual(t, ctx, seq, mock, []balanceEntry{})
	mock.AssertExpectations(t)
}
