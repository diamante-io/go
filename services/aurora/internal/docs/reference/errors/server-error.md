---
title: Internal Server Error
replacement: https://developers.diamcircle.org/api/errors/http-status-codes/standard/
---

If there's an internal error within Aurora, Aurora will return a
`server_error` response.  This response is a catch-all, and can refer to many
possible errors in the Aurora server: a configuration mistake, a database
connection error, etc. This is analogous to a
[HTTP 500 Error](https://developer.mozilla.org/en-US/docs/Web/HTTP/Response_codes).

Aurora does not expose information such as stack traces or raw error messages
to a client, as doing so may reveal sensitive configuration data such as secret
keys. If you are encountering this error on a server you control, please check the
Aurora log files for more details. The logs should contain detailed
information to help you discover the root issue.

If you are encountering this error on the public Diamcircle infrastructure, please
report an error on [Aurora's issue tracker](https://github.com/diamcircle/go/issues)
and include as much information about the request that triggered the response
as you can (especially the time of the request).

## Attributes

As with all errors Aurora returns, `server_error` follows the
[Problem Details for HTTP APIs](https://tools.ietf.org/html/draft-ietf-appsawg-http-problem-00)
draft specification guide and thus has the following attributes:

| Attribute   | Type   | Description                                                                     |
| ----------- | ------ | ------------------------------------------------------------------------------- |
| `type`      | URL    | The identifier for the error.  This is a URL that can be visited in the browser.|
| `title`     | String | A short title describing the error.                                             |
| `status`    | Number | An HTTP status code that maps to the error.                                     |
| `detail`    | String | A more detailed description of the error.                                       |

## Examples
```json
{
  "type": "https://diamcircle.org/aurora-errors/server_error",
  "title": "Internal Server Error",
  "status": 500,
  "details": "An error occurred while processing this request. This is usually due to a bug within the server software. Trying this request again may succeed if the bug is transient, otherwise please report this issue to the issue tracker at: https://github.com/diamcircle/go/issues. Please include this response in your issue."
}
```

## Related

- [Not Implemented](./not-implemented.md)
