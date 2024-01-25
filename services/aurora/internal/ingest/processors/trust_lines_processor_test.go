//lint:file-ignore U1001 Ignore all unused code, staticcheck doesn't understand testify/suite

package processors

import (
	"context"
	"testing"

	"github.com/guregu/null"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/diamcircle/go/ingest"
	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/xdr"
)

var trustLineIssuer = xdr.MustAddress("GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H")

func TestTrustLinesProcessorTestSuiteState(t *testing.T) {
	suite.Run(t, new(TrustLinesProcessorTestSuiteState))
}

type TrustLinesProcessorTestSuiteState struct {
	suite.Suite
	ctx       context.Context
	processor *TrustLinesProcessor
	mockQ     *history.MockQTrustLines
}

func (s *TrustLinesProcessorTestSuiteState) SetupTest() {
	s.ctx = context.Background()
	s.mockQ = &history.MockQTrustLines{}
	s.processor = NewTrustLinesProcessor(s.mockQ)
}

func (s *TrustLinesProcessorTestSuiteState) TearDownTest() {
	s.Assert().NoError(s.processor.Commit(s.ctx))
	s.mockQ.AssertExpectations(s.T())
}

func (s *TrustLinesProcessorTestSuiteState) TestCreateTrustLine() {
	trustLine := xdr.TrustLineEntry{
		AccountId: xdr.MustAddress("GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"),
		Asset:     xdr.MustNewCreditAsset("EUR", trustLineIssuer.Address()).ToTrustLineAsset(),
		Flags:     xdr.Uint32(xdr.TrustLineFlagsAuthorizedFlag),
	}
	lastModifiedLedgerSeq := xdr.Uint32(123)

	poolShareTrustLine := xdr.TrustLineEntry{
		AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
		Asset: xdr.TrustLineAsset{
			Type:            xdr.AssetTypeAssetTypePoolShare,
			LiquidityPoolId: &xdr.PoolId{1, 2, 3, 4},
		},
		Balance: 12365,
		Limit:   123659,
		Flags:   xdr.Uint32(xdr.TrustLineFlagsAuthorizedFlag),
	}

	err := s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeTrustline,
		Pre:  nil,
		Post: &xdr.LedgerEntry{
			Data: xdr.LedgerEntryData{
				Type:      xdr.LedgerEntryTypeTrustline,
				TrustLine: &trustLine,
			},
			LastModifiedLedgerSeq: lastModifiedLedgerSeq,
		},
	})
	s.Assert().NoError(err)

	err = s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeTrustline,
		Pre:  nil,
		Post: &xdr.LedgerEntry{
			Data: xdr.LedgerEntryData{
				Type:      xdr.LedgerEntryTypeTrustline,
				TrustLine: &poolShareTrustLine,
			},
			LastModifiedLedgerSeq: lastModifiedLedgerSeq,
		},
	})
	s.Assert().NoError(err)

	s.mockQ.On(
		"UpsertTrustLines",
		s.ctx,
		mock.AnythingOfType("[]history.TrustLine"),
	).Run(func(args mock.Arguments) {
		arg := args.Get(1).([]history.TrustLine)
		s.Assert().ElementsMatch(
			[]history.TrustLine{
				{
					AccountID:          trustLine.AccountId.Address(),
					AssetType:          xdr.AssetTypeAssetTypeCreditAlphanum4,
					AssetIssuer:        trustLineIssuer.Address(),
					AssetCode:          "EUR",
					Balance:            int64(trustLine.Balance),
					LedgerKey:          "AAAAAQAAAAAdBJqAD9qPq+j2nRDdjdp5KVoUh8riPkNO9ato7BNs8wAAAAFFVVIAAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3",
					Limit:              int64(trustLine.Limit),
					LiquidityPoolID:    "",
					BuyingLiabilities:  int64(trustLine.Liabilities().Buying),
					SellingLiabilities: int64(trustLine.Liabilities().Selling),
					Flags:              uint32(trustLine.Flags),
					LastModifiedLedger: uint32(lastModifiedLedgerSeq),
					Sponsor:            null.String{},
				},
				{
					AccountID:          poolShareTrustLine.AccountId.Address(),
					AssetType:          xdr.AssetTypeAssetTypePoolShare,
					Balance:            int64(poolShareTrustLine.Balance),
					LedgerKey:          "AAAAAQAAAAC2LgFRDBZ3J52nLm30kq2iMgrO7dYzYAN3hvjtf1IHWgAAAAMBAgMEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==",
					Limit:              int64(poolShareTrustLine.Limit),
					LiquidityPoolID:    "0102030400000000000000000000000000000000000000000000000000000000",
					Flags:              uint32(poolShareTrustLine.Flags),
					LastModifiedLedger: uint32(lastModifiedLedgerSeq),
				},
			},
			arg)
	}).Return(nil).Once()
}

