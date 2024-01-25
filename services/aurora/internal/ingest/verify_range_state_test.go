//lint:file-ignore U1001 Ignore all unused code, staticcheck doesn't understand testify/suite

package ingest

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"testing"

	"github.com/guregu/null"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/diamcircle/go/historyarchive"
	"github.com/diamcircle/go/ingest"
	"github.com/diamcircle/go/ingest/ledgerbackend"
	"github.com/diamcircle/go/keypair"
	"github.com/diamcircle/go/services/aurora/internal/db2"
	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/services/aurora/internal/ingest/processors"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/xdr"
)

func TestVerifyRangeStateTestSuite(t *testing.T) {
	suite.Run(t, new(VerifyRangeStateTestSuite))
}

type VerifyRangeStateTestSuite struct {
	suite.Suite
	ctx            context.Context
	ledgerBackend  *ledgerbackend.MockDatabaseBackend
	historyQ       *mockDBQ
	historyAdapter *mockHistoryArchiveAdapter
	runner         *mockProcessorsRunner
	system         *system
}

func (s *VerifyRangeStateTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.ledgerBackend = &ledgerbackend.MockDatabaseBackend{}
	s.historyQ = &mockDBQ{}
	s.historyAdapter = &mockHistoryArchiveAdapter{}
	s.runner = &mockProcessorsRunner{}
	s.system = &system{
		ctx:               s.ctx,
		historyQ:          s.historyQ,
		historyAdapter:    s.historyAdapter,
		ledgerBackend:     s.ledgerBackend,
		runner:            s.runner,
		checkpointManager: historyarchive.NewCheckpointManager(64),
	}
	s.system.initMetrics()

	s.historyQ.On("Rollback").Return(nil).Once()
}

func (s *VerifyRangeStateTestSuite) TearDownTest() {
	t := s.T()
	s.historyQ.AssertExpectations(t)
	s.historyAdapter.AssertExpectations(t)
	s.runner.AssertExpectations(t)
}

func (s *VerifyRangeStateTestSuite) TestInvalidRange() {
	// Recreate mock in this single test to remove Rollback assertion.
	*s.historyQ = mockDBQ{}

	next, err := verifyRangeState{fromLedger: 0, toLedger: 0}.run(s.system)
	s.Assert().Error(err)
	s.Assert().EqualError(err, "invalid range: [0, 0]")
	s.Assert().Equal(
		transition{node: stopState{}, sleepDuration: 0},
		next,
	)

	next, err = verifyRangeState{fromLedger: 0, toLedger: 100}.run(s.system)
	s.Assert().Error(err)
	s.Assert().EqualError(err, "invalid range: [0, 100]")
	s.Assert().Equal(
		transition{node: stopState{}, sleepDuration: 0},
		next,
	)

	next, err = verifyRangeState{fromLedger: 100, toLedger: 0}.run(s.system)
	s.Assert().Error(err)
	s.Assert().EqualError(err, "invalid range: [100, 0]")
	s.Assert().Equal(
		transition{node: stopState{}, sleepDuration: 0},
		next,
	)

	next, err = verifyRangeState{fromLedger: 100, toLedger: 99}.run(s.system)
	s.Assert().Error(err)
	s.Assert().EqualError(err, "invalid range: [100, 99]")
	s.Assert().Equal(
		transition{node: stopState{}, sleepDuration: 0},
		next,
	)
}

func (s *VerifyRangeStateTestSuite) TestBeginReturnsError() {
	// Recreate mock in this single test to remove Rollback assertion.
	*s.historyQ = mockDBQ{}
	s.historyQ.On("Begin").Return(errors.New("my error")).Once()

	next, err := verifyRangeState{fromLedger: 100, toLedger: 200}.run(s.system)
	s.Assert().Error(err)
	s.Assert().EqualError(err, "Error starting a transaction: my error")
	s.Assert().Equal(
		transition{node: stopState{}, sleepDuration: 0},
		next,
	)
}

func (s *VerifyRangeStateTestSuite) TestGetLastLedgerIngestReturnsError() {
	s.historyQ.On("Begin").Return(nil).Once()
	s.historyQ.On("GetLastLedgerIngest", s.ctx).Return(uint32(0), errors.New("my error")).Once()

	next, err := verifyRangeState{fromLedger: 100, toLedger: 200}.run(s.system)
	s.Assert().Error(err)
	s.Assert().EqualError(err, "Error getting last ingested ledger: my error")
	s.Assert().Equal(
		transition{node: stopState{}, sleepDuration: 0},
		next,
	)
}

