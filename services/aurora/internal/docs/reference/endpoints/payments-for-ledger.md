---
title: Payments for Ledger
clientData:
  laboratoryUrl: https://www.diamcircle.org/laboratory/#explorer?resource=payments&endpoint=for_ledger
replacement: https://developers.diamcircle.org/api/resources/ledgers/payments/
---

This endpoint represents all payment-related [operations](../resources/operation.md) that are part
of a valid [transactions](../resources/transaction.md) in a given [ledger](../resources/ledger.md).

The operations that can be returned in by this endpoint are:
- `create_account`
- `payment`
- `path_payment`
- `account_merge`

## Request

```
GET /ledgers/{id}/payments{?cursor,limit,order,include_failed}
```

### Arguments

|  name  |  notes  | description | example |
| ------ | ------- | ----------- | ------- |
| `id` | required, number | Ledger ID | `696960` |
| `?cursor` | optional, default _null_ | A paging token, specifying where to start returning records from. | `12884905984` |
| `?order`  | optional, string, default `asc` | The order in which to return rows, "asc" or "desc". | `asc` |
| `?limit`  | optional, number, default `10` | Maximum number of records to return. | `200` |
| `?include_failed` | optional, bool, default: `false` | Set to `true` to include payments of failed transactions in results. | `true` |
| `?join` | optional, string, default: _null_ | Set to `transactions` to include the transactions which created each of the payments in the response. | `transactions` |

### curl Example Request

```sh
curl "https://aurora-testnet.diamcircle.org/ledgers/696960/payments?limit=1"
```

### JavaScript Example Request

```javascript
var diamcircleSdk = require('diamcircle-sdk')
var server = new diamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

server.payments()
  .forLedger("696960")
  .call()
  .then(function (paymentResult) {
    console.log(paymentResult)
  })
  .catch(function (err) {
    console.log(err)
  })
```

## Response

This endpoint responds with a list of payment operations in a given ledger. See [operation
resource](../resources/operation.md) for more information about operations (and payment
operations).

### Example Response

```json
{
  "_links": {
    "self": {
      "href": "https://aurora-testnet.diamcircle.org/ledgers/696960/payments?cursor=&limit=1&order=asc"
    },
    "next": {
      "href": "https://aurora-testnet.diamcircle.org/ledgers/696960/payments?cursor=2993420406628353&limit=1&order=asc"
    },
    "prev": {
      "href": "https://aurora-testnet.diamcircle.org/ledgers/696960/payments?cursor=2993420406628353&limit=1&order=desc"
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
- [not_found](../errors/not-found.md): A `not_found` error will be returned if there is no ledger whose ID matches the `id` argument.
