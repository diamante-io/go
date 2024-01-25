---
title: Response Format
replacement: https://developers.diamcircle.org/api/introduction/response-format/
---

Rather than using a fully custom way of representing the resources we expose in
Aurora, we use [HAL](http://stateless.co/hal_specification.html). HAL is a
hypermedia format in JSON that remains simple while giving us a couple of
benefits such as simpler client integration for several languages. See [this
wiki page](https://github.com/mikekelly/hal_specification/wiki/Libraries) for a
list of libraries.

## Attributes, Links, Embedded Resources

At its simplest, a HAL response is just a JSON object with a couple of reserved
property names:  `_links` is used for expressing links and `_embedded` is used
for bundling other HAL objects with the response.  Other than links and embedded
objects, **HAL is just JSON**.

### Links

HAL is a hypermedia format, like HTML, in that it has a mechanism to express
links between documents.  Let's look at a simple example:

```json
{
  "_links": {
    "self": {
      "href": "/ledgers/1"
    },
    "transactions": {
      "href": "/ledgers/1/transactions{?cursor,limit,order}",
      "templated": true
    }
  },
  "id": "43cf4db3741a7d6c2322e7b646320ce9d7b099a0b3501734dcf70e74a8a4e637",
  "hash": "43cf4db3741a7d6c2322e7b646320ce9d7b099a0b3501734dcf70e74a8a4e637",
  "prev_hash": "",
  "sequence": 1,
  "transaction_count": 0,
  "operation_count": 0,
  "closed_at": "0001-01-01T00:00:00Z",
  "total_coins": "100000000000.0000000",
  "fee_pool": "0.0000000",
  "base_fee_in_stroops": 100,
  "base_reserve_in_stroops": 100000000,
  "max_tx_set_size": 50
}
```

The above response is for the genesis ledger of the Diamcircle test network, and
the links in the `_links` attribute provide links to other relavant resources in
Aurora.  Notice the object beneath the `transactions` key.  The key of each
link specifies that links relation to the current resource, and in this case
`transactions` means "Transactions that occurred in this ledger".  Logically,
you should expect that resource to respond with a collection of transactions
with all of the results having a `ledger_sequence` attribute equal to 1.

The `transactions` link is also _templated_, which means that the `href`
attribute of the link is actually a URI template, as specified by  [RFC
6570](https://tools.ietf.org/html/rfc6570).  We use URI templates to show you
what parameters a given resource can take. You must evaluate the template to a
valid URI before navigating to it.

## Pages

Pages represent a subset of a larger collection of objects.  
As an example, it would be unfeasible to provide the
[All Transactions](../reference/endpoints/transactions-all.md) endpoint without paging.  
Over time there will be millions of transactions in the Diamcircle network's ledger
and returning them all over a single request would be unfeasible.

Read more about paging in following docs:
- [Page](../reference/resources/page.md)
- [Paging](./paging.md)
