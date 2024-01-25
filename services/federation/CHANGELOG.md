# Changelog

All notable changes to this project will be documented in this
file.  This project adheres to [Semantic Versioning](http://semver.org/).

As this project is pre 1.0, breaking changes may happen for minor version
bumps.  A breaking change will get clearly notified in this log.

## Unreleased

* Dropped support for Go 1.12.
* Dropped support for Go 1.13.
* Log User-Agent header in request logs.

## [v0.3.0] - 2019-11-20

### Changed

- BREAKING CHANGE: MySQL is no longer supported. To migrate your data to postgresql use any of the tools provided [here](https://wiki.postgresql.org/wiki/Converting_from_other_Databases_to_PostgreSQL#MySQL).
- Add `ReadTimeout` to HTTP server configuration to fix potential DoS vector.
- Dropped support for Go 1.10, 1.11.

## [v0.2.1] - 2017-02-14

### Changed

- BREAKING CHANGE: The `url` database configuration has been renamed to `dsn` to more accurately reflect its content.

### Fixed

- TLS support re-enabled.

### Added

- Reverse federation is now optional.
- Logging:  http requests will be logged at the "Info" log level

## [v0.2.0] - 2016-08-17

Initial release after import from https://github.com/diamcircle/federation

[Unreleased]: https://go/compare/federation-v0.2.0...master
