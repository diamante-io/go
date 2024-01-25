//lint:file-ignore S1005 The issue should be fixed in xdrgen. Unfortunately, there's no way to ignore a single file in staticcheck.
//lint:file-ignore U1000 fmtTest is not needed anywhere, should be removed in xdrgen.
// Package xdr is generated from:
//
//  xdr/diamcircle-SCP.x
//  xdr/diamcircle-ledger-entries.x
//  xdr/diamcircle-ledger.x
//  xdr/diamcircle-overlay.x
//  xdr/diamcircle-transaction.x
//  xdr/diamcircle-types.x
//
// DO NOT EDIT or your changes may be overwritten
package xdr

import (
	"bytes"
	"encoding"
	"fmt"
	"io"

	"github.com/diamcircle/go-xdr/xdr3"
)

type xdrType interface {
	xdrType()
}

type decoderFrom interface {
	DecodeFrom(d *xdr.Decoder) (int, error)
}

// Unmarshal reads an xdr element from `r` into `v`.
func Unmarshal(r io.Reader, v interface{}) (int, error) {
	if decodable, ok := v.(decoderFrom); ok {
		d := xdr.NewDecoder(r)
		return decodable.DecodeFrom(d)
	}
	// delegate to xdr package's Unmarshal
	return xdr.Unmarshal(r, v)
}

// Marshal writes an xdr element `v` into `w`.
func Marshal(w io.Writer, v interface{}) (int, error) {
	if _, ok := v.(xdrType); ok {
		if bm, ok := v.(encoding.BinaryMarshaler); ok {
			b, err := bm.MarshalBinary()
			if err != nil {
				return 0, err
			}
			return w.Write(b)
		}
	}
	// delegate to xdr package's Marshal
	return xdr.Marshal(w, v)
}

// Value is an XDR Typedef defines as:
//
//   typedef opaque Value<>;
//
type Value []byte

// EncodeTo encodes this value using the Encoder.
func (s Value) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeOpaque(s[:]); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Value)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Value) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	(*s), nTmp, err = d.DecodeOpaque(0)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Value: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Value) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Value) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Value)(nil)
	_ encoding.BinaryUnmarshaler = (*Value)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s Value) xdrType() {}

var _ xdrType = (*Value)(nil)

// ScpBallot is an XDR Struct defines as:
//
//   struct SCPBallot
//    {
//        uint32 counter; // n
//        Value value;    // x
//    };
//
type ScpBallot struct {
	Counter Uint32
	Value   Value
}

// EncodeTo encodes this value using the Encoder.
func (s *ScpBallot) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Counter.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Value.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*ScpBallot)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *ScpBallot) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Counter.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.Value.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Value: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ScpBallot) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ScpBallot) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ScpBallot)(nil)
	_ encoding.BinaryUnmarshaler = (*ScpBallot)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ScpBallot) xdrType() {}

var _ xdrType = (*ScpBallot)(nil)

// ScpStatementType is an XDR Enum defines as:
//
//   enum SCPStatementType
//    {
//        SCP_ST_PREPARE = 0,
//        SCP_ST_CONFIRM = 1,
//        SCP_ST_EXTERNALIZE = 2,
//        SCP_ST_NOMINATE = 3
//    };
//
type ScpStatementType int32

const (
	ScpStatementTypeScpStPrepare     ScpStatementType = 0
	ScpStatementTypeScpStConfirm     ScpStatementType = 1
	ScpStatementTypeScpStExternalize ScpStatementType = 2
	ScpStatementTypeScpStNominate    ScpStatementType = 3
)

var scpStatementTypeMap = map[int32]string{
	0: "ScpStatementTypeScpStPrepare",
	1: "ScpStatementTypeScpStConfirm",
	2: "ScpStatementTypeScpStExternalize",
	3: "ScpStatementTypeScpStNominate",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ScpStatementType
func (e ScpStatementType) ValidEnum(v int32) bool {
	_, ok := scpStatementTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e ScpStatementType) String() string {
	name, _ := scpStatementTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e ScpStatementType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := scpStatementTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid ScpStatementType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*ScpStatementType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *ScpStatementType) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding ScpStatementType: %s", err)
	}
	if _, ok := scpStatementTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid ScpStatementType enum value", v)
	}
	*e = ScpStatementType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ScpStatementType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ScpStatementType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ScpStatementType)(nil)
	_ encoding.BinaryUnmarshaler = (*ScpStatementType)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ScpStatementType) xdrType() {}

var _ xdrType = (*ScpStatementType)(nil)

// ScpNomination is an XDR Struct defines as:
//
//   struct SCPNomination
//    {
//        Hash quorumSetHash; // D
//        Value votes<>;      // X
//        Value accepted<>;   // Y
//    };
//
type ScpNomination struct {
	QuorumSetHash Hash
	Votes         []Value
	Accepted      []Value
}

// EncodeTo encodes this value using the Encoder.
func (s *ScpNomination) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.QuorumSetHash.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeUint(uint32(len(s.Votes))); err != nil {
		return err
	}
	for i := 0; i < len(s.Votes); i++ {
		if err = s.Votes[i].EncodeTo(e); err != nil {
			return err
		}
	}
	if _, err = e.EncodeUint(uint32(len(s.Accepted))); err != nil {
		return err
	}
	for i := 0; i < len(s.Accepted); i++ {
		if err = s.Accepted[i].EncodeTo(e); err != nil {
			return err
		}
	}
	return nil
}

var _ decoderFrom = (*ScpNomination)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *ScpNomination) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.QuorumSetHash.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Hash: %s", err)
	}
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Value: %s", err)
	}
	s.Votes = nil
	if l > 0 {
		s.Votes = make([]Value, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.Votes[i].DecodeFrom(d)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding Value: %s", err)
			}
		}
	}
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Value: %s", err)
	}
	s.Accepted = nil
	if l > 0 {
		s.Accepted = make([]Value, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.Accepted[i].DecodeFrom(d)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding Value: %s", err)
			}
		}
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ScpNomination) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ScpNomination) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ScpNomination)(nil)
	_ encoding.BinaryUnmarshaler = (*ScpNomination)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ScpNomination) xdrType() {}

var _ xdrType = (*ScpNomination)(nil)

// ScpStatementPrepare is an XDR NestedStruct defines as:
//
//   struct
//            {
//                Hash quorumSetHash;       // D
//                SCPBallot ballot;         // b
//                SCPBallot* prepared;      // p
//                SCPBallot* preparedPrime; // p'
//                uint32 nC;                // c.n
//                uint32 nH;                // h.n
//            }
//
type ScpStatementPrepare struct {
	QuorumSetHash Hash
	Ballot        ScpBallot
	Prepared      *ScpBallot
	PreparedPrime *ScpBallot
	NC            Uint32
	NH            Uint32
}

// EncodeTo encodes this value using the Encoder.
func (s *ScpStatementPrepare) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.QuorumSetHash.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Ballot.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeBool(s.Prepared != nil); err != nil {
		return err
	}
	if s.Prepared != nil {
		if err = (*s.Prepared).EncodeTo(e); err != nil {
			return err
		}
	}
	if _, err = e.EncodeBool(s.PreparedPrime != nil); err != nil {
		return err
	}
	if s.PreparedPrime != nil {
		if err = (*s.PreparedPrime).EncodeTo(e); err != nil {
			return err
		}
	}
	if err = s.NC.EncodeTo(e); err != nil {
		return err
	}
	if err = s.NH.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*ScpStatementPrepare)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *ScpStatementPrepare) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.QuorumSetHash.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Hash: %s", err)
	}
	nTmp, err = s.Ballot.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ScpBallot: %s", err)
	}
	var b bool
	b, nTmp, err = d.DecodeBool()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ScpBallot: %s", err)
	}
	s.Prepared = nil
	if b {
		s.Prepared = new(ScpBallot)
		nTmp, err = s.Prepared.DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ScpBallot: %s", err)
		}
	}
	b, nTmp, err = d.DecodeBool()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ScpBallot: %s", err)
	}
	s.PreparedPrime = nil
	if b {
		s.PreparedPrime = new(ScpBallot)
		nTmp, err = s.PreparedPrime.DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ScpBallot: %s", err)
		}
	}
	nTmp, err = s.NC.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.NH.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ScpStatementPrepare) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ScpStatementPrepare) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ScpStatementPrepare)(nil)
	_ encoding.BinaryUnmarshaler = (*ScpStatementPrepare)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ScpStatementPrepare) xdrType() {}

var _ xdrType = (*ScpStatementPrepare)(nil)

// ScpStatementConfirm is an XDR NestedStruct defines as:
//
//   struct
//            {
//                SCPBallot ballot;   // b
//                uint32 nPrepared;   // p.n
//                uint32 nCommit;     // c.n
//                uint32 nH;          // h.n
//                Hash quorumSetHash; // D
//            }
//
type ScpStatementConfirm struct {
	Ballot        ScpBallot
	NPrepared     Uint32
	NCommit       Uint32
	NH            Uint32
	QuorumSetHash Hash
}

// EncodeTo encodes this value using the Encoder.
func (s *ScpStatementConfirm) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Ballot.EncodeTo(e); err != nil {
		return err
	}
	if err = s.NPrepared.EncodeTo(e); err != nil {
		return err
	}
	if err = s.NCommit.EncodeTo(e); err != nil {
		return err
	}
	if err = s.NH.EncodeTo(e); err != nil {
		return err
	}
	if err = s.QuorumSetHash.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*ScpStatementConfirm)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *ScpStatementConfirm) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Ballot.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ScpBallot: %s", err)
	}
	nTmp, err = s.NPrepared.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.NCommit.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.NH.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.QuorumSetHash.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Hash: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ScpStatementConfirm) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ScpStatementConfirm) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ScpStatementConfirm)(nil)
	_ encoding.BinaryUnmarshaler = (*ScpStatementConfirm)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ScpStatementConfirm) xdrType() {}

var _ xdrType = (*ScpStatementConfirm)(nil)

// ScpStatementExternalize is an XDR NestedStruct defines as:
//
//   struct
//            {
//                SCPBallot commit;         // c
//                uint32 nH;                // h.n
//                Hash commitQuorumSetHash; // D used before EXTERNALIZE
//            }
//
type ScpStatementExternalize struct {
	Commit              ScpBallot
	NH                  Uint32
	CommitQuorumSetHash Hash
}

// EncodeTo encodes this value using the Encoder.
func (s *ScpStatementExternalize) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Commit.EncodeTo(e); err != nil {
		return err
	}
	if err = s.NH.EncodeTo(e); err != nil {
		return err
	}
	if err = s.CommitQuorumSetHash.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*ScpStatementExternalize)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *ScpStatementExternalize) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Commit.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ScpBallot: %s", err)
	}
	nTmp, err = s.NH.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.CommitQuorumSetHash.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Hash: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ScpStatementExternalize) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ScpStatementExternalize) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ScpStatementExternalize)(nil)
	_ encoding.BinaryUnmarshaler = (*ScpStatementExternalize)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ScpStatementExternalize) xdrType() {}

var _ xdrType = (*ScpStatementExternalize)(nil)

// ScpStatementPledges is an XDR NestedUnion defines as:
//
//   union switch (SCPStatementType type)
//        {
//        case SCP_ST_PREPARE:
//            struct
//            {
//                Hash quorumSetHash;       // D
//                SCPBallot ballot;         // b
//                SCPBallot* prepared;      // p
//                SCPBallot* preparedPrime; // p'
//                uint32 nC;                // c.n
//                uint32 nH;                // h.n
//            } prepare;
//        case SCP_ST_CONFIRM:
//            struct
//            {
//                SCPBallot ballot;   // b
//                uint32 nPrepared;   // p.n
//                uint32 nCommit;     // c.n
//                uint32 nH;          // h.n
//                Hash quorumSetHash; // D
//            } confirm;
//        case SCP_ST_EXTERNALIZE:
//            struct
//            {
//                SCPBallot commit;         // c
//                uint32 nH;                // h.n
//                Hash commitQuorumSetHash; // D used before EXTERNALIZE
//            } externalize;
//        case SCP_ST_NOMINATE:
//            SCPNomination nominate;
//        }
//
type ScpStatementPledges struct {
	Type        ScpStatementType
	Prepare     *ScpStatementPrepare
	Confirm     *ScpStatementConfirm
	Externalize *ScpStatementExternalize
	Nominate    *ScpNomination
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ScpStatementPledges) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ScpStatementPledges
func (u ScpStatementPledges) ArmForSwitch(sw int32) (string, bool) {
	switch ScpStatementType(sw) {
	case ScpStatementTypeScpStPrepare:
		return "Prepare", true
	case ScpStatementTypeScpStConfirm:
		return "Confirm", true
	case ScpStatementTypeScpStExternalize:
		return "Externalize", true
	case ScpStatementTypeScpStNominate:
		return "Nominate", true
	}
	return "-", false
}

// NewScpStatementPledges creates a new  ScpStatementPledges.
func NewScpStatementPledges(aType ScpStatementType, value interface{}) (result ScpStatementPledges, err error) {
	result.Type = aType
	switch ScpStatementType(aType) {
	case ScpStatementTypeScpStPrepare:
		tv, ok := value.(ScpStatementPrepare)
		if !ok {
			err = fmt.Errorf("invalid value, must be ScpStatementPrepare")
			return
		}
		result.Prepare = &tv
	case ScpStatementTypeScpStConfirm:
		tv, ok := value.(ScpStatementConfirm)
		if !ok {
			err = fmt.Errorf("invalid value, must be ScpStatementConfirm")
			return
		}
		result.Confirm = &tv
	case ScpStatementTypeScpStExternalize:
		tv, ok := value.(ScpStatementExternalize)
		if !ok {
			err = fmt.Errorf("invalid value, must be ScpStatementExternalize")
			return
		}
		result.Externalize = &tv
	case ScpStatementTypeScpStNominate:
		tv, ok := value.(ScpNomination)
		if !ok {
			err = fmt.Errorf("invalid value, must be ScpNomination")
			return
		}
		result.Nominate = &tv
	}
	return
}

// MustPrepare retrieves the Prepare value from the union,
// panicing if the value is not set.
func (u ScpStatementPledges) MustPrepare() ScpStatementPrepare {
	val, ok := u.GetPrepare()

	if !ok {
		panic("arm Prepare is not set")
	}

	return val
}

// GetPrepare retrieves the Prepare value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ScpStatementPledges) GetPrepare() (result ScpStatementPrepare, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Prepare" {
		result = *u.Prepare
		ok = true
	}

	return
}

// MustConfirm retrieves the Confirm value from the union,
// panicing if the value is not set.
func (u ScpStatementPledges) MustConfirm() ScpStatementConfirm {
	val, ok := u.GetConfirm()

	if !ok {
		panic("arm Confirm is not set")
	}

	return val
}

// GetConfirm retrieves the Confirm value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ScpStatementPledges) GetConfirm() (result ScpStatementConfirm, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Confirm" {
		result = *u.Confirm
		ok = true
	}

	return
}

// MustExternalize retrieves the Externalize value from the union,
// panicing if the value is not set.
func (u ScpStatementPledges) MustExternalize() ScpStatementExternalize {
	val, ok := u.GetExternalize()

	if !ok {
		panic("arm Externalize is not set")
	}

	return val
}

// GetExternalize retrieves the Externalize value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ScpStatementPledges) GetExternalize() (result ScpStatementExternalize, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Externalize" {
		result = *u.Externalize
		ok = true
	}

	return
}

// MustNominate retrieves the Nominate value from the union,
// panicing if the value is not set.
func (u ScpStatementPledges) MustNominate() ScpNomination {
	val, ok := u.GetNominate()

	if !ok {
		panic("arm Nominate is not set")
	}

	return val
}

// GetNominate retrieves the Nominate value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ScpStatementPledges) GetNominate() (result ScpNomination, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Nominate" {
		result = *u.Nominate
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u ScpStatementPledges) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch ScpStatementType(u.Type) {
	case ScpStatementTypeScpStPrepare:
		if err = (*u.Prepare).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case ScpStatementTypeScpStConfirm:
		if err = (*u.Confirm).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case ScpStatementTypeScpStExternalize:
		if err = (*u.Externalize).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case ScpStatementTypeScpStNominate:
		if err = (*u.Nominate).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (ScpStatementType) switch value '%d' is not valid for union ScpStatementPledges", u.Type)
}

var _ decoderFrom = (*ScpStatementPledges)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *ScpStatementPledges) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ScpStatementType: %s", err)
	}
	switch ScpStatementType(u.Type) {
	case ScpStatementTypeScpStPrepare:
		u.Prepare = new(ScpStatementPrepare)
		nTmp, err = (*u.Prepare).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ScpStatementPrepare: %s", err)
		}
		return n, nil
	case ScpStatementTypeScpStConfirm:
		u.Confirm = new(ScpStatementConfirm)
		nTmp, err = (*u.Confirm).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ScpStatementConfirm: %s", err)
		}
		return n, nil
	case ScpStatementTypeScpStExternalize:
		u.Externalize = new(ScpStatementExternalize)
		nTmp, err = (*u.Externalize).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ScpStatementExternalize: %s", err)
		}
		return n, nil
	case ScpStatementTypeScpStNominate:
		u.Nominate = new(ScpNomination)
		nTmp, err = (*u.Nominate).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ScpNomination: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union ScpStatementPledges has invalid Type (ScpStatementType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ScpStatementPledges) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ScpStatementPledges) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ScpStatementPledges)(nil)
	_ encoding.BinaryUnmarshaler = (*ScpStatementPledges)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ScpStatementPledges) xdrType() {}

var _ xdrType = (*ScpStatementPledges)(nil)

// ScpStatement is an XDR Struct defines as:
//
//   struct SCPStatement
//    {
//        NodeID nodeID;    // v
//        uint64 slotIndex; // i
//
//        union switch (SCPStatementType type)
//        {
//        case SCP_ST_PREPARE:
//            struct
//            {
//                Hash quorumSetHash;       // D
//                SCPBallot ballot;         // b
//                SCPBallot* prepared;      // p
//                SCPBallot* preparedPrime; // p'
//                uint32 nC;                // c.n
//                uint32 nH;                // h.n
//            } prepare;
//        case SCP_ST_CONFIRM:
//            struct
//            {
//                SCPBallot ballot;   // b
//                uint32 nPrepared;   // p.n
//                uint32 nCommit;     // c.n
//                uint32 nH;          // h.n
//                Hash quorumSetHash; // D
//            } confirm;
//        case SCP_ST_EXTERNALIZE:
//            struct
//            {
//                SCPBallot commit;         // c
//                uint32 nH;                // h.n
//                Hash commitQuorumSetHash; // D used before EXTERNALIZE
//            } externalize;
//        case SCP_ST_NOMINATE:
//            SCPNomination nominate;
//        }
//        pledges;
//    };
//
type ScpStatement struct {
	NodeId    NodeId
	SlotIndex Uint64
	Pledges   ScpStatementPledges
}

// EncodeTo encodes this value using the Encoder.
func (s *ScpStatement) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.NodeId.EncodeTo(e); err != nil {
		return err
	}
	if err = s.SlotIndex.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Pledges.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*ScpStatement)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *ScpStatement) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.NodeId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding NodeId: %s", err)
	}
	nTmp, err = s.SlotIndex.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint64: %s", err)
	}
	nTmp, err = s.Pledges.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ScpStatementPledges: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ScpStatement) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ScpStatement) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ScpStatement)(nil)
	_ encoding.BinaryUnmarshaler = (*ScpStatement)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ScpStatement) xdrType() {}

var _ xdrType = (*ScpStatement)(nil)

// ScpEnvelope is an XDR Struct defines as:
//
//   struct SCPEnvelope
//    {
//        SCPStatement statement;
//        Signature signature;
//    };
//
type ScpEnvelope struct {
	Statement ScpStatement
	Signature Signature
}

// EncodeTo encodes this value using the Encoder.
func (s *ScpEnvelope) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Statement.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Signature.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*ScpEnvelope)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *ScpEnvelope) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Statement.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ScpStatement: %s", err)
	}
	nTmp, err = s.Signature.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Signature: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ScpEnvelope) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ScpEnvelope) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ScpEnvelope)(nil)
	_ encoding.BinaryUnmarshaler = (*ScpEnvelope)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ScpEnvelope) xdrType() {}

var _ xdrType = (*ScpEnvelope)(nil)

// ScpQuorumSet is an XDR Struct defines as:
//
//   struct SCPQuorumSet
//    {
//        uint32 threshold;
//        NodeID validators<>;
//        SCPQuorumSet innerSets<>;
//    };
//
type ScpQuorumSet struct {
	Threshold  Uint32
	Validators []NodeId
	InnerSets  []ScpQuorumSet
}

// EncodeTo encodes this value using the Encoder.
func (s *ScpQuorumSet) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Threshold.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeUint(uint32(len(s.Validators))); err != nil {
		return err
	}
	for i := 0; i < len(s.Validators); i++ {
		if err = s.Validators[i].EncodeTo(e); err != nil {
			return err
		}
	}
	if _, err = e.EncodeUint(uint32(len(s.InnerSets))); err != nil {
		return err
	}
	for i := 0; i < len(s.InnerSets); i++ {
		if err = s.InnerSets[i].EncodeTo(e); err != nil {
			return err
		}
	}
	return nil
}

var _ decoderFrom = (*ScpQuorumSet)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *ScpQuorumSet) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Threshold.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding NodeId: %s", err)
	}
	s.Validators = nil
	if l > 0 {
		s.Validators = make([]NodeId, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.Validators[i].DecodeFrom(d)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding NodeId: %s", err)
			}
		}
	}
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ScpQuorumSet: %s", err)
	}
	s.InnerSets = nil
	if l > 0 {
		s.InnerSets = make([]ScpQuorumSet, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.InnerSets[i].DecodeFrom(d)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding ScpQuorumSet: %s", err)
			}
		}
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ScpQuorumSet) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ScpQuorumSet) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ScpQuorumSet)(nil)
	_ encoding.BinaryUnmarshaler = (*ScpQuorumSet)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ScpQuorumSet) xdrType() {}

var _ xdrType = (*ScpQuorumSet)(nil)

// AccountId is an XDR Typedef defines as:
//
//   typedef PublicKey AccountID;
//
type AccountId PublicKey

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u AccountId) SwitchFieldName() string {
	return PublicKey(u).SwitchFieldName()
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PublicKey
func (u AccountId) ArmForSwitch(sw int32) (string, bool) {
	return PublicKey(u).ArmForSwitch(sw)
}

// NewAccountId creates a new  AccountId.
func NewAccountId(aType PublicKeyType, value interface{}) (result AccountId, err error) {
	u, err := NewPublicKey(aType, value)
	result = AccountId(u)
	return
}

// MustEd25519 retrieves the Ed25519 value from the union,
// panicing if the value is not set.
func (u AccountId) MustEd25519() Uint256 {
	return PublicKey(u).MustEd25519()
}

// GetEd25519 retrieves the Ed25519 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u AccountId) GetEd25519() (result Uint256, ok bool) {
	return PublicKey(u).GetEd25519()
}

// EncodeTo encodes this value using the Encoder.
func (s AccountId) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = PublicKey(s).EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*AccountId)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *AccountId) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = (*PublicKey)(s).DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding PublicKey: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s AccountId) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AccountId) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*AccountId)(nil)
	_ encoding.BinaryUnmarshaler = (*AccountId)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s AccountId) xdrType() {}

var _ xdrType = (*AccountId)(nil)

// Thresholds is an XDR Typedef defines as:
//
//   typedef opaque Thresholds[4];
//
type Thresholds [4]byte

// XDRMaxSize implements the Sized interface for Thresholds
func (e Thresholds) XDRMaxSize() int {
	return 4
}

// EncodeTo encodes this value using the Encoder.
func (s *Thresholds) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeFixedOpaque(s[:]); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Thresholds)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Thresholds) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = d.DecodeFixedOpaqueInplace(s[:])
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Thresholds: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Thresholds) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Thresholds) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Thresholds)(nil)
	_ encoding.BinaryUnmarshaler = (*Thresholds)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s Thresholds) xdrType() {}

var _ xdrType = (*Thresholds)(nil)

// String32 is an XDR Typedef defines as:
//
//   typedef string string32<32>;
//
type String32 string

// XDRMaxSize implements the Sized interface for String32
func (e String32) XDRMaxSize() int {
	return 32
}

// EncodeTo encodes this value using the Encoder.
func (s String32) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeString(string(s)); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*String32)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *String32) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	var v string
	v, nTmp, err = d.DecodeString(32)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding String32: %s", err)
	}
	*s = String32(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s String32) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *String32) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*String32)(nil)
	_ encoding.BinaryUnmarshaler = (*String32)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s String32) xdrType() {}

var _ xdrType = (*String32)(nil)

// String64 is an XDR Typedef defines as:
//
//   typedef string string64<64>;
//
type String64 string

// XDRMaxSize implements the Sized interface for String64
func (e String64) XDRMaxSize() int {
	return 64
}

// EncodeTo encodes this value using the Encoder.
func (s String64) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeString(string(s)); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*String64)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *String64) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	var v string
	v, nTmp, err = d.DecodeString(64)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding String64: %s", err)
	}
	*s = String64(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s String64) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *String64) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*String64)(nil)
	_ encoding.BinaryUnmarshaler = (*String64)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s String64) xdrType() {}

var _ xdrType = (*String64)(nil)

// SequenceNumber is an XDR Typedef defines as:
//
//   typedef int64 SequenceNumber;
//
type SequenceNumber Int64

// EncodeTo encodes this value using the Encoder.
func (s SequenceNumber) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = Int64(s).EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*SequenceNumber)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *SequenceNumber) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = (*Int64)(s).DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s SequenceNumber) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SequenceNumber) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*SequenceNumber)(nil)
	_ encoding.BinaryUnmarshaler = (*SequenceNumber)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s SequenceNumber) xdrType() {}

var _ xdrType = (*SequenceNumber)(nil)

// TimePoint is an XDR Typedef defines as:
//
//   typedef uint64 TimePoint;
//
type TimePoint Uint64

// EncodeTo encodes this value using the Encoder.
func (s TimePoint) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = Uint64(s).EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*TimePoint)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *TimePoint) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = (*Uint64)(s).DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint64: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TimePoint) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TimePoint) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TimePoint)(nil)
	_ encoding.BinaryUnmarshaler = (*TimePoint)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TimePoint) xdrType() {}

var _ xdrType = (*TimePoint)(nil)

// DataValue is an XDR Typedef defines as:
//
//   typedef opaque DataValue<64>;
//
type DataValue []byte

// XDRMaxSize implements the Sized interface for DataValue
func (e DataValue) XDRMaxSize() int {
	return 64
}

// EncodeTo encodes this value using the Encoder.
func (s DataValue) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeOpaque(s[:]); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*DataValue)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *DataValue) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	(*s), nTmp, err = d.DecodeOpaque(64)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding DataValue: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s DataValue) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *DataValue) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*DataValue)(nil)
	_ encoding.BinaryUnmarshaler = (*DataValue)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s DataValue) xdrType() {}

var _ xdrType = (*DataValue)(nil)

// PoolId is an XDR Typedef defines as:
//
//   typedef Hash PoolID;
//
type PoolId Hash

// EncodeTo encodes this value using the Encoder.
func (s *PoolId) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = (*Hash)(s).EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*PoolId)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *PoolId) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = (*Hash)(s).DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Hash: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s PoolId) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *PoolId) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*PoolId)(nil)
	_ encoding.BinaryUnmarshaler = (*PoolId)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s PoolId) xdrType() {}

var _ xdrType = (*PoolId)(nil)

// AssetCode4 is an XDR Typedef defines as:
//
//   typedef opaque AssetCode4[4];
//
type AssetCode4 [4]byte

// XDRMaxSize implements the Sized interface for AssetCode4
func (e AssetCode4) XDRMaxSize() int {
	return 4
}

// EncodeTo encodes this value using the Encoder.
func (s *AssetCode4) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeFixedOpaque(s[:]); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*AssetCode4)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *AssetCode4) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = d.DecodeFixedOpaqueInplace(s[:])
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AssetCode4: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s AssetCode4) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AssetCode4) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*AssetCode4)(nil)
	_ encoding.BinaryUnmarshaler = (*AssetCode4)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s AssetCode4) xdrType() {}

var _ xdrType = (*AssetCode4)(nil)

// AssetCode12 is an XDR Typedef defines as:
//
//   typedef opaque AssetCode12[12];
//
type AssetCode12 [12]byte

// XDRMaxSize implements the Sized interface for AssetCode12
func (e AssetCode12) XDRMaxSize() int {
	return 12
}

// EncodeTo encodes this value using the Encoder.
func (s *AssetCode12) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeFixedOpaque(s[:]); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*AssetCode12)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *AssetCode12) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = d.DecodeFixedOpaqueInplace(s[:])
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AssetCode12: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s AssetCode12) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AssetCode12) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*AssetCode12)(nil)
	_ encoding.BinaryUnmarshaler = (*AssetCode12)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s AssetCode12) xdrType() {}

var _ xdrType = (*AssetCode12)(nil)

// AssetType is an XDR Enum defines as:
//
//   enum AssetType
//    {
//        ASSET_TYPE_NATIVE = 0,
//        ASSET_TYPE_CREDIT_ALPHANUM4 = 1,
//        ASSET_TYPE_CREDIT_ALPHANUM12 = 2,
//        ASSET_TYPE_POOL_SHARE = 3
//    };
//
type AssetType int32

const (
	AssetTypeAssetTypeNative           AssetType = 0
	AssetTypeAssetTypeCreditAlphanum4  AssetType = 1
	AssetTypeAssetTypeCreditAlphanum12 AssetType = 2
	AssetTypeAssetTypePoolShare        AssetType = 3
)

var assetTypeMap = map[int32]string{
	0: "AssetTypeAssetTypeNative",
	1: "AssetTypeAssetTypeCreditAlphanum4",
	2: "AssetTypeAssetTypeCreditAlphanum12",
	3: "AssetTypeAssetTypePoolShare",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for AssetType
func (e AssetType) ValidEnum(v int32) bool {
	_, ok := assetTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e AssetType) String() string {
	name, _ := assetTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e AssetType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := assetTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid AssetType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*AssetType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *AssetType) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding AssetType: %s", err)
	}
	if _, ok := assetTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid AssetType enum value", v)
	}
	*e = AssetType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s AssetType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AssetType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*AssetType)(nil)
	_ encoding.BinaryUnmarshaler = (*AssetType)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s AssetType) xdrType() {}

var _ xdrType = (*AssetType)(nil)

// AssetCode is an XDR Union defines as:
//
//   union AssetCode switch (AssetType type)
//    {
//    case ASSET_TYPE_CREDIT_ALPHANUM4:
//        AssetCode4 assetCode4;
//
//    case ASSET_TYPE_CREDIT_ALPHANUM12:
//        AssetCode12 assetCode12;
//
//        // add other asset types here in the future
//    };
//
type AssetCode struct {
	Type        AssetType
	AssetCode4  *AssetCode4
	AssetCode12 *AssetCode12
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u AssetCode) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of AssetCode
func (u AssetCode) ArmForSwitch(sw int32) (string, bool) {
	switch AssetType(sw) {
	case AssetTypeAssetTypeCreditAlphanum4:
		return "AssetCode4", true
	case AssetTypeAssetTypeCreditAlphanum12:
		return "AssetCode12", true
	}
	return "-", false
}

// NewAssetCode creates a new  AssetCode.
func NewAssetCode(aType AssetType, value interface{}) (result AssetCode, err error) {
	result.Type = aType
	switch AssetType(aType) {
	case AssetTypeAssetTypeCreditAlphanum4:
		tv, ok := value.(AssetCode4)
		if !ok {
			err = fmt.Errorf("invalid value, must be AssetCode4")
			return
		}
		result.AssetCode4 = &tv
	case AssetTypeAssetTypeCreditAlphanum12:
		tv, ok := value.(AssetCode12)
		if !ok {
			err = fmt.Errorf("invalid value, must be AssetCode12")
			return
		}
		result.AssetCode12 = &tv
	}
	return
}

// MustAssetCode4 retrieves the AssetCode4 value from the union,
// panicing if the value is not set.
func (u AssetCode) MustAssetCode4() AssetCode4 {
	val, ok := u.GetAssetCode4()

	if !ok {
		panic("arm AssetCode4 is not set")
	}

	return val
}

// GetAssetCode4 retrieves the AssetCode4 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u AssetCode) GetAssetCode4() (result AssetCode4, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AssetCode4" {
		result = *u.AssetCode4
		ok = true
	}

	return
}

// MustAssetCode12 retrieves the AssetCode12 value from the union,
// panicing if the value is not set.
func (u AssetCode) MustAssetCode12() AssetCode12 {
	val, ok := u.GetAssetCode12()

	if !ok {
		panic("arm AssetCode12 is not set")
	}

	return val
}

// GetAssetCode12 retrieves the AssetCode12 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u AssetCode) GetAssetCode12() (result AssetCode12, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AssetCode12" {
		result = *u.AssetCode12
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u AssetCode) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch AssetType(u.Type) {
	case AssetTypeAssetTypeCreditAlphanum4:
		if err = (*u.AssetCode4).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case AssetTypeAssetTypeCreditAlphanum12:
		if err = (*u.AssetCode12).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (AssetType) switch value '%d' is not valid for union AssetCode", u.Type)
}

var _ decoderFrom = (*AssetCode)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *AssetCode) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AssetType: %s", err)
	}
	switch AssetType(u.Type) {
	case AssetTypeAssetTypeCreditAlphanum4:
		u.AssetCode4 = new(AssetCode4)
		nTmp, err = (*u.AssetCode4).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding AssetCode4: %s", err)
		}
		return n, nil
	case AssetTypeAssetTypeCreditAlphanum12:
		u.AssetCode12 = new(AssetCode12)
		nTmp, err = (*u.AssetCode12).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding AssetCode12: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union AssetCode has invalid Type (AssetType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s AssetCode) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AssetCode) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*AssetCode)(nil)
	_ encoding.BinaryUnmarshaler = (*AssetCode)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s AssetCode) xdrType() {}

var _ xdrType = (*AssetCode)(nil)

// AlphaNum4 is an XDR Struct defines as:
//
//   struct AlphaNum4
//    {
//        AssetCode4 assetCode;
//        AccountID issuer;
//    };
//
type AlphaNum4 struct {
	AssetCode AssetCode4
	Issuer    AccountId
}

// EncodeTo encodes this value using the Encoder.
func (s *AlphaNum4) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.AssetCode.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Issuer.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*AlphaNum4)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *AlphaNum4) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.AssetCode.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AssetCode4: %s", err)
	}
	nTmp, err = s.Issuer.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountId: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s AlphaNum4) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AlphaNum4) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*AlphaNum4)(nil)
	_ encoding.BinaryUnmarshaler = (*AlphaNum4)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s AlphaNum4) xdrType() {}

var _ xdrType = (*AlphaNum4)(nil)

// AlphaNum12 is an XDR Struct defines as:
//
//   struct AlphaNum12
//    {
//        AssetCode12 assetCode;
//        AccountID issuer;
//    };
//
type AlphaNum12 struct {
	AssetCode AssetCode12
	Issuer    AccountId
}

// EncodeTo encodes this value using the Encoder.
func (s *AlphaNum12) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.AssetCode.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Issuer.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*AlphaNum12)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *AlphaNum12) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.AssetCode.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AssetCode12: %s", err)
	}
	nTmp, err = s.Issuer.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountId: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s AlphaNum12) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AlphaNum12) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*AlphaNum12)(nil)
	_ encoding.BinaryUnmarshaler = (*AlphaNum12)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s AlphaNum12) xdrType() {}

var _ xdrType = (*AlphaNum12)(nil)

// Asset is an XDR Union defines as:
//
//   union Asset switch (AssetType type)
//    {
//    case ASSET_TYPE_NATIVE: // Not credit
//        void;
//
//    case ASSET_TYPE_CREDIT_ALPHANUM4:
//        AlphaNum4 alphaNum4;
//
//    case ASSET_TYPE_CREDIT_ALPHANUM12:
//        AlphaNum12 alphaNum12;
//
//        // add other asset types here in the future
//    };
//
type Asset struct {
	Type       AssetType
	AlphaNum4  *AlphaNum4
	AlphaNum12 *AlphaNum12
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u Asset) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of Asset
func (u Asset) ArmForSwitch(sw int32) (string, bool) {
	switch AssetType(sw) {
	case AssetTypeAssetTypeNative:
		return "", true
	case AssetTypeAssetTypeCreditAlphanum4:
		return "AlphaNum4", true
	case AssetTypeAssetTypeCreditAlphanum12:
		return "AlphaNum12", true
	}
	return "-", false
}

// NewAsset creates a new  Asset.
func NewAsset(aType AssetType, value interface{}) (result Asset, err error) {
	result.Type = aType
	switch AssetType(aType) {
	case AssetTypeAssetTypeNative:
		// void
	case AssetTypeAssetTypeCreditAlphanum4:
		tv, ok := value.(AlphaNum4)
		if !ok {
			err = fmt.Errorf("invalid value, must be AlphaNum4")
			return
		}
		result.AlphaNum4 = &tv
	case AssetTypeAssetTypeCreditAlphanum12:
		tv, ok := value.(AlphaNum12)
		if !ok {
			err = fmt.Errorf("invalid value, must be AlphaNum12")
			return
		}
		result.AlphaNum12 = &tv
	}
	return
}

// MustAlphaNum4 retrieves the AlphaNum4 value from the union,
// panicing if the value is not set.
func (u Asset) MustAlphaNum4() AlphaNum4 {
	val, ok := u.GetAlphaNum4()

	if !ok {
		panic("arm AlphaNum4 is not set")
	}

	return val
}

// GetAlphaNum4 retrieves the AlphaNum4 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u Asset) GetAlphaNum4() (result AlphaNum4, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AlphaNum4" {
		result = *u.AlphaNum4
		ok = true
	}

	return
}

// MustAlphaNum12 retrieves the AlphaNum12 value from the union,
// panicing if the value is not set.
func (u Asset) MustAlphaNum12() AlphaNum12 {
	val, ok := u.GetAlphaNum12()

	if !ok {
		panic("arm AlphaNum12 is not set")
	}

	return val
}

// GetAlphaNum12 retrieves the AlphaNum12 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u Asset) GetAlphaNum12() (result AlphaNum12, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AlphaNum12" {
		result = *u.AlphaNum12
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u Asset) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch AssetType(u.Type) {
	case AssetTypeAssetTypeNative:
		// Void
		return nil
	case AssetTypeAssetTypeCreditAlphanum4:
		if err = (*u.AlphaNum4).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case AssetTypeAssetTypeCreditAlphanum12:
		if err = (*u.AlphaNum12).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (AssetType) switch value '%d' is not valid for union Asset", u.Type)
}

var _ decoderFrom = (*Asset)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *Asset) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AssetType: %s", err)
	}
	switch AssetType(u.Type) {
	case AssetTypeAssetTypeNative:
		// Void
		return n, nil
	case AssetTypeAssetTypeCreditAlphanum4:
		u.AlphaNum4 = new(AlphaNum4)
		nTmp, err = (*u.AlphaNum4).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding AlphaNum4: %s", err)
		}
		return n, nil
	case AssetTypeAssetTypeCreditAlphanum12:
		u.AlphaNum12 = new(AlphaNum12)
		nTmp, err = (*u.AlphaNum12).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding AlphaNum12: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union Asset has invalid Type (AssetType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Asset) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Asset) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Asset)(nil)
	_ encoding.BinaryUnmarshaler = (*Asset)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s Asset) xdrType() {}

var _ xdrType = (*Asset)(nil)

// Price is an XDR Struct defines as:
//
//   struct Price
//    {
//        int32 n; // numerator
//        int32 d; // denominator
//    };
//
type Price struct {
	N Int32
	D Int32
}

// EncodeTo encodes this value using the Encoder.
func (s *Price) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.N.EncodeTo(e); err != nil {
		return err
	}
	if err = s.D.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Price)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Price) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.N.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int32: %s", err)
	}
	nTmp, err = s.D.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int32: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Price) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Price) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Price)(nil)
	_ encoding.BinaryUnmarshaler = (*Price)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s Price) xdrType() {}

var _ xdrType = (*Price)(nil)

// Liabilities is an XDR Struct defines as:
//
//   struct Liabilities
//    {
//        int64 buying;
//        int64 selling;
//    };
//
type Liabilities struct {
	Buying  Int64
	Selling Int64
}

// EncodeTo encodes this value using the Encoder.
func (s *Liabilities) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Buying.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Selling.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Liabilities)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Liabilities) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Buying.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	nTmp, err = s.Selling.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Liabilities) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Liabilities) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Liabilities)(nil)
	_ encoding.BinaryUnmarshaler = (*Liabilities)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s Liabilities) xdrType() {}

var _ xdrType = (*Liabilities)(nil)

// ThresholdIndexes is an XDR Enum defines as:
//
//   enum ThresholdIndexes
//    {
//        THRESHOLD_MASTER_WEIGHT = 0,
//        THRESHOLD_LOW = 1,
//        THRESHOLD_MED = 2,
//        THRESHOLD_HIGH = 3
//    };
//
type ThresholdIndexes int32

const (
	ThresholdIndexesThresholdMasterWeight ThresholdIndexes = 0
	ThresholdIndexesThresholdLow          ThresholdIndexes = 1
	ThresholdIndexesThresholdMed          ThresholdIndexes = 2
	ThresholdIndexesThresholdHigh         ThresholdIndexes = 3
)

var thresholdIndexesMap = map[int32]string{
	0: "ThresholdIndexesThresholdMasterWeight",
	1: "ThresholdIndexesThresholdLow",
	2: "ThresholdIndexesThresholdMed",
	3: "ThresholdIndexesThresholdHigh",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ThresholdIndexes
func (e ThresholdIndexes) ValidEnum(v int32) bool {
	_, ok := thresholdIndexesMap[v]
	return ok
}

// String returns the name of `e`
func (e ThresholdIndexes) String() string {
	name, _ := thresholdIndexesMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e ThresholdIndexes) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := thresholdIndexesMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid ThresholdIndexes enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*ThresholdIndexes)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *ThresholdIndexes) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding ThresholdIndexes: %s", err)
	}
	if _, ok := thresholdIndexesMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid ThresholdIndexes enum value", v)
	}
	*e = ThresholdIndexes(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ThresholdIndexes) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ThresholdIndexes) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ThresholdIndexes)(nil)
	_ encoding.BinaryUnmarshaler = (*ThresholdIndexes)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ThresholdIndexes) xdrType() {}

var _ xdrType = (*ThresholdIndexes)(nil)

// LedgerEntryType is an XDR Enum defines as:
//
//   enum LedgerEntryType
//    {
//        ACCOUNT = 0,
//        TRUSTLINE = 1,
//        OFFER = 2,
//        DATA = 3,
//        CLAIMABLE_BALANCE = 4,
//        LIQUIDITY_POOL = 5
//    };
//
type LedgerEntryType int32

const (
	LedgerEntryTypeAccount          LedgerEntryType = 0
	LedgerEntryTypeTrustline        LedgerEntryType = 1
	LedgerEntryTypeOffer            LedgerEntryType = 2
	LedgerEntryTypeData             LedgerEntryType = 3
	LedgerEntryTypeClaimableBalance LedgerEntryType = 4
	LedgerEntryTypeLiquidityPool    LedgerEntryType = 5
)

var ledgerEntryTypeMap = map[int32]string{
	0: "LedgerEntryTypeAccount",
	1: "LedgerEntryTypeTrustline",
	2: "LedgerEntryTypeOffer",
	3: "LedgerEntryTypeData",
	4: "LedgerEntryTypeClaimableBalance",
	5: "LedgerEntryTypeLiquidityPool",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for LedgerEntryType
func (e LedgerEntryType) ValidEnum(v int32) bool {
	_, ok := ledgerEntryTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e LedgerEntryType) String() string {
	name, _ := ledgerEntryTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e LedgerEntryType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := ledgerEntryTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid LedgerEntryType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*LedgerEntryType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *LedgerEntryType) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding LedgerEntryType: %s", err)
	}
	if _, ok := ledgerEntryTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid LedgerEntryType enum value", v)
	}
	*e = LedgerEntryType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerEntryType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerEntryType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerEntryType)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerEntryType)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LedgerEntryType) xdrType() {}

var _ xdrType = (*LedgerEntryType)(nil)

// Signer is an XDR Struct defines as:
//
//   struct Signer
//    {
//        SignerKey key;
//        uint32 weight; // really only need 1 byte
//    };
//
type Signer struct {
	Key    SignerKey
	Weight Uint32
}

// EncodeTo encodes this value using the Encoder.
func (s *Signer) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Key.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Weight.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Signer)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Signer) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Key.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SignerKey: %s", err)
	}
	nTmp, err = s.Weight.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Signer) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Signer) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Signer)(nil)
	_ encoding.BinaryUnmarshaler = (*Signer)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s Signer) xdrType() {}

var _ xdrType = (*Signer)(nil)

// AccountFlags is an XDR Enum defines as:
//
//   enum AccountFlags
//    { // masks for each flag
//
//        // Flags set on issuer accounts
//        // TrustLines are created with authorized set to "false" requiring
//        // the issuer to set it for each TrustLine
//        AUTH_REQUIRED_FLAG = 0x1,
//        // If set, the authorized flag in TrustLines can be cleared
//        // otherwise, authorization cannot be revoked
//        AUTH_REVOCABLE_FLAG = 0x2,
//        // Once set, causes all AUTH_* flags to be read-only
//        AUTH_IMMUTABLE_FLAG = 0x4,
//        // Trustlines are created with clawback enabled set to "true",
//        // and claimable balances created from those trustlines are created
//        // with clawback enabled set to "true"
//        AUTH_CLAWBACK_ENABLED_FLAG = 0x8
//    };
//
type AccountFlags int32

const (
	AccountFlagsAuthRequiredFlag        AccountFlags = 1
	AccountFlagsAuthRevocableFlag       AccountFlags = 2
	AccountFlagsAuthImmutableFlag       AccountFlags = 4
	AccountFlagsAuthClawbackEnabledFlag AccountFlags = 8
)

var accountFlagsMap = map[int32]string{
	1: "AccountFlagsAuthRequiredFlag",
	2: "AccountFlagsAuthRevocableFlag",
	4: "AccountFlagsAuthImmutableFlag",
	8: "AccountFlagsAuthClawbackEnabledFlag",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for AccountFlags
func (e AccountFlags) ValidEnum(v int32) bool {
	_, ok := accountFlagsMap[v]
	return ok
}

// String returns the name of `e`
func (e AccountFlags) String() string {
	name, _ := accountFlagsMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e AccountFlags) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := accountFlagsMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid AccountFlags enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*AccountFlags)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *AccountFlags) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding AccountFlags: %s", err)
	}
	if _, ok := accountFlagsMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid AccountFlags enum value", v)
	}
	*e = AccountFlags(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s AccountFlags) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AccountFlags) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*AccountFlags)(nil)
	_ encoding.BinaryUnmarshaler = (*AccountFlags)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s AccountFlags) xdrType() {}

var _ xdrType = (*AccountFlags)(nil)

// MaskAccountFlags is an XDR Const defines as:
//
//   const MASK_ACCOUNT_FLAGS = 0x7;
//
const MaskAccountFlags = 0x7

// MaskAccountFlagsV17 is an XDR Const defines as:
//
//   const MASK_ACCOUNT_FLAGS_V17 = 0xF;
//
const MaskAccountFlagsV17 = 0xF

// MaxSigners is an XDR Const defines as:
//
//   const MAX_SIGNERS = 20;
//
const MaxSigners = 20

// SponsorshipDescriptor is an XDR Typedef defines as:
//
//   typedef AccountID* SponsorshipDescriptor;
//
type SponsorshipDescriptor = *AccountId

// AccountEntryExtensionV2Ext is an XDR NestedUnion defines as:
//
//   union switch (int v)
//        {
//        case 0:
//            void;
//        }
//
type AccountEntryExtensionV2Ext struct {
	V int32
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u AccountEntryExtensionV2Ext) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of AccountEntryExtensionV2Ext
func (u AccountEntryExtensionV2Ext) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	}
	return "-", false
}

// NewAccountEntryExtensionV2Ext creates a new  AccountEntryExtensionV2Ext.
func NewAccountEntryExtensionV2Ext(v int32, value interface{}) (result AccountEntryExtensionV2Ext, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u AccountEntryExtensionV2Ext) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union AccountEntryExtensionV2Ext", u.V)
}

var _ decoderFrom = (*AccountEntryExtensionV2Ext)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *AccountEntryExtensionV2Ext) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %s", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	}
	return n, fmt.Errorf("union AccountEntryExtensionV2Ext has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s AccountEntryExtensionV2Ext) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AccountEntryExtensionV2Ext) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*AccountEntryExtensionV2Ext)(nil)
	_ encoding.BinaryUnmarshaler = (*AccountEntryExtensionV2Ext)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s AccountEntryExtensionV2Ext) xdrType() {}

var _ xdrType = (*AccountEntryExtensionV2Ext)(nil)

// AccountEntryExtensionV2 is an XDR Struct defines as:
//
//   struct AccountEntryExtensionV2
//    {
//        uint32 numSponsored;
//        uint32 numSponsoring;
//        SponsorshipDescriptor signerSponsoringIDs<MAX_SIGNERS>;
//
//        union switch (int v)
//        {
//        case 0:
//            void;
//        }
//        ext;
//    };
//
type AccountEntryExtensionV2 struct {
	NumSponsored        Uint32
	NumSponsoring       Uint32
	SignerSponsoringIDs []SponsorshipDescriptor `xdrmaxsize:"20"`
	Ext                 AccountEntryExtensionV2Ext
}

// EncodeTo encodes this value using the Encoder.
func (s *AccountEntryExtensionV2) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.NumSponsored.EncodeTo(e); err != nil {
		return err
	}
	if err = s.NumSponsoring.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeUint(uint32(len(s.SignerSponsoringIDs))); err != nil {
		return err
	}
	for i := 0; i < len(s.SignerSponsoringIDs); i++ {
		if _, err = e.EncodeBool(s.SignerSponsoringIDs[i] != nil); err != nil {
			return err
		}
		if s.SignerSponsoringIDs[i] != nil {
			if err = s.SignerSponsoringIDs[i].EncodeTo(e); err != nil {
				return err
			}
		}
	}
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*AccountEntryExtensionV2)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *AccountEntryExtensionV2) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.NumSponsored.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.NumSponsoring.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SponsorshipDescriptor: %s", err)
	}
	if l > 20 {
		return n, fmt.Errorf("decoding SponsorshipDescriptor: data size (%d) exceeds size limit (20)", l)
	}
	s.SignerSponsoringIDs = nil
	if l > 0 {
		s.SignerSponsoringIDs = make([]SponsorshipDescriptor, l)
		for i := uint32(0); i < l; i++ {
			var eb bool
			eb, nTmp, err = d.DecodeBool()
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding SponsorshipDescriptor: %s", err)
			}
			s.SignerSponsoringIDs[i] = nil
			if eb {
				s.SignerSponsoringIDs[i] = new(AccountId)
				nTmp, err = s.SignerSponsoringIDs[i].DecodeFrom(d)
				n += nTmp
				if err != nil {
					return n, fmt.Errorf("decoding SponsorshipDescriptor: %s", err)
				}
			}
		}
	}
	nTmp, err = s.Ext.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountEntryExtensionV2Ext: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s AccountEntryExtensionV2) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AccountEntryExtensionV2) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*AccountEntryExtensionV2)(nil)
	_ encoding.BinaryUnmarshaler = (*AccountEntryExtensionV2)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s AccountEntryExtensionV2) xdrType() {}

var _ xdrType = (*AccountEntryExtensionV2)(nil)

// AccountEntryExtensionV1Ext is an XDR NestedUnion defines as:
//
//   union switch (int v)
//        {
//        case 0:
//            void;
//        case 2:
//            AccountEntryExtensionV2 v2;
//        }
//
type AccountEntryExtensionV1Ext struct {
	V  int32
	V2 *AccountEntryExtensionV2
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u AccountEntryExtensionV1Ext) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of AccountEntryExtensionV1Ext
func (u AccountEntryExtensionV1Ext) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	case 2:
		return "V2", true
	}
	return "-", false
}

// NewAccountEntryExtensionV1Ext creates a new  AccountEntryExtensionV1Ext.
func NewAccountEntryExtensionV1Ext(v int32, value interface{}) (result AccountEntryExtensionV1Ext, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	case 2:
		tv, ok := value.(AccountEntryExtensionV2)
		if !ok {
			err = fmt.Errorf("invalid value, must be AccountEntryExtensionV2")
			return
		}
		result.V2 = &tv
	}
	return
}

// MustV2 retrieves the V2 value from the union,
// panicing if the value is not set.
func (u AccountEntryExtensionV1Ext) MustV2() AccountEntryExtensionV2 {
	val, ok := u.GetV2()

	if !ok {
		panic("arm V2 is not set")
	}

	return val
}

// GetV2 retrieves the V2 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u AccountEntryExtensionV1Ext) GetV2() (result AccountEntryExtensionV2, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.V))

	if armName == "V2" {
		result = *u.V2
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u AccountEntryExtensionV1Ext) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	case 2:
		if err = (*u.V2).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union AccountEntryExtensionV1Ext", u.V)
}

var _ decoderFrom = (*AccountEntryExtensionV1Ext)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *AccountEntryExtensionV1Ext) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %s", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	case 2:
		u.V2 = new(AccountEntryExtensionV2)
		nTmp, err = (*u.V2).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding AccountEntryExtensionV2: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union AccountEntryExtensionV1Ext has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s AccountEntryExtensionV1Ext) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AccountEntryExtensionV1Ext) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*AccountEntryExtensionV1Ext)(nil)
	_ encoding.BinaryUnmarshaler = (*AccountEntryExtensionV1Ext)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s AccountEntryExtensionV1Ext) xdrType() {}

var _ xdrType = (*AccountEntryExtensionV1Ext)(nil)

// AccountEntryExtensionV1 is an XDR Struct defines as:
//
//   struct AccountEntryExtensionV1
//    {
//        Liabilities liabilities;
//
//        union switch (int v)
//        {
//        case 0:
//            void;
//        case 2:
//            AccountEntryExtensionV2 v2;
//        }
//        ext;
//    };
//
type AccountEntryExtensionV1 struct {
	Liabilities Liabilities
	Ext         AccountEntryExtensionV1Ext
}

// EncodeTo encodes this value using the Encoder.
func (s *AccountEntryExtensionV1) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Liabilities.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*AccountEntryExtensionV1)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *AccountEntryExtensionV1) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Liabilities.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Liabilities: %s", err)
	}
	nTmp, err = s.Ext.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountEntryExtensionV1Ext: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s AccountEntryExtensionV1) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AccountEntryExtensionV1) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*AccountEntryExtensionV1)(nil)
	_ encoding.BinaryUnmarshaler = (*AccountEntryExtensionV1)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s AccountEntryExtensionV1) xdrType() {}

var _ xdrType = (*AccountEntryExtensionV1)(nil)

// AccountEntryExt is an XDR NestedUnion defines as:
//
//   union switch (int v)
//        {
//        case 0:
//            void;
//        case 1:
//            AccountEntryExtensionV1 v1;
//        }
//
type AccountEntryExt struct {
	V  int32
	V1 *AccountEntryExtensionV1
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u AccountEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of AccountEntryExt
func (u AccountEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	case 1:
		return "V1", true
	}
	return "-", false
}

// NewAccountEntryExt creates a new  AccountEntryExt.
func NewAccountEntryExt(v int32, value interface{}) (result AccountEntryExt, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	case 1:
		tv, ok := value.(AccountEntryExtensionV1)
		if !ok {
			err = fmt.Errorf("invalid value, must be AccountEntryExtensionV1")
			return
		}
		result.V1 = &tv
	}
	return
}

// MustV1 retrieves the V1 value from the union,
// panicing if the value is not set.
func (u AccountEntryExt) MustV1() AccountEntryExtensionV1 {
	val, ok := u.GetV1()

	if !ok {
		panic("arm V1 is not set")
	}

	return val
}

// GetV1 retrieves the V1 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u AccountEntryExt) GetV1() (result AccountEntryExtensionV1, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.V))

	if armName == "V1" {
		result = *u.V1
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u AccountEntryExt) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	case 1:
		if err = (*u.V1).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union AccountEntryExt", u.V)
}

var _ decoderFrom = (*AccountEntryExt)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *AccountEntryExt) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %s", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	case 1:
		u.V1 = new(AccountEntryExtensionV1)
		nTmp, err = (*u.V1).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding AccountEntryExtensionV1: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union AccountEntryExt has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s AccountEntryExt) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AccountEntryExt) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*AccountEntryExt)(nil)
	_ encoding.BinaryUnmarshaler = (*AccountEntryExt)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s AccountEntryExt) xdrType() {}

var _ xdrType = (*AccountEntryExt)(nil)

// AccountEntry is an XDR Struct defines as:
//
//   struct AccountEntry
//    {
//        AccountID accountID;      // master public key for this account
//        int64 balance;            // in stroops
//        SequenceNumber seqNum;    // last sequence number used for this account
//        uint32 numSubEntries;     // number of sub-entries this account has
//                                  // drives the reserve
//        AccountID* inflationDest; // Account to vote for during inflation
//        uint32 flags;             // see AccountFlags
//
//        string32 homeDomain; // can be used for reverse federation and memo lookup
//
//        // fields used for signatures
//        // thresholds stores unsigned bytes: [weight of master|low|medium|high]
//        Thresholds thresholds;
//
//        Signer signers<MAX_SIGNERS>; // possible signers for this account
//
//        // reserved for future use
//        union switch (int v)
//        {
//        case 0:
//            void;
//        case 1:
//            AccountEntryExtensionV1 v1;
//        }
//        ext;
//    };
//
type AccountEntry struct {
	AccountId     AccountId
	Balance       Int64
	SeqNum        SequenceNumber
	NumSubEntries Uint32
	InflationDest *AccountId
	Flags         Uint32
	HomeDomain    String32
	Thresholds    Thresholds
	Signers       []Signer `xdrmaxsize:"20"`
	Ext           AccountEntryExt
}

// EncodeTo encodes this value using the Encoder.
func (s *AccountEntry) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.AccountId.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Balance.EncodeTo(e); err != nil {
		return err
	}
	if err = s.SeqNum.EncodeTo(e); err != nil {
		return err
	}
	if err = s.NumSubEntries.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeBool(s.InflationDest != nil); err != nil {
		return err
	}
	if s.InflationDest != nil {
		if err = (*s.InflationDest).EncodeTo(e); err != nil {
			return err
		}
	}
	if err = s.Flags.EncodeTo(e); err != nil {
		return err
	}
	if err = s.HomeDomain.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Thresholds.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeUint(uint32(len(s.Signers))); err != nil {
		return err
	}
	for i := 0; i < len(s.Signers); i++ {
		if err = s.Signers[i].EncodeTo(e); err != nil {
			return err
		}
	}
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*AccountEntry)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *AccountEntry) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.AccountId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountId: %s", err)
	}
	nTmp, err = s.Balance.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	nTmp, err = s.SeqNum.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SequenceNumber: %s", err)
	}
	nTmp, err = s.NumSubEntries.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	var b bool
	b, nTmp, err = d.DecodeBool()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountId: %s", err)
	}
	s.InflationDest = nil
	if b {
		s.InflationDest = new(AccountId)
		nTmp, err = s.InflationDest.DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding AccountId: %s", err)
		}
	}
	nTmp, err = s.Flags.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.HomeDomain.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding String32: %s", err)
	}
	nTmp, err = s.Thresholds.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Thresholds: %s", err)
	}
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Signer: %s", err)
	}
	if l > 20 {
		return n, fmt.Errorf("decoding Signer: data size (%d) exceeds size limit (20)", l)
	}
	s.Signers = nil
	if l > 0 {
		s.Signers = make([]Signer, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.Signers[i].DecodeFrom(d)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding Signer: %s", err)
			}
		}
	}
	nTmp, err = s.Ext.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountEntryExt: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s AccountEntry) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AccountEntry) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*AccountEntry)(nil)
	_ encoding.BinaryUnmarshaler = (*AccountEntry)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s AccountEntry) xdrType() {}

var _ xdrType = (*AccountEntry)(nil)

// TrustLineFlags is an XDR Enum defines as:
//
//   enum TrustLineFlags
//    {
//        // issuer has authorized account to perform transactions with its credit
//        AUTHORIZED_FLAG = 1,
//        // issuer has authorized account to maintain and reduce liabilities for its
//        // credit
//        AUTHORIZED_TO_MAINTAIN_LIABILITIES_FLAG = 2,
//        // issuer has specified that it may clawback its credit, and that claimable
//        // balances created with its credit may also be clawed back
//        TRUSTLINE_CLAWBACK_ENABLED_FLAG = 4
//    };
//
type TrustLineFlags int32

const (
	TrustLineFlagsAuthorizedFlag                      TrustLineFlags = 1
	TrustLineFlagsAuthorizedToMaintainLiabilitiesFlag TrustLineFlags = 2
	TrustLineFlagsTrustlineClawbackEnabledFlag        TrustLineFlags = 4
)

var trustLineFlagsMap = map[int32]string{
	1: "TrustLineFlagsAuthorizedFlag",
	2: "TrustLineFlagsAuthorizedToMaintainLiabilitiesFlag",
	4: "TrustLineFlagsTrustlineClawbackEnabledFlag",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for TrustLineFlags
func (e TrustLineFlags) ValidEnum(v int32) bool {
	_, ok := trustLineFlagsMap[v]
	return ok
}

// String returns the name of `e`
func (e TrustLineFlags) String() string {
	name, _ := trustLineFlagsMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e TrustLineFlags) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := trustLineFlagsMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid TrustLineFlags enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*TrustLineFlags)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *TrustLineFlags) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding TrustLineFlags: %s", err)
	}
	if _, ok := trustLineFlagsMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid TrustLineFlags enum value", v)
	}
	*e = TrustLineFlags(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TrustLineFlags) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TrustLineFlags) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TrustLineFlags)(nil)
	_ encoding.BinaryUnmarshaler = (*TrustLineFlags)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TrustLineFlags) xdrType() {}

var _ xdrType = (*TrustLineFlags)(nil)

// MaskTrustlineFlags is an XDR Const defines as:
//
//   const MASK_TRUSTLINE_FLAGS = 1;
//
const MaskTrustlineFlags = 1

// MaskTrustlineFlagsV13 is an XDR Const defines as:
//
//   const MASK_TRUSTLINE_FLAGS_V13 = 3;
//
const MaskTrustlineFlagsV13 = 3

// MaskTrustlineFlagsV17 is an XDR Const defines as:
//
//   const MASK_TRUSTLINE_FLAGS_V17 = 7;
//
const MaskTrustlineFlagsV17 = 7

// LiquidityPoolType is an XDR Enum defines as:
//
//   enum LiquidityPoolType
//    {
//        LIQUIDITY_POOL_CONSTANT_PRODUCT = 0
//    };
//
type LiquidityPoolType int32

const (
	LiquidityPoolTypeLiquidityPoolConstantProduct LiquidityPoolType = 0
)

var liquidityPoolTypeMap = map[int32]string{
	0: "LiquidityPoolTypeLiquidityPoolConstantProduct",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for LiquidityPoolType
func (e LiquidityPoolType) ValidEnum(v int32) bool {
	_, ok := liquidityPoolTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e LiquidityPoolType) String() string {
	name, _ := liquidityPoolTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e LiquidityPoolType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := liquidityPoolTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid LiquidityPoolType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*LiquidityPoolType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *LiquidityPoolType) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding LiquidityPoolType: %s", err)
	}
	if _, ok := liquidityPoolTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid LiquidityPoolType enum value", v)
	}
	*e = LiquidityPoolType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LiquidityPoolType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LiquidityPoolType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LiquidityPoolType)(nil)
	_ encoding.BinaryUnmarshaler = (*LiquidityPoolType)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LiquidityPoolType) xdrType() {}

var _ xdrType = (*LiquidityPoolType)(nil)

// TrustLineAsset is an XDR Union defines as:
//
//   union TrustLineAsset switch (AssetType type)
//    {
//    case ASSET_TYPE_NATIVE: // Not credit
//        void;
//
//    case ASSET_TYPE_CREDIT_ALPHANUM4:
//        AlphaNum4 alphaNum4;
//
//    case ASSET_TYPE_CREDIT_ALPHANUM12:
//        AlphaNum12 alphaNum12;
//
//    case ASSET_TYPE_POOL_SHARE:
//        PoolID liquidityPoolID;
//
//        // add other asset types here in the future
//    };
//
type TrustLineAsset struct {
	Type            AssetType
	AlphaNum4       *AlphaNum4
	AlphaNum12      *AlphaNum12
	LiquidityPoolId *PoolId
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TrustLineAsset) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TrustLineAsset
func (u TrustLineAsset) ArmForSwitch(sw int32) (string, bool) {
	switch AssetType(sw) {
	case AssetTypeAssetTypeNative:
		return "", true
	case AssetTypeAssetTypeCreditAlphanum4:
		return "AlphaNum4", true
	case AssetTypeAssetTypeCreditAlphanum12:
		return "AlphaNum12", true
	case AssetTypeAssetTypePoolShare:
		return "LiquidityPoolId", true
	}
	return "-", false
}

// NewTrustLineAsset creates a new  TrustLineAsset.
func NewTrustLineAsset(aType AssetType, value interface{}) (result TrustLineAsset, err error) {
	result.Type = aType
	switch AssetType(aType) {
	case AssetTypeAssetTypeNative:
		// void
	case AssetTypeAssetTypeCreditAlphanum4:
		tv, ok := value.(AlphaNum4)
		if !ok {
			err = fmt.Errorf("invalid value, must be AlphaNum4")
			return
		}
		result.AlphaNum4 = &tv
	case AssetTypeAssetTypeCreditAlphanum12:
		tv, ok := value.(AlphaNum12)
		if !ok {
			err = fmt.Errorf("invalid value, must be AlphaNum12")
			return
		}
		result.AlphaNum12 = &tv
	case AssetTypeAssetTypePoolShare:
		tv, ok := value.(PoolId)
		if !ok {
			err = fmt.Errorf("invalid value, must be PoolId")
			return
		}
		result.LiquidityPoolId = &tv
	}
	return
}

// MustAlphaNum4 retrieves the AlphaNum4 value from the union,
// panicing if the value is not set.
func (u TrustLineAsset) MustAlphaNum4() AlphaNum4 {
	val, ok := u.GetAlphaNum4()

	if !ok {
		panic("arm AlphaNum4 is not set")
	}

	return val
}

// GetAlphaNum4 retrieves the AlphaNum4 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u TrustLineAsset) GetAlphaNum4() (result AlphaNum4, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AlphaNum4" {
		result = *u.AlphaNum4
		ok = true
	}

	return
}

// MustAlphaNum12 retrieves the AlphaNum12 value from the union,
// panicing if the value is not set.
func (u TrustLineAsset) MustAlphaNum12() AlphaNum12 {
	val, ok := u.GetAlphaNum12()

	if !ok {
		panic("arm AlphaNum12 is not set")
	}

	return val
}

// GetAlphaNum12 retrieves the AlphaNum12 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u TrustLineAsset) GetAlphaNum12() (result AlphaNum12, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AlphaNum12" {
		result = *u.AlphaNum12
		ok = true
	}

	return
}

// MustLiquidityPoolId retrieves the LiquidityPoolId value from the union,
// panicing if the value is not set.
func (u TrustLineAsset) MustLiquidityPoolId() PoolId {
	val, ok := u.GetLiquidityPoolId()

	if !ok {
		panic("arm LiquidityPoolId is not set")
	}

	return val
}

// GetLiquidityPoolId retrieves the LiquidityPoolId value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u TrustLineAsset) GetLiquidityPoolId() (result PoolId, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "LiquidityPoolId" {
		result = *u.LiquidityPoolId
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u TrustLineAsset) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch AssetType(u.Type) {
	case AssetTypeAssetTypeNative:
		// Void
		return nil
	case AssetTypeAssetTypeCreditAlphanum4:
		if err = (*u.AlphaNum4).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case AssetTypeAssetTypeCreditAlphanum12:
		if err = (*u.AlphaNum12).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case AssetTypeAssetTypePoolShare:
		if err = (*u.LiquidityPoolId).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (AssetType) switch value '%d' is not valid for union TrustLineAsset", u.Type)
}

var _ decoderFrom = (*TrustLineAsset)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *TrustLineAsset) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AssetType: %s", err)
	}
	switch AssetType(u.Type) {
	case AssetTypeAssetTypeNative:
		// Void
		return n, nil
	case AssetTypeAssetTypeCreditAlphanum4:
		u.AlphaNum4 = new(AlphaNum4)
		nTmp, err = (*u.AlphaNum4).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding AlphaNum4: %s", err)
		}
		return n, nil
	case AssetTypeAssetTypeCreditAlphanum12:
		u.AlphaNum12 = new(AlphaNum12)
		nTmp, err = (*u.AlphaNum12).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding AlphaNum12: %s", err)
		}
		return n, nil
	case AssetTypeAssetTypePoolShare:
		u.LiquidityPoolId = new(PoolId)
		nTmp, err = (*u.LiquidityPoolId).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding PoolId: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union TrustLineAsset has invalid Type (AssetType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TrustLineAsset) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TrustLineAsset) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TrustLineAsset)(nil)
	_ encoding.BinaryUnmarshaler = (*TrustLineAsset)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TrustLineAsset) xdrType() {}

var _ xdrType = (*TrustLineAsset)(nil)

// TrustLineEntryExtensionV2Ext is an XDR NestedUnion defines as:
//
//   union switch (int v)
//        {
//        case 0:
//            void;
//        }
//
type TrustLineEntryExtensionV2Ext struct {
	V int32
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TrustLineEntryExtensionV2Ext) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TrustLineEntryExtensionV2Ext
func (u TrustLineEntryExtensionV2Ext) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	}
	return "-", false
}

// NewTrustLineEntryExtensionV2Ext creates a new  TrustLineEntryExtensionV2Ext.
func NewTrustLineEntryExtensionV2Ext(v int32, value interface{}) (result TrustLineEntryExtensionV2Ext, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u TrustLineEntryExtensionV2Ext) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union TrustLineEntryExtensionV2Ext", u.V)
}

var _ decoderFrom = (*TrustLineEntryExtensionV2Ext)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *TrustLineEntryExtensionV2Ext) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %s", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	}
	return n, fmt.Errorf("union TrustLineEntryExtensionV2Ext has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TrustLineEntryExtensionV2Ext) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TrustLineEntryExtensionV2Ext) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TrustLineEntryExtensionV2Ext)(nil)
	_ encoding.BinaryUnmarshaler = (*TrustLineEntryExtensionV2Ext)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TrustLineEntryExtensionV2Ext) xdrType() {}

var _ xdrType = (*TrustLineEntryExtensionV2Ext)(nil)

// TrustLineEntryExtensionV2 is an XDR Struct defines as:
//
//   struct TrustLineEntryExtensionV2
//    {
//        int32 liquidityPoolUseCount;
//
//        union switch (int v)
//        {
//        case 0:
//            void;
//        }
//        ext;
//    };
//
type TrustLineEntryExtensionV2 struct {
	LiquidityPoolUseCount Int32
	Ext                   TrustLineEntryExtensionV2Ext
}

// EncodeTo encodes this value using the Encoder.
func (s *TrustLineEntryExtensionV2) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.LiquidityPoolUseCount.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*TrustLineEntryExtensionV2)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *TrustLineEntryExtensionV2) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.LiquidityPoolUseCount.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int32: %s", err)
	}
	nTmp, err = s.Ext.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TrustLineEntryExtensionV2Ext: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TrustLineEntryExtensionV2) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TrustLineEntryExtensionV2) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TrustLineEntryExtensionV2)(nil)
	_ encoding.BinaryUnmarshaler = (*TrustLineEntryExtensionV2)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TrustLineEntryExtensionV2) xdrType() {}

var _ xdrType = (*TrustLineEntryExtensionV2)(nil)

// TrustLineEntryV1Ext is an XDR NestedUnion defines as:
//
//   union switch (int v)
//                {
//                case 0:
//                    void;
//                case 2:
//                    TrustLineEntryExtensionV2 v2;
//                }
//
type TrustLineEntryV1Ext struct {
	V  int32
	V2 *TrustLineEntryExtensionV2
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TrustLineEntryV1Ext) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TrustLineEntryV1Ext
func (u TrustLineEntryV1Ext) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	case 2:
		return "V2", true
	}
	return "-", false
}

// NewTrustLineEntryV1Ext creates a new  TrustLineEntryV1Ext.
func NewTrustLineEntryV1Ext(v int32, value interface{}) (result TrustLineEntryV1Ext, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	case 2:
		tv, ok := value.(TrustLineEntryExtensionV2)
		if !ok {
			err = fmt.Errorf("invalid value, must be TrustLineEntryExtensionV2")
			return
		}
		result.V2 = &tv
	}
	return
}

// MustV2 retrieves the V2 value from the union,
// panicing if the value is not set.
func (u TrustLineEntryV1Ext) MustV2() TrustLineEntryExtensionV2 {
	val, ok := u.GetV2()

	if !ok {
		panic("arm V2 is not set")
	}

	return val
}

// GetV2 retrieves the V2 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u TrustLineEntryV1Ext) GetV2() (result TrustLineEntryExtensionV2, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.V))

	if armName == "V2" {
		result = *u.V2
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u TrustLineEntryV1Ext) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	case 2:
		if err = (*u.V2).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union TrustLineEntryV1Ext", u.V)
}

var _ decoderFrom = (*TrustLineEntryV1Ext)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *TrustLineEntryV1Ext) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %s", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	case 2:
		u.V2 = new(TrustLineEntryExtensionV2)
		nTmp, err = (*u.V2).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding TrustLineEntryExtensionV2: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union TrustLineEntryV1Ext has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TrustLineEntryV1Ext) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TrustLineEntryV1Ext) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TrustLineEntryV1Ext)(nil)
	_ encoding.BinaryUnmarshaler = (*TrustLineEntryV1Ext)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TrustLineEntryV1Ext) xdrType() {}

var _ xdrType = (*TrustLineEntryV1Ext)(nil)

// TrustLineEntryV1 is an XDR NestedStruct defines as:
//
//   struct
//            {
//                Liabilities liabilities;
//
//                union switch (int v)
//                {
//                case 0:
//                    void;
//                case 2:
//                    TrustLineEntryExtensionV2 v2;
//                }
//                ext;
//            }
//
type TrustLineEntryV1 struct {
	Liabilities Liabilities
	Ext         TrustLineEntryV1Ext
}

// EncodeTo encodes this value using the Encoder.
func (s *TrustLineEntryV1) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Liabilities.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*TrustLineEntryV1)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *TrustLineEntryV1) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Liabilities.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Liabilities: %s", err)
	}
	nTmp, err = s.Ext.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TrustLineEntryV1Ext: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TrustLineEntryV1) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TrustLineEntryV1) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TrustLineEntryV1)(nil)
	_ encoding.BinaryUnmarshaler = (*TrustLineEntryV1)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TrustLineEntryV1) xdrType() {}

var _ xdrType = (*TrustLineEntryV1)(nil)

// TrustLineEntryExt is an XDR NestedUnion defines as:
//
//   union switch (int v)
//        {
//        case 0:
//            void;
//        case 1:
//            struct
//            {
//                Liabilities liabilities;
//
//                union switch (int v)
//                {
//                case 0:
//                    void;
//                case 2:
//                    TrustLineEntryExtensionV2 v2;
//                }
//                ext;
//            } v1;
//        }
//
type TrustLineEntryExt struct {
	V  int32
	V1 *TrustLineEntryV1
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TrustLineEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TrustLineEntryExt
func (u TrustLineEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	case 1:
		return "V1", true
	}
	return "-", false
}

// NewTrustLineEntryExt creates a new  TrustLineEntryExt.
func NewTrustLineEntryExt(v int32, value interface{}) (result TrustLineEntryExt, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	case 1:
		tv, ok := value.(TrustLineEntryV1)
		if !ok {
			err = fmt.Errorf("invalid value, must be TrustLineEntryV1")
			return
		}
		result.V1 = &tv
	}
	return
}

// MustV1 retrieves the V1 value from the union,
// panicing if the value is not set.
func (u TrustLineEntryExt) MustV1() TrustLineEntryV1 {
	val, ok := u.GetV1()

	if !ok {
		panic("arm V1 is not set")
	}

	return val
}

// GetV1 retrieves the V1 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u TrustLineEntryExt) GetV1() (result TrustLineEntryV1, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.V))

	if armName == "V1" {
		result = *u.V1
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u TrustLineEntryExt) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	case 1:
		if err = (*u.V1).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union TrustLineEntryExt", u.V)
}

var _ decoderFrom = (*TrustLineEntryExt)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *TrustLineEntryExt) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %s", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	case 1:
		u.V1 = new(TrustLineEntryV1)
		nTmp, err = (*u.V1).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding TrustLineEntryV1: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union TrustLineEntryExt has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TrustLineEntryExt) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TrustLineEntryExt) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TrustLineEntryExt)(nil)
	_ encoding.BinaryUnmarshaler = (*TrustLineEntryExt)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TrustLineEntryExt) xdrType() {}

var _ xdrType = (*TrustLineEntryExt)(nil)

// TrustLineEntry is an XDR Struct defines as:
//
//   struct TrustLineEntry
//    {
//        AccountID accountID;  // account this trustline belongs to
//        TrustLineAsset asset; // type of asset (with issuer)
//        int64 balance;        // how much of this asset the user has.
//                              // Asset defines the unit for this;
//
//        int64 limit;  // balance cannot be above this
//        uint32 flags; // see TrustLineFlags
//
//        // reserved for future use
//        union switch (int v)
//        {
//        case 0:
//            void;
//        case 1:
//            struct
//            {
//                Liabilities liabilities;
//
//                union switch (int v)
//                {
//                case 0:
//                    void;
//                case 2:
//                    TrustLineEntryExtensionV2 v2;
//                }
//                ext;
//            } v1;
//        }
//        ext;
//    };
//
type TrustLineEntry struct {
	AccountId AccountId
	Asset     TrustLineAsset
	Balance   Int64
	Limit     Int64
	Flags     Uint32
	Ext       TrustLineEntryExt
}

// EncodeTo encodes this value using the Encoder.
func (s *TrustLineEntry) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.AccountId.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Asset.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Balance.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Limit.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Flags.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*TrustLineEntry)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *TrustLineEntry) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.AccountId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountId: %s", err)
	}
	nTmp, err = s.Asset.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TrustLineAsset: %s", err)
	}
	nTmp, err = s.Balance.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	nTmp, err = s.Limit.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	nTmp, err = s.Flags.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.Ext.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TrustLineEntryExt: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TrustLineEntry) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TrustLineEntry) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TrustLineEntry)(nil)
	_ encoding.BinaryUnmarshaler = (*TrustLineEntry)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TrustLineEntry) xdrType() {}

var _ xdrType = (*TrustLineEntry)(nil)

// OfferEntryFlags is an XDR Enum defines as:
//
//   enum OfferEntryFlags
//    {
//        // issuer has authorized account to perform transactions with its credit
//        PASSIVE_FLAG = 1
//    };
//
type OfferEntryFlags int32

const (
	OfferEntryFlagsPassiveFlag OfferEntryFlags = 1
)

var offerEntryFlagsMap = map[int32]string{
	1: "OfferEntryFlagsPassiveFlag",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for OfferEntryFlags
func (e OfferEntryFlags) ValidEnum(v int32) bool {
	_, ok := offerEntryFlagsMap[v]
	return ok
}

// String returns the name of `e`
func (e OfferEntryFlags) String() string {
	name, _ := offerEntryFlagsMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e OfferEntryFlags) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := offerEntryFlagsMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid OfferEntryFlags enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*OfferEntryFlags)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *OfferEntryFlags) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding OfferEntryFlags: %s", err)
	}
	if _, ok := offerEntryFlagsMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid OfferEntryFlags enum value", v)
	}
	*e = OfferEntryFlags(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s OfferEntryFlags) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *OfferEntryFlags) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*OfferEntryFlags)(nil)
	_ encoding.BinaryUnmarshaler = (*OfferEntryFlags)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s OfferEntryFlags) xdrType() {}

var _ xdrType = (*OfferEntryFlags)(nil)

// MaskOfferentryFlags is an XDR Const defines as:
//
//   const MASK_OFFERENTRY_FLAGS = 1;
//
const MaskOfferentryFlags = 1

// OfferEntryExt is an XDR NestedUnion defines as:
//
//   union switch (int v)
//        {
//        case 0:
//            void;
//        }
//
type OfferEntryExt struct {
	V int32
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u OfferEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of OfferEntryExt
func (u OfferEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	}
	return "-", false
}

// NewOfferEntryExt creates a new  OfferEntryExt.
func NewOfferEntryExt(v int32, value interface{}) (result OfferEntryExt, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u OfferEntryExt) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union OfferEntryExt", u.V)
}

var _ decoderFrom = (*OfferEntryExt)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *OfferEntryExt) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %s", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	}
	return n, fmt.Errorf("union OfferEntryExt has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s OfferEntryExt) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *OfferEntryExt) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*OfferEntryExt)(nil)
	_ encoding.BinaryUnmarshaler = (*OfferEntryExt)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s OfferEntryExt) xdrType() {}

var _ xdrType = (*OfferEntryExt)(nil)

// OfferEntry is an XDR Struct defines as:
//
//   struct OfferEntry
//    {
//        AccountID sellerID;
//        int64 offerID;
//        Asset selling; // A
//        Asset buying;  // B
//        int64 amount;  // amount of A
//
//        /* price for this offer:
//            price of A in terms of B
//            price=AmountB/AmountA=priceNumerator/priceDenominator
//            price is after fees
//        */
//        Price price;
//        uint32 flags; // see OfferEntryFlags
//
//        // reserved for future use
//        union switch (int v)
//        {
//        case 0:
//            void;
//        }
//        ext;
//    };
//
type OfferEntry struct {
	SellerId AccountId
	OfferId  Int64
	Selling  Asset
	Buying   Asset
	Amount   Int64
	Price    Price
	Flags    Uint32
	Ext      OfferEntryExt
}

// EncodeTo encodes this value using the Encoder.
func (s *OfferEntry) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.SellerId.EncodeTo(e); err != nil {
		return err
	}
	if err = s.OfferId.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Selling.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Buying.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Amount.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Price.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Flags.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*OfferEntry)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *OfferEntry) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.SellerId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountId: %s", err)
	}
	nTmp, err = s.OfferId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	nTmp, err = s.Selling.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Asset: %s", err)
	}
	nTmp, err = s.Buying.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Asset: %s", err)
	}
	nTmp, err = s.Amount.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	nTmp, err = s.Price.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Price: %s", err)
	}
	nTmp, err = s.Flags.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.Ext.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding OfferEntryExt: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s OfferEntry) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *OfferEntry) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*OfferEntry)(nil)
	_ encoding.BinaryUnmarshaler = (*OfferEntry)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s OfferEntry) xdrType() {}

var _ xdrType = (*OfferEntry)(nil)

// DataEntryExt is an XDR NestedUnion defines as:
//
//   union switch (int v)
//        {
//        case 0:
//            void;
//        }
//
type DataEntryExt struct {
	V int32
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u DataEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of DataEntryExt
func (u DataEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	}
	return "-", false
}

// NewDataEntryExt creates a new  DataEntryExt.
func NewDataEntryExt(v int32, value interface{}) (result DataEntryExt, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u DataEntryExt) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union DataEntryExt", u.V)
}

var _ decoderFrom = (*DataEntryExt)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *DataEntryExt) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %s", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	}
	return n, fmt.Errorf("union DataEntryExt has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s DataEntryExt) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *DataEntryExt) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*DataEntryExt)(nil)
	_ encoding.BinaryUnmarshaler = (*DataEntryExt)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s DataEntryExt) xdrType() {}

var _ xdrType = (*DataEntryExt)(nil)

// DataEntry is an XDR Struct defines as:
//
//   struct DataEntry
//    {
//        AccountID accountID; // account this data belongs to
//        string64 dataName;
//        DataValue dataValue;
//
//        // reserved for future use
//        union switch (int v)
//        {
//        case 0:
//            void;
//        }
//        ext;
//    };
//
type DataEntry struct {
	AccountId AccountId
	DataName  String64
	DataValue DataValue
	Ext       DataEntryExt
}

// EncodeTo encodes this value using the Encoder.
func (s *DataEntry) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.AccountId.EncodeTo(e); err != nil {
		return err
	}
	if err = s.DataName.EncodeTo(e); err != nil {
		return err
	}
	if err = s.DataValue.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*DataEntry)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *DataEntry) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.AccountId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountId: %s", err)
	}
	nTmp, err = s.DataName.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding String64: %s", err)
	}
	nTmp, err = s.DataValue.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding DataValue: %s", err)
	}
	nTmp, err = s.Ext.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding DataEntryExt: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s DataEntry) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *DataEntry) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*DataEntry)(nil)
	_ encoding.BinaryUnmarshaler = (*DataEntry)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s DataEntry) xdrType() {}

var _ xdrType = (*DataEntry)(nil)

// ClaimPredicateType is an XDR Enum defines as:
//
//   enum ClaimPredicateType
//    {
//        CLAIM_PREDICATE_UNCONDITIONAL = 0,
//        CLAIM_PREDICATE_AND = 1,
//        CLAIM_PREDICATE_OR = 2,
//        CLAIM_PREDICATE_NOT = 3,
//        CLAIM_PREDICATE_BEFORE_ABSOLUTE_TIME = 4,
//        CLAIM_PREDICATE_BEFORE_RELATIVE_TIME = 5
//    };
//
type ClaimPredicateType int32

const (
	ClaimPredicateTypeClaimPredicateUnconditional      ClaimPredicateType = 0
	ClaimPredicateTypeClaimPredicateAnd                ClaimPredicateType = 1
	ClaimPredicateTypeClaimPredicateOr                 ClaimPredicateType = 2
	ClaimPredicateTypeClaimPredicateNot                ClaimPredicateType = 3
	ClaimPredicateTypeClaimPredicateBeforeAbsoluteTime ClaimPredicateType = 4
	ClaimPredicateTypeClaimPredicateBeforeRelativeTime ClaimPredicateType = 5
)

var claimPredicateTypeMap = map[int32]string{
	0: "ClaimPredicateTypeClaimPredicateUnconditional",
	1: "ClaimPredicateTypeClaimPredicateAnd",
	2: "ClaimPredicateTypeClaimPredicateOr",
	3: "ClaimPredicateTypeClaimPredicateNot",
	4: "ClaimPredicateTypeClaimPredicateBeforeAbsoluteTime",
	5: "ClaimPredicateTypeClaimPredicateBeforeRelativeTime",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ClaimPredicateType
func (e ClaimPredicateType) ValidEnum(v int32) bool {
	_, ok := claimPredicateTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e ClaimPredicateType) String() string {
	name, _ := claimPredicateTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e ClaimPredicateType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := claimPredicateTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid ClaimPredicateType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*ClaimPredicateType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *ClaimPredicateType) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding ClaimPredicateType: %s", err)
	}
	if _, ok := claimPredicateTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid ClaimPredicateType enum value", v)
	}
	*e = ClaimPredicateType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ClaimPredicateType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ClaimPredicateType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ClaimPredicateType)(nil)
	_ encoding.BinaryUnmarshaler = (*ClaimPredicateType)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ClaimPredicateType) xdrType() {}

var _ xdrType = (*ClaimPredicateType)(nil)

// ClaimPredicate is an XDR Union defines as:
//
//   union ClaimPredicate switch (ClaimPredicateType type)
//    {
//    case CLAIM_PREDICATE_UNCONDITIONAL:
//        void;
//    case CLAIM_PREDICATE_AND:
//        ClaimPredicate andPredicates<2>;
//    case CLAIM_PREDICATE_OR:
//        ClaimPredicate orPredicates<2>;
//    case CLAIM_PREDICATE_NOT:
//        ClaimPredicate* notPredicate;
//    case CLAIM_PREDICATE_BEFORE_ABSOLUTE_TIME:
//        int64 absBefore; // Predicate will be true if closeTime < absBefore
//    case CLAIM_PREDICATE_BEFORE_RELATIVE_TIME:
//        int64 relBefore; // Seconds since closeTime of the ledger in which the
//                         // ClaimableBalanceEntry was created
//    };
//
type ClaimPredicate struct {
	Type          ClaimPredicateType
	AndPredicates *[]ClaimPredicate `xdrmaxsize:"2"`
	OrPredicates  *[]ClaimPredicate `xdrmaxsize:"2"`
	NotPredicate  **ClaimPredicate
	AbsBefore     *Int64
	RelBefore     *Int64
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ClaimPredicate) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ClaimPredicate
func (u ClaimPredicate) ArmForSwitch(sw int32) (string, bool) {
	switch ClaimPredicateType(sw) {
	case ClaimPredicateTypeClaimPredicateUnconditional:
		return "", true
	case ClaimPredicateTypeClaimPredicateAnd:
		return "AndPredicates", true
	case ClaimPredicateTypeClaimPredicateOr:
		return "OrPredicates", true
	case ClaimPredicateTypeClaimPredicateNot:
		return "NotPredicate", true
	case ClaimPredicateTypeClaimPredicateBeforeAbsoluteTime:
		return "AbsBefore", true
	case ClaimPredicateTypeClaimPredicateBeforeRelativeTime:
		return "RelBefore", true
	}
	return "-", false
}

// NewClaimPredicate creates a new  ClaimPredicate.
func NewClaimPredicate(aType ClaimPredicateType, value interface{}) (result ClaimPredicate, err error) {
	result.Type = aType
	switch ClaimPredicateType(aType) {
	case ClaimPredicateTypeClaimPredicateUnconditional:
		// void
	case ClaimPredicateTypeClaimPredicateAnd:
		tv, ok := value.([]ClaimPredicate)
		if !ok {
			err = fmt.Errorf("invalid value, must be []ClaimPredicate")
			return
		}
		result.AndPredicates = &tv
	case ClaimPredicateTypeClaimPredicateOr:
		tv, ok := value.([]ClaimPredicate)
		if !ok {
			err = fmt.Errorf("invalid value, must be []ClaimPredicate")
			return
		}
		result.OrPredicates = &tv
	case ClaimPredicateTypeClaimPredicateNot:
		tv, ok := value.(*ClaimPredicate)
		if !ok {
			err = fmt.Errorf("invalid value, must be *ClaimPredicate")
			return
		}
		result.NotPredicate = &tv
	case ClaimPredicateTypeClaimPredicateBeforeAbsoluteTime:
		tv, ok := value.(Int64)
		if !ok {
			err = fmt.Errorf("invalid value, must be Int64")
			return
		}
		result.AbsBefore = &tv
	case ClaimPredicateTypeClaimPredicateBeforeRelativeTime:
		tv, ok := value.(Int64)
		if !ok {
			err = fmt.Errorf("invalid value, must be Int64")
			return
		}
		result.RelBefore = &tv
	}
	return
}

// MustAndPredicates retrieves the AndPredicates value from the union,
// panicing if the value is not set.
func (u ClaimPredicate) MustAndPredicates() []ClaimPredicate {
	val, ok := u.GetAndPredicates()

	if !ok {
		panic("arm AndPredicates is not set")
	}

	return val
}

// GetAndPredicates retrieves the AndPredicates value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ClaimPredicate) GetAndPredicates() (result []ClaimPredicate, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AndPredicates" {
		result = *u.AndPredicates
		ok = true
	}

	return
}

// MustOrPredicates retrieves the OrPredicates value from the union,
// panicing if the value is not set.
func (u ClaimPredicate) MustOrPredicates() []ClaimPredicate {
	val, ok := u.GetOrPredicates()

	if !ok {
		panic("arm OrPredicates is not set")
	}

	return val
}

// GetOrPredicates retrieves the OrPredicates value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ClaimPredicate) GetOrPredicates() (result []ClaimPredicate, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "OrPredicates" {
		result = *u.OrPredicates
		ok = true
	}

	return
}

// MustNotPredicate retrieves the NotPredicate value from the union,
// panicing if the value is not set.
func (u ClaimPredicate) MustNotPredicate() *ClaimPredicate {
	val, ok := u.GetNotPredicate()

	if !ok {
		panic("arm NotPredicate is not set")
	}

	return val
}

// GetNotPredicate retrieves the NotPredicate value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ClaimPredicate) GetNotPredicate() (result *ClaimPredicate, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "NotPredicate" {
		result = *u.NotPredicate
		ok = true
	}

	return
}

// MustAbsBefore retrieves the AbsBefore value from the union,
// panicing if the value is not set.
func (u ClaimPredicate) MustAbsBefore() Int64 {
	val, ok := u.GetAbsBefore()

	if !ok {
		panic("arm AbsBefore is not set")
	}

	return val
}

// GetAbsBefore retrieves the AbsBefore value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ClaimPredicate) GetAbsBefore() (result Int64, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AbsBefore" {
		result = *u.AbsBefore
		ok = true
	}

	return
}

// MustRelBefore retrieves the RelBefore value from the union,
// panicing if the value is not set.
func (u ClaimPredicate) MustRelBefore() Int64 {
	val, ok := u.GetRelBefore()

	if !ok {
		panic("arm RelBefore is not set")
	}

	return val
}

// GetRelBefore retrieves the RelBefore value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ClaimPredicate) GetRelBefore() (result Int64, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "RelBefore" {
		result = *u.RelBefore
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u ClaimPredicate) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch ClaimPredicateType(u.Type) {
	case ClaimPredicateTypeClaimPredicateUnconditional:
		// Void
		return nil
	case ClaimPredicateTypeClaimPredicateAnd:
		if _, err = e.EncodeUint(uint32(len((*u.AndPredicates)))); err != nil {
			return err
		}
		for i := 0; i < len((*u.AndPredicates)); i++ {
			if err = (*u.AndPredicates)[i].EncodeTo(e); err != nil {
				return err
			}
		}
		return nil
	case ClaimPredicateTypeClaimPredicateOr:
		if _, err = e.EncodeUint(uint32(len((*u.OrPredicates)))); err != nil {
			return err
		}
		for i := 0; i < len((*u.OrPredicates)); i++ {
			if err = (*u.OrPredicates)[i].EncodeTo(e); err != nil {
				return err
			}
		}
		return nil
	case ClaimPredicateTypeClaimPredicateNot:
		if _, err = e.EncodeBool((*u.NotPredicate) != nil); err != nil {
			return err
		}
		if (*u.NotPredicate) != nil {
			if err = (*(*u.NotPredicate)).EncodeTo(e); err != nil {
				return err
			}
		}
		return nil
	case ClaimPredicateTypeClaimPredicateBeforeAbsoluteTime:
		if err = (*u.AbsBefore).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case ClaimPredicateTypeClaimPredicateBeforeRelativeTime:
		if err = (*u.RelBefore).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (ClaimPredicateType) switch value '%d' is not valid for union ClaimPredicate", u.Type)
}

var _ decoderFrom = (*ClaimPredicate)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *ClaimPredicate) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ClaimPredicateType: %s", err)
	}
	switch ClaimPredicateType(u.Type) {
	case ClaimPredicateTypeClaimPredicateUnconditional:
		// Void
		return n, nil
	case ClaimPredicateTypeClaimPredicateAnd:
		u.AndPredicates = new([]ClaimPredicate)
		var l uint32
		l, nTmp, err = d.DecodeUint()
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ClaimPredicate: %s", err)
		}
		if l > 2 {
			return n, fmt.Errorf("decoding ClaimPredicate: data size (%d) exceeds size limit (2)", l)
		}
		(*u.AndPredicates) = nil
		if l > 0 {
			(*u.AndPredicates) = make([]ClaimPredicate, l)
			for i := uint32(0); i < l; i++ {
				nTmp, err = (*u.AndPredicates)[i].DecodeFrom(d)
				n += nTmp
				if err != nil {
					return n, fmt.Errorf("decoding ClaimPredicate: %s", err)
				}
			}
		}
		return n, nil
	case ClaimPredicateTypeClaimPredicateOr:
		u.OrPredicates = new([]ClaimPredicate)
		var l uint32
		l, nTmp, err = d.DecodeUint()
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ClaimPredicate: %s", err)
		}
		if l > 2 {
			return n, fmt.Errorf("decoding ClaimPredicate: data size (%d) exceeds size limit (2)", l)
		}
		(*u.OrPredicates) = nil
		if l > 0 {
			(*u.OrPredicates) = make([]ClaimPredicate, l)
			for i := uint32(0); i < l; i++ {
				nTmp, err = (*u.OrPredicates)[i].DecodeFrom(d)
				n += nTmp
				if err != nil {
					return n, fmt.Errorf("decoding ClaimPredicate: %s", err)
				}
			}
		}
		return n, nil
	case ClaimPredicateTypeClaimPredicateNot:
		u.NotPredicate = new(*ClaimPredicate)
		var b bool
		b, nTmp, err = d.DecodeBool()
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ClaimPredicate: %s", err)
		}
		(*u.NotPredicate) = nil
		if b {
			(*u.NotPredicate) = new(ClaimPredicate)
			nTmp, err = (*u.NotPredicate).DecodeFrom(d)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding ClaimPredicate: %s", err)
			}
		}
		return n, nil
	case ClaimPredicateTypeClaimPredicateBeforeAbsoluteTime:
		u.AbsBefore = new(Int64)
		nTmp, err = (*u.AbsBefore).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Int64: %s", err)
		}
		return n, nil
	case ClaimPredicateTypeClaimPredicateBeforeRelativeTime:
		u.RelBefore = new(Int64)
		nTmp, err = (*u.RelBefore).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Int64: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union ClaimPredicate has invalid Type (ClaimPredicateType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ClaimPredicate) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ClaimPredicate) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ClaimPredicate)(nil)
	_ encoding.BinaryUnmarshaler = (*ClaimPredicate)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ClaimPredicate) xdrType() {}

var _ xdrType = (*ClaimPredicate)(nil)

// ClaimantType is an XDR Enum defines as:
//
//   enum ClaimantType
//    {
//        CLAIMANT_TYPE_V0 = 0
//    };
//
type ClaimantType int32

const (
	ClaimantTypeClaimantTypeV0 ClaimantType = 0
)

var claimantTypeMap = map[int32]string{
	0: "ClaimantTypeClaimantTypeV0",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ClaimantType
func (e ClaimantType) ValidEnum(v int32) bool {
	_, ok := claimantTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e ClaimantType) String() string {
	name, _ := claimantTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e ClaimantType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := claimantTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid ClaimantType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*ClaimantType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *ClaimantType) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding ClaimantType: %s", err)
	}
	if _, ok := claimantTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid ClaimantType enum value", v)
	}
	*e = ClaimantType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ClaimantType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ClaimantType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ClaimantType)(nil)
	_ encoding.BinaryUnmarshaler = (*ClaimantType)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ClaimantType) xdrType() {}

var _ xdrType = (*ClaimantType)(nil)

// ClaimantV0 is an XDR NestedStruct defines as:
//
//   struct
//        {
//            AccountID destination;    // The account that can use this condition
//            ClaimPredicate predicate; // Claimable if predicate is true
//        }
//
type ClaimantV0 struct {
	Destination AccountId
	Predicate   ClaimPredicate
}

// EncodeTo encodes this value using the Encoder.
func (s *ClaimantV0) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Destination.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Predicate.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*ClaimantV0)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *ClaimantV0) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Destination.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountId: %s", err)
	}
	nTmp, err = s.Predicate.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ClaimPredicate: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ClaimantV0) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ClaimantV0) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ClaimantV0)(nil)
	_ encoding.BinaryUnmarshaler = (*ClaimantV0)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ClaimantV0) xdrType() {}

var _ xdrType = (*ClaimantV0)(nil)

// Claimant is an XDR Union defines as:
//
//   union Claimant switch (ClaimantType type)
//    {
//    case CLAIMANT_TYPE_V0:
//        struct
//        {
//            AccountID destination;    // The account that can use this condition
//            ClaimPredicate predicate; // Claimable if predicate is true
//        } v0;
//    };
//
type Claimant struct {
	Type ClaimantType
	V0   *ClaimantV0
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u Claimant) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of Claimant
func (u Claimant) ArmForSwitch(sw int32) (string, bool) {
	switch ClaimantType(sw) {
	case ClaimantTypeClaimantTypeV0:
		return "V0", true
	}
	return "-", false
}

// NewClaimant creates a new  Claimant.
func NewClaimant(aType ClaimantType, value interface{}) (result Claimant, err error) {
	result.Type = aType
	switch ClaimantType(aType) {
	case ClaimantTypeClaimantTypeV0:
		tv, ok := value.(ClaimantV0)
		if !ok {
			err = fmt.Errorf("invalid value, must be ClaimantV0")
			return
		}
		result.V0 = &tv
	}
	return
}

// MustV0 retrieves the V0 value from the union,
// panicing if the value is not set.
func (u Claimant) MustV0() ClaimantV0 {
	val, ok := u.GetV0()

	if !ok {
		panic("arm V0 is not set")
	}

	return val
}

// GetV0 retrieves the V0 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u Claimant) GetV0() (result ClaimantV0, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "V0" {
		result = *u.V0
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u Claimant) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch ClaimantType(u.Type) {
	case ClaimantTypeClaimantTypeV0:
		if err = (*u.V0).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (ClaimantType) switch value '%d' is not valid for union Claimant", u.Type)
}

var _ decoderFrom = (*Claimant)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *Claimant) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ClaimantType: %s", err)
	}
	switch ClaimantType(u.Type) {
	case ClaimantTypeClaimantTypeV0:
		u.V0 = new(ClaimantV0)
		nTmp, err = (*u.V0).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ClaimantV0: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union Claimant has invalid Type (ClaimantType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Claimant) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Claimant) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Claimant)(nil)
	_ encoding.BinaryUnmarshaler = (*Claimant)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s Claimant) xdrType() {}

var _ xdrType = (*Claimant)(nil)

// ClaimableBalanceIdType is an XDR Enum defines as:
//
//   enum ClaimableBalanceIDType
//    {
//        CLAIMABLE_BALANCE_ID_TYPE_V0 = 0
//    };
//
type ClaimableBalanceIdType int32

const (
	ClaimableBalanceIdTypeClaimableBalanceIdTypeV0 ClaimableBalanceIdType = 0
)

var claimableBalanceIdTypeMap = map[int32]string{
	0: "ClaimableBalanceIdTypeClaimableBalanceIdTypeV0",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ClaimableBalanceIdType
func (e ClaimableBalanceIdType) ValidEnum(v int32) bool {
	_, ok := claimableBalanceIdTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e ClaimableBalanceIdType) String() string {
	name, _ := claimableBalanceIdTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e ClaimableBalanceIdType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := claimableBalanceIdTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid ClaimableBalanceIdType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*ClaimableBalanceIdType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *ClaimableBalanceIdType) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding ClaimableBalanceIdType: %s", err)
	}
	if _, ok := claimableBalanceIdTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid ClaimableBalanceIdType enum value", v)
	}
	*e = ClaimableBalanceIdType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ClaimableBalanceIdType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ClaimableBalanceIdType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ClaimableBalanceIdType)(nil)
	_ encoding.BinaryUnmarshaler = (*ClaimableBalanceIdType)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ClaimableBalanceIdType) xdrType() {}

var _ xdrType = (*ClaimableBalanceIdType)(nil)

// ClaimableBalanceId is an XDR Union defines as:
//
//   union ClaimableBalanceID switch (ClaimableBalanceIDType type)
//    {
//    case CLAIMABLE_BALANCE_ID_TYPE_V0:
//        Hash v0;
//    };
//
type ClaimableBalanceId struct {
	Type ClaimableBalanceIdType
	V0   *Hash
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ClaimableBalanceId) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ClaimableBalanceId
func (u ClaimableBalanceId) ArmForSwitch(sw int32) (string, bool) {
	switch ClaimableBalanceIdType(sw) {
	case ClaimableBalanceIdTypeClaimableBalanceIdTypeV0:
		return "V0", true
	}
	return "-", false
}

// NewClaimableBalanceId creates a new  ClaimableBalanceId.
func NewClaimableBalanceId(aType ClaimableBalanceIdType, value interface{}) (result ClaimableBalanceId, err error) {
	result.Type = aType
	switch ClaimableBalanceIdType(aType) {
	case ClaimableBalanceIdTypeClaimableBalanceIdTypeV0:
		tv, ok := value.(Hash)
		if !ok {
			err = fmt.Errorf("invalid value, must be Hash")
			return
		}
		result.V0 = &tv
	}
	return
}

// MustV0 retrieves the V0 value from the union,
// panicing if the value is not set.
func (u ClaimableBalanceId) MustV0() Hash {
	val, ok := u.GetV0()

	if !ok {
		panic("arm V0 is not set")
	}

	return val
}

// GetV0 retrieves the V0 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ClaimableBalanceId) GetV0() (result Hash, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "V0" {
		result = *u.V0
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u ClaimableBalanceId) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch ClaimableBalanceIdType(u.Type) {
	case ClaimableBalanceIdTypeClaimableBalanceIdTypeV0:
		if err = (*u.V0).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (ClaimableBalanceIdType) switch value '%d' is not valid for union ClaimableBalanceId", u.Type)
}

var _ decoderFrom = (*ClaimableBalanceId)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *ClaimableBalanceId) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ClaimableBalanceIdType: %s", err)
	}
	switch ClaimableBalanceIdType(u.Type) {
	case ClaimableBalanceIdTypeClaimableBalanceIdTypeV0:
		u.V0 = new(Hash)
		nTmp, err = (*u.V0).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Hash: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union ClaimableBalanceId has invalid Type (ClaimableBalanceIdType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ClaimableBalanceId) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ClaimableBalanceId) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ClaimableBalanceId)(nil)
	_ encoding.BinaryUnmarshaler = (*ClaimableBalanceId)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ClaimableBalanceId) xdrType() {}

var _ xdrType = (*ClaimableBalanceId)(nil)

// ClaimableBalanceFlags is an XDR Enum defines as:
//
//   enum ClaimableBalanceFlags
//    {
//        // If set, the issuer account of the asset held by the claimable balance may
//        // clawback the claimable balance
//        CLAIMABLE_BALANCE_CLAWBACK_ENABLED_FLAG = 0x1
//    };
//
type ClaimableBalanceFlags int32

const (
	ClaimableBalanceFlagsClaimableBalanceClawbackEnabledFlag ClaimableBalanceFlags = 1
)

var claimableBalanceFlagsMap = map[int32]string{
	1: "ClaimableBalanceFlagsClaimableBalanceClawbackEnabledFlag",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ClaimableBalanceFlags
func (e ClaimableBalanceFlags) ValidEnum(v int32) bool {
	_, ok := claimableBalanceFlagsMap[v]
	return ok
}

// String returns the name of `e`
func (e ClaimableBalanceFlags) String() string {
	name, _ := claimableBalanceFlagsMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e ClaimableBalanceFlags) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := claimableBalanceFlagsMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid ClaimableBalanceFlags enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*ClaimableBalanceFlags)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *ClaimableBalanceFlags) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding ClaimableBalanceFlags: %s", err)
	}
	if _, ok := claimableBalanceFlagsMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid ClaimableBalanceFlags enum value", v)
	}
	*e = ClaimableBalanceFlags(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ClaimableBalanceFlags) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ClaimableBalanceFlags) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ClaimableBalanceFlags)(nil)
	_ encoding.BinaryUnmarshaler = (*ClaimableBalanceFlags)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ClaimableBalanceFlags) xdrType() {}

var _ xdrType = (*ClaimableBalanceFlags)(nil)

// MaskClaimableBalanceFlags is an XDR Const defines as:
//
//   const MASK_CLAIMABLE_BALANCE_FLAGS = 0x1;
//
const MaskClaimableBalanceFlags = 0x1

// ClaimableBalanceEntryExtensionV1Ext is an XDR NestedUnion defines as:
//
//   union switch (int v)
//        {
//        case 0:
//            void;
//        }
//
type ClaimableBalanceEntryExtensionV1Ext struct {
	V int32
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ClaimableBalanceEntryExtensionV1Ext) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ClaimableBalanceEntryExtensionV1Ext
func (u ClaimableBalanceEntryExtensionV1Ext) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	}
	return "-", false
}

// NewClaimableBalanceEntryExtensionV1Ext creates a new  ClaimableBalanceEntryExtensionV1Ext.
func NewClaimableBalanceEntryExtensionV1Ext(v int32, value interface{}) (result ClaimableBalanceEntryExtensionV1Ext, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u ClaimableBalanceEntryExtensionV1Ext) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union ClaimableBalanceEntryExtensionV1Ext", u.V)
}

var _ decoderFrom = (*ClaimableBalanceEntryExtensionV1Ext)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *ClaimableBalanceEntryExtensionV1Ext) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %s", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	}
	return n, fmt.Errorf("union ClaimableBalanceEntryExtensionV1Ext has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ClaimableBalanceEntryExtensionV1Ext) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ClaimableBalanceEntryExtensionV1Ext) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ClaimableBalanceEntryExtensionV1Ext)(nil)
	_ encoding.BinaryUnmarshaler = (*ClaimableBalanceEntryExtensionV1Ext)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ClaimableBalanceEntryExtensionV1Ext) xdrType() {}

var _ xdrType = (*ClaimableBalanceEntryExtensionV1Ext)(nil)

// ClaimableBalanceEntryExtensionV1 is an XDR Struct defines as:
//
//   struct ClaimableBalanceEntryExtensionV1
//    {
//        union switch (int v)
//        {
//        case 0:
//            void;
//        }
//        ext;
//
//        uint32 flags; // see ClaimableBalanceFlags
//    };
//
type ClaimableBalanceEntryExtensionV1 struct {
	Ext   ClaimableBalanceEntryExtensionV1Ext
	Flags Uint32
}

// EncodeTo encodes this value using the Encoder.
func (s *ClaimableBalanceEntryExtensionV1) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Flags.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*ClaimableBalanceEntryExtensionV1)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *ClaimableBalanceEntryExtensionV1) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Ext.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ClaimableBalanceEntryExtensionV1Ext: %s", err)
	}
	nTmp, err = s.Flags.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ClaimableBalanceEntryExtensionV1) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ClaimableBalanceEntryExtensionV1) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ClaimableBalanceEntryExtensionV1)(nil)
	_ encoding.BinaryUnmarshaler = (*ClaimableBalanceEntryExtensionV1)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ClaimableBalanceEntryExtensionV1) xdrType() {}

var _ xdrType = (*ClaimableBalanceEntryExtensionV1)(nil)

// ClaimableBalanceEntryExt is an XDR NestedUnion defines as:
//
//   union switch (int v)
//        {
//        case 0:
//            void;
//        case 1:
//            ClaimableBalanceEntryExtensionV1 v1;
//        }
//
type ClaimableBalanceEntryExt struct {
	V  int32
	V1 *ClaimableBalanceEntryExtensionV1
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ClaimableBalanceEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ClaimableBalanceEntryExt
func (u ClaimableBalanceEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	case 1:
		return "V1", true
	}
	return "-", false
}

// NewClaimableBalanceEntryExt creates a new  ClaimableBalanceEntryExt.
func NewClaimableBalanceEntryExt(v int32, value interface{}) (result ClaimableBalanceEntryExt, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	case 1:
		tv, ok := value.(ClaimableBalanceEntryExtensionV1)
		if !ok {
			err = fmt.Errorf("invalid value, must be ClaimableBalanceEntryExtensionV1")
			return
		}
		result.V1 = &tv
	}
	return
}

// MustV1 retrieves the V1 value from the union,
// panicing if the value is not set.
func (u ClaimableBalanceEntryExt) MustV1() ClaimableBalanceEntryExtensionV1 {
	val, ok := u.GetV1()

	if !ok {
		panic("arm V1 is not set")
	}

	return val
}

// GetV1 retrieves the V1 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ClaimableBalanceEntryExt) GetV1() (result ClaimableBalanceEntryExtensionV1, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.V))

	if armName == "V1" {
		result = *u.V1
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u ClaimableBalanceEntryExt) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	case 1:
		if err = (*u.V1).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union ClaimableBalanceEntryExt", u.V)
}

var _ decoderFrom = (*ClaimableBalanceEntryExt)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *ClaimableBalanceEntryExt) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %s", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	case 1:
		u.V1 = new(ClaimableBalanceEntryExtensionV1)
		nTmp, err = (*u.V1).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ClaimableBalanceEntryExtensionV1: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union ClaimableBalanceEntryExt has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ClaimableBalanceEntryExt) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ClaimableBalanceEntryExt) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ClaimableBalanceEntryExt)(nil)
	_ encoding.BinaryUnmarshaler = (*ClaimableBalanceEntryExt)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ClaimableBalanceEntryExt) xdrType() {}

var _ xdrType = (*ClaimableBalanceEntryExt)(nil)

// ClaimableBalanceEntry is an XDR Struct defines as:
//
//   struct ClaimableBalanceEntry
//    {
//        // Unique identifier for this ClaimableBalanceEntry
//        ClaimableBalanceID balanceID;
//
//        // List of claimants with associated predicate
//        Claimant claimants<10>;
//
//        // Any asset including native
//        Asset asset;
//
//        // Amount of asset
//        int64 amount;
//
//        // reserved for future use
//        union switch (int v)
//        {
//        case 0:
//            void;
//        case 1:
//            ClaimableBalanceEntryExtensionV1 v1;
//        }
//        ext;
//    };
//
type ClaimableBalanceEntry struct {
	BalanceId ClaimableBalanceId
	Claimants []Claimant `xdrmaxsize:"10"`
	Asset     Asset
	Amount    Int64
	Ext       ClaimableBalanceEntryExt
}

// EncodeTo encodes this value using the Encoder.
func (s *ClaimableBalanceEntry) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.BalanceId.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeUint(uint32(len(s.Claimants))); err != nil {
		return err
	}
	for i := 0; i < len(s.Claimants); i++ {
		if err = s.Claimants[i].EncodeTo(e); err != nil {
			return err
		}
	}
	if err = s.Asset.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Amount.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*ClaimableBalanceEntry)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *ClaimableBalanceEntry) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.BalanceId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ClaimableBalanceId: %s", err)
	}
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Claimant: %s", err)
	}
	if l > 10 {
		return n, fmt.Errorf("decoding Claimant: data size (%d) exceeds size limit (10)", l)
	}
	s.Claimants = nil
	if l > 0 {
		s.Claimants = make([]Claimant, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.Claimants[i].DecodeFrom(d)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding Claimant: %s", err)
			}
		}
	}
	nTmp, err = s.Asset.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Asset: %s", err)
	}
	nTmp, err = s.Amount.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	nTmp, err = s.Ext.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ClaimableBalanceEntryExt: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ClaimableBalanceEntry) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ClaimableBalanceEntry) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ClaimableBalanceEntry)(nil)
	_ encoding.BinaryUnmarshaler = (*ClaimableBalanceEntry)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ClaimableBalanceEntry) xdrType() {}

var _ xdrType = (*ClaimableBalanceEntry)(nil)

// LiquidityPoolConstantProductParameters is an XDR Struct defines as:
//
//   struct LiquidityPoolConstantProductParameters
//    {
//        Asset assetA; // assetA < assetB
//        Asset assetB;
//        int32 fee;    // Fee is in basis points, so the actual rate is (fee/100)%
//    };
//
type LiquidityPoolConstantProductParameters struct {
	AssetA Asset
	AssetB Asset
	Fee    Int32
}

// EncodeTo encodes this value using the Encoder.
func (s *LiquidityPoolConstantProductParameters) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.AssetA.EncodeTo(e); err != nil {
		return err
	}
	if err = s.AssetB.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Fee.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*LiquidityPoolConstantProductParameters)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *LiquidityPoolConstantProductParameters) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.AssetA.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Asset: %s", err)
	}
	nTmp, err = s.AssetB.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Asset: %s", err)
	}
	nTmp, err = s.Fee.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int32: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LiquidityPoolConstantProductParameters) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LiquidityPoolConstantProductParameters) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LiquidityPoolConstantProductParameters)(nil)
	_ encoding.BinaryUnmarshaler = (*LiquidityPoolConstantProductParameters)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LiquidityPoolConstantProductParameters) xdrType() {}

var _ xdrType = (*LiquidityPoolConstantProductParameters)(nil)

// LiquidityPoolEntryConstantProduct is an XDR NestedStruct defines as:
//
//   struct
//            {
//                LiquidityPoolConstantProductParameters params;
//
//                int64 reserveA;        // amount of A in the pool
//                int64 reserveB;        // amount of B in the pool
//                int64 totalPoolShares; // total number of pool shares issued
//                int64 poolSharesTrustLineCount; // number of trust lines for the associated pool shares
//            }
//
type LiquidityPoolEntryConstantProduct struct {
	Params                   LiquidityPoolConstantProductParameters
	ReserveA                 Int64
	ReserveB                 Int64
	TotalPoolShares          Int64
	PoolSharesTrustLineCount Int64
}

// EncodeTo encodes this value using the Encoder.
func (s *LiquidityPoolEntryConstantProduct) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Params.EncodeTo(e); err != nil {
		return err
	}
	if err = s.ReserveA.EncodeTo(e); err != nil {
		return err
	}
	if err = s.ReserveB.EncodeTo(e); err != nil {
		return err
	}
	if err = s.TotalPoolShares.EncodeTo(e); err != nil {
		return err
	}
	if err = s.PoolSharesTrustLineCount.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*LiquidityPoolEntryConstantProduct)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *LiquidityPoolEntryConstantProduct) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Params.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LiquidityPoolConstantProductParameters: %s", err)
	}
	nTmp, err = s.ReserveA.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	nTmp, err = s.ReserveB.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	nTmp, err = s.TotalPoolShares.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	nTmp, err = s.PoolSharesTrustLineCount.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LiquidityPoolEntryConstantProduct) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LiquidityPoolEntryConstantProduct) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LiquidityPoolEntryConstantProduct)(nil)
	_ encoding.BinaryUnmarshaler = (*LiquidityPoolEntryConstantProduct)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LiquidityPoolEntryConstantProduct) xdrType() {}

var _ xdrType = (*LiquidityPoolEntryConstantProduct)(nil)

// LiquidityPoolEntryBody is an XDR NestedUnion defines as:
//
//   union switch (LiquidityPoolType type)
//        {
//        case LIQUIDITY_POOL_CONSTANT_PRODUCT:
//            struct
//            {
//                LiquidityPoolConstantProductParameters params;
//
//                int64 reserveA;        // amount of A in the pool
//                int64 reserveB;        // amount of B in the pool
//                int64 totalPoolShares; // total number of pool shares issued
//                int64 poolSharesTrustLineCount; // number of trust lines for the associated pool shares
//            } constantProduct;
//        }
//
type LiquidityPoolEntryBody struct {
	Type            LiquidityPoolType
	ConstantProduct *LiquidityPoolEntryConstantProduct
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LiquidityPoolEntryBody) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LiquidityPoolEntryBody
func (u LiquidityPoolEntryBody) ArmForSwitch(sw int32) (string, bool) {
	switch LiquidityPoolType(sw) {
	case LiquidityPoolTypeLiquidityPoolConstantProduct:
		return "ConstantProduct", true
	}
	return "-", false
}

// NewLiquidityPoolEntryBody creates a new  LiquidityPoolEntryBody.
func NewLiquidityPoolEntryBody(aType LiquidityPoolType, value interface{}) (result LiquidityPoolEntryBody, err error) {
	result.Type = aType
	switch LiquidityPoolType(aType) {
	case LiquidityPoolTypeLiquidityPoolConstantProduct:
		tv, ok := value.(LiquidityPoolEntryConstantProduct)
		if !ok {
			err = fmt.Errorf("invalid value, must be LiquidityPoolEntryConstantProduct")
			return
		}
		result.ConstantProduct = &tv
	}
	return
}

// MustConstantProduct retrieves the ConstantProduct value from the union,
// panicing if the value is not set.
func (u LiquidityPoolEntryBody) MustConstantProduct() LiquidityPoolEntryConstantProduct {
	val, ok := u.GetConstantProduct()

	if !ok {
		panic("arm ConstantProduct is not set")
	}

	return val
}

// GetConstantProduct retrieves the ConstantProduct value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LiquidityPoolEntryBody) GetConstantProduct() (result LiquidityPoolEntryConstantProduct, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ConstantProduct" {
		result = *u.ConstantProduct
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u LiquidityPoolEntryBody) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch LiquidityPoolType(u.Type) {
	case LiquidityPoolTypeLiquidityPoolConstantProduct:
		if err = (*u.ConstantProduct).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (LiquidityPoolType) switch value '%d' is not valid for union LiquidityPoolEntryBody", u.Type)
}

var _ decoderFrom = (*LiquidityPoolEntryBody)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *LiquidityPoolEntryBody) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LiquidityPoolType: %s", err)
	}
	switch LiquidityPoolType(u.Type) {
	case LiquidityPoolTypeLiquidityPoolConstantProduct:
		u.ConstantProduct = new(LiquidityPoolEntryConstantProduct)
		nTmp, err = (*u.ConstantProduct).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding LiquidityPoolEntryConstantProduct: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union LiquidityPoolEntryBody has invalid Type (LiquidityPoolType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LiquidityPoolEntryBody) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LiquidityPoolEntryBody) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LiquidityPoolEntryBody)(nil)
	_ encoding.BinaryUnmarshaler = (*LiquidityPoolEntryBody)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LiquidityPoolEntryBody) xdrType() {}

var _ xdrType = (*LiquidityPoolEntryBody)(nil)

// LiquidityPoolEntry is an XDR Struct defines as:
//
//   struct LiquidityPoolEntry
//    {
//        PoolID liquidityPoolID;
//
//        union switch (LiquidityPoolType type)
//        {
//        case LIQUIDITY_POOL_CONSTANT_PRODUCT:
//            struct
//            {
//                LiquidityPoolConstantProductParameters params;
//
//                int64 reserveA;        // amount of A in the pool
//                int64 reserveB;        // amount of B in the pool
//                int64 totalPoolShares; // total number of pool shares issued
//                int64 poolSharesTrustLineCount; // number of trust lines for the associated pool shares
//            } constantProduct;
//        }
//        body;
//    };
//
type LiquidityPoolEntry struct {
	LiquidityPoolId PoolId
	Body            LiquidityPoolEntryBody
}

// EncodeTo encodes this value using the Encoder.
func (s *LiquidityPoolEntry) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.LiquidityPoolId.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Body.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*LiquidityPoolEntry)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *LiquidityPoolEntry) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.LiquidityPoolId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding PoolId: %s", err)
	}
	nTmp, err = s.Body.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LiquidityPoolEntryBody: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LiquidityPoolEntry) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LiquidityPoolEntry) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LiquidityPoolEntry)(nil)
	_ encoding.BinaryUnmarshaler = (*LiquidityPoolEntry)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LiquidityPoolEntry) xdrType() {}

var _ xdrType = (*LiquidityPoolEntry)(nil)

// LedgerEntryExtensionV1Ext is an XDR NestedUnion defines as:
//
//   union switch (int v)
//        {
//        case 0:
//            void;
//        }
//
type LedgerEntryExtensionV1Ext struct {
	V int32
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerEntryExtensionV1Ext) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerEntryExtensionV1Ext
func (u LedgerEntryExtensionV1Ext) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	}
	return "-", false
}

// NewLedgerEntryExtensionV1Ext creates a new  LedgerEntryExtensionV1Ext.
func NewLedgerEntryExtensionV1Ext(v int32, value interface{}) (result LedgerEntryExtensionV1Ext, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u LedgerEntryExtensionV1Ext) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union LedgerEntryExtensionV1Ext", u.V)
}

var _ decoderFrom = (*LedgerEntryExtensionV1Ext)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *LedgerEntryExtensionV1Ext) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %s", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	}
	return n, fmt.Errorf("union LedgerEntryExtensionV1Ext has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerEntryExtensionV1Ext) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerEntryExtensionV1Ext) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerEntryExtensionV1Ext)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerEntryExtensionV1Ext)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LedgerEntryExtensionV1Ext) xdrType() {}

var _ xdrType = (*LedgerEntryExtensionV1Ext)(nil)

// LedgerEntryExtensionV1 is an XDR Struct defines as:
//
//   struct LedgerEntryExtensionV1
//    {
//        SponsorshipDescriptor sponsoringID;
//
//        union switch (int v)
//        {
//        case 0:
//            void;
//        }
//        ext;
//    };
//
type LedgerEntryExtensionV1 struct {
	SponsoringId SponsorshipDescriptor
	Ext          LedgerEntryExtensionV1Ext
}

// EncodeTo encodes this value using the Encoder.
func (s *LedgerEntryExtensionV1) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeBool(s.SponsoringId != nil); err != nil {
		return err
	}
	if s.SponsoringId != nil {
		if err = (*s.SponsoringId).EncodeTo(e); err != nil {
			return err
		}
	}
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*LedgerEntryExtensionV1)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *LedgerEntryExtensionV1) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	var b bool
	b, nTmp, err = d.DecodeBool()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SponsorshipDescriptor: %s", err)
	}
	s.SponsoringId = nil
	if b {
		s.SponsoringId = new(AccountId)
		nTmp, err = s.SponsoringId.DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding SponsorshipDescriptor: %s", err)
		}
	}
	nTmp, err = s.Ext.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LedgerEntryExtensionV1Ext: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerEntryExtensionV1) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerEntryExtensionV1) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerEntryExtensionV1)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerEntryExtensionV1)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LedgerEntryExtensionV1) xdrType() {}

var _ xdrType = (*LedgerEntryExtensionV1)(nil)

// LedgerEntryData is an XDR NestedUnion defines as:
//
//   union switch (LedgerEntryType type)
//        {
//        case ACCOUNT:
//            AccountEntry account;
//        case TRUSTLINE:
//            TrustLineEntry trustLine;
//        case OFFER:
//            OfferEntry offer;
//        case DATA:
//            DataEntry data;
//        case CLAIMABLE_BALANCE:
//            ClaimableBalanceEntry claimableBalance;
//        case LIQUIDITY_POOL:
//            LiquidityPoolEntry liquidityPool;
//        }
//
type LedgerEntryData struct {
	Type             LedgerEntryType
	Account          *AccountEntry
	TrustLine        *TrustLineEntry
	Offer            *OfferEntry
	Data             *DataEntry
	ClaimableBalance *ClaimableBalanceEntry
	LiquidityPool    *LiquidityPoolEntry
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerEntryData) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerEntryData
func (u LedgerEntryData) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerEntryType(sw) {
	case LedgerEntryTypeAccount:
		return "Account", true
	case LedgerEntryTypeTrustline:
		return "TrustLine", true
	case LedgerEntryTypeOffer:
		return "Offer", true
	case LedgerEntryTypeData:
		return "Data", true
	case LedgerEntryTypeClaimableBalance:
		return "ClaimableBalance", true
	case LedgerEntryTypeLiquidityPool:
		return "LiquidityPool", true
	}
	return "-", false
}

// NewLedgerEntryData creates a new  LedgerEntryData.
func NewLedgerEntryData(aType LedgerEntryType, value interface{}) (result LedgerEntryData, err error) {
	result.Type = aType
	switch LedgerEntryType(aType) {
	case LedgerEntryTypeAccount:
		tv, ok := value.(AccountEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be AccountEntry")
			return
		}
		result.Account = &tv
	case LedgerEntryTypeTrustline:
		tv, ok := value.(TrustLineEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be TrustLineEntry")
			return
		}
		result.TrustLine = &tv
	case LedgerEntryTypeOffer:
		tv, ok := value.(OfferEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be OfferEntry")
			return
		}
		result.Offer = &tv
	case LedgerEntryTypeData:
		tv, ok := value.(DataEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be DataEntry")
			return
		}
		result.Data = &tv
	case LedgerEntryTypeClaimableBalance:
		tv, ok := value.(ClaimableBalanceEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be ClaimableBalanceEntry")
			return
		}
		result.ClaimableBalance = &tv
	case LedgerEntryTypeLiquidityPool:
		tv, ok := value.(LiquidityPoolEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be LiquidityPoolEntry")
			return
		}
		result.LiquidityPool = &tv
	}
	return
}

// MustAccount retrieves the Account value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustAccount() AccountEntry {
	val, ok := u.GetAccount()

	if !ok {
		panic("arm Account is not set")
	}

	return val
}

// GetAccount retrieves the Account value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetAccount() (result AccountEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Account" {
		result = *u.Account
		ok = true
	}

	return
}

// MustTrustLine retrieves the TrustLine value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustTrustLine() TrustLineEntry {
	val, ok := u.GetTrustLine()

	if !ok {
		panic("arm TrustLine is not set")
	}

	return val
}

// GetTrustLine retrieves the TrustLine value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetTrustLine() (result TrustLineEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "TrustLine" {
		result = *u.TrustLine
		ok = true
	}

	return
}

// MustOffer retrieves the Offer value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustOffer() OfferEntry {
	val, ok := u.GetOffer()

	if !ok {
		panic("arm Offer is not set")
	}

	return val
}

// GetOffer retrieves the Offer value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetOffer() (result OfferEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Offer" {
		result = *u.Offer
		ok = true
	}

	return
}

// MustData retrieves the Data value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustData() DataEntry {
	val, ok := u.GetData()

	if !ok {
		panic("arm Data is not set")
	}

	return val
}

// GetData retrieves the Data value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetData() (result DataEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Data" {
		result = *u.Data
		ok = true
	}

	return
}

// MustClaimableBalance retrieves the ClaimableBalance value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustClaimableBalance() ClaimableBalanceEntry {
	val, ok := u.GetClaimableBalance()

	if !ok {
		panic("arm ClaimableBalance is not set")
	}

	return val
}

// GetClaimableBalance retrieves the ClaimableBalance value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetClaimableBalance() (result ClaimableBalanceEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ClaimableBalance" {
		result = *u.ClaimableBalance
		ok = true
	}

	return
}

// MustLiquidityPool retrieves the LiquidityPool value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustLiquidityPool() LiquidityPoolEntry {
	val, ok := u.GetLiquidityPool()

	if !ok {
		panic("arm LiquidityPool is not set")
	}

	return val
}

// GetLiquidityPool retrieves the LiquidityPool value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetLiquidityPool() (result LiquidityPoolEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "LiquidityPool" {
		result = *u.LiquidityPool
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u LedgerEntryData) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch LedgerEntryType(u.Type) {
	case LedgerEntryTypeAccount:
		if err = (*u.Account).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case LedgerEntryTypeTrustline:
		if err = (*u.TrustLine).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case LedgerEntryTypeOffer:
		if err = (*u.Offer).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case LedgerEntryTypeData:
		if err = (*u.Data).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case LedgerEntryTypeClaimableBalance:
		if err = (*u.ClaimableBalance).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case LedgerEntryTypeLiquidityPool:
		if err = (*u.LiquidityPool).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (LedgerEntryType) switch value '%d' is not valid for union LedgerEntryData", u.Type)
}

var _ decoderFrom = (*LedgerEntryData)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *LedgerEntryData) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LedgerEntryType: %s", err)
	}
	switch LedgerEntryType(u.Type) {
	case LedgerEntryTypeAccount:
		u.Account = new(AccountEntry)
		nTmp, err = (*u.Account).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding AccountEntry: %s", err)
		}
		return n, nil
	case LedgerEntryTypeTrustline:
		u.TrustLine = new(TrustLineEntry)
		nTmp, err = (*u.TrustLine).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding TrustLineEntry: %s", err)
		}
		return n, nil
	case LedgerEntryTypeOffer:
		u.Offer = new(OfferEntry)
		nTmp, err = (*u.Offer).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding OfferEntry: %s", err)
		}
		return n, nil
	case LedgerEntryTypeData:
		u.Data = new(DataEntry)
		nTmp, err = (*u.Data).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding DataEntry: %s", err)
		}
		return n, nil
	case LedgerEntryTypeClaimableBalance:
		u.ClaimableBalance = new(ClaimableBalanceEntry)
		nTmp, err = (*u.ClaimableBalance).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ClaimableBalanceEntry: %s", err)
		}
		return n, nil
	case LedgerEntryTypeLiquidityPool:
		u.LiquidityPool = new(LiquidityPoolEntry)
		nTmp, err = (*u.LiquidityPool).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding LiquidityPoolEntry: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union LedgerEntryData has invalid Type (LedgerEntryType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerEntryData) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerEntryData) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerEntryData)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerEntryData)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LedgerEntryData) xdrType() {}

var _ xdrType = (*LedgerEntryData)(nil)

// LedgerEntryExt is an XDR NestedUnion defines as:
//
//   union switch (int v)
//        {
//        case 0:
//            void;
//        case 1:
//            LedgerEntryExtensionV1 v1;
//        }
//
type LedgerEntryExt struct {
	V  int32
	V1 *LedgerEntryExtensionV1
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerEntryExt
func (u LedgerEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	case 1:
		return "V1", true
	}
	return "-", false
}

// NewLedgerEntryExt creates a new  LedgerEntryExt.
func NewLedgerEntryExt(v int32, value interface{}) (result LedgerEntryExt, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	case 1:
		tv, ok := value.(LedgerEntryExtensionV1)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerEntryExtensionV1")
			return
		}
		result.V1 = &tv
	}
	return
}

// MustV1 retrieves the V1 value from the union,
// panicing if the value is not set.
func (u LedgerEntryExt) MustV1() LedgerEntryExtensionV1 {
	val, ok := u.GetV1()

	if !ok {
		panic("arm V1 is not set")
	}

	return val
}

// GetV1 retrieves the V1 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryExt) GetV1() (result LedgerEntryExtensionV1, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.V))

	if armName == "V1" {
		result = *u.V1
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u LedgerEntryExt) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	case 1:
		if err = (*u.V1).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union LedgerEntryExt", u.V)
}

var _ decoderFrom = (*LedgerEntryExt)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *LedgerEntryExt) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %s", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	case 1:
		u.V1 = new(LedgerEntryExtensionV1)
		nTmp, err = (*u.V1).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding LedgerEntryExtensionV1: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union LedgerEntryExt has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerEntryExt) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerEntryExt) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerEntryExt)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerEntryExt)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LedgerEntryExt) xdrType() {}

var _ xdrType = (*LedgerEntryExt)(nil)

// LedgerEntry is an XDR Struct defines as:
//
//   struct LedgerEntry
//    {
//        uint32 lastModifiedLedgerSeq; // ledger the LedgerEntry was last changed
//
//        union switch (LedgerEntryType type)
//        {
//        case ACCOUNT:
//            AccountEntry account;
//        case TRUSTLINE:
//            TrustLineEntry trustLine;
//        case OFFER:
//            OfferEntry offer;
//        case DATA:
//            DataEntry data;
//        case CLAIMABLE_BALANCE:
//            ClaimableBalanceEntry claimableBalance;
//        case LIQUIDITY_POOL:
//            LiquidityPoolEntry liquidityPool;
//        }
//        data;
//
//        // reserved for future use
//        union switch (int v)
//        {
//        case 0:
//            void;
//        case 1:
//            LedgerEntryExtensionV1 v1;
//        }
//        ext;
//    };
//
type LedgerEntry struct {
	LastModifiedLedgerSeq Uint32
	Data                  LedgerEntryData
	Ext                   LedgerEntryExt
}

// EncodeTo encodes this value using the Encoder.
func (s *LedgerEntry) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.LastModifiedLedgerSeq.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Data.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*LedgerEntry)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *LedgerEntry) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.LastModifiedLedgerSeq.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.Data.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LedgerEntryData: %s", err)
	}
	nTmp, err = s.Ext.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LedgerEntryExt: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerEntry) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerEntry) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerEntry)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerEntry)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LedgerEntry) xdrType() {}

var _ xdrType = (*LedgerEntry)(nil)

// LedgerKeyAccount is an XDR NestedStruct defines as:
//
//   struct
//        {
//            AccountID accountID;
//        }
//
type LedgerKeyAccount struct {
	AccountId AccountId
}

// EncodeTo encodes this value using the Encoder.
func (s *LedgerKeyAccount) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.AccountId.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*LedgerKeyAccount)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *LedgerKeyAccount) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.AccountId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountId: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerKeyAccount) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerKeyAccount) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerKeyAccount)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerKeyAccount)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LedgerKeyAccount) xdrType() {}

var _ xdrType = (*LedgerKeyAccount)(nil)

// LedgerKeyTrustLine is an XDR NestedStruct defines as:
//
//   struct
//        {
//            AccountID accountID;
//            TrustLineAsset asset;
//        }
//
type LedgerKeyTrustLine struct {
	AccountId AccountId
	Asset     TrustLineAsset
}

// EncodeTo encodes this value using the Encoder.
func (s *LedgerKeyTrustLine) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.AccountId.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Asset.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*LedgerKeyTrustLine)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *LedgerKeyTrustLine) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.AccountId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountId: %s", err)
	}
	nTmp, err = s.Asset.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TrustLineAsset: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerKeyTrustLine) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerKeyTrustLine) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerKeyTrustLine)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerKeyTrustLine)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LedgerKeyTrustLine) xdrType() {}

var _ xdrType = (*LedgerKeyTrustLine)(nil)

// LedgerKeyOffer is an XDR NestedStruct defines as:
//
//   struct
//        {
//            AccountID sellerID;
//            int64 offerID;
//        }
//
type LedgerKeyOffer struct {
	SellerId AccountId
	OfferId  Int64
}

// EncodeTo encodes this value using the Encoder.
func (s *LedgerKeyOffer) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.SellerId.EncodeTo(e); err != nil {
		return err
	}
	if err = s.OfferId.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*LedgerKeyOffer)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *LedgerKeyOffer) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.SellerId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountId: %s", err)
	}
	nTmp, err = s.OfferId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerKeyOffer) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerKeyOffer) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerKeyOffer)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerKeyOffer)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LedgerKeyOffer) xdrType() {}

var _ xdrType = (*LedgerKeyOffer)(nil)

// LedgerKeyData is an XDR NestedStruct defines as:
//
//   struct
//        {
//            AccountID accountID;
//            string64 dataName;
//        }
//
type LedgerKeyData struct {
	AccountId AccountId
	DataName  String64
}

// EncodeTo encodes this value using the Encoder.
func (s *LedgerKeyData) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.AccountId.EncodeTo(e); err != nil {
		return err
	}
	if err = s.DataName.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*LedgerKeyData)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *LedgerKeyData) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.AccountId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountId: %s", err)
	}
	nTmp, err = s.DataName.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding String64: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerKeyData) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerKeyData) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerKeyData)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerKeyData)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LedgerKeyData) xdrType() {}

var _ xdrType = (*LedgerKeyData)(nil)

// LedgerKeyClaimableBalance is an XDR NestedStruct defines as:
//
//   struct
//        {
//            ClaimableBalanceID balanceID;
//        }
//
type LedgerKeyClaimableBalance struct {
	BalanceId ClaimableBalanceId
}

// EncodeTo encodes this value using the Encoder.
func (s *LedgerKeyClaimableBalance) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.BalanceId.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*LedgerKeyClaimableBalance)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *LedgerKeyClaimableBalance) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.BalanceId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ClaimableBalanceId: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerKeyClaimableBalance) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerKeyClaimableBalance) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerKeyClaimableBalance)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerKeyClaimableBalance)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LedgerKeyClaimableBalance) xdrType() {}

var _ xdrType = (*LedgerKeyClaimableBalance)(nil)

// LedgerKeyLiquidityPool is an XDR NestedStruct defines as:
//
//   struct
//        {
//            PoolID liquidityPoolID;
//        }
//
type LedgerKeyLiquidityPool struct {
	LiquidityPoolId PoolId
}

// EncodeTo encodes this value using the Encoder.
func (s *LedgerKeyLiquidityPool) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.LiquidityPoolId.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*LedgerKeyLiquidityPool)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *LedgerKeyLiquidityPool) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.LiquidityPoolId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding PoolId: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerKeyLiquidityPool) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerKeyLiquidityPool) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerKeyLiquidityPool)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerKeyLiquidityPool)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LedgerKeyLiquidityPool) xdrType() {}

var _ xdrType = (*LedgerKeyLiquidityPool)(nil)

// LedgerKey is an XDR Union defines as:
//
//   union LedgerKey switch (LedgerEntryType type)
//    {
//    case ACCOUNT:
//        struct
//        {
//            AccountID accountID;
//        } account;
//
//    case TRUSTLINE:
//        struct
//        {
//            AccountID accountID;
//            TrustLineAsset asset;
//        } trustLine;
//
//    case OFFER:
//        struct
//        {
//            AccountID sellerID;
//            int64 offerID;
//        } offer;
//
//    case DATA:
//        struct
//        {
//            AccountID accountID;
//            string64 dataName;
//        } data;
//
//    case CLAIMABLE_BALANCE:
//        struct
//        {
//            ClaimableBalanceID balanceID;
//        } claimableBalance;
//
//    case LIQUIDITY_POOL:
//        struct
//        {
//            PoolID liquidityPoolID;
//        } liquidityPool;
//    };
//
type LedgerKey struct {
	Type             LedgerEntryType
	Account          *LedgerKeyAccount
	TrustLine        *LedgerKeyTrustLine
	Offer            *LedgerKeyOffer
	Data             *LedgerKeyData
	ClaimableBalance *LedgerKeyClaimableBalance
	LiquidityPool    *LedgerKeyLiquidityPool
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKey) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKey
func (u LedgerKey) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerEntryType(sw) {
	case LedgerEntryTypeAccount:
		return "Account", true
	case LedgerEntryTypeTrustline:
		return "TrustLine", true
	case LedgerEntryTypeOffer:
		return "Offer", true
	case LedgerEntryTypeData:
		return "Data", true
	case LedgerEntryTypeClaimableBalance:
		return "ClaimableBalance", true
	case LedgerEntryTypeLiquidityPool:
		return "LiquidityPool", true
	}
	return "-", false
}

// NewLedgerKey creates a new  LedgerKey.
func NewLedgerKey(aType LedgerEntryType, value interface{}) (result LedgerKey, err error) {
	result.Type = aType
	switch LedgerEntryType(aType) {
	case LedgerEntryTypeAccount:
		tv, ok := value.(LedgerKeyAccount)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyAccount")
			return
		}
		result.Account = &tv
	case LedgerEntryTypeTrustline:
		tv, ok := value.(LedgerKeyTrustLine)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyTrustLine")
			return
		}
		result.TrustLine = &tv
	case LedgerEntryTypeOffer:
		tv, ok := value.(LedgerKeyOffer)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyOffer")
			return
		}
		result.Offer = &tv
	case LedgerEntryTypeData:
		tv, ok := value.(LedgerKeyData)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyData")
			return
		}
		result.Data = &tv
	case LedgerEntryTypeClaimableBalance:
		tv, ok := value.(LedgerKeyClaimableBalance)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyClaimableBalance")
			return
		}
		result.ClaimableBalance = &tv
	case LedgerEntryTypeLiquidityPool:
		tv, ok := value.(LedgerKeyLiquidityPool)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyLiquidityPool")
			return
		}
		result.LiquidityPool = &tv
	}
	return
}

// MustAccount retrieves the Account value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustAccount() LedgerKeyAccount {
	val, ok := u.GetAccount()

	if !ok {
		panic("arm Account is not set")
	}

	return val
}

// GetAccount retrieves the Account value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetAccount() (result LedgerKeyAccount, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Account" {
		result = *u.Account
		ok = true
	}

	return
}

// MustTrustLine retrieves the TrustLine value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustTrustLine() LedgerKeyTrustLine {
	val, ok := u.GetTrustLine()

	if !ok {
		panic("arm TrustLine is not set")
	}

	return val
}

// GetTrustLine retrieves the TrustLine value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetTrustLine() (result LedgerKeyTrustLine, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "TrustLine" {
		result = *u.TrustLine
		ok = true
	}

	return
}

// MustOffer retrieves the Offer value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustOffer() LedgerKeyOffer {
	val, ok := u.GetOffer()

	if !ok {
		panic("arm Offer is not set")
	}

	return val
}

// GetOffer retrieves the Offer value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetOffer() (result LedgerKeyOffer, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Offer" {
		result = *u.Offer
		ok = true
	}

	return
}

// MustData retrieves the Data value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustData() LedgerKeyData {
	val, ok := u.GetData()

	if !ok {
		panic("arm Data is not set")
	}

	return val
}

// GetData retrieves the Data value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetData() (result LedgerKeyData, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Data" {
		result = *u.Data
		ok = true
	}

	return
}

// MustClaimableBalance retrieves the ClaimableBalance value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustClaimableBalance() LedgerKeyClaimableBalance {
	val, ok := u.GetClaimableBalance()

	if !ok {
		panic("arm ClaimableBalance is not set")
	}

	return val
}

// GetClaimableBalance retrieves the ClaimableBalance value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetClaimableBalance() (result LedgerKeyClaimableBalance, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ClaimableBalance" {
		result = *u.ClaimableBalance
		ok = true
	}

	return
}

// MustLiquidityPool retrieves the LiquidityPool value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustLiquidityPool() LedgerKeyLiquidityPool {
	val, ok := u.GetLiquidityPool()

	if !ok {
		panic("arm LiquidityPool is not set")
	}

	return val
}

// GetLiquidityPool retrieves the LiquidityPool value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetLiquidityPool() (result LedgerKeyLiquidityPool, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "LiquidityPool" {
		result = *u.LiquidityPool
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u LedgerKey) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch LedgerEntryType(u.Type) {
	case LedgerEntryTypeAccount:
		if err = (*u.Account).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case LedgerEntryTypeTrustline:
		if err = (*u.TrustLine).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case LedgerEntryTypeOffer:
		if err = (*u.Offer).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case LedgerEntryTypeData:
		if err = (*u.Data).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case LedgerEntryTypeClaimableBalance:
		if err = (*u.ClaimableBalance).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case LedgerEntryTypeLiquidityPool:
		if err = (*u.LiquidityPool).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (LedgerEntryType) switch value '%d' is not valid for union LedgerKey", u.Type)
}

var _ decoderFrom = (*LedgerKey)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *LedgerKey) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LedgerEntryType: %s", err)
	}
	switch LedgerEntryType(u.Type) {
	case LedgerEntryTypeAccount:
		u.Account = new(LedgerKeyAccount)
		nTmp, err = (*u.Account).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding LedgerKeyAccount: %s", err)
		}
		return n, nil
	case LedgerEntryTypeTrustline:
		u.TrustLine = new(LedgerKeyTrustLine)
		nTmp, err = (*u.TrustLine).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding LedgerKeyTrustLine: %s", err)
		}
		return n, nil
	case LedgerEntryTypeOffer:
		u.Offer = new(LedgerKeyOffer)
		nTmp, err = (*u.Offer).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding LedgerKeyOffer: %s", err)
		}
		return n, nil
	case LedgerEntryTypeData:
		u.Data = new(LedgerKeyData)
		nTmp, err = (*u.Data).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding LedgerKeyData: %s", err)
		}
		return n, nil
	case LedgerEntryTypeClaimableBalance:
		u.ClaimableBalance = new(LedgerKeyClaimableBalance)
		nTmp, err = (*u.ClaimableBalance).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding LedgerKeyClaimableBalance: %s", err)
		}
		return n, nil
	case LedgerEntryTypeLiquidityPool:
		u.LiquidityPool = new(LedgerKeyLiquidityPool)
		nTmp, err = (*u.LiquidityPool).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding LedgerKeyLiquidityPool: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union LedgerKey has invalid Type (LedgerEntryType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerKey) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerKey) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerKey)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerKey)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LedgerKey) xdrType() {}

var _ xdrType = (*LedgerKey)(nil)

// EnvelopeType is an XDR Enum defines as:
//
//   enum EnvelopeType
//    {
//        ENVELOPE_TYPE_TX_V0 = 0,
//        ENVELOPE_TYPE_SCP = 1,
//        ENVELOPE_TYPE_TX = 2,
//        ENVELOPE_TYPE_AUTH = 3,
//        ENVELOPE_TYPE_SCPVALUE = 4,
//        ENVELOPE_TYPE_TX_FEE_BUMP = 5,
//        ENVELOPE_TYPE_OP_ID = 6,
//        ENVELOPE_TYPE_POOL_REVOKE_OP_ID = 7
//    };
//
type EnvelopeType int32

const (
	EnvelopeTypeEnvelopeTypeTxV0           EnvelopeType = 0
	EnvelopeTypeEnvelopeTypeScp            EnvelopeType = 1
	EnvelopeTypeEnvelopeTypeTx             EnvelopeType = 2
	EnvelopeTypeEnvelopeTypeAuth           EnvelopeType = 3
	EnvelopeTypeEnvelopeTypeScpvalue       EnvelopeType = 4
	EnvelopeTypeEnvelopeTypeTxFeeBump      EnvelopeType = 5
	EnvelopeTypeEnvelopeTypeOpId           EnvelopeType = 6
	EnvelopeTypeEnvelopeTypePoolRevokeOpId EnvelopeType = 7
)

var envelopeTypeMap = map[int32]string{
	0: "EnvelopeTypeEnvelopeTypeTxV0",
	1: "EnvelopeTypeEnvelopeTypeScp",
	2: "EnvelopeTypeEnvelopeTypeTx",
	3: "EnvelopeTypeEnvelopeTypeAuth",
	4: "EnvelopeTypeEnvelopeTypeScpvalue",
	5: "EnvelopeTypeEnvelopeTypeTxFeeBump",
	6: "EnvelopeTypeEnvelopeTypeOpId",
	7: "EnvelopeTypeEnvelopeTypePoolRevokeOpId",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for EnvelopeType
func (e EnvelopeType) ValidEnum(v int32) bool {
	_, ok := envelopeTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e EnvelopeType) String() string {
	name, _ := envelopeTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e EnvelopeType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := envelopeTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid EnvelopeType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*EnvelopeType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *EnvelopeType) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding EnvelopeType: %s", err)
	}
	if _, ok := envelopeTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid EnvelopeType enum value", v)
	}
	*e = EnvelopeType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s EnvelopeType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *EnvelopeType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*EnvelopeType)(nil)
	_ encoding.BinaryUnmarshaler = (*EnvelopeType)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s EnvelopeType) xdrType() {}

var _ xdrType = (*EnvelopeType)(nil)

// UpgradeType is an XDR Typedef defines as:
//
//   typedef opaque UpgradeType<128>;
//
type UpgradeType []byte

// XDRMaxSize implements the Sized interface for UpgradeType
func (e UpgradeType) XDRMaxSize() int {
	return 128
}

// EncodeTo encodes this value using the Encoder.
func (s UpgradeType) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeOpaque(s[:]); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*UpgradeType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *UpgradeType) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	(*s), nTmp, err = d.DecodeOpaque(128)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding UpgradeType: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s UpgradeType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *UpgradeType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*UpgradeType)(nil)
	_ encoding.BinaryUnmarshaler = (*UpgradeType)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s UpgradeType) xdrType() {}

var _ xdrType = (*UpgradeType)(nil)

// diamcircleValueType is an XDR Enum defines as:
//
//   enum diamcircleValueType
//    {
//        diamcircle_VALUE_BASIC = 0,
//        diamcircle_VALUE_SIGNED = 1
//    };
//
type diamcircleValueType int32

const (
	diamcircleValueTypediamcircleValueBasic  diamcircleValueType = 0
	diamcircleValueTypediamcircleValueSigned diamcircleValueType = 1
)

var diamcircleValueTypeMap = map[int32]string{
	0: "diamcircleValueTypediamcircleValueBasic",
	1: "diamcircleValueTypediamcircleValueSigned",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for diamcircleValueType
func (e diamcircleValueType) ValidEnum(v int32) bool {
	_, ok := diamcircleValueTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e diamcircleValueType) String() string {
	name, _ := diamcircleValueTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e diamcircleValueType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := diamcircleValueTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid diamcircleValueType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*diamcircleValueType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *diamcircleValueType) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding diamcircleValueType: %s", err)
	}
	if _, ok := diamcircleValueTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid diamcircleValueType enum value", v)
	}
	*e = diamcircleValueType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s diamcircleValueType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *diamcircleValueType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*diamcircleValueType)(nil)
	_ encoding.BinaryUnmarshaler = (*diamcircleValueType)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s diamcircleValueType) xdrType() {}

var _ xdrType = (*diamcircleValueType)(nil)

// LedgerCloseValueSignature is an XDR Struct defines as:
//
//   struct LedgerCloseValueSignature
//    {
//        NodeID nodeID;       // which node introduced the value
//        Signature signature; // nodeID's signature
//    };
//
type LedgerCloseValueSignature struct {
	NodeId    NodeId
	Signature Signature
}

// EncodeTo encodes this value using the Encoder.
func (s *LedgerCloseValueSignature) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.NodeId.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Signature.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*LedgerCloseValueSignature)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *LedgerCloseValueSignature) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.NodeId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding NodeId: %s", err)
	}
	nTmp, err = s.Signature.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Signature: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerCloseValueSignature) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerCloseValueSignature) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerCloseValueSignature)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerCloseValueSignature)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LedgerCloseValueSignature) xdrType() {}

var _ xdrType = (*LedgerCloseValueSignature)(nil)

// diamcircleValueExt is an XDR NestedUnion defines as:
//
//   union switch (diamcircleValueType v)
//        {
//        case diamcircle_VALUE_BASIC:
//            void;
//        case diamcircle_VALUE_SIGNED:
//            LedgerCloseValueSignature lcValueSignature;
//        }
//
type diamcircleValueExt struct {
	V                diamcircleValueType
	LcValueSignature *LedgerCloseValueSignature
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u diamcircleValueExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of diamcircleValueExt
func (u diamcircleValueExt) ArmForSwitch(sw int32) (string, bool) {
	switch diamcircleValueType(sw) {
	case diamcircleValueTypediamcircleValueBasic:
		return "", true
	case diamcircleValueTypediamcircleValueSigned:
		return "LcValueSignature", true
	}
	return "-", false
}

// NewdiamcircleValueExt creates a new  diamcircleValueExt.
func NewdiamcircleValueExt(v diamcircleValueType, value interface{}) (result diamcircleValueExt, err error) {
	result.V = v
	switch diamcircleValueType(v) {
	case diamcircleValueTypediamcircleValueBasic:
		// void
	case diamcircleValueTypediamcircleValueSigned:
		tv, ok := value.(LedgerCloseValueSignature)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerCloseValueSignature")
			return
		}
		result.LcValueSignature = &tv
	}
	return
}

// MustLcValueSignature retrieves the LcValueSignature value from the union,
// panicing if the value is not set.
func (u diamcircleValueExt) MustLcValueSignature() LedgerCloseValueSignature {
	val, ok := u.GetLcValueSignature()

	if !ok {
		panic("arm LcValueSignature is not set")
	}

	return val
}

// GetLcValueSignature retrieves the LcValueSignature value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u diamcircleValueExt) GetLcValueSignature() (result LedgerCloseValueSignature, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.V))

	if armName == "LcValueSignature" {
		result = *u.LcValueSignature
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u diamcircleValueExt) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.V.EncodeTo(e); err != nil {
		return err
	}
	switch diamcircleValueType(u.V) {
	case diamcircleValueTypediamcircleValueBasic:
		// Void
		return nil
	case diamcircleValueTypediamcircleValueSigned:
		if err = (*u.LcValueSignature).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("V (diamcircleValueType) switch value '%d' is not valid for union diamcircleValueExt", u.V)
}

var _ decoderFrom = (*diamcircleValueExt)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *diamcircleValueExt) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.V.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding diamcircleValueType: %s", err)
	}
	switch diamcircleValueType(u.V) {
	case diamcircleValueTypediamcircleValueBasic:
		// Void
		return n, nil
	case diamcircleValueTypediamcircleValueSigned:
		u.LcValueSignature = new(LedgerCloseValueSignature)
		nTmp, err = (*u.LcValueSignature).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding LedgerCloseValueSignature: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union diamcircleValueExt has invalid V (diamcircleValueType) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s diamcircleValueExt) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *diamcircleValueExt) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*diamcircleValueExt)(nil)
	_ encoding.BinaryUnmarshaler = (*diamcircleValueExt)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s diamcircleValueExt) xdrType() {}

var _ xdrType = (*diamcircleValueExt)(nil)

// diamcircleValue is an XDR Struct defines as:
//
//   struct diamcircleValue
//    {
//        Hash txSetHash;      // transaction set to apply to previous ledger
//        TimePoint closeTime; // network close time
//
//        // upgrades to apply to the previous ledger (usually empty)
//        // this is a vector of encoded 'LedgerUpgrade' so that nodes can drop
//        // unknown steps during consensus if needed.
//        // see notes below on 'LedgerUpgrade' for more detail
//        // max size is dictated by number of upgrade types (+ room for future)
//        UpgradeType upgrades<6>;
//
//        // reserved for future use
//        union switch (diamcircleValueType v)
//        {
//        case diamcircle_VALUE_BASIC:
//            void;
//        case diamcircle_VALUE_SIGNED:
//            LedgerCloseValueSignature lcValueSignature;
//        }
//        ext;
//    };
//
type diamcircleValue struct {
	TxSetHash Hash
	CloseTime TimePoint
	Upgrades  []UpgradeType `xdrmaxsize:"6"`
	Ext       diamcircleValueExt
}

// EncodeTo encodes this value using the Encoder.
func (s *diamcircleValue) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.TxSetHash.EncodeTo(e); err != nil {
		return err
	}
	if err = s.CloseTime.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeUint(uint32(len(s.Upgrades))); err != nil {
		return err
	}
	for i := 0; i < len(s.Upgrades); i++ {
		if err = s.Upgrades[i].EncodeTo(e); err != nil {
			return err
		}
	}
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*diamcircleValue)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *diamcircleValue) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.TxSetHash.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Hash: %s", err)
	}
	nTmp, err = s.CloseTime.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TimePoint: %s", err)
	}
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding UpgradeType: %s", err)
	}
	if l > 6 {
		return n, fmt.Errorf("decoding UpgradeType: data size (%d) exceeds size limit (6)", l)
	}
	s.Upgrades = nil
	if l > 0 {
		s.Upgrades = make([]UpgradeType, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.Upgrades[i].DecodeFrom(d)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding UpgradeType: %s", err)
			}
		}
	}
	nTmp, err = s.Ext.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding diamcircleValueExt: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s diamcircleValue) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *diamcircleValue) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*diamcircleValue)(nil)
	_ encoding.BinaryUnmarshaler = (*diamcircleValue)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s diamcircleValue) xdrType() {}

var _ xdrType = (*diamcircleValue)(nil)

// MaskLedgerHeaderFlags is an XDR Const defines as:
//
//   const MASK_LEDGER_HEADER_FLAGS = 0x7;
//
const MaskLedgerHeaderFlags = 0x7

// LedgerHeaderFlags is an XDR Enum defines as:
//
//   enum LedgerHeaderFlags
//    {
//        DISABLE_LIQUIDITY_POOL_TRADING_FLAG = 0x1,
//        DISABLE_LIQUIDITY_POOL_DEPOSIT_FLAG = 0x2,
//        DISABLE_LIQUIDITY_POOL_WITHDRAWAL_FLAG = 0x4
//    };
//
type LedgerHeaderFlags int32

const (
	LedgerHeaderFlagsDisableLiquidityPoolTradingFlag    LedgerHeaderFlags = 1
	LedgerHeaderFlagsDisableLiquidityPoolDepositFlag    LedgerHeaderFlags = 2
	LedgerHeaderFlagsDisableLiquidityPoolWithdrawalFlag LedgerHeaderFlags = 4
)

var ledgerHeaderFlagsMap = map[int32]string{
	1: "LedgerHeaderFlagsDisableLiquidityPoolTradingFlag",
	2: "LedgerHeaderFlagsDisableLiquidityPoolDepositFlag",
	4: "LedgerHeaderFlagsDisableLiquidityPoolWithdrawalFlag",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for LedgerHeaderFlags
func (e LedgerHeaderFlags) ValidEnum(v int32) bool {
	_, ok := ledgerHeaderFlagsMap[v]
	return ok
}

// String returns the name of `e`
func (e LedgerHeaderFlags) String() string {
	name, _ := ledgerHeaderFlagsMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e LedgerHeaderFlags) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := ledgerHeaderFlagsMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid LedgerHeaderFlags enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*LedgerHeaderFlags)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *LedgerHeaderFlags) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding LedgerHeaderFlags: %s", err)
	}
	if _, ok := ledgerHeaderFlagsMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid LedgerHeaderFlags enum value", v)
	}
	*e = LedgerHeaderFlags(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerHeaderFlags) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerHeaderFlags) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerHeaderFlags)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerHeaderFlags)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LedgerHeaderFlags) xdrType() {}

var _ xdrType = (*LedgerHeaderFlags)(nil)

// LedgerHeaderExtensionV1Ext is an XDR NestedUnion defines as:
//
//   union switch (int v)
//        {
//        case 0:
//            void;
//        }
//
type LedgerHeaderExtensionV1Ext struct {
	V int32
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerHeaderExtensionV1Ext) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerHeaderExtensionV1Ext
func (u LedgerHeaderExtensionV1Ext) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	}
	return "-", false
}

// NewLedgerHeaderExtensionV1Ext creates a new  LedgerHeaderExtensionV1Ext.
func NewLedgerHeaderExtensionV1Ext(v int32, value interface{}) (result LedgerHeaderExtensionV1Ext, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u LedgerHeaderExtensionV1Ext) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union LedgerHeaderExtensionV1Ext", u.V)
}

var _ decoderFrom = (*LedgerHeaderExtensionV1Ext)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *LedgerHeaderExtensionV1Ext) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %s", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	}
	return n, fmt.Errorf("union LedgerHeaderExtensionV1Ext has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerHeaderExtensionV1Ext) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerHeaderExtensionV1Ext) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerHeaderExtensionV1Ext)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerHeaderExtensionV1Ext)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LedgerHeaderExtensionV1Ext) xdrType() {}

var _ xdrType = (*LedgerHeaderExtensionV1Ext)(nil)

// LedgerHeaderExtensionV1 is an XDR Struct defines as:
//
//   struct LedgerHeaderExtensionV1
//    {
//        uint32 flags; // LedgerHeaderFlags
//
//        union switch (int v)
//        {
//        case 0:
//            void;
//        }
//        ext;
//    };
//
type LedgerHeaderExtensionV1 struct {
	Flags Uint32
	Ext   LedgerHeaderExtensionV1Ext
}

// EncodeTo encodes this value using the Encoder.
func (s *LedgerHeaderExtensionV1) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Flags.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*LedgerHeaderExtensionV1)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *LedgerHeaderExtensionV1) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Flags.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.Ext.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LedgerHeaderExtensionV1Ext: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerHeaderExtensionV1) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerHeaderExtensionV1) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerHeaderExtensionV1)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerHeaderExtensionV1)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LedgerHeaderExtensionV1) xdrType() {}

var _ xdrType = (*LedgerHeaderExtensionV1)(nil)

// LedgerHeaderExt is an XDR NestedUnion defines as:
//
//   union switch (int v)
//        {
//        case 0:
//            void;
//        case 1:
//            LedgerHeaderExtensionV1 v1;
//        }
//
type LedgerHeaderExt struct {
	V  int32
	V1 *LedgerHeaderExtensionV1
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerHeaderExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerHeaderExt
func (u LedgerHeaderExt) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	case 1:
		return "V1", true
	}
	return "-", false
}

// NewLedgerHeaderExt creates a new  LedgerHeaderExt.
func NewLedgerHeaderExt(v int32, value interface{}) (result LedgerHeaderExt, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	case 1:
		tv, ok := value.(LedgerHeaderExtensionV1)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerHeaderExtensionV1")
			return
		}
		result.V1 = &tv
	}
	return
}

// MustV1 retrieves the V1 value from the union,
// panicing if the value is not set.
func (u LedgerHeaderExt) MustV1() LedgerHeaderExtensionV1 {
	val, ok := u.GetV1()

	if !ok {
		panic("arm V1 is not set")
	}

	return val
}

// GetV1 retrieves the V1 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerHeaderExt) GetV1() (result LedgerHeaderExtensionV1, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.V))

	if armName == "V1" {
		result = *u.V1
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u LedgerHeaderExt) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	case 1:
		if err = (*u.V1).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union LedgerHeaderExt", u.V)
}

var _ decoderFrom = (*LedgerHeaderExt)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *LedgerHeaderExt) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %s", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	case 1:
		u.V1 = new(LedgerHeaderExtensionV1)
		nTmp, err = (*u.V1).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding LedgerHeaderExtensionV1: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union LedgerHeaderExt has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerHeaderExt) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerHeaderExt) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerHeaderExt)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerHeaderExt)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LedgerHeaderExt) xdrType() {}

var _ xdrType = (*LedgerHeaderExt)(nil)

// LedgerHeader is an XDR Struct defines as:
//
//   struct LedgerHeader
//    {
//        uint32 ledgerVersion;    // the protocol version of the ledger
//        Hash previousLedgerHash; // hash of the previous ledger header
//        diamcircleValue scpValue;   // what consensus agreed to
//        Hash txSetResultHash;    // the TransactionResultSet that led to this ledger
//        Hash bucketListHash;     // hash of the ledger state
//
//        uint32 ledgerSeq; // sequence number of this ledger
//
//        int64 totalCoins; // total number of stroops in existence.
//                          // 10,000,000 stroops in 1 XLM
//
//        int64 feePool;       // fees burned since last inflation run
//        uint32 inflationSeq; // inflation sequence number
//
//        uint64 idPool; // last used global ID, used for generating objects
//
//        uint32 baseFee;     // base fee per operation in stroops
//        uint32 baseReserve; // account base reserve in stroops
//
//        uint32 maxTxSetSize; // maximum size a transaction set can be
//
//        Hash skipList[4]; // hashes of ledgers in the past. allows you to jump back
//                          // in time without walking the chain back ledger by ledger
//                          // each slot contains the oldest ledger that is mod of
//                          // either 50  5000  50000 or 500000 depending on index
//                          // skipList[0] mod(50), skipList[1] mod(5000), etc
//
//        // reserved for future use
//        union switch (int v)
//        {
//        case 0:
//            void;
//        case 1:
//            LedgerHeaderExtensionV1 v1;
//        }
//        ext;
//    };
//
type LedgerHeader struct {
	LedgerVersion      Uint32
	PreviousLedgerHash Hash
	ScpValue           diamcircleValue
	TxSetResultHash    Hash
	BucketListHash     Hash
	LedgerSeq          Uint32
	TotalCoins         Int64
	FeePool            Int64
	InflationSeq       Uint32
	IdPool             Uint64
	BaseFee            Uint32
	BaseReserve        Uint32
	MaxTxSetSize       Uint32
	SkipList           [4]Hash
	Ext                LedgerHeaderExt
}

// EncodeTo encodes this value using the Encoder.
func (s *LedgerHeader) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.LedgerVersion.EncodeTo(e); err != nil {
		return err
	}
	if err = s.PreviousLedgerHash.EncodeTo(e); err != nil {
		return err
	}
	if err = s.ScpValue.EncodeTo(e); err != nil {
		return err
	}
	if err = s.TxSetResultHash.EncodeTo(e); err != nil {
		return err
	}
	if err = s.BucketListHash.EncodeTo(e); err != nil {
		return err
	}
	if err = s.LedgerSeq.EncodeTo(e); err != nil {
		return err
	}
	if err = s.TotalCoins.EncodeTo(e); err != nil {
		return err
	}
	if err = s.FeePool.EncodeTo(e); err != nil {
		return err
	}
	if err = s.InflationSeq.EncodeTo(e); err != nil {
		return err
	}
	if err = s.IdPool.EncodeTo(e); err != nil {
		return err
	}
	if err = s.BaseFee.EncodeTo(e); err != nil {
		return err
	}
	if err = s.BaseReserve.EncodeTo(e); err != nil {
		return err
	}
	if err = s.MaxTxSetSize.EncodeTo(e); err != nil {
		return err
	}
	for i := 0; i < len(s.SkipList); i++ {
		if err = s.SkipList[i].EncodeTo(e); err != nil {
			return err
		}
	}
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*LedgerHeader)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *LedgerHeader) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.LedgerVersion.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.PreviousLedgerHash.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Hash: %s", err)
	}
	nTmp, err = s.ScpValue.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding diamcircleValue: %s", err)
	}
	nTmp, err = s.TxSetResultHash.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Hash: %s", err)
	}
	nTmp, err = s.BucketListHash.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Hash: %s", err)
	}
	nTmp, err = s.LedgerSeq.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.TotalCoins.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	nTmp, err = s.FeePool.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	nTmp, err = s.InflationSeq.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.IdPool.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint64: %s", err)
	}
	nTmp, err = s.BaseFee.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.BaseReserve.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.MaxTxSetSize.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	for i := 0; i < len(s.SkipList); i++ {
		nTmp, err = s.SkipList[i].DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Hash: %s", err)
		}
	}
	nTmp, err = s.Ext.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LedgerHeaderExt: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerHeader) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerHeader) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerHeader)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerHeader)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LedgerHeader) xdrType() {}

var _ xdrType = (*LedgerHeader)(nil)

// LedgerUpgradeType is an XDR Enum defines as:
//
//   enum LedgerUpgradeType
//    {
//        LEDGER_UPGRADE_VERSION = 1,
//        LEDGER_UPGRADE_BASE_FEE = 2,
//        LEDGER_UPGRADE_MAX_TX_SET_SIZE = 3,
//        LEDGER_UPGRADE_BASE_RESERVE = 4,
//        LEDGER_UPGRADE_FLAGS = 5
//    };
//
type LedgerUpgradeType int32

const (
	LedgerUpgradeTypeLedgerUpgradeVersion      LedgerUpgradeType = 1
	LedgerUpgradeTypeLedgerUpgradeBaseFee      LedgerUpgradeType = 2
	LedgerUpgradeTypeLedgerUpgradeMaxTxSetSize LedgerUpgradeType = 3
	LedgerUpgradeTypeLedgerUpgradeBaseReserve  LedgerUpgradeType = 4
	LedgerUpgradeTypeLedgerUpgradeFlags        LedgerUpgradeType = 5
)

var ledgerUpgradeTypeMap = map[int32]string{
	1: "LedgerUpgradeTypeLedgerUpgradeVersion",
	2: "LedgerUpgradeTypeLedgerUpgradeBaseFee",
	3: "LedgerUpgradeTypeLedgerUpgradeMaxTxSetSize",
	4: "LedgerUpgradeTypeLedgerUpgradeBaseReserve",
	5: "LedgerUpgradeTypeLedgerUpgradeFlags",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for LedgerUpgradeType
func (e LedgerUpgradeType) ValidEnum(v int32) bool {
	_, ok := ledgerUpgradeTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e LedgerUpgradeType) String() string {
	name, _ := ledgerUpgradeTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e LedgerUpgradeType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := ledgerUpgradeTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid LedgerUpgradeType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*LedgerUpgradeType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *LedgerUpgradeType) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding LedgerUpgradeType: %s", err)
	}
	if _, ok := ledgerUpgradeTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid LedgerUpgradeType enum value", v)
	}
	*e = LedgerUpgradeType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerUpgradeType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerUpgradeType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerUpgradeType)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerUpgradeType)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LedgerUpgradeType) xdrType() {}

var _ xdrType = (*LedgerUpgradeType)(nil)

// LedgerUpgrade is an XDR Union defines as:
//
//   union LedgerUpgrade switch (LedgerUpgradeType type)
//    {
//    case LEDGER_UPGRADE_VERSION:
//        uint32 newLedgerVersion; // update ledgerVersion
//    case LEDGER_UPGRADE_BASE_FEE:
//        uint32 newBaseFee; // update baseFee
//    case LEDGER_UPGRADE_MAX_TX_SET_SIZE:
//        uint32 newMaxTxSetSize; // update maxTxSetSize
//    case LEDGER_UPGRADE_BASE_RESERVE:
//        uint32 newBaseReserve; // update baseReserve
//    case LEDGER_UPGRADE_FLAGS:
//        uint32 newFlags; // update flags
//    };
//
type LedgerUpgrade struct {
	Type             LedgerUpgradeType
	NewLedgerVersion *Uint32
	NewBaseFee       *Uint32
	NewMaxTxSetSize  *Uint32
	NewBaseReserve   *Uint32
	NewFlags         *Uint32
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerUpgrade) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerUpgrade
func (u LedgerUpgrade) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerUpgradeType(sw) {
	case LedgerUpgradeTypeLedgerUpgradeVersion:
		return "NewLedgerVersion", true
	case LedgerUpgradeTypeLedgerUpgradeBaseFee:
		return "NewBaseFee", true
	case LedgerUpgradeTypeLedgerUpgradeMaxTxSetSize:
		return "NewMaxTxSetSize", true
	case LedgerUpgradeTypeLedgerUpgradeBaseReserve:
		return "NewBaseReserve", true
	case LedgerUpgradeTypeLedgerUpgradeFlags:
		return "NewFlags", true
	}
	return "-", false
}

// NewLedgerUpgrade creates a new  LedgerUpgrade.
func NewLedgerUpgrade(aType LedgerUpgradeType, value interface{}) (result LedgerUpgrade, err error) {
	result.Type = aType
	switch LedgerUpgradeType(aType) {
	case LedgerUpgradeTypeLedgerUpgradeVersion:
		tv, ok := value.(Uint32)
		if !ok {
			err = fmt.Errorf("invalid value, must be Uint32")
			return
		}
		result.NewLedgerVersion = &tv
	case LedgerUpgradeTypeLedgerUpgradeBaseFee:
		tv, ok := value.(Uint32)
		if !ok {
			err = fmt.Errorf("invalid value, must be Uint32")
			return
		}
		result.NewBaseFee = &tv
	case LedgerUpgradeTypeLedgerUpgradeMaxTxSetSize:
		tv, ok := value.(Uint32)
		if !ok {
			err = fmt.Errorf("invalid value, must be Uint32")
			return
		}
		result.NewMaxTxSetSize = &tv
	case LedgerUpgradeTypeLedgerUpgradeBaseReserve:
		tv, ok := value.(Uint32)
		if !ok {
			err = fmt.Errorf("invalid value, must be Uint32")
			return
		}
		result.NewBaseReserve = &tv
	case LedgerUpgradeTypeLedgerUpgradeFlags:
		tv, ok := value.(Uint32)
		if !ok {
			err = fmt.Errorf("invalid value, must be Uint32")
			return
		}
		result.NewFlags = &tv
	}
	return
}

// MustNewLedgerVersion retrieves the NewLedgerVersion value from the union,
// panicing if the value is not set.
func (u LedgerUpgrade) MustNewLedgerVersion() Uint32 {
	val, ok := u.GetNewLedgerVersion()

	if !ok {
		panic("arm NewLedgerVersion is not set")
	}

	return val
}

// GetNewLedgerVersion retrieves the NewLedgerVersion value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerUpgrade) GetNewLedgerVersion() (result Uint32, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "NewLedgerVersion" {
		result = *u.NewLedgerVersion
		ok = true
	}

	return
}

// MustNewBaseFee retrieves the NewBaseFee value from the union,
// panicing if the value is not set.
func (u LedgerUpgrade) MustNewBaseFee() Uint32 {
	val, ok := u.GetNewBaseFee()

	if !ok {
		panic("arm NewBaseFee is not set")
	}

	return val
}

// GetNewBaseFee retrieves the NewBaseFee value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerUpgrade) GetNewBaseFee() (result Uint32, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "NewBaseFee" {
		result = *u.NewBaseFee
		ok = true
	}

	return
}

// MustNewMaxTxSetSize retrieves the NewMaxTxSetSize value from the union,
// panicing if the value is not set.
func (u LedgerUpgrade) MustNewMaxTxSetSize() Uint32 {
	val, ok := u.GetNewMaxTxSetSize()

	if !ok {
		panic("arm NewMaxTxSetSize is not set")
	}

	return val
}

// GetNewMaxTxSetSize retrieves the NewMaxTxSetSize value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerUpgrade) GetNewMaxTxSetSize() (result Uint32, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "NewMaxTxSetSize" {
		result = *u.NewMaxTxSetSize
		ok = true
	}

	return
}

// MustNewBaseReserve retrieves the NewBaseReserve value from the union,
// panicing if the value is not set.
func (u LedgerUpgrade) MustNewBaseReserve() Uint32 {
	val, ok := u.GetNewBaseReserve()

	if !ok {
		panic("arm NewBaseReserve is not set")
	}

	return val
}

// GetNewBaseReserve retrieves the NewBaseReserve value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerUpgrade) GetNewBaseReserve() (result Uint32, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "NewBaseReserve" {
		result = *u.NewBaseReserve
		ok = true
	}

	return
}

// MustNewFlags retrieves the NewFlags value from the union,
// panicing if the value is not set.
func (u LedgerUpgrade) MustNewFlags() Uint32 {
	val, ok := u.GetNewFlags()

	if !ok {
		panic("arm NewFlags is not set")
	}

	return val
}

// GetNewFlags retrieves the NewFlags value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerUpgrade) GetNewFlags() (result Uint32, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "NewFlags" {
		result = *u.NewFlags
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u LedgerUpgrade) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch LedgerUpgradeType(u.Type) {
	case LedgerUpgradeTypeLedgerUpgradeVersion:
		if err = (*u.NewLedgerVersion).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case LedgerUpgradeTypeLedgerUpgradeBaseFee:
		if err = (*u.NewBaseFee).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case LedgerUpgradeTypeLedgerUpgradeMaxTxSetSize:
		if err = (*u.NewMaxTxSetSize).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case LedgerUpgradeTypeLedgerUpgradeBaseReserve:
		if err = (*u.NewBaseReserve).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case LedgerUpgradeTypeLedgerUpgradeFlags:
		if err = (*u.NewFlags).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (LedgerUpgradeType) switch value '%d' is not valid for union LedgerUpgrade", u.Type)
}

var _ decoderFrom = (*LedgerUpgrade)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *LedgerUpgrade) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LedgerUpgradeType: %s", err)
	}
	switch LedgerUpgradeType(u.Type) {
	case LedgerUpgradeTypeLedgerUpgradeVersion:
		u.NewLedgerVersion = new(Uint32)
		nTmp, err = (*u.NewLedgerVersion).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Uint32: %s", err)
		}
		return n, nil
	case LedgerUpgradeTypeLedgerUpgradeBaseFee:
		u.NewBaseFee = new(Uint32)
		nTmp, err = (*u.NewBaseFee).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Uint32: %s", err)
		}
		return n, nil
	case LedgerUpgradeTypeLedgerUpgradeMaxTxSetSize:
		u.NewMaxTxSetSize = new(Uint32)
		nTmp, err = (*u.NewMaxTxSetSize).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Uint32: %s", err)
		}
		return n, nil
	case LedgerUpgradeTypeLedgerUpgradeBaseReserve:
		u.NewBaseReserve = new(Uint32)
		nTmp, err = (*u.NewBaseReserve).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Uint32: %s", err)
		}
		return n, nil
	case LedgerUpgradeTypeLedgerUpgradeFlags:
		u.NewFlags = new(Uint32)
		nTmp, err = (*u.NewFlags).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Uint32: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union LedgerUpgrade has invalid Type (LedgerUpgradeType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerUpgrade) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerUpgrade) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerUpgrade)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerUpgrade)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LedgerUpgrade) xdrType() {}

var _ xdrType = (*LedgerUpgrade)(nil)

// BucketEntryType is an XDR Enum defines as:
//
//   enum BucketEntryType
//    {
//        METAENTRY =
//            -1, // At-and-after protocol 11: bucket metadata, should come first.
//        LIVEENTRY = 0, // Before protocol 11: created-or-updated;
//                       // At-and-after protocol 11: only updated.
//        DEADENTRY = 1,
//        INITENTRY = 2 // At-and-after protocol 11: only created.
//    };
//
type BucketEntryType int32

const (
	BucketEntryTypeMetaentry BucketEntryType = -1
	BucketEntryTypeLiveentry BucketEntryType = 0
	BucketEntryTypeDeadentry BucketEntryType = 1
	BucketEntryTypeInitentry BucketEntryType = 2
)

var bucketEntryTypeMap = map[int32]string{
	-1: "BucketEntryTypeMetaentry",
	0:  "BucketEntryTypeLiveentry",
	1:  "BucketEntryTypeDeadentry",
	2:  "BucketEntryTypeInitentry",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for BucketEntryType
func (e BucketEntryType) ValidEnum(v int32) bool {
	_, ok := bucketEntryTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e BucketEntryType) String() string {
	name, _ := bucketEntryTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e BucketEntryType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := bucketEntryTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid BucketEntryType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*BucketEntryType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *BucketEntryType) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding BucketEntryType: %s", err)
	}
	if _, ok := bucketEntryTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid BucketEntryType enum value", v)
	}
	*e = BucketEntryType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s BucketEntryType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *BucketEntryType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*BucketEntryType)(nil)
	_ encoding.BinaryUnmarshaler = (*BucketEntryType)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s BucketEntryType) xdrType() {}

var _ xdrType = (*BucketEntryType)(nil)

// BucketMetadataExt is an XDR NestedUnion defines as:
//
//   union switch (int v)
//        {
//        case 0:
//            void;
//        }
//
type BucketMetadataExt struct {
	V int32
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u BucketMetadataExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of BucketMetadataExt
func (u BucketMetadataExt) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	}
	return "-", false
}

// NewBucketMetadataExt creates a new  BucketMetadataExt.
func NewBucketMetadataExt(v int32, value interface{}) (result BucketMetadataExt, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u BucketMetadataExt) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union BucketMetadataExt", u.V)
}

var _ decoderFrom = (*BucketMetadataExt)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *BucketMetadataExt) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %s", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	}
	return n, fmt.Errorf("union BucketMetadataExt has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s BucketMetadataExt) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *BucketMetadataExt) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*BucketMetadataExt)(nil)
	_ encoding.BinaryUnmarshaler = (*BucketMetadataExt)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s BucketMetadataExt) xdrType() {}

var _ xdrType = (*BucketMetadataExt)(nil)

// BucketMetadata is an XDR Struct defines as:
//
//   struct BucketMetadata
//    {
//        // Indicates the protocol version used to create / merge this bucket.
//        uint32 ledgerVersion;
//
//        // reserved for future use
//        union switch (int v)
//        {
//        case 0:
//            void;
//        }
//        ext;
//    };
//
type BucketMetadata struct {
	LedgerVersion Uint32
	Ext           BucketMetadataExt
}

// EncodeTo encodes this value using the Encoder.
func (s *BucketMetadata) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.LedgerVersion.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*BucketMetadata)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *BucketMetadata) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.LedgerVersion.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.Ext.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding BucketMetadataExt: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s BucketMetadata) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *BucketMetadata) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*BucketMetadata)(nil)
	_ encoding.BinaryUnmarshaler = (*BucketMetadata)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s BucketMetadata) xdrType() {}

var _ xdrType = (*BucketMetadata)(nil)

// BucketEntry is an XDR Union defines as:
//
//   union BucketEntry switch (BucketEntryType type)
//    {
//    case LIVEENTRY:
//    case INITENTRY:
//        LedgerEntry liveEntry;
//
//    case DEADENTRY:
//        LedgerKey deadEntry;
//    case METAENTRY:
//        BucketMetadata metaEntry;
//    };
//
type BucketEntry struct {
	Type      BucketEntryType
	LiveEntry *LedgerEntry
	DeadEntry *LedgerKey
	MetaEntry *BucketMetadata
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u BucketEntry) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of BucketEntry
func (u BucketEntry) ArmForSwitch(sw int32) (string, bool) {
	switch BucketEntryType(sw) {
	case BucketEntryTypeLiveentry:
		return "LiveEntry", true
	case BucketEntryTypeInitentry:
		return "LiveEntry", true
	case BucketEntryTypeDeadentry:
		return "DeadEntry", true
	case BucketEntryTypeMetaentry:
		return "MetaEntry", true
	}
	return "-", false
}

// NewBucketEntry creates a new  BucketEntry.
func NewBucketEntry(aType BucketEntryType, value interface{}) (result BucketEntry, err error) {
	result.Type = aType
	switch BucketEntryType(aType) {
	case BucketEntryTypeLiveentry:
		tv, ok := value.(LedgerEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerEntry")
			return
		}
		result.LiveEntry = &tv
	case BucketEntryTypeInitentry:
		tv, ok := value.(LedgerEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerEntry")
			return
		}
		result.LiveEntry = &tv
	case BucketEntryTypeDeadentry:
		tv, ok := value.(LedgerKey)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKey")
			return
		}
		result.DeadEntry = &tv
	case BucketEntryTypeMetaentry:
		tv, ok := value.(BucketMetadata)
		if !ok {
			err = fmt.Errorf("invalid value, must be BucketMetadata")
			return
		}
		result.MetaEntry = &tv
	}
	return
}

// MustLiveEntry retrieves the LiveEntry value from the union,
// panicing if the value is not set.
func (u BucketEntry) MustLiveEntry() LedgerEntry {
	val, ok := u.GetLiveEntry()

	if !ok {
		panic("arm LiveEntry is not set")
	}

	return val
}

// GetLiveEntry retrieves the LiveEntry value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u BucketEntry) GetLiveEntry() (result LedgerEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "LiveEntry" {
		result = *u.LiveEntry
		ok = true
	}

	return
}

// MustDeadEntry retrieves the DeadEntry value from the union,
// panicing if the value is not set.
func (u BucketEntry) MustDeadEntry() LedgerKey {
	val, ok := u.GetDeadEntry()

	if !ok {
		panic("arm DeadEntry is not set")
	}

	return val
}

// GetDeadEntry retrieves the DeadEntry value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u BucketEntry) GetDeadEntry() (result LedgerKey, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "DeadEntry" {
		result = *u.DeadEntry
		ok = true
	}

	return
}

// MustMetaEntry retrieves the MetaEntry value from the union,
// panicing if the value is not set.
func (u BucketEntry) MustMetaEntry() BucketMetadata {
	val, ok := u.GetMetaEntry()

	if !ok {
		panic("arm MetaEntry is not set")
	}

	return val
}

// GetMetaEntry retrieves the MetaEntry value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u BucketEntry) GetMetaEntry() (result BucketMetadata, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "MetaEntry" {
		result = *u.MetaEntry
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u BucketEntry) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch BucketEntryType(u.Type) {
	case BucketEntryTypeLiveentry:
		if err = (*u.LiveEntry).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case BucketEntryTypeInitentry:
		if err = (*u.LiveEntry).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case BucketEntryTypeDeadentry:
		if err = (*u.DeadEntry).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case BucketEntryTypeMetaentry:
		if err = (*u.MetaEntry).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (BucketEntryType) switch value '%d' is not valid for union BucketEntry", u.Type)
}

var _ decoderFrom = (*BucketEntry)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *BucketEntry) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding BucketEntryType: %s", err)
	}
	switch BucketEntryType(u.Type) {
	case BucketEntryTypeLiveentry:
		u.LiveEntry = new(LedgerEntry)
		nTmp, err = (*u.LiveEntry).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding LedgerEntry: %s", err)
		}
		return n, nil
	case BucketEntryTypeInitentry:
		u.LiveEntry = new(LedgerEntry)
		nTmp, err = (*u.LiveEntry).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding LedgerEntry: %s", err)
		}
		return n, nil
	case BucketEntryTypeDeadentry:
		u.DeadEntry = new(LedgerKey)
		nTmp, err = (*u.DeadEntry).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding LedgerKey: %s", err)
		}
		return n, nil
	case BucketEntryTypeMetaentry:
		u.MetaEntry = new(BucketMetadata)
		nTmp, err = (*u.MetaEntry).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding BucketMetadata: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union BucketEntry has invalid Type (BucketEntryType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s BucketEntry) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *BucketEntry) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*BucketEntry)(nil)
	_ encoding.BinaryUnmarshaler = (*BucketEntry)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s BucketEntry) xdrType() {}

var _ xdrType = (*BucketEntry)(nil)

// TransactionSet is an XDR Struct defines as:
//
//   struct TransactionSet
//    {
//        Hash previousLedgerHash;
//        TransactionEnvelope txs<>;
//    };
//
type TransactionSet struct {
	PreviousLedgerHash Hash
	Txs                []TransactionEnvelope
}

// EncodeTo encodes this value using the Encoder.
func (s *TransactionSet) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.PreviousLedgerHash.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeUint(uint32(len(s.Txs))); err != nil {
		return err
	}
	for i := 0; i < len(s.Txs); i++ {
		if err = s.Txs[i].EncodeTo(e); err != nil {
			return err
		}
	}
	return nil
}

var _ decoderFrom = (*TransactionSet)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *TransactionSet) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.PreviousLedgerHash.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Hash: %s", err)
	}
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TransactionEnvelope: %s", err)
	}
	s.Txs = nil
	if l > 0 {
		s.Txs = make([]TransactionEnvelope, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.Txs[i].DecodeFrom(d)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding TransactionEnvelope: %s", err)
			}
		}
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionSet) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionSet) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionSet)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionSet)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TransactionSet) xdrType() {}

var _ xdrType = (*TransactionSet)(nil)

// TransactionResultPair is an XDR Struct defines as:
//
//   struct TransactionResultPair
//    {
//        Hash transactionHash;
//        TransactionResult result; // result for the transaction
//    };
//
type TransactionResultPair struct {
	TransactionHash Hash
	Result          TransactionResult
}

// EncodeTo encodes this value using the Encoder.
func (s *TransactionResultPair) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.TransactionHash.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Result.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*TransactionResultPair)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *TransactionResultPair) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.TransactionHash.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Hash: %s", err)
	}
	nTmp, err = s.Result.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TransactionResult: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionResultPair) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionResultPair) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionResultPair)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionResultPair)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TransactionResultPair) xdrType() {}

var _ xdrType = (*TransactionResultPair)(nil)

// TransactionResultSet is an XDR Struct defines as:
//
//   struct TransactionResultSet
//    {
//        TransactionResultPair results<>;
//    };
//
type TransactionResultSet struct {
	Results []TransactionResultPair
}

// EncodeTo encodes this value using the Encoder.
func (s *TransactionResultSet) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeUint(uint32(len(s.Results))); err != nil {
		return err
	}
	for i := 0; i < len(s.Results); i++ {
		if err = s.Results[i].EncodeTo(e); err != nil {
			return err
		}
	}
	return nil
}

var _ decoderFrom = (*TransactionResultSet)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *TransactionResultSet) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TransactionResultPair: %s", err)
	}
	s.Results = nil
	if l > 0 {
		s.Results = make([]TransactionResultPair, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.Results[i].DecodeFrom(d)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding TransactionResultPair: %s", err)
			}
		}
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionResultSet) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionResultSet) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionResultSet)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionResultSet)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TransactionResultSet) xdrType() {}

var _ xdrType = (*TransactionResultSet)(nil)

// TransactionHistoryEntryExt is an XDR NestedUnion defines as:
//
//   union switch (int v)
//        {
//        case 0:
//            void;
//        }
//
type TransactionHistoryEntryExt struct {
	V int32
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TransactionHistoryEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TransactionHistoryEntryExt
func (u TransactionHistoryEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	}
	return "-", false
}

// NewTransactionHistoryEntryExt creates a new  TransactionHistoryEntryExt.
func NewTransactionHistoryEntryExt(v int32, value interface{}) (result TransactionHistoryEntryExt, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u TransactionHistoryEntryExt) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union TransactionHistoryEntryExt", u.V)
}

var _ decoderFrom = (*TransactionHistoryEntryExt)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *TransactionHistoryEntryExt) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %s", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	}
	return n, fmt.Errorf("union TransactionHistoryEntryExt has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionHistoryEntryExt) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionHistoryEntryExt) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionHistoryEntryExt)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionHistoryEntryExt)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TransactionHistoryEntryExt) xdrType() {}

var _ xdrType = (*TransactionHistoryEntryExt)(nil)

// TransactionHistoryEntry is an XDR Struct defines as:
//
//   struct TransactionHistoryEntry
//    {
//        uint32 ledgerSeq;
//        TransactionSet txSet;
//
//        // reserved for future use
//        union switch (int v)
//        {
//        case 0:
//            void;
//        }
//        ext;
//    };
//
type TransactionHistoryEntry struct {
	LedgerSeq Uint32
	TxSet     TransactionSet
	Ext       TransactionHistoryEntryExt
}

// EncodeTo encodes this value using the Encoder.
func (s *TransactionHistoryEntry) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.LedgerSeq.EncodeTo(e); err != nil {
		return err
	}
	if err = s.TxSet.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*TransactionHistoryEntry)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *TransactionHistoryEntry) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.LedgerSeq.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.TxSet.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TransactionSet: %s", err)
	}
	nTmp, err = s.Ext.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TransactionHistoryEntryExt: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionHistoryEntry) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionHistoryEntry) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionHistoryEntry)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionHistoryEntry)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TransactionHistoryEntry) xdrType() {}

var _ xdrType = (*TransactionHistoryEntry)(nil)

// TransactionHistoryResultEntryExt is an XDR NestedUnion defines as:
//
//   union switch (int v)
//        {
//        case 0:
//            void;
//        }
//
type TransactionHistoryResultEntryExt struct {
	V int32
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TransactionHistoryResultEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TransactionHistoryResultEntryExt
func (u TransactionHistoryResultEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	}
	return "-", false
}

// NewTransactionHistoryResultEntryExt creates a new  TransactionHistoryResultEntryExt.
func NewTransactionHistoryResultEntryExt(v int32, value interface{}) (result TransactionHistoryResultEntryExt, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u TransactionHistoryResultEntryExt) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union TransactionHistoryResultEntryExt", u.V)
}

var _ decoderFrom = (*TransactionHistoryResultEntryExt)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *TransactionHistoryResultEntryExt) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %s", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	}
	return n, fmt.Errorf("union TransactionHistoryResultEntryExt has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionHistoryResultEntryExt) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionHistoryResultEntryExt) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionHistoryResultEntryExt)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionHistoryResultEntryExt)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TransactionHistoryResultEntryExt) xdrType() {}

var _ xdrType = (*TransactionHistoryResultEntryExt)(nil)

// TransactionHistoryResultEntry is an XDR Struct defines as:
//
//   struct TransactionHistoryResultEntry
//    {
//        uint32 ledgerSeq;
//        TransactionResultSet txResultSet;
//
//        // reserved for future use
//        union switch (int v)
//        {
//        case 0:
//            void;
//        }
//        ext;
//    };
//
type TransactionHistoryResultEntry struct {
	LedgerSeq   Uint32
	TxResultSet TransactionResultSet
	Ext         TransactionHistoryResultEntryExt
}

// EncodeTo encodes this value using the Encoder.
func (s *TransactionHistoryResultEntry) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.LedgerSeq.EncodeTo(e); err != nil {
		return err
	}
	if err = s.TxResultSet.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*TransactionHistoryResultEntry)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *TransactionHistoryResultEntry) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.LedgerSeq.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.TxResultSet.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TransactionResultSet: %s", err)
	}
	nTmp, err = s.Ext.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TransactionHistoryResultEntryExt: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionHistoryResultEntry) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionHistoryResultEntry) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionHistoryResultEntry)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionHistoryResultEntry)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TransactionHistoryResultEntry) xdrType() {}

var _ xdrType = (*TransactionHistoryResultEntry)(nil)

// LedgerHeaderHistoryEntryExt is an XDR NestedUnion defines as:
//
//   union switch (int v)
//        {
//        case 0:
//            void;
//        }
//
type LedgerHeaderHistoryEntryExt struct {
	V int32
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerHeaderHistoryEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerHeaderHistoryEntryExt
func (u LedgerHeaderHistoryEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	}
	return "-", false
}

// NewLedgerHeaderHistoryEntryExt creates a new  LedgerHeaderHistoryEntryExt.
func NewLedgerHeaderHistoryEntryExt(v int32, value interface{}) (result LedgerHeaderHistoryEntryExt, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u LedgerHeaderHistoryEntryExt) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union LedgerHeaderHistoryEntryExt", u.V)
}

var _ decoderFrom = (*LedgerHeaderHistoryEntryExt)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *LedgerHeaderHistoryEntryExt) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %s", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	}
	return n, fmt.Errorf("union LedgerHeaderHistoryEntryExt has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerHeaderHistoryEntryExt) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerHeaderHistoryEntryExt) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerHeaderHistoryEntryExt)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerHeaderHistoryEntryExt)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LedgerHeaderHistoryEntryExt) xdrType() {}

var _ xdrType = (*LedgerHeaderHistoryEntryExt)(nil)

// LedgerHeaderHistoryEntry is an XDR Struct defines as:
//
//   struct LedgerHeaderHistoryEntry
//    {
//        Hash hash;
//        LedgerHeader header;
//
//        // reserved for future use
//        union switch (int v)
//        {
//        case 0:
//            void;
//        }
//        ext;
//    };
//
type LedgerHeaderHistoryEntry struct {
	Hash   Hash
	Header LedgerHeader
	Ext    LedgerHeaderHistoryEntryExt
}

// EncodeTo encodes this value using the Encoder.
func (s *LedgerHeaderHistoryEntry) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Hash.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Header.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*LedgerHeaderHistoryEntry)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *LedgerHeaderHistoryEntry) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Hash.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Hash: %s", err)
	}
	nTmp, err = s.Header.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LedgerHeader: %s", err)
	}
	nTmp, err = s.Ext.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LedgerHeaderHistoryEntryExt: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerHeaderHistoryEntry) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerHeaderHistoryEntry) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerHeaderHistoryEntry)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerHeaderHistoryEntry)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LedgerHeaderHistoryEntry) xdrType() {}

var _ xdrType = (*LedgerHeaderHistoryEntry)(nil)

// LedgerScpMessages is an XDR Struct defines as:
//
//   struct LedgerSCPMessages
//    {
//        uint32 ledgerSeq;
//        SCPEnvelope messages<>;
//    };
//
type LedgerScpMessages struct {
	LedgerSeq Uint32
	Messages  []ScpEnvelope
}

// EncodeTo encodes this value using the Encoder.
func (s *LedgerScpMessages) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.LedgerSeq.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeUint(uint32(len(s.Messages))); err != nil {
		return err
	}
	for i := 0; i < len(s.Messages); i++ {
		if err = s.Messages[i].EncodeTo(e); err != nil {
			return err
		}
	}
	return nil
}

var _ decoderFrom = (*LedgerScpMessages)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *LedgerScpMessages) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.LedgerSeq.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ScpEnvelope: %s", err)
	}
	s.Messages = nil
	if l > 0 {
		s.Messages = make([]ScpEnvelope, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.Messages[i].DecodeFrom(d)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding ScpEnvelope: %s", err)
			}
		}
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerScpMessages) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerScpMessages) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerScpMessages)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerScpMessages)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LedgerScpMessages) xdrType() {}

var _ xdrType = (*LedgerScpMessages)(nil)

// ScpHistoryEntryV0 is an XDR Struct defines as:
//
//   struct SCPHistoryEntryV0
//    {
//        SCPQuorumSet quorumSets<>; // additional quorum sets used by ledgerMessages
//        LedgerSCPMessages ledgerMessages;
//    };
//
type ScpHistoryEntryV0 struct {
	QuorumSets     []ScpQuorumSet
	LedgerMessages LedgerScpMessages
}

// EncodeTo encodes this value using the Encoder.
func (s *ScpHistoryEntryV0) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeUint(uint32(len(s.QuorumSets))); err != nil {
		return err
	}
	for i := 0; i < len(s.QuorumSets); i++ {
		if err = s.QuorumSets[i].EncodeTo(e); err != nil {
			return err
		}
	}
	if err = s.LedgerMessages.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*ScpHistoryEntryV0)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *ScpHistoryEntryV0) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ScpQuorumSet: %s", err)
	}
	s.QuorumSets = nil
	if l > 0 {
		s.QuorumSets = make([]ScpQuorumSet, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.QuorumSets[i].DecodeFrom(d)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding ScpQuorumSet: %s", err)
			}
		}
	}
	nTmp, err = s.LedgerMessages.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LedgerScpMessages: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ScpHistoryEntryV0) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ScpHistoryEntryV0) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ScpHistoryEntryV0)(nil)
	_ encoding.BinaryUnmarshaler = (*ScpHistoryEntryV0)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ScpHistoryEntryV0) xdrType() {}

var _ xdrType = (*ScpHistoryEntryV0)(nil)

// ScpHistoryEntry is an XDR Union defines as:
//
//   union SCPHistoryEntry switch (int v)
//    {
//    case 0:
//        SCPHistoryEntryV0 v0;
//    };
//
type ScpHistoryEntry struct {
	V  int32
	V0 *ScpHistoryEntryV0
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ScpHistoryEntry) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ScpHistoryEntry
func (u ScpHistoryEntry) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "V0", true
	}
	return "-", false
}

// NewScpHistoryEntry creates a new  ScpHistoryEntry.
func NewScpHistoryEntry(v int32, value interface{}) (result ScpHistoryEntry, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		tv, ok := value.(ScpHistoryEntryV0)
		if !ok {
			err = fmt.Errorf("invalid value, must be ScpHistoryEntryV0")
			return
		}
		result.V0 = &tv
	}
	return
}

// MustV0 retrieves the V0 value from the union,
// panicing if the value is not set.
func (u ScpHistoryEntry) MustV0() ScpHistoryEntryV0 {
	val, ok := u.GetV0()

	if !ok {
		panic("arm V0 is not set")
	}

	return val
}

// GetV0 retrieves the V0 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ScpHistoryEntry) GetV0() (result ScpHistoryEntryV0, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.V))

	if armName == "V0" {
		result = *u.V0
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u ScpHistoryEntry) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		if err = (*u.V0).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union ScpHistoryEntry", u.V)
}

var _ decoderFrom = (*ScpHistoryEntry)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *ScpHistoryEntry) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %s", err)
	}
	switch int32(u.V) {
	case 0:
		u.V0 = new(ScpHistoryEntryV0)
		nTmp, err = (*u.V0).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ScpHistoryEntryV0: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union ScpHistoryEntry has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ScpHistoryEntry) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ScpHistoryEntry) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ScpHistoryEntry)(nil)
	_ encoding.BinaryUnmarshaler = (*ScpHistoryEntry)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ScpHistoryEntry) xdrType() {}

var _ xdrType = (*ScpHistoryEntry)(nil)

// LedgerEntryChangeType is an XDR Enum defines as:
//
//   enum LedgerEntryChangeType
//    {
//        LEDGER_ENTRY_CREATED = 0, // entry was added to the ledger
//        LEDGER_ENTRY_UPDATED = 1, // entry was modified in the ledger
//        LEDGER_ENTRY_REMOVED = 2, // entry was removed from the ledger
//        LEDGER_ENTRY_STATE = 3    // value of the entry
//    };
//
type LedgerEntryChangeType int32

const (
	LedgerEntryChangeTypeLedgerEntryCreated LedgerEntryChangeType = 0
	LedgerEntryChangeTypeLedgerEntryUpdated LedgerEntryChangeType = 1
	LedgerEntryChangeTypeLedgerEntryRemoved LedgerEntryChangeType = 2
	LedgerEntryChangeTypeLedgerEntryState   LedgerEntryChangeType = 3
)

var ledgerEntryChangeTypeMap = map[int32]string{
	0: "LedgerEntryChangeTypeLedgerEntryCreated",
	1: "LedgerEntryChangeTypeLedgerEntryUpdated",
	2: "LedgerEntryChangeTypeLedgerEntryRemoved",
	3: "LedgerEntryChangeTypeLedgerEntryState",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for LedgerEntryChangeType
func (e LedgerEntryChangeType) ValidEnum(v int32) bool {
	_, ok := ledgerEntryChangeTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e LedgerEntryChangeType) String() string {
	name, _ := ledgerEntryChangeTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e LedgerEntryChangeType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := ledgerEntryChangeTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid LedgerEntryChangeType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*LedgerEntryChangeType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *LedgerEntryChangeType) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding LedgerEntryChangeType: %s", err)
	}
	if _, ok := ledgerEntryChangeTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid LedgerEntryChangeType enum value", v)
	}
	*e = LedgerEntryChangeType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerEntryChangeType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerEntryChangeType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerEntryChangeType)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerEntryChangeType)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LedgerEntryChangeType) xdrType() {}

var _ xdrType = (*LedgerEntryChangeType)(nil)

// LedgerEntryChange is an XDR Union defines as:
//
//   union LedgerEntryChange switch (LedgerEntryChangeType type)
//    {
//    case LEDGER_ENTRY_CREATED:
//        LedgerEntry created;
//    case LEDGER_ENTRY_UPDATED:
//        LedgerEntry updated;
//    case LEDGER_ENTRY_REMOVED:
//        LedgerKey removed;
//    case LEDGER_ENTRY_STATE:
//        LedgerEntry state;
//    };
//
type LedgerEntryChange struct {
	Type    LedgerEntryChangeType
	Created *LedgerEntry
	Updated *LedgerEntry
	Removed *LedgerKey
	State   *LedgerEntry
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerEntryChange) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerEntryChange
func (u LedgerEntryChange) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerEntryChangeType(sw) {
	case LedgerEntryChangeTypeLedgerEntryCreated:
		return "Created", true
	case LedgerEntryChangeTypeLedgerEntryUpdated:
		return "Updated", true
	case LedgerEntryChangeTypeLedgerEntryRemoved:
		return "Removed", true
	case LedgerEntryChangeTypeLedgerEntryState:
		return "State", true
	}
	return "-", false
}

// NewLedgerEntryChange creates a new  LedgerEntryChange.
func NewLedgerEntryChange(aType LedgerEntryChangeType, value interface{}) (result LedgerEntryChange, err error) {
	result.Type = aType
	switch LedgerEntryChangeType(aType) {
	case LedgerEntryChangeTypeLedgerEntryCreated:
		tv, ok := value.(LedgerEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerEntry")
			return
		}
		result.Created = &tv
	case LedgerEntryChangeTypeLedgerEntryUpdated:
		tv, ok := value.(LedgerEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerEntry")
			return
		}
		result.Updated = &tv
	case LedgerEntryChangeTypeLedgerEntryRemoved:
		tv, ok := value.(LedgerKey)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKey")
			return
		}
		result.Removed = &tv
	case LedgerEntryChangeTypeLedgerEntryState:
		tv, ok := value.(LedgerEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerEntry")
			return
		}
		result.State = &tv
	}
	return
}

// MustCreated retrieves the Created value from the union,
// panicing if the value is not set.
func (u LedgerEntryChange) MustCreated() LedgerEntry {
	val, ok := u.GetCreated()

	if !ok {
		panic("arm Created is not set")
	}

	return val
}

// GetCreated retrieves the Created value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryChange) GetCreated() (result LedgerEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Created" {
		result = *u.Created
		ok = true
	}

	return
}

// MustUpdated retrieves the Updated value from the union,
// panicing if the value is not set.
func (u LedgerEntryChange) MustUpdated() LedgerEntry {
	val, ok := u.GetUpdated()

	if !ok {
		panic("arm Updated is not set")
	}

	return val
}

// GetUpdated retrieves the Updated value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryChange) GetUpdated() (result LedgerEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Updated" {
		result = *u.Updated
		ok = true
	}

	return
}

// MustRemoved retrieves the Removed value from the union,
// panicing if the value is not set.
func (u LedgerEntryChange) MustRemoved() LedgerKey {
	val, ok := u.GetRemoved()

	if !ok {
		panic("arm Removed is not set")
	}

	return val
}

// GetRemoved retrieves the Removed value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryChange) GetRemoved() (result LedgerKey, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Removed" {
		result = *u.Removed
		ok = true
	}

	return
}

// MustState retrieves the State value from the union,
// panicing if the value is not set.
func (u LedgerEntryChange) MustState() LedgerEntry {
	val, ok := u.GetState()

	if !ok {
		panic("arm State is not set")
	}

	return val
}

// GetState retrieves the State value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryChange) GetState() (result LedgerEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "State" {
		result = *u.State
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u LedgerEntryChange) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch LedgerEntryChangeType(u.Type) {
	case LedgerEntryChangeTypeLedgerEntryCreated:
		if err = (*u.Created).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case LedgerEntryChangeTypeLedgerEntryUpdated:
		if err = (*u.Updated).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case LedgerEntryChangeTypeLedgerEntryRemoved:
		if err = (*u.Removed).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case LedgerEntryChangeTypeLedgerEntryState:
		if err = (*u.State).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (LedgerEntryChangeType) switch value '%d' is not valid for union LedgerEntryChange", u.Type)
}

var _ decoderFrom = (*LedgerEntryChange)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *LedgerEntryChange) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LedgerEntryChangeType: %s", err)
	}
	switch LedgerEntryChangeType(u.Type) {
	case LedgerEntryChangeTypeLedgerEntryCreated:
		u.Created = new(LedgerEntry)
		nTmp, err = (*u.Created).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding LedgerEntry: %s", err)
		}
		return n, nil
	case LedgerEntryChangeTypeLedgerEntryUpdated:
		u.Updated = new(LedgerEntry)
		nTmp, err = (*u.Updated).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding LedgerEntry: %s", err)
		}
		return n, nil
	case LedgerEntryChangeTypeLedgerEntryRemoved:
		u.Removed = new(LedgerKey)
		nTmp, err = (*u.Removed).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding LedgerKey: %s", err)
		}
		return n, nil
	case LedgerEntryChangeTypeLedgerEntryState:
		u.State = new(LedgerEntry)
		nTmp, err = (*u.State).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding LedgerEntry: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union LedgerEntryChange has invalid Type (LedgerEntryChangeType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerEntryChange) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerEntryChange) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerEntryChange)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerEntryChange)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LedgerEntryChange) xdrType() {}

var _ xdrType = (*LedgerEntryChange)(nil)

// LedgerEntryChanges is an XDR Typedef defines as:
//
//   typedef LedgerEntryChange LedgerEntryChanges<>;
//
type LedgerEntryChanges []LedgerEntryChange

// EncodeTo encodes this value using the Encoder.
func (s LedgerEntryChanges) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeUint(uint32(len(s))); err != nil {
		return err
	}
	for i := 0; i < len(s); i++ {
		if err = s[i].EncodeTo(e); err != nil {
			return err
		}
	}
	return nil
}

var _ decoderFrom = (*LedgerEntryChanges)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *LedgerEntryChanges) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LedgerEntryChange: %s", err)
	}
	(*s) = nil
	if l > 0 {
		(*s) = make([]LedgerEntryChange, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = (*s)[i].DecodeFrom(d)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding LedgerEntryChange: %s", err)
			}
		}
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerEntryChanges) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerEntryChanges) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerEntryChanges)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerEntryChanges)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LedgerEntryChanges) xdrType() {}

var _ xdrType = (*LedgerEntryChanges)(nil)

// OperationMeta is an XDR Struct defines as:
//
//   struct OperationMeta
//    {
//        LedgerEntryChanges changes;
//    };
//
type OperationMeta struct {
	Changes LedgerEntryChanges
}

// EncodeTo encodes this value using the Encoder.
func (s *OperationMeta) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Changes.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*OperationMeta)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *OperationMeta) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Changes.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LedgerEntryChanges: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s OperationMeta) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *OperationMeta) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*OperationMeta)(nil)
	_ encoding.BinaryUnmarshaler = (*OperationMeta)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s OperationMeta) xdrType() {}

var _ xdrType = (*OperationMeta)(nil)

// TransactionMetaV1 is an XDR Struct defines as:
//
//   struct TransactionMetaV1
//    {
//        LedgerEntryChanges txChanges; // tx level changes if any
//        OperationMeta operations<>;   // meta for each operation
//    };
//
type TransactionMetaV1 struct {
	TxChanges  LedgerEntryChanges
	Operations []OperationMeta
}

// EncodeTo encodes this value using the Encoder.
func (s *TransactionMetaV1) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.TxChanges.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeUint(uint32(len(s.Operations))); err != nil {
		return err
	}
	for i := 0; i < len(s.Operations); i++ {
		if err = s.Operations[i].EncodeTo(e); err != nil {
			return err
		}
	}
	return nil
}

var _ decoderFrom = (*TransactionMetaV1)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *TransactionMetaV1) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.TxChanges.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LedgerEntryChanges: %s", err)
	}
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding OperationMeta: %s", err)
	}
	s.Operations = nil
	if l > 0 {
		s.Operations = make([]OperationMeta, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.Operations[i].DecodeFrom(d)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding OperationMeta: %s", err)
			}
		}
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionMetaV1) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionMetaV1) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionMetaV1)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionMetaV1)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TransactionMetaV1) xdrType() {}

var _ xdrType = (*TransactionMetaV1)(nil)

// TransactionMetaV2 is an XDR Struct defines as:
//
//   struct TransactionMetaV2
//    {
//        LedgerEntryChanges txChangesBefore; // tx level changes before operations
//                                            // are applied if any
//        OperationMeta operations<>;         // meta for each operation
//        LedgerEntryChanges txChangesAfter;  // tx level changes after operations are
//                                            // applied if any
//    };
//
type TransactionMetaV2 struct {
	TxChangesBefore LedgerEntryChanges
	Operations      []OperationMeta
	TxChangesAfter  LedgerEntryChanges
}

// EncodeTo encodes this value using the Encoder.
func (s *TransactionMetaV2) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.TxChangesBefore.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeUint(uint32(len(s.Operations))); err != nil {
		return err
	}
	for i := 0; i < len(s.Operations); i++ {
		if err = s.Operations[i].EncodeTo(e); err != nil {
			return err
		}
	}
	if err = s.TxChangesAfter.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*TransactionMetaV2)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *TransactionMetaV2) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.TxChangesBefore.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LedgerEntryChanges: %s", err)
	}
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding OperationMeta: %s", err)
	}
	s.Operations = nil
	if l > 0 {
		s.Operations = make([]OperationMeta, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.Operations[i].DecodeFrom(d)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding OperationMeta: %s", err)
			}
		}
	}
	nTmp, err = s.TxChangesAfter.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LedgerEntryChanges: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionMetaV2) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionMetaV2) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionMetaV2)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionMetaV2)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TransactionMetaV2) xdrType() {}

var _ xdrType = (*TransactionMetaV2)(nil)

// TransactionMeta is an XDR Union defines as:
//
//   union TransactionMeta switch (int v)
//    {
//    case 0:
//        OperationMeta operations<>;
//    case 1:
//        TransactionMetaV1 v1;
//    case 2:
//        TransactionMetaV2 v2;
//    };
//
type TransactionMeta struct {
	V          int32
	Operations *[]OperationMeta
	V1         *TransactionMetaV1
	V2         *TransactionMetaV2
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TransactionMeta) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TransactionMeta
func (u TransactionMeta) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "Operations", true
	case 1:
		return "V1", true
	case 2:
		return "V2", true
	}
	return "-", false
}

// NewTransactionMeta creates a new  TransactionMeta.
func NewTransactionMeta(v int32, value interface{}) (result TransactionMeta, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		tv, ok := value.([]OperationMeta)
		if !ok {
			err = fmt.Errorf("invalid value, must be []OperationMeta")
			return
		}
		result.Operations = &tv
	case 1:
		tv, ok := value.(TransactionMetaV1)
		if !ok {
			err = fmt.Errorf("invalid value, must be TransactionMetaV1")
			return
		}
		result.V1 = &tv
	case 2:
		tv, ok := value.(TransactionMetaV2)
		if !ok {
			err = fmt.Errorf("invalid value, must be TransactionMetaV2")
			return
		}
		result.V2 = &tv
	}
	return
}

// MustOperations retrieves the Operations value from the union,
// panicing if the value is not set.
func (u TransactionMeta) MustOperations() []OperationMeta {
	val, ok := u.GetOperations()

	if !ok {
		panic("arm Operations is not set")
	}

	return val
}

// GetOperations retrieves the Operations value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u TransactionMeta) GetOperations() (result []OperationMeta, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.V))

	if armName == "Operations" {
		result = *u.Operations
		ok = true
	}

	return
}

// MustV1 retrieves the V1 value from the union,
// panicing if the value is not set.
func (u TransactionMeta) MustV1() TransactionMetaV1 {
	val, ok := u.GetV1()

	if !ok {
		panic("arm V1 is not set")
	}

	return val
}

// GetV1 retrieves the V1 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u TransactionMeta) GetV1() (result TransactionMetaV1, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.V))

	if armName == "V1" {
		result = *u.V1
		ok = true
	}

	return
}

// MustV2 retrieves the V2 value from the union,
// panicing if the value is not set.
func (u TransactionMeta) MustV2() TransactionMetaV2 {
	val, ok := u.GetV2()

	if !ok {
		panic("arm V2 is not set")
	}

	return val
}

// GetV2 retrieves the V2 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u TransactionMeta) GetV2() (result TransactionMetaV2, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.V))

	if armName == "V2" {
		result = *u.V2
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u TransactionMeta) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		if _, err = e.EncodeUint(uint32(len((*u.Operations)))); err != nil {
			return err
		}
		for i := 0; i < len((*u.Operations)); i++ {
			if err = (*u.Operations)[i].EncodeTo(e); err != nil {
				return err
			}
		}
		return nil
	case 1:
		if err = (*u.V1).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case 2:
		if err = (*u.V2).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union TransactionMeta", u.V)
}

var _ decoderFrom = (*TransactionMeta)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *TransactionMeta) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %s", err)
	}
	switch int32(u.V) {
	case 0:
		u.Operations = new([]OperationMeta)
		var l uint32
		l, nTmp, err = d.DecodeUint()
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding OperationMeta: %s", err)
		}
		(*u.Operations) = nil
		if l > 0 {
			(*u.Operations) = make([]OperationMeta, l)
			for i := uint32(0); i < l; i++ {
				nTmp, err = (*u.Operations)[i].DecodeFrom(d)
				n += nTmp
				if err != nil {
					return n, fmt.Errorf("decoding OperationMeta: %s", err)
				}
			}
		}
		return n, nil
	case 1:
		u.V1 = new(TransactionMetaV1)
		nTmp, err = (*u.V1).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding TransactionMetaV1: %s", err)
		}
		return n, nil
	case 2:
		u.V2 = new(TransactionMetaV2)
		nTmp, err = (*u.V2).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding TransactionMetaV2: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union TransactionMeta has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionMeta) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionMeta) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionMeta)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionMeta)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TransactionMeta) xdrType() {}

var _ xdrType = (*TransactionMeta)(nil)

// TransactionResultMeta is an XDR Struct defines as:
//
//   struct TransactionResultMeta
//    {
//        TransactionResultPair result;
//        LedgerEntryChanges feeProcessing;
//        TransactionMeta txApplyProcessing;
//    };
//
type TransactionResultMeta struct {
	Result            TransactionResultPair
	FeeProcessing     LedgerEntryChanges
	TxApplyProcessing TransactionMeta
}

// EncodeTo encodes this value using the Encoder.
func (s *TransactionResultMeta) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Result.EncodeTo(e); err != nil {
		return err
	}
	if err = s.FeeProcessing.EncodeTo(e); err != nil {
		return err
	}
	if err = s.TxApplyProcessing.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*TransactionResultMeta)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *TransactionResultMeta) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Result.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TransactionResultPair: %s", err)
	}
	nTmp, err = s.FeeProcessing.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LedgerEntryChanges: %s", err)
	}
	nTmp, err = s.TxApplyProcessing.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TransactionMeta: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionResultMeta) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionResultMeta) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionResultMeta)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionResultMeta)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TransactionResultMeta) xdrType() {}

var _ xdrType = (*TransactionResultMeta)(nil)

// UpgradeEntryMeta is an XDR Struct defines as:
//
//   struct UpgradeEntryMeta
//    {
//        LedgerUpgrade upgrade;
//        LedgerEntryChanges changes;
//    };
//
type UpgradeEntryMeta struct {
	Upgrade LedgerUpgrade
	Changes LedgerEntryChanges
}

// EncodeTo encodes this value using the Encoder.
func (s *UpgradeEntryMeta) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Upgrade.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Changes.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*UpgradeEntryMeta)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *UpgradeEntryMeta) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Upgrade.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LedgerUpgrade: %s", err)
	}
	nTmp, err = s.Changes.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LedgerEntryChanges: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s UpgradeEntryMeta) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *UpgradeEntryMeta) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*UpgradeEntryMeta)(nil)
	_ encoding.BinaryUnmarshaler = (*UpgradeEntryMeta)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s UpgradeEntryMeta) xdrType() {}

var _ xdrType = (*UpgradeEntryMeta)(nil)

// LedgerCloseMetaV0 is an XDR Struct defines as:
//
//   struct LedgerCloseMetaV0
//    {
//        LedgerHeaderHistoryEntry ledgerHeader;
//        // NB: txSet is sorted in "Hash order"
//        TransactionSet txSet;
//
//        // NB: transactions are sorted in apply order here
//        // fees for all transactions are processed first
//        // followed by applying transactions
//        TransactionResultMeta txProcessing<>;
//
//        // upgrades are applied last
//        UpgradeEntryMeta upgradesProcessing<>;
//
//        // other misc information attached to the ledger close
//        SCPHistoryEntry scpInfo<>;
//    };
//
type LedgerCloseMetaV0 struct {
	LedgerHeader       LedgerHeaderHistoryEntry
	TxSet              TransactionSet
	TxProcessing       []TransactionResultMeta
	UpgradesProcessing []UpgradeEntryMeta
	ScpInfo            []ScpHistoryEntry
}

// EncodeTo encodes this value using the Encoder.
func (s *LedgerCloseMetaV0) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.LedgerHeader.EncodeTo(e); err != nil {
		return err
	}
	if err = s.TxSet.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeUint(uint32(len(s.TxProcessing))); err != nil {
		return err
	}
	for i := 0; i < len(s.TxProcessing); i++ {
		if err = s.TxProcessing[i].EncodeTo(e); err != nil {
			return err
		}
	}
	if _, err = e.EncodeUint(uint32(len(s.UpgradesProcessing))); err != nil {
		return err
	}
	for i := 0; i < len(s.UpgradesProcessing); i++ {
		if err = s.UpgradesProcessing[i].EncodeTo(e); err != nil {
			return err
		}
	}
	if _, err = e.EncodeUint(uint32(len(s.ScpInfo))); err != nil {
		return err
	}
	for i := 0; i < len(s.ScpInfo); i++ {
		if err = s.ScpInfo[i].EncodeTo(e); err != nil {
			return err
		}
	}
	return nil
}

var _ decoderFrom = (*LedgerCloseMetaV0)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *LedgerCloseMetaV0) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.LedgerHeader.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LedgerHeaderHistoryEntry: %s", err)
	}
	nTmp, err = s.TxSet.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TransactionSet: %s", err)
	}
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TransactionResultMeta: %s", err)
	}
	s.TxProcessing = nil
	if l > 0 {
		s.TxProcessing = make([]TransactionResultMeta, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.TxProcessing[i].DecodeFrom(d)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding TransactionResultMeta: %s", err)
			}
		}
	}
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding UpgradeEntryMeta: %s", err)
	}
	s.UpgradesProcessing = nil
	if l > 0 {
		s.UpgradesProcessing = make([]UpgradeEntryMeta, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.UpgradesProcessing[i].DecodeFrom(d)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding UpgradeEntryMeta: %s", err)
			}
		}
	}
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ScpHistoryEntry: %s", err)
	}
	s.ScpInfo = nil
	if l > 0 {
		s.ScpInfo = make([]ScpHistoryEntry, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.ScpInfo[i].DecodeFrom(d)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding ScpHistoryEntry: %s", err)
			}
		}
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerCloseMetaV0) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerCloseMetaV0) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerCloseMetaV0)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerCloseMetaV0)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LedgerCloseMetaV0) xdrType() {}

var _ xdrType = (*LedgerCloseMetaV0)(nil)

// LedgerCloseMeta is an XDR Union defines as:
//
//   union LedgerCloseMeta switch (int v)
//    {
//    case 0:
//        LedgerCloseMetaV0 v0;
//    };
//
type LedgerCloseMeta struct {
	V  int32
	V0 *LedgerCloseMetaV0
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerCloseMeta) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerCloseMeta
func (u LedgerCloseMeta) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "V0", true
	}
	return "-", false
}

// NewLedgerCloseMeta creates a new  LedgerCloseMeta.
func NewLedgerCloseMeta(v int32, value interface{}) (result LedgerCloseMeta, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		tv, ok := value.(LedgerCloseMetaV0)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerCloseMetaV0")
			return
		}
		result.V0 = &tv
	}
	return
}

// MustV0 retrieves the V0 value from the union,
// panicing if the value is not set.
func (u LedgerCloseMeta) MustV0() LedgerCloseMetaV0 {
	val, ok := u.GetV0()

	if !ok {
		panic("arm V0 is not set")
	}

	return val
}

// GetV0 retrieves the V0 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerCloseMeta) GetV0() (result LedgerCloseMetaV0, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.V))

	if armName == "V0" {
		result = *u.V0
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u LedgerCloseMeta) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		if err = (*u.V0).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union LedgerCloseMeta", u.V)
}

var _ decoderFrom = (*LedgerCloseMeta)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *LedgerCloseMeta) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %s", err)
	}
	switch int32(u.V) {
	case 0:
		u.V0 = new(LedgerCloseMetaV0)
		nTmp, err = (*u.V0).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding LedgerCloseMetaV0: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union LedgerCloseMeta has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerCloseMeta) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerCloseMeta) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerCloseMeta)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerCloseMeta)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LedgerCloseMeta) xdrType() {}

var _ xdrType = (*LedgerCloseMeta)(nil)

// ErrorCode is an XDR Enum defines as:
//
//   enum ErrorCode
//    {
//        ERR_MISC = 0, // Unspecific error
//        ERR_DATA = 1, // Malformed data
//        ERR_CONF = 2, // Misconfiguration error
//        ERR_AUTH = 3, // Authentication failure
//        ERR_LOAD = 4  // System overloaded
//    };
//
type ErrorCode int32

const (
	ErrorCodeErrMisc ErrorCode = 0
	ErrorCodeErrData ErrorCode = 1
	ErrorCodeErrConf ErrorCode = 2
	ErrorCodeErrAuth ErrorCode = 3
	ErrorCodeErrLoad ErrorCode = 4
)

var errorCodeMap = map[int32]string{
	0: "ErrorCodeErrMisc",
	1: "ErrorCodeErrData",
	2: "ErrorCodeErrConf",
	3: "ErrorCodeErrAuth",
	4: "ErrorCodeErrLoad",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ErrorCode
func (e ErrorCode) ValidEnum(v int32) bool {
	_, ok := errorCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ErrorCode) String() string {
	name, _ := errorCodeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e ErrorCode) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := errorCodeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid ErrorCode enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*ErrorCode)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *ErrorCode) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding ErrorCode: %s", err)
	}
	if _, ok := errorCodeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid ErrorCode enum value", v)
	}
	*e = ErrorCode(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ErrorCode) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ErrorCode) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ErrorCode)(nil)
	_ encoding.BinaryUnmarshaler = (*ErrorCode)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ErrorCode) xdrType() {}

var _ xdrType = (*ErrorCode)(nil)

// Error is an XDR Struct defines as:
//
//   struct Error
//    {
//        ErrorCode code;
//        string msg<100>;
//    };
//
type Error struct {
	Code ErrorCode
	Msg  string `xdrmaxsize:"100"`
}

// EncodeTo encodes this value using the Encoder.
func (s *Error) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Code.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeString(string(s.Msg)); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Error)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Error) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Code.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ErrorCode: %s", err)
	}
	s.Msg, nTmp, err = d.DecodeString(100)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Msg: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Error) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Error) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Error)(nil)
	_ encoding.BinaryUnmarshaler = (*Error)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s Error) xdrType() {}

var _ xdrType = (*Error)(nil)

// AuthCert is an XDR Struct defines as:
//
//   struct AuthCert
//    {
//        Curve25519Public pubkey;
//        uint64 expiration;
//        Signature sig;
//    };
//
type AuthCert struct {
	Pubkey     Curve25519Public
	Expiration Uint64
	Sig        Signature
}

// EncodeTo encodes this value using the Encoder.
func (s *AuthCert) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Pubkey.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Expiration.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Sig.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*AuthCert)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *AuthCert) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Pubkey.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Curve25519Public: %s", err)
	}
	nTmp, err = s.Expiration.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint64: %s", err)
	}
	nTmp, err = s.Sig.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Signature: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s AuthCert) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AuthCert) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*AuthCert)(nil)
	_ encoding.BinaryUnmarshaler = (*AuthCert)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s AuthCert) xdrType() {}

var _ xdrType = (*AuthCert)(nil)

// Hello is an XDR Struct defines as:
//
//   struct Hello
//    {
//        uint32 ledgerVersion;
//        uint32 overlayVersion;
//        uint32 overlayMinVersion;
//        Hash networkID;
//        string versionStr<100>;
//        int listeningPort;
//        NodeID peerID;
//        AuthCert cert;
//        uint256 nonce;
//    };
//
type Hello struct {
	LedgerVersion     Uint32
	OverlayVersion    Uint32
	OverlayMinVersion Uint32
	NetworkId         Hash
	VersionStr        string `xdrmaxsize:"100"`
	ListeningPort     int32
	PeerId            NodeId
	Cert              AuthCert
	Nonce             Uint256
}

// EncodeTo encodes this value using the Encoder.
func (s *Hello) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.LedgerVersion.EncodeTo(e); err != nil {
		return err
	}
	if err = s.OverlayVersion.EncodeTo(e); err != nil {
		return err
	}
	if err = s.OverlayMinVersion.EncodeTo(e); err != nil {
		return err
	}
	if err = s.NetworkId.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeString(string(s.VersionStr)); err != nil {
		return err
	}
	if _, err = e.EncodeInt(int32(s.ListeningPort)); err != nil {
		return err
	}
	if err = s.PeerId.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Cert.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Nonce.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Hello)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Hello) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.LedgerVersion.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.OverlayVersion.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.OverlayMinVersion.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.NetworkId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Hash: %s", err)
	}
	s.VersionStr, nTmp, err = d.DecodeString(100)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding VersionStr: %s", err)
	}
	s.ListeningPort, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %s", err)
	}
	nTmp, err = s.PeerId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding NodeId: %s", err)
	}
	nTmp, err = s.Cert.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AuthCert: %s", err)
	}
	nTmp, err = s.Nonce.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint256: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Hello) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Hello) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Hello)(nil)
	_ encoding.BinaryUnmarshaler = (*Hello)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s Hello) xdrType() {}

var _ xdrType = (*Hello)(nil)

// Auth is an XDR Struct defines as:
//
//   struct Auth
//    {
//        // Empty message, just to confirm
//        // establishment of MAC keys.
//        int unused;
//    };
//
type Auth struct {
	Unused int32
}

// EncodeTo encodes this value using the Encoder.
func (s *Auth) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(s.Unused)); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Auth)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Auth) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	s.Unused, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Auth) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Auth) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Auth)(nil)
	_ encoding.BinaryUnmarshaler = (*Auth)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s Auth) xdrType() {}

var _ xdrType = (*Auth)(nil)

// IpAddrType is an XDR Enum defines as:
//
//   enum IPAddrType
//    {
//        IPv4 = 0,
//        IPv6 = 1
//    };
//
type IpAddrType int32

const (
	IpAddrTypeIPv4 IpAddrType = 0
	IpAddrTypeIPv6 IpAddrType = 1
)

var ipAddrTypeMap = map[int32]string{
	0: "IpAddrTypeIPv4",
	1: "IpAddrTypeIPv6",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for IpAddrType
func (e IpAddrType) ValidEnum(v int32) bool {
	_, ok := ipAddrTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e IpAddrType) String() string {
	name, _ := ipAddrTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e IpAddrType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := ipAddrTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid IpAddrType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*IpAddrType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *IpAddrType) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding IpAddrType: %s", err)
	}
	if _, ok := ipAddrTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid IpAddrType enum value", v)
	}
	*e = IpAddrType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s IpAddrType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *IpAddrType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*IpAddrType)(nil)
	_ encoding.BinaryUnmarshaler = (*IpAddrType)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s IpAddrType) xdrType() {}

var _ xdrType = (*IpAddrType)(nil)

// PeerAddressIp is an XDR NestedUnion defines as:
//
//   union switch (IPAddrType type)
//        {
//        case IPv4:
//            opaque ipv4[4];
//        case IPv6:
//            opaque ipv6[16];
//        }
//
type PeerAddressIp struct {
	Type IpAddrType
	Ipv4 *[4]byte  `xdrmaxsize:"4"`
	Ipv6 *[16]byte `xdrmaxsize:"16"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u PeerAddressIp) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PeerAddressIp
func (u PeerAddressIp) ArmForSwitch(sw int32) (string, bool) {
	switch IpAddrType(sw) {
	case IpAddrTypeIPv4:
		return "Ipv4", true
	case IpAddrTypeIPv6:
		return "Ipv6", true
	}
	return "-", false
}

// NewPeerAddressIp creates a new  PeerAddressIp.
func NewPeerAddressIp(aType IpAddrType, value interface{}) (result PeerAddressIp, err error) {
	result.Type = aType
	switch IpAddrType(aType) {
	case IpAddrTypeIPv4:
		tv, ok := value.([4]byte)
		if !ok {
			err = fmt.Errorf("invalid value, must be [4]byte")
			return
		}
		result.Ipv4 = &tv
	case IpAddrTypeIPv6:
		tv, ok := value.([16]byte)
		if !ok {
			err = fmt.Errorf("invalid value, must be [16]byte")
			return
		}
		result.Ipv6 = &tv
	}
	return
}

// MustIpv4 retrieves the Ipv4 value from the union,
// panicing if the value is not set.
func (u PeerAddressIp) MustIpv4() [4]byte {
	val, ok := u.GetIpv4()

	if !ok {
		panic("arm Ipv4 is not set")
	}

	return val
}

// GetIpv4 retrieves the Ipv4 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u PeerAddressIp) GetIpv4() (result [4]byte, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Ipv4" {
		result = *u.Ipv4
		ok = true
	}

	return
}

// MustIpv6 retrieves the Ipv6 value from the union,
// panicing if the value is not set.
func (u PeerAddressIp) MustIpv6() [16]byte {
	val, ok := u.GetIpv6()

	if !ok {
		panic("arm Ipv6 is not set")
	}

	return val
}

// GetIpv6 retrieves the Ipv6 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u PeerAddressIp) GetIpv6() (result [16]byte, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Ipv6" {
		result = *u.Ipv6
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u PeerAddressIp) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch IpAddrType(u.Type) {
	case IpAddrTypeIPv4:
		if _, err = e.EncodeFixedOpaque((*u.Ipv4)[:]); err != nil {
			return err
		}
		return nil
	case IpAddrTypeIPv6:
		if _, err = e.EncodeFixedOpaque((*u.Ipv6)[:]); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (IpAddrType) switch value '%d' is not valid for union PeerAddressIp", u.Type)
}

var _ decoderFrom = (*PeerAddressIp)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *PeerAddressIp) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding IpAddrType: %s", err)
	}
	switch IpAddrType(u.Type) {
	case IpAddrTypeIPv4:
		u.Ipv4 = new([4]byte)
		nTmp, err = d.DecodeFixedOpaqueInplace((*u.Ipv4)[:])
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Ipv4: %s", err)
		}
		return n, nil
	case IpAddrTypeIPv6:
		u.Ipv6 = new([16]byte)
		nTmp, err = d.DecodeFixedOpaqueInplace((*u.Ipv6)[:])
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Ipv6: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union PeerAddressIp has invalid Type (IpAddrType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s PeerAddressIp) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *PeerAddressIp) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*PeerAddressIp)(nil)
	_ encoding.BinaryUnmarshaler = (*PeerAddressIp)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s PeerAddressIp) xdrType() {}

var _ xdrType = (*PeerAddressIp)(nil)

// PeerAddress is an XDR Struct defines as:
//
//   struct PeerAddress
//    {
//        union switch (IPAddrType type)
//        {
//        case IPv4:
//            opaque ipv4[4];
//        case IPv6:
//            opaque ipv6[16];
//        }
//        ip;
//        uint32 port;
//        uint32 numFailures;
//    };
//
type PeerAddress struct {
	Ip          PeerAddressIp
	Port        Uint32
	NumFailures Uint32
}

// EncodeTo encodes this value using the Encoder.
func (s *PeerAddress) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Ip.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Port.EncodeTo(e); err != nil {
		return err
	}
	if err = s.NumFailures.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*PeerAddress)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *PeerAddress) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Ip.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding PeerAddressIp: %s", err)
	}
	nTmp, err = s.Port.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.NumFailures.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s PeerAddress) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *PeerAddress) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*PeerAddress)(nil)
	_ encoding.BinaryUnmarshaler = (*PeerAddress)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s PeerAddress) xdrType() {}

var _ xdrType = (*PeerAddress)(nil)

// MessageType is an XDR Enum defines as:
//
//   enum MessageType
//    {
//        ERROR_MSG = 0,
//        AUTH = 2,
//        DONT_HAVE = 3,
//
//        GET_PEERS = 4, // gets a list of peers this guy knows about
//        PEERS = 5,
//
//        GET_TX_SET = 6, // gets a particular txset by hash
//        TX_SET = 7,
//
//        TRANSACTION = 8, // pass on a tx you have heard about
//
//        // SCP
//        GET_SCP_QUORUMSET = 9,
//        SCP_QUORUMSET = 10,
//        SCP_MESSAGE = 11,
//        GET_SCP_STATE = 12,
//
//        // new messages
//        HELLO = 13,
//
//        SURVEY_REQUEST = 14,
//        SURVEY_RESPONSE = 15
//    };
//
type MessageType int32

const (
	MessageTypeErrorMsg        MessageType = 0
	MessageTypeAuth            MessageType = 2
	MessageTypeDontHave        MessageType = 3
	MessageTypeGetPeers        MessageType = 4
	MessageTypePeers           MessageType = 5
	MessageTypeGetTxSet        MessageType = 6
	MessageTypeTxSet           MessageType = 7
	MessageTypeTransaction     MessageType = 8
	MessageTypeGetScpQuorumset MessageType = 9
	MessageTypeScpQuorumset    MessageType = 10
	MessageTypeScpMessage      MessageType = 11
	MessageTypeGetScpState     MessageType = 12
	MessageTypeHello           MessageType = 13
	MessageTypeSurveyRequest   MessageType = 14
	MessageTypeSurveyResponse  MessageType = 15
)

var messageTypeMap = map[int32]string{
	0:  "MessageTypeErrorMsg",
	2:  "MessageTypeAuth",
	3:  "MessageTypeDontHave",
	4:  "MessageTypeGetPeers",
	5:  "MessageTypePeers",
	6:  "MessageTypeGetTxSet",
	7:  "MessageTypeTxSet",
	8:  "MessageTypeTransaction",
	9:  "MessageTypeGetScpQuorumset",
	10: "MessageTypeScpQuorumset",
	11: "MessageTypeScpMessage",
	12: "MessageTypeGetScpState",
	13: "MessageTypeHello",
	14: "MessageTypeSurveyRequest",
	15: "MessageTypeSurveyResponse",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for MessageType
func (e MessageType) ValidEnum(v int32) bool {
	_, ok := messageTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e MessageType) String() string {
	name, _ := messageTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e MessageType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := messageTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid MessageType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*MessageType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *MessageType) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding MessageType: %s", err)
	}
	if _, ok := messageTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid MessageType enum value", v)
	}
	*e = MessageType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s MessageType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *MessageType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*MessageType)(nil)
	_ encoding.BinaryUnmarshaler = (*MessageType)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s MessageType) xdrType() {}

var _ xdrType = (*MessageType)(nil)

// DontHave is an XDR Struct defines as:
//
//   struct DontHave
//    {
//        MessageType type;
//        uint256 reqHash;
//    };
//
type DontHave struct {
	Type    MessageType
	ReqHash Uint256
}

// EncodeTo encodes this value using the Encoder.
func (s *DontHave) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Type.EncodeTo(e); err != nil {
		return err
	}
	if err = s.ReqHash.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*DontHave)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *DontHave) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding MessageType: %s", err)
	}
	nTmp, err = s.ReqHash.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint256: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s DontHave) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *DontHave) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*DontHave)(nil)
	_ encoding.BinaryUnmarshaler = (*DontHave)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s DontHave) xdrType() {}

var _ xdrType = (*DontHave)(nil)

// SurveyMessageCommandType is an XDR Enum defines as:
//
//   enum SurveyMessageCommandType
//    {
//        SURVEY_TOPOLOGY = 0
//    };
//
type SurveyMessageCommandType int32

const (
	SurveyMessageCommandTypeSurveyTopology SurveyMessageCommandType = 0
)

var surveyMessageCommandTypeMap = map[int32]string{
	0: "SurveyMessageCommandTypeSurveyTopology",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for SurveyMessageCommandType
func (e SurveyMessageCommandType) ValidEnum(v int32) bool {
	_, ok := surveyMessageCommandTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e SurveyMessageCommandType) String() string {
	name, _ := surveyMessageCommandTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e SurveyMessageCommandType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := surveyMessageCommandTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid SurveyMessageCommandType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*SurveyMessageCommandType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *SurveyMessageCommandType) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding SurveyMessageCommandType: %s", err)
	}
	if _, ok := surveyMessageCommandTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid SurveyMessageCommandType enum value", v)
	}
	*e = SurveyMessageCommandType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s SurveyMessageCommandType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SurveyMessageCommandType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*SurveyMessageCommandType)(nil)
	_ encoding.BinaryUnmarshaler = (*SurveyMessageCommandType)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s SurveyMessageCommandType) xdrType() {}

var _ xdrType = (*SurveyMessageCommandType)(nil)

// SurveyRequestMessage is an XDR Struct defines as:
//
//   struct SurveyRequestMessage
//    {
//        NodeID surveyorPeerID;
//        NodeID surveyedPeerID;
//        uint32 ledgerNum;
//        Curve25519Public encryptionKey;
//        SurveyMessageCommandType commandType;
//    };
//
type SurveyRequestMessage struct {
	SurveyorPeerId NodeId
	SurveyedPeerId NodeId
	LedgerNum      Uint32
	EncryptionKey  Curve25519Public
	CommandType    SurveyMessageCommandType
}

// EncodeTo encodes this value using the Encoder.
func (s *SurveyRequestMessage) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.SurveyorPeerId.EncodeTo(e); err != nil {
		return err
	}
	if err = s.SurveyedPeerId.EncodeTo(e); err != nil {
		return err
	}
	if err = s.LedgerNum.EncodeTo(e); err != nil {
		return err
	}
	if err = s.EncryptionKey.EncodeTo(e); err != nil {
		return err
	}
	if err = s.CommandType.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*SurveyRequestMessage)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *SurveyRequestMessage) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.SurveyorPeerId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding NodeId: %s", err)
	}
	nTmp, err = s.SurveyedPeerId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding NodeId: %s", err)
	}
	nTmp, err = s.LedgerNum.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.EncryptionKey.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Curve25519Public: %s", err)
	}
	nTmp, err = s.CommandType.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SurveyMessageCommandType: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s SurveyRequestMessage) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SurveyRequestMessage) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*SurveyRequestMessage)(nil)
	_ encoding.BinaryUnmarshaler = (*SurveyRequestMessage)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s SurveyRequestMessage) xdrType() {}

var _ xdrType = (*SurveyRequestMessage)(nil)

// SignedSurveyRequestMessage is an XDR Struct defines as:
//
//   struct SignedSurveyRequestMessage
//    {
//        Signature requestSignature;
//        SurveyRequestMessage request;
//    };
//
type SignedSurveyRequestMessage struct {
	RequestSignature Signature
	Request          SurveyRequestMessage
}

// EncodeTo encodes this value using the Encoder.
func (s *SignedSurveyRequestMessage) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.RequestSignature.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Request.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*SignedSurveyRequestMessage)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *SignedSurveyRequestMessage) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.RequestSignature.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Signature: %s", err)
	}
	nTmp, err = s.Request.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SurveyRequestMessage: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s SignedSurveyRequestMessage) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SignedSurveyRequestMessage) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*SignedSurveyRequestMessage)(nil)
	_ encoding.BinaryUnmarshaler = (*SignedSurveyRequestMessage)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s SignedSurveyRequestMessage) xdrType() {}

var _ xdrType = (*SignedSurveyRequestMessage)(nil)

// EncryptedBody is an XDR Typedef defines as:
//
//   typedef opaque EncryptedBody<64000>;
//
type EncryptedBody []byte

// XDRMaxSize implements the Sized interface for EncryptedBody
func (e EncryptedBody) XDRMaxSize() int {
	return 64000
}

// EncodeTo encodes this value using the Encoder.
func (s EncryptedBody) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeOpaque(s[:]); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*EncryptedBody)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *EncryptedBody) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	(*s), nTmp, err = d.DecodeOpaque(64000)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding EncryptedBody: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s EncryptedBody) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *EncryptedBody) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*EncryptedBody)(nil)
	_ encoding.BinaryUnmarshaler = (*EncryptedBody)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s EncryptedBody) xdrType() {}

var _ xdrType = (*EncryptedBody)(nil)

// SurveyResponseMessage is an XDR Struct defines as:
//
//   struct SurveyResponseMessage
//    {
//        NodeID surveyorPeerID;
//        NodeID surveyedPeerID;
//        uint32 ledgerNum;
//        SurveyMessageCommandType commandType;
//        EncryptedBody encryptedBody;
//    };
//
type SurveyResponseMessage struct {
	SurveyorPeerId NodeId
	SurveyedPeerId NodeId
	LedgerNum      Uint32
	CommandType    SurveyMessageCommandType
	EncryptedBody  EncryptedBody
}

// EncodeTo encodes this value using the Encoder.
func (s *SurveyResponseMessage) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.SurveyorPeerId.EncodeTo(e); err != nil {
		return err
	}
	if err = s.SurveyedPeerId.EncodeTo(e); err != nil {
		return err
	}
	if err = s.LedgerNum.EncodeTo(e); err != nil {
		return err
	}
	if err = s.CommandType.EncodeTo(e); err != nil {
		return err
	}
	if err = s.EncryptedBody.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*SurveyResponseMessage)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *SurveyResponseMessage) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.SurveyorPeerId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding NodeId: %s", err)
	}
	nTmp, err = s.SurveyedPeerId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding NodeId: %s", err)
	}
	nTmp, err = s.LedgerNum.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.CommandType.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SurveyMessageCommandType: %s", err)
	}
	nTmp, err = s.EncryptedBody.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding EncryptedBody: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s SurveyResponseMessage) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SurveyResponseMessage) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*SurveyResponseMessage)(nil)
	_ encoding.BinaryUnmarshaler = (*SurveyResponseMessage)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s SurveyResponseMessage) xdrType() {}

var _ xdrType = (*SurveyResponseMessage)(nil)

// SignedSurveyResponseMessage is an XDR Struct defines as:
//
//   struct SignedSurveyResponseMessage
//    {
//        Signature responseSignature;
//        SurveyResponseMessage response;
//    };
//
type SignedSurveyResponseMessage struct {
	ResponseSignature Signature
	Response          SurveyResponseMessage
}

// EncodeTo encodes this value using the Encoder.
func (s *SignedSurveyResponseMessage) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.ResponseSignature.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Response.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*SignedSurveyResponseMessage)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *SignedSurveyResponseMessage) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.ResponseSignature.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Signature: %s", err)
	}
	nTmp, err = s.Response.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SurveyResponseMessage: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s SignedSurveyResponseMessage) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SignedSurveyResponseMessage) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*SignedSurveyResponseMessage)(nil)
	_ encoding.BinaryUnmarshaler = (*SignedSurveyResponseMessage)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s SignedSurveyResponseMessage) xdrType() {}

var _ xdrType = (*SignedSurveyResponseMessage)(nil)

// PeerStats is an XDR Struct defines as:
//
//   struct PeerStats
//    {
//        NodeID id;
//        string versionStr<100>;
//        uint64 messagesRead;
//        uint64 messagesWritten;
//        uint64 bytesRead;
//        uint64 bytesWritten;
//        uint64 secondsConnected;
//
//        uint64 uniqueFloodBytesRecv;
//        uint64 duplicateFloodBytesRecv;
//        uint64 uniqueFetchBytesRecv;
//        uint64 duplicateFetchBytesRecv;
//
//        uint64 uniqueFloodMessageRecv;
//        uint64 duplicateFloodMessageRecv;
//        uint64 uniqueFetchMessageRecv;
//        uint64 duplicateFetchMessageRecv;
//    };
//
type PeerStats struct {
	Id                        NodeId
	VersionStr                string `xdrmaxsize:"100"`
	MessagesRead              Uint64
	MessagesWritten           Uint64
	BytesRead                 Uint64
	BytesWritten              Uint64
	SecondsConnected          Uint64
	UniqueFloodBytesRecv      Uint64
	DuplicateFloodBytesRecv   Uint64
	UniqueFetchBytesRecv      Uint64
	DuplicateFetchBytesRecv   Uint64
	UniqueFloodMessageRecv    Uint64
	DuplicateFloodMessageRecv Uint64
	UniqueFetchMessageRecv    Uint64
	DuplicateFetchMessageRecv Uint64
}

// EncodeTo encodes this value using the Encoder.
func (s *PeerStats) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Id.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeString(string(s.VersionStr)); err != nil {
		return err
	}
	if err = s.MessagesRead.EncodeTo(e); err != nil {
		return err
	}
	if err = s.MessagesWritten.EncodeTo(e); err != nil {
		return err
	}
	if err = s.BytesRead.EncodeTo(e); err != nil {
		return err
	}
	if err = s.BytesWritten.EncodeTo(e); err != nil {
		return err
	}
	if err = s.SecondsConnected.EncodeTo(e); err != nil {
		return err
	}
	if err = s.UniqueFloodBytesRecv.EncodeTo(e); err != nil {
		return err
	}
	if err = s.DuplicateFloodBytesRecv.EncodeTo(e); err != nil {
		return err
	}
	if err = s.UniqueFetchBytesRecv.EncodeTo(e); err != nil {
		return err
	}
	if err = s.DuplicateFetchBytesRecv.EncodeTo(e); err != nil {
		return err
	}
	if err = s.UniqueFloodMessageRecv.EncodeTo(e); err != nil {
		return err
	}
	if err = s.DuplicateFloodMessageRecv.EncodeTo(e); err != nil {
		return err
	}
	if err = s.UniqueFetchMessageRecv.EncodeTo(e); err != nil {
		return err
	}
	if err = s.DuplicateFetchMessageRecv.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*PeerStats)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *PeerStats) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Id.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding NodeId: %s", err)
	}
	s.VersionStr, nTmp, err = d.DecodeString(100)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding VersionStr: %s", err)
	}
	nTmp, err = s.MessagesRead.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint64: %s", err)
	}
	nTmp, err = s.MessagesWritten.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint64: %s", err)
	}
	nTmp, err = s.BytesRead.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint64: %s", err)
	}
	nTmp, err = s.BytesWritten.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint64: %s", err)
	}
	nTmp, err = s.SecondsConnected.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint64: %s", err)
	}
	nTmp, err = s.UniqueFloodBytesRecv.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint64: %s", err)
	}
	nTmp, err = s.DuplicateFloodBytesRecv.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint64: %s", err)
	}
	nTmp, err = s.UniqueFetchBytesRecv.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint64: %s", err)
	}
	nTmp, err = s.DuplicateFetchBytesRecv.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint64: %s", err)
	}
	nTmp, err = s.UniqueFloodMessageRecv.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint64: %s", err)
	}
	nTmp, err = s.DuplicateFloodMessageRecv.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint64: %s", err)
	}
	nTmp, err = s.UniqueFetchMessageRecv.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint64: %s", err)
	}
	nTmp, err = s.DuplicateFetchMessageRecv.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint64: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s PeerStats) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *PeerStats) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*PeerStats)(nil)
	_ encoding.BinaryUnmarshaler = (*PeerStats)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s PeerStats) xdrType() {}

var _ xdrType = (*PeerStats)(nil)

// PeerStatList is an XDR Typedef defines as:
//
//   typedef PeerStats PeerStatList<25>;
//
type PeerStatList []PeerStats

// XDRMaxSize implements the Sized interface for PeerStatList
func (e PeerStatList) XDRMaxSize() int {
	return 25
}

// EncodeTo encodes this value using the Encoder.
func (s PeerStatList) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeUint(uint32(len(s))); err != nil {
		return err
	}
	for i := 0; i < len(s); i++ {
		if err = s[i].EncodeTo(e); err != nil {
			return err
		}
	}
	return nil
}

var _ decoderFrom = (*PeerStatList)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *PeerStatList) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding PeerStats: %s", err)
	}
	if l > 25 {
		return n, fmt.Errorf("decoding PeerStats: data size (%d) exceeds size limit (25)", l)
	}
	(*s) = nil
	if l > 0 {
		(*s) = make([]PeerStats, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = (*s)[i].DecodeFrom(d)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding PeerStats: %s", err)
			}
		}
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s PeerStatList) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *PeerStatList) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*PeerStatList)(nil)
	_ encoding.BinaryUnmarshaler = (*PeerStatList)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s PeerStatList) xdrType() {}

var _ xdrType = (*PeerStatList)(nil)

// TopologyResponseBody is an XDR Struct defines as:
//
//   struct TopologyResponseBody
//    {
//        PeerStatList inboundPeers;
//        PeerStatList outboundPeers;
//
//        uint32 totalInboundPeerCount;
//        uint32 totalOutboundPeerCount;
//    };
//
type TopologyResponseBody struct {
	InboundPeers           PeerStatList
	OutboundPeers          PeerStatList
	TotalInboundPeerCount  Uint32
	TotalOutboundPeerCount Uint32
}

// EncodeTo encodes this value using the Encoder.
func (s *TopologyResponseBody) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.InboundPeers.EncodeTo(e); err != nil {
		return err
	}
	if err = s.OutboundPeers.EncodeTo(e); err != nil {
		return err
	}
	if err = s.TotalInboundPeerCount.EncodeTo(e); err != nil {
		return err
	}
	if err = s.TotalOutboundPeerCount.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*TopologyResponseBody)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *TopologyResponseBody) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.InboundPeers.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding PeerStatList: %s", err)
	}
	nTmp, err = s.OutboundPeers.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding PeerStatList: %s", err)
	}
	nTmp, err = s.TotalInboundPeerCount.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.TotalOutboundPeerCount.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TopologyResponseBody) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TopologyResponseBody) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TopologyResponseBody)(nil)
	_ encoding.BinaryUnmarshaler = (*TopologyResponseBody)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TopologyResponseBody) xdrType() {}

var _ xdrType = (*TopologyResponseBody)(nil)

// SurveyResponseBody is an XDR Union defines as:
//
//   union SurveyResponseBody switch (SurveyMessageCommandType type)
//    {
//    case SURVEY_TOPOLOGY:
//        TopologyResponseBody topologyResponseBody;
//    };
//
type SurveyResponseBody struct {
	Type                 SurveyMessageCommandType
	TopologyResponseBody *TopologyResponseBody
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u SurveyResponseBody) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of SurveyResponseBody
func (u SurveyResponseBody) ArmForSwitch(sw int32) (string, bool) {
	switch SurveyMessageCommandType(sw) {
	case SurveyMessageCommandTypeSurveyTopology:
		return "TopologyResponseBody", true
	}
	return "-", false
}

// NewSurveyResponseBody creates a new  SurveyResponseBody.
func NewSurveyResponseBody(aType SurveyMessageCommandType, value interface{}) (result SurveyResponseBody, err error) {
	result.Type = aType
	switch SurveyMessageCommandType(aType) {
	case SurveyMessageCommandTypeSurveyTopology:
		tv, ok := value.(TopologyResponseBody)
		if !ok {
			err = fmt.Errorf("invalid value, must be TopologyResponseBody")
			return
		}
		result.TopologyResponseBody = &tv
	}
	return
}

// MustTopologyResponseBody retrieves the TopologyResponseBody value from the union,
// panicing if the value is not set.
func (u SurveyResponseBody) MustTopologyResponseBody() TopologyResponseBody {
	val, ok := u.GetTopologyResponseBody()

	if !ok {
		panic("arm TopologyResponseBody is not set")
	}

	return val
}

// GetTopologyResponseBody retrieves the TopologyResponseBody value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u SurveyResponseBody) GetTopologyResponseBody() (result TopologyResponseBody, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "TopologyResponseBody" {
		result = *u.TopologyResponseBody
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u SurveyResponseBody) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch SurveyMessageCommandType(u.Type) {
	case SurveyMessageCommandTypeSurveyTopology:
		if err = (*u.TopologyResponseBody).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (SurveyMessageCommandType) switch value '%d' is not valid for union SurveyResponseBody", u.Type)
}

var _ decoderFrom = (*SurveyResponseBody)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *SurveyResponseBody) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SurveyMessageCommandType: %s", err)
	}
	switch SurveyMessageCommandType(u.Type) {
	case SurveyMessageCommandTypeSurveyTopology:
		u.TopologyResponseBody = new(TopologyResponseBody)
		nTmp, err = (*u.TopologyResponseBody).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding TopologyResponseBody: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union SurveyResponseBody has invalid Type (SurveyMessageCommandType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s SurveyResponseBody) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SurveyResponseBody) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*SurveyResponseBody)(nil)
	_ encoding.BinaryUnmarshaler = (*SurveyResponseBody)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s SurveyResponseBody) xdrType() {}

var _ xdrType = (*SurveyResponseBody)(nil)

// diamcircleMessage is an XDR Union defines as:
//
//   union diamcircleMessage switch (MessageType type)
//    {
//    case ERROR_MSG:
//        Error error;
//    case HELLO:
//        Hello hello;
//    case AUTH:
//        Auth auth;
//    case DONT_HAVE:
//        DontHave dontHave;
//    case GET_PEERS:
//        void;
//    case PEERS:
//        PeerAddress peers<100>;
//
//    case GET_TX_SET:
//        uint256 txSetHash;
//    case TX_SET:
//        TransactionSet txSet;
//
//    case TRANSACTION:
//        TransactionEnvelope transaction;
//
//    case SURVEY_REQUEST:
//        SignedSurveyRequestMessage signedSurveyRequestMessage;
//
//    case SURVEY_RESPONSE:
//        SignedSurveyResponseMessage signedSurveyResponseMessage;
//
//    // SCP
//    case GET_SCP_QUORUMSET:
//        uint256 qSetHash;
//    case SCP_QUORUMSET:
//        SCPQuorumSet qSet;
//    case SCP_MESSAGE:
//        SCPEnvelope envelope;
//    case GET_SCP_STATE:
//        uint32 getSCPLedgerSeq; // ledger seq requested ; if 0, requests the latest
//    };
//
type diamcircleMessage struct {
	Type                        MessageType
	Error                       *Error
	Hello                       *Hello
	Auth                        *Auth
	DontHave                    *DontHave
	Peers                       *[]PeerAddress `xdrmaxsize:"100"`
	TxSetHash                   *Uint256
	TxSet                       *TransactionSet
	Transaction                 *TransactionEnvelope
	SignedSurveyRequestMessage  *SignedSurveyRequestMessage
	SignedSurveyResponseMessage *SignedSurveyResponseMessage
	QSetHash                    *Uint256
	QSet                        *ScpQuorumSet
	Envelope                    *ScpEnvelope
	GetScpLedgerSeq             *Uint32
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u diamcircleMessage) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of diamcircleMessage
func (u diamcircleMessage) ArmForSwitch(sw int32) (string, bool) {
	switch MessageType(sw) {
	case MessageTypeErrorMsg:
		return "Error", true
	case MessageTypeHello:
		return "Hello", true
	case MessageTypeAuth:
		return "Auth", true
	case MessageTypeDontHave:
		return "DontHave", true
	case MessageTypeGetPeers:
		return "", true
	case MessageTypePeers:
		return "Peers", true
	case MessageTypeGetTxSet:
		return "TxSetHash", true
	case MessageTypeTxSet:
		return "TxSet", true
	case MessageTypeTransaction:
		return "Transaction", true
	case MessageTypeSurveyRequest:
		return "SignedSurveyRequestMessage", true
	case MessageTypeSurveyResponse:
		return "SignedSurveyResponseMessage", true
	case MessageTypeGetScpQuorumset:
		return "QSetHash", true
	case MessageTypeScpQuorumset:
		return "QSet", true
	case MessageTypeScpMessage:
		return "Envelope", true
	case MessageTypeGetScpState:
		return "GetScpLedgerSeq", true
	}
	return "-", false
}

// NewdiamcircleMessage creates a new  diamcircleMessage.
func NewdiamcircleMessage(aType MessageType, value interface{}) (result diamcircleMessage, err error) {
	result.Type = aType
	switch MessageType(aType) {
	case MessageTypeErrorMsg:
		tv, ok := value.(Error)
		if !ok {
			err = fmt.Errorf("invalid value, must be Error")
			return
		}
		result.Error = &tv
	case MessageTypeHello:
		tv, ok := value.(Hello)
		if !ok {
			err = fmt.Errorf("invalid value, must be Hello")
			return
		}
		result.Hello = &tv
	case MessageTypeAuth:
		tv, ok := value.(Auth)
		if !ok {
			err = fmt.Errorf("invalid value, must be Auth")
			return
		}
		result.Auth = &tv
	case MessageTypeDontHave:
		tv, ok := value.(DontHave)
		if !ok {
			err = fmt.Errorf("invalid value, must be DontHave")
			return
		}
		result.DontHave = &tv
	case MessageTypeGetPeers:
		// void
	case MessageTypePeers:
		tv, ok := value.([]PeerAddress)
		if !ok {
			err = fmt.Errorf("invalid value, must be []PeerAddress")
			return
		}
		result.Peers = &tv
	case MessageTypeGetTxSet:
		tv, ok := value.(Uint256)
		if !ok {
			err = fmt.Errorf("invalid value, must be Uint256")
			return
		}
		result.TxSetHash = &tv
	case MessageTypeTxSet:
		tv, ok := value.(TransactionSet)
		if !ok {
			err = fmt.Errorf("invalid value, must be TransactionSet")
			return
		}
		result.TxSet = &tv
	case MessageTypeTransaction:
		tv, ok := value.(TransactionEnvelope)
		if !ok {
			err = fmt.Errorf("invalid value, must be TransactionEnvelope")
			return
		}
		result.Transaction = &tv
	case MessageTypeSurveyRequest:
		tv, ok := value.(SignedSurveyRequestMessage)
		if !ok {
			err = fmt.Errorf("invalid value, must be SignedSurveyRequestMessage")
			return
		}
		result.SignedSurveyRequestMessage = &tv
	case MessageTypeSurveyResponse:
		tv, ok := value.(SignedSurveyResponseMessage)
		if !ok {
			err = fmt.Errorf("invalid value, must be SignedSurveyResponseMessage")
			return
		}
		result.SignedSurveyResponseMessage = &tv
	case MessageTypeGetScpQuorumset:
		tv, ok := value.(Uint256)
		if !ok {
			err = fmt.Errorf("invalid value, must be Uint256")
			return
		}
		result.QSetHash = &tv
	case MessageTypeScpQuorumset:
		tv, ok := value.(ScpQuorumSet)
		if !ok {
			err = fmt.Errorf("invalid value, must be ScpQuorumSet")
			return
		}
		result.QSet = &tv
	case MessageTypeScpMessage:
		tv, ok := value.(ScpEnvelope)
		if !ok {
			err = fmt.Errorf("invalid value, must be ScpEnvelope")
			return
		}
		result.Envelope = &tv
	case MessageTypeGetScpState:
		tv, ok := value.(Uint32)
		if !ok {
			err = fmt.Errorf("invalid value, must be Uint32")
			return
		}
		result.GetScpLedgerSeq = &tv
	}
	return
}

// MustError retrieves the Error value from the union,
// panicing if the value is not set.
func (u diamcircleMessage) MustError() Error {
	val, ok := u.GetError()

	if !ok {
		panic("arm Error is not set")
	}

	return val
}

// GetError retrieves the Error value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u diamcircleMessage) GetError() (result Error, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Error" {
		result = *u.Error
		ok = true
	}

	return
}

// MustHello retrieves the Hello value from the union,
// panicing if the value is not set.
func (u diamcircleMessage) MustHello() Hello {
	val, ok := u.GetHello()

	if !ok {
		panic("arm Hello is not set")
	}

	return val
}

// GetHello retrieves the Hello value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u diamcircleMessage) GetHello() (result Hello, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Hello" {
		result = *u.Hello
		ok = true
	}

	return
}

// MustAuth retrieves the Auth value from the union,
// panicing if the value is not set.
func (u diamcircleMessage) MustAuth() Auth {
	val, ok := u.GetAuth()

	if !ok {
		panic("arm Auth is not set")
	}

	return val
}

// GetAuth retrieves the Auth value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u diamcircleMessage) GetAuth() (result Auth, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Auth" {
		result = *u.Auth
		ok = true
	}

	return
}

// MustDontHave retrieves the DontHave value from the union,
// panicing if the value is not set.
func (u diamcircleMessage) MustDontHave() DontHave {
	val, ok := u.GetDontHave()

	if !ok {
		panic("arm DontHave is not set")
	}

	return val
}

// GetDontHave retrieves the DontHave value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u diamcircleMessage) GetDontHave() (result DontHave, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "DontHave" {
		result = *u.DontHave
		ok = true
	}

	return
}

// MustPeers retrieves the Peers value from the union,
// panicing if the value is not set.
func (u diamcircleMessage) MustPeers() []PeerAddress {
	val, ok := u.GetPeers()

	if !ok {
		panic("arm Peers is not set")
	}

	return val
}

// GetPeers retrieves the Peers value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u diamcircleMessage) GetPeers() (result []PeerAddress, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Peers" {
		result = *u.Peers
		ok = true
	}

	return
}

// MustTxSetHash retrieves the TxSetHash value from the union,
// panicing if the value is not set.
func (u diamcircleMessage) MustTxSetHash() Uint256 {
	val, ok := u.GetTxSetHash()

	if !ok {
		panic("arm TxSetHash is not set")
	}

	return val
}

// GetTxSetHash retrieves the TxSetHash value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u diamcircleMessage) GetTxSetHash() (result Uint256, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "TxSetHash" {
		result = *u.TxSetHash
		ok = true
	}

	return
}

// MustTxSet retrieves the TxSet value from the union,
// panicing if the value is not set.
func (u diamcircleMessage) MustTxSet() TransactionSet {
	val, ok := u.GetTxSet()

	if !ok {
		panic("arm TxSet is not set")
	}

	return val
}

// GetTxSet retrieves the TxSet value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u diamcircleMessage) GetTxSet() (result TransactionSet, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "TxSet" {
		result = *u.TxSet
		ok = true
	}

	return
}

// MustTransaction retrieves the Transaction value from the union,
// panicing if the value is not set.
func (u diamcircleMessage) MustTransaction() TransactionEnvelope {
	val, ok := u.GetTransaction()

	if !ok {
		panic("arm Transaction is not set")
	}

	return val
}

// GetTransaction retrieves the Transaction value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u diamcircleMessage) GetTransaction() (result TransactionEnvelope, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Transaction" {
		result = *u.Transaction
		ok = true
	}

	return
}

// MustSignedSurveyRequestMessage retrieves the SignedSurveyRequestMessage value from the union,
// panicing if the value is not set.
func (u diamcircleMessage) MustSignedSurveyRequestMessage() SignedSurveyRequestMessage {
	val, ok := u.GetSignedSurveyRequestMessage()

	if !ok {
		panic("arm SignedSurveyRequestMessage is not set")
	}

	return val
}

// GetSignedSurveyRequestMessage retrieves the SignedSurveyRequestMessage value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u diamcircleMessage) GetSignedSurveyRequestMessage() (result SignedSurveyRequestMessage, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "SignedSurveyRequestMessage" {
		result = *u.SignedSurveyRequestMessage
		ok = true
	}

	return
}

// MustSignedSurveyResponseMessage retrieves the SignedSurveyResponseMessage value from the union,
// panicing if the value is not set.
func (u diamcircleMessage) MustSignedSurveyResponseMessage() SignedSurveyResponseMessage {
	val, ok := u.GetSignedSurveyResponseMessage()

	if !ok {
		panic("arm SignedSurveyResponseMessage is not set")
	}

	return val
}

// GetSignedSurveyResponseMessage retrieves the SignedSurveyResponseMessage value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u diamcircleMessage) GetSignedSurveyResponseMessage() (result SignedSurveyResponseMessage, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "SignedSurveyResponseMessage" {
		result = *u.SignedSurveyResponseMessage
		ok = true
	}

	return
}

// MustQSetHash retrieves the QSetHash value from the union,
// panicing if the value is not set.
func (u diamcircleMessage) MustQSetHash() Uint256 {
	val, ok := u.GetQSetHash()

	if !ok {
		panic("arm QSetHash is not set")
	}

	return val
}

// GetQSetHash retrieves the QSetHash value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u diamcircleMessage) GetQSetHash() (result Uint256, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "QSetHash" {
		result = *u.QSetHash
		ok = true
	}

	return
}

// MustQSet retrieves the QSet value from the union,
// panicing if the value is not set.
func (u diamcircleMessage) MustQSet() ScpQuorumSet {
	val, ok := u.GetQSet()

	if !ok {
		panic("arm QSet is not set")
	}

	return val
}

// GetQSet retrieves the QSet value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u diamcircleMessage) GetQSet() (result ScpQuorumSet, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "QSet" {
		result = *u.QSet
		ok = true
	}

	return
}

// MustEnvelope retrieves the Envelope value from the union,
// panicing if the value is not set.
func (u diamcircleMessage) MustEnvelope() ScpEnvelope {
	val, ok := u.GetEnvelope()

	if !ok {
		panic("arm Envelope is not set")
	}

	return val
}

// GetEnvelope retrieves the Envelope value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u diamcircleMessage) GetEnvelope() (result ScpEnvelope, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Envelope" {
		result = *u.Envelope
		ok = true
	}

	return
}

// MustGetScpLedgerSeq retrieves the GetScpLedgerSeq value from the union,
// panicing if the value is not set.
func (u diamcircleMessage) MustGetScpLedgerSeq() Uint32 {
	val, ok := u.GetGetScpLedgerSeq()

	if !ok {
		panic("arm GetScpLedgerSeq is not set")
	}

	return val
}

// GetGetScpLedgerSeq retrieves the GetScpLedgerSeq value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u diamcircleMessage) GetGetScpLedgerSeq() (result Uint32, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "GetScpLedgerSeq" {
		result = *u.GetScpLedgerSeq
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u diamcircleMessage) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch MessageType(u.Type) {
	case MessageTypeErrorMsg:
		if err = (*u.Error).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case MessageTypeHello:
		if err = (*u.Hello).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case MessageTypeAuth:
		if err = (*u.Auth).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case MessageTypeDontHave:
		if err = (*u.DontHave).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case MessageTypeGetPeers:
		// Void
		return nil
	case MessageTypePeers:
		if _, err = e.EncodeUint(uint32(len((*u.Peers)))); err != nil {
			return err
		}
		for i := 0; i < len((*u.Peers)); i++ {
			if err = (*u.Peers)[i].EncodeTo(e); err != nil {
				return err
			}
		}
		return nil
	case MessageTypeGetTxSet:
		if err = (*u.TxSetHash).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case MessageTypeTxSet:
		if err = (*u.TxSet).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case MessageTypeTransaction:
		if err = (*u.Transaction).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case MessageTypeSurveyRequest:
		if err = (*u.SignedSurveyRequestMessage).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case MessageTypeSurveyResponse:
		if err = (*u.SignedSurveyResponseMessage).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case MessageTypeGetScpQuorumset:
		if err = (*u.QSetHash).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case MessageTypeScpQuorumset:
		if err = (*u.QSet).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case MessageTypeScpMessage:
		if err = (*u.Envelope).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case MessageTypeGetScpState:
		if err = (*u.GetScpLedgerSeq).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (MessageType) switch value '%d' is not valid for union diamcircleMessage", u.Type)
}

var _ decoderFrom = (*diamcircleMessage)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *diamcircleMessage) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding MessageType: %s", err)
	}
	switch MessageType(u.Type) {
	case MessageTypeErrorMsg:
		u.Error = new(Error)
		nTmp, err = (*u.Error).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Error: %s", err)
		}
		return n, nil
	case MessageTypeHello:
		u.Hello = new(Hello)
		nTmp, err = (*u.Hello).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Hello: %s", err)
		}
		return n, nil
	case MessageTypeAuth:
		u.Auth = new(Auth)
		nTmp, err = (*u.Auth).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Auth: %s", err)
		}
		return n, nil
	case MessageTypeDontHave:
		u.DontHave = new(DontHave)
		nTmp, err = (*u.DontHave).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding DontHave: %s", err)
		}
		return n, nil
	case MessageTypeGetPeers:
		// Void
		return n, nil
	case MessageTypePeers:
		u.Peers = new([]PeerAddress)
		var l uint32
		l, nTmp, err = d.DecodeUint()
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding PeerAddress: %s", err)
		}
		if l > 100 {
			return n, fmt.Errorf("decoding PeerAddress: data size (%d) exceeds size limit (100)", l)
		}
		(*u.Peers) = nil
		if l > 0 {
			(*u.Peers) = make([]PeerAddress, l)
			for i := uint32(0); i < l; i++ {
				nTmp, err = (*u.Peers)[i].DecodeFrom(d)
				n += nTmp
				if err != nil {
					return n, fmt.Errorf("decoding PeerAddress: %s", err)
				}
			}
		}
		return n, nil
	case MessageTypeGetTxSet:
		u.TxSetHash = new(Uint256)
		nTmp, err = (*u.TxSetHash).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Uint256: %s", err)
		}
		return n, nil
	case MessageTypeTxSet:
		u.TxSet = new(TransactionSet)
		nTmp, err = (*u.TxSet).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding TransactionSet: %s", err)
		}
		return n, nil
	case MessageTypeTransaction:
		u.Transaction = new(TransactionEnvelope)
		nTmp, err = (*u.Transaction).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding TransactionEnvelope: %s", err)
		}
		return n, nil
	case MessageTypeSurveyRequest:
		u.SignedSurveyRequestMessage = new(SignedSurveyRequestMessage)
		nTmp, err = (*u.SignedSurveyRequestMessage).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding SignedSurveyRequestMessage: %s", err)
		}
		return n, nil
	case MessageTypeSurveyResponse:
		u.SignedSurveyResponseMessage = new(SignedSurveyResponseMessage)
		nTmp, err = (*u.SignedSurveyResponseMessage).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding SignedSurveyResponseMessage: %s", err)
		}
		return n, nil
	case MessageTypeGetScpQuorumset:
		u.QSetHash = new(Uint256)
		nTmp, err = (*u.QSetHash).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Uint256: %s", err)
		}
		return n, nil
	case MessageTypeScpQuorumset:
		u.QSet = new(ScpQuorumSet)
		nTmp, err = (*u.QSet).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ScpQuorumSet: %s", err)
		}
		return n, nil
	case MessageTypeScpMessage:
		u.Envelope = new(ScpEnvelope)
		nTmp, err = (*u.Envelope).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ScpEnvelope: %s", err)
		}
		return n, nil
	case MessageTypeGetScpState:
		u.GetScpLedgerSeq = new(Uint32)
		nTmp, err = (*u.GetScpLedgerSeq).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Uint32: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union diamcircleMessage has invalid Type (MessageType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s diamcircleMessage) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *diamcircleMessage) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*diamcircleMessage)(nil)
	_ encoding.BinaryUnmarshaler = (*diamcircleMessage)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s diamcircleMessage) xdrType() {}

var _ xdrType = (*diamcircleMessage)(nil)

// AuthenticatedMessageV0 is an XDR NestedStruct defines as:
//
//   struct
//        {
//            uint64 sequence;
//            diamcircleMessage message;
//            HmacSha256Mac mac;
//        }
//
type AuthenticatedMessageV0 struct {
	Sequence Uint64
	Message  diamcircleMessage
	Mac      HmacSha256Mac
}

// EncodeTo encodes this value using the Encoder.
func (s *AuthenticatedMessageV0) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Sequence.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Message.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Mac.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*AuthenticatedMessageV0)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *AuthenticatedMessageV0) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Sequence.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint64: %s", err)
	}
	nTmp, err = s.Message.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding diamcircleMessage: %s", err)
	}
	nTmp, err = s.Mac.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding HmacSha256Mac: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s AuthenticatedMessageV0) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AuthenticatedMessageV0) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*AuthenticatedMessageV0)(nil)
	_ encoding.BinaryUnmarshaler = (*AuthenticatedMessageV0)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s AuthenticatedMessageV0) xdrType() {}

var _ xdrType = (*AuthenticatedMessageV0)(nil)

// AuthenticatedMessage is an XDR Union defines as:
//
//   union AuthenticatedMessage switch (uint32 v)
//    {
//    case 0:
//        struct
//        {
//            uint64 sequence;
//            diamcircleMessage message;
//            HmacSha256Mac mac;
//        } v0;
//    };
//
type AuthenticatedMessage struct {
	V  Uint32
	V0 *AuthenticatedMessageV0
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u AuthenticatedMessage) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of AuthenticatedMessage
func (u AuthenticatedMessage) ArmForSwitch(sw int32) (string, bool) {
	switch Uint32(sw) {
	case 0:
		return "V0", true
	}
	return "-", false
}

// NewAuthenticatedMessage creates a new  AuthenticatedMessage.
func NewAuthenticatedMessage(v Uint32, value interface{}) (result AuthenticatedMessage, err error) {
	result.V = v
	switch Uint32(v) {
	case 0:
		tv, ok := value.(AuthenticatedMessageV0)
		if !ok {
			err = fmt.Errorf("invalid value, must be AuthenticatedMessageV0")
			return
		}
		result.V0 = &tv
	}
	return
}

// MustV0 retrieves the V0 value from the union,
// panicing if the value is not set.
func (u AuthenticatedMessage) MustV0() AuthenticatedMessageV0 {
	val, ok := u.GetV0()

	if !ok {
		panic("arm V0 is not set")
	}

	return val
}

// GetV0 retrieves the V0 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u AuthenticatedMessage) GetV0() (result AuthenticatedMessageV0, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.V))

	if armName == "V0" {
		result = *u.V0
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u AuthenticatedMessage) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.V.EncodeTo(e); err != nil {
		return err
	}
	switch Uint32(u.V) {
	case 0:
		if err = (*u.V0).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("V (Uint32) switch value '%d' is not valid for union AuthenticatedMessage", u.V)
}

var _ decoderFrom = (*AuthenticatedMessage)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *AuthenticatedMessage) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.V.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	switch Uint32(u.V) {
	case 0:
		u.V0 = new(AuthenticatedMessageV0)
		nTmp, err = (*u.V0).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding AuthenticatedMessageV0: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union AuthenticatedMessage has invalid V (Uint32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s AuthenticatedMessage) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AuthenticatedMessage) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*AuthenticatedMessage)(nil)
	_ encoding.BinaryUnmarshaler = (*AuthenticatedMessage)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s AuthenticatedMessage) xdrType() {}

var _ xdrType = (*AuthenticatedMessage)(nil)

// LiquidityPoolParameters is an XDR Union defines as:
//
//   union LiquidityPoolParameters switch (LiquidityPoolType type)
//    {
//    case LIQUIDITY_POOL_CONSTANT_PRODUCT:
//        LiquidityPoolConstantProductParameters constantProduct;
//    };
//
type LiquidityPoolParameters struct {
	Type            LiquidityPoolType
	ConstantProduct *LiquidityPoolConstantProductParameters
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LiquidityPoolParameters) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LiquidityPoolParameters
func (u LiquidityPoolParameters) ArmForSwitch(sw int32) (string, bool) {
	switch LiquidityPoolType(sw) {
	case LiquidityPoolTypeLiquidityPoolConstantProduct:
		return "ConstantProduct", true
	}
	return "-", false
}

// NewLiquidityPoolParameters creates a new  LiquidityPoolParameters.
func NewLiquidityPoolParameters(aType LiquidityPoolType, value interface{}) (result LiquidityPoolParameters, err error) {
	result.Type = aType
	switch LiquidityPoolType(aType) {
	case LiquidityPoolTypeLiquidityPoolConstantProduct:
		tv, ok := value.(LiquidityPoolConstantProductParameters)
		if !ok {
			err = fmt.Errorf("invalid value, must be LiquidityPoolConstantProductParameters")
			return
		}
		result.ConstantProduct = &tv
	}
	return
}

// MustConstantProduct retrieves the ConstantProduct value from the union,
// panicing if the value is not set.
func (u LiquidityPoolParameters) MustConstantProduct() LiquidityPoolConstantProductParameters {
	val, ok := u.GetConstantProduct()

	if !ok {
		panic("arm ConstantProduct is not set")
	}

	return val
}

// GetConstantProduct retrieves the ConstantProduct value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LiquidityPoolParameters) GetConstantProduct() (result LiquidityPoolConstantProductParameters, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ConstantProduct" {
		result = *u.ConstantProduct
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u LiquidityPoolParameters) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch LiquidityPoolType(u.Type) {
	case LiquidityPoolTypeLiquidityPoolConstantProduct:
		if err = (*u.ConstantProduct).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (LiquidityPoolType) switch value '%d' is not valid for union LiquidityPoolParameters", u.Type)
}

var _ decoderFrom = (*LiquidityPoolParameters)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *LiquidityPoolParameters) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LiquidityPoolType: %s", err)
	}
	switch LiquidityPoolType(u.Type) {
	case LiquidityPoolTypeLiquidityPoolConstantProduct:
		u.ConstantProduct = new(LiquidityPoolConstantProductParameters)
		nTmp, err = (*u.ConstantProduct).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding LiquidityPoolConstantProductParameters: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union LiquidityPoolParameters has invalid Type (LiquidityPoolType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LiquidityPoolParameters) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LiquidityPoolParameters) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LiquidityPoolParameters)(nil)
	_ encoding.BinaryUnmarshaler = (*LiquidityPoolParameters)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LiquidityPoolParameters) xdrType() {}

var _ xdrType = (*LiquidityPoolParameters)(nil)

// MuxedAccountMed25519 is an XDR NestedStruct defines as:
//
//   struct
//        {
//            uint64 id;
//            uint256 ed25519;
//        }
//
type MuxedAccountMed25519 struct {
	Id      Uint64
	Ed25519 Uint256
}

// EncodeTo encodes this value using the Encoder.
func (s *MuxedAccountMed25519) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Id.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Ed25519.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*MuxedAccountMed25519)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *MuxedAccountMed25519) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Id.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint64: %s", err)
	}
	nTmp, err = s.Ed25519.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint256: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s MuxedAccountMed25519) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *MuxedAccountMed25519) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*MuxedAccountMed25519)(nil)
	_ encoding.BinaryUnmarshaler = (*MuxedAccountMed25519)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s MuxedAccountMed25519) xdrType() {}

var _ xdrType = (*MuxedAccountMed25519)(nil)

// MuxedAccount is an XDR Union defines as:
//
//   union MuxedAccount switch (CryptoKeyType type)
//    {
//    case KEY_TYPE_ED25519:
//        uint256 ed25519;
//    case KEY_TYPE_MUXED_ED25519:
//        struct
//        {
//            uint64 id;
//            uint256 ed25519;
//        } med25519;
//    };
//
type MuxedAccount struct {
	Type     CryptoKeyType
	Ed25519  *Uint256
	Med25519 *MuxedAccountMed25519
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u MuxedAccount) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of MuxedAccount
func (u MuxedAccount) ArmForSwitch(sw int32) (string, bool) {
	switch CryptoKeyType(sw) {
	case CryptoKeyTypeKeyTypeEd25519:
		return "Ed25519", true
	case CryptoKeyTypeKeyTypeMuxedEd25519:
		return "Med25519", true
	}
	return "-", false
}

// NewMuxedAccount creates a new  MuxedAccount.
func NewMuxedAccount(aType CryptoKeyType, value interface{}) (result MuxedAccount, err error) {
	result.Type = aType
	switch CryptoKeyType(aType) {
	case CryptoKeyTypeKeyTypeEd25519:
		tv, ok := value.(Uint256)
		if !ok {
			err = fmt.Errorf("invalid value, must be Uint256")
			return
		}
		result.Ed25519 = &tv
	case CryptoKeyTypeKeyTypeMuxedEd25519:
		tv, ok := value.(MuxedAccountMed25519)
		if !ok {
			err = fmt.Errorf("invalid value, must be MuxedAccountMed25519")
			return
		}
		result.Med25519 = &tv
	}
	return
}

// MustEd25519 retrieves the Ed25519 value from the union,
// panicing if the value is not set.
func (u MuxedAccount) MustEd25519() Uint256 {
	val, ok := u.GetEd25519()

	if !ok {
		panic("arm Ed25519 is not set")
	}

	return val
}

// GetEd25519 retrieves the Ed25519 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u MuxedAccount) GetEd25519() (result Uint256, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Ed25519" {
		result = *u.Ed25519
		ok = true
	}

	return
}

// MustMed25519 retrieves the Med25519 value from the union,
// panicing if the value is not set.
func (u MuxedAccount) MustMed25519() MuxedAccountMed25519 {
	val, ok := u.GetMed25519()

	if !ok {
		panic("arm Med25519 is not set")
	}

	return val
}

// GetMed25519 retrieves the Med25519 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u MuxedAccount) GetMed25519() (result MuxedAccountMed25519, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Med25519" {
		result = *u.Med25519
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u MuxedAccount) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch CryptoKeyType(u.Type) {
	case CryptoKeyTypeKeyTypeEd25519:
		if err = (*u.Ed25519).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case CryptoKeyTypeKeyTypeMuxedEd25519:
		if err = (*u.Med25519).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (CryptoKeyType) switch value '%d' is not valid for union MuxedAccount", u.Type)
}

var _ decoderFrom = (*MuxedAccount)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *MuxedAccount) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding CryptoKeyType: %s", err)
	}
	switch CryptoKeyType(u.Type) {
	case CryptoKeyTypeKeyTypeEd25519:
		u.Ed25519 = new(Uint256)
		nTmp, err = (*u.Ed25519).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Uint256: %s", err)
		}
		return n, nil
	case CryptoKeyTypeKeyTypeMuxedEd25519:
		u.Med25519 = new(MuxedAccountMed25519)
		nTmp, err = (*u.Med25519).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding MuxedAccountMed25519: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union MuxedAccount has invalid Type (CryptoKeyType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s MuxedAccount) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *MuxedAccount) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*MuxedAccount)(nil)
	_ encoding.BinaryUnmarshaler = (*MuxedAccount)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s MuxedAccount) xdrType() {}

var _ xdrType = (*MuxedAccount)(nil)

// DecoratedSignature is an XDR Struct defines as:
//
//   struct DecoratedSignature
//    {
//        SignatureHint hint;  // last 4 bytes of the public key, used as a hint
//        Signature signature; // actual signature
//    };
//
type DecoratedSignature struct {
	Hint      SignatureHint
	Signature Signature
}

// EncodeTo encodes this value using the Encoder.
func (s *DecoratedSignature) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Hint.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Signature.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*DecoratedSignature)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *DecoratedSignature) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Hint.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SignatureHint: %s", err)
	}
	nTmp, err = s.Signature.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Signature: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s DecoratedSignature) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *DecoratedSignature) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*DecoratedSignature)(nil)
	_ encoding.BinaryUnmarshaler = (*DecoratedSignature)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s DecoratedSignature) xdrType() {}

var _ xdrType = (*DecoratedSignature)(nil)

// OperationType is an XDR Enum defines as:
//
//   enum OperationType
//    {
//        CREATE_ACCOUNT = 0,
//        PAYMENT = 1,
//        PATH_PAYMENT_STRICT_RECEIVE = 2,
//        MANAGE_SELL_OFFER = 3,
//        CREATE_PASSIVE_SELL_OFFER = 4,
//        SET_OPTIONS = 5,
//        CHANGE_TRUST = 6,
//        ALLOW_TRUST = 7,
//        ACCOUNT_MERGE = 8,
//        INFLATION = 9,
//        MANAGE_DATA = 10,
//        BUMP_SEQUENCE = 11,
//        MANAGE_BUY_OFFER = 12,
//        PATH_PAYMENT_STRICT_SEND = 13,
//        CREATE_CLAIMABLE_BALANCE = 14,
//        CLAIM_CLAIMABLE_BALANCE = 15,
//        BEGIN_SPONSORING_FUTURE_RESERVES = 16,
//        END_SPONSORING_FUTURE_RESERVES = 17,
//        REVOKE_SPONSORSHIP = 18,
//        CLAWBACK = 19,
//        CLAWBACK_CLAIMABLE_BALANCE = 20,
//        SET_TRUST_LINE_FLAGS = 21,
//        LIQUIDITY_POOL_DEPOSIT = 22,
//        LIQUIDITY_POOL_WITHDRAW = 23
//    };
//
type OperationType int32

const (
	OperationTypeCreateAccount                 OperationType = 0
	OperationTypePayment                       OperationType = 1
	OperationTypePathPaymentStrictReceive      OperationType = 2
	OperationTypeManageSellOffer               OperationType = 3
	OperationTypeCreatePassiveSellOffer        OperationType = 4
	OperationTypeSetOptions                    OperationType = 5
	OperationTypeChangeTrust                   OperationType = 6
	OperationTypeAllowTrust                    OperationType = 7
	OperationTypeAccountMerge                  OperationType = 8
	OperationTypeInflation                     OperationType = 9
	OperationTypeManageData                    OperationType = 10
	OperationTypeBumpSequence                  OperationType = 11
	OperationTypeManageBuyOffer                OperationType = 12
	OperationTypePathPaymentStrictSend         OperationType = 13
	OperationTypeCreateClaimableBalance        OperationType = 14
	OperationTypeClaimClaimableBalance         OperationType = 15
	OperationTypeBeginSponsoringFutureReserves OperationType = 16
	OperationTypeEndSponsoringFutureReserves   OperationType = 17
	OperationTypeRevokeSponsorship             OperationType = 18
	OperationTypeClawback                      OperationType = 19
	OperationTypeClawbackClaimableBalance      OperationType = 20
	OperationTypeSetTrustLineFlags             OperationType = 21
	OperationTypeLiquidityPoolDeposit          OperationType = 22
	OperationTypeLiquidityPoolWithdraw         OperationType = 23
)

var operationTypeMap = map[int32]string{
	0:  "OperationTypeCreateAccount",
	1:  "OperationTypePayment",
	2:  "OperationTypePathPaymentStrictReceive",
	3:  "OperationTypeManageSellOffer",
	4:  "OperationTypeCreatePassiveSellOffer",
	5:  "OperationTypeSetOptions",
	6:  "OperationTypeChangeTrust",
	7:  "OperationTypeAllowTrust",
	8:  "OperationTypeAccountMerge",
	9:  "OperationTypeInflation",
	10: "OperationTypeManageData",
	11: "OperationTypeBumpSequence",
	12: "OperationTypeManageBuyOffer",
	13: "OperationTypePathPaymentStrictSend",
	14: "OperationTypeCreateClaimableBalance",
	15: "OperationTypeClaimClaimableBalance",
	16: "OperationTypeBeginSponsoringFutureReserves",
	17: "OperationTypeEndSponsoringFutureReserves",
	18: "OperationTypeRevokeSponsorship",
	19: "OperationTypeClawback",
	20: "OperationTypeClawbackClaimableBalance",
	21: "OperationTypeSetTrustLineFlags",
	22: "OperationTypeLiquidityPoolDeposit",
	23: "OperationTypeLiquidityPoolWithdraw",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for OperationType
func (e OperationType) ValidEnum(v int32) bool {
	_, ok := operationTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e OperationType) String() string {
	name, _ := operationTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e OperationType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := operationTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid OperationType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*OperationType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *OperationType) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding OperationType: %s", err)
	}
	if _, ok := operationTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid OperationType enum value", v)
	}
	*e = OperationType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s OperationType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *OperationType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*OperationType)(nil)
	_ encoding.BinaryUnmarshaler = (*OperationType)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s OperationType) xdrType() {}

var _ xdrType = (*OperationType)(nil)

// CreateAccountOp is an XDR Struct defines as:
//
//   struct CreateAccountOp
//    {
//        AccountID destination; // account to create
//        int64 startingBalance; // amount they end up with
//    };
//
type CreateAccountOp struct {
	Destination     AccountId
	StartingBalance Int64
}

// EncodeTo encodes this value using the Encoder.
func (s *CreateAccountOp) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Destination.EncodeTo(e); err != nil {
		return err
	}
	if err = s.StartingBalance.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*CreateAccountOp)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *CreateAccountOp) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Destination.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountId: %s", err)
	}
	nTmp, err = s.StartingBalance.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s CreateAccountOp) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *CreateAccountOp) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*CreateAccountOp)(nil)
	_ encoding.BinaryUnmarshaler = (*CreateAccountOp)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s CreateAccountOp) xdrType() {}

var _ xdrType = (*CreateAccountOp)(nil)

// PaymentOp is an XDR Struct defines as:
//
//   struct PaymentOp
//    {
//        MuxedAccount destination; // recipient of the payment
//        Asset asset;              // what they end up with
//        int64 amount;             // amount they end up with
//    };
//
type PaymentOp struct {
	Destination MuxedAccount
	Asset       Asset
	Amount      Int64
}

// EncodeTo encodes this value using the Encoder.
func (s *PaymentOp) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Destination.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Asset.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Amount.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*PaymentOp)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *PaymentOp) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Destination.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding MuxedAccount: %s", err)
	}
	nTmp, err = s.Asset.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Asset: %s", err)
	}
	nTmp, err = s.Amount.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s PaymentOp) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *PaymentOp) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*PaymentOp)(nil)
	_ encoding.BinaryUnmarshaler = (*PaymentOp)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s PaymentOp) xdrType() {}

var _ xdrType = (*PaymentOp)(nil)

// PathPaymentStrictReceiveOp is an XDR Struct defines as:
//
//   struct PathPaymentStrictReceiveOp
//    {
//        Asset sendAsset; // asset we pay with
//        int64 sendMax;   // the maximum amount of sendAsset to
//                         // send (excluding fees).
//                         // The operation will fail if can't be met
//
//        MuxedAccount destination; // recipient of the payment
//        Asset destAsset;          // what they end up with
//        int64 destAmount;         // amount they end up with
//
//        Asset path<5>; // additional hops it must go through to get there
//    };
//
type PathPaymentStrictReceiveOp struct {
	SendAsset   Asset
	SendMax     Int64
	Destination MuxedAccount
	DestAsset   Asset
	DestAmount  Int64
	Path        []Asset `xdrmaxsize:"5"`
}

// EncodeTo encodes this value using the Encoder.
func (s *PathPaymentStrictReceiveOp) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.SendAsset.EncodeTo(e); err != nil {
		return err
	}
	if err = s.SendMax.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Destination.EncodeTo(e); err != nil {
		return err
	}
	if err = s.DestAsset.EncodeTo(e); err != nil {
		return err
	}
	if err = s.DestAmount.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeUint(uint32(len(s.Path))); err != nil {
		return err
	}
	for i := 0; i < len(s.Path); i++ {
		if err = s.Path[i].EncodeTo(e); err != nil {
			return err
		}
	}
	return nil
}

var _ decoderFrom = (*PathPaymentStrictReceiveOp)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *PathPaymentStrictReceiveOp) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.SendAsset.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Asset: %s", err)
	}
	nTmp, err = s.SendMax.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	nTmp, err = s.Destination.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding MuxedAccount: %s", err)
	}
	nTmp, err = s.DestAsset.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Asset: %s", err)
	}
	nTmp, err = s.DestAmount.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Asset: %s", err)
	}
	if l > 5 {
		return n, fmt.Errorf("decoding Asset: data size (%d) exceeds size limit (5)", l)
	}
	s.Path = nil
	if l > 0 {
		s.Path = make([]Asset, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.Path[i].DecodeFrom(d)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding Asset: %s", err)
			}
		}
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s PathPaymentStrictReceiveOp) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *PathPaymentStrictReceiveOp) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*PathPaymentStrictReceiveOp)(nil)
	_ encoding.BinaryUnmarshaler = (*PathPaymentStrictReceiveOp)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s PathPaymentStrictReceiveOp) xdrType() {}

var _ xdrType = (*PathPaymentStrictReceiveOp)(nil)

// PathPaymentStrictSendOp is an XDR Struct defines as:
//
//   struct PathPaymentStrictSendOp
//    {
//        Asset sendAsset;  // asset we pay with
//        int64 sendAmount; // amount of sendAsset to send (excluding fees)
//
//        MuxedAccount destination; // recipient of the payment
//        Asset destAsset;          // what they end up with
//        int64 destMin;            // the minimum amount of dest asset to
//                                  // be received
//                                  // The operation will fail if it can't be met
//
//        Asset path<5>; // additional hops it must go through to get there
//    };
//
type PathPaymentStrictSendOp struct {
	SendAsset   Asset
	SendAmount  Int64
	Destination MuxedAccount
	DestAsset   Asset
	DestMin     Int64
	Path        []Asset `xdrmaxsize:"5"`
}

// EncodeTo encodes this value using the Encoder.
func (s *PathPaymentStrictSendOp) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.SendAsset.EncodeTo(e); err != nil {
		return err
	}
	if err = s.SendAmount.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Destination.EncodeTo(e); err != nil {
		return err
	}
	if err = s.DestAsset.EncodeTo(e); err != nil {
		return err
	}
	if err = s.DestMin.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeUint(uint32(len(s.Path))); err != nil {
		return err
	}
	for i := 0; i < len(s.Path); i++ {
		if err = s.Path[i].EncodeTo(e); err != nil {
			return err
		}
	}
	return nil
}

var _ decoderFrom = (*PathPaymentStrictSendOp)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *PathPaymentStrictSendOp) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.SendAsset.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Asset: %s", err)
	}
	nTmp, err = s.SendAmount.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	nTmp, err = s.Destination.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding MuxedAccount: %s", err)
	}
	nTmp, err = s.DestAsset.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Asset: %s", err)
	}
	nTmp, err = s.DestMin.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Asset: %s", err)
	}
	if l > 5 {
		return n, fmt.Errorf("decoding Asset: data size (%d) exceeds size limit (5)", l)
	}
	s.Path = nil
	if l > 0 {
		s.Path = make([]Asset, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.Path[i].DecodeFrom(d)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding Asset: %s", err)
			}
		}
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s PathPaymentStrictSendOp) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *PathPaymentStrictSendOp) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*PathPaymentStrictSendOp)(nil)
	_ encoding.BinaryUnmarshaler = (*PathPaymentStrictSendOp)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s PathPaymentStrictSendOp) xdrType() {}

var _ xdrType = (*PathPaymentStrictSendOp)(nil)

// ManageSellOfferOp is an XDR Struct defines as:
//
//   struct ManageSellOfferOp
//    {
//        Asset selling;
//        Asset buying;
//        int64 amount; // amount being sold. if set to 0, delete the offer
//        Price price;  // price of thing being sold in terms of what you are buying
//
//        // 0=create a new offer, otherwise edit an existing offer
//        int64 offerID;
//    };
//
type ManageSellOfferOp struct {
	Selling Asset
	Buying  Asset
	Amount  Int64
	Price   Price
	OfferId Int64
}

// EncodeTo encodes this value using the Encoder.
func (s *ManageSellOfferOp) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Selling.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Buying.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Amount.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Price.EncodeTo(e); err != nil {
		return err
	}
	if err = s.OfferId.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*ManageSellOfferOp)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *ManageSellOfferOp) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Selling.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Asset: %s", err)
	}
	nTmp, err = s.Buying.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Asset: %s", err)
	}
	nTmp, err = s.Amount.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	nTmp, err = s.Price.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Price: %s", err)
	}
	nTmp, err = s.OfferId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ManageSellOfferOp) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ManageSellOfferOp) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ManageSellOfferOp)(nil)
	_ encoding.BinaryUnmarshaler = (*ManageSellOfferOp)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ManageSellOfferOp) xdrType() {}

var _ xdrType = (*ManageSellOfferOp)(nil)

// ManageBuyOfferOp is an XDR Struct defines as:
//
//   struct ManageBuyOfferOp
//    {
//        Asset selling;
//        Asset buying;
//        int64 buyAmount; // amount being bought. if set to 0, delete the offer
//        Price price;     // price of thing being bought in terms of what you are
//                         // selling
//
//        // 0=create a new offer, otherwise edit an existing offer
//        int64 offerID;
//    };
//
type ManageBuyOfferOp struct {
	Selling   Asset
	Buying    Asset
	BuyAmount Int64
	Price     Price
	OfferId   Int64
}

// EncodeTo encodes this value using the Encoder.
func (s *ManageBuyOfferOp) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Selling.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Buying.EncodeTo(e); err != nil {
		return err
	}
	if err = s.BuyAmount.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Price.EncodeTo(e); err != nil {
		return err
	}
	if err = s.OfferId.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*ManageBuyOfferOp)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *ManageBuyOfferOp) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Selling.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Asset: %s", err)
	}
	nTmp, err = s.Buying.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Asset: %s", err)
	}
	nTmp, err = s.BuyAmount.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	nTmp, err = s.Price.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Price: %s", err)
	}
	nTmp, err = s.OfferId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ManageBuyOfferOp) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ManageBuyOfferOp) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ManageBuyOfferOp)(nil)
	_ encoding.BinaryUnmarshaler = (*ManageBuyOfferOp)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ManageBuyOfferOp) xdrType() {}

var _ xdrType = (*ManageBuyOfferOp)(nil)

// CreatePassiveSellOfferOp is an XDR Struct defines as:
//
//   struct CreatePassiveSellOfferOp
//    {
//        Asset selling; // A
//        Asset buying;  // B
//        int64 amount;  // amount taker gets
//        Price price;   // cost of A in terms of B
//    };
//
type CreatePassiveSellOfferOp struct {
	Selling Asset
	Buying  Asset
	Amount  Int64
	Price   Price
}

// EncodeTo encodes this value using the Encoder.
func (s *CreatePassiveSellOfferOp) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Selling.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Buying.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Amount.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Price.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*CreatePassiveSellOfferOp)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *CreatePassiveSellOfferOp) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Selling.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Asset: %s", err)
	}
	nTmp, err = s.Buying.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Asset: %s", err)
	}
	nTmp, err = s.Amount.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	nTmp, err = s.Price.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Price: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s CreatePassiveSellOfferOp) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *CreatePassiveSellOfferOp) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*CreatePassiveSellOfferOp)(nil)
	_ encoding.BinaryUnmarshaler = (*CreatePassiveSellOfferOp)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s CreatePassiveSellOfferOp) xdrType() {}

var _ xdrType = (*CreatePassiveSellOfferOp)(nil)

// SetOptionsOp is an XDR Struct defines as:
//
//   struct SetOptionsOp
//    {
//        AccountID* inflationDest; // sets the inflation destination
//
//        uint32* clearFlags; // which flags to clear
//        uint32* setFlags;   // which flags to set
//
//        // account threshold manipulation
//        uint32* masterWeight; // weight of the master account
//        uint32* lowThreshold;
//        uint32* medThreshold;
//        uint32* highThreshold;
//
//        string32* homeDomain; // sets the home domain
//
//        // Add, update or remove a signer for the account
//        // signer is deleted if the weight is 0
//        Signer* signer;
//    };
//
type SetOptionsOp struct {
	InflationDest *AccountId
	ClearFlags    *Uint32
	SetFlags      *Uint32
	MasterWeight  *Uint32
	LowThreshold  *Uint32
	MedThreshold  *Uint32
	HighThreshold *Uint32
	HomeDomain    *String32
	Signer        *Signer
}

// EncodeTo encodes this value using the Encoder.
func (s *SetOptionsOp) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeBool(s.InflationDest != nil); err != nil {
		return err
	}
	if s.InflationDest != nil {
		if err = (*s.InflationDest).EncodeTo(e); err != nil {
			return err
		}
	}
	if _, err = e.EncodeBool(s.ClearFlags != nil); err != nil {
		return err
	}
	if s.ClearFlags != nil {
		if err = (*s.ClearFlags).EncodeTo(e); err != nil {
			return err
		}
	}
	if _, err = e.EncodeBool(s.SetFlags != nil); err != nil {
		return err
	}
	if s.SetFlags != nil {
		if err = (*s.SetFlags).EncodeTo(e); err != nil {
			return err
		}
	}
	if _, err = e.EncodeBool(s.MasterWeight != nil); err != nil {
		return err
	}
	if s.MasterWeight != nil {
		if err = (*s.MasterWeight).EncodeTo(e); err != nil {
			return err
		}
	}
	if _, err = e.EncodeBool(s.LowThreshold != nil); err != nil {
		return err
	}
	if s.LowThreshold != nil {
		if err = (*s.LowThreshold).EncodeTo(e); err != nil {
			return err
		}
	}
	if _, err = e.EncodeBool(s.MedThreshold != nil); err != nil {
		return err
	}
	if s.MedThreshold != nil {
		if err = (*s.MedThreshold).EncodeTo(e); err != nil {
			return err
		}
	}
	if _, err = e.EncodeBool(s.HighThreshold != nil); err != nil {
		return err
	}
	if s.HighThreshold != nil {
		if err = (*s.HighThreshold).EncodeTo(e); err != nil {
			return err
		}
	}
	if _, err = e.EncodeBool(s.HomeDomain != nil); err != nil {
		return err
	}
	if s.HomeDomain != nil {
		if err = (*s.HomeDomain).EncodeTo(e); err != nil {
			return err
		}
	}
	if _, err = e.EncodeBool(s.Signer != nil); err != nil {
		return err
	}
	if s.Signer != nil {
		if err = (*s.Signer).EncodeTo(e); err != nil {
			return err
		}
	}
	return nil
}

var _ decoderFrom = (*SetOptionsOp)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *SetOptionsOp) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	var b bool
	b, nTmp, err = d.DecodeBool()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountId: %s", err)
	}
	s.InflationDest = nil
	if b {
		s.InflationDest = new(AccountId)
		nTmp, err = s.InflationDest.DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding AccountId: %s", err)
		}
	}
	b, nTmp, err = d.DecodeBool()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	s.ClearFlags = nil
	if b {
		s.ClearFlags = new(Uint32)
		nTmp, err = s.ClearFlags.DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Uint32: %s", err)
		}
	}
	b, nTmp, err = d.DecodeBool()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	s.SetFlags = nil
	if b {
		s.SetFlags = new(Uint32)
		nTmp, err = s.SetFlags.DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Uint32: %s", err)
		}
	}
	b, nTmp, err = d.DecodeBool()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	s.MasterWeight = nil
	if b {
		s.MasterWeight = new(Uint32)
		nTmp, err = s.MasterWeight.DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Uint32: %s", err)
		}
	}
	b, nTmp, err = d.DecodeBool()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	s.LowThreshold = nil
	if b {
		s.LowThreshold = new(Uint32)
		nTmp, err = s.LowThreshold.DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Uint32: %s", err)
		}
	}
	b, nTmp, err = d.DecodeBool()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	s.MedThreshold = nil
	if b {
		s.MedThreshold = new(Uint32)
		nTmp, err = s.MedThreshold.DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Uint32: %s", err)
		}
	}
	b, nTmp, err = d.DecodeBool()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	s.HighThreshold = nil
	if b {
		s.HighThreshold = new(Uint32)
		nTmp, err = s.HighThreshold.DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Uint32: %s", err)
		}
	}
	b, nTmp, err = d.DecodeBool()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding String32: %s", err)
	}
	s.HomeDomain = nil
	if b {
		s.HomeDomain = new(String32)
		nTmp, err = s.HomeDomain.DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding String32: %s", err)
		}
	}
	b, nTmp, err = d.DecodeBool()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Signer: %s", err)
	}
	s.Signer = nil
	if b {
		s.Signer = new(Signer)
		nTmp, err = s.Signer.DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Signer: %s", err)
		}
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s SetOptionsOp) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SetOptionsOp) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*SetOptionsOp)(nil)
	_ encoding.BinaryUnmarshaler = (*SetOptionsOp)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s SetOptionsOp) xdrType() {}

var _ xdrType = (*SetOptionsOp)(nil)

// ChangeTrustAsset is an XDR Union defines as:
//
//   union ChangeTrustAsset switch (AssetType type)
//    {
//    case ASSET_TYPE_NATIVE: // Not credit
//        void;
//
//    case ASSET_TYPE_CREDIT_ALPHANUM4:
//        AlphaNum4 alphaNum4;
//
//    case ASSET_TYPE_CREDIT_ALPHANUM12:
//        AlphaNum12 alphaNum12;
//
//    case ASSET_TYPE_POOL_SHARE:
//        LiquidityPoolParameters liquidityPool;
//
//        // add other asset types here in the future
//    };
//
type ChangeTrustAsset struct {
	Type          AssetType
	AlphaNum4     *AlphaNum4
	AlphaNum12    *AlphaNum12
	LiquidityPool *LiquidityPoolParameters
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ChangeTrustAsset) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ChangeTrustAsset
func (u ChangeTrustAsset) ArmForSwitch(sw int32) (string, bool) {
	switch AssetType(sw) {
	case AssetTypeAssetTypeNative:
		return "", true
	case AssetTypeAssetTypeCreditAlphanum4:
		return "AlphaNum4", true
	case AssetTypeAssetTypeCreditAlphanum12:
		return "AlphaNum12", true
	case AssetTypeAssetTypePoolShare:
		return "LiquidityPool", true
	}
	return "-", false
}

// NewChangeTrustAsset creates a new  ChangeTrustAsset.
func NewChangeTrustAsset(aType AssetType, value interface{}) (result ChangeTrustAsset, err error) {
	result.Type = aType
	switch AssetType(aType) {
	case AssetTypeAssetTypeNative:
		// void
	case AssetTypeAssetTypeCreditAlphanum4:
		tv, ok := value.(AlphaNum4)
		if !ok {
			err = fmt.Errorf("invalid value, must be AlphaNum4")
			return
		}
		result.AlphaNum4 = &tv
	case AssetTypeAssetTypeCreditAlphanum12:
		tv, ok := value.(AlphaNum12)
		if !ok {
			err = fmt.Errorf("invalid value, must be AlphaNum12")
			return
		}
		result.AlphaNum12 = &tv
	case AssetTypeAssetTypePoolShare:
		tv, ok := value.(LiquidityPoolParameters)
		if !ok {
			err = fmt.Errorf("invalid value, must be LiquidityPoolParameters")
			return
		}
		result.LiquidityPool = &tv
	}
	return
}

// MustAlphaNum4 retrieves the AlphaNum4 value from the union,
// panicing if the value is not set.
func (u ChangeTrustAsset) MustAlphaNum4() AlphaNum4 {
	val, ok := u.GetAlphaNum4()

	if !ok {
		panic("arm AlphaNum4 is not set")
	}

	return val
}

// GetAlphaNum4 retrieves the AlphaNum4 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ChangeTrustAsset) GetAlphaNum4() (result AlphaNum4, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AlphaNum4" {
		result = *u.AlphaNum4
		ok = true
	}

	return
}

// MustAlphaNum12 retrieves the AlphaNum12 value from the union,
// panicing if the value is not set.
func (u ChangeTrustAsset) MustAlphaNum12() AlphaNum12 {
	val, ok := u.GetAlphaNum12()

	if !ok {
		panic("arm AlphaNum12 is not set")
	}

	return val
}

// GetAlphaNum12 retrieves the AlphaNum12 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ChangeTrustAsset) GetAlphaNum12() (result AlphaNum12, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AlphaNum12" {
		result = *u.AlphaNum12
		ok = true
	}

	return
}

// MustLiquidityPool retrieves the LiquidityPool value from the union,
// panicing if the value is not set.
func (u ChangeTrustAsset) MustLiquidityPool() LiquidityPoolParameters {
	val, ok := u.GetLiquidityPool()

	if !ok {
		panic("arm LiquidityPool is not set")
	}

	return val
}

// GetLiquidityPool retrieves the LiquidityPool value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ChangeTrustAsset) GetLiquidityPool() (result LiquidityPoolParameters, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "LiquidityPool" {
		result = *u.LiquidityPool
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u ChangeTrustAsset) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch AssetType(u.Type) {
	case AssetTypeAssetTypeNative:
		// Void
		return nil
	case AssetTypeAssetTypeCreditAlphanum4:
		if err = (*u.AlphaNum4).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case AssetTypeAssetTypeCreditAlphanum12:
		if err = (*u.AlphaNum12).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case AssetTypeAssetTypePoolShare:
		if err = (*u.LiquidityPool).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (AssetType) switch value '%d' is not valid for union ChangeTrustAsset", u.Type)
}

var _ decoderFrom = (*ChangeTrustAsset)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *ChangeTrustAsset) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AssetType: %s", err)
	}
	switch AssetType(u.Type) {
	case AssetTypeAssetTypeNative:
		// Void
		return n, nil
	case AssetTypeAssetTypeCreditAlphanum4:
		u.AlphaNum4 = new(AlphaNum4)
		nTmp, err = (*u.AlphaNum4).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding AlphaNum4: %s", err)
		}
		return n, nil
	case AssetTypeAssetTypeCreditAlphanum12:
		u.AlphaNum12 = new(AlphaNum12)
		nTmp, err = (*u.AlphaNum12).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding AlphaNum12: %s", err)
		}
		return n, nil
	case AssetTypeAssetTypePoolShare:
		u.LiquidityPool = new(LiquidityPoolParameters)
		nTmp, err = (*u.LiquidityPool).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding LiquidityPoolParameters: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union ChangeTrustAsset has invalid Type (AssetType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ChangeTrustAsset) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ChangeTrustAsset) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ChangeTrustAsset)(nil)
	_ encoding.BinaryUnmarshaler = (*ChangeTrustAsset)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ChangeTrustAsset) xdrType() {}

var _ xdrType = (*ChangeTrustAsset)(nil)

// ChangeTrustOp is an XDR Struct defines as:
//
//   struct ChangeTrustOp
//    {
//        ChangeTrustAsset line;
//
//        // if limit is set to 0, deletes the trust line
//        int64 limit;
//    };
//
type ChangeTrustOp struct {
	Line  ChangeTrustAsset
	Limit Int64
}

// EncodeTo encodes this value using the Encoder.
func (s *ChangeTrustOp) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Line.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Limit.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*ChangeTrustOp)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *ChangeTrustOp) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Line.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ChangeTrustAsset: %s", err)
	}
	nTmp, err = s.Limit.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ChangeTrustOp) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ChangeTrustOp) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ChangeTrustOp)(nil)
	_ encoding.BinaryUnmarshaler = (*ChangeTrustOp)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ChangeTrustOp) xdrType() {}

var _ xdrType = (*ChangeTrustOp)(nil)

// AllowTrustOp is an XDR Struct defines as:
//
//   struct AllowTrustOp
//    {
//        AccountID trustor;
//        AssetCode asset;
//
//        // One of 0, AUTHORIZED_FLAG, or AUTHORIZED_TO_MAINTAIN_LIABILITIES_FLAG
//        uint32 authorize;
//    };
//
type AllowTrustOp struct {
	Trustor   AccountId
	Asset     AssetCode
	Authorize Uint32
}

// EncodeTo encodes this value using the Encoder.
func (s *AllowTrustOp) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Trustor.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Asset.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Authorize.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*AllowTrustOp)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *AllowTrustOp) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Trustor.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountId: %s", err)
	}
	nTmp, err = s.Asset.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AssetCode: %s", err)
	}
	nTmp, err = s.Authorize.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s AllowTrustOp) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AllowTrustOp) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*AllowTrustOp)(nil)
	_ encoding.BinaryUnmarshaler = (*AllowTrustOp)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s AllowTrustOp) xdrType() {}

var _ xdrType = (*AllowTrustOp)(nil)

// ManageDataOp is an XDR Struct defines as:
//
//   struct ManageDataOp
//    {
//        string64 dataName;
//        DataValue* dataValue; // set to null to clear
//    };
//
type ManageDataOp struct {
	DataName  String64
	DataValue *DataValue
}

// EncodeTo encodes this value using the Encoder.
func (s *ManageDataOp) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.DataName.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeBool(s.DataValue != nil); err != nil {
		return err
	}
	if s.DataValue != nil {
		if err = (*s.DataValue).EncodeTo(e); err != nil {
			return err
		}
	}
	return nil
}

var _ decoderFrom = (*ManageDataOp)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *ManageDataOp) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.DataName.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding String64: %s", err)
	}
	var b bool
	b, nTmp, err = d.DecodeBool()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding DataValue: %s", err)
	}
	s.DataValue = nil
	if b {
		s.DataValue = new(DataValue)
		nTmp, err = s.DataValue.DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding DataValue: %s", err)
		}
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ManageDataOp) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ManageDataOp) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ManageDataOp)(nil)
	_ encoding.BinaryUnmarshaler = (*ManageDataOp)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ManageDataOp) xdrType() {}

var _ xdrType = (*ManageDataOp)(nil)

// BumpSequenceOp is an XDR Struct defines as:
//
//   struct BumpSequenceOp
//    {
//        SequenceNumber bumpTo;
//    };
//
type BumpSequenceOp struct {
	BumpTo SequenceNumber
}

// EncodeTo encodes this value using the Encoder.
func (s *BumpSequenceOp) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.BumpTo.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*BumpSequenceOp)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *BumpSequenceOp) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.BumpTo.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SequenceNumber: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s BumpSequenceOp) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *BumpSequenceOp) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*BumpSequenceOp)(nil)
	_ encoding.BinaryUnmarshaler = (*BumpSequenceOp)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s BumpSequenceOp) xdrType() {}

var _ xdrType = (*BumpSequenceOp)(nil)

// CreateClaimableBalanceOp is an XDR Struct defines as:
//
//   struct CreateClaimableBalanceOp
//    {
//        Asset asset;
//        int64 amount;
//        Claimant claimants<10>;
//    };
//
type CreateClaimableBalanceOp struct {
	Asset     Asset
	Amount    Int64
	Claimants []Claimant `xdrmaxsize:"10"`
}

// EncodeTo encodes this value using the Encoder.
func (s *CreateClaimableBalanceOp) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Asset.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Amount.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeUint(uint32(len(s.Claimants))); err != nil {
		return err
	}
	for i := 0; i < len(s.Claimants); i++ {
		if err = s.Claimants[i].EncodeTo(e); err != nil {
			return err
		}
	}
	return nil
}

var _ decoderFrom = (*CreateClaimableBalanceOp)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *CreateClaimableBalanceOp) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Asset.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Asset: %s", err)
	}
	nTmp, err = s.Amount.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Claimant: %s", err)
	}
	if l > 10 {
		return n, fmt.Errorf("decoding Claimant: data size (%d) exceeds size limit (10)", l)
	}
	s.Claimants = nil
	if l > 0 {
		s.Claimants = make([]Claimant, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.Claimants[i].DecodeFrom(d)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding Claimant: %s", err)
			}
		}
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s CreateClaimableBalanceOp) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *CreateClaimableBalanceOp) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*CreateClaimableBalanceOp)(nil)
	_ encoding.BinaryUnmarshaler = (*CreateClaimableBalanceOp)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s CreateClaimableBalanceOp) xdrType() {}

var _ xdrType = (*CreateClaimableBalanceOp)(nil)

// ClaimClaimableBalanceOp is an XDR Struct defines as:
//
//   struct ClaimClaimableBalanceOp
//    {
//        ClaimableBalanceID balanceID;
//    };
//
type ClaimClaimableBalanceOp struct {
	BalanceId ClaimableBalanceId
}

// EncodeTo encodes this value using the Encoder.
func (s *ClaimClaimableBalanceOp) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.BalanceId.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*ClaimClaimableBalanceOp)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *ClaimClaimableBalanceOp) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.BalanceId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ClaimableBalanceId: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ClaimClaimableBalanceOp) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ClaimClaimableBalanceOp) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ClaimClaimableBalanceOp)(nil)
	_ encoding.BinaryUnmarshaler = (*ClaimClaimableBalanceOp)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ClaimClaimableBalanceOp) xdrType() {}

var _ xdrType = (*ClaimClaimableBalanceOp)(nil)

// BeginSponsoringFutureReservesOp is an XDR Struct defines as:
//
//   struct BeginSponsoringFutureReservesOp
//    {
//        AccountID sponsoredID;
//    };
//
type BeginSponsoringFutureReservesOp struct {
	SponsoredId AccountId
}

// EncodeTo encodes this value using the Encoder.
func (s *BeginSponsoringFutureReservesOp) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.SponsoredId.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*BeginSponsoringFutureReservesOp)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *BeginSponsoringFutureReservesOp) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.SponsoredId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountId: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s BeginSponsoringFutureReservesOp) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *BeginSponsoringFutureReservesOp) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*BeginSponsoringFutureReservesOp)(nil)
	_ encoding.BinaryUnmarshaler = (*BeginSponsoringFutureReservesOp)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s BeginSponsoringFutureReservesOp) xdrType() {}

var _ xdrType = (*BeginSponsoringFutureReservesOp)(nil)

// RevokeSponsorshipType is an XDR Enum defines as:
//
//   enum RevokeSponsorshipType
//    {
//        REVOKE_SPONSORSHIP_LEDGER_ENTRY = 0,
//        REVOKE_SPONSORSHIP_SIGNER = 1
//    };
//
type RevokeSponsorshipType int32

const (
	RevokeSponsorshipTypeRevokeSponsorshipLedgerEntry RevokeSponsorshipType = 0
	RevokeSponsorshipTypeRevokeSponsorshipSigner      RevokeSponsorshipType = 1
)

var revokeSponsorshipTypeMap = map[int32]string{
	0: "RevokeSponsorshipTypeRevokeSponsorshipLedgerEntry",
	1: "RevokeSponsorshipTypeRevokeSponsorshipSigner",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for RevokeSponsorshipType
func (e RevokeSponsorshipType) ValidEnum(v int32) bool {
	_, ok := revokeSponsorshipTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e RevokeSponsorshipType) String() string {
	name, _ := revokeSponsorshipTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e RevokeSponsorshipType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := revokeSponsorshipTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid RevokeSponsorshipType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*RevokeSponsorshipType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *RevokeSponsorshipType) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding RevokeSponsorshipType: %s", err)
	}
	if _, ok := revokeSponsorshipTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid RevokeSponsorshipType enum value", v)
	}
	*e = RevokeSponsorshipType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s RevokeSponsorshipType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *RevokeSponsorshipType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*RevokeSponsorshipType)(nil)
	_ encoding.BinaryUnmarshaler = (*RevokeSponsorshipType)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s RevokeSponsorshipType) xdrType() {}

var _ xdrType = (*RevokeSponsorshipType)(nil)

// RevokeSponsorshipOpSigner is an XDR NestedStruct defines as:
//
//   struct
//        {
//            AccountID accountID;
//            SignerKey signerKey;
//        }
//
type RevokeSponsorshipOpSigner struct {
	AccountId AccountId
	SignerKey SignerKey
}

// EncodeTo encodes this value using the Encoder.
func (s *RevokeSponsorshipOpSigner) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.AccountId.EncodeTo(e); err != nil {
		return err
	}
	if err = s.SignerKey.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*RevokeSponsorshipOpSigner)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *RevokeSponsorshipOpSigner) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.AccountId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountId: %s", err)
	}
	nTmp, err = s.SignerKey.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SignerKey: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s RevokeSponsorshipOpSigner) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *RevokeSponsorshipOpSigner) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*RevokeSponsorshipOpSigner)(nil)
	_ encoding.BinaryUnmarshaler = (*RevokeSponsorshipOpSigner)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s RevokeSponsorshipOpSigner) xdrType() {}

var _ xdrType = (*RevokeSponsorshipOpSigner)(nil)

// RevokeSponsorshipOp is an XDR Union defines as:
//
//   union RevokeSponsorshipOp switch (RevokeSponsorshipType type)
//    {
//    case REVOKE_SPONSORSHIP_LEDGER_ENTRY:
//        LedgerKey ledgerKey;
//    case REVOKE_SPONSORSHIP_SIGNER:
//        struct
//        {
//            AccountID accountID;
//            SignerKey signerKey;
//        } signer;
//    };
//
type RevokeSponsorshipOp struct {
	Type      RevokeSponsorshipType
	LedgerKey *LedgerKey
	Signer    *RevokeSponsorshipOpSigner
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u RevokeSponsorshipOp) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of RevokeSponsorshipOp
func (u RevokeSponsorshipOp) ArmForSwitch(sw int32) (string, bool) {
	switch RevokeSponsorshipType(sw) {
	case RevokeSponsorshipTypeRevokeSponsorshipLedgerEntry:
		return "LedgerKey", true
	case RevokeSponsorshipTypeRevokeSponsorshipSigner:
		return "Signer", true
	}
	return "-", false
}

// NewRevokeSponsorshipOp creates a new  RevokeSponsorshipOp.
func NewRevokeSponsorshipOp(aType RevokeSponsorshipType, value interface{}) (result RevokeSponsorshipOp, err error) {
	result.Type = aType
	switch RevokeSponsorshipType(aType) {
	case RevokeSponsorshipTypeRevokeSponsorshipLedgerEntry:
		tv, ok := value.(LedgerKey)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKey")
			return
		}
		result.LedgerKey = &tv
	case RevokeSponsorshipTypeRevokeSponsorshipSigner:
		tv, ok := value.(RevokeSponsorshipOpSigner)
		if !ok {
			err = fmt.Errorf("invalid value, must be RevokeSponsorshipOpSigner")
			return
		}
		result.Signer = &tv
	}
	return
}

// MustLedgerKey retrieves the LedgerKey value from the union,
// panicing if the value is not set.
func (u RevokeSponsorshipOp) MustLedgerKey() LedgerKey {
	val, ok := u.GetLedgerKey()

	if !ok {
		panic("arm LedgerKey is not set")
	}

	return val
}

// GetLedgerKey retrieves the LedgerKey value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u RevokeSponsorshipOp) GetLedgerKey() (result LedgerKey, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "LedgerKey" {
		result = *u.LedgerKey
		ok = true
	}

	return
}

// MustSigner retrieves the Signer value from the union,
// panicing if the value is not set.
func (u RevokeSponsorshipOp) MustSigner() RevokeSponsorshipOpSigner {
	val, ok := u.GetSigner()

	if !ok {
		panic("arm Signer is not set")
	}

	return val
}

// GetSigner retrieves the Signer value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u RevokeSponsorshipOp) GetSigner() (result RevokeSponsorshipOpSigner, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Signer" {
		result = *u.Signer
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u RevokeSponsorshipOp) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch RevokeSponsorshipType(u.Type) {
	case RevokeSponsorshipTypeRevokeSponsorshipLedgerEntry:
		if err = (*u.LedgerKey).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case RevokeSponsorshipTypeRevokeSponsorshipSigner:
		if err = (*u.Signer).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (RevokeSponsorshipType) switch value '%d' is not valid for union RevokeSponsorshipOp", u.Type)
}

var _ decoderFrom = (*RevokeSponsorshipOp)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *RevokeSponsorshipOp) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding RevokeSponsorshipType: %s", err)
	}
	switch RevokeSponsorshipType(u.Type) {
	case RevokeSponsorshipTypeRevokeSponsorshipLedgerEntry:
		u.LedgerKey = new(LedgerKey)
		nTmp, err = (*u.LedgerKey).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding LedgerKey: %s", err)
		}
		return n, nil
	case RevokeSponsorshipTypeRevokeSponsorshipSigner:
		u.Signer = new(RevokeSponsorshipOpSigner)
		nTmp, err = (*u.Signer).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding RevokeSponsorshipOpSigner: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union RevokeSponsorshipOp has invalid Type (RevokeSponsorshipType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s RevokeSponsorshipOp) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *RevokeSponsorshipOp) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*RevokeSponsorshipOp)(nil)
	_ encoding.BinaryUnmarshaler = (*RevokeSponsorshipOp)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s RevokeSponsorshipOp) xdrType() {}

var _ xdrType = (*RevokeSponsorshipOp)(nil)

// ClawbackOp is an XDR Struct defines as:
//
//   struct ClawbackOp
//    {
//        Asset asset;
//        MuxedAccount from;
//        int64 amount;
//    };
//
type ClawbackOp struct {
	Asset  Asset
	From   MuxedAccount
	Amount Int64
}

// EncodeTo encodes this value using the Encoder.
func (s *ClawbackOp) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Asset.EncodeTo(e); err != nil {
		return err
	}
	if err = s.From.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Amount.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*ClawbackOp)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *ClawbackOp) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Asset.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Asset: %s", err)
	}
	nTmp, err = s.From.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding MuxedAccount: %s", err)
	}
	nTmp, err = s.Amount.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ClawbackOp) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ClawbackOp) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ClawbackOp)(nil)
	_ encoding.BinaryUnmarshaler = (*ClawbackOp)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ClawbackOp) xdrType() {}

var _ xdrType = (*ClawbackOp)(nil)

// ClawbackClaimableBalanceOp is an XDR Struct defines as:
//
//   struct ClawbackClaimableBalanceOp
//    {
//        ClaimableBalanceID balanceID;
//    };
//
type ClawbackClaimableBalanceOp struct {
	BalanceId ClaimableBalanceId
}

// EncodeTo encodes this value using the Encoder.
func (s *ClawbackClaimableBalanceOp) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.BalanceId.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*ClawbackClaimableBalanceOp)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *ClawbackClaimableBalanceOp) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.BalanceId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ClaimableBalanceId: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ClawbackClaimableBalanceOp) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ClawbackClaimableBalanceOp) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ClawbackClaimableBalanceOp)(nil)
	_ encoding.BinaryUnmarshaler = (*ClawbackClaimableBalanceOp)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ClawbackClaimableBalanceOp) xdrType() {}

var _ xdrType = (*ClawbackClaimableBalanceOp)(nil)

// SetTrustLineFlagsOp is an XDR Struct defines as:
//
//   struct SetTrustLineFlagsOp
//    {
//        AccountID trustor;
//        Asset asset;
//
//        uint32 clearFlags; // which flags to clear
//        uint32 setFlags;   // which flags to set
//    };
//
type SetTrustLineFlagsOp struct {
	Trustor    AccountId
	Asset      Asset
	ClearFlags Uint32
	SetFlags   Uint32
}

// EncodeTo encodes this value using the Encoder.
func (s *SetTrustLineFlagsOp) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Trustor.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Asset.EncodeTo(e); err != nil {
		return err
	}
	if err = s.ClearFlags.EncodeTo(e); err != nil {
		return err
	}
	if err = s.SetFlags.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*SetTrustLineFlagsOp)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *SetTrustLineFlagsOp) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Trustor.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountId: %s", err)
	}
	nTmp, err = s.Asset.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Asset: %s", err)
	}
	nTmp, err = s.ClearFlags.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.SetFlags.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s SetTrustLineFlagsOp) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SetTrustLineFlagsOp) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*SetTrustLineFlagsOp)(nil)
	_ encoding.BinaryUnmarshaler = (*SetTrustLineFlagsOp)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s SetTrustLineFlagsOp) xdrType() {}

var _ xdrType = (*SetTrustLineFlagsOp)(nil)

// LiquidityPoolFeeV18 is an XDR Const defines as:
//
//   const LIQUIDITY_POOL_FEE_V18 = 30;
//
const LiquidityPoolFeeV18 = 30

// LiquidityPoolDepositOp is an XDR Struct defines as:
//
//   struct LiquidityPoolDepositOp
//    {
//        PoolID liquidityPoolID;
//        int64 maxAmountA;     // maximum amount of first asset to deposit
//        int64 maxAmountB;     // maximum amount of second asset to deposit
//        Price minPrice;       // minimum depositA/depositB
//        Price maxPrice;       // maximum depositA/depositB
//    };
//
type LiquidityPoolDepositOp struct {
	LiquidityPoolId PoolId
	MaxAmountA      Int64
	MaxAmountB      Int64
	MinPrice        Price
	MaxPrice        Price
}

// EncodeTo encodes this value using the Encoder.
func (s *LiquidityPoolDepositOp) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.LiquidityPoolId.EncodeTo(e); err != nil {
		return err
	}
	if err = s.MaxAmountA.EncodeTo(e); err != nil {
		return err
	}
	if err = s.MaxAmountB.EncodeTo(e); err != nil {
		return err
	}
	if err = s.MinPrice.EncodeTo(e); err != nil {
		return err
	}
	if err = s.MaxPrice.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*LiquidityPoolDepositOp)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *LiquidityPoolDepositOp) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.LiquidityPoolId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding PoolId: %s", err)
	}
	nTmp, err = s.MaxAmountA.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	nTmp, err = s.MaxAmountB.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	nTmp, err = s.MinPrice.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Price: %s", err)
	}
	nTmp, err = s.MaxPrice.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Price: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LiquidityPoolDepositOp) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LiquidityPoolDepositOp) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LiquidityPoolDepositOp)(nil)
	_ encoding.BinaryUnmarshaler = (*LiquidityPoolDepositOp)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LiquidityPoolDepositOp) xdrType() {}

var _ xdrType = (*LiquidityPoolDepositOp)(nil)

// LiquidityPoolWithdrawOp is an XDR Struct defines as:
//
//   struct LiquidityPoolWithdrawOp
//    {
//        PoolID liquidityPoolID;
//        int64 amount;         // amount of pool shares to withdraw
//        int64 minAmountA;     // minimum amount of first asset to withdraw
//        int64 minAmountB;     // minimum amount of second asset to withdraw
//    };
//
type LiquidityPoolWithdrawOp struct {
	LiquidityPoolId PoolId
	Amount          Int64
	MinAmountA      Int64
	MinAmountB      Int64
}

// EncodeTo encodes this value using the Encoder.
func (s *LiquidityPoolWithdrawOp) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.LiquidityPoolId.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Amount.EncodeTo(e); err != nil {
		return err
	}
	if err = s.MinAmountA.EncodeTo(e); err != nil {
		return err
	}
	if err = s.MinAmountB.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*LiquidityPoolWithdrawOp)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *LiquidityPoolWithdrawOp) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.LiquidityPoolId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding PoolId: %s", err)
	}
	nTmp, err = s.Amount.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	nTmp, err = s.MinAmountA.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	nTmp, err = s.MinAmountB.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LiquidityPoolWithdrawOp) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LiquidityPoolWithdrawOp) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LiquidityPoolWithdrawOp)(nil)
	_ encoding.BinaryUnmarshaler = (*LiquidityPoolWithdrawOp)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LiquidityPoolWithdrawOp) xdrType() {}

var _ xdrType = (*LiquidityPoolWithdrawOp)(nil)

// OperationBody is an XDR NestedUnion defines as:
//
//   union switch (OperationType type)
//        {
//        case CREATE_ACCOUNT:
//            CreateAccountOp createAccountOp;
//        case PAYMENT:
//            PaymentOp paymentOp;
//        case PATH_PAYMENT_STRICT_RECEIVE:
//            PathPaymentStrictReceiveOp pathPaymentStrictReceiveOp;
//        case MANAGE_SELL_OFFER:
//            ManageSellOfferOp manageSellOfferOp;
//        case CREATE_PASSIVE_SELL_OFFER:
//            CreatePassiveSellOfferOp createPassiveSellOfferOp;
//        case SET_OPTIONS:
//            SetOptionsOp setOptionsOp;
//        case CHANGE_TRUST:
//            ChangeTrustOp changeTrustOp;
//        case ALLOW_TRUST:
//            AllowTrustOp allowTrustOp;
//        case ACCOUNT_MERGE:
//            MuxedAccount destination;
//        case INFLATION:
//            void;
//        case MANAGE_DATA:
//            ManageDataOp manageDataOp;
//        case BUMP_SEQUENCE:
//            BumpSequenceOp bumpSequenceOp;
//        case MANAGE_BUY_OFFER:
//            ManageBuyOfferOp manageBuyOfferOp;
//        case PATH_PAYMENT_STRICT_SEND:
//            PathPaymentStrictSendOp pathPaymentStrictSendOp;
//        case CREATE_CLAIMABLE_BALANCE:
//            CreateClaimableBalanceOp createClaimableBalanceOp;
//        case CLAIM_CLAIMABLE_BALANCE:
//            ClaimClaimableBalanceOp claimClaimableBalanceOp;
//        case BEGIN_SPONSORING_FUTURE_RESERVES:
//            BeginSponsoringFutureReservesOp beginSponsoringFutureReservesOp;
//        case END_SPONSORING_FUTURE_RESERVES:
//            void;
//        case REVOKE_SPONSORSHIP:
//            RevokeSponsorshipOp revokeSponsorshipOp;
//        case CLAWBACK:
//            ClawbackOp clawbackOp;
//        case CLAWBACK_CLAIMABLE_BALANCE:
//            ClawbackClaimableBalanceOp clawbackClaimableBalanceOp;
//        case SET_TRUST_LINE_FLAGS:
//            SetTrustLineFlagsOp setTrustLineFlagsOp;
//        case LIQUIDITY_POOL_DEPOSIT:
//            LiquidityPoolDepositOp liquidityPoolDepositOp;
//        case LIQUIDITY_POOL_WITHDRAW:
//            LiquidityPoolWithdrawOp liquidityPoolWithdrawOp;
//        }
//
type OperationBody struct {
	Type                            OperationType
	CreateAccountOp                 *CreateAccountOp
	PaymentOp                       *PaymentOp
	PathPaymentStrictReceiveOp      *PathPaymentStrictReceiveOp
	ManageSellOfferOp               *ManageSellOfferOp
	CreatePassiveSellOfferOp        *CreatePassiveSellOfferOp
	SetOptionsOp                    *SetOptionsOp
	ChangeTrustOp                   *ChangeTrustOp
	AllowTrustOp                    *AllowTrustOp
	Destination                     *MuxedAccount
	ManageDataOp                    *ManageDataOp
	BumpSequenceOp                  *BumpSequenceOp
	ManageBuyOfferOp                *ManageBuyOfferOp
	PathPaymentStrictSendOp         *PathPaymentStrictSendOp
	CreateClaimableBalanceOp        *CreateClaimableBalanceOp
	ClaimClaimableBalanceOp         *ClaimClaimableBalanceOp
	BeginSponsoringFutureReservesOp *BeginSponsoringFutureReservesOp
	RevokeSponsorshipOp             *RevokeSponsorshipOp
	ClawbackOp                      *ClawbackOp
	ClawbackClaimableBalanceOp      *ClawbackClaimableBalanceOp
	SetTrustLineFlagsOp             *SetTrustLineFlagsOp
	LiquidityPoolDepositOp          *LiquidityPoolDepositOp
	LiquidityPoolWithdrawOp         *LiquidityPoolWithdrawOp
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u OperationBody) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of OperationBody
func (u OperationBody) ArmForSwitch(sw int32) (string, bool) {
	switch OperationType(sw) {
	case OperationTypeCreateAccount:
		return "CreateAccountOp", true
	case OperationTypePayment:
		return "PaymentOp", true
	case OperationTypePathPaymentStrictReceive:
		return "PathPaymentStrictReceiveOp", true
	case OperationTypeManageSellOffer:
		return "ManageSellOfferOp", true
	case OperationTypeCreatePassiveSellOffer:
		return "CreatePassiveSellOfferOp", true
	case OperationTypeSetOptions:
		return "SetOptionsOp", true
	case OperationTypeChangeTrust:
		return "ChangeTrustOp", true
	case OperationTypeAllowTrust:
		return "AllowTrustOp", true
	case OperationTypeAccountMerge:
		return "Destination", true
	case OperationTypeInflation:
		return "", true
	case OperationTypeManageData:
		return "ManageDataOp", true
	case OperationTypeBumpSequence:
		return "BumpSequenceOp", true
	case OperationTypeManageBuyOffer:
		return "ManageBuyOfferOp", true
	case OperationTypePathPaymentStrictSend:
		return "PathPaymentStrictSendOp", true
	case OperationTypeCreateClaimableBalance:
		return "CreateClaimableBalanceOp", true
	case OperationTypeClaimClaimableBalance:
		return "ClaimClaimableBalanceOp", true
	case OperationTypeBeginSponsoringFutureReserves:
		return "BeginSponsoringFutureReservesOp", true
	case OperationTypeEndSponsoringFutureReserves:
		return "", true
	case OperationTypeRevokeSponsorship:
		return "RevokeSponsorshipOp", true
	case OperationTypeClawback:
		return "ClawbackOp", true
	case OperationTypeClawbackClaimableBalance:
		return "ClawbackClaimableBalanceOp", true
	case OperationTypeSetTrustLineFlags:
		return "SetTrustLineFlagsOp", true
	case OperationTypeLiquidityPoolDeposit:
		return "LiquidityPoolDepositOp", true
	case OperationTypeLiquidityPoolWithdraw:
		return "LiquidityPoolWithdrawOp", true
	}
	return "-", false
}

// NewOperationBody creates a new  OperationBody.
func NewOperationBody(aType OperationType, value interface{}) (result OperationBody, err error) {
	result.Type = aType
	switch OperationType(aType) {
	case OperationTypeCreateAccount:
		tv, ok := value.(CreateAccountOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be CreateAccountOp")
			return
		}
		result.CreateAccountOp = &tv
	case OperationTypePayment:
		tv, ok := value.(PaymentOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be PaymentOp")
			return
		}
		result.PaymentOp = &tv
	case OperationTypePathPaymentStrictReceive:
		tv, ok := value.(PathPaymentStrictReceiveOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be PathPaymentStrictReceiveOp")
			return
		}
		result.PathPaymentStrictReceiveOp = &tv
	case OperationTypeManageSellOffer:
		tv, ok := value.(ManageSellOfferOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageSellOfferOp")
			return
		}
		result.ManageSellOfferOp = &tv
	case OperationTypeCreatePassiveSellOffer:
		tv, ok := value.(CreatePassiveSellOfferOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be CreatePassiveSellOfferOp")
			return
		}
		result.CreatePassiveSellOfferOp = &tv
	case OperationTypeSetOptions:
		tv, ok := value.(SetOptionsOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be SetOptionsOp")
			return
		}
		result.SetOptionsOp = &tv
	case OperationTypeChangeTrust:
		tv, ok := value.(ChangeTrustOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be ChangeTrustOp")
			return
		}
		result.ChangeTrustOp = &tv
	case OperationTypeAllowTrust:
		tv, ok := value.(AllowTrustOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be AllowTrustOp")
			return
		}
		result.AllowTrustOp = &tv
	case OperationTypeAccountMerge:
		tv, ok := value.(MuxedAccount)
		if !ok {
			err = fmt.Errorf("invalid value, must be MuxedAccount")
			return
		}
		result.Destination = &tv
	case OperationTypeInflation:
		// void
	case OperationTypeManageData:
		tv, ok := value.(ManageDataOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageDataOp")
			return
		}
		result.ManageDataOp = &tv
	case OperationTypeBumpSequence:
		tv, ok := value.(BumpSequenceOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be BumpSequenceOp")
			return
		}
		result.BumpSequenceOp = &tv
	case OperationTypeManageBuyOffer:
		tv, ok := value.(ManageBuyOfferOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageBuyOfferOp")
			return
		}
		result.ManageBuyOfferOp = &tv
	case OperationTypePathPaymentStrictSend:
		tv, ok := value.(PathPaymentStrictSendOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be PathPaymentStrictSendOp")
			return
		}
		result.PathPaymentStrictSendOp = &tv
	case OperationTypeCreateClaimableBalance:
		tv, ok := value.(CreateClaimableBalanceOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be CreateClaimableBalanceOp")
			return
		}
		result.CreateClaimableBalanceOp = &tv
	case OperationTypeClaimClaimableBalance:
		tv, ok := value.(ClaimClaimableBalanceOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be ClaimClaimableBalanceOp")
			return
		}
		result.ClaimClaimableBalanceOp = &tv
	case OperationTypeBeginSponsoringFutureReserves:
		tv, ok := value.(BeginSponsoringFutureReservesOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be BeginSponsoringFutureReservesOp")
			return
		}
		result.BeginSponsoringFutureReservesOp = &tv
	case OperationTypeEndSponsoringFutureReserves:
		// void
	case OperationTypeRevokeSponsorship:
		tv, ok := value.(RevokeSponsorshipOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be RevokeSponsorshipOp")
			return
		}
		result.RevokeSponsorshipOp = &tv
	case OperationTypeClawback:
		tv, ok := value.(ClawbackOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be ClawbackOp")
			return
		}
		result.ClawbackOp = &tv
	case OperationTypeClawbackClaimableBalance:
		tv, ok := value.(ClawbackClaimableBalanceOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be ClawbackClaimableBalanceOp")
			return
		}
		result.ClawbackClaimableBalanceOp = &tv
	case OperationTypeSetTrustLineFlags:
		tv, ok := value.(SetTrustLineFlagsOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be SetTrustLineFlagsOp")
			return
		}
		result.SetTrustLineFlagsOp = &tv
	case OperationTypeLiquidityPoolDeposit:
		tv, ok := value.(LiquidityPoolDepositOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be LiquidityPoolDepositOp")
			return
		}
		result.LiquidityPoolDepositOp = &tv
	case OperationTypeLiquidityPoolWithdraw:
		tv, ok := value.(LiquidityPoolWithdrawOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be LiquidityPoolWithdrawOp")
			return
		}
		result.LiquidityPoolWithdrawOp = &tv
	}
	return
}

// MustCreateAccountOp retrieves the CreateAccountOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustCreateAccountOp() CreateAccountOp {
	val, ok := u.GetCreateAccountOp()

	if !ok {
		panic("arm CreateAccountOp is not set")
	}

	return val
}

// GetCreateAccountOp retrieves the CreateAccountOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetCreateAccountOp() (result CreateAccountOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "CreateAccountOp" {
		result = *u.CreateAccountOp
		ok = true
	}

	return
}

// MustPaymentOp retrieves the PaymentOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustPaymentOp() PaymentOp {
	val, ok := u.GetPaymentOp()

	if !ok {
		panic("arm PaymentOp is not set")
	}

	return val
}

// GetPaymentOp retrieves the PaymentOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetPaymentOp() (result PaymentOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "PaymentOp" {
		result = *u.PaymentOp
		ok = true
	}

	return
}

// MustPathPaymentStrictReceiveOp retrieves the PathPaymentStrictReceiveOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustPathPaymentStrictReceiveOp() PathPaymentStrictReceiveOp {
	val, ok := u.GetPathPaymentStrictReceiveOp()

	if !ok {
		panic("arm PathPaymentStrictReceiveOp is not set")
	}

	return val
}

// GetPathPaymentStrictReceiveOp retrieves the PathPaymentStrictReceiveOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetPathPaymentStrictReceiveOp() (result PathPaymentStrictReceiveOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "PathPaymentStrictReceiveOp" {
		result = *u.PathPaymentStrictReceiveOp
		ok = true
	}

	return
}

// MustManageSellOfferOp retrieves the ManageSellOfferOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustManageSellOfferOp() ManageSellOfferOp {
	val, ok := u.GetManageSellOfferOp()

	if !ok {
		panic("arm ManageSellOfferOp is not set")
	}

	return val
}

// GetManageSellOfferOp retrieves the ManageSellOfferOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetManageSellOfferOp() (result ManageSellOfferOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageSellOfferOp" {
		result = *u.ManageSellOfferOp
		ok = true
	}

	return
}

// MustCreatePassiveSellOfferOp retrieves the CreatePassiveSellOfferOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustCreatePassiveSellOfferOp() CreatePassiveSellOfferOp {
	val, ok := u.GetCreatePassiveSellOfferOp()

	if !ok {
		panic("arm CreatePassiveSellOfferOp is not set")
	}

	return val
}

// GetCreatePassiveSellOfferOp retrieves the CreatePassiveSellOfferOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetCreatePassiveSellOfferOp() (result CreatePassiveSellOfferOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "CreatePassiveSellOfferOp" {
		result = *u.CreatePassiveSellOfferOp
		ok = true
	}

	return
}

// MustSetOptionsOp retrieves the SetOptionsOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustSetOptionsOp() SetOptionsOp {
	val, ok := u.GetSetOptionsOp()

	if !ok {
		panic("arm SetOptionsOp is not set")
	}

	return val
}

// GetSetOptionsOp retrieves the SetOptionsOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetSetOptionsOp() (result SetOptionsOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "SetOptionsOp" {
		result = *u.SetOptionsOp
		ok = true
	}

	return
}

// MustChangeTrustOp retrieves the ChangeTrustOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustChangeTrustOp() ChangeTrustOp {
	val, ok := u.GetChangeTrustOp()

	if !ok {
		panic("arm ChangeTrustOp is not set")
	}

	return val
}

// GetChangeTrustOp retrieves the ChangeTrustOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetChangeTrustOp() (result ChangeTrustOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ChangeTrustOp" {
		result = *u.ChangeTrustOp
		ok = true
	}

	return
}

// MustAllowTrustOp retrieves the AllowTrustOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustAllowTrustOp() AllowTrustOp {
	val, ok := u.GetAllowTrustOp()

	if !ok {
		panic("arm AllowTrustOp is not set")
	}

	return val
}

// GetAllowTrustOp retrieves the AllowTrustOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetAllowTrustOp() (result AllowTrustOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AllowTrustOp" {
		result = *u.AllowTrustOp
		ok = true
	}

	return
}

// MustDestination retrieves the Destination value from the union,
// panicing if the value is not set.
func (u OperationBody) MustDestination() MuxedAccount {
	val, ok := u.GetDestination()

	if !ok {
		panic("arm Destination is not set")
	}

	return val
}

// GetDestination retrieves the Destination value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetDestination() (result MuxedAccount, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Destination" {
		result = *u.Destination
		ok = true
	}

	return
}

// MustManageDataOp retrieves the ManageDataOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustManageDataOp() ManageDataOp {
	val, ok := u.GetManageDataOp()

	if !ok {
		panic("arm ManageDataOp is not set")
	}

	return val
}

// GetManageDataOp retrieves the ManageDataOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetManageDataOp() (result ManageDataOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageDataOp" {
		result = *u.ManageDataOp
		ok = true
	}

	return
}

// MustBumpSequenceOp retrieves the BumpSequenceOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustBumpSequenceOp() BumpSequenceOp {
	val, ok := u.GetBumpSequenceOp()

	if !ok {
		panic("arm BumpSequenceOp is not set")
	}

	return val
}

// GetBumpSequenceOp retrieves the BumpSequenceOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetBumpSequenceOp() (result BumpSequenceOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "BumpSequenceOp" {
		result = *u.BumpSequenceOp
		ok = true
	}

	return
}

// MustManageBuyOfferOp retrieves the ManageBuyOfferOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustManageBuyOfferOp() ManageBuyOfferOp {
	val, ok := u.GetManageBuyOfferOp()

	if !ok {
		panic("arm ManageBuyOfferOp is not set")
	}

	return val
}

// GetManageBuyOfferOp retrieves the ManageBuyOfferOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetManageBuyOfferOp() (result ManageBuyOfferOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageBuyOfferOp" {
		result = *u.ManageBuyOfferOp
		ok = true
	}

	return
}

// MustPathPaymentStrictSendOp retrieves the PathPaymentStrictSendOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustPathPaymentStrictSendOp() PathPaymentStrictSendOp {
	val, ok := u.GetPathPaymentStrictSendOp()

	if !ok {
		panic("arm PathPaymentStrictSendOp is not set")
	}

	return val
}

// GetPathPaymentStrictSendOp retrieves the PathPaymentStrictSendOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetPathPaymentStrictSendOp() (result PathPaymentStrictSendOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "PathPaymentStrictSendOp" {
		result = *u.PathPaymentStrictSendOp
		ok = true
	}

	return
}

// MustCreateClaimableBalanceOp retrieves the CreateClaimableBalanceOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustCreateClaimableBalanceOp() CreateClaimableBalanceOp {
	val, ok := u.GetCreateClaimableBalanceOp()

	if !ok {
		panic("arm CreateClaimableBalanceOp is not set")
	}

	return val
}

// GetCreateClaimableBalanceOp retrieves the CreateClaimableBalanceOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetCreateClaimableBalanceOp() (result CreateClaimableBalanceOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "CreateClaimableBalanceOp" {
		result = *u.CreateClaimableBalanceOp
		ok = true
	}

	return
}

// MustClaimClaimableBalanceOp retrieves the ClaimClaimableBalanceOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustClaimClaimableBalanceOp() ClaimClaimableBalanceOp {
	val, ok := u.GetClaimClaimableBalanceOp()

	if !ok {
		panic("arm ClaimClaimableBalanceOp is not set")
	}

	return val
}

// GetClaimClaimableBalanceOp retrieves the ClaimClaimableBalanceOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetClaimClaimableBalanceOp() (result ClaimClaimableBalanceOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ClaimClaimableBalanceOp" {
		result = *u.ClaimClaimableBalanceOp
		ok = true
	}

	return
}

// MustBeginSponsoringFutureReservesOp retrieves the BeginSponsoringFutureReservesOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustBeginSponsoringFutureReservesOp() BeginSponsoringFutureReservesOp {
	val, ok := u.GetBeginSponsoringFutureReservesOp()

	if !ok {
		panic("arm BeginSponsoringFutureReservesOp is not set")
	}

	return val
}

// GetBeginSponsoringFutureReservesOp retrieves the BeginSponsoringFutureReservesOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetBeginSponsoringFutureReservesOp() (result BeginSponsoringFutureReservesOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "BeginSponsoringFutureReservesOp" {
		result = *u.BeginSponsoringFutureReservesOp
		ok = true
	}

	return
}

// MustRevokeSponsorshipOp retrieves the RevokeSponsorshipOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustRevokeSponsorshipOp() RevokeSponsorshipOp {
	val, ok := u.GetRevokeSponsorshipOp()

	if !ok {
		panic("arm RevokeSponsorshipOp is not set")
	}

	return val
}

// GetRevokeSponsorshipOp retrieves the RevokeSponsorshipOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetRevokeSponsorshipOp() (result RevokeSponsorshipOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "RevokeSponsorshipOp" {
		result = *u.RevokeSponsorshipOp
		ok = true
	}

	return
}

// MustClawbackOp retrieves the ClawbackOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustClawbackOp() ClawbackOp {
	val, ok := u.GetClawbackOp()

	if !ok {
		panic("arm ClawbackOp is not set")
	}

	return val
}

// GetClawbackOp retrieves the ClawbackOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetClawbackOp() (result ClawbackOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ClawbackOp" {
		result = *u.ClawbackOp
		ok = true
	}

	return
}

// MustClawbackClaimableBalanceOp retrieves the ClawbackClaimableBalanceOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustClawbackClaimableBalanceOp() ClawbackClaimableBalanceOp {
	val, ok := u.GetClawbackClaimableBalanceOp()

	if !ok {
		panic("arm ClawbackClaimableBalanceOp is not set")
	}

	return val
}

// GetClawbackClaimableBalanceOp retrieves the ClawbackClaimableBalanceOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetClawbackClaimableBalanceOp() (result ClawbackClaimableBalanceOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ClawbackClaimableBalanceOp" {
		result = *u.ClawbackClaimableBalanceOp
		ok = true
	}

	return
}

// MustSetTrustLineFlagsOp retrieves the SetTrustLineFlagsOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustSetTrustLineFlagsOp() SetTrustLineFlagsOp {
	val, ok := u.GetSetTrustLineFlagsOp()

	if !ok {
		panic("arm SetTrustLineFlagsOp is not set")
	}

	return val
}

// GetSetTrustLineFlagsOp retrieves the SetTrustLineFlagsOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetSetTrustLineFlagsOp() (result SetTrustLineFlagsOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "SetTrustLineFlagsOp" {
		result = *u.SetTrustLineFlagsOp
		ok = true
	}

	return
}

// MustLiquidityPoolDepositOp retrieves the LiquidityPoolDepositOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustLiquidityPoolDepositOp() LiquidityPoolDepositOp {
	val, ok := u.GetLiquidityPoolDepositOp()

	if !ok {
		panic("arm LiquidityPoolDepositOp is not set")
	}

	return val
}

// GetLiquidityPoolDepositOp retrieves the LiquidityPoolDepositOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetLiquidityPoolDepositOp() (result LiquidityPoolDepositOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "LiquidityPoolDepositOp" {
		result = *u.LiquidityPoolDepositOp
		ok = true
	}

	return
}

// MustLiquidityPoolWithdrawOp retrieves the LiquidityPoolWithdrawOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustLiquidityPoolWithdrawOp() LiquidityPoolWithdrawOp {
	val, ok := u.GetLiquidityPoolWithdrawOp()

	if !ok {
		panic("arm LiquidityPoolWithdrawOp is not set")
	}

	return val
}

// GetLiquidityPoolWithdrawOp retrieves the LiquidityPoolWithdrawOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetLiquidityPoolWithdrawOp() (result LiquidityPoolWithdrawOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "LiquidityPoolWithdrawOp" {
		result = *u.LiquidityPoolWithdrawOp
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u OperationBody) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch OperationType(u.Type) {
	case OperationTypeCreateAccount:
		if err = (*u.CreateAccountOp).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypePayment:
		if err = (*u.PaymentOp).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypePathPaymentStrictReceive:
		if err = (*u.PathPaymentStrictReceiveOp).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeManageSellOffer:
		if err = (*u.ManageSellOfferOp).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeCreatePassiveSellOffer:
		if err = (*u.CreatePassiveSellOfferOp).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeSetOptions:
		if err = (*u.SetOptionsOp).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeChangeTrust:
		if err = (*u.ChangeTrustOp).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeAllowTrust:
		if err = (*u.AllowTrustOp).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeAccountMerge:
		if err = (*u.Destination).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeInflation:
		// Void
		return nil
	case OperationTypeManageData:
		if err = (*u.ManageDataOp).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeBumpSequence:
		if err = (*u.BumpSequenceOp).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeManageBuyOffer:
		if err = (*u.ManageBuyOfferOp).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypePathPaymentStrictSend:
		if err = (*u.PathPaymentStrictSendOp).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeCreateClaimableBalance:
		if err = (*u.CreateClaimableBalanceOp).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeClaimClaimableBalance:
		if err = (*u.ClaimClaimableBalanceOp).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeBeginSponsoringFutureReserves:
		if err = (*u.BeginSponsoringFutureReservesOp).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeEndSponsoringFutureReserves:
		// Void
		return nil
	case OperationTypeRevokeSponsorship:
		if err = (*u.RevokeSponsorshipOp).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeClawback:
		if err = (*u.ClawbackOp).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeClawbackClaimableBalance:
		if err = (*u.ClawbackClaimableBalanceOp).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeSetTrustLineFlags:
		if err = (*u.SetTrustLineFlagsOp).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeLiquidityPoolDeposit:
		if err = (*u.LiquidityPoolDepositOp).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeLiquidityPoolWithdraw:
		if err = (*u.LiquidityPoolWithdrawOp).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (OperationType) switch value '%d' is not valid for union OperationBody", u.Type)
}

var _ decoderFrom = (*OperationBody)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *OperationBody) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding OperationType: %s", err)
	}
	switch OperationType(u.Type) {
	case OperationTypeCreateAccount:
		u.CreateAccountOp = new(CreateAccountOp)
		nTmp, err = (*u.CreateAccountOp).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding CreateAccountOp: %s", err)
		}
		return n, nil
	case OperationTypePayment:
		u.PaymentOp = new(PaymentOp)
		nTmp, err = (*u.PaymentOp).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding PaymentOp: %s", err)
		}
		return n, nil
	case OperationTypePathPaymentStrictReceive:
		u.PathPaymentStrictReceiveOp = new(PathPaymentStrictReceiveOp)
		nTmp, err = (*u.PathPaymentStrictReceiveOp).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding PathPaymentStrictReceiveOp: %s", err)
		}
		return n, nil
	case OperationTypeManageSellOffer:
		u.ManageSellOfferOp = new(ManageSellOfferOp)
		nTmp, err = (*u.ManageSellOfferOp).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ManageSellOfferOp: %s", err)
		}
		return n, nil
	case OperationTypeCreatePassiveSellOffer:
		u.CreatePassiveSellOfferOp = new(CreatePassiveSellOfferOp)
		nTmp, err = (*u.CreatePassiveSellOfferOp).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding CreatePassiveSellOfferOp: %s", err)
		}
		return n, nil
	case OperationTypeSetOptions:
		u.SetOptionsOp = new(SetOptionsOp)
		nTmp, err = (*u.SetOptionsOp).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding SetOptionsOp: %s", err)
		}
		return n, nil
	case OperationTypeChangeTrust:
		u.ChangeTrustOp = new(ChangeTrustOp)
		nTmp, err = (*u.ChangeTrustOp).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ChangeTrustOp: %s", err)
		}
		return n, nil
	case OperationTypeAllowTrust:
		u.AllowTrustOp = new(AllowTrustOp)
		nTmp, err = (*u.AllowTrustOp).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding AllowTrustOp: %s", err)
		}
		return n, nil
	case OperationTypeAccountMerge:
		u.Destination = new(MuxedAccount)
		nTmp, err = (*u.Destination).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding MuxedAccount: %s", err)
		}
		return n, nil
	case OperationTypeInflation:
		// Void
		return n, nil
	case OperationTypeManageData:
		u.ManageDataOp = new(ManageDataOp)
		nTmp, err = (*u.ManageDataOp).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ManageDataOp: %s", err)
		}
		return n, nil
	case OperationTypeBumpSequence:
		u.BumpSequenceOp = new(BumpSequenceOp)
		nTmp, err = (*u.BumpSequenceOp).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding BumpSequenceOp: %s", err)
		}
		return n, nil
	case OperationTypeManageBuyOffer:
		u.ManageBuyOfferOp = new(ManageBuyOfferOp)
		nTmp, err = (*u.ManageBuyOfferOp).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ManageBuyOfferOp: %s", err)
		}
		return n, nil
	case OperationTypePathPaymentStrictSend:
		u.PathPaymentStrictSendOp = new(PathPaymentStrictSendOp)
		nTmp, err = (*u.PathPaymentStrictSendOp).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding PathPaymentStrictSendOp: %s", err)
		}
		return n, nil
	case OperationTypeCreateClaimableBalance:
		u.CreateClaimableBalanceOp = new(CreateClaimableBalanceOp)
		nTmp, err = (*u.CreateClaimableBalanceOp).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding CreateClaimableBalanceOp: %s", err)
		}
		return n, nil
	case OperationTypeClaimClaimableBalance:
		u.ClaimClaimableBalanceOp = new(ClaimClaimableBalanceOp)
		nTmp, err = (*u.ClaimClaimableBalanceOp).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ClaimClaimableBalanceOp: %s", err)
		}
		return n, nil
	case OperationTypeBeginSponsoringFutureReserves:
		u.BeginSponsoringFutureReservesOp = new(BeginSponsoringFutureReservesOp)
		nTmp, err = (*u.BeginSponsoringFutureReservesOp).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding BeginSponsoringFutureReservesOp: %s", err)
		}
		return n, nil
	case OperationTypeEndSponsoringFutureReserves:
		// Void
		return n, nil
	case OperationTypeRevokeSponsorship:
		u.RevokeSponsorshipOp = new(RevokeSponsorshipOp)
		nTmp, err = (*u.RevokeSponsorshipOp).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding RevokeSponsorshipOp: %s", err)
		}
		return n, nil
	case OperationTypeClawback:
		u.ClawbackOp = new(ClawbackOp)
		nTmp, err = (*u.ClawbackOp).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ClawbackOp: %s", err)
		}
		return n, nil
	case OperationTypeClawbackClaimableBalance:
		u.ClawbackClaimableBalanceOp = new(ClawbackClaimableBalanceOp)
		nTmp, err = (*u.ClawbackClaimableBalanceOp).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ClawbackClaimableBalanceOp: %s", err)
		}
		return n, nil
	case OperationTypeSetTrustLineFlags:
		u.SetTrustLineFlagsOp = new(SetTrustLineFlagsOp)
		nTmp, err = (*u.SetTrustLineFlagsOp).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding SetTrustLineFlagsOp: %s", err)
		}
		return n, nil
	case OperationTypeLiquidityPoolDeposit:
		u.LiquidityPoolDepositOp = new(LiquidityPoolDepositOp)
		nTmp, err = (*u.LiquidityPoolDepositOp).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding LiquidityPoolDepositOp: %s", err)
		}
		return n, nil
	case OperationTypeLiquidityPoolWithdraw:
		u.LiquidityPoolWithdrawOp = new(LiquidityPoolWithdrawOp)
		nTmp, err = (*u.LiquidityPoolWithdrawOp).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding LiquidityPoolWithdrawOp: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union OperationBody has invalid Type (OperationType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s OperationBody) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *OperationBody) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*OperationBody)(nil)
	_ encoding.BinaryUnmarshaler = (*OperationBody)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s OperationBody) xdrType() {}

var _ xdrType = (*OperationBody)(nil)

// Operation is an XDR Struct defines as:
//
//   struct Operation
//    {
//        // sourceAccount is the account used to run the operation
//        // if not set, the runtime defaults to "sourceAccount" specified at
//        // the transaction level
//        MuxedAccount* sourceAccount;
//
//        union switch (OperationType type)
//        {
//        case CREATE_ACCOUNT:
//            CreateAccountOp createAccountOp;
//        case PAYMENT:
//            PaymentOp paymentOp;
//        case PATH_PAYMENT_STRICT_RECEIVE:
//            PathPaymentStrictReceiveOp pathPaymentStrictReceiveOp;
//        case MANAGE_SELL_OFFER:
//            ManageSellOfferOp manageSellOfferOp;
//        case CREATE_PASSIVE_SELL_OFFER:
//            CreatePassiveSellOfferOp createPassiveSellOfferOp;
//        case SET_OPTIONS:
//            SetOptionsOp setOptionsOp;
//        case CHANGE_TRUST:
//            ChangeTrustOp changeTrustOp;
//        case ALLOW_TRUST:
//            AllowTrustOp allowTrustOp;
//        case ACCOUNT_MERGE:
//            MuxedAccount destination;
//        case INFLATION:
//            void;
//        case MANAGE_DATA:
//            ManageDataOp manageDataOp;
//        case BUMP_SEQUENCE:
//            BumpSequenceOp bumpSequenceOp;
//        case MANAGE_BUY_OFFER:
//            ManageBuyOfferOp manageBuyOfferOp;
//        case PATH_PAYMENT_STRICT_SEND:
//            PathPaymentStrictSendOp pathPaymentStrictSendOp;
//        case CREATE_CLAIMABLE_BALANCE:
//            CreateClaimableBalanceOp createClaimableBalanceOp;
//        case CLAIM_CLAIMABLE_BALANCE:
//            ClaimClaimableBalanceOp claimClaimableBalanceOp;
//        case BEGIN_SPONSORING_FUTURE_RESERVES:
//            BeginSponsoringFutureReservesOp beginSponsoringFutureReservesOp;
//        case END_SPONSORING_FUTURE_RESERVES:
//            void;
//        case REVOKE_SPONSORSHIP:
//            RevokeSponsorshipOp revokeSponsorshipOp;
//        case CLAWBACK:
//            ClawbackOp clawbackOp;
//        case CLAWBACK_CLAIMABLE_BALANCE:
//            ClawbackClaimableBalanceOp clawbackClaimableBalanceOp;
//        case SET_TRUST_LINE_FLAGS:
//            SetTrustLineFlagsOp setTrustLineFlagsOp;
//        case LIQUIDITY_POOL_DEPOSIT:
//            LiquidityPoolDepositOp liquidityPoolDepositOp;
//        case LIQUIDITY_POOL_WITHDRAW:
//            LiquidityPoolWithdrawOp liquidityPoolWithdrawOp;
//        }
//        body;
//    };
//
type Operation struct {
	SourceAccount *MuxedAccount
	Body          OperationBody
}

// EncodeTo encodes this value using the Encoder.
func (s *Operation) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeBool(s.SourceAccount != nil); err != nil {
		return err
	}
	if s.SourceAccount != nil {
		if err = (*s.SourceAccount).EncodeTo(e); err != nil {
			return err
		}
	}
	if err = s.Body.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Operation)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Operation) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	var b bool
	b, nTmp, err = d.DecodeBool()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding MuxedAccount: %s", err)
	}
	s.SourceAccount = nil
	if b {
		s.SourceAccount = new(MuxedAccount)
		nTmp, err = s.SourceAccount.DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding MuxedAccount: %s", err)
		}
	}
	nTmp, err = s.Body.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding OperationBody: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Operation) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Operation) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Operation)(nil)
	_ encoding.BinaryUnmarshaler = (*Operation)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s Operation) xdrType() {}

var _ xdrType = (*Operation)(nil)

// HashIdPreimageOperationId is an XDR NestedStruct defines as:
//
//   struct
//        {
//            AccountID sourceAccount;
//            SequenceNumber seqNum;
//            uint32 opNum;
//        }
//
type HashIdPreimageOperationId struct {
	SourceAccount AccountId
	SeqNum        SequenceNumber
	OpNum         Uint32
}

// EncodeTo encodes this value using the Encoder.
func (s *HashIdPreimageOperationId) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.SourceAccount.EncodeTo(e); err != nil {
		return err
	}
	if err = s.SeqNum.EncodeTo(e); err != nil {
		return err
	}
	if err = s.OpNum.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*HashIdPreimageOperationId)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *HashIdPreimageOperationId) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.SourceAccount.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountId: %s", err)
	}
	nTmp, err = s.SeqNum.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SequenceNumber: %s", err)
	}
	nTmp, err = s.OpNum.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s HashIdPreimageOperationId) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *HashIdPreimageOperationId) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*HashIdPreimageOperationId)(nil)
	_ encoding.BinaryUnmarshaler = (*HashIdPreimageOperationId)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s HashIdPreimageOperationId) xdrType() {}

var _ xdrType = (*HashIdPreimageOperationId)(nil)

// HashIdPreimageRevokeId is an XDR NestedStruct defines as:
//
//   struct
//        {
//            AccountID sourceAccount;
//            SequenceNumber seqNum;
//            uint32 opNum;
//            PoolID liquidityPoolID;
//            Asset asset;
//        }
//
type HashIdPreimageRevokeId struct {
	SourceAccount   AccountId
	SeqNum          SequenceNumber
	OpNum           Uint32
	LiquidityPoolId PoolId
	Asset           Asset
}

// EncodeTo encodes this value using the Encoder.
func (s *HashIdPreimageRevokeId) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.SourceAccount.EncodeTo(e); err != nil {
		return err
	}
	if err = s.SeqNum.EncodeTo(e); err != nil {
		return err
	}
	if err = s.OpNum.EncodeTo(e); err != nil {
		return err
	}
	if err = s.LiquidityPoolId.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Asset.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*HashIdPreimageRevokeId)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *HashIdPreimageRevokeId) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.SourceAccount.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountId: %s", err)
	}
	nTmp, err = s.SeqNum.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SequenceNumber: %s", err)
	}
	nTmp, err = s.OpNum.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.LiquidityPoolId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding PoolId: %s", err)
	}
	nTmp, err = s.Asset.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Asset: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s HashIdPreimageRevokeId) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *HashIdPreimageRevokeId) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*HashIdPreimageRevokeId)(nil)
	_ encoding.BinaryUnmarshaler = (*HashIdPreimageRevokeId)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s HashIdPreimageRevokeId) xdrType() {}

var _ xdrType = (*HashIdPreimageRevokeId)(nil)

// HashIdPreimage is an XDR Union defines as:
//
//   union HashIDPreimage switch (EnvelopeType type)
//    {
//    case ENVELOPE_TYPE_OP_ID:
//        struct
//        {
//            AccountID sourceAccount;
//            SequenceNumber seqNum;
//            uint32 opNum;
//        } operationID;
//    case ENVELOPE_TYPE_POOL_REVOKE_OP_ID:
//        struct
//        {
//            AccountID sourceAccount;
//            SequenceNumber seqNum;
//            uint32 opNum;
//            PoolID liquidityPoolID;
//            Asset asset;
//        } revokeID;
//    };
//
type HashIdPreimage struct {
	Type        EnvelopeType
	OperationId *HashIdPreimageOperationId
	RevokeId    *HashIdPreimageRevokeId
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u HashIdPreimage) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of HashIdPreimage
func (u HashIdPreimage) ArmForSwitch(sw int32) (string, bool) {
	switch EnvelopeType(sw) {
	case EnvelopeTypeEnvelopeTypeOpId:
		return "OperationId", true
	case EnvelopeTypeEnvelopeTypePoolRevokeOpId:
		return "RevokeId", true
	}
	return "-", false
}

// NewHashIdPreimage creates a new  HashIdPreimage.
func NewHashIdPreimage(aType EnvelopeType, value interface{}) (result HashIdPreimage, err error) {
	result.Type = aType
	switch EnvelopeType(aType) {
	case EnvelopeTypeEnvelopeTypeOpId:
		tv, ok := value.(HashIdPreimageOperationId)
		if !ok {
			err = fmt.Errorf("invalid value, must be HashIdPreimageOperationId")
			return
		}
		result.OperationId = &tv
	case EnvelopeTypeEnvelopeTypePoolRevokeOpId:
		tv, ok := value.(HashIdPreimageRevokeId)
		if !ok {
			err = fmt.Errorf("invalid value, must be HashIdPreimageRevokeId")
			return
		}
		result.RevokeId = &tv
	}
	return
}

// MustOperationId retrieves the OperationId value from the union,
// panicing if the value is not set.
func (u HashIdPreimage) MustOperationId() HashIdPreimageOperationId {
	val, ok := u.GetOperationId()

	if !ok {
		panic("arm OperationId is not set")
	}

	return val
}

// GetOperationId retrieves the OperationId value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u HashIdPreimage) GetOperationId() (result HashIdPreimageOperationId, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "OperationId" {
		result = *u.OperationId
		ok = true
	}

	return
}

// MustRevokeId retrieves the RevokeId value from the union,
// panicing if the value is not set.
func (u HashIdPreimage) MustRevokeId() HashIdPreimageRevokeId {
	val, ok := u.GetRevokeId()

	if !ok {
		panic("arm RevokeId is not set")
	}

	return val
}

// GetRevokeId retrieves the RevokeId value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u HashIdPreimage) GetRevokeId() (result HashIdPreimageRevokeId, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "RevokeId" {
		result = *u.RevokeId
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u HashIdPreimage) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch EnvelopeType(u.Type) {
	case EnvelopeTypeEnvelopeTypeOpId:
		if err = (*u.OperationId).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case EnvelopeTypeEnvelopeTypePoolRevokeOpId:
		if err = (*u.RevokeId).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (EnvelopeType) switch value '%d' is not valid for union HashIdPreimage", u.Type)
}

var _ decoderFrom = (*HashIdPreimage)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *HashIdPreimage) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding EnvelopeType: %s", err)
	}
	switch EnvelopeType(u.Type) {
	case EnvelopeTypeEnvelopeTypeOpId:
		u.OperationId = new(HashIdPreimageOperationId)
		nTmp, err = (*u.OperationId).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding HashIdPreimageOperationId: %s", err)
		}
		return n, nil
	case EnvelopeTypeEnvelopeTypePoolRevokeOpId:
		u.RevokeId = new(HashIdPreimageRevokeId)
		nTmp, err = (*u.RevokeId).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding HashIdPreimageRevokeId: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union HashIdPreimage has invalid Type (EnvelopeType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s HashIdPreimage) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *HashIdPreimage) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*HashIdPreimage)(nil)
	_ encoding.BinaryUnmarshaler = (*HashIdPreimage)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s HashIdPreimage) xdrType() {}

var _ xdrType = (*HashIdPreimage)(nil)

// MemoType is an XDR Enum defines as:
//
//   enum MemoType
//    {
//        MEMO_NONE = 0,
//        MEMO_TEXT = 1,
//        MEMO_ID = 2,
//        MEMO_HASH = 3,
//        MEMO_RETURN = 4
//    };
//
type MemoType int32

const (
	MemoTypeMemoNone   MemoType = 0
	MemoTypeMemoText   MemoType = 1
	MemoTypeMemoId     MemoType = 2
	MemoTypeMemoHash   MemoType = 3
	MemoTypeMemoReturn MemoType = 4
)

var memoTypeMap = map[int32]string{
	0: "MemoTypeMemoNone",
	1: "MemoTypeMemoText",
	2: "MemoTypeMemoId",
	3: "MemoTypeMemoHash",
	4: "MemoTypeMemoReturn",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for MemoType
func (e MemoType) ValidEnum(v int32) bool {
	_, ok := memoTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e MemoType) String() string {
	name, _ := memoTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e MemoType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := memoTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid MemoType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*MemoType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *MemoType) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding MemoType: %s", err)
	}
	if _, ok := memoTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid MemoType enum value", v)
	}
	*e = MemoType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s MemoType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *MemoType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*MemoType)(nil)
	_ encoding.BinaryUnmarshaler = (*MemoType)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s MemoType) xdrType() {}

var _ xdrType = (*MemoType)(nil)

// Memo is an XDR Union defines as:
//
//   union Memo switch (MemoType type)
//    {
//    case MEMO_NONE:
//        void;
//    case MEMO_TEXT:
//        string text<28>;
//    case MEMO_ID:
//        uint64 id;
//    case MEMO_HASH:
//        Hash hash; // the hash of what to pull from the content server
//    case MEMO_RETURN:
//        Hash retHash; // the hash of the tx you are rejecting
//    };
//
type Memo struct {
	Type    MemoType
	Text    *string `xdrmaxsize:"28"`
	Id      *Uint64
	Hash    *Hash
	RetHash *Hash
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u Memo) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of Memo
func (u Memo) ArmForSwitch(sw int32) (string, bool) {
	switch MemoType(sw) {
	case MemoTypeMemoNone:
		return "", true
	case MemoTypeMemoText:
		return "Text", true
	case MemoTypeMemoId:
		return "Id", true
	case MemoTypeMemoHash:
		return "Hash", true
	case MemoTypeMemoReturn:
		return "RetHash", true
	}
	return "-", false
}

// NewMemo creates a new  Memo.
func NewMemo(aType MemoType, value interface{}) (result Memo, err error) {
	result.Type = aType
	switch MemoType(aType) {
	case MemoTypeMemoNone:
		// void
	case MemoTypeMemoText:
		tv, ok := value.(string)
		if !ok {
			err = fmt.Errorf("invalid value, must be string")
			return
		}
		result.Text = &tv
	case MemoTypeMemoId:
		tv, ok := value.(Uint64)
		if !ok {
			err = fmt.Errorf("invalid value, must be Uint64")
			return
		}
		result.Id = &tv
	case MemoTypeMemoHash:
		tv, ok := value.(Hash)
		if !ok {
			err = fmt.Errorf("invalid value, must be Hash")
			return
		}
		result.Hash = &tv
	case MemoTypeMemoReturn:
		tv, ok := value.(Hash)
		if !ok {
			err = fmt.Errorf("invalid value, must be Hash")
			return
		}
		result.RetHash = &tv
	}
	return
}

// MustText retrieves the Text value from the union,
// panicing if the value is not set.
func (u Memo) MustText() string {
	val, ok := u.GetText()

	if !ok {
		panic("arm Text is not set")
	}

	return val
}

// GetText retrieves the Text value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u Memo) GetText() (result string, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Text" {
		result = *u.Text
		ok = true
	}

	return
}

// MustId retrieves the Id value from the union,
// panicing if the value is not set.
func (u Memo) MustId() Uint64 {
	val, ok := u.GetId()

	if !ok {
		panic("arm Id is not set")
	}

	return val
}

// GetId retrieves the Id value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u Memo) GetId() (result Uint64, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Id" {
		result = *u.Id
		ok = true
	}

	return
}

// MustHash retrieves the Hash value from the union,
// panicing if the value is not set.
func (u Memo) MustHash() Hash {
	val, ok := u.GetHash()

	if !ok {
		panic("arm Hash is not set")
	}

	return val
}

// GetHash retrieves the Hash value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u Memo) GetHash() (result Hash, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Hash" {
		result = *u.Hash
		ok = true
	}

	return
}

// MustRetHash retrieves the RetHash value from the union,
// panicing if the value is not set.
func (u Memo) MustRetHash() Hash {
	val, ok := u.GetRetHash()

	if !ok {
		panic("arm RetHash is not set")
	}

	return val
}

// GetRetHash retrieves the RetHash value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u Memo) GetRetHash() (result Hash, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "RetHash" {
		result = *u.RetHash
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u Memo) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch MemoType(u.Type) {
	case MemoTypeMemoNone:
		// Void
		return nil
	case MemoTypeMemoText:
		if _, err = e.EncodeString(string((*u.Text))); err != nil {
			return err
		}
		return nil
	case MemoTypeMemoId:
		if err = (*u.Id).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case MemoTypeMemoHash:
		if err = (*u.Hash).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case MemoTypeMemoReturn:
		if err = (*u.RetHash).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (MemoType) switch value '%d' is not valid for union Memo", u.Type)
}

var _ decoderFrom = (*Memo)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *Memo) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding MemoType: %s", err)
	}
	switch MemoType(u.Type) {
	case MemoTypeMemoNone:
		// Void
		return n, nil
	case MemoTypeMemoText:
		u.Text = new(string)
		(*u.Text), nTmp, err = d.DecodeString(28)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Text: %s", err)
		}
		return n, nil
	case MemoTypeMemoId:
		u.Id = new(Uint64)
		nTmp, err = (*u.Id).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Uint64: %s", err)
		}
		return n, nil
	case MemoTypeMemoHash:
		u.Hash = new(Hash)
		nTmp, err = (*u.Hash).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Hash: %s", err)
		}
		return n, nil
	case MemoTypeMemoReturn:
		u.RetHash = new(Hash)
		nTmp, err = (*u.RetHash).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Hash: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union Memo has invalid Type (MemoType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Memo) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Memo) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Memo)(nil)
	_ encoding.BinaryUnmarshaler = (*Memo)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s Memo) xdrType() {}

var _ xdrType = (*Memo)(nil)

// TimeBounds is an XDR Struct defines as:
//
//   struct TimeBounds
//    {
//        TimePoint minTime;
//        TimePoint maxTime; // 0 here means no maxTime
//    };
//
type TimeBounds struct {
	MinTime TimePoint
	MaxTime TimePoint
}

// EncodeTo encodes this value using the Encoder.
func (s *TimeBounds) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.MinTime.EncodeTo(e); err != nil {
		return err
	}
	if err = s.MaxTime.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*TimeBounds)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *TimeBounds) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.MinTime.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TimePoint: %s", err)
	}
	nTmp, err = s.MaxTime.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TimePoint: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TimeBounds) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TimeBounds) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TimeBounds)(nil)
	_ encoding.BinaryUnmarshaler = (*TimeBounds)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TimeBounds) xdrType() {}

var _ xdrType = (*TimeBounds)(nil)

// MaxOpsPerTx is an XDR Const defines as:
//
//   const MAX_OPS_PER_TX = 100;
//
const MaxOpsPerTx = 100

// TransactionV0Ext is an XDR NestedUnion defines as:
//
//   union switch (int v)
//        {
//        case 0:
//            void;
//        }
//
type TransactionV0Ext struct {
	V int32
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TransactionV0Ext) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TransactionV0Ext
func (u TransactionV0Ext) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	}
	return "-", false
}

// NewTransactionV0Ext creates a new  TransactionV0Ext.
func NewTransactionV0Ext(v int32, value interface{}) (result TransactionV0Ext, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u TransactionV0Ext) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union TransactionV0Ext", u.V)
}

var _ decoderFrom = (*TransactionV0Ext)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *TransactionV0Ext) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %s", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	}
	return n, fmt.Errorf("union TransactionV0Ext has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionV0Ext) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionV0Ext) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionV0Ext)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionV0Ext)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TransactionV0Ext) xdrType() {}

var _ xdrType = (*TransactionV0Ext)(nil)

// TransactionV0 is an XDR Struct defines as:
//
//   struct TransactionV0
//    {
//        uint256 sourceAccountEd25519;
//        uint32 fee;
//        SequenceNumber seqNum;
//        TimeBounds* timeBounds;
//        Memo memo;
//        Operation operations<MAX_OPS_PER_TX>;
//        union switch (int v)
//        {
//        case 0:
//            void;
//        }
//        ext;
//    };
//
type TransactionV0 struct {
	SourceAccountEd25519 Uint256
	Fee                  Uint32
	SeqNum               SequenceNumber
	TimeBounds           *TimeBounds
	Memo                 Memo
	Operations           []Operation `xdrmaxsize:"100"`
	Ext                  TransactionV0Ext
}

// EncodeTo encodes this value using the Encoder.
func (s *TransactionV0) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.SourceAccountEd25519.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Fee.EncodeTo(e); err != nil {
		return err
	}
	if err = s.SeqNum.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeBool(s.TimeBounds != nil); err != nil {
		return err
	}
	if s.TimeBounds != nil {
		if err = (*s.TimeBounds).EncodeTo(e); err != nil {
			return err
		}
	}
	if err = s.Memo.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeUint(uint32(len(s.Operations))); err != nil {
		return err
	}
	for i := 0; i < len(s.Operations); i++ {
		if err = s.Operations[i].EncodeTo(e); err != nil {
			return err
		}
	}
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*TransactionV0)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *TransactionV0) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.SourceAccountEd25519.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint256: %s", err)
	}
	nTmp, err = s.Fee.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.SeqNum.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SequenceNumber: %s", err)
	}
	var b bool
	b, nTmp, err = d.DecodeBool()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TimeBounds: %s", err)
	}
	s.TimeBounds = nil
	if b {
		s.TimeBounds = new(TimeBounds)
		nTmp, err = s.TimeBounds.DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding TimeBounds: %s", err)
		}
	}
	nTmp, err = s.Memo.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Memo: %s", err)
	}
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Operation: %s", err)
	}
	if l > 100 {
		return n, fmt.Errorf("decoding Operation: data size (%d) exceeds size limit (100)", l)
	}
	s.Operations = nil
	if l > 0 {
		s.Operations = make([]Operation, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.Operations[i].DecodeFrom(d)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding Operation: %s", err)
			}
		}
	}
	nTmp, err = s.Ext.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TransactionV0Ext: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionV0) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionV0) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionV0)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionV0)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TransactionV0) xdrType() {}

var _ xdrType = (*TransactionV0)(nil)

// TransactionV0Envelope is an XDR Struct defines as:
//
//   struct TransactionV0Envelope
//    {
//        TransactionV0 tx;
//        /* Each decorated signature is a signature over the SHA256 hash of
//         * a TransactionSignaturePayload */
//        DecoratedSignature signatures<20>;
//    };
//
type TransactionV0Envelope struct {
	Tx         TransactionV0
	Signatures []DecoratedSignature `xdrmaxsize:"20"`
}

// EncodeTo encodes this value using the Encoder.
func (s *TransactionV0Envelope) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Tx.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeUint(uint32(len(s.Signatures))); err != nil {
		return err
	}
	for i := 0; i < len(s.Signatures); i++ {
		if err = s.Signatures[i].EncodeTo(e); err != nil {
			return err
		}
	}
	return nil
}

var _ decoderFrom = (*TransactionV0Envelope)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *TransactionV0Envelope) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Tx.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TransactionV0: %s", err)
	}
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding DecoratedSignature: %s", err)
	}
	if l > 20 {
		return n, fmt.Errorf("decoding DecoratedSignature: data size (%d) exceeds size limit (20)", l)
	}
	s.Signatures = nil
	if l > 0 {
		s.Signatures = make([]DecoratedSignature, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.Signatures[i].DecodeFrom(d)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding DecoratedSignature: %s", err)
			}
		}
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionV0Envelope) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionV0Envelope) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionV0Envelope)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionV0Envelope)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TransactionV0Envelope) xdrType() {}

var _ xdrType = (*TransactionV0Envelope)(nil)

// TransactionExt is an XDR NestedUnion defines as:
//
//   union switch (int v)
//        {
//        case 0:
//            void;
//        }
//
type TransactionExt struct {
	V int32
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TransactionExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TransactionExt
func (u TransactionExt) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	}
	return "-", false
}

// NewTransactionExt creates a new  TransactionExt.
func NewTransactionExt(v int32, value interface{}) (result TransactionExt, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u TransactionExt) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union TransactionExt", u.V)
}

var _ decoderFrom = (*TransactionExt)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *TransactionExt) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %s", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	}
	return n, fmt.Errorf("union TransactionExt has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionExt) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionExt) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionExt)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionExt)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TransactionExt) xdrType() {}

var _ xdrType = (*TransactionExt)(nil)

// Transaction is an XDR Struct defines as:
//
//   struct Transaction
//    {
//        // account used to run the transaction
//        MuxedAccount sourceAccount;
//
//        // the fee the sourceAccount will pay
//        uint32 fee;
//
//        // sequence number to consume in the account
//        SequenceNumber seqNum;
//
//        // validity range (inclusive) for the last ledger close time
//        TimeBounds* timeBounds;
//
//        Memo memo;
//
//        Operation operations<MAX_OPS_PER_TX>;
//
//        // reserved for future use
//        union switch (int v)
//        {
//        case 0:
//            void;
//        }
//        ext;
//    };
//
type Transaction struct {
	SourceAccount MuxedAccount
	Fee           Uint32
	SeqNum        SequenceNumber
	TimeBounds    *TimeBounds
	Memo          Memo
	Operations    []Operation `xdrmaxsize:"100"`
	Ext           TransactionExt
}

// EncodeTo encodes this value using the Encoder.
func (s *Transaction) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.SourceAccount.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Fee.EncodeTo(e); err != nil {
		return err
	}
	if err = s.SeqNum.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeBool(s.TimeBounds != nil); err != nil {
		return err
	}
	if s.TimeBounds != nil {
		if err = (*s.TimeBounds).EncodeTo(e); err != nil {
			return err
		}
	}
	if err = s.Memo.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeUint(uint32(len(s.Operations))); err != nil {
		return err
	}
	for i := 0; i < len(s.Operations); i++ {
		if err = s.Operations[i].EncodeTo(e); err != nil {
			return err
		}
	}
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Transaction)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Transaction) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.SourceAccount.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding MuxedAccount: %s", err)
	}
	nTmp, err = s.Fee.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %s", err)
	}
	nTmp, err = s.SeqNum.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SequenceNumber: %s", err)
	}
	var b bool
	b, nTmp, err = d.DecodeBool()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TimeBounds: %s", err)
	}
	s.TimeBounds = nil
	if b {
		s.TimeBounds = new(TimeBounds)
		nTmp, err = s.TimeBounds.DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding TimeBounds: %s", err)
		}
	}
	nTmp, err = s.Memo.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Memo: %s", err)
	}
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Operation: %s", err)
	}
	if l > 100 {
		return n, fmt.Errorf("decoding Operation: data size (%d) exceeds size limit (100)", l)
	}
	s.Operations = nil
	if l > 0 {
		s.Operations = make([]Operation, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.Operations[i].DecodeFrom(d)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding Operation: %s", err)
			}
		}
	}
	nTmp, err = s.Ext.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TransactionExt: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Transaction) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Transaction) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Transaction)(nil)
	_ encoding.BinaryUnmarshaler = (*Transaction)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s Transaction) xdrType() {}

var _ xdrType = (*Transaction)(nil)

// TransactionV1Envelope is an XDR Struct defines as:
//
//   struct TransactionV1Envelope
//    {
//        Transaction tx;
//        /* Each decorated signature is a signature over the SHA256 hash of
//         * a TransactionSignaturePayload */
//        DecoratedSignature signatures<20>;
//    };
//
type TransactionV1Envelope struct {
	Tx         Transaction
	Signatures []DecoratedSignature `xdrmaxsize:"20"`
}

// EncodeTo encodes this value using the Encoder.
func (s *TransactionV1Envelope) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Tx.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeUint(uint32(len(s.Signatures))); err != nil {
		return err
	}
	for i := 0; i < len(s.Signatures); i++ {
		if err = s.Signatures[i].EncodeTo(e); err != nil {
			return err
		}
	}
	return nil
}

var _ decoderFrom = (*TransactionV1Envelope)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *TransactionV1Envelope) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Tx.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Transaction: %s", err)
	}
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding DecoratedSignature: %s", err)
	}
	if l > 20 {
		return n, fmt.Errorf("decoding DecoratedSignature: data size (%d) exceeds size limit (20)", l)
	}
	s.Signatures = nil
	if l > 0 {
		s.Signatures = make([]DecoratedSignature, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.Signatures[i].DecodeFrom(d)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding DecoratedSignature: %s", err)
			}
		}
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionV1Envelope) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionV1Envelope) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionV1Envelope)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionV1Envelope)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TransactionV1Envelope) xdrType() {}

var _ xdrType = (*TransactionV1Envelope)(nil)

// FeeBumpTransactionInnerTx is an XDR NestedUnion defines as:
//
//   union switch (EnvelopeType type)
//        {
//        case ENVELOPE_TYPE_TX:
//            TransactionV1Envelope v1;
//        }
//
type FeeBumpTransactionInnerTx struct {
	Type EnvelopeType
	V1   *TransactionV1Envelope
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u FeeBumpTransactionInnerTx) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of FeeBumpTransactionInnerTx
func (u FeeBumpTransactionInnerTx) ArmForSwitch(sw int32) (string, bool) {
	switch EnvelopeType(sw) {
	case EnvelopeTypeEnvelopeTypeTx:
		return "V1", true
	}
	return "-", false
}

// NewFeeBumpTransactionInnerTx creates a new  FeeBumpTransactionInnerTx.
func NewFeeBumpTransactionInnerTx(aType EnvelopeType, value interface{}) (result FeeBumpTransactionInnerTx, err error) {
	result.Type = aType
	switch EnvelopeType(aType) {
	case EnvelopeTypeEnvelopeTypeTx:
		tv, ok := value.(TransactionV1Envelope)
		if !ok {
			err = fmt.Errorf("invalid value, must be TransactionV1Envelope")
			return
		}
		result.V1 = &tv
	}
	return
}

// MustV1 retrieves the V1 value from the union,
// panicing if the value is not set.
func (u FeeBumpTransactionInnerTx) MustV1() TransactionV1Envelope {
	val, ok := u.GetV1()

	if !ok {
		panic("arm V1 is not set")
	}

	return val
}

// GetV1 retrieves the V1 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u FeeBumpTransactionInnerTx) GetV1() (result TransactionV1Envelope, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "V1" {
		result = *u.V1
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u FeeBumpTransactionInnerTx) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch EnvelopeType(u.Type) {
	case EnvelopeTypeEnvelopeTypeTx:
		if err = (*u.V1).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (EnvelopeType) switch value '%d' is not valid for union FeeBumpTransactionInnerTx", u.Type)
}

var _ decoderFrom = (*FeeBumpTransactionInnerTx)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *FeeBumpTransactionInnerTx) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding EnvelopeType: %s", err)
	}
	switch EnvelopeType(u.Type) {
	case EnvelopeTypeEnvelopeTypeTx:
		u.V1 = new(TransactionV1Envelope)
		nTmp, err = (*u.V1).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding TransactionV1Envelope: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union FeeBumpTransactionInnerTx has invalid Type (EnvelopeType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s FeeBumpTransactionInnerTx) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *FeeBumpTransactionInnerTx) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*FeeBumpTransactionInnerTx)(nil)
	_ encoding.BinaryUnmarshaler = (*FeeBumpTransactionInnerTx)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s FeeBumpTransactionInnerTx) xdrType() {}

var _ xdrType = (*FeeBumpTransactionInnerTx)(nil)

// FeeBumpTransactionExt is an XDR NestedUnion defines as:
//
//   union switch (int v)
//        {
//        case 0:
//            void;
//        }
//
type FeeBumpTransactionExt struct {
	V int32
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u FeeBumpTransactionExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of FeeBumpTransactionExt
func (u FeeBumpTransactionExt) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	}
	return "-", false
}

// NewFeeBumpTransactionExt creates a new  FeeBumpTransactionExt.
func NewFeeBumpTransactionExt(v int32, value interface{}) (result FeeBumpTransactionExt, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u FeeBumpTransactionExt) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union FeeBumpTransactionExt", u.V)
}

var _ decoderFrom = (*FeeBumpTransactionExt)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *FeeBumpTransactionExt) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %s", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	}
	return n, fmt.Errorf("union FeeBumpTransactionExt has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s FeeBumpTransactionExt) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *FeeBumpTransactionExt) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*FeeBumpTransactionExt)(nil)
	_ encoding.BinaryUnmarshaler = (*FeeBumpTransactionExt)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s FeeBumpTransactionExt) xdrType() {}

var _ xdrType = (*FeeBumpTransactionExt)(nil)

// FeeBumpTransaction is an XDR Struct defines as:
//
//   struct FeeBumpTransaction
//    {
//        MuxedAccount feeSource;
//        int64 fee;
//        union switch (EnvelopeType type)
//        {
//        case ENVELOPE_TYPE_TX:
//            TransactionV1Envelope v1;
//        }
//        innerTx;
//        union switch (int v)
//        {
//        case 0:
//            void;
//        }
//        ext;
//    };
//
type FeeBumpTransaction struct {
	FeeSource MuxedAccount
	Fee       Int64
	InnerTx   FeeBumpTransactionInnerTx
	Ext       FeeBumpTransactionExt
}

// EncodeTo encodes this value using the Encoder.
func (s *FeeBumpTransaction) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.FeeSource.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Fee.EncodeTo(e); err != nil {
		return err
	}
	if err = s.InnerTx.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*FeeBumpTransaction)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *FeeBumpTransaction) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.FeeSource.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding MuxedAccount: %s", err)
	}
	nTmp, err = s.Fee.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	nTmp, err = s.InnerTx.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding FeeBumpTransactionInnerTx: %s", err)
	}
	nTmp, err = s.Ext.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding FeeBumpTransactionExt: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s FeeBumpTransaction) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *FeeBumpTransaction) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*FeeBumpTransaction)(nil)
	_ encoding.BinaryUnmarshaler = (*FeeBumpTransaction)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s FeeBumpTransaction) xdrType() {}

var _ xdrType = (*FeeBumpTransaction)(nil)

// FeeBumpTransactionEnvelope is an XDR Struct defines as:
//
//   struct FeeBumpTransactionEnvelope
//    {
//        FeeBumpTransaction tx;
//        /* Each decorated signature is a signature over the SHA256 hash of
//         * a TransactionSignaturePayload */
//        DecoratedSignature signatures<20>;
//    };
//
type FeeBumpTransactionEnvelope struct {
	Tx         FeeBumpTransaction
	Signatures []DecoratedSignature `xdrmaxsize:"20"`
}

// EncodeTo encodes this value using the Encoder.
func (s *FeeBumpTransactionEnvelope) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Tx.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeUint(uint32(len(s.Signatures))); err != nil {
		return err
	}
	for i := 0; i < len(s.Signatures); i++ {
		if err = s.Signatures[i].EncodeTo(e); err != nil {
			return err
		}
	}
	return nil
}

var _ decoderFrom = (*FeeBumpTransactionEnvelope)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *FeeBumpTransactionEnvelope) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Tx.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding FeeBumpTransaction: %s", err)
	}
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding DecoratedSignature: %s", err)
	}
	if l > 20 {
		return n, fmt.Errorf("decoding DecoratedSignature: data size (%d) exceeds size limit (20)", l)
	}
	s.Signatures = nil
	if l > 0 {
		s.Signatures = make([]DecoratedSignature, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.Signatures[i].DecodeFrom(d)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding DecoratedSignature: %s", err)
			}
		}
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s FeeBumpTransactionEnvelope) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *FeeBumpTransactionEnvelope) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*FeeBumpTransactionEnvelope)(nil)
	_ encoding.BinaryUnmarshaler = (*FeeBumpTransactionEnvelope)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s FeeBumpTransactionEnvelope) xdrType() {}

var _ xdrType = (*FeeBumpTransactionEnvelope)(nil)

// TransactionEnvelope is an XDR Union defines as:
//
//   union TransactionEnvelope switch (EnvelopeType type)
//    {
//    case ENVELOPE_TYPE_TX_V0:
//        TransactionV0Envelope v0;
//    case ENVELOPE_TYPE_TX:
//        TransactionV1Envelope v1;
//    case ENVELOPE_TYPE_TX_FEE_BUMP:
//        FeeBumpTransactionEnvelope feeBump;
//    };
//
type TransactionEnvelope struct {
	Type    EnvelopeType
	V0      *TransactionV0Envelope
	V1      *TransactionV1Envelope
	FeeBump *FeeBumpTransactionEnvelope
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TransactionEnvelope) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TransactionEnvelope
func (u TransactionEnvelope) ArmForSwitch(sw int32) (string, bool) {
	switch EnvelopeType(sw) {
	case EnvelopeTypeEnvelopeTypeTxV0:
		return "V0", true
	case EnvelopeTypeEnvelopeTypeTx:
		return "V1", true
	case EnvelopeTypeEnvelopeTypeTxFeeBump:
		return "FeeBump", true
	}
	return "-", false
}

// NewTransactionEnvelope creates a new  TransactionEnvelope.
func NewTransactionEnvelope(aType EnvelopeType, value interface{}) (result TransactionEnvelope, err error) {
	result.Type = aType
	switch EnvelopeType(aType) {
	case EnvelopeTypeEnvelopeTypeTxV0:
		tv, ok := value.(TransactionV0Envelope)
		if !ok {
			err = fmt.Errorf("invalid value, must be TransactionV0Envelope")
			return
		}
		result.V0 = &tv
	case EnvelopeTypeEnvelopeTypeTx:
		tv, ok := value.(TransactionV1Envelope)
		if !ok {
			err = fmt.Errorf("invalid value, must be TransactionV1Envelope")
			return
		}
		result.V1 = &tv
	case EnvelopeTypeEnvelopeTypeTxFeeBump:
		tv, ok := value.(FeeBumpTransactionEnvelope)
		if !ok {
			err = fmt.Errorf("invalid value, must be FeeBumpTransactionEnvelope")
			return
		}
		result.FeeBump = &tv
	}
	return
}

// MustV0 retrieves the V0 value from the union,
// panicing if the value is not set.
func (u TransactionEnvelope) MustV0() TransactionV0Envelope {
	val, ok := u.GetV0()

	if !ok {
		panic("arm V0 is not set")
	}

	return val
}

// GetV0 retrieves the V0 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u TransactionEnvelope) GetV0() (result TransactionV0Envelope, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "V0" {
		result = *u.V0
		ok = true
	}

	return
}

// MustV1 retrieves the V1 value from the union,
// panicing if the value is not set.
func (u TransactionEnvelope) MustV1() TransactionV1Envelope {
	val, ok := u.GetV1()

	if !ok {
		panic("arm V1 is not set")
	}

	return val
}

// GetV1 retrieves the V1 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u TransactionEnvelope) GetV1() (result TransactionV1Envelope, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "V1" {
		result = *u.V1
		ok = true
	}

	return
}

// MustFeeBump retrieves the FeeBump value from the union,
// panicing if the value is not set.
func (u TransactionEnvelope) MustFeeBump() FeeBumpTransactionEnvelope {
	val, ok := u.GetFeeBump()

	if !ok {
		panic("arm FeeBump is not set")
	}

	return val
}

// GetFeeBump retrieves the FeeBump value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u TransactionEnvelope) GetFeeBump() (result FeeBumpTransactionEnvelope, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "FeeBump" {
		result = *u.FeeBump
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u TransactionEnvelope) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch EnvelopeType(u.Type) {
	case EnvelopeTypeEnvelopeTypeTxV0:
		if err = (*u.V0).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case EnvelopeTypeEnvelopeTypeTx:
		if err = (*u.V1).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case EnvelopeTypeEnvelopeTypeTxFeeBump:
		if err = (*u.FeeBump).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (EnvelopeType) switch value '%d' is not valid for union TransactionEnvelope", u.Type)
}

var _ decoderFrom = (*TransactionEnvelope)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *TransactionEnvelope) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding EnvelopeType: %s", err)
	}
	switch EnvelopeType(u.Type) {
	case EnvelopeTypeEnvelopeTypeTxV0:
		u.V0 = new(TransactionV0Envelope)
		nTmp, err = (*u.V0).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding TransactionV0Envelope: %s", err)
		}
		return n, nil
	case EnvelopeTypeEnvelopeTypeTx:
		u.V1 = new(TransactionV1Envelope)
		nTmp, err = (*u.V1).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding TransactionV1Envelope: %s", err)
		}
		return n, nil
	case EnvelopeTypeEnvelopeTypeTxFeeBump:
		u.FeeBump = new(FeeBumpTransactionEnvelope)
		nTmp, err = (*u.FeeBump).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding FeeBumpTransactionEnvelope: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union TransactionEnvelope has invalid Type (EnvelopeType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionEnvelope) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionEnvelope) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionEnvelope)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionEnvelope)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TransactionEnvelope) xdrType() {}

var _ xdrType = (*TransactionEnvelope)(nil)

// TransactionSignaturePayloadTaggedTransaction is an XDR NestedUnion defines as:
//
//   union switch (EnvelopeType type)
//        {
//        // Backwards Compatibility: Use ENVELOPE_TYPE_TX to sign ENVELOPE_TYPE_TX_V0
//        case ENVELOPE_TYPE_TX:
//            Transaction tx;
//        case ENVELOPE_TYPE_TX_FEE_BUMP:
//            FeeBumpTransaction feeBump;
//        }
//
type TransactionSignaturePayloadTaggedTransaction struct {
	Type    EnvelopeType
	Tx      *Transaction
	FeeBump *FeeBumpTransaction
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TransactionSignaturePayloadTaggedTransaction) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TransactionSignaturePayloadTaggedTransaction
func (u TransactionSignaturePayloadTaggedTransaction) ArmForSwitch(sw int32) (string, bool) {
	switch EnvelopeType(sw) {
	case EnvelopeTypeEnvelopeTypeTx:
		return "Tx", true
	case EnvelopeTypeEnvelopeTypeTxFeeBump:
		return "FeeBump", true
	}
	return "-", false
}

// NewTransactionSignaturePayloadTaggedTransaction creates a new  TransactionSignaturePayloadTaggedTransaction.
func NewTransactionSignaturePayloadTaggedTransaction(aType EnvelopeType, value interface{}) (result TransactionSignaturePayloadTaggedTransaction, err error) {
	result.Type = aType
	switch EnvelopeType(aType) {
	case EnvelopeTypeEnvelopeTypeTx:
		tv, ok := value.(Transaction)
		if !ok {
			err = fmt.Errorf("invalid value, must be Transaction")
			return
		}
		result.Tx = &tv
	case EnvelopeTypeEnvelopeTypeTxFeeBump:
		tv, ok := value.(FeeBumpTransaction)
		if !ok {
			err = fmt.Errorf("invalid value, must be FeeBumpTransaction")
			return
		}
		result.FeeBump = &tv
	}
	return
}

// MustTx retrieves the Tx value from the union,
// panicing if the value is not set.
func (u TransactionSignaturePayloadTaggedTransaction) MustTx() Transaction {
	val, ok := u.GetTx()

	if !ok {
		panic("arm Tx is not set")
	}

	return val
}

// GetTx retrieves the Tx value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u TransactionSignaturePayloadTaggedTransaction) GetTx() (result Transaction, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Tx" {
		result = *u.Tx
		ok = true
	}

	return
}

// MustFeeBump retrieves the FeeBump value from the union,
// panicing if the value is not set.
func (u TransactionSignaturePayloadTaggedTransaction) MustFeeBump() FeeBumpTransaction {
	val, ok := u.GetFeeBump()

	if !ok {
		panic("arm FeeBump is not set")
	}

	return val
}

// GetFeeBump retrieves the FeeBump value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u TransactionSignaturePayloadTaggedTransaction) GetFeeBump() (result FeeBumpTransaction, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "FeeBump" {
		result = *u.FeeBump
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u TransactionSignaturePayloadTaggedTransaction) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch EnvelopeType(u.Type) {
	case EnvelopeTypeEnvelopeTypeTx:
		if err = (*u.Tx).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case EnvelopeTypeEnvelopeTypeTxFeeBump:
		if err = (*u.FeeBump).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (EnvelopeType) switch value '%d' is not valid for union TransactionSignaturePayloadTaggedTransaction", u.Type)
}

var _ decoderFrom = (*TransactionSignaturePayloadTaggedTransaction)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *TransactionSignaturePayloadTaggedTransaction) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding EnvelopeType: %s", err)
	}
	switch EnvelopeType(u.Type) {
	case EnvelopeTypeEnvelopeTypeTx:
		u.Tx = new(Transaction)
		nTmp, err = (*u.Tx).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Transaction: %s", err)
		}
		return n, nil
	case EnvelopeTypeEnvelopeTypeTxFeeBump:
		u.FeeBump = new(FeeBumpTransaction)
		nTmp, err = (*u.FeeBump).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding FeeBumpTransaction: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union TransactionSignaturePayloadTaggedTransaction has invalid Type (EnvelopeType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionSignaturePayloadTaggedTransaction) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionSignaturePayloadTaggedTransaction) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionSignaturePayloadTaggedTransaction)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionSignaturePayloadTaggedTransaction)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TransactionSignaturePayloadTaggedTransaction) xdrType() {}

var _ xdrType = (*TransactionSignaturePayloadTaggedTransaction)(nil)

// TransactionSignaturePayload is an XDR Struct defines as:
//
//   struct TransactionSignaturePayload
//    {
//        Hash networkId;
//        union switch (EnvelopeType type)
//        {
//        // Backwards Compatibility: Use ENVELOPE_TYPE_TX to sign ENVELOPE_TYPE_TX_V0
//        case ENVELOPE_TYPE_TX:
//            Transaction tx;
//        case ENVELOPE_TYPE_TX_FEE_BUMP:
//            FeeBumpTransaction feeBump;
//        }
//        taggedTransaction;
//    };
//
type TransactionSignaturePayload struct {
	NetworkId         Hash
	TaggedTransaction TransactionSignaturePayloadTaggedTransaction
}

// EncodeTo encodes this value using the Encoder.
func (s *TransactionSignaturePayload) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.NetworkId.EncodeTo(e); err != nil {
		return err
	}
	if err = s.TaggedTransaction.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*TransactionSignaturePayload)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *TransactionSignaturePayload) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.NetworkId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Hash: %s", err)
	}
	nTmp, err = s.TaggedTransaction.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TransactionSignaturePayloadTaggedTransaction: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionSignaturePayload) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionSignaturePayload) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionSignaturePayload)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionSignaturePayload)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TransactionSignaturePayload) xdrType() {}

var _ xdrType = (*TransactionSignaturePayload)(nil)

// ClaimAtomType is an XDR Enum defines as:
//
//   enum ClaimAtomType
//    {
//        CLAIM_ATOM_TYPE_V0 = 0,
//        CLAIM_ATOM_TYPE_ORDER_BOOK = 1,
//        CLAIM_ATOM_TYPE_LIQUIDITY_POOL = 2
//    };
//
type ClaimAtomType int32

const (
	ClaimAtomTypeClaimAtomTypeV0            ClaimAtomType = 0
	ClaimAtomTypeClaimAtomTypeOrderBook     ClaimAtomType = 1
	ClaimAtomTypeClaimAtomTypeLiquidityPool ClaimAtomType = 2
)

var claimAtomTypeMap = map[int32]string{
	0: "ClaimAtomTypeClaimAtomTypeV0",
	1: "ClaimAtomTypeClaimAtomTypeOrderBook",
	2: "ClaimAtomTypeClaimAtomTypeLiquidityPool",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ClaimAtomType
func (e ClaimAtomType) ValidEnum(v int32) bool {
	_, ok := claimAtomTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e ClaimAtomType) String() string {
	name, _ := claimAtomTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e ClaimAtomType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := claimAtomTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid ClaimAtomType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*ClaimAtomType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *ClaimAtomType) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding ClaimAtomType: %s", err)
	}
	if _, ok := claimAtomTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid ClaimAtomType enum value", v)
	}
	*e = ClaimAtomType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ClaimAtomType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ClaimAtomType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ClaimAtomType)(nil)
	_ encoding.BinaryUnmarshaler = (*ClaimAtomType)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ClaimAtomType) xdrType() {}

var _ xdrType = (*ClaimAtomType)(nil)

// ClaimOfferAtomV0 is an XDR Struct defines as:
//
//   struct ClaimOfferAtomV0
//    {
//        // emitted to identify the offer
//        uint256 sellerEd25519; // Account that owns the offer
//        int64 offerID;
//
//        // amount and asset taken from the owner
//        Asset assetSold;
//        int64 amountSold;
//
//        // amount and asset sent to the owner
//        Asset assetBought;
//        int64 amountBought;
//    };
//
type ClaimOfferAtomV0 struct {
	SellerEd25519 Uint256
	OfferId       Int64
	AssetSold     Asset
	AmountSold    Int64
	AssetBought   Asset
	AmountBought  Int64
}

// EncodeTo encodes this value using the Encoder.
func (s *ClaimOfferAtomV0) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.SellerEd25519.EncodeTo(e); err != nil {
		return err
	}
	if err = s.OfferId.EncodeTo(e); err != nil {
		return err
	}
	if err = s.AssetSold.EncodeTo(e); err != nil {
		return err
	}
	if err = s.AmountSold.EncodeTo(e); err != nil {
		return err
	}
	if err = s.AssetBought.EncodeTo(e); err != nil {
		return err
	}
	if err = s.AmountBought.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*ClaimOfferAtomV0)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *ClaimOfferAtomV0) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.SellerEd25519.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint256: %s", err)
	}
	nTmp, err = s.OfferId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	nTmp, err = s.AssetSold.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Asset: %s", err)
	}
	nTmp, err = s.AmountSold.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	nTmp, err = s.AssetBought.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Asset: %s", err)
	}
	nTmp, err = s.AmountBought.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ClaimOfferAtomV0) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ClaimOfferAtomV0) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ClaimOfferAtomV0)(nil)
	_ encoding.BinaryUnmarshaler = (*ClaimOfferAtomV0)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ClaimOfferAtomV0) xdrType() {}

var _ xdrType = (*ClaimOfferAtomV0)(nil)

// ClaimOfferAtom is an XDR Struct defines as:
//
//   struct ClaimOfferAtom
//    {
//        // emitted to identify the offer
//        AccountID sellerID; // Account that owns the offer
//        int64 offerID;
//
//        // amount and asset taken from the owner
//        Asset assetSold;
//        int64 amountSold;
//
//        // amount and asset sent to the owner
//        Asset assetBought;
//        int64 amountBought;
//    };
//
type ClaimOfferAtom struct {
	SellerId     AccountId
	OfferId      Int64
	AssetSold    Asset
	AmountSold   Int64
	AssetBought  Asset
	AmountBought Int64
}

// EncodeTo encodes this value using the Encoder.
func (s *ClaimOfferAtom) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.SellerId.EncodeTo(e); err != nil {
		return err
	}
	if err = s.OfferId.EncodeTo(e); err != nil {
		return err
	}
	if err = s.AssetSold.EncodeTo(e); err != nil {
		return err
	}
	if err = s.AmountSold.EncodeTo(e); err != nil {
		return err
	}
	if err = s.AssetBought.EncodeTo(e); err != nil {
		return err
	}
	if err = s.AmountBought.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*ClaimOfferAtom)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *ClaimOfferAtom) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.SellerId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountId: %s", err)
	}
	nTmp, err = s.OfferId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	nTmp, err = s.AssetSold.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Asset: %s", err)
	}
	nTmp, err = s.AmountSold.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	nTmp, err = s.AssetBought.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Asset: %s", err)
	}
	nTmp, err = s.AmountBought.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ClaimOfferAtom) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ClaimOfferAtom) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ClaimOfferAtom)(nil)
	_ encoding.BinaryUnmarshaler = (*ClaimOfferAtom)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ClaimOfferAtom) xdrType() {}

var _ xdrType = (*ClaimOfferAtom)(nil)

// ClaimLiquidityAtom is an XDR Struct defines as:
//
//   struct ClaimLiquidityAtom
//    {
//        PoolID liquidityPoolID;
//
//        // amount and asset taken from the pool
//        Asset assetSold;
//        int64 amountSold;
//
//        // amount and asset sent to the pool
//        Asset assetBought;
//        int64 amountBought;
//    };
//
type ClaimLiquidityAtom struct {
	LiquidityPoolId PoolId
	AssetSold       Asset
	AmountSold      Int64
	AssetBought     Asset
	AmountBought    Int64
}

// EncodeTo encodes this value using the Encoder.
func (s *ClaimLiquidityAtom) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.LiquidityPoolId.EncodeTo(e); err != nil {
		return err
	}
	if err = s.AssetSold.EncodeTo(e); err != nil {
		return err
	}
	if err = s.AmountSold.EncodeTo(e); err != nil {
		return err
	}
	if err = s.AssetBought.EncodeTo(e); err != nil {
		return err
	}
	if err = s.AmountBought.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*ClaimLiquidityAtom)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *ClaimLiquidityAtom) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.LiquidityPoolId.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding PoolId: %s", err)
	}
	nTmp, err = s.AssetSold.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Asset: %s", err)
	}
	nTmp, err = s.AmountSold.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	nTmp, err = s.AssetBought.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Asset: %s", err)
	}
	nTmp, err = s.AmountBought.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ClaimLiquidityAtom) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ClaimLiquidityAtom) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ClaimLiquidityAtom)(nil)
	_ encoding.BinaryUnmarshaler = (*ClaimLiquidityAtom)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ClaimLiquidityAtom) xdrType() {}

var _ xdrType = (*ClaimLiquidityAtom)(nil)

// ClaimAtom is an XDR Union defines as:
//
//   union ClaimAtom switch (ClaimAtomType type)
//    {
//    case CLAIM_ATOM_TYPE_V0:
//        ClaimOfferAtomV0 v0;
//    case CLAIM_ATOM_TYPE_ORDER_BOOK:
//        ClaimOfferAtom orderBook;
//    case CLAIM_ATOM_TYPE_LIQUIDITY_POOL:
//        ClaimLiquidityAtom liquidityPool;
//    };
//
type ClaimAtom struct {
	Type          ClaimAtomType
	V0            *ClaimOfferAtomV0
	OrderBook     *ClaimOfferAtom
	LiquidityPool *ClaimLiquidityAtom
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ClaimAtom) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ClaimAtom
func (u ClaimAtom) ArmForSwitch(sw int32) (string, bool) {
	switch ClaimAtomType(sw) {
	case ClaimAtomTypeClaimAtomTypeV0:
		return "V0", true
	case ClaimAtomTypeClaimAtomTypeOrderBook:
		return "OrderBook", true
	case ClaimAtomTypeClaimAtomTypeLiquidityPool:
		return "LiquidityPool", true
	}
	return "-", false
}

// NewClaimAtom creates a new  ClaimAtom.
func NewClaimAtom(aType ClaimAtomType, value interface{}) (result ClaimAtom, err error) {
	result.Type = aType
	switch ClaimAtomType(aType) {
	case ClaimAtomTypeClaimAtomTypeV0:
		tv, ok := value.(ClaimOfferAtomV0)
		if !ok {
			err = fmt.Errorf("invalid value, must be ClaimOfferAtomV0")
			return
		}
		result.V0 = &tv
	case ClaimAtomTypeClaimAtomTypeOrderBook:
		tv, ok := value.(ClaimOfferAtom)
		if !ok {
			err = fmt.Errorf("invalid value, must be ClaimOfferAtom")
			return
		}
		result.OrderBook = &tv
	case ClaimAtomTypeClaimAtomTypeLiquidityPool:
		tv, ok := value.(ClaimLiquidityAtom)
		if !ok {
			err = fmt.Errorf("invalid value, must be ClaimLiquidityAtom")
			return
		}
		result.LiquidityPool = &tv
	}
	return
}

// MustV0 retrieves the V0 value from the union,
// panicing if the value is not set.
func (u ClaimAtom) MustV0() ClaimOfferAtomV0 {
	val, ok := u.GetV0()

	if !ok {
		panic("arm V0 is not set")
	}

	return val
}

// GetV0 retrieves the V0 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ClaimAtom) GetV0() (result ClaimOfferAtomV0, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "V0" {
		result = *u.V0
		ok = true
	}

	return
}

// MustOrderBook retrieves the OrderBook value from the union,
// panicing if the value is not set.
func (u ClaimAtom) MustOrderBook() ClaimOfferAtom {
	val, ok := u.GetOrderBook()

	if !ok {
		panic("arm OrderBook is not set")
	}

	return val
}

// GetOrderBook retrieves the OrderBook value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ClaimAtom) GetOrderBook() (result ClaimOfferAtom, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "OrderBook" {
		result = *u.OrderBook
		ok = true
	}

	return
}

// MustLiquidityPool retrieves the LiquidityPool value from the union,
// panicing if the value is not set.
func (u ClaimAtom) MustLiquidityPool() ClaimLiquidityAtom {
	val, ok := u.GetLiquidityPool()

	if !ok {
		panic("arm LiquidityPool is not set")
	}

	return val
}

// GetLiquidityPool retrieves the LiquidityPool value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ClaimAtom) GetLiquidityPool() (result ClaimLiquidityAtom, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "LiquidityPool" {
		result = *u.LiquidityPool
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u ClaimAtom) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch ClaimAtomType(u.Type) {
	case ClaimAtomTypeClaimAtomTypeV0:
		if err = (*u.V0).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case ClaimAtomTypeClaimAtomTypeOrderBook:
		if err = (*u.OrderBook).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case ClaimAtomTypeClaimAtomTypeLiquidityPool:
		if err = (*u.LiquidityPool).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (ClaimAtomType) switch value '%d' is not valid for union ClaimAtom", u.Type)
}

var _ decoderFrom = (*ClaimAtom)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *ClaimAtom) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ClaimAtomType: %s", err)
	}
	switch ClaimAtomType(u.Type) {
	case ClaimAtomTypeClaimAtomTypeV0:
		u.V0 = new(ClaimOfferAtomV0)
		nTmp, err = (*u.V0).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ClaimOfferAtomV0: %s", err)
		}
		return n, nil
	case ClaimAtomTypeClaimAtomTypeOrderBook:
		u.OrderBook = new(ClaimOfferAtom)
		nTmp, err = (*u.OrderBook).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ClaimOfferAtom: %s", err)
		}
		return n, nil
	case ClaimAtomTypeClaimAtomTypeLiquidityPool:
		u.LiquidityPool = new(ClaimLiquidityAtom)
		nTmp, err = (*u.LiquidityPool).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ClaimLiquidityAtom: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union ClaimAtom has invalid Type (ClaimAtomType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ClaimAtom) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ClaimAtom) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ClaimAtom)(nil)
	_ encoding.BinaryUnmarshaler = (*ClaimAtom)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ClaimAtom) xdrType() {}

var _ xdrType = (*ClaimAtom)(nil)

// CreateAccountResultCode is an XDR Enum defines as:
//
//   enum CreateAccountResultCode
//    {
//        // codes considered as "success" for the operation
//        CREATE_ACCOUNT_SUCCESS = 0, // account was created
//
//        // codes considered as "failure" for the operation
//        CREATE_ACCOUNT_MALFORMED = -1,   // invalid destination
//        CREATE_ACCOUNT_UNDERFUNDED = -2, // not enough funds in source account
//        CREATE_ACCOUNT_LOW_RESERVE =
//            -3, // would create an account below the min reserve
//        CREATE_ACCOUNT_ALREADY_EXIST = -4 // account already exists
//    };
//
type CreateAccountResultCode int32

const (
	CreateAccountResultCodeCreateAccountSuccess      CreateAccountResultCode = 0
	CreateAccountResultCodeCreateAccountMalformed    CreateAccountResultCode = -1
	CreateAccountResultCodeCreateAccountUnderfunded  CreateAccountResultCode = -2
	CreateAccountResultCodeCreateAccountLowReserve   CreateAccountResultCode = -3
	CreateAccountResultCodeCreateAccountAlreadyExist CreateAccountResultCode = -4
)

var createAccountResultCodeMap = map[int32]string{
	0:  "CreateAccountResultCodeCreateAccountSuccess",
	-1: "CreateAccountResultCodeCreateAccountMalformed",
	-2: "CreateAccountResultCodeCreateAccountUnderfunded",
	-3: "CreateAccountResultCodeCreateAccountLowReserve",
	-4: "CreateAccountResultCodeCreateAccountAlreadyExist",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for CreateAccountResultCode
func (e CreateAccountResultCode) ValidEnum(v int32) bool {
	_, ok := createAccountResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e CreateAccountResultCode) String() string {
	name, _ := createAccountResultCodeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e CreateAccountResultCode) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := createAccountResultCodeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid CreateAccountResultCode enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*CreateAccountResultCode)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *CreateAccountResultCode) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding CreateAccountResultCode: %s", err)
	}
	if _, ok := createAccountResultCodeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid CreateAccountResultCode enum value", v)
	}
	*e = CreateAccountResultCode(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s CreateAccountResultCode) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *CreateAccountResultCode) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*CreateAccountResultCode)(nil)
	_ encoding.BinaryUnmarshaler = (*CreateAccountResultCode)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s CreateAccountResultCode) xdrType() {}

var _ xdrType = (*CreateAccountResultCode)(nil)

// CreateAccountResult is an XDR Union defines as:
//
//   union CreateAccountResult switch (CreateAccountResultCode code)
//    {
//    case CREATE_ACCOUNT_SUCCESS:
//        void;
//    default:
//        void;
//    };
//
type CreateAccountResult struct {
	Code CreateAccountResultCode
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u CreateAccountResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of CreateAccountResult
func (u CreateAccountResult) ArmForSwitch(sw int32) (string, bool) {
	switch CreateAccountResultCode(sw) {
	case CreateAccountResultCodeCreateAccountSuccess:
		return "", true
	default:
		return "", true
	}
}

// NewCreateAccountResult creates a new  CreateAccountResult.
func NewCreateAccountResult(code CreateAccountResultCode, value interface{}) (result CreateAccountResult, err error) {
	result.Code = code
	switch CreateAccountResultCode(code) {
	case CreateAccountResultCodeCreateAccountSuccess:
		// void
	default:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u CreateAccountResult) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Code.EncodeTo(e); err != nil {
		return err
	}
	switch CreateAccountResultCode(u.Code) {
	case CreateAccountResultCodeCreateAccountSuccess:
		// Void
		return nil
	default:
		// Void
		return nil
	}
}

var _ decoderFrom = (*CreateAccountResult)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *CreateAccountResult) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Code.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding CreateAccountResultCode: %s", err)
	}
	switch CreateAccountResultCode(u.Code) {
	case CreateAccountResultCodeCreateAccountSuccess:
		// Void
		return n, nil
	default:
		// Void
		return n, nil
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s CreateAccountResult) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *CreateAccountResult) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*CreateAccountResult)(nil)
	_ encoding.BinaryUnmarshaler = (*CreateAccountResult)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s CreateAccountResult) xdrType() {}

var _ xdrType = (*CreateAccountResult)(nil)

// PaymentResultCode is an XDR Enum defines as:
//
//   enum PaymentResultCode
//    {
//        // codes considered as "success" for the operation
//        PAYMENT_SUCCESS = 0, // payment successfully completed
//
//        // codes considered as "failure" for the operation
//        PAYMENT_MALFORMED = -1,          // bad input
//        PAYMENT_UNDERFUNDED = -2,        // not enough funds in source account
//        PAYMENT_SRC_NO_TRUST = -3,       // no trust line on source account
//        PAYMENT_SRC_NOT_AUTHORIZED = -4, // source not authorized to transfer
//        PAYMENT_NO_DESTINATION = -5,     // destination account does not exist
//        PAYMENT_NO_TRUST = -6,       // destination missing a trust line for asset
//        PAYMENT_NOT_AUTHORIZED = -7, // destination not authorized to hold asset
//        PAYMENT_LINE_FULL = -8,      // destination would go above their limit
//        PAYMENT_NO_ISSUER = -9       // missing issuer on asset
//    };
//
type PaymentResultCode int32

const (
	PaymentResultCodePaymentSuccess          PaymentResultCode = 0
	PaymentResultCodePaymentMalformed        PaymentResultCode = -1
	PaymentResultCodePaymentUnderfunded      PaymentResultCode = -2
	PaymentResultCodePaymentSrcNoTrust       PaymentResultCode = -3
	PaymentResultCodePaymentSrcNotAuthorized PaymentResultCode = -4
	PaymentResultCodePaymentNoDestination    PaymentResultCode = -5
	PaymentResultCodePaymentNoTrust          PaymentResultCode = -6
	PaymentResultCodePaymentNotAuthorized    PaymentResultCode = -7
	PaymentResultCodePaymentLineFull         PaymentResultCode = -8
	PaymentResultCodePaymentNoIssuer         PaymentResultCode = -9
)

var paymentResultCodeMap = map[int32]string{
	0:  "PaymentResultCodePaymentSuccess",
	-1: "PaymentResultCodePaymentMalformed",
	-2: "PaymentResultCodePaymentUnderfunded",
	-3: "PaymentResultCodePaymentSrcNoTrust",
	-4: "PaymentResultCodePaymentSrcNotAuthorized",
	-5: "PaymentResultCodePaymentNoDestination",
	-6: "PaymentResultCodePaymentNoTrust",
	-7: "PaymentResultCodePaymentNotAuthorized",
	-8: "PaymentResultCodePaymentLineFull",
	-9: "PaymentResultCodePaymentNoIssuer",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for PaymentResultCode
func (e PaymentResultCode) ValidEnum(v int32) bool {
	_, ok := paymentResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e PaymentResultCode) String() string {
	name, _ := paymentResultCodeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e PaymentResultCode) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := paymentResultCodeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid PaymentResultCode enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*PaymentResultCode)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *PaymentResultCode) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding PaymentResultCode: %s", err)
	}
	if _, ok := paymentResultCodeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid PaymentResultCode enum value", v)
	}
	*e = PaymentResultCode(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s PaymentResultCode) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *PaymentResultCode) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*PaymentResultCode)(nil)
	_ encoding.BinaryUnmarshaler = (*PaymentResultCode)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s PaymentResultCode) xdrType() {}

var _ xdrType = (*PaymentResultCode)(nil)

// PaymentResult is an XDR Union defines as:
//
//   union PaymentResult switch (PaymentResultCode code)
//    {
//    case PAYMENT_SUCCESS:
//        void;
//    default:
//        void;
//    };
//
type PaymentResult struct {
	Code PaymentResultCode
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u PaymentResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PaymentResult
func (u PaymentResult) ArmForSwitch(sw int32) (string, bool) {
	switch PaymentResultCode(sw) {
	case PaymentResultCodePaymentSuccess:
		return "", true
	default:
		return "", true
	}
}

// NewPaymentResult creates a new  PaymentResult.
func NewPaymentResult(code PaymentResultCode, value interface{}) (result PaymentResult, err error) {
	result.Code = code
	switch PaymentResultCode(code) {
	case PaymentResultCodePaymentSuccess:
		// void
	default:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u PaymentResult) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Code.EncodeTo(e); err != nil {
		return err
	}
	switch PaymentResultCode(u.Code) {
	case PaymentResultCodePaymentSuccess:
		// Void
		return nil
	default:
		// Void
		return nil
	}
}

var _ decoderFrom = (*PaymentResult)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *PaymentResult) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Code.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding PaymentResultCode: %s", err)
	}
	switch PaymentResultCode(u.Code) {
	case PaymentResultCodePaymentSuccess:
		// Void
		return n, nil
	default:
		// Void
		return n, nil
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s PaymentResult) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *PaymentResult) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*PaymentResult)(nil)
	_ encoding.BinaryUnmarshaler = (*PaymentResult)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s PaymentResult) xdrType() {}

var _ xdrType = (*PaymentResult)(nil)

// PathPaymentStrictReceiveResultCode is an XDR Enum defines as:
//
//   enum PathPaymentStrictReceiveResultCode
//    {
//        // codes considered as "success" for the operation
//        PATH_PAYMENT_STRICT_RECEIVE_SUCCESS = 0, // success
//
//        // codes considered as "failure" for the operation
//        PATH_PAYMENT_STRICT_RECEIVE_MALFORMED = -1, // bad input
//        PATH_PAYMENT_STRICT_RECEIVE_UNDERFUNDED =
//            -2, // not enough funds in source account
//        PATH_PAYMENT_STRICT_RECEIVE_SRC_NO_TRUST =
//            -3, // no trust line on source account
//        PATH_PAYMENT_STRICT_RECEIVE_SRC_NOT_AUTHORIZED =
//            -4, // source not authorized to transfer
//        PATH_PAYMENT_STRICT_RECEIVE_NO_DESTINATION =
//            -5, // destination account does not exist
//        PATH_PAYMENT_STRICT_RECEIVE_NO_TRUST =
//            -6, // dest missing a trust line for asset
//        PATH_PAYMENT_STRICT_RECEIVE_NOT_AUTHORIZED =
//            -7, // dest not authorized to hold asset
//        PATH_PAYMENT_STRICT_RECEIVE_LINE_FULL =
//            -8, // dest would go above their limit
//        PATH_PAYMENT_STRICT_RECEIVE_NO_ISSUER = -9, // missing issuer on one asset
//        PATH_PAYMENT_STRICT_RECEIVE_TOO_FEW_OFFERS =
//            -10, // not enough offers to satisfy path
//        PATH_PAYMENT_STRICT_RECEIVE_OFFER_CROSS_SELF =
//            -11, // would cross one of its own offers
//        PATH_PAYMENT_STRICT_RECEIVE_OVER_SENDMAX = -12 // could not satisfy sendmax
//    };
//
type PathPaymentStrictReceiveResultCode int32

const (
	PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveSuccess          PathPaymentStrictReceiveResultCode = 0
	PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveMalformed        PathPaymentStrictReceiveResultCode = -1
	PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveUnderfunded      PathPaymentStrictReceiveResultCode = -2
	PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveSrcNoTrust       PathPaymentStrictReceiveResultCode = -3
	PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveSrcNotAuthorized PathPaymentStrictReceiveResultCode = -4
	PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveNoDestination    PathPaymentStrictReceiveResultCode = -5
	PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveNoTrust          PathPaymentStrictReceiveResultCode = -6
	PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveNotAuthorized    PathPaymentStrictReceiveResultCode = -7
	PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveLineFull         PathPaymentStrictReceiveResultCode = -8
	PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveNoIssuer         PathPaymentStrictReceiveResultCode = -9
	PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveTooFewOffers     PathPaymentStrictReceiveResultCode = -10
	PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveOfferCrossSelf   PathPaymentStrictReceiveResultCode = -11
	PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveOverSendmax      PathPaymentStrictReceiveResultCode = -12
)

var pathPaymentStrictReceiveResultCodeMap = map[int32]string{
	0:   "PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveSuccess",
	-1:  "PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveMalformed",
	-2:  "PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveUnderfunded",
	-3:  "PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveSrcNoTrust",
	-4:  "PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveSrcNotAuthorized",
	-5:  "PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveNoDestination",
	-6:  "PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveNoTrust",
	-7:  "PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveNotAuthorized",
	-8:  "PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveLineFull",
	-9:  "PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveNoIssuer",
	-10: "PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveTooFewOffers",
	-11: "PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveOfferCrossSelf",
	-12: "PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveOverSendmax",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for PathPaymentStrictReceiveResultCode
func (e PathPaymentStrictReceiveResultCode) ValidEnum(v int32) bool {
	_, ok := pathPaymentStrictReceiveResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e PathPaymentStrictReceiveResultCode) String() string {
	name, _ := pathPaymentStrictReceiveResultCodeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e PathPaymentStrictReceiveResultCode) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := pathPaymentStrictReceiveResultCodeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid PathPaymentStrictReceiveResultCode enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*PathPaymentStrictReceiveResultCode)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *PathPaymentStrictReceiveResultCode) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding PathPaymentStrictReceiveResultCode: %s", err)
	}
	if _, ok := pathPaymentStrictReceiveResultCodeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid PathPaymentStrictReceiveResultCode enum value", v)
	}
	*e = PathPaymentStrictReceiveResultCode(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s PathPaymentStrictReceiveResultCode) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *PathPaymentStrictReceiveResultCode) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*PathPaymentStrictReceiveResultCode)(nil)
	_ encoding.BinaryUnmarshaler = (*PathPaymentStrictReceiveResultCode)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s PathPaymentStrictReceiveResultCode) xdrType() {}

var _ xdrType = (*PathPaymentStrictReceiveResultCode)(nil)

// SimplePaymentResult is an XDR Struct defines as:
//
//   struct SimplePaymentResult
//    {
//        AccountID destination;
//        Asset asset;
//        int64 amount;
//    };
//
type SimplePaymentResult struct {
	Destination AccountId
	Asset       Asset
	Amount      Int64
}

// EncodeTo encodes this value using the Encoder.
func (s *SimplePaymentResult) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Destination.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Asset.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Amount.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*SimplePaymentResult)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *SimplePaymentResult) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Destination.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountId: %s", err)
	}
	nTmp, err = s.Asset.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Asset: %s", err)
	}
	nTmp, err = s.Amount.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s SimplePaymentResult) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SimplePaymentResult) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*SimplePaymentResult)(nil)
	_ encoding.BinaryUnmarshaler = (*SimplePaymentResult)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s SimplePaymentResult) xdrType() {}

var _ xdrType = (*SimplePaymentResult)(nil)

// PathPaymentStrictReceiveResultSuccess is an XDR NestedStruct defines as:
//
//   struct
//        {
//            ClaimAtom offers<>;
//            SimplePaymentResult last;
//        }
//
type PathPaymentStrictReceiveResultSuccess struct {
	Offers []ClaimAtom
	Last   SimplePaymentResult
}

// EncodeTo encodes this value using the Encoder.
func (s *PathPaymentStrictReceiveResultSuccess) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeUint(uint32(len(s.Offers))); err != nil {
		return err
	}
	for i := 0; i < len(s.Offers); i++ {
		if err = s.Offers[i].EncodeTo(e); err != nil {
			return err
		}
	}
	if err = s.Last.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*PathPaymentStrictReceiveResultSuccess)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *PathPaymentStrictReceiveResultSuccess) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ClaimAtom: %s", err)
	}
	s.Offers = nil
	if l > 0 {
		s.Offers = make([]ClaimAtom, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.Offers[i].DecodeFrom(d)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding ClaimAtom: %s", err)
			}
		}
	}
	nTmp, err = s.Last.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SimplePaymentResult: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s PathPaymentStrictReceiveResultSuccess) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *PathPaymentStrictReceiveResultSuccess) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*PathPaymentStrictReceiveResultSuccess)(nil)
	_ encoding.BinaryUnmarshaler = (*PathPaymentStrictReceiveResultSuccess)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s PathPaymentStrictReceiveResultSuccess) xdrType() {}

var _ xdrType = (*PathPaymentStrictReceiveResultSuccess)(nil)

// PathPaymentStrictReceiveResult is an XDR Union defines as:
//
//   union PathPaymentStrictReceiveResult switch (
//        PathPaymentStrictReceiveResultCode code)
//    {
//    case PATH_PAYMENT_STRICT_RECEIVE_SUCCESS:
//        struct
//        {
//            ClaimAtom offers<>;
//            SimplePaymentResult last;
//        } success;
//    case PATH_PAYMENT_STRICT_RECEIVE_NO_ISSUER:
//        Asset noIssuer; // the asset that caused the error
//    default:
//        void;
//    };
//
type PathPaymentStrictReceiveResult struct {
	Code     PathPaymentStrictReceiveResultCode
	Success  *PathPaymentStrictReceiveResultSuccess
	NoIssuer *Asset
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u PathPaymentStrictReceiveResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PathPaymentStrictReceiveResult
func (u PathPaymentStrictReceiveResult) ArmForSwitch(sw int32) (string, bool) {
	switch PathPaymentStrictReceiveResultCode(sw) {
	case PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveSuccess:
		return "Success", true
	case PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveNoIssuer:
		return "NoIssuer", true
	default:
		return "", true
	}
}

// NewPathPaymentStrictReceiveResult creates a new  PathPaymentStrictReceiveResult.
func NewPathPaymentStrictReceiveResult(code PathPaymentStrictReceiveResultCode, value interface{}) (result PathPaymentStrictReceiveResult, err error) {
	result.Code = code
	switch PathPaymentStrictReceiveResultCode(code) {
	case PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveSuccess:
		tv, ok := value.(PathPaymentStrictReceiveResultSuccess)
		if !ok {
			err = fmt.Errorf("invalid value, must be PathPaymentStrictReceiveResultSuccess")
			return
		}
		result.Success = &tv
	case PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveNoIssuer:
		tv, ok := value.(Asset)
		if !ok {
			err = fmt.Errorf("invalid value, must be Asset")
			return
		}
		result.NoIssuer = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u PathPaymentStrictReceiveResult) MustSuccess() PathPaymentStrictReceiveResultSuccess {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u PathPaymentStrictReceiveResult) GetSuccess() (result PathPaymentStrictReceiveResultSuccess, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// MustNoIssuer retrieves the NoIssuer value from the union,
// panicing if the value is not set.
func (u PathPaymentStrictReceiveResult) MustNoIssuer() Asset {
	val, ok := u.GetNoIssuer()

	if !ok {
		panic("arm NoIssuer is not set")
	}

	return val
}

// GetNoIssuer retrieves the NoIssuer value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u PathPaymentStrictReceiveResult) GetNoIssuer() (result Asset, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "NoIssuer" {
		result = *u.NoIssuer
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u PathPaymentStrictReceiveResult) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Code.EncodeTo(e); err != nil {
		return err
	}
	switch PathPaymentStrictReceiveResultCode(u.Code) {
	case PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveSuccess:
		if err = (*u.Success).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveNoIssuer:
		if err = (*u.NoIssuer).EncodeTo(e); err != nil {
			return err
		}
		return nil
	default:
		// Void
		return nil
	}
}

var _ decoderFrom = (*PathPaymentStrictReceiveResult)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *PathPaymentStrictReceiveResult) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Code.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding PathPaymentStrictReceiveResultCode: %s", err)
	}
	switch PathPaymentStrictReceiveResultCode(u.Code) {
	case PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveSuccess:
		u.Success = new(PathPaymentStrictReceiveResultSuccess)
		nTmp, err = (*u.Success).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding PathPaymentStrictReceiveResultSuccess: %s", err)
		}
		return n, nil
	case PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveNoIssuer:
		u.NoIssuer = new(Asset)
		nTmp, err = (*u.NoIssuer).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Asset: %s", err)
		}
		return n, nil
	default:
		// Void
		return n, nil
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s PathPaymentStrictReceiveResult) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *PathPaymentStrictReceiveResult) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*PathPaymentStrictReceiveResult)(nil)
	_ encoding.BinaryUnmarshaler = (*PathPaymentStrictReceiveResult)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s PathPaymentStrictReceiveResult) xdrType() {}

var _ xdrType = (*PathPaymentStrictReceiveResult)(nil)

// PathPaymentStrictSendResultCode is an XDR Enum defines as:
//
//   enum PathPaymentStrictSendResultCode
//    {
//        // codes considered as "success" for the operation
//        PATH_PAYMENT_STRICT_SEND_SUCCESS = 0, // success
//
//        // codes considered as "failure" for the operation
//        PATH_PAYMENT_STRICT_SEND_MALFORMED = -1, // bad input
//        PATH_PAYMENT_STRICT_SEND_UNDERFUNDED =
//            -2, // not enough funds in source account
//        PATH_PAYMENT_STRICT_SEND_SRC_NO_TRUST =
//            -3, // no trust line on source account
//        PATH_PAYMENT_STRICT_SEND_SRC_NOT_AUTHORIZED =
//            -4, // source not authorized to transfer
//        PATH_PAYMENT_STRICT_SEND_NO_DESTINATION =
//            -5, // destination account does not exist
//        PATH_PAYMENT_STRICT_SEND_NO_TRUST =
//            -6, // dest missing a trust line for asset
//        PATH_PAYMENT_STRICT_SEND_NOT_AUTHORIZED =
//            -7, // dest not authorized to hold asset
//        PATH_PAYMENT_STRICT_SEND_LINE_FULL = -8, // dest would go above their limit
//        PATH_PAYMENT_STRICT_SEND_NO_ISSUER = -9, // missing issuer on one asset
//        PATH_PAYMENT_STRICT_SEND_TOO_FEW_OFFERS =
//            -10, // not enough offers to satisfy path
//        PATH_PAYMENT_STRICT_SEND_OFFER_CROSS_SELF =
//            -11, // would cross one of its own offers
//        PATH_PAYMENT_STRICT_SEND_UNDER_DESTMIN = -12 // could not satisfy destMin
//    };
//
type PathPaymentStrictSendResultCode int32

const (
	PathPaymentStrictSendResultCodePathPaymentStrictSendSuccess          PathPaymentStrictSendResultCode = 0
	PathPaymentStrictSendResultCodePathPaymentStrictSendMalformed        PathPaymentStrictSendResultCode = -1
	PathPaymentStrictSendResultCodePathPaymentStrictSendUnderfunded      PathPaymentStrictSendResultCode = -2
	PathPaymentStrictSendResultCodePathPaymentStrictSendSrcNoTrust       PathPaymentStrictSendResultCode = -3
	PathPaymentStrictSendResultCodePathPaymentStrictSendSrcNotAuthorized PathPaymentStrictSendResultCode = -4
	PathPaymentStrictSendResultCodePathPaymentStrictSendNoDestination    PathPaymentStrictSendResultCode = -5
	PathPaymentStrictSendResultCodePathPaymentStrictSendNoTrust          PathPaymentStrictSendResultCode = -6
	PathPaymentStrictSendResultCodePathPaymentStrictSendNotAuthorized    PathPaymentStrictSendResultCode = -7
	PathPaymentStrictSendResultCodePathPaymentStrictSendLineFull         PathPaymentStrictSendResultCode = -8
	PathPaymentStrictSendResultCodePathPaymentStrictSendNoIssuer         PathPaymentStrictSendResultCode = -9
	PathPaymentStrictSendResultCodePathPaymentStrictSendTooFewOffers     PathPaymentStrictSendResultCode = -10
	PathPaymentStrictSendResultCodePathPaymentStrictSendOfferCrossSelf   PathPaymentStrictSendResultCode = -11
	PathPaymentStrictSendResultCodePathPaymentStrictSendUnderDestmin     PathPaymentStrictSendResultCode = -12
)

var pathPaymentStrictSendResultCodeMap = map[int32]string{
	0:   "PathPaymentStrictSendResultCodePathPaymentStrictSendSuccess",
	-1:  "PathPaymentStrictSendResultCodePathPaymentStrictSendMalformed",
	-2:  "PathPaymentStrictSendResultCodePathPaymentStrictSendUnderfunded",
	-3:  "PathPaymentStrictSendResultCodePathPaymentStrictSendSrcNoTrust",
	-4:  "PathPaymentStrictSendResultCodePathPaymentStrictSendSrcNotAuthorized",
	-5:  "PathPaymentStrictSendResultCodePathPaymentStrictSendNoDestination",
	-6:  "PathPaymentStrictSendResultCodePathPaymentStrictSendNoTrust",
	-7:  "PathPaymentStrictSendResultCodePathPaymentStrictSendNotAuthorized",
	-8:  "PathPaymentStrictSendResultCodePathPaymentStrictSendLineFull",
	-9:  "PathPaymentStrictSendResultCodePathPaymentStrictSendNoIssuer",
	-10: "PathPaymentStrictSendResultCodePathPaymentStrictSendTooFewOffers",
	-11: "PathPaymentStrictSendResultCodePathPaymentStrictSendOfferCrossSelf",
	-12: "PathPaymentStrictSendResultCodePathPaymentStrictSendUnderDestmin",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for PathPaymentStrictSendResultCode
func (e PathPaymentStrictSendResultCode) ValidEnum(v int32) bool {
	_, ok := pathPaymentStrictSendResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e PathPaymentStrictSendResultCode) String() string {
	name, _ := pathPaymentStrictSendResultCodeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e PathPaymentStrictSendResultCode) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := pathPaymentStrictSendResultCodeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid PathPaymentStrictSendResultCode enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*PathPaymentStrictSendResultCode)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *PathPaymentStrictSendResultCode) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding PathPaymentStrictSendResultCode: %s", err)
	}
	if _, ok := pathPaymentStrictSendResultCodeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid PathPaymentStrictSendResultCode enum value", v)
	}
	*e = PathPaymentStrictSendResultCode(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s PathPaymentStrictSendResultCode) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *PathPaymentStrictSendResultCode) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*PathPaymentStrictSendResultCode)(nil)
	_ encoding.BinaryUnmarshaler = (*PathPaymentStrictSendResultCode)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s PathPaymentStrictSendResultCode) xdrType() {}

var _ xdrType = (*PathPaymentStrictSendResultCode)(nil)

// PathPaymentStrictSendResultSuccess is an XDR NestedStruct defines as:
//
//   struct
//        {
//            ClaimAtom offers<>;
//            SimplePaymentResult last;
//        }
//
type PathPaymentStrictSendResultSuccess struct {
	Offers []ClaimAtom
	Last   SimplePaymentResult
}

// EncodeTo encodes this value using the Encoder.
func (s *PathPaymentStrictSendResultSuccess) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeUint(uint32(len(s.Offers))); err != nil {
		return err
	}
	for i := 0; i < len(s.Offers); i++ {
		if err = s.Offers[i].EncodeTo(e); err != nil {
			return err
		}
	}
	if err = s.Last.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*PathPaymentStrictSendResultSuccess)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *PathPaymentStrictSendResultSuccess) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ClaimAtom: %s", err)
	}
	s.Offers = nil
	if l > 0 {
		s.Offers = make([]ClaimAtom, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.Offers[i].DecodeFrom(d)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding ClaimAtom: %s", err)
			}
		}
	}
	nTmp, err = s.Last.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SimplePaymentResult: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s PathPaymentStrictSendResultSuccess) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *PathPaymentStrictSendResultSuccess) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*PathPaymentStrictSendResultSuccess)(nil)
	_ encoding.BinaryUnmarshaler = (*PathPaymentStrictSendResultSuccess)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s PathPaymentStrictSendResultSuccess) xdrType() {}

var _ xdrType = (*PathPaymentStrictSendResultSuccess)(nil)

// PathPaymentStrictSendResult is an XDR Union defines as:
//
//   union PathPaymentStrictSendResult switch (PathPaymentStrictSendResultCode code)
//    {
//    case PATH_PAYMENT_STRICT_SEND_SUCCESS:
//        struct
//        {
//            ClaimAtom offers<>;
//            SimplePaymentResult last;
//        } success;
//    case PATH_PAYMENT_STRICT_SEND_NO_ISSUER:
//        Asset noIssuer; // the asset that caused the error
//    default:
//        void;
//    };
//
type PathPaymentStrictSendResult struct {
	Code     PathPaymentStrictSendResultCode
	Success  *PathPaymentStrictSendResultSuccess
	NoIssuer *Asset
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u PathPaymentStrictSendResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PathPaymentStrictSendResult
func (u PathPaymentStrictSendResult) ArmForSwitch(sw int32) (string, bool) {
	switch PathPaymentStrictSendResultCode(sw) {
	case PathPaymentStrictSendResultCodePathPaymentStrictSendSuccess:
		return "Success", true
	case PathPaymentStrictSendResultCodePathPaymentStrictSendNoIssuer:
		return "NoIssuer", true
	default:
		return "", true
	}
}

// NewPathPaymentStrictSendResult creates a new  PathPaymentStrictSendResult.
func NewPathPaymentStrictSendResult(code PathPaymentStrictSendResultCode, value interface{}) (result PathPaymentStrictSendResult, err error) {
	result.Code = code
	switch PathPaymentStrictSendResultCode(code) {
	case PathPaymentStrictSendResultCodePathPaymentStrictSendSuccess:
		tv, ok := value.(PathPaymentStrictSendResultSuccess)
		if !ok {
			err = fmt.Errorf("invalid value, must be PathPaymentStrictSendResultSuccess")
			return
		}
		result.Success = &tv
	case PathPaymentStrictSendResultCodePathPaymentStrictSendNoIssuer:
		tv, ok := value.(Asset)
		if !ok {
			err = fmt.Errorf("invalid value, must be Asset")
			return
		}
		result.NoIssuer = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u PathPaymentStrictSendResult) MustSuccess() PathPaymentStrictSendResultSuccess {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u PathPaymentStrictSendResult) GetSuccess() (result PathPaymentStrictSendResultSuccess, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// MustNoIssuer retrieves the NoIssuer value from the union,
// panicing if the value is not set.
func (u PathPaymentStrictSendResult) MustNoIssuer() Asset {
	val, ok := u.GetNoIssuer()

	if !ok {
		panic("arm NoIssuer is not set")
	}

	return val
}

// GetNoIssuer retrieves the NoIssuer value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u PathPaymentStrictSendResult) GetNoIssuer() (result Asset, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "NoIssuer" {
		result = *u.NoIssuer
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u PathPaymentStrictSendResult) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Code.EncodeTo(e); err != nil {
		return err
	}
	switch PathPaymentStrictSendResultCode(u.Code) {
	case PathPaymentStrictSendResultCodePathPaymentStrictSendSuccess:
		if err = (*u.Success).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case PathPaymentStrictSendResultCodePathPaymentStrictSendNoIssuer:
		if err = (*u.NoIssuer).EncodeTo(e); err != nil {
			return err
		}
		return nil
	default:
		// Void
		return nil
	}
}

var _ decoderFrom = (*PathPaymentStrictSendResult)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *PathPaymentStrictSendResult) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Code.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding PathPaymentStrictSendResultCode: %s", err)
	}
	switch PathPaymentStrictSendResultCode(u.Code) {
	case PathPaymentStrictSendResultCodePathPaymentStrictSendSuccess:
		u.Success = new(PathPaymentStrictSendResultSuccess)
		nTmp, err = (*u.Success).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding PathPaymentStrictSendResultSuccess: %s", err)
		}
		return n, nil
	case PathPaymentStrictSendResultCodePathPaymentStrictSendNoIssuer:
		u.NoIssuer = new(Asset)
		nTmp, err = (*u.NoIssuer).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Asset: %s", err)
		}
		return n, nil
	default:
		// Void
		return n, nil
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s PathPaymentStrictSendResult) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *PathPaymentStrictSendResult) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*PathPaymentStrictSendResult)(nil)
	_ encoding.BinaryUnmarshaler = (*PathPaymentStrictSendResult)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s PathPaymentStrictSendResult) xdrType() {}

var _ xdrType = (*PathPaymentStrictSendResult)(nil)

// ManageSellOfferResultCode is an XDR Enum defines as:
//
//   enum ManageSellOfferResultCode
//    {
//        // codes considered as "success" for the operation
//        MANAGE_SELL_OFFER_SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//        MANAGE_SELL_OFFER_MALFORMED = -1, // generated offer would be invalid
//        MANAGE_SELL_OFFER_SELL_NO_TRUST =
//            -2,                              // no trust line for what we're selling
//        MANAGE_SELL_OFFER_BUY_NO_TRUST = -3, // no trust line for what we're buying
//        MANAGE_SELL_OFFER_SELL_NOT_AUTHORIZED = -4, // not authorized to sell
//        MANAGE_SELL_OFFER_BUY_NOT_AUTHORIZED = -5,  // not authorized to buy
//        MANAGE_SELL_OFFER_LINE_FULL = -6, // can't receive more of what it's buying
//        MANAGE_SELL_OFFER_UNDERFUNDED = -7, // doesn't hold what it's trying to sell
//        MANAGE_SELL_OFFER_CROSS_SELF =
//            -8, // would cross an offer from the same user
//        MANAGE_SELL_OFFER_SELL_NO_ISSUER = -9, // no issuer for what we're selling
//        MANAGE_SELL_OFFER_BUY_NO_ISSUER = -10, // no issuer for what we're buying
//
//        // update errors
//        MANAGE_SELL_OFFER_NOT_FOUND =
//            -11, // offerID does not match an existing offer
//
//        MANAGE_SELL_OFFER_LOW_RESERVE =
//            -12 // not enough funds to create a new Offer
//    };
//
type ManageSellOfferResultCode int32

const (
	ManageSellOfferResultCodeManageSellOfferSuccess           ManageSellOfferResultCode = 0
	ManageSellOfferResultCodeManageSellOfferMalformed         ManageSellOfferResultCode = -1
	ManageSellOfferResultCodeManageSellOfferSellNoTrust       ManageSellOfferResultCode = -2
	ManageSellOfferResultCodeManageSellOfferBuyNoTrust        ManageSellOfferResultCode = -3
	ManageSellOfferResultCodeManageSellOfferSellNotAuthorized ManageSellOfferResultCode = -4
	ManageSellOfferResultCodeManageSellOfferBuyNotAuthorized  ManageSellOfferResultCode = -5
	ManageSellOfferResultCodeManageSellOfferLineFull          ManageSellOfferResultCode = -6
	ManageSellOfferResultCodeManageSellOfferUnderfunded       ManageSellOfferResultCode = -7
	ManageSellOfferResultCodeManageSellOfferCrossSelf         ManageSellOfferResultCode = -8
	ManageSellOfferResultCodeManageSellOfferSellNoIssuer      ManageSellOfferResultCode = -9
	ManageSellOfferResultCodeManageSellOfferBuyNoIssuer       ManageSellOfferResultCode = -10
	ManageSellOfferResultCodeManageSellOfferNotFound          ManageSellOfferResultCode = -11
	ManageSellOfferResultCodeManageSellOfferLowReserve        ManageSellOfferResultCode = -12
)

var manageSellOfferResultCodeMap = map[int32]string{
	0:   "ManageSellOfferResultCodeManageSellOfferSuccess",
	-1:  "ManageSellOfferResultCodeManageSellOfferMalformed",
	-2:  "ManageSellOfferResultCodeManageSellOfferSellNoTrust",
	-3:  "ManageSellOfferResultCodeManageSellOfferBuyNoTrust",
	-4:  "ManageSellOfferResultCodeManageSellOfferSellNotAuthorized",
	-5:  "ManageSellOfferResultCodeManageSellOfferBuyNotAuthorized",
	-6:  "ManageSellOfferResultCodeManageSellOfferLineFull",
	-7:  "ManageSellOfferResultCodeManageSellOfferUnderfunded",
	-8:  "ManageSellOfferResultCodeManageSellOfferCrossSelf",
	-9:  "ManageSellOfferResultCodeManageSellOfferSellNoIssuer",
	-10: "ManageSellOfferResultCodeManageSellOfferBuyNoIssuer",
	-11: "ManageSellOfferResultCodeManageSellOfferNotFound",
	-12: "ManageSellOfferResultCodeManageSellOfferLowReserve",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageSellOfferResultCode
func (e ManageSellOfferResultCode) ValidEnum(v int32) bool {
	_, ok := manageSellOfferResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageSellOfferResultCode) String() string {
	name, _ := manageSellOfferResultCodeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e ManageSellOfferResultCode) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := manageSellOfferResultCodeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid ManageSellOfferResultCode enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*ManageSellOfferResultCode)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *ManageSellOfferResultCode) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding ManageSellOfferResultCode: %s", err)
	}
	if _, ok := manageSellOfferResultCodeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid ManageSellOfferResultCode enum value", v)
	}
	*e = ManageSellOfferResultCode(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ManageSellOfferResultCode) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ManageSellOfferResultCode) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ManageSellOfferResultCode)(nil)
	_ encoding.BinaryUnmarshaler = (*ManageSellOfferResultCode)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ManageSellOfferResultCode) xdrType() {}

var _ xdrType = (*ManageSellOfferResultCode)(nil)

// ManageOfferEffect is an XDR Enum defines as:
//
//   enum ManageOfferEffect
//    {
//        MANAGE_OFFER_CREATED = 0,
//        MANAGE_OFFER_UPDATED = 1,
//        MANAGE_OFFER_DELETED = 2
//    };
//
type ManageOfferEffect int32

const (
	ManageOfferEffectManageOfferCreated ManageOfferEffect = 0
	ManageOfferEffectManageOfferUpdated ManageOfferEffect = 1
	ManageOfferEffectManageOfferDeleted ManageOfferEffect = 2
)

var manageOfferEffectMap = map[int32]string{
	0: "ManageOfferEffectManageOfferCreated",
	1: "ManageOfferEffectManageOfferUpdated",
	2: "ManageOfferEffectManageOfferDeleted",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageOfferEffect
func (e ManageOfferEffect) ValidEnum(v int32) bool {
	_, ok := manageOfferEffectMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageOfferEffect) String() string {
	name, _ := manageOfferEffectMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e ManageOfferEffect) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := manageOfferEffectMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid ManageOfferEffect enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*ManageOfferEffect)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *ManageOfferEffect) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding ManageOfferEffect: %s", err)
	}
	if _, ok := manageOfferEffectMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid ManageOfferEffect enum value", v)
	}
	*e = ManageOfferEffect(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ManageOfferEffect) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ManageOfferEffect) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ManageOfferEffect)(nil)
	_ encoding.BinaryUnmarshaler = (*ManageOfferEffect)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ManageOfferEffect) xdrType() {}

var _ xdrType = (*ManageOfferEffect)(nil)

// ManageOfferSuccessResultOffer is an XDR NestedUnion defines as:
//
//   union switch (ManageOfferEffect effect)
//        {
//        case MANAGE_OFFER_CREATED:
//        case MANAGE_OFFER_UPDATED:
//            OfferEntry offer;
//        default:
//            void;
//        }
//
type ManageOfferSuccessResultOffer struct {
	Effect ManageOfferEffect
	Offer  *OfferEntry
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageOfferSuccessResultOffer) SwitchFieldName() string {
	return "Effect"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageOfferSuccessResultOffer
func (u ManageOfferSuccessResultOffer) ArmForSwitch(sw int32) (string, bool) {
	switch ManageOfferEffect(sw) {
	case ManageOfferEffectManageOfferCreated:
		return "Offer", true
	case ManageOfferEffectManageOfferUpdated:
		return "Offer", true
	default:
		return "", true
	}
}

// NewManageOfferSuccessResultOffer creates a new  ManageOfferSuccessResultOffer.
func NewManageOfferSuccessResultOffer(effect ManageOfferEffect, value interface{}) (result ManageOfferSuccessResultOffer, err error) {
	result.Effect = effect
	switch ManageOfferEffect(effect) {
	case ManageOfferEffectManageOfferCreated:
		tv, ok := value.(OfferEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be OfferEntry")
			return
		}
		result.Offer = &tv
	case ManageOfferEffectManageOfferUpdated:
		tv, ok := value.(OfferEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be OfferEntry")
			return
		}
		result.Offer = &tv
	default:
		// void
	}
	return
}

// MustOffer retrieves the Offer value from the union,
// panicing if the value is not set.
func (u ManageOfferSuccessResultOffer) MustOffer() OfferEntry {
	val, ok := u.GetOffer()

	if !ok {
		panic("arm Offer is not set")
	}

	return val
}

// GetOffer retrieves the Offer value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ManageOfferSuccessResultOffer) GetOffer() (result OfferEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Effect))

	if armName == "Offer" {
		result = *u.Offer
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u ManageOfferSuccessResultOffer) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Effect.EncodeTo(e); err != nil {
		return err
	}
	switch ManageOfferEffect(u.Effect) {
	case ManageOfferEffectManageOfferCreated:
		if err = (*u.Offer).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case ManageOfferEffectManageOfferUpdated:
		if err = (*u.Offer).EncodeTo(e); err != nil {
			return err
		}
		return nil
	default:
		// Void
		return nil
	}
}

var _ decoderFrom = (*ManageOfferSuccessResultOffer)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *ManageOfferSuccessResultOffer) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Effect.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ManageOfferEffect: %s", err)
	}
	switch ManageOfferEffect(u.Effect) {
	case ManageOfferEffectManageOfferCreated:
		u.Offer = new(OfferEntry)
		nTmp, err = (*u.Offer).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding OfferEntry: %s", err)
		}
		return n, nil
	case ManageOfferEffectManageOfferUpdated:
		u.Offer = new(OfferEntry)
		nTmp, err = (*u.Offer).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding OfferEntry: %s", err)
		}
		return n, nil
	default:
		// Void
		return n, nil
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ManageOfferSuccessResultOffer) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ManageOfferSuccessResultOffer) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ManageOfferSuccessResultOffer)(nil)
	_ encoding.BinaryUnmarshaler = (*ManageOfferSuccessResultOffer)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ManageOfferSuccessResultOffer) xdrType() {}

var _ xdrType = (*ManageOfferSuccessResultOffer)(nil)

// ManageOfferSuccessResult is an XDR Struct defines as:
//
//   struct ManageOfferSuccessResult
//    {
//        // offers that got claimed while creating this offer
//        ClaimAtom offersClaimed<>;
//
//        union switch (ManageOfferEffect effect)
//        {
//        case MANAGE_OFFER_CREATED:
//        case MANAGE_OFFER_UPDATED:
//            OfferEntry offer;
//        default:
//            void;
//        }
//        offer;
//    };
//
type ManageOfferSuccessResult struct {
	OffersClaimed []ClaimAtom
	Offer         ManageOfferSuccessResultOffer
}

// EncodeTo encodes this value using the Encoder.
func (s *ManageOfferSuccessResult) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeUint(uint32(len(s.OffersClaimed))); err != nil {
		return err
	}
	for i := 0; i < len(s.OffersClaimed); i++ {
		if err = s.OffersClaimed[i].EncodeTo(e); err != nil {
			return err
		}
	}
	if err = s.Offer.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*ManageOfferSuccessResult)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *ManageOfferSuccessResult) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ClaimAtom: %s", err)
	}
	s.OffersClaimed = nil
	if l > 0 {
		s.OffersClaimed = make([]ClaimAtom, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.OffersClaimed[i].DecodeFrom(d)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding ClaimAtom: %s", err)
			}
		}
	}
	nTmp, err = s.Offer.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ManageOfferSuccessResultOffer: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ManageOfferSuccessResult) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ManageOfferSuccessResult) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ManageOfferSuccessResult)(nil)
	_ encoding.BinaryUnmarshaler = (*ManageOfferSuccessResult)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ManageOfferSuccessResult) xdrType() {}

var _ xdrType = (*ManageOfferSuccessResult)(nil)

// ManageSellOfferResult is an XDR Union defines as:
//
//   union ManageSellOfferResult switch (ManageSellOfferResultCode code)
//    {
//    case MANAGE_SELL_OFFER_SUCCESS:
//        ManageOfferSuccessResult success;
//    default:
//        void;
//    };
//
type ManageSellOfferResult struct {
	Code    ManageSellOfferResultCode
	Success *ManageOfferSuccessResult
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageSellOfferResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageSellOfferResult
func (u ManageSellOfferResult) ArmForSwitch(sw int32) (string, bool) {
	switch ManageSellOfferResultCode(sw) {
	case ManageSellOfferResultCodeManageSellOfferSuccess:
		return "Success", true
	default:
		return "", true
	}
}

// NewManageSellOfferResult creates a new  ManageSellOfferResult.
func NewManageSellOfferResult(code ManageSellOfferResultCode, value interface{}) (result ManageSellOfferResult, err error) {
	result.Code = code
	switch ManageSellOfferResultCode(code) {
	case ManageSellOfferResultCodeManageSellOfferSuccess:
		tv, ok := value.(ManageOfferSuccessResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageOfferSuccessResult")
			return
		}
		result.Success = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u ManageSellOfferResult) MustSuccess() ManageOfferSuccessResult {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ManageSellOfferResult) GetSuccess() (result ManageOfferSuccessResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u ManageSellOfferResult) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Code.EncodeTo(e); err != nil {
		return err
	}
	switch ManageSellOfferResultCode(u.Code) {
	case ManageSellOfferResultCodeManageSellOfferSuccess:
		if err = (*u.Success).EncodeTo(e); err != nil {
			return err
		}
		return nil
	default:
		// Void
		return nil
	}
}

var _ decoderFrom = (*ManageSellOfferResult)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *ManageSellOfferResult) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Code.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ManageSellOfferResultCode: %s", err)
	}
	switch ManageSellOfferResultCode(u.Code) {
	case ManageSellOfferResultCodeManageSellOfferSuccess:
		u.Success = new(ManageOfferSuccessResult)
		nTmp, err = (*u.Success).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ManageOfferSuccessResult: %s", err)
		}
		return n, nil
	default:
		// Void
		return n, nil
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ManageSellOfferResult) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ManageSellOfferResult) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ManageSellOfferResult)(nil)
	_ encoding.BinaryUnmarshaler = (*ManageSellOfferResult)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ManageSellOfferResult) xdrType() {}

var _ xdrType = (*ManageSellOfferResult)(nil)

// ManageBuyOfferResultCode is an XDR Enum defines as:
//
//   enum ManageBuyOfferResultCode
//    {
//        // codes considered as "success" for the operation
//        MANAGE_BUY_OFFER_SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//        MANAGE_BUY_OFFER_MALFORMED = -1,     // generated offer would be invalid
//        MANAGE_BUY_OFFER_SELL_NO_TRUST = -2, // no trust line for what we're selling
//        MANAGE_BUY_OFFER_BUY_NO_TRUST = -3,  // no trust line for what we're buying
//        MANAGE_BUY_OFFER_SELL_NOT_AUTHORIZED = -4, // not authorized to sell
//        MANAGE_BUY_OFFER_BUY_NOT_AUTHORIZED = -5,  // not authorized to buy
//        MANAGE_BUY_OFFER_LINE_FULL = -6,   // can't receive more of what it's buying
//        MANAGE_BUY_OFFER_UNDERFUNDED = -7, // doesn't hold what it's trying to sell
//        MANAGE_BUY_OFFER_CROSS_SELF = -8, // would cross an offer from the same user
//        MANAGE_BUY_OFFER_SELL_NO_ISSUER = -9, // no issuer for what we're selling
//        MANAGE_BUY_OFFER_BUY_NO_ISSUER = -10, // no issuer for what we're buying
//
//        // update errors
//        MANAGE_BUY_OFFER_NOT_FOUND =
//            -11, // offerID does not match an existing offer
//
//        MANAGE_BUY_OFFER_LOW_RESERVE = -12 // not enough funds to create a new Offer
//    };
//
type ManageBuyOfferResultCode int32

const (
	ManageBuyOfferResultCodeManageBuyOfferSuccess           ManageBuyOfferResultCode = 0
	ManageBuyOfferResultCodeManageBuyOfferMalformed         ManageBuyOfferResultCode = -1
	ManageBuyOfferResultCodeManageBuyOfferSellNoTrust       ManageBuyOfferResultCode = -2
	ManageBuyOfferResultCodeManageBuyOfferBuyNoTrust        ManageBuyOfferResultCode = -3
	ManageBuyOfferResultCodeManageBuyOfferSellNotAuthorized ManageBuyOfferResultCode = -4
	ManageBuyOfferResultCodeManageBuyOfferBuyNotAuthorized  ManageBuyOfferResultCode = -5
	ManageBuyOfferResultCodeManageBuyOfferLineFull          ManageBuyOfferResultCode = -6
	ManageBuyOfferResultCodeManageBuyOfferUnderfunded       ManageBuyOfferResultCode = -7
	ManageBuyOfferResultCodeManageBuyOfferCrossSelf         ManageBuyOfferResultCode = -8
	ManageBuyOfferResultCodeManageBuyOfferSellNoIssuer      ManageBuyOfferResultCode = -9
	ManageBuyOfferResultCodeManageBuyOfferBuyNoIssuer       ManageBuyOfferResultCode = -10
	ManageBuyOfferResultCodeManageBuyOfferNotFound          ManageBuyOfferResultCode = -11
	ManageBuyOfferResultCodeManageBuyOfferLowReserve        ManageBuyOfferResultCode = -12
)

var manageBuyOfferResultCodeMap = map[int32]string{
	0:   "ManageBuyOfferResultCodeManageBuyOfferSuccess",
	-1:  "ManageBuyOfferResultCodeManageBuyOfferMalformed",
	-2:  "ManageBuyOfferResultCodeManageBuyOfferSellNoTrust",
	-3:  "ManageBuyOfferResultCodeManageBuyOfferBuyNoTrust",
	-4:  "ManageBuyOfferResultCodeManageBuyOfferSellNotAuthorized",
	-5:  "ManageBuyOfferResultCodeManageBuyOfferBuyNotAuthorized",
	-6:  "ManageBuyOfferResultCodeManageBuyOfferLineFull",
	-7:  "ManageBuyOfferResultCodeManageBuyOfferUnderfunded",
	-8:  "ManageBuyOfferResultCodeManageBuyOfferCrossSelf",
	-9:  "ManageBuyOfferResultCodeManageBuyOfferSellNoIssuer",
	-10: "ManageBuyOfferResultCodeManageBuyOfferBuyNoIssuer",
	-11: "ManageBuyOfferResultCodeManageBuyOfferNotFound",
	-12: "ManageBuyOfferResultCodeManageBuyOfferLowReserve",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageBuyOfferResultCode
func (e ManageBuyOfferResultCode) ValidEnum(v int32) bool {
	_, ok := manageBuyOfferResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageBuyOfferResultCode) String() string {
	name, _ := manageBuyOfferResultCodeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e ManageBuyOfferResultCode) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := manageBuyOfferResultCodeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid ManageBuyOfferResultCode enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*ManageBuyOfferResultCode)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *ManageBuyOfferResultCode) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding ManageBuyOfferResultCode: %s", err)
	}
	if _, ok := manageBuyOfferResultCodeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid ManageBuyOfferResultCode enum value", v)
	}
	*e = ManageBuyOfferResultCode(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ManageBuyOfferResultCode) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ManageBuyOfferResultCode) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ManageBuyOfferResultCode)(nil)
	_ encoding.BinaryUnmarshaler = (*ManageBuyOfferResultCode)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ManageBuyOfferResultCode) xdrType() {}

var _ xdrType = (*ManageBuyOfferResultCode)(nil)

// ManageBuyOfferResult is an XDR Union defines as:
//
//   union ManageBuyOfferResult switch (ManageBuyOfferResultCode code)
//    {
//    case MANAGE_BUY_OFFER_SUCCESS:
//        ManageOfferSuccessResult success;
//    default:
//        void;
//    };
//
type ManageBuyOfferResult struct {
	Code    ManageBuyOfferResultCode
	Success *ManageOfferSuccessResult
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageBuyOfferResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageBuyOfferResult
func (u ManageBuyOfferResult) ArmForSwitch(sw int32) (string, bool) {
	switch ManageBuyOfferResultCode(sw) {
	case ManageBuyOfferResultCodeManageBuyOfferSuccess:
		return "Success", true
	default:
		return "", true
	}
}

// NewManageBuyOfferResult creates a new  ManageBuyOfferResult.
func NewManageBuyOfferResult(code ManageBuyOfferResultCode, value interface{}) (result ManageBuyOfferResult, err error) {
	result.Code = code
	switch ManageBuyOfferResultCode(code) {
	case ManageBuyOfferResultCodeManageBuyOfferSuccess:
		tv, ok := value.(ManageOfferSuccessResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageOfferSuccessResult")
			return
		}
		result.Success = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u ManageBuyOfferResult) MustSuccess() ManageOfferSuccessResult {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ManageBuyOfferResult) GetSuccess() (result ManageOfferSuccessResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u ManageBuyOfferResult) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Code.EncodeTo(e); err != nil {
		return err
	}
	switch ManageBuyOfferResultCode(u.Code) {
	case ManageBuyOfferResultCodeManageBuyOfferSuccess:
		if err = (*u.Success).EncodeTo(e); err != nil {
			return err
		}
		return nil
	default:
		// Void
		return nil
	}
}

var _ decoderFrom = (*ManageBuyOfferResult)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *ManageBuyOfferResult) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Code.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ManageBuyOfferResultCode: %s", err)
	}
	switch ManageBuyOfferResultCode(u.Code) {
	case ManageBuyOfferResultCodeManageBuyOfferSuccess:
		u.Success = new(ManageOfferSuccessResult)
		nTmp, err = (*u.Success).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ManageOfferSuccessResult: %s", err)
		}
		return n, nil
	default:
		// Void
		return n, nil
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ManageBuyOfferResult) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ManageBuyOfferResult) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ManageBuyOfferResult)(nil)
	_ encoding.BinaryUnmarshaler = (*ManageBuyOfferResult)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ManageBuyOfferResult) xdrType() {}

var _ xdrType = (*ManageBuyOfferResult)(nil)

// SetOptionsResultCode is an XDR Enum defines as:
//
//   enum SetOptionsResultCode
//    {
//        // codes considered as "success" for the operation
//        SET_OPTIONS_SUCCESS = 0,
//        // codes considered as "failure" for the operation
//        SET_OPTIONS_LOW_RESERVE = -1,      // not enough funds to add a signer
//        SET_OPTIONS_TOO_MANY_SIGNERS = -2, // max number of signers already reached
//        SET_OPTIONS_BAD_FLAGS = -3,        // invalid combination of clear/set flags
//        SET_OPTIONS_INVALID_INFLATION = -4,      // inflation account does not exist
//        SET_OPTIONS_CANT_CHANGE = -5,            // can no longer change this option
//        SET_OPTIONS_UNKNOWN_FLAG = -6,           // can't set an unknown flag
//        SET_OPTIONS_THRESHOLD_OUT_OF_RANGE = -7, // bad value for weight/threshold
//        SET_OPTIONS_BAD_SIGNER = -8,             // signer cannot be masterkey
//        SET_OPTIONS_INVALID_HOME_DOMAIN = -9,    // malformed home domain
//        SET_OPTIONS_AUTH_REVOCABLE_REQUIRED =
//            -10 // auth revocable is required for clawback
//    };
//
type SetOptionsResultCode int32

const (
	SetOptionsResultCodeSetOptionsSuccess               SetOptionsResultCode = 0
	SetOptionsResultCodeSetOptionsLowReserve            SetOptionsResultCode = -1
	SetOptionsResultCodeSetOptionsTooManySigners        SetOptionsResultCode = -2
	SetOptionsResultCodeSetOptionsBadFlags              SetOptionsResultCode = -3
	SetOptionsResultCodeSetOptionsInvalidInflation      SetOptionsResultCode = -4
	SetOptionsResultCodeSetOptionsCantChange            SetOptionsResultCode = -5
	SetOptionsResultCodeSetOptionsUnknownFlag           SetOptionsResultCode = -6
	SetOptionsResultCodeSetOptionsThresholdOutOfRange   SetOptionsResultCode = -7
	SetOptionsResultCodeSetOptionsBadSigner             SetOptionsResultCode = -8
	SetOptionsResultCodeSetOptionsInvalidHomeDomain     SetOptionsResultCode = -9
	SetOptionsResultCodeSetOptionsAuthRevocableRequired SetOptionsResultCode = -10
)

var setOptionsResultCodeMap = map[int32]string{
	0:   "SetOptionsResultCodeSetOptionsSuccess",
	-1:  "SetOptionsResultCodeSetOptionsLowReserve",
	-2:  "SetOptionsResultCodeSetOptionsTooManySigners",
	-3:  "SetOptionsResultCodeSetOptionsBadFlags",
	-4:  "SetOptionsResultCodeSetOptionsInvalidInflation",
	-5:  "SetOptionsResultCodeSetOptionsCantChange",
	-6:  "SetOptionsResultCodeSetOptionsUnknownFlag",
	-7:  "SetOptionsResultCodeSetOptionsThresholdOutOfRange",
	-8:  "SetOptionsResultCodeSetOptionsBadSigner",
	-9:  "SetOptionsResultCodeSetOptionsInvalidHomeDomain",
	-10: "SetOptionsResultCodeSetOptionsAuthRevocableRequired",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for SetOptionsResultCode
func (e SetOptionsResultCode) ValidEnum(v int32) bool {
	_, ok := setOptionsResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e SetOptionsResultCode) String() string {
	name, _ := setOptionsResultCodeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e SetOptionsResultCode) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := setOptionsResultCodeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid SetOptionsResultCode enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*SetOptionsResultCode)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *SetOptionsResultCode) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding SetOptionsResultCode: %s", err)
	}
	if _, ok := setOptionsResultCodeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid SetOptionsResultCode enum value", v)
	}
	*e = SetOptionsResultCode(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s SetOptionsResultCode) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SetOptionsResultCode) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*SetOptionsResultCode)(nil)
	_ encoding.BinaryUnmarshaler = (*SetOptionsResultCode)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s SetOptionsResultCode) xdrType() {}

var _ xdrType = (*SetOptionsResultCode)(nil)

// SetOptionsResult is an XDR Union defines as:
//
//   union SetOptionsResult switch (SetOptionsResultCode code)
//    {
//    case SET_OPTIONS_SUCCESS:
//        void;
//    default:
//        void;
//    };
//
type SetOptionsResult struct {
	Code SetOptionsResultCode
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u SetOptionsResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of SetOptionsResult
func (u SetOptionsResult) ArmForSwitch(sw int32) (string, bool) {
	switch SetOptionsResultCode(sw) {
	case SetOptionsResultCodeSetOptionsSuccess:
		return "", true
	default:
		return "", true
	}
}

// NewSetOptionsResult creates a new  SetOptionsResult.
func NewSetOptionsResult(code SetOptionsResultCode, value interface{}) (result SetOptionsResult, err error) {
	result.Code = code
	switch SetOptionsResultCode(code) {
	case SetOptionsResultCodeSetOptionsSuccess:
		// void
	default:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u SetOptionsResult) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Code.EncodeTo(e); err != nil {
		return err
	}
	switch SetOptionsResultCode(u.Code) {
	case SetOptionsResultCodeSetOptionsSuccess:
		// Void
		return nil
	default:
		// Void
		return nil
	}
}

var _ decoderFrom = (*SetOptionsResult)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *SetOptionsResult) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Code.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SetOptionsResultCode: %s", err)
	}
	switch SetOptionsResultCode(u.Code) {
	case SetOptionsResultCodeSetOptionsSuccess:
		// Void
		return n, nil
	default:
		// Void
		return n, nil
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s SetOptionsResult) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SetOptionsResult) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*SetOptionsResult)(nil)
	_ encoding.BinaryUnmarshaler = (*SetOptionsResult)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s SetOptionsResult) xdrType() {}

var _ xdrType = (*SetOptionsResult)(nil)

// ChangeTrustResultCode is an XDR Enum defines as:
//
//   enum ChangeTrustResultCode
//    {
//        // codes considered as "success" for the operation
//        CHANGE_TRUST_SUCCESS = 0,
//        // codes considered as "failure" for the operation
//        CHANGE_TRUST_MALFORMED = -1,     // bad input
//        CHANGE_TRUST_NO_ISSUER = -2,     // could not find issuer
//        CHANGE_TRUST_INVALID_LIMIT = -3, // cannot drop limit below balance
//                                         // cannot create with a limit of 0
//        CHANGE_TRUST_LOW_RESERVE =
//            -4, // not enough funds to create a new trust line,
//        CHANGE_TRUST_SELF_NOT_ALLOWED = -5, // trusting self is not allowed
//        CHANGE_TRUST_TRUST_LINE_MISSING = -6, // Asset trustline is missing for pool
//        CHANGE_TRUST_CANNOT_DELETE = -7, // Asset trustline is still referenced in a pool
//        CHANGE_TRUST_NOT_AUTH_MAINTAIN_LIABILITIES = -8 // Asset trustline is deauthorized
//    };
//
type ChangeTrustResultCode int32

const (
	ChangeTrustResultCodeChangeTrustSuccess                    ChangeTrustResultCode = 0
	ChangeTrustResultCodeChangeTrustMalformed                  ChangeTrustResultCode = -1
	ChangeTrustResultCodeChangeTrustNoIssuer                   ChangeTrustResultCode = -2
	ChangeTrustResultCodeChangeTrustInvalidLimit               ChangeTrustResultCode = -3
	ChangeTrustResultCodeChangeTrustLowReserve                 ChangeTrustResultCode = -4
	ChangeTrustResultCodeChangeTrustSelfNotAllowed             ChangeTrustResultCode = -5
	ChangeTrustResultCodeChangeTrustTrustLineMissing           ChangeTrustResultCode = -6
	ChangeTrustResultCodeChangeTrustCannotDelete               ChangeTrustResultCode = -7
	ChangeTrustResultCodeChangeTrustNotAuthMaintainLiabilities ChangeTrustResultCode = -8
)

var changeTrustResultCodeMap = map[int32]string{
	0:  "ChangeTrustResultCodeChangeTrustSuccess",
	-1: "ChangeTrustResultCodeChangeTrustMalformed",
	-2: "ChangeTrustResultCodeChangeTrustNoIssuer",
	-3: "ChangeTrustResultCodeChangeTrustInvalidLimit",
	-4: "ChangeTrustResultCodeChangeTrustLowReserve",
	-5: "ChangeTrustResultCodeChangeTrustSelfNotAllowed",
	-6: "ChangeTrustResultCodeChangeTrustTrustLineMissing",
	-7: "ChangeTrustResultCodeChangeTrustCannotDelete",
	-8: "ChangeTrustResultCodeChangeTrustNotAuthMaintainLiabilities",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ChangeTrustResultCode
func (e ChangeTrustResultCode) ValidEnum(v int32) bool {
	_, ok := changeTrustResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ChangeTrustResultCode) String() string {
	name, _ := changeTrustResultCodeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e ChangeTrustResultCode) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := changeTrustResultCodeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid ChangeTrustResultCode enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*ChangeTrustResultCode)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *ChangeTrustResultCode) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding ChangeTrustResultCode: %s", err)
	}
	if _, ok := changeTrustResultCodeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid ChangeTrustResultCode enum value", v)
	}
	*e = ChangeTrustResultCode(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ChangeTrustResultCode) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ChangeTrustResultCode) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ChangeTrustResultCode)(nil)
	_ encoding.BinaryUnmarshaler = (*ChangeTrustResultCode)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ChangeTrustResultCode) xdrType() {}

var _ xdrType = (*ChangeTrustResultCode)(nil)

// ChangeTrustResult is an XDR Union defines as:
//
//   union ChangeTrustResult switch (ChangeTrustResultCode code)
//    {
//    case CHANGE_TRUST_SUCCESS:
//        void;
//    default:
//        void;
//    };
//
type ChangeTrustResult struct {
	Code ChangeTrustResultCode
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ChangeTrustResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ChangeTrustResult
func (u ChangeTrustResult) ArmForSwitch(sw int32) (string, bool) {
	switch ChangeTrustResultCode(sw) {
	case ChangeTrustResultCodeChangeTrustSuccess:
		return "", true
	default:
		return "", true
	}
}

// NewChangeTrustResult creates a new  ChangeTrustResult.
func NewChangeTrustResult(code ChangeTrustResultCode, value interface{}) (result ChangeTrustResult, err error) {
	result.Code = code
	switch ChangeTrustResultCode(code) {
	case ChangeTrustResultCodeChangeTrustSuccess:
		// void
	default:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u ChangeTrustResult) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Code.EncodeTo(e); err != nil {
		return err
	}
	switch ChangeTrustResultCode(u.Code) {
	case ChangeTrustResultCodeChangeTrustSuccess:
		// Void
		return nil
	default:
		// Void
		return nil
	}
}

var _ decoderFrom = (*ChangeTrustResult)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *ChangeTrustResult) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Code.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ChangeTrustResultCode: %s", err)
	}
	switch ChangeTrustResultCode(u.Code) {
	case ChangeTrustResultCodeChangeTrustSuccess:
		// Void
		return n, nil
	default:
		// Void
		return n, nil
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ChangeTrustResult) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ChangeTrustResult) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ChangeTrustResult)(nil)
	_ encoding.BinaryUnmarshaler = (*ChangeTrustResult)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ChangeTrustResult) xdrType() {}

var _ xdrType = (*ChangeTrustResult)(nil)

// AllowTrustResultCode is an XDR Enum defines as:
//
//   enum AllowTrustResultCode
//    {
//        // codes considered as "success" for the operation
//        ALLOW_TRUST_SUCCESS = 0,
//        // codes considered as "failure" for the operation
//        ALLOW_TRUST_MALFORMED = -1,     // asset is not ASSET_TYPE_ALPHANUM
//        ALLOW_TRUST_NO_TRUST_LINE = -2, // trustor does not have a trustline
//                                        // source account does not require trust
//        ALLOW_TRUST_TRUST_NOT_REQUIRED = -3,
//        ALLOW_TRUST_CANT_REVOKE = -4,     // source account can't revoke trust,
//        ALLOW_TRUST_SELF_NOT_ALLOWED = -5, // trusting self is not allowed
//        ALLOW_TRUST_LOW_RESERVE = -6 // claimable balances can't be created
//                                     // on revoke due to low reserves
//    };
//
type AllowTrustResultCode int32

const (
	AllowTrustResultCodeAllowTrustSuccess          AllowTrustResultCode = 0
	AllowTrustResultCodeAllowTrustMalformed        AllowTrustResultCode = -1
	AllowTrustResultCodeAllowTrustNoTrustLine      AllowTrustResultCode = -2
	AllowTrustResultCodeAllowTrustTrustNotRequired AllowTrustResultCode = -3
	AllowTrustResultCodeAllowTrustCantRevoke       AllowTrustResultCode = -4
	AllowTrustResultCodeAllowTrustSelfNotAllowed   AllowTrustResultCode = -5
	AllowTrustResultCodeAllowTrustLowReserve       AllowTrustResultCode = -6
)

var allowTrustResultCodeMap = map[int32]string{
	0:  "AllowTrustResultCodeAllowTrustSuccess",
	-1: "AllowTrustResultCodeAllowTrustMalformed",
	-2: "AllowTrustResultCodeAllowTrustNoTrustLine",
	-3: "AllowTrustResultCodeAllowTrustTrustNotRequired",
	-4: "AllowTrustResultCodeAllowTrustCantRevoke",
	-5: "AllowTrustResultCodeAllowTrustSelfNotAllowed",
	-6: "AllowTrustResultCodeAllowTrustLowReserve",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for AllowTrustResultCode
func (e AllowTrustResultCode) ValidEnum(v int32) bool {
	_, ok := allowTrustResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e AllowTrustResultCode) String() string {
	name, _ := allowTrustResultCodeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e AllowTrustResultCode) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := allowTrustResultCodeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid AllowTrustResultCode enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*AllowTrustResultCode)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *AllowTrustResultCode) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding AllowTrustResultCode: %s", err)
	}
	if _, ok := allowTrustResultCodeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid AllowTrustResultCode enum value", v)
	}
	*e = AllowTrustResultCode(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s AllowTrustResultCode) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AllowTrustResultCode) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*AllowTrustResultCode)(nil)
	_ encoding.BinaryUnmarshaler = (*AllowTrustResultCode)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s AllowTrustResultCode) xdrType() {}

var _ xdrType = (*AllowTrustResultCode)(nil)

// AllowTrustResult is an XDR Union defines as:
//
//   union AllowTrustResult switch (AllowTrustResultCode code)
//    {
//    case ALLOW_TRUST_SUCCESS:
//        void;
//    default:
//        void;
//    };
//
type AllowTrustResult struct {
	Code AllowTrustResultCode
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u AllowTrustResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of AllowTrustResult
func (u AllowTrustResult) ArmForSwitch(sw int32) (string, bool) {
	switch AllowTrustResultCode(sw) {
	case AllowTrustResultCodeAllowTrustSuccess:
		return "", true
	default:
		return "", true
	}
}

// NewAllowTrustResult creates a new  AllowTrustResult.
func NewAllowTrustResult(code AllowTrustResultCode, value interface{}) (result AllowTrustResult, err error) {
	result.Code = code
	switch AllowTrustResultCode(code) {
	case AllowTrustResultCodeAllowTrustSuccess:
		// void
	default:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u AllowTrustResult) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Code.EncodeTo(e); err != nil {
		return err
	}
	switch AllowTrustResultCode(u.Code) {
	case AllowTrustResultCodeAllowTrustSuccess:
		// Void
		return nil
	default:
		// Void
		return nil
	}
}

var _ decoderFrom = (*AllowTrustResult)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *AllowTrustResult) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Code.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AllowTrustResultCode: %s", err)
	}
	switch AllowTrustResultCode(u.Code) {
	case AllowTrustResultCodeAllowTrustSuccess:
		// Void
		return n, nil
	default:
		// Void
		return n, nil
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s AllowTrustResult) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AllowTrustResult) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*AllowTrustResult)(nil)
	_ encoding.BinaryUnmarshaler = (*AllowTrustResult)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s AllowTrustResult) xdrType() {}

var _ xdrType = (*AllowTrustResult)(nil)

// AccountMergeResultCode is an XDR Enum defines as:
//
//   enum AccountMergeResultCode
//    {
//        // codes considered as "success" for the operation
//        ACCOUNT_MERGE_SUCCESS = 0,
//        // codes considered as "failure" for the operation
//        ACCOUNT_MERGE_MALFORMED = -1,       // can't merge onto itself
//        ACCOUNT_MERGE_NO_ACCOUNT = -2,      // destination does not exist
//        ACCOUNT_MERGE_IMMUTABLE_SET = -3,   // source account has AUTH_IMMUTABLE set
//        ACCOUNT_MERGE_HAS_SUB_ENTRIES = -4, // account has trust lines/offers
//        ACCOUNT_MERGE_SEQNUM_TOO_FAR = -5,  // sequence number is over max allowed
//        ACCOUNT_MERGE_DEST_FULL = -6,       // can't add source balance to
//                                            // destination balance
//        ACCOUNT_MERGE_IS_SPONSOR = -7       // can't merge account that is a sponsor
//    };
//
type AccountMergeResultCode int32

const (
	AccountMergeResultCodeAccountMergeSuccess       AccountMergeResultCode = 0
	AccountMergeResultCodeAccountMergeMalformed     AccountMergeResultCode = -1
	AccountMergeResultCodeAccountMergeNoAccount     AccountMergeResultCode = -2
	AccountMergeResultCodeAccountMergeImmutableSet  AccountMergeResultCode = -3
	AccountMergeResultCodeAccountMergeHasSubEntries AccountMergeResultCode = -4
	AccountMergeResultCodeAccountMergeSeqnumTooFar  AccountMergeResultCode = -5
	AccountMergeResultCodeAccountMergeDestFull      AccountMergeResultCode = -6
	AccountMergeResultCodeAccountMergeIsSponsor     AccountMergeResultCode = -7
)

var accountMergeResultCodeMap = map[int32]string{
	0:  "AccountMergeResultCodeAccountMergeSuccess",
	-1: "AccountMergeResultCodeAccountMergeMalformed",
	-2: "AccountMergeResultCodeAccountMergeNoAccount",
	-3: "AccountMergeResultCodeAccountMergeImmutableSet",
	-4: "AccountMergeResultCodeAccountMergeHasSubEntries",
	-5: "AccountMergeResultCodeAccountMergeSeqnumTooFar",
	-6: "AccountMergeResultCodeAccountMergeDestFull",
	-7: "AccountMergeResultCodeAccountMergeIsSponsor",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for AccountMergeResultCode
func (e AccountMergeResultCode) ValidEnum(v int32) bool {
	_, ok := accountMergeResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e AccountMergeResultCode) String() string {
	name, _ := accountMergeResultCodeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e AccountMergeResultCode) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := accountMergeResultCodeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid AccountMergeResultCode enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*AccountMergeResultCode)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *AccountMergeResultCode) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding AccountMergeResultCode: %s", err)
	}
	if _, ok := accountMergeResultCodeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid AccountMergeResultCode enum value", v)
	}
	*e = AccountMergeResultCode(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s AccountMergeResultCode) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AccountMergeResultCode) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*AccountMergeResultCode)(nil)
	_ encoding.BinaryUnmarshaler = (*AccountMergeResultCode)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s AccountMergeResultCode) xdrType() {}

var _ xdrType = (*AccountMergeResultCode)(nil)

// AccountMergeResult is an XDR Union defines as:
//
//   union AccountMergeResult switch (AccountMergeResultCode code)
//    {
//    case ACCOUNT_MERGE_SUCCESS:
//        int64 sourceAccountBalance; // how much got transferred from source account
//    default:
//        void;
//    };
//
type AccountMergeResult struct {
	Code                 AccountMergeResultCode
	SourceAccountBalance *Int64
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u AccountMergeResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of AccountMergeResult
func (u AccountMergeResult) ArmForSwitch(sw int32) (string, bool) {
	switch AccountMergeResultCode(sw) {
	case AccountMergeResultCodeAccountMergeSuccess:
		return "SourceAccountBalance", true
	default:
		return "", true
	}
}

// NewAccountMergeResult creates a new  AccountMergeResult.
func NewAccountMergeResult(code AccountMergeResultCode, value interface{}) (result AccountMergeResult, err error) {
	result.Code = code
	switch AccountMergeResultCode(code) {
	case AccountMergeResultCodeAccountMergeSuccess:
		tv, ok := value.(Int64)
		if !ok {
			err = fmt.Errorf("invalid value, must be Int64")
			return
		}
		result.SourceAccountBalance = &tv
	default:
		// void
	}
	return
}

// MustSourceAccountBalance retrieves the SourceAccountBalance value from the union,
// panicing if the value is not set.
func (u AccountMergeResult) MustSourceAccountBalance() Int64 {
	val, ok := u.GetSourceAccountBalance()

	if !ok {
		panic("arm SourceAccountBalance is not set")
	}

	return val
}

// GetSourceAccountBalance retrieves the SourceAccountBalance value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u AccountMergeResult) GetSourceAccountBalance() (result Int64, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "SourceAccountBalance" {
		result = *u.SourceAccountBalance
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u AccountMergeResult) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Code.EncodeTo(e); err != nil {
		return err
	}
	switch AccountMergeResultCode(u.Code) {
	case AccountMergeResultCodeAccountMergeSuccess:
		if err = (*u.SourceAccountBalance).EncodeTo(e); err != nil {
			return err
		}
		return nil
	default:
		// Void
		return nil
	}
}

var _ decoderFrom = (*AccountMergeResult)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *AccountMergeResult) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Code.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountMergeResultCode: %s", err)
	}
	switch AccountMergeResultCode(u.Code) {
	case AccountMergeResultCodeAccountMergeSuccess:
		u.SourceAccountBalance = new(Int64)
		nTmp, err = (*u.SourceAccountBalance).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Int64: %s", err)
		}
		return n, nil
	default:
		// Void
		return n, nil
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s AccountMergeResult) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AccountMergeResult) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*AccountMergeResult)(nil)
	_ encoding.BinaryUnmarshaler = (*AccountMergeResult)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s AccountMergeResult) xdrType() {}

var _ xdrType = (*AccountMergeResult)(nil)

// InflationResultCode is an XDR Enum defines as:
//
//   enum InflationResultCode
//    {
//        // codes considered as "success" for the operation
//        INFLATION_SUCCESS = 0,
//        // codes considered as "failure" for the operation
//        INFLATION_NOT_TIME = -1
//    };
//
type InflationResultCode int32

const (
	InflationResultCodeInflationSuccess InflationResultCode = 0
	InflationResultCodeInflationNotTime InflationResultCode = -1
)

var inflationResultCodeMap = map[int32]string{
	0:  "InflationResultCodeInflationSuccess",
	-1: "InflationResultCodeInflationNotTime",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for InflationResultCode
func (e InflationResultCode) ValidEnum(v int32) bool {
	_, ok := inflationResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e InflationResultCode) String() string {
	name, _ := inflationResultCodeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e InflationResultCode) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := inflationResultCodeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid InflationResultCode enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*InflationResultCode)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *InflationResultCode) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding InflationResultCode: %s", err)
	}
	if _, ok := inflationResultCodeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid InflationResultCode enum value", v)
	}
	*e = InflationResultCode(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s InflationResultCode) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *InflationResultCode) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*InflationResultCode)(nil)
	_ encoding.BinaryUnmarshaler = (*InflationResultCode)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s InflationResultCode) xdrType() {}

var _ xdrType = (*InflationResultCode)(nil)

// InflationPayout is an XDR Struct defines as:
//
//   struct InflationPayout // or use PaymentResultAtom to limit types?
//    {
//        AccountID destination;
//        int64 amount;
//    };
//
type InflationPayout struct {
	Destination AccountId
	Amount      Int64
}

// EncodeTo encodes this value using the Encoder.
func (s *InflationPayout) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Destination.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Amount.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*InflationPayout)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *InflationPayout) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.Destination.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountId: %s", err)
	}
	nTmp, err = s.Amount.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s InflationPayout) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *InflationPayout) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*InflationPayout)(nil)
	_ encoding.BinaryUnmarshaler = (*InflationPayout)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s InflationPayout) xdrType() {}

var _ xdrType = (*InflationPayout)(nil)

// InflationResult is an XDR Union defines as:
//
//   union InflationResult switch (InflationResultCode code)
//    {
//    case INFLATION_SUCCESS:
//        InflationPayout payouts<>;
//    default:
//        void;
//    };
//
type InflationResult struct {
	Code    InflationResultCode
	Payouts *[]InflationPayout
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u InflationResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of InflationResult
func (u InflationResult) ArmForSwitch(sw int32) (string, bool) {
	switch InflationResultCode(sw) {
	case InflationResultCodeInflationSuccess:
		return "Payouts", true
	default:
		return "", true
	}
}

// NewInflationResult creates a new  InflationResult.
func NewInflationResult(code InflationResultCode, value interface{}) (result InflationResult, err error) {
	result.Code = code
	switch InflationResultCode(code) {
	case InflationResultCodeInflationSuccess:
		tv, ok := value.([]InflationPayout)
		if !ok {
			err = fmt.Errorf("invalid value, must be []InflationPayout")
			return
		}
		result.Payouts = &tv
	default:
		// void
	}
	return
}

// MustPayouts retrieves the Payouts value from the union,
// panicing if the value is not set.
func (u InflationResult) MustPayouts() []InflationPayout {
	val, ok := u.GetPayouts()

	if !ok {
		panic("arm Payouts is not set")
	}

	return val
}

// GetPayouts retrieves the Payouts value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u InflationResult) GetPayouts() (result []InflationPayout, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Payouts" {
		result = *u.Payouts
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u InflationResult) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Code.EncodeTo(e); err != nil {
		return err
	}
	switch InflationResultCode(u.Code) {
	case InflationResultCodeInflationSuccess:
		if _, err = e.EncodeUint(uint32(len((*u.Payouts)))); err != nil {
			return err
		}
		for i := 0; i < len((*u.Payouts)); i++ {
			if err = (*u.Payouts)[i].EncodeTo(e); err != nil {
				return err
			}
		}
		return nil
	default:
		// Void
		return nil
	}
}

var _ decoderFrom = (*InflationResult)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *InflationResult) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Code.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding InflationResultCode: %s", err)
	}
	switch InflationResultCode(u.Code) {
	case InflationResultCodeInflationSuccess:
		u.Payouts = new([]InflationPayout)
		var l uint32
		l, nTmp, err = d.DecodeUint()
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding InflationPayout: %s", err)
		}
		(*u.Payouts) = nil
		if l > 0 {
			(*u.Payouts) = make([]InflationPayout, l)
			for i := uint32(0); i < l; i++ {
				nTmp, err = (*u.Payouts)[i].DecodeFrom(d)
				n += nTmp
				if err != nil {
					return n, fmt.Errorf("decoding InflationPayout: %s", err)
				}
			}
		}
		return n, nil
	default:
		// Void
		return n, nil
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s InflationResult) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *InflationResult) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*InflationResult)(nil)
	_ encoding.BinaryUnmarshaler = (*InflationResult)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s InflationResult) xdrType() {}

var _ xdrType = (*InflationResult)(nil)

// ManageDataResultCode is an XDR Enum defines as:
//
//   enum ManageDataResultCode
//    {
//        // codes considered as "success" for the operation
//        MANAGE_DATA_SUCCESS = 0,
//        // codes considered as "failure" for the operation
//        MANAGE_DATA_NOT_SUPPORTED_YET =
//            -1, // The network hasn't moved to this protocol change yet
//        MANAGE_DATA_NAME_NOT_FOUND =
//            -2, // Trying to remove a Data Entry that isn't there
//        MANAGE_DATA_LOW_RESERVE = -3, // not enough funds to create a new Data Entry
//        MANAGE_DATA_INVALID_NAME = -4 // Name not a valid string
//    };
//
type ManageDataResultCode int32

const (
	ManageDataResultCodeManageDataSuccess         ManageDataResultCode = 0
	ManageDataResultCodeManageDataNotSupportedYet ManageDataResultCode = -1
	ManageDataResultCodeManageDataNameNotFound    ManageDataResultCode = -2
	ManageDataResultCodeManageDataLowReserve      ManageDataResultCode = -3
	ManageDataResultCodeManageDataInvalidName     ManageDataResultCode = -4
)

var manageDataResultCodeMap = map[int32]string{
	0:  "ManageDataResultCodeManageDataSuccess",
	-1: "ManageDataResultCodeManageDataNotSupportedYet",
	-2: "ManageDataResultCodeManageDataNameNotFound",
	-3: "ManageDataResultCodeManageDataLowReserve",
	-4: "ManageDataResultCodeManageDataInvalidName",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageDataResultCode
func (e ManageDataResultCode) ValidEnum(v int32) bool {
	_, ok := manageDataResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageDataResultCode) String() string {
	name, _ := manageDataResultCodeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e ManageDataResultCode) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := manageDataResultCodeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid ManageDataResultCode enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*ManageDataResultCode)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *ManageDataResultCode) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding ManageDataResultCode: %s", err)
	}
	if _, ok := manageDataResultCodeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid ManageDataResultCode enum value", v)
	}
	*e = ManageDataResultCode(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ManageDataResultCode) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ManageDataResultCode) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ManageDataResultCode)(nil)
	_ encoding.BinaryUnmarshaler = (*ManageDataResultCode)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ManageDataResultCode) xdrType() {}

var _ xdrType = (*ManageDataResultCode)(nil)

// ManageDataResult is an XDR Union defines as:
//
//   union ManageDataResult switch (ManageDataResultCode code)
//    {
//    case MANAGE_DATA_SUCCESS:
//        void;
//    default:
//        void;
//    };
//
type ManageDataResult struct {
	Code ManageDataResultCode
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageDataResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageDataResult
func (u ManageDataResult) ArmForSwitch(sw int32) (string, bool) {
	switch ManageDataResultCode(sw) {
	case ManageDataResultCodeManageDataSuccess:
		return "", true
	default:
		return "", true
	}
}

// NewManageDataResult creates a new  ManageDataResult.
func NewManageDataResult(code ManageDataResultCode, value interface{}) (result ManageDataResult, err error) {
	result.Code = code
	switch ManageDataResultCode(code) {
	case ManageDataResultCodeManageDataSuccess:
		// void
	default:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u ManageDataResult) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Code.EncodeTo(e); err != nil {
		return err
	}
	switch ManageDataResultCode(u.Code) {
	case ManageDataResultCodeManageDataSuccess:
		// Void
		return nil
	default:
		// Void
		return nil
	}
}

var _ decoderFrom = (*ManageDataResult)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *ManageDataResult) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Code.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ManageDataResultCode: %s", err)
	}
	switch ManageDataResultCode(u.Code) {
	case ManageDataResultCodeManageDataSuccess:
		// Void
		return n, nil
	default:
		// Void
		return n, nil
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ManageDataResult) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ManageDataResult) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ManageDataResult)(nil)
	_ encoding.BinaryUnmarshaler = (*ManageDataResult)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ManageDataResult) xdrType() {}

var _ xdrType = (*ManageDataResult)(nil)

// BumpSequenceResultCode is an XDR Enum defines as:
//
//   enum BumpSequenceResultCode
//    {
//        // codes considered as "success" for the operation
//        BUMP_SEQUENCE_SUCCESS = 0,
//        // codes considered as "failure" for the operation
//        BUMP_SEQUENCE_BAD_SEQ = -1 // `bumpTo` is not within bounds
//    };
//
type BumpSequenceResultCode int32

const (
	BumpSequenceResultCodeBumpSequenceSuccess BumpSequenceResultCode = 0
	BumpSequenceResultCodeBumpSequenceBadSeq  BumpSequenceResultCode = -1
)

var bumpSequenceResultCodeMap = map[int32]string{
	0:  "BumpSequenceResultCodeBumpSequenceSuccess",
	-1: "BumpSequenceResultCodeBumpSequenceBadSeq",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for BumpSequenceResultCode
func (e BumpSequenceResultCode) ValidEnum(v int32) bool {
	_, ok := bumpSequenceResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e BumpSequenceResultCode) String() string {
	name, _ := bumpSequenceResultCodeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e BumpSequenceResultCode) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := bumpSequenceResultCodeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid BumpSequenceResultCode enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*BumpSequenceResultCode)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *BumpSequenceResultCode) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding BumpSequenceResultCode: %s", err)
	}
	if _, ok := bumpSequenceResultCodeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid BumpSequenceResultCode enum value", v)
	}
	*e = BumpSequenceResultCode(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s BumpSequenceResultCode) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *BumpSequenceResultCode) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*BumpSequenceResultCode)(nil)
	_ encoding.BinaryUnmarshaler = (*BumpSequenceResultCode)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s BumpSequenceResultCode) xdrType() {}

var _ xdrType = (*BumpSequenceResultCode)(nil)

// BumpSequenceResult is an XDR Union defines as:
//
//   union BumpSequenceResult switch (BumpSequenceResultCode code)
//    {
//    case BUMP_SEQUENCE_SUCCESS:
//        void;
//    default:
//        void;
//    };
//
type BumpSequenceResult struct {
	Code BumpSequenceResultCode
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u BumpSequenceResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of BumpSequenceResult
func (u BumpSequenceResult) ArmForSwitch(sw int32) (string, bool) {
	switch BumpSequenceResultCode(sw) {
	case BumpSequenceResultCodeBumpSequenceSuccess:
		return "", true
	default:
		return "", true
	}
}

// NewBumpSequenceResult creates a new  BumpSequenceResult.
func NewBumpSequenceResult(code BumpSequenceResultCode, value interface{}) (result BumpSequenceResult, err error) {
	result.Code = code
	switch BumpSequenceResultCode(code) {
	case BumpSequenceResultCodeBumpSequenceSuccess:
		// void
	default:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u BumpSequenceResult) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Code.EncodeTo(e); err != nil {
		return err
	}
	switch BumpSequenceResultCode(u.Code) {
	case BumpSequenceResultCodeBumpSequenceSuccess:
		// Void
		return nil
	default:
		// Void
		return nil
	}
}

var _ decoderFrom = (*BumpSequenceResult)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *BumpSequenceResult) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Code.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding BumpSequenceResultCode: %s", err)
	}
	switch BumpSequenceResultCode(u.Code) {
	case BumpSequenceResultCodeBumpSequenceSuccess:
		// Void
		return n, nil
	default:
		// Void
		return n, nil
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s BumpSequenceResult) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *BumpSequenceResult) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*BumpSequenceResult)(nil)
	_ encoding.BinaryUnmarshaler = (*BumpSequenceResult)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s BumpSequenceResult) xdrType() {}

var _ xdrType = (*BumpSequenceResult)(nil)

// CreateClaimableBalanceResultCode is an XDR Enum defines as:
//
//   enum CreateClaimableBalanceResultCode
//    {
//        CREATE_CLAIMABLE_BALANCE_SUCCESS = 0,
//        CREATE_CLAIMABLE_BALANCE_MALFORMED = -1,
//        CREATE_CLAIMABLE_BALANCE_LOW_RESERVE = -2,
//        CREATE_CLAIMABLE_BALANCE_NO_TRUST = -3,
//        CREATE_CLAIMABLE_BALANCE_NOT_AUTHORIZED = -4,
//        CREATE_CLAIMABLE_BALANCE_UNDERFUNDED = -5
//    };
//
type CreateClaimableBalanceResultCode int32

const (
	CreateClaimableBalanceResultCodeCreateClaimableBalanceSuccess       CreateClaimableBalanceResultCode = 0
	CreateClaimableBalanceResultCodeCreateClaimableBalanceMalformed     CreateClaimableBalanceResultCode = -1
	CreateClaimableBalanceResultCodeCreateClaimableBalanceLowReserve    CreateClaimableBalanceResultCode = -2
	CreateClaimableBalanceResultCodeCreateClaimableBalanceNoTrust       CreateClaimableBalanceResultCode = -3
	CreateClaimableBalanceResultCodeCreateClaimableBalanceNotAuthorized CreateClaimableBalanceResultCode = -4
	CreateClaimableBalanceResultCodeCreateClaimableBalanceUnderfunded   CreateClaimableBalanceResultCode = -5
)

var createClaimableBalanceResultCodeMap = map[int32]string{
	0:  "CreateClaimableBalanceResultCodeCreateClaimableBalanceSuccess",
	-1: "CreateClaimableBalanceResultCodeCreateClaimableBalanceMalformed",
	-2: "CreateClaimableBalanceResultCodeCreateClaimableBalanceLowReserve",
	-3: "CreateClaimableBalanceResultCodeCreateClaimableBalanceNoTrust",
	-4: "CreateClaimableBalanceResultCodeCreateClaimableBalanceNotAuthorized",
	-5: "CreateClaimableBalanceResultCodeCreateClaimableBalanceUnderfunded",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for CreateClaimableBalanceResultCode
func (e CreateClaimableBalanceResultCode) ValidEnum(v int32) bool {
	_, ok := createClaimableBalanceResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e CreateClaimableBalanceResultCode) String() string {
	name, _ := createClaimableBalanceResultCodeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e CreateClaimableBalanceResultCode) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := createClaimableBalanceResultCodeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid CreateClaimableBalanceResultCode enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*CreateClaimableBalanceResultCode)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *CreateClaimableBalanceResultCode) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding CreateClaimableBalanceResultCode: %s", err)
	}
	if _, ok := createClaimableBalanceResultCodeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid CreateClaimableBalanceResultCode enum value", v)
	}
	*e = CreateClaimableBalanceResultCode(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s CreateClaimableBalanceResultCode) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *CreateClaimableBalanceResultCode) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*CreateClaimableBalanceResultCode)(nil)
	_ encoding.BinaryUnmarshaler = (*CreateClaimableBalanceResultCode)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s CreateClaimableBalanceResultCode) xdrType() {}

var _ xdrType = (*CreateClaimableBalanceResultCode)(nil)

// CreateClaimableBalanceResult is an XDR Union defines as:
//
//   union CreateClaimableBalanceResult switch (
//        CreateClaimableBalanceResultCode code)
//    {
//    case CREATE_CLAIMABLE_BALANCE_SUCCESS:
//        ClaimableBalanceID balanceID;
//    default:
//        void;
//    };
//
type CreateClaimableBalanceResult struct {
	Code      CreateClaimableBalanceResultCode
	BalanceId *ClaimableBalanceId
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u CreateClaimableBalanceResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of CreateClaimableBalanceResult
func (u CreateClaimableBalanceResult) ArmForSwitch(sw int32) (string, bool) {
	switch CreateClaimableBalanceResultCode(sw) {
	case CreateClaimableBalanceResultCodeCreateClaimableBalanceSuccess:
		return "BalanceId", true
	default:
		return "", true
	}
}

// NewCreateClaimableBalanceResult creates a new  CreateClaimableBalanceResult.
func NewCreateClaimableBalanceResult(code CreateClaimableBalanceResultCode, value interface{}) (result CreateClaimableBalanceResult, err error) {
	result.Code = code
	switch CreateClaimableBalanceResultCode(code) {
	case CreateClaimableBalanceResultCodeCreateClaimableBalanceSuccess:
		tv, ok := value.(ClaimableBalanceId)
		if !ok {
			err = fmt.Errorf("invalid value, must be ClaimableBalanceId")
			return
		}
		result.BalanceId = &tv
	default:
		// void
	}
	return
}

// MustBalanceId retrieves the BalanceId value from the union,
// panicing if the value is not set.
func (u CreateClaimableBalanceResult) MustBalanceId() ClaimableBalanceId {
	val, ok := u.GetBalanceId()

	if !ok {
		panic("arm BalanceId is not set")
	}

	return val
}

// GetBalanceId retrieves the BalanceId value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u CreateClaimableBalanceResult) GetBalanceId() (result ClaimableBalanceId, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "BalanceId" {
		result = *u.BalanceId
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u CreateClaimableBalanceResult) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Code.EncodeTo(e); err != nil {
		return err
	}
	switch CreateClaimableBalanceResultCode(u.Code) {
	case CreateClaimableBalanceResultCodeCreateClaimableBalanceSuccess:
		if err = (*u.BalanceId).EncodeTo(e); err != nil {
			return err
		}
		return nil
	default:
		// Void
		return nil
	}
}

var _ decoderFrom = (*CreateClaimableBalanceResult)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *CreateClaimableBalanceResult) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Code.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding CreateClaimableBalanceResultCode: %s", err)
	}
	switch CreateClaimableBalanceResultCode(u.Code) {
	case CreateClaimableBalanceResultCodeCreateClaimableBalanceSuccess:
		u.BalanceId = new(ClaimableBalanceId)
		nTmp, err = (*u.BalanceId).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ClaimableBalanceId: %s", err)
		}
		return n, nil
	default:
		// Void
		return n, nil
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s CreateClaimableBalanceResult) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *CreateClaimableBalanceResult) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*CreateClaimableBalanceResult)(nil)
	_ encoding.BinaryUnmarshaler = (*CreateClaimableBalanceResult)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s CreateClaimableBalanceResult) xdrType() {}

var _ xdrType = (*CreateClaimableBalanceResult)(nil)

// ClaimClaimableBalanceResultCode is an XDR Enum defines as:
//
//   enum ClaimClaimableBalanceResultCode
//    {
//        CLAIM_CLAIMABLE_BALANCE_SUCCESS = 0,
//        CLAIM_CLAIMABLE_BALANCE_DOES_NOT_EXIST = -1,
//        CLAIM_CLAIMABLE_BALANCE_CANNOT_CLAIM = -2,
//        CLAIM_CLAIMABLE_BALANCE_LINE_FULL = -3,
//        CLAIM_CLAIMABLE_BALANCE_NO_TRUST = -4,
//        CLAIM_CLAIMABLE_BALANCE_NOT_AUTHORIZED = -5
//
//    };
//
type ClaimClaimableBalanceResultCode int32

const (
	ClaimClaimableBalanceResultCodeClaimClaimableBalanceSuccess       ClaimClaimableBalanceResultCode = 0
	ClaimClaimableBalanceResultCodeClaimClaimableBalanceDoesNotExist  ClaimClaimableBalanceResultCode = -1
	ClaimClaimableBalanceResultCodeClaimClaimableBalanceCannotClaim   ClaimClaimableBalanceResultCode = -2
	ClaimClaimableBalanceResultCodeClaimClaimableBalanceLineFull      ClaimClaimableBalanceResultCode = -3
	ClaimClaimableBalanceResultCodeClaimClaimableBalanceNoTrust       ClaimClaimableBalanceResultCode = -4
	ClaimClaimableBalanceResultCodeClaimClaimableBalanceNotAuthorized ClaimClaimableBalanceResultCode = -5
)

var claimClaimableBalanceResultCodeMap = map[int32]string{
	0:  "ClaimClaimableBalanceResultCodeClaimClaimableBalanceSuccess",
	-1: "ClaimClaimableBalanceResultCodeClaimClaimableBalanceDoesNotExist",
	-2: "ClaimClaimableBalanceResultCodeClaimClaimableBalanceCannotClaim",
	-3: "ClaimClaimableBalanceResultCodeClaimClaimableBalanceLineFull",
	-4: "ClaimClaimableBalanceResultCodeClaimClaimableBalanceNoTrust",
	-5: "ClaimClaimableBalanceResultCodeClaimClaimableBalanceNotAuthorized",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ClaimClaimableBalanceResultCode
func (e ClaimClaimableBalanceResultCode) ValidEnum(v int32) bool {
	_, ok := claimClaimableBalanceResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ClaimClaimableBalanceResultCode) String() string {
	name, _ := claimClaimableBalanceResultCodeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e ClaimClaimableBalanceResultCode) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := claimClaimableBalanceResultCodeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid ClaimClaimableBalanceResultCode enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*ClaimClaimableBalanceResultCode)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *ClaimClaimableBalanceResultCode) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding ClaimClaimableBalanceResultCode: %s", err)
	}
	if _, ok := claimClaimableBalanceResultCodeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid ClaimClaimableBalanceResultCode enum value", v)
	}
	*e = ClaimClaimableBalanceResultCode(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ClaimClaimableBalanceResultCode) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ClaimClaimableBalanceResultCode) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ClaimClaimableBalanceResultCode)(nil)
	_ encoding.BinaryUnmarshaler = (*ClaimClaimableBalanceResultCode)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ClaimClaimableBalanceResultCode) xdrType() {}

var _ xdrType = (*ClaimClaimableBalanceResultCode)(nil)

// ClaimClaimableBalanceResult is an XDR Union defines as:
//
//   union ClaimClaimableBalanceResult switch (ClaimClaimableBalanceResultCode code)
//    {
//    case CLAIM_CLAIMABLE_BALANCE_SUCCESS:
//        void;
//    default:
//        void;
//    };
//
type ClaimClaimableBalanceResult struct {
	Code ClaimClaimableBalanceResultCode
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ClaimClaimableBalanceResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ClaimClaimableBalanceResult
func (u ClaimClaimableBalanceResult) ArmForSwitch(sw int32) (string, bool) {
	switch ClaimClaimableBalanceResultCode(sw) {
	case ClaimClaimableBalanceResultCodeClaimClaimableBalanceSuccess:
		return "", true
	default:
		return "", true
	}
}

// NewClaimClaimableBalanceResult creates a new  ClaimClaimableBalanceResult.
func NewClaimClaimableBalanceResult(code ClaimClaimableBalanceResultCode, value interface{}) (result ClaimClaimableBalanceResult, err error) {
	result.Code = code
	switch ClaimClaimableBalanceResultCode(code) {
	case ClaimClaimableBalanceResultCodeClaimClaimableBalanceSuccess:
		// void
	default:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u ClaimClaimableBalanceResult) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Code.EncodeTo(e); err != nil {
		return err
	}
	switch ClaimClaimableBalanceResultCode(u.Code) {
	case ClaimClaimableBalanceResultCodeClaimClaimableBalanceSuccess:
		// Void
		return nil
	default:
		// Void
		return nil
	}
}

var _ decoderFrom = (*ClaimClaimableBalanceResult)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *ClaimClaimableBalanceResult) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Code.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ClaimClaimableBalanceResultCode: %s", err)
	}
	switch ClaimClaimableBalanceResultCode(u.Code) {
	case ClaimClaimableBalanceResultCodeClaimClaimableBalanceSuccess:
		// Void
		return n, nil
	default:
		// Void
		return n, nil
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ClaimClaimableBalanceResult) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ClaimClaimableBalanceResult) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ClaimClaimableBalanceResult)(nil)
	_ encoding.BinaryUnmarshaler = (*ClaimClaimableBalanceResult)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ClaimClaimableBalanceResult) xdrType() {}

var _ xdrType = (*ClaimClaimableBalanceResult)(nil)

// BeginSponsoringFutureReservesResultCode is an XDR Enum defines as:
//
//   enum BeginSponsoringFutureReservesResultCode
//    {
//        // codes considered as "success" for the operation
//        BEGIN_SPONSORING_FUTURE_RESERVES_SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//        BEGIN_SPONSORING_FUTURE_RESERVES_MALFORMED = -1,
//        BEGIN_SPONSORING_FUTURE_RESERVES_ALREADY_SPONSORED = -2,
//        BEGIN_SPONSORING_FUTURE_RESERVES_RECURSIVE = -3
//    };
//
type BeginSponsoringFutureReservesResultCode int32

const (
	BeginSponsoringFutureReservesResultCodeBeginSponsoringFutureReservesSuccess          BeginSponsoringFutureReservesResultCode = 0
	BeginSponsoringFutureReservesResultCodeBeginSponsoringFutureReservesMalformed        BeginSponsoringFutureReservesResultCode = -1
	BeginSponsoringFutureReservesResultCodeBeginSponsoringFutureReservesAlreadySponsored BeginSponsoringFutureReservesResultCode = -2
	BeginSponsoringFutureReservesResultCodeBeginSponsoringFutureReservesRecursive        BeginSponsoringFutureReservesResultCode = -3
)

var beginSponsoringFutureReservesResultCodeMap = map[int32]string{
	0:  "BeginSponsoringFutureReservesResultCodeBeginSponsoringFutureReservesSuccess",
	-1: "BeginSponsoringFutureReservesResultCodeBeginSponsoringFutureReservesMalformed",
	-2: "BeginSponsoringFutureReservesResultCodeBeginSponsoringFutureReservesAlreadySponsored",
	-3: "BeginSponsoringFutureReservesResultCodeBeginSponsoringFutureReservesRecursive",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for BeginSponsoringFutureReservesResultCode
func (e BeginSponsoringFutureReservesResultCode) ValidEnum(v int32) bool {
	_, ok := beginSponsoringFutureReservesResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e BeginSponsoringFutureReservesResultCode) String() string {
	name, _ := beginSponsoringFutureReservesResultCodeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e BeginSponsoringFutureReservesResultCode) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := beginSponsoringFutureReservesResultCodeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid BeginSponsoringFutureReservesResultCode enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*BeginSponsoringFutureReservesResultCode)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *BeginSponsoringFutureReservesResultCode) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding BeginSponsoringFutureReservesResultCode: %s", err)
	}
	if _, ok := beginSponsoringFutureReservesResultCodeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid BeginSponsoringFutureReservesResultCode enum value", v)
	}
	*e = BeginSponsoringFutureReservesResultCode(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s BeginSponsoringFutureReservesResultCode) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *BeginSponsoringFutureReservesResultCode) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*BeginSponsoringFutureReservesResultCode)(nil)
	_ encoding.BinaryUnmarshaler = (*BeginSponsoringFutureReservesResultCode)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s BeginSponsoringFutureReservesResultCode) xdrType() {}

var _ xdrType = (*BeginSponsoringFutureReservesResultCode)(nil)

// BeginSponsoringFutureReservesResult is an XDR Union defines as:
//
//   union BeginSponsoringFutureReservesResult switch (
//        BeginSponsoringFutureReservesResultCode code)
//    {
//    case BEGIN_SPONSORING_FUTURE_RESERVES_SUCCESS:
//        void;
//    default:
//        void;
//    };
//
type BeginSponsoringFutureReservesResult struct {
	Code BeginSponsoringFutureReservesResultCode
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u BeginSponsoringFutureReservesResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of BeginSponsoringFutureReservesResult
func (u BeginSponsoringFutureReservesResult) ArmForSwitch(sw int32) (string, bool) {
	switch BeginSponsoringFutureReservesResultCode(sw) {
	case BeginSponsoringFutureReservesResultCodeBeginSponsoringFutureReservesSuccess:
		return "", true
	default:
		return "", true
	}
}

// NewBeginSponsoringFutureReservesResult creates a new  BeginSponsoringFutureReservesResult.
func NewBeginSponsoringFutureReservesResult(code BeginSponsoringFutureReservesResultCode, value interface{}) (result BeginSponsoringFutureReservesResult, err error) {
	result.Code = code
	switch BeginSponsoringFutureReservesResultCode(code) {
	case BeginSponsoringFutureReservesResultCodeBeginSponsoringFutureReservesSuccess:
		// void
	default:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u BeginSponsoringFutureReservesResult) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Code.EncodeTo(e); err != nil {
		return err
	}
	switch BeginSponsoringFutureReservesResultCode(u.Code) {
	case BeginSponsoringFutureReservesResultCodeBeginSponsoringFutureReservesSuccess:
		// Void
		return nil
	default:
		// Void
		return nil
	}
}

var _ decoderFrom = (*BeginSponsoringFutureReservesResult)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *BeginSponsoringFutureReservesResult) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Code.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding BeginSponsoringFutureReservesResultCode: %s", err)
	}
	switch BeginSponsoringFutureReservesResultCode(u.Code) {
	case BeginSponsoringFutureReservesResultCodeBeginSponsoringFutureReservesSuccess:
		// Void
		return n, nil
	default:
		// Void
		return n, nil
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s BeginSponsoringFutureReservesResult) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *BeginSponsoringFutureReservesResult) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*BeginSponsoringFutureReservesResult)(nil)
	_ encoding.BinaryUnmarshaler = (*BeginSponsoringFutureReservesResult)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s BeginSponsoringFutureReservesResult) xdrType() {}

var _ xdrType = (*BeginSponsoringFutureReservesResult)(nil)

// EndSponsoringFutureReservesResultCode is an XDR Enum defines as:
//
//   enum EndSponsoringFutureReservesResultCode
//    {
//        // codes considered as "success" for the operation
//        END_SPONSORING_FUTURE_RESERVES_SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//        END_SPONSORING_FUTURE_RESERVES_NOT_SPONSORED = -1
//    };
//
type EndSponsoringFutureReservesResultCode int32

const (
	EndSponsoringFutureReservesResultCodeEndSponsoringFutureReservesSuccess      EndSponsoringFutureReservesResultCode = 0
	EndSponsoringFutureReservesResultCodeEndSponsoringFutureReservesNotSponsored EndSponsoringFutureReservesResultCode = -1
)

var endSponsoringFutureReservesResultCodeMap = map[int32]string{
	0:  "EndSponsoringFutureReservesResultCodeEndSponsoringFutureReservesSuccess",
	-1: "EndSponsoringFutureReservesResultCodeEndSponsoringFutureReservesNotSponsored",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for EndSponsoringFutureReservesResultCode
func (e EndSponsoringFutureReservesResultCode) ValidEnum(v int32) bool {
	_, ok := endSponsoringFutureReservesResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e EndSponsoringFutureReservesResultCode) String() string {
	name, _ := endSponsoringFutureReservesResultCodeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e EndSponsoringFutureReservesResultCode) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := endSponsoringFutureReservesResultCodeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid EndSponsoringFutureReservesResultCode enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*EndSponsoringFutureReservesResultCode)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *EndSponsoringFutureReservesResultCode) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding EndSponsoringFutureReservesResultCode: %s", err)
	}
	if _, ok := endSponsoringFutureReservesResultCodeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid EndSponsoringFutureReservesResultCode enum value", v)
	}
	*e = EndSponsoringFutureReservesResultCode(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s EndSponsoringFutureReservesResultCode) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *EndSponsoringFutureReservesResultCode) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*EndSponsoringFutureReservesResultCode)(nil)
	_ encoding.BinaryUnmarshaler = (*EndSponsoringFutureReservesResultCode)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s EndSponsoringFutureReservesResultCode) xdrType() {}

var _ xdrType = (*EndSponsoringFutureReservesResultCode)(nil)

// EndSponsoringFutureReservesResult is an XDR Union defines as:
//
//   union EndSponsoringFutureReservesResult switch (
//        EndSponsoringFutureReservesResultCode code)
//    {
//    case END_SPONSORING_FUTURE_RESERVES_SUCCESS:
//        void;
//    default:
//        void;
//    };
//
type EndSponsoringFutureReservesResult struct {
	Code EndSponsoringFutureReservesResultCode
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u EndSponsoringFutureReservesResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of EndSponsoringFutureReservesResult
func (u EndSponsoringFutureReservesResult) ArmForSwitch(sw int32) (string, bool) {
	switch EndSponsoringFutureReservesResultCode(sw) {
	case EndSponsoringFutureReservesResultCodeEndSponsoringFutureReservesSuccess:
		return "", true
	default:
		return "", true
	}
}

// NewEndSponsoringFutureReservesResult creates a new  EndSponsoringFutureReservesResult.
func NewEndSponsoringFutureReservesResult(code EndSponsoringFutureReservesResultCode, value interface{}) (result EndSponsoringFutureReservesResult, err error) {
	result.Code = code
	switch EndSponsoringFutureReservesResultCode(code) {
	case EndSponsoringFutureReservesResultCodeEndSponsoringFutureReservesSuccess:
		// void
	default:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u EndSponsoringFutureReservesResult) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Code.EncodeTo(e); err != nil {
		return err
	}
	switch EndSponsoringFutureReservesResultCode(u.Code) {
	case EndSponsoringFutureReservesResultCodeEndSponsoringFutureReservesSuccess:
		// Void
		return nil
	default:
		// Void
		return nil
	}
}

var _ decoderFrom = (*EndSponsoringFutureReservesResult)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *EndSponsoringFutureReservesResult) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Code.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding EndSponsoringFutureReservesResultCode: %s", err)
	}
	switch EndSponsoringFutureReservesResultCode(u.Code) {
	case EndSponsoringFutureReservesResultCodeEndSponsoringFutureReservesSuccess:
		// Void
		return n, nil
	default:
		// Void
		return n, nil
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s EndSponsoringFutureReservesResult) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *EndSponsoringFutureReservesResult) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*EndSponsoringFutureReservesResult)(nil)
	_ encoding.BinaryUnmarshaler = (*EndSponsoringFutureReservesResult)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s EndSponsoringFutureReservesResult) xdrType() {}

var _ xdrType = (*EndSponsoringFutureReservesResult)(nil)

// RevokeSponsorshipResultCode is an XDR Enum defines as:
//
//   enum RevokeSponsorshipResultCode
//    {
//        // codes considered as "success" for the operation
//        REVOKE_SPONSORSHIP_SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//        REVOKE_SPONSORSHIP_DOES_NOT_EXIST = -1,
//        REVOKE_SPONSORSHIP_NOT_SPONSOR = -2,
//        REVOKE_SPONSORSHIP_LOW_RESERVE = -3,
//        REVOKE_SPONSORSHIP_ONLY_TRANSFERABLE = -4,
//        REVOKE_SPONSORSHIP_MALFORMED = -5
//    };
//
type RevokeSponsorshipResultCode int32

const (
	RevokeSponsorshipResultCodeRevokeSponsorshipSuccess          RevokeSponsorshipResultCode = 0
	RevokeSponsorshipResultCodeRevokeSponsorshipDoesNotExist     RevokeSponsorshipResultCode = -1
	RevokeSponsorshipResultCodeRevokeSponsorshipNotSponsor       RevokeSponsorshipResultCode = -2
	RevokeSponsorshipResultCodeRevokeSponsorshipLowReserve       RevokeSponsorshipResultCode = -3
	RevokeSponsorshipResultCodeRevokeSponsorshipOnlyTransferable RevokeSponsorshipResultCode = -4
	RevokeSponsorshipResultCodeRevokeSponsorshipMalformed        RevokeSponsorshipResultCode = -5
)

var revokeSponsorshipResultCodeMap = map[int32]string{
	0:  "RevokeSponsorshipResultCodeRevokeSponsorshipSuccess",
	-1: "RevokeSponsorshipResultCodeRevokeSponsorshipDoesNotExist",
	-2: "RevokeSponsorshipResultCodeRevokeSponsorshipNotSponsor",
	-3: "RevokeSponsorshipResultCodeRevokeSponsorshipLowReserve",
	-4: "RevokeSponsorshipResultCodeRevokeSponsorshipOnlyTransferable",
	-5: "RevokeSponsorshipResultCodeRevokeSponsorshipMalformed",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for RevokeSponsorshipResultCode
func (e RevokeSponsorshipResultCode) ValidEnum(v int32) bool {
	_, ok := revokeSponsorshipResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e RevokeSponsorshipResultCode) String() string {
	name, _ := revokeSponsorshipResultCodeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e RevokeSponsorshipResultCode) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := revokeSponsorshipResultCodeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid RevokeSponsorshipResultCode enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*RevokeSponsorshipResultCode)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *RevokeSponsorshipResultCode) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding RevokeSponsorshipResultCode: %s", err)
	}
	if _, ok := revokeSponsorshipResultCodeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid RevokeSponsorshipResultCode enum value", v)
	}
	*e = RevokeSponsorshipResultCode(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s RevokeSponsorshipResultCode) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *RevokeSponsorshipResultCode) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*RevokeSponsorshipResultCode)(nil)
	_ encoding.BinaryUnmarshaler = (*RevokeSponsorshipResultCode)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s RevokeSponsorshipResultCode) xdrType() {}

var _ xdrType = (*RevokeSponsorshipResultCode)(nil)

// RevokeSponsorshipResult is an XDR Union defines as:
//
//   union RevokeSponsorshipResult switch (RevokeSponsorshipResultCode code)
//    {
//    case REVOKE_SPONSORSHIP_SUCCESS:
//        void;
//    default:
//        void;
//    };
//
type RevokeSponsorshipResult struct {
	Code RevokeSponsorshipResultCode
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u RevokeSponsorshipResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of RevokeSponsorshipResult
func (u RevokeSponsorshipResult) ArmForSwitch(sw int32) (string, bool) {
	switch RevokeSponsorshipResultCode(sw) {
	case RevokeSponsorshipResultCodeRevokeSponsorshipSuccess:
		return "", true
	default:
		return "", true
	}
}

// NewRevokeSponsorshipResult creates a new  RevokeSponsorshipResult.
func NewRevokeSponsorshipResult(code RevokeSponsorshipResultCode, value interface{}) (result RevokeSponsorshipResult, err error) {
	result.Code = code
	switch RevokeSponsorshipResultCode(code) {
	case RevokeSponsorshipResultCodeRevokeSponsorshipSuccess:
		// void
	default:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u RevokeSponsorshipResult) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Code.EncodeTo(e); err != nil {
		return err
	}
	switch RevokeSponsorshipResultCode(u.Code) {
	case RevokeSponsorshipResultCodeRevokeSponsorshipSuccess:
		// Void
		return nil
	default:
		// Void
		return nil
	}
}

var _ decoderFrom = (*RevokeSponsorshipResult)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *RevokeSponsorshipResult) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Code.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding RevokeSponsorshipResultCode: %s", err)
	}
	switch RevokeSponsorshipResultCode(u.Code) {
	case RevokeSponsorshipResultCodeRevokeSponsorshipSuccess:
		// Void
		return n, nil
	default:
		// Void
		return n, nil
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s RevokeSponsorshipResult) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *RevokeSponsorshipResult) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*RevokeSponsorshipResult)(nil)
	_ encoding.BinaryUnmarshaler = (*RevokeSponsorshipResult)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s RevokeSponsorshipResult) xdrType() {}

var _ xdrType = (*RevokeSponsorshipResult)(nil)

// ClawbackResultCode is an XDR Enum defines as:
//
//   enum ClawbackResultCode
//    {
//        // codes considered as "success" for the operation
//        CLAWBACK_SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//        CLAWBACK_MALFORMED = -1,
//        CLAWBACK_NOT_CLAWBACK_ENABLED = -2,
//        CLAWBACK_NO_TRUST = -3,
//        CLAWBACK_UNDERFUNDED = -4
//    };
//
type ClawbackResultCode int32

const (
	ClawbackResultCodeClawbackSuccess            ClawbackResultCode = 0
	ClawbackResultCodeClawbackMalformed          ClawbackResultCode = -1
	ClawbackResultCodeClawbackNotClawbackEnabled ClawbackResultCode = -2
	ClawbackResultCodeClawbackNoTrust            ClawbackResultCode = -3
	ClawbackResultCodeClawbackUnderfunded        ClawbackResultCode = -4
)

var clawbackResultCodeMap = map[int32]string{
	0:  "ClawbackResultCodeClawbackSuccess",
	-1: "ClawbackResultCodeClawbackMalformed",
	-2: "ClawbackResultCodeClawbackNotClawbackEnabled",
	-3: "ClawbackResultCodeClawbackNoTrust",
	-4: "ClawbackResultCodeClawbackUnderfunded",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ClawbackResultCode
func (e ClawbackResultCode) ValidEnum(v int32) bool {
	_, ok := clawbackResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ClawbackResultCode) String() string {
	name, _ := clawbackResultCodeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e ClawbackResultCode) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := clawbackResultCodeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid ClawbackResultCode enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*ClawbackResultCode)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *ClawbackResultCode) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding ClawbackResultCode: %s", err)
	}
	if _, ok := clawbackResultCodeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid ClawbackResultCode enum value", v)
	}
	*e = ClawbackResultCode(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ClawbackResultCode) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ClawbackResultCode) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ClawbackResultCode)(nil)
	_ encoding.BinaryUnmarshaler = (*ClawbackResultCode)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ClawbackResultCode) xdrType() {}

var _ xdrType = (*ClawbackResultCode)(nil)

// ClawbackResult is an XDR Union defines as:
//
//   union ClawbackResult switch (ClawbackResultCode code)
//    {
//    case CLAWBACK_SUCCESS:
//        void;
//    default:
//        void;
//    };
//
type ClawbackResult struct {
	Code ClawbackResultCode
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ClawbackResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ClawbackResult
func (u ClawbackResult) ArmForSwitch(sw int32) (string, bool) {
	switch ClawbackResultCode(sw) {
	case ClawbackResultCodeClawbackSuccess:
		return "", true
	default:
		return "", true
	}
}

// NewClawbackResult creates a new  ClawbackResult.
func NewClawbackResult(code ClawbackResultCode, value interface{}) (result ClawbackResult, err error) {
	result.Code = code
	switch ClawbackResultCode(code) {
	case ClawbackResultCodeClawbackSuccess:
		// void
	default:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u ClawbackResult) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Code.EncodeTo(e); err != nil {
		return err
	}
	switch ClawbackResultCode(u.Code) {
	case ClawbackResultCodeClawbackSuccess:
		// Void
		return nil
	default:
		// Void
		return nil
	}
}

var _ decoderFrom = (*ClawbackResult)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *ClawbackResult) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Code.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ClawbackResultCode: %s", err)
	}
	switch ClawbackResultCode(u.Code) {
	case ClawbackResultCodeClawbackSuccess:
		// Void
		return n, nil
	default:
		// Void
		return n, nil
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ClawbackResult) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ClawbackResult) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ClawbackResult)(nil)
	_ encoding.BinaryUnmarshaler = (*ClawbackResult)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ClawbackResult) xdrType() {}

var _ xdrType = (*ClawbackResult)(nil)

// ClawbackClaimableBalanceResultCode is an XDR Enum defines as:
//
//   enum ClawbackClaimableBalanceResultCode
//    {
//        // codes considered as "success" for the operation
//        CLAWBACK_CLAIMABLE_BALANCE_SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//        CLAWBACK_CLAIMABLE_BALANCE_DOES_NOT_EXIST = -1,
//        CLAWBACK_CLAIMABLE_BALANCE_NOT_ISSUER = -2,
//        CLAWBACK_CLAIMABLE_BALANCE_NOT_CLAWBACK_ENABLED = -3
//    };
//
type ClawbackClaimableBalanceResultCode int32

const (
	ClawbackClaimableBalanceResultCodeClawbackClaimableBalanceSuccess            ClawbackClaimableBalanceResultCode = 0
	ClawbackClaimableBalanceResultCodeClawbackClaimableBalanceDoesNotExist       ClawbackClaimableBalanceResultCode = -1
	ClawbackClaimableBalanceResultCodeClawbackClaimableBalanceNotIssuer          ClawbackClaimableBalanceResultCode = -2
	ClawbackClaimableBalanceResultCodeClawbackClaimableBalanceNotClawbackEnabled ClawbackClaimableBalanceResultCode = -3
)

var clawbackClaimableBalanceResultCodeMap = map[int32]string{
	0:  "ClawbackClaimableBalanceResultCodeClawbackClaimableBalanceSuccess",
	-1: "ClawbackClaimableBalanceResultCodeClawbackClaimableBalanceDoesNotExist",
	-2: "ClawbackClaimableBalanceResultCodeClawbackClaimableBalanceNotIssuer",
	-3: "ClawbackClaimableBalanceResultCodeClawbackClaimableBalanceNotClawbackEnabled",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ClawbackClaimableBalanceResultCode
func (e ClawbackClaimableBalanceResultCode) ValidEnum(v int32) bool {
	_, ok := clawbackClaimableBalanceResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ClawbackClaimableBalanceResultCode) String() string {
	name, _ := clawbackClaimableBalanceResultCodeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e ClawbackClaimableBalanceResultCode) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := clawbackClaimableBalanceResultCodeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid ClawbackClaimableBalanceResultCode enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*ClawbackClaimableBalanceResultCode)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *ClawbackClaimableBalanceResultCode) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding ClawbackClaimableBalanceResultCode: %s", err)
	}
	if _, ok := clawbackClaimableBalanceResultCodeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid ClawbackClaimableBalanceResultCode enum value", v)
	}
	*e = ClawbackClaimableBalanceResultCode(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ClawbackClaimableBalanceResultCode) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ClawbackClaimableBalanceResultCode) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ClawbackClaimableBalanceResultCode)(nil)
	_ encoding.BinaryUnmarshaler = (*ClawbackClaimableBalanceResultCode)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ClawbackClaimableBalanceResultCode) xdrType() {}

var _ xdrType = (*ClawbackClaimableBalanceResultCode)(nil)

// ClawbackClaimableBalanceResult is an XDR Union defines as:
//
//   union ClawbackClaimableBalanceResult switch (
//        ClawbackClaimableBalanceResultCode code)
//    {
//    case CLAWBACK_CLAIMABLE_BALANCE_SUCCESS:
//        void;
//    default:
//        void;
//    };
//
type ClawbackClaimableBalanceResult struct {
	Code ClawbackClaimableBalanceResultCode
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ClawbackClaimableBalanceResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ClawbackClaimableBalanceResult
func (u ClawbackClaimableBalanceResult) ArmForSwitch(sw int32) (string, bool) {
	switch ClawbackClaimableBalanceResultCode(sw) {
	case ClawbackClaimableBalanceResultCodeClawbackClaimableBalanceSuccess:
		return "", true
	default:
		return "", true
	}
}

// NewClawbackClaimableBalanceResult creates a new  ClawbackClaimableBalanceResult.
func NewClawbackClaimableBalanceResult(code ClawbackClaimableBalanceResultCode, value interface{}) (result ClawbackClaimableBalanceResult, err error) {
	result.Code = code
	switch ClawbackClaimableBalanceResultCode(code) {
	case ClawbackClaimableBalanceResultCodeClawbackClaimableBalanceSuccess:
		// void
	default:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u ClawbackClaimableBalanceResult) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Code.EncodeTo(e); err != nil {
		return err
	}
	switch ClawbackClaimableBalanceResultCode(u.Code) {
	case ClawbackClaimableBalanceResultCodeClawbackClaimableBalanceSuccess:
		// Void
		return nil
	default:
		// Void
		return nil
	}
}

var _ decoderFrom = (*ClawbackClaimableBalanceResult)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *ClawbackClaimableBalanceResult) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Code.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ClawbackClaimableBalanceResultCode: %s", err)
	}
	switch ClawbackClaimableBalanceResultCode(u.Code) {
	case ClawbackClaimableBalanceResultCodeClawbackClaimableBalanceSuccess:
		// Void
		return n, nil
	default:
		// Void
		return n, nil
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ClawbackClaimableBalanceResult) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ClawbackClaimableBalanceResult) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ClawbackClaimableBalanceResult)(nil)
	_ encoding.BinaryUnmarshaler = (*ClawbackClaimableBalanceResult)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s ClawbackClaimableBalanceResult) xdrType() {}

var _ xdrType = (*ClawbackClaimableBalanceResult)(nil)

// SetTrustLineFlagsResultCode is an XDR Enum defines as:
//
//   enum SetTrustLineFlagsResultCode
//    {
//        // codes considered as "success" for the operation
//        SET_TRUST_LINE_FLAGS_SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//        SET_TRUST_LINE_FLAGS_MALFORMED = -1,
//        SET_TRUST_LINE_FLAGS_NO_TRUST_LINE = -2,
//        SET_TRUST_LINE_FLAGS_CANT_REVOKE = -3,
//        SET_TRUST_LINE_FLAGS_INVALID_STATE = -4,
//        SET_TRUST_LINE_FLAGS_LOW_RESERVE = -5 // claimable balances can't be created
//                                              // on revoke due to low reserves
//    };
//
type SetTrustLineFlagsResultCode int32

const (
	SetTrustLineFlagsResultCodeSetTrustLineFlagsSuccess      SetTrustLineFlagsResultCode = 0
	SetTrustLineFlagsResultCodeSetTrustLineFlagsMalformed    SetTrustLineFlagsResultCode = -1
	SetTrustLineFlagsResultCodeSetTrustLineFlagsNoTrustLine  SetTrustLineFlagsResultCode = -2
	SetTrustLineFlagsResultCodeSetTrustLineFlagsCantRevoke   SetTrustLineFlagsResultCode = -3
	SetTrustLineFlagsResultCodeSetTrustLineFlagsInvalidState SetTrustLineFlagsResultCode = -4
	SetTrustLineFlagsResultCodeSetTrustLineFlagsLowReserve   SetTrustLineFlagsResultCode = -5
)

var setTrustLineFlagsResultCodeMap = map[int32]string{
	0:  "SetTrustLineFlagsResultCodeSetTrustLineFlagsSuccess",
	-1: "SetTrustLineFlagsResultCodeSetTrustLineFlagsMalformed",
	-2: "SetTrustLineFlagsResultCodeSetTrustLineFlagsNoTrustLine",
	-3: "SetTrustLineFlagsResultCodeSetTrustLineFlagsCantRevoke",
	-4: "SetTrustLineFlagsResultCodeSetTrustLineFlagsInvalidState",
	-5: "SetTrustLineFlagsResultCodeSetTrustLineFlagsLowReserve",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for SetTrustLineFlagsResultCode
func (e SetTrustLineFlagsResultCode) ValidEnum(v int32) bool {
	_, ok := setTrustLineFlagsResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e SetTrustLineFlagsResultCode) String() string {
	name, _ := setTrustLineFlagsResultCodeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e SetTrustLineFlagsResultCode) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := setTrustLineFlagsResultCodeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid SetTrustLineFlagsResultCode enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*SetTrustLineFlagsResultCode)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *SetTrustLineFlagsResultCode) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding SetTrustLineFlagsResultCode: %s", err)
	}
	if _, ok := setTrustLineFlagsResultCodeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid SetTrustLineFlagsResultCode enum value", v)
	}
	*e = SetTrustLineFlagsResultCode(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s SetTrustLineFlagsResultCode) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SetTrustLineFlagsResultCode) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*SetTrustLineFlagsResultCode)(nil)
	_ encoding.BinaryUnmarshaler = (*SetTrustLineFlagsResultCode)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s SetTrustLineFlagsResultCode) xdrType() {}

var _ xdrType = (*SetTrustLineFlagsResultCode)(nil)

// SetTrustLineFlagsResult is an XDR Union defines as:
//
//   union SetTrustLineFlagsResult switch (SetTrustLineFlagsResultCode code)
//    {
//    case SET_TRUST_LINE_FLAGS_SUCCESS:
//        void;
//    default:
//        void;
//    };
//
type SetTrustLineFlagsResult struct {
	Code SetTrustLineFlagsResultCode
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u SetTrustLineFlagsResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of SetTrustLineFlagsResult
func (u SetTrustLineFlagsResult) ArmForSwitch(sw int32) (string, bool) {
	switch SetTrustLineFlagsResultCode(sw) {
	case SetTrustLineFlagsResultCodeSetTrustLineFlagsSuccess:
		return "", true
	default:
		return "", true
	}
}

// NewSetTrustLineFlagsResult creates a new  SetTrustLineFlagsResult.
func NewSetTrustLineFlagsResult(code SetTrustLineFlagsResultCode, value interface{}) (result SetTrustLineFlagsResult, err error) {
	result.Code = code
	switch SetTrustLineFlagsResultCode(code) {
	case SetTrustLineFlagsResultCodeSetTrustLineFlagsSuccess:
		// void
	default:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u SetTrustLineFlagsResult) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Code.EncodeTo(e); err != nil {
		return err
	}
	switch SetTrustLineFlagsResultCode(u.Code) {
	case SetTrustLineFlagsResultCodeSetTrustLineFlagsSuccess:
		// Void
		return nil
	default:
		// Void
		return nil
	}
}

var _ decoderFrom = (*SetTrustLineFlagsResult)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *SetTrustLineFlagsResult) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Code.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SetTrustLineFlagsResultCode: %s", err)
	}
	switch SetTrustLineFlagsResultCode(u.Code) {
	case SetTrustLineFlagsResultCodeSetTrustLineFlagsSuccess:
		// Void
		return n, nil
	default:
		// Void
		return n, nil
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s SetTrustLineFlagsResult) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SetTrustLineFlagsResult) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*SetTrustLineFlagsResult)(nil)
	_ encoding.BinaryUnmarshaler = (*SetTrustLineFlagsResult)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s SetTrustLineFlagsResult) xdrType() {}

var _ xdrType = (*SetTrustLineFlagsResult)(nil)

// LiquidityPoolDepositResultCode is an XDR Enum defines as:
//
//   enum LiquidityPoolDepositResultCode
//    {
//        // codes considered as "success" for the operation
//        LIQUIDITY_POOL_DEPOSIT_SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//        LIQUIDITY_POOL_DEPOSIT_MALFORMED = -1,      // bad input
//        LIQUIDITY_POOL_DEPOSIT_NO_TRUST = -2,       // no trust line for one of the
//                                                    // assets
//        LIQUIDITY_POOL_DEPOSIT_NOT_AUTHORIZED = -3, // not authorized for one of the
//                                                    // assets
//        LIQUIDITY_POOL_DEPOSIT_UNDERFUNDED = -4,    // not enough balance for one of
//                                                    // the assets
//        LIQUIDITY_POOL_DEPOSIT_LINE_FULL = -5,      // pool share trust line doesn't
//                                                    // have sufficient limit
//        LIQUIDITY_POOL_DEPOSIT_BAD_PRICE = -6,      // deposit price outside bounds
//        LIQUIDITY_POOL_DEPOSIT_POOL_FULL = -7       // pool reserves are full
//    };
//
type LiquidityPoolDepositResultCode int32

const (
	LiquidityPoolDepositResultCodeLiquidityPoolDepositSuccess       LiquidityPoolDepositResultCode = 0
	LiquidityPoolDepositResultCodeLiquidityPoolDepositMalformed     LiquidityPoolDepositResultCode = -1
	LiquidityPoolDepositResultCodeLiquidityPoolDepositNoTrust       LiquidityPoolDepositResultCode = -2
	LiquidityPoolDepositResultCodeLiquidityPoolDepositNotAuthorized LiquidityPoolDepositResultCode = -3
	LiquidityPoolDepositResultCodeLiquidityPoolDepositUnderfunded   LiquidityPoolDepositResultCode = -4
	LiquidityPoolDepositResultCodeLiquidityPoolDepositLineFull      LiquidityPoolDepositResultCode = -5
	LiquidityPoolDepositResultCodeLiquidityPoolDepositBadPrice      LiquidityPoolDepositResultCode = -6
	LiquidityPoolDepositResultCodeLiquidityPoolDepositPoolFull      LiquidityPoolDepositResultCode = -7
)

var liquidityPoolDepositResultCodeMap = map[int32]string{
	0:  "LiquidityPoolDepositResultCodeLiquidityPoolDepositSuccess",
	-1: "LiquidityPoolDepositResultCodeLiquidityPoolDepositMalformed",
	-2: "LiquidityPoolDepositResultCodeLiquidityPoolDepositNoTrust",
	-3: "LiquidityPoolDepositResultCodeLiquidityPoolDepositNotAuthorized",
	-4: "LiquidityPoolDepositResultCodeLiquidityPoolDepositUnderfunded",
	-5: "LiquidityPoolDepositResultCodeLiquidityPoolDepositLineFull",
	-6: "LiquidityPoolDepositResultCodeLiquidityPoolDepositBadPrice",
	-7: "LiquidityPoolDepositResultCodeLiquidityPoolDepositPoolFull",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for LiquidityPoolDepositResultCode
func (e LiquidityPoolDepositResultCode) ValidEnum(v int32) bool {
	_, ok := liquidityPoolDepositResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e LiquidityPoolDepositResultCode) String() string {
	name, _ := liquidityPoolDepositResultCodeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e LiquidityPoolDepositResultCode) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := liquidityPoolDepositResultCodeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid LiquidityPoolDepositResultCode enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*LiquidityPoolDepositResultCode)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *LiquidityPoolDepositResultCode) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding LiquidityPoolDepositResultCode: %s", err)
	}
	if _, ok := liquidityPoolDepositResultCodeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid LiquidityPoolDepositResultCode enum value", v)
	}
	*e = LiquidityPoolDepositResultCode(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LiquidityPoolDepositResultCode) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LiquidityPoolDepositResultCode) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LiquidityPoolDepositResultCode)(nil)
	_ encoding.BinaryUnmarshaler = (*LiquidityPoolDepositResultCode)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LiquidityPoolDepositResultCode) xdrType() {}

var _ xdrType = (*LiquidityPoolDepositResultCode)(nil)

// LiquidityPoolDepositResult is an XDR Union defines as:
//
//   union LiquidityPoolDepositResult switch (
//        LiquidityPoolDepositResultCode code)
//    {
//    case LIQUIDITY_POOL_DEPOSIT_SUCCESS:
//        void;
//    default:
//        void;
//    };
//
type LiquidityPoolDepositResult struct {
	Code LiquidityPoolDepositResultCode
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LiquidityPoolDepositResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LiquidityPoolDepositResult
func (u LiquidityPoolDepositResult) ArmForSwitch(sw int32) (string, bool) {
	switch LiquidityPoolDepositResultCode(sw) {
	case LiquidityPoolDepositResultCodeLiquidityPoolDepositSuccess:
		return "", true
	default:
		return "", true
	}
}

// NewLiquidityPoolDepositResult creates a new  LiquidityPoolDepositResult.
func NewLiquidityPoolDepositResult(code LiquidityPoolDepositResultCode, value interface{}) (result LiquidityPoolDepositResult, err error) {
	result.Code = code
	switch LiquidityPoolDepositResultCode(code) {
	case LiquidityPoolDepositResultCodeLiquidityPoolDepositSuccess:
		// void
	default:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u LiquidityPoolDepositResult) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Code.EncodeTo(e); err != nil {
		return err
	}
	switch LiquidityPoolDepositResultCode(u.Code) {
	case LiquidityPoolDepositResultCodeLiquidityPoolDepositSuccess:
		// Void
		return nil
	default:
		// Void
		return nil
	}
}

var _ decoderFrom = (*LiquidityPoolDepositResult)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *LiquidityPoolDepositResult) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Code.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LiquidityPoolDepositResultCode: %s", err)
	}
	switch LiquidityPoolDepositResultCode(u.Code) {
	case LiquidityPoolDepositResultCodeLiquidityPoolDepositSuccess:
		// Void
		return n, nil
	default:
		// Void
		return n, nil
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LiquidityPoolDepositResult) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LiquidityPoolDepositResult) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LiquidityPoolDepositResult)(nil)
	_ encoding.BinaryUnmarshaler = (*LiquidityPoolDepositResult)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LiquidityPoolDepositResult) xdrType() {}

var _ xdrType = (*LiquidityPoolDepositResult)(nil)

// LiquidityPoolWithdrawResultCode is an XDR Enum defines as:
//
//   enum LiquidityPoolWithdrawResultCode
//    {
//        // codes considered as "success" for the operation
//        LIQUIDITY_POOL_WITHDRAW_SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//        LIQUIDITY_POOL_WITHDRAW_MALFORMED = -1,      // bad input
//        LIQUIDITY_POOL_WITHDRAW_NO_TRUST = -2,       // no trust line for one of the
//                                                     // assets
//        LIQUIDITY_POOL_WITHDRAW_UNDERFUNDED = -3,    // not enough balance of the
//                                                     // pool share
//        LIQUIDITY_POOL_WITHDRAW_LINE_FULL = -4,      // would go above limit for one
//                                                     // of the assets
//        LIQUIDITY_POOL_WITHDRAW_UNDER_MINIMUM = -5   // didn't withdraw enough
//    };
//
type LiquidityPoolWithdrawResultCode int32

const (
	LiquidityPoolWithdrawResultCodeLiquidityPoolWithdrawSuccess      LiquidityPoolWithdrawResultCode = 0
	LiquidityPoolWithdrawResultCodeLiquidityPoolWithdrawMalformed    LiquidityPoolWithdrawResultCode = -1
	LiquidityPoolWithdrawResultCodeLiquidityPoolWithdrawNoTrust      LiquidityPoolWithdrawResultCode = -2
	LiquidityPoolWithdrawResultCodeLiquidityPoolWithdrawUnderfunded  LiquidityPoolWithdrawResultCode = -3
	LiquidityPoolWithdrawResultCodeLiquidityPoolWithdrawLineFull     LiquidityPoolWithdrawResultCode = -4
	LiquidityPoolWithdrawResultCodeLiquidityPoolWithdrawUnderMinimum LiquidityPoolWithdrawResultCode = -5
)

var liquidityPoolWithdrawResultCodeMap = map[int32]string{
	0:  "LiquidityPoolWithdrawResultCodeLiquidityPoolWithdrawSuccess",
	-1: "LiquidityPoolWithdrawResultCodeLiquidityPoolWithdrawMalformed",
	-2: "LiquidityPoolWithdrawResultCodeLiquidityPoolWithdrawNoTrust",
	-3: "LiquidityPoolWithdrawResultCodeLiquidityPoolWithdrawUnderfunded",
	-4: "LiquidityPoolWithdrawResultCodeLiquidityPoolWithdrawLineFull",
	-5: "LiquidityPoolWithdrawResultCodeLiquidityPoolWithdrawUnderMinimum",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for LiquidityPoolWithdrawResultCode
func (e LiquidityPoolWithdrawResultCode) ValidEnum(v int32) bool {
	_, ok := liquidityPoolWithdrawResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e LiquidityPoolWithdrawResultCode) String() string {
	name, _ := liquidityPoolWithdrawResultCodeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e LiquidityPoolWithdrawResultCode) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := liquidityPoolWithdrawResultCodeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid LiquidityPoolWithdrawResultCode enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*LiquidityPoolWithdrawResultCode)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *LiquidityPoolWithdrawResultCode) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding LiquidityPoolWithdrawResultCode: %s", err)
	}
	if _, ok := liquidityPoolWithdrawResultCodeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid LiquidityPoolWithdrawResultCode enum value", v)
	}
	*e = LiquidityPoolWithdrawResultCode(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LiquidityPoolWithdrawResultCode) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LiquidityPoolWithdrawResultCode) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LiquidityPoolWithdrawResultCode)(nil)
	_ encoding.BinaryUnmarshaler = (*LiquidityPoolWithdrawResultCode)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LiquidityPoolWithdrawResultCode) xdrType() {}

var _ xdrType = (*LiquidityPoolWithdrawResultCode)(nil)

// LiquidityPoolWithdrawResult is an XDR Union defines as:
//
//   union LiquidityPoolWithdrawResult switch (
//        LiquidityPoolWithdrawResultCode code)
//    {
//    case LIQUIDITY_POOL_WITHDRAW_SUCCESS:
//        void;
//    default:
//        void;
//    };
//
type LiquidityPoolWithdrawResult struct {
	Code LiquidityPoolWithdrawResultCode
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LiquidityPoolWithdrawResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LiquidityPoolWithdrawResult
func (u LiquidityPoolWithdrawResult) ArmForSwitch(sw int32) (string, bool) {
	switch LiquidityPoolWithdrawResultCode(sw) {
	case LiquidityPoolWithdrawResultCodeLiquidityPoolWithdrawSuccess:
		return "", true
	default:
		return "", true
	}
}

// NewLiquidityPoolWithdrawResult creates a new  LiquidityPoolWithdrawResult.
func NewLiquidityPoolWithdrawResult(code LiquidityPoolWithdrawResultCode, value interface{}) (result LiquidityPoolWithdrawResult, err error) {
	result.Code = code
	switch LiquidityPoolWithdrawResultCode(code) {
	case LiquidityPoolWithdrawResultCodeLiquidityPoolWithdrawSuccess:
		// void
	default:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u LiquidityPoolWithdrawResult) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Code.EncodeTo(e); err != nil {
		return err
	}
	switch LiquidityPoolWithdrawResultCode(u.Code) {
	case LiquidityPoolWithdrawResultCodeLiquidityPoolWithdrawSuccess:
		// Void
		return nil
	default:
		// Void
		return nil
	}
}

var _ decoderFrom = (*LiquidityPoolWithdrawResult)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *LiquidityPoolWithdrawResult) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Code.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LiquidityPoolWithdrawResultCode: %s", err)
	}
	switch LiquidityPoolWithdrawResultCode(u.Code) {
	case LiquidityPoolWithdrawResultCodeLiquidityPoolWithdrawSuccess:
		// Void
		return n, nil
	default:
		// Void
		return n, nil
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LiquidityPoolWithdrawResult) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LiquidityPoolWithdrawResult) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LiquidityPoolWithdrawResult)(nil)
	_ encoding.BinaryUnmarshaler = (*LiquidityPoolWithdrawResult)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s LiquidityPoolWithdrawResult) xdrType() {}

var _ xdrType = (*LiquidityPoolWithdrawResult)(nil)

// OperationResultCode is an XDR Enum defines as:
//
//   enum OperationResultCode
//    {
//        opINNER = 0, // inner object result is valid
//
//        opBAD_AUTH = -1,            // too few valid signatures / wrong network
//        opNO_ACCOUNT = -2,          // source account was not found
//        opNOT_SUPPORTED = -3,       // operation not supported at this time
//        opTOO_MANY_SUBENTRIES = -4, // max number of subentries already reached
//        opEXCEEDED_WORK_LIMIT = -5, // operation did too much work
//        opTOO_MANY_SPONSORING = -6  // account is sponsoring too many entries
//    };
//
type OperationResultCode int32

const (
	OperationResultCodeOpInner             OperationResultCode = 0
	OperationResultCodeOpBadAuth           OperationResultCode = -1
	OperationResultCodeOpNoAccount         OperationResultCode = -2
	OperationResultCodeOpNotSupported      OperationResultCode = -3
	OperationResultCodeOpTooManySubentries OperationResultCode = -4
	OperationResultCodeOpExceededWorkLimit OperationResultCode = -5
	OperationResultCodeOpTooManySponsoring OperationResultCode = -6
)

var operationResultCodeMap = map[int32]string{
	0:  "OperationResultCodeOpInner",
	-1: "OperationResultCodeOpBadAuth",
	-2: "OperationResultCodeOpNoAccount",
	-3: "OperationResultCodeOpNotSupported",
	-4: "OperationResultCodeOpTooManySubentries",
	-5: "OperationResultCodeOpExceededWorkLimit",
	-6: "OperationResultCodeOpTooManySponsoring",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for OperationResultCode
func (e OperationResultCode) ValidEnum(v int32) bool {
	_, ok := operationResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e OperationResultCode) String() string {
	name, _ := operationResultCodeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e OperationResultCode) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := operationResultCodeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid OperationResultCode enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*OperationResultCode)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *OperationResultCode) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding OperationResultCode: %s", err)
	}
	if _, ok := operationResultCodeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid OperationResultCode enum value", v)
	}
	*e = OperationResultCode(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s OperationResultCode) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *OperationResultCode) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*OperationResultCode)(nil)
	_ encoding.BinaryUnmarshaler = (*OperationResultCode)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s OperationResultCode) xdrType() {}

var _ xdrType = (*OperationResultCode)(nil)

// OperationResultTr is an XDR NestedUnion defines as:
//
//   union switch (OperationType type)
//        {
//        case CREATE_ACCOUNT:
//            CreateAccountResult createAccountResult;
//        case PAYMENT:
//            PaymentResult paymentResult;
//        case PATH_PAYMENT_STRICT_RECEIVE:
//            PathPaymentStrictReceiveResult pathPaymentStrictReceiveResult;
//        case MANAGE_SELL_OFFER:
//            ManageSellOfferResult manageSellOfferResult;
//        case CREATE_PASSIVE_SELL_OFFER:
//            ManageSellOfferResult createPassiveSellOfferResult;
//        case SET_OPTIONS:
//            SetOptionsResult setOptionsResult;
//        case CHANGE_TRUST:
//            ChangeTrustResult changeTrustResult;
//        case ALLOW_TRUST:
//            AllowTrustResult allowTrustResult;
//        case ACCOUNT_MERGE:
//            AccountMergeResult accountMergeResult;
//        case INFLATION:
//            InflationResult inflationResult;
//        case MANAGE_DATA:
//            ManageDataResult manageDataResult;
//        case BUMP_SEQUENCE:
//            BumpSequenceResult bumpSeqResult;
//        case MANAGE_BUY_OFFER:
//            ManageBuyOfferResult manageBuyOfferResult;
//        case PATH_PAYMENT_STRICT_SEND:
//            PathPaymentStrictSendResult pathPaymentStrictSendResult;
//        case CREATE_CLAIMABLE_BALANCE:
//            CreateClaimableBalanceResult createClaimableBalanceResult;
//        case CLAIM_CLAIMABLE_BALANCE:
//            ClaimClaimableBalanceResult claimClaimableBalanceResult;
//        case BEGIN_SPONSORING_FUTURE_RESERVES:
//            BeginSponsoringFutureReservesResult beginSponsoringFutureReservesResult;
//        case END_SPONSORING_FUTURE_RESERVES:
//            EndSponsoringFutureReservesResult endSponsoringFutureReservesResult;
//        case REVOKE_SPONSORSHIP:
//            RevokeSponsorshipResult revokeSponsorshipResult;
//        case CLAWBACK:
//            ClawbackResult clawbackResult;
//        case CLAWBACK_CLAIMABLE_BALANCE:
//            ClawbackClaimableBalanceResult clawbackClaimableBalanceResult;
//        case SET_TRUST_LINE_FLAGS:
//            SetTrustLineFlagsResult setTrustLineFlagsResult;
//        case LIQUIDITY_POOL_DEPOSIT:
//            LiquidityPoolDepositResult liquidityPoolDepositResult;
//        case LIQUIDITY_POOL_WITHDRAW:
//            LiquidityPoolWithdrawResult liquidityPoolWithdrawResult;
//        }
//
type OperationResultTr struct {
	Type                                OperationType
	CreateAccountResult                 *CreateAccountResult
	PaymentResult                       *PaymentResult
	PathPaymentStrictReceiveResult      *PathPaymentStrictReceiveResult
	ManageSellOfferResult               *ManageSellOfferResult
	CreatePassiveSellOfferResult        *ManageSellOfferResult
	SetOptionsResult                    *SetOptionsResult
	ChangeTrustResult                   *ChangeTrustResult
	AllowTrustResult                    *AllowTrustResult
	AccountMergeResult                  *AccountMergeResult
	InflationResult                     *InflationResult
	ManageDataResult                    *ManageDataResult
	BumpSeqResult                       *BumpSequenceResult
	ManageBuyOfferResult                *ManageBuyOfferResult
	PathPaymentStrictSendResult         *PathPaymentStrictSendResult
	CreateClaimableBalanceResult        *CreateClaimableBalanceResult
	ClaimClaimableBalanceResult         *ClaimClaimableBalanceResult
	BeginSponsoringFutureReservesResult *BeginSponsoringFutureReservesResult
	EndSponsoringFutureReservesResult   *EndSponsoringFutureReservesResult
	RevokeSponsorshipResult             *RevokeSponsorshipResult
	ClawbackResult                      *ClawbackResult
	ClawbackClaimableBalanceResult      *ClawbackClaimableBalanceResult
	SetTrustLineFlagsResult             *SetTrustLineFlagsResult
	LiquidityPoolDepositResult          *LiquidityPoolDepositResult
	LiquidityPoolWithdrawResult         *LiquidityPoolWithdrawResult
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u OperationResultTr) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of OperationResultTr
func (u OperationResultTr) ArmForSwitch(sw int32) (string, bool) {
	switch OperationType(sw) {
	case OperationTypeCreateAccount:
		return "CreateAccountResult", true
	case OperationTypePayment:
		return "PaymentResult", true
	case OperationTypePathPaymentStrictReceive:
		return "PathPaymentStrictReceiveResult", true
	case OperationTypeManageSellOffer:
		return "ManageSellOfferResult", true
	case OperationTypeCreatePassiveSellOffer:
		return "CreatePassiveSellOfferResult", true
	case OperationTypeSetOptions:
		return "SetOptionsResult", true
	case OperationTypeChangeTrust:
		return "ChangeTrustResult", true
	case OperationTypeAllowTrust:
		return "AllowTrustResult", true
	case OperationTypeAccountMerge:
		return "AccountMergeResult", true
	case OperationTypeInflation:
		return "InflationResult", true
	case OperationTypeManageData:
		return "ManageDataResult", true
	case OperationTypeBumpSequence:
		return "BumpSeqResult", true
	case OperationTypeManageBuyOffer:
		return "ManageBuyOfferResult", true
	case OperationTypePathPaymentStrictSend:
		return "PathPaymentStrictSendResult", true
	case OperationTypeCreateClaimableBalance:
		return "CreateClaimableBalanceResult", true
	case OperationTypeClaimClaimableBalance:
		return "ClaimClaimableBalanceResult", true
	case OperationTypeBeginSponsoringFutureReserves:
		return "BeginSponsoringFutureReservesResult", true
	case OperationTypeEndSponsoringFutureReserves:
		return "EndSponsoringFutureReservesResult", true
	case OperationTypeRevokeSponsorship:
		return "RevokeSponsorshipResult", true
	case OperationTypeClawback:
		return "ClawbackResult", true
	case OperationTypeClawbackClaimableBalance:
		return "ClawbackClaimableBalanceResult", true
	case OperationTypeSetTrustLineFlags:
		return "SetTrustLineFlagsResult", true
	case OperationTypeLiquidityPoolDeposit:
		return "LiquidityPoolDepositResult", true
	case OperationTypeLiquidityPoolWithdraw:
		return "LiquidityPoolWithdrawResult", true
	}
	return "-", false
}

// NewOperationResultTr creates a new  OperationResultTr.
func NewOperationResultTr(aType OperationType, value interface{}) (result OperationResultTr, err error) {
	result.Type = aType
	switch OperationType(aType) {
	case OperationTypeCreateAccount:
		tv, ok := value.(CreateAccountResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be CreateAccountResult")
			return
		}
		result.CreateAccountResult = &tv
	case OperationTypePayment:
		tv, ok := value.(PaymentResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be PaymentResult")
			return
		}
		result.PaymentResult = &tv
	case OperationTypePathPaymentStrictReceive:
		tv, ok := value.(PathPaymentStrictReceiveResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be PathPaymentStrictReceiveResult")
			return
		}
		result.PathPaymentStrictReceiveResult = &tv
	case OperationTypeManageSellOffer:
		tv, ok := value.(ManageSellOfferResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageSellOfferResult")
			return
		}
		result.ManageSellOfferResult = &tv
	case OperationTypeCreatePassiveSellOffer:
		tv, ok := value.(ManageSellOfferResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageSellOfferResult")
			return
		}
		result.CreatePassiveSellOfferResult = &tv
	case OperationTypeSetOptions:
		tv, ok := value.(SetOptionsResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be SetOptionsResult")
			return
		}
		result.SetOptionsResult = &tv
	case OperationTypeChangeTrust:
		tv, ok := value.(ChangeTrustResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ChangeTrustResult")
			return
		}
		result.ChangeTrustResult = &tv
	case OperationTypeAllowTrust:
		tv, ok := value.(AllowTrustResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be AllowTrustResult")
			return
		}
		result.AllowTrustResult = &tv
	case OperationTypeAccountMerge:
		tv, ok := value.(AccountMergeResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be AccountMergeResult")
			return
		}
		result.AccountMergeResult = &tv
	case OperationTypeInflation:
		tv, ok := value.(InflationResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be InflationResult")
			return
		}
		result.InflationResult = &tv
	case OperationTypeManageData:
		tv, ok := value.(ManageDataResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageDataResult")
			return
		}
		result.ManageDataResult = &tv
	case OperationTypeBumpSequence:
		tv, ok := value.(BumpSequenceResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be BumpSequenceResult")
			return
		}
		result.BumpSeqResult = &tv
	case OperationTypeManageBuyOffer:
		tv, ok := value.(ManageBuyOfferResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageBuyOfferResult")
			return
		}
		result.ManageBuyOfferResult = &tv
	case OperationTypePathPaymentStrictSend:
		tv, ok := value.(PathPaymentStrictSendResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be PathPaymentStrictSendResult")
			return
		}
		result.PathPaymentStrictSendResult = &tv
	case OperationTypeCreateClaimableBalance:
		tv, ok := value.(CreateClaimableBalanceResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be CreateClaimableBalanceResult")
			return
		}
		result.CreateClaimableBalanceResult = &tv
	case OperationTypeClaimClaimableBalance:
		tv, ok := value.(ClaimClaimableBalanceResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ClaimClaimableBalanceResult")
			return
		}
		result.ClaimClaimableBalanceResult = &tv
	case OperationTypeBeginSponsoringFutureReserves:
		tv, ok := value.(BeginSponsoringFutureReservesResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be BeginSponsoringFutureReservesResult")
			return
		}
		result.BeginSponsoringFutureReservesResult = &tv
	case OperationTypeEndSponsoringFutureReserves:
		tv, ok := value.(EndSponsoringFutureReservesResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be EndSponsoringFutureReservesResult")
			return
		}
		result.EndSponsoringFutureReservesResult = &tv
	case OperationTypeRevokeSponsorship:
		tv, ok := value.(RevokeSponsorshipResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be RevokeSponsorshipResult")
			return
		}
		result.RevokeSponsorshipResult = &tv
	case OperationTypeClawback:
		tv, ok := value.(ClawbackResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ClawbackResult")
			return
		}
		result.ClawbackResult = &tv
	case OperationTypeClawbackClaimableBalance:
		tv, ok := value.(ClawbackClaimableBalanceResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ClawbackClaimableBalanceResult")
			return
		}
		result.ClawbackClaimableBalanceResult = &tv
	case OperationTypeSetTrustLineFlags:
		tv, ok := value.(SetTrustLineFlagsResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be SetTrustLineFlagsResult")
			return
		}
		result.SetTrustLineFlagsResult = &tv
	case OperationTypeLiquidityPoolDeposit:
		tv, ok := value.(LiquidityPoolDepositResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be LiquidityPoolDepositResult")
			return
		}
		result.LiquidityPoolDepositResult = &tv
	case OperationTypeLiquidityPoolWithdraw:
		tv, ok := value.(LiquidityPoolWithdrawResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be LiquidityPoolWithdrawResult")
			return
		}
		result.LiquidityPoolWithdrawResult = &tv
	}
	return
}

// MustCreateAccountResult retrieves the CreateAccountResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustCreateAccountResult() CreateAccountResult {
	val, ok := u.GetCreateAccountResult()

	if !ok {
		panic("arm CreateAccountResult is not set")
	}

	return val
}

// GetCreateAccountResult retrieves the CreateAccountResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetCreateAccountResult() (result CreateAccountResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "CreateAccountResult" {
		result = *u.CreateAccountResult
		ok = true
	}

	return
}

// MustPaymentResult retrieves the PaymentResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustPaymentResult() PaymentResult {
	val, ok := u.GetPaymentResult()

	if !ok {
		panic("arm PaymentResult is not set")
	}

	return val
}

// GetPaymentResult retrieves the PaymentResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetPaymentResult() (result PaymentResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "PaymentResult" {
		result = *u.PaymentResult
		ok = true
	}

	return
}

// MustPathPaymentStrictReceiveResult retrieves the PathPaymentStrictReceiveResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustPathPaymentStrictReceiveResult() PathPaymentStrictReceiveResult {
	val, ok := u.GetPathPaymentStrictReceiveResult()

	if !ok {
		panic("arm PathPaymentStrictReceiveResult is not set")
	}

	return val
}

// GetPathPaymentStrictReceiveResult retrieves the PathPaymentStrictReceiveResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetPathPaymentStrictReceiveResult() (result PathPaymentStrictReceiveResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "PathPaymentStrictReceiveResult" {
		result = *u.PathPaymentStrictReceiveResult
		ok = true
	}

	return
}

// MustManageSellOfferResult retrieves the ManageSellOfferResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustManageSellOfferResult() ManageSellOfferResult {
	val, ok := u.GetManageSellOfferResult()

	if !ok {
		panic("arm ManageSellOfferResult is not set")
	}

	return val
}

// GetManageSellOfferResult retrieves the ManageSellOfferResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetManageSellOfferResult() (result ManageSellOfferResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageSellOfferResult" {
		result = *u.ManageSellOfferResult
		ok = true
	}

	return
}

// MustCreatePassiveSellOfferResult retrieves the CreatePassiveSellOfferResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustCreatePassiveSellOfferResult() ManageSellOfferResult {
	val, ok := u.GetCreatePassiveSellOfferResult()

	if !ok {
		panic("arm CreatePassiveSellOfferResult is not set")
	}

	return val
}

// GetCreatePassiveSellOfferResult retrieves the CreatePassiveSellOfferResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetCreatePassiveSellOfferResult() (result ManageSellOfferResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "CreatePassiveSellOfferResult" {
		result = *u.CreatePassiveSellOfferResult
		ok = true
	}

	return
}

// MustSetOptionsResult retrieves the SetOptionsResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustSetOptionsResult() SetOptionsResult {
	val, ok := u.GetSetOptionsResult()

	if !ok {
		panic("arm SetOptionsResult is not set")
	}

	return val
}

// GetSetOptionsResult retrieves the SetOptionsResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetSetOptionsResult() (result SetOptionsResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "SetOptionsResult" {
		result = *u.SetOptionsResult
		ok = true
	}

	return
}

// MustChangeTrustResult retrieves the ChangeTrustResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustChangeTrustResult() ChangeTrustResult {
	val, ok := u.GetChangeTrustResult()

	if !ok {
		panic("arm ChangeTrustResult is not set")
	}

	return val
}

// GetChangeTrustResult retrieves the ChangeTrustResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetChangeTrustResult() (result ChangeTrustResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ChangeTrustResult" {
		result = *u.ChangeTrustResult
		ok = true
	}

	return
}

// MustAllowTrustResult retrieves the AllowTrustResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustAllowTrustResult() AllowTrustResult {
	val, ok := u.GetAllowTrustResult()

	if !ok {
		panic("arm AllowTrustResult is not set")
	}

	return val
}

// GetAllowTrustResult retrieves the AllowTrustResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetAllowTrustResult() (result AllowTrustResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AllowTrustResult" {
		result = *u.AllowTrustResult
		ok = true
	}

	return
}

// MustAccountMergeResult retrieves the AccountMergeResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustAccountMergeResult() AccountMergeResult {
	val, ok := u.GetAccountMergeResult()

	if !ok {
		panic("arm AccountMergeResult is not set")
	}

	return val
}

// GetAccountMergeResult retrieves the AccountMergeResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetAccountMergeResult() (result AccountMergeResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AccountMergeResult" {
		result = *u.AccountMergeResult
		ok = true
	}

	return
}

// MustInflationResult retrieves the InflationResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustInflationResult() InflationResult {
	val, ok := u.GetInflationResult()

	if !ok {
		panic("arm InflationResult is not set")
	}

	return val
}

// GetInflationResult retrieves the InflationResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetInflationResult() (result InflationResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "InflationResult" {
		result = *u.InflationResult
		ok = true
	}

	return
}

// MustManageDataResult retrieves the ManageDataResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustManageDataResult() ManageDataResult {
	val, ok := u.GetManageDataResult()

	if !ok {
		panic("arm ManageDataResult is not set")
	}

	return val
}

// GetManageDataResult retrieves the ManageDataResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetManageDataResult() (result ManageDataResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageDataResult" {
		result = *u.ManageDataResult
		ok = true
	}

	return
}

// MustBumpSeqResult retrieves the BumpSeqResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustBumpSeqResult() BumpSequenceResult {
	val, ok := u.GetBumpSeqResult()

	if !ok {
		panic("arm BumpSeqResult is not set")
	}

	return val
}

// GetBumpSeqResult retrieves the BumpSeqResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetBumpSeqResult() (result BumpSequenceResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "BumpSeqResult" {
		result = *u.BumpSeqResult
		ok = true
	}

	return
}

// MustManageBuyOfferResult retrieves the ManageBuyOfferResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustManageBuyOfferResult() ManageBuyOfferResult {
	val, ok := u.GetManageBuyOfferResult()

	if !ok {
		panic("arm ManageBuyOfferResult is not set")
	}

	return val
}

// GetManageBuyOfferResult retrieves the ManageBuyOfferResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetManageBuyOfferResult() (result ManageBuyOfferResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageBuyOfferResult" {
		result = *u.ManageBuyOfferResult
		ok = true
	}

	return
}

// MustPathPaymentStrictSendResult retrieves the PathPaymentStrictSendResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustPathPaymentStrictSendResult() PathPaymentStrictSendResult {
	val, ok := u.GetPathPaymentStrictSendResult()

	if !ok {
		panic("arm PathPaymentStrictSendResult is not set")
	}

	return val
}

// GetPathPaymentStrictSendResult retrieves the PathPaymentStrictSendResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetPathPaymentStrictSendResult() (result PathPaymentStrictSendResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "PathPaymentStrictSendResult" {
		result = *u.PathPaymentStrictSendResult
		ok = true
	}

	return
}

// MustCreateClaimableBalanceResult retrieves the CreateClaimableBalanceResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustCreateClaimableBalanceResult() CreateClaimableBalanceResult {
	val, ok := u.GetCreateClaimableBalanceResult()

	if !ok {
		panic("arm CreateClaimableBalanceResult is not set")
	}

	return val
}

// GetCreateClaimableBalanceResult retrieves the CreateClaimableBalanceResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetCreateClaimableBalanceResult() (result CreateClaimableBalanceResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "CreateClaimableBalanceResult" {
		result = *u.CreateClaimableBalanceResult
		ok = true
	}

	return
}

// MustClaimClaimableBalanceResult retrieves the ClaimClaimableBalanceResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustClaimClaimableBalanceResult() ClaimClaimableBalanceResult {
	val, ok := u.GetClaimClaimableBalanceResult()

	if !ok {
		panic("arm ClaimClaimableBalanceResult is not set")
	}

	return val
}

// GetClaimClaimableBalanceResult retrieves the ClaimClaimableBalanceResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetClaimClaimableBalanceResult() (result ClaimClaimableBalanceResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ClaimClaimableBalanceResult" {
		result = *u.ClaimClaimableBalanceResult
		ok = true
	}

	return
}

// MustBeginSponsoringFutureReservesResult retrieves the BeginSponsoringFutureReservesResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustBeginSponsoringFutureReservesResult() BeginSponsoringFutureReservesResult {
	val, ok := u.GetBeginSponsoringFutureReservesResult()

	if !ok {
		panic("arm BeginSponsoringFutureReservesResult is not set")
	}

	return val
}

// GetBeginSponsoringFutureReservesResult retrieves the BeginSponsoringFutureReservesResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetBeginSponsoringFutureReservesResult() (result BeginSponsoringFutureReservesResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "BeginSponsoringFutureReservesResult" {
		result = *u.BeginSponsoringFutureReservesResult
		ok = true
	}

	return
}

// MustEndSponsoringFutureReservesResult retrieves the EndSponsoringFutureReservesResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustEndSponsoringFutureReservesResult() EndSponsoringFutureReservesResult {
	val, ok := u.GetEndSponsoringFutureReservesResult()

	if !ok {
		panic("arm EndSponsoringFutureReservesResult is not set")
	}

	return val
}

// GetEndSponsoringFutureReservesResult retrieves the EndSponsoringFutureReservesResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetEndSponsoringFutureReservesResult() (result EndSponsoringFutureReservesResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "EndSponsoringFutureReservesResult" {
		result = *u.EndSponsoringFutureReservesResult
		ok = true
	}

	return
}

// MustRevokeSponsorshipResult retrieves the RevokeSponsorshipResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustRevokeSponsorshipResult() RevokeSponsorshipResult {
	val, ok := u.GetRevokeSponsorshipResult()

	if !ok {
		panic("arm RevokeSponsorshipResult is not set")
	}

	return val
}

// GetRevokeSponsorshipResult retrieves the RevokeSponsorshipResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetRevokeSponsorshipResult() (result RevokeSponsorshipResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "RevokeSponsorshipResult" {
		result = *u.RevokeSponsorshipResult
		ok = true
	}

	return
}

// MustClawbackResult retrieves the ClawbackResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustClawbackResult() ClawbackResult {
	val, ok := u.GetClawbackResult()

	if !ok {
		panic("arm ClawbackResult is not set")
	}

	return val
}

// GetClawbackResult retrieves the ClawbackResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetClawbackResult() (result ClawbackResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ClawbackResult" {
		result = *u.ClawbackResult
		ok = true
	}

	return
}

// MustClawbackClaimableBalanceResult retrieves the ClawbackClaimableBalanceResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustClawbackClaimableBalanceResult() ClawbackClaimableBalanceResult {
	val, ok := u.GetClawbackClaimableBalanceResult()

	if !ok {
		panic("arm ClawbackClaimableBalanceResult is not set")
	}

	return val
}

// GetClawbackClaimableBalanceResult retrieves the ClawbackClaimableBalanceResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetClawbackClaimableBalanceResult() (result ClawbackClaimableBalanceResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ClawbackClaimableBalanceResult" {
		result = *u.ClawbackClaimableBalanceResult
		ok = true
	}

	return
}

// MustSetTrustLineFlagsResult retrieves the SetTrustLineFlagsResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustSetTrustLineFlagsResult() SetTrustLineFlagsResult {
	val, ok := u.GetSetTrustLineFlagsResult()

	if !ok {
		panic("arm SetTrustLineFlagsResult is not set")
	}

	return val
}

// GetSetTrustLineFlagsResult retrieves the SetTrustLineFlagsResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetSetTrustLineFlagsResult() (result SetTrustLineFlagsResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "SetTrustLineFlagsResult" {
		result = *u.SetTrustLineFlagsResult
		ok = true
	}

	return
}

// MustLiquidityPoolDepositResult retrieves the LiquidityPoolDepositResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustLiquidityPoolDepositResult() LiquidityPoolDepositResult {
	val, ok := u.GetLiquidityPoolDepositResult()

	if !ok {
		panic("arm LiquidityPoolDepositResult is not set")
	}

	return val
}

// GetLiquidityPoolDepositResult retrieves the LiquidityPoolDepositResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetLiquidityPoolDepositResult() (result LiquidityPoolDepositResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "LiquidityPoolDepositResult" {
		result = *u.LiquidityPoolDepositResult
		ok = true
	}

	return
}

// MustLiquidityPoolWithdrawResult retrieves the LiquidityPoolWithdrawResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustLiquidityPoolWithdrawResult() LiquidityPoolWithdrawResult {
	val, ok := u.GetLiquidityPoolWithdrawResult()

	if !ok {
		panic("arm LiquidityPoolWithdrawResult is not set")
	}

	return val
}

// GetLiquidityPoolWithdrawResult retrieves the LiquidityPoolWithdrawResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetLiquidityPoolWithdrawResult() (result LiquidityPoolWithdrawResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "LiquidityPoolWithdrawResult" {
		result = *u.LiquidityPoolWithdrawResult
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u OperationResultTr) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch OperationType(u.Type) {
	case OperationTypeCreateAccount:
		if err = (*u.CreateAccountResult).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypePayment:
		if err = (*u.PaymentResult).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypePathPaymentStrictReceive:
		if err = (*u.PathPaymentStrictReceiveResult).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeManageSellOffer:
		if err = (*u.ManageSellOfferResult).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeCreatePassiveSellOffer:
		if err = (*u.CreatePassiveSellOfferResult).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeSetOptions:
		if err = (*u.SetOptionsResult).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeChangeTrust:
		if err = (*u.ChangeTrustResult).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeAllowTrust:
		if err = (*u.AllowTrustResult).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeAccountMerge:
		if err = (*u.AccountMergeResult).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeInflation:
		if err = (*u.InflationResult).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeManageData:
		if err = (*u.ManageDataResult).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeBumpSequence:
		if err = (*u.BumpSeqResult).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeManageBuyOffer:
		if err = (*u.ManageBuyOfferResult).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypePathPaymentStrictSend:
		if err = (*u.PathPaymentStrictSendResult).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeCreateClaimableBalance:
		if err = (*u.CreateClaimableBalanceResult).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeClaimClaimableBalance:
		if err = (*u.ClaimClaimableBalanceResult).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeBeginSponsoringFutureReserves:
		if err = (*u.BeginSponsoringFutureReservesResult).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeEndSponsoringFutureReserves:
		if err = (*u.EndSponsoringFutureReservesResult).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeRevokeSponsorship:
		if err = (*u.RevokeSponsorshipResult).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeClawback:
		if err = (*u.ClawbackResult).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeClawbackClaimableBalance:
		if err = (*u.ClawbackClaimableBalanceResult).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeSetTrustLineFlags:
		if err = (*u.SetTrustLineFlagsResult).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeLiquidityPoolDeposit:
		if err = (*u.LiquidityPoolDepositResult).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeLiquidityPoolWithdraw:
		if err = (*u.LiquidityPoolWithdrawResult).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (OperationType) switch value '%d' is not valid for union OperationResultTr", u.Type)
}

var _ decoderFrom = (*OperationResultTr)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *OperationResultTr) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding OperationType: %s", err)
	}
	switch OperationType(u.Type) {
	case OperationTypeCreateAccount:
		u.CreateAccountResult = new(CreateAccountResult)
		nTmp, err = (*u.CreateAccountResult).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding CreateAccountResult: %s", err)
		}
		return n, nil
	case OperationTypePayment:
		u.PaymentResult = new(PaymentResult)
		nTmp, err = (*u.PaymentResult).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding PaymentResult: %s", err)
		}
		return n, nil
	case OperationTypePathPaymentStrictReceive:
		u.PathPaymentStrictReceiveResult = new(PathPaymentStrictReceiveResult)
		nTmp, err = (*u.PathPaymentStrictReceiveResult).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding PathPaymentStrictReceiveResult: %s", err)
		}
		return n, nil
	case OperationTypeManageSellOffer:
		u.ManageSellOfferResult = new(ManageSellOfferResult)
		nTmp, err = (*u.ManageSellOfferResult).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ManageSellOfferResult: %s", err)
		}
		return n, nil
	case OperationTypeCreatePassiveSellOffer:
		u.CreatePassiveSellOfferResult = new(ManageSellOfferResult)
		nTmp, err = (*u.CreatePassiveSellOfferResult).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ManageSellOfferResult: %s", err)
		}
		return n, nil
	case OperationTypeSetOptions:
		u.SetOptionsResult = new(SetOptionsResult)
		nTmp, err = (*u.SetOptionsResult).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding SetOptionsResult: %s", err)
		}
		return n, nil
	case OperationTypeChangeTrust:
		u.ChangeTrustResult = new(ChangeTrustResult)
		nTmp, err = (*u.ChangeTrustResult).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ChangeTrustResult: %s", err)
		}
		return n, nil
	case OperationTypeAllowTrust:
		u.AllowTrustResult = new(AllowTrustResult)
		nTmp, err = (*u.AllowTrustResult).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding AllowTrustResult: %s", err)
		}
		return n, nil
	case OperationTypeAccountMerge:
		u.AccountMergeResult = new(AccountMergeResult)
		nTmp, err = (*u.AccountMergeResult).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding AccountMergeResult: %s", err)
		}
		return n, nil
	case OperationTypeInflation:
		u.InflationResult = new(InflationResult)
		nTmp, err = (*u.InflationResult).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding InflationResult: %s", err)
		}
		return n, nil
	case OperationTypeManageData:
		u.ManageDataResult = new(ManageDataResult)
		nTmp, err = (*u.ManageDataResult).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ManageDataResult: %s", err)
		}
		return n, nil
	case OperationTypeBumpSequence:
		u.BumpSeqResult = new(BumpSequenceResult)
		nTmp, err = (*u.BumpSeqResult).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding BumpSequenceResult: %s", err)
		}
		return n, nil
	case OperationTypeManageBuyOffer:
		u.ManageBuyOfferResult = new(ManageBuyOfferResult)
		nTmp, err = (*u.ManageBuyOfferResult).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ManageBuyOfferResult: %s", err)
		}
		return n, nil
	case OperationTypePathPaymentStrictSend:
		u.PathPaymentStrictSendResult = new(PathPaymentStrictSendResult)
		nTmp, err = (*u.PathPaymentStrictSendResult).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding PathPaymentStrictSendResult: %s", err)
		}
		return n, nil
	case OperationTypeCreateClaimableBalance:
		u.CreateClaimableBalanceResult = new(CreateClaimableBalanceResult)
		nTmp, err = (*u.CreateClaimableBalanceResult).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding CreateClaimableBalanceResult: %s", err)
		}
		return n, nil
	case OperationTypeClaimClaimableBalance:
		u.ClaimClaimableBalanceResult = new(ClaimClaimableBalanceResult)
		nTmp, err = (*u.ClaimClaimableBalanceResult).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ClaimClaimableBalanceResult: %s", err)
		}
		return n, nil
	case OperationTypeBeginSponsoringFutureReserves:
		u.BeginSponsoringFutureReservesResult = new(BeginSponsoringFutureReservesResult)
		nTmp, err = (*u.BeginSponsoringFutureReservesResult).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding BeginSponsoringFutureReservesResult: %s", err)
		}
		return n, nil
	case OperationTypeEndSponsoringFutureReserves:
		u.EndSponsoringFutureReservesResult = new(EndSponsoringFutureReservesResult)
		nTmp, err = (*u.EndSponsoringFutureReservesResult).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding EndSponsoringFutureReservesResult: %s", err)
		}
		return n, nil
	case OperationTypeRevokeSponsorship:
		u.RevokeSponsorshipResult = new(RevokeSponsorshipResult)
		nTmp, err = (*u.RevokeSponsorshipResult).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding RevokeSponsorshipResult: %s", err)
		}
		return n, nil
	case OperationTypeClawback:
		u.ClawbackResult = new(ClawbackResult)
		nTmp, err = (*u.ClawbackResult).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ClawbackResult: %s", err)
		}
		return n, nil
	case OperationTypeClawbackClaimableBalance:
		u.ClawbackClaimableBalanceResult = new(ClawbackClaimableBalanceResult)
		nTmp, err = (*u.ClawbackClaimableBalanceResult).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ClawbackClaimableBalanceResult: %s", err)
		}
		return n, nil
	case OperationTypeSetTrustLineFlags:
		u.SetTrustLineFlagsResult = new(SetTrustLineFlagsResult)
		nTmp, err = (*u.SetTrustLineFlagsResult).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding SetTrustLineFlagsResult: %s", err)
		}
		return n, nil
	case OperationTypeLiquidityPoolDeposit:
		u.LiquidityPoolDepositResult = new(LiquidityPoolDepositResult)
		nTmp, err = (*u.LiquidityPoolDepositResult).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding LiquidityPoolDepositResult: %s", err)
		}
		return n, nil
	case OperationTypeLiquidityPoolWithdraw:
		u.LiquidityPoolWithdrawResult = new(LiquidityPoolWithdrawResult)
		nTmp, err = (*u.LiquidityPoolWithdrawResult).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding LiquidityPoolWithdrawResult: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union OperationResultTr has invalid Type (OperationType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s OperationResultTr) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *OperationResultTr) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*OperationResultTr)(nil)
	_ encoding.BinaryUnmarshaler = (*OperationResultTr)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s OperationResultTr) xdrType() {}

var _ xdrType = (*OperationResultTr)(nil)

// OperationResult is an XDR Union defines as:
//
//   union OperationResult switch (OperationResultCode code)
//    {
//    case opINNER:
//        union switch (OperationType type)
//        {
//        case CREATE_ACCOUNT:
//            CreateAccountResult createAccountResult;
//        case PAYMENT:
//            PaymentResult paymentResult;
//        case PATH_PAYMENT_STRICT_RECEIVE:
//            PathPaymentStrictReceiveResult pathPaymentStrictReceiveResult;
//        case MANAGE_SELL_OFFER:
//            ManageSellOfferResult manageSellOfferResult;
//        case CREATE_PASSIVE_SELL_OFFER:
//            ManageSellOfferResult createPassiveSellOfferResult;
//        case SET_OPTIONS:
//            SetOptionsResult setOptionsResult;
//        case CHANGE_TRUST:
//            ChangeTrustResult changeTrustResult;
//        case ALLOW_TRUST:
//            AllowTrustResult allowTrustResult;
//        case ACCOUNT_MERGE:
//            AccountMergeResult accountMergeResult;
//        case INFLATION:
//            InflationResult inflationResult;
//        case MANAGE_DATA:
//            ManageDataResult manageDataResult;
//        case BUMP_SEQUENCE:
//            BumpSequenceResult bumpSeqResult;
//        case MANAGE_BUY_OFFER:
//            ManageBuyOfferResult manageBuyOfferResult;
//        case PATH_PAYMENT_STRICT_SEND:
//            PathPaymentStrictSendResult pathPaymentStrictSendResult;
//        case CREATE_CLAIMABLE_BALANCE:
//            CreateClaimableBalanceResult createClaimableBalanceResult;
//        case CLAIM_CLAIMABLE_BALANCE:
//            ClaimClaimableBalanceResult claimClaimableBalanceResult;
//        case BEGIN_SPONSORING_FUTURE_RESERVES:
//            BeginSponsoringFutureReservesResult beginSponsoringFutureReservesResult;
//        case END_SPONSORING_FUTURE_RESERVES:
//            EndSponsoringFutureReservesResult endSponsoringFutureReservesResult;
//        case REVOKE_SPONSORSHIP:
//            RevokeSponsorshipResult revokeSponsorshipResult;
//        case CLAWBACK:
//            ClawbackResult clawbackResult;
//        case CLAWBACK_CLAIMABLE_BALANCE:
//            ClawbackClaimableBalanceResult clawbackClaimableBalanceResult;
//        case SET_TRUST_LINE_FLAGS:
//            SetTrustLineFlagsResult setTrustLineFlagsResult;
//        case LIQUIDITY_POOL_DEPOSIT:
//            LiquidityPoolDepositResult liquidityPoolDepositResult;
//        case LIQUIDITY_POOL_WITHDRAW:
//            LiquidityPoolWithdrawResult liquidityPoolWithdrawResult;
//        }
//        tr;
//    default:
//        void;
//    };
//
type OperationResult struct {
	Code OperationResultCode
	Tr   *OperationResultTr
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u OperationResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of OperationResult
func (u OperationResult) ArmForSwitch(sw int32) (string, bool) {
	switch OperationResultCode(sw) {
	case OperationResultCodeOpInner:
		return "Tr", true
	default:
		return "", true
	}
}

// NewOperationResult creates a new  OperationResult.
func NewOperationResult(code OperationResultCode, value interface{}) (result OperationResult, err error) {
	result.Code = code
	switch OperationResultCode(code) {
	case OperationResultCodeOpInner:
		tv, ok := value.(OperationResultTr)
		if !ok {
			err = fmt.Errorf("invalid value, must be OperationResultTr")
			return
		}
		result.Tr = &tv
	default:
		// void
	}
	return
}

// MustTr retrieves the Tr value from the union,
// panicing if the value is not set.
func (u OperationResult) MustTr() OperationResultTr {
	val, ok := u.GetTr()

	if !ok {
		panic("arm Tr is not set")
	}

	return val
}

// GetTr retrieves the Tr value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResult) GetTr() (result OperationResultTr, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Tr" {
		result = *u.Tr
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u OperationResult) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Code.EncodeTo(e); err != nil {
		return err
	}
	switch OperationResultCode(u.Code) {
	case OperationResultCodeOpInner:
		if err = (*u.Tr).EncodeTo(e); err != nil {
			return err
		}
		return nil
	default:
		// Void
		return nil
	}
}

var _ decoderFrom = (*OperationResult)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *OperationResult) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Code.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding OperationResultCode: %s", err)
	}
	switch OperationResultCode(u.Code) {
	case OperationResultCodeOpInner:
		u.Tr = new(OperationResultTr)
		nTmp, err = (*u.Tr).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding OperationResultTr: %s", err)
		}
		return n, nil
	default:
		// Void
		return n, nil
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s OperationResult) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *OperationResult) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*OperationResult)(nil)
	_ encoding.BinaryUnmarshaler = (*OperationResult)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s OperationResult) xdrType() {}

var _ xdrType = (*OperationResult)(nil)

// TransactionResultCode is an XDR Enum defines as:
//
//   enum TransactionResultCode
//    {
//        txFEE_BUMP_INNER_SUCCESS = 1, // fee bump inner transaction succeeded
//        txSUCCESS = 0,                // all operations succeeded
//
//        txFAILED = -1, // one of the operations failed (none were applied)
//
//        txTOO_EARLY = -2,         // ledger closeTime before minTime
//        txTOO_LATE = -3,          // ledger closeTime after maxTime
//        txMISSING_OPERATION = -4, // no operation was specified
//        txBAD_SEQ = -5,           // sequence number does not match source account
//
//        txBAD_AUTH = -6,             // too few valid signatures / wrong network
//        txINSUFFICIENT_BALANCE = -7, // fee would bring account below reserve
//        txNO_ACCOUNT = -8,           // source account not found
//        txINSUFFICIENT_FEE = -9,     // fee is too small
//        txBAD_AUTH_EXTRA = -10,      // unused signatures attached to transaction
//        txINTERNAL_ERROR = -11,      // an unknown error occurred
//
//        txNOT_SUPPORTED = -12,         // transaction type not supported
//        txFEE_BUMP_INNER_FAILED = -13, // fee bump inner transaction failed
//        txBAD_SPONSORSHIP = -14        // sponsorship not confirmed
//    };
//
type TransactionResultCode int32

const (
	TransactionResultCodeTxFeeBumpInnerSuccess TransactionResultCode = 1
	TransactionResultCodeTxSuccess             TransactionResultCode = 0
	TransactionResultCodeTxFailed              TransactionResultCode = -1
	TransactionResultCodeTxTooEarly            TransactionResultCode = -2
	TransactionResultCodeTxTooLate             TransactionResultCode = -3
	TransactionResultCodeTxMissingOperation    TransactionResultCode = -4
	TransactionResultCodeTxBadSeq              TransactionResultCode = -5
	TransactionResultCodeTxBadAuth             TransactionResultCode = -6
	TransactionResultCodeTxInsufficientBalance TransactionResultCode = -7
	TransactionResultCodeTxNoAccount           TransactionResultCode = -8
	TransactionResultCodeTxInsufficientFee     TransactionResultCode = -9
	TransactionResultCodeTxBadAuthExtra        TransactionResultCode = -10
	TransactionResultCodeTxInternalError       TransactionResultCode = -11
	TransactionResultCodeTxNotSupported        TransactionResultCode = -12
	TransactionResultCodeTxFeeBumpInnerFailed  TransactionResultCode = -13
	TransactionResultCodeTxBadSponsorship      TransactionResultCode = -14
)

var transactionResultCodeMap = map[int32]string{
	1:   "TransactionResultCodeTxFeeBumpInnerSuccess",
	0:   "TransactionResultCodeTxSuccess",
	-1:  "TransactionResultCodeTxFailed",
	-2:  "TransactionResultCodeTxTooEarly",
	-3:  "TransactionResultCodeTxTooLate",
	-4:  "TransactionResultCodeTxMissingOperation",
	-5:  "TransactionResultCodeTxBadSeq",
	-6:  "TransactionResultCodeTxBadAuth",
	-7:  "TransactionResultCodeTxInsufficientBalance",
	-8:  "TransactionResultCodeTxNoAccount",
	-9:  "TransactionResultCodeTxInsufficientFee",
	-10: "TransactionResultCodeTxBadAuthExtra",
	-11: "TransactionResultCodeTxInternalError",
	-12: "TransactionResultCodeTxNotSupported",
	-13: "TransactionResultCodeTxFeeBumpInnerFailed",
	-14: "TransactionResultCodeTxBadSponsorship",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for TransactionResultCode
func (e TransactionResultCode) ValidEnum(v int32) bool {
	_, ok := transactionResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e TransactionResultCode) String() string {
	name, _ := transactionResultCodeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e TransactionResultCode) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := transactionResultCodeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid TransactionResultCode enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*TransactionResultCode)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *TransactionResultCode) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding TransactionResultCode: %s", err)
	}
	if _, ok := transactionResultCodeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid TransactionResultCode enum value", v)
	}
	*e = TransactionResultCode(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionResultCode) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionResultCode) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionResultCode)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionResultCode)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TransactionResultCode) xdrType() {}

var _ xdrType = (*TransactionResultCode)(nil)

// InnerTransactionResultResult is an XDR NestedUnion defines as:
//
//   union switch (TransactionResultCode code)
//        {
//        // txFEE_BUMP_INNER_SUCCESS is not included
//        case txSUCCESS:
//        case txFAILED:
//            OperationResult results<>;
//        case txTOO_EARLY:
//        case txTOO_LATE:
//        case txMISSING_OPERATION:
//        case txBAD_SEQ:
//        case txBAD_AUTH:
//        case txINSUFFICIENT_BALANCE:
//        case txNO_ACCOUNT:
//        case txINSUFFICIENT_FEE:
//        case txBAD_AUTH_EXTRA:
//        case txINTERNAL_ERROR:
//        case txNOT_SUPPORTED:
//        // txFEE_BUMP_INNER_FAILED is not included
//        case txBAD_SPONSORSHIP:
//            void;
//        }
//
type InnerTransactionResultResult struct {
	Code    TransactionResultCode
	Results *[]OperationResult
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u InnerTransactionResultResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of InnerTransactionResultResult
func (u InnerTransactionResultResult) ArmForSwitch(sw int32) (string, bool) {
	switch TransactionResultCode(sw) {
	case TransactionResultCodeTxSuccess:
		return "Results", true
	case TransactionResultCodeTxFailed:
		return "Results", true
	case TransactionResultCodeTxTooEarly:
		return "", true
	case TransactionResultCodeTxTooLate:
		return "", true
	case TransactionResultCodeTxMissingOperation:
		return "", true
	case TransactionResultCodeTxBadSeq:
		return "", true
	case TransactionResultCodeTxBadAuth:
		return "", true
	case TransactionResultCodeTxInsufficientBalance:
		return "", true
	case TransactionResultCodeTxNoAccount:
		return "", true
	case TransactionResultCodeTxInsufficientFee:
		return "", true
	case TransactionResultCodeTxBadAuthExtra:
		return "", true
	case TransactionResultCodeTxInternalError:
		return "", true
	case TransactionResultCodeTxNotSupported:
		return "", true
	case TransactionResultCodeTxBadSponsorship:
		return "", true
	}
	return "-", false
}

// NewInnerTransactionResultResult creates a new  InnerTransactionResultResult.
func NewInnerTransactionResultResult(code TransactionResultCode, value interface{}) (result InnerTransactionResultResult, err error) {
	result.Code = code
	switch TransactionResultCode(code) {
	case TransactionResultCodeTxSuccess:
		tv, ok := value.([]OperationResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be []OperationResult")
			return
		}
		result.Results = &tv
	case TransactionResultCodeTxFailed:
		tv, ok := value.([]OperationResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be []OperationResult")
			return
		}
		result.Results = &tv
	case TransactionResultCodeTxTooEarly:
		// void
	case TransactionResultCodeTxTooLate:
		// void
	case TransactionResultCodeTxMissingOperation:
		// void
	case TransactionResultCodeTxBadSeq:
		// void
	case TransactionResultCodeTxBadAuth:
		// void
	case TransactionResultCodeTxInsufficientBalance:
		// void
	case TransactionResultCodeTxNoAccount:
		// void
	case TransactionResultCodeTxInsufficientFee:
		// void
	case TransactionResultCodeTxBadAuthExtra:
		// void
	case TransactionResultCodeTxInternalError:
		// void
	case TransactionResultCodeTxNotSupported:
		// void
	case TransactionResultCodeTxBadSponsorship:
		// void
	}
	return
}

// MustResults retrieves the Results value from the union,
// panicing if the value is not set.
func (u InnerTransactionResultResult) MustResults() []OperationResult {
	val, ok := u.GetResults()

	if !ok {
		panic("arm Results is not set")
	}

	return val
}

// GetResults retrieves the Results value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u InnerTransactionResultResult) GetResults() (result []OperationResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Results" {
		result = *u.Results
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u InnerTransactionResultResult) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Code.EncodeTo(e); err != nil {
		return err
	}
	switch TransactionResultCode(u.Code) {
	case TransactionResultCodeTxSuccess:
		if _, err = e.EncodeUint(uint32(len((*u.Results)))); err != nil {
			return err
		}
		for i := 0; i < len((*u.Results)); i++ {
			if err = (*u.Results)[i].EncodeTo(e); err != nil {
				return err
			}
		}
		return nil
	case TransactionResultCodeTxFailed:
		if _, err = e.EncodeUint(uint32(len((*u.Results)))); err != nil {
			return err
		}
		for i := 0; i < len((*u.Results)); i++ {
			if err = (*u.Results)[i].EncodeTo(e); err != nil {
				return err
			}
		}
		return nil
	case TransactionResultCodeTxTooEarly:
		// Void
		return nil
	case TransactionResultCodeTxTooLate:
		// Void
		return nil
	case TransactionResultCodeTxMissingOperation:
		// Void
		return nil
	case TransactionResultCodeTxBadSeq:
		// Void
		return nil
	case TransactionResultCodeTxBadAuth:
		// Void
		return nil
	case TransactionResultCodeTxInsufficientBalance:
		// Void
		return nil
	case TransactionResultCodeTxNoAccount:
		// Void
		return nil
	case TransactionResultCodeTxInsufficientFee:
		// Void
		return nil
	case TransactionResultCodeTxBadAuthExtra:
		// Void
		return nil
	case TransactionResultCodeTxInternalError:
		// Void
		return nil
	case TransactionResultCodeTxNotSupported:
		// Void
		return nil
	case TransactionResultCodeTxBadSponsorship:
		// Void
		return nil
	}
	return fmt.Errorf("Code (TransactionResultCode) switch value '%d' is not valid for union InnerTransactionResultResult", u.Code)
}

var _ decoderFrom = (*InnerTransactionResultResult)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *InnerTransactionResultResult) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Code.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TransactionResultCode: %s", err)
	}
	switch TransactionResultCode(u.Code) {
	case TransactionResultCodeTxSuccess:
		u.Results = new([]OperationResult)
		var l uint32
		l, nTmp, err = d.DecodeUint()
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding OperationResult: %s", err)
		}
		(*u.Results) = nil
		if l > 0 {
			(*u.Results) = make([]OperationResult, l)
			for i := uint32(0); i < l; i++ {
				nTmp, err = (*u.Results)[i].DecodeFrom(d)
				n += nTmp
				if err != nil {
					return n, fmt.Errorf("decoding OperationResult: %s", err)
				}
			}
		}
		return n, nil
	case TransactionResultCodeTxFailed:
		u.Results = new([]OperationResult)
		var l uint32
		l, nTmp, err = d.DecodeUint()
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding OperationResult: %s", err)
		}
		(*u.Results) = nil
		if l > 0 {
			(*u.Results) = make([]OperationResult, l)
			for i := uint32(0); i < l; i++ {
				nTmp, err = (*u.Results)[i].DecodeFrom(d)
				n += nTmp
				if err != nil {
					return n, fmt.Errorf("decoding OperationResult: %s", err)
				}
			}
		}
		return n, nil
	case TransactionResultCodeTxTooEarly:
		// Void
		return n, nil
	case TransactionResultCodeTxTooLate:
		// Void
		return n, nil
	case TransactionResultCodeTxMissingOperation:
		// Void
		return n, nil
	case TransactionResultCodeTxBadSeq:
		// Void
		return n, nil
	case TransactionResultCodeTxBadAuth:
		// Void
		return n, nil
	case TransactionResultCodeTxInsufficientBalance:
		// Void
		return n, nil
	case TransactionResultCodeTxNoAccount:
		// Void
		return n, nil
	case TransactionResultCodeTxInsufficientFee:
		// Void
		return n, nil
	case TransactionResultCodeTxBadAuthExtra:
		// Void
		return n, nil
	case TransactionResultCodeTxInternalError:
		// Void
		return n, nil
	case TransactionResultCodeTxNotSupported:
		// Void
		return n, nil
	case TransactionResultCodeTxBadSponsorship:
		// Void
		return n, nil
	}
	return n, fmt.Errorf("union InnerTransactionResultResult has invalid Code (TransactionResultCode) switch value '%d'", u.Code)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s InnerTransactionResultResult) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *InnerTransactionResultResult) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*InnerTransactionResultResult)(nil)
	_ encoding.BinaryUnmarshaler = (*InnerTransactionResultResult)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s InnerTransactionResultResult) xdrType() {}

var _ xdrType = (*InnerTransactionResultResult)(nil)

// InnerTransactionResultExt is an XDR NestedUnion defines as:
//
//   union switch (int v)
//        {
//        case 0:
//            void;
//        }
//
type InnerTransactionResultExt struct {
	V int32
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u InnerTransactionResultExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of InnerTransactionResultExt
func (u InnerTransactionResultExt) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	}
	return "-", false
}

// NewInnerTransactionResultExt creates a new  InnerTransactionResultExt.
func NewInnerTransactionResultExt(v int32, value interface{}) (result InnerTransactionResultExt, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u InnerTransactionResultExt) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union InnerTransactionResultExt", u.V)
}

var _ decoderFrom = (*InnerTransactionResultExt)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *InnerTransactionResultExt) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %s", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	}
	return n, fmt.Errorf("union InnerTransactionResultExt has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s InnerTransactionResultExt) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *InnerTransactionResultExt) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*InnerTransactionResultExt)(nil)
	_ encoding.BinaryUnmarshaler = (*InnerTransactionResultExt)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s InnerTransactionResultExt) xdrType() {}

var _ xdrType = (*InnerTransactionResultExt)(nil)

// InnerTransactionResult is an XDR Struct defines as:
//
//   struct InnerTransactionResult
//    {
//        // Always 0. Here for binary compatibility.
//        int64 feeCharged;
//
//        union switch (TransactionResultCode code)
//        {
//        // txFEE_BUMP_INNER_SUCCESS is not included
//        case txSUCCESS:
//        case txFAILED:
//            OperationResult results<>;
//        case txTOO_EARLY:
//        case txTOO_LATE:
//        case txMISSING_OPERATION:
//        case txBAD_SEQ:
//        case txBAD_AUTH:
//        case txINSUFFICIENT_BALANCE:
//        case txNO_ACCOUNT:
//        case txINSUFFICIENT_FEE:
//        case txBAD_AUTH_EXTRA:
//        case txINTERNAL_ERROR:
//        case txNOT_SUPPORTED:
//        // txFEE_BUMP_INNER_FAILED is not included
//        case txBAD_SPONSORSHIP:
//            void;
//        }
//        result;
//
//        // reserved for future use
//        union switch (int v)
//        {
//        case 0:
//            void;
//        }
//        ext;
//    };
//
type InnerTransactionResult struct {
	FeeCharged Int64
	Result     InnerTransactionResultResult
	Ext        InnerTransactionResultExt
}

// EncodeTo encodes this value using the Encoder.
func (s *InnerTransactionResult) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.FeeCharged.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Result.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*InnerTransactionResult)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *InnerTransactionResult) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.FeeCharged.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	nTmp, err = s.Result.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding InnerTransactionResultResult: %s", err)
	}
	nTmp, err = s.Ext.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding InnerTransactionResultExt: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s InnerTransactionResult) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *InnerTransactionResult) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*InnerTransactionResult)(nil)
	_ encoding.BinaryUnmarshaler = (*InnerTransactionResult)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s InnerTransactionResult) xdrType() {}

var _ xdrType = (*InnerTransactionResult)(nil)

// InnerTransactionResultPair is an XDR Struct defines as:
//
//   struct InnerTransactionResultPair
//    {
//        Hash transactionHash;          // hash of the inner transaction
//        InnerTransactionResult result; // result for the inner transaction
//    };
//
type InnerTransactionResultPair struct {
	TransactionHash Hash
	Result          InnerTransactionResult
}

// EncodeTo encodes this value using the Encoder.
func (s *InnerTransactionResultPair) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.TransactionHash.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Result.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*InnerTransactionResultPair)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *InnerTransactionResultPair) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.TransactionHash.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Hash: %s", err)
	}
	nTmp, err = s.Result.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding InnerTransactionResult: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s InnerTransactionResultPair) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *InnerTransactionResultPair) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*InnerTransactionResultPair)(nil)
	_ encoding.BinaryUnmarshaler = (*InnerTransactionResultPair)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s InnerTransactionResultPair) xdrType() {}

var _ xdrType = (*InnerTransactionResultPair)(nil)

// TransactionResultResult is an XDR NestedUnion defines as:
//
//   union switch (TransactionResultCode code)
//        {
//        case txFEE_BUMP_INNER_SUCCESS:
//        case txFEE_BUMP_INNER_FAILED:
//            InnerTransactionResultPair innerResultPair;
//        case txSUCCESS:
//        case txFAILED:
//            OperationResult results<>;
//        default:
//            void;
//        }
//
type TransactionResultResult struct {
	Code            TransactionResultCode
	InnerResultPair *InnerTransactionResultPair
	Results         *[]OperationResult
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TransactionResultResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TransactionResultResult
func (u TransactionResultResult) ArmForSwitch(sw int32) (string, bool) {
	switch TransactionResultCode(sw) {
	case TransactionResultCodeTxFeeBumpInnerSuccess:
		return "InnerResultPair", true
	case TransactionResultCodeTxFeeBumpInnerFailed:
		return "InnerResultPair", true
	case TransactionResultCodeTxSuccess:
		return "Results", true
	case TransactionResultCodeTxFailed:
		return "Results", true
	default:
		return "", true
	}
}

// NewTransactionResultResult creates a new  TransactionResultResult.
func NewTransactionResultResult(code TransactionResultCode, value interface{}) (result TransactionResultResult, err error) {
	result.Code = code
	switch TransactionResultCode(code) {
	case TransactionResultCodeTxFeeBumpInnerSuccess:
		tv, ok := value.(InnerTransactionResultPair)
		if !ok {
			err = fmt.Errorf("invalid value, must be InnerTransactionResultPair")
			return
		}
		result.InnerResultPair = &tv
	case TransactionResultCodeTxFeeBumpInnerFailed:
		tv, ok := value.(InnerTransactionResultPair)
		if !ok {
			err = fmt.Errorf("invalid value, must be InnerTransactionResultPair")
			return
		}
		result.InnerResultPair = &tv
	case TransactionResultCodeTxSuccess:
		tv, ok := value.([]OperationResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be []OperationResult")
			return
		}
		result.Results = &tv
	case TransactionResultCodeTxFailed:
		tv, ok := value.([]OperationResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be []OperationResult")
			return
		}
		result.Results = &tv
	default:
		// void
	}
	return
}

// MustInnerResultPair retrieves the InnerResultPair value from the union,
// panicing if the value is not set.
func (u TransactionResultResult) MustInnerResultPair() InnerTransactionResultPair {
	val, ok := u.GetInnerResultPair()

	if !ok {
		panic("arm InnerResultPair is not set")
	}

	return val
}

// GetInnerResultPair retrieves the InnerResultPair value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u TransactionResultResult) GetInnerResultPair() (result InnerTransactionResultPair, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "InnerResultPair" {
		result = *u.InnerResultPair
		ok = true
	}

	return
}

// MustResults retrieves the Results value from the union,
// panicing if the value is not set.
func (u TransactionResultResult) MustResults() []OperationResult {
	val, ok := u.GetResults()

	if !ok {
		panic("arm Results is not set")
	}

	return val
}

// GetResults retrieves the Results value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u TransactionResultResult) GetResults() (result []OperationResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Results" {
		result = *u.Results
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u TransactionResultResult) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Code.EncodeTo(e); err != nil {
		return err
	}
	switch TransactionResultCode(u.Code) {
	case TransactionResultCodeTxFeeBumpInnerSuccess:
		if err = (*u.InnerResultPair).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case TransactionResultCodeTxFeeBumpInnerFailed:
		if err = (*u.InnerResultPair).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case TransactionResultCodeTxSuccess:
		if _, err = e.EncodeUint(uint32(len((*u.Results)))); err != nil {
			return err
		}
		for i := 0; i < len((*u.Results)); i++ {
			if err = (*u.Results)[i].EncodeTo(e); err != nil {
				return err
			}
		}
		return nil
	case TransactionResultCodeTxFailed:
		if _, err = e.EncodeUint(uint32(len((*u.Results)))); err != nil {
			return err
		}
		for i := 0; i < len((*u.Results)); i++ {
			if err = (*u.Results)[i].EncodeTo(e); err != nil {
				return err
			}
		}
		return nil
	default:
		// Void
		return nil
	}
}

var _ decoderFrom = (*TransactionResultResult)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *TransactionResultResult) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Code.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TransactionResultCode: %s", err)
	}
	switch TransactionResultCode(u.Code) {
	case TransactionResultCodeTxFeeBumpInnerSuccess:
		u.InnerResultPair = new(InnerTransactionResultPair)
		nTmp, err = (*u.InnerResultPair).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding InnerTransactionResultPair: %s", err)
		}
		return n, nil
	case TransactionResultCodeTxFeeBumpInnerFailed:
		u.InnerResultPair = new(InnerTransactionResultPair)
		nTmp, err = (*u.InnerResultPair).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding InnerTransactionResultPair: %s", err)
		}
		return n, nil
	case TransactionResultCodeTxSuccess:
		u.Results = new([]OperationResult)
		var l uint32
		l, nTmp, err = d.DecodeUint()
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding OperationResult: %s", err)
		}
		(*u.Results) = nil
		if l > 0 {
			(*u.Results) = make([]OperationResult, l)
			for i := uint32(0); i < l; i++ {
				nTmp, err = (*u.Results)[i].DecodeFrom(d)
				n += nTmp
				if err != nil {
					return n, fmt.Errorf("decoding OperationResult: %s", err)
				}
			}
		}
		return n, nil
	case TransactionResultCodeTxFailed:
		u.Results = new([]OperationResult)
		var l uint32
		l, nTmp, err = d.DecodeUint()
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding OperationResult: %s", err)
		}
		(*u.Results) = nil
		if l > 0 {
			(*u.Results) = make([]OperationResult, l)
			for i := uint32(0); i < l; i++ {
				nTmp, err = (*u.Results)[i].DecodeFrom(d)
				n += nTmp
				if err != nil {
					return n, fmt.Errorf("decoding OperationResult: %s", err)
				}
			}
		}
		return n, nil
	default:
		// Void
		return n, nil
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionResultResult) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionResultResult) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionResultResult)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionResultResult)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TransactionResultResult) xdrType() {}

var _ xdrType = (*TransactionResultResult)(nil)

// TransactionResultExt is an XDR NestedUnion defines as:
//
//   union switch (int v)
//        {
//        case 0:
//            void;
//        }
//
type TransactionResultExt struct {
	V int32
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TransactionResultExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TransactionResultExt
func (u TransactionResultExt) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	}
	return "-", false
}

// NewTransactionResultExt creates a new  TransactionResultExt.
func NewTransactionResultExt(v int32, value interface{}) (result TransactionResultExt, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u TransactionResultExt) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union TransactionResultExt", u.V)
}

var _ decoderFrom = (*TransactionResultExt)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *TransactionResultExt) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %s", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	}
	return n, fmt.Errorf("union TransactionResultExt has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionResultExt) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionResultExt) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionResultExt)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionResultExt)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TransactionResultExt) xdrType() {}

var _ xdrType = (*TransactionResultExt)(nil)

// TransactionResult is an XDR Struct defines as:
//
//   struct TransactionResult
//    {
//        int64 feeCharged; // actual fee charged for the transaction
//
//        union switch (TransactionResultCode code)
//        {
//        case txFEE_BUMP_INNER_SUCCESS:
//        case txFEE_BUMP_INNER_FAILED:
//            InnerTransactionResultPair innerResultPair;
//        case txSUCCESS:
//        case txFAILED:
//            OperationResult results<>;
//        default:
//            void;
//        }
//        result;
//
//        // reserved for future use
//        union switch (int v)
//        {
//        case 0:
//            void;
//        }
//        ext;
//    };
//
type TransactionResult struct {
	FeeCharged Int64
	Result     TransactionResultResult
	Ext        TransactionResultExt
}

// EncodeTo encodes this value using the Encoder.
func (s *TransactionResult) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.FeeCharged.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Result.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*TransactionResult)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *TransactionResult) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = s.FeeCharged.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %s", err)
	}
	nTmp, err = s.Result.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TransactionResultResult: %s", err)
	}
	nTmp, err = s.Ext.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TransactionResultExt: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionResult) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionResult) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionResult)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionResult)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s TransactionResult) xdrType() {}

var _ xdrType = (*TransactionResult)(nil)

// Hash is an XDR Typedef defines as:
//
//   typedef opaque Hash[32];
//
type Hash [32]byte

// XDRMaxSize implements the Sized interface for Hash
func (e Hash) XDRMaxSize() int {
	return 32
}

// EncodeTo encodes this value using the Encoder.
func (s *Hash) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeFixedOpaque(s[:]); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Hash)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Hash) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = d.DecodeFixedOpaqueInplace(s[:])
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Hash: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Hash) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Hash) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Hash)(nil)
	_ encoding.BinaryUnmarshaler = (*Hash)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s Hash) xdrType() {}

var _ xdrType = (*Hash)(nil)

// Uint256 is an XDR Typedef defines as:
//
//   typedef opaque uint256[32];
//
type Uint256 [32]byte

// XDRMaxSize implements the Sized interface for Uint256
func (e Uint256) XDRMaxSize() int {
	return 32
}

// EncodeTo encodes this value using the Encoder.
func (s *Uint256) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeFixedOpaque(s[:]); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Uint256)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Uint256) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = d.DecodeFixedOpaqueInplace(s[:])
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint256: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Uint256) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Uint256) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Uint256)(nil)
	_ encoding.BinaryUnmarshaler = (*Uint256)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s Uint256) xdrType() {}

var _ xdrType = (*Uint256)(nil)

// Uint32 is an XDR Typedef defines as:
//
//   typedef unsigned int uint32;
//
type Uint32 uint32

// EncodeTo encodes this value using the Encoder.
func (s Uint32) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeUint(uint32(s)); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Uint32)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Uint32) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	var v uint32
	v, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Unsigned int: %s", err)
	}
	*s = Uint32(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Uint32) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Uint32) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Uint32)(nil)
	_ encoding.BinaryUnmarshaler = (*Uint32)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s Uint32) xdrType() {}

var _ xdrType = (*Uint32)(nil)

// Int32 is an XDR Typedef defines as:
//
//   typedef int int32;
//
type Int32 int32

// EncodeTo encodes this value using the Encoder.
func (s Int32) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(s)); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Int32)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Int32) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	var v int32
	v, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %s", err)
	}
	*s = Int32(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Int32) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Int32) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Int32)(nil)
	_ encoding.BinaryUnmarshaler = (*Int32)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s Int32) xdrType() {}

var _ xdrType = (*Int32)(nil)

// Uint64 is an XDR Typedef defines as:
//
//   typedef unsigned hyper uint64;
//
type Uint64 uint64

// EncodeTo encodes this value using the Encoder.
func (s Uint64) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeUhyper(uint64(s)); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Uint64)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Uint64) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	var v uint64
	v, nTmp, err = d.DecodeUhyper()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Unsigned hyper: %s", err)
	}
	*s = Uint64(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Uint64) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Uint64) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Uint64)(nil)
	_ encoding.BinaryUnmarshaler = (*Uint64)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s Uint64) xdrType() {}

var _ xdrType = (*Uint64)(nil)

// Int64 is an XDR Typedef defines as:
//
//   typedef hyper int64;
//
type Int64 int64

// EncodeTo encodes this value using the Encoder.
func (s Int64) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeHyper(int64(s)); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Int64)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Int64) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	var v int64
	v, nTmp, err = d.DecodeHyper()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Hyper: %s", err)
	}
	*s = Int64(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Int64) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Int64) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Int64)(nil)
	_ encoding.BinaryUnmarshaler = (*Int64)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s Int64) xdrType() {}

var _ xdrType = (*Int64)(nil)

// CryptoKeyType is an XDR Enum defines as:
//
//   enum CryptoKeyType
//    {
//        KEY_TYPE_ED25519 = 0,
//        KEY_TYPE_PRE_AUTH_TX = 1,
//        KEY_TYPE_HASH_X = 2,
//        // MUXED enum values for supported type are derived from the enum values
//        // above by ORing them with 0x100
//        KEY_TYPE_MUXED_ED25519 = 0x100
//    };
//
type CryptoKeyType int32

const (
	CryptoKeyTypeKeyTypeEd25519      CryptoKeyType = 0
	CryptoKeyTypeKeyTypePreAuthTx    CryptoKeyType = 1
	CryptoKeyTypeKeyTypeHashX        CryptoKeyType = 2
	CryptoKeyTypeKeyTypeMuxedEd25519 CryptoKeyType = 256
)

var cryptoKeyTypeMap = map[int32]string{
	0:   "CryptoKeyTypeKeyTypeEd25519",
	1:   "CryptoKeyTypeKeyTypePreAuthTx",
	2:   "CryptoKeyTypeKeyTypeHashX",
	256: "CryptoKeyTypeKeyTypeMuxedEd25519",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for CryptoKeyType
func (e CryptoKeyType) ValidEnum(v int32) bool {
	_, ok := cryptoKeyTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e CryptoKeyType) String() string {
	name, _ := cryptoKeyTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e CryptoKeyType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := cryptoKeyTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid CryptoKeyType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*CryptoKeyType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *CryptoKeyType) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding CryptoKeyType: %s", err)
	}
	if _, ok := cryptoKeyTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid CryptoKeyType enum value", v)
	}
	*e = CryptoKeyType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s CryptoKeyType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *CryptoKeyType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*CryptoKeyType)(nil)
	_ encoding.BinaryUnmarshaler = (*CryptoKeyType)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s CryptoKeyType) xdrType() {}

var _ xdrType = (*CryptoKeyType)(nil)

// PublicKeyType is an XDR Enum defines as:
//
//   enum PublicKeyType
//    {
//        PUBLIC_KEY_TYPE_ED25519 = KEY_TYPE_ED25519
//    };
//
type PublicKeyType int32

const (
	PublicKeyTypePublicKeyTypeEd25519 PublicKeyType = 0
)

var publicKeyTypeMap = map[int32]string{
	0: "PublicKeyTypePublicKeyTypeEd25519",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for PublicKeyType
func (e PublicKeyType) ValidEnum(v int32) bool {
	_, ok := publicKeyTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e PublicKeyType) String() string {
	name, _ := publicKeyTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e PublicKeyType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := publicKeyTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid PublicKeyType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*PublicKeyType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *PublicKeyType) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding PublicKeyType: %s", err)
	}
	if _, ok := publicKeyTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid PublicKeyType enum value", v)
	}
	*e = PublicKeyType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s PublicKeyType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *PublicKeyType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*PublicKeyType)(nil)
	_ encoding.BinaryUnmarshaler = (*PublicKeyType)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s PublicKeyType) xdrType() {}

var _ xdrType = (*PublicKeyType)(nil)

// SignerKeyType is an XDR Enum defines as:
//
//   enum SignerKeyType
//    {
//        SIGNER_KEY_TYPE_ED25519 = KEY_TYPE_ED25519,
//        SIGNER_KEY_TYPE_PRE_AUTH_TX = KEY_TYPE_PRE_AUTH_TX,
//        SIGNER_KEY_TYPE_HASH_X = KEY_TYPE_HASH_X
//    };
//
type SignerKeyType int32

const (
	SignerKeyTypeSignerKeyTypeEd25519   SignerKeyType = 0
	SignerKeyTypeSignerKeyTypePreAuthTx SignerKeyType = 1
	SignerKeyTypeSignerKeyTypeHashX     SignerKeyType = 2
)

var signerKeyTypeMap = map[int32]string{
	0: "SignerKeyTypeSignerKeyTypeEd25519",
	1: "SignerKeyTypeSignerKeyTypePreAuthTx",
	2: "SignerKeyTypeSignerKeyTypeHashX",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for SignerKeyType
func (e SignerKeyType) ValidEnum(v int32) bool {
	_, ok := signerKeyTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e SignerKeyType) String() string {
	name, _ := signerKeyTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e SignerKeyType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := signerKeyTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid SignerKeyType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*SignerKeyType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *SignerKeyType) DecodeFrom(d *xdr.Decoder) (int, error) {
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding SignerKeyType: %s", err)
	}
	if _, ok := signerKeyTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid SignerKeyType enum value", v)
	}
	*e = SignerKeyType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s SignerKeyType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SignerKeyType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*SignerKeyType)(nil)
	_ encoding.BinaryUnmarshaler = (*SignerKeyType)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s SignerKeyType) xdrType() {}

var _ xdrType = (*SignerKeyType)(nil)

// PublicKey is an XDR Union defines as:
//
//   union PublicKey switch (PublicKeyType type)
//    {
//    case PUBLIC_KEY_TYPE_ED25519:
//        uint256 ed25519;
//    };
//
type PublicKey struct {
	Type    PublicKeyType
	Ed25519 *Uint256
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u PublicKey) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PublicKey
func (u PublicKey) ArmForSwitch(sw int32) (string, bool) {
	switch PublicKeyType(sw) {
	case PublicKeyTypePublicKeyTypeEd25519:
		return "Ed25519", true
	}
	return "-", false
}

// NewPublicKey creates a new  PublicKey.
func NewPublicKey(aType PublicKeyType, value interface{}) (result PublicKey, err error) {
	result.Type = aType
	switch PublicKeyType(aType) {
	case PublicKeyTypePublicKeyTypeEd25519:
		tv, ok := value.(Uint256)
		if !ok {
			err = fmt.Errorf("invalid value, must be Uint256")
			return
		}
		result.Ed25519 = &tv
	}
	return
}

// MustEd25519 retrieves the Ed25519 value from the union,
// panicing if the value is not set.
func (u PublicKey) MustEd25519() Uint256 {
	val, ok := u.GetEd25519()

	if !ok {
		panic("arm Ed25519 is not set")
	}

	return val
}

// GetEd25519 retrieves the Ed25519 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u PublicKey) GetEd25519() (result Uint256, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Ed25519" {
		result = *u.Ed25519
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u PublicKey) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch PublicKeyType(u.Type) {
	case PublicKeyTypePublicKeyTypeEd25519:
		if err = (*u.Ed25519).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (PublicKeyType) switch value '%d' is not valid for union PublicKey", u.Type)
}

var _ decoderFrom = (*PublicKey)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *PublicKey) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding PublicKeyType: %s", err)
	}
	switch PublicKeyType(u.Type) {
	case PublicKeyTypePublicKeyTypeEd25519:
		u.Ed25519 = new(Uint256)
		nTmp, err = (*u.Ed25519).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Uint256: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union PublicKey has invalid Type (PublicKeyType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s PublicKey) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *PublicKey) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*PublicKey)(nil)
	_ encoding.BinaryUnmarshaler = (*PublicKey)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s PublicKey) xdrType() {}

var _ xdrType = (*PublicKey)(nil)

// SignerKey is an XDR Union defines as:
//
//   union SignerKey switch (SignerKeyType type)
//    {
//    case SIGNER_KEY_TYPE_ED25519:
//        uint256 ed25519;
//    case SIGNER_KEY_TYPE_PRE_AUTH_TX:
//        /* SHA-256 Hash of TransactionSignaturePayload structure */
//        uint256 preAuthTx;
//    case SIGNER_KEY_TYPE_HASH_X:
//        /* Hash of random 256 bit preimage X */
//        uint256 hashX;
//    };
//
type SignerKey struct {
	Type      SignerKeyType
	Ed25519   *Uint256
	PreAuthTx *Uint256
	HashX     *Uint256
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u SignerKey) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of SignerKey
func (u SignerKey) ArmForSwitch(sw int32) (string, bool) {
	switch SignerKeyType(sw) {
	case SignerKeyTypeSignerKeyTypeEd25519:
		return "Ed25519", true
	case SignerKeyTypeSignerKeyTypePreAuthTx:
		return "PreAuthTx", true
	case SignerKeyTypeSignerKeyTypeHashX:
		return "HashX", true
	}
	return "-", false
}

// NewSignerKey creates a new  SignerKey.
func NewSignerKey(aType SignerKeyType, value interface{}) (result SignerKey, err error) {
	result.Type = aType
	switch SignerKeyType(aType) {
	case SignerKeyTypeSignerKeyTypeEd25519:
		tv, ok := value.(Uint256)
		if !ok {
			err = fmt.Errorf("invalid value, must be Uint256")
			return
		}
		result.Ed25519 = &tv
	case SignerKeyTypeSignerKeyTypePreAuthTx:
		tv, ok := value.(Uint256)
		if !ok {
			err = fmt.Errorf("invalid value, must be Uint256")
			return
		}
		result.PreAuthTx = &tv
	case SignerKeyTypeSignerKeyTypeHashX:
		tv, ok := value.(Uint256)
		if !ok {
			err = fmt.Errorf("invalid value, must be Uint256")
			return
		}
		result.HashX = &tv
	}
	return
}

// MustEd25519 retrieves the Ed25519 value from the union,
// panicing if the value is not set.
func (u SignerKey) MustEd25519() Uint256 {
	val, ok := u.GetEd25519()

	if !ok {
		panic("arm Ed25519 is not set")
	}

	return val
}

// GetEd25519 retrieves the Ed25519 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u SignerKey) GetEd25519() (result Uint256, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Ed25519" {
		result = *u.Ed25519
		ok = true
	}

	return
}

// MustPreAuthTx retrieves the PreAuthTx value from the union,
// panicing if the value is not set.
func (u SignerKey) MustPreAuthTx() Uint256 {
	val, ok := u.GetPreAuthTx()

	if !ok {
		panic("arm PreAuthTx is not set")
	}

	return val
}

// GetPreAuthTx retrieves the PreAuthTx value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u SignerKey) GetPreAuthTx() (result Uint256, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "PreAuthTx" {
		result = *u.PreAuthTx
		ok = true
	}

	return
}

// MustHashX retrieves the HashX value from the union,
// panicing if the value is not set.
func (u SignerKey) MustHashX() Uint256 {
	val, ok := u.GetHashX()

	if !ok {
		panic("arm HashX is not set")
	}

	return val
}

// GetHashX retrieves the HashX value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u SignerKey) GetHashX() (result Uint256, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "HashX" {
		result = *u.HashX
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u SignerKey) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch SignerKeyType(u.Type) {
	case SignerKeyTypeSignerKeyTypeEd25519:
		if err = (*u.Ed25519).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case SignerKeyTypeSignerKeyTypePreAuthTx:
		if err = (*u.PreAuthTx).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case SignerKeyTypeSignerKeyTypeHashX:
		if err = (*u.HashX).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (SignerKeyType) switch value '%d' is not valid for union SignerKey", u.Type)
}

var _ decoderFrom = (*SignerKey)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *SignerKey) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SignerKeyType: %s", err)
	}
	switch SignerKeyType(u.Type) {
	case SignerKeyTypeSignerKeyTypeEd25519:
		u.Ed25519 = new(Uint256)
		nTmp, err = (*u.Ed25519).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Uint256: %s", err)
		}
		return n, nil
	case SignerKeyTypeSignerKeyTypePreAuthTx:
		u.PreAuthTx = new(Uint256)
		nTmp, err = (*u.PreAuthTx).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Uint256: %s", err)
		}
		return n, nil
	case SignerKeyTypeSignerKeyTypeHashX:
		u.HashX = new(Uint256)
		nTmp, err = (*u.HashX).DecodeFrom(d)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Uint256: %s", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union SignerKey has invalid Type (SignerKeyType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s SignerKey) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SignerKey) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*SignerKey)(nil)
	_ encoding.BinaryUnmarshaler = (*SignerKey)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s SignerKey) xdrType() {}

var _ xdrType = (*SignerKey)(nil)

// Signature is an XDR Typedef defines as:
//
//   typedef opaque Signature<64>;
//
type Signature []byte

// XDRMaxSize implements the Sized interface for Signature
func (e Signature) XDRMaxSize() int {
	return 64
}

// EncodeTo encodes this value using the Encoder.
func (s Signature) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeOpaque(s[:]); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Signature)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Signature) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	(*s), nTmp, err = d.DecodeOpaque(64)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Signature: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Signature) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Signature) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Signature)(nil)
	_ encoding.BinaryUnmarshaler = (*Signature)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s Signature) xdrType() {}

var _ xdrType = (*Signature)(nil)

// SignatureHint is an XDR Typedef defines as:
//
//   typedef opaque SignatureHint[4];
//
type SignatureHint [4]byte

// XDRMaxSize implements the Sized interface for SignatureHint
func (e SignatureHint) XDRMaxSize() int {
	return 4
}

// EncodeTo encodes this value using the Encoder.
func (s *SignatureHint) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeFixedOpaque(s[:]); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*SignatureHint)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *SignatureHint) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = d.DecodeFixedOpaqueInplace(s[:])
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SignatureHint: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s SignatureHint) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SignatureHint) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*SignatureHint)(nil)
	_ encoding.BinaryUnmarshaler = (*SignatureHint)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s SignatureHint) xdrType() {}

var _ xdrType = (*SignatureHint)(nil)

// NodeId is an XDR Typedef defines as:
//
//   typedef PublicKey NodeID;
//
type NodeId PublicKey

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u NodeId) SwitchFieldName() string {
	return PublicKey(u).SwitchFieldName()
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PublicKey
func (u NodeId) ArmForSwitch(sw int32) (string, bool) {
	return PublicKey(u).ArmForSwitch(sw)
}

// NewNodeId creates a new  NodeId.
func NewNodeId(aType PublicKeyType, value interface{}) (result NodeId, err error) {
	u, err := NewPublicKey(aType, value)
	result = NodeId(u)
	return
}

// MustEd25519 retrieves the Ed25519 value from the union,
// panicing if the value is not set.
func (u NodeId) MustEd25519() Uint256 {
	return PublicKey(u).MustEd25519()
}

// GetEd25519 retrieves the Ed25519 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u NodeId) GetEd25519() (result Uint256, ok bool) {
	return PublicKey(u).GetEd25519()
}

// EncodeTo encodes this value using the Encoder.
func (s NodeId) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = PublicKey(s).EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*NodeId)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *NodeId) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = (*PublicKey)(s).DecodeFrom(d)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding PublicKey: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s NodeId) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *NodeId) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*NodeId)(nil)
	_ encoding.BinaryUnmarshaler = (*NodeId)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s NodeId) xdrType() {}

var _ xdrType = (*NodeId)(nil)

// Curve25519Secret is an XDR Struct defines as:
//
//   struct Curve25519Secret
//    {
//        opaque key[32];
//    };
//
type Curve25519Secret struct {
	Key [32]byte `xdrmaxsize:"32"`
}

// EncodeTo encodes this value using the Encoder.
func (s *Curve25519Secret) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeFixedOpaque(s.Key[:]); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Curve25519Secret)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Curve25519Secret) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = d.DecodeFixedOpaqueInplace(s.Key[:])
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Key: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Curve25519Secret) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Curve25519Secret) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Curve25519Secret)(nil)
	_ encoding.BinaryUnmarshaler = (*Curve25519Secret)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s Curve25519Secret) xdrType() {}

var _ xdrType = (*Curve25519Secret)(nil)

// Curve25519Public is an XDR Struct defines as:
//
//   struct Curve25519Public
//    {
//        opaque key[32];
//    };
//
type Curve25519Public struct {
	Key [32]byte `xdrmaxsize:"32"`
}

// EncodeTo encodes this value using the Encoder.
func (s *Curve25519Public) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeFixedOpaque(s.Key[:]); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Curve25519Public)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Curve25519Public) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = d.DecodeFixedOpaqueInplace(s.Key[:])
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Key: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Curve25519Public) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Curve25519Public) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Curve25519Public)(nil)
	_ encoding.BinaryUnmarshaler = (*Curve25519Public)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s Curve25519Public) xdrType() {}

var _ xdrType = (*Curve25519Public)(nil)

// HmacSha256Key is an XDR Struct defines as:
//
//   struct HmacSha256Key
//    {
//        opaque key[32];
//    };
//
type HmacSha256Key struct {
	Key [32]byte `xdrmaxsize:"32"`
}

// EncodeTo encodes this value using the Encoder.
func (s *HmacSha256Key) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeFixedOpaque(s.Key[:]); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*HmacSha256Key)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *HmacSha256Key) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = d.DecodeFixedOpaqueInplace(s.Key[:])
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Key: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s HmacSha256Key) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *HmacSha256Key) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*HmacSha256Key)(nil)
	_ encoding.BinaryUnmarshaler = (*HmacSha256Key)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s HmacSha256Key) xdrType() {}

var _ xdrType = (*HmacSha256Key)(nil)

// HmacSha256Mac is an XDR Struct defines as:
//
//   struct HmacSha256Mac
//    {
//        opaque mac[32];
//    };
//
type HmacSha256Mac struct {
	Mac [32]byte `xdrmaxsize:"32"`
}

// EncodeTo encodes this value using the Encoder.
func (s *HmacSha256Mac) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeFixedOpaque(s.Mac[:]); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*HmacSha256Mac)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *HmacSha256Mac) DecodeFrom(d *xdr.Decoder) (int, error) {
	var err error
	var n, nTmp int
	nTmp, err = d.DecodeFixedOpaqueInplace(s.Mac[:])
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Mac: %s", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s HmacSha256Mac) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *HmacSha256Mac) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	d := xdr.NewDecoder(r)
	_, err := s.DecodeFrom(d)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*HmacSha256Mac)(nil)
	_ encoding.BinaryUnmarshaler = (*HmacSha256Mac)(nil)
)

// xdrType signals that this type is an type representing
// representing XDR values defined by this package.
func (s HmacSha256Mac) xdrType() {}

var _ xdrType = (*HmacSha256Mac)(nil)

var fmtTest = fmt.Sprint("this is a dummy usage of fmt")
