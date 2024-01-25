---
title: Rate Limit Exceeded
replacement: https://developers.diamcircle.org/api/errors/http-status-codes/standard/
---

When a single user makes too many requests to aurora in a one hour time frame, aurora returns a
`rate_limit_exceeded` error. By default, aurora allows 3600 requests per hour -- an average of one
request per second. This is analogous to a
[HTTP 429 Error](https://developer.mozilla.org/en-US/docs/Web/HTTP/Response_codes).

If you are encountering this error, please reduce your request speed. Here are some strategies for
doing so:
* For collection endpoints, try specifying larger page sizes.
* Try streaming responses to watch for new data instead of pulling data every time.
* Cache immutable data, such as transaction details, locally.

See the [Rate Limiting Guide](../../reference/rate-limiting.md) for more info.

## Attributes

As with all errors aurora returns, `rate_limit_exceeded` follows the
[Problem Details for HTTP APIs](https://tools.ietf.org/html/draft-ietf-appsawg-http-problem-00)
draft specification guide and thus has the following attributes:

| Attribute   | Type   | Description                                                                     |
| ----------- | ------ | ------------------------------------------------------------------------------- |
| `type`      | URL    | The identifier for the error.  This is a URL that can be visited in the browser.|
| `title`     | String | A short title describing the error.                                             |
| `status`    | Number | An HTTP status code that maps to the error.                                     |
| `detail`    | String | A more detailed description of the error.                                       |

## Example

```json
{
  "type": "https://diamcircle.org/aurora-errors/rate_limit_exceeded",
  "title": "Rate Limit Exceeded",
  "status": 429,
  "details": "The rate limit for the requesting IP address is over its alloted limit.  The allowed limit and requests left per time period are communicated to clients via the http response headers 'X-RateLimit-*' headers."
}
```
