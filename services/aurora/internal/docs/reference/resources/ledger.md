---
title: Ledger
replacement: https://developers.diamcircle.org/api/resources/ledgers/
---

A **ledger** resource contains information about a given ledger.

To learn more about the concept of ledgers in the Diamcircle network, take a look at the [Diamcircle ledger concept guide](https://www.diamcircle.org/developers/learn/concepts/ledger.html).

## Attributes

| Attribute                    | Type   |                                                                                                                              |
|------------------------------|--------|------------------------------------------------------------------------------------------------------------------------------|
| id                           | string | The id is a unique identifier for this ledger.                                                                               |
| paging_token                 | number | A [paging token](./page.md) suitable for use as a `cursor` parameter.                                                        |
| hash                         | string | A hex-encoded, lowercase SHA-256 hash of the ledger's [XDR](../../learn/xdr.md)-encoded form.                                |
| prev_hash                    | string | The hash of the ledger that chronologically came before this one.                                                            |
| sequence                     | number | Sequence number of this ledger, suitable for use as the as the :id parameter for url templates that require a ledger number. |
| successful_transaction_count | number | The number of successful transactions in this ledger.                                                                        |
| failed_transaction_count     | number | The number of failed transactions in this ledger.                                                                            |
| operation_count              | number | The number of operations applied in this ledger.                                                                             |
| tx_set_operation_count       | number | The number of operations in this ledger. This number includes operations from failed and successful transactions.            |
| closed_at                    | string | An [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) formatted string of when this ledger was closed.                       |
| total_coins                  | string | The total number of lumens in circulation.                                                                                   |
| fee_pool                     | string | The sum of all transaction fees *(in lumens)* since the last inflation operation. They are redistributed during [inflation]. |
| base_fee                     | number | The [fee] the network charges per operation in a transaction.                                                                |
| base_reserve                 | string | The [reserve][fee] the network uses when calculating an account's minimum balance.                                           |
| max_tx_set_size              | number | The maximum number of transactions validators have agreed to process in a given ledger.                                      |
| protocol_version             | number | The protocol version that the diamcircle network was running when this ledger was committed.                                    |
| header_xdr                   | string | A base64 encoded string of the raw `LedgerHeader` xdr struct for this ledger.                                                |
| base_fee_in_stroops          | number | The [fee] the network charges per operation in a transaction.  Expressed in stroops.                                         |
| base_reserve_in_stroops      | number | The [reserve][fee] the network uses when calculating an account's minimum balance. Expressed in stroops.                     |

## Links
|              | Example                                           | Relation                        | templated |
|--------------|---------------------------------------------------|---------------------------------|-----------|
| self         | `/ledgers/500`                                    |                                 |           |
| effects      | `/ledgers/500/effects/{?cursor,limit,order}`      | The effects in this transaction | true      |
| operations   | `/ledgers/500/operations/{?cursor,limit,order}`   | The operations in this ledger   | true      |
| transactions | `/ledgers/500/transactions/{?cursor,limit,order}` | The transactions in this ledger | true      |


## Example

```json
{
  "_links": {
    "effects": {
      "href": "/ledgers/500/effects/{?cursor,limit,order}",
      "templated": true
    },
    "operations": {
      "href": "/ledgers/500/operations/{?cursor,limit,order}",
      "templated": true
    },
    "self": {
      "href": "/ledgers/500"
    },
    "transactions": {
      "href": "/ledgers/500/transactions/{?cursor,limit,order}",
      "templated": true
    }
  },
  "id": "689f00d4824b8e69330bf4ad7eb10092ff2f8fdb76d4668a41eebb9469ef7f30",
  "paging_token": "2147483648000",
  "hash": "689f00d4824b8e69330bf4ad7eb10092ff2f8fdb76d4668a41eebb9469ef7f30",
  "prev_hash": "b608e110c7cc58200c912140f121af50dc5ef407aabd53b76e1741080aca1cf0",
  "sequence": 500,
  "transaction_count": 0,
  "successful_transaction_count": 0,
  "failed_transaction_count": 0,
  "operation_count": 0,
  "tx_set_operation_count": 0,
  "closed_at": "2015-07-09T21:39:28Z",
  "total_coins": "100000000000.0000000",
  "fee_pool": "0.0025600",
  "base_fee": 100,
  "base_reserve": "10.0000000",
  "max_tx_set_size": 50,
  "protocol_version": 8,
  "header_xdr": "...",
  "base_fee_in_stroops": 100,
  "base_reserve_in_stroops": 100000000
}
```

## Endpoints
| Resource                | Type       | Resource URI Template              |
|-------------------------|------------|------------------------------------|
| [All ledgers](../endpoints/ledgers-all.md)         | Collection | `/ledgers`                         |
| [Single Ledger](../endpoints/ledgers-single.md)       | Single     | `/ledgers/:id`                     |
| [Ledger Transactions](../endpoints/transactions-for-ledger.md) | Collection | `/ledgers/:ledger_id/transactions` |
| [Ledger Operations](../endpoints/operations-for-ledger.md)   | Collection | `/ledgers/:ledger_id/operations`   |
| [Ledger Payments](../endpoints/payments-for-ledger.md)     | Collection | `/ledgers/:ledger_id/payments`     |
| [Ledger Effects](../endpoints/effects-for-ledger.md)      | Collection | `/ledgers/:ledger_id/effects`      |



[inflation]: https://www.diamcircle.org/developers/learn/concepts/inflation.html
[fee]: https://www.diamcircle.org/developers/learn/concepts/fees.html
