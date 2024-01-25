---
title: Operations for Transaction
clientData:
  laboratoryUrl: https://www.diamcircle.org/laboratory/#explorer?resource=operations&endpoint=for_transaction
replacement: https://developers.diamcircle.org/api/resources/transactions/operations/
---

This endpoint represents successful [operations](../resources/operation.md) that are part of a given [transaction](../resources/transaction.md).

### Warning - failed transactions

The "Operations for Transaction" endpoint returns a list of operations in a successful or failed
transaction. Make sure to always check the operation status in this endpoint using
`transaction_successful` field!

## Request

```
GET /transactions/{hash}/operations{?cursor,limit,order}
```

## Arguments

| name | notes | description | example |
| ---- | ----- | ----------- | ------- |
| `hash` | required, string | A transaction hash, hex-encoded, lowercase. | `4a3365180521e16b478d9f0c9198b97a9434fc9cb07b34f83ecc32fc54d0ca8a` |
| `?cursor` | optional, default _null_ | A paging token, specifying where to start returning records from. | `12884905984` |
| `?order` | optional, string, default `asc` | The order in which to return rows, "asc" or "desc". | `asc` |
| `?limit` | optional, number, default `10` | Maximum number of records to return. | `200` |
| `?join` | optional, string, default: _null_ | Set to `transactions` to include the transactions which created each of the operations in the response. | `transactions` |

### curl Example Request

```sh
curl "https://aurora-testnet.diamcircle.org/transactions/4a3365180521e16b478d9f0c9198b97a9434fc9cb07b34f83ecc32fc54d0ca8a/operations?limit=1"
```

### JavaScript Example Request

```javascript
var diamcircleSdk = require('diamcircle-sdk');
var server = new diamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

server.operations()
  .forTransaction("4a3365180521e16b478d9f0c9198b97a9434fc9cb07b34f83ecc32fc54d0ca8a")
  .call()
  .then(function (operationsResult) {
    console.log(operationsResult.records);
  })
  .catch(function (err) {
    console.log(err)
  })
```

## Response

This endpoint responds with a list of operations that are part of a given transaction. See [operation resource](../resources/operation.md) for reference.

### Example Response

```json
{
  "_links": {
    "self": {
      "href": "https://aurora-testnet.diamcircle.org/transactions/4a3365180521e16b478d9f0c9198b97a9434fc9cb07b34f83ecc32fc54d0ca8a/operations?cursor=&limit=10&order=asc"
    },
    "next": {
      "href": "https://aurora-testnet.diamcircle.org/transactions/4a3365180521e16b478d9f0c9198b97a9434fc9cb07b34f83ecc32fc54d0ca8a/operations?cursor=2927608622747649&limit=10&order=asc"
    },
    "prev": {
      "href": "https://aurora-testnet.diamcircle.org/transactions/4a3365180521e16b478d9f0c9198b97a9434fc9cb07b34f83ecc32fc54d0ca8a/operations?cursor=2927608622747649&limit=10&order=desc"
    }
  },
  "_embedded": {
    "records": [
      {
        "_links": {
          "self": {
            "href": "https://aurora-testnet.diamcircle.org/operations/2927608622747649"
          },
          "transaction": {
            "href": "https://aurora-testnet.diamcircle.org/transactions/4a3365180521e16b478d9f0c9198b97a9434fc9cb07b34f83ecc32fc54d0ca8a"
          },
          "effects": {
            "href": "https://aurora-testnet.diamcircle.org/operations/2927608622747649/effects"
          },
          "succeeds": {
            "href": "https://aurora-testnet.diamcircle.org/effects?order=desc&cursor=2927608622747649"
          },
          "precedes": {
            "href": "https://aurora-testnet.diamcircle.org/effects?order=asc&cursor=2927608622747649"
          }
        },
        "id": "2927608622747649",
        "paging_token": "2927608622747649",
        "transaction_successful": true,
        "source_account": "GCGXZPH2QNKJP4GI2J77EFQQUMP3NYY4PCUZ4UPKHR2XYBKRUYKQ2DS6",
        "type": "payment",
        "type_i": 1,
        "created_at": "2019-04-08T21:59:27Z",
        "transaction_hash": "4a3365180521e16b478d9f0c9198b97a9434fc9cb07b34f83ecc32fc54d0ca8a",
        "asset_type": "native",
        "from": "GCGXZPH2QNKJP4GI2J77EFQQUMP3NYY4PCUZ4UPKHR2XYBKRUYKQ2DS6",
        "to": "GDGEQS64ISS6Y2KDM5V67B6LXALJX4E7VE4MIA54NANSUX5MKGKBZM5G",
        "amount": "404.0000000"
      }
    ]
  }
}
```

## Possible Errors

- The [standard errors](../errors.md#Standard-Errors).
- [not_found](../errors/not-found.md): A `not_found` error will be returned if there is no account whose ID matches the `hash` argument.
