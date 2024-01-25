---
title: Rate Limiting
replacement: https://developers.diamcircle.org/api/introduction/rate-limiting/
---

In order to provide service stability, Aurora limits the number of requests a
client (single IP) can perform within a one hour window.  By default this is set to 3600
requests per hour—an average of one request per second. Also, while streaming
every update of the stream (what happens every time there's a new ledger) is
counted. Ex. if there were 12 new ledgers in a minute, 12 requests will be
subtracted from the limit.

Aurora is using [GCRA](https://brandur.org/rate-limiting#gcra) algorithm.

## Response headers for rate limiting

Every response from Aurora sets advisory headers to inform clients of their
standing with rate limiting system:

|          Header         |                               Description                                |
| ----------------------- | ------------------------------------------------------------------------ |
| `X-RateLimit-Limit`     | The maximum number of requests that the current client can make in one hour. |
| `X-RateLimit-Remaining` | The number of remaining requests for the current window.                 |
| `X-RateLimit-Reset`     | Seconds until a new window starts.                                        |

In addition, a `Retry-After` header will be set when the current client is being
throttled.
