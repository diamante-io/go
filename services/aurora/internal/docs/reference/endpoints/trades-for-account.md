---
title: Trades for Account
clientData:
  laboratoryUrl: https://www.diamcircle.org/laboratory/#explorer?resource=trades&endpoint=for_account
replacement: https://developers.diamcircle.org/api/resources/accounts/trades/
---

This endpoint represents all [trades](../resources/trade.md) that affect a given [account](../resources/account.md).

This endpoint can also be used in [streaming](../streaming.md) mode, making it possible to listen for new trades that affect the given account as they occur on the diamcircle network.
If called in streaming mode aurora will start at the earliest known trade unless a `cursor` is set. In that case it will start from the `cursor`. You can also set `cursor` value to `now` to only stream trades created since your request time.

## Request

```
GET /accounts/{account_id}/trades{?cursor,limit,order}
```

### Arguments

| name | notes | description | example |
| ---- | ----- | ----------- | ------- |
| `account_id` | required, string | ID of an account | GBYTR4MC5JAX4ALGUBJD7EIKZVM7CUGWKXIUJMRSMK573XH2O7VAK3SR |
| `?cursor` | optional, any, default _null_ | A paging token, specifying where to start returning records from. When streaming this can be set to `now` to stream object created since your request time. | 12884905984 |
| `?order`  | optional, string, default `asc` | The order in which to return rows, "asc" or "desc". | `asc` |
| `?limit`  | optional, number, default: `10` | Maximum number of records to return. | `200` |

### curl Example Request

```sh
curl "https://aurora-testnet.diamcircle.org/accounts/GBYTR4MC5JAX4ALGUBJD7EIKZVM7CUGWKXIUJMRSMK573XH2O7VAK3SR/trades?limit=1"
```

### JavaScript Example Request

```javascript
var diamcircleSdk = require('diamcircle-sdk');
var server = new diamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

server.trades()
  .forAccount("GBYTR4MC5JAX4ALGUBJD7EIKZVM7CUGWKXIUJMRSMK573XH2O7VAK3SR")
  .call()
  .then(function (accountResult) {
    console.log(accountResult);
  })
  .catch(function (err) {
    console.error(err);
  })
```


## Response

This endpoint responds with a list of trades that changed a given account's state. See the [trade resource](../resources/trade.md) for reference.

### Example Response
```json
{
  "_links": {
    "self": {
      "href": "/accounts/GBYTR4MC5JAX4ALGUBJD7EIKZVM7CUGWKXIUJMRSMK573XH2O7VAK3SR/trades?cursor=\u0026limit=1\u0026order=asc"
    },
    "next": {
      "href": "/accounts/GBYTR4MC5JAX4ALGUBJD7EIKZVM7CUGWKXIUJMRSMK573XH2O7VAK3SR/trades?cursor=940258535411713-0\u0026limit=1\u0026order=asc"
    },
    "prev": {
      "href": "/accounts/GBYTR4MC5JAX4ALGUBJD7EIKZVM7CUGWKXIUJMRSMK573XH2O7VAK3SR/trades?cursor=940258535411713-0\u0026limit=1\u0026order=desc"
    }
  },
  "_embedded": {
    "records": [
      {
        "_links": {
          "self": {
            "href": ""
          },
          "base": {
            "href": "/accounts/GBYTR4MC5JAX4ALGUBJD7EIKZVM7CUGWKXIUJMRSMK573XH2O7VAK3SR"
          },
          "counter": {
            "href": "/accounts/GBOOAYCAJIN7YCUUAHEQJJARNQMRUP4P2WXVO6P4KAMAB27NGA3CYTZU"
          },
          "operation": {
            "href": "/operations/940258535411713"
          }
        },
        "id": "940258535411713-0",
        "paging_token": "940258535411713-0",
        "ledger_close_time": "2017-03-30T13:20:41Z",
        "offer_id": "8",
        "base_offer_id": "8",
        "base_account": "GBYTR4MC5JAX4ALGUBJD7EIKZVM7CUGWKXIUJMRSMK573XH2O7VAK3SR",
        "base_amount": "1.0000000",
        "base_asset_type": "credit_alphanum4",
        "base_asset_code": "BTC",
        "base_asset_issuer": "GB6FN4C7ZLWKENAOZDLZOQHNIOK4RDMV6EKLR53LWCHEBR6LVXOEKDZH",
        "counter_offer_id": "4611686044197195777",
        "counter_account": "GBOOAYCAJIN7YCUUAHEQJJARNQMRUP4P2WXVO6P4KAMAB27NGA3CYTZU",
        "counter_amount": "1.0000000",
        "counter_asset_type": "native",
        "base_is_seller": true,
        "price": {
          "n": 1,
          "d": 1
        }
      }
    ]
  }
}
```

## Example Streaming Event
```
{ 
  _links: 
    { self: { href: '' },
      base: { href: '/accounts/GDICGE2CFCNM3ZWRUVOWDJB2RAO667UE7WOSJJ2Z3IMISUA7CJZCE3KO' },
      counter: { href: '/accounts/GBILENMVJPVPEPXUPUPRBUEAME5OUQWAHIGZAX7TQX65NIQW3G3DGUYX' },
      operation: { href: '/operations/47274327069954049' } },
  id: '47274327069954049-0',
  paging_token: '47274327069954049-0',
  ledger_close_time: '2018-09-12T00:00:34Z',
  offer_id: '711437',
  base_account: 'GDICGE2CFCNM3ZWRUVOWDJB2RAO667UE7WOSJJ2Z3IMISUA7CJZCE3KO',
  base_amount: '13.0000000',
  base_asset_type: 'native',
  counter_account: 'GBILENMVJPVPEPXUPUPRBUEAME5OUQWAHIGZAX7TQX65NIQW3G3DGUYX',
  counter_amount: '13.0000000',
  counter_asset_type: 'credit_alphanum4',
  counter_asset_code: 'CNY',
  counter_asset_issuer: 'GAREELUB43IRHWEASCFBLKHURCGMHE5IF6XSE7EXDLACYHGRHM43RFOX',
  base_is_seller: true,
  price: { n: 1, d: 1 } 
}
```

## Possible Errors

- The [standard errors](../errors.md#Standard-Errors).
- [not_found](../errors/not-found.md): A `not_found` error will be returned if there is no account whose ID matches the `account_id` argument.
