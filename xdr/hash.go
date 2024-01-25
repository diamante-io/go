package xdr

import "encoding/hex"

func (h Hash) HexString() string {
	return hex.EncodeToString(h[:])
}
