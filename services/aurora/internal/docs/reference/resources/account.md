---
title: Account
replacement: https://developers.diamcircle.org/api/resources/accounts/
---

In the diamcircle network, users interact using **accounts** which can be controlled by a corresponding keypair that can authorize transactions. One can create a new account with the [Create Account](./operation.md#create-account) operation.

To learn more about the concept of accounts in the diamcircle network, take a look at the [diamcircle account concept guide](https://www.diamcircle.org/developers/learn/concepts/accounts.html).

When aurora returns information about an account it uses the following format:

## Attributes
| Attribute      | Type             | Description                                                                                                                                  |
|----------------|------------------|------------------------------------------------------------------------------------------------------------------------                      |
| id             | string           | The canonical id of this account, suitable for use as the :id parameter for url templates that require an account's ID.                      |
| account_id     | string           | The account's public key encoded into a base32 string representation.                                                                        |
| sequence       | number           | The current sequence number that can be used when submitting a transaction from this account.                                                |
| subentry_count | number           | The number of [account subentries](https://www.diamcircle.org/developers/guides/concepts/ledger.html#ledger-entries).                           |
| balances       | array of objects | An array of the native asset or credits this account holds.                                                                                  |
| thresholds     | object           | An object of account thresholds.                                                                                                             |
| flags          | object           | The flags denote the enabling/disabling of certain asset issuer privileges.                                                                  |
| signers        | array of objects | An array of [account signers](https://www.diamcircle.org/developers/guides/concepts/multi-sig.html#additional-signing-keys) with their weights. |
| data           | object           | An array of account [data](./data.md) fields.                                                                                                |

### Signer Object
| Attribute  | Type   | Description                                                                                                      |
|------------|--------|------------------------------------------------------------------------------------------------------------------|
| weight     | number | The numerical weight of a signer, necessary to determine whether a transaction meets the threshold requirements. |
| key        | string | Different depending on the type of the signer.                                                                   |
| type       | string | See below.                                                                                                       |

### Possible Signer Types
| Type               | Description                                                                                                                                                                                                         |
|--------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| ed25519_public_key | A normal diamcircle public key.                                                                                                                                                                                        |
| sha256_hash        | The SHA256 hash of some arbitrary `x`. Adding a signature of this type allows anyone who knows `x` to sign a transaction from this account. *Note: Once this transaction is broadcast, `x` will be known publicly.* |
| preauth_tx         | The hash of a pre-authorized transaction. This signer is automatically removed from the account when a matching transaction is properly applied.                                                                    |
### Balances Object
| Attribute           | Type             |                                                                                                                                                                                                                                                               |
|---------------------|------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| balance             | string           | How much of an asset is owned.                                                                                                                                                                                                                                |
| buying_liabilities  | string           | The total amount of an asset offered to buy aggregated over all offers owned by this account.                                                                                                                                                                 |
| selling_liabilities | string           | The total amount of an asset offered to sell aggregated over all offers owned by this account.                                                                                                                                                                |
| limit               | optional, number | The maximum amount of an asset that this account is willing to accept (this is specified when an account opens a trustline).                                                                                                                                  |
| asset_type          | string           | Either native, credit_alphanum4, or credit_alphanum12.                                                                                                                                                                                                        |
| asset_code          | optional, string | The code for the asset.                                                                                                                                                                                                                                       |
| asset_issuer        | optional, string | The diamcircle address of the given asset's issuer.                                                                                                                                                                                                              |
| is_authorized       | optional, bool   | The trustline status for an `auth_required` asset.  If true, the issuer of the asset has granted the account permission to send, receive, buy, or sell the asset.  If false, the issuer has not, so the account cannot send, receive, buy, or sell the asset. |

### Flag Object
| Attribute      | Type |                                                                                                                                |
|----------------|------|--------------------------------------------------------------------------------------------------------------------------------|
| auth_immutable | bool | With this setting, none of the following authorization flags can be changed.                                                   |
| auth_required  | bool | With this setting, an anchor must approve anyone who wants to hold its asset.                                                  |
| auth_revocable | bool | With this setting, an anchor can set the authorize flag of an existing trustline to freeze the assets held by an asset holder. |

### Threshold Object
| Attribute      | Type   |                                                                                                                                                                                                                                                                                                                                                                           |
|----------------|--------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| low_threshold  | number | The weight required for a valid transaction including the [Allow Trust][allow_trust] and [Bump Sequence][bump_seq] operations.                                                                                                                                                                                                                                            |
| med_threshold  | number | The weight required for a valid transaction including the [Create Account][create_acc], [Payment][payment], [Path Payment Strict Send][path_payment_send], [Path Payment Strict Receive][path_payment_receive], [Manage Buy Offer][manage_buy_offer], [Manage Sell Offer][manage_sell_offer], [Create Passive Sell Offer][passive_sell_offer], [Change Trust][change_trust], [Inflation][inflation], and [Manage Data][manage_data] operations. |
| high_threshold | number | The weight required for a valid transaction including the [Account Merge][account_merge] and [Set Options]() operations.                                                                                                                                                                                                                                                  |

[account_merge]: https://www.diamcircle.org/developers/guides/concepts/list-of-operations.html#account-merge
[allow_trust]: https://www.diamcircle.org/developers/guides/concepts/list-of-operations.html#allow-trust
[bump_seq]: https://www.diamcircle.org/developers/guides/concepts/list-of-operations.html#bump-sequence
[change_trust]: https://www.diamcircle.org/developers/guides/concepts/list-of-operations.html#change-trust
[create_acc]: https://www.diamcircle.org/developers/guides/concepts/list-of-operations.html#create-account
[inflation]: https://www.diamcircle.org/developers/guides/concepts/list-of-operations.html#inflation
[manage_data]: https://www.diamcircle.org/developers/guides/concepts/list-of-operations.html#manage-data
[manage_buy_offer]: https://www.diamcircle.org/developers/guides/concepts/list-of-operations.html#manage-buy-offer
[manage_sell_offer]: https://www.diamcircle.org/developers/guides/concepts/list-of-operations.html#manage-sell-offer
[passive_sell_offer]: https://www.diamcircle.org/developers/guides/concepts/list-of-operations.html#create-passive-sell-offer
[path_payment_receive]: https://www.diamcircle.org/developers/guides/concepts/list-of-operations.html#path-payment-strict-receive
[path_payment_send]: https://www.diamcircle.org/developers/guides/concepts/list-of-operations.html#path-payment-strict-send
[payment]: https://www.diamcircle.org/developers/guides/concepts/list-of-operations.html#payment
[set_options]: https://www.diamcircle.org/developers/guides/concepts/list-of-operations.html#set-options

## Links
| rel          | Example                                                                                                 | Description                                                  | `templated` |
|--------------|---------------------------------------------------------------------------------------------------------|--------------------------------------------------------------|-------------|
| data         | `/accounts/GAOEWNUEKXKNGB2AAOX6S6FEP6QKCFTU7KJH647XTXQXTMOAUATX2VF5/data/{key}`                         | [Data fields](./data.md) related to this account             | true        |
| effects      | `/accounts/GAOEWNUEKXKNGB2AAOX6S6FEP6QKCFTU7KJH647XTXQXTMOAUATX2VF5/effects/{?cursor,limit,order}`      | The [effects](./effect.md) related to this account           | true        |
| offers       | `/accounts/GAOEWNUEKXKNGB2AAOX6S6FEP6QKCFTU7KJH647XTXQXTMOAUATX2VF5/offers/{?cursor,limit,order}`       | The [offers](./offer.md) related to this account             | true        |
| operations   | `/accounts/GAOEWNUEKXKNGB2AAOX6S6FEP6QKCFTU7KJH647XTXQXTMOAUATX2VF5/operations/{?cursor,limit,order}`   | The [operations](./operation.md) related to this account     | true        |
| payments     | `/accounts/GAOEWNUEKXKNGB2AAOX6S6FEP6QKCFTU7KJH647XTXQXTMOAUATX2VF5/payments/{?cursor,limit,order}`     | The [payments](./payment.md) related to this account         | true        |
| trades       | `/accounts/GAOEWNUEKXKNGB2AAOX6S6FEP6QKCFTU7KJH647XTXQXTMOAUATX2VF5/trades/{?cursor,limit,order}`       | The [trades](./trade.md) related to this account             | true        |
| transactions | `/accounts/GAOEWNUEKXKNGB2AAOX6S6FEP6QKCFTU7KJH647XTXQXTMOAUATX2VF5/transactions/{?cursor,limit,order}` | The [transactions](./transaction.md) related to this account | true        |

## Example

```json
{
  "_links": {
    "self": {
      "href": "https://aurora-testnet.diamcircle.org/accounts/GBWRID7MPYUDBTNQPEHUN4XOBVVDPJOHYXAVW3UTOD2RG7BDAY6O3PHW"
    },
    "transactions": {
      "href": "https://aurora-testnet.diamcircle.org/accounts/GBWRID7MPYUDBTNQPEHUN4XOBVVDPJOHYXAVW3UTOD2RG7BDAY6O3PHW/transactions{?cursor,limit,order}",
      "templated": true
    },
    "operations": {
      "href": "https://aurora-testnet.diamcircle.org/accounts/GBWRID7MPYUDBTNQPEHUN4XOBVVDPJOHYXAVW3UTOD2RG7BDAY6O3PHW/operations{?cursor,limit,order}",
      "templated": true
    },
    "payments": {
      "href": "https://aurora-testnet.diamcircle.org/accounts/GBWRID7MPYUDBTNQPEHUN4XOBVVDPJOHYXAVW3UTOD2RG7BDAY6O3PHW/payments{?cursor,limit,order}",
      "templated": true
    },
    "effects": {
      "href": "https://aurora-testnet.diamcircle.org/accounts/GBWRID7MPYUDBTNQPEHUN4XOBVVDPJOHYXAVW3UTOD2RG7BDAY6O3PHW/effects{?cursor,limit,order}",
      "templated": true
    },
    "offers": {
      "href": "https://aurora-testnet.diamcircle.org/accounts/GBWRID7MPYUDBTNQPEHUN4XOBVVDPJOHYXAVW3UTOD2RG7BDAY6O3PHW/offers{?cursor,limit,order}",
      "templated": true
    },
    "trades": {
      "href": "https://aurora-testnet.diamcircle.org/accounts/GBWRID7MPYUDBTNQPEHUN4XOBVVDPJOHYXAVW3UTOD2RG7BDAY6O3PHW/trades{?cursor,limit,order}",
      "templated": true
    },
    "data": {
      "href": "https://aurora-testnet.diamcircle.org/accounts/GBWRID7MPYUDBTNQPEHUN4XOBVVDPJOHYXAVW3UTOD2RG7BDAY6O3PHW/data/{key}",
      "templated": true
    }
  },
  "id": "GBWRID7MPYUDBTNQPEHUN4XOBVVDPJOHYXAVW3UTOD2RG7BDAY6O3PHW",
  "paging_token": "",
  "account_id": "GBWRID7MPYUDBTNQPEHUN4XOBVVDPJOHYXAVW3UTOD2RG7BDAY6O3PHW",
  "sequence": "43692723777044483",
  "subentry_count": 3,
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
      "public_key": "XCPNCUKYDHPMMH6TMHK73K5VP5A6ZTQ2L7Q74JR3TDANNFB3TMRS5OKG",
      "weight": 1,
      "key": "XCPNCUKYDHPMMH6TMHK73K5VP5A6ZTQ2L7Q74JR3TDANNFB3TMRS5OKG",
      "type": "sha256_hash"
    },
    {
      "public_key": "TABGGIW6EXOVOSNJ2O27U2DUX7RWHSRBGOKQLGYDTOXPANEX6LXBX7O7",
      "weight": 1,
      "key": "TABGGIW6EXOVOSNJ2O27U2DUX7RWHSRBGOKQLGYDTOXPANEX6LXBX7O7",
      "type": "preauth_tx"
    },
    {
      "public_key": "GBWRID7MPYUDBTNQPEHUN4XOBVVDPJOHYXAVW3UTOD2RG7BDAY6O3PHW",
      "weight": 1,
      "key": "GBWRID7MPYUDBTNQPEHUN4XOBVVDPJOHYXAVW3UTOD2RG7BDAY6O3PHW",
      "type": "ed25519_public_key"
    }
  ],
  "data": {}
}
```

## Endpoints
| Resource                                                         | Type       | Resource URI Template                |
|------------------------------------------------------------------|------------|--------------------------------------|
| [Account Details](../endpoints/accounts-single.md)               | Single     | `/accounts/:id`                      |
| [Account Data](../endpoints/data-for-account.md)                 | Single     | `/accounts/:id/data/:key`            |
| [Account Transactions](../endpoints/transactions-for-account.md) | Collection | `/accounts/:account_id/transactions` |
| [Account Operations](../endpoints/operations-for-account.md)     | Collection | `/accounts/:account_id/operations`   |
| [Account Payments](../endpoints/payments-for-account.md)         | Collection | `/accounts/:account_id/payments`     |
| [Account Effects](../endpoints/effects-for-account.md)           | Collection | `/accounts/:account_id/effects`      |
| [Account Offers](../endpoints/offers-for-account.md)             | Collection | `/accounts/:account_id/offers`       |