func (s *TrustLinesProcessorTestSuiteState) TestCreateTrustLineUnauthorized() {
	trustLine := xdr.TrustLineEntry{
		AccountId: xdr.MustAddress("GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"),
		Asset:     xdr.MustNewCreditAsset("EUR", trustLineIssuer.Address()).ToTrustLineAsset(),
	}
	lastModifiedLedgerSeq := xdr.Uint32(123)
	err := s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeTrustline,
		Pre:  nil,
		Post: &xdr.LedgerEntry{
			Data: xdr.LedgerEntryData{
				Type:      xdr.LedgerEntryTypeTrustline,
				TrustLine: &trustLine,
			},
			LastModifiedLedgerSeq: lastModifiedLedgerSeq,
		},
	})
	s.Assert().NoError(err)

	s.mockQ.On(
		"UpsertTrustLines",
		s.ctx,
		[]history.TrustLine{
			{
				AccountID:          trustLine.AccountId.Address(),
				AssetType:          xdr.AssetTypeAssetTypeCreditAlphanum4,
				AssetIssuer:        trustLineIssuer.Address(),
				AssetCode:          "EUR",
				Balance:            int64(trustLine.Balance),
				LedgerKey:          "AAAAAQAAAAAdBJqAD9qPq+j2nRDdjdp5KVoUh8riPkNO9ato7BNs8wAAAAFFVVIAAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3",
				Limit:              int64(trustLine.Limit),
				LiquidityPoolID:    "",
				BuyingLiabilities:  int64(trustLine.Liabilities().Buying),
				SellingLiabilities: int64(trustLine.Liabilities().Selling),
				Flags:              uint32(trustLine.Flags),
				LastModifiedLedger: uint32(lastModifiedLedgerSeq),
				Sponsor:            null.String{},
			},
		},
	).Return(nil).Once()
}

func TestTrustLinesProcessorTestSuiteLedger(t *testing.T) {
	suite.Run(t, new(TrustLinesProcessorTestSuiteLedger))
}

type TrustLinesProcessorTestSuiteLedger struct {
	suite.Suite
	ctx       context.Context
	processor *TrustLinesProcessor
	mockQ     *history.MockQTrustLines
}

func (s *TrustLinesProcessorTestSuiteLedger) SetupTest() {
	s.ctx = context.Background()
	s.mockQ = &history.MockQTrustLines{}
	s.processor = NewTrustLinesProcessor(s.mockQ)
}

func (s *TrustLinesProcessorTestSuiteLedger) TearDownTest() {
	s.mockQ.AssertExpectations(s.T())
}

func (s *TrustLinesProcessorTestSuiteLedger) TestNoIngestUpdateState() {
	// Nothing processed, assertions in TearDownTest.
	s.Assert().NoError(s.processor.Commit(s.ctx))
}

