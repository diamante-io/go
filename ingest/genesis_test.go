package ingest

import (
	"testing"

	"github.com/diamcircle/go/xdr"
	"github.com/stretchr/testify/assert"
)

func TestGenesisLeaderStateReader(t *testing.T) {
	change := GenesisChange("Public Global Diamcircle Network ; September 2015")
	assert.Equal(t, xdr.LedgerEntryTypeAccount, change.Type)
	assert.Equal(t, xdr.Uint32(1), change.Post.LastModifiedLedgerSeq)
	account := change.Post.Data.MustAccount()
	assert.Equal(t, "GAAZI4TCR3TY5OJHCTJC2A4QSY6CJWJH5IAJTGKIN2ER7LBNVKOCCWN7", account.AccountId.Address())
	assert.Equal(t, xdr.SequenceNumber(0), account.SeqNum)
	assert.Equal(t, xdr.Int64(1000000000000000000), account.Balance)
	assert.Equal(t, xdr.Thresholds{1, 0, 0, 0}, account.Thresholds)
}
