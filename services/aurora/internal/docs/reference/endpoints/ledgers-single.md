---
title: Ledger Details
clientData:
  laboratoryUrl: https://www.diamcircle.org/laboratory/#explorer?resource=ledgers&endpoint=single
replacement: https://developers.diamcircle.org/api/resources/ledgers/single/
---

The ledger details endpoint provides information on a single [ledger](../resources/ledger.md).

## Request

```
GET /ledgers/{sequence}
```

### Arguments

|  name  |  notes  | description | example |
| ------ | ------- | ----------- | ------- |
| `sequence` | required, number | Ledger Sequence | `69859` |

### curl Example Request

```sh
curl "https://aurora-testnet.diamcircle.org/ledgers/69859"
```

### JavaScript Example Request

```js
var diamcircleSdk = require('diamcircle-sdk')
var server = new diamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

server.ledgers()
  .ledger('69858')
  .call()
  .then(function(ledgerResult) {
    console.log(ledgerResult)
  })
  .catch(function(err) {
    console.log(err)
  })

```
## Response

This endpoint responds with a single Ledger.  See [ledger resource](../resources/ledger.md) for reference.

### Example Response

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
  "tx_set_operation_count": 0,
  "closed_at": "2015-07-20T15:51:52Z",
  "total_coins": "100000000000.0000000",
  "fee_pool": "0.0025600",
  "base_fee_in_stroops": 100,
  "base_reserve_in_stroops": 100000000,
  "max_tx_set_size": 50
}
```

## Errors

- The [standard errors](../errors.md#Standard-Errors).
- [not_found](../errors/not-found.md): A `not_found` error will be returned if there is no ledger whose sequence number matches the `sequence` argument.
