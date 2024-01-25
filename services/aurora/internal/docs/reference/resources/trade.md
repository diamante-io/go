---
title: Trade
---

A trade represents a fulfilled offer.  For example, let's say that there exists an offer to sell 9 `foo_bank/EUR` for 3 `baz_exchange/BTC` and you make an offer to buy 3 `foo_bank/EUR` for 1 `baz_exchange/BTC`.  Since your offer and the existing one cross, a trade happens.  After the trade completes:

- you are 3 `foo_bank/EUR` richer and 1 `baz_exchange/BTC` poorer
- the maker of the other offer is 1 `baz_exchange/BTC` richer and 3 `foo_bank/EUR` poorer
- your offer is completely fulfilled and no longer exists
- the other offer is partially fulfilled and becomes an offer to sell 6 `foo_bank/EUR` for 2 `baz_exchange/BTC`.  The price of that offer doesn't change, but the amount does.

Trades can also be caused by successful [path payments](https://www.diamcircle.org/developers/learn/concepts/exchange.html), because path payments involve fulfilling offers.

Payments are one-way in that afterwards, the source account has a smaller balance and the destination account of the payment has a bigger one.  Trades are two-way; both accounts increase and decrease their balances.

A trade occurs between two parties - `base` and `counter`. Which is either arbitrary or determined by the calling query.

## Attributes
| Attribute    | Type             |                                                                                                                        |
|--------------|------------------|------------------------------------------------------------------------------------------------------------------------|
| id | string | The ID of this trade. |
| paging_token | string | A [paging token](./page.md) suitable for use as a `cursor` parameter.|
| ledger_close_time | string | An ISO 8601 formatted string of when the ledger with this trade was closed.|
| offer_id | string | DEPRECATED. the sell offer id.
| base_account | string | base party of this trade|
| base_offer_id | string | the base offer id. If this offer was immediately fully consumed this will be a synthetic id
| base_amount | string | amount of base asset that was moved from `base_account` to `counter_account`|
| base_asset_type | string | type of base asset|
| base_asset_code | string | code of base asset|
| base_asset_issuer | string | issuer of base asset|
| counter_offer_id | string | the counter offer id. If this offer was immediately fully consumed this will be a synthetic id
| counter_account | string | counter party of this trade|
| counter_amount | string | amount of counter asset that was moved from `counter_account` to `base_account`|
| counter_asset_type | string | type of counter asset|
| counter_asset_code | string | code of counter asset|
| counter_asset_issuer | string | issuer of counter asset|
| price | object | original offer price, expressed as a rational number. example: {n:7, d:3}
| base_is_seller | boolean | indicates which party of the trade made the sell offer|

#### Price Object
Price is a precise representation of a bid/ask offer.

|    Attribute     |  Type  |                                                                                                                                |
| ---------------- | ------ | ------------------------------------------------------------------------------------------------------------------------------ |
| n               | number | The numerator.   |
| d              | number | The denominator.  |

Thus to get price you would take n / d.

#### Synthetic Offer Ids
Offer ids in the aurora trade resource (base_offer_id, counter_offer_id) are synthetic and don't always reflect the respective diamcircle-core offer ids. This is due to the fact that diamcircle-core does not assign offer ids when an offer gets filled immediately. In these cases, Aurora synthetically generates an offer id for the buying offer, based on the total order id of the offer operation. This allows wallets to aggregate historical trades based on offer ids without adding special handling for edge cases. The exact encoding can be found [here](https://github.com/diamcircle/go/blob/master/services/aurora/internal/db2/history/synt_offer_id.go). 

## Links

| rel          | Example                                                                                           | Description                                                | `templated` |
|--------------|---------------------------------------------------------------------------------------------------|------------------------------------------------------------|-------------|
| base      | `/accounts/{base_account}`      | Link to details about the base account| true        |
| counter | `/accounts/{counter_account}`      | Link to details about the counter account | true        |
| operation | `/operation/{operation_id}` | Link to the operation of the assets bought and sold. | true |

## Endpoints

| Resource                 | Type       | Resource URI Template                |
|--------------------------|------------|--------------------------------------|
| [Trades](../endpoints/trades.md)       | Collection | `/trades`       |
| [Account Trades](../endpoints/trades-for-account.md) | Collection | `/accounts/:account_id/trades`      |
