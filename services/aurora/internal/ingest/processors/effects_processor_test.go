//lint:file-ignore U1001 Ignore all unused code, staticcheck doesn't understand testify/suite

package processors

import (
	"context"
	"encoding/hex"
	"testing"

	"go/protocols/aurora/base"

	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"go/ingest"
	"go/services/aurora/internal/db2/history"
	. "go/services/aurora/internal/test/transactions"
	"go/services/aurora/internal/toid"
	"go/support/errors"
	"go/xdr"
)

type EffectsProcessorTestSuiteLedger struct {
	suite.Suite
	ctx                    context.Context
	processor              *EffectProcessor
	mockQ                  *history.MockQEffects
	mockBatchInsertBuilder *history.MockEffectBatchInsertBuilder

	firstTx     ingest.LedgerTransaction
	secondTx    ingest.LedgerTransaction
	thirdTx     ingest.LedgerTransaction
	failedTx    ingest.LedgerTransaction
	firstTxID   int64
	secondTxID  int64
	thirdTxID   int64
	failedTxID  int64
	sequence    uint32
	addresses   []string
	addressToID map[string]int64
	txs         []ingest.LedgerTransaction
}

func TestEffectsProcessorTestSuiteLedger(t *testing.T) {
	suite.Run(t, new(EffectsProcessorTestSuiteLedger))
}

func (s *EffectsProcessorTestSuiteLedger) SetupTest() {
	s.ctx = context.Background()
	s.mockQ = &history.MockQEffects{}
	s.mockBatchInsertBuilder = &history.MockEffectBatchInsertBuilder{}

	s.sequence = uint32(20)

	s.addresses = []string{
		"GANFZDRBCNTUXIODCJEYMACPMCSZEVE4WZGZ3CZDZ3P2SXK4KH75IK6Y",
		"GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
		"GCQZP3IU7XU6EJ63JZXKCQOYT2RNXN3HB5CNHENNUEUHSMA4VUJJJSEN",
	}

	s.firstTx = BuildLedgerTransaction(
		s.Suite.T(),
		TestTransaction{
			Index:         1,
			EnvelopeXDR:   "AAAAAKGX7RT96eIn205uoUHYnqLbt2cPRNORraEoeTAcrRKUAAAAZAAAADkAAAABAAAAAAAAAAAAAAABAAAAAAAAAAsAAABF2WS4AAAAAAAAAAABHK0SlAAAAEDq0JVhKNIq9ag0sR+R/cv3d9tEuaYEm2BazIzILRdGj9alaVMZBhxoJ3ZIpP3rraCJzyoKZO+p5HBVe10a2+UG",
			ResultXDR:     "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAALAAAAAAAAAAA=",
			MetaXDR:       "AAAAAQAAAAIAAAADAAAAOgAAAAAAAAAAoZftFP3p4ifbTm6hQdieotu3Zw9E05GtoSh5MBytEpQAAAACVAvjnAAAADkAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAOgAAAAAAAAAAoZftFP3p4ifbTm6hQdieotu3Zw9E05GtoSh5MBytEpQAAAACVAvjnAAAADkAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAA6AAAAAAAAAAChl+0U/eniJ9tObqFB2J6i27dnD0TTka2hKHkwHK0SlAAAAAJUC+OcAAAAOQAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAA6AAAAAAAAAAChl+0U/eniJ9tObqFB2J6i27dnD0TTka2hKHkwHK0SlAAAAAJUC+OcAAAARdlkuAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			FeeChangesXDR: "AAAAAgAAAAMAAAA5AAAAAAAAAAChl+0U/eniJ9tObqFB2J6i27dnD0TTka2hKHkwHK0SlAAAAAJUC+QAAAAAOQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAA6AAAAAAAAAAChl+0U/eniJ9tObqFB2J6i27dnD0TTka2hKHkwHK0SlAAAAAJUC+OcAAAAOQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			Hash:          "829d53f2dceebe10af8007564b0aefde819b95734ad431df84270651e7ed8a90",
		},
	)
	s.firstTxID = toid.New(int32(s.sequence), 1, 0).ToInt64()

	s.secondTx = BuildLedgerTransaction(
		s.Suite.T(),
		TestTransaction{
			Index:         2,
			EnvelopeXDR:   "AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAAaAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAoZftFP3p4ifbTm6hQdieotu3Zw9E05GtoSh5MBytEpQAAAACVAvkAAAAAAAAAAABVvwF9wAAAEDHU95E9wxgETD8TqxUrkgC0/7XHyNDts6Q5huRHfDRyRcoHdv7aMp/sPvC3RPkXjOMjgbKJUX7SgExUeYB5f8F",
			ResultXDR:     "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=",
			MetaXDR:       "AAAAAQAAAAIAAAADAAAAOQAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcLGrZY9dZxbAAAAAAAAAAZAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAOQAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcLGrZY9dZxbAAAAAAAAAAaAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAA5AAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsatlj11nFsAAAAAAAAABoAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAA5AAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsatlahyo1sAAAAAAAAABoAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAA5AAAAAAAAAAChl+0U/eniJ9tObqFB2J6i27dnD0TTka2hKHkwHK0SlAAAAAJUC+QAAAAAOQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			FeeChangesXDR: "AAAAAgAAAAMAAAA3AAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsatlj11nHQAAAAAAAAABkAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAA5AAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsatlj11nFsAAAAAAAAABkAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			Hash:          "0e5bd332291e3098e49886df2cdb9b5369a5f9e0a9973f0d9e1a9489c6581ba2",
		},
	)

	s.secondTxID = toid.New(int32(s.sequence), 2, 0).ToInt64()

	s.thirdTx = BuildLedgerTransaction(
		s.Suite.T(),
		TestTransaction{
			Index:         3,
			EnvelopeXDR:   "AAAAABpcjiETZ0uhwxJJhgBPYKWSVJy2TZ2LI87fqV1cUf/UAAAAZAAAADcAAAABAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAGlyOIRNnS6HDEkmGAE9gpZJUnLZNnYsjzt+pXVxR/9QAAAAAAAAAAAX14QAAAAAAAAAAAVxR/9QAAABAK6pcXYMzAEmH08CZ1LWmvtNDKauhx+OImtP/Lk4hVTMJRVBOebVs5WEPj9iSrgGT0EswuDCZ2i5AEzwgGof9Ag==",
			ResultXDR:     "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=",
			MetaXDR:       "AAAAAQAAAAIAAAADAAAAOAAAAAAAAAAAGlyOIRNnS6HDEkmGAE9gpZJUnLZNnYsjzt+pXVxR/9QAAAACVAvjnAAAADcAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAOAAAAAAAAAAAGlyOIRNnS6HDEkmGAE9gpZJUnLZNnYsjzt+pXVxR/9QAAAACVAvjnAAAADcAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAA==",
			FeeChangesXDR: "AAAAAgAAAAMAAAA3AAAAAAAAAAAaXI4hE2dLocMSSYYAT2ClklSctk2diyPO36ldXFH/1AAAAAJUC+QAAAAANwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAA4AAAAAAAAAAAaXI4hE2dLocMSSYYAT2ClklSctk2diyPO36ldXFH/1AAAAAJUC+OcAAAANwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			Hash:          "2a805712c6d10f9e74bb0ccf54ae92a2b4b1e586451fe8133a2433816f6b567c",
		},
	)
	s.thirdTxID = toid.New(int32(s.sequence), 3, 0).ToInt64()

	s.failedTx = BuildLedgerTransaction(
		s.Suite.T(),
		TestTransaction{
			Index:         4,
			EnvelopeXDR:   "AAAAAPCq/iehD2ASJorqlTyEt0usn2WG3yF4w9xBkgd4itu6AAAAZAAMpboAADNGAAAAAAAAAAAAAAABAAAAAAAAAAMAAAABVEVTVAAAAAAObS6P1g8rj8sCVzRQzYgHhWFkbh1oV+1s47LFPstSpQAAAAAAAAACVAvkAAAAAfcAAAD6AAAAAAAAAAAAAAAAAAAAAXiK27oAAABAHHk5mvM6xBRsvu3RBvzzPIb8GpXaL2M7InPn65LIhFJ2RnHIYrpP6ufZc6SUtKqChNRaN4qw5rjwFXNezmrBCw==",
			ResultXDR:     "AAAAAAAAAGT/////AAAAAQAAAAAAAAAD////+QAAAAA=",
			MetaXDR:       "AAAAAQAAAAIAAAADABDLGAAAAAAAAAAA8Kr+J6EPYBImiuqVPIS3S6yfZYbfIXjD3EGSB3iK27oAAAB2ucIg2AAMpboAADNFAAAA4wAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAQAABHT9ws4fAAAAAAAAAAAAAAAAAAAAAAAAAAEAEMsYAAAAAAAAAADwqv4noQ9gEiaK6pU8hLdLrJ9lht8heMPcQZIHeIrbugAAAHa5wiDYAAylugAAM0YAAADjAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAABAAAEdP3Czh8AAAAAAAAAAAAAAAAAAAAAAAAAAA==",
			FeeChangesXDR: "AAAAAgAAAAMAEMsCAAAAAAAAAADwqv4noQ9gEiaK6pU8hLdLrJ9lht8heMPcQZIHeIrbugAAAHa5wiE8AAylugAAM0UAAADjAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAABAAAEdP3Czh8AAAAAAAAAAAAAAAAAAAAAAAAAAQAQyxgAAAAAAAAAAPCq/iehD2ASJorqlTyEt0usn2WG3yF4w9xBkgd4itu6AAAAdrnCINgADKW6AAAzRQAAAOMAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAEAAAR0/cLOHwAAAAAAAAAAAAAAAAAAAAA=",
			Hash:          "24206737a02f7f855c46e367418e38c223f897792c76bbfb948e1b0dbd695f8b",
		},
	)
	s.failedTxID = toid.New(int32(s.sequence), 4, 0).ToInt64()

	s.addressToID = map[string]int64{
		s.addresses[0]: 2,
		s.addresses[1]: 20,
		s.addresses[2]: 200,
	}

	s.processor = NewEffectProcessor(
		s.mockQ,
		20,
	)

	s.txs = []ingest.LedgerTransaction{
		s.firstTx,
		s.secondTx,
		s.thirdTx,
	}
}

func (s *EffectsProcessorTestSuiteLedger) TearDownTest() {
	s.mockQ.AssertExpectations(s.T())
}

func (s *EffectsProcessorTestSuiteLedger) mockSuccessfulEffectBatchAdds() {
	s.mockBatchInsertBuilder.On(
		"Add",
		s.ctx,
		s.addressToID[s.addresses[2]],
		null.String{},
		toid.New(int32(s.sequence), 1, 1).ToInt64(),
		uint32(1),
		history.EffectSequenceBumped,
		[]byte("{\"new_seq\":300000000000}"),
	).Return(nil).Once()
	s.mockBatchInsertBuilder.On(
		"Add",
		s.ctx,
		s.addressToID[s.addresses[2]],
		null.String{},
		toid.New(int32(s.sequence), 2, 1).ToInt64(),
		uint32(1),
		history.EffectAccountCreated,
		[]byte("{\"starting_balance\":\"1000.0000000\"}"),
	).Return(nil).Once()
	s.mockBatchInsertBuilder.On(
		"Add",
		s.ctx,
		s.addressToID[s.addresses[1]],
		null.String{},
		toid.New(int32(s.sequence), 2, 1).ToInt64(),
		uint32(2),
		history.EffectAccountDebited,
		[]byte("{\"amount\":\"1000.0000000\",\"asset_type\":\"native\"}"),
	).Return(nil).Once()
	s.mockBatchInsertBuilder.On(
		"Add",
		s.ctx,
		s.addressToID[s.addresses[2]],
		null.String{},
		toid.New(int32(s.sequence), 2, 1).ToInt64(),
		uint32(3),
		history.EffectSignerCreated,
		[]byte("{\"public_key\":\"GCQZP3IU7XU6EJ63JZXKCQOYT2RNXN3HB5CNHENNUEUHSMA4VUJJJSEN\",\"weight\":1}"),
	).Return(nil).Once()

	s.mockBatchInsertBuilder.On(
		"Add",
		s.ctx,
		s.addressToID[s.addresses[0]],
		null.String{},
		toid.New(int32(s.sequence), 3, 1).ToInt64(),
		uint32(1),
		history.EffectAccountCredited,
		[]byte("{\"amount\":\"10.0000000\",\"asset_type\":\"native\"}"),
	).Return(nil).Once()

	s.mockBatchInsertBuilder.On(
		"Add",
		s.ctx,
		s.addressToID[s.addresses[0]],
		null.String{},
		toid.New(int32(s.sequence), 3, 1).ToInt64(),
		uint32(2),
		history.EffectAccountDebited,
		[]byte("{\"amount\":\"10.0000000\",\"asset_type\":\"native\"}"),
	).Return(nil).Once()
}

func (s *EffectsProcessorTestSuiteLedger) mockSuccessfulCreateAccounts() {
	s.mockQ.On(
		"CreateAccounts",
		s.ctx,
		mock.AnythingOfType("[]string"),
		maxBatchSize,
	).Run(func(args mock.Arguments) {
		arg := args.Get(1).([]string)
		s.Assert().ElementsMatch(s.addresses, arg)
	}).Return(s.addressToID, nil).Once()
}

func (s *EffectsProcessorTestSuiteLedger) TestEmptyEffects() {
	err := s.processor.Commit(context.Background())
	s.Assert().NoError(err)
}

func (s *EffectsProcessorTestSuiteLedger) TestIngestEffectsSucceeds() {
	s.mockSuccessfulCreateAccounts()
	s.mockQ.On("NewEffectBatchInsertBuilder", maxBatchSize).
		Return(s.mockBatchInsertBuilder).Once()

	s.mockSuccessfulEffectBatchAdds()

	s.mockBatchInsertBuilder.On("Exec", s.ctx).Return(nil).Once()

	for _, tx := range s.txs {
		err := s.processor.ProcessTransaction(s.ctx, tx)
		s.Assert().NoError(err)
	}
	err := s.processor.Commit(s.ctx)
	s.Assert().NoError(err)
}

func (s *EffectsProcessorTestSuiteLedger) TestCreateAccountsFails() {
	s.mockQ.On("CreateAccounts", s.ctx, mock.AnythingOfType("[]string"), maxBatchSize).
		Return(s.addressToID, errors.New("transient error")).Once()

	for _, tx := range s.txs {
		err := s.processor.ProcessTransaction(s.ctx, tx)
		s.Assert().NoError(err)
	}
	err := s.processor.Commit(s.ctx)
	s.Assert().EqualError(err, "Could not create account ids: transient error")
}

func (s *EffectsProcessorTestSuiteLedger) TestBatchAddFails() {
	s.mockSuccessfulCreateAccounts()
	s.mockQ.On("NewEffectBatchInsertBuilder", maxBatchSize).
		Return(s.mockBatchInsertBuilder).Once()

	s.mockBatchInsertBuilder.On(
		"Add", s.ctx,
		s.addressToID[s.addresses[2]],
		null.String{},
		toid.New(int32(s.sequence), 1, 1).ToInt64(),
		uint32(1),
		history.EffectSequenceBumped,
		[]byte("{\"new_seq\":300000000000}"),
	).Return(errors.New("transient error")).Once()
	for _, tx := range s.txs {
		err := s.processor.ProcessTransaction(s.ctx, tx)
		s.Assert().NoError(err)
	}
	err := s.processor.Commit(s.ctx)
	s.Assert().EqualError(err, "could not insert operation effect in db: transient error")
}

