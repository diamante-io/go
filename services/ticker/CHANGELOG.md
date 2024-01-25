## Unreleased

* Dropped support for Go 1.12.
* Dropped support for Go 1.13.


## [v1.2.0] - 2019-11-20
- Add `ReadTimeout` to Ticker HTTP server configuration to fix potential DoS vector.
- Added nested `"issuer_detail"` field to `/assets.json`.
- Dropped support for Go 1.10, 1.11.


## [v1.1.0] - 2019-07-22

- Added support for running the ticker on Diamcircle's Test Network, by using the `--testnet` CLI flag.
- The ticker now retries requests to Aurora if it gets rate-limited.
- Minor bug fixes and performance improvements.


## [v1.0.0] - 2019-05-20

Initial release of the ticker.
