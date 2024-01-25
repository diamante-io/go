---
title: Strict Send Payment Paths
clientData:
  laboratoryUrl: https://www.diamcircle.org/laboratory/#explorer?resource=paths&endpoint=all
replacement: https://developers.diamcircle.org/api/aggregations/paths/strict-send/
---

The diamcircle Network allows payments to be made across assets through _path payments_.  A path
payment specifies a series of assets to route a payment through, from source asset (the asset
debited from the payer) to destination asset (the asset credited to the payee).

A [Path Payment Strict Send](../../../guides/concepts/list-of-operations.html#path-payment-strict-send) allows a user to specify the amount of the asset to send. The amount received will vary based on offers in the order books.


A path payment strict send search is specified using:

- The destination account id or destination assets.
- The source asset.
- The source amount.

As part of the search, aurora will load a list of assets available to the source account id or use the assets passed in the request and will find any payment paths from those source assets to the desired destination asset. The source's amount parameter will be used to determine if a given path can satisfy a payment of the desired amount.

## Request

```
https://aurora-testnet.diamcircle.org/paths/strict-send?&source_amount={sa}&source_asset_type={at}&source_asset_code={ac}&source_asset_issuer={ai}&destination_account={da}
```

## Arguments

| name | notes | description | example |
| ---- | ----- | ----------- | ------- |
| `?source_amount` | string | The amount, denominated in the source asset, that any returned path should be able to satisfy | `10.1` |
| `?source_asset_type` | string | The type of the source asset | `credit_alphanum4` |
| `?source_asset_code` | string, required if `source_asset_type` is not `native`, string | The source asset code, if source_asset_type is not "native" | `USD` |
| `?source_asset_issuer` | string, required if `source_asset_type` is not `native`, string | The issuer for the source asset, if source_asset_type is not "native" | `GAEDTJ4PPEFVW5XV2S7LUXBEHNQMX5Q2GM562RJGOQG7GVCE5H3HIB4V` |
| `?destination_account` | string optional | The destination account that any returned path should use | `GAEDTJ4PPEFVW5XV2S7LUXBEHNQMX5Q2GM562RJGOQG7GVCE5H3HIB4V` |
| `?destination_assets` | string optional | A comma separated list of assets. Any returned path must use an asset included in this list  | `USD:GAEDTJ4PPEFVW5XV2S7LUXBEHNQMX5Q2GM562RJGOQG7GVCE5H3HIB4V,native` |

The endpoint will not allow requests which provide both a `destination_account` and `destination_assets` parameter. All requests must provide one or the other.
The assets in `destination_assets` are expected to be encoded using the following format:

XLM should be represented as `"native"`. Issued assets should be represented as `"Code:IssuerAccountID"`. `"Code"` must consist of alphanumeric ASCII characters.


### curl Example Request

```sh
curl "https://aurora-testnet.diamcircle.org/paths/strict-send?&source_amount=10&source_asset_type=native&destination_assets=MXN:GC2GFGZ5CZCFCDJSQF3YYEAYBOS3ZREXJSPU7LUJ7JU3LP3BQNHY7YKS"
```

### JavaScript Example Request

```javascript
var diamcircleSdk = require('diamcircle-sdk');
var server = new diamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

var sourceAsset = diamcircleSdk.Asset.native();
var sourceAmount = "20";
var destinationAsset = new diamcircleSdk.Asset(
  'USD',
  'GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN'
)


server.strictSendPaths(sourceAsset, sourceAmount, [destinationAsset])
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
