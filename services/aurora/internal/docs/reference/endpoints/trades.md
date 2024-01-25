---
title: Trades
replacement: https://developers.diamcircle.org/api/resources/trades/
---

People on the diamcircle network can make [offers](../resources/offer.md) to buy or sell assets. When
an offer is fully or partially fulfilled, a [trade](../resources/trade.md) happens.

Trades can be filtered for a specific orderbook, defined by an asset pair: `base` and `counter`.

This endpoint can also be used in [streaming](../streaming.md) mode, making it possible to listen
for new trades as they occur on the diamcircle network.

If called in streaming mode aurora will start at the earliest known trade unless a `cursor` is
set. In that case it will start from the `cursor`. You can also set `cursor` value to `now` to only
stream trades created since your request time.

## Request

```
GET /trades?base_asset_type={base_asset_type}&base_asset_code={base_asset_code}&base_asset_issuer={base_asset_issuer}&counter_asset_type={counter_asset_type}&counter_asset_code={counter_asset_code}&counter_asset_issuer={counter_asset_issuer}
```

### Arguments

| name | notes | description | example |
| ---- | ----- | ----------- | ------- |
| `base_asset_type` | optional, string | Type of base asset | `native` |
| `base_asset_code` | optional, string | Code of base asset, not required if type is `native` | `USD` |
| `base_asset_issuer` | optional, string | Issuer of base asset, not required if type is `native` | 'GA2HGBJIJKI6O4XEM7CZWY5PS6GKSXL6D34ERAJYQSPYA6X6AI7HYW36' |
| `counter_asset_type` | optional, string | Type of counter asset  | `credit_alphanum4` |
| `counter_asset_code` | optional, string | Code of counter asset, not required if type is `native` | `BTC` |
| `counter_asset_issuer` | optional, string | Issuer of counter asset, not required if type is `native` | 'GD6VWBXI6NY3AOOR55RLVQ4MNIDSXE5JSAVXUTF35FRRI72LYPI3WL6Z' |
| `offer_id` | optional, string | filter for by a specific offer id | `283606` |
| `?cursor` | optional, any, default _null_ | A paging token, specifying where to start returning records from. | `12884905984` |
| `?order`  | optional, string, default `asc` | The order, in terms of timeline, in which to return rows, "asc" or "desc". | `asc` |
| `?limit`  | optional, number, default: `10` | Maximum number of records to return. | `200` |

### curl Example Request
```sh
curl https://aurora.diamcircle.org/trades?base_asset_type=native&counter_asset_code=SLT&counter_asset_issuer=GCKA6K5PCQ6PNF5RQBF7PQDJWRHO6UOGFMRLK3DYHDOI244V47XKQ4GP&counter_asset_type=credit_alphanum4&limit=2&order=desc
```

### JavaScript Example Request

```javascript
var diamcircleSdk = require('diamcircle-sdk');
var server = new diamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

server.trades()
  .call()
  .then(function (tradesResult) {
    console.log(tradesResult.records);
  })
  .catch(function (err) {
    console.log(err)
  })
```

### JavaScript Example Streaming Request

```javascript
var diamcircleSdk = require('diamcircle-sdk')
var server = new diamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

var tradesHandler = function (tradeResponse) {
  console.log(tradeResponse);
};

var es = server.trades()
  .cursor('now')
  .stream({
  onmessage: tradesHandler
})
```

## Response

The list of trades. `base` and `counter` in the records will match the asset pair filter order. If an asset pair is not specified, the order is arbitrary.

