---
title: All Ledgers
clientData:
  laboratoryUrl: https://www.diamcircle.org/laboratory/#explorer?resource=ledgers&endpoint=all
replacement: https://developers.diamcircle.org/api/resources/ledgers/
---

This endpoint represents all [ledgers](../resources/ledger.md).
This endpoint can also be used in [streaming](../streaming.md) mode so it is possible to use it to get notifications as ledgers are closed by the diamcircle network.
If called in streaming mode aurora will start at the earliest known ledger unless a `cursor` is set. In that case it will start from the `cursor`. You can also set `cursor` value to `now` to only stream ledgers created since your request time.

## Request

```
GET /ledgers{?cursor,limit,order}
```

### Arguments

| name | notes | description | example |
| ---- | ----- | ----------- | ------- |
| `?cursor` | optional, any, default _null_ | A paging token, specifying where to start returning records from. When streaming this can be set to `now` to stream object created since your request time. | `12884905984` |
| `?order`  | optional, string, default `asc` | The order in which to return rows, "asc" or "desc". | `asc` |
| `?limit`  | optional, number, default: `10` | Maximum number of records to return. | `200` |

### curl Example Request

```sh
# Retrieve the 200 latest ledgers, ordered chronologically
curl "https://aurora-testnet.diamcircle.org/ledgers?limit=200&order=desc"
```

### JavaScript Example Request

```javascript
server.ledgers()
  .call()
  .then(function (ledgerResult) {
    // page 1
    console.log(ledgerResult.records)
    return ledgerResult.next()
  })
  .then(function (ledgerResult) {
    // page 2
    console.log(ledgerResult.records)
  })
  .catch(function(err) {
    console.log(err)
  })
```


### JavaScript Streaming Example

```javascript
var diamcircleSdk = require('diamcircle-sdk')
var server = new diamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

var ledgerHandler = function (ledgerResponse) {
  console.log(ledgerResponse);
};

var es = server.ledgers()
  .cursor('now')
  .stream({
    onmessage: ledgerHandler
})
```

## Response

This endpoint responds with a list of ledgers.  See [ledger resource](../resources/ledger.md) for reference.

### Example Response

```json
{
  "_embedded": {
    "records": [
      {
        "_links": {
          "effects": {
            "href": "/ledgers/1/effects/{?cursor,limit,order}",
            "templated": true
          },
          "operations": {
            "href": "/ledgers/1/operations/{?cursor,limit,order}",
            "templated": true
          },
          "self": {
            "href": "/ledgers/1"
          },
          "transactions": {
            "href": "/ledgers/1/transactions/{?cursor,limit,order}",
            "templated": true
          }
        },
        "id": "e8e10918f9c000c73119abe54cf089f59f9015cc93c49ccf00b5e8b9afb6e6b1",
        "paging_token": "4294967296",
        "hash": "e8e10918f9c000c73119abe54cf089f59f9015cc93c49ccf00b5e8b9afb6e6b1",
        "sequence": 1,
        "transaction_count": 0,
        "successful_transaction_count": 0,
        "failed_transaction_count": 0,
        "operation_count": 0,
        "tx_set_operation_count": 0,
        "closed_at": "1970-01-01T00:00:00Z",
        "total_coins": "100000000000.0000000",
        "fee_pool": "0.0000000",
        "base_fee_in_stroops": 100,
        "base_reserve_in_stroops": 100000000,
        "max_tx_set_size": 50
      },
      {
        "_links": {
          "effects": {
            "href": "/ledgers/2/effects/{?cursor,limit,order}",
            "templated": true
          },
          "operations": {
            "href": "/ledgers/2/operations/{?cursor,limit,order}",
            "templated": true
          },
          "self": {
            "href": "/ledgers/2"
          },
          "transactions": {
            "href": "/ledgers/2/transactions/{?cursor,limit,order}",
            "templated": true
          }
        },
        "id": "e12e5809ab8c59d8256e691cb48a024dd43960bc15902d9661cd627962b2bc71",
        "paging_token": "8589934592",
        "hash": "e12e5809ab8c59d8256e691cb48a024dd43960bc15902d9661cd627962b2bc71",
        "prev_hash": "e8e10918f9c000c73119abe54cf089f59f9015cc93c49ccf00b5e8b9afb6e6b1",
        "sequence": 2,
        "transaction_count": 0,
        "successful_transaction_count": 0,
        "failed_transaction_count": 0,
        "operation_count": 0,
        "closed_at": "2015-07-16T23:49:00Z",
        "total_coins": "100000000000.0000000",
        "fee_pool": "0.0000000",
        "base_fee_in_stroops": 100,
        "base_reserve_in_stroops": 100000000,
        "max_tx_set_size": 100
      }
    ]
  },
  "_links": {
    "next": {
      "href": "/ledgers?order=asc&limit=2&cursor=8589934592"
    },
    "prev": {
      "href": "/ledgers?order=desc&limit=2&cursor=4294967296"
    },
    "self": {
      "href": "/ledgers?order=asc&limit=2&cursor="
    }
  }
}
```

### Example Streaming Event

```json
{
  "_links": {
    "effects": {
      "href": "/ledgers/69859/effects/{?cursor,limit,order}",
      "templated": true
    },
    "operations": {
      "href": "/ledgers/69859/operations/{?cursor,limit,order}",
      "templated": true
    },
    "self": {
      "href": "/ledgers/69859"
    },
    "transactions": {
      "href": "/ledgers/69859/transactions/{?cursor,limit,order}",
      "templated": true
    }
  },
  "id": "4db1e4f145e9ee75162040d26284795e0697e2e84084624e7c6c723ebbf80118",
  "paging_token": "300042120331264",
  "hash": "4db1e4f145e9ee75162040d26284795e0697e2e84084624e7c6c723ebbf80118",
  "prev_hash": "4b0b8bace3b2438b2404776ce57643966855487ba6384724a3c664c7aa4cd9e4",
  "sequence": 69859,
  "transaction_count": 0,
  "successful_transaction_count": 0,
  "failed_transaction_count": 0,
  "operation_count": 0,
  "closed_at": "2015-07-20T15:51:52Z",
  "total_coins": "100000000000.0000000",
  "fee_pool": "0.0025600",
  "base_fee_in_stroops": 100,
  "base_reserve_in_stroops": "100000000",
  "max_tx_set_size": 50
}
```

## Errors

- The [standard errors](../errors.md#standard-errors).
