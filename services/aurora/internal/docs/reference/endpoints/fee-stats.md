---
title: Fee Stats
clientData:
  laboratoryUrl:
replacement: https://developers.diamcircle.org/api/aggregations/fee-stats/
---

This endpoint gives useful information about per-operation fee stats in the last 5 ledgers. It can be used to
predict a fee set on the transaction that will be submitted to the network.

## Request

```
GET /fee_stats
```

### curl Example Request

```sh
curl "https://aurora-testnet.diamcircle.org/fee_stats"
```

## Response

Response contains the following fields:

| Field | |
| - | - |
| last_ledger | Last ledger sequence number |
| last_ledger_base_fee | Base fee as defined in the last ledger |
| ledger_capacity_usage | Average capacity usage over the last 5 ledgers. (0 is no usage, 1.0 is completely full ledgers) |
| fee_charged      | fee charged object |
| max_fee          | max fee object |

### Fee Charged Object

Information about the fee charged for transactions in the last 5 ledgers.

| Field | |
| - | - |
| min | Minimum fee charged over the last 5 ledgers. |
| mode | Mode fee charged over the last 5 ledgers. |
| p10 | 10th percentile fee charged over the last 5 ledgers. |
| p20 | 20th percentile fee charged over the last 5 ledgers. |
| p30 | 30th percentile fee charged over the last 5 ledgers. |
| p40 | 40th percentile fee charged over the last 5 ledgers. |
| p50 | 50th percentile fee charged over the last 5 ledgers. |
| p60 | 60th percentile fee charged over the last 5 ledgers. |
| p70 | 70th percentile fee charged over the last 5 ledgers. |
| p80 | 80th percentile fee charged over the last 5 ledgers. |
| p90 | 90th percentile fee charged over the last 5 ledgers. |
| p95 | 95th percentile fee charged over the last 5 ledgers. |
| p99 | 99th percentile fee charged over the last 5 ledgers. |

Note: The difference between `fee_charged` and `max_fee` is that the former
represents the actual fee paid for the transaction while `max_fee` represents
the maximum bid the transaction creator was willing to pay for the transaction.

### Max Fee Object

Information about max fee bid for transactions over the last 5 ledgers.

| Field | |
| - | - |
| min | Minimum (lowest) value of the maximum fee bid over the last 5 ledgers. |
| mode | Mode max fee over the last 5 ledgers. |
| p10 | 10th percentile max fee over the last 5 ledgers. |
| p20 | 20th percentile max fee over the last 5 ledgers. |
| p30 | 30th percentile max fee over the last 5 ledgers. |
| p40 | 40th percentile max fee over the last 5 ledgers. |
| p50 | 50th percentile max fee over the last 5 ledgers. |
| p60 | 60th percentile max fee over the last 5 ledgers. |
| p70 | 70th percentile max fee over the last 5 ledgers. |
| p80 | 80th percentile max fee over the last 5 ledgers. |
| p90 | 90th percentile max fee over the last 5 ledgers. |
| p95 | 95th percentile max fee over the last 5 ledgers. |
| p99 | 99th percentile max fee over the last 5 ledgers. |


### Example Response

```json
{
  "last_ledger": "22606298",
  "last_ledger_base_fee": "100",
  "ledger_capacity_usage": "0.97",
  "fee_charged": {
    "max": "100",
    "min": "100",
    "mode": "100",
    "p10": "100",
    "p20": "100",
    "p30": "100",
    "p40": "100",
    "p50": "100",
    "p60": "100",
    "p70": "100",
    "p80": "100",
    "p90": "100",
    "p95": "100",
    "p99": "100"
  },
  "max_fee": {
    "max": "100000",
    "min": "100",
    "mode": "100",
    "p10": "100",
    "p20": "100",
    "p30": "100",
    "p40": "100",
    "p50": "100",
    "p60": "100",
    "p70": "100",
    "p80": "100",
    "p90": "15000",
    "p95": "100000",
    "p99": "100000"
  }
}
```

## Possible Errors

- The [standard errors](../errors.md#standard-errors).
