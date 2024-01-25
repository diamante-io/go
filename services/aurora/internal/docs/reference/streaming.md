---
title: Streaming
replacement: https://developers.diamcircle.org/api/introduction/streaming/
---

## Streaming

Certain endpoints in aurora can be called in streaming mode using Server-Sent Events. This mode will keep the connection to aurora open and aurora will continue to return responses as ledgers close. All parameters for the endpoints that allow this mode are the same. The way a caller initiates this mode is by setting `Accept: text/event-stream` in the HTTP header when you make the request.
You can read an example of using the streaming mode in the [Follow Received Payments](./tutorials/follow-received-payments.md) tutorial.

Endpoints that currently support streaming:
* [Account](./endpoints/accounts-single.md)
* [Effects](./endpoints/effects-all.md)
* [Ledgers](./endpoints/ledgers-all.md)
* [Offers](./endpoints/offers-for-account.md)
* [Operations](./endpoints/operations-all.md)
* [Orderbook](./endpoints/orderbook-details.md)
* [Payments](./endpoints/payments-all.md)
* [Transactions](./endpoints/transactions-all.md)
* [Trades](./endpoints/trades.md)
