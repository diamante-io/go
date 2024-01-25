package history

import (
	"testing"

	"go/services/aurora/internal/test"

	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
)

var (
	data1 = Data{
		AccountID: "GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB",
		Name:      "test data",
		// This also tests if base64 encoding is working as 0 is invalid UTF-8 byte
		Value:              []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		LastModifiedLedger: 1234,
	}

	data2 = Data{
		AccountID:          "GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB",
		Name:               "test data2",
		Value:              []byte{10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
		LastModifiedLedger: 1234,
		Sponsor:            null.StringFrom("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
	}
)

func TestInsertAccountData(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &Q{tt.auroraSession()}

	err := q.UpsertAccountData(tt.Ctx, []Data{data1})
	assert.NoError(t, err)

	err = q.UpsertAccountData(tt.Ctx, []Data{data2})
	assert.NoError(t, err)

	keys := []AccountDataKey{
		{AccountID: data1.AccountID, DataName: data1.Name},
		{AccountID: data2.AccountID, DataName: data2.Name},
	}

	datas, err := q.GetAccountDataByKeys(tt.Ctx, keys)
	assert.NoError(t, err)
	assert.Len(t, datas, 2)

	tt.Assert.Equal(data1.Name, datas[0].Name)
	tt.Assert.Equal(data1.Value, datas[0].Value)
	tt.Assert.True(datas[0].Sponsor.IsZero())

	tt.Assert.Equal(data2.Name, datas[1].Name)
	tt.Assert.Equal(data2.Value, datas[1].Value)
	tt.Assert.Equal("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML", datas[1].Sponsor.String)
}

func TestUpdateAccountData(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &Q{tt.auroraSession()}

	err := q.UpsertAccountData(tt.Ctx, []Data{data1})
	assert.NoError(t, err)

	modifiedData := data1
	value2 := make([]byte, len(modifiedData.Value))
	copy(value2, modifiedData.Value)
	value2[0] = 1
	modifiedData.Value = value2

	err = q.UpsertAccountData(tt.Ctx, []Data{modifiedData})
	assert.NoError(t, err)

	keys := []AccountDataKey{
		{AccountID: data1.AccountID, DataName: data1.Name},
	}
	datas, err := q.GetAccountDataByKeys(tt.Ctx, keys)
	assert.NoError(t, err)
	assert.Len(t, datas, 1)

	tt.Assert.Equal(modifiedData.Name, datas[0].Name)
	tt.Assert.Equal(modifiedData.Value, datas[0].Value)
	tt.Assert.Equal(uint32(1234), datas[0].LastModifiedLedger)
}

func TestRemoveAccountData(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &Q{tt.auroraSession()}

	err := q.UpsertAccountData(tt.Ctx, []Data{data1})
	assert.NoError(t, err)

	key := AccountDataKey{AccountID: data1.AccountID, DataName: data1.Name}
	rows, err := q.RemoveAccountData(tt.Ctx, []AccountDataKey{key})
	assert.NoError(t, err)
	tt.Assert.Equal(int64(1), rows)

	datas, err := q.GetAccountDataByKeys(tt.Ctx, []AccountDataKey{key})
	assert.NoError(t, err)
	assert.Len(t, datas, 0)

	// Doesn't exist anymore
	rows, err = q.RemoveAccountData(tt.Ctx, []AccountDataKey{key})
	assert.NoError(t, err)
	tt.Assert.Equal(int64(0), rows)
}

func TestGetAccountDataByAccountsID(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &Q{tt.auroraSession()}

	err := q.UpsertAccountData(tt.Ctx, []Data{data1})
	assert.NoError(t, err)
	err = q.UpsertAccountData(tt.Ctx, []Data{data2})
	assert.NoError(t, err)

	ids := []string{
		data1.AccountID,
		data2.AccountID,
	}
	datas, err := q.GetAccountDataByAccountsID(tt.Ctx, ids)
	assert.NoError(t, err)
	assert.Len(t, datas, 2)

	tt.Assert.Equal(data1.Name, datas[0].Name)
	tt.Assert.Equal(data1.Value, datas[0].Value)

	tt.Assert.Equal(data2.Name, datas[1].Name)
	tt.Assert.Equal(data2.Value, datas[1].Value)
}

func TestGetAccountDataByAccountID(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &Q{tt.auroraSession()}

	err := q.UpsertAccountData(tt.Ctx, []Data{data1})
	assert.NoError(t, err)
	err = q.UpsertAccountData(tt.Ctx, []Data{data2})
	assert.NoError(t, err)

	records, err := q.GetAccountDataByAccountID(tt.Ctx, data1.AccountID)
	assert.NoError(t, err)
	assert.Len(t, records, 2)

	tt.Assert.Equal(data1.Name, records[0].Name)
	tt.Assert.Equal(data1.Value, records[0].Value)

	tt.Assert.Equal(data2.Name, records[1].Name)
	tt.Assert.Equal(data2.Value, records[1].Value)
}

func TestGetAccountDataByName(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &Q{tt.auroraSession()}

	err := q.UpsertAccountData(tt.Ctx, []Data{data1})
	assert.NoError(t, err)
	err = q.UpsertAccountData(tt.Ctx, []Data{data2})
	assert.NoError(t, err)

	record, err := q.GetAccountDataByName(tt.Ctx, data1.AccountID, data1.Name)
	assert.NoError(t, err)
	tt.Assert.Equal(data1.Name, record.Name)
	tt.Assert.Equal(data1.Value, record.Value)

	record, err = q.GetAccountDataByName(tt.Ctx, data1.AccountID, data2.Name)
	assert.NoError(t, err)
	tt.Assert.Equal(data2.Name, record.Name)
	tt.Assert.Equal(data2.Value, record.Value)

}
