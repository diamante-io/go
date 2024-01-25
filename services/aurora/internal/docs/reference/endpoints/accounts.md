---
title: Accounts
replacement: https://developers.diamcircle.org/api/resources/accounts/
---

This endpoint allows filtering accounts who have a given `signer` or have a trustline to an `asset`. The result is a list of [accounts](../resources/account.md).

To find all accounts who are trustees to an asset, pass the query parameter `asset` using the canonical representation for an issued assets which is `Code:IssuerAccountID`. Read more about canonical representation of assets in [SEP-0011](https://github.com/diamcircle/diamcircle-protocol/blob/0c675fb3a482183dcf0f5db79c12685acf82a95c/ecosystem/sep-0011.md#values).

### Notes
- The default behavior when filtering by `asset` is to return accounts with `authorized` and `unauthorized` trustlines.

## Request

```
GET /accounts{?signer,asset,cursor,limit,order}
```

### Arguments

| name | notes | description | example |
| ---- | ----- | ----------- | ------- |
| `?signer` | optional, string | Account ID | GD42RQNXTRIW6YR3E2HXV5T2AI27LBRHOERV2JIYNFMXOBA234SWLQQB |
| `?asset` | optional, string | An issued asset represented as "Code:IssuerAccountID". | `USD:GAEDTJ4PPEFVW5XV2S7LUXBEHNQMX5Q2GM562RJGOQG7GVCE5H3HIB4V,native` |
| `?cursor` | optional, default _null_ | A paging token, specifying where to start returning records from. | `GA2HGBJIJKI6O4XEM7CZWY5PS6GKSXL6D34ERAJYQSPYA6X6AI7HYW36` |
| `?order` | optional, string, default `asc` | The order in which to return rows, "asc" or "desc". | `asc` |
| `?limit` | optional, number, default `10` | Maximum number of records to return. | `200` |

### curl Example Request

```sh
curl "https://aurora-testnet.diamcircle.org/accounts?signer=GBPOFUJUHOFTZHMZ63H5GE6NX5KVKQRD6N3I2E5AL3T2UG7HSLPLXN2K"
```

<!-- ### JavaScript Example Request -->

<!-- ```javascript -->
<!-- var DiamcircleSdk = require('diamcircle-sdk'); -->
<!-- var server = new DiamcircleSdk.Server('https://aurora-testnet.diamcircle.org'); -->

<!-- server.accounts(asset: asset) -->
<!--   .call() -->
<!--   .then(function (accountResult) { -->
<!--     console.log(accountResult); -->
<!--   }) -->
<!--   .catch(function (err) { -->
<!--     console.error(err); -->
<!--   }) -->
<!-- ``` -->

## Response

This endpoint responds with the details of all accounts matching the filters. See [account resource](../resources/account.md) for reference.

### Example Response
```json
{
  "_links": {
    "self": {
      "href": "https://aurora-testnet.diamcircle.org/accounts?cursor=\u0026limit=10\u0026order=asc\u0026signer=GBPOFUJUHOFTZHMZ63H5GE6NX5KVKQRD6N3I2E5AL3T2UG7HSLPLXN2K"
    },
    "next": {
      "href": "https://aurora-testnet.diamcircle.org/accounts?cursor=GDRREYWHQWJDICNH4SAH4TT2JRBYRPTDYIMLK4UWBDT3X3ZVVYT6I4UQ\u0026limit=10\u0026order=asc\u0026signer=GBPOFUJUHOFTZHMZ63H5GE6NX5KVKQRD6N3I2E5AL3T2UG7HSLPLXN2K"
    },
    "prev": {
      "href": "https://aurora-testnet.diamcircle.org/accounts?cursor=GDRREYWHQWJDICNH4SAH4TT2JRBYRPTDYIMLK4UWBDT3X3ZVVYT6I4UQ\u0026limit=10\u0026order=desc\u0026signer=GBPOFUJUHOFTZHMZ63H5GE6NX5KVKQRD6N3I2E5AL3T2UG7HSLPLXN2K"
    }
  },
  "_embedded": {
    "records": [
      {
        "_links": {
          "self": {
            "href": "https://aurora-testnet.diamcircle.org/accounts/GD42RQNXTRIW6YR3E2HXV5T2AI27LBRHOERV2JIYNFMXOBA234SWLQQB"
          },
          "transactions": {
            "href": "https://aurora-testnet.diamcircle.org/accounts/GD42RQNXTRIW6YR3E2HXV5T2AI27LBRHOERV2JIYNFMXOBA234SWLQQB/transactions{?cursor,limit,order}",
            "templated": true
          },
          "operations": {
            "href": "https://aurora-testnet.diamcircle.org/accounts/GD42RQNXTRIW6YR3E2HXV5T2AI27LBRHOERV2JIYNFMXOBA234SWLQQB/operations{?cursor,limit,order}",
            "templated": true
          },
          "payments": {
            "href": "https://aurora-testnet.diamcircle.org/accounts/GD42RQNXTRIW6YR3E2HXV5T2AI27LBRHOERV2JIYNFMXOBA234SWLQQB/payments{?cursor,limit,order}",
            "templated": true
          },
          "effects": {
            "href": "https://aurora-testnet.diamcircle.org/accounts/GD42RQNXTRIW6YR3E2HXV5T2AI27LBRHOERV2JIYNFMXOBA234SWLQQB/effects{?cursor,limit,order}",
            "templated": true
          },
          "offers": {
            "href": "https://aurora-testnet.diamcircle.org/accounts/GD42RQNXTRIW6YR3E2HXV5T2AI27LBRHOERV2JIYNFMXOBA234SWLQQB/offers{?cursor,limit,order}",
            "templated": true
          },
          "trades": {
            "href": "https://aurora-testnet.diamcircle.org/accounts/GD42RQNXTRIW6YR3E2HXV5T2AI27LBRHOERV2JIYNFMXOBA234SWLQQB/trades{?cursor,limit,order}",
            "templated": true
          },
          "data": {
            "href": "https://aurora-testnet.diamcircle.org/accounts/GD42RQNXTRIW6YR3E2HXV5T2AI27LBRHOERV2JIYNFMXOBA234SWLQQB/data/{key}",
            "templated": true
          }
        },
        "id": "GD42RQNXTRIW6YR3E2HXV5T2AI27LBRHOERV2JIYNFMXOBA234SWLQQB",
        "paging_token": "",
        "account_id": "GD42RQNXTRIW6YR3E2HXV5T2AI27LBRHOERV2JIYNFMXOBA234SWLQQB",
        "sequence": 7275146318446606,
        "last_modified_ledger": 22379074,
        "subentry_count": 4,
        "thresholds": {
          "low_threshold": 0,
          "med_threshold": 0,
          "high_threshold": 0
        },
        "flags": {
          "auth_required": false,
          "auth_revocable": false,
          "auth_immutable": false
        },
        "balances": [
          {
            "balance": "1000000.0000000",
            "limit": "922337203685.4775807",
            "buying_liabilities": "0.0000000",
            "selling_liabilities": "0.0000000",
            "last_modified_ledger": 632070,
            "asset_type": "credit_alphanum4",
            "asset_code": "FOO",
            "asset_issuer": "GAGLYFZJMN5HEULSTH5CIGPOPAVUYPG5YSWIYDJMAPIECYEBPM2TA3QR",
            "is_authorized": true
          },
          {
            "balance": "10000.0000000",
            "buying_liabilities": "0.0000000",
            "selling_liabilities": "0.0000000",
            "asset_type": "native"
          }
        ],
        "signers": [
          {
            "public_key": "GDLEPBJBC2VSKJCLJB264F2WDK63X4NKOG774A3QWVH2U6PERGDPUCS4",
            "weight": 1,
            "key": "GDLEPBJBC2VSKJCLJB264F2WDK63X4NKOG774A3QWVH2U6PERGDPUCS4",
            "type": "ed25519_public_key"
          },
          {
            "public_key": "GBPOFUJUHOFTZHMZ63H5GE6NX5KVKQRD6N3I2E5AL3T2UG7HSLPLXN2K",
            "weight": 1,
            "key": "GBPOFUJUHOFTZHMZ63H5GE6NX5KVKQRD6N3I2E5AL3T2UG7HSLPLXN2K",
            "type": "sha256_hash"
          },
          {
            "public_key": "GDUDIN23QQTB23Q3Q6GUL6I7CEAQY4CWCFVRXFWPF4UJAQO47SPUFCXG",
            "weight": 1,
            "key": "GDUDIN23QQTB23Q3Q6GUL6I7CEAQY4CWCFVRXFWPF4UJAQO47SPUFCXG",
            "type": "preauth_tx"
          },
          {
            "public_key": "GD42RQNXTRIW6YR3E2HXV5T2AI27LBRHOERV2JIYNFMXOBA234SWLQQB",
            "weight": 1,
            "key": "GD42RQNXTRIW6YR3E2HXV5T2AI27LBRHOERV2JIYNFMXOBA234SWLQQB",
            "type": "ed25519_public_key"
          }
        ],
        "data": {
          "best_friend": "c3Ryb29weQ=="
        }
      }
    ]
  }
}
```

## Possible Errors

- The [standard errors](../errors.md#Standard-Errors).
