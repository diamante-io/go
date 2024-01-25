package xdr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func createLegacyTx() TransactionEnvelope {
	return TransactionEnvelope{
		Type: EnvelopeTypeEnvelopeTypeTxV0,
		V0: &TransactionV0Envelope{
			Tx: TransactionV0{
				SourceAccountEd25519: Uint256{1, 2, 3},
				Fee:                  99,
				Memo: Memo{
					Type: MemoTypeMemoNone,
				},
				SeqNum: 33,
				TimeBounds: &TimeBounds{
					MinTime: 1,
					MaxTime: 2,
				},
				Operations: []Operation{
					{
						Body: OperationBody{
							BumpSequenceOp: &BumpSequenceOp{
								BumpTo: 34,
							},
						},
					},
				},
			},
			Signatures: []DecoratedSignature{
				{
					Hint:      SignatureHint{1, 1, 1, 1},
					Signature: Signature{10, 10, 10},
				},
			},
		},
	}
}

func createTx() TransactionEnvelope {
	return TransactionEnvelope{
		Type: EnvelopeTypeEnvelopeTypeTx,
		V1: &TransactionV1Envelope{
			Tx: Transaction{
				SourceAccount: MuxedAccount{
					Type: CryptoKeyTypeKeyTypeEd25519,
					Ed25519: &Uint256{
						3, 3, 3,
					},
				},
				Fee: 99,
				Memo: Memo{
					Type: MemoTypeMemoHash,
					Hash: &Hash{1, 1, 1},
				},
				SeqNum: 97,
				TimeBounds: &TimeBounds{
					MinTime: 2,
					MaxTime: 4,
				},
				Operations: []Operation{
					{
						Body: OperationBody{
							BumpSequenceOp: &BumpSequenceOp{
								BumpTo: 98,
							},
						},
					},
				},
			},
			Signatures: []DecoratedSignature{
				{
					Hint:      SignatureHint{2, 2, 2, 2},
					Signature: Signature{20, 20, 20},
				},
			},
		},
	}
}

func createFeeBumpTx() TransactionEnvelope {
	return TransactionEnvelope{
		Type: EnvelopeTypeEnvelopeTypeTxFeeBump,
		FeeBump: &FeeBumpTransactionEnvelope{
			Tx: FeeBumpTransaction{
				FeeSource: MuxedAccount{
					Type:    CryptoKeyTypeKeyTypeEd25519,
					Ed25519: &Uint256{2, 2, 2},
				},
				Fee: 776,
				InnerTx: FeeBumpTransactionInnerTx{
					Type: EnvelopeTypeEnvelopeTypeTx,
					V1:   createTx().V1,
				},
			},
			Signatures: []DecoratedSignature{
				{
					Hint:      SignatureHint{3, 3, 3, 3},
					Signature: Signature{30, 30, 30},
				},
			},
		},
	}
}

func TestIsFeeBump(t *testing.T) {
	legacyTx := createLegacyTx()
	tx := createTx()
	feeBumpTx := createFeeBumpTx()

	assert.False(t, legacyTx.IsFeeBump())
	assert.False(t, tx.IsFeeBump())

	assert.True(t, feeBumpTx.IsFeeBump())
}

func TestFeeBumpAccount(t *testing.T) {
	legacyTx := createLegacyTx()
	tx := createTx()
	feeBumpTx := createFeeBumpTx()

	assert.Panics(t, func() {
		tx.FeeBumpAccount()
	})
	assert.Panics(t, func() {
		legacyTx.FeeBumpAccount()
	})

	account := feeBumpTx.FeeBumpAccount()
	assert.Equal(t, feeBumpTx.FeeBump.Tx.FeeSource, account)
}

func TestFeeBumpFee(t *testing.T) {
	legacyTx := createLegacyTx()
	tx := createTx()
	feeBumpTx := createFeeBumpTx()

	assert.Panics(t, func() {
		tx.FeeBumpFee()
	})
	assert.Panics(t, func() {
		legacyTx.FeeBumpFee()
	})

	fee := feeBumpTx.FeeBumpFee()
	assert.Equal(t, int64(feeBumpTx.FeeBump.Tx.Fee), fee)
}

func TestFeeBumpSignatures(t *testing.T) {
	legacyTx := createLegacyTx()
	tx := createTx()
	feeBumpTx := createFeeBumpTx()

	assert.Panics(t, func() {
		tx.FeeBumpSignatures()
	})
	assert.Panics(t, func() {
		legacyTx.FeeBumpSignatures()
	})

	sigs := feeBumpTx.FeeBumpSignatures()
	assert.Equal(t, feeBumpTx.FeeBump.Signatures, sigs)
}

