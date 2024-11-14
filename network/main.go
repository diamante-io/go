// Package network contains functions that deal with diamcircle network passphrases
// and IDs.
package network

import (
	"bytes"

	"strings"

	"github.com/diamcircle/go/hash"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/xdr"
)

const (
	// PublicNetworkPassphrase is the pass phrase used for every transaction intended for the public diamcircle network
	PublicNetworkPassphrase = "Diamante MainNet; SEP 2022"
	// TestNetworkPassphrase is the pass phrase used for every transaction intended for the SDF-run test network
	TestNetworkPassphrase = "Diamante Testnet 2024" // Diamante Testnet //"Pyramid ; October 2022" //
)

// ID returns the network ID derived from the provided passphrase.  This value
// also happens to be the raw (i.e. not strkey encoded) secret key for the root
// account of the network.
func ID(passphrase string) [32]byte {
	return hash.Hash([]byte(passphrase))
}

// HashTransactionInEnvelope derives the network specific hash for the transaction
// contained in the provided envelope using the network identified by the supplied passphrase.
// The resulting hash is the value that can be signed by diamcircle secret key to
// authorize the transaction identified by the hash to diamcircle validators.
func HashTransactionInEnvelope(envelope xdr.TransactionEnvelope, passphrase string) ([32]byte, error) {
	var hash [32]byte
	var err error
	switch envelope.Type {
	case xdr.EnvelopeTypeEnvelopeTypeTx:
		hash, err = HashTransaction(envelope.V1.Tx, passphrase)
	case xdr.EnvelopeTypeEnvelopeTypeTxV0:
		hash, err = HashTransactionV0(envelope.V0.Tx, passphrase)
	case xdr.EnvelopeTypeEnvelopeTypeTxFeeBump:
		hash, err = HashFeeBumpTransaction(envelope.FeeBump.Tx, passphrase)
	default:
		err = errors.New("invalid transaction type")
	}
	return hash, err
}

// HashTransaction derives the network specific hash for the provided
// transaction using the network identified by the supplied passphrase.  The
// resulting hash is the value that can be signed by diamcircle secret key to
// authorize the transaction identified by the hash to diamcircle validators.
func HashTransaction(tx xdr.Transaction, passphrase string) ([32]byte, error) {
	taggedTx := xdr.TransactionSignaturePayloadTaggedTransaction{
		Type: xdr.EnvelopeTypeEnvelopeTypeTx,
		Tx:   &tx,
	}
	return hashTx(taggedTx, passphrase)
}

// HashFeeBumpTransaction derives the network specific hash for the provided
// fee bump transaction using the network identified by the supplied passphrase.  The
// resulting hash is the value that can be signed by diamcircle secret key to
// authorize the transaction identified by the hash to diamcircle validators.
func HashFeeBumpTransaction(tx xdr.FeeBumpTransaction, passphrase string) ([32]byte, error) {
	taggedTx := xdr.TransactionSignaturePayloadTaggedTransaction{
		Type:    xdr.EnvelopeTypeEnvelopeTypeTxFeeBump,
		FeeBump: &tx,
	}
	return hashTx(taggedTx, passphrase)
}

// HashTransactionV0 derives the network specific hash for the provided
// legacy transaction using the network identified by the supplied passphrase.  The
// resulting hash is the value that can be signed by diamcircle secret key to
// authorize the transaction identified by the hash to diamcircle validators.
func HashTransactionV0(tx xdr.TransactionV0, passphrase string) ([32]byte, error) {
	sa, err := xdr.NewMuxedAccount(xdr.CryptoKeyTypeKeyTypeEd25519, tx.SourceAccountEd25519)
	if err != nil {
		return [32]byte{}, err
	}
	v1Tx := xdr.Transaction{
		SourceAccount: sa,
		Fee:           tx.Fee,
		Memo:          tx.Memo,
		Operations:    tx.Operations,
		SeqNum:        tx.SeqNum,
		TimeBounds:    tx.TimeBounds,
	}
	return HashTransaction(v1Tx, passphrase)
}

func hashTx(
	tx xdr.TransactionSignaturePayloadTaggedTransaction,
	passphrase string,
) ([32]byte, error) {
	if strings.TrimSpace(passphrase) == "" {
		return [32]byte{}, errors.New("empty network passphrase")
	}

	var txBytes bytes.Buffer
	payload := xdr.TransactionSignaturePayload{
		NetworkId:         ID(passphrase),
		TaggedTransaction: tx,
	}

	_, err := xdr.Marshal(&txBytes, payload)
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "marshal tx failed")
	}

	return hash.Hash(txBytes.Bytes()), nil
}
