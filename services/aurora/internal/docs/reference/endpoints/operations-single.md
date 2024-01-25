---
title: Operation Details
clientData:
  laboratoryUrl: https://www.diamcircle.org/laboratory/#explorer?resource=operations&endpoint=single
replacement: https://developers.diamcircle.org/api/resources/operations/single/
---

The operation details endpoint provides information on a single
[operation](../resources/operation.md). The operation ID provided in the `id` argument specifies
which operation to load.

### Warning - failed transactions

Operations can be part of successful or failed transactions (failed transactions are also included
in diamcircle ledger). Always check operation status using `transaction_successful` field!

## Request

```
GET /operations/{id}
```

### Arguments

| name | notes | description | example |
| ---- | ----- | ----------- | ------- |
| `id` | required, number | An operation ID. | 2927608622747649 |
| `?join` | optional, string, default: _null_ | Set to `transactions` to include the transactions which created each of the operations in the response. | `transactions` |

### curl Example Request

```sh
curlhttps://diamtestnet.diamcircle.io/operations/2927608622747649
```

### JavaScript Example Request

```javascript
var diamcircleSdk = require('diamcircle-sdk');
var server = new diamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

server.operations()
  .operation('2927608622747649')
  .call()
  .then(function (operationsResult) {
    console.log(operationsResult)
  })
  .catch(function (err) {
    console.log(err)
  })
```

## Response

This endpoint responds with a single Operation.  See [operation resource](../resources/operation.md) for reference.

### Example Response

```json
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
```

## Possible Errors

- The [standard errors](../errors.md#Standard-Errors).
- [not_found](../errors/not-found.md): A `not_found` error will be returned if the
  there is no operation that matches the ID argument, i.e. the operation does not exist.
