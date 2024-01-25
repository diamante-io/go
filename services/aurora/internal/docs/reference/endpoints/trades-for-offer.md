---
title: Trades for Offer
clientData:
  laboratoryUrl: https://www.diamcircle.org/laboratory/#explorer?resource=trades&endpoint=for_offer
---

This endpoint represents all [trades](../resources/trade.md) for a given [offer](../resources/offer.md).

This endpoint can also be used in [streaming](../streaming.md) mode, making it possible to listen for new trades for the given offer as they occur on the Diamcircle network.
If called in streaming mode Aurora will start at the earliest known trade unless a `cursor` is set. In that case it will start from the `cursor`. You can also set `cursor` value to `now` to only stream trades created since your request time.
## Request

```
GET /offers/{offer_id}/trades{?cursor,limit,order}
```

### Arguments

| name | notes | description | example |
| ---- | ----- | ----------- | ------- |
| `offer_id` | required, number | ID of an offer | 323223 |
| `?cursor` | optional, any, default _null_ | A paging token, specifying where to start returning records from. | 12884905984 |
| `?order`  | optional, string, default `asc` | The order in which to return rows, "asc" or "desc". | `asc` |
| `?limit`  | optional, number, default: `10` | Maximum number of records to return. | `200` |

### curl Example Request

```sh
curl "https://aurora-testnet.diamcircle.org/offers/323223/trades"
```

### JavaScript Example Request

```js
var DiamcircleSdk = require('diamcircle-sdk');
var server = new DiamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

server.trades()
    .forOffer(323223)
    .call()
    .then(function (tradesResult) {
      console.log(tradesResult);
    })
    .catch(function (err) {
      console.error(err);
    })
```


## Response

This endpoint responds with a list of trades that consumed a given offer. See the [trade resource](../resources/trade.md) for reference.

### Example Response
```json
{
  "_links": {
    "self": {
      "href": "https://aurora-testnet.diamcircle.org/offers/323223/trades?cursor=\u0026limit=10\u0026order=asc"
    },
    "next": {
      "href": "https://aurora-testnet.diamcircle.org/offers/323223/trades?cursor=35789107779080193-0\u0026limit=10\u0026order=asc"
    },
    "prev": {
      "href": "https://aurora-testnet.diamcircle.org/offers/323223/trades?cursor=35789107779080193-0\u0026limit=10\u0026order=desc"
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
            "href": "https://aurora-testnet.diamcircle.org/accounts/GDRCFIQAUEFUQ6GXF5DPRO2M77E4UB7RW7EWI2FTKOW7CWYKZCHSI75K"
          },
          "counter": {
            "href": "https://aurora-testnet.diamcircle.org/accounts/GCUD7CBKTQI4D7ZR7IKHMGXZKKVABML7XFBHV4AIYBOEN5UQFZ5DSPPT"
          },
          "operation": {
            "href": "https://aurora-testnet.diamcircle.org/operations/35789107779080193"
          }
        },
        "id": "35789107779080193-0",
        "paging_token": "35789107779080193-0",
        "ledger_close_time": "2018-04-08T05:58:37Z",
        "base_offer_id": "323223",
        "base_account": "GDRCFIQAUEFUQ6GXF5DPRO2M77E4UB7RW7EWI2FTKOW7CWYKZCHSI75K",
        "base_amount": "912.6607285",
        "base_asset_type": "native",
        "counter_offer_id": "4611686044197195777",   
        "counter_account": "GCUD7CBKTQI4D7ZR7IKHMGXZKKVABML7XFBHV4AIYBOEN5UQFZ5DSPPT",
        "counter_amount": "16.5200719",
        "counter_asset_type": "credit_alphanum4",
        "counter_asset_code": "CM10",
        "counter_asset_issuer": "GBUJJAYHS64L4RDHPLURQJUKSHHPINSAYXYVMWPEF4LECHDKB2EFMKBX",
        "base_is_seller": true,
        "price": {
          "n": 18101,
          "d": 1000000
        }
      }
    ]
  }
}
```

## Example Streaming Event
```cgo
{ _links: 
   { self: { href: '' },
     base: 
      { href: '/accounts/GDJNMHET4DTS7HUHU7IG5DB274OSMHUYA7TRRKOD6ZABHPUW5YWJ4SUD' },
     counter: 
      { href: '/accounts/GCALYDRCCJEUPMV24TAX2N2N3IBX7NUUYZNM7I5FQS5GIEQ4A7EVKUOP' },
     operation: { href: '/operations/47261068505915393' } },
  id: '47261068505915393-0',
  paging_token: '47261068505915393-0',
  ledger_close_time: '2018-09-11T19:42:04Z',
  offer_id: '734529',
  base_account: 'GDJNMHET4DTS7HUHU7IG5DB274OSMHUYA7TRRKOD6ZABHPUW5YWJ4SUD',
  base_amount: '0.0175999',
  base_asset_type: 'credit_alphanum4',
  base_asset_code: 'BOC',
  base_asset_issuer: 'GCTS32RGWRH6RJM62UVZ4UT5ZN5L6B2D3LPGO6Z2NM2EOGVQA7TA6SKO',
  counter_account: 'GCALYDRCCJEUPMV24TAX2N2N3IBX7NUUYZNM7I5FQS5GIEQ4A7EVKUOP',
  counter_amount: '0.0199998',
  counter_asset_type: 'credit_alphanum4',
  counter_asset_code: 'ABC',
  counter_asset_issuer: 'GCTS32RGWRH6RJM62UVZ4UT5ZN5L6B2D3LPGO6Z2NM2EOGVQA7TA6SKO',
  base_is_seller: true,
  price: { n: 2840909, d: 2500000 }
}
```
## Possible Errors

- The [standard errors](../errors.md#Standard-Errors).
- [not_found](../errors/not-found.md): A `not_found` error will be returned if there is no offer whose ID matches the `offer_id` argument.
