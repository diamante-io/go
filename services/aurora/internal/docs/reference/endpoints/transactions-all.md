---
title: All Transactions
clientData:
  laboratoryUrl: https://www.diamcircle.org/laboratory/#explorer?resource=transactions&endpoint=all
replacement: https://developers.diamcircle.org/api/resources/transactions/single/
---

This endpoint represents all successful [transactions](../resources/transaction.md).
Please note that this endpoint returns failed transactions that are included in the ledger if
`include_failed` parameter is `true` and Aurora is ingesting failed transactions.
This endpoint can also be used in [streaming](../streaming.md) mode. This makes it possible to use
it to listen for new transactions as they get made in the Diamcircle network. If called in streaming
mode Aurora will start at the earliest known transaction unless a `cursor` is set. In that case it
will start from the `cursor`. You can also set `cursor` value to `now` to only stream transaction
created since your request time.

## Request

```
GET /transactions{?cursor,limit,order,include_failed}
```

### Arguments

| name | notes | description | example |
| ---- | ----- | ----------- | ------- |
| `?cursor` | optional, any, default _null_ | A paging token, specifying where to start returning records from. When streaming this can be set to `now` to stream object created since your request time. | `12884905984` |
| `?order`  | optional, string, default `asc` | The order in which to return rows, "asc" or "desc". | `asc` |
| `?limit`  | optional, number, default: `10` | Maximum number of records to return. | `200` |
| `?include_failed` | optional, bool, default: `false` | Set to `true` to include failed transactions in results. | `true` |

### curl Example Request

```sh
# Retrieve the 200 latest transactions, ordered chronologically:
curl "https://aurora-testnet.diamcircle.org/transactions?limit=200&order=desc"
```

### JavaScript Example Request

```javascript
var DiamcircleSdk = require('diamcircle-sdk');
var server = new DiamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

server.transactions()
  .call()
  .then(function (transactionResult) {
    //page 1
    console.log(transactionResult.records);
    return transactionResult.next();
  })
  .then(function (transactionResult) {
    console.log(transactionResult.records);
  })
  .catch(function (err) {
    console.log(err)
  })
```

### JavaScript Streaming Example

```javascript
var DiamcircleSdk = require('diamcircle-sdk')
var server = new DiamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

var txHandler = function (txResponse) {
  console.log(txResponse);
};

var es = server.transactions()
  .cursor('now')
  .stream({
      onmessage: txHandler
  })
```

## Response

If called normally this endpoint responds with a [page](../resources/page.md) of transactions.
If called in streaming mode the transaction resources are returned individually.
See [transaction resource](../resources/transaction.md) for reference.

### Example Response

