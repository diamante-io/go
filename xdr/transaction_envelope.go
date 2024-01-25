package xdr

// IsFeeBump returns true if the transaction envelope is a fee bump transaction
func (e TransactionEnvelope) IsFeeBump() bool {
	return e.Type == EnvelopeTypeEnvelopeTypeTxFeeBump
}

// FeeBumpAccount returns the account paying for the fee bump transaction
func (e TransactionEnvelope) FeeBumpAccount() MuxedAccount {
	return e.MustFeeBump().Tx.FeeSource
}

// FeeBumpFee returns the fee defined for the fee bump transaction
func (e TransactionEnvelope) FeeBumpFee() int64 {
	return int64(e.MustFeeBump().Tx.Fee)
}

// FeeBumpSignatures returns the list of signatures for the fee bump transaction
func (e TransactionEnvelope) FeeBumpSignatures() []DecoratedSignature {
	return e.MustFeeBump().Signatures
}

// SourceAccount returns the source account for the transaction
// If the transaction envelope is for a fee bump transaction, SourceAccount()
// returns the source account of the inner transaction
func (e TransactionEnvelope) SourceAccount() MuxedAccount {
	switch e.Type {
	case EnvelopeTypeEnvelopeTypeTxFeeBump:
		return e.FeeBump.Tx.InnerTx.V1.Tx.SourceAccount
	case EnvelopeTypeEnvelopeTypeTx:
		return e.V1.Tx.SourceAccount
	case EnvelopeTypeEnvelopeTypeTxV0:
		return MuxedAccount{
			Type:    CryptoKeyTypeKeyTypeEd25519,
			Ed25519: &e.V0.Tx.SourceAccountEd25519,
		}
	default:
		panic("unsupported transaction type: " + e.Type.String())
	}
}

// Fee returns the fee defined for the transaction envelope
// If the transaction envelope is for a fee bump transaction, Fee()
// returns the fee defined in the inner transaction
func (e TransactionEnvelope) Fee() uint32 {
	switch e.Type {
	case EnvelopeTypeEnvelopeTypeTxFeeBump:
		return uint32(e.FeeBump.Tx.InnerTx.V1.Tx.Fee)
	case EnvelopeTypeEnvelopeTypeTx:
		return uint32(e.V1.Tx.Fee)
	case EnvelopeTypeEnvelopeTypeTxV0:
		return uint32(e.V0.Tx.Fee)
	default:
		panic("unsupported transaction type: " + e.Type.String())
	}
}

// Signatures returns the list of signatures included in the transaction envelope
// If the transaction envelope is for a fee bump transaction, Signatures()
// returns the signatures for the inner transaction
func (e TransactionEnvelope) Signatures() []DecoratedSignature {
	switch e.Type {
	case EnvelopeTypeEnvelopeTypeTxFeeBump:
		return e.FeeBump.Tx.InnerTx.V1.Signatures
	case EnvelopeTypeEnvelopeTypeTx:
		return e.V1.Signatures
	case EnvelopeTypeEnvelopeTypeTxV0:
		return e.V0.Signatures
	default:
		panic("unsupported transaction type: " + e.Type.String())
	}
}

// SeqNum returns the sequence number set in the transaction envelope
// Note for fee bump transactions, SeqNum() returns the sequence number
// of the inner transaction
func (e TransactionEnvelope) SeqNum() int64 {
	switch e.Type {
	case EnvelopeTypeEnvelopeTypeTxFeeBump:
		return int64(e.FeeBump.Tx.InnerTx.V1.Tx.SeqNum)
	case EnvelopeTypeEnvelopeTypeTx:
		return int64(e.V1.Tx.SeqNum)
	case EnvelopeTypeEnvelopeTypeTxV0:
		return int64(e.V0.Tx.SeqNum)
	default:
		panic("unsupported transaction type: " + e.Type.String())
	}
}

// TimeBounds returns the time bounds set in the transaction envelope
// Note for fee bump transactions, TimeBounds() returns the time bounds
// of the inner transaction
func (e TransactionEnvelope) TimeBounds() *TimeBounds {
	switch e.Type {
	case EnvelopeTypeEnvelopeTypeTxFeeBump:
		return e.FeeBump.Tx.InnerTx.V1.Tx.TimeBounds
	case EnvelopeTypeEnvelopeTypeTx:
		return e.V1.Tx.TimeBounds
	case EnvelopeTypeEnvelopeTypeTxV0:
		return e.V0.Tx.TimeBounds
	default:
		panic("unsupported transaction type: " + e.Type.String())
	}
}

// Operations returns the operations set in the transaction envelope
// Note for fee bump transactions, Operations() returns the operations
// of the inner transaction
func (e TransactionEnvelope) Operations() []Operation {
	switch e.Type {
	case EnvelopeTypeEnvelopeTypeTxFeeBump:
		return e.FeeBump.Tx.InnerTx.V1.Tx.Operations
	case EnvelopeTypeEnvelopeTypeTx:
		return e.V1.Tx.Operations
	case EnvelopeTypeEnvelopeTypeTxV0:
		return e.V0.Tx.Operations
	default:
		panic("unsupported transaction type: " + e.Type.String())
	}
}

// Memo returns the memo set in the transaction envelope
// Note for fee bump transactions, Memo() returns the memo
// of the inner transaction
func (e TransactionEnvelope) Memo() Memo {
	switch e.Type {
	case EnvelopeTypeEnvelopeTypeTxFeeBump:
		return e.FeeBump.Tx.InnerTx.V1.Tx.Memo
	case EnvelopeTypeEnvelopeTypeTx:
		return e.V1.Tx.Memo
	case EnvelopeTypeEnvelopeTypeTxV0:
		return e.V0.Tx.Memo
	default:
		panic("unsupported transaction type: " + e.Type.String())
	}
}
