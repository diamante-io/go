package ingest

import (
	"testing"

	"github.com/diamcircle/go/xdr"
	"github.com/stretchr/testify/assert"
)

func TestChangeAccountChangedExceptSignersInvalidType(t *testing.T) {
	change := Change{
		Type: xdr.LedgerEntryTypeOffer,
	}

	assert.Panics(t, func() {
		change.AccountChangedExceptSigners()
	})
}

func TestFeeMetaAndOperationsChangesSeparate(t *testing.T) {
	tx := LedgerTransaction{
		FeeChanges: xdr.LedgerEntryChanges{
			xdr.LedgerEntryChange{
				Type: xdr.LedgerEntryChangeTypeLedgerEntryState,
				State: &xdr.LedgerEntry{
					Data: xdr.LedgerEntryData{
						Type: xdr.LedgerEntryTypeAccount,
						Account: &xdr.AccountEntry{
							AccountId: xdr.MustAddress("GAHK7EEG2WWHVKDNT4CEQFZGKF2LGDSW2IVM4S5DP42RBW3K6BTODB4A"),
							Balance:   100,
						},
					},
				},
			},
			xdr.LedgerEntryChange{
				Type: xdr.LedgerEntryChangeTypeLedgerEntryUpdated,
				Updated: &xdr.LedgerEntry{
					Data: xdr.LedgerEntryData{
						Type: xdr.LedgerEntryTypeAccount,
						Account: &xdr.AccountEntry{
							AccountId: xdr.MustAddress("GAHK7EEG2WWHVKDNT4CEQFZGKF2LGDSW2IVM4S5DP42RBW3K6BTODB4A"),
							Balance:   200,
						},
					},
				},
			},
		},
		UnsafeMeta: xdr.TransactionMeta{
			V: 1,
			V1: &xdr.TransactionMetaV1{
				Operations: []xdr.OperationMeta{
					{
						Changes: xdr.LedgerEntryChanges{
							xdr.LedgerEntryChange{
								Type: xdr.LedgerEntryChangeTypeLedgerEntryState,
								State: &xdr.LedgerEntry{
									Data: xdr.LedgerEntryData{
										Type: xdr.LedgerEntryTypeAccount,
										Account: &xdr.AccountEntry{
											AccountId: xdr.MustAddress("GAHK7EEG2WWHVKDNT4CEQFZGKF2LGDSW2IVM4S5DP42RBW3K6BTODB4A"),
											Balance:   300,
										},
									},
								},
							},
							xdr.LedgerEntryChange{
								Type: xdr.LedgerEntryChangeTypeLedgerEntryUpdated,
								Updated: &xdr.LedgerEntry{
									Data: xdr.LedgerEntryData{
										Type: xdr.LedgerEntryTypeAccount,
										Account: &xdr.AccountEntry{
											AccountId: xdr.MustAddress("GAHK7EEG2WWHVKDNT4CEQFZGKF2LGDSW2IVM4S5DP42RBW3K6BTODB4A"),
											Balance:   400,
										},
									},
								},
							},
						},
					},
				},
			},
		}}

	feeChanges := tx.GetFeeChanges()
	assert.Len(t, feeChanges, 1)
	assert.Equal(t, feeChanges[0].Pre.Data.MustAccount().Balance, xdr.Int64(100))
	assert.Equal(t, feeChanges[0].Post.Data.MustAccount().Balance, xdr.Int64(200))

	metaChanges, err := tx.GetChanges()
	assert.NoError(t, err)
	assert.Len(t, metaChanges, 1)
	assert.Equal(t, metaChanges[0].Pre.Data.MustAccount().Balance, xdr.Int64(300))
	assert.Equal(t, metaChanges[0].Post.Data.MustAccount().Balance, xdr.Int64(400))

	operationChanges, err := tx.GetOperationChanges(0)
	assert.NoError(t, err)
	assert.Len(t, operationChanges, 1)
	assert.Equal(t, operationChanges[0].Pre.Data.MustAccount().Balance, xdr.Int64(300))
	assert.Equal(t, operationChanges[0].Post.Data.MustAccount().Balance, xdr.Int64(400))

	// Ignore operation meta if tx result is txInternalError
	// https://github.com/diamcircle/go/issues/2111
	tx.Result.Result.Result.Code = xdr.TransactionResultCodeTxInternalError
	metaChanges, err = tx.GetChanges()
	assert.NoError(t, err)
	assert.Len(t, metaChanges, 0)

	operationChanges, err = tx.GetOperationChanges(0)
	assert.NoError(t, err)
	assert.Len(t, operationChanges, 0)
}