### Example Response
```json
{
  "_links": {
    "self": {
      "href": "https://aurora-testnet.diamcircle.org/trades?cursor=&limit=10&order=asc"
    },
    "next": {
      "href": "https://aurora-testnet.diamcircle.org/trades?cursor=6025839120434-0&limit=10&order=asc"
    },
    "prev": {
      "href": "https://aurora-testnet.diamcircle.org/trades?cursor=6012954218535-0&limit=10&order=desc"
    }
  },
  "_embedded": {
    "records": [
      {
        "_links": {
          "self": {
            "href": ""
          },
          "base": {
            "href": "https://aurora-testnet.diamcircle.org/accounts/GAQHWQYBBW272OOXNQMMLCA5WY2XAZPODGB7Q3S5OKKIXVESKO55ZQ7C"
          },
          "counter": {
            "href": "https://aurora-testnet.diamcircle.org/accounts/GCYN7MI6VXVRP74KR6MKBAW2ELLCXL6QCY5H4YQ62HVWZWMCE6Y232UC"
          },
          "operation": {
            "href": "https://aurora-testnet.diamcircle.org/operations/6012954218535"
          }
        },
        "id": "6012954218535-0",
        "paging_token": "6012954218535-0",
        "ledger_close_time": "2019-02-27T11:54:53Z",
        "offer_id": "37",
        "base_offer_id": "4611692031381606439",
        "base_account": "GAQHWQYBBW272OOXNQMMLCA5WY2XAZPODGB7Q3S5OKKIXVESKO55ZQ7C",
        "base_amount": "25.6687300",
        "base_asset_type": "credit_alphanum4",
        "base_asset_code": "DSQ",
        "base_asset_issuer": "GBDQPTQJDATT7Z7EO4COS4IMYXH44RDLLI6N6WIL5BZABGMUOVMLWMQF",
        "counter_offer_id": "37",
        "counter_account": "GCYN7MI6VXVRP74KR6MKBAW2ELLCXL6QCY5H4YQ62HVWZWMCE6Y232UC",
        "counter_amount": "1.0265563",
        "counter_asset_type": "credit_alphanum4",
        "counter_asset_code": "USD",
        "counter_asset_issuer": "GAA4MFNZGUPJAVLWWG6G5XZJFZDHLKQNG3Q6KB24BAD6JHNNVXDCF4XG",
        "base_is_seller": false,
        "price": {
          "n": 10000000,
          "d": 250046977
        }
      },
      {
        "_links": {
          "self": {
            "href": ""
          },
          "base": {
            "href": "https://aurora-testnet.diamcircle.org/accounts/GAQHWQYBBW272OOXNQMMLCA5WY2XAZPODGB7Q3S5OKKIXVESKO55ZQ7C"
          },
          "counter": {
            "href": "https://aurora-testnet.diamcircle.org/accounts/GCYN7MI6VXVRP74KR6MKBAW2ELLCXL6QCY5H4YQ62HVWZWMCE6Y232UC"
          },
          "operation": {
            "href": "https://aurora-testnet.diamcircle.org/operations/6025839120385"
          }
        },
        "id": "6025839120385-0",
        "paging_token": "6025839120385-0",
        "ledger_close_time": "2019-02-27T11:55:09Z",
        "offer_id": "1",
        "base_offer_id": "4611692044266508289",
        "base_account": "GAQHWQYBBW272OOXNQMMLCA5WY2XAZPODGB7Q3S5OKKIXVESKO55ZQ7C",
        "base_amount": "1434.4442973",
        "base_asset_type": "credit_alphanum4",
        "base_asset_code": "DSQ",
        "base_asset_issuer": "GBDQPTQJDATT7Z7EO4COS4IMYXH44RDLLI6N6WIL5BZABGMUOVMLWMQF",
        "counter_offer_id": "1",
        "counter_account": "GCYN7MI6VXVRP74KR6MKBAW2ELLCXL6QCY5H4YQ62HVWZWMCE6Y232UC",
        "counter_amount": "0.5622050",
        "counter_asset_type": "credit_alphanum4",
        "counter_asset_code": "SXRT",
        "counter_asset_issuer": "GAIOQ3UYK5NYIZY5ZFAG4JBN4O37NAVFKZM5YDYEB6YEFBZSZ5KDCUFO",
        "base_is_seller": false,
        "price": {
          "n": 642706,
          "d": 1639839483
        }
      }
    ]
  }
}
```

## Possible Errors

- The [standard errors](../errors.md#standard-errors).
