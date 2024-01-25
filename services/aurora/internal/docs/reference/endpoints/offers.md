---
title: Offers
replacement: https://developers.diamcircle.org/api/resources/offers/list/
---

People on the Diamcircle network can make [offers](../resources/offer.md) to buy or sell assets. This
endpoint represents all the current offers, allowing filtering by `seller`, `selling_asset` or `buying_asset`.

## Request

```
GET /offers{?selling_asset_type,selling_asset_issuer,selling_asset_code,buying_asset_type,buying_asset_issuer,buying_asset_code,seller,cursor,limit,order}
```

### Arguments

| name | notes | description | example |
| ---- | ----- | ----------- | ------- |
| `?seller` | optional, string | Account ID of the offer creator  | `GA2HGBJIJKI6O4XEM7CZWY5PS6GKSXL6D34ERAJYQSPYA6X6AI7HYW36` |
| `?selling` | optional, string | Asset being sold | `native` or `EUR:GD6VWBXI6NY3AOOR55RLVQ4MNIDSXE5JSAVXUTF35FRRI72LYPI3WL6Z` |
| `?buying` | optional, string | Asset being bought | `native` or `USD:GD6VWBXI6NY3AOOR55RLVQ4MNIDSXE5JSAVXUTF35FRRI72LYPI3WL6Z` |
| `?cursor` | optional, any, default _null_ | A paging token, specifying where to start returning records from. | `12884905984` |
| `?order`  | optional, string, default `asc` | The order in which to return rows, "asc" or "desc". | `asc` |
| `?limit`  | optional, number, default: `10` | Maximum number of records to return. | `200` |

### curl Example Request

```sh
curl "https://aurora-testnet.diamcircle.org/offers{?selling_asset_type,selling_asset_issuer,selling_asset_code,buying_asset_type,buying_asset_issuer,buying_asset_code,seller,cursor,limit,order}"
```

<!-- ### JavaScript Example Request -->

<!-- ```javascript -->
<!-- var DiamcircleSdk = require('diamcircle-sdk'); -->
<!-- var server = new DiamcircleSdk.Server('https://aurora-testnet.diamcircle.org'); -->

<!-- server.offers('accounts', 'GBYUUJHG6F4EPJGNLERINATVQLNDOFRUD7SGJZ26YZLG5PAYLG7XUSGF') -->
<!--   .call() -->
<!--   .then(function (offerResult) { -->
<!--     console.log(offerResult); -->
<!--   }) -->
<!--   .catch(function (err) { -->
<!--     console.error(err); -->
<!--   }) -->
<!-- ``` -->

<!-- ### JavaScript Streaming Example -->

<!-- ```javascript -->
<!-- var DiamcircleSdk = require('diamcircle-sdk') -->
<!-- var server = new DiamcircleSdk.Server('https://aurora-testnet.diamcircle.org'); -->

<!-- var offerHandler = function (offerResponse) { -->
<!--   console.log(offerResponse); -->
<!-- }; -->

<!-- var es = server.offers('accounts', 'GBYUUJHG6F4EPJGNLERINATVQLNDOFRUD7SGJZ26YZLG5PAYLG7XUSGF') -->
<!--   .cursor('now') -->
<!--   .stream({ -->
<!--     onmessage: offerHandler -->
<!--   }) -->
<!-- ``` -->

## Response

The list of offers.

### Example Response

```json
{
  "_links": {
    "self": {
      "href": "https://aurora-testnet.diamcircle.org/offers?cursor=&limit=10&order=asc"
    },
    "next": {
      "href": "https://aurora-testnet.diamcircle.org/offers?cursor=5443256&limit=10&order=asc"
    },
    "prev": {
      "href": "https://aurora-testnet.diamcircle.org/offers?cursor=5443256&limit=10&order=desc"
    }
  },
  "_embedded": {
    "records": [
      {
        "_links": {
          "self": {
            "href": "https://aurora-testnet.diamcircle.org/offers/5443256"
          },
          "offer_maker": {
            "href": "https://aurora-testnet.diamcircle.org/"
          }
        },
        "id": "5443256",
        "paging_token": "5443256",
        "seller": "GBYUUJHG6F4EPJGNLERINATVQLNDOFRUD7SGJZ26YZLG5PAYLG7XUSGF",
        "selling": {
          "asset_type": "native"
        },
        "buying": {
          "asset_type": "credit_alphanum4",
          "asset_code": "FOO",
          "asset_issuer": "GAGLYFZJMN5HEULSTH5CIGPOPAVUYPG5YSWIYDJMAPIECYEBPM2TA3QR"
        },
        "amount": "10.0000000",
        "price_r": {
          "n": 1,
          "d": 1
        },
        "price": "1.0000000",
        "last_modified_ledger": 694974,
        "last_modified_time": "2019-04-09T17:14:22Z"
      }
    ]
  }
}
```

## Possible Errors

- The [standard errors](../errors.md#standard-errors).