func TestFailedTransactionOperationChangesMeta(t *testing.T) {
	testCases := []struct {
		desc string
		meta xdr.TransactionMeta
	}{
		{
			desc: "V0",
			meta: xdr.TransactionMeta{
				Operations: &[]xdr.OperationMeta{},
			},
		},
		{
			desc: "V1",
			meta: xdr.TransactionMeta{
				V:  1,
				V1: &xdr.TransactionMetaV1{},
			},
		},
		{
			desc: "V2",
			meta: xdr.TransactionMeta{
				V:  2,
				V2: &xdr.TransactionMetaV2{},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tx := LedgerTransaction{
				Result: xdr.TransactionResultPair{
					Result: xdr.TransactionResult{
						Result: xdr.TransactionResultResult{
							Code: xdr.TransactionResultCodeTxFailed,
						},
					},
				},
				UnsafeMeta: tc.meta,
			}

			operationChanges, err := tx.GetOperationChanges(0)
			if tx.UnsafeMeta.V == 0 {
				assert.Error(t, err)
				assert.EqualError(t, err, "TransactionMeta.V=0 not supported")
			} else {
				assert.NoError(t, err)
				assert.Len(t, operationChanges, 0)
			}
		})
	}
}
func TestMetaV2Order(t *testing.T) {
	tx := LedgerTransaction{
		UnsafeMeta: xdr.TransactionMeta{
			V: 2,
			V2: &xdr.TransactionMetaV2{
				TxChangesBefore: xdr.LedgerEntryChanges{
					xdr.LedgerEntryChange{
						Type: xdr.LedgerEntryChangeTypeLedgerEntryState,
						State: &xdr.LedgerEntry{
							Data: xdr.LedgerEntryData{
								Type: xdr.LedgerEntryTypeAccount,
								Account: &xdr.AccountEntry{
									AccountId: xdr.MustAddress("GACMZD5VJXTRLKVET72CETCYKELPNCOTTBDC6DHFEUPLG5DHEK534JQX"),
									Balance:   100,
								},
							},
						},
					},
					xdr.LedgerEntryChange{
						Type: xdr.LedgerEntryChangeTypeLedgerEntryUpdated,
						Updated: &xdr.LedgerEntry{
							Data: xdr.LedgerEntryData{
								Type: xdr.LedgerEntryTypeAccount,
								Account: &xdr.AccountEntry{
									AccountId: xdr.MustAddress("GACMZD5VJXTRLKVET72CETCYKELPNCOTTBDC6DHFEUPLG5DHEK534JQX"),
									Balance:   200,
								},
							},
						},
					},
					xdr.LedgerEntryChange{
						Type: xdr.LedgerEntryChangeTypeLedgerEntryState,
						State: &xdr.LedgerEntry{
							Data: xdr.LedgerEntryData{
								Type: xdr.LedgerEntryTypeAccount,
								Account: &xdr.AccountEntry{
									AccountId: xdr.MustAddress("GAHK7EEG2WWHVKDNT4CEQFZGKF2LGDSW2IVM4S5DP42RBW3K6BTODB4A"),
									Balance:   100,
								},
							},
						},
					},
					xdr.LedgerEntryChange{
						Type: xdr.LedgerEntryChangeTypeLedgerEntryUpdated,
						Updated: &xdr.LedgerEntry{
							Data: xdr.LedgerEntryData{
								Type: xdr.LedgerEntryTypeAccount,
								Account: &xdr.AccountEntry{
									AccountId: xdr.MustAddress("GAHK7EEG2WWHVKDNT4CEQFZGKF2LGDSW2IVM4S5DP42RBW3K6BTODB4A"),
									Balance:   200,
								},
							},
						},
					},
				},
				Operations: []xdr.OperationMeta{
					{
						Changes: xdr.LedgerEntryChanges{
							xdr.LedgerEntryChange{
								Type: xdr.LedgerEntryChangeTypeLedgerEntryState,
								State: &xdr.LedgerEntry{
									Data: xdr.LedgerEntryData{
										Type: xdr.LedgerEntryTypeAccount,
										Account: &xdr.AccountEntry{
											AccountId: xdr.MustAddress("GAHK7EEG2WWHVKDNT4CEQFZGKF2LGDSW2IVM4S5DP42RBW3K6BTODB4A"),
											Balance:   300,
										},
									},
								},
							},
							xdr.LedgerEntryChange{
								Type: xdr.LedgerEntryChangeTypeLedgerEntryUpdated,
								Updated: &xdr.LedgerEntry{
									Data: xdr.LedgerEntryData{
										Type: xdr.LedgerEntryTypeAccount,
										Account: &xdr.AccountEntry{
											AccountId: xdr.MustAddress("GAHK7EEG2WWHVKDNT4CEQFZGKF2LGDSW2IVM4S5DP42RBW3K6BTODB4A"),
											Balance:   400,
										},
									},
								},
							},
						},
					},
				},
				TxChangesAfter: xdr.LedgerEntryChanges{
					xdr.LedgerEntryChange{
						Type: xdr.LedgerEntryChangeTypeLedgerEntryState,
						State: &xdr.LedgerEntry{
							Data: xdr.LedgerEntryData{
								Type: xdr.LedgerEntryTypeAccount,
								Account: &xdr.AccountEntry{
									AccountId: xdr.MustAddress("GACMZD5VJXTRLKVET72CETCYKELPNCOTTBDC6DHFEUPLG5DHEK534JQX"),
									Balance:   300,
								},
							},
						},
					},
					xdr.LedgerEntryChange{
						Type: xdr.LedgerEntryChangeTypeLedgerEntryUpdated,
						Updated: &xdr.LedgerEntry{
							Data: xdr.LedgerEntryData{
								Type: xdr.LedgerEntryTypeAccount,
								Account: &xdr.AccountEntry{
									AccountId: xdr.MustAddress("GACMZD5VJXTRLKVET72CETCYKELPNCOTTBDC6DHFEUPLG5DHEK534JQX"),
									Balance:   400,
								},
							},
						},
					},
				},
			},
		}}

	metaChanges, err := tx.GetChanges()
	assert.NoError(t, err)
	assert.Len(t, metaChanges, 4)

	change := metaChanges[0]
	id := change.Pre.Data.MustAccount().AccountId
	assert.Equal(t, id.Address(), "GACMZD5VJXTRLKVET72CETCYKELPNCOTTBDC6DHFEUPLG5DHEK534JQX")
	assert.Equal(t, change.Pre.Data.MustAccount().Balance, xdr.Int64(100))
	assert.Equal(t, change.Post.Data.MustAccount().Balance, xdr.Int64(200))

	change = metaChanges[1]
	id = change.Pre.Data.MustAccount().AccountId
	assert.Equal(t, id.Address(), "GAHK7EEG2WWHVKDNT4CEQFZGKF2LGDSW2IVM4S5DP42RBW3K6BTODB4A")
	assert.Equal(t, change.Pre.Data.MustAccount().Balance, xdr.Int64(100))
	assert.Equal(t, change.Post.Data.MustAccount().Balance, xdr.Int64(200))

	change = metaChanges[2]
	id = change.Pre.Data.MustAccount().AccountId
	assert.Equal(t, id.Address(), "GAHK7EEG2WWHVKDNT4CEQFZGKF2LGDSW2IVM4S5DP42RBW3K6BTODB4A")
	assert.Equal(t, change.Pre.Data.MustAccount().Balance, xdr.Int64(300))
	assert.Equal(t, change.Post.Data.MustAccount().Balance, xdr.Int64(400))

	change = metaChanges[3]
	id = change.Pre.Data.MustAccount().AccountId
	assert.Equal(t, id.Address(), "GACMZD5VJXTRLKVET72CETCYKELPNCOTTBDC6DHFEUPLG5DHEK534JQX")
	assert.Equal(t, change.Pre.Data.MustAccount().Balance, xdr.Int64(300))
	assert.Equal(t, change.Post.Data.MustAccount().Balance, xdr.Int64(400))

	operationChanges, err := tx.GetOperationChanges(0)
	assert.NoError(t, err)
	assert.Len(t, operationChanges, 1)

	// Ignore operations meta and txChangesAfter if txInternalError
	// https://github.com/diamcircle/go/issues/2111
	tx.Result.Result.Result.Code = xdr.TransactionResultCodeTxInternalError
	metaChanges, err = tx.GetChanges()
	assert.NoError(t, err)
	assert.Len(t, metaChanges, 2)

	change = metaChanges[0]
	id = change.Pre.Data.MustAccount().AccountId
	assert.Equal(t, id.Address(), "GACMZD5VJXTRLKVET72CETCYKELPNCOTTBDC6DHFEUPLG5DHEK534JQX")
	assert.Equal(t, change.Pre.Data.MustAccount().Balance, xdr.Int64(100))
	assert.Equal(t, change.Post.Data.MustAccount().Balance, xdr.Int64(200))

	change = metaChanges[1]
	id = change.Pre.Data.MustAccount().AccountId
	assert.Equal(t, id.Address(), "GAHK7EEG2WWHVKDNT4CEQFZGKF2LGDSW2IVM4S5DP42RBW3K6BTODB4A")
	assert.Equal(t, change.Pre.Data.MustAccount().Balance, xdr.Int64(100))
	assert.Equal(t, change.Post.Data.MustAccount().Balance, xdr.Int64(200))

	operationChanges, err = tx.GetOperationChanges(0)
	assert.NoError(t, err)
	assert.Len(t, operationChanges, 0)

}

