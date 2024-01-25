---
title: Payments for Account
clientData:
  laboratoryUrl: https://www.diamcircle.org/laboratory/#explorer?resource=payments&endpoint=for_account
replacement: https://developers.diamcircle.org/api/resources/accounts/payments/
---

This endpoint responds with a collection of payment-related operations where the given
[account](../resources/account.md) was either the sender or receiver.

This endpoint can also be used in [streaming](../streaming.md) mode so it is possible to use it to
listen for new payments to or from an account as they get made in the diamcircle network.
If called in streaming mode aurora will start at the earliest known payment unless a `cursor` is
set. In that case it will start from the `cursor`. You can also set `cursor` value to `now` to only
stream payments created since your request time.

The operations that can be returned in by this endpoint are:
- `create_account`
- `payment`
- `path_payment`
- `account_merge`

## Request

```
GET /accounts/{id}/payments{?cursor,limit,order}
```

### Arguments

| name | notes | description | example |
| ---- | ----- | ----------- | ------- |
| `id` | required, string | The account id of the account used to constrain results. | `GCEZWKCA5VLDNRLN3RPRJMRZOX3Z6G5CHCGSNFHEYVXM3XOJMDS674JZ` |
| `?cursor` | optional, default _null_ | A payment paging token specifying from where to begin results. When streaming this can be set to `now` to stream object created since your request time. | `8589934592` |
| `?limit` | optional, number, default `10` | Specifies the count of records at most to return. | `200` |
| `?order` | optional, string, default `asc` | Specifies order of returned results. `asc` means older payments first, `desc` mean newer payments first. | `desc` |
| `?include_failed` | optional, bool, default: `false` | Set to `true` to include payments of failed transactions in results. | `true` |
| `?join` | optional, string, default: _null_ | Set to `transactions` to include the transactions which created each of the payments in the response. | `transactions` |

### curl Example Request

```bash
# Retrieve the 25 latest payments for a specific account.
curl "https://aurora-testnet.diamcircle.org/accounts/GCEZWKCA5VLDNRLN3RPRJMRZOX3Z6G5CHCGSNFHEYVXM3XOJMDS674JZ/payments?limit=25&order=desc"
```

### JavaScript Example Request

```javascript
var diamcircleSdk = require('diamcircle-sdk');
var server = new diamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

server.payments()
  .forAccount("GCEZWKCA5VLDNRLN3RPRJMRZOX3Z6G5CHCGSNFHEYVXM3XOJMDS674JZ")
  .call()
  .then(function (accountResult) {
    console.log(accountResult);
  })
  .catch(function (err) {
    console.error(err);
  })
```

### JavaScript Streaming Example

```javascript
var diamcircleSdk = require('diamcircle-sdk')
var server = new diamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

var paymentHandler = function (paymentResponse) {
  console.log(paymentResponse);
};

var es = server.payments()
  .forAccount("GCEZWKCA5VLDNRLN3RPRJMRZOX3Z6G5CHCGSNFHEYVXM3XOJMDS674JZ")
  .cursor('now')
  .stream({
    onmessage: paymentHandler
  })
```

## Response

This endpoint responds with a [page](../resources/page.md) of [payment operations](../resources/operation.md).

### Example Response

```json
{"_embedded": {
  "records": [
    {
      "_links": {
        "self": {
          "href": "/operations/12884905984"
        },
        "transaction": {
          "href": "/transaction/6391dd190f15f7d1665ba53c63842e368f485651a53d8d852ed442a446d1c69a"
        },
        "precedes": {
          "href": "/accounts/GCEZWKCA5VLDNRLN3RPRJMRZOX3Z6G5CHCGSNFHEYVXM3XOJMDS674JZ/payments?cursor=12884905984&order=asc{?limit}",
          "templated": true
        },
        "succeeds": {
          "href": "/accounts/GCEZWKCA5VLDNRLN3RPRJMRZOX3Z6G5CHCGSNFHEYVXM3XOJMDS674JZ/payments?cursor=12884905984&order=desc{?limit}",
          "templated": true
        }
      },
      "id": 12884905984,
      "paging_token": "12884905984",
      "type_i": 0,
      "type": "payment",
      "sender": "GCEZWKCA5VLDNRLN3RPRJMRZOX3Z6G5CHCGSNFHEYVXM3XOJMDS674JZ",
      "receiver": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU",
      "asset": {
        "code": "XLM"
      },
      "amount": 1000000000,
      "amount_f": 100.00
    }
  ]
},
"_links": {
  "next": {
    "href": "/accounts/GCEZWKCA5VLDNRLN3RPRJMRZOX3Z6G5CHCGSNFHEYVXM3XOJMDS674JZ/payments?cursor=12884905984&order=asc{?limit}",
    "templated": true
  },
  "self": {
    "href": "/accounts/GCEZWKCA5VLDNRLN3RPRJMRZOX3Z6G5CHCGSNFHEYVXM3XOJMDS674JZ/payments"
  }
}
}
```

### Example Streaming Event

```json
{
  "_links": {
    "effects": {
      "href": "/operations/77309415424/effects/{?cursor,limit,order}",
      "templated": true
    },
    "precedes": {
      "href": "/operations?cursor=77309415424&order=asc"
    },
    "self": {
      "href": "/operations/77309415424"
    },
    "succeeds": {
      "href": "/operations?cursor=77309415424&order=desc"
    },
    "transactions": {
      "href": "/transactions/77309415424"
    }
  },
  "account": "GBIA4FH6TV64KSPDAJCNUQSM7PFL4ILGUVJDPCLUOPJ7ONMKBBVUQHRO",
  "funder": "GCEZWKCA5VLDNRLN3RPRJMRZOX3Z6G5CHCGSNFHEYVXM3XOJMDS674JZ",
  "id": 77309415424,
  "paging_token": "77309415424",
  "starting_balance": 1e+14,
  "type_i": 0,
  "type": "create_account"
}
```

## Possible Errors

- The [standard errors](../errors.md#standard-errors).
