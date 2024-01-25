package actions

import (
	"database/sql"
	"math"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/services/aurora/internal/test"
	"github.com/stretchr/testify/assert"

	protocol "github.com/diamcircle/go/protocols/aurora"
)

type intObject int

func (i intObject) Equals(other StreamableObjectResponse) bool {
	return i == other.(intObject)
}

func TestOrderBookResponseEquals(t *testing.T) {
	for _, testCase := range []struct {
		name     string
		response protocol.OrderBookSummary
		other    StreamableObjectResponse
		expected bool
	}{
		{
			"empty orderbook summary",
			protocol.OrderBookSummary{},
			OrderBookResponse{},
			true,
		},
		{
			"types don't match",
			protocol.OrderBookSummary{},
			intObject(0),
			false,
		},
		{
			"buying asset doesn't match",
			protocol.OrderBookSummary{
				Buying: protocol.Asset{
					Type: "native",
				},
				Selling: protocol.Asset{
					Type: "native",
				},
			},
			OrderBookResponse{
				protocol.OrderBookSummary{
					Buying: protocol.Asset{
						Type:   "credit_alphanum4",
						Code:   "USD",
						Issuer: "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU",
					},
					Selling: protocol.Asset{
						Type: "native",
					},
				},
			},
			false,
		},
		{
			"selling asset doesn't match",
			protocol.OrderBookSummary{
				Selling: protocol.Asset{
					Type: "native",
				},
				Buying: protocol.Asset{
					Type: "native",
				},
			},
			OrderBookResponse{
				protocol.OrderBookSummary{
					Selling: protocol.Asset{
						Type:   "credit_alphanum4",
						Code:   "USD",
						Issuer: "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU",
					},
					Buying: protocol.Asset{
						Type: "native",
					},
				},
			},
			false,
		},
		{
			"bid lengths don't match",
			protocol.OrderBookSummary{
				Selling: protocol.Asset{
					Type:   "credit_alphanum4",
					Code:   "USD",
					Issuer: "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU",
				},
				Buying: protocol.Asset{
					Type: "native",
				},
				Bids: []protocol.PriceLevel{
					{
						PriceR: protocol.Price{N: 1, D: 2},
						Price:  "0.5",
						Amount: "123",
					},
					{
						PriceR: protocol.Price{N: 1, D: 1},
						Price:  "1.0",
						Amount: "123",
					},
				},
			},
			OrderBookResponse{
				protocol.OrderBookSummary{
					Selling: protocol.Asset{
						Type:   "credit_alphanum4",
						Code:   "USD",
						Issuer: "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU",
					},
					Buying: protocol.Asset{
						Type: "native",
					},
					Bids: []protocol.PriceLevel{
						{
							PriceR: protocol.Price{N: 1, D: 2},
							Price:  "0.5",
							Amount: "123",
						},
					},
				},
			},
			false,
		},
		{
			"ask lengths don't match",
			protocol.OrderBookSummary{
				Selling: protocol.Asset{
					Type:   "credit_alphanum4",
					Code:   "USD",
					Issuer: "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU",
				},
				Buying: protocol.Asset{
					Type: "native",
				},
				Asks: []protocol.PriceLevel{
					{
						PriceR: protocol.Price{N: 1, D: 2},
						Price:  "0.5",
						Amount: "123",
					},
					{
						PriceR: protocol.Price{N: 1, D: 1},
						Price:  "1.0",
						Amount: "123",
					},
				},
			},
			OrderBookResponse{
				protocol.OrderBookSummary{
					Selling: protocol.Asset{
						Type:   "credit_alphanum4",
						Code:   "USD",
						Issuer: "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU",
					},
					Buying: protocol.Asset{
						Type: "native",
					},
					Asks: []protocol.PriceLevel{
						{
							PriceR: protocol.Price{N: 1, D: 2},
							Price:  "0.5",
							Amount: "123",
						},
					},
				},
			},
			false,
		},
		{
			"bids don't match",
			protocol.OrderBookSummary{
				Selling: protocol.Asset{
					Type:   "credit_alphanum4",
					Code:   "USD",
					Issuer: "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU",
				},
				Buying: protocol.Asset{
					Type: "native",
				},
				Bids: []protocol.PriceLevel{
					{
						PriceR: protocol.Price{N: 1, D: 2},
						Price:  "0.5",
						Amount: "123",
					},
					{
						PriceR: protocol.Price{N: 1, D: 1},
						Price:  "1.0",
						Amount: "123",
					},
				},
			},
			OrderBookResponse{
				protocol.OrderBookSummary{
					Selling: protocol.Asset{
						Type:   "credit_alphanum4",
						Code:   "USD",
						Issuer: "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU",
					},
					Buying: protocol.Asset{
						Type: "native",
					},
					Bids: []protocol.PriceLevel{
						{
							PriceR: protocol.Price{N: 1, D: 2},
							Price:  "0.5",
							Amount: "123",
						},
						{
							PriceR: protocol.Price{N: 2, D: 1},
							Price:  "2.0",
							Amount: "123",
						},
					},
				},
			},
			false,
		},
		{
			"asks don't match",
			protocol.OrderBookSummary{
				Selling: protocol.Asset{
					Type:   "credit_alphanum4",
					Code:   "USD",
					Issuer: "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU",
				},
				Buying: protocol.Asset{
					Type: "native",
				},
				Asks: []protocol.PriceLevel{
					{
						PriceR: protocol.Price{N: 1, D: 2},
						Price:  "0.5",
						Amount: "123",
					},
					{
						PriceR: protocol.Price{N: 1, D: 1},
						Price:  "1.0",
						Amount: "123",
					},
				},
			},
			OrderBookResponse{
				protocol.OrderBookSummary{
					Selling: protocol.Asset{
						Type:   "credit_alphanum4",
						Code:   "USD",
						Issuer: "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU",
					},
					Buying: protocol.Asset{
						Type: "native",
					},
					Asks: []protocol.PriceLevel{
						{
							PriceR: protocol.Price{N: 1, D: 2},
							Price:  "0.5",
							Amount: "123",
						},
						{
							PriceR: protocol.Price{N: 1, D: 1},
							Price:  "1.0",
							Amount: "12",
						},
					},
				},
			},
			false,
		},
		{
			"orderbook summaries match",
			protocol.OrderBookSummary{
				Selling: protocol.Asset{
					Type:   "credit_alphanum4",
					Code:   "USD",
					Issuer: "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU",
				},
				Buying: protocol.Asset{
					Type: "native",
				},
				Asks: []protocol.PriceLevel{
					{
						PriceR: protocol.Price{N: 1, D: 2},
						Price:  "0.5",
						Amount: "123",
					},
					{
						PriceR: protocol.Price{N: 1, D: 1},
						Price:  "1.0",
						Amount: "123",
					},
				},
				Bids: []protocol.PriceLevel{
					{
						PriceR: protocol.Price{N: 1, D: 3},
						Price:  "0.33",
						Amount: "13",
					},
				},
			},
			OrderBookResponse{
				protocol.OrderBookSummary{
					Selling: protocol.Asset{
						Type:   "credit_alphanum4",
						Code:   "USD",
						Issuer: "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU",
					},
					Buying: protocol.Asset{
						Type: "native",
					},
					Bids: []protocol.PriceLevel{
						{
							PriceR: protocol.Price{N: 1, D: 3},
							Price:  "0.33",
							Amount: "13",
						},
					},
					Asks: []protocol.PriceLevel{
						{
							PriceR: protocol.Price{N: 1, D: 2},
							Price:  "0.5",
							Amount: "123",
						},
						{
							PriceR: protocol.Price{N: 1, D: 1},
							Price:  "1.0",
							Amount: "123",
						},
					},
				},
			},
			true,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			equals := (OrderBookResponse{testCase.response}).Equals(testCase.other)
			if equals != testCase.expected {
				t.Fatalf("expected %v but got %v", testCase.expected, equals)
			}
		})
	}
}