func TestMetaV0(t *testing.T) {
	tx := LedgerTransaction{
		UnsafeMeta: xdr.TransactionMeta{
			V: 0,
		}}

	_, err := tx.GetChanges()
	assert.Error(t, err)
	assert.EqualError(t, err, "TransactionMeta.V=0 not supported")
}

func TestChangeAccountChangedExceptSignersLastModifiedLedgerSeq(t *testing.T) {
	change := Change{
		Type: xdr.LedgerEntryTypeAccount,
		Pre: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 10,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
				},
			},
		},
		Post: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 11,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
				},
			},
		},
	}
	changed, err := change.AccountChangedExceptSigners()
	assert.NoError(t, err)
	assert.True(t, changed)
}

func TestChangeAccountChangedExceptSignersNoPre(t *testing.T) {
	change := Change{
		Type: xdr.LedgerEntryTypeAccount,
		Pre:  nil,
		Post: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 10,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
				},
			},
		},
	}
	changed, err := change.AccountChangedExceptSigners()
	assert.NoError(t, err)
	assert.True(t, changed)
}

func TestChangeAccountChangedExceptSignersNoPost(t *testing.T) {
	change := Change{
		Type: xdr.LedgerEntryTypeAccount,
		Pre: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 10,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
				},
			},
		},
		Post: nil,
	}
	changed, err := change.AccountChangedExceptSigners()
	assert.NoError(t, err)
	assert.True(t, changed)
}

