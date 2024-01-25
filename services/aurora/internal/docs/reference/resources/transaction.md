---
title: Transaction
replacement: https://developers.diamcircle.org/api/resources/transactions/
---

**Transactions** are the basic unit of change in the diamcircle Network.

A transaction is a grouping of [operations](./operation.md).

To learn more about the concept of transactions in the diamcircle network, take a look at the [diamcircle transactions concept guide](https://www.diamcircle.org/developers/learn/concepts/transactions.html).

## Attributes

| Attribute               | Type                     |                                                                                                                                |
|-------------------------|--------------------------|--------------------------------------------------------------------------------------------------------------------------------|
| id                      | string                   | The canonical id of this transaction, suitable for use as the :id parameter for url templates that require a transaction's ID. |
| paging_token            | string                   | A [paging token](./page.md) suitable for use as the `cursor` parameter to transaction collection resources.                    |
| successful              | bool                     | Indicates if transaction was successful or not.                                                                                |
| hash                    | string                   | A hex-encoded, lowercase SHA-256 hash of the transaction's [XDR](../../learn/xdr.md)-encoded form.                             |
| ledger                  | number                   | Sequence number of the ledger in which this transaction was applied.                                                           |
| created_at              | ISO8601 string           |                                                                                                                                |
| fee_account             | string                   | The account which paid for the transaction fees                                                                                |
| source_account          | string                   |                                                                                                                                |
| source_account_sequence | string                   |                                                                                                                                |
| max_fee                 | number                   | The the maximum fee the fee account was willing to pay.                                                                        |
| fee_charged             | number                   | The fee paid by the fee account of this transaction when the transaction was applied to the ledger.                            |
| operation_count         | number                   | The number of operations that are contained within this transaction.                                                           |
| envelope_xdr            | string                   | A base64 encoded string of the raw `TransactionEnvelope` xdr struct for this transaction                                       |
| result_xdr              | string                   | A base64 encoded string of the raw `TransactionResult` xdr struct for this transaction                                         |
| result_meta_xdr         | string                   | A base64 encoded string of the raw `TransactionMeta` xdr struct for this transaction                                           |
| fee_meta_xdr            | string                   | A base64 encoded string of the raw `LedgerEntryChanges` xdr struct produced by taking fees for this transaction.               |
| memo_type               | string                   | The type of memo set in the transaction. Possible values are `none`, `text`, `id`, `hash`, and `return`.                       |
| memo                    | string                   | The string representation of the memo set in the transaction. When `memo_type` is `id`, the `memo` is a decimal string representation of an unsigned 64 bit integer. When `memo_type` is `hash` or `return`, the `memo` is a base64 encoded string. When `memo_type` is `text`, the `memo` is a unicode string. However, if the original memo byte sequence in the transaction XDR is not valid unicode, aurora will replace any invalid byte sequences with the utf-8 replacement character. Note this field is only present when `memo_type` is not `none`. |
| memo_bytes              | string                   | A base64 encoded string of the memo bytes set in the transaction's xdr envelope. Note this field is only present when `memo_type` is `text`. |
| signatures              | string[]                 | An array of signatures used to sign this transaction                                                                           |
| valid_after             | RFC3339 date-time string |                                                                                                                                |
| valid_before            | RFC3339 date-time string |                                                                                                                                |
| fee_bump_transaction    | object                   | This object is only present if the transaction is a fee bump transaction or is wrapped by a fee bump transaction. The object has two fields: `hash` (the hash of the fee bump transaction) and `signatures` (the signatures present in the fee bump transaction envelope)                                                                                                                               |
| inner_transaction       | object                   | This object is only present if the transaction is a fee bump transaction or is wrapped by a fee bump transaction. The object has three fields: `hash` (the hash of the inner transaction wrapped by the fee bump transaction), `max_fee` (the max fee set in the inner transaction), and `signatures` (the signatures present in the inner transaction envelope)                                                                                                                               |

## Links

| rel        | Example                                                                                                                                              | Description                                                                                |
|------------|------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------|
| self       | `https://aurora-testnet.diamcircle.org/transactions/cb9a25394acb6fe0d1d9bdea5afc01cafe2c6fde59a96ddceb2564a65780a81f`                                  |                                                                                            |
| account    | `https://aurora-testnet.diamcircle.org/accounts/GCDLRUXOD6KA53G5ILL435TZAISNLPS4EKIHSOVY3MVD3DVJ333NO4DT`                                              | The source [account](../endpoints/accounts-single.md) for this transaction.                          |
| ledger     | `https://aurora-testnet.diamcircle.org/ledgers/2352988`                                                                                                | The [ledger](../endpoints/ledgers-single.md) in which this transaction was applied.                  |
| operations | `https://aurora-testnet.diamcircle.org/transactions/cb9a25394acb6fe0d1d9bdea5afc01cafe2c6fde59a96ddceb2564a65780a81f/operations{?cursor,limit,order}"` | [Operations](../endpoints/operations-for-transaction.md) included in this transaction.               |
| effects    | `https://aurora-testnet.diamcircle.org/transactions/cb9a25394acb6fe0d1d9bdea5afc01cafe2c6fde59a96ddceb2564a65780a81f/effects{?cursor,limit,order}"`    | [Effects](../endpoints/effects-for-transaction.md) which resulted by operations in this transaction. |
| precedes   | `https://aurora-testnet.diamcircle.org/transactions?order=asc&cursor=10106006507900928`                                                                | A collection of transactions that occur after this transaction.                            |
| succeeds   | `https://aurora-testnet.diamcircle.org/transactions?order=desc&cursor=10106006507900928`                                                               | A collection of transactions that occur before this transaction.                           |

## Example

```json
{
  "_links": {
    "self": {
      "href": "https://aurora-testnet.diamcircle.org/transactions/cb9a25394acb6fe0d1d9bdea5afc01cafe2c6fde59a96ddceb2564a65780a81f"
    },
    "account": {
      "href": "https://aurora-testnet.diamcircle.org/accounts/GCDLRUXOD6KA53G5ILL435TZAISNLPS4EKIHSOVY3MVD3DVJ333NO4DT"
    },
    "ledger": {
      "href": "https://aurora-testnet.diamcircle.org/ledgers/2352988"
    },
    "operations": {
      "href": "https://aurora-testnet.diamcircle.org/transactions/cb9a25394acb6fe0d1d9bdea5afc01cafe2c6fde59a96ddceb2564a65780a81f/operations{?cursor,limit,order}",
      "templated": true
    },
    "effects": {
      "href": "https://aurora-testnet.diamcircle.org/transactions/cb9a25394acb6fe0d1d9bdea5afc01cafe2c6fde59a96ddceb2564a65780a81f/effects{?cursor,limit,order}",
      "templated": true
    },
    "precedes": {
      "href": "https://aurora-testnet.diamcircle.org/transactions?order=asc&cursor=10106006507900928"
    },
    "succeeds": {
      "href": "https://aurora-testnet.diamcircle.org/transactions?order=desc&cursor=10106006507900928"
    }
  },
  "id": "cb9a25394acb6fe0d1d9bdea5afc01cafe2c6fde59a96ddceb2564a65780a81f",
  "paging_token": "10106006507900928",
  "successful": true,
  "hash": "cb9a25394acb6fe0d1d9bdea5afc01cafe2c6fde59a96ddceb2564a65780a81f",
  "ledger": 2352988,
  "created_at": "2019-02-21T21:44:13Z",
  "source_account": "GCDLRUXOD6KA53G5ILL435TZAISNLPS4EKIHSOVY3MVD3DVJ333NO4DT",
  "fee_account": "GCDLRUXOD6KA53G5ILL435TZAISNLPS4EKIHSOVY3MVD3DVJ333NO4DT",
  "source_account_sequence": "10105916313567234",
  "max_fee": 100,
  "fee_charged":100,
  "operation_count": 1,
  "envelope_xdr": "AAAAAIa40u4flA7s3ULXzfZ5AiTVvlwikHk6uNsqPY6p3vbXAAAAZAAj50cAAAACAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAABAAAAB2Fmc2RmYXMAAAAAAQAAAAAAAAABAAAAAIa40u4flA7s3ULXzfZ5AiTVvlwikHk6uNsqPY6p3vbXAAAAAAAAAAEqBfIAAAAAAAAAAAGp3vbXAAAAQKElK3CoNo1f8fWIGeJm98lw2AaFiyVVFhx3uFok0XVW3MHV9MubtEhfA+n1iLPrxmzHtHfmZsumWk+sOEQlSwI=",
  "result_xdr": "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=",
  "result_meta_xdr": "AAAAAQAAAAIAAAADACPnXAAAAAAAAAAAhrjS7h+UDuzdQtfN9nkCJNW+XCKQeTq42yo9jqne9tcAAAAXSHbnOAAj50cAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABACPnXAAAAAAAAAAAhrjS7h+UDuzdQtfN9nkCJNW+XCKQeTq42yo9jqne9tcAAAAXSHbnOAAj50cAAAACAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAA==",
  "fee_meta_xdr": "AAAAAgAAAAMAI+dTAAAAAAAAAACGuNLuH5QO7N1C1832eQIk1b5cIpB5OrjbKj2Oqd721wAAABdIduecACPnRwAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAI+dcAAAAAAAAAACGuNLuH5QO7N1C1832eQIk1b5cIpB5OrjbKj2Oqd721wAAABdIduc4ACPnRwAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
  "memo_type": "text",
  "memo": "afsdfas",
  "valid_after": "1970-01-01T00:00:00Z",
  "signatures": [
    "oSUrcKg2jV/x9YgZ4mb3yXDYBoWLJVUWHHe4WiTRdVbcwdX0y5u0SF8D6fWIs+vGbMe0d+Zmy6ZaT6w4RCVLAg=="
  ]
}
```

## Endpoints

| Resource                                               | Type       | Resource URI Template                |
|--------------------------------------------------------|------------|--------------------------------------|
| [All Transactions](../endpoints/transactions-all.md)             | Collection | `/transactions` (`GET`)              |
| [Post Transaction](../endpoints/transactions-create.md)          | Action     | `/transactions`  (`POST`)            |
| [Transaction Details](../endpoints/transactions-single.md)       | Single     | `/transactions/:id`                  |
| [Account Transactions](../endpoints/transactions-for-account.md) | Collection | `/accounts/:account_id/transactions` |
| [Ledger Transactions](../endpoints/transactions-for-ledger.md)   | Collection | `/ledgers/:ledger_id/transactions`   |


## Submitting transactions
To submit a new transaction to diamcircle network, it must first be built and signed locally. Then you can submit a hex representation of your transaction’s [XDR](../xdr.md) to the `/transactions` endpoint. Read more about submitting transactions in [Post Transaction](../endpoints/transactions-create.md) doc.