func (s *VerifyRangeStateTestSuite) TestGetLastLedgerIngestNonEmpty() {
	s.historyQ.On("Begin").Return(nil).Once()
	s.historyQ.On("GetLastLedgerIngest", s.ctx).Return(uint32(100), nil).Once()

	next, err := verifyRangeState{fromLedger: 100, toLedger: 200}.run(s.system)
	s.Assert().Error(err)
	s.Assert().EqualError(err, "Database not empty")
	s.Assert().Equal(
		transition{node: stopState{}, sleepDuration: 0},
		next,
	)
}

func (s *VerifyRangeStateTestSuite) TestRunHistoryArchiveIngestionReturnsError() {
	s.historyQ.On("Begin").Return(nil).Once()
	s.historyQ.On("GetLastLedgerIngest", s.ctx).Return(uint32(0), nil).Once()
	s.ledgerBackend.On("PrepareRange", s.ctx, ledgerbackend.BoundedRange(100, 200)).Return(nil).Once()

	meta := xdr.LedgerCloseMeta{
		V0: &xdr.LedgerCloseMetaV0{
			LedgerHeader: xdr.LedgerHeaderHistoryEntry{
				Header: xdr.LedgerHeader{
					LedgerSeq:      xdr.Uint32(100),
					LedgerVersion:  xdr.Uint32(MaxSupportedProtocolVersion),
					BucketListHash: xdr.Hash{1, 2, 3},
				},
			},
		},
	}
	s.ledgerBackend.On("GetLedger", s.ctx, uint32(100)).Return(meta, nil).Once()
	s.runner.On("RunHistoryArchiveIngestion", uint32(100), MaxSupportedProtocolVersion, xdr.Hash{1, 2, 3}).Return(ingest.StatsChangeProcessorResults{}, errors.New("my error")).Once()

	next, err := verifyRangeState{fromLedger: 100, toLedger: 200}.run(s.system)
	s.Assert().Error(err)
	s.Assert().EqualError(err, "Error ingesting history archive: my error")
	s.Assert().Equal(
		transition{node: stopState{}, sleepDuration: 0},
		next,
	)
}

func (s *VerifyRangeStateTestSuite) TestSuccess() {
	s.historyQ.On("Begin").Return(nil).Once()
	s.historyQ.On("GetLastLedgerIngest", s.ctx).Return(uint32(0), nil).Once()
	s.ledgerBackend.On("PrepareRange", s.ctx, ledgerbackend.BoundedRange(100, 200)).Return(nil).Once()

	meta := xdr.LedgerCloseMeta{
		V0: &xdr.LedgerCloseMetaV0{
			LedgerHeader: xdr.LedgerHeaderHistoryEntry{
				Header: xdr.LedgerHeader{
					LedgerSeq:      xdr.Uint32(100),
					LedgerVersion:  xdr.Uint32(MaxSupportedProtocolVersion),
					BucketListHash: xdr.Hash{1, 2, 3},
				},
			},
		},
	}
	s.ledgerBackend.On("GetLedger", s.ctx, uint32(100)).Return(meta, nil).Once()
	s.runner.On("RunHistoryArchiveIngestion", uint32(100), MaxSupportedProtocolVersion, xdr.Hash{1, 2, 3}).Return(ingest.StatsChangeProcessorResults{}, nil).Once()

	s.historyQ.On("UpdateLastLedgerIngest", s.ctx, uint32(100)).Return(nil).Once()
	s.historyQ.On("Commit").Return(nil).Once()

	for i := uint32(101); i <= 200; i++ {
		s.historyQ.On("Begin").Return(nil).Once()

		meta := xdr.LedgerCloseMeta{
			V0: &xdr.LedgerCloseMetaV0{
				LedgerHeader: xdr.LedgerHeaderHistoryEntry{
					Header: xdr.LedgerHeader{
						LedgerSeq: xdr.Uint32(i),
					},
				},
			},
		}
		s.ledgerBackend.On("GetLedger", s.ctx, uint32(i)).Return(meta, nil).Once()

		s.runner.On("RunAllProcessorsOnLedger", meta).Return(
			ledgerStats{},
			nil,
		).Once()
		s.historyQ.On("UpdateLastLedgerIngest", s.ctx, i).Return(nil).Once()
		s.historyQ.On("Commit").Return(nil).Once()
	}

	s.historyQ.On("RebuildTradeAggregationBuckets", s.ctx, uint32(100), uint32(200)).Return(nil).Once()

	next, err := verifyRangeState{fromLedger: 100, toLedger: 200}.run(s.system)
	s.Assert().NoError(err)
	s.Assert().Equal(
		transition{node: stopState{}, sleepDuration: 0},
		next,
	)
}

