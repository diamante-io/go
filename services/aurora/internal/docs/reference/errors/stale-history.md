---
title: Stale History
replacement: https://developers.diamcircle.org/api/errors/http-status-codes/aurora-specific/
---

A aurora server may be configured to reject historical requests when the history is known to be
further out of date than the configured threshold. In such cases, this error is returned.  To
resolve this error (provided you are the aurora instance's operator) please ensure that the
ingestion system is running correctly and importing new ledgers. This error returns a
[HTTP 503 Error](https://developer.mozilla.org/en-US/docs/Web/HTTP/Response_codes).

## Attributes

As with all errors Aurora returns, `stale_history` follows the
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

  "type": "https://diamcircle.org/aurora-errors/stale_history",
  "title": "Historical DB Is Too Stale",
  "status": 503,
  "detail": "This aurora instance is configured to reject client requests when it can determine that the history database is lagging too far behind the connected instance of diamcircle-core.  If you operate this server, please ensure that the ingestion system is properly running."
}
```

## Related

- [Internal Server Error](./server-error.md)
