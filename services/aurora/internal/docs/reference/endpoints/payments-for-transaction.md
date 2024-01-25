---
title: Payments for Transaction
clientData:
  laboratoryUrl: https://www.diamcircle.org/laboratory/#explorer?resource=payments&endpoint=for_transaction
---

This endpoint represents all payment-related [operations](../resources/operation.md) that are part
of a given [transaction](../resources/transaction.md).

The operations that can be returned in by this endpoint are:
- `create_account`
- `payment`
- `path_payment`
- `account_merge`

### Warning - failed transactions

"Payments for Transaction" endpoint returns list of payments of successful or failed transactions
(that are also included in Diamcircle ledger). Always check the payment status in this endpoint using
`transaction_successful` field!

## Request

```
GET /transactions/{hash}/payments{?cursor,limit,order}
```

### Arguments

| name | notes | description | example |
| ---- | ----- | ----------- | ------- |
| `hash` | required, string | A transaction hash, hex-encoded, lowercase. | `f65278b36875d170e865853838da400515f59ca23836f072e8d62cac18b803e5` |
| `?cursor` | optional, default _null_ | A paging token, specifying where to start returning records from. | `12884905984` |
| `?order` | optional, string, default `asc` | The order in which to return rows, "asc" or "desc". | `asc` |
| `?limit` | optional, number, default `10` | Maximum number of records to return. | `200` |
| `?join` | optional, string, default: _null_ | Set to `transactions` to include the transactions which created each of the payments in the response. | `transactions` |

### curl Example Request

```sh
curl "https://aurora-testnet.diamcircle.org/transactions/f65278b36875d170e865853838da400515f59ca23836f072e8d62cac18b803e5/payments"
```

### JavaScript Example Request

```javascript
var DiamcircleSdk = require('diamcircle-sdk');
var server = new DiamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

server.payments()
  .forTransaction("f65278b36875d170e865853838da400515f59ca23836f072e8d62cac18b803e5")
  .call()
  .then(function (paymentResult) {
    console.log(paymentResult.records);
  })
  .catch(function (err) {
    console.log(err);
  })
```

## Response

This endpoint responds with a list of payments operations that are part of a given transaction. See
[operation resource](../resources/operation.md) for more information about operations (and payment
operations).

### Example Response

```json
{
  "_links": {
    "self": {
      "href": "https://aurora-testnet.diamcircle.org/transactions/f65278b36875d170e865853838da400515f59ca23836f072e8d62cac18b803e5/payments?cursor=&limit=10&order=asc"
    },
    "next": {
      "href": "https://aurora-testnet.diamcircle.org/transactions/f65278b36875d170e865853838da400515f59ca23836f072e8d62cac18b803e5/payments?cursor=2993420406628353&limit=10&order=asc"
    },
    "prev": {
      "href": "https://aurora-testnet.diamcircle.org/transactions/f65278b36875d170e865853838da400515f59ca23836f072e8d62cac18b803e5/payments?cursor=2993420406628353&limit=10&order=desc"
    }
  },
  "_embedded": {
    "records": [
      {
        "_links": {
          "self": {
            "href": "https://aurora-testnet.diamcircle.org/operations/2993420406628353"
          },
          "transaction": {
            "href": "https://aurora-testnet.diamcircle.org/transactions/f65278b36875d170e865853838da400515f59ca23836f072e8d62cac18b803e5"
          },
          "effects": {
            "href": "https://aurora-testnet.diamcircle.org/operations/2993420406628353/effects"
          },
          "succeeds": {
            "href": "https://aurora-testnet.diamcircle.org/effects?order=desc&cursor=2993420406628353"
          },
          "precedes": {
            "href": "https://aurora-testnet.diamcircle.org/effects?order=asc&cursor=2993420406628353"
          }
        },
        "id": "2993420406628353",
        "paging_token": "2993420406628353",
        "transaction_successful": true,
        "source_account": "GAYB4GWPX2HUWR5QE7YX77QY6TSNFZIJZTYX2TDRW6YX6332BGD5SEAK",
        "type": "payment",
        "type_i": 1,
        "created_at": "2019-04-09T20:00:54Z",
        "transaction_hash": "f65278b36875d170e865853838da400515f59ca23836f072e8d62cac18b803e5",
        "asset_type": "native",
        "from": "GAYB4GWPX2HUWR5QE7YX77QY6TSNFZIJZTYX2TDRW6YX6332BGD5SEAK",
        "to": "GDGEQS64ISS6Y2KDM5V67B6LXALJX4E7VE4MIA54NANSUX5MKGKBZM5G",
        "amount": "293.0000000"
      }
    ]
  }
}
```

## Possible Errors

- The [standard errors](../errors.md#Standard-Errors).
- [not_found](../errors/not-found.md): A `not_found` error will be returned if there is no
  transaction whose ID matches the `hash` argument.