// Bartek: looks like this test really tests the state verifier. Instead, I think we should just ensure
// data is passed so a single account would be enough to test if the FSM state works correctly.
func (s *VerifyRangeStateTestSuite) TestSuccessWithVerify() {
	s.historyQ.On("Begin").Return(nil).Once()
	s.historyQ.On("GetLastLedgerIngest", s.ctx).Return(uint32(0), nil).Once()
	s.ledgerBackend.On("PrepareRange", s.ctx, ledgerbackend.BoundedRange(100, 110)).Return(nil).Once()

	meta := xdr.LedgerCloseMeta{
		V0: &xdr.LedgerCloseMetaV0{
			LedgerHeader: xdr.LedgerHeaderHistoryEntry{
				Header: xdr.LedgerHeader{
					LedgerSeq:      xdr.Uint32(100),
					LedgerVersion:  xdr.Uint32(MaxSupportedProtocolVersion),
					BucketListHash: xdr.Hash{1, 2, 3},
				},
			},
		},
	}
	s.ledgerBackend.On("GetLedger", s.ctx, uint32(100)).Return(meta, nil).Once()
	s.runner.On("RunHistoryArchiveIngestion", uint32(100), MaxSupportedProtocolVersion, xdr.Hash{1, 2, 3}).Return(ingest.StatsChangeProcessorResults{}, nil).Once()

	s.historyQ.On("UpdateLastLedgerIngest", s.ctx, uint32(100)).Return(nil).Once()
	s.historyQ.On("Commit").Return(nil).Once()

	for i := uint32(101); i <= 110; i++ {
		s.historyQ.On("Begin").Return(nil).Once()

		meta := xdr.LedgerCloseMeta{
			V0: &xdr.LedgerCloseMetaV0{
				LedgerHeader: xdr.LedgerHeaderHistoryEntry{
					Header: xdr.LedgerHeader{
						LedgerSeq: xdr.Uint32(i),
					},
				},
			},
		}
		s.ledgerBackend.On("GetLedger", s.ctx, uint32(i)).Return(meta, nil).Once()

		s.runner.On("RunAllProcessorsOnLedger", meta).Return(
			ledgerStats{},
			nil,
		).Once()
		s.historyQ.On("UpdateLastLedgerIngest", s.ctx, i).Return(nil).Once()
		s.historyQ.On("Commit").Return(nil).Once()
	}

	s.historyQ.On("RebuildTradeAggregationBuckets", s.ctx, uint32(100), uint32(110)).Return(nil).Once()

	clonedQ := &mockDBQ{}
	s.historyQ.On("CloneIngestionQ").Return(clonedQ).Once()

	clonedQ.On("BeginTx", mock.AnythingOfType("*sql.TxOptions")).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*sql.TxOptions)
		s.Assert().Equal(sql.LevelRepeatableRead, arg.Isolation)
		s.Assert().True(arg.ReadOnly)
	}).Return(nil).Once()
	clonedQ.On("Rollback").Return(nil).Once()
	clonedQ.On("GetLastLedgerIngestNonBlocking", s.ctx).Return(uint32(63), nil).Once()
	mockChangeReader := &ingest.MockChangeReader{}
	mockChangeReader.On("Close").Return(nil).Once()
	mockAccountID := "GACMZD5VJXTRLKVET72CETCYKELPNCOTTBDC6DHFEUPLG5DHEK534JQX"
	sponsor := "GAREDQSXC7QZYJLVMTU7XZW4LSILQ4M5U4GNLO523LEWZ3JBRC5E4HLE"
	signers := []string{
		"GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML",
		"GA25GQLHJU3LPEJXEIAXK23AWEA5GWDUGRSHTQHDFT6HXHVMRULSQJUJ",
		"GC6G3EQFKOKIIZFTJQSCHTSXBVC4XO3I64F5IBRQNS3E5SW3MO3KWGMT",
	}
	accountChange := ingest.Change{
		Type: xdr.LedgerEntryTypeAccount,
		Pre:  nil,
		Post: &xdr.LedgerEntry{
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId:  xdr.MustAddress(mockAccountID),
					Balance:    xdr.Int64(600),
					Thresholds: [4]byte{1, 0, 0, 0},
					Signers: []xdr.Signer{
						{
							Key:    xdr.MustSigner(signers[0]),
							Weight: 1,
						},
						{
							Key:    xdr.MustSigner(signers[1]),
							Weight: 2,
						},
						{
							Key:    xdr.MustSigner(signers[2]),
							Weight: 3,
						},
					},
					Ext: xdr.AccountEntryExt{
						V: 1,
						V1: &xdr.AccountEntryExtensionV1{
							Liabilities: xdr.Liabilities{
								Buying:  1,
								Selling: 1,
							},
							Ext: xdr.AccountEntryExtensionV1Ext{
								V: 2,
								V2: &xdr.AccountEntryExtensionV2{
									NumSponsored:  xdr.Uint32(0),
									NumSponsoring: xdr.Uint32(2),
									SignerSponsoringIDs: []xdr.SponsorshipDescriptor{
										nil,
										xdr.MustAddressPtr(mockAccountID),
										xdr.MustAddressPtr(sponsor),
									},
								},
							},
						},
					},
				},
			},
			LastModifiedLedgerSeq: xdr.Uint32(62),
			Ext: xdr.LedgerEntryExt{
				V: 1,
				V1: &xdr.LedgerEntryExtensionV1{
					SponsoringId: xdr.MustAddressPtr("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
				},
			},
		},
	}
	offerChange := ingest.Change{
		Type: xdr.LedgerEntryTypeOffer,
		Pre:  nil,
		Post: &xdr.LedgerEntry{
			Data: xdr.LedgerEntryData{
				Type:  xdr.LedgerEntryTypeOffer,
				Offer: &eurOffer,
			},
			LastModifiedLedgerSeq: xdr.Uint32(62),
			Ext: xdr.LedgerEntryExt{
				V: 1,
				V1: &xdr.LedgerEntryExtensionV1{
					SponsoringId: xdr.MustAddressPtr("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
				},
			},
		},
	}
	balanceID := xdr.ClaimableBalanceId{
		Type: xdr.ClaimableBalanceIdTypeClaimableBalanceIdTypeV0,
		V0:   &xdr.Hash{1, 2, 3},
	}
	balanceIDStr, err := xdr.MarshalHex(balanceID)
	s.Assert().NoError(err)
	claimableBalance := history.ClaimableBalance{
		BalanceID:          balanceIDStr,
		Asset:              xdr.MustNewNativeAsset(),
		Amount:             10,
		LastModifiedLedger: 62,
		Claimants: []history.Claimant{
			{
				Destination: mockAccountID,
				Predicate: xdr.ClaimPredicate{
					Type: xdr.ClaimPredicateTypeClaimPredicateUnconditional,
				},
			},
		},
		Sponsor: null.StringFrom("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
		Flags:   uint32(xdr.ClaimableBalanceFlagsClaimableBalanceClawbackEnabledFlag),
	}
	claimableBalanceChange := ingest.Change{
		Type: xdr.LedgerEntryTypeClaimableBalance,
		Pre:  nil,
		Post: &xdr.LedgerEntry{
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeClaimableBalance,
				ClaimableBalance: &xdr.ClaimableBalanceEntry{
					BalanceId: balanceID,
					Claimants: []xdr.Claimant{
						{
							Type: xdr.ClaimantTypeClaimantTypeV0,
							V0: &xdr.ClaimantV0{
								Destination: xdr.MustAddress(claimableBalance.Claimants[0].Destination),
								Predicate:   claimableBalance.Claimants[0].Predicate,
							},
						},
					},
					Asset:  claimableBalance.Asset,
					Amount: claimableBalance.Amount,
					Ext: xdr.ClaimableBalanceEntryExt{
						V: 1,
						V1: &xdr.ClaimableBalanceEntryExtensionV1{
							Flags: xdr.Uint32(xdr.ClaimableBalanceFlagsClaimableBalanceClawbackEnabledFlag),
						},
					},
				},
			},
			LastModifiedLedgerSeq: xdr.Uint32(62),
			Ext: xdr.LedgerEntryExt{
				V: 1,
				V1: &xdr.LedgerEntryExtensionV1{
					SponsoringId: xdr.MustAddressPtr("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
				},
			},
		},
	}
	liquidityPool := history.LiquidityPool{
		PoolID:         "cafebabedeadbeef000000000000000000000000000000000000000000000000",
		Type:           xdr.LiquidityPoolTypeLiquidityPoolConstantProduct,
		Fee:            34,
		TrustlineCount: 52115,
		ShareCount:     412241,
		AssetReserves: []history.LiquidityPoolAssetReserve{
			{
				Asset:   xdr.MustNewCreditAsset("USD", "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
				Reserve: 450,
			},
			{
				Asset:   xdr.MustNewNativeAsset(),
				Reserve: 123,
			},
		},
		LastModifiedLedger: 62,
	}
	liquidityPoolChange := ingest.Change{
		Type: xdr.LedgerEntryTypeLiquidityPool,
		Pre:  nil,
		Post: &xdr.LedgerEntry{
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeLiquidityPool,
				LiquidityPool: &xdr.LiquidityPoolEntry{
					LiquidityPoolId: xdr.PoolId{0xca, 0xfe, 0xba, 0xbe, 0xde, 0xad, 0xbe, 0xef},
					Body: xdr.LiquidityPoolEntryBody{
						Type: xdr.LiquidityPoolTypeLiquidityPoolConstantProduct,
						ConstantProduct: &xdr.LiquidityPoolEntryConstantProduct{
							Params: xdr.LiquidityPoolConstantProductParameters{
								AssetA: xdr.MustNewCreditAsset("USD", "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
								AssetB: xdr.MustNewNativeAsset(),
								Fee:    34,
							},
							ReserveA:                 450,
							ReserveB:                 123,
							TotalPoolShares:          412241,
							PoolSharesTrustLineCount: 52115,
						},
					},
				},
			},
			LastModifiedLedgerSeq: xdr.Uint32(62),
		},
	}

	mockChangeReader.On("Read").Return(accountChange, nil).Once()
	mockChangeReader.On("Read").Return(offerChange, nil).Once()
	mockChangeReader.On("Read").Return(claimableBalanceChange, nil).Once()
	mockChangeReader.On("Read").Return(liquidityPoolChange, nil).Once()
	mockChangeReader.On("Read").Return(ingest.Change{}, io.EOF).Once()
	mockChangeReader.On("Read").Return(ingest.Change{}, io.EOF).Once()
	s.historyAdapter.On("GetState", s.ctx, uint32(63)).Return(mockChangeReader, nil).Once()
	mockAccount := history.AccountEntry{
		AccountID:          mockAccountID,
		Balance:            600,
		LastModifiedLedger: 62,
		MasterWeight:       1,
		NumSponsored:       0,
		NumSponsoring:      2,
		BuyingLiabilities:  1,
		SellingLiabilities: 1,
		Sponsor:            null.StringFrom("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
	}

	clonedQ.MockQAccounts.On("GetAccountsByIDs", s.ctx, []string{mockAccountID}).Return([]history.AccountEntry{mockAccount}, nil).Once()
	clonedQ.MockQSigners.On("SignersForAccounts", s.ctx, []string{mockAccountID}).Return([]history.AccountSigner{
		{
			Account: mockAccountID,
			Signer:  mockAccountID,
			Weight:  1,
		},
		{
			Account: mockAccountID,
			Signer:  signers[0],
			Weight:  1,
		},
		{
			Account: mockAccountID,
			Signer:  signers[2],
			Weight:  3,
			Sponsor: null.StringFrom(sponsor),
		},
		{
			Account: mockAccountID,
			Signer:  signers[1],
			Weight:  2,
			Sponsor: null.StringFrom(mockAccountID),
		},
	}, nil).Once()
	clonedQ.MockQSigners.On("CountAccounts", s.ctx).Return(1, nil).Once()
	mockOffer := history.Offer{
		SellerID:           eurOffer.SellerId.Address(),
		OfferID:            int64(eurOffer.OfferId),
		SellingAsset:       eurOffer.Selling,
		BuyingAsset:        eurOffer.Buying,
		Amount:             int64(eurOffer.Amount),
		Pricen:             int32(eurOffer.Price.N),
		Priced:             int32(eurOffer.Price.D),
		Price:              float64(eurOffer.Price.N) / float64(eurOffer.Price.N),
		Flags:              int32(eurOffer.Flags),
		LastModifiedLedger: 62,
		Sponsor:            null.StringFrom("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
	}
	clonedQ.MockQOffers.On("GetOffersByIDs", s.ctx, []int64{int64(eurOffer.OfferId)}).Return([]history.Offer{mockOffer}, nil).Once()
	clonedQ.MockQOffers.On("CountOffers", s.ctx).Return(1, nil).Once()
	// TODO: add accounts data, trustlines and asset stats
	clonedQ.MockQData.On("CountAccountsData", s.ctx).Return(0, nil).Once()
	clonedQ.MockQAssetStats.On("CountTrustLines", s.ctx).Return(0, nil).Once()
	clonedQ.MockQAssetStats.On("GetAssetStats", s.ctx, "", "", db2.PageQuery{
		Order: "asc",
		Limit: assetStatsBatchSize,
	}).Return([]history.ExpAssetStat{
		// Created by liquidity pool:
		{
			AssetType:   xdr.AssetTypeAssetTypeCreditAlphanum4,
			AssetCode:   "USD",
			AssetIssuer: "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML",
			Accounts: history.ExpAssetStatAccounts{
				LiquidityPools: 1,
			},
			Balances: history.ExpAssetStatBalances{
				Authorized:                      "0",
				AuthorizedToMaintainLiabilities: "0",
				ClaimableBalances:               "0",
				LiquidityPools:                  "450",
				Unauthorized:                    "0",
			},
			Amount: "0",
		}}, nil).Once()
	clonedQ.MockQAssetStats.On("GetAssetStats", s.ctx, "", "", db2.PageQuery{
		Cursor: "USD_GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML_credit_alphanum4",
		Order:  "asc",
		Limit:  assetStatsBatchSize,
	}).Return([]history.ExpAssetStat{}, nil).Once()

	clonedQ.MockQClaimableBalances.On("CountClaimableBalances", s.ctx).Return(1, nil).Once()
	clonedQ.MockQClaimableBalances.
		On("GetClaimableBalancesByID", s.ctx, []string{balanceIDStr}).
		Return([]history.ClaimableBalance{claimableBalance}, nil).Once()

	clonedQ.MockQLiquidityPools.On("CountLiquidityPools", s.ctx).Return(1, nil).Once()
	clonedQ.MockQLiquidityPools.
		On("GetLiquidityPoolsByID", s.ctx, []string{liquidityPool.PoolID}).
		Return([]history.LiquidityPool{liquidityPool}, nil).Once()

	next, err := verifyRangeState{
		fromLedger: 100, toLedger: 110, verifyState: true,
	}.run(s.system)
	s.Assert().NoError(err)
	s.Assert().Equal(
		transition{node: stopState{}, sleepDuration: 0},
		next,
	)
	clonedQ.AssertExpectations(s.T())
}

func (s *VerifyRangeStateTestSuite) TestVerifyFailsWhenAssetStatsMismatch() {
	set := processors.AssetStatSet{}

	trustLineIssuer := xdr.MustAddress("GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H")
	set.AddTrustline(
		ingest.Change{
			Post: &xdr.LedgerEntry{
				Data: xdr.LedgerEntryData{
					TrustLine: &xdr.TrustLineEntry{
						AccountId: xdr.MustAddress(keypair.MustRandom().Address()),
						Balance:   123,
						Asset:     xdr.MustNewCreditAsset("EUR", trustLineIssuer.Address()).ToTrustLineAsset(),
						Flags:     xdr.Uint32(xdr.TrustLineFlagsAuthorizedToMaintainLiabilitiesFlag),
					},
				},
			},
		},
	)

	stat := history.ExpAssetStat{
		AssetType:   xdr.AssetTypeAssetTypeCreditAlphanum4,
		AssetCode:   "EUR",
		AssetIssuer: trustLineIssuer.Address(),
		Accounts: history.ExpAssetStatAccounts{
			Unauthorized: 1,
		},
		Balances: history.ExpAssetStatBalances{
			Authorized:                      "0",
			AuthorizedToMaintainLiabilities: "0",
			Unauthorized:                    "123",
		},
		Amount:      "0",
		NumAccounts: 0,
	}

	s.historyQ.MockQAssetStats.On("GetAssetStats", s.ctx, "", "", db2.PageQuery{
		Order: "asc",
		Limit: assetStatsBatchSize,
	}).Return([]history.ExpAssetStat{stat}, nil).Once()
	s.historyQ.MockQAssetStats.On("GetAssetStats", s.ctx, "", "", db2.PageQuery{
		Cursor: stat.PagingToken(),
		Order:  "asc",
		Limit:  assetStatsBatchSize,
	}).Return([]history.ExpAssetStat{}, nil).Once()

	err := checkAssetStats(s.ctx, set, s.historyQ)
	s.Assert().Contains(err.Error(), fmt.Sprintf("db asset stat with code EUR issuer %s does not match asset stat from HAS", trustLineIssuer.Address()))

	// Satisfy the mock
	s.historyQ.Rollback()
}
