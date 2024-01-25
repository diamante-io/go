---
title: Payment Path
replacement: https://developers.diamcircle.org/api/aggregations/paths/
---

A **path** resource contains information about a payment path.  A path can be used by code to populate necessary fields on path payment operation, such as `path` and `sendMax`.


## Attributes
| Attribute                | Type             |                                                                                                                                |
|--------------------------|------------------|--------------------------------------------------------------------------------------------------------------------------------|
| path                     | array of objects            | An array of assets that represents the intermediary assets this path hops through                                               |
| source_amount            | string           | An estimated cost for making a payment of destination_amount on this path. Suitable for use in a path payments `sendMax` field |
| destination_amount       | string           | The destination amount specified in the search that found this path                                                            |
| destination_asset_type   | string           | The type for the destination asset specified in the search that found this path                                                |
| destination_asset_code   | optional, string | The code for the destination asset specified in the search that found this path                                                |
| destination_asset_issuer | optional, string | The issuer for the destination asset specified in the search that found this path                                              |
| source_asset_type        | string           | The type for the source asset specified in the search that found this path                                                     |
| source_asset_code        | optional, string | The code for the source asset specified in the search that found this path                                                     |
| source_asset_issuer      | optional, string | The issuer for the source asset specified in the search that found this path                                                   |

#### Asset Object
| Attribute    | Type             |                                                                                                                        |
|--------------|------------------|------------------------------------------------------------------------------------------------------------------------
| asset_code     | optional, string           | The code for the asset.                       |
| asset_type     | string           | Either native, credit_alphanum4, or credit_alphanum12.                        |
| asset_issuer     | optional, string           | The diamcircle address of the given asset's issuer.  |

## Example

```json
{
		"destination_amount": "20.0000000",
		"destination_asset_code": "EUR",
		"destination_asset_issuer": "GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN",
		"destination_asset_type": "credit_alphanum4",
		"path": [
				{
						"asset_code": "1",
						"asset_issuer": "GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN",
						"asset_type": "credit_alphanum4"
				}
		],
		"source_amount": "20.0000000",
		"source_asset_code": "USD",
		"source_asset_issuer": "GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN",
		"source_asset_type": "credit_alphanum4"
}
```

## Endpoints
| Resource                                 | Type       | Resource URI Template |
|------------------------------------------|------------|-----------------------|
| [Find Payment Paths](../endpoints/path-finding.md) | Collection | `/paths`              |
