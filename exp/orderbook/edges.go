package orderbook

import (
	"sort"

	"go/xdr"
)

// edgeSet maintains a mapping of assets to a set of venues, which
// is composed of a sorted lists of offers and, optionally, a liquidity pool.
// The offers are sorted by ascending price (in terms of the buying asset).
type edgeSet []edge

type edge struct {
	key   int32
	value Venues
}

func (e edgeSet) find(key int32) int {
	for i := 0; i < len(e); i++ {
		if e[i].key == key {
			return i
		}
	}
	return -1
}

// addOffer will insert the given offer into the edge set
func (e edgeSet) addOffer(key int32, offer xdr.OfferEntry) edgeSet {
	// The list of offers in a venue is sorted by cheapest to most expensive
	// price to convert buyingAsset to sellingAsset
	i := e.find(key)
	if i < 0 {
		return append(e, edge{key: key, value: Venues{offers: []xdr.OfferEntry{offer}}})
	}

	offers := e[i].value.offers
	// find the smallest i such that Price of offers[i] >  Price of offer
	insertIndex := sort.Search(len(offers), func(j int) bool {
		return offer.Price.Cheaper(offers[j].Price)
	})

	// then insert it into the slice (taken from Method 2 at
	// https://github.com/golang/go/wiki/SliceTricks#insert).
	offers = append(offers, xdr.OfferEntry{})          // add to end
	copy(offers[insertIndex+1:], offers[insertIndex:]) // shift right by 1
	offers[insertIndex] = offer                        // insert

	e[i].value = Venues{offers: offers, pool: e[i].value.pool}
	return e
}

// addPool makes `pool` a viable venue at `key`.
func (e edgeSet) addPool(key int32, pool liquidityPool) edgeSet {
	i := e.find(key)
	if i < 0 {
		return append(e, edge{key: key, value: Venues{pool: pool}})
	}
	e[i].value.pool = pool
	return e
}

// removeOffer will delete the given offer from the edge set, returning whether
// or not the given offer was actually found.
func (e edgeSet) removeOffer(key int32, offerID xdr.Int64) (edgeSet, bool) {
	i := e.find(key)
	if i < 0 {
		return e, false
	}

	offers := e[i].value.offers
	updatedOffers := offers
	contains := false
	for i, offer := range offers {
		if offer.OfferId != offerID {
			continue
		}

		// remove the entry in the slice at this location (taken from
		// https://github.com/golang/go/wiki/SliceTricks#cut).
		updatedOffers = append(offers[:i], offers[i+1:]...)
		contains = true
		break
	}

	if !contains {
		return e, false
	}

	if len(updatedOffers) == 0 && e[i].value.pool.Body.ConstantProduct == nil {
		return append(e[:i], e[i+1:]...), true
	}
	e[i].value.offers = updatedOffers
	return e, true
}

func (e edgeSet) removePool(key int32) edgeSet {
	i := e.find(key)
	if i < 0 {
		return e
	}

	if len(e[i].value.offers) == 0 {
		return append(e[:i], e[i+1:]...)
	}

	e[i].value = Venues{offers: e[i].value.offers}
	return e
}
