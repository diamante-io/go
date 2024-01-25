package txnbuild

import (
	"strconv"

	pricepkg "github.com/diamcircle/go/price"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/xdr"
)

type price struct {
	n int
	d int
	s string
}

func (p *price) parse(s string) error {
	if len(s) == 0 {
		return errors.New("cannot parse price from empty string")
	}

	xdrPrice, err := pricepkg.Parse(s)
	if err != nil {
		return errors.Wrap(err, "failed to parse price from string")
	}

	if len(p.s) > 0 {
		inverse, err := pricepkg.Parse(p.s)
		if err == nil && xdrPrice == inverse {
			return nil
		}
	}

	p.n = int(xdrPrice.N)
	p.d = int(xdrPrice.D)
	p.s = s
	return nil
}

func (p *price) fromXDR(xdrPrice xdr.Price) {
	n := int(xdrPrice.N)
	d := int(xdrPrice.D)
	if n == p.n && d == p.d {
		return
	}
	p.n = n
	p.d = d
	v := float64(n) / float64(d)
	// The special precision -1 uses the smallest number of digits
	// necessary such that ParseFloat will return f exactly.
	p.s = strconv.FormatFloat(v, 'f', -1, 32)
}

func priceFromXDR(xdrPrice xdr.Price) price {
	var p price
	p.fromXDR(xdrPrice)
	return p
}

func (p price) string() string {
	return p.s
}

func (p price) toXDR() xdr.Price {
	return xdr.Price{
		N: xdr.Int32(p.n),
		D: xdr.Int32(p.d),
	}
}
