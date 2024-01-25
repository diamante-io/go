package xdr

import (
	"math/big"
)

// String returns a string represenation of `p`
func (p Price) String() string {
	return big.NewRat(int64(p.N), int64(p.D)).FloatString(7)
}

// Equal returns whether the price's value is the same,
// taking into account denormalized representation
// (e.g. Price{1, 2}.EqualValue(Price{2,4}) == true )
func (p Price) Equal(q Price) bool {
	// See the Cheaper() method for the reasoning behind this:
	return uint64(p.N)*uint64(q.D) == uint64(q.N)*uint64(p.D)
}

// Cheaper indicates if the Price's value is lower,
// taking into account denormalized representation
// (e.g. Price{1, 2}.Cheaper(Price{2,4}) == false )
func (p Price) Cheaper(q Price) bool {
	// To avoid float precision issues when naively comparing Price.N/Price.D,
	// we use the cross product instead:
	//
	// Price of p <  Price of q
	//  <==>
	// (p.N / p.D) < (q.N / q.D)
	//  <==>
	// (p.N / p.D) * (p.D * q.D) < (q.N / q.D) * (p.D * q.D)
	//  <==>
	// p.N * q.D < q.N * p.D
	return uint64(p.N)*uint64(q.D) < uint64(q.N)*uint64(p.D)
}

// Normalize sets Price to its rational canonical form
func (p *Price) Normalize() {
	r := big.NewRat(int64(p.N), int64(p.D))
	p.N = Int32(r.Num().Int64())
	p.D = Int32(r.Denom().Int64())
}

// Invert inverts Price.
func (p *Price) Invert() {
	p.N, p.D = p.D, p.N
}
