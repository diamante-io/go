package processors

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go/ingest"
	"go/xdr"
)

func TestParticipantsForTransaction(t *testing.T) {
	var envelope xdr.TransactionEnvelope
	var meta xdr.TransactionMeta
	var feeChanges xdr.LedgerEntryChanges
	assert.NoError(
		t,
		xdr.SafeUnmarshalBase64(
			"AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAABAAAAAQAAAAAAAABkAAAAAF4L0vAAAAAAAAAAAQAAAAAAAAAAAAAAAC6N7oJcJiUzTWRDL98Bj3fVrJUB19wFvCzEHh8nn/IOAAAAAlQL5AAAAAAAAAAAAVb8BfcAAABA8CyjzEXXVTMwnZTAbHfJeq2HCFzAWkU98ds2ZXFqjXR4EiN0YDSAb/pJwXc0TjMa//SiX83UvUFSqLa8hOXICQ==",
			&envelope,
		),
	)
	assert.NoError(
		t,
		xdr.SafeUnmarshalBase64(
			"AAAAAQAAAAIAAAADAAAAAwAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4Lazp2P/nAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4Lazp2P/nAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAADAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/+cAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAADAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrFTWBucAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAADAAAAAAAAAAAuje6CXCYlM01kQy/fAY931ayVAdfcBbwsxB4fJ5/yDgAAAAJUC+QAAAAAAwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			&meta,
		),
	)
	assert.NoError(
		t,
		xdr.SafeUnmarshalBase64(
			"AAAAAgAAAAMAAAABAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnZAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAADAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/+cAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			&feeChanges,
		),
	)

	particpants, err := participantsForTransaction(
		3,
		ingest.LedgerTransaction{
			Index:      1,
			Envelope:   envelope,
			FeeChanges: feeChanges,
			UnsafeMeta: meta,
			Result: xdr.TransactionResultPair{
				Result: xdr.TransactionResult{
					Result: xdr.TransactionResultResult{
						Code: xdr.TransactionResultCodeTxSuccess,
					},
				},
			},
		},
	)
	assert.NoError(t, err)
	assert.Len(t, particpants, 2)
	assert.Contains(
		t,
		particpants,
		xdr.MustAddress("GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H"),
	)
	assert.Contains(
		t,
		particpants,
		xdr.MustAddress("GAXI33UCLQTCKM2NMRBS7XYBR535LLEVAHL5YBN4FTCB4HZHT7ZA5CVK"),
	)
}
