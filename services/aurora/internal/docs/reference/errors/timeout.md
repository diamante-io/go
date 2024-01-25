---
title: Timeout
replacement: https://developers.diamcircle.org/api/errors/http-status-codes/aurora-specific/
---

If you are encountering this error it means that either:

* Aurora has not received a confirmation from the Diamcircle Core server that the transaction you are
  trying to submit to the network was included in a ledger in a timely manner.
* Aurora has not sent a response to a reverse-proxy before a specified amount of time has elapsed.

The former case may happen because there was no room for your transaction for 3 consecutive
ledgers. This is because Diamcircle Core removes each submitted transaction from a queue. To solve
this you can:

* Keep resubmitting the same transaction (with the same sequence number) and wait until it finally
  is added to a new ledger.
* Increase the [fee](../../../guides/concepts/fees.md) in order to prioritize the transaction.

This error returns a
[HTTP 504 Error](https://developer.mozilla.org/en-US/docs/Web/HTTP/Response_codes).

## Attributes

As with all errors Aurora returns, `timeout` follows the
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
  "type": "https://diamcircle.org/aurora-errors/timeout",
  "title": "Timeout",
  "status": 504,
  "detail": "Your request timed out before completing.  Please try your request again. If you are submitting a transaction make sure you are sending exactly the same transaction (with the same sequence number)."
}
```

## Related

- [Not Acceptable](./not-acceptable.md)
- [Transaction Failed](./transaction-failed.md)
