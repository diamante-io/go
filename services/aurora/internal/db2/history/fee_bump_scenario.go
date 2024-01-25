package history

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/guregu/null"

	"github.com/diamcircle/go/ingest"
	"github.com/diamcircle/go/network"
	"github.com/diamcircle/go/services/aurora/internal/test"
	"github.com/diamcircle/go/services/aurora/internal/toid"
	"github.com/diamcircle/go/xdr"
	"github.com/stretchr/testify/assert"
)

func ledgerToMap(ledger Ledger) map[string]interface{} {
	return map[string]interface{}{
		"importer_version":             ledger.ImporterVersion,
		"id":                           ledger.TotalOrderID.ID,
		"sequence":                     ledger.Sequence,
		"ledger_hash":                  ledger.LedgerHash,
		"previous_ledger_hash":         ledger.PreviousLedgerHash,
		"total_coins":                  ledger.TotalCoins,
		"fee_pool":                     ledger.FeePool,
		"base_fee":                     ledger.BaseFee,
		"base_reserve":                 ledger.BaseReserve,
		"max_tx_set_size":              ledger.MaxTxSetSize,
		"closed_at":                    ledger.ClosedAt,
		"created_at":                   ledger.CreatedAt,
		"updated_at":                   ledger.UpdatedAt,
		"transaction_count":            ledger.SuccessfulTransactionCount,
		"successful_transaction_count": ledger.SuccessfulTransactionCount,
		"failed_transaction_count":     ledger.FailedTransactionCount,
		"operation_count":              ledger.OperationCount,
		"protocol_version":             ledger.ProtocolVersion,
		"ledger_header":                ledger.LedgerHeaderXDR,
	}
}

type testTransaction struct {
	index         uint32
	envelopeXDR   string
	resultXDR     string
	feeChangesXDR string
	metaXDR       string
	hash          string
}

func buildLedgerTransaction(t *testing.T, tx testTransaction) ingest.LedgerTransaction {
	transaction := ingest.LedgerTransaction{
		Index:      tx.index,
		Envelope:   xdr.TransactionEnvelope{},
		Result:     xdr.TransactionResultPair{},
		FeeChanges: xdr.LedgerEntryChanges{},
		UnsafeMeta: xdr.TransactionMeta{},
	}

	tt := assert.New(t)

	err := xdr.SafeUnmarshalBase64(tx.envelopeXDR, &transaction.Envelope)
	tt.NoError(err)
	err = xdr.SafeUnmarshalBase64(tx.resultXDR, &transaction.Result.Result)
	tt.NoError(err)
	err = xdr.SafeUnmarshalBase64(tx.metaXDR, &transaction.UnsafeMeta)
	tt.NoError(err)
	err = xdr.SafeUnmarshalBase64(tx.feeChangesXDR, &transaction.FeeChanges)
	tt.NoError(err)

	_, err = hex.Decode(transaction.Result.TransactionHash[:], []byte(tx.hash))
	tt.NoError(err)

	return transaction
}

// FeeBumpFixture contains the data inserted into the database
// when running FeeBumpScenario
type FeeBumpFixture struct {
	Ledger            Ledger
	Envelope          xdr.TransactionEnvelope
	Transaction       Transaction
	NormalTransaction Transaction
	OuterHash         string
	InnerHash         string
}

