package ingest

import (
	"context"

	"go/ingest/ledgerbackend"
	"go/keypair"
	logpkg "go/support/log"
	"go/xdr"
)

type fakeLedgerBackend struct {
	numTransactions       int
	changesPerTransaction int
}

func (fakeLedgerBackend) GetLatestLedgerSequence(ctx context.Context) (uint32, error) {
	return 1, nil
}

func (fakeLedgerBackend) PrepareRange(ctx context.Context, r ledgerbackend.Range) error {
	return nil
}

func (fakeLedgerBackend) IsPrepared(ctx context.Context, r ledgerbackend.Range) (bool, error) {
	return true, nil
}

func fakeAccount() xdr.LedgerEntryChange {
	account := keypair.MustRandom().Address()
	return xdr.LedgerEntryChange{
		Type: xdr.LedgerEntryChangeTypeLedgerEntryCreated,
		Created: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 1,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress(account),
					Balance:   xdr.Int64(100),
				},
			},
		},
	}
}

func fakeAccountData() xdr.LedgerEntryChange {
	account := keypair.MustRandom().Address()
	return xdr.LedgerEntryChange{
		Type: xdr.LedgerEntryChangeTypeLedgerEntryCreated,
		Created: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 1,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeData,
				Data: &xdr.DataEntry{
					AccountId: xdr.MustAddress(account),
					DataName:  "test-name",
					DataValue: xdr.DataValue("test"),
				},
			},
		},
	}
}

func fakeTrustline() xdr.LedgerEntryChange {
	account := keypair.MustRandom().Address()
	return xdr.LedgerEntryChange{
		Type: xdr.LedgerEntryChangeTypeLedgerEntryCreated,
		Created: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 1,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeTrustline,
				TrustLine: &xdr.TrustLineEntry{
					AccountId: xdr.MustAddress(account),
					Balance:   123,
					Asset:     xdr.MustNewCreditAsset("usd", account).ToTrustLineAsset(),
				},
			},
		},
	}
}

func fakeOffer(offerID int64) xdr.LedgerEntryChange {
	account := keypair.MustRandom().Address()
	return xdr.LedgerEntryChange{
		Type: xdr.LedgerEntryChangeTypeLedgerEntryCreated,
		Created: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 1,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeOffer,
				Offer: &xdr.OfferEntry{
					SellerId: xdr.MustAddress(account),
					OfferId:  xdr.Int64(offerID),
					Amount:   213,
					Buying:   xdr.MustNewCreditAsset("usd", account),
					Price:    xdr.Price{N: 1, D: 1},
					Selling:  xdr.MustNewCreditAsset("eur", account),
				},
			},
		},
	}
}

func (f fakeLedgerBackend) getLedgerAsync(ctx context.Context, sequence uint32) (bool, xdr.LedgerCloseMeta, error) {
	ledgerCloseMeta := xdr.LedgerCloseMeta{
		V0: &xdr.LedgerCloseMetaV0{
			LedgerHeader: xdr.LedgerHeaderHistoryEntry{
				Hash: xdr.Hash{1, 2, 3, 4, 5, 6},
				Header: xdr.LedgerHeader{
					LedgerVersion: 7,
					LedgerSeq:     xdr.Uint32(1),
				},
			},
			TxSet: xdr.TransactionSet{
				Txs: make([]xdr.TransactionEnvelope, f.numTransactions),
			},
			TxProcessing:       make([]xdr.TransactionResultMeta, f.numTransactions),
			UpgradesProcessing: []xdr.UpgradeEntryMeta{},
		},
	}

	logger := log.WithField("sequence", sequence)
	logger.Info("Creating fake ledger")
	var offers, trustlines, accounts, accountData, total int

	for i := 0; i < f.numTransactions; i++ {
		var results []xdr.OperationResult
		ledgerCloseMeta.V0.TxProcessing[i] = xdr.TransactionResultMeta{
			Result: xdr.TransactionResultPair{
				TransactionHash: xdr.Hash{1, byte(i % 256), byte((i / 256) % 256)},
				Result: xdr.TransactionResult{
					Result: xdr.TransactionResultResult{
						Code:    xdr.TransactionResultCodeTxSuccess,
						Results: &results,
					},
				},
			},
			//FeeProcessing:     nil,
			TxApplyProcessing: xdr.TransactionMeta{
				V: 1,
				V1: &xdr.TransactionMetaV1{
					Operations: []xdr.OperationMeta{},
				},
			},
		}
		aid := xdr.MustAddress("GAHK7EEG2WWHVKDNT4CEQFZGKF2LGDSW2IVM4S5DP42RBW3K6BTODB4A")
		ledgerCloseMeta.V0.TxSet.Txs[i] = xdr.TransactionEnvelope{
			Type: xdr.EnvelopeTypeEnvelopeTypeTx,
			V1: &xdr.TransactionV1Envelope{
				Tx: xdr.Transaction{
					SourceAccount: aid.ToMuxedAccount(),
				},
			},
		}
		feeChanges := xdr.LedgerEntryChanges{}
		for j := 0; j < f.changesPerTransaction; j++ {
			var change xdr.LedgerEntryChange
			switch total % 4 {
			case 0:
				change = fakeAccount()
				accounts++
			case 1:
				change = fakeAccountData()
				accountData++
			case 2:
				offers++
				change = fakeOffer(int64(offers))
			case 3:
				change = fakeTrustline()
				trustlines++
			}
			total++
			feeChanges = append(feeChanges, change)

			if total%logFrequency == 0 {
				curHeap, sysHeap := getMemStats()
				logger.WithFields(logpkg.F{
					"currentHeapSizeMB": curHeap,
					"systemHeapSizeMB":  sysHeap,
					"accounts":          accounts,
					"offfers":           offers,
					"trustlines":        trustlines,
					"accountData":       accountData,
					"totalChanges":      total,
				}).Info("Adding changes to fake ledger")
			}
		}

		ledgerCloseMeta.V0.TxProcessing[i].FeeProcessing = feeChanges
	}

	curHeap, sysHeap := getMemStats()
	logger.WithFields(logpkg.F{
		"currentHeapSizeMB": curHeap,
		"systemHeapSizeMB":  sysHeap,
		"accounts":          accounts,
		"offfers":           offers,
		"trustlines":        trustlines,
		"accountData":       accountData,
		"totalChanges":      total,
	}).Info("Finished creating fake ledger")
	return true, ledgerCloseMeta, nil
}

func (f *fakeLedgerBackend) GetLedger(ctx context.Context, sequence uint32) (xdr.LedgerCloseMeta, error) {
	_, meta, err := f.getLedgerAsync(ctx, sequence)
	return meta, err
}

func (fakeLedgerBackend) Close() error {
	return nil
}