func (s *TrustLinesProcessorTestSuiteLedger) TestInsertTrustLine() {
	// should be ignored because it's not an trust line type
	err := s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeAccount,
		Pre:  nil,
		Post: &xdr.LedgerEntry{
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId:  xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					Thresholds: [4]byte{1, 1, 1, 1},
				},
			},
		},
	})
	s.Assert().NoError(err)

	// add trust line
	trustLine := xdr.TrustLineEntry{
		AccountId: xdr.MustAddress("GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"),
		Asset:     xdr.MustNewCreditAsset("EUR", trustLineIssuer.Address()).ToTrustLineAsset(),
		Balance:   0,
		Flags:     xdr.Uint32(xdr.TrustLineFlagsAuthorizedFlag),
	}
	unauthorizedTrustline := xdr.TrustLineEntry{
		AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
		Asset:     xdr.MustNewCreditAsset("USD", trustLineIssuer.Address()).ToTrustLineAsset(),
		Balance:   0,
	}
	lastModifiedLedgerSeq := xdr.Uint32(1234)

	err = s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeTrustline,
		Pre:  nil,
		Post: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: lastModifiedLedgerSeq,
			Data: xdr.LedgerEntryData{
				Type:      xdr.LedgerEntryTypeTrustline,
				TrustLine: &trustLine,
			},
		},
	})
	s.Assert().NoError(err)

	err = s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeTrustline,
		Pre:  nil,
		Post: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: lastModifiedLedgerSeq,
			Data: xdr.LedgerEntryData{
				Type:      xdr.LedgerEntryTypeTrustline,
				TrustLine: &unauthorizedTrustline,
			},
		},
	})
	s.Assert().NoError(err)

	updatedTrustLine := xdr.TrustLineEntry{
		AccountId: xdr.MustAddress("GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"),
		Asset:     xdr.MustNewCreditAsset("EUR", trustLineIssuer.Address()).ToTrustLineAsset(),
		Balance:   10,
		Flags:     xdr.Uint32(xdr.TrustLineFlagsAuthorizedFlag),
	}
	updatedUnauthorizedTrustline := xdr.TrustLineEntry{
		AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
		Asset:     xdr.MustNewCreditAsset("USD", trustLineIssuer.Address()).ToTrustLineAsset(),
		Balance:   10,
	}

	err = s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeTrustline,
		Pre: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: lastModifiedLedgerSeq,
			Data: xdr.LedgerEntryData{
				Type:      xdr.LedgerEntryTypeTrustline,
				TrustLine: &trustLine,
			},
		},
		Post: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: lastModifiedLedgerSeq,
			Data: xdr.LedgerEntryData{
				Type:      xdr.LedgerEntryTypeTrustline,
				TrustLine: &updatedTrustLine,
			},
		},
	})
	s.Assert().NoError(err)

	err = s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeTrustline,
		Pre: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: lastModifiedLedgerSeq,
			Data: xdr.LedgerEntryData{
				Type:      xdr.LedgerEntryTypeTrustline,
				TrustLine: &unauthorizedTrustline,
			},
		},
		Post: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: lastModifiedLedgerSeq,
			Data: xdr.LedgerEntryData{
				Type:      xdr.LedgerEntryTypeTrustline,
				TrustLine: &updatedUnauthorizedTrustline,
			},
		},
	})
	s.Assert().NoError(err)

	s.mockQ.On(
		"UpsertTrustLines",
		s.ctx,
		mock.AnythingOfType("[]history.TrustLine"),
	).Run(func(args mock.Arguments) {
		arg := args.Get(1).([]history.TrustLine)
		s.Assert().ElementsMatch(
			[]history.TrustLine{
				{
					AccountID:          updatedTrustLine.AccountId.Address(),
					AssetType:          updatedTrustLine.Asset.Type,
					AssetIssuer:        trustLineIssuer.Address(),
					AssetCode:          "EUR",
					Balance:            int64(updatedTrustLine.Balance),
					LedgerKey:          "AAAAAQAAAAAdBJqAD9qPq+j2nRDdjdp5KVoUh8riPkNO9ato7BNs8wAAAAFFVVIAAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3",
					Limit:              int64(updatedTrustLine.Limit),
					LiquidityPoolID:    "",
					BuyingLiabilities:  int64(updatedTrustLine.Liabilities().Buying),
					SellingLiabilities: int64(updatedTrustLine.Liabilities().Selling),
					Flags:              uint32(updatedTrustLine.Flags),
					LastModifiedLedger: uint32(lastModifiedLedgerSeq),
					Sponsor:            null.String{},
				},
				{
					AccountID:          updatedUnauthorizedTrustline.AccountId.Address(),
					AssetType:          updatedUnauthorizedTrustline.Asset.Type,
					AssetIssuer:        trustLineIssuer.Address(),
					AssetCode:          "USD",
					Balance:            int64(updatedUnauthorizedTrustline.Balance),
					LedgerKey:          "AAAAAQAAAAC2LgFRDBZ3J52nLm30kq2iMgrO7dYzYAN3hvjtf1IHWgAAAAFVU0QAAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3",
					Limit:              int64(updatedUnauthorizedTrustline.Limit),
					LiquidityPoolID:    "",
					BuyingLiabilities:  int64(updatedUnauthorizedTrustline.Liabilities().Buying),
					SellingLiabilities: int64(updatedUnauthorizedTrustline.Liabilities().Selling),
					Flags:              uint32(updatedUnauthorizedTrustline.Flags),
					LastModifiedLedger: uint32(lastModifiedLedgerSeq),
					Sponsor:            null.String{},
				},
			},
			arg,
		)
	}).Return(nil).Once()
	s.Assert().NoError(s.processor.Commit(s.ctx))
}

