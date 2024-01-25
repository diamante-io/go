---
title: Post Transaction
clientData:
  laboratoryUrl: https://www.diamcircle.org/laboratory/#explorer?resource=transactions&endpoint=create
replacement: https://developers.diamcircle.org/api/resources/transactions/post/
---

Posts a new [transaction](../resources/transaction.md) to the diamcircle Network.
Note that creating a valid transaction and signing it properly is the
responsibility of your client library.

Transaction submission and the subsequent validation and inclusion into the
diamcircle Network's ledger is a [complicated and asynchronous
process](https://www.diamcircle.org/developers/learn/concepts/transactions.html#life-cycle).
To reduce the complexity, aurora manages these asynchronous processes for the
client and will wait to hear results from the diamcircle Network before returning
an HTTP response to a client.

Transaction submission to aurora aims to be
[idempotent](https://en.wikipedia.org/wiki/Idempotence#Computer_science_meaning):
a client can submit a given transaction to aurora more than once and aurora
will behave the same each time.  If the transaction has already been
successfully applied to the ledger, aurora will simply return the saved result
and not attempt to submit the transaction again. Only in cases where a
transaction's status is unknown (and thus will have a chance of being included
into a ledger) will a resubmission to the network occur.

Information about [building transactions](https://www.diamcircle.org/developers/js-diamcircle-base/reference/building-transactions.html) in JavaScript.

### Timeout

If you are encountering this error it means that either:

* aurora has not received a confirmation from the Core server that the transaction you are trying to submit to the network was included in a ledger in a timely manner or:
* aurora has not sent a response to a reverse-proxy before in a specified time.

The former case may happen because there was no room for your transaction in the 3 consecutive ledgers. In such case, Core server removes a transaction from a queue. To solve this you can either:

* Keep resubmitting the same transaction (with the same sequence number) and wait until it finally is added to a new ledger or:
* Increase the [fee](../../../guides/concepts/fees.html).

## Request

```
POST /transactions
```

### Arguments

| name | loc  |  notes   |         example        | description |
| ---- | ---- | -------- | ---------------------- | ----------- |
| `tx` | body | required | `AAAAAO`....`f4yDBA==` | Base64 representation of transaction envelope [XDR](../xdr.md) |


### curl Example Request

```sh
curl -X POST \
     -F "tx=AAAAAOo1QK/3upA74NLkdq4Io3DQAQZPi4TVhuDnvCYQTKIVAAAACgAAH8AAAAABAAAAAAAAAAAAAAABAAAAAQAAAADqNUCv97qQO+DS5HauCKNw0AEGT4uE1Ybg57wmEEyiFQAAAAEAAAAAZc2EuuEa2W1PAKmaqVquHuzUMHaEiRs//+ODOfgWiz8AAAAAAAAAAAAAA+gAAAAAAAAAARBMohUAAABAPnnZL8uPlS+c/AM02r4EbxnZuXmP6pQHvSGmxdOb0SzyfDB2jUKjDtL+NC7zcMIyw4NjTa9Ebp4lvONEf4yDBA==" \
  "https://aurora-testnet.diamcircle.org/transactions"
```

## Response

A successful response (i.e. any response with a successful HTTP response code)
indicates that the transaction was successful and has been included into the
ledger.

If the transaction failed or errored, then an error response will be returned. Please see the errors section below.

### Attributes

The response will include all fields from the [transaction resource](../resources/transaction.md).

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

- The [standard errors](../errors.md#Standard_Errors).
- [transaction_failed](../errors/transaction-failed.md): The transaction failed and could not be applied to the ledger.
- [transaction_malformed](../errors/transaction-malformed.md): The transaction could not be decoded and was not submitted to the network.
- [timeout](../errors/timeout.md): No response from the Core server in a timely manner. Please check "Timeout" section above.
