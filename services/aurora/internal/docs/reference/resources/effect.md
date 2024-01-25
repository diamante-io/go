---
title: Effect
replacement: https://developers.diamcircle.org/api/resources/effects/
---

A successful operation will yield zero or more **effects**.  These effects
represent specific changes that occur in the ledger, but are not necessarily
directly reflected in the [ledger](https://www.diamcircle.org/developers/learn/concepts/ledger.html) or [history](https://github.com/diamcircle/diamcircle-core/blob/master/docs/history.md), as [transactions](https://www.diamcircle.org/developers/learn/concepts/transactions.html) and [operations](https://www.diamcircle.org/developers/learn/concepts/operations.html) are.

## Effect types

We can distinguish 6 effect groups:
- Account effects
- Signer effects
- Trustline effects
- Trading effects
- Data effects
- Misc effects

### Account effects

| Type                                  | Operation                                            |
|---------------------------------------|------------------------------------------------------|
| Account Created                       | create_account                                       |
| Account Removed                       | merge_account                                        |
| Account Credited                      | create_account, payment, path_payment, merge_account |
| Account Debited                       | create_account, payment, path_payment, merge_account |
| Account Thresholds Updated            | set_options                                          |
| Account Home Domain Updated           | set_options                                          |
| Account Flags Updated                 | set_options                                          |
| Account Inflation Destination Updated | set_options                                          |

### Signer effects

| Type           | Operation   |
|----------------|-------------|
| Signer Created | set_options |
| Signer Removed | set_options |
| Signer Updated | set_options |

### Trustline effects

| Type                   | Operation                 |
|------------------------|---------------------------|
| Trustline Created      | change_trust              |
| Trustline Removed      | change_trust              |
| Trustline Updated      | change_trust, allow_trust |
| Trustline Authorized   | allow_trust               |
| Trustline Deauthorized | allow_trust               |

### Trading effects

| Type          | Operation                                                                    |
|---------------|------------------------------------------------------------------------------|
| Offer Created | manage_buy_offer, manage_sell_offer, create_passive_sell_offer               |
| Offer Removed | manage_buy_offer, manage_sell_offer, create_passive_sell_offer, path_payment |
| Offer Updated | manage_buy_offer, manage_sell_offer, create_passive_sell_offer, path_payment |
| Trade         | manage_buy_offer, manage_sell_offer, create_passive_sell_offer, path_payment |
### Data effects

| Type         | Operation   |
|--------------|-------------|
| Data Created | manage_data |
| Data Removed | manage_data |
| Data Updated | manage_data |
### Misc effects

| Type            | Operation     |
|-----------------|---------------|
| Sequence Bumped | bump_sequence |

## Attributes

Attributes depend on effect type.

## Links

| rel       | Example                                                       | Relation                          |
|-----------|---------------------------------------------------------------|-----------------------------------|
| self      | `/effects?order=asc\u0026limit=1`                             |                                   |
| prev      | `/effects?order=desc\u0026limit=1\u0026cursor=141733924865-1` |                                   |
| next      | `/effects?order=asc\u0026limit=1\u0026cursor=141733924865-1`  |                                   |
| operation | `/operations/141733924865`                                    | Operation that created the effect |

## Example

```json
{
  "_embedded": {
    "records": [
      {
        "_links": {
          "operation": {
            "href": "/operations/141733924865"
          },
          "precedes": {
            "href": "/effects?cursor=141733924865-1\u0026order=asc"
          },
          "succeeds": {
            "href": "/effects?cursor=141733924865-1\u0026order=desc"
          }
        },
        "account": "GBS43BF24ENNS3KPACUZVKK2VYPOZVBQO2CISGZ777RYGOPYC2FT6S3K",
        "paging_token": "141733924865-1",
        "starting_balance": "10000000.0",
        "type_i": 0,
        "type": "account_created"
      }
    ]
  },
  "_links": {
    "next": {
      "href": "/effects?order=asc\u0026limit=1\u0026cursor=141733924865-1"
    },
    "prev": {
      "href": "/effects?order=desc\u0026limit=1\u0026cursor=141733924865-1"
    },
    "self": {
      "href": "/effects?order=asc\u0026limit=1\u0026cursor="
    }
  }
}
```

## Endpoints

| Resource                                                                                                                                   | Type       | Resource URI Template           |
|--------------------------------------------------------------------------------------------------------------------------------------------|------------|---------------------------------|
| [All Effects](https://go/blob/master/services/aurora/internal/docs/reference/endpoints/effects-all.md)                 | Collection | `/effects`                      |
| [Operation Effects](https://go/blob/master/services/aurora/internal/docs/reference/endpoints/effects-for-operation.md) | Collection | `/operations/:id/effects`       |
| [Account Effects](https://go/blob/master/services/aurora/internal/docs/reference/endpoints/effects-for-account.md)     | Collection | `/accounts/:account_id/effects` |
| [Ledger Effects](https://go/blob/master/services/aurora/internal/docs/reference/endpoints/effects-for-ledger.md)       | Collection | `/ledgers/:ledger_id/effects`   |