func TestChangeAccountChangedExceptSignersMasterKeyRemoved(t *testing.T) {
	change := Change{
		Type: xdr.LedgerEntryTypeAccount,
		Pre: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 10,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					// Master weight = 1
					Thresholds: [4]byte{1, 1, 1, 1},
				},
			},
		},
		Post: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 10,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					// Master weight = 0
					Thresholds: [4]byte{0, 1, 1, 1},
				},
			},
		},
	}

	changed, err := change.AccountChangedExceptSigners()
	assert.NoError(t, err)
	assert.True(t, changed)
}

func TestChangeAccountChangedExceptSignersSignerChange(t *testing.T) {
	change := Change{
		Type: xdr.LedgerEntryTypeAccount,
		Pre: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 10,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					Signers: []xdr.Signer{
						{
							Key:    xdr.MustSigner("GCCCU34WDY2RATQTOOQKY6SZWU6J5DONY42SWGW2CIXGW4LICAGNRZKX"),
							Weight: 1,
						},
					},
				},
			},
		},
		Post: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 10,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					Signers: []xdr.Signer{
						{
							Key:    xdr.MustSigner("GCCCU34WDY2RATQTOOQKY6SZWU6J5DONY42SWGW2CIXGW4LICAGNRZKX"),
							Weight: 2,
						},
					},
				},
			},
		},
	}

	changed, err := change.AccountChangedExceptSigners()
	assert.NoError(t, err)
	assert.False(t, changed)
}

