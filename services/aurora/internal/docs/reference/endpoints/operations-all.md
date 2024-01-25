---
title: All Operations
clientData:
  laboratoryUrl: https://www.diamcircle.org/laboratory/#explorer?resource=operations&endpoint=all
replacement: https://developers.diamcircle.org/api/resources/operations/
---

This endpoint represents [operations](../resources/operation.md) that are part of successfully validated [transactions](../resources/transaction.md).
Please note that this endpoint returns operations that are part of failed transactions if `include_failed` parameter is `true`
and Aurora is ingesting failed transactions.
This endpoint can also be used in [streaming](../streaming.md) mode so it is possible to use it to listen as operations are processed in the Diamcircle network.
If called in streaming mode Aurora will start at the earliest known operation unless a `cursor` is set. In that case it will start from the `cursor`. You can also set `cursor` value to `now` to only stream operations created since your request time.

## Request

```
GET /operations{?cursor,limit,order,include_failed}
```

### Arguments

| name | notes | description | example |
| ---- | ----- | ----------- | ------- |
| `?cursor` | optional, any, default _null_ | A paging token, specifying where to start returning records from. When streaming this can be set to `now` to stream object created since your request time. | `12884905984` |
| `?order`  | optional, string, default `asc` | The order in which to return rows, "asc" or "desc". | `asc` |
| `?limit`  | optional, number, default: `10` | Maximum number of records to return. | `200` |
| `?include_failed` | optional, bool, default: `false` | Set to `true` to include operations of failed transactions in results. | `true` |
| `?join` | optional, string, default: _null_ | Set to `transactions` to include the transactions which created each of the operations in the response. | `transactions` |

### curl Example Request

```sh
curl "https://aurora-testnet.diamcircle.org/operations?limit=200&order=desc"
```

### JavaScript Example Request

```js
var DiamcircleSdk = require('diamcircle-sdk');
var server = new DiamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

server.operations()
  .call()
  .then(function (operationsResult) {
    //page 1
    console.log(operationsResult.records)
    return operationsResult.next()
  })
  .then(function (operationsResult) {
    //page 2
    console.log(operationsResult.records)
  })
  .catch(function (err) {
    console.log(err)
  })
```

### JavaScript Streaming Example

```javascript
var DiamcircleSdk = require('diamcircle-sdk')
var server = new DiamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

var operationHandler = function (operationResponse) {
  console.log(operationResponse);
};

var es = server.operations()
  .cursor('now')
  .stream({
    onmessage: operationHandler
  })
```

## Response

This endpoint responds with a list of operations. See [operation resource](../resources/operation.md) for reference.

### Example Response

```json
{
  "_embedded": {
    "records": [
      {
        "_links": {
          "effects": {
            "href": "/operations/77309415424/effects/{?cursor,limit,order}",
            "templated": true
          },
          "precedes": {
            "href": "/operations?cursor=77309415424&order=asc"
          },
          "self": {
            "href": "/operations/77309415424"
          },
          "succeeds": {
            "href": "/operations?cursor=77309415424&order=desc"
          },
          "transactions": {
            "href": "/transactions/77309415424"
          }
        },
        "account": "GBIA4FH6TV64KSPDAJCNUQSM7PFL4ILGUVJDPCLUOPJ7ONMKBBVUQHRO",
        "funder": "GCEZWKCA5VLDNRLN3RPRJMRZOX3Z6G5CHCGSNFHEYVXM3XOJMDS674JZ",
        "id": 77309415424,
        "paging_token": "77309415424",
        "starting_balance": "1000.0000000",
        "transaction_successful": true,
        "type_i": 0,
        "type": "create_account"
      },
      {
        "_links": {
          "effects": {
            "href": "/operations/463856472064/effects/{?cursor,limit,order}",
            "templated": true
          },
          "precedes": {
            "href": "/operations?cursor=463856472064&order=asc"
          },
          "self": {
            "href": "/operations/463856472064"
          },
          "succeeds": {
            "href": "/operations?cursor=463856472064&order=desc"
          },
          "transactions": {
            "href": "/transactions/463856472064"
          }
        },
        "account": "GC2ADYAIPKYQRGGUFYBV2ODJ54PY6VZUPKNCWWNX2C7FCJYKU4ZZNKVL",
        "funder": "GBIA4FH6TV64KSPDAJCNUQSM7PFL4ILGUVJDPCLUOPJ7ONMKBBVUQHRO",
        "id": 463856472064,
        "paging_token": "463856472064",
        "starting_balance": "1000.0000000",
        "transaction_successful": true,
        "type_i": 0,
        "type": "create_account"
      }
    ]
  },
  "_links": {
    "next": {
      "href": "/operations?order=asc&limit=2&cursor=463856472064"
    },
    "prev": {
      "href": "/operations?order=desc&limit=2&cursor=77309415424"
    },
    "self": {
      "href": "/operations?order=asc&limit=2&cursor="
    }
  }
}
```

### Example Streaming Event

```json
{
  "_links": {
    "effects": {
      "href": "/operations/77309415424/effects/{?cursor,limit,order}",
      "templated": true
    },
    "precedes": {
      "href": "/operations?cursor=77309415424&order=asc"
    },
    "self": {
      "href": "/operations/77309415424"
    },
    "succeeds": {
      "href": "/operations?cursor=77309415424&order=desc"
    },
    "transactions": {
      "href": "/transactions/77309415424"
    }
  },
  "account": "GBIA4FH6TV64KSPDAJCNUQSM7PFL4ILGUVJDPCLUOPJ7ONMKBBVUQHRO",
  "funder": "GCEZWKCA5VLDNRLN3RPRJMRZOX3Z6G5CHCGSNFHEYVXM3XOJMDS674JZ",
  "id": 77309415424,
  "paging_token": "77309415424",
  "starting_balance": "1000.0000000",
  "transaction_successful": true,
  "type_i": 0,
  "type": "create_account"
}
```

## Possible Errors

- The [standard errors](../errors.md#standard-errors).
