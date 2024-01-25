package aurora

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Account Tests
// An example account to be used in all the Account tests
var exampleAccount = Account{
	Data: map[string]string{
		"test":    "aGVsbG8=",
		"invalid": "a_*&^*",
	},
	Sequence: "3002985298788353",
}

func TestAccount_IncrementSequenceNumber(t *testing.T) {
	seqNum, err := exampleAccount.IncrementSequenceNumber()

	assert.Nil(t, err)
	assert.Equal(t, "3002985298788354", exampleAccount.Sequence, "sequence number string was incremented")
	assert.Equal(t, int64(3002985298788354), seqNum, "incremented sequence number is correct value/type")
}

func TestAccount_GetData(t *testing.T) {
	decoded, err := exampleAccount.GetData("test")
	assert.Nil(t, err)
	assert.Equal(t, string(decoded), "hello", "returns decoded value when key exists")

	decoded, err = exampleAccount.GetData("test2")
	assert.Nil(t, err)
	assert.Equal(t, len(decoded), 0, "returns empty slice if key doesn't exist")

	_, err = exampleAccount.GetData("invalid")
	assert.NotNil(t, err, "returns error slice if value is invalid")
}

func TestAccount_MustGetData(t *testing.T) {
	decoded := exampleAccount.MustGetData("test")
	assert.Equal(t, string(decoded), "hello", "returns decoded value when the key exists")

	decoded = exampleAccount.MustGetData("test2")
	assert.Equal(t, len(decoded), 0, "returns empty slice if key doesn't exist")

	assert.Panics(t, func() { exampleAccount.MustGetData("invalid") }, "panics on invalid input")
}

// Transaction Tests
func TestTransactionJSONMarshal(t *testing.T) {
	transaction := Transaction{
		ID:         "12345",
		MaxFee:     11,
		FeeCharged: 10,
		MemoType:   "text",
		Memo:       "",
	}
	marshaledTransaction, marshalErr := json.Marshal(transaction)
	assert.Nil(t, marshalErr)
	var result Transaction
	json.Unmarshal(marshaledTransaction, &result)
	assert.Equal(t, result, transaction, "data matches original input")
}

// Test that a typical friendbot fund response can unmarshal to the Transaction
// type. The auroraclient uses the Transaction type for friendbot responses
// also, but their response is a slimmed down version of the full transaction
// response. This test confirms there are no errors unmarshalling that slimmed
// down version.
func TestTransactionUnmarshalsFriendbotFund(t *testing.T) {
	friendbotFundResponse := `{
  "_links": {
    "transaction": {
      "href": "https://aurora-testnet.diamcircle.org/transactions/94e42f65d3ff5f30669b6109c2ce3e82c0e592c52004e3b41bb30e24df33954e"
    }
  },
  "hash": "94e42f65d3ff5f30669b6109c2ce3e82c0e592c52004e3b41bb30e24df33954e",
  "ledger": 8269,
  "envelope_xdr": "AAAAAgAAAAD2Leuk4afNVCYqxbN03yPH6kgKe/o2yiOd3CQNkpkpQwABhqAAAAFSAAAACQAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAABAAAAABB90WssODNIgi6BHveqzxTRmIpvAFRyVNM+Hm2GVuCcAAAAAAAAAABW9+rbvt6YXwwXyFszptQFlfzzFMrWObLiJmBhOzNblAAAABdIdugAAAAAAAAAAAKSmSlDAAAAQHWNbXOoVQqH0YJRr8LAtpalV+NoXb8Tv/ETkPNv2NignhN8seUSde8m2HLNLHOo+5W34BXfxfBmDXgZn8yHkwSGVuCcAAAAQDQLh1UAxYZ27sIxyYgyYFo8IUbTiANWadUJUR7K0q1eY6Q5J/BFfNlf6UqLqJ5zd8uI3TXCaBNJDkiQc1ZLEg4=",
  "result_xdr": "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=",
  "result_meta_xdr": "AAAAAgAAAAIAAAADAAAgTQAAAAAAAAAA9i3rpOGnzVQmKsWzdN8jx+pICnv6NsojndwkDZKZKUMAAAAAPDNbbAAAAVIAAAAIAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAgTQAAAAAAAAAA9i3rpOGnzVQmKsWzdN8jx+pICnv6NsojndwkDZKZKUMAAAAAPDNbbAAAAVIAAAAJAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAACBMAAAAAAAAAAAQfdFrLDgzSIIugR73qs8U0ZiKbwBUclTTPh5thlbgnAFg09HQY/uMAAAA2wAAAAoAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAACBNAAAAAAAAAAAQfdFrLDgzSIIugR73qs8U0ZiKbwBUclTTPh5thlbgnAFg07qH7ROMAAAA2wAAAAoAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAACBNAAAAAAAAAABW9+rbvt6YXwwXyFszptQFlfzzFMrWObLiJmBhOzNblAAAABdIdugAAAAgTQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAA="
}`
	transaction := Transaction{}
	err := json.Unmarshal([]byte(friendbotFundResponse), &transaction)
	assert.Nil(t, err)
}

