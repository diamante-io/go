---
title: Transaction Details
clientData:
  laboratoryUrl: https://www.diamcircle.org/laboratory/#explorer?resource=transactions&endpoint=single
replacement: https://developers.diamcircle.org/api/resources/transactions/single/
---

The transaction details endpoint provides information on a single
[transaction](../resources/transaction.md). The transaction hash provided in the `hash` argument
specifies which transaction to load.

### Warning - failed transactions

Transaction can be successful or failed (failed transactions are also included in diamcircle ledger).
Always check it's status using `successful` field!

## Request

```
GET /transactions/{hash}
```

### Arguments

|  name  |  notes  | description | example |
| ------ | ------- | ----------- | ------- |
| `hash` | required, string | A transaction hash, hex-encoded, lowercase. | 264226cb06af3b86299031884175155e67a02e0a8ad0b3ab3a88b409a8c09d5c |

### curl Example Request

```sh
curl "https://aurora-testnet.diamcircle.org/transactions/264226cb06af3b86299031884175155e67a02e0a8ad0b3ab3a88b409a8c09d5c"
```

### JavaScript Example Request

```javascript
var diamcircleSdk = require('diamcircle-sdk');
var server = new diamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

server.transactions()
  .transaction("264226cb06af3b86299031884175155e67a02e0a8ad0b3ab3a88b409a8c09d5c")
  .call()
  .then(function (transactionResult) {
    console.log(transactionResult)
  })
  .catch(function (err) {
    console.log(err)
  })
```

## Response

This endpoint responds with a single Transaction.  See [transaction resource](../resources/transaction.md) for reference.

### Example Response

```json
{
  "_links": {
    "self": {
      "href": "https://aurora-testnet.diamcircle.org/transactions/264226cb06af3b86299031884175155e67a02e0a8ad0b3ab3a88b409a8c09d5c"
    },
    "account": {
      "href": "https://aurora-testnet.diamcircle.org/accounts/GAIH3ULLFQ4DGSECF2AR555KZ4KNDGEKN4AFI4SU2M7B43MGK3QJZNSR"
    },
    "ledger": {
      "href": "https://aurora-testnet.diamcircle.org/ledgers/697121"
    },
    "operations": {
      "href": "https://aurora-testnet.diamcircle.org/transactions/264226cb06af3b86299031884175155e67a02e0a8ad0b3ab3a88b409a8c09d5c/operations{?cursor,limit,order}",
      "templated": true
    },
    "effects": {
      "href": "https://aurora-testnet.diamcircle.org/transactions/264226cb06af3b86299031884175155e67a02e0a8ad0b3ab3a88b409a8c09d5c/effects{?cursor,limit,order}",
      "templated": true
    },
    "precedes": {
      "href": "https://aurora-testnet.diamcircle.org/transactions?order=asc&cursor=2994111896358912"
    },
    "succeeds": {
      "href": "https://aurora-testnet.diamcircle.org/transactions?order=desc&cursor=2994111896358912"
    }
  },
  "id": "264226cb06af3b86299031884175155e67a02e0a8ad0b3ab3a88b409a8c09d5c",
  "paging_token": "2994111896358912",
  "successful": true,
  "hash": "264226cb06af3b86299031884175155e67a02e0a8ad0b3ab3a88b409a8c09d5c",
  "ledger": 697121,
  "created_at": "2019-04-09T20:14:25Z",
  "source_account": "GAIH3ULLFQ4DGSECF2AR555KZ4KNDGEKN4AFI4SU2M7B43MGK3QJZNSR",
  "fee_account": "GAIH3ULLFQ4DGSECF2AR555KZ4KNDGEKN4AFI4SU2M7B43MGK3QJZNSR",
  "source_account_sequence": "4660039994869",
  "fee_charged": 100,
  "max_fee": 100,
  "operation_count": 1,
  "envelope_xdr": "AAAAABB90WssODNIgi6BHveqzxTRmIpvAFRyVNM+Hm2GVuCcAAAAZAAABD0AB031AAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAFIMRkFZ9gZifhRSlklQpsz/9P04Earv0dzS3MkIM1cYAAAAXSHboAAAAAAAAAAABhlbgnAAAAEA+biIjrDy8yi+SvhFElIdWGBRYlDscnSSHkPchePy2JYDJn4wvJYDBumXI7/NmttUey3+cGWbBFfnnWh1H5EoD",
  "result_xdr": "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=",
  "result_meta_xdr": "AAAAAQAAAAIAAAADAAqjIQAAAAAAAAAAEH3Rayw4M0iCLoEe96rPFNGYim8AVHJU0z4ebYZW4JwBOLmYhGq/IAAABD0AB030AAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAqjIQAAAAAAAAAAEH3Rayw4M0iCLoEe96rPFNGYim8AVHJU0z4ebYZW4JwBOLmYhGq/IAAABD0AB031AAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMACqMhAAAAAAAAAAAQfdFrLDgzSIIugR73qs8U0ZiKbwBUclTTPh5thlbgnAE4uZiEar8gAAAEPQAHTfUAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEACqMhAAAAAAAAAAAQfdFrLDgzSIIugR73qs8U0ZiKbwBUclTTPh5thlbgnAE4uYE789cgAAAEPQAHTfUAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAACqMhAAAAAAAAAAAUgxGQVn2BmJ+FFKWSVCmzP/0/TgRqu/R3NLcyQgzVxgAAABdIdugAAAqjIQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
  "fee_meta_xdr": "AAAAAgAAAAMACqMgAAAAAAAAAAAQfdFrLDgzSIIugR73qs8U0ZiKbwBUclTTPh5thlbgnAE4uZiEar+EAAAEPQAHTfQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEACqMhAAAAAAAAAAAQfdFrLDgzSIIugR73qs8U0ZiKbwBUclTTPh5thlbgnAE4uZiEar8gAAAEPQAHTfQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
  "memo_type": "none",
  "signatures": [
    "Pm4iI6w8vMovkr4RRJSHVhgUWJQ7HJ0kh5D3IXj8tiWAyZ+MLyWAwbplyO/zZrbVHst/nBlmwRX551odR+RKAw=="
  ]
}
```

## Possible Errors

- The [standard errors](../errors.md#Standard-Errors).
- [not_found](../errors/not-found.md): A `not_found` error will be returned if there is no
  transaction whose ID matches the `hash` argument.
