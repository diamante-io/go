---
title: Effects for Ledger
clientData:
  laboratoryUrl: https://www.diamcircle.org/laboratory/#explorer?resource=effects&endpoint=for_ledger
---

Effects are the specific ways that the ledger was changed by any operation.

This endpoint represents all [effects](../resources/effect.md) that occurred in the given [ledger](../resources/ledger.md).

## Request

```
GET /ledgers/{sequence}/effects{?cursor,limit,order}
```

## Arguments

| name | notes | description | example |
| ---- | ----- | ----------- | ------- |
| `sequence` | required, number | Ledger Sequence Number | `680777` |
| `?cursor` | optional, default _null_ | A paging token, specifying where to start returning records from. | `12884905984` |
| `?order` | optional, string, default `asc` | The order in which to return rows, "asc" or "desc". | `asc` |
| `?limit` | optional, number, default `10` | Maximum number of records to return. | `200` |

### curl Example Request

```sh
curl "https://aurora-testnet.diamcircle.org/ledgers/680777/effects?limit=1"
```

### JavaScript Example Request

```javascript
var DiamcircleSdk = require('diamcircle-sdk');
var server = new DiamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

server.effects()
  .forLedger("680777")
  .call()
  .then(function (effectResults) {
    //page 1
    console.log(effectResults.records)
  })
  .catch(function (err) {
    console.log(err)
  })

```

## Response

This endpoint responds with a list of effects that occurred in the ledger. See [effect resource](../resources/effect.md) for reference.

### Example Response

```json
{
  "_links": {
    "self": {
      "href": "https://aurora-testnet.diamcircle.org/ledgers/680777/effects?cursor=&limit=10&order=asc"
    },
    "next": {
      "href": "https://aurora-testnet.diamcircle.org/ledgers/680777/effects?cursor=2923914950873089-3&limit=10&order=asc"
    },
    "prev": {
      "href": "https://aurora-testnet.diamcircle.org/ledgers/680777/effects?cursor=2923914950873089-1&limit=10&order=desc"
    }
  },
  "_embedded": {
    "records": [
      {
        "_links": {
          "operation": {
            "href": "https://aurora-testnet.diamcircle.org/operations/2923914950873089"
          },
          "succeeds": {
            "href": "https://aurora-testnet.diamcircle.org/effects?order=desc&cursor=2923914950873089-1"
          },
          "precedes": {
            "href": "https://aurora-testnet.diamcircle.org/effects?order=asc&cursor=2923914950873089-1"
          }
        },
        "id": "0002923914950873089-0000000001",
        "paging_token": "2923914950873089-1",
        "account": "GC4ALQ3GTT5BTHTOULHCJGAT4P3MUSPLU4OEE74BAVIJ6K443O6RVLRT",
        "type": "account_created",
        "type_i": 0,
        "created_at": "2019-04-08T20:47:22Z",
        "starting_balance": "10000.0000000"
      },
      {
        "_links": {
          "operation": {
            "href": "https://aurora-testnet.diamcircle.org/operations/2923914950873089"
          },
          "succeeds": {
            "href": "https://aurora-testnet.diamcircle.org/effects?order=desc&cursor=2923914950873089-2"
          },
          "precedes": {
            "href": "https://aurora-testnet.diamcircle.org/effects?order=asc&cursor=2923914950873089-2"
          }
        },
        "id": "0002923914950873089-0000000002",
        "paging_token": "2923914950873089-2",
        "account": "GAIH3ULLFQ4DGSECF2AR555KZ4KNDGEKN4AFI4SU2M7B43MGK3QJZNSR",
        "type": "account_debited",
        "type_i": 3,
        "created_at": "2019-04-08T20:47:22Z",
        "asset_type": "native",
        "amount": "10000.0000000"
      },
      {
        "_links": {
          "operation": {
            "href": "https://aurora-testnet.diamcircle.org/operations/2923914950873089"
          },
          "succeeds": {
            "href": "https://aurora-testnet.diamcircle.org/effects?order=desc&cursor=2923914950873089-3"
          },
          "precedes": {
            "href": "https://aurora-testnet.diamcircle.org/effects?order=asc&cursor=2923914950873089-3"
          }
        },
        "id": "0002923914950873089-0000000003",
        "paging_token": "2923914950873089-3",
        "account": "GC4ALQ3GTT5BTHTOULHCJGAT4P3MUSPLU4OEE74BAVIJ6K443O6RVLRT",
        "type": "signer_created",
        "type_i": 10,
        "created_at": "2019-04-08T20:47:22Z",
        "weight": 1,
        "public_key": "GC4ALQ3GTT5BTHTOULHCJGAT4P3MUSPLU4OEE74BAVIJ6K443O6RVLRT",
        "key": ""
      }
    ]
  }
}
```

## Possible Errors

- The [standard errors](../errors.md#Standard-Errors).
- [not_found](../errors/not-found.md): A `not_found` error will be returned if there are no effects for a given ledger.
