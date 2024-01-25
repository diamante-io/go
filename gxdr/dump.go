package gxdr

import (
	"bytes"
	"encoding"

	goxdr "github.com/xdrpp/goxdr/xdr"
)

// Dump serializes the given goxdr value into binary.
func Dump(v goxdr.XdrType) []byte {
	var buf bytes.Buffer
	writer := goxdr.XdrOut{Out: &buf}
	writer.Marshal("", v)
	return buf.Bytes()
}

// Convert serializes the given goxdr value into another destination value
// which supports binary unmarshalling.
//
// This function can be used to convert github.com/xdrpp/goxdr/xdr values into
// equivalent https://go-xdr values.
func Convert(src goxdr.XdrType, dest encoding.BinaryUnmarshaler) error {
	return dest.UnmarshalBinary(Dump(src))
}
