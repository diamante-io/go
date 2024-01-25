package randxdr

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"

	goxdr "github.com/xdrpp/goxdr/xdr"
)

type randMarshaller struct {
	useTag       bool
	tag          uint32
	rand         *rand.Rand
	maxBytesSize uint32
	maxVecLen    uint32
	presets      []Preset
}

func (*randMarshaller) Sprintf(f string, args ...interface{}) string {
	return fmt.Sprintf(f, args...)
}

func (rm *randMarshaller) randomKey(m interface{}) int32 {
	keys := reflect.ValueOf(m).MapKeys()
	// the keys of a map in golang are returned in random order
	// here we sort the keys to ensure the selection is
	// deterministic for the same rand seed
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Int() < keys[j].Int()
	})
	return int32(keys[rm.rand.Intn(len(keys))].Int())
}

func (rm *randMarshaller) applyPreset(field string, i goxdr.XdrType) bool {
	for _, preset := range rm.presets {
		if preset.Selector(field, i) {
			switch goxdr.XdrBaseType(i).(type) {
			case goxdr.XdrEnum, goxdr.XdrNum32:
				rm.useTag = false
			}

			preset.Setter(rm, field, i)
			return true
		}
	}
	return false
}

// Marshal populates a given goxdr.XdrType with random values.
//
// Every complex goxdr.XdrType has functions like XdrRecurse() which
// allow you to visit subfields of the complex type. That is how
// randMarshaller is able to populate subfields of a complex type with
// random values.
//
// Note that randMarshaller is stateful because of how union types are handled.
// Therefore Marshal() should not be used concurrently.
func (rm *randMarshaller) Marshal(field string, i goxdr.XdrType) {
	if rm.applyPreset(field, i) {
		return
	}

	switch t := goxdr.XdrBaseType(i).(type) {
	case goxdr.XdrVarBytes:
		bound := t.XdrBound()
		if bound > rm.maxBytesSize {
			bound = rm.maxBytesSize
		}
		bound++
		bs := make([]byte, rm.rand.Uint32()%bound)
		rm.rand.Read(bs)
		t.SetByteSlice(bs)
	case goxdr.XdrBytes:
		// t.GetByteSlice() returns the underlying byte slice
		// rm.rand.Read() will fill that byte slice with random values
		rm.rand.Read(t.GetByteSlice())
	case goxdr.XdrVec:
		bound := t.XdrBound()
		if bound > rm.maxVecLen {
			bound = rm.maxVecLen
		}
		bound++
		vecLen := rm.rand.Uint32() % bound
		t.SetVecLen(vecLen)
		t.XdrMarshalN(rm, field, vecLen)
	case goxdr.XdrPtr:
		present := rm.rand.Uint32()&1 == 1
		t.SetPresent(present)
		t.XdrMarshalValue(rm, field)
	case *goxdr.XdrBool:
		t.SetU32(rm.rand.Uint32() & 1)
	case goxdr.XdrEnum:
		if rm.useTag {
			rm.useTag = false
			t.SetU32(rm.tag)
		} else {
			t.SetU32(uint32(rm.randomKey(t.XdrEnumNames())))
		}
	case goxdr.XdrNum32:
		if rm.useTag {
			rm.useTag = false
			t.SetU32(rm.tag)
		} else {
			t.SetU32(rm.rand.Uint32())
		}
	case goxdr.XdrNum64:
		t.SetU64(rm.rand.Uint64())
	case goxdr.XdrUnion:
		// If we have an XDR union we need to set the tag of the union.
		// However, there is no SetTag() function in the goxdr.XdrUnion interface.
		// We must rely on these two facts:
		// * when XdrRecurse() is called on a union, the first field which will be marshalled is the tag field
		// * the tag field can be one of two types: uint32 or enum
		if m := t.XdrValidTags(); m != nil {
			rm.tag = uint32(rm.randomKey(m))
			rm.useTag = true
			// The next field the marshaller will visit is the tag field.
			// Once the tag is set, we need to toggle rm.useTag to false.
		}
		t.XdrRecurse(rm, field)
	case goxdr.XdrAggregate:
		t.XdrRecurse(rm, field)
	default:
		panic(fmt.Sprintf("field %s has unexpected xdr type %v", field, t))
	}
}