```json
{
  "_embedded": {
    "records": [
      {
        "_links": {
          "account": {
            "href": "/accounts/GBS43BF24ENNS3KPACUZVKK2VYPOZVBQO2CISGZ777RYGOPYC2FT6S3K"
          },
          "effects": {
            "href": "/transactions/fa78cb43d72171fdb2c6376be12d57daa787b1fa1a9fdd0e9453e1f41ee5f15a/effects{?cursor,limit,order}",
            "templated": true
          },
          "ledger": {
            "href": "/ledgers/146970"
          },
          "operations": {
            "href": "/transactions/fa78cb43d72171fdb2c6376be12d57daa787b1fa1a9fdd0e9453e1f41ee5f15a/operations{?cursor,limit,order}",
            "templated": true
          },
          "precedes": {
            "href": "/transactions?cursor=631231343497216\u0026order=asc"
          },
          "self": {
            "href": "/transactions/fa78cb43d72171fdb2c6376be12d57daa787b1fa1a9fdd0e9453e1f41ee5f15a"
          },
          "succeeds": {
            "href": "/transactions?cursor=631231343497216\u0026order=desc"
          }
        },
        "id": "fa78cb43d72171fdb2c6376be12d57daa787b1fa1a9fdd0e9453e1f41ee5f15a",
        "paging_token": "631231343497216",
        "successful": true,
        "hash": "fa78cb43d72171fdb2c6376be12d57daa787b1fa1a9fdd0e9453e1f41ee5f15a",
        "ledger": 146970,
        "created_at": "2015-09-24T10:07:09Z",
        "account": "GBS43BF24ENNS3KPACUZVKK2VYPOZVBQO2CISGZ777RYGOPYC2FT6S3K",
        "account_sequence": 279172874343,
        "max_fee": 100,
        "fee_charged": 100,
        "operation_count": 1,
        "envelope_xdr": "AAAAAGXNhLrhGtltTwCpmqlarh7s1DB2hIkbP//jgzn4Fos/AAAACgAAAEEAAABnAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAA2ddmTOFAgr21Crs2RXRGLhiAKxicZb/IERyEZL/Y2kUAAAAXSHboAAAAAAAAAAAB+BaLPwAAAECDEEZmzbgBr5fc3mfJsCjWPDtL6H8/vf16me121CC09ONyWJZnw0PUvp4qusmRwC6ZKfLDdk8F3Rq41s+yOgQD",
        "result_xdr": "AAAAAAAAAAoAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=",
        "result_meta_xdr": "AAAAAAAAAAEAAAACAAAAAAACPhoAAAAAAAAAANnXZkzhQIK9tQq7NkV0Ri4YgCsYnGW/yBEchGS/2NpFAAAAF0h26AAAAj4aAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAQACPhoAAAAAAAAAAGXNhLrhGtltTwCpmqlarh7s1DB2hIkbP//jgzn4Fos/AABT8kS2c/oAAABBAAAAZwAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAA"
      },
      {
        "_links": {
          "account": {
            "href": "/accounts/GBS43BF24ENNS3KPACUZVKK2VYPOZVBQO2CISGZ777RYGOPYC2FT6S3K"
          },
          "effects": {
            "href": "/transactions/90ad6cfc9b0911bdbf202cace78ae7ecf50989c424288670dadb69bf8237c1b3/effects{?cursor,limit,order}",
            "templated": true
          },
          "ledger": {
            "href": "/ledgers/144798"
          },
          "operations": {
            "href": "/transactions/90ad6cfc9b0911bdbf202cace78ae7ecf50989c424288670dadb69bf8237c1b3/operations{?cursor,limit,order}",
            "templated": true
          },
          "precedes": {
            "href": "/transactions?cursor=621902674530304\u0026order=asc"
          },
          "self": {
            "href": "/transactions/90ad6cfc9b0911bdbf202cace78ae7ecf50989c424288670dadb69bf8237c1b3"
          },
          "succeeds": {
            "href": "/transactions?cursor=621902674530304\u0026order=desc"
          }
        },
        "id": "90ad6cfc9b0911bdbf202cace78ae7ecf50989c424288670dadb69bf8237c1b3",
        "paging_token": "621902674530304",
        "successful": false,
        "hash": "90ad6cfc9b0911bdbf202cace78ae7ecf50989c424288670dadb69bf8237c1b3",
        "ledger": 144798,
        "created_at": "2015-09-24T07:49:38Z",
        "account": "GBS43BF24ENNS3KPACUZVKK2VYPOZVBQO2CISGZ777RYGOPYC2FT6S3K",
        "account_sequence": 279172874342,
        "max_fee": 100,
        "fee_charged": 100,
        "operation_count": 1,
        "envelope_xdr": "AAAAAGXNhLrhGtltTwCpmqlarh7s1DB2hIkbP//jgzn4Fos/AAAACgAAAEEAAABmAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAMPT7P7buwqnMueFS4NV10vE2q3C/mcAy4jx03/RdSGsAAAAXSHboAAAAAAAAAAAB+BaLPwAAAEBPWWMNSWyPBbQlhRheXyvAFDVx1rnf68fdDOUHPdDIkHdUczBpzvCjpdgwhQ2NYOX5ga1ZgOIWLy789YNnuIcL",
        "result_xdr": "AAAAAAAAAAoAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=",
        "result_meta_xdr": "AAAAAAAAAAEAAAACAAAAAAACNZ4AAAAAAAAAADD0+z+27sKpzLnhUuDVddLxNqtwv5nAMuI8dN/0XUhrAAAAF0h26AAAAjWeAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAQACNZ4AAAAAAAAAAGXNhLrhGtltTwCpmqlarh7s1DB2hIkbP//jgzn4Fos/AABUCY0tXAQAAABBAAAAZgAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAA"
      }
    ]
  },
  "_links": {
    "next": {
      "href": "/transactions?order=desc\u0026limit=2\u0026cursor=621902674530304"
    },
    "prev": {
      "href": "/transactions?order=asc\u0026limit=2\u0026cursor=631231343497216"
    },
    "self": {
      "href": "/transactions?order=desc\u0026limit=2\u0026cursor="
    }
  }
}
```

## Possible Errors

- The [standard errors](../errors.md#standard-errors).
