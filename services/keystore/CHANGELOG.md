## Unreleased

- Dropped support for Go 1.12.
* Dropped support for Go 1.13.

## [v1.2.0] - 2019-11-20

- Add `ReadTimeout` to Keystore HTTP server configuration to fix potential DoS vector.
- Dropped support for Go 1.10, 1.11.

## [v1.1.0] - 2019-08-21

Keystore has new interface for managing keys blob.
Please refer to the [spec](https://github.com/diamcircle/go/blob/bcaf3d55229df822b155442633adc230294588b4/services/keystore/spec.md) for the new changes.
Note that the data you previously store will be wiped out. Be sure to save the
data that's important to you before upgrading to this version.

## [v1.0.0] - 2019-06-18

Initial release of the keystore.