// FeeBumpScenario creates a ledger containing a fee bump transaction,
// an operation, and an effect
func FeeBumpScenario(tt *test.T, q *Q, successful bool) FeeBumpFixture {
	fixture := FeeBumpFixture{}
	sequence := uint32(123)
	fixture.Ledger = Ledger{
		Sequence:                   int32(sequence),
		LedgerHash:                 "4db1e4f145e9ee75162040d26284795e0697e2e84084624e7c6c723ebbf80118",
		PreviousLedgerHash:         null.NewString("4b0b8bace3b2438b2404776ce57643966855487ba6384724a3c664c7aa4cd9e4", true),
		TotalOrderID:               TotalOrderID{toid.New(int32(69859), 0, 0).ToInt64()},
		ImporterVersion:            321,
		TransactionCount:           12,
		SuccessfulTransactionCount: new(int32),
		FailedTransactionCount:     new(int32),
		OperationCount:             23,
		TotalCoins:                 23451,
		FeePool:                    213,
		BaseReserve:                687,
		MaxTxSetSize:               345,
		ProtocolVersion:            12,
		BaseFee:                    100,
		ClosedAt:                   time.Now().UTC().Truncate(time.Second),
		LedgerHeaderXDR:            null.NewString("temp", true),
	}
	*fixture.Ledger.SuccessfulTransactionCount = 1
	*fixture.Ledger.FailedTransactionCount = 0
	_, err := q.Exec(context.Background(), sq.Insert("history_ledgers").SetMap(ledgerToMap(fixture.Ledger)))
	tt.Assert.NoError(err)

	fixture.Envelope = xdr.TransactionEnvelope{
		Type: xdr.EnvelopeTypeEnvelopeTypeTxFeeBump,
		FeeBump: &xdr.FeeBumpTransactionEnvelope{
			Tx: xdr.FeeBumpTransaction{
				FeeSource: xdr.MuxedAccount{
					Type:    xdr.CryptoKeyTypeKeyTypeEd25519,
					Ed25519: &xdr.Uint256{2, 2, 2},
				},
				Fee: 776,
				InnerTx: xdr.FeeBumpTransactionInnerTx{
					Type: xdr.EnvelopeTypeEnvelopeTypeTx,
					V1: &xdr.TransactionV1Envelope{
						Tx: xdr.Transaction{
							SourceAccount: xdr.MuxedAccount{
								Type: xdr.CryptoKeyTypeKeyTypeEd25519,
								Ed25519: &xdr.Uint256{
									3, 3, 3,
								},
							},
							Fee: 99,
							Memo: xdr.Memo{
								Type: xdr.MemoTypeMemoNone,
							},
							SeqNum: 97,
							TimeBounds: &xdr.TimeBounds{
								MinTime: 2,
								MaxTime: 4,
							},
							Operations: []xdr.Operation{
								{
									Body: xdr.OperationBody{
										Type: xdr.OperationTypeBumpSequence,
										BumpSequenceOp: &xdr.BumpSequenceOp{
											BumpTo: 98,
										},
									},
								},
							},
						},
						Signatures: []xdr.DecoratedSignature{
							{
								Hint:      xdr.SignatureHint{2, 2, 2, 2},
								Signature: xdr.Signature{20, 20, 20},
							},
						},
					},
				},
			},
			Signatures: []xdr.DecoratedSignature{
				{
					Hint:      xdr.SignatureHint{3, 3, 3, 3},
					Signature: xdr.Signature{30, 30, 30},
				},
			},
		},
	}
	envelopeXDR, err := xdr.MarshalBase64(fixture.Envelope)
	tt.Assert.NoError(err)

	innerHash, err := network.HashTransaction(
		fixture.Envelope.FeeBump.Tx.InnerTx.V1.Tx,
		"Test SDF Network ; September 2015",
	)
	tt.Assert.NoError(err)
	fixture.InnerHash = hex.EncodeToString(innerHash[:])

	outerHash, err := network.HashFeeBumpTransaction(
		fixture.Envelope.FeeBump.Tx,
		"Test SDF Network ; September 2015",
	)
	tt.Assert.NoError(err)
	fixture.OuterHash = hex.EncodeToString(outerHash[:])

	tt.Assert.NotEqual(fixture.InnerHash, fixture.OuterHash)

	resultPair := xdr.TransactionResultPair{
		TransactionHash: xdr.Hash(outerHash),
		Result: xdr.TransactionResult{
			FeeCharged: 123,
			Result: xdr.TransactionResultResult{
				Code: xdr.TransactionResultCodeTxFeeBumpInnerSuccess,
				InnerResultPair: &xdr.InnerTransactionResultPair{
					TransactionHash: xdr.Hash(innerHash),
					Result: xdr.InnerTransactionResult{
						Result: xdr.InnerTransactionResultResult{
							Code: xdr.TransactionResultCodeTxSuccess,
							Results: &[]xdr.OperationResult{
								{
									Tr: &xdr.OperationResultTr{
										Type: xdr.OperationTypeBumpSequence,
										BumpSeqResult: &xdr.BumpSequenceResult{
											Code: xdr.BumpSequenceResultCodeBumpSequenceSuccess,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	if !successful {
		resultPair.Result.Result.Code = xdr.TransactionResultCodeTxFeeBumpInnerFailed
		resultPair.Result.Result.InnerResultPair.Result.Result.Code = xdr.TransactionResultCodeTxBadAuth
	}

	resultXDR, err := xdr.MarshalBase64(resultPair.Result)
	tt.Assert.NoError(err)

	feeBumpTransaction := buildLedgerTransaction(tt.T, testTransaction{
		index:         1,
		envelopeXDR:   envelopeXDR,
		resultXDR:     resultXDR,
		feeChangesXDR: "AAAAAA==",
		metaXDR:       "AAAAAQAAAAAAAAAA",
		hash:          fixture.OuterHash,
	})
	normalTransaction := buildLedgerTransaction(tt.T, testTransaction{
		index:         2,
		envelopeXDR:   "AAAAACiSTRmpH6bHC6Ekna5e82oiGY5vKDEEUgkq9CB//t+rAAAAyAEXUhsAADDRAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAABAAAACXRlc3QgbWVtbwAAAAAAAAEAAAAAAAAACwEXUhsAAFfhAAAAAAAAAAA=",
		resultXDR:     "AAAAAAAAASwAAAAAAAAAAwAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAFAAAAAAAAAAA=",
		feeChangesXDR: "AAAAAA==",
		metaXDR:       "AAAAAQAAAAAAAAAA",
		hash:          "edba3051b2f2d9b713e8a08709d631eccb72c59864ff3c564c68792271bb24a7",
	})
	ctx := context.Background()
	insertBuilder := q.NewTransactionBatchInsertBuilder(2)
	// include both fee bump and normal transaction in the same batch
	// to make sure both kinds of transactions can be inserted using a single exec statement
	tt.Assert.NoError(insertBuilder.Add(ctx, feeBumpTransaction, sequence))
	tt.Assert.NoError(insertBuilder.Add(ctx, normalTransaction, sequence))
	tt.Assert.NoError(insertBuilder.Exec(ctx))

	account := fixture.Envelope.SourceAccount().ToAccountId()
	feeBumpAccount := fixture.Envelope.FeeBumpAccount().ToAccountId()

	opBuilder := q.NewOperationBatchInsertBuilder(1)
	details, err := json.Marshal(map[string]string{
		"bump_to": "98",
	})
	tt.Assert.NoError(err)

	tt.Assert.NoError(opBuilder.Add(
		ctx,
		toid.New(fixture.Ledger.Sequence, 1, 1).ToInt64(),
		toid.New(fixture.Ledger.Sequence, 1, 0).ToInt64(),
		1,
		xdr.OperationTypeBumpSequence,
		details,
		account.Address(),
		null.String{},
	))
	tt.Assert.NoError(opBuilder.Exec(ctx))

	effectBuilder := q.NewEffectBatchInsertBuilder(2)
	details, err = json.Marshal(map[string]interface{}{"new_seq": 98})
	tt.Assert.NoError(err)

	accounIDs, err := q.CreateAccounts(ctx, []string{account.Address()}, 1)
	tt.Assert.NoError(err)

	err = effectBuilder.Add(
		ctx,
		accounIDs[account.Address()],
		null.String{},
		toid.New(fixture.Ledger.Sequence, 1, 1).ToInt64(),
		1,
		EffectSequenceBumped,
		details,
	)
	tt.Assert.NoError(err)
	tt.Assert.NoError(effectBuilder.Exec(ctx))

	fixture.Transaction = Transaction{
		TransactionWithoutLedger: TransactionWithoutLedger{
			TotalOrderID:         TotalOrderID{528280981504},
			TransactionHash:      fixture.OuterHash,
			LedgerSequence:       fixture.Ledger.Sequence,
			ApplicationOrder:     1,
			Account:              account.Address(),
			AccountSequence:      "97",
			MaxFee:               int64(fixture.Envelope.Fee()),
			FeeCharged:           int64(resultPair.Result.FeeCharged),
			OperationCount:       1,
			TxEnvelope:           envelopeXDR,
			TxResult:             resultXDR,
			TxFeeMeta:            "AAAAAA==",
			TxMeta:               "AAAAAQAAAAAAAAAA",
			MemoType:             "none",
			Memo:                 null.NewString("", false),
			TimeBounds:           TimeBounds{Lower: null.IntFrom(2), Upper: null.IntFrom(4)},
			Signatures:           signatures(fixture.Envelope.FeeBumpSignatures()),
			InnerSignatures:      signatures(fixture.Envelope.Signatures()),
			Successful:           successful,
			NewMaxFee:            null.IntFrom(int64(fixture.Envelope.FeeBumpFee())),
			InnerTransactionHash: null.StringFrom(fixture.InnerHash),
			FeeAccount:           null.StringFrom(feeBumpAccount.Address()),
		},
	}

	fixture.NormalTransaction = Transaction{
		TransactionWithoutLedger: TransactionWithoutLedger{
			TotalOrderID:     TotalOrderID{528280981504},
			TransactionHash:  "edba3051b2f2d9b713e8a08709d631eccb72c59864ff3c564c68792271bb24a7",
			LedgerSequence:   fixture.Ledger.Sequence,
			ApplicationOrder: 1,
			Account:          "GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
			AccountSequence:  "78621794419880145",
			MaxFee:           200,
			FeeCharged:       300,
			OperationCount:   1,
			TxEnvelope:       "AAAAACiSTRmpH6bHC6Ekna5e82oiGY5vKDEEUgkq9CB//t+rAAAAyAEXUhsAADDRAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAABAAAACXRlc3QgbWVtbwAAAAAAAAEAAAAAAAAACwEXUhsAAFfhAAAAAAAAAAA=",
			TxResult:         "AAAAAAAAASwAAAAAAAAAAwAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAFAAAAAAAAAAA=",
			TxFeeMeta:        "AAAAAA==",
			TxMeta:           "AAAAAQAAAAAAAAAA",
			MemoType:         "text",
			Memo:             null.NewString("test memo", true),
			Successful:       successful,
		},
	}

	results, err := q.TransactionsByIDs(ctx, fixture.Transaction.ID, fixture.NormalTransaction.ID)
	tt.Assert.NoError(err)

	fixture.Transaction.CreatedAt = results[fixture.Transaction.ID].CreatedAt
	fixture.Transaction.UpdatedAt = results[fixture.Transaction.ID].UpdatedAt
	fixture.Transaction.LedgerCloseTime = results[fixture.Transaction.ID].LedgerCloseTime

	fixture.NormalTransaction.CreatedAt = results[fixture.NormalTransaction.ID].CreatedAt
	fixture.NormalTransaction.UpdatedAt = results[fixture.NormalTransaction.ID].UpdatedAt
	fixture.NormalTransaction.LedgerCloseTime = results[fixture.NormalTransaction.ID].LedgerCloseTime

	return fixture
}