func TestSourceAccount(t *testing.T) {
	legacyTx := createLegacyTx()
	tx := createTx()
	feeBumpTx := createFeeBumpTx()

	assert.Equal(
		t,
		CryptoKeyTypeKeyTypeEd25519,
		legacyTx.SourceAccount().Type,
	)
	assert.Equal(
		t,
		legacyTx.V0.Tx.SourceAccountEd25519,
		*legacyTx.SourceAccount().Ed25519,
	)

	assert.Equal(
		t,
		tx.V1.Tx.SourceAccount,
		tx.SourceAccount(),
	)

	assert.Equal(
		t,
		feeBumpTx.FeeBump.Tx.InnerTx.V1.Tx.SourceAccount,
		feeBumpTx.SourceAccount(),
	)
}

func TestFee(t *testing.T) {
	legacyTx := createLegacyTx()
	tx := createTx()
	feeBumpTx := createFeeBumpTx()

	assert.Equal(
		t,
		uint32(legacyTx.V0.Tx.Fee),
		legacyTx.Fee(),
	)

	assert.Equal(
		t,
		uint32(tx.V1.Tx.Fee),
		tx.Fee(),
	)

	assert.Equal(
		t,
		uint32(feeBumpTx.FeeBump.Tx.InnerTx.V1.Tx.Fee),
		feeBumpTx.Fee(),
	)
}

func TestSignatures(t *testing.T) {
	legacyTx := createLegacyTx()
	tx := createTx()
	feeBumpTx := createFeeBumpTx()

	assert.Equal(
		t,
		legacyTx.V0.Signatures,
		legacyTx.Signatures(),
	)

	assert.Equal(
		t,
		tx.V1.Signatures,
		tx.Signatures(),
	)

	assert.Equal(
		t,
		feeBumpTx.FeeBump.Tx.InnerTx.V1.Signatures,
		feeBumpTx.Signatures(),
	)
}

func TestSeqNum(t *testing.T) {
	legacyTx := createLegacyTx()
	tx := createTx()
	feeBumpTx := createFeeBumpTx()

	assert.Equal(
		t,
		int64(legacyTx.V0.Tx.SeqNum),
		legacyTx.SeqNum(),
	)

	assert.Equal(
		t,
		int64(tx.V1.Tx.SeqNum),
		tx.SeqNum(),
	)

	assert.Equal(
		t,
		int64(feeBumpTx.FeeBump.Tx.InnerTx.V1.Tx.SeqNum),
		feeBumpTx.SeqNum(),
	)
}

func TestTimeBounds(t *testing.T) {
	legacyTx := createLegacyTx()
	tx := createTx()
	feeBumpTx := createFeeBumpTx()

	assert.Equal(
		t,
		legacyTx.V0.Tx.TimeBounds,
		legacyTx.TimeBounds(),
	)

	assert.Equal(
		t,
		tx.V1.Tx.TimeBounds,
		tx.TimeBounds(),
	)

	assert.Equal(
		t,
		feeBumpTx.FeeBump.Tx.InnerTx.V1.Tx.TimeBounds,
		feeBumpTx.TimeBounds(),
	)
}

func TestOperations(t *testing.T) {
	legacyTx := createLegacyTx()
	tx := createTx()
	feeBumpTx := createFeeBumpTx()

	assert.Equal(
		t,
		legacyTx.V0.Tx.Operations,
		legacyTx.Operations(),
	)

	assert.Equal(
		t,
		tx.V1.Tx.Operations,
		tx.Operations(),
	)

	assert.Equal(
		t,
		feeBumpTx.FeeBump.Tx.InnerTx.V1.Tx.Operations,
		feeBumpTx.Operations(),
	)
}

func TestMemo(t *testing.T) {
	legacyTx := createLegacyTx()
	tx := createTx()
	feeBumpTx := createFeeBumpTx()

	assert.Equal(
		t,
		legacyTx.V0.Tx.Memo,
		legacyTx.Memo(),
	)

	assert.Equal(
		t,
		tx.V1.Tx.Memo,
		tx.Memo(),
	)

	assert.Equal(
		t,
		feeBumpTx.FeeBump.Tx.InnerTx.V1.Tx.Memo,
		feeBumpTx.Memo(),
	)
}
