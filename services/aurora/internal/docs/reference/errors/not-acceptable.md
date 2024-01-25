---
title: Not Acceptable
replacement: https://developers.diamcircle.org/api/errors/http-status-codes/standard/
---

When your client only accepts certain formats of data from aurora and aurora cannot fulfill that
request, aurora will return a `not_acceptable` error. This is analogous to a
[HTTP 406 Error](https://developer.mozilla.org/en-US/docs/Web/HTTP/Response_codes).

For example, if your client only accepts an XML response (`Accept: application/xml`), aurora will
respond with a `not_acceptable` error.

If you are encountering this error, please check to make sure the criteria for content you’ll
accept is correct.

## Attributes

As with all errors aurora returns, `not_acceptable` follows the
[Problem Details for HTTP APIs](https://tools.ietf.org/html/draft-ietf-appsawg-http-problem-00)
draft specification guide and thus has the following attributes:

| Attribute   | Type   | Description                                                                     |
| ----------- | ------ | ------------------------------------------------------------------------------- |
| `type`      | URL    | The identifier for the error.  This is a URL that can be visited in the browser.|
| `title`     | String | A short title describing the error.                                             |
| `status`    | Number | An HTTP status code that maps to the error.                                     |
| `detail`    | String | A more detailed description of the error.                                       |

## Example

```shell
$ curl -X GET -H "Accept: application/xml" "https://aurora-testnet.diamcircle.org/accounts/GALWEV6GY73RJ255JC7XUOZ2L7WZ5JJDTKATB2MUK7F3S67DVT2A6R5G"
{
  "type": "https://diamcircle.org/aurora-errors/not_acceptable",
  "title": "An acceptable response content-type could not be provided for this request",
  "status": 406
}
```

## Related

- [Not Found](./not-found.md)
