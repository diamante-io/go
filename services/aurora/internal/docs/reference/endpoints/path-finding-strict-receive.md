---
title: Strict Receive Payment Paths
clientData:
  laboratoryUrl: https://www.diamcircle.org/laboratory/#explorer?resource=paths&endpoint=all
replacement: https://developers.diamcircle.org/api/aggregations/paths/strict-receive/
---

The Diamcircle Network allows payments to be made across assets through _path payments_.  A path
payment specifies a series of assets to route a payment through, from source asset (the asset
debited from the payer) to destination asset (the asset credited to the payee).

A [Path Payment Strict Receive](../../../guides/concepts/list-of-operations.html#path-payment-strict-receive) allows a user to specify the *amount of the asset received*. The amount sent varies based on offers in the order books.  If you would like to search for a path specifying the amount to be sent, use the [Find Payment Paths (Strict Send)](./path-finding-strict-send.html).

A strict receive path search is specified using:

- The source account id or source assets.
- The asset and amount that the destination account should receive.

As part of the search, aurora will load a list of assets available to the source account id and
will find any payment paths from those source assets to the desired destination asset. The search's
amount parameter will be used to determine if a given path can satisfy a payment of the
desired amount.

## Request

```
GET /paths/strict-receive?source_account={sa}&destination_asset_type={at}&destination_asset_code={ac}&destination_asset_issuer={di}&destination_amount={amount}&destination_account={da}
```

## Arguments

| name | notes | description | example |
| ---- | ----- | ----------- | ------- |
| `?source_account` | string | The sender's account id. Any returned path must use an asset that the sender has a trustline to. | `GARSFJNXJIHO6ULUBK3DBYKVSIZE7SC72S5DYBCHU7DKL22UXKVD7MXP` |
| `?source_assets` | string | A comma separated list of assets. Any returned path must use an asset included in this list | `USD:GAEDTJ4PPEFVW5XV2S7LUXBEHNQMX5Q2GM562RJGOQG7GVCE5H3HIB4V,native` |
| `?destination_account` | string | The destination account that any returned path should use | `GAEDTJ4PPEFVW5XV2S7LUXBEHNQMX5Q2GM562RJGOQG7GVCE5H3HIB4V` |
| `?destination_asset_type` | string | The type of the destination asset | `credit_alphanum4` |
| `?destination_asset_code` | required if `destination_asset_type` is not `native`, string | The destination asset code, if destination_asset_type is not "native" | `USD` |
| `?destination_asset_issuer` | required if `destination_asset_type` is not `native`, string | The issuer for the destination asset, if destination_asset_type is not "native" | `GAEDTJ4PPEFVW5XV2S7LUXBEHNQMX5Q2GM562RJGOQG7GVCE5H3HIB4V` |
| `?destination_amount` | string | The amount, denominated in the destination asset, that any returned path should be able to satisfy | `10.1` |

The endpoint will not allow requests which provide both a `source_account` and a `source_assets` parameter. All requests must provide one or the other.
The assets in `source_assets` are expected to be encoded using the following format:

XLM should be represented as `"native"`. Issued assets should be represented as `"Code:IssuerAccountID"`. `"Code"` must consist of alphanumeric ASCII characters.


### curl Example Request

```sh
curl "https://aurora-testnet.diamcircle.org/paths/strict-receive?destination_account=GAEDTJ4PPEFVW5XV2S7LUXBEHNQMX5Q2GM562RJGOQG7GVCE5H3HIB4V&source_account=GARSFJNXJIHO6ULUBK3DBYKVSIZE7SC72S5DYBCHU7DKL22UXKVD7MXP&destination_asset_type=native&destination_amount=20"
```

### JavaScript Example Request

```javascript
var DiamcircleSdk = require('diamcircle-sdk');
var server = new DiamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

var sourceAccount = "GARSFJNXJIHO6ULUBK3DBYKVSIZE7SC72S5DYBCHU7DKL22UXKVD7MXP";
var destinationAsset = DiamcircleSdk.Asset.native();
var destinationAmount = "20";

server.paths(sourceAccount, destinationAsset, destinationAmount)
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
