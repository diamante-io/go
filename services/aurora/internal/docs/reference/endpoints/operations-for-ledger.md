---
title: Operations for Ledger
clientData:
  laboratoryUrl: https://www.diamcircle.org/laboratory/#explorer?resource=operations&endpoint=for_ledger
replacement: https://developers.diamcircle.org/api/resources/ledgers/operations/
---

This endpoint returns successful [operations](../resources/operation.md) that occurred in a given [ledger](../resources/ledger.md).

## Request

```
GET /ledgers/{sequence}/operations{?cursor,limit,order,include_failed}
```

### Arguments

| name | notes | description | example |
| ---- | ----- | ----------- | ------- |
| `sequence` | required, number | Ledger Sequence | `681637` |
| `?cursor` | optional, default _null_ | A paging token, specifying where to start returning records from. | `12884905984` |
| `?order` | optional, string, default `asc` | The order in which to return rows, "asc" or "desc". | `asc` |
| `?limit` | optional, number, default `10` | Maximum number of records to return. | `200` |
| `?include_failed` | optional, bool, default: `false` | Set to `true` to include operations of failed transactions in results. | `true` |
| `?join` | optional, string, default: _null_ | Set to `transactions` to include the transactions which created each of the operations in the response. | `transactions` |

### curl Example Request

```sh
curl "https://aurora-testnet.diamcircle.org/ledgers/681637/operations?limit=1"
```

### JavaScript Example Request

```javascript
var DiamcircleSdk = require('diamcircle-sdk');
var server = new DiamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

server.operations()
  .forLedger("681637")
  .call()
  .then(function (operationsResult) {
    console.log(operationsResult.records);
  })
  .catch(function (err) {
    console.log(err)
  })
```

## Response

This endpoint responds with a list of operations in a given ledger.  See [operation resource](../resources/operation.md) for reference.

### Example Response

```json
{
  "_links": {
    "self": {
      "href": "https://aurora-testnet.diamcircle.org/ledgers/681637/operations?cursor=&limit=10&order=asc"
    },
    "next": {
      "href": "https://aurora-testnet.diamcircle.org/ledgers/681637/operations?cursor=2927608622751745&limit=10&order=asc"
    },
    "prev": {
      "href": "https://aurora-testnet.diamcircle.org/ledgers/681637/operations?cursor=2927608622747649&limit=10&order=desc"
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
      },
      {
        "_links": {
          "self": {
            "href": "https://aurora-testnet.diamcircle.org/operations/2927608622751745"
          },
          "transaction": {
            "href": "https://aurora-testnet.diamcircle.org/transactions/fdabcee816bd439dd1d20bcb0abab5aa939c15cca5fccc1db060ba6096a5e0ed"
          },
          "effects": {
            "href": "https://aurora-testnet.diamcircle.org/operations/2927608622751745/effects"
          },
          "succeeds": {
            "href": "https://aurora-testnet.diamcircle.org/effects?order=desc&cursor=2927608622751745"
          },
          "precedes": {
            "href": "https://aurora-testnet.diamcircle.org/effects?order=asc&cursor=2927608622751745"
          }
        },
        "id": "2927608622751745",
        "paging_token": "2927608622751745",
        "transaction_successful": true,
        "source_account": "GAIH3ULLFQ4DGSECF2AR555KZ4KNDGEKN4AFI4SU2M7B43MGK3QJZNSR",
        "type": "create_account",
        "type_i": 0,
        "created_at": "2019-04-08T21:59:27Z",
        "transaction_hash": "fdabcee816bd439dd1d20bcb0abab5aa939c15cca5fccc1db060ba6096a5e0ed",
        "starting_balance": "10000.0000000",
        "funder": "GAIH3ULLFQ4DGSECF2AR555KZ4KNDGEKN4AFI4SU2M7B43MGK3QJZNSR",
        "account": "GCD5UL3DHC5TQRQVJKFTM66CLFTHGULOQ2HEAXNSA2JWUGBCT36BP55F"
      }
    ]
  }
}
```

## Possible Errors

- The [standard errors](../errors.md#Standard-Errors).
- [not_found](../errors/not-found.md): A `not_found` error will be returned if there is no ledger whose ID matches the `id` argument.
