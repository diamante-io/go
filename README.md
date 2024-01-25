<div align="center">
<a href="https://diamcircle.org"><img alt="diamcircle" src="https://github.com/diamcircle/.github/raw/master/diamcircle-logo.png" width="558" /></a>
<br/>
<strong>Creating equitable access to the global financial system</strong>
<h1>diamcircle Go Monorepo</h1>
</div>
<p align="center">
<a href="https://circleci.com/gh/diamcircle/go"><img alt="Build Status" src="https://circleci.com/gh/diamcircle/go.svg?style=shield" /></a>
<a href="https://godoc.org/go"><img alt="GoDoc" src="https://godoc.org/go?status.svg" /></a>
<a href="https://goreportcard.com/report/go"><img alt="Go Report Card" src="https://goreportcard.com/badge/go" /></a>
</p>

This repo is the home for all of the public Go code produced by the [diamcircle Development Foundation].

This repo contains various tools and services that you can use and deploy, as well as the SDK you can use to develop applications that integrate with the diamcircle network.

## Package Index

* [aurora Server](services/aurora): Full-featured API server for diamcircle network
* [Go aurora SDK - auroraclient](clients/auroraclient): Client for aurora server (queries and transaction submission)
* [Go aurora SDK - txnbuild](txnbuild): Construct diamcircle transactions and operations
* [Ticker](services/ticker): An API server that provides statistics about assets and markets on the diamcircle network
* [Keystore](services/keystore): An API server that is used to store and manage encrypted keys for diamcircle client applications
* Servers for Anchors & Financial Institutions
  * [Compliance Server](services/compliance): Allows financial institutions to exchange KYC information
  * [Federation Server](services/federation): Allows organizations to provide addresses for users (`jane*examplebank.com`)

## Dependencies

This repository is officially supported on the last two releases of Go.

It depends on a [number of external dependencies](./go.mod), and uses Go [Modules](https://github.com/golang/go/wiki/Modules) to manage them. Running any `go` command will automatically download dependencies required for that operation.

You can choose to checkout this repository into a [GOPATH](https://github.com/golang/go/wiki/GOPATH) or into any directory.

## Directory Layout

In addition to the other top-level packages, there are a few special directories that contain specific types of packages:

* **clients** contains packages that provide client packages to the various diamcircle services.
* **exp** contains experimental packages.  Use at your own risk.
* **handlers** contains packages that provide pluggable implementors of `http.Handler` that make it easier to incorporate portions of the diamcircle protocol into your own http server. 
* **support** contains packages that are not intended for consumption outside of diamcircle's other packages.  Packages that provide common infrastructure for use in our services and tools should go here, such as `db` or `log`. 
* **support/scripts** contains single-file go programs and bash scripts used to support the development of this repo. 
* **services** contains packages that compile to applications that are long-running processes (such as API servers).
* **tools** contains packages that compile to command line applications.

Each of these directories have their own README file that explain further the nature of their contents.

### Other packages

In addition to the packages described above, this repository contains various packages related to working with the diamcircle network from a go program.  It's recommended that you use [godoc](https://godoc.org/go#pkg-subdirectories) to browse the documentation for each.


## Package source layout

While much of the code in individual packages is organized based upon different developers' personal preferences, many of the packages follow a simple convention for organizing the declarations inside of a package that aim to aid in your ability to find code.

In each package, there may be one or more of a set of common files:

- *errors.go*: This file should contains declarations (both types and vars) for errors that are used by the package.
- *example_test.go*: This file should contains example tests, as described at https://blog.golang.org/examples.
- *main.go/internal.go* (**deprecated**): Older packages may have a `main.go` (public symbols) or `internal.go` (private symbols).  These files contain, respectively, the exported and unexported vars, consts, types and funcs for the package. New packages do not follow this pattern, and instead follow the standard Go convention to co-locate structs and their methods in the same files. 
- *main.go* (**new convention**): If present, this file contains a `main` function as part of an executable `main` package.

In addition to the above files, a package often has files that contains code that is specific to one declared type.  This file uses the snake case form of the type name (for example `loggly_hook.go` would correspond to the type `LogglyHook`).  This file should contain method declarations, interface implementation assertions and any other declarations that are tied solely to that type.

Each non-test file can have a test counterpart like normal, whose name ends with `_test.go`.  The common files described above also have their own test counterparts... for example `internal_test.go` should contains tests that test unexported behavior and more commonly test helpers that are unexported.

Generally, file contents are sorted by exported/unexported, then declaration type  (ordered as consts, vars, types, then funcs), then finally alphabetically.

### Test helpers

Often, we provide test packages that aid in the creation of tests that interact with our other packages.  For example, the `support/db` package has the `support/db/dbtest` package underneath it that contains elements that make it easier to test code that accesses a SQL database.  We've found that this pattern of having a separate test package maximizes flexibility and simplifies package dependencies.

### Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for more details.

### Developing

See [DEVELOPING.md](DEVELOPING.md) for helpful instructions for getting started developing code in this repository.

[diamcircle Development Foundation]: https://diamcircle.org
