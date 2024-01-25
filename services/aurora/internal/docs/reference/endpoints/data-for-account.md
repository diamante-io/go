---
title: Data for Account
clientData:
  laboratoryUrl: https://www.diamcircle.org/laboratory/#explorer?resource=data&endpoint=for_account
replacement: https://developers.diamcircle.org/api/resources/accounts/data/
---

This endpoint represents a single [data](../resources/data.md) associated with a given [account](../resources/account.md).

## Request

```
GET /accounts/{account}/data/{key}
```

### Arguments

| name     | notes                          | description                                                      | example                                                   |
| ------   | -------                        | -----------                                                      | -------                                                   |
| `key`| required, string               | Key name | `user-id`|

### curl Example Request

```sh
curl "https://aurora-testnet.diamcircle.org/accounts/GA2HGBJIJKI6O4XEM7CZWY5PS6GKSXL6D34ERAJYQSPYA6X6AI7HYW36/data/user-id"
```

### JavaScript Example Request

```javascript
var diamcircleSdk = require('diamcircle-sdk');
var server = new diamcircleSdk.Server('https://aurora-testnet.diamcircle.org');

server.accounts()
  .accountId("GAKLBGHNHFQ3BMUYG5KU4BEWO6EYQHZHAXEWC33W34PH2RBHZDSQBD75")
  .call()
  .then(function (account) {
    return account.data({key: 'user-id'})
  })
  .then(function(dataValue) {
    console.log(dataValue)
  })
  .catch(function (err) {
    console.log(err)
  })
```

## Response

This endpoint responds with a value of the data field for the given account. See [data resource](../resources/data.md) for reference.

### Example Response

```json
{
  "value": "MTAw"
}
```

## Possible Errors

- The [standard errors](../errors.md#Standard-Errors).
- [not_found](../errors/not-found.md): A `not_found` error will be returned if there is no account whose ID matches the `account` argument or there is no data field with a given key.
