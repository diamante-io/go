---
title: Find Payment Paths
clientData:
  laboratoryUrl: https://www.diamcircle.org/laboratory/#explorer?resource=paths&endpoint=all
replacement: https://developers.diamcircle.org/api/aggregations/paths/
---

**Note**: This endpoint will be deprecated, use [/path/strict-receive](./path-finding-strict-receive.md) instead. There are no differences between both endpoints, `/paths` is an alias for `/path/strict-receive`.


The Diamcircle Network allows payments to be made across assets through _path payments_.  A path
payment specifies a series of assets to route a payment through, from source asset (the asset
debited from the payer) to destination asset (the asset credited to the payee).

A path search is specified using:

- The destination account id
- The source account id
- The asset and amount that the destination account should receive

As part of the search, aurora will load a list of assets available to the source account id and
will find any payment paths from those source assets to the desired destination asset. The search's
amount parameter will be used to determine if there a given path can satisfy a payment of the
desired amount.

## Request

```
GET /paths?destination_account={da}&source_account={sa}&destination_asset_type={at}&destination_asset_code={ac}&destination_asset_issuer={di}&destination_amount={amount}
```

## Arguments

| name | notes | description | example |
| ---- | ----- | ----------- | ------- |
| `?destination_account` | string | The destination account that any returned path should use | `GAEDTJ4PPEFVW5XV2S7LUXBEHNQMX5Q2GM562RJGOQG7GVCE5H3HIB4V` |
| `?destination_asset_type` | string | The type of the destination asset | `credit_alphanum4` |
| `?destination_asset_code` | string | The destination asset code, if destination_asset_type is not "native" | `USD` |
| `?destination_asset_issuer` | string | The issuer for the destination, if destination_asset_type is not "native" | `GAEDTJ4PPEFVW5XV2S7LUXBEHNQMX5Q2GM562RJGOQG7GVCE5H3HIB4V` |
| `?destination_amount` | string | The amount, denominated in the destination asset, that any returned path should be able to satisfy | `10.1` |
| `?source_account` | string | The sender's account id. Any returned path must use a source that the sender can hold | `GARSFJNXJIHO6ULUBK3DBYKVSIZE7SC72S5DYBCHU7DKL22UXKVD7MXP` |
| `?source_assets` | string | A comma separated list of assets. Any returned path must use a source included in this list  | `USD:GAEDTJ4PPEFVW5XV2S7LUXBEHNQMX5Q2GM562RJGOQG7GVCE5H3HIB4V,native` |

The endpoint will not allow requests which provide both a `source_account` and a `source_assets` parameter. All requests must provide one or the other.
The assets in `source_assets` are expected to be encoded using the following format:

The native asset should be represented as `"native"`. Issued assets should be represented as `"Code:IssuerAccountID"`. `"Code"` must consist of alphanumeric ASCII characters.


### curl Example Request

```sh
curl "https://aurora-testnet.diamcircle.org/paths?destination_account=GAEDTJ4PPEFVW5XV2S7LUXBEHNQMX5Q2GM562RJGOQG7GVCE5H3HIB4V&source_account=GARSFJNXJIHO6ULUBK3DBYKVSIZE7SC72S5DYBCHU7DKL22UXKVD7MXP&destination_asset_type=native&destination_amount=20"
```

### JavaScript Example Request

```javascript
var DiamcircleSdk = require('diamcircle-sdk');
var server = new DiamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

var source_account = "GARSFJNXJIHO6ULUBK3DBYKVSIZE7SC72S5DYBCHU7DKL22UXKVD7MXP";
var destination_account = "GAEDTJ4PPEFVW5XV2S7LUXBEHNQMX5Q2GM562RJGOQG7GVCE5H3HIB4V";
var destination_asset = DiamcircleSdk.Asset.native();
var destination_amount = "20";

server.paths(source_account, destination_account, destination_asset, destination_amount)
  .call()
  .then(function (pathResult) {
    console.log(pathResult.records);
  })
  .catch(function (err) {
    console.log(err)
  })
```

## Response

This endpoint responds with a page of path resources.  See [path resource](../resources/path.md) for reference.

### Example Response

```json
{
  "_embedded": {
    "records": [
      {
        "source_asset_type": "credit_alphanum4",
        "source_asset_code": "FOO",
        "source_asset_issuer": "GAGLYFZJMN5HEULSTH5CIGPOPAVUYPG5YSWIYDJMAPIECYEBPM2TA3QR",
        "source_amount": "100.0000000",
        "destination_asset_type": "credit_alphanum4",
        "destination_asset_code": "FOO",
        "destination_asset_issuer": "GAGLYFZJMN5HEULSTH5CIGPOPAVUYPG5YSWIYDJMAPIECYEBPM2TA3QR",
        "destination_amount": "100.0000000",
        "path": []
      }
    ]
  }
}
```

## Possible Errors

- The [standard errors](../errors.md#Standard-Errors).
- [not_found](../errors/not-found.md): A `not_found` error will be returned if no paths could be found to fulfill this payment request
