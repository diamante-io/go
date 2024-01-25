package xdr

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLedgerSequence(t *testing.T) {
	l := LedgerCloseMeta{
		V0: &LedgerCloseMetaV0{
			LedgerHeader: LedgerHeaderHistoryEntry{
				Header: LedgerHeader{LedgerSeq: 23},
			},
		},
	}
	assert.Equal(t, uint32(23), l.LedgerSequence())
}
