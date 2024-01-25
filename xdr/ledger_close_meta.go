package xdr

func (l LedgerCloseMeta) LedgerSequence() uint32 {
	return uint32(l.MustV0().LedgerHeader.Header.LedgerSeq)
}

func (l LedgerCloseMeta) LedgerHash() Hash {
	return l.MustV0().LedgerHeader.Hash
}

func (l LedgerCloseMeta) PreviousLedgerHash() Hash {
	return l.MustV0().LedgerHeader.Header.PreviousLedgerHash
}

func (l LedgerCloseMeta) ProtocolVersion() uint32 {
	return uint32(l.MustV0().LedgerHeader.Header.LedgerVersion)
}

func (l LedgerCloseMeta) BucketListHash() Hash {
	return l.MustV0().LedgerHeader.Header.BucketListHash
}