func TestChangeAccountChangedExceptSignersNoChanges(t *testing.T) {
	inflationDest := xdr.MustAddress("GBAH2GBLJB54JAROJ3FVO4ZTTJJI3XKOBTMJOZFUJ3UHYIVNJTLJUYFY")
	change := Change{
		Type: xdr.LedgerEntryTypeAccount,
		Pre: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 10,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId:     xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					Balance:       1000,
					SeqNum:        432894732,
					NumSubEntries: 2,
					InflationDest: &inflationDest,
					Flags:         4,
					HomeDomain:    "diamcircle.org",
					Thresholds:    [4]byte{1, 1, 1, 1},
					Signers: []xdr.Signer{
						{
							Key:    xdr.MustSigner("GCCCU34WDY2RATQTOOQKY6SZWU6J5DONY42SWGW2CIXGW4LICAGNRZKX"),
							Weight: 1,
						},
					},
					Ext: xdr.AccountEntryExt{
						V: 1,
						V1: &xdr.AccountEntryExtensionV1{
							Liabilities: xdr.Liabilities{
								Buying:  10,
								Selling: 20,
							},
						},
					},
				},
			},
		},
		Post: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 10,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId:     xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					Balance:       1000,
					SeqNum:        432894732,
					NumSubEntries: 2,
					InflationDest: &inflationDest,
					Flags:         4,
					HomeDomain:    "diamcircle.org",
					Thresholds:    [4]byte{1, 1, 1, 1},
					Signers: []xdr.Signer{
						{
							Key:    xdr.MustSigner("GCCCU34WDY2RATQTOOQKY6SZWU6J5DONY42SWGW2CIXGW4LICAGNRZKX"),
							Weight: 1,
						},
					},
					Ext: xdr.AccountEntryExt{
						V: 1,
						V1: &xdr.AccountEntryExtensionV1{
							Liabilities: xdr.Liabilities{
								Buying:  10,
								Selling: 20,
							},
						},
					},
				},
			},
		},
	}

	changed, err := change.AccountChangedExceptSigners()
	assert.NoError(t, err)
	assert.False(t, changed)

	// Make sure pre and post not modified
	assert.NotNil(t, change.Pre.Data.Account.Signers)
	assert.Len(t, change.Pre.Data.Account.Signers, 1)

	assert.NotNil(t, change.Post.Data.Account.Signers)
	assert.Len(t, change.Post.Data.Account.Signers, 1)
}

func TestChangeAccountSignersChangedInvalidType(t *testing.T) {
	change := Change{
		Type: xdr.LedgerEntryTypeOffer,
	}

	assert.Panics(t, func() {
		change.AccountSignersChanged()
	})
}

func TestChangeAccountSignersChangedNoPre(t *testing.T) {
	change := Change{
		Type: xdr.LedgerEntryTypeAccount,
		Pre:  nil,
		Post: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 10,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
				},
			},
		},
	}

	assert.True(t, change.AccountSignersChanged())
}

func TestChangeAccountSignersChangedNoPostMasterKey(t *testing.T) {
	change := Change{
		Type: xdr.LedgerEntryTypeAccount,
		Pre: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 10,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					// Master weight = 1
					Thresholds: [4]byte{1, 1, 1, 1},
				},
			},
		},
		Post: nil,
	}

	assert.True(t, change.AccountSignersChanged())
}

