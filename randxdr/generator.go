package randxdr

import (
	"math/rand"

	goxdr "github.com/xdrpp/goxdr/xdr"
)

// Generator generates random XDR values.
type Generator struct {
	// MaxBytesSize configures the upper bound limit for variable length
	// opaque data and variable length strings
	// https://tools.ietf.org/html/rfc4506#section-4.10
	MaxBytesSize uint32
	// MaxVecLen configures the upper bound limit for variable length arrays
	// https://tools.ietf.org/html/rfc4506#section-4.13
	MaxVecLen uint32
	// Source is the rand.Source which is used by the Generator to create
	// random values
	Source rand.Source
}

const (
	// DefaultMaxBytesSize is the MaxBytesSize value in the Generator returned by NewGenerator()
	DefaultMaxBytesSize = 1024
	// DefaultMaxVecLen is the MaxVecLen value in the Generator returned by NewGenerator()
	DefaultMaxVecLen = 10
	// DefaultSeed is the seed for the Source value in the Generator returned by NewGenerator()
	DefaultSeed = 99
)

// NewGenerator returns a new Generator instance configured with default settings.
// The returned Generator is deterministic but it is not thread-safe.
func NewGenerator() Generator {
	return Generator{
		MaxBytesSize: DefaultMaxBytesSize,
		MaxVecLen:    DefaultMaxVecLen,
		// rand.NewSource returns a source which is *not* safe for concurrent use
		Source: rand.NewSource(DefaultSeed),
	}
}

// Next modifies the given shape and populates it with random value fields.
func (g Generator) Next(shape goxdr.XdrType, presets []Preset) {
	marshaller := &randMarshaller{
		rand:         rand.New(g.Source),
		maxBytesSize: g.MaxBytesSize,
		maxVecLen:    g.MaxVecLen,
		presets:      presets,
	}
	marshaller.Marshal("", shape)
}
