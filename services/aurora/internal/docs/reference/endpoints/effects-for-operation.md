---
title: Effects for Operation
clientData:
  laboratoryUrl: https://www.diamcircle.org/laboratory/#explorer?resource=effects&endpoint=for_operation
replacement: https://developers.diamcircle.org/api/resources/operations/effects/
---

This endpoint represents all [effects](../resources/effect.md) that occurred as a result of a given [operation](../resources/operation.md).

## Request

```
GET /operations/{id}/effects{?cursor,limit,order}
```

### Arguments

| name | notes | description | example |
| ---- | ----- | ----------- | ------- |
| `id` | required, number | An operation ID | `1919197546291201` |
| `?cursor` | optional, default _null_ | A paging token, specifying where to start returning records from. | `12884905984` |
| `?order` | optional, string, default `asc` | The order in which to return rows, "asc" or "desc". | `asc` |
| `?limit` | optional, number, default `10` | Maximum number of records to return. | `200` |

### curl Example Request

```sh
curl "https://aurora-testnet.diamcircle.org/operations/1919197546291201/effects"
```

### JavaScript Example Request

```javascript
var diamcircleSdk = require('diamcircle-sdk');
var server = new diamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

server.effects()
  .forOperation("1919197546291201")
  .call()
  .then(function (effectResults) {
    // page 1
    console.log(effectResults.records)
  })
  .catch(function (err) {
    console.log(err)
  })

```

## Response

This endpoint responds with a list of effects on the ledger as a result of a given operation. See [effect resource](../resources/effect.md) for reference.

### Example Response

```json
{
  "_links": {
    "self": {
      "href": "https://aurora-testnet.diamcircle.org/operations/1919197546291201/effects?cursor=&limit=10&order=asc"
    },
    "next": {
      "href": "https://aurora-testnet.diamcircle.org/operations/1919197546291201/effects?cursor=1919197546291201-3&limit=10&order=asc"
    },
    "prev": {
      "href": "https://aurora-testnet.diamcircle.org/operations/1919197546291201/effects?cursor=1919197546291201-1&limit=10&order=desc"
    }
  },
  "_embedded": {
    "records": [
      {
        "_links": {
          "operation": {
            "href": "https://aurora-testnet.diamcircle.org/operations/1919197546291201"
          },
          "succeeds": {
            "href": "https://aurora-testnet.diamcircle.org/effects?order=desc&cursor=1919197546291201-1"
          },
          "precedes": {
            "href": "https://aurora-testnet.diamcircle.org/effects?order=asc&cursor=1919197546291201-1"
          }
        },
        "id": "0001919197546291201-0000000001",
        "paging_token": "1919197546291201-1",
        "account": "GBYUUJHG6F4EPJGNLERINATVQLNDOFRUD7SGJZ26YZLG5PAYLG7XUSGF",
        "type": "account_created",
        "type_i": 0,
        "created_at": "2019-03-25T22:43:38Z",
        "starting_balance": "10000.0000000"
      },
      {
        "_links": {
          "operation": {
            "href": "https://aurora-testnet.diamcircle.org/operations/1919197546291201"
          },
          "succeeds": {
            "href": "https://aurora-testnet.diamcircle.org/effects?order=desc&cursor=1919197546291201-2"
          },
          "precedes": {
            "href": "https://aurora-testnet.diamcircle.org/effects?order=asc&cursor=1919197546291201-2"
          }
        },
        "id": "0001919197546291201-0000000002",
        "paging_token": "1919197546291201-2",
        "account": "GAIH3ULLFQ4DGSECF2AR555KZ4KNDGEKN4AFI4SU2M7B43MGK3QJZNSR",
        "type": "account_debited",
        "type_i": 3,
        "created_at": "2019-03-25T22:43:38Z",
        "asset_type": "native",
        "amount": "10000.0000000"
      },
      {
        "_links": {
          "operation": {
            "href": "https://aurora-testnet.diamcircle.org/operations/1919197546291201"
          },
          "succeeds": {
            "href": "https://aurora-testnet.diamcircle.org/effects?order=desc&cursor=1919197546291201-3"
          },
          "precedes": {
            "href": "https://aurora-testnet.diamcircle.org/effects?order=asc&cursor=1919197546291201-3"
          }
        },
        "id": "0001919197546291201-0000000003",
        "paging_token": "1919197546291201-3",
        "account": "GBYUUJHG6F4EPJGNLERINATVQLNDOFRUD7SGJZ26YZLG5PAYLG7XUSGF",
        "type": "signer_created",
        "type_i": 10,
        "created_at": "2019-03-25T22:43:38Z",
        "weight": 1,
        "public_key": "GBYUUJHG6F4EPJGNLERINATVQLNDOFRUD7SGJZ26YZLG5PAYLG7XUSGF",
        "key": ""
      }
    ]
  }
}
```

## Possible Errors

- The [standard errors](../errors.md#Standard-Errors).
- [not_found](../errors/not-found.md): A `not_found` errors will be returned if there are no effects for operation whose ID matches the `id` argument.