func TestChangeAccountSignersChangedNoPostNoMasterKey(t *testing.T) {
	change := Change{
		Type: xdr.LedgerEntryTypeAccount,
		Pre: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 10,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					// Master weight = 0
					Thresholds: [4]byte{0, 1, 1, 1},
				},
			},
		},
		Post: nil,
	}

	// Account being merge can still have signers so they will be removed.
	assert.True(t, change.AccountSignersChanged())
}

func TestChangeAccountSignersChangedMasterKeyRemoved(t *testing.T) {
	change := Change{
		Type: xdr.LedgerEntryTypeAccount,
		Pre: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 10,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					// Master weight = 1
					Thresholds: [4]byte{1, 1, 1, 1},
				},
			},
		},
		Post: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 10,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					// Master weight = 0
					Thresholds: [4]byte{0, 1, 1, 1},
				},
			},
		},
	}

	assert.True(t, change.AccountSignersChanged())
}

func TestChangeAccountSignersChangedMasterKeyAdded(t *testing.T) {
	change := Change{
		Type: xdr.LedgerEntryTypeAccount,
		Pre: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 10,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					// Master weight = 0
					Thresholds: [4]byte{0, 1, 1, 1},
				},
			},
		},
		Post: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 10,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					// Master weight = 1
					Thresholds: [4]byte{1, 1, 1, 1},
				},
			},
		},
	}

	assert.True(t, change.AccountSignersChanged())
}

func TestChangeAccountSignersChangedSignerAdded(t *testing.T) {
	change := Change{
		Type: xdr.LedgerEntryTypeAccount,
		Pre: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 10,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					Signers:   []xdr.Signer{},
				},
			},
		},
		Post: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 10,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					Signers: []xdr.Signer{
						{
							Key:    xdr.MustSigner("GCCCU34WDY2RATQTOOQKY6SZWU6J5DONY42SWGW2CIXGW4LICAGNRZKX"),
							Weight: 1,
						},
					},
				},
			},
		},
	}

	assert.True(t, change.AccountSignersChanged())
}

func TestChangeAccountSignersChangedSignerRemoved(t *testing.T) {
	change := Change{
		Type: xdr.LedgerEntryTypeAccount,
		Pre: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 10,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					Signers: []xdr.Signer{
						{
							Key:    xdr.MustSigner("GCCCU34WDY2RATQTOOQKY6SZWU6J5DONY42SWGW2CIXGW4LICAGNRZKX"),
							Weight: 1,
						},
					},
				},
			},
		},
		Post: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 10,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					Signers:   []xdr.Signer{},
				},
			},
		},
	}

	assert.True(t, change.AccountSignersChanged())
}

func TestChangeAccountSignersChangedSignerWeightChanged(t *testing.T) {
	change := Change{
		Type: xdr.LedgerEntryTypeAccount,
		Pre: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 10,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					Signers: []xdr.Signer{
						{
							Key:    xdr.MustSigner("GCCCU34WDY2RATQTOOQKY6SZWU6J5DONY42SWGW2CIXGW4LICAGNRZKX"),
							Weight: 1,
						},
					},
				},
			},
		},
		Post: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 10,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					Signers: []xdr.Signer{
						{
							Key:    xdr.MustSigner("GCCCU34WDY2RATQTOOQKY6SZWU6J5DONY42SWGW2CIXGW4LICAGNRZKX"),
							Weight: 2,
						},
					},
				},
			},
		},
	}

	assert.True(t, change.AccountSignersChanged())
}

