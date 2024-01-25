package actions

import (
	"net/http"

	protocol "go/protocols/aurora"
	"go/services/aurora/internal/context"
	"go/services/aurora/internal/db2/history"
	"go/services/aurora/internal/resourceadapter"
	"go/support/render/problem"
)

// StreamableObjectResponse is an interface for objects returned by streamable object endpoints
// A streamable object endpoint is an SSE endpoint which returns a single JSON object response
// instead of a page of items.
type StreamableObjectResponse interface {
	Equals(other StreamableObjectResponse) bool
}

// OrderBookResponse is the response for the /order_book endpoint
// OrderBookResponse implements StreamableObjectResponse
type OrderBookResponse struct {
	protocol.OrderBookSummary
}

func priceLevelsEqual(a, b []protocol.PriceLevel) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// Equals returns true if the OrderBookResponse is equal to `other`
func (o OrderBookResponse) Equals(other StreamableObjectResponse) bool {
	otherOrderBook, ok := other.(OrderBookResponse)
	if !ok {
		return false
	}
	return otherOrderBook.Selling == o.Selling &&
		otherOrderBook.Buying == o.Buying &&
		priceLevelsEqual(otherOrderBook.Bids, o.Bids) &&
		priceLevelsEqual(otherOrderBook.Asks, o.Asks)
}

var invalidOrderBook = problem.P{
	Type:   "invalid_order_book",
	Title:  "Invalid Order Book Parameters",
	Status: http.StatusBadRequest,
	Detail: "The parameters that specify what order book to view are invalid in some way. " +
		"Please ensure that your type parameters (selling_asset_type and buying_asset_type) are one the " +
		"following valid values: native, credit_alphanum4, credit_alphanum12.  Also ensure that you " +
		"have specified selling_asset_code and selling_asset_issuer if selling_asset_type is not 'native', as well " +
		"as buying_asset_code and buying_asset_issuer if buying_asset_type is not 'native'",
}

// GetOrderbookHandler is the action handler for the /order_book endpoint
type GetOrderbookHandler struct {
}

func convertPriceLevels(src []history.PriceLevel) []protocol.PriceLevel {
	result := make([]protocol.PriceLevel, len(src))
	for i, l := range src {
		result[i] = protocol.PriceLevel{
			PriceR: protocol.Price{
				N: l.Pricen,
				D: l.Priced,
			},
			Price:  l.Pricef,
			Amount: l.Amount,
		}
	}

	return result
}

// GetResource implements the /order_book endpoint
func (handler GetOrderbookHandler) GetResource(w HeaderWriter, r *http.Request) (StreamableObjectResponse, error) {
	selling, err := getAsset(r, "selling_")
	if err != nil {
		return nil, invalidOrderBook
	}
	buying, err := getAsset(r, "buying_")
	if err != nil {
		return nil, invalidOrderBook
	}
	limit, err := getLimit(r, "limit", 20, 200)
	if err != nil {
		return nil, invalidOrderBook
	}

	historyQ, err := context.HistoryQFromRequest(r)
	if err != nil {
		return nil, err
	}

	summary, err := historyQ.GetOrderBookSummary(r.Context(), selling, buying, int(limit))
	if err != nil {
		return nil, err
	}

	var response OrderBookResponse
	if err := resourceadapter.PopulateAsset(r.Context(), &response.Selling, selling); err != nil {
		return nil, err
	}
	if err := resourceadapter.PopulateAsset(r.Context(), &response.Buying, buying); err != nil {
		return nil, err
	}
	response.Bids = convertPriceLevels(summary.Bids)
	response.Asks = convertPriceLevels(summary.Asks)

	return response, nil
}
