---
title: Not Implemented
replacement: https://developers.diamcircle.org/api/errors/http-status-codes/standard/
---

If your [request method](http://www.w3.org/Protocols/rfc2616/rfc2616-sec9.html) is not supported by
aurora, aurora will return a `not_implemented` error. Likewise, if functionality that is intended
but does not exist (thus reserving the endpoint for future use), it will also return a
`not_implemented` error. This is analogous to a
[HTTP 501 Error](https://developer.mozilla.org/en-US/docs/Web/HTTP/Response_codes).

If you are encountering this error, aurora does not have the functionality you are requesting.

## Attributes

As with all errors aurora returns, `not_implemented` follows the
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
$ curl -X GET "https://aurora-testnet.diamcircle.org/offers/1234"
{
  "type": "https://diamcircle.org/aurora-errors/not_implemented",
  "title": "Resource Not Yet Implemented",
  "status": 404,
  "detail": "While the requested URL is expected to eventually point to a valid resource, the work to implement the resource has not yet been completed."
}
```

## Related

- [Server Error](./server-error.md)