func TestOrderbookGetResourceValidation(t *testing.T) {
	handler := GetOrderbookHandler{}

	var eurAssetType, eurAssetCode, eurAssetIssuer string
	if err := eurAsset.Extract(&eurAssetType, &eurAssetCode, &eurAssetIssuer); err != nil {
		t.Fatalf("cound not extract eur asset: %v", err)
	}
	var usdAssetType, usdAssetCode, usdAssetIssuer string
	if err := eurAsset.Extract(&usdAssetType, &usdAssetCode, &usdAssetIssuer); err != nil {
		t.Fatalf("cound not extract usd asset: %v", err)
	}

	for _, testCase := range []struct {
		name        string
		queryParams map[string]string
	}{
		{
			"missing all params",
			map[string]string{},
		},
		{
			"missing buying asset",
			map[string]string{
				"selling_asset_type":   eurAssetType,
				"selling_asset_code":   eurAssetCode,
				"selling_asset_issuer": eurAssetIssuer,
				"limit":                "25",
			},
		},
		{
			"missing selling asset",
			map[string]string{
				"buying_asset_type":   eurAssetType,
				"buying_asset_code":   eurAssetCode,
				"buying_asset_issuer": eurAssetIssuer,
				"limit":               "25",
			},
		},
		{
			"invalid buying asset",
			map[string]string{
				"buying_asset_type":    "invalid",
				"buying_asset_code":    eurAssetCode,
				"buying_asset_issuer":  eurAssetIssuer,
				"selling_asset_type":   usdAssetType,
				"selling_asset_code":   usdAssetCode,
				"selling_asset_issuer": usdAssetIssuer,
				"limit":                "25",
			},
		},
		{
			"invalid selling asset",
			map[string]string{
				"buying_asset_type":    eurAssetType,
				"buying_asset_code":    eurAssetCode,
				"buying_asset_issuer":  eurAssetIssuer,
				"selling_asset_type":   "invalid",
				"selling_asset_code":   usdAssetCode,
				"selling_asset_issuer": usdAssetIssuer,
				"limit":                "25",
			},
		},
		{
			"limit is not a number",
			map[string]string{
				"buying_asset_type":    eurAssetType,
				"buying_asset_code":    eurAssetCode,
				"buying_asset_issuer":  eurAssetIssuer,
				"selling_asset_type":   usdAssetType,
				"selling_asset_code":   usdAssetCode,
				"selling_asset_issuer": usdAssetIssuer,
				"limit":                "avcdef",
			},
		},
		{
			"limit is negative",
			map[string]string{
				"buying_asset_type":    eurAssetType,
				"buying_asset_code":    eurAssetCode,
				"buying_asset_issuer":  eurAssetIssuer,
				"selling_asset_type":   usdAssetType,
				"selling_asset_code":   usdAssetCode,
				"selling_asset_issuer": usdAssetIssuer,
				"limit":                "-1",
			},
		},
		{
			"limit is too high",
			map[string]string{
				"buying_asset_type":    eurAssetType,
				"buying_asset_code":    eurAssetCode,
				"buying_asset_issuer":  eurAssetIssuer,
				"selling_asset_type":   usdAssetType,
				"selling_asset_code":   usdAssetCode,
				"selling_asset_issuer": usdAssetIssuer,
				"limit":                "20000",
			},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			r := makeRequest(t, testCase.queryParams, map[string]string{}, nil)
			w := httptest.NewRecorder()
			_, err := handler.GetResource(w, r)
			if err == nil || err.Error() != invalidOrderBook.Error() {
				t.Fatalf("expected error %v but got %v", invalidOrderBook, err)
			}
			if lastLedger := w.Header().Get(LastLedgerHeaderName); lastLedger != "" {
				t.Fatalf("expected last ledger to be not set but got %v", lastLedger)
			}
		})
	}
}

