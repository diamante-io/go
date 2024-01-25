---
title: Overview
---

The Go SDK is a set of packages for interacting with most aspects of the Diamcircle ecosystem. The primary component is the Aurora SDK, which provides convenient access to Aurora services. There are also packages for other Diamcircle services such as [TOML support](https://github.com/diamcircle/diamcircle-protocol/blob/master/ecosystem/sep-0001.md) and [federation](https://github.com/diamcircle/diamcircle-protocol/blob/master/ecosystem/sep-0002.md).

## Aurora SDK

The Aurora SDK is composed of two complementary libraries: `txnbuild` + `auroraclient`.
The `txnbuild` ([source](https://github.com/diamcircle/go/tree/master/txnbuild), [docs](https://godoc.org/github.com/diamcircle/go/txnbuild)) package enables the construction, signing and encoding of Diamcircle [transactions](https://developers.diamcircle.org/docs/glossary/transactions/) and [operations](https://developers.diamcircle.org/docs/start/list-of-operations/) in Go. The `auroraclient` ([source](https://github.com/diamcircle/go/tree/master/clients/auroraclient), [docs](https://godoc.org/github.com/diamcircle/go/clients/auroraclient)) package provides a web client for interfacing with [Aurora](https://developers.diamcircle.org/docs/start/introduction/) server REST endpoints to retrieve ledger information, and to submit transactions built with `txnbuild`.

## List of major SDK packages

- `auroraclient` ([source](https://github.com/diamcircle/go/tree/master/clients/auroraclient), [docs](https://godoc.org/github.com/diamcircle/go/clients/auroraclient)) - programmatic client access to Aurora
- `txnbuild` ([source](https://github.com/diamcircle/go/tree/master/txnbuild), [docs](https://godoc.org/github.com/diamcircle/go/txnbuild)) - construction, signing and encoding of Diamcircle transactions and operations
- `diamcircletoml` ([source](https://github.com/diamcircle/go/tree/master/clients/diamcircletoml), [docs](https://godoc.org/github.com/diamcircle/go/clients/diamcircletoml)) - parse [Diamcircle.toml](../../guides/concepts/diamcircle-toml.md) files from the internet
- `federation` ([source](https://godoc.org/github.com/diamcircle/go/clients/federation)) - resolve federation addresses  into diamcircle account IDs, suitable for use within a transaction