func getRevokeSponsorshipMeta(t *testing.T) (string, []effect) {
	source := xdr.MustAddress("GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY")
	firstSigner := xdr.MustAddress("GCQZP3IU7XU6EJ63JZXKCQOYT2RNXN3HB5CNHENNUEUHSMA4VUJJJSEN")
	secondSigner := xdr.MustAddress("GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H")
	thirdSigner := xdr.MustAddress("GACMZD5VJXTRLKVET72CETCYKELPNCOTTBDC6DHFEUPLG5DHEK534JQX")
	formerSponsor := xdr.MustAddress("GAHK7EEG2WWHVKDNT4CEQFZGKF2LGDSW2IVM4S5DP42RBW3K6BTODB4A")
	oldSponsor := xdr.MustAddress("GANFZDRBCNTUXIODCJEYMACPMCSZEVE4WZGZ3CZDZ3P2SXK4KH75IK6Y")
	updatedSponsor := xdr.MustAddress("GAHK7EEG2WWHVKDNT4CEQFZGKF2LGDSW2IVM4S5DP42RBW3K6BTODB4A")
	newSponsor := xdr.MustAddress("GDEOVUDLCYTO46D6GD6WH7BFESPBV5RACC6F6NUFCIRU7PL2XONQHVGJ")

	expectedEffects := []effect{
		{
			address:     source.Address(),
			operationID: 249108107265,
			details: map[string]interface{}{
				"sponsor": newSponsor.Address(),
				"signer":  thirdSigner.Address(),
			},
			effectType: history.EffectSignerSponsorshipCreated,
			order:      1,
		},
		{
			address:     source.Address(),
			operationID: 249108107265,
			details: map[string]interface{}{
				"former_sponsor": oldSponsor.Address(),
				"new_sponsor":    updatedSponsor.Address(),
				"signer":         secondSigner.Address(),
			},
			effectType: history.EffectSignerSponsorshipUpdated,
			order:      2,
		},
		{
			address:     source.Address(),
			operationID: 249108107265,
			details: map[string]interface{}{
				"former_sponsor": formerSponsor.Address(),
				"signer":         firstSigner.Address(),
			},
			effectType: history.EffectSignerSponsorshipRemoved,
			order:      3,
		},
	}

	accountSignersMeta := &xdr.TransactionMeta{
		V: 1,
		V1: &xdr.TransactionMetaV1{
			TxChanges: xdr.LedgerEntryChanges{},
			Operations: []xdr.OperationMeta{
				{
					Changes: xdr.LedgerEntryChanges{
						{
							Type: xdr.LedgerEntryChangeTypeLedgerEntryState,
							State: &xdr.LedgerEntry{
								LastModifiedLedgerSeq: 0x39,
								Data: xdr.LedgerEntryData{
									Type: xdr.LedgerEntryTypeAccount,
									Account: &xdr.AccountEntry{
										AccountId:     source,
										Balance:       800152367009533292,
										SeqNum:        26,
										InflationDest: &source,
										Thresholds:    xdr.Thresholds{0x1, 0x0, 0x0, 0x0},
										Signers: []xdr.Signer{
											{
												Key: xdr.SignerKey{
													Type:    xdr.SignerKeyTypeSignerKeyTypeEd25519,
													Ed25519: firstSigner.Ed25519,
												},
												Weight: 10,
											},
											{
												Key: xdr.SignerKey{
													Type:    xdr.SignerKeyTypeSignerKeyTypeEd25519,
													Ed25519: secondSigner.Ed25519,
												},
												Weight: 10,
											},
											{
												Key: xdr.SignerKey{
													Type:    xdr.SignerKeyTypeSignerKeyTypeEd25519,
													Ed25519: thirdSigner.Ed25519,
												},
												Weight: 10,
											},
										},
										Ext: xdr.AccountEntryExt{
											V: 1,
											V1: &xdr.AccountEntryExtensionV1{
												Liabilities: xdr.Liabilities{},
												Ext: xdr.AccountEntryExtensionV1Ext{
													V: 2,
													V2: &xdr.AccountEntryExtensionV2{
														NumSponsored:  0,
														NumSponsoring: 0,
														SignerSponsoringIDs: []xdr.SponsorshipDescriptor{
															&formerSponsor,
															&oldSponsor,
															nil,
														},
													},
												},
											},
										},
									},
								},
							},
						},
						{
							Type: xdr.LedgerEntryChangeTypeLedgerEntryUpdated,
							Updated: &xdr.LedgerEntry{
								LastModifiedLedgerSeq: 0x39,
								Data: xdr.LedgerEntryData{
									Type: xdr.LedgerEntryTypeAccount,
									Account: &xdr.AccountEntry{
										AccountId:     source,
										Balance:       800152367009533292,
										SeqNum:        26,
										InflationDest: &source,
										Thresholds:    xdr.Thresholds{0x1, 0x0, 0x0, 0x0},
										Signers: []xdr.Signer{
											{
												Key: xdr.SignerKey{
													Type:    xdr.SignerKeyTypeSignerKeyTypeEd25519,
													Ed25519: secondSigner.Ed25519,
												},
												Weight: 10,
											},
											{
												Key: xdr.SignerKey{
													Type:    xdr.SignerKeyTypeSignerKeyTypeEd25519,
													Ed25519: thirdSigner.Ed25519,
												},
												Weight: 10,
											},
										},
										Ext: xdr.AccountEntryExt{
											V: 1,
											V1: &xdr.AccountEntryExtensionV1{
												Liabilities: xdr.Liabilities{},
												Ext: xdr.AccountEntryExtensionV1Ext{
													V: 2,
													V2: &xdr.AccountEntryExtensionV2{
														NumSponsored:  0,
														NumSponsoring: 0,
														SignerSponsoringIDs: []xdr.SponsorshipDescriptor{
															&updatedSponsor,
															&newSponsor,
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	b64, err := xdr.MarshalBase64(accountSignersMeta)
	assert.NoError(t, err)

	return b64, expectedEffects
}

func TestEffectsCoversAllOperationTypes(t *testing.T) {
	for typ, s := range xdr.OperationTypeToStringMap {
		op := xdr.Operation{
			Body: xdr.OperationBody{
				Type: xdr.OperationType(typ),
			},
		}
		operation := transactionOperationWrapper{
			index: 0,
			transaction: ingest.LedgerTransaction{
				UnsafeMeta: xdr.TransactionMeta{
					V:  2,
					V2: &xdr.TransactionMetaV2{},
				},
			},
			operation:      op,
			ledgerSequence: 1,
		}
		// calling effects should either panic (because the operation field is set to nil)
		// or not error
		func() {
			var err error
			defer func() {
				err2 := recover()
				if err != nil {
					assert.NotContains(t, err.Error(), "Unknown operation type")
				}
				assert.True(t, err2 != nil || err == nil, s)
			}()
			_, err = operation.effects()
		}()
	}

	// make sure the check works for an unknown operation type
	op := xdr.Operation{
		Body: xdr.OperationBody{
			Type: xdr.OperationType(20000),
		},
	}
	operation := transactionOperationWrapper{
		index: 0,
		transaction: ingest.LedgerTransaction{
			UnsafeMeta: xdr.TransactionMeta{
				V:  2,
				V2: &xdr.TransactionMetaV2{},
			},
		},
		operation:      op,
		ledgerSequence: 1,
	}
	// calling effects should error due to the unknown operation
	_, err := operation.effects()
	assert.Contains(t, err.Error(), "Unknown operation type")
}

func TestOperationEffects(t *testing.T) {

	sourceAID := xdr.MustAddress("GD3MMHD2YZWL5RAUWG6O3RMA5HTZYM7S3JLSZ2Z35JNJAWTDIKXY737V")
	sourceAccount := xdr.MuxedAccount{
		Type: xdr.CryptoKeyTypeKeyTypeMuxedEd25519,
		Med25519: &xdr.MuxedAccountMed25519{
			Id:      0xcafebabe,
			Ed25519: *sourceAID.Ed25519,
		},
	}
	destAID := xdr.MustAddress("GDEOVUDLCYTO46D6GD6WH7BFESPBV5RACC6F6NUFCIRU7PL2XONQHVGJ")
	dest := xdr.MuxedAccount{
		Type: xdr.CryptoKeyTypeKeyTypeMuxedEd25519,
		Med25519: &xdr.MuxedAccountMed25519{
			Id:      0xcafebabe,
			Ed25519: *destAID.Ed25519,
		},
	}
	strictPaymentWithMuxedAccountsTx := xdr.TransactionEnvelope{
		Type: xdr.EnvelopeTypeEnvelopeTypeTx,
		V1: &xdr.TransactionV1Envelope{
			Tx: xdr.Transaction{
				SourceAccount: sourceAccount,
				Fee:           100,
				SeqNum:        3684420515004429,
				Operations: []xdr.Operation{
					{
						Body: xdr.OperationBody{
							Type: xdr.OperationTypePathPaymentStrictSend,
							PathPaymentStrictSendOp: &xdr.PathPaymentStrictSendOp{
								SendAsset: xdr.Asset{
									Type: xdr.AssetTypeAssetTypeCreditAlphanum4,
									AlphaNum4: &xdr.AlphaNum4{
										AssetCode: xdr.AssetCode4{66, 82, 76, 0},
										Issuer:    xdr.MustAddress("GCXI6Q73J7F6EUSBZTPW4G4OUGVDHABPYF2U4KO7MVEX52OH5VMVUCRF"),
									},
								},
								SendAmount:  300000,
								Destination: dest,
								DestAsset: xdr.Asset{
									Type: 1,
									AlphaNum4: &xdr.AlphaNum4{
										AssetCode: xdr.AssetCode4{65, 82, 83, 0},
										Issuer:    xdr.MustAddress("GCXI6Q73J7F6EUSBZTPW4G4OUGVDHABPYF2U4KO7MVEX52OH5VMVUCRF"),
									},
								},
								DestMin: 10000000,
								Path: []xdr.Asset{
									{
										Type: xdr.AssetTypeAssetTypeCreditAlphanum4,
										AlphaNum4: &xdr.AlphaNum4{
											AssetCode: xdr.AssetCode4{65, 82, 83, 0},
											Issuer:    xdr.MustAddress("GCXI6Q73J7F6EUSBZTPW4G4OUGVDHABPYF2U4KO7MVEX52OH5VMVUCRF"),
										},
									},
								},
							},
						},
					},
				},
			},
			Signatures: []xdr.DecoratedSignature{
				{
					Hint:      xdr.SignatureHint{99, 66, 175, 143},
					Signature: xdr.Signature{244, 107, 139, 92, 189, 156, 207, 79, 84, 56, 2, 70, 75, 22, 237, 50, 100, 242, 159, 177, 27, 240, 66, 122, 182, 45, 189, 78, 5, 127, 26, 61, 179, 238, 229, 76, 32, 206, 122, 13, 154, 133, 148, 149, 29, 250, 48, 132, 44, 86, 163, 56, 32, 44, 75, 87, 226, 251, 76, 4, 59, 182, 132, 8},
				},
			},
		},
	}
	strictPaymentWithMuxedAccountsTxBase64, err := xdr.MarshalBase64(strictPaymentWithMuxedAccountsTx)
	assert.NoError(t, err)

	creator := xdr.MustAddress("GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H")
	created := xdr.MustAddress("GCQZP3IU7XU6EJ63JZXKCQOYT2RNXN3HB5CNHENNUEUHSMA4VUJJJSEN")
	sponsor := xdr.MustAddress("GAHK7EEG2WWHVKDNT4CEQFZGKF2LGDSW2IVM4S5DP42RBW3K6BTODB4A")
	sponsor2 := xdr.MustAddress("GACMZD5VJXTRLKVET72CETCYKELPNCOTTBDC6DHFEUPLG5DHEK534JQX")
	createAccountMeta := &xdr.TransactionMeta{
		V: 1,
		V1: &xdr.TransactionMetaV1{
			TxChanges: xdr.LedgerEntryChanges{
				{
					Type: 3,
					State: &xdr.LedgerEntry{
						LastModifiedLedgerSeq: 0x39,
						Data: xdr.LedgerEntryData{
							Type: 0,
							Account: &xdr.AccountEntry{
								AccountId:     creator,
								Balance:       800152377009533292,
								SeqNum:        25,
								InflationDest: &creator,
								Thresholds:    xdr.Thresholds{0x1, 0x0, 0x0, 0x0},
							},
						},
					},
				},
				{
					Type: 1,
					Updated: &xdr.LedgerEntry{
						LastModifiedLedgerSeq: 0x39,
						Data: xdr.LedgerEntryData{
							Type: 0,
							Account: &xdr.AccountEntry{
								AccountId:     creator,
								Balance:       800152377009533292,
								SeqNum:        26,
								InflationDest: &creator,
							},
						},
						Ext: xdr.LedgerEntryExt{},
					},
				},
			},
			Operations: []xdr.OperationMeta{
				{
					Changes: xdr.LedgerEntryChanges{
						{
							Type: xdr.LedgerEntryChangeTypeLedgerEntryState,
							State: &xdr.LedgerEntry{
								LastModifiedLedgerSeq: 0x39,
								Data: xdr.LedgerEntryData{
									Type: xdr.LedgerEntryTypeAccount,
									Account: &xdr.AccountEntry{
										AccountId:     creator,
										Balance:       800152367009533292,
										SeqNum:        26,
										InflationDest: &creator,
										Thresholds:    xdr.Thresholds{0x1, 0x0, 0x0, 0x0},
									},
								},
								Ext: xdr.LedgerEntryExt{
									V: 1,
									V1: &xdr.LedgerEntryExtensionV1{
										SponsoringId: &sponsor2,
									},
								},
							},
						},
						{
							Type: xdr.LedgerEntryChangeTypeLedgerEntryRemoved,
							Removed: &xdr.LedgerKey{
								Type: xdr.LedgerEntryTypeAccount,
								Account: &xdr.LedgerKeyAccount{
									AccountId: created,
								},
							},
						},
						{
							Type: xdr.LedgerEntryChangeTypeLedgerEntryState,
							State: &xdr.LedgerEntry{
								LastModifiedLedgerSeq: 0x39,
								Data: xdr.LedgerEntryData{
									Type: xdr.LedgerEntryTypeAccount,
									Account: &xdr.AccountEntry{
										AccountId:     creator,
										Balance:       800152367009533292,
										SeqNum:        26,
										InflationDest: &creator,
										Thresholds:    xdr.Thresholds{0x1, 0x0, 0x0, 0x0},
									},
								},
								Ext: xdr.LedgerEntryExt{
									V: 1,
									V1: &xdr.LedgerEntryExtensionV1{
										SponsoringId: &sponsor,
									},
								},
							},
						},
						{
							Type: xdr.LedgerEntryChangeTypeLedgerEntryUpdated,
							Updated: &xdr.LedgerEntry{
								LastModifiedLedgerSeq: 0x39,
								Data: xdr.LedgerEntryData{
									Type: xdr.LedgerEntryTypeAccount,
									Account: &xdr.AccountEntry{
										AccountId:     creator,
										Balance:       800152367009533292,
										SeqNum:        26,
										InflationDest: &creator,
										Thresholds:    xdr.Thresholds{0x1, 0x0, 0x0, 0x0},
									},
								},
								Ext: xdr.LedgerEntryExt{
									V: 1,
									V1: &xdr.LedgerEntryExtensionV1{
										SponsoringId: &sponsor2,
									},
								},
							},
						},
						{
							Type: xdr.LedgerEntryChangeTypeLedgerEntryState,
							State: &xdr.LedgerEntry{
								LastModifiedLedgerSeq: 0x39,
								Data: xdr.LedgerEntryData{
									Type: xdr.LedgerEntryTypeAccount,
									Account: &xdr.AccountEntry{
										AccountId:     creator,
										Balance:       800152377009533292,
										SeqNum:        26,
										InflationDest: &creator,
										Thresholds:    xdr.Thresholds{0x1, 0x0, 0x0, 0x0},
									},
								},
								Ext: xdr.LedgerEntryExt{
									V: 1,
									V1: &xdr.LedgerEntryExtensionV1{
										SponsoringId: &sponsor,
									},
								},
							},
						},
						{
							Type: xdr.LedgerEntryChangeTypeLedgerEntryUpdated,
							Updated: &xdr.LedgerEntry{
								LastModifiedLedgerSeq: 0x39,
								Data: xdr.LedgerEntryData{
									Type: xdr.LedgerEntryTypeAccount,
									Account: &xdr.AccountEntry{
										AccountId:     creator,
										Balance:       800152367009533292,
										SeqNum:        26,
										InflationDest: &creator,
										Thresholds:    xdr.Thresholds{0x1, 0x0, 0x0, 0x0},
									},
								},
								Ext: xdr.LedgerEntryExt{
									V: 1,
									V1: &xdr.LedgerEntryExtensionV1{
										SponsoringId: &sponsor,
									},
								},
							},
						},
						{
							Type: xdr.LedgerEntryChangeTypeLedgerEntryCreated,
							Created: &xdr.LedgerEntry{
								LastModifiedLedgerSeq: 0x39,
								Data: xdr.LedgerEntryData{
									Type: xdr.LedgerEntryTypeAccount,
									Account: &xdr.AccountEntry{
										AccountId:  created,
										Balance:    10000000000,
										SeqNum:     244813135872,
										Thresholds: xdr.Thresholds{0x1, 0x0, 0x0, 0x0},
									},
								},
								Ext: xdr.LedgerEntryExt{
									V: 1,
									V1: &xdr.LedgerEntryExtensionV1{
										SponsoringId: &sponsor,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	createAccountMetaB64, err := xdr.MarshalBase64(createAccountMeta)
	assert.NoError(t, err)
	assert.NoError(t, err)

	revokeSponsorshipMeta, revokeSponsorshipEffects := getRevokeSponsorshipMeta(t)

	testCases := []struct {
		desc          string
		envelopeXDR   string
		resultXDR     string
		metaXDR       string
		feeChangesXDR string
		hash          string
		index         uint32
		sequence      uint32
		expected      []effect
	}{
		{
			desc:          "createAccount",
			envelopeXDR:   "AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAAaAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAoZftFP3p4ifbTm6hQdieotu3Zw9E05GtoSh5MBytEpQAAAACVAvkAAAAAAAAAAABVvwF9wAAAEDHU95E9wxgETD8TqxUrkgC0/7XHyNDts6Q5huRHfDRyRcoHdv7aMp/sPvC3RPkXjOMjgbKJUX7SgExUeYB5f8F",
			resultXDR:     "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=",
			metaXDR:       createAccountMetaB64,
			feeChangesXDR: "AAAAAgAAAAMAAAA3AAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsatlj11nHQAAAAAAAAABkAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAA5AAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsatlj11nFsAAAAAAAAABkAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			hash:          "0e5bd332291e3098e49886df2cdb9b5369a5f9e0a9973f0d9e1a9489c6581ba2",
			index:         0,
			sequence:      57,
			expected: []effect{
				{
					address:     "GCQZP3IU7XU6EJ63JZXKCQOYT2RNXN3HB5CNHENNUEUHSMA4VUJJJSEN",
					operationID: int64(244813139969),
					details: map[string]interface{}{
						"starting_balance": "1000.0000000",
					},
					effectType: history.EffectAccountCreated,
					order:      uint32(1),
				},
				{
					address:     "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
					operationID: int64(244813139969),
					details: map[string]interface{}{
						"amount":     "1000.0000000",
						"asset_type": "native",
					},
					effectType: history.EffectAccountDebited,
					order:      uint32(2),
				},
				{
					address:     "GCQZP3IU7XU6EJ63JZXKCQOYT2RNXN3HB5CNHENNUEUHSMA4VUJJJSEN",
					operationID: int64(244813139969),
					details: map[string]interface{}{
						"public_key": "GCQZP3IU7XU6EJ63JZXKCQOYT2RNXN3HB5CNHENNUEUHSMA4VUJJJSEN",
						"weight":     1,
					},
					effectType: history.EffectSignerCreated,
					order:      uint32(3),
				},
				{
					address:     "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
					operationID: int64(244813139969),
					details: map[string]interface{}{
						"former_sponsor": "GACMZD5VJXTRLKVET72CETCYKELPNCOTTBDC6DHFEUPLG5DHEK534JQX",
					},
					effectType: history.EffectAccountSponsorshipRemoved,
					order:      uint32(4),
				},
				{
					address:     "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
					operationID: int64(244813139969),
					details: map[string]interface{}{
						"former_sponsor": "GAHK7EEG2WWHVKDNT4CEQFZGKF2LGDSW2IVM4S5DP42RBW3K6BTODB4A",
						"new_sponsor":    "GACMZD5VJXTRLKVET72CETCYKELPNCOTTBDC6DHFEUPLG5DHEK534JQX",
					},
					effectType: history.EffectAccountSponsorshipUpdated,
					order:      uint32(5),
				},
				{
					address:     "GCQZP3IU7XU6EJ63JZXKCQOYT2RNXN3HB5CNHENNUEUHSMA4VUJJJSEN",
					operationID: int64(244813139969),
					details: map[string]interface{}{
						"sponsor": "GAHK7EEG2WWHVKDNT4CEQFZGKF2LGDSW2IVM4S5DP42RBW3K6BTODB4A",
					},
					effectType: history.EffectAccountSponsorshipCreated,
					order:      uint32(6),
				},
			},
		},
		{
			desc:          "payment",
			envelopeXDR:   "AAAAABpcjiETZ0uhwxJJhgBPYKWSVJy2TZ2LI87fqV1cUf/UAAAAZAAAADcAAAABAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAGlyOIRNnS6HDEkmGAE9gpZJUnLZNnYsjzt+pXVxR/9QAAAAAAAAAAAX14QAAAAAAAAAAAVxR/9QAAABAK6pcXYMzAEmH08CZ1LWmvtNDKauhx+OImtP/Lk4hVTMJRVBOebVs5WEPj9iSrgGT0EswuDCZ2i5AEzwgGof9Ag==",
			resultXDR:     "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=",
			metaXDR:       "AAAAAQAAAAIAAAADAAAAOAAAAAAAAAAAGlyOIRNnS6HDEkmGAE9gpZJUnLZNnYsjzt+pXVxR/9QAAAACVAvjnAAAADcAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAOAAAAAAAAAAAGlyOIRNnS6HDEkmGAE9gpZJUnLZNnYsjzt+pXVxR/9QAAAACVAvjnAAAADcAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAA==",
			feeChangesXDR: "AAAAAgAAAAMAAAA3AAAAAAAAAAAaXI4hE2dLocMSSYYAT2ClklSctk2diyPO36ldXFH/1AAAAAJUC+QAAAAANwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAA4AAAAAAAAAAAaXI4hE2dLocMSSYYAT2ClklSctk2diyPO36ldXFH/1AAAAAJUC+OcAAAANwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			hash:          "2a805712c6d10f9e74bb0ccf54ae92a2b4b1e586451fe8133a2433816f6b567c",
			index:         0,
			sequence:      56,
			expected: []effect{
				{
					address: "GANFZDRBCNTUXIODCJEYMACPMCSZEVE4WZGZ3CZDZ3P2SXK4KH75IK6Y",
					details: map[string]interface{}{
						"amount":     "10.0000000",
						"asset_type": "native",
					},
					effectType:  history.EffectAccountCredited,
					operationID: int64(240518172673),
					order:       uint32(1),
				},
				{
					address: "GANFZDRBCNTUXIODCJEYMACPMCSZEVE4WZGZ3CZDZ3P2SXK4KH75IK6Y",
					details: map[string]interface{}{
						"amount":     "10.0000000",
						"asset_type": "native",
					},
					effectType:  history.EffectAccountDebited,
					operationID: int64(240518172673),
					order:       uint32(2),
				},
			},
		},
		{
			desc:          "pathPaymentStrictSend",
			envelopeXDR:   "AAAAAPbGHHrGbL7EFLG87cWA6eecM/LaVyzrO+pakFpjQq+PAAAAZAANFvYAAAANAAAAAAAAAAAAAAABAAAAAAAAAA0AAAABQlJMAAAAAACuj0P7T8viUkHM324bjqGqM4AvwXVOKd9lSX7px+1ZWgAAAAAABJPgAAAAAMjq0GsWJu54fjD9Y/wlJJ4a9iAQvF82hRIjT716u5sDAAAAAUFSUwAAAAAAro9D+0/L4lJBzN9uG46hqjOAL8F1TinfZUl+6cftWVoAAAAAAJiWgAAAAAEAAAABQVJTAAAAAACuj0P7T8viUkHM324bjqGqM4AvwXVOKd9lSX7px+1ZWgAAAAAAAAABY0KvjwAAAED0a4tcvZzPT1Q4AkZLFu0yZPKfsRvwQnq2Lb1OBX8aPbPu5UwgznoNmoWUlR36MIQsVqM4ICxLV+L7TAQ7toQI",
			resultXDR:     "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAANAAAAAAAAAAEAAAAAyOrQaxYm7nh+MP1j/CUknhr2IBC8XzaFEiNPvXq7mwMAAAAAAJmwQAAAAAFBUlMAAAAAAK6PQ/tPy+JSQczfbhuOoaozgC/BdU4p32VJfunH7VlaAAAAAACYloAAAAABQlJMAAAAAACuj0P7T8viUkHM324bjqGqM4AvwXVOKd9lSX7px+1ZWgAAAAAABJPgAAAAAMjq0GsWJu54fjD9Y/wlJJ4a9iAQvF82hRIjT716u5sDAAAAAUFSUwAAAAAAro9D+0/L4lJBzN9uG46hqjOAL8F1TinfZUl+6cftWVoAAAAAAJiWgAAAAAA=",
			metaXDR:       "AAAAAQAAAAIAAAADAA0aVQAAAAAAAAAA9sYcesZsvsQUsbztxYDp55wz8tpXLOs76lqQWmNCr48AAAAXSHbi7AANFvYAAAAMAAAAAwAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAA0aVQAAAAAAAAAA9sYcesZsvsQUsbztxYDp55wz8tpXLOs76lqQWmNCr48AAAAXSHbi7AANFvYAAAANAAAAAwAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAACAAAAAMADRo0AAAAAQAAAAD2xhx6xmy+xBSxvO3FgOnnnDPy2lcs6zvqWpBaY0KvjwAAAAFCUkwAAAAAAK6PQ/tPy+JSQczfbhuOoaozgC/BdU4p32VJfunH7VlaAAAAAB22gaB//////////wAAAAEAAAABAAAAAAC3GwAAAAAAAAAAAAAAAAAAAAAAAAAAAQANGlUAAAABAAAAAPbGHHrGbL7EFLG87cWA6eecM/LaVyzrO+pakFpjQq+PAAAAAUJSTAAAAAAAro9D+0/L4lJBzN9uG46hqjOAL8F1TinfZUl+6cftWVoAAAAAHbHtwH//////////AAAAAQAAAAEAAAAAALcbAAAAAAAAAAAAAAAAAAAAAAAAAAADAA0aNAAAAAIAAAAAyOrQaxYm7nh+MP1j/CUknhr2IBC8XzaFEiNPvXq7mwMAAAAAAJmwQAAAAAFBUlMAAAAAAK6PQ/tPy+JSQczfbhuOoaozgC/BdU4p32VJfunH7VlaAAAAAUJSTAAAAAAAro9D+0/L4lJBzN9uG46hqjOAL8F1TinfZUl+6cftWVoAAAAAFNyTgAAAAAMAAABkAAAAAAAAAAAAAAAAAAAAAQANGlUAAAACAAAAAMjq0GsWJu54fjD9Y/wlJJ4a9iAQvF82hRIjT716u5sDAAAAAACZsEAAAAABQVJTAAAAAACuj0P7T8viUkHM324bjqGqM4AvwXVOKd9lSX7px+1ZWgAAAAFCUkwAAAAAAK6PQ/tPy+JSQczfbhuOoaozgC/BdU4p32VJfunH7VlaAAAAABRD/QAAAAADAAAAZAAAAAAAAAAAAAAAAAAAAAMADRo0AAAAAQAAAADI6tBrFibueH4w/WP8JSSeGvYgELxfNoUSI0+9erubAwAAAAFCUkwAAAAAAK6PQ/tPy+JSQczfbhuOoaozgC/BdU4p32VJfunH7VlaAAAAAB3kSGB//////////wAAAAEAAAABAAAAAACgN6AAAAAAAAAAAAAAAAAAAAAAAAAAAQANGlUAAAABAAAAAMjq0GsWJu54fjD9Y/wlJJ4a9iAQvF82hRIjT716u5sDAAAAAUJSTAAAAAAAro9D+0/L4lJBzN9uG46hqjOAL8F1TinfZUl+6cftWVoAAAAAHejcQH//////////AAAAAQAAAAEAAAAAAJujwAAAAAAAAAAAAAAAAAAAAAAAAAADAA0aNAAAAAEAAAAAyOrQaxYm7nh+MP1j/CUknhr2IBC8XzaFEiNPvXq7mwMAAAABQVJTAAAAAACuj0P7T8viUkHM324bjqGqM4AvwXVOKd9lSX7px+1ZWgAAAAB2BGcAf/////////8AAAABAAAAAQAAAAAAAAAAAAAAABTck4AAAAAAAAAAAAAAAAEADRpVAAAAAQAAAADI6tBrFibueH4w/WP8JSSeGvYgELxfNoUSI0+9erubAwAAAAFBUlMAAAAAAK6PQ/tPy+JSQczfbhuOoaozgC/BdU4p32VJfunH7VlaAAAAAHYEZwB//////////wAAAAEAAAABAAAAAAAAAAAAAAAAFEP9AAAAAAAAAAAA",
			feeChangesXDR: "AAAAAgAAAAMADRpIAAAAAAAAAAD2xhx6xmy+xBSxvO3FgOnnnDPy2lcs6zvqWpBaY0KvjwAAABdIduNQAA0W9gAAAAwAAAADAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEADRpVAAAAAAAAAAD2xhx6xmy+xBSxvO3FgOnnnDPy2lcs6zvqWpBaY0KvjwAAABdIduLsAA0W9gAAAAwAAAADAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			hash:          "96415ac1d2f79621b26b1568f963fd8dd6c50c20a22c7428cefbfe9dee867588",
			index:         0,
			sequence:      20,
			expected: []effect{
				{
					address: "GDEOVUDLCYTO46D6GD6WH7BFESPBV5RACC6F6NUFCIRU7PL2XONQHVGJ",
					details: map[string]interface{}{
						"amount":       "1.0000000",
						"asset_code":   "ARS",
						"asset_type":   "credit_alphanum4",
						"asset_issuer": "GCXI6Q73J7F6EUSBZTPW4G4OUGVDHABPYF2U4KO7MVEX52OH5VMVUCRF",
					},
					effectType:  history.EffectAccountCredited,
					operationID: int64(85899350017),
					order:       uint32(1),
				},
				{
					address: "GD3MMHD2YZWL5RAUWG6O3RMA5HTZYM7S3JLSZ2Z35JNJAWTDIKXY737V",
					details: map[string]interface{}{
						"amount":       "0.0300000",
						"asset_code":   "BRL",
						"asset_type":   "credit_alphanum4",
						"asset_issuer": "GCXI6Q73J7F6EUSBZTPW4G4OUGVDHABPYF2U4KO7MVEX52OH5VMVUCRF",
					},
					effectType:  history.EffectAccountDebited,
					operationID: int64(85899350017),
					order:       uint32(2),
				},
				{
					address: "GD3MMHD2YZWL5RAUWG6O3RMA5HTZYM7S3JLSZ2Z35JNJAWTDIKXY737V",
					details: map[string]interface{}{
						"seller":              "GDEOVUDLCYTO46D6GD6WH7BFESPBV5RACC6F6NUFCIRU7PL2XONQHVGJ",
						"offer_id":            xdr.Int64(10072128),
						"sold_amount":         "0.0300000",
						"bought_amount":       "1.0000000",
						"sold_asset_code":     "BRL",
						"sold_asset_type":     "credit_alphanum4",
						"bought_asset_code":   "ARS",
						"bought_asset_type":   "credit_alphanum4",
						"sold_asset_issuer":   "GCXI6Q73J7F6EUSBZTPW4G4OUGVDHABPYF2U4KO7MVEX52OH5VMVUCRF",
						"bought_asset_issuer": "GCXI6Q73J7F6EUSBZTPW4G4OUGVDHABPYF2U4KO7MVEX52OH5VMVUCRF",
					},
					effectType:  history.EffectTrade,
					operationID: int64(85899350017),
					order:       uint32(3),
				},
				{
					address: "GDEOVUDLCYTO46D6GD6WH7BFESPBV5RACC6F6NUFCIRU7PL2XONQHVGJ",
					details: map[string]interface{}{
						"seller":              "GD3MMHD2YZWL5RAUWG6O3RMA5HTZYM7S3JLSZ2Z35JNJAWTDIKXY737V",
						"offer_id":            xdr.Int64(10072128),
						"sold_amount":         "1.0000000",
						"bought_amount":       "0.0300000",
						"sold_asset_code":     "ARS",
						"sold_asset_type":     "credit_alphanum4",
						"bought_asset_code":   "BRL",
						"bought_asset_type":   "credit_alphanum4",
						"sold_asset_issuer":   "GCXI6Q73J7F6EUSBZTPW4G4OUGVDHABPYF2U4KO7MVEX52OH5VMVUCRF",
						"bought_asset_issuer": "GCXI6Q73J7F6EUSBZTPW4G4OUGVDHABPYF2U4KO7MVEX52OH5VMVUCRF",
					},
					effectType:  history.EffectTrade,
					operationID: int64(85899350017),
					order:       uint32(4),
				},
			},
		},
		{
			desc:          "pathPaymentStrictSend with muxed accounts",
			envelopeXDR:   strictPaymentWithMuxedAccountsTxBase64,
			resultXDR:     "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAANAAAAAAAAAAEAAAAAyOrQaxYm7nh+MP1j/CUknhr2IBC8XzaFEiNPvXq7mwMAAAAAAJmwQAAAAAFBUlMAAAAAAK6PQ/tPy+JSQczfbhuOoaozgC/BdU4p32VJfunH7VlaAAAAAACYloAAAAABQlJMAAAAAACuj0P7T8viUkHM324bjqGqM4AvwXVOKd9lSX7px+1ZWgAAAAAABJPgAAAAAMjq0GsWJu54fjD9Y/wlJJ4a9iAQvF82hRIjT716u5sDAAAAAUFSUwAAAAAAro9D+0/L4lJBzN9uG46hqjOAL8F1TinfZUl+6cftWVoAAAAAAJiWgAAAAAA=",
			metaXDR:       "AAAAAQAAAAIAAAADAA0aVQAAAAAAAAAA9sYcesZsvsQUsbztxYDp55wz8tpXLOs76lqQWmNCr48AAAAXSHbi7AANFvYAAAAMAAAAAwAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAA0aVQAAAAAAAAAA9sYcesZsvsQUsbztxYDp55wz8tpXLOs76lqQWmNCr48AAAAXSHbi7AANFvYAAAANAAAAAwAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAACAAAAAMADRo0AAAAAQAAAAD2xhx6xmy+xBSxvO3FgOnnnDPy2lcs6zvqWpBaY0KvjwAAAAFCUkwAAAAAAK6PQ/tPy+JSQczfbhuOoaozgC/BdU4p32VJfunH7VlaAAAAAB22gaB//////////wAAAAEAAAABAAAAAAC3GwAAAAAAAAAAAAAAAAAAAAAAAAAAAQANGlUAAAABAAAAAPbGHHrGbL7EFLG87cWA6eecM/LaVyzrO+pakFpjQq+PAAAAAUJSTAAAAAAAro9D+0/L4lJBzN9uG46hqjOAL8F1TinfZUl+6cftWVoAAAAAHbHtwH//////////AAAAAQAAAAEAAAAAALcbAAAAAAAAAAAAAAAAAAAAAAAAAAADAA0aNAAAAAIAAAAAyOrQaxYm7nh+MP1j/CUknhr2IBC8XzaFEiNPvXq7mwMAAAAAAJmwQAAAAAFBUlMAAAAAAK6PQ/tPy+JSQczfbhuOoaozgC/BdU4p32VJfunH7VlaAAAAAUJSTAAAAAAAro9D+0/L4lJBzN9uG46hqjOAL8F1TinfZUl+6cftWVoAAAAAFNyTgAAAAAMAAABkAAAAAAAAAAAAAAAAAAAAAQANGlUAAAACAAAAAMjq0GsWJu54fjD9Y/wlJJ4a9iAQvF82hRIjT716u5sDAAAAAACZsEAAAAABQVJTAAAAAACuj0P7T8viUkHM324bjqGqM4AvwXVOKd9lSX7px+1ZWgAAAAFCUkwAAAAAAK6PQ/tPy+JSQczfbhuOoaozgC/BdU4p32VJfunH7VlaAAAAABRD/QAAAAADAAAAZAAAAAAAAAAAAAAAAAAAAAMADRo0AAAAAQAAAADI6tBrFibueH4w/WP8JSSeGvYgELxfNoUSI0+9erubAwAAAAFCUkwAAAAAAK6PQ/tPy+JSQczfbhuOoaozgC/BdU4p32VJfunH7VlaAAAAAB3kSGB//////////wAAAAEAAAABAAAAAACgN6AAAAAAAAAAAAAAAAAAAAAAAAAAAQANGlUAAAABAAAAAMjq0GsWJu54fjD9Y/wlJJ4a9iAQvF82hRIjT716u5sDAAAAAUJSTAAAAAAAro9D+0/L4lJBzN9uG46hqjOAL8F1TinfZUl+6cftWVoAAAAAHejcQH//////////AAAAAQAAAAEAAAAAAJujwAAAAAAAAAAAAAAAAAAAAAAAAAADAA0aNAAAAAEAAAAAyOrQaxYm7nh+MP1j/CUknhr2IBC8XzaFEiNPvXq7mwMAAAABQVJTAAAAAACuj0P7T8viUkHM324bjqGqM4AvwXVOKd9lSX7px+1ZWgAAAAB2BGcAf/////////8AAAABAAAAAQAAAAAAAAAAAAAAABTck4AAAAAAAAAAAAAAAAEADRpVAAAAAQAAAADI6tBrFibueH4w/WP8JSSeGvYgELxfNoUSI0+9erubAwAAAAFBUlMAAAAAAK6PQ/tPy+JSQczfbhuOoaozgC/BdU4p32VJfunH7VlaAAAAAHYEZwB//////////wAAAAEAAAABAAAAAAAAAAAAAAAAFEP9AAAAAAAAAAAA",
			feeChangesXDR: "AAAAAgAAAAMADRpIAAAAAAAAAAD2xhx6xmy+xBSxvO3FgOnnnDPy2lcs6zvqWpBaY0KvjwAAABdIduNQAA0W9gAAAAwAAAADAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEADRpVAAAAAAAAAAD2xhx6xmy+xBSxvO3FgOnnnDPy2lcs6zvqWpBaY0KvjwAAABdIduLsAA0W9gAAAAwAAAADAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			hash:          "96415ac1d2f79621b26b1568f963fd8dd6c50c20a22c7428cefbfe9dee867588",
			index:         0,
			sequence:      20,
			expected: []effect{
				{
					address:      "GDEOVUDLCYTO46D6GD6WH7BFESPBV5RACC6F6NUFCIRU7PL2XONQHVGJ",
					addressMuxed: null.StringFrom("MDEOVUDLCYTO46D6GD6WH7BFESPBV5RACC6F6NUFCIRU7PL2XONQGAAAAAAMV7V2X24II"),
					details: map[string]interface{}{
						"amount":       "1.0000000",
						"asset_code":   "ARS",
						"asset_type":   "credit_alphanum4",
						"asset_issuer": "GCXI6Q73J7F6EUSBZTPW4G4OUGVDHABPYF2U4KO7MVEX52OH5VMVUCRF",
					},
					effectType:  history.EffectAccountCredited,
					operationID: int64(85899350017),
					order:       uint32(1),
				},
				{
					address:      "GD3MMHD2YZWL5RAUWG6O3RMA5HTZYM7S3JLSZ2Z35JNJAWTDIKXY737V",
					addressMuxed: null.StringFrom("MD3MMHD2YZWL5RAUWG6O3RMA5HTZYM7S3JLSZ2Z35JNJAWTDIKXY6AAAAAAMV7V2XZY4C"),
					details: map[string]interface{}{
						"amount":       "0.0300000",
						"asset_code":   "BRL",
						"asset_type":   "credit_alphanum4",
						"asset_issuer": "GCXI6Q73J7F6EUSBZTPW4G4OUGVDHABPYF2U4KO7MVEX52OH5VMVUCRF",
					},
					effectType:  history.EffectAccountDebited,
					operationID: int64(85899350017),
					order:       uint32(2),
				},
				{
					address:      "GD3MMHD2YZWL5RAUWG6O3RMA5HTZYM7S3JLSZ2Z35JNJAWTDIKXY737V",
					addressMuxed: null.StringFrom("MD3MMHD2YZWL5RAUWG6O3RMA5HTZYM7S3JLSZ2Z35JNJAWTDIKXY6AAAAAAMV7V2XZY4C"),
					details: map[string]interface{}{
						"seller":              "GDEOVUDLCYTO46D6GD6WH7BFESPBV5RACC6F6NUFCIRU7PL2XONQHVGJ",
						"offer_id":            xdr.Int64(10072128),
						"sold_amount":         "0.0300000",
						"bought_amount":       "1.0000000",
						"sold_asset_code":     "BRL",
						"sold_asset_type":     "credit_alphanum4",
						"bought_asset_code":   "ARS",
						"bought_asset_type":   "credit_alphanum4",
						"sold_asset_issuer":   "GCXI6Q73J7F6EUSBZTPW4G4OUGVDHABPYF2U4KO7MVEX52OH5VMVUCRF",
						"bought_asset_issuer": "GCXI6Q73J7F6EUSBZTPW4G4OUGVDHABPYF2U4KO7MVEX52OH5VMVUCRF",
					},
					effectType:  history.EffectTrade,
					operationID: int64(85899350017),
					order:       uint32(3),
				},
				{
					address: "GDEOVUDLCYTO46D6GD6WH7BFESPBV5RACC6F6NUFCIRU7PL2XONQHVGJ",
					details: map[string]interface{}{
						"seller":              "GD3MMHD2YZWL5RAUWG6O3RMA5HTZYM7S3JLSZ2Z35JNJAWTDIKXY737V",
						"seller_muxed":        "MD3MMHD2YZWL5RAUWG6O3RMA5HTZYM7S3JLSZ2Z35JNJAWTDIKXY6AAAAAAMV7V2XZY4C",
						"seller_muxed_id":     uint64(0xcafebabe),
						"offer_id":            xdr.Int64(10072128),
						"sold_amount":         "1.0000000",
						"bought_amount":       "0.0300000",
						"sold_asset_code":     "ARS",
						"sold_asset_type":     "credit_alphanum4",
						"bought_asset_code":   "BRL",
						"bought_asset_type":   "credit_alphanum4",
						"sold_asset_issuer":   "GCXI6Q73J7F6EUSBZTPW4G4OUGVDHABPYF2U4KO7MVEX52OH5VMVUCRF",
						"bought_asset_issuer": "GCXI6Q73J7F6EUSBZTPW4G4OUGVDHABPYF2U4KO7MVEX52OH5VMVUCRF",
					},
					effectType:  history.EffectTrade,
					operationID: int64(85899350017),
					order:       uint32(4),
				},
			},
		},
		{
			desc:          "manageSellOffer - without claims",
			envelopeXDR:   "AAAAAC7C83M2T23Bu4kdQGqdfboZgjcxsJ2lBT23ifoRVFexAAAAZAAAABAAAAACAAAAAAAAAAAAAAABAAAAAAAAAAMAAAAAAAAAAVVTRAAAAAAALsLzczZPbcG7iR1Aap19uhmCNzGwnaUFPbeJ+hFUV7EAAAAA7msoAAAAAAEAAAACAAAAAAAAAAAAAAAAAAAAARFUV7EAAABALuai5QxceFbtAiC5nkntNVnvSPeWR+C+FgplPAdRgRS+PPESpUiSCyuiwuhmvuDw7kwxn+A6E0M4ca1s2qzMAg==",
			resultXDR:     "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAADAAAAAAAAAAAAAAAAAAAAAC7C83M2T23Bu4kdQGqdfboZgjcxsJ2lBT23ifoRVFexAAAAAAAAAAEAAAAAAAAAAVVTRAAAAAAALsLzczZPbcG7iR1Aap19uhmCNzGwnaUFPbeJ+hFUV7EAAAAA7msoAAAAAAEAAAACAAAAAAAAAAAAAAAA",
			metaXDR:       "AAAAAQAAAAIAAAADAAAAEgAAAAAAAAAALsLzczZPbcG7iR1Aap19uhmCNzGwnaUFPbeJ+hFUV7EAAAACVAvi1AAAABAAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAEgAAAAAAAAAALsLzczZPbcG7iR1Aap19uhmCNzGwnaUFPbeJ+hFUV7EAAAACVAvi1AAAABAAAAACAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAASAAAAAAAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsQAAAAJUC+LUAAAAEAAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAASAAAAAAAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsQAAAAJUC+LUAAAAEAAAAAIAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAABAAAAAAAAAAAAAAAA7msoAAAAAAAAAAAAAAAAAAAAABIAAAACAAAAAC7C83M2T23Bu4kdQGqdfboZgjcxsJ2lBT23ifoRVFexAAAAAAAAAAEAAAAAAAAAAVVTRAAAAAAALsLzczZPbcG7iR1Aap19uhmCNzGwnaUFPbeJ+hFUV7EAAAAA7msoAAAAAAEAAAACAAAAAAAAAAAAAAAA",
			feeChangesXDR: "AAAAAgAAAAMAAAASAAAAAAAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsQAAAAJUC+OcAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAASAAAAAAAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsQAAAAJUC+M4AAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			hash:          "ca756d1519ceda79f8722042b12cea7ba004c3bd961adb62b59f88a867f86eb3",
			index:         0,
			sequence:      56,
			expected:      []effect{},
		},
		{
			desc:          "manageSellOffer - with claims",
			envelopeXDR:   "AAAAAPrjQnnOn4RqMmOSDwYfEMVtJuC4VR9fKvPfEtM7DS7VAAAAZAAMDl8AAAADAAAAAAAAAAAAAAABAAAAAAAAAAMAAAAAAAAAAVNUUgAAAAAASYK2XlJiUiNav1waFVDq1fzoualYC4UNFqThKBroJe0AAAACVAvkAAAAAGMAAADIAAAAAAAAAAAAAAAAAAAAATsNLtUAAABABmA0aLobgdSrjIrus94Y8PWeD6dDfl7Sya12t2uZasJFI7mZ+yowE1enUMzC/cAhDTypK8QuH2EVXPQC3xpYDA==",
			resultXDR:     "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAADAAAAAAAAAAEAAAAADkfaGg9y56NND7n4CRcr4R4fvivwAcMd4ZrCm4jAe5AAAAAAAI0f+AAAAAFTVFIAAAAAAEmCtl5SYlIjWr9cGhVQ6tX86LmpWAuFDRak4Sga6CXtAAAAAS0Il1oAAAAAAAAAAlQL4/8AAAACAAAAAA==",
			metaXDR:       "AAAAAQAAAAIAAAADAAxMfwAAAAAAAAAA+uNCec6fhGoyY5IPBh8QxW0m4LhVH18q898S0zsNLtUAAAAU9GsC1QAMDl8AAAACAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAxMfwAAAAAAAAAA+uNCec6fhGoyY5IPBh8QxW0m4LhVH18q898S0zsNLtUAAAAU9GsC1QAMDl8AAAADAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAACgAAAAMADEx+AAAAAgAAAAAOR9oaD3Lno00PufgJFyvhHh++K/ABwx3hmsKbiMB7kAAAAAAAjR/4AAAAAVNUUgAAAAAASYK2XlJiUiNav1waFVDq1fzoualYC4UNFqThKBroJe0AAAAAAAAAA2L6BdYAAABjAAAAMgAAAAAAAAAAAAAAAAAAAAEADEx/AAAAAgAAAAAOR9oaD3Lno00PufgJFyvhHh++K/ABwx3hmsKbiMB7kAAAAAAAjR/4AAAAAVNUUgAAAAAASYK2XlJiUiNav1waFVDq1fzoualYC4UNFqThKBroJe0AAAAAAAAAAjXxbnwAAABjAAAAMgAAAAAAAAAAAAAAAAAAAAMADEx+AAAAAAAAAAAOR9oaD3Lno00PufgJFyvhHh++K/ABwx3hmsKbiMB7kAAAABnMMdMvAAwOZQAAAAIAAAACAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAABAAAABrSdIAkAAAAAAAAAAAAAAAAAAAAAAAAAAQAMTH8AAAAAAAAAAA5H2hoPcuejTQ+5+AkXK+EeH74r8AHDHeGawpuIwHuQAAAAHCA9ty4ADA5lAAAAAgAAAAIAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAEAAAAEYJE8CgAAAAAAAAAAAAAAAAAAAAAAAAADAAxMfgAAAAEAAAAADkfaGg9y56NND7n4CRcr4R4fvivwAcMd4ZrCm4jAe5AAAAABU1RSAAAAAABJgrZeUmJSI1q/XBoVUOrV/Oi5qVgLhQ0WpOEoGugl7QAAABYDWSXWf/////////8AAAABAAAAAQAAAAAAAAAAAAAAA2L6BdYAAAAAAAAAAAAAAAEADEx/AAAAAQAAAAAOR9oaD3Lno00PufgJFyvhHh++K/ABwx3hmsKbiMB7kAAAAAFTVFIAAAAAAEmCtl5SYlIjWr9cGhVQ6tX86LmpWAuFDRak4Sga6CXtAAAAFNZQjnx//////////wAAAAEAAAABAAAAAAAAAAAAAAACNfFufAAAAAAAAAAAAAAAAwAMDnEAAAABAAAAAPrjQnnOn4RqMmOSDwYfEMVtJuC4VR9fKvPfEtM7DS7VAAAAAVNUUgAAAAAASYK2XlJiUiNav1waFVDq1fzoualYC4UNFqThKBroJe0AAAAYdX9/Wn//////////AAAAAQAAAAAAAAAAAAAAAQAMTH8AAAABAAAAAPrjQnnOn4RqMmOSDwYfEMVtJuC4VR9fKvPfEtM7DS7VAAAAAVNUUgAAAAAASYK2XlJiUiNav1waFVDq1fzoualYC4UNFqThKBroJe0AAAAZoogWtH//////////AAAAAQAAAAAAAAAAAAAAAwAMTH8AAAAAAAAAAPrjQnnOn4RqMmOSDwYfEMVtJuC4VR9fKvPfEtM7DS7VAAAAFPRrAtUADA5fAAAAAwAAAAEAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAQAMTH8AAAAAAAAAAPrjQnnOn4RqMmOSDwYfEMVtJuC4VR9fKvPfEtM7DS7VAAAAEqBfHtYADA5fAAAAAwAAAAEAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAA",
			feeChangesXDR: "AAAAAgAAAAMADA5xAAAAAAAAAAD640J5zp+EajJjkg8GHxDFbSbguFUfXyrz3xLTOw0u1QAAABT0awM5AAwOXwAAAAIAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEADEx/AAAAAAAAAAD640J5zp+EajJjkg8GHxDFbSbguFUfXyrz3xLTOw0u1QAAABT0awLVAAwOXwAAAAIAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			hash:          "ef62da32b6b3eb3c4534dac2be1088387fb93b0093b47e113073c1431fac9db7",
			index:         0,
			sequence:      56,
			expected: []effect{
				{
					address: "GD5OGQTZZ2PYI2RSMOJA6BQ7CDCW2JXAXBKR6XZK6PPRFUZ3BUXNLFKP",
					details: map[string]interface{}{
						"seller":              "GAHEPWQ2B5ZOPI2NB647QCIXFPQR4H56FPYADQY54GNMFG4IYB5ZAJ5H",
						"offer_id":            xdr.Int64(9248760),
						"sold_amount":         "999.9999999",
						"bought_amount":       "505.0505050",
						"sold_asset_type":     "native",
						"bought_asset_code":   "STR",
						"bought_asset_type":   "credit_alphanum4",
						"bought_asset_issuer": "GBEYFNS6KJRFEI22X5OBUFKQ5LK7Z2FZVFMAXBINC2SOCKA25AS62PUN",
					},
					effectType:  history.EffectTrade,
					operationID: int64(240518172673),
					order:       uint32(1),
				},
				{
					address: "GAHEPWQ2B5ZOPI2NB647QCIXFPQR4H56FPYADQY54GNMFG4IYB5ZAJ5H",
					details: map[string]interface{}{
						"seller":            "GD5OGQTZZ2PYI2RSMOJA6BQ7CDCW2JXAXBKR6XZK6PPRFUZ3BUXNLFKP",
						"offer_id":          xdr.Int64(9248760),
						"sold_amount":       "505.0505050",
						"bought_amount":     "999.9999999",
						"sold_asset_code":   "STR",
						"sold_asset_type":   "credit_alphanum4",
						"bought_asset_type": "native",
						"sold_asset_issuer": "GBEYFNS6KJRFEI22X5OBUFKQ5LK7Z2FZVFMAXBINC2SOCKA25AS62PUN",
					},
					effectType:  history.EffectTrade,
					operationID: int64(240518172673),
					order:       uint32(2),
				},
			},
		},
		{
			desc:          "manageBuyOffer - with claims",
			envelopeXDR:   "AAAAAEotqBM9oOzudkkctgQlY/PHS0rFcxVasWQVnSytiuBEAAAAZAANIfEAAAADAAAAAAAAAAAAAAABAAAAAAAAAAwAAAAAAAAAAlRYVGFscGhhNAAAAAAAAABKLagTPaDs7nZJHLYEJWPzx0tKxXMVWrFkFZ0srYrgRAAAAAB3NZQAAAAAAQAAAAEAAAAAAAAAAAAAAAAAAAABrYrgRAAAAEAh57TBifjJuUPj1TI7zIvaAZmyRjWLY4ktc0F16Knmy4Fw07L7cC5vCwjn4ZXyrgr9bpEGhv4oN6znbPpNLQUH",
			resultXDR:     "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAMAAAAAAAAAAEAAAAAgbI9jY68fYXd6+DwMcZQQIYCK4HsKKvqnR5o+1IdVoUAAAAAAJovcgAAAAJUWFRhbHBoYTQAAAAAAAAASi2oEz2g7O52SRy2BCVj88dLSsVzFVqxZBWdLK2K4EQAAAAAdzWUAAAAAAAAAAAAdzWUAAAAAAIAAAAA",
			metaXDR:       "AAAAAQAAAAIAAAADAA0pGAAAAAAAAAAASi2oEz2g7O52SRy2BCVj88dLSsVzFVqxZBWdLK2K4EQAAAAXSHbm1AANIfEAAAACAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAA0pGAAAAAAAAAAASi2oEz2g7O52SRy2BCVj88dLSsVzFVqxZBWdLK2K4EQAAAAXSHbm1AANIfEAAAADAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAACAAAAAMADSkYAAAAAAAAAABKLagTPaDs7nZJHLYEJWPzx0tKxXMVWrFkFZ0srYrgRAAAABdIdubUAA0h8QAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEADSkYAAAAAAAAAABKLagTPaDs7nZJHLYEJWPzx0tKxXMVWrFkFZ0srYrgRAAAABbRQVLUAA0h8QAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMADSjEAAAAAgAAAACBsj2Njrx9hd3r4PAxxlBAhgIrgewoq+qdHmj7Uh1WhQAAAAAAmi9yAAAAAlRYVGFscGhhNAAAAAAAAABKLagTPaDs7nZJHLYEJWPzx0tKxXMVWrFkFZ0srYrgRAAAAAAAAAAAstBeAAAAAAEAAAABAAAAAAAAAAAAAAAAAAAAAQANKRgAAAACAAAAAIGyPY2OvH2F3evg8DHGUECGAiuB7Cir6p0eaPtSHVaFAAAAAACaL3IAAAACVFhUYWxwaGE0AAAAAAAAAEotqBM9oOzudkkctgQlY/PHS0rFcxVasWQVnSytiuBEAAAAAAAAAAA7msoAAAAAAQAAAAEAAAAAAAAAAAAAAAAAAAADAA0oxAAAAAAAAAAAgbI9jY68fYXd6+DwMcZQQIYCK4HsKKvqnR5o+1IdVoUAAAAZJU0xXAANGSMAAAARAAAABAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAQADMowLgdQAAAAAAAAAAAAAAAAAAAAAAAAAAAEADSkYAAAAAAAAAACBsj2Njrx9hd3r4PAxxlBAhgIrgewoq+qdHmj7Uh1WhQAAABmcgsVcAA0ZIwAAABEAAAAEAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAABAAMyi5RMQAAAAAAAAAAAAAAAAAAAAAAAAAAAAwANKMQAAAABAAAAAIGyPY2OvH2F3evg8DHGUECGAiuB7Cir6p0eaPtSHVaFAAAAAlRYVGFscGhhNAAAAAAAAABKLagTPaDs7nZJHLYEJWPzx0tKxXMVWrFkFZ0srYrgRAAACRatNxoAf/////////8AAAABAAAAAQAAAAAAAAAAAAAAALLQXgAAAAAAAAAAAAAAAAEADSkYAAAAAQAAAACBsj2Njrx9hd3r4PAxxlBAhgIrgewoq+qdHmj7Uh1WhQAAAAJUWFRhbHBoYTQAAAAAAAAASi2oEz2g7O52SRy2BCVj88dLSsVzFVqxZBWdLK2K4EQAAAkWNgGGAH//////////AAAAAQAAAAEAAAAAAAAAAAAAAAA7msoAAAAAAAAAAAA=",
			feeChangesXDR: "AAAAAgAAAAMADSSgAAAAAAAAAABKLagTPaDs7nZJHLYEJWPzx0tKxXMVWrFkFZ0srYrgRAAAABdIduc4AA0h8QAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEADSkYAAAAAAAAAABKLagTPaDs7nZJHLYEJWPzx0tKxXMVWrFkFZ0srYrgRAAAABdIdubUAA0h8QAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			hash:          "9caa91eec6e29730f4aabafb60898a8ecedd3bf67b8628e6e32066fbba9bec5d",
			index:         0,
			sequence:      56,
			expected: []effect{
				{
					address: "GBFC3KATHWQOZ3TWJEOLMBBFMPZ4OS2KYVZRKWVRMQKZ2LFNRLQEIRCV",
					details: map[string]interface{}{
						"seller":              "GCA3EPMNR26H3BO55PQPAMOGKBAIMARLQHWCRK7KTUPGR62SDVLIL7D6",
						"offer_id":            xdr.Int64(10104690),
						"sold_amount":         "200.0000000",
						"bought_amount":       "200.0000000",
						"sold_asset_type":     "native",
						"bought_asset_code":   "TXTalpha4",
						"bought_asset_type":   "credit_alphanum12",
						"bought_asset_issuer": "GBFC3KATHWQOZ3TWJEOLMBBFMPZ4OS2KYVZRKWVRMQKZ2LFNRLQEIRCV",
					},
					effectType:  history.EffectTrade,
					operationID: int64(240518172673),
					order:       uint32(1),
				},
				{
					address: "GCA3EPMNR26H3BO55PQPAMOGKBAIMARLQHWCRK7KTUPGR62SDVLIL7D6",
					details: map[string]interface{}{
						"seller":            "GBFC3KATHWQOZ3TWJEOLMBBFMPZ4OS2KYVZRKWVRMQKZ2LFNRLQEIRCV",
						"offer_id":          xdr.Int64(10104690),
						"sold_amount":       "200.0000000",
						"bought_amount":     "200.0000000",
						"sold_asset_code":   "TXTalpha4",
						"sold_asset_type":   "credit_alphanum12",
						"bought_asset_type": "native",
						"sold_asset_issuer": "GBFC3KATHWQOZ3TWJEOLMBBFMPZ4OS2KYVZRKWVRMQKZ2LFNRLQEIRCV",
					},
					effectType:  history.EffectTrade,
					operationID: int64(240518172673),
					order:       uint32(2),
				},
			},
		},
		{
			desc:          "createPassiveSellOffer",
			envelopeXDR:   "AAAAAAHwZwJPu1TJhQGgsLRXBzcIeySkeGXzEqh0W9AHWvFDAAAAZAAN3tMAAAACAAAAAQAAAAAAAAAAAAAAAF4FBqwAAAAAAAAAAQAAAAAAAAAEAAAAAAAAAAFDT1AAAAAAALly/iTceP/82O3aZAmd8hyqUjYAANfc5RfN0/iibCtTAAAAADuaygAAAAAJAAAACgAAAAAAAAABB1rxQwAAAEDz2JIw8Z3Owoc5c2tsiY3kzOYUmh32155u00Xs+RYxO5fL0ApYd78URHcYCbe0R32YmuLTfefWQStR3RfhqKAL",
			resultXDR:     "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAADAAAAAAAAAAEAAAAAMgQ65fmCczzuwmU3oQLivASzvZdhzjhJOQ6C+xTSDu8AAAAAAKMvZgAAAAFDT1AAAAAAALly/iTceP/82O3aZAmd8hyqUjYAANfc5RfN0/iibCtTAAAA6NSlEAAAAAAAAAAAADuaygAAAAACAAAAAA==",
			metaXDR:       "AAAAAQAAAAIAAAADAA3fGgAAAAAAAAAAAfBnAk+7VMmFAaCwtFcHNwh7JKR4ZfMSqHRb0Ada8UMAAAAXSHbnOAAN3tMAAAABAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAA3fGgAAAAAAAAAAAfBnAk+7VMmFAaCwtFcHNwh7JKR4ZfMSqHRb0Ada8UMAAAAXSHbnOAAN3tMAAAACAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAACgAAAAMADd72AAAAAgAAAAAyBDrl+YJzPO7CZTehAuK8BLO9l2HOOEk5DoL7FNIO7wAAAAAAoy9mAAAAAUNPUAAAAAAAuXL+JNx4//zY7dpkCZ3yHKpSNgAA19zlF83T+KJsK1MAAAAAAAAA6NSlEAAAAAABAAAD6AAAAAAAAAAAAAAAAAAAAAIAAAACAAAAADIEOuX5gnM87sJlN6EC4rwEs72XYc44STkOgvsU0g7vAAAAAACjL2YAAAADAA3fGQAAAAAAAAAAMgQ65fmCczzuwmU3oQLivASzvZdhzjhJOQ6C+xTSDu8AAAAXSHbkfAAIGHsAAAAJAAAAAwAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAQAAAAB3NZQAAAAAAAAAAAAAAAAAAAAAAAAAAAEADd8aAAAAAAAAAAAyBDrl+YJzPO7CZTehAuK8BLO9l2HOOEk5DoL7FNIO7wAAABeEEa58AAgYewAAAAkAAAACAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAABAAAAADuaygAAAAAAAAAAAAAAAAAAAAAAAAAAAwAN3xkAAAABAAAAADIEOuX5gnM87sJlN6EC4rwEs72XYc44STkOgvsU0g7vAAAAAUNPUAAAAAAAuXL+JNx4//zY7dpkCZ3yHKpSNgAA19zlF83T+KJsK1MAABI3mQjsAH//////////AAAAAQAAAAEAAAAAAAAAAAAAAdGpSiAAAAAAAAAAAAAAAAABAA3fGgAAAAEAAAAAMgQ65fmCczzuwmU3oQLivASzvZdhzjhJOQ6C+xTSDu8AAAABQ09QAAAAAAC5cv4k3Hj//Njt2mQJnfIcqlI2AADX3OUXzdP4omwrUwAAEU7EY9wAf/////////8AAAABAAAAAQAAAAAAAAAAAAAA6NSlEAAAAAAAAAAAAAAAAAMADd7UAAAAAQAAAAAB8GcCT7tUyYUBoLC0Vwc3CHskpHhl8xKodFvQB1rxQwAAAAFDT1AAAAAAALly/iTceP/82O3aZAmd8hyqUjYAANfc5RfN0/iibCtTAAAAAAAAAAB//////////wAAAAEAAAAAAAAAAAAAAAEADd8aAAAAAQAAAAAB8GcCT7tUyYUBoLC0Vwc3CHskpHhl8xKodFvQB1rxQwAAAAFDT1AAAAAAALly/iTceP/82O3aZAmd8hyqUjYAANfc5RfN0/iibCtTAAAA6NSlEAB//////////wAAAAEAAAAAAAAAAAAAAAMADd8aAAAAAAAAAAAB8GcCT7tUyYUBoLC0Vwc3CHskpHhl8xKodFvQB1rxQwAAABdIduc4AA3e0wAAAAIAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEADd8aAAAAAAAAAAAB8GcCT7tUyYUBoLC0Vwc3CHskpHhl8xKodFvQB1rxQwAAABcM3B04AA3e0wAAAAIAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			feeChangesXDR: "AAAAAgAAAAMADd7UAAAAAAAAAAAB8GcCT7tUyYUBoLC0Vwc3CHskpHhl8xKodFvQB1rxQwAAABdIduecAA3e0wAAAAEAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEADd8aAAAAAAAAAAAB8GcCT7tUyYUBoLC0Vwc3CHskpHhl8xKodFvQB1rxQwAAABdIduc4AA3e0wAAAAEAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			hash:          "e4b286344ae1c863ab15773ddf6649b08fe031383135194f8613a3a475c41a5a",
			index:         0,
			sequence:      56,
			expected: []effect{
				{
					address: "GAA7AZYCJ65VJSMFAGQLBNCXA43QQ6ZEUR4GL4YSVB2FXUAHLLYUHIO5",
					details: map[string]interface{}{
						"bought_amount":       "100000.0000000",
						"bought_asset_code":   "COP",
						"bought_asset_issuer": "GC4XF7RE3R4P77GY5XNGICM56IOKUURWAAANPXHFC7G5H6FCNQVVH3OH",
						"bought_asset_type":   "credit_alphanum4",
						"offer_id":            xdr.Int64(10694502),
						"seller":              "GAZAIOXF7GBHGPHOYJSTPIIC4K6AJM55S5Q44OCJHEHIF6YU2IHO6VHU",
						"sold_amount":         "100.0000000",
						"sold_asset_type":     "native",
					},
					effectType:  history.EffectTrade,
					operationID: int64(240518172673),
					order:       uint32(1),
				},
				{
					address: "GAZAIOXF7GBHGPHOYJSTPIIC4K6AJM55S5Q44OCJHEHIF6YU2IHO6VHU",
					details: map[string]interface{}{
						"bought_amount":     "100.0000000",
						"bought_asset_type": "native",
						"offer_id":          xdr.Int64(10694502),
						"seller":            "GAA7AZYCJ65VJSMFAGQLBNCXA43QQ6ZEUR4GL4YSVB2FXUAHLLYUHIO5",
						"sold_amount":       "100000.0000000",
						"sold_asset_code":   "COP",
						"sold_asset_issuer": "GC4XF7RE3R4P77GY5XNGICM56IOKUURWAAANPXHFC7G5H6FCNQVVH3OH",
						"sold_asset_type":   "credit_alphanum4",
					},
					effectType:  history.EffectTrade,
					operationID: int64(240518172673),
					order:       uint32(2),
				},
			},
		},
		{
			desc:          "setOption",
			envelopeXDR:   "AAAAALly/iTceP/82O3aZAmd8hyqUjYAANfc5RfN0/iibCtTAAAAZAAIGHoAAAAHAAAAAQAAAAAAAAAAAAAAAF4FFtcAAAAAAAAAAQAAAAAAAAAFAAAAAQAAAAAge0MBDbX9OddsGMWIHbY1cGXuGYP4bl1ylIvUklO73AAAAAEAAAACAAAAAQAAAAEAAAABAAAAAwAAAAEAAAABAAAAAQAAAAIAAAABAAAAAwAAAAEAAAAVaHR0cHM6Ly93d3cuaG9tZS5vcmcvAAAAAAAAAQAAAAAge0MBDbX9OddsGMWIHbY1cGXuGYP4bl1ylIvUklO73AAAAAIAAAAAAAAAAaJsK1MAAABAiQjCxE53GjInjJtvNr6gdhztRi0GWOZKlUS2KZBLjX3n2N/y7RRNt7B1ZuFcZAxrnxWHD/fF2XcrEwFAuf4TDA==",
			resultXDR:     "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=",
			metaXDR:       "AAAAAQAAAAIAAAADAA3iDQAAAAAAAAAAuXL+JNx4//zY7dpkCZ3yHKpSNgAA19zlF83T+KJsK1MAAAAXSHblRAAIGHoAAAAGAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAA3iDQAAAAAAAAAAuXL+JNx4//zY7dpkCZ3yHKpSNgAA19zlF83T+KJsK1MAAAAXSHblRAAIGHoAAAAHAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMADeINAAAAAAAAAAC5cv4k3Hj//Njt2mQJnfIcqlI2AADX3OUXzdP4omwrUwAAABdIduVEAAgYegAAAAcAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEADeINAAAAAAAAAAC5cv4k3Hj//Njt2mQJnfIcqlI2AADX3OUXzdP4omwrUwAAABdIduVEAAgYegAAAAcAAAABAAAAAQAAAAAge0MBDbX9OddsGMWIHbY1cGXuGYP4bl1ylIvUklO73AAAAAEAAAAVaHR0cHM6Ly93d3cuaG9tZS5vcmcvAAAAAwECAwAAAAEAAAAAIHtDAQ21/TnXbBjFiB22NXBl7hmD+G5dcpSL1JJTu9wAAAACAAAAAAAAAAA=",
			feeChangesXDR: "AAAAAgAAAAMADd8YAAAAAAAAAAC5cv4k3Hj//Njt2mQJnfIcqlI2AADX3OUXzdP4omwrUwAAABdIduWoAAgYegAAAAYAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEADeINAAAAAAAAAAC5cv4k3Hj//Njt2mQJnfIcqlI2AADX3OUXzdP4omwrUwAAABdIduVEAAgYegAAAAYAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			hash:          "e76b7b0133690fbfb2de8fa9ca2273cb4f2e29447e0cf0e14a5f82d0daa48760",
			index:         0,
			sequence:      56,
			expected: []effect{
				{
					address: "GC4XF7RE3R4P77GY5XNGICM56IOKUURWAAANPXHFC7G5H6FCNQVVH3OH",
					details: map[string]interface{}{
						"home_domain": "https://www.home.org/",
					},
					effectType:  history.EffectAccountHomeDomainUpdated,
					operationID: int64(240518172673),
					order:       uint32(1),
				},
				{
					address: "GC4XF7RE3R4P77GY5XNGICM56IOKUURWAAANPXHFC7G5H6FCNQVVH3OH",
					details: map[string]interface{}{
						"high_threshold": xdr.Uint32(3),
						"low_threshold":  xdr.Uint32(1),
						"med_threshold":  xdr.Uint32(2),
					},
					effectType:  history.EffectAccountThresholdsUpdated,
					operationID: int64(240518172673),
					order:       uint32(2),
				},
				{
					address: "GC4XF7RE3R4P77GY5XNGICM56IOKUURWAAANPXHFC7G5H6FCNQVVH3OH",
					details: map[string]interface{}{
						"auth_required_flag":  true,
						"auth_revocable_flag": false,
					},
					effectType:  history.EffectAccountFlagsUpdated,
					operationID: int64(240518172673),
					order:       uint32(3),
				},
				{
					address: "GC4XF7RE3R4P77GY5XNGICM56IOKUURWAAANPXHFC7G5H6FCNQVVH3OH",
					details: map[string]interface{}{
						"inflation_destination": "GAQHWQYBBW272OOXNQMMLCA5WY2XAZPODGB7Q3S5OKKIXVESKO55ZQ7C",
					},
					effectType:  history.EffectAccountInflationDestinationUpdated,
					operationID: int64(240518172673),
					order:       uint32(4),
				},
				{
					address: "GC4XF7RE3R4P77GY5XNGICM56IOKUURWAAANPXHFC7G5H6FCNQVVH3OH",
					details: map[string]interface{}{
						"public_key": "GC4XF7RE3R4P77GY5XNGICM56IOKUURWAAANPXHFC7G5H6FCNQVVH3OH",
						"weight":     int32(3),
					},
					effectType:  history.EffectSignerUpdated,
					operationID: int64(240518172673),
					order:       uint32(5),
				},
				{
					address: "GC4XF7RE3R4P77GY5XNGICM56IOKUURWAAANPXHFC7G5H6FCNQVVH3OH",
					details: map[string]interface{}{
						"public_key": "GAQHWQYBBW272OOXNQMMLCA5WY2XAZPODGB7Q3S5OKKIXVESKO55ZQ7C",
						"weight":     int32(2),
					},
					effectType:  history.EffectSignerCreated,
					operationID: int64(240518172673),
					order:       uint32(6),
				},
			},
		},
		{
			desc:          "changeTrust - trustline created",
			envelopeXDR:   "AAAAAKturFHJX/eRt5gM6qIXAMbaXvlImqLysA6Qr9tLemxfAAAAZAAAACYAAAABAAAAAAAAAAAAAAABAAAAAAAAAAYAAAABVVNEAAAAAAD5Jjibq+Rf5jsUyQ2/tGzCwiRg0Zd5nj9jARA1Skjz+H//////////AAAAAAAAAAFLemxfAAAAQKN8LftAafeoAGmvpsEokqm47jAuqw4g1UWjmL0j6QPm1jxoalzDwDS3W+N2HOHdjSJlEQaTxGBfQKHhr6nNsAA=",
			resultXDR:     "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=",
			metaXDR:       "AAAAAQAAAAIAAAADAAAAKAAAAAAAAAAAq26sUclf95G3mAzqohcAxtpe+UiaovKwDpCv20t6bF8AAAACVAvjOAAAACYAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAKAAAAAAAAAAAq26sUclf95G3mAzqohcAxtpe+UiaovKwDpCv20t6bF8AAAACVAvjOAAAACYAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAoAAAAAAAAAACrbqxRyV/3kbeYDOqiFwDG2l75SJqi8rAOkK/bS3psXwAAAAJUC+M4AAAAJgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAoAAAAAAAAAACrbqxRyV/3kbeYDOqiFwDG2l75SJqi8rAOkK/bS3psXwAAAAJUC+M4AAAAJgAAAAEAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAoAAAAAQAAAACrbqxRyV/3kbeYDOqiFwDG2l75SJqi8rAOkK/bS3psXwAAAAFVU0QAAAAAAPkmOJur5F/mOxTJDb+0bMLCJGDRl3meP2MBEDVKSPP4AAAAAAAAAAB//////////wAAAAAAAAAAAAAAAA==",
			feeChangesXDR: "AAAAAgAAAAMAAAAmAAAAAAAAAACrbqxRyV/3kbeYDOqiFwDG2l75SJqi8rAOkK/bS3psXwAAAAJUC+QAAAAAJgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAoAAAAAAAAAACrbqxRyV/3kbeYDOqiFwDG2l75SJqi8rAOkK/bS3psXwAAAAJUC+OcAAAAJgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			hash:          "6fa467b53f5386d77ad35c2502ed2cd3dd8b460a5be22b6b2818b81bcd3ed2da",
			index:         0,
			sequence:      40,
			expected: []effect{
				{
					address:     "GCVW5LCRZFP7PENXTAGOVIQXADDNUXXZJCNKF4VQB2IK7W2LPJWF73UG",
					effectType:  history.EffectTrustlineCreated,
					operationID: int64(171798695937),
					order:       uint32(1),
					details: map[string]interface{}{
						"limit":        "922337203685.4775807",
						"asset_code":   "USD",
						"asset_type":   "credit_alphanum4",
						"asset_issuer": "GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF",
					},
				},
			},
		},
		{
			desc:          "changeTrust - trustline removed",
			envelopeXDR:   "AAAAABwDSftLnTVAHpKUGYPZfTJr6rIm5Z5IqDHVBFuTI3ubAAAAZAARM9kAAAADAAAAAQAAAAAAAAAAAAAAAF4XMm8AAAAAAAAAAQAAAAAAAAAGAAAAAk9DSVRva2VuAAAAAAAAAABJxf/HoI4oaD9CLBvECRhG9GPMNa/65PTI9N7F37o4nwAAAAAAAAAAAAAAAAAAAAGTI3ubAAAAQMHTFPeyHA+W2EYHVDut4dQ18zvF+47SsTPaePwZUaCgw/A3tKDx7sO7R8xlI3GwKQl91Ljmm1dbvAONU9nk/AQ=",
			resultXDR:     "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAGAAAAAAAAAAA=",
			metaXDR:       "AAAAAQAAAAIAAAADABEz3wAAAAAAAAAAHANJ+0udNUAekpQZg9l9MmvqsiblnkioMdUEW5Mje5sAAAAXSHbm1AARM9kAAAACAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABABEz3wAAAAAAAAAAHANJ+0udNUAekpQZg9l9MmvqsiblnkioMdUEW5Mje5sAAAAXSHbm1AARM9kAAAADAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAETPeAAAAAQAAAAAcA0n7S501QB6SlBmD2X0ya+qyJuWeSKgx1QRbkyN7mwAAAAJPQ0lUb2tlbgAAAAAAAAAAScX/x6COKGg/QiwbxAkYRvRjzDWv+uT0yPTexd+6OJ8AAAAAAAAAAH//////////AAAAAQAAAAAAAAAAAAAAAgAAAAEAAAAAHANJ+0udNUAekpQZg9l9MmvqsiblnkioMdUEW5Mje5sAAAACT0NJVG9rZW4AAAAAAAAAAEnF/8egjihoP0IsG8QJGEb0Y8w1r/rk9Mj03sXfujifAAAAAwARM98AAAAAAAAAABwDSftLnTVAHpKUGYPZfTJr6rIm5Z5IqDHVBFuTI3ubAAAAF0h25tQAETPZAAAAAwAAAAEAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAQARM98AAAAAAAAAABwDSftLnTVAHpKUGYPZfTJr6rIm5Z5IqDHVBFuTI3ubAAAAF0h25tQAETPZAAAAAwAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAA",
			feeChangesXDR: "AAAAAgAAAAMAETPeAAAAAAAAAAAcA0n7S501QB6SlBmD2X0ya+qyJuWeSKgx1QRbkyN7mwAAABdIduc4ABEz2QAAAAIAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAETPfAAAAAAAAAAAcA0n7S501QB6SlBmD2X0ya+qyJuWeSKgx1QRbkyN7mwAAABdIdubUABEz2QAAAAIAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			hash:          "0f1e93ed9a83edb01ad8ccab67fd59dc7a513c413a8d5a580c5eb7a9c44f2844",
			index:         0,
			sequence:      40,
			expected: []effect{
				{
					address:     "GAOAGSP3JOOTKQA6SKKBTA6ZPUZGX2VSE3SZ4SFIGHKQIW4TEN5ZX3WW",
					effectType:  history.EffectTrustlineRemoved,
					operationID: int64(171798695937),
					order:       uint32(1),
					details: map[string]interface{}{
						"limit":        "0.0000000",
						"asset_code":   "OCIToken",
						"asset_type":   "credit_alphanum12",
						"asset_issuer": "GBE4L76HUCHCQ2B7IIWBXRAJDBDPIY6MGWX7VZHUZD2N5RO7XI4J6GTJ",
					},
				},
			},
		},
		{
			desc:          "changeTrust - trustline updated",
			envelopeXDR:   "AAAAAHHbEhVipyZ2k4byyCZkS1Bdvpj7faBChuYo8S/Rt89UAAAAZAAQuJIAAAAHAAAAAQAAAAAAAAAAAAAAAF4XVskAAAAAAAAAAQAAAAAAAAAGAAAAAlRFU1RBU1NFVAAAAAAAAAA7JUkkD+tgCi2xTVyEcs4WZXOA0l7w2orZg/bghXOgkAAAAAA7msoAAAAAAAAAAAHRt89UAAAAQOCi2ylqRvvRzZaCFjGkLYFk7DCjJA5uZ1nXo8FaPCRl2LZczoMbc46sZIlHh0ENzk7fKjFnRPMo8XAirrrf2go=",
			resultXDR:     "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAGAAAAAAAAAAA=",
			metaXDR:       "AAAAAQAAAAIAAAADABE6jwAAAAAAAAAAcdsSFWKnJnaThvLIJmRLUF2+mPt9oEKG5ijxL9G3z1QAAAAAO5rHRAAQuJIAAAAGAAAAAgAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABABE6jwAAAAAAAAAAcdsSFWKnJnaThvLIJmRLUF2+mPt9oEKG5ijxL9G3z1QAAAAAO5rHRAAQuJIAAAAHAAAAAgAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAETqAAAAAAQAAAABx2xIVYqcmdpOG8sgmZEtQXb6Y+32gQobmKPEv0bfPVAAAAAJURVNUQVNTRVQAAAAAAAAAOyVJJA/rYAotsU1chHLOFmVzgNJe8NqK2YP24IVzoJAAAAAAO5rKAAAAAAA7msoAAAAAAQAAAAAAAAAAAAAAAQAROo8AAAABAAAAAHHbEhVipyZ2k4byyCZkS1Bdvpj7faBChuYo8S/Rt89UAAAAAlRFU1RBU1NFVAAAAAAAAAA7JUkkD+tgCi2xTVyEcs4WZXOA0l7w2orZg/bghXOgkAAAAAA7msoAAAAAADuaygAAAAABAAAAAAAAAAA=",
			feeChangesXDR: "AAAAAgAAAAMAETp/AAAAAAAAAABx2xIVYqcmdpOG8sgmZEtQXb6Y+32gQobmKPEv0bfPVAAAAAA7mseoABC4kgAAAAYAAAACAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAETqPAAAAAAAAAABx2xIVYqcmdpOG8sgmZEtQXb6Y+32gQobmKPEv0bfPVAAAAAA7msdEABC4kgAAAAYAAAACAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			hash:          "dc8d4714d7db3d0e27ae07f629bc72f1605fc24a2d178af04edbb602592791aa",
			index:         0,
			sequence:      40,
			expected: []effect{
				{
					address:     "GBY5WEQVMKTSM5UTQ3ZMQJTEJNIF3PUY7N62AQUG4YUPCL6RW7HVJARI",
					effectType:  history.EffectTrustlineUpdated,
					operationID: int64(171798695937),
					order:       uint32(1),
					details: map[string]interface{}{
						"limit":        "100.0000000",
						"asset_code":   "TESTASSET",
						"asset_type":   "credit_alphanum12",
						"asset_issuer": "GA5SKSJEB7VWACRNWFGVZBDSZYLGK44A2JPPBWUK3GB7NYEFOOQJAC2B",
					},
				},
			},
		},
		{
			desc:          "allowTrust",
			envelopeXDR:   "AAAAAPkmOJur5F/mOxTJDb+0bMLCJGDRl3meP2MBEDVKSPP4AAAAZAAAACYAAAACAAAAAAAAAAAAAAABAAAAAAAAAAcAAAAAq26sUclf95G3mAzqohcAxtpe+UiaovKwDpCv20t6bF8AAAABVVNEAAAAAAEAAAAAAAAAAUpI8/gAAABA6O2fe1gQBwoO0fMNNEUKH0QdVXVjEWbN5VL51DmRUedYMMXtbX5JKVSzla2kIGvWgls1dXuXHZY/IOlaK01rBQ==",
			resultXDR:     "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=",
			metaXDR:       "AAAAAQAAAAIAAAADAAAAKQAAAAAAAAAA+SY4m6vkX+Y7FMkNv7RswsIkYNGXeZ4/YwEQNUpI8/gAAAACVAvi1AAAACYAAAABAAAAAAAAAAAAAAADAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAKQAAAAAAAAAA+SY4m6vkX+Y7FMkNv7RswsIkYNGXeZ4/YwEQNUpI8/gAAAACVAvi1AAAACYAAAACAAAAAAAAAAAAAAADAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAAoAAAAAQAAAACrbqxRyV/3kbeYDOqiFwDG2l75SJqi8rAOkK/bS3psXwAAAAFVU0QAAAAAAPkmOJur5F/mOxTJDb+0bMLCJGDRl3meP2MBEDVKSPP4AAAAAAAAAAB//////////wAAAAAAAAAAAAAAAAAAAAEAAAApAAAAAQAAAACrbqxRyV/3kbeYDOqiFwDG2l75SJqi8rAOkK/bS3psXwAAAAFVU0QAAAAAAPkmOJur5F/mOxTJDb+0bMLCJGDRl3meP2MBEDVKSPP4AAAAAAAAAAB//////////wAAAAEAAAAAAAAAAA==",
			feeChangesXDR: "AAAAAgAAAAMAAAAnAAAAAAAAAAD5Jjibq+Rf5jsUyQ2/tGzCwiRg0Zd5nj9jARA1Skjz+AAAAAJUC+OcAAAAJgAAAAEAAAAAAAAAAAAAAAMAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAApAAAAAAAAAAD5Jjibq+Rf5jsUyQ2/tGzCwiRg0Zd5nj9jARA1Skjz+AAAAAJUC+M4AAAAJgAAAAEAAAAAAAAAAAAAAAMAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			hash:          "6d2e30fd57492bf2e2b132e1bc91a548a369189bebf77eb2b3d829121a9d2c50",
			index:         0,
			sequence:      41,
			expected: []effect{
				{
					address:     "GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF",
					effectType:  history.EffectTrustlineAuthorized,
					operationID: int64(176093663233),
					order:       uint32(1),
					details: map[string]interface{}{
						"trustor":      "GCVW5LCRZFP7PENXTAGOVIQXADDNUXXZJCNKF4VQB2IK7W2LPJWF73UG",
						"asset_code":   "USD",
						"asset_type":   "credit_alphanum4",
						"asset_issuer": "GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF",
					},
				},
				{
					address:     "GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF",
					effectType:  history.EffectTrustlineFlagsUpdated,
					order:       uint32(2),
					operationID: int64(176093663233),
					details: map[string]interface{}{
						"asset_code":      "USD",
						"asset_issuer":    "GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF",
						"asset_type":      "credit_alphanum4",
						"authorized_flag": true,
						"trustor":         "GCVW5LCRZFP7PENXTAGOVIQXADDNUXXZJCNKF4VQB2IK7W2LPJWF73UG",
					},
				},
			},
		},
		{
			desc:          "accountMerge (Destination)",
			envelopeXDR:   "AAAAAI77mqNTy9VPgmgn+//uvjP8VJxJ1FHQ4jCrYS+K4+HvAAAAZAAAACsAAAABAAAAAAAAAAAAAAABAAAAAAAAAAgAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAYrj4e8AAABA3jJ7wBrRpsrcnqBQWjyzwvVz2v5UJ56G60IhgsaWQFSf+7om462KToc+HJ27aLVOQ83dGh1ivp+VIuREJq/SBw==",
			resultXDR:     "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAIAAAAAAAAAAJUC+OcAAAAAA==",
			metaXDR:       "AAAAAQAAAAIAAAADAAAALAAAAAAAAAAAjvuao1PL1U+CaCf7/+6+M/xUnEnUUdDiMKthL4rj4e8AAAACVAvjnAAAACsAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAALAAAAAAAAAAAjvuao1PL1U+CaCf7/+6+M/xUnEnUUdDiMKthL4rj4e8AAAACVAvjnAAAACsAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAArAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtonM3Az4AAAAAAAAABIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAsAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtowg5/CUAAAAAAAAABIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMAAAAsAAAAAAAAAACO+5qjU8vVT4JoJ/v/7r4z/FScSdRR0OIwq2EviuPh7wAAAAJUC+OcAAAAKwAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAIAAAAAAAAAAI77mqNTy9VPgmgn+//uvjP8VJxJ1FHQ4jCrYS+K4+Hv",
			feeChangesXDR: "AAAAAgAAAAMAAAArAAAAAAAAAACO+5qjU8vVT4JoJ/v/7r4z/FScSdRR0OIwq2EviuPh7wAAAAJUC+QAAAAAKwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAsAAAAAAAAAACO+5qjU8vVT4JoJ/v/7r4z/FScSdRR0OIwq2EviuPh7wAAAAJUC+OcAAAAKwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			hash:          "e0773d07aba23d11e6a06b021682294be1f9f202a2926827022539662ce2c7fc",
			index:         0,
			sequence:      44,
			expected: []effect{
				{
					address:     "GCHPXGVDKPF5KT4CNAT7X77OXYZ7YVE4JHKFDUHCGCVWCL4K4PQ67KKZ",
					effectType:  history.EffectAccountDebited,
					operationID: int64(188978565121),
					order:       uint32(1),
					details: map[string]interface{}{
						"amount":     "999.9999900",
						"asset_type": "native",
					},
				},
				{
					address:     "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
					effectType:  history.EffectAccountCredited,
					operationID: int64(188978565121),
					order:       uint32(2),
					details: map[string]interface{}{
						"amount":     "999.9999900",
						"asset_type": "native",
					},
				},
				{
					address:     "GCHPXGVDKPF5KT4CNAT7X77OXYZ7YVE4JHKFDUHCGCVWCL4K4PQ67KKZ",
					effectType:  history.EffectAccountRemoved,
					operationID: int64(188978565121),
					order:       uint32(3),
					details:     map[string]interface{}{},
				},
			},
		},
		{
			desc:          "inflation",
			envelopeXDR:   "AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAAVAAAAAAAAAAAAAAABAAAAAAAAAAkAAAAAAAAAAVb8BfcAAABABUHuXY+MTgW/wDv5+NDVh9fw4meszxeXO98HEQfgXVeCZ7eObCI2orSGUNA/SK6HV9/uTVSxIQQWIso1QoxHBQ==",
			resultXDR:     "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAJAAAAAAAAAAIAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAIrEjCYwXAAAAADj3dgEQp1N5U3fBSOCx/nr5XtiCmNJ2oMJZMx+MYK3JwAAIrEjfceLAAAAAA==",
			metaXDR:       "AAAAAQAAAAIAAAADAAAALwAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcLGiubZdPvaAAAAAAAAAAUAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAALwAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcLGiubZdPvaAAAAAAAAAAVAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAAvAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsaK5tl0+9oAAAAAAAAABUAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAvAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsatl/x+h/EAAAAAAAAABUAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMAAAAuAAAAAAAAAADj3dgEQp1N5U3fBSOCx/nr5XtiCmNJ2oMJZMx+MYK3JwLGivC7E/+cAAAALQAAAAEAAAAAAAAAAQAAAADj3dgEQp1N5U3fBSOCx/nr5XtiCmNJ2oMJZMx+MYK3JwAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAvAAAAAAAAAADj3dgEQp1N5U3fBSOCx/nr5XtiCmNJ2oMJZMx+MYK3JwLGraHekccnAAAALQAAAAEAAAAAAAAAAQAAAADj3dgEQp1N5U3fBSOCx/nr5XtiCmNJ2oMJZMx+MYK3JwAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			feeChangesXDR: "AAAAAgAAAAMAAAAuAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsaK5tl0+/MAAAAAAAAABQAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAvAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsaK5tl0+9oAAAAAAAAABQAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			hash:          "ea93efd8c2f4e45c0318c69ec958623a0e4374f40d569eec124d43c8a54d6256",
			index:         0,
			sequence:      47,
			expected: []effect{
				{
					address:     "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
					effectType:  history.EffectAccountCredited,
					operationID: int64(201863467009),
					order:       uint32(1),
					details: map[string]interface{}{
						"amount":     "15257676.9536092",
						"asset_type": "native",
					},
				},
				{
					address:     "GDR53WAEIKOU3ZKN34CSHAWH7HV6K63CBJRUTWUDBFSMY7RRQK3SPKOS",
					effectType:  history.EffectAccountCredited,
					operationID: int64(201863467009),
					order:       uint32(2),
					details: map[string]interface{}{
						"amount":     "3814420.0001419",
						"asset_type": "native",
					},
				},
			},
		},
		{
			desc:          "manageData - data created",
			envelopeXDR:   "AAAAADEhMVDHiYXdz5z8l73XGyrQ2RN85ZRW1uLsCNQumfsZAAAAZAAAADAAAAACAAAAAAAAAAAAAAABAAAAAAAAAAoAAAAFbmFtZTIAAAAAAAABAAAABDU2NzgAAAAAAAAAAS6Z+xkAAABAjxgnTRBCa0n1efZocxpEjXeITQ5sEYTVd9fowuto2kPw5eFwgVnz6OrKJwCRt5L8ylmWiATXVI3Zyfi3yTKqBA==",
			resultXDR:     "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAKAAAAAAAAAAA=",
			metaXDR:       "AAAAAQAAAAIAAAADAAAAMQAAAAAAAAAAMSExUMeJhd3PnPyXvdcbKtDZE3zllFbW4uwI1C6Z+xkAAAACVAvi1AAAADAAAAABAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAMQAAAAAAAAAAMSExUMeJhd3PnPyXvdcbKtDZE3zllFbW4uwI1C6Z+xkAAAACVAvi1AAAADAAAAACAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAxAAAAAAAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAJUC+LUAAAAMAAAAAIAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAxAAAAAAAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAJUC+LUAAAAMAAAAAIAAAACAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAxAAAAAwAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAVuYW1lMgAAAAAAAAQ1Njc4AAAAAAAAAAA=",
			feeChangesXDR: "AAAAAgAAAAMAAAAxAAAAAAAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAJUC+OcAAAAMAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAxAAAAAAAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAJUC+M4AAAAMAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			hash:          "e4609180751e7702466a8845857df43e4d154ec84b6bad62ce507fe12f1daf99",
			index:         0,
			sequence:      49,
			expected: []effect{
				{
					address:     "GAYSCMKQY6EYLXOPTT6JPPOXDMVNBWITPTSZIVWW4LWARVBOTH5RTLAD",
					effectType:  history.EffectDataCreated,
					operationID: int64(210453401601),
					order:       uint32(1),
					details: map[string]interface{}{
						"name":  xdr.String64("name2"),
						"value": "NTY3OA==",
					},
				},
			},
		},
		{
			desc:          "manageData - data removed",
			envelopeXDR:   "AAAAALly/iTceP/82O3aZAmd8hyqUjYAANfc5RfN0/iibCtTAAAAZAAIGHoAAAAKAAAAAQAAAAAAAAAAAAAAAF4XaMIAAAAAAAAAAQAAAAAAAAAKAAAABWhlbGxvAAAAAAAAAAAAAAAAAAABomwrUwAAAEDyu3HI9bdkzNBs4UgTjVmYt3LQ0CC/6a8yWBmz8OiKeY/RJ9wJvV9/m0JWGtFWbPOXWBg/Pj3ttgKMiHh9TKoF",
			resultXDR:     "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAKAAAAAAAAAAA=",
			metaXDR:       "AAAAAQAAAAIAAAADABE92wAAAAAAAAAAuXL+JNx4//zY7dpkCZ3yHKpSNgAA19zlF83T+KJsK1MAAAAXSHbkGAAIGHoAAAAJAAAAAgAAAAEAAAAAIHtDAQ21/TnXbBjFiB22NXBl7hmD+G5dcpSL1JJTu9wAAAABAAAAFWh0dHBzOi8vd3d3LmhvbWUub3JnLwAAAAMBAgMAAAABAAAAACB7QwENtf0512wYxYgdtjVwZe4Zg/huXXKUi9SSU7vcAAAAAgAAAAAAAAAAAAAAAQARPdsAAAAAAAAAALly/iTceP/82O3aZAmd8hyqUjYAANfc5RfN0/iibCtTAAAAF0h25BgACBh6AAAACgAAAAIAAAABAAAAACB7QwENtf0512wYxYgdtjVwZe4Zg/huXXKUi9SSU7vcAAAAAQAAABVodHRwczovL3d3dy5ob21lLm9yZy8AAAADAQIDAAAAAQAAAAAge0MBDbX9OddsGMWIHbY1cGXuGYP4bl1ylIvUklO73AAAAAIAAAAAAAAAAAAAAAEAAAAEAAAAAwARPcsAAAADAAAAALly/iTceP/82O3aZAmd8hyqUjYAANfc5RfN0/iibCtTAAAABWhlbGxvAAAAAAAAAAAAAAAAAAAAAAAAAgAAAAMAAAAAuXL+JNx4//zY7dpkCZ3yHKpSNgAA19zlF83T+KJsK1MAAAAFaGVsbG8AAAAAAAADABE92wAAAAAAAAAAuXL+JNx4//zY7dpkCZ3yHKpSNgAA19zlF83T+KJsK1MAAAAXSHbkGAAIGHoAAAAKAAAAAgAAAAEAAAAAIHtDAQ21/TnXbBjFiB22NXBl7hmD+G5dcpSL1JJTu9wAAAABAAAAFWh0dHBzOi8vd3d3LmhvbWUub3JnLwAAAAMBAgMAAAABAAAAACB7QwENtf0512wYxYgdtjVwZe4Zg/huXXKUi9SSU7vcAAAAAgAAAAAAAAAAAAAAAQARPdsAAAAAAAAAALly/iTceP/82O3aZAmd8hyqUjYAANfc5RfN0/iibCtTAAAAF0h25BgACBh6AAAACgAAAAEAAAABAAAAACB7QwENtf0512wYxYgdtjVwZe4Zg/huXXKUi9SSU7vcAAAAAQAAABVodHRwczovL3d3dy5ob21lLm9yZy8AAAADAQIDAAAAAQAAAAAge0MBDbX9OddsGMWIHbY1cGXuGYP4bl1ylIvUklO73AAAAAIAAAAAAAAAAA==",
			feeChangesXDR: "AAAAAgAAAAMAET3LAAAAAAAAAAC5cv4k3Hj//Njt2mQJnfIcqlI2AADX3OUXzdP4omwrUwAAABdIduR8AAgYegAAAAkAAAACAAAAAQAAAAAge0MBDbX9OddsGMWIHbY1cGXuGYP4bl1ylIvUklO73AAAAAEAAAAVaHR0cHM6Ly93d3cuaG9tZS5vcmcvAAAAAwECAwAAAAEAAAAAIHtDAQ21/TnXbBjFiB22NXBl7hmD+G5dcpSL1JJTu9wAAAACAAAAAAAAAAAAAAABABE92wAAAAAAAAAAuXL+JNx4//zY7dpkCZ3yHKpSNgAA19zlF83T+KJsK1MAAAAXSHbkGAAIGHoAAAAJAAAAAgAAAAEAAAAAIHtDAQ21/TnXbBjFiB22NXBl7hmD+G5dcpSL1JJTu9wAAAABAAAAFWh0dHBzOi8vd3d3LmhvbWUub3JnLwAAAAMBAgMAAAABAAAAACB7QwENtf0512wYxYgdtjVwZe4Zg/huXXKUi9SSU7vcAAAAAgAAAAAAAAAA",
			hash:          "397b208adb3d484d14ddd3237422baae0b6bd1e8feb3c970147bc6bcc493d112",
			index:         0,
			sequence:      49,
			expected: []effect{
				{
					address:     "GC4XF7RE3R4P77GY5XNGICM56IOKUURWAAANPXHFC7G5H6FCNQVVH3OH",
					effectType:  history.EffectDataRemoved,
					operationID: int64(210453401601),
					order:       uint32(1),
					details: map[string]interface{}{
						"name": xdr.String64("hello"),
					},
				},
			},
		},
		{
			desc:          "manageData - data updated",
			envelopeXDR:   "AAAAAKO5w1Op9wij5oMFtCTUoGO9YgewUKQyeIw1g/L0mMP+AAAAZAAALbYAADNjAAAAAQAAAAAAAAAAAAAAAF4WVfgAAAAAAAAAAQAAAAEAAAAAOO6NdKTWKbGao6zsPag+izHxq3eUPLiwjREobLhQAmQAAAAKAAAAOEdDUjNUUTJUVkgzUVJJN0dRTUMzSUpHVVVCUjMyWVFIV0JJS0lNVFlSUTJZSDRYVVREQjc1VUtFAAAAAQAAABQxNTc4NTIxMjA0XzI5MzI5MDI3OAAAAAAAAAAC0oPafQAAAEAcsS0iq/t8i+p85xwLsRy8JpRNEeqobEC5yuhO9ouVf3PE0VjLqv8sDd0St4qbtXU5fqlHd49R9CR+z7tiRLEB9JjD/gAAAEBmaa9sGxQhEhrakzXcSNpMbR4nox/Ha0p/1sI4tabNEzjgYLwKMn1U9tIdVvKKDwE22jg+CI2FlPJ3+FJPmKUA",
			resultXDR:     "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAKAAAAAAAAAAA=",
			metaXDR:       "AAAAAQAAAAIAAAADABEK2wAAAAAAAAAAo7nDU6n3CKPmgwW0JNSgY71iB7BQpDJ4jDWD8vSYw/4AAAAXSGLVVAAALbYAADNiAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABABEK2wAAAAAAAAAAo7nDU6n3CKPmgwW0JNSgY71iB7BQpDJ4jDWD8vSYw/4AAAAXSGLVVAAALbYAADNjAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAEQqbAAAAAwAAAAA47o10pNYpsZqjrOw9qD6LMfGrd5Q8uLCNEShsuFACZAAAADhHQ1IzVFEyVFZIM1FSSTdHUU1DM0lKR1VVQlIzMllRSFdCSUtJTVRZUlEyWUg0WFVUREI3NVVLRQAAABQxNTc4NTIwODU4XzI1MjM5MTc2OAAAAAAAAAAAAAAAAQARCtsAAAADAAAAADjujXSk1imxmqOs7D2oPosx8at3lDy4sI0RKGy4UAJkAAAAOEdDUjNUUTJUVkgzUVJJN0dRTUMzSUpHVVVCUjMyWVFIV0JJS0lNVFlSUTJZSDRYVVREQjc1VUtFAAAAFDE1Nzg1MjEyMDRfMjkzMjkwMjc4AAAAAAAAAAA=",
			feeChangesXDR: "AAAAAgAAAAMAEQqbAAAAAAAAAACjucNTqfcIo+aDBbQk1KBjvWIHsFCkMniMNYPy9JjD/gAAABdIYtW4AAAttgAAM2IAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAEQrbAAAAAAAAAACjucNTqfcIo+aDBbQk1KBjvWIHsFCkMniMNYPy9JjD/gAAABdIYtVUAAAttgAAM2IAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			hash:          "c60b74a14b628d06d3683db8b36ce81344967ac13bc433124bcef44115fbb257",
			index:         0,
			sequence:      49,
			expected: []effect{
				{
					address:     "GA4O5DLUUTLCTMM2UOWOYPNIH2FTD4NLO6KDZOFQRUISQ3FYKABGJLPC",
					effectType:  history.EffectDataUpdated,
					operationID: int64(210453401601),
					order:       uint32(1),
					details: map[string]interface{}{
						"name":  xdr.String64("GCR3TQ2TVH3QRI7GQMC3IJGUUBR32YQHWBIKIMTYRQ2YH4XUTDB75UKE"),
						"value": "MTU3ODUyMTIwNF8yOTMyOTAyNzg=",
					},
				},
			},
		},
		{
			desc:          "bumpSequence - new_seq is the same as current sequence",
			envelopeXDR:   "AAAAAKGX7RT96eIn205uoUHYnqLbt2cPRNORraEoeTAcrRKUAAAAZAAAAEXZZLgDAAAAAAAAAAAAAAABAAAAAAAAAAsAAABF2WS4AwAAAAAAAAABHK0SlAAAAECcI6ex0Dq6YAh6aK14jHxuAvhvKG2+NuzboAKrfYCaC1ZSQ77BYH/5MghPX97JO9WXV17ehNK7d0umxBgaJj8A",
			resultXDR:     "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAALAAAAAAAAAAA=",
			metaXDR:       "AAAAAQAAAAIAAAADAAAAPQAAAAAAAAAAoZftFP3p4ifbTm6hQdieotu3Zw9E05GtoSh5MBytEpQAAAACVAvicAAAAEXZZLgCAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAPQAAAAAAAAAAoZftFP3p4ifbTm6hQdieotu3Zw9E05GtoSh5MBytEpQAAAACVAvicAAAAEXZZLgDAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAA==",
			feeChangesXDR: "AAAAAgAAAAMAAAA8AAAAAAAAAAChl+0U/eniJ9tObqFB2J6i27dnD0TTka2hKHkwHK0SlAAAAAJUC+LUAAAARdlkuAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAA9AAAAAAAAAAChl+0U/eniJ9tObqFB2J6i27dnD0TTka2hKHkwHK0SlAAAAAJUC+JwAAAARdlkuAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			hash:          "bc11b5c41de791369fd85fa1ccf01c35c20df5f98ff2f75d02ead61bfd520e21",
			index:         0,
			sequence:      61,
			expected:      []effect{},
		},
		{

			desc:          "bumpSequence - new_seq is lower than current sequence",
			envelopeXDR:   "AAAAAKGX7RT96eIn205uoUHYnqLbt2cPRNORraEoeTAcrRKUAAAAZAAAAEXZZLgCAAAAAAAAAAAAAAABAAAAAAAAAAsAAABF2WS4AQAAAAAAAAABHK0SlAAAAEC4H7TDntOUXDMg4MfoCPlbLRQZH7VwNpUHMvtnRWqWIiY/qnYYu0bvgYUVtoFOOeqElRKLYqtOW3Fz9iKl0WQJ",
			resultXDR:     "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAALAAAAAAAAAAA=",
			metaXDR:       "AAAAAQAAAAIAAAADAAAAPAAAAAAAAAAAoZftFP3p4ifbTm6hQdieotu3Zw9E05GtoSh5MBytEpQAAAACVAvi1AAAAEXZZLgBAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAPAAAAAAAAAAAoZftFP3p4ifbTm6hQdieotu3Zw9E05GtoSh5MBytEpQAAAACVAvi1AAAAEXZZLgCAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAA==",
			feeChangesXDR: "AAAAAgAAAAMAAAA7AAAAAAAAAAChl+0U/eniJ9tObqFB2J6i27dnD0TTka2hKHkwHK0SlAAAAAJUC+M4AAAARdlkuAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAA8AAAAAAAAAAChl+0U/eniJ9tObqFB2J6i27dnD0TTka2hKHkwHK0SlAAAAAJUC+LUAAAARdlkuAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			hash:          "c8132b95c0063cafd20b26d27f06c12e688609d2d9d3724b840821e861870b8e",
			index:         0,
			sequence:      60,
			expected:      []effect{},
		},
		{

			desc:          "bumpSequence - new_seq is higher than current sequence",
			envelopeXDR:   "AAAAAKGX7RT96eIn205uoUHYnqLbt2cPRNORraEoeTAcrRKUAAAAZAAAADkAAAABAAAAAAAAAAAAAAABAAAAAAAAAAsAAABF2WS4AAAAAAAAAAABHK0SlAAAAEDq0JVhKNIq9ag0sR+R/cv3d9tEuaYEm2BazIzILRdGj9alaVMZBhxoJ3ZIpP3rraCJzyoKZO+p5HBVe10a2+UG",
			resultXDR:     "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAALAAAAAAAAAAA=",
			metaXDR:       "AAAAAQAAAAIAAAADAAAAOgAAAAAAAAAAoZftFP3p4ifbTm6hQdieotu3Zw9E05GtoSh5MBytEpQAAAACVAvjnAAAADkAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAOgAAAAAAAAAAoZftFP3p4ifbTm6hQdieotu3Zw9E05GtoSh5MBytEpQAAAACVAvjnAAAADkAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAA6AAAAAAAAAAChl+0U/eniJ9tObqFB2J6i27dnD0TTka2hKHkwHK0SlAAAAAJUC+OcAAAAOQAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAA6AAAAAAAAAAChl+0U/eniJ9tObqFB2J6i27dnD0TTka2hKHkwHK0SlAAAAAJUC+OcAAAARdlkuAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			feeChangesXDR: "AAAAAgAAAAMAAAA5AAAAAAAAAAChl+0U/eniJ9tObqFB2J6i27dnD0TTka2hKHkwHK0SlAAAAAJUC+QAAAAAOQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAA6AAAAAAAAAAChl+0U/eniJ9tObqFB2J6i27dnD0TTka2hKHkwHK0SlAAAAAJUC+OcAAAAOQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
			hash:          "829d53f2dceebe10af8007564b0aefde819b95734ad431df84270651e7ed8a90",
			index:         0,
			sequence:      58,
			expected: []effect{
				{
					address:     "GCQZP3IU7XU6EJ63JZXKCQOYT2RNXN3HB5CNHENNUEUHSMA4VUJJJSEN",
					effectType:  history.EffectSequenceBumped,
					operationID: int64(249108107265),
					order:       uint32(1),
					details: map[string]interface{}{
						"new_seq": xdr.SequenceNumber(300000000000),
					},
				},
			},
		},
		{
			desc:          "revokeSponsorship (signer)",
			envelopeXDR:   getRevokeSponsorshipEnvelopeXDR(t),
			resultXDR:     "AAAAAAAAAAAAAAAAAAAAAAAAAAA=",
			metaXDR:       revokeSponsorshipMeta,
			feeChangesXDR: "AAAAAA==",
			hash:          "a41d1c8cdf515203ac5a10d945d5023325076b23dbe7d65ae402cd5f8cd9f891",
			index:         0,
			sequence:      58,
			expected:      revokeSponsorshipEffects,
		},
		{
			desc:          "Failed transaction",
			envelopeXDR:   "AAAAAPCq/iehD2ASJorqlTyEt0usn2WG3yF4w9xBkgd4itu6AAAAZAAMpboAADNGAAAAAAAAAAAAAAABAAAAAAAAAAMAAAABVEVTVAAAAAAObS6P1g8rj8sCVzRQzYgHhWFkbh1oV+1s47LFPstSpQAAAAAAAAACVAvkAAAAAfcAAAD6AAAAAAAAAAAAAAAAAAAAAXiK27oAAABAHHk5mvM6xBRsvu3RBvzzPIb8GpXaL2M7InPn65LIhFJ2RnHIYrpP6ufZc6SUtKqChNRaN4qw5rjwFXNezmrBCw==",
			resultXDR:     "AAAAAAAAAGT/////AAAAAQAAAAAAAAAD////+QAAAAA=",
			metaXDR:       "AAAAAQAAAAIAAAADABDLGAAAAAAAAAAA8Kr+J6EPYBImiuqVPIS3S6yfZYbfIXjD3EGSB3iK27oAAAB2ucIg2AAMpboAADNFAAAA4wAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAQAABHT9ws4fAAAAAAAAAAAAAAAAAAAAAAAAAAEAEMsYAAAAAAAAAADwqv4noQ9gEiaK6pU8hLdLrJ9lht8heMPcQZIHeIrbugAAAHa5wiDYAAylugAAM0YAAADjAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAABAAAEdP3Czh8AAAAAAAAAAAAAAAAAAAAAAAAAAA==",
			feeChangesXDR: "AAAAAgAAAAMAEMsCAAAAAAAAAADwqv4noQ9gEiaK6pU8hLdLrJ9lht8heMPcQZIHeIrbugAAAHa5wiE8AAylugAAM0UAAADjAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAABAAAEdP3Czh8AAAAAAAAAAAAAAAAAAAAAAAAAAQAQyxgAAAAAAAAAAPCq/iehD2ASJorqlTyEt0usn2WG3yF4w9xBkgd4itu6AAAAdrnCINgADKW6AAAzRQAAAOMAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAEAAAR0/cLOHwAAAAAAAAAAAAAAAAAAAAA=",
			hash:          "24206737a02f7f855c46e367418e38c223f897792c76bbfb948e1b0dbd695f8b",
			index:         0,
			sequence:      58,
			expected:      []effect{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tt := assert.New(t)
			transaction := BuildLedgerTransaction(
				t,
				TestTransaction{
					Index:         1,
					EnvelopeXDR:   tc.envelopeXDR,
					ResultXDR:     tc.resultXDR,
					MetaXDR:       tc.metaXDR,
					FeeChangesXDR: tc.feeChangesXDR,
					Hash:          tc.hash,
				},
			)

			operation := transactionOperationWrapper{
				index:          tc.index,
				transaction:    transaction,
				operation:      transaction.Envelope.Operations()[tc.index],
				ledgerSequence: tc.sequence,
			}

			effects, err := operation.effects()
			tt.NoError(err)
			tt.Equal(tc.expected, effects)
		})
	}
}

func TestOperationEffectsSetOptionsSignersOrder(t *testing.T) {
	tt := assert.New(t)
	transaction := ingest.LedgerTransaction{
		UnsafeMeta: createTransactionMeta([]xdr.OperationMeta{
			{
				Changes: []xdr.LedgerEntryChange{
					// State
					{
						Type: xdr.LedgerEntryChangeTypeLedgerEntryState,
						State: &xdr.LedgerEntry{
							Data: xdr.LedgerEntryData{
								Type: xdr.LedgerEntryTypeAccount,
								Account: &xdr.AccountEntry{
									AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
									Signers: []xdr.Signer{
										{
											Key:    xdr.MustSigner("GCBBDQLCTNASZJ3MTKAOYEOWRGSHDFAJVI7VPZUOP7KXNHYR3HP2BUKV"),
											Weight: 10,
										},
										{
											Key:    xdr.MustSigner("GCAHY6JSXQFKWKP6R7U5JPXDVNV4DJWOWRFLY3Y6YPBF64QRL4BPFDNS"),
											Weight: 10,
										},
									},
								},
							},
						},
					},
					// Updated
					{
						Type: xdr.LedgerEntryChangeTypeLedgerEntryUpdated,
						Updated: &xdr.LedgerEntry{
							Data: xdr.LedgerEntryData{
								Type: xdr.LedgerEntryTypeAccount,
								Account: &xdr.AccountEntry{
									AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
									Signers: []xdr.Signer{
										{
											Key:    xdr.MustSigner("GCBBDQLCTNASZJ3MTKAOYEOWRGSHDFAJVI7VPZUOP7KXNHYR3HP2BUKV"),
											Weight: 16,
										},
										{
											Key:    xdr.MustSigner("GCAHY6JSXQFKWKP6R7U5JPXDVNV4DJWOWRFLY3Y6YPBF64QRL4BPFDNS"),
											Weight: 15,
										},
										{
											Key:    xdr.MustSigner("GCR3TQ2TVH3QRI7GQMC3IJGUUBR32YQHWBIKIMTYRQ2YH4XUTDB75UKE"),
											Weight: 14,
										},
										{
											Key:    xdr.MustSigner("GA4O5DLUUTLCTMM2UOWOYPNIH2FTD4NLO6KDZOFQRUISQ3FYKABGJLPC"),
											Weight: 17,
										},
									},
								},
							},
						},
					},
				},
			},
		}),
	}
	transaction.Index = 1
	transaction.Envelope.Type = xdr.EnvelopeTypeEnvelopeTypeTx
	aid := xdr.MustAddress("GCBBDQLCTNASZJ3MTKAOYEOWRGSHDFAJVI7VPZUOP7KXNHYR3HP2BUKV")
	transaction.Envelope.V1 = &xdr.TransactionV1Envelope{
		Tx: xdr.Transaction{
			SourceAccount: aid.ToMuxedAccount(),
		},
	}

	operation := transactionOperationWrapper{
		index:       0,
		transaction: transaction,
		operation: xdr.Operation{
			Body: xdr.OperationBody{
				Type:         xdr.OperationTypeSetOptions,
				SetOptionsOp: &xdr.SetOptionsOp{},
			},
		},
		ledgerSequence: 46,
	}

	effects, err := operation.effects()
	tt.NoError(err)
	expected := []effect{
		{
			address:     "GCBBDQLCTNASZJ3MTKAOYEOWRGSHDFAJVI7VPZUOP7KXNHYR3HP2BUKV",
			operationID: int64(197568499713),
			details: map[string]interface{}{
				"public_key": "GCAHY6JSXQFKWKP6R7U5JPXDVNV4DJWOWRFLY3Y6YPBF64QRL4BPFDNS",
				"weight":     int32(15),
			},
			effectType: history.EffectSignerUpdated,
			order:      uint32(1),
		},
		{
			address:     "GCBBDQLCTNASZJ3MTKAOYEOWRGSHDFAJVI7VPZUOP7KXNHYR3HP2BUKV",
			operationID: int64(197568499713),
			details: map[string]interface{}{
				"public_key": "GCBBDQLCTNASZJ3MTKAOYEOWRGSHDFAJVI7VPZUOP7KXNHYR3HP2BUKV",
				"weight":     int32(16),
			},
			effectType: history.EffectSignerUpdated,
			order:      uint32(2),
		},
		{
			address:     "GCBBDQLCTNASZJ3MTKAOYEOWRGSHDFAJVI7VPZUOP7KXNHYR3HP2BUKV",
			operationID: int64(197568499713),
			details: map[string]interface{}{
				"public_key": "GA4O5DLUUTLCTMM2UOWOYPNIH2FTD4NLO6KDZOFQRUISQ3FYKABGJLPC",
				"weight":     int32(17),
			},
			effectType: history.EffectSignerCreated,
			order:      uint32(3),
		},
		{
			address:     "GCBBDQLCTNASZJ3MTKAOYEOWRGSHDFAJVI7VPZUOP7KXNHYR3HP2BUKV",
			operationID: int64(197568499713),
			details: map[string]interface{}{
				"public_key": "GCR3TQ2TVH3QRI7GQMC3IJGUUBR32YQHWBIKIMTYRQ2YH4XUTDB75UKE",
				"weight":     int32(14),
			},
			effectType: history.EffectSignerCreated,
			order:      uint32(4),
		},
	}
	tt.Equal(expected, effects)
}

// Regression for https://go/issues/2136
func TestOperationEffectsSetOptionsSignersNoUpdated(t *testing.T) {
	tt := assert.New(t)
	transaction := ingest.LedgerTransaction{
		UnsafeMeta: createTransactionMeta([]xdr.OperationMeta{
			{
				Changes: []xdr.LedgerEntryChange{
					// State
					{
						Type: xdr.LedgerEntryChangeTypeLedgerEntryState,
						State: &xdr.LedgerEntry{
							Data: xdr.LedgerEntryData{
								Type: xdr.LedgerEntryTypeAccount,
								Account: &xdr.AccountEntry{
									AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
									Signers: []xdr.Signer{
										{
											Key:    xdr.MustSigner("GCBBDQLCTNASZJ3MTKAOYEOWRGSHDFAJVI7VPZUOP7KXNHYR3HP2BUKV"),
											Weight: 10,
										},
										{
											Key:    xdr.MustSigner("GCAHY6JSXQFKWKP6R7U5JPXDVNV4DJWOWRFLY3Y6YPBF64QRL4BPFDNS"),
											Weight: 10,
										},
										{
											Key:    xdr.MustSigner("GA4O5DLUUTLCTMM2UOWOYPNIH2FTD4NLO6KDZOFQRUISQ3FYKABGJLPC"),
											Weight: 17,
										},
									},
								},
							},
						},
					},
					// Updated
					{
						Type: xdr.LedgerEntryChangeTypeLedgerEntryUpdated,
						Updated: &xdr.LedgerEntry{
							Data: xdr.LedgerEntryData{
								Type: xdr.LedgerEntryTypeAccount,
								Account: &xdr.AccountEntry{
									AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
									Signers: []xdr.Signer{
										{
											Key:    xdr.MustSigner("GCBBDQLCTNASZJ3MTKAOYEOWRGSHDFAJVI7VPZUOP7KXNHYR3HP2BUKV"),
											Weight: 16,
										},
										{
											Key:    xdr.MustSigner("GCAHY6JSXQFKWKP6R7U5JPXDVNV4DJWOWRFLY3Y6YPBF64QRL4BPFDNS"),
											Weight: 10,
										},
										{
											Key:    xdr.MustSigner("GCR3TQ2TVH3QRI7GQMC3IJGUUBR32YQHWBIKIMTYRQ2YH4XUTDB75UKE"),
											Weight: 14,
										},
									},
								},
							},
						},
					},
				},
			},
		}),
	}
	transaction.Index = 1
	transaction.Envelope.Type = xdr.EnvelopeTypeEnvelopeTypeTx
	aid := xdr.MustAddress("GCBBDQLCTNASZJ3MTKAOYEOWRGSHDFAJVI7VPZUOP7KXNHYR3HP2BUKV")
	transaction.Envelope.V1 = &xdr.TransactionV1Envelope{
		Tx: xdr.Transaction{
			SourceAccount: aid.ToMuxedAccount(),
		},
	}

	operation := transactionOperationWrapper{
		index:       0,
		transaction: transaction,
		operation: xdr.Operation{
			Body: xdr.OperationBody{
				Type:         xdr.OperationTypeSetOptions,
				SetOptionsOp: &xdr.SetOptionsOp{},
			},
		},
		ledgerSequence: 46,
	}

	effects, err := operation.effects()
	tt.NoError(err)
	expected := []effect{
		{
			address:     "GCBBDQLCTNASZJ3MTKAOYEOWRGSHDFAJVI7VPZUOP7KXNHYR3HP2BUKV",
			operationID: int64(197568499713),
			details: map[string]interface{}{
				"public_key": "GA4O5DLUUTLCTMM2UOWOYPNIH2FTD4NLO6KDZOFQRUISQ3FYKABGJLPC",
			},
			effectType: history.EffectSignerRemoved,
			order:      uint32(1),
		},
		{
			address:     "GCBBDQLCTNASZJ3MTKAOYEOWRGSHDFAJVI7VPZUOP7KXNHYR3HP2BUKV",
			operationID: int64(197568499713),
			details: map[string]interface{}{
				"public_key": "GCBBDQLCTNASZJ3MTKAOYEOWRGSHDFAJVI7VPZUOP7KXNHYR3HP2BUKV",
				"weight":     int32(16),
			},
			effectType: history.EffectSignerUpdated,
			order:      uint32(2),
		},
		{
			address:     "GCBBDQLCTNASZJ3MTKAOYEOWRGSHDFAJVI7VPZUOP7KXNHYR3HP2BUKV",
			operationID: int64(197568499713),
			details: map[string]interface{}{
				"public_key": "GCR3TQ2TVH3QRI7GQMC3IJGUUBR32YQHWBIKIMTYRQ2YH4XUTDB75UKE",
				"weight":     int32(14),
			},
			effectType: history.EffectSignerCreated,
			order:      uint32(3),
		},
	}
	tt.Equal(expected, effects)
}

func TestOperationRegressionAccountTrustItself(t *testing.T) {
	tt := assert.New(t)
	// NOTE:  when an account trusts itself, the transaction is successful but
	// no ledger entries are actually modified.
	transaction := ingest.LedgerTransaction{
		UnsafeMeta: createTransactionMeta([]xdr.OperationMeta{}),
	}
	transaction.Index = 1
	transaction.Envelope.Type = xdr.EnvelopeTypeEnvelopeTypeTx
	aid := xdr.MustAddress("GCBBDQLCTNASZJ3MTKAOYEOWRGSHDFAJVI7VPZUOP7KXNHYR3HP2BUKV")
	transaction.Envelope.V1 = &xdr.TransactionV1Envelope{
		Tx: xdr.Transaction{
			SourceAccount: aid.ToMuxedAccount(),
		},
	}
	operation := transactionOperationWrapper{
		index:       0,
		transaction: transaction,
		operation: xdr.Operation{
			Body: xdr.OperationBody{
				Type: xdr.OperationTypeChangeTrust,
				ChangeTrustOp: &xdr.ChangeTrustOp{
					Line:  xdr.MustNewCreditAsset("COP", "GCBBDQLCTNASZJ3MTKAOYEOWRGSHDFAJVI7VPZUOP7KXNHYR3HP2BUKV").ToChangeTrustAsset(),
					Limit: xdr.Int64(1000),
				},
			},
		},
		ledgerSequence: 46,
	}

	effects, err := operation.effects()
	tt.NoError(err)
	tt.Equal([]effect{}, effects)
}

func TestOperationEffectsAllowTrustAuthorizedToMaintainLiabilities(t *testing.T) {
	tt := assert.New(t)
	asset := xdr.Asset{}
	allowTrustAsset, err := asset.ToAssetCode("COP")
	tt.NoError(err)
	aid := xdr.MustAddress("GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD")
	source := aid.ToMuxedAccount()
	op := xdr.Operation{
		SourceAccount: &source,
		Body: xdr.OperationBody{
			Type: xdr.OperationTypeAllowTrust,
			AllowTrustOp: &xdr.AllowTrustOp{
				Trustor:   xdr.MustAddress("GDQNY3PBOJOKYZSRMK2S7LHHGWZIUISD4QORETLMXEWXBI7KFZZMKTL3"),
				Asset:     allowTrustAsset,
				Authorize: xdr.Uint32(xdr.TrustLineFlagsAuthorizedToMaintainLiabilitiesFlag),
			},
		},
	}

	operation := transactionOperationWrapper{
		index: 0,
		transaction: ingest.LedgerTransaction{
			UnsafeMeta: xdr.TransactionMeta{
				V:  2,
				V2: &xdr.TransactionMetaV2{},
			},
		},
		operation:      op,
		ledgerSequence: 1,
	}

	effects, err := operation.effects()
	tt.NoError(err)

	expected := []effect{
		{
			address:     "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
			operationID: 4294967297,
			details: map[string]interface{}{
				"asset_code":   "COP",
				"asset_issuer": "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
				"asset_type":   "credit_alphanum4",
				"trustor":      "GDQNY3PBOJOKYZSRMK2S7LHHGWZIUISD4QORETLMXEWXBI7KFZZMKTL3",
			},
			effectType: history.EffectTrustlineAuthorizedToMaintainLiabilities,
			order:      uint32(1),
		},
		{
			address:     "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
			operationID: int64(4294967297),
			details: map[string]interface{}{
				"asset_code":                        "COP",
				"asset_issuer":                      "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
				"asset_type":                        "credit_alphanum4",
				"authorized_to_maintain_liabilites": true,
				"trustor":                           "GDQNY3PBOJOKYZSRMK2S7LHHGWZIUISD4QORETLMXEWXBI7KFZZMKTL3",
			},
			effectType: history.EffectTrustlineFlagsUpdated,
			order:      uint32(2),
		},
	}
	tt.Equal(expected, effects)
}

func TestOperationEffectsClawback(t *testing.T) {
	tt := assert.New(t)
	aid := xdr.MustAddress("GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD")
	source := aid.ToMuxedAccount()
	op := xdr.Operation{
		SourceAccount: &source,
		Body: xdr.OperationBody{
			Type: xdr.OperationTypeClawback,
			ClawbackOp: &xdr.ClawbackOp{
				Asset:  xdr.MustNewCreditAsset("COP", "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD"),
				From:   xdr.MustMuxedAddress("GDQNY3PBOJOKYZSRMK2S7LHHGWZIUISD4QORETLMXEWXBI7KFZZMKTL3"),
				Amount: 34,
			},
		},
	}

	operation := transactionOperationWrapper{
		index: 0,
		transaction: ingest.LedgerTransaction{
			UnsafeMeta: xdr.TransactionMeta{
				V:  2,
				V2: &xdr.TransactionMetaV2{},
			},
		},
		operation:      op,
		ledgerSequence: 1,
	}

	effects, err := operation.effects()
	tt.NoError(err)

	expected := []effect{
		{
			address:     "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
			operationID: 4294967297,
			details: map[string]interface{}{
				"asset_code":   "COP",
				"asset_issuer": "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
				"asset_type":   "credit_alphanum4",
				"amount":       "0.0000034",
			},
			effectType: history.EffectAccountCredited,
			order:      uint32(1),
		},
		{
			address:     "GDQNY3PBOJOKYZSRMK2S7LHHGWZIUISD4QORETLMXEWXBI7KFZZMKTL3",
			operationID: 4294967297,
			details: map[string]interface{}{
				"asset_code":   "COP",
				"asset_issuer": "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
				"asset_type":   "credit_alphanum4",
				"amount":       "0.0000034",
			},
			effectType: history.EffectAccountDebited,
			order:      uint32(2),
		},
	}
	tt.Equal(expected, effects)
}

func TestOperationEffectsClawbackClaimableBalance(t *testing.T) {
	tt := assert.New(t)
	aid := xdr.MustAddress("GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD")
	source := aid.ToMuxedAccount()
	var balanceID xdr.ClaimableBalanceId
	xdr.SafeUnmarshalBase64("AAAAANoNV9p9SFDn/BDSqdDrxzH3r7QFdMAzlbF9SRSbkfW+", &balanceID)
	op := xdr.Operation{
		SourceAccount: &source,
		Body: xdr.OperationBody{
			Type: xdr.OperationTypeClawbackClaimableBalance,
			ClawbackClaimableBalanceOp: &xdr.ClawbackClaimableBalanceOp{
				BalanceId: balanceID,
			},
		},
	}

	operation := transactionOperationWrapper{
		index: 0,
		transaction: ingest.LedgerTransaction{
			UnsafeMeta: xdr.TransactionMeta{
				V:  2,
				V2: &xdr.TransactionMetaV2{},
			},
		},
		operation:      op,
		ledgerSequence: 1,
	}

	effects, err := operation.effects()
	tt.NoError(err)

	expected := []effect{
		{
			address:     "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
			operationID: 4294967297,
			details: map[string]interface{}{
				"balance_id": "00000000da0d57da7d4850e7fc10d2a9d0ebc731f7afb40574c03395b17d49149b91f5be",
			},
			effectType: history.EffectClaimableBalanceClawedBack,
			order:      uint32(1),
		},
	}
	tt.Equal(expected, effects)
}

func TestOperationEffectsSetTrustLineFlags(t *testing.T) {
	tt := assert.New(t)
	aid := xdr.MustAddress("GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD")
	source := aid.ToMuxedAccount()
	trustor := xdr.MustAddress("GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY")
	setFlags := xdr.Uint32(xdr.TrustLineFlagsAuthorizedToMaintainLiabilitiesFlag)
	clearFlags := xdr.Uint32(xdr.TrustLineFlagsTrustlineClawbackEnabledFlag | xdr.TrustLineFlagsAuthorizedFlag)
	op := xdr.Operation{
		SourceAccount: &source,
		Body: xdr.OperationBody{
			Type: xdr.OperationTypeSetTrustLineFlags,
			SetTrustLineFlagsOp: &xdr.SetTrustLineFlagsOp{
				Trustor:    trustor,
				Asset:      xdr.MustNewCreditAsset("USD", "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD"),
				ClearFlags: clearFlags,
				SetFlags:   setFlags,
			},
		},
	}

	operation := transactionOperationWrapper{
		index: 0,
		transaction: ingest.LedgerTransaction{
			UnsafeMeta: xdr.TransactionMeta{
				V:  2,
				V2: &xdr.TransactionMetaV2{},
			},
		},
		operation:      op,
		ledgerSequence: 1,
	}

	effects, err := operation.effects()
	tt.NoError(err)

	expected := []effect{
		{
			address:     "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
			operationID: 4294967297,
			details: map[string]interface{}{
				"asset_code":                        "USD",
				"asset_issuer":                      "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
				"asset_type":                        "credit_alphanum4",
				"authorized_flag":                   false,
				"authorized_to_maintain_liabilites": true,
				"clawback_enabled_flag":             false,
				"trustor":                           "GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
			},
			effectType: history.EffectTrustlineFlagsUpdated,
			order:      uint32(1),
		},
	}
	tt.Equal(expected, effects)
}

type CreateClaimableBalanceEffectsTestSuite struct {
	suite.Suite
	ops []xdr.Operation
	tx  ingest.LedgerTransaction
}

func (s *CreateClaimableBalanceEffectsTestSuite) SetupTest() {
	aid := xdr.MustAddress("GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD")
	source := aid.ToMuxedAccount()
	var balanceIDOp0, balanceIDOp1 xdr.ClaimableBalanceId
	xdr.SafeUnmarshalBase64("AAAAANoNV9p9SFDn/BDSqdDrxzH3r7QFdMAzlbF9SRSbkfW+", &balanceIDOp0)
	xdr.SafeUnmarshalBase64("AAAAALHcX0PDa9UefSAzitC6vQOUr802phH8OF2ahLzg6j1D", &balanceIDOp1)
	cb0 := xdr.ClaimableBalanceEntry{
		BalanceId: balanceIDOp0,
		Amount:    xdr.Int64(100000000),
		Asset:     xdr.MustNewNativeAsset(),
		Claimants: []xdr.Claimant{
			{
				Type: xdr.ClaimantTypeClaimantTypeV0,
				V0: &xdr.ClaimantV0{
					Destination: xdr.MustAddress("GD5OVB6FKDV7P7SOJ5UB2BPLBL4XGSHPYHINR5355SY3RSXLT2BZWAKY"),

					Predicate: xdr.ClaimPredicate{
						Type: xdr.ClaimPredicateTypeClaimPredicateUnconditional,
					},
				},
			},
		},
		Ext: xdr.ClaimableBalanceEntryExt{
			V: 1,
			V1: &xdr.ClaimableBalanceEntryExtensionV1{
				Flags: xdr.Uint32(xdr.ClaimableBalanceFlagsClaimableBalanceClawbackEnabledFlag),
			},
		},
	}
	relBefore := xdr.Int64(1000)
	cb1 := xdr.ClaimableBalanceEntry{
		BalanceId: balanceIDOp1,
		Amount:    xdr.Int64(200000000),
		Asset:     xdr.MustNewCreditAsset("USD", "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD"),
		Claimants: []xdr.Claimant{
			{
				Type: xdr.ClaimantTypeClaimantTypeV0,
				V0: &xdr.ClaimantV0{
					Destination: xdr.MustAddress("GDMQUXK7ZUCWM5472ZU3YLDP4BMJLQQ76DEMNYDEY2ODEEGGRKLEWGW2"),
					Predicate: xdr.ClaimPredicate{
						Type:      xdr.ClaimPredicateTypeClaimPredicateBeforeRelativeTime,
						RelBefore: &relBefore,
					},
				},
			},
			{
				Type: xdr.ClaimantTypeClaimantTypeV0,
				V0: &xdr.ClaimantV0{
					Destination: xdr.MustAddress("GDQNY3PBOJOKYZSRMK2S7LHHGWZIUISD4QORETLMXEWXBI7KFZZMKTL3"),
					Predicate: xdr.ClaimPredicate{
						Type: xdr.ClaimPredicateTypeClaimPredicateUnconditional,
					},
				},
			},
		},
	}
	s.ops = []xdr.Operation{
		{
			SourceAccount: &source,
			Body: xdr.OperationBody{
				Type: xdr.OperationTypeCreateClaimableBalance,
				CreateClaimableBalanceOp: &xdr.CreateClaimableBalanceOp{
					Amount:    cb0.Amount,
					Asset:     cb0.Asset,
					Claimants: cb0.Claimants,
				},
			},
		},
		{
			SourceAccount: &source,
			Body: xdr.OperationBody{
				Type: xdr.OperationTypeCreateClaimableBalance,
				CreateClaimableBalanceOp: &xdr.CreateClaimableBalanceOp{
					Amount:    cb1.Amount,
					Asset:     cb1.Asset,
					Claimants: cb1.Claimants,
				},
			},
		},
	}

	s.tx = ingest.LedgerTransaction{
		Index: 0,
		Envelope: xdr.TransactionEnvelope{
			Type: xdr.EnvelopeTypeEnvelopeTypeTx,
			V1: &xdr.TransactionV1Envelope{
				Tx: xdr.Transaction{
					Operations: s.ops,
				},
			},
		},
		UnsafeMeta: xdr.TransactionMeta{
			V: 2,
			V2: &xdr.TransactionMetaV2{
				Operations: []xdr.OperationMeta{
					{
						Changes: []xdr.LedgerEntryChange{
							{
								Type: xdr.LedgerEntryChangeTypeLedgerEntryCreated,
								Created: &xdr.LedgerEntry{
									Data: xdr.LedgerEntryData{
										Type:             xdr.LedgerEntryTypeClaimableBalance,
										ClaimableBalance: &cb0,
									},
								},
							},
						},
					},
					{
						Changes: []xdr.LedgerEntryChange{
							{
								Type: xdr.LedgerEntryChangeTypeLedgerEntryCreated,
								Created: &xdr.LedgerEntry{
									Data: xdr.LedgerEntryData{
										Type:             xdr.LedgerEntryTypeClaimableBalance,
										ClaimableBalance: &cb1,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
func (s *CreateClaimableBalanceEffectsTestSuite) TestEffects() {
	relBefore := xdr.Int64(1000)
	testCases := []struct {
		desc     string
		op       xdr.Operation
		expected []effect
	}{
		{
			desc: "claimable balance with native asset",
			op:   s.ops[0],
			expected: []effect{
				{
					address: "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
					details: map[string]interface{}{
						"asset":      "native",
						"amount":     "10.0000000",
						"balance_id": "00000000da0d57da7d4850e7fc10d2a9d0ebc731f7afb40574c03395b17d49149b91f5be",
						"claimable_balance_clawback_enabled_flag": true,
					},
					effectType:  history.EffectClaimableBalanceCreated,
					operationID: int64(4294967297),
					order:       uint32(1),
				},
				{
					address: "GD5OVB6FKDV7P7SOJ5UB2BPLBL4XGSHPYHINR5355SY3RSXLT2BZWAKY",
					details: map[string]interface{}{
						"asset":      "native",
						"amount":     "10.0000000",
						"balance_id": "00000000da0d57da7d4850e7fc10d2a9d0ebc731f7afb40574c03395b17d49149b91f5be",
						"predicate":  xdr.ClaimPredicate{Type: xdr.ClaimPredicateTypeClaimPredicateUnconditional},
					},
					effectType:  history.EffectClaimableBalanceClaimantCreated,
					operationID: int64(4294967297),
					order:       uint32(2),
				},
				{
					address: "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
					details: map[string]interface{}{
						"amount":     "10.0000000",
						"asset_type": "native",
					},
					effectType:  history.EffectAccountDebited,
					operationID: int64(4294967297),
					order:       uint32(3),
				},
			},
		},
		{
			desc: "claimable balance with issued asset",
			op:   s.ops[1],
			expected: []effect{
				{
					address: "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
					details: map[string]interface{}{
						"asset":      "USD:GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
						"amount":     "20.0000000",
						"balance_id": "00000000b1dc5f43c36bd51e7d20338ad0babd0394afcd36a611fc385d9a84bce0ea3d43",
					},
					effectType:  history.EffectClaimableBalanceCreated,
					operationID: int64(4294967298),
					order:       uint32(1),
				},
				{
					address: "GDMQUXK7ZUCWM5472ZU3YLDP4BMJLQQ76DEMNYDEY2ODEEGGRKLEWGW2",
					details: map[string]interface{}{
						"asset":      "USD:GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
						"amount":     "20.0000000",
						"balance_id": "00000000b1dc5f43c36bd51e7d20338ad0babd0394afcd36a611fc385d9a84bce0ea3d43",
						// Make sure data ingested from op body (rel_before)
						"predicate": xdr.ClaimPredicate{
							Type:      xdr.ClaimPredicateTypeClaimPredicateBeforeRelativeTime,
							RelBefore: &relBefore,
						},
					},
					effectType:  history.EffectClaimableBalanceClaimantCreated,
					operationID: int64(4294967298),
					order:       uint32(2),
				},
				{
					address: "GDQNY3PBOJOKYZSRMK2S7LHHGWZIUISD4QORETLMXEWXBI7KFZZMKTL3",
					details: map[string]interface{}{
						"asset":      "USD:GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
						"amount":     "20.0000000",
						"balance_id": "00000000b1dc5f43c36bd51e7d20338ad0babd0394afcd36a611fc385d9a84bce0ea3d43",
						"predicate":  xdr.ClaimPredicate{Type: xdr.ClaimPredicateTypeClaimPredicateUnconditional},
					},
					effectType:  history.EffectClaimableBalanceClaimantCreated,
					operationID: int64(4294967298),
					order:       uint32(3),
				},
				{
					address: "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
					details: map[string]interface{}{
						"amount":       "20.0000000",
						"asset_code":   "USD",
						"asset_type":   "credit_alphanum4",
						"asset_issuer": "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
					},
					effectType:  history.EffectAccountDebited,
					operationID: int64(4294967298),
					order:       uint32(4),
				},
			},
		},
	}
	for i, tc := range testCases {
		s.T().Run(tc.desc, func(t *testing.T) {
			operation := transactionOperationWrapper{
				index:          uint32(i),
				transaction:    s.tx,
				operation:      tc.op,
				ledgerSequence: 1,
			}

			effects, err := operation.effects()
			s.Assert().NoError(err)
			s.Assert().Equal(tc.expected, effects)
		})
	}
}

func TestCreateClaimableBalanceEffectsTestSuite(t *testing.T) {
	suite.Run(t, new(CreateClaimableBalanceEffectsTestSuite))
}

type ClaimClaimableBalanceEffectsTestSuite struct {
	suite.Suite
	ops []xdr.Operation
	tx  ingest.LedgerTransaction
}

func (s *ClaimClaimableBalanceEffectsTestSuite) SetupTest() {
	var balanceIDOp1, balanceIDOp1Meta, balanceIDOp2, balanceIDOp2Meta xdr.ClaimableBalanceId
	xdr.SafeUnmarshalBase64("AAAAANoNV9p9SFDn/BDSqdDrxzH3r7QFdMAzlbF9SRSbkfW+", &balanceIDOp1)
	xdr.SafeUnmarshalBase64("AAAAANoNV9p9SFDn/BDSqdDrxzH3r7QFdMAzlbF9SRSbkfW+", &balanceIDOp1Meta)
	xdr.SafeUnmarshalBase64("AAAAALHcX0PDa9UefSAzitC6vQOUr802phH8OF2ahLzg6j1D", &balanceIDOp2)
	xdr.SafeUnmarshalBase64("AAAAALHcX0PDa9UefSAzitC6vQOUr802phH8OF2ahLzg6j1D", &balanceIDOp2Meta)

	aid := xdr.MustAddress("GD5OVB6FKDV7P7SOJ5UB2BPLBL4XGSHPYHINR5355SY3RSXLT2BZWAKY")
	claimant1 := aid.ToMuxedAccount()
	aid = xdr.MustAddress("GDMQUXK7ZUCWM5472ZU3YLDP4BMJLQQ76DEMNYDEY2ODEEGGRKLEWGW2")
	claimant2 := aid.ToMuxedAccount()
	s.ops = []xdr.Operation{
		{
			SourceAccount: &claimant1,
			Body: xdr.OperationBody{
				Type: xdr.OperationTypeClaimClaimableBalance,
				ClaimClaimableBalanceOp: &xdr.ClaimClaimableBalanceOp{
					BalanceId: balanceIDOp1,
				},
			},
		},
		{
			SourceAccount: &claimant2,
			Body: xdr.OperationBody{
				Type: xdr.OperationTypeClaimClaimableBalance,
				ClaimClaimableBalanceOp: &xdr.ClaimClaimableBalanceOp{
					BalanceId: balanceIDOp2,
				},
			},
		},
	}

	s.tx = ingest.LedgerTransaction{
		Index: 0,
		Envelope: xdr.TransactionEnvelope{
			Type: xdr.EnvelopeTypeEnvelopeTypeTx,
			V1: &xdr.TransactionV1Envelope{
				Tx: xdr.Transaction{
					Operations: s.ops,
				},
			},
		},
		Result: xdr.TransactionResultPair{
			Result: xdr.TransactionResult{
				Result: xdr.TransactionResultResult{
					Results: &[]xdr.OperationResult{
						{
							Code: xdr.OperationResultCodeOpInner,
							Tr: &xdr.OperationResultTr{
								Type: xdr.OperationTypeClaimClaimableBalance,
								ClaimClaimableBalanceResult: &xdr.ClaimClaimableBalanceResult{
									Code: xdr.ClaimClaimableBalanceResultCodeClaimClaimableBalanceSuccess,
								},
							},
						},
						{
							Code: xdr.OperationResultCodeOpInner,
							Tr: &xdr.OperationResultTr{
								Type: xdr.OperationTypeClaimClaimableBalance,
								ClaimClaimableBalanceResult: &xdr.ClaimClaimableBalanceResult{
									Code: xdr.ClaimClaimableBalanceResultCodeClaimClaimableBalanceSuccess,
								},
							},
						},
					},
				},
			},
		},
		FeeChanges: xdr.LedgerEntryChanges{},
		UnsafeMeta: xdr.TransactionMeta{
			V: 2,
			V2: &xdr.TransactionMetaV2{
				Operations: []xdr.OperationMeta{
					// op1
					{
						Changes: xdr.LedgerEntryChanges{
							xdr.LedgerEntryChange{
								Type: xdr.LedgerEntryChangeTypeLedgerEntryState,
								State: &xdr.LedgerEntry{
									Data: xdr.LedgerEntryData{
										Type: xdr.LedgerEntryTypeClaimableBalance,
										ClaimableBalance: &xdr.ClaimableBalanceEntry{
											BalanceId: balanceIDOp1Meta,
											Amount:    xdr.Int64(100000000),
											Asset:     xdr.MustNewNativeAsset(),
											Claimants: []xdr.Claimant{
												{
													Type: xdr.ClaimantTypeClaimantTypeV0,
													V0: &xdr.ClaimantV0{
														Destination: xdr.MustAddress("GD5OVB6FKDV7P7SOJ5UB2BPLBL4XGSHPYHINR5355SY3RSXLT2BZWAKY"),

														Predicate: xdr.ClaimPredicate{
															Type: xdr.ClaimPredicateTypeClaimPredicateUnconditional,
														},
													},
												},
											},
											Ext: xdr.ClaimableBalanceEntryExt{
												V: 1,
												V1: &xdr.ClaimableBalanceEntryExtensionV1{
													Flags: xdr.Uint32(xdr.ClaimableBalanceFlagsClaimableBalanceClawbackEnabledFlag),
												},
											},
										},
									},
								},
							},
							xdr.LedgerEntryChange{
								Type: xdr.LedgerEntryChangeTypeLedgerEntryRemoved,
								Removed: &xdr.LedgerKey{
									Type: xdr.LedgerEntryTypeClaimableBalance,
									ClaimableBalance: &xdr.LedgerKeyClaimableBalance{
										BalanceId: balanceIDOp1Meta,
									},
								},
							},
						},
					},
					// op2
					{
						Changes: xdr.LedgerEntryChanges{
							xdr.LedgerEntryChange{
								Type: xdr.LedgerEntryChangeTypeLedgerEntryState,
								State: &xdr.LedgerEntry{
									Data: xdr.LedgerEntryData{
										Type: xdr.LedgerEntryTypeClaimableBalance,
										ClaimableBalance: &xdr.ClaimableBalanceEntry{
											BalanceId: balanceIDOp2Meta,
											Amount:    xdr.Int64(200000000),
											Asset:     xdr.MustNewCreditAsset("USD", "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD"),
											Claimants: []xdr.Claimant{
												{
													Type: xdr.ClaimantTypeClaimantTypeV0,
													V0: &xdr.ClaimantV0{
														Destination: xdr.MustAddress("GDMQUXK7ZUCWM5472ZU3YLDP4BMJLQQ76DEMNYDEY2ODEEGGRKLEWGW2"),
														Predicate: xdr.ClaimPredicate{
															Type: xdr.ClaimPredicateTypeClaimPredicateUnconditional,
														},
													},
												},
												{
													Type: xdr.ClaimantTypeClaimantTypeV0,
													V0: &xdr.ClaimantV0{
														Destination: xdr.MustAddress("GDQNY3PBOJOKYZSRMK2S7LHHGWZIUISD4QORETLMXEWXBI7KFZZMKTL3"),
														Predicate: xdr.ClaimPredicate{
															Type: xdr.ClaimPredicateTypeClaimPredicateUnconditional,
														},
													},
												},
											},
										},
									},
								},
							},
							xdr.LedgerEntryChange{
								Type: xdr.LedgerEntryChangeTypeLedgerEntryRemoved,
								Removed: &xdr.LedgerKey{
									Type: xdr.LedgerEntryTypeClaimableBalance,
									ClaimableBalance: &xdr.LedgerKeyClaimableBalance{
										BalanceId: balanceIDOp2Meta,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
func (s *ClaimClaimableBalanceEffectsTestSuite) TestEffects() {
	testCases := []struct {
		desc     string
		op       xdr.Operation
		expected []effect
	}{
		{
			desc: "claimable balance with native asset",
			op:   s.ops[0],
			expected: []effect{
				{
					address: "GD5OVB6FKDV7P7SOJ5UB2BPLBL4XGSHPYHINR5355SY3RSXLT2BZWAKY",
					details: map[string]interface{}{
						"asset":      "native",
						"amount":     "10.0000000",
						"balance_id": "00000000da0d57da7d4850e7fc10d2a9d0ebc731f7afb40574c03395b17d49149b91f5be",
						"claimable_balance_clawback_enabled_flag": true,
					},
					effectType:  history.EffectClaimableBalanceClaimed,
					operationID: int64(4294967297),
					order:       uint32(1),
				},
				{
					address: "GD5OVB6FKDV7P7SOJ5UB2BPLBL4XGSHPYHINR5355SY3RSXLT2BZWAKY",
					details: map[string]interface{}{
						"asset_type": "native",
						"amount":     "10.0000000",
					},
					effectType:  history.EffectAccountCredited,
					operationID: int64(4294967297),
					order:       uint32(2),
				},
			},
		},
		{
			desc: "claimable balance with issued asset",
			op:   s.ops[1],
			expected: []effect{
				{
					address: "GDMQUXK7ZUCWM5472ZU3YLDP4BMJLQQ76DEMNYDEY2ODEEGGRKLEWGW2",
					details: map[string]interface{}{
						"asset":      "USD:GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
						"amount":     "20.0000000",
						"balance_id": "00000000b1dc5f43c36bd51e7d20338ad0babd0394afcd36a611fc385d9a84bce0ea3d43",
					},
					effectType:  history.EffectClaimableBalanceClaimed,
					operationID: int64(4294967298),
					order:       uint32(1),
				},
				{
					address: "GDMQUXK7ZUCWM5472ZU3YLDP4BMJLQQ76DEMNYDEY2ODEEGGRKLEWGW2",
					details: map[string]interface{}{
						"amount":       "20.0000000",
						"asset_code":   "USD",
						"asset_type":   "credit_alphanum4",
						"asset_issuer": "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
					},
					effectType:  history.EffectAccountCredited,
					operationID: int64(4294967298),
					order:       uint32(2),
				},
			},
		},
	}
	for i, tc := range testCases {
		s.T().Run(tc.desc, func(t *testing.T) {
			operation := transactionOperationWrapper{
				index:          uint32(i),
				transaction:    s.tx,
				operation:      tc.op,
				ledgerSequence: 1,
			}

			effects, err := operation.effects()
			s.Assert().NoError(err)
			s.Assert().Equal(tc.expected, effects)
		})
	}
}

func TestClaimClaimableBalanceEffectsTestSuite(t *testing.T) {
	suite.Run(t, new(ClaimClaimableBalanceEffectsTestSuite))
}

func TestTrustlineSponsorhipEffects(t *testing.T) {
	source := xdr.MustMuxedAddress("GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY")
	usdAsset := xdr.MustNewCreditAsset("USD", "GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY")
	poolIDStr := "19cc788419412926a11049b9fb1f87906b8f02bc6bf8f73d8fd347ede0b79fa5"
	var poolID xdr.PoolId
	poolIDBytes, err := hex.DecodeString(poolIDStr)
	assert.NoError(t, err)
	copy(poolID[:], poolIDBytes)
	baseAssetTrustLineEntry := xdr.LedgerEntry{
		LastModifiedLedgerSeq: 20,
		Data: xdr.LedgerEntryData{
			Type: xdr.LedgerEntryTypeTrustline,
			TrustLine: &xdr.TrustLineEntry{
				AccountId: source.ToAccountId(),
				Asset:     usdAsset.ToTrustLineAsset(),
				Balance:   100,
				Limit:     1000,
				Flags:     0,
			},
		},
	}
	baseLiquidityPoolTrustLineEntry := xdr.LedgerEntry{
		LastModifiedLedgerSeq: 20,
		Data: xdr.LedgerEntryData{
			Type: xdr.LedgerEntryTypeTrustline,
			TrustLine: &xdr.TrustLineEntry{
				AccountId: source.ToAccountId(),
				Asset: xdr.TrustLineAsset{
					Type:            xdr.AssetTypeAssetTypePoolShare,
					LiquidityPoolId: &poolID,
				},
				Balance: 100,
				Limit:   1000,
				Flags:   0,
			},
		},
	}

	sponsor1 := xdr.MustAddress("GDMQUXK7ZUCWM5472ZU3YLDP4BMJLQQ76DEMNYDEY2ODEEGGRKLEWGW2")
	sponsor2 := xdr.MustAddress("GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD")
	withSponsor := func(le *xdr.LedgerEntry, accID *xdr.AccountId) *xdr.LedgerEntry {
		le2 := *le
		le2.Ext = xdr.LedgerEntryExt{
			V: 1,
			V1: &xdr.LedgerEntryExtensionV1{
				SponsoringId: accID,
			},
		}
		return &le2
	}

	changes := xdr.LedgerEntryChanges{
		// create asset sponsorship
		xdr.LedgerEntryChange{
			Type:  xdr.LedgerEntryChangeTypeLedgerEntryState,
			State: &baseAssetTrustLineEntry,
		},
		xdr.LedgerEntryChange{
			Type:    xdr.LedgerEntryChangeTypeLedgerEntryUpdated,
			Updated: withSponsor(&baseAssetTrustLineEntry, &sponsor1),
		},
		// update asset sponsorship
		xdr.LedgerEntryChange{
			Type:  xdr.LedgerEntryChangeTypeLedgerEntryState,
			State: withSponsor(&baseAssetTrustLineEntry, &sponsor1),
		},
		xdr.LedgerEntryChange{
			Type:    xdr.LedgerEntryChangeTypeLedgerEntryUpdated,
			Updated: withSponsor(&baseAssetTrustLineEntry, &sponsor2),
		},
		// remove asset sponsorship
		xdr.LedgerEntryChange{
			Type:  xdr.LedgerEntryChangeTypeLedgerEntryState,
			State: withSponsor(&baseAssetTrustLineEntry, &sponsor2),
		},
		xdr.LedgerEntryChange{
			Type:    xdr.LedgerEntryChangeTypeLedgerEntryUpdated,
			Updated: &baseAssetTrustLineEntry,
		},

		// create liquidity pool sponsorship
		xdr.LedgerEntryChange{
			Type:  xdr.LedgerEntryChangeTypeLedgerEntryState,
			State: &baseLiquidityPoolTrustLineEntry,
		},
		xdr.LedgerEntryChange{
			Type:    xdr.LedgerEntryChangeTypeLedgerEntryUpdated,
			Updated: withSponsor(&baseLiquidityPoolTrustLineEntry, &sponsor1),
		},
		// update liquidity pool sponsorship
		xdr.LedgerEntryChange{
			Type:  xdr.LedgerEntryChangeTypeLedgerEntryState,
			State: withSponsor(&baseLiquidityPoolTrustLineEntry, &sponsor1),
		},
		xdr.LedgerEntryChange{
			Type:    xdr.LedgerEntryChangeTypeLedgerEntryUpdated,
			Updated: withSponsor(&baseLiquidityPoolTrustLineEntry, &sponsor2),
		},
		// remove liquidity pool sponsorship
		xdr.LedgerEntryChange{
			Type:  xdr.LedgerEntryChangeTypeLedgerEntryState,
			State: withSponsor(&baseLiquidityPoolTrustLineEntry, &sponsor2),
		},
		xdr.LedgerEntryChange{
			Type:    xdr.LedgerEntryChangeTypeLedgerEntryUpdated,
			Updated: &baseLiquidityPoolTrustLineEntry,
		},
	}
	expected := []effect{
		{
			effectType:  history.EffectTrustlineSponsorshipCreated,
			order:       1,
			address:     "GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
			operationID: 4294967297,
			details: map[string]interface{}{
				"asset": "USD:GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
				// `asset_type` set in `Effect.UnmarshalDetails` to prevent reingestion
				"sponsor": "GDMQUXK7ZUCWM5472ZU3YLDP4BMJLQQ76DEMNYDEY2ODEEGGRKLEWGW2",
			},
		},
		{
			effectType:  history.EffectTrustlineSponsorshipUpdated,
			order:       2,
			address:     "GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
			operationID: 4294967297,
			details: map[string]interface{}{
				"asset": "USD:GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
				// `asset_type` set in `Effect.UnmarshalDetails` to prevent reingestion
				"former_sponsor": "GDMQUXK7ZUCWM5472ZU3YLDP4BMJLQQ76DEMNYDEY2ODEEGGRKLEWGW2",
				"new_sponsor":    "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
			},
		},
		{
			effectType:  history.EffectTrustlineSponsorshipRemoved,
			order:       3,
			address:     "GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
			operationID: 4294967297,
			details: map[string]interface{}{
				"asset": "USD:GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
				// `asset_type` set in `Effect.UnmarshalDetails` to prevent reingestion
				"former_sponsor": "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
			},
		},
		{
			effectType:  history.EffectTrustlineSponsorshipCreated,
			order:       4,
			address:     "GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
			operationID: 4294967297,
			details: map[string]interface{}{
				"liquidity_pool_id": poolIDStr,
				"asset_type":        "liquidity_pool",
				"sponsor":           "GDMQUXK7ZUCWM5472ZU3YLDP4BMJLQQ76DEMNYDEY2ODEEGGRKLEWGW2",
			},
		},
		{
			effectType:  history.EffectTrustlineSponsorshipUpdated,
			order:       5,
			address:     "GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
			operationID: 4294967297,
			details: map[string]interface{}{
				"liquidity_pool_id": poolIDStr,
				"asset_type":        "liquidity_pool",
				"former_sponsor":    "GDMQUXK7ZUCWM5472ZU3YLDP4BMJLQQ76DEMNYDEY2ODEEGGRKLEWGW2",
				"new_sponsor":       "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
			},
		},
		{
			effectType:  history.EffectTrustlineSponsorshipRemoved,
			order:       6,
			address:     "GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
			operationID: 4294967297,
			details: map[string]interface{}{
				"liquidity_pool_id": poolIDStr,
				"asset_type":        "liquidity_pool",
				"former_sponsor":    "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD",
			},
		},
	}

	// pick an operation with no intrinsic effects
	// (the sponsosrhip effects are obtained from the changes, so it doesn't matter)
	phonyOp := xdr.Operation{
		Body: xdr.OperationBody{
			Type: xdr.OperationTypeEndSponsoringFutureReserves,
		},
	}
	tx := ingest.LedgerTransaction{
		Index: 0,
		Envelope: xdr.TransactionEnvelope{
			Type: xdr.EnvelopeTypeEnvelopeTypeTx,
			V1: &xdr.TransactionV1Envelope{
				Tx: xdr.Transaction{
					SourceAccount: source,
					Operations:    []xdr.Operation{phonyOp},
				},
			},
		},
		UnsafeMeta: xdr.TransactionMeta{
			V: 2,
			V2: &xdr.TransactionMetaV2{
				Operations: []xdr.OperationMeta{{Changes: changes}},
			},
		},
	}

	operation := transactionOperationWrapper{
		index:          0,
		transaction:    tx,
		operation:      phonyOp,
		ledgerSequence: 1,
	}

	effects, err := operation.effects()
	assert.NoError(t, err)
	assert.Equal(t, expected, effects)

}

func TestLiquidityPoolEffects(t *testing.T) {
	source := xdr.MustMuxedAddress("GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY")
	usdAsset := xdr.MustNewCreditAsset("USD", "GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY")
	poolIDStr := "ea4e3e63a95fd840c1394f195722ffdcb2d0d4f0a26589c6ab557d81e6b0bf9d"
	var poolID xdr.PoolId
	poolIDBytes, err := hex.DecodeString(poolIDStr)
	assert.NoError(t, err)
	copy(poolID[:], poolIDBytes)
	baseLiquidityPoolEntry := xdr.LiquidityPoolEntry{
		LiquidityPoolId: poolID,
		Body: xdr.LiquidityPoolEntryBody{
			Type: xdr.LiquidityPoolTypeLiquidityPoolConstantProduct,
			ConstantProduct: &xdr.LiquidityPoolEntryConstantProduct{
				Params: xdr.LiquidityPoolConstantProductParameters{
					AssetA: xdr.MustNewNativeAsset(),
					AssetB: usdAsset,
					Fee:    20,
				},
				ReserveA:                 200,
				ReserveB:                 100,
				TotalPoolShares:          1000,
				PoolSharesTrustLineCount: 10,
			},
		},
	}
	baseState := xdr.LedgerEntryChange{
		Type: xdr.LedgerEntryChangeTypeLedgerEntryState,
		State: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 20,
			Data: xdr.LedgerEntryData{
				Type:          xdr.LedgerEntryTypeLiquidityPool,
				LiquidityPool: &baseLiquidityPoolEntry,
			},
		},
	}
	updateState := func(cp xdr.LiquidityPoolEntryConstantProduct) xdr.LedgerEntryChange {
		return xdr.LedgerEntryChange{
			Type: xdr.LedgerEntryChangeTypeLedgerEntryUpdated,
			Updated: &xdr.LedgerEntry{
				LastModifiedLedgerSeq: 20,
				Data: xdr.LedgerEntryData{
					Type: xdr.LedgerEntryTypeLiquidityPool,
					LiquidityPool: &xdr.LiquidityPoolEntry{
						LiquidityPoolId: poolID,
						Body: xdr.LiquidityPoolEntryBody{
							Type:            xdr.LiquidityPoolTypeLiquidityPoolConstantProduct,
							ConstantProduct: &cp,
						},
					},
				},
			},
		}
	}

	testCases := []struct {
		desc     string
		op       xdr.OperationBody
		result   xdr.OperationResult
		changes  xdr.LedgerEntryChanges
		expected []effect
	}{
		{
			desc: "liquidity pool creation",
			op: xdr.OperationBody{
				Type: xdr.OperationTypeChangeTrust,
				ChangeTrustOp: &xdr.ChangeTrustOp{
					Line: xdr.ChangeTrustAsset{
						Type: xdr.AssetTypeAssetTypePoolShare,
						LiquidityPool: &xdr.LiquidityPoolParameters{
							Type:            xdr.LiquidityPoolTypeLiquidityPoolConstantProduct,
							ConstantProduct: &baseLiquidityPoolEntry.Body.ConstantProduct.Params,
						},
					},
					Limit: 1000,
				},
			},
			changes: xdr.LedgerEntryChanges{
				xdr.LedgerEntryChange{
					Type: xdr.LedgerEntryChangeTypeLedgerEntryCreated,
					Created: &xdr.LedgerEntry{
						LastModifiedLedgerSeq: 20,
						Data: xdr.LedgerEntryData{
							Type:          xdr.LedgerEntryTypeLiquidityPool,
							LiquidityPool: &baseLiquidityPoolEntry,
						},
					},
				},
				xdr.LedgerEntryChange{
					Type: xdr.LedgerEntryChangeTypeLedgerEntryCreated,
					Created: &xdr.LedgerEntry{
						LastModifiedLedgerSeq: 20,
						Data: xdr.LedgerEntryData{
							Type: xdr.LedgerEntryTypeTrustline,
							TrustLine: &xdr.TrustLineEntry{
								AccountId: source.ToAccountId(),
								Asset: xdr.TrustLineAsset{
									Type:            xdr.AssetTypeAssetTypePoolShare,
									LiquidityPoolId: &poolID,
								},
							},
						},
					},
				},
			},
			expected: []effect{
				{
					effectType:  history.EffectTrustlineCreated,
					order:       1,
					address:     "GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
					operationID: 4294967297,
					details: map[string]interface{}{
						"asset_type":        "liquidity_pool_shares",
						"limit":             "0.0001000",
						"liquidity_pool_id": poolIDStr,
					},
				},
				{
					effectType:  history.EffectLiquidityPoolCreated,
					order:       2,
					address:     "GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
					operationID: 4294967297,
					details: map[string]interface{}{
						"liquidity_pool": map[string]interface{}{
							"fee_bp": uint32(20),
							"id":     poolIDStr,
							"reserves": []base.AssetAmount{
								{
									"native",
									"0.0000200",
								},
								{
									"USD:GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
									"0.0000100",
								},
							},
							"total_shares":     "0.0001000",
							"total_trustlines": "10",
							"type":             "constant_product",
						},
					},
				},
			},
		},
		{
			desc: "liquidity pool deposit",
			op: xdr.OperationBody{
				Type: xdr.OperationTypeLiquidityPoolDeposit,
				LiquidityPoolDepositOp: &xdr.LiquidityPoolDepositOp{
					LiquidityPoolId: poolID,
					MaxAmountA:      100,
					MaxAmountB:      200,
					MinPrice: xdr.Price{
						N: 50,
						D: 3,
					},
					MaxPrice: xdr.Price{
						N: 100,
						D: 2,
					},
				},
			},
			changes: xdr.LedgerEntryChanges{
				baseState,
				updateState(xdr.LiquidityPoolEntryConstantProduct{

					Params:                   baseLiquidityPoolEntry.Body.ConstantProduct.Params,
					ReserveA:                 baseLiquidityPoolEntry.Body.ConstantProduct.ReserveA + 50,
					ReserveB:                 baseLiquidityPoolEntry.Body.ConstantProduct.ReserveB + 60,
					TotalPoolShares:          baseLiquidityPoolEntry.Body.ConstantProduct.TotalPoolShares + 10,
					PoolSharesTrustLineCount: baseLiquidityPoolEntry.Body.ConstantProduct.PoolSharesTrustLineCount,
				}),
			},
			expected: []effect{
				{
					effectType:  history.EffectLiquidityPoolDeposited,
					order:       1,
					address:     "GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
					operationID: 4294967297,
					details: map[string]interface{}{
						"liquidity_pool": map[string]interface{}{
							"fee_bp": uint32(20),
							"id":     poolIDStr,
							"reserves": []base.AssetAmount{
								{
									"native",
									"0.0000250",
								},
								{
									"USD:GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
									"0.0000160",
								},
							},
							"total_shares":     "0.0001010",
							"total_trustlines": "10",
							"type":             "constant_product",
						},
						"reserves_deposited": []base.AssetAmount{
							{
								"native",
								"0.0000050",
							},
							{
								"USD:GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
								"0.0000060",
							},
						},
						"shares_received": "0.0000010",
					},
				},
			},
		},
		{
			desc: "liquidity pool withdrawal",
			op: xdr.OperationBody{
				Type: xdr.OperationTypeLiquidityPoolWithdraw,
				LiquidityPoolWithdrawOp: &xdr.LiquidityPoolWithdrawOp{
					LiquidityPoolId: poolID,
					Amount:          10,
					MinAmountA:      10,
					MinAmountB:      5,
				},
			},
			changes: xdr.LedgerEntryChanges{
				baseState,
				updateState(xdr.LiquidityPoolEntryConstantProduct{

					Params:                   baseLiquidityPoolEntry.Body.ConstantProduct.Params,
					ReserveA:                 baseLiquidityPoolEntry.Body.ConstantProduct.ReserveA - 11,
					ReserveB:                 baseLiquidityPoolEntry.Body.ConstantProduct.ReserveB - 6,
					TotalPoolShares:          baseLiquidityPoolEntry.Body.ConstantProduct.TotalPoolShares - 10,
					PoolSharesTrustLineCount: baseLiquidityPoolEntry.Body.ConstantProduct.PoolSharesTrustLineCount,
				}),
			},
			expected: []effect{
				{
					effectType:  history.EffectLiquidityPoolWithdrew,
					order:       1,
					address:     "GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
					operationID: 4294967297,
					details: map[string]interface{}{
						"liquidity_pool": map[string]interface{}{
							"fee_bp": uint32(20),
							"id":     poolIDStr,
							"reserves": []base.AssetAmount{
								{
									"native",
									"0.0000189",
								},
								{
									"USD:GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
									"0.0000094",
								},
							},
							"total_shares":     "0.0000990",
							"total_trustlines": "10",
							"type":             "constant_product",
						},
						"reserves_received": []base.AssetAmount{
							{
								"native",
								"0.0000011",
							},
							{
								"USD:GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
								"0.0000006",
							},
						},
						"shares_redeemed": "0.0000010",
					},
				},
			},
		},
		{
			desc: "liquidity pool trade",
			op: xdr.OperationBody{
				Type: xdr.OperationTypePathPaymentStrictSend,
				PathPaymentStrictSendOp: &xdr.PathPaymentStrictSendOp{
					SendAsset:   xdr.MustNewNativeAsset(),
					SendAmount:  10,
					Destination: xdr.MustMuxedAddress("GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY"),
					DestAsset:   usdAsset,
					DestMin:     5,
					Path:        nil,
				},
			},
			changes: xdr.LedgerEntryChanges{
				baseState,
				updateState(xdr.LiquidityPoolEntryConstantProduct{

					Params:                   baseLiquidityPoolEntry.Body.ConstantProduct.Params,
					ReserveA:                 baseLiquidityPoolEntry.Body.ConstantProduct.ReserveA - 11,
					ReserveB:                 baseLiquidityPoolEntry.Body.ConstantProduct.ReserveB - 6,
					TotalPoolShares:          baseLiquidityPoolEntry.Body.ConstantProduct.TotalPoolShares - 10,
					PoolSharesTrustLineCount: baseLiquidityPoolEntry.Body.ConstantProduct.PoolSharesTrustLineCount,
				}),
			},
			result: xdr.OperationResult{
				Code: xdr.OperationResultCodeOpInner,
				Tr: &xdr.OperationResultTr{
					Type: xdr.OperationTypePathPaymentStrictSend,
					PathPaymentStrictSendResult: &xdr.PathPaymentStrictSendResult{
						Code: xdr.PathPaymentStrictSendResultCodePathPaymentStrictSendSuccess,
						Success: &xdr.PathPaymentStrictSendResultSuccess{
							Last: xdr.SimplePaymentResult{
								Destination: xdr.MustAddress("GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY"),
								Asset:       xdr.MustNewCreditAsset("USD", "GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY"),
								Amount:      5,
							},
							Offers: []xdr.ClaimAtom{
								{
									Type: xdr.ClaimAtomTypeClaimAtomTypeLiquidityPool,
									LiquidityPool: &xdr.ClaimLiquidityAtom{
										LiquidityPoolId: poolID,
										AssetSold:       xdr.MustNewNativeAsset(),
										AmountSold:      10,
										AssetBought:     xdr.MustNewCreditAsset("USD", "GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY"),
										AmountBought:    5,
									},
								},
							},
						},
					},
				},
			},
			expected: []effect{
				{
					effectType:  history.EffectAccountCredited,
					order:       1,
					address:     "GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
					operationID: 4294967297,
					details: map[string]interface{}{
						"amount":       "0.0000005",
						"asset_code":   "USD",
						"asset_issuer": "GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
						"asset_type":   "credit_alphanum4",
					},
				},
				{
					effectType:  history.EffectAccountDebited,
					order:       2,
					address:     "GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
					operationID: 4294967297,
					details: map[string]interface{}{
						"amount":     "0.0000010",
						"asset_type": "native",
					},
				},
				{
					effectType:  history.EffectLiquidityPoolTrade,
					order:       3,
					address:     "GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
					operationID: 4294967297,
					details: map[string]interface{}{
						"bought": map[string]string{
							"amount": "0.0000005",
							"asset":  "USD:GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
						},
						"liquidity_pool": map[string]interface{}{
							"fee_bp": uint32(20),
							"id":     poolIDStr,
							"reserves": []base.AssetAmount{
								{
									"native",
									"0.0000189",
								},
								{
									"USD:GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
									"0.0000094",
								},
							},
							"total_shares":     "0.0000990",
							"total_trustlines": "10",
							"type":             "constant_product",
						},
						"sold": map[string]string{
							"amount": "0.0000010",
							"asset":  "native",
						},
					},
				},
			},
		},
		{
			desc: "liquidity pool revocation",
			// Deauthorize an asset
			//
			// This scenario assumes that the asset being deauthorized is also part of a liquidity pool trustline
			// from the same account. This results in a revocation (with a claimable balance being created).
			//
			// This scenario also assumes that the liquidity pool trustline was the last one, cause a liquidity pool removal.
			op: xdr.OperationBody{
				Type: xdr.OperationTypeSetTrustLineFlags,
				SetTrustLineFlagsOp: &xdr.SetTrustLineFlagsOp{
					Trustor:    xdr.MustAddress("GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY"),
					Asset:      usdAsset,
					ClearFlags: xdr.Uint32(xdr.TrustLineFlagsAuthorizedFlag),
				},
			},
			changes: xdr.LedgerEntryChanges{
				// Asset trustline update
				{
					Type: xdr.LedgerEntryChangeTypeLedgerEntryState,
					State: &xdr.LedgerEntry{
						LastModifiedLedgerSeq: 20,
						Data: xdr.LedgerEntryData{
							Type: xdr.LedgerEntryTypeTrustline,
							TrustLine: &xdr.TrustLineEntry{
								AccountId: xdr.MustAddress("GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY"),
								Asset:     usdAsset.ToTrustLineAsset(),
								Balance:   5,
								Limit:     100,
								Flags:     xdr.Uint32(xdr.TrustLineFlagsAuthorizedFlag),
							},
						},
					},
				},
				{
					Type: xdr.LedgerEntryChangeTypeLedgerEntryUpdated,
					Updated: &xdr.LedgerEntry{
						LastModifiedLedgerSeq: 20,
						Data: xdr.LedgerEntryData{
							Type: xdr.LedgerEntryTypeTrustline,
							TrustLine: &xdr.TrustLineEntry{
								AccountId: xdr.MustAddress("GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY"),
								Asset:     usdAsset.ToTrustLineAsset(),
								Balance:   5,
								Limit:     100,
								Flags:     0,
							},
						},
					},
				},
				// Liquidity pool trustline removal
				{
					Type: xdr.LedgerEntryChangeTypeLedgerEntryState,
					State: &xdr.LedgerEntry{
						LastModifiedLedgerSeq: 20,
						Data: xdr.LedgerEntryData{
							Type: xdr.LedgerEntryTypeTrustline,
							TrustLine: &xdr.TrustLineEntry{
								AccountId: xdr.MustAddress("GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY"),
								Asset: xdr.TrustLineAsset{
									Type:            xdr.AssetTypeAssetTypePoolShare,
									LiquidityPoolId: &poolID,
								},
								Balance: 1000,
								Limit:   2000,
								Flags:   xdr.Uint32(xdr.TrustLineFlagsAuthorizedFlag),
							},
						},
					},
				},
				{
					Type: xdr.LedgerEntryChangeTypeLedgerEntryRemoved,
					Removed: &xdr.LedgerKey{
						Type: xdr.LedgerEntryTypeTrustline,
						TrustLine: &xdr.LedgerKeyTrustLine{
							AccountId: xdr.MustAddress("GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY"),
							Asset: xdr.TrustLineAsset{
								Type:            xdr.AssetTypeAssetTypePoolShare,
								LiquidityPoolId: &poolID,
							},
						},
					},
				},
				// create claimable balance for USD asset as part of the revocation (in reality there would probably be another claimable
				// balance crested for the native asset, but let's keep this simple)
				{
					Type: xdr.LedgerEntryChangeTypeLedgerEntryCreated,
					Created: &xdr.LedgerEntry{
						LastModifiedLedgerSeq: 20,
						Data: xdr.LedgerEntryData{
							Type: xdr.LedgerEntryTypeClaimableBalance,
							ClaimableBalance: &xdr.ClaimableBalanceEntry{
								BalanceId: xdr.ClaimableBalanceId{
									Type: xdr.ClaimableBalanceIdTypeClaimableBalanceIdTypeV0,
									V0:   &xdr.Hash{0xa, 0xb},
								},
								Claimants: []xdr.Claimant{
									{
										Type: xdr.ClaimantTypeClaimantTypeV0,
										V0: &xdr.ClaimantV0{
											Destination: xdr.MustAddress("GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY"),
											Predicate: xdr.ClaimPredicate{
												Type: xdr.ClaimPredicateTypeClaimPredicateUnconditional,
											},
										},
									},
								},
								Asset:  usdAsset,
								Amount: 100,
							},
						},
					},
				},
				// Liquidity pool removal
				baseState,
				{
					Type: xdr.LedgerEntryChangeTypeLedgerEntryRemoved,
					Removed: &xdr.LedgerKey{
						Type: xdr.LedgerEntryTypeLiquidityPool,
						LiquidityPool: &xdr.LedgerKeyLiquidityPool{
							LiquidityPoolId: poolID,
						},
					},
				},
			},
			expected: []effect{
				{
					effectType:  history.EffectTrustlineFlagsUpdated,
					order:       1,
					address:     "GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
					operationID: 4294967297,
					details: map[string]interface{}{
						"asset_code":      "USD",
						"asset_issuer":    "GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
						"asset_type":      "credit_alphanum4",
						"authorized_flag": false,
						"trustor":         "GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
					},
				},
				{
					effectType:  history.EffectClaimableBalanceCreated,
					order:       2,
					address:     "GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
					operationID: 4294967297,
					details: map[string]interface{}{
						"amount":     "0.0000100",
						"asset":      "USD:GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
						"balance_id": "000000000a0b000000000000000000000000000000000000000000000000000000000000",
					},
				},
				{
					effectType:  history.EffectClaimableBalanceClaimantCreated,
					order:       3,
					address:     "GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
					operationID: 4294967297,
					details: map[string]interface{}{
						"amount":     "0.0000100",
						"asset":      "USD:GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
						"balance_id": "000000000a0b000000000000000000000000000000000000000000000000000000000000",
						"predicate":  xdr.ClaimPredicate{},
					},
				},
				{
					effectType:  history.EffectLiquidityPoolRevoked,
					order:       4,
					address:     "GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
					operationID: 4294967297,
					details: map[string]interface{}{
						"liquidity_pool": map[string]interface{}{
							"fee_bp": uint32(20),
							"id":     poolIDStr,
							"reserves": []base.AssetAmount{
								{
									"native",
									"0.0000200",
								},
								{
									"USD:GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
									"0.0000100",
								},
							},
							"total_shares":     "0.0001000",
							"total_trustlines": "10",
							"type":             "constant_product",
						},
						"reserves_revoked": []map[string]string{
							{
								"amount":               "0.0000100",
								"asset":                "USD:GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
								"claimable_balance_id": "000000000a0b000000000000000000000000000000000000000000000000000000000000",
							},
						},
						"shares_revoked": "0.0001000",
					},
				},
				{
					effectType:  history.EffectLiquidityPoolRemoved,
					order:       5,
					address:     "GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
					operationID: 4294967297,
					details: map[string]interface{}{
						"liquidity_pool_id": poolIDStr,
					},
				},
			},
		},
	}
	for _, tc := range testCases {

		op := xdr.Operation{Body: tc.op}
		tx := ingest.LedgerTransaction{
			Index: 0,
			Envelope: xdr.TransactionEnvelope{
				Type: xdr.EnvelopeTypeEnvelopeTypeTx,
				V1: &xdr.TransactionV1Envelope{
					Tx: xdr.Transaction{
						SourceAccount: source,
						Operations:    []xdr.Operation{op},
					},
				},
			},
			Result: xdr.TransactionResultPair{
				Result: xdr.TransactionResult{
					Result: xdr.TransactionResultResult{
						Results: &[]xdr.OperationResult{
							tc.result,
						},
					},
				},
			},
			UnsafeMeta: xdr.TransactionMeta{
				V: 2,
				V2: &xdr.TransactionMetaV2{
					Operations: []xdr.OperationMeta{{Changes: tc.changes}},
				},
			},
		}

		t.Run(tc.desc, func(t *testing.T) {
			operation := transactionOperationWrapper{
				index:          0,
				transaction:    tx,
				operation:      op,
				ledgerSequence: 1,
			}

			effects, err := operation.effects()
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, effects)
		})
	}

}
