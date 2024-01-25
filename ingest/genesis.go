package ingest

import (
	"go/amount"
	"go/keypair"
	"go/xdr"
)

// GenesisChange returns the Change occurring at the genesis ledger (ledgerseq = 1)..
func GenesisChange(networkPassPhrase string) Change {
	masterKeyPair := keypair.Master(networkPassPhrase)

	masterAccountEntry := xdr.LedgerEntry{
		LastModifiedLedgerSeq: 1,
		Data: xdr.LedgerEntryData{
			Type: xdr.LedgerEntryTypeAccount,
			Account: &xdr.AccountEntry{
				AccountId: xdr.MustAddress(masterKeyPair.Address()),
				// 100B
				Balance:    amount.MustParse("100000000000"),
				SeqNum:     0,
				Thresholds: xdr.Thresholds{1, 0, 0, 0},
			},
		},
	}

	return Change{
		Type: masterAccountEntry.Data.Type,
		Post: &masterAccountEntry,
	}
}