func (s *TrustLinesProcessorTestSuiteLedger) TestUpdateTrustLine() {
	lastModifiedLedgerSeq := xdr.Uint32(1234)

	trustLine := xdr.TrustLineEntry{
		AccountId: xdr.MustAddress("GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"),
		Asset:     xdr.MustNewCreditAsset("EUR", trustLineIssuer.Address()).ToTrustLineAsset(),
		Balance:   0,
		Flags:     xdr.Uint32(xdr.TrustLineFlagsAuthorizedFlag),
	}
	updatedTrustLine := xdr.TrustLineEntry{
		AccountId: xdr.MustAddress("GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"),
		Asset:     xdr.MustNewCreditAsset("EUR", trustLineIssuer.Address()).ToTrustLineAsset(),
		Balance:   10,
		Flags:     xdr.Uint32(xdr.TrustLineFlagsAuthorizedFlag),
	}

	err := s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeTrustline,
		Pre: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: lastModifiedLedgerSeq,
			Data: xdr.LedgerEntryData{
				Type:      xdr.LedgerEntryTypeTrustline,
				TrustLine: &trustLine,
			},
		},
		Post: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: lastModifiedLedgerSeq,
			Data: xdr.LedgerEntryData{
				Type:      xdr.LedgerEntryTypeTrustline,
				TrustLine: &updatedTrustLine,
			},
		},
	})
	s.Assert().NoError(err)

	s.mockQ.On(
		"UpsertTrustLines",
		s.ctx,
		[]history.TrustLine{
			{
				AccountID:          updatedTrustLine.AccountId.Address(),
				AssetType:          updatedTrustLine.Asset.Type,
				AssetIssuer:        trustLineIssuer.Address(),
				AssetCode:          "EUR",
				Balance:            int64(updatedTrustLine.Balance),
				LedgerKey:          "AAAAAQAAAAAdBJqAD9qPq+j2nRDdjdp5KVoUh8riPkNO9ato7BNs8wAAAAFFVVIAAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3",
				Limit:              int64(updatedTrustLine.Limit),
				LiquidityPoolID:    "",
				BuyingLiabilities:  int64(updatedTrustLine.Liabilities().Buying),
				SellingLiabilities: int64(updatedTrustLine.Liabilities().Selling),
				Flags:              uint32(updatedTrustLine.Flags),
				LastModifiedLedger: uint32(lastModifiedLedgerSeq),
				Sponsor:            null.String{},
			},
		},
	).Return(nil).Once()
	s.Assert().NoError(s.processor.Commit(s.ctx))
}

