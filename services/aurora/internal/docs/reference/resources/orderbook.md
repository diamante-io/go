---
title: Orderbook
replacement: https://developers.diamcircle.org/api/aggregations/order-books/
---

[Orderbooks](https://www.diamcircle.org/developers/learn/concepts/exchange.html) are collections of offers for each issuer and currency pairs.  Let's say you wanted to exchange EUR issued by a particular bank for BTC issued by a particular exchange.  You would look at the orderbook and see who is buying `foo_bank/EUR` and selling `baz_exchange/BTC` and at what prices.

## Attributes
| Attribute    | Type             |                                                                                                                        |
|--------------|------------------|------------------------------------------------------------------------------------------------------------------------|
| bids | object     |  Array of {`price_r`, `price`, `amount`} objects (see [offers](./offer.md)).  These represent prices and amounts accounts are willing to buy for the given `selling` and `buying` pair. |
| asks | object |  Array of {`price_r`, `price`, `amount`} objects (see [offers](./offer.md)).  These represent prices and amounts accounts are willing to sell for the given `selling` and `buying` pair.|
| base | [Asset](http://diamcircle.org/developers/learn/concepts/assets.html) | The Asset this offer wants to sell.|
| counter | [Asset](http://diamcircle.org/developers/learn/concepts/assets.html) | The Asset this offer wants to buy.|

#### Bid Object
|    Attribute     |  Type  |                                                                                                                                |
| ---------------- | ------ | ------------------------------------------------------------------------------------------------------------------------------ |
| price_r              | object | An object of a number numerator and number denominator that represents the bid price. |
| price               | string | The bid price of the asset. A number representing the decimal form of price_r |
| amount              | string | The amount of asset bid offer.  |

#### Ask Object
|    Attribute     |  Type  |                                                                                                                                |
| ---------------- | ------ | ------------------------------------------------------------------------------------------------------------------------------ |
| price_r              | object | An object of a number numerator and number denominator that represents the ask price. |
| price               | string | The ask price of the asset. A number representing the decimal form of price_r |
| amount              | string | The amount of asset ask offer.  |

#### Price_r Object
Price_r is a more precise representation of a bid/ask offer.

|    Attribute     |  Type  |                                                                                                                                |
| ---------------- | ------ | ------------------------------------------------------------------------------------------------------------------------------ |
| n               | number | The numerator.   |
| d              | number | The denominator.  |

Thus to get price you would take n / d.

## Links

This resource has no links.


## Endpoints

| Resource                 | Type       | Resource URI Template                |
|--------------------------|------------|--------------------------------------|
| [Orderbook Details](../endpoints/orderbook-details.md)       | Single | `/orderbook?{orderbook_params}`       |
| [Trades](../endpoints/trades.md)   | Collection | `/trades?{orderbook_params}`       |
