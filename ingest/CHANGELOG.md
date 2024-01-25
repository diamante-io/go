# Changelog

All notable changes to this project will be documented in this file. This project adheres to [Semantic Versioning](http://semver.org/).


## Unreleased

* Let filewatcher use binary hash instead of timestamp to detect core version update [4050](https://github.com/diamcircle/go/pull/4050)

### New Features
* **Performance improvement**: the Captive Core backend now reuses bucket files whenever it finds existing ones in the corresponding `--captive-core-storage-path` (introduced in [v2.0](#v2.0.0)) rather than generating a one-time temporary sub-directory ([#3670](https://github.com/diamcircle/go/pull/3670)). Note that taking advantage of this feature requires [Diamcircle-Core v17.1.0](https://github.com/diamcircle/diamcircle-core/releases/tag/v17.1.0) or later.

### Bug Fixes
* The Diamcircle Core runner now parses logs from its underlying subprocess better [#3746](https://github.com/diamcircle/go/pull/3746).


## v2.0.0

This release is related to the release of [Aurora v2.3.0](https://github.com/diamcircle/go/releases/tag/aurora-v2.3.0) and introduces some breaking changes to the `ingest` package for those building their own tools.

### Breaking Changes
- Many APIs now require a `context.Context` parameter, allowing you to interact with the backends and control calls in a more finely-controlled manner. This includes the readers (`ChangeReader` et al.) as well as the backends themselves (`CaptiveDiamcircleCore` et al.).

- **`GetLedger()` always blocks** now, even for an `UnboundedRange`.

- The `CaptiveCoreBackend` now requires an all-inclusive `CaptiveCoreToml` object to configure Captive Core rather than an assortment of individual parameters. This object can be built from a TOML file (see `NewCaptiveCoreTomlFromFile`) or from parameters (see `NewCaptiveCoreToml`) as was done before.

- `LedgerTransaction.Meta` has been renamed to `UnsafeMeta` to highlight that users should be careful when interacting with it.

- Remote Captive Core no longer includes the `present` field in the ledger response JSON.

### New Features
- `NewLedgerChangeReaderFromLedgerCloseMeta` and `NewLedgerTransactionReaderFromLedgerCloseMeta` are new ways to construct readers from a particular single ledger.

### Other Changes
- The remote Captive Core client timeout has doubled.

- Captive Core now creates a temporary directory (`captive-core-...`) in the specified storage path (current directory by default) that it cleans it up on shutdown rather than in the OS's temp directory.
