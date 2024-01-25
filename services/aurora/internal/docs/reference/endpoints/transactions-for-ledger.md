---
title: Transactions for Ledger
clientData:
  laboratoryUrl: https://www.diamcircle.org/laboratory/#explorer?resource=transactions&endpoint=for_ledger
replacement: https://developers.diamcircle.org/api/resources/ledgers/transactions/
---

This endpoint represents successful [transactions](../resources/transaction.md) in a given [ledger](../resources/ledger.md).

## Request

```
GET /ledgers/{id}/transactions{?cursor,limit,order,include_failed}
```

### Arguments

|  name  |  notes  | description | example |
| ------ | ------- | ----------- | ------- |
| `id` | required, number | Ledger ID | `697121` |
| `?cursor` | optional, default _null_ | A paging token, specifying where to start returning records from. | `12884905984` |
| `?order`  | optional, string, default `asc` | The order in which to return rows, "asc" or "desc". | `asc` |
| `?limit`  | optional, number, default `10` | Maximum number of records to return. | `200` |
| `?include_failed` | optional, bool, default: `false` | Set to `true` to include failed transactions in results. | `true` |

### curl Example Request

```sh
curl "https://aurora-testnet.diamcircle.org/ledgers/697121/transactions?limit=1"
```

### JavaScript Example Request

```javascript
var DiamcircleSdk = require('diamcircle-sdk');
var server = new DiamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

server.transactions()
  .forLedger("697121")
  .limit("1")
  .call()
  .then(function (accountResults) {
    console.log(accountResults.records)
  })
  .catch(function (err) {
    console.log(err)
  })
```

## Response

This endpoint responds with a list of transactions in a given ledger. See [transaction
resource](../resources/transaction.md) for reference.

### Example Response

