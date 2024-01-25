---
title: Offers for Account
clientData:
  laboratoryUrl: https://www.diamcircle.org/laboratory/#explorer?resource=offers&endpoint=for_account
replacement: https://developers.diamcircle.org/api/resources/accounts/offers/
---

People on the diamcircle network can make [offers](../resources/offer.md) to buy or sell assets. This
endpoint represents all the offers a particular account makes.

This endpoint can also be used in [streaming](../streaming.md) mode so it is possible to use it to
listen as offers are processed in the diamcircle network. If called in streaming mode aurora will
start at the earliest known offer unless a `cursor` is set. In that case it will start from the
`cursor`. You can also set `cursor` value to `now` to only stream offers created since your request
time.

## Request

```
GET /accounts/{account}/offers{?cursor,limit,order}
```

### Arguments

| name | notes | description | example |
| ---- | ----- | ----------- | ------- |
| `account` | required, string | Account ID | `GBYUUJHG6F4EPJGNLERINATVQLNDOFRUD7SGJZ26YZLG5PAYLG7XUSGF` |
| `?cursor` | optional, any, default _null_ | A paging token, specifying where to start returning records from. | `12884905984` |
| `?order`  | optional, string, default `asc` | The order in which to return rows, "asc" or "desc". | `asc` |
| `?limit`  | optional, number, default: `10` | Maximum number of records to return. | `200` |

### curl Example Request

```sh
curl "https://aurora-testnet.diamcircle.org/accounts/GBYUUJHG6F4EPJGNLERINATVQLNDOFRUD7SGJZ26YZLG5PAYLG7XUSGF/offers"
```

### JavaScript Example Request

```javascript
var diamcircleSdk = require('diamcircle-sdk');
var server = new diamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

server.offers('accounts', 'GBYUUJHG6F4EPJGNLERINATVQLNDOFRUD7SGJZ26YZLG5PAYLG7XUSGF')
  .call()
  .then(function (offerResult) {
    console.log(offerResult);
  })
  .catch(function (err) {
    console.error(err);
  })
```

### JavaScript Streaming Example

```javascript
var diamcircleSdk = require('diamcircle-sdk')
var server = new diamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

var offerHandler = function (offerResponse) {
  console.log(offerResponse);
};

var es = server.offers('accounts', 'GBYUUJHG6F4EPJGNLERINATVQLNDOFRUD7SGJZ26YZLG5PAYLG7XUSGF')
  .cursor('now')
  .stream({
    onmessage: offerHandler
  })
```

## Response

The list of offers.

**Note:** a response of 200 with an empty records array may either mean there are no offers for
`account_id` or `account_id` does not exist.

### Example Response

```json
{
  "_links": {
    "self": {
      "href": "https://aurora-testnet.diamcircle.org/accounts/GBYUUJHG6F4EPJGNLERINATVQLNDOFRUD7SGJZ26YZLG5PAYLG7XUSGF/offers?cursor=&limit=10&order=asc"
    },
    "next": {
      "href": "https://aurora-testnet.diamcircle.org/accounts/GBYUUJHG6F4EPJGNLERINATVQLNDOFRUD7SGJZ26YZLG5PAYLG7XUSGF/offers?cursor=5443256&limit=10&order=asc"
    },
    "prev": {
      "href": "https://aurora-testnet.diamcircle.org/accounts/GBYUUJHG6F4EPJGNLERINATVQLNDOFRUD7SGJZ26YZLG5PAYLG7XUSGF/offers?cursor=5443256&limit=10&order=desc"
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
            "href": "https://aurora-testnet.diamcircle.org/accounts/GBYUUJHG6F4EPJGNLERINATVQLNDOFRUD7SGJZ26YZLG5PAYLG7XUSGF"
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
