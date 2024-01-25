---
title: Overview
---

The Go SDK is a set of packages for interacting with most aspects of the diamcircle ecosystem. The primary component is the aurora SDK, which provides convenient access to aurora services. There are also packages for other diamcircle services such as [TOML support](https://github.com/diamcircle/diamcircle-protocol/blob/master/ecosystem/sep-0001.md) and [federation](https://github.com/diamcircle/diamcircle-protocol/blob/master/ecosystem/sep-0002.md).

## aurora SDK

The aurora SDK is composed of two complementary libraries: `txnbuild` + `auroraclient`.
The `txnbuild` ([source](https://go/tree/master/txnbuild), [docs](https://godoc.org/go/txnbuild)) package enables the construction, signing and encoding of diamcircle [transactions](https://developers.diamcircle.org/docs/glossary/transactions/) and [operations](https://developers.diamcircle.org/docs/start/list-of-operations/) in Go. The `auroraclient` ([source](https://go/tree/master/clients/auroraclient), [docs](https://godoc.org/go/clients/auroraclient)) package provides a web client for interfacing with [aurora](https://developers.diamcircle.org/docs/start/introduction/) server REST endpoints to retrieve ledger information, and to submit transactions built with `txnbuild`.

## List of major SDK packages

- `auroraclient` ([source](https://go/tree/master/clients/auroraclient), [docs](https://godoc.org/go/clients/auroraclient)) - programmatic client access to aurora
- `txnbuild` ([source](https://go/tree/master/txnbuild), [docs](https://godoc.org/go/txnbuild)) - construction, signing and encoding of diamcircle transactions and operations
- `diamcircletoml` ([source](https://go/tree/master/clients/diamcircletoml), [docs](https://godoc.org/go/clients/diamcircletoml)) - parse [diamcircle.toml](../../guides/concepts/diamcircle-toml.md) files from the internet
- `federation` ([source](https://godoc.org/go/clients/federation)) - resolve federation addresses  into diamcircle account IDs, suitable for use within a transaction

