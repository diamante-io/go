---
title: Page
replacement: https://developers.diamcircle.org/api/introduction/pagination/page-arguments/
---

Pages represent a subset of a larger collection of objects.
As an example, it would be unfeasible to provide the
[All Transactions](../endpoints/transactions-all.md) endpoint without paging.  Over time there
will be millions of transactions in the diamcircle network's ledger and returning
them all over a single request would be unfeasible.

## Attributes

A page itself exposes no attributes.  It is merely a container for embedded
records and some links to aid in iterating the entire collection the page is
part of.

## Cursor
A `cursor` is a number that points to a specific location in a collection of resources.

The `cursor` attribute itself is an opaque value meaning that users should not try to parse it.

## Embedded Resources

A page contains an embedded set of `records`, regardless of the contained resource.

## Links

A page provides a couple of links to ease in iteration.

|      |                        Example                         |           Relation           |
| ---- | ------------------------------------------------------ | ---------------------------- |
| self | `/transactions`                                        |                              |
| prev | `/transactions?cursor=12884905984&order=desc&limit=10` | The previous page of results |
| next | `/transactions?cursor=12884905984&order=asc&limit=10`  | The next page of results     |

## Example

```json
{
  "_embedded": {
    "records": [
      {
        "_links": {
          "self": {
            "href": "/operations/12884905984"
          },
          "transaction": {
            "href": "/transaction/6391dd190f15f7d1665ba53c63842e368f485651a53d8d852ed442a446d1c69a"
          },
          "precedes": {
            "href": "/account/GCEZWKCA5VLDNRLN3RPRJMRZOX3Z6G5CHCGSNFHEYVXM3XOJMDS674JZ/payments?cursor=12884905984&order=asc{?limit}",
            "templated": true
          },
          "succeeds": {
            "href": "/account/GCEZWKCA5VLDNRLN3RPRJMRZOX3Z6G5CHCGSNFHEYVXM3XOJMDS674JZ/payments?cursor=12884905984&order=desc{?limit}",
            "templated": true
          }
        },
        "id": 12884905984,
        "paging_token": "12884905984",
        "type_i": 0,
        "type": "payment",
        "sender": "GCEZWKCA5VLDNRLN3RPRJMRZOX3Z6G5CHCGSNFHEYVXM3XOJMDS674JZ",
        "receiver": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU",
        "asset": {
          "code": "XLM"
        },
        "amount": 1000000000,
        "amount_f": 100.00
      }
    ]
  },
  "_links": {
    "next": {
      "href": "/account/GCEZWKCA5VLDNRLN3RPRJMRZOX3Z6G5CHCGSNFHEYVXM3XOJMDS674JZ/payments?cursor=12884905984&order=asc&limit=100"
    },
    "prev": {
      "href": "/account/GCEZWKCA5VLDNRLN3RPRJMRZOX3Z6G5CHCGSNFHEYVXM3XOJMDS674JZ/payments?cursor=12884905984&order=desc&limit=100"
    },
    "self": {
      "href": "/account/GCEZWKCA5VLDNRLN3RPRJMRZOX3Z6G5CHCGSNFHEYVXM3XOJMDS674JZ/payments?limit=100"
    }
  }
}

```

## Endpoints

Any endpoint that provides a collection of resources will represent them as pages.

