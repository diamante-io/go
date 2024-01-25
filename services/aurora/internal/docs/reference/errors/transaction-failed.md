---
title: Transaction Failed
replacement: https://developers.diamcircle.org/api/errors/http-status-codes/aurora-specific/
---

The `transaction_failed` error occurs when a client submits a transaction that was well-formed but
was not included into the ledger due to some other failure. For example, a transaction may fail if:

- The source account for transaction cannot pay the minimum fee.
- The sequence number is incorrect.
- One of the contained operations has failed such as a payment operation that overdraws on the
  paying account.

In almost every case, this error indicates that the transaction submitted in the initial request
will never succeed. There is one exception: a transaction that fails with the `tx_bad_seq` result
code (as expressed in the `result_code` field of the error) may become valid in the future if the
sequence number it used was too high.

This error returns a
[HTTP 400 Error](https://developer.mozilla.org/en-US/docs/Web/HTTP/Response_codes).

## Attributes

As with all errors aurora returns, `transaction_failed` follows the
[Problem Details for HTTP APIs](https://tools.ietf.org/html/draft-ietf-appsawg-http-problem-00)
draft specification guide and thus has the following attributes:

| Attribute   | Type   | Description                                                                     |
| ----------- | ------ | ------------------------------------------------------------------------------- |
| `type`      | URL    | The identifier for the error.  This is a URL that can be visited in the browser.|
| `title`     | String | A short title describing the error.                                             |
| `status`    | Number | An HTTP status code that maps to the error.                                     |
| `detail`    | String | A more detailed description of the error.                                       |

In addition, the following additional data is provided in the `extras` field of the error:

| Attribute                  | Type   | Description                                                                                                                 |
|----------------------------|--------|-----------------------------------------------------------------------------------------------------------------------------|
| `envelope_xdr`             | String | A base64-encoded representation of the TransactionEnvelope XDR whose failure triggered this response.                       |
| `result_xdr`               | String | A base64-encoded representation of the TransactionResult XDR returned by diamcircle-core when submitting this transaction.     |
| `result_codes.transaction` | String | The transaction result code returned by diamcircle Core.                                                                       |
| `result_codes.operations`  | Array  | An array of strings, representing the operation result codes for each operation in the submitted transaction, if available. |


## Examples

### No Source Account
```json
{
  "type": "https://diamcircle.org/aurora-errors/transaction_failed",
  "title": "Transaction Failed",
  "status": 400,
  "detail": "The transaction failed when submitted to the diamcircle network. The `extras.result_codes` field on this response contains further details.  Descriptions of each code can be found at: https://www.diamcircle.org/developers/learn/concepts/list-of-operations.html",
  "extras": {
    "envelope_xdr": "AAAAANNVpdQ9vctZdAJ67sFmNe1KDzaj51dAdkW3vKKM51H3AAAAZAAAAABJlgLSAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAA01Wl1D29y1l0AnruwWY17UoPNqPnV0B2Rbe8ooznUfcAAAAAAAAAAAL68IAAAAAAAAAAAA==",
    "result_codes": {
      "transaction": "tx_no_source_account"
    },
    "result_xdr": "AAAAAAAAAAD////4AAAAAA=="
  }
}
```

### Bad Authentication
```json
{
  "type": "https://diamcircle.org/aurora-errors/transaction_failed",
  "title": "Transaction Failed",
  "status": 400,
  "detail": "The transaction failed when submitted to the diamcircle network. The `extras.result_codes` field on this response contains further details.  Descriptions of each code can be found at: https://www.diamcircle.org/developers/learn/concepts/list-of-operations.html",
  "extras": {
    "envelope_xdr": "AAAAAPORy3CoX6ox2ilbeiVjBA5WlpCSZRcjZ7VE9Wf4QVk7AAAAZAAAQz0AAAACAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAA85HLcKhfqjHaKVt6JWMEDlaWkJJlFyNntUT1Z/hBWTsAAAAAAAAAAAL68IAAAAAAAAAAARN17BEAAABAA9Ad7OKc7y60NT/JuobaHOfmuq8KbZqcV6G/es94u9yT84fi0aI7tJsFMOyy8cZ4meY3Nn908OU+KfRWV40UCw==",
    "result_codes": {
      "transaction": "tx_bad_auth"
    },
    "result_xdr": "AAAAAAAAAGT////6AAAAAA=="
  }
}
```

## Related

- [Transaction Malformed](./transaction-malformed.md)