func (s *TrustLinesProcessorTestSuiteLedger) TestUpdateTrustLineAuthorization() {
	lastModifiedLedgerSeq := xdr.Uint32(1234)

	trustLine := xdr.TrustLineEntry{
		AccountId: xdr.MustAddress("GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"),
		Asset:     xdr.MustNewCreditAsset("EUR", trustLineIssuer.Address()).ToTrustLineAsset(),
		Balance:   100,
	}
	updatedTrustLine := xdr.TrustLineEntry{
		AccountId: xdr.MustAddress("GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"),
		Asset:     xdr.MustNewCreditAsset("EUR", trustLineIssuer.Address()).ToTrustLineAsset(),
		Balance:   10,
		Flags:     xdr.Uint32(xdr.TrustLineFlagsAuthorizedFlag),
	}

	otherTrustLine := xdr.TrustLineEntry{
		AccountId: xdr.MustAddress("GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"),
		Asset:     xdr.MustNewCreditAsset("USD", trustLineIssuer.Address()).ToTrustLineAsset(),
		Balance:   100,
		Flags:     xdr.Uint32(xdr.TrustLineFlagsAuthorizedFlag),
	}
	otherUpdatedTrustLine := xdr.TrustLineEntry{
		AccountId: xdr.MustAddress("GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"),
		Asset:     xdr.MustNewCreditAsset("USD", trustLineIssuer.Address()).ToTrustLineAsset(),
		Balance:   10,
	}

	err := s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeTrustline,
		Pre: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: lastModifiedLedgerSeq - 1,
			Data: xdr.LedgerEntryData{
				Type:      xdr.LedgerEntryTypeTrustline,
				TrustLine: &trustLine,
			},
		},
		Post: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: lastModifiedLedgerSeq,
			Data: xdr.LedgerEntryData{
				Type:      xdr.LedgerEntryTypeTrustline,
				TrustLine: &updatedTrustLine,
			},
		},
	})
	s.Assert().NoError(err)

	err = s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeTrustline,
		Pre: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: lastModifiedLedgerSeq - 1,
			Data: xdr.LedgerEntryData{
				Type:      xdr.LedgerEntryTypeTrustline,
				TrustLine: &otherTrustLine,
			},
		},
		Post: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: lastModifiedLedgerSeq,
			Data: xdr.LedgerEntryData{
				Type:      xdr.LedgerEntryTypeTrustline,
				TrustLine: &otherUpdatedTrustLine,
			},
		},
	})
	s.Assert().NoError(err)

	s.mockQ.On(
		"UpsertTrustLines",
		s.ctx,
		mock.AnythingOfType("[]history.TrustLine"),
	).Run(func(args mock.Arguments) {
		arg := args.Get(1).([]history.TrustLine)
		s.Assert().ElementsMatch(
			[]history.TrustLine{
				{
					AccountID:          updatedTrustLine.AccountId.Address(),
					AssetType:          updatedTrustLine.Asset.Type,
					AssetIssuer:        trustLineIssuer.Address(),
					AssetCode:          "EUR",
					Balance:            int64(updatedTrustLine.Balance),
					LedgerKey:          "AAAAAQAAAAAdBJqAD9qPq+j2nRDdjdp5KVoUh8riPkNO9ato7BNs8wAAAAFFVVIAAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3",
					Limit:              int64(updatedTrustLine.Limit),
					LiquidityPoolID:    "",
					BuyingLiabilities:  int64(updatedTrustLine.Liabilities().Buying),
					SellingLiabilities: int64(updatedTrustLine.Liabilities().Selling),
					Flags:              uint32(updatedTrustLine.Flags),
					LastModifiedLedger: uint32(lastModifiedLedgerSeq),
					Sponsor:            null.String{},
				},
				{
					AccountID:          otherUpdatedTrustLine.AccountId.Address(),
					AssetType:          otherUpdatedTrustLine.Asset.Type,
					AssetIssuer:        trustLineIssuer.Address(),
					AssetCode:          "USD",
					Balance:            int64(otherUpdatedTrustLine.Balance),
					LedgerKey:          "AAAAAQAAAAAdBJqAD9qPq+j2nRDdjdp5KVoUh8riPkNO9ato7BNs8wAAAAFVU0QAAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3",
					Limit:              int64(otherUpdatedTrustLine.Limit),
					LiquidityPoolID:    "",
					BuyingLiabilities:  int64(otherUpdatedTrustLine.Liabilities().Buying),
					SellingLiabilities: int64(otherUpdatedTrustLine.Liabilities().Selling),
					Flags:              uint32(otherUpdatedTrustLine.Flags),
					LastModifiedLedger: uint32(lastModifiedLedgerSeq),
					Sponsor:            null.String{},
				},
			},
			arg,
		)
	}).Return(nil).Once()
	s.Assert().NoError(s.processor.Commit(s.ctx))
}

