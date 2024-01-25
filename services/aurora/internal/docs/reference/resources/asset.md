---
title: Asset
replacement: https://developers.diamcircle.org/api/resources/assets/
---

**Assets** are the units that are traded on the Diamcircle Network.

An asset consists of an type, code, and issuer.

To learn more about the concept of assets in the Diamcircle network, take a look at the [Diamcircle assets concept guide](https://www.diamcircle.org/developers/guides/concepts/assets.html).

## Attributes

|    Attribute     |  Type  |                                                                                                                                |
| ---------------- | ------ | ------------------------------------------------------------------------------------------------------------------------------ |
| asset_type               | string | The type of this asset: "credit_alphanum4", or "credit_alphanum12". |
| asset_code               | string | The code of this asset.   |
| asset_issuer             | string | The issuer of this asset. |
| accounts                 | object | The number of accounts and claimable balances holding this asset. Accounts are summarized by each state of the trust line flags. |
| balances                 | object | The number of units of credit issued, summarized by each state of the trust line flags, or if they are in a claimable balance. |
| flags                    | object | The flags denote the enabling/disabling of certain asset issuer privileges. |
| paging_token             | string | A [paging token](./page.md) suitable for use as the `cursor` parameter to transaction collection resources.                   |

#### Flag Object
|    Attribute     |  Type  |                                                                                                                                |
| ---------------- | ------ | ------------------------------------------------------------------------------------------------------------------------------ |
| auth_immutable             | bool | With this setting, none of the following authorization flags can be changed. |
| auth_required              | bool | With this setting, an anchor must approve anyone who wants to hold its asset.  |
| auth_revocable             | bool | With this setting, an anchor can set the authorize flag of an existing trustline to freeze the assets held by an asset holder.  |

## Links
| rel          | Example                                                                                           | Description                                                
|--------------|---------------------------------------------------------------------------------------------------|------------------------------------------------------------
| toml  | `https://www.diamcircle.org/.well-known/diamcircle.toml`| Link to the TOML file for this issuer |

## Example

```json
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
    "authorized": "100.0000000",
    "authorized_to_maintain_liabilities": "50.0000000",
    "unauthorized": "25.0000000",
    "claimable_balances": "12.5000000"
  },
  "flags": {
    "auth_required": false,
    "auth_revocable": false
  }
}
```

## Endpoints

|  Resource                                |    Type    |    Resource URI Template     |
| ---------------------------------------- | ---------- | ---------------------------- |
| [All Assets](../endpoints/assets-all.md) | Collection | `/assets` (`GET`)            |
