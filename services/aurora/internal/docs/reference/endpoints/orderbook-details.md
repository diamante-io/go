---
title: Orderbook Details
clientData:
  laboratoryUrl: https://www.diamcircle.org/laboratory/#explorer?resource=order_book&endpoint=details
replacement: https://developers.diamcircle.org/api/aggregations/order-books/
---

People on the Diamcircle network can make [offers](../resources/offer.md) to buy or sell assets.
These offers are summarized by the assets being bought and sold in
[orderbooks](../resources/orderbook.md).

Aurora will return, for each orderbook, a summary of the orderbook and the bids and asks
associated with that orderbook.

This endpoint can also be used in [streaming](../streaming.md) mode so it is possible to use it to
listen as offers are processed in the Diamcircle network.  If called in streaming mode Aurora will
start at the earliest known offer unless a `cursor` is set. In that case it will start from the
`cursor`. You can also set `cursor` value to `now` to only stream offers created since your request
time.

## Request

```
GET /order_book?selling_asset_type={selling_asset_type}&selling_asset_code={selling_asset_code}&selling_asset_issuer={selling_asset_issuer}&buying_asset_type={buying_asset_type}&buying_asset_code={buying_asset_code}&buying_asset_issuer={buying_asset_issuer}&limit={limit}
```

### Arguments

| name | notes | description | example |
| ---- | ----- | ----------- | ------- |
| `selling_asset_type` | required, string | Type of the Asset being sold | `native` |
| `selling_asset_code` | optional, string | Code of the Asset being sold | `USD` |
| `selling_asset_issuer` | optional, string | Account ID of the issuer of the Asset being sold | `GA2HGBJIJKI6O4XEM7CZWY5PS6GKSXL6D34ERAJYQSPYA6X6AI7HYW36` |
| `buying_asset_type` | required, string | Type of the Asset being bought | `credit_alphanum4` |
| `buying_asset_code` | optional, string | Code of the Asset being bought | `BTC` |
| `buying_asset_issuer` | optional, string | Account ID of the issuer of the Asset being bought | `GD6VWBXI6NY3AOOR55RLVQ4MNIDSXE5JSAVXUTF35FRRI72LYPI3WL6Z` |
| `limit` | optional, string | Limit the number of items returned | `20` |

### curl Example Request

```sh
curl "https://aurora-testnet.diamcircle.org/order_book?selling_asset_type=native&buying_asset_type=credit_alphanum4&buying_asset_code=FOO&buying_asset_issuer=GBAUUA74H4XOQYRSOW2RZUA4QL5PB37U3JS5NE3RTB2ELJVMIF5RLMAG&limit=20"
```

### JavaScript Example Request

```javascript
var DiamcircleSdk = require('diamcircle-sdk');
var server = new DiamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

server.orderbook(new DiamcircleSdk.Asset.native(), new DiamcircleSdk.Asset('FOO', 'GBAUUA74H4XOQYRSOW2RZUA4QL5PB37U3JS5NE3RTB2ELJVMIF5RLMAG'))
  .call()
  .then(function(resp) {
    console.log(resp);
  })
  .catch(function(err) {
    console.log(err);
  })
```

### JavaScript Streaming Example

```javascript
var DiamcircleSdk = require('diamcircle-sdk')
var server = new DiamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

var orderbookHandler = function (orderbookResponse) {
  console.log(orderbookResponse);
};

var es = server.orderbook(new DiamcircleSdk.Asset.native(), new DiamcircleSdk.Asset('FOO', 'GBAUUA74H4XOQYRSOW2RZUA4QL5PB37U3JS5NE3RTB2ELJVMIF5RLMAG'))
  .cursor('now')
  .stream({
    onmessage: orderbookHandler
  })
```

## Response

The summary of the orderbook and its bids and asks.

## Example Response
```json
{
  "bids": [
    {
      "price_r": {
        "n": 100000000,
        "d": 12953367
      },
      "price": "7.7200005",
      "amount": "12.0000000"
    }
  ],
  "asks": [
    {
      "price_r": {
        "n": 194,
        "d": 25
      },
      "price": "7.7600000",
      "amount": "238.4804125"
    }
  ],
  "base": {
    "asset_type": "native"
  },
  "counter": {
    "asset_type": "credit_alphanum4",
    "asset_code": "FOO",
    "asset_issuer": "GBAUUA74H4XOQYRSOW2RZUA4QL5PB37U3JS5NE3RTB2ELJVMIF5RLMAG"
  }
}
```

## Possible Errors

- The [standard errors](../errors.md#standard-errors).