func (s *TrustLinesProcessorTestSuiteLedger) TestRemoveTrustLine() {
	unauthorizedTrustLine := xdr.TrustLineEntry{
		AccountId: xdr.MustAddress("GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"),
		Asset:     xdr.MustNewCreditAsset("USD", trustLineIssuer.Address()).ToTrustLineAsset(),
		Balance:   0,
	}

	err := s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeTrustline,
		Pre: &xdr.LedgerEntry{
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeTrustline,
				TrustLine: &xdr.TrustLineEntry{
					AccountId: xdr.MustAddress("GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"),
					Asset:     xdr.MustNewCreditAsset("EUR", trustLineIssuer.Address()).ToTrustLineAsset(),
					Balance:   0,
					Flags:     xdr.Uint32(xdr.TrustLineFlagsAuthorizedFlag),
				},
			},
		},
		Post: nil,
	})
	s.Assert().NoError(err)

	err = s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeTrustline,
		Pre: &xdr.LedgerEntry{
			Data: xdr.LedgerEntryData{
				Type:      xdr.LedgerEntryTypeTrustline,
				TrustLine: &unauthorizedTrustLine,
			},
		},
		Post: nil,
	})
	s.Assert().NoError(err)

	lkStr1, err := xdr.LedgerKey{
		Type: xdr.LedgerEntryTypeTrustline,
		TrustLine: &xdr.LedgerKeyTrustLine{
			AccountId: xdr.MustAddress("GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"),
			Asset:     xdr.MustNewCreditAsset("EUR", trustLineIssuer.Address()).ToTrustLineAsset(),
		},
	}.MarshalBinaryBase64()

	lkStr2, err := xdr.LedgerKey{
		Type: xdr.LedgerEntryTypeTrustline,
		TrustLine: &xdr.LedgerKeyTrustLine{
			AccountId: xdr.MustAddress("GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"),
			Asset:     xdr.MustNewCreditAsset("USD", trustLineIssuer.Address()).ToTrustLineAsset(),
		},
	}.MarshalBinaryBase64()
	s.Assert().NoError(err)
	s.mockQ.On(
		"RemoveTrustLines", s.ctx, mock.Anything,
	).Run(func(args mock.Arguments) {
		// To fix order issue due to using ChangeCompactor
		ledgerKeys := args.Get(1).([]string)
		s.Assert().ElementsMatch(
			ledgerKeys,
			[]string{lkStr1, lkStr2},
		)
	}).Return(int64(2), nil).Once()

	s.Assert().NoError(s.processor.Commit(s.ctx))
}

