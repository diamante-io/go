---
title: Account Details
clientData:
  laboratoryUrl: https://www.diamcircle.org/laboratory/#explorer?resource=accounts&endpoint=single
replacement: https://developers.diamcircle.org/api/resources/accounts/single/
---

Returns information and links relating to a single [account](../resources/account.md).

The balances section in the returned JSON will also list all the
[trustlines](https://www.diamcircle.org/developers/learn/concepts/assets.html) this account
established. Note this will only return trustlines that have the necessary authorization to work.
Meaning if an account `A` trusts another account `B` that has the
[authorization required](https://www.diamcircle.org/developers/guides/concepts/accounts.html#flags)
flag set, the trustline won't show up until account `B`
[allows](https://www.diamcircle.org/developers/guides/concepts/list-of-operations.html#allow-trust)
account `A` to hold its assets.

## Request

```
GET /accounts/{account}
```

### Arguments

| name | notes | description | example |
| ---- | ----- | ----------- | ------- |
| `account` | required, string | Account ID | GD42RQNXTRIW6YR3E2HXV5T2AI27LBRHOERV2JIYNFMXOBA234SWLQQB |

### curl Example Request

```sh
curl "https://aurora-testnet.diamcircle.org/accounts/GD42RQNXTRIW6YR3E2HXV5T2AI27LBRHOERV2JIYNFMXOBA234SWLQQB"
```

### JavaScript Example Request

```javascript
var DiamcircleSdk = require('diamcircle-sdk');
var server = new DiamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

server.accounts()
  .accountId("GD42RQNXTRIW6YR3E2HXV5T2AI27LBRHOERV2JIYNFMXOBA234SWLQQB")
  .call()
  .then(function (accountResult) {
    console.log(accountResult);
  })
  .catch(function (err) {
    console.error(err);
  })
```

## Response

This endpoint responds with the details of a single account for a given ID. See [account resource](../resources/account.md) for reference.

### Example Response
```json
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
      "is_authorized": true,
      "asset_type": "credit_alphanum4",
      "asset_code": "FOO",
      "asset_issuer": "GAGLYFZJMN5HEULSTH5CIGPOPAVUYPG5YSWIYDJMAPIECYEBPM2TA3QR"
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
```

## Possible Errors

- The [standard errors](../errors.md#Standard-Errors).
- [not_found](../errors/not-found.md): A `not_found` error will be returned if there is no account whose ID matches the `account` argument.
