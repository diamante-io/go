# auroraclient


`auroraclient` is a [Diamcircle Go SDK](https://developers.diamcircle.org/api/) package that provides client access to a aurora server. It supports all endpoints exposed by the [aurora API](https://developers.diamcircle.org/api/introduction/).

This project is maintained by the Diamcircle Development Foundation.

## Getting Started
This library is aimed at developers building Go applications that interact with the [Diamcircle network](https://www.diamcircle.org/). It allows users to query the network and submit transactions to the network. The recommended transaction builder for Go programmers is [txnbuild](https://github.com/diamcircle/go/tree/master/txnbuild). Together, these two libraries provide a complete Diamcircle SDK.

* The [auroraclient API reference](https://godoc.org/github.com/diamcircle/go/clients/auroraclient).
* The [txnbuild API reference](https://godoc.org/github.com/diamcircle/go/txnbuild).

### Prerequisites
* Go (this repository is officially supported on the last two releases of Go)
* [Modules](https://github.com/golang/go/wiki/Modules) to manage dependencies

### Installing
* `go get github.com/diamcircle/go/clients/auroraclient`

### Usage

``` golang
    ...
    import hClient "github.com/diamcircle/go/clients/auroraclient"
    ...

    // Use the default pubnet client
    client := hClient.DefaultPublicNetClient

    // Create an account request
    accountRequest := hClient.AccountRequest{AccountID: "GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU"}

    // Load the account detail from the network
    account, err := client.AccountDetail(accountRequest)
    if err != nil {
        fmt.Println(err)
        return
    }
    // Account contains information about the diamcircle account
    fmt.Print(account)
```
For more examples, refer to the [documentation](https://godoc.org/github.com/diamcircle/go/clients/auroraclient).

## Running the tests
Run the unit tests from the package directory: `go test`

## Contributing
Please read [Code of Conduct](https://github.com/diamcircle/.github/blob/master/CODE_OF_CONDUCT.md) to understand this project's communication rules.

To submit improvements and fixes to this library, please see [CONTRIBUTING](../CONTRIBUTING.md).

## License
This project is licensed under the Apache License - see the [LICENSE](../../LICENSE-APACHE.txt) file for details.
