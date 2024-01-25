---
title: All Assets
clientData:
  laboratoryUrl:
replacement: https://developers.diamcircle.org/api/resources/assets/
---

This endpoint represents all [assets](../resources/asset.md).
It will give you all the assets in the system along with various statistics about each.

## Request

```
GET /assets{?asset_code,asset_issuer,cursor,limit,order}
```

### Arguments

| name | notes | description | example |
| ---- | ----- | ----------- | ------- |
| `?asset_code` | optional, string, default _null_ | Code of the Asset to filter by | `USD` |
| `?asset_issuer` | optional, string, default _null_ | Issuer of the Asset to filter by | `GA2HGBJIJKI6O4XEM7CZWY5PS6GKSXL6D34ERAJYQSPYA6X6AI7HYW36` |
| `?cursor` | optional, any, default _null_ | A paging token, specifying where to start returning records from. | `1` |
| `?order` | optional, string, default `asc` | The order in which to return rows, "asc" or "desc", ordered by asset_code then by asset_issuer. | `asc` |
| `?limit` | optional, number, default: `10` | Maximum number of records to return. | `200` |

### curl Example Request

```sh
# Retrieve the 200 assets, ordered alphabetically:
curl "https://aurora-testnet.diamcircle.org/assets?limit=200"
```

### JavaScript Example Request

```javascript
var diamcircleSdk = require('diamcircle-sdk');
var server = new diamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

server.assets()
  .call()
  .then(function (result) {
    console.log(result.records);
  })
  .catch(function (err) {
    console.log(err)
  })
```

## Response

If called normally this endpoint responds with a [page](../resources/page.md) of assets.

### Example Response

```json
{
  "_links": {
    "self": {
      "href": "/assets?order=asc\u0026limit=10\u0026cursor="
    },
    "next": {
      "href": "/assets?order=asc\u0026limit=10\u0026cursor=3"
    },
    "prev": {
      "href": "/assets?order=desc\u0026limit=10\u0026cursor=1"
    }
  },
  "_embedded": {
    "records": [
      {
        "_links": {
          "toml": {
            "href": "https://www.diamcircle.org/.well-known/diamcircle.toml"
          }
        },
        "asset_type": "credit_alphanum12",
        "asset_code": "BANANA",
        "asset_issuer": "GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN",
        "paging_token": "BANANA_GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN_credit_alphanum4",
        "accounts": {
          "authorized": 2126,
          "authorized_to_maintain_liabilities": 32,
          "unauthorized": 5,
          "claimable_balances": 18
        },
        "balances": {
          "authorized": "10000.0000000",
          "authorized_to_maintain_liabilities": "3000.0000000",
          "unauthorized": "4000.0000000",
          "claimable_balances": "2380.0000000"
        },
        "flags": {
          "auth_required": true,
          "auth_revocable": false
        }
      },
      {
        "_links": {
          "toml": {
            "href": "https://www.diamcircle.org/.well-known/diamcircle.toml"
          }
        },
        "asset_type": "credit_alphanum4",
        "asset_code": "BTC",
        "asset_issuer": "GBAUUA74H4XOQYRSOW2RZUA4QL5PB37U3JS5NE3RTB2ELJVMIF5RLMAG",
        "paging_token": "BTC_GBAUUA74H4XOQYRSOW2RZUA4QL5PB37U3JS5NE3RTB2ELJVMIF5RLMAG_credit_alphanum4",
        "accounts": {
          "authorized": 32,
          "authorized_to_maintain_liabilities": 124,
          "unauthorized": 6,
          "claimable_balances": 18
        },
        "balances": {
          "authorized": "5000.0000000",
          "authorized_to_maintain_liabilities": "8000.0000000",
          "unauthorized": "2000.0000000",
          "claimable_balances": "1200.0000000"
        },
        "flags": {
          "auth_required": false,
          "auth_revocable": false
        }
      },
      {
        "_links": {
          "toml": {
            "href": "https://www.diamcircle.org/.well-known/diamcircle.toml"
          }
        },
        "asset_type": "credit_alphanum4",
        "asset_code": "USD",
        "asset_issuer": "GBAUUA74H4XOQYRSOW2RZUA4QL5PB37U3JS5NE3RTB2ELJVMIF5RLMAG",
        "paging_token": "USD_GBAUUA74H4XOQYRSOW2RZUA4QL5PB37U3JS5NE3RTB2ELJVMIF5RLMAG_credit_alphanum4",
        "accounts": {
          "authorized": 91547871,
          "authorized_to_maintain_liabilities": 45773935,
          "unauthorized": 22886967,
          "claimable_balances": 11443483
        },
        "balances": {
          "authorized": "1000000000.0000000",
          "authorized_to_maintain_liabilities": "500000000.0000000",
          "unauthorized": "250000000.0000000",
          "claimable_balances": "12500000.0000000"
        },
        "flags": {
          "auth_required": false,
          "auth_revocable": false
        }
      }
    ]
  }
}
```

## Possible Errors

- The [standard errors](../errors.md#standard-errors).