func (s *TrustLinesProcessorTestSuiteLedger) TestProcessUpgradeChange() {
	// add trust line
	lastModifiedLedgerSeq := xdr.Uint32(1234)
	trustLine := xdr.TrustLineEntry{
		AccountId: xdr.MustAddress("GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"),
		Asset:     xdr.MustNewCreditAsset("EUR", trustLineIssuer.Address()).ToTrustLineAsset(),
		Balance:   0,
		Flags:     xdr.Uint32(xdr.TrustLineFlagsAuthorizedFlag),
	}

	err := s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeTrustline,
		Post: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: lastModifiedLedgerSeq,
			Data: xdr.LedgerEntryData{
				Type:      xdr.LedgerEntryTypeTrustline,
				TrustLine: &trustLine,
			},
		},
	})
	s.Assert().NoError(err)

	updatedTrustLine := xdr.TrustLineEntry{
		AccountId: xdr.MustAddress("GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"),
		Asset:     xdr.MustNewCreditAsset("EUR", trustLineIssuer.Address()).ToTrustLineAsset(),
		Balance:   10,
		Flags:     xdr.Uint32(xdr.TrustLineFlagsAuthorizedFlag),
	}

	err = s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeTrustline,
		Pre: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: lastModifiedLedgerSeq,
			Data: xdr.LedgerEntryData{
				Type:      xdr.LedgerEntryTypeTrustline,
				TrustLine: &trustLine,
			},
		},
		Post: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: lastModifiedLedgerSeq,
			Data: xdr.LedgerEntryData{
				Type:      xdr.LedgerEntryTypeTrustline,
				TrustLine: &updatedTrustLine,
			},
		},
	})
	s.Assert().NoError(err)

	s.mockQ.On(
		"UpsertTrustLines", s.ctx,
		[]history.TrustLine{
			{
				AccountID:          updatedTrustLine.AccountId.Address(),
				AssetType:          updatedTrustLine.Asset.Type,
				AssetIssuer:        trustLineIssuer.Address(),
				AssetCode:          "EUR",
				Balance:            int64(updatedTrustLine.Balance),
				LedgerKey:          "AAAAAQAAAAAdBJqAD9qPq+j2nRDdjdp5KVoUh8riPkNO9ato7BNs8wAAAAFFVVIAAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3",
				Limit:              int64(updatedTrustLine.Limit),
				LiquidityPoolID:    "",
				BuyingLiabilities:  int64(updatedTrustLine.Liabilities().Buying),
				SellingLiabilities: int64(updatedTrustLine.Liabilities().Selling),
				Flags:              uint32(updatedTrustLine.Flags),
				LastModifiedLedger: uint32(lastModifiedLedgerSeq),
				Sponsor:            null.String{},
			},
		},
	).Return(nil).Once()
	s.Assert().NoError(s.processor.Commit(s.ctx))
}

func (s *TrustLinesProcessorTestSuiteLedger) TestRemoveTrustlineNoRowsAffected() {
	err := s.processor.ProcessChange(s.ctx, ingest.Change{
		Type: xdr.LedgerEntryTypeTrustline,
		Pre: &xdr.LedgerEntry{
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeTrustline,
				TrustLine: &xdr.TrustLineEntry{
					AccountId: xdr.MustAddress("GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"),
					Asset:     xdr.MustNewCreditAsset("EUR", trustLineIssuer.Address()).ToTrustLineAsset(),
					Balance:   0,
					Flags:     xdr.Uint32(xdr.TrustLineFlagsAuthorizedFlag),
				},
			},
		},
		Post: nil,
	})
	s.Assert().NoError(err)

	lkStr, err := xdr.LedgerKey{
		Type: xdr.LedgerEntryTypeTrustline,
		TrustLine: &xdr.LedgerKeyTrustLine{
			AccountId: xdr.MustAddress("GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"),
			Asset:     xdr.MustNewCreditAsset("EUR", trustLineIssuer.Address()).ToTrustLineAsset(),
		},
	}.MarshalBinaryBase64()
	s.Assert().NoError(err)

	s.mockQ.On(
		"RemoveTrustLines", s.ctx, []string{lkStr},
	).Return(int64(0), nil).Once()

	err = s.processor.Commit(s.ctx)
	s.Assert().Error(err)
	s.Assert().IsType(ingest.StateError{}, errors.Cause(err))
	s.Assert().EqualError(err, "0 rows affected when removing 1 trust lines")
}