func TestChangeAccountSignersChangedSponsorAdded(t *testing.T) {
	sponsor, err := xdr.AddressToAccountId("GBADGWKHSUFOC4C7E3KXKINZSRX5KPHUWHH67UGJU77LEORGVLQ3BN3B")
	assert.NoError(t, err)

	change := Change{
		Type: xdr.LedgerEntryTypeAccount,
		Pre: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 10,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					Signers: []xdr.Signer{
						{
							Key:    xdr.MustSigner("GCCCU34WDY2RATQTOOQKY6SZWU6J5DONY42SWGW2CIXGW4LICAGNRZKX"),
							Weight: 1,
						},
					},
				},
			},
		},
		Post: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 10,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					Signers: []xdr.Signer{
						{
							Key:    xdr.MustSigner("GCCCU34WDY2RATQTOOQKY6SZWU6J5DONY42SWGW2CIXGW4LICAGNRZKX"),
							Weight: 1,
						},
					},
					Ext: xdr.AccountEntryExt{
						V1: &xdr.AccountEntryExtensionV1{
							Ext: xdr.AccountEntryExtensionV1Ext{
								V2: &xdr.AccountEntryExtensionV2{
									SignerSponsoringIDs: []xdr.SponsorshipDescriptor{
										&sponsor,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	assert.True(t, change.AccountSignersChanged())
}

func TestChangeAccountSignersChangedSponsorRemoved(t *testing.T) {
	sponsor, err := xdr.AddressToAccountId("GBADGWKHSUFOC4C7E3KXKINZSRX5KPHUWHH67UGJU77LEORGVLQ3BN3B")
	assert.NoError(t, err)

	change := Change{
		Type: xdr.LedgerEntryTypeAccount,
		Pre: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 10,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					Signers: []xdr.Signer{
						{
							Key:    xdr.MustSigner("GCCCU34WDY2RATQTOOQKY6SZWU6J5DONY42SWGW2CIXGW4LICAGNRZKX"),
							Weight: 1,
						},
					},
					Ext: xdr.AccountEntryExt{
						V1: &xdr.AccountEntryExtensionV1{
							Ext: xdr.AccountEntryExtensionV1Ext{
								V2: &xdr.AccountEntryExtensionV2{
									SignerSponsoringIDs: []xdr.SponsorshipDescriptor{
										&sponsor,
									},
								},
							},
						},
					},
				},
			},
		},
		Post: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 10,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					Signers: []xdr.Signer{
						{
							Key:    xdr.MustSigner("GCCCU34WDY2RATQTOOQKY6SZWU6J5DONY42SWGW2CIXGW4LICAGNRZKX"),
							Weight: 1,
						},
					},
				},
			},
		},
	}

	assert.True(t, change.AccountSignersChanged())
}

func TestChangeAccountSignersChangedSponsorChanged(t *testing.T) {
	sponsor, err := xdr.AddressToAccountId("GBADGWKHSUFOC4C7E3KXKINZSRX5KPHUWHH67UGJU77LEORGVLQ3BN3B")
	assert.NoError(t, err)

	newSponsor, err := xdr.AddressToAccountId("GB2Y6D5QFDJSCR6GSBO5D2LOLGZI4RVPRGZSSPLIFWNJZ7SL73TOMXAQ")
	assert.NoError(t, err)

	change := Change{
		Type: xdr.LedgerEntryTypeAccount,
		Pre: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 10,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					Signers: []xdr.Signer{
						{
							Key:    xdr.MustSigner("GCCCU34WDY2RATQTOOQKY6SZWU6J5DONY42SWGW2CIXGW4LICAGNRZKX"),
							Weight: 1,
						},
					},
					Ext: xdr.AccountEntryExt{
						V1: &xdr.AccountEntryExtensionV1{
							Ext: xdr.AccountEntryExtensionV1Ext{
								V2: &xdr.AccountEntryExtensionV2{
									SignerSponsoringIDs: []xdr.SponsorshipDescriptor{
										&sponsor,
									},
								},
							},
						},
					},
				},
			},
		},
		Post: &xdr.LedgerEntry{
			LastModifiedLedgerSeq: 10,
			Data: xdr.LedgerEntryData{
				Type: xdr.LedgerEntryTypeAccount,
				Account: &xdr.AccountEntry{
					AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
					Signers: []xdr.Signer{
						{
							Key:    xdr.MustSigner("GCCCU34WDY2RATQTOOQKY6SZWU6J5DONY42SWGW2CIXGW4LICAGNRZKX"),
							Weight: 1,
						},
					},
					Ext: xdr.AccountEntryExt{
						V1: &xdr.AccountEntryExtensionV1{
							Ext: xdr.AccountEntryExtensionV1Ext{
								V2: &xdr.AccountEntryExtensionV2{
									SignerSponsoringIDs: []xdr.SponsorshipDescriptor{
										&newSponsor,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	assert.True(t, change.AccountSignersChanged())
}