```json
{
  "_links": {
    "self": {
      "href": "https://aurora-testnet.diamcircle.org/ledgers/697121/transactions?cursor=&limit=10&order=asc"
    },
    "next": {
      "href": "https://aurora-testnet.diamcircle.org/ledgers/697121/transactions?cursor=2994111896367104&limit=10&order=asc"
    },
    "prev": {
      "href": "https://aurora-testnet.diamcircle.org/ledgers/697121/transactions?cursor=2994111896358912&limit=10&order=desc"
    }
  },
  "_embedded": {
    "records": [
      {
        "_links": {
          "self": {
            "href": "https://aurora-testnet.diamcircle.org/transactions/264226cb06af3b86299031884175155e67a02e0a8ad0b3ab3a88b409a8c09d5c"
          },
          "account": {
            "href": "https://aurora-testnet.diamcircle.org/accounts/GAIH3ULLFQ4DGSECF2AR555KZ4KNDGEKN4AFI4SU2M7B43MGK3QJZNSR"
          },
          "ledger": {
            "href": "https://aurora-testnet.diamcircle.org/ledgers/697121"
          },
          "operations": {
            "href": "https://aurora-testnet.diamcircle.org/transactions/264226cb06af3b86299031884175155e67a02e0a8ad0b3ab3a88b409a8c09d5c/operations{?cursor,limit,order}",
            "templated": true
          },
          "effects": {
            "href": "https://aurora-testnet.diamcircle.org/transactions/264226cb06af3b86299031884175155e67a02e0a8ad0b3ab3a88b409a8c09d5c/effects{?cursor,limit,order}",
            "templated": true
          },
          "precedes": {
            "href": "https://aurora-testnet.diamcircle.org/transactions?order=asc&cursor=2994111896358912"
          },
          "succeeds": {
            "href": "https://aurora-testnet.diamcircle.org/transactions?order=desc&cursor=2994111896358912"
          }
        },
        "id": "264226cb06af3b86299031884175155e67a02e0a8ad0b3ab3a88b409a8c09d5c",
        "paging_token": "2994111896358912",
        "successful": true,
        "hash": "264226cb06af3b86299031884175155e67a02e0a8ad0b3ab3a88b409a8c09d5c",
        "ledger": 697121,
        "created_at": "2019-04-09T20:14:25Z",
        "source_account": "GAIH3ULLFQ4DGSECF2AR555KZ4KNDGEKN4AFI4SU2M7B43MGK3QJZNSR",
        "source_account_sequence": "4660039994869",
        "fee_charged": 100,
        "max_fee": 100,
        "operation_count": 1,
        "envelope_xdr": "AAAAABB90WssODNIgi6BHveqzxTRmIpvAFRyVNM+Hm2GVuCcAAAAZAAABD0AB031AAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAFIMRkFZ9gZifhRSlklQpsz/9P04Earv0dzS3MkIM1cYAAAAXSHboAAAAAAAAAAABhlbgnAAAAEA+biIjrDy8yi+SvhFElIdWGBRYlDscnSSHkPchePy2JYDJn4wvJYDBumXI7/NmttUey3+cGWbBFfnnWh1H5EoD",
        "result_xdr": "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=",
        "result_meta_xdr": "AAAAAQAAAAIAAAADAAqjIQAAAAAAAAAAEH3Rayw4M0iCLoEe96rPFNGYim8AVHJU0z4ebYZW4JwBOLmYhGq/IAAABD0AB030AAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAqjIQAAAAAAAAAAEH3Rayw4M0iCLoEe96rPFNGYim8AVHJU0z4ebYZW4JwBOLmYhGq/IAAABD0AB031AAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMACqMhAAAAAAAAAAAQfdFrLDgzSIIugR73qs8U0ZiKbwBUclTTPh5thlbgnAE4uZiEar8gAAAEPQAHTfUAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEACqMhAAAAAAAAAAAQfdFrLDgzSIIugR73qs8U0ZiKbwBUclTTPh5thlbgnAE4uYE789cgAAAEPQAHTfUAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAACqMhAAAAAAAAAAAUgxGQVn2BmJ+FFKWSVCmzP/0/TgRqu/R3NLcyQgzVxgAAABdIdugAAAqjIQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
        "fee_meta_xdr": "AAAAAgAAAAMACqMgAAAAAAAAAAAQfdFrLDgzSIIugR73qs8U0ZiKbwBUclTTPh5thlbgnAE4uZiEar+EAAAEPQAHTfQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEACqMhAAAAAAAAAAAQfdFrLDgzSIIugR73qs8U0ZiKbwBUclTTPh5thlbgnAE4uZiEar8gAAAEPQAHTfQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
        "memo_type": "none",
        "signatures": [
          "Pm4iI6w8vMovkr4RRJSHVhgUWJQ7HJ0kh5D3IXj8tiWAyZ+MLyWAwbplyO/zZrbVHst/nBlmwRX551odR+RKAw=="
        ]
      },
      {
        "memo": "2A1V6J5703G47XHY",
        "_links": {
          "self": {
            "href": "https://aurora-testnet.diamcircle.org/transactions/f175108e5c64619705b112a99fa32884dfa0511d9a8986aade87905b08eabe5b"
          },
          "account": {
            "href": "https://aurora-testnet.diamcircle.org/accounts/GAZ4A54KE6MTMXYEPM7T3IDLZWGNCCKB5ME422NZ3MAMTHWWP37RPEBW"
          },
          "ledger": {
            "href": "https://aurora-testnet.diamcircle.org/ledgers/697121"
          },
          "operations": {
            "href": "https://aurora-testnet.diamcircle.org/transactions/f175108e5c64619705b112a99fa32884dfa0511d9a8986aade87905b08eabe5b/operations{?cursor,limit,order}",
            "templated": true
          },
          "effects": {
            "href": "https://aurora-testnet.diamcircle.org/transactions/f175108e5c64619705b112a99fa32884dfa0511d9a8986aade87905b08eabe5b/effects{?cursor,limit,order}",
            "templated": true
          },
          "precedes": {
            "href": "https://aurora-testnet.diamcircle.org/transactions?order=asc&cursor=2994111896363008"
          },
          "succeeds": {
            "href": "https://aurora-testnet.diamcircle.org/transactions?order=desc&cursor=2994111896363008"
          }
        },
        "id": "f175108e5c64619705b112a99fa32884dfa0511d9a8986aade87905b08eabe5b",
        "paging_token": "2994111896363008",
        "successful": true,
        "hash": "f175108e5c64619705b112a99fa32884dfa0511d9a8986aade87905b08eabe5b",
        "ledger": 697121,
        "created_at": "2019-04-09T20:14:25Z",
        "source_account": "GAZ4A54KE6MTMXYEPM7T3IDLZWGNCCKB5ME422NZ3MAMTHWWP37RPEBW",
        "source_account_sequence": "2994107601387521",
        "fee_charged": 100,
        "max_fee": 100,
        "operation_count": 1,
        "envelope_xdr": "AAAAADPAd4onmTZfBHs/PaBrzYzRCUHrCc1pudsAyZ7Wfv8XAAAAZAAKoyAAAAABAAAAAAAAAAEAAAAQMkExVjZKNTcwM0c0N1hIWQAAAAEAAAABAAAAADPAd4onmTZfBHs/PaBrzYzRCUHrCc1pudsAyZ7Wfv8XAAAAAQAAAADMSEvcRKXsaUNna++Hy7gWm/CfqTjEA7xoGypfrFGUHAAAAAAAAAAAhFKDAAAAAAAAAAAB1n7/FwAAAEBJdXuYg13Glzx1RinVCXd/cc1usrhU/0f5HFZ7lyIR8kS3T6PRrW78TQDNqXz+ukUiPwlB1A8MqxoW/SAL5FIB",
        "result_xdr": "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=",
        "result_meta_xdr": "AAAAAQAAAAIAAAADAAqjIQAAAAAAAAAAM8B3iieZNl8Eez89oGvNjNEJQesJzWm52wDJntZ+/xcAAAAXSHbnnAAKoyAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAqjIQAAAAAAAAAAM8B3iieZNl8Eez89oGvNjNEJQesJzWm52wDJntZ+/xcAAAAXSHbnnAAKoyAAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMACqMgAAAAAAAAAADMSEvcRKXsaUNna++Hy7gWm/CfqTjEA7xoGypfrFGUHAAAAANViducAABeBgAAoRQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEACqMhAAAAAAAAAADMSEvcRKXsaUNna++Hy7gWm/CfqTjEA7xoGypfrFGUHAAAAAPZ3F6cAABeBgAAoRQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMACqMhAAAAAAAAAAAzwHeKJ5k2XwR7Pz2ga82M0QlB6wnNabnbAMme1n7/FwAAABdIduecAAqjIAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEACqMhAAAAAAAAAAAzwHeKJ5k2XwR7Pz2ga82M0QlB6wnNabnbAMme1n7/FwAAABbEJGScAAqjIAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
        "fee_meta_xdr": "AAAAAgAAAAMACqMgAAAAAAAAAAAzwHeKJ5k2XwR7Pz2ga82M0QlB6wnNabnbAMme1n7/FwAAABdIdugAAAqjIAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEACqMhAAAAAAAAAAAzwHeKJ5k2XwR7Pz2ga82M0QlB6wnNabnbAMme1n7/FwAAABdIduecAAqjIAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
        "memo_type": "text",
        "signatures": [
          "SXV7mINdxpc8dUYp1Ql3f3HNbrK4VP9H+RxWe5ciEfJEt0+j0a1u/E0Azal8/rpFIj8JQdQPDKsaFv0gC+RSAQ=="
        ]
      },
      {
        "memo": "WHALE",
        "_links": {
          "self": {
            "href": "https://aurora-testnet.diamcircle.org/transactions/83b6ebf4b3aec5b36cab14ae0f438a23487746857903a9e0bb002564b4641e25"
          },
          "account": {
            "href": "https://aurora-testnet.diamcircle.org/accounts/GABRMXDIJCTDSMPC67J64NSAMWRSYXVCXYTXVFC73DTHBKELHNKWANXP"
          },
          "ledger": {
            "href": "https://aurora-testnet.diamcircle.org/ledgers/697121"
          },
          "operations": {
            "href": "https://aurora-testnet.diamcircle.org/transactions/83b6ebf4b3aec5b36cab14ae0f438a23487746857903a9e0bb002564b4641e25/operations{?cursor,limit,order}",
            "templated": true
          },
          "effects": {
            "href": "https://aurora-testnet.diamcircle.org/transactions/83b6ebf4b3aec5b36cab14ae0f438a23487746857903a9e0bb002564b4641e25/effects{?cursor,limit,order}",
            "templated": true
          },
          "precedes": {
            "href": "https://aurora-testnet.diamcircle.org/transactions?order=asc&cursor=2994111896367104"
          },
          "succeeds": {
            "href": "https://aurora-testnet.diamcircle.org/transactions?order=desc&cursor=2994111896367104"
          }
        },
        "id": "83b6ebf4b3aec5b36cab14ae0f438a23487746857903a9e0bb002564b4641e25",
        "paging_token": "2994111896367104",
        "successful": true,
        "hash": "83b6ebf4b3aec5b36cab14ae0f438a23487746857903a9e0bb002564b4641e25",
        "ledger": 697121,
        "created_at": "2019-04-09T20:14:25Z",
        "source_account": "GABRMXDIJCTDSMPC67J64NSAMWRSYXVCXYTXVFC73DTHBKELHNKWANXP",
        "source_account_sequence": "122518237256298",
        "fee_charged": 100,
        "max_fee": 100,
        "operation_count": 1,
        "envelope_xdr": "AAAAAAMWXGhIpjkx4vfT7jZAZaMsXqK+J3qUX9jmcKiLO1VgAAAAZAAAb24AAppqAAAAAQAAAAAAAAAAAAAAAFys/kkAAAABAAAABVdIQUxFAAAAAAAAAQAAAAAAAAAAAAAAAKrN4k6edFMb0WEyPzEEjWUAji0pvvALw+BAH4OnekA5AAAAAAcnDgAAAAAAAAAAAYs7VWAAAABAYd9uIm+TjIcAjTU90YJoNg/r+6PU3Uss7ewUb1w3yMa+HyoSvDq8sDz/SYmDBH7F+0ACIeBF4kkVEKVBJMh0AQ==",
        "result_xdr": "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=",
        "result_meta_xdr": "AAAAAQAAAAIAAAADAAqjIQAAAAAAAAAAAxZcaEimOTHi99PuNkBloyxeor4nepRf2OZwqIs7VWAAJBMYWVFGqAAAb24AApppAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAqjIQAAAAAAAAAAAxZcaEimOTHi99PuNkBloyxeor4nepRf2OZwqIs7VWAAJBMYWVFGqAAAb24AAppqAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMACqMhAAAAAAAAAAADFlxoSKY5MeL30+42QGWjLF6ivid6lF/Y5nCoiztVYAAkExhZUUaoAABvbgACmmoAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEACqMhAAAAAAAAAAADFlxoSKY5MeL30+42QGWjLF6ivid6lF/Y5nCoiztVYAAkExhSKjioAABvbgACmmoAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAACqMhAAAAAAAAAACqzeJOnnRTG9FhMj8xBI1lAI4tKb7wC8PgQB+Dp3pAOQAAAAAHJw4AAAqjIQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
        "fee_meta_xdr": "AAAAAgAAAAMACqMfAAAAAAAAAAADFlxoSKY5MeL30+42QGWjLF6ivid6lF/Y5nCoiztVYAAkExhZUUcMAABvbgACmmkAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEACqMhAAAAAAAAAAADFlxoSKY5MeL30+42QGWjLF6ivid6lF/Y5nCoiztVYAAkExhZUUaoAABvbgACmmkAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==",
        "memo_type": "text",
        "signatures": [
          "Yd9uIm+TjIcAjTU90YJoNg/r+6PU3Uss7ewUb1w3yMa+HyoSvDq8sDz/SYmDBH7F+0ACIeBF4kkVEKVBJMh0AQ=="
        ],
        "valid_after": "1970-01-01T00:00:00Z",
        "valid_before": "2019-04-09T20:19:21Z"
      }
    ]
  }
}
```

## Possible Errors

- The [standard errors](../errors.md#Standard-Errors).
- [not_found](../errors/not-found.md): A `not_found` error will be returned if there is no ledgers whose sequence matches the `id` argument.
