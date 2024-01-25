---
title: Transaction Malformed
replacement: https://developers.diamcircle.org/api/errors/http-status-codes/aurora-specific/
---

When you submit a malformed transaction to aurora, aurora will return a `transaction_malformed`
error. There are many ways in which a transaction could be malformed, including:

- You submitted an empty string.
- Your base64-encoded string is invalid.
- Your [XDR](../xdr.md) structure is invalid.
- You have leftover bytes in your [XDR](../xdr.md) structure.

If you are encountering this error, please check the contents of the transaction you are
submitting. This error returns a
[HTTP 400 Error](https://developer.mozilla.org/en-US/docs/Web/HTTP/Response_codes).

## Attributes

As with all errors aurora returns, `transaction_malformed` follows the
[Problem Details for HTTP APIs](https://tools.ietf.org/html/draft-ietf-appsawg-http-problem-00)
draft specification guide and thus has the following attributes:

| Attribute   | Type   | Description                                                                     |
| ----------- | ------ | ------------------------------------------------------------------------------- |
| `type`      | URL    | The identifier for the error.  This is a URL that can be visited in the browser.|
| `title`     | String | A short title describing the error.                                             |
| `status`    | Number | An HTTP status code that maps to the error.                                     |
| `detail`    | String | A more detailed description of the error.                                       |

In addition, the following additional data is provided in the `extras` field of the error:

| Attribute      | Type   | Description                                        |
|----------------|--------|----------------------------------------------------|
| `envelope_xdr` | String | The submitted data that was malformed in some way. |

## Example

```json
{
  "type": "https://diamcircle.org/aurora-errors/transaction_malformed",
  "title": "Transaction Malformed",
  "status": 400,
  "detail": "aurora could not decode the transaction envelope in this request. A transaction should be an XDR TransactionEnvelope struct encoded using base64.  The envelope read from this request is echoed in the `extras.envelope_xdr` field of this response for your convenience.",
  "extras": {
    "envelope_xdr": "BBBBBPORy3CoX6ox2ilbeiVjBA5WlpCSZRcjZ7VE9Wf4QVk7AAAAZAAAQz0AAAACAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAA85HLcKhfqjHaKVt6JWMEDlaWkJJlFyNntUT1Z/hBWTsAAAAAAAAAAAL68IAAAAAAAAAAARN17BEAAABAA9Ad7OKc7y60NT/JuobaHOfmuq8KbZqcV6G/es94u9yT84fi0aI7tJsFMOyy8cZ4meY3Nn908OU+KfRWV40UCw=="
  }
}
```

## Related

- [Bad Request](./bad-request.md)