func TestOrderbookGetResource(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &history.Q{tt.AuroraSession()}

	var eurAssetType, eurAssetCode, eurAssetIssuer string
	if err := eurAsset.Extract(&eurAssetType, &eurAssetCode, &eurAssetIssuer); err != nil {
		t.Fatalf("cound not extract eur asset: %v", err)
	}

	empty := OrderBookResponse{
		OrderBookSummary: protocol.OrderBookSummary{
			Bids: []protocol.PriceLevel{},
			Asks: []protocol.PriceLevel{},
			Selling: protocol.Asset{
				Type: "native",
			},
			Buying: protocol.Asset{
				Type:   eurAssetType,
				Code:   eurAssetCode,
				Issuer: eurAssetIssuer,
			},
		},
	}

	sellEurOffer := history.Offer{
		SellerID: seller.Address(),
		OfferID:  int64(15),

		BuyingAsset:  nativeAsset,
		SellingAsset: eurAsset,

		Amount:             int64(500),
		Pricen:             int32(2),
		Priced:             int32(1),
		Price:              float64(2),
		Flags:              2,
		LastModifiedLedger: uint32(4),
	}

	otherEurOffer := history.Offer{
		SellerID: seller.Address(),
		OfferID:  int64(16),

		BuyingAsset:  eurAsset,
		SellingAsset: nativeAsset,

		Amount:             int64(math.MaxInt64),
		Pricen:             int32(2),
		Priced:             int32(1),
		Price:              float64(2),
		Flags:              2,
		LastModifiedLedger: uint32(1234),
	}

	nonCanonicalPriceTwoEurOffer := history.Offer{
		SellerID: seller.Address(),
		OfferID:  int64(7),

		BuyingAsset:  eurAsset,
		SellingAsset: nativeAsset,

		Amount:             int64(500),
		Pricen:             int32(2 * 15),
		Priced:             int32(1 * 15),
		Price:              float64(2),
		Flags:              2,
		LastModifiedLedger: uint32(4),
	}

	threeEurOffer := history.Offer{
		SellerID: seller.Address(),
		OfferID:  int64(20),

		BuyingAsset:  eurAsset,
		SellingAsset: nativeAsset,

		Amount:             int64(500),
		Pricen:             int32(3),
		Priced:             int32(1),
		Price:              float64(3),
		Flags:              2,
		LastModifiedLedger: uint32(1234),
	}

	otherSellEurOffer := history.Offer{
		SellerID: seller.Address(),
		OfferID:  int64(17),

		BuyingAsset:  nativeAsset,
		SellingAsset: eurAsset,

		Amount:             int64(500),
		Pricen:             int32(5),
		Priced:             int32(9),
		Price:              float64(5) / float64(9),
		Flags:              2,
		LastModifiedLedger: uint32(1234),
	}

	offers := []history.Offer{
		twoEurOffer,
		otherEurOffer,
		nonCanonicalPriceTwoEurOffer,
		threeEurOffer,
		sellEurOffer,
		otherSellEurOffer,
	}

	assert.NoError(t, q.TruncateTables(tt.Ctx, []string{"offers"}))
	assert.NoError(t, q.UpsertOffers(tt.Ctx, offers))

	assert.NoError(t, q.BeginTx(&sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	}))
	defer q.Rollback()

	fullResponse := empty
	fullResponse.Asks = []protocol.PriceLevel{
		{
			PriceR: protocol.Price{N: int32(twoEurOffer.Pricen), D: int32(twoEurOffer.Priced)},
			Price:  "2.0000000",
			Amount: "922337203685.4776807",
		},
		{
			PriceR: protocol.Price{N: int32(threeEurOffer.Pricen), D: int32(threeEurOffer.Priced)},
			Price:  "3.0000000",
			Amount: "0.0000500",
		},
	}
	fullResponse.Bids = []protocol.PriceLevel{
		{
			PriceR: protocol.Price{N: int32(otherSellEurOffer.Priced), D: int32(otherSellEurOffer.Pricen)},
			Price:  "1.8000000",
			Amount: "0.0000500",
		},
		{
			PriceR: protocol.Price{N: int32(sellEurOffer.Priced), D: int32(sellEurOffer.Pricen)},
			Price:  "0.5000000",
			Amount: "0.0000500",
		},
	}

	limitResponse := empty
	limitResponse.Asks = []protocol.PriceLevel{
		{
			PriceR: protocol.Price{N: int32(twoEurOffer.Pricen), D: int32(twoEurOffer.Priced)},
			Price:  "2.0000000",
			Amount: "922337203685.4776807",
		},
	}
	limitResponse.Bids = []protocol.PriceLevel{
		{
			PriceR: protocol.Price{N: int32(otherSellEurOffer.Priced), D: int32(otherSellEurOffer.Pricen)},
			Price:  "1.8000000",
			Amount: "0.0000500",
		},
	}

	for _, testCase := range []struct {
		name     string
		limit    int
		expected OrderBookResponse
	}{

		{
			"full orderbook",
			10,
			fullResponse,
		},
		{
			"limit request",
			1,
			limitResponse,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			handler := GetOrderbookHandler{}
			r := makeRequest(
				t,
				map[string]string{
					"buying_asset_type":   eurAssetType,
					"buying_asset_code":   eurAssetCode,
					"buying_asset_issuer": eurAssetIssuer,
					"selling_asset_type":  "native",
					"limit":               strconv.Itoa(testCase.limit),
				},
				map[string]string{},
				q,
			)
			w := httptest.NewRecorder()
			response, err := handler.GetResource(w, r)
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			if !response.Equals(testCase.expected) {
				t.Fatalf("expected %v but got %v", testCase.expected, response)
			}
		})
	}
}