func TestTransactionEmptyMemoText(t *testing.T) {
	transaction := Transaction{
		MemoType:  "text",
		Memo:      "",
		MemoBytes: "",
	}
	marshaledTransaction, marshalErr := json.Marshal(transaction)
	assert.Nil(t, marshalErr)
	var result struct {
		Memo      *string
		MemoBytes *string `json:"memo_bytes"`
	}
	json.Unmarshal(marshaledTransaction, &result)
	assert.NotNil(t, result.Memo, "memo field is present even if input memo was empty string")
	assert.NotNil(t, result.MemoBytes, "memo_bytes field is present even if input memo was empty string")
}

func TestTransactionMemoTypeNone(t *testing.T) {
	transaction := Transaction{
		MemoType: "none",
	}
	marshaledTransaction, marshalErr := json.Marshal(transaction)
	assert.Nil(t, marshalErr)
	var result struct {
		Memo *string
	}
	json.Unmarshal(marshaledTransaction, &result)
	assert.Nil(t, result.Memo, "no memo field is present when memo input type was `none`")
}

func TestTransactionUnmarshalJSON(t *testing.T) {
	const feesAsInt64s = `{
        "memo": "MzUyODFmNThmZjkxMGNiMTVhYWQ1NjM2ZGIyNzUzZTA=",
        "_links": {
          "self": {
            "href": "https://diamtestnet.diamcircle.io/transactionsctions/998605ace4a0b89293cf729cf216405f29c1ce5d44d6a40232982a4bdccda033"
          },
          "account": {
            "href": "https://mainnet.diamcircle.io/GBUYDJH3AOPFFND3L54DUDWIHOMYKUONDV4RAHOHDBNN2D5N4BPPWDQ3"
          },
          "ledger": {
            "href": "https://aurora.diamcircle.org/ledgers/29113108"
          },
          "operations": {
            "href": "https://diamtestnet.diamcircle.io/transactionsctions/998605ace4a0b89293cf729cf216405f29c1ce5d44d6a40232982a4bdccda033/operations{?cursor,limit,order}",
            "templated": true
          },
          "effects": {
            "href": "https://diamtestnet.diamcircle.io/transactions/998605ace4a0b89293cf729cf216405f29c1ce5d44d6a40232982a4bdccda033/effects{?cursor,limit,order}",
            "templated": true
          },
          "precedes": {
            "href": "https://diamtestnet.diamcircle.io/transactions?order=asc\u0026cursor=125039846745419776"
          },
          "succeeds": {
            "href": "https://diamtestnet.diamcircle.io/transactions?order=desc\u0026cursor=125039846745419776"
          },
          "transaction": {
            "href": "https://diamtestnet.diamcircle.io/transactions/998605ace4a0b89293cf729cf216405f29c1ce5d44d6a40232982a4bdccda033"
          }
        },
        "id": "998605ace4a0b89293cf729cf216405f29c1ce5d44d6a40232982a4bdccda033",
        "paging_token": "125039846745419776",
        "successful": true,
        "hash": "998605ace4a0b89293cf729cf216405f29c1ce5d44d6a40232982a4bdccda033",
        "ledger": 29113108,
        "created_at": "2020-04-10T17:03:18Z",
        "source_account": "GBUYDJH3AOPFFND3L54DUDWIHOMYKUONDV4RAHOHDBNN2D5N4BPPWDQ3",
        "source_account_sequence": "113942901088600162",
        "fee_account": "GBUYDJH3AOPFFND3L54DUDWIHOMYKUONDV4RAHOHDBNN2D5N4BPPWDQ3",
        "fee_charged": 3000000000,
        "max_fee": 2500000000,
        "operation_count": 1,
        "envelope_xdr": "AAAAAGmBpPsDnlK0e194Og7IO5mFUc0deRAdxxha3Q+t4F77AAAAZAGUzncAEEBiAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAADMzUyODFmNThmZjkxMGNiMTVhYWQ1NjM2ZGIyNzUzZTAAAAABAAAAAQAAAADhZHiqD/Q3uSTgjYEWGVRfCCHYvFmeqJU12G9SkzJYEQAAAAEAAAAAP29uBulc9ouSoH62BRypPhD6zcLWoS5sj7CHf5SJ15MAAAABTk9ETAAAAAB1jYLXrFzNBOWCoPnZSHI3PJAhHtc1TrCaiPuZwSf5pgAAAAAAAAABAAAAAAAAAALw9Tl2AAAAQOknEHs7ZaPNVlXMU0uOtT+0TVo9kW/jDuNxN40FdJDic0p23V4lxOfPGCgQwBgTehqCIEzCMQ4LkbfzkdgkFAut4F77AAAAQKtFmT73srS8RHeQgWWia8mb+TrLCr1CJbK+MAKGdUnb4s4JBOKUjHhqQLrs7GCkJ3wOpgTbtW8VpwNedCJhFQ0=",
        "result_xdr": "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=",
        "result_meta_xdr": "AAAAAQAAAAIAAAADAbw7FAAAAAAAAAAAaYGk+wOeUrR7X3g6Dsg7mYVRzR15EB3HGFrdD63gXvsAAAAA0gRBOAGUzncAEEBhAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAbw7FAAAAAAAAAAAaYGk+wOeUrR7X3g6Dsg7mYVRzR15EB3HGFrdD63gXvsAAAAA0gRBOAGUzncAEEBiAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMBvDsUAAAAAQAAAAA/b24G6Vz2i5KgfrYFHKk+EPrNwtahLmyPsId/lInXkwAAAAFOT0RMAAAAAHWNgtesXM0E5YKg+dlIcjc8kCEe1zVOsJqI+5nBJ/mmAAAAAMUrlJIASptjhEiAAAAAAAEAAAAAAAAAAAAAAAEBvDsUAAAAAQAAAAA/b24G6Vz2i5KgfrYFHKk+EPrNwtahLmyPsId/lInXkwAAAAFOT0RMAAAAAHWNgtesXM0E5YKg+dlIcjc8kCEe1zVOsJqI+5nBJ/mmAAAAAMUrlJMASptjhEiAAAAAAAEAAAAAAAAAAAAAAAMBvDsUAAAAAQAAAADhZHiqD/Q3uSTgjYEWGVRfCCHYvFmeqJU12G9SkzJYEQAAAAFOT0RMAAAAAHWNgtesXM0E5YKg+dlIcjc8kCEe1zVOsJqI+5nBJ/mmAAAAKdrbU7B//////////wAAAAEAAAAAAAAAAAAAAAEBvDsUAAAAAQAAAADhZHiqD/Q3uSTgjYEWGVRfCCHYvFmeqJU12G9SkzJYEQAAAAFOT0RMAAAAAHWNgtesXM0E5YKg+dlIcjc8kCEe1zVOsJqI+5nBJ/mmAAAAKdrbU69//////////wAAAAEAAAAAAAAAAA==",
        "fee_meta_xdr": "AAAAAgAAAAMBvDsTAAAAAAAAAABpgaT7A55StHtfeDoOyDuZhVHNHXkQHccYWt0PreBe+wAAAADSBEGcAZTOdwAQQGEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEBvDsUAAAAAAAAAABpgaT7A55StHtfeDoOyDuZhVHNHXkQHccYWt0PreBe+wAAAADSBEE4AZTOdwAQQGEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
        "memo_type": "hash",
        "signatures": [
          "6ScQeztlo81WVcxTS461P7RNWj2Rb+MO43E3jQV0kOJzSnbdXiXE588YKBDAGBN6GoIgTMIxDguRt/OR2CQUCw==",
          "q0WZPveytLxEd5CBZaJryZv5OssKvUIlsr4wAoZ1SdvizgkE4pSMeGpAuuzsYKQnfA6mBNu1bxWnA150ImEVDQ=="
        ],
        "valid_after": "1970-01-01T00:00:00Z"
      }`

	const feesAsStrings = `{
        "memo": "MzUyODFmNThmZjkxMGNiMTVhYWQ1NjM2ZGIyNzUzZTA=",
        "_links": {
          "self": {
            "href": "https://diamtestnet.diamcircle.io/transactions/998605ace4a0b89293cf729cf216405f29c1ce5d44d6a40232982a4bdccda033"
          },
          "account": {
            "href": "https://mainnet.diamcircle.io/GBUYDJH3AOPFFND3L54DUDWIHOMYKUONDV4RAHOHDBNN2D5N4BPPWDQ3"
          },
          "ledger": {
            "href": "https://aurora.diamcircle.org/ledgers/29113108"
          },
          "operations": {
            "href": "https://diamtestnet.diamcircle.io/transactions/998605ace4a0b89293cf729cf216405f29c1ce5d44d6a40232982a4bdccda033/operations{?cursor,limit,order}",
            "templated": true
          },
          "effects": {
            "href": "https://diamtestnet.diamcircle.io/transactions/998605ace4a0b89293cf729cf216405f29c1ce5d44d6a40232982a4bdccda033/effects{?cursor,limit,order}",
            "templated": true
          },
          "precedes": {
            "href": "https://diamtestnet.diamcircle.io/transactions?order=asc\u0026cursor=125039846745419776"
          },
          "succeeds": {
            "href": "https://diamtestnet.diamcircle.io/transactions?order=desc\u0026cursor=125039846745419776"
          },
          "transaction": {
            "href": "https://diamtestnet.diamcircle.io/transactions/998605ace4a0b89293cf729cf216405f29c1ce5d44d6a40232982a4bdccda033"
          }
        },
        "id": "998605ace4a0b89293cf729cf216405f29c1ce5d44d6a40232982a4bdccda033",
        "paging_token": "125039846745419776",
        "successful": true,
        "hash": "998605ace4a0b89293cf729cf216405f29c1ce5d44d6a40232982a4bdccda033",
        "ledger": 29113108,
        "created_at": "2020-04-10T17:03:18Z",
        "source_account": "GBUYDJH3AOPFFND3L54DUDWIHOMYKUONDV4RAHOHDBNN2D5N4BPPWDQ3",
        "source_account_sequence": "113942901088600162",
        "fee_account": "GBUYDJH3AOPFFND3L54DUDWIHOMYKUONDV4RAHOHDBNN2D5N4BPPWDQ3",
        "fee_charged": "3000000000",
        "max_fee": "2500000000",
        "operation_count": 1,
        "envelope_xdr": "AAAAAGmBpPsDnlK0e194Og7IO5mFUc0deRAdxxha3Q+t4F77AAAAZAGUzncAEEBiAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAADMzUyODFmNThmZjkxMGNiMTVhYWQ1NjM2ZGIyNzUzZTAAAAABAAAAAQAAAADhZHiqD/Q3uSTgjYEWGVRfCCHYvFmeqJU12G9SkzJYEQAAAAEAAAAAP29uBulc9ouSoH62BRypPhD6zcLWoS5sj7CHf5SJ15MAAAABTk9ETAAAAAB1jYLXrFzNBOWCoPnZSHI3PJAhHtc1TrCaiPuZwSf5pgAAAAAAAAABAAAAAAAAAALw9Tl2AAAAQOknEHs7ZaPNVlXMU0uOtT+0TVo9kW/jDuNxN40FdJDic0p23V4lxOfPGCgQwBgTehqCIEzCMQ4LkbfzkdgkFAut4F77AAAAQKtFmT73srS8RHeQgWWia8mb+TrLCr1CJbK+MAKGdUnb4s4JBOKUjHhqQLrs7GCkJ3wOpgTbtW8VpwNedCJhFQ0=",
        "result_xdr": "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=",
        "result_meta_xdr": "AAAAAQAAAAIAAAADAbw7FAAAAAAAAAAAaYGk+wOeUrR7X3g6Dsg7mYVRzR15EB3HGFrdD63gXvsAAAAA0gRBOAGUzncAEEBhAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAbw7FAAAAAAAAAAAaYGk+wOeUrR7X3g6Dsg7mYVRzR15EB3HGFrdD63gXvsAAAAA0gRBOAGUzncAEEBiAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMBvDsUAAAAAQAAAAA/b24G6Vz2i5KgfrYFHKk+EPrNwtahLmyPsId/lInXkwAAAAFOT0RMAAAAAHWNgtesXM0E5YKg+dlIcjc8kCEe1zVOsJqI+5nBJ/mmAAAAAMUrlJIASptjhEiAAAAAAAEAAAAAAAAAAAAAAAEBvDsUAAAAAQAAAAA/b24G6Vz2i5KgfrYFHKk+EPrNwtahLmyPsId/lInXkwAAAAFOT0RMAAAAAHWNgtesXM0E5YKg+dlIcjc8kCEe1zVOsJqI+5nBJ/mmAAAAAMUrlJMASptjhEiAAAAAAAEAAAAAAAAAAAAAAAMBvDsUAAAAAQAAAADhZHiqD/Q3uSTgjYEWGVRfCCHYvFmeqJU12G9SkzJYEQAAAAFOT0RMAAAAAHWNgtesXM0E5YKg+dlIcjc8kCEe1zVOsJqI+5nBJ/mmAAAAKdrbU7B//////////wAAAAEAAAAAAAAAAAAAAAEBvDsUAAAAAQAAAADhZHiqD/Q3uSTgjYEWGVRfCCHYvFmeqJU12G9SkzJYEQAAAAFOT0RMAAAAAHWNgtesXM0E5YKg+dlIcjc8kCEe1zVOsJqI+5nBJ/mmAAAAKdrbU69//////////wAAAAEAAAAAAAAAAA==",
        "fee_meta_xdr": "AAAAAgAAAAMBvDsTAAAAAAAAAABpgaT7A55StHtfeDoOyDuZhVHNHXkQHccYWt0PreBe+wAAAADSBEGcAZTOdwAQQGEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEBvDsUAAAAAAAAAABpgaT7A55StHtfeDoOyDuZhVHNHXkQHccYWt0PreBe+wAAAADSBEE4AZTOdwAQQGEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
        "memo_type": "hash",
        "signatures": [
          "6ScQeztlo81WVcxTS461P7RNWj2Rb+MO43E3jQV0kOJzSnbdXiXE588YKBDAGBN6GoIgTMIxDguRt/OR2CQUCw==",
          "q0WZPveytLxEd5CBZaJryZv5OssKvUIlsr4wAoZ1SdvizgkE4pSMeGpAuuzsYKQnfA6mBNu1bxWnA150ImEVDQ=="
        ],
        "valid_after": "1970-01-01T00:00:00Z"
      }`

	var parsedFeesAsInts, parsedFeesAsStrings Transaction
	assert.NoError(t, json.Unmarshal([]byte(feesAsInt64s), &parsedFeesAsInts))
	assert.NoError(t, json.Unmarshal([]byte(feesAsStrings), &parsedFeesAsStrings))
	assert.Equal(t, parsedFeesAsInts, parsedFeesAsStrings)
	assert.Equal(t, int64(2500000000), parsedFeesAsInts.MaxFee)
	assert.Equal(t, int64(3000000000), parsedFeesAsInts.FeeCharged)
}

func TestTradeAggregation_PagingToken(t *testing.T) {
	ta := TradeAggregation{Timestamp: 64}
	assert.Equal(t, "64", ta.PagingToken())
}
