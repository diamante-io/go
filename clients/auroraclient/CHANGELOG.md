# Changelog

All notable changes to this project will be documented in this
file.  This project adheres to [Semantic Versioning](http://semver.org/).


## Unreleased


## [8.0.0-beta.0](https://github.com/diamcircle/go/releases/tag/auroraclient-v8.0.0-beta.0) - 2021-10-04

**This release adds support for Protocol 18.**

* The restriction that `Fund` can only be called on the DefaultTestNetClient has been removed. Any `auroraclient.Client` may now call Fund. Aurora instances not supporting funding will error with a resource not found error.
* Change `AccountRequest` to accept `Sponsor` and `LiquidityPool` filters
* Change `EffectRequest`, `TransactionRequest`, and `OperationRequest` to accept a `ForLiquidityPool` filter
* Change `TradeRequest` to accept both a `ForLiquidityPool` filter or a `TradeType` filter
* Add `LiquidityPoolsRequest` for getting details about liquidity pools, with an optional `Reserves` field to filter by pools' reserve asset(s).
* Add `LiquidityPoolRequest` for getting details about a _specific_ liquidity pool via the `LiquidityPoolID` filter.


## [v7.1.1](https://github.com/diamcircle/go/releases/tag/auroraclient-v7.1.1) - 2021-06-25

* Added transaction and operation result codes to the auroraclient.Error string for easy glancing at string only errors for underlying cause.
* Fix bug in the transaction submission where requests with large transaction payloads fail with an HTTP 414 URI Too Long error ([#3643](https://github.com/diamcircle/go/pull/3643)).
* Fix a data race in `Client.fixAuroraURL`([#3690](https://github.com/diamcircle/go/pull/3690)).
* Fix bug which occurs when parsing operations involving muxed accounts ([#3722](https://github.com/diamcircle/go/pull/3722)).

## [v7.0.0](https://github.com/diamcircle/go/releases/tag/auroraclient-v7.0.0) - 2021-05-15

None

## [v6.0.0](https://github.com/diamcircle/go/releases/tag/auroraclient-v6.0.0) - 2021-02-22

None

## [v5.0.0](https://github.com/diamcircle/go/releases/tag/auroraclient-v5.0.0) - 2020-11-12

None

## [v4.2.0](https://github.com/diamcircle/go/releases/tag/auroraclient-v4.2.0) - 2020-11-11

None

## [v4.1.0](https://github.com/diamcircle/go/releases/tag/auroraclient-v4.1.0) - 2020-10-16

None

## [v4.0.1](https://github.com/diamcircle/go/releases/tag/auroraclient-v4.0.1) - 2020-10-02

None

## [v4.0.0](https://github.com/diamcircle/go/releases/tag/auroraclient-v4.0.0) - 2020-09-29

Added new client methods and effects supporting [Protocol 14](https://github.com/diamcircle/go/issues/3035).

* New client methods
  * `ClaimableBalances(req ClaimableBalanceRequest)` - returns details about available claimable balances, possibly filtered to a specific sponsor or other parameters.
  * `ClaimableBalance(balanceID string)` - returns details about a *specific*, unique claimable balance.
* New effects:
  * `ClaimableBalance{Created,Updated,Removed}`
  * `ClaimabeBalanceSponsorship{Created,Updated,Removed}`
  * `AccountSponsorship{Created,Updated,Removed}`
  * `TrustlineSponsorship{Created,Updated,Removed}`
  * `Data{Created,Updated,Removed}`
  * `DataSponsorship{Created,Updated,Removed}`
  * `SignerSponsorship{Created,Updated,Removed}`
* Removed JSON variant of `GET /metrics`, both in the server and client code. It's using Prometheus format by default now.
* Added `NextAccountsPage`.
* Fixed `Fund` function that consistently errored.
* Added support for Go 1.15.

### Breaking changes

* Dropped support for Go 1.13.

## [v3.0.0](https://github.com/diamcircle/go/releases/tag/auroraclient-v3.0.0) - 2020-04-28

### Breaking changes

- The type for the following fields in the `Transaction` struct have changed from `int32` to `int64`:
  - `FeeCharged`
  - `MaxFee`
- The `TransactionSuccess` Aurora response has been removed. Instead, all submit transaction functions return with a full Aurora `Transaction` response on success.
- The `GetSequenceNumber()` and `IncrementSequenceNumber()` functions on the `Account` struct now return `int64` values instead of `xdr.SequenceNumber` values.

### Add

- Add `IsNotFoundError`
- Add `client.SubmitFeeBumpTransaction` and `client.SubmitFeeBumpTransactionWithOptions` for submitting [fee bump transactions](https://github.com/diamcircle/diamcircle-protocol/blob/master/core/cap-0015.md). Note that fee bump transactions will only be accepted by Diamcircle Core once Protocol 13 is enabled.

### Updates

- To support [CAP0018](https://github.com/diamcircle/diamcircle-protocol/blob/master/core/cap-0018.md) Fine-Grained Control of Authorization:
  - There is a new effect `TrustlineAuthorizedToMaintainLiabilities` which occurs when a trustline has been authorized to maintain liabilities.
  - The `AuthorizeToMaintainLiabilities` boolean field was added to the `AllowTrust` operation struct.
- These fields were added to the `Transaction` struct to support fee bump transactions:
  - `FeeAccount` (the account which paid the transaction fee)
  - `FeeBumpTransaction` (only present in Protocol 13 fee bump transactions)
  - `InnerTransaction` (only present in Protocol 13 fee bump transactions).
- `Transaction` has a new `MemoBytes` field which is populated when `MemoType` is equal to `text`. `MemoBytes` stores the base 64 encoding of the memo bytes set in the transaction envelope.
- Fixed a bug where AuroraTimeOut has misleading units of time by:
  - Removed AuroraTimeOut (seconds)
  - Added AuroraTimeout (nanoseconds)

### Remove

- Dropped support for Go 1.12.

## [v2.2.0](https://github.com/diamcircle/go/releases/tag/auroraclient-v2.2.0) - 2020-03-26

### Added

- Add `client.SubmitTransactionWithOptions` with support for [SEP0029](https://github.com/diamcircle/diamcircle-protocol/blob/master/ecosystem/sep-0029.md).
    If any of the operations included in `client.SubmitTransactionWithOptions` is of type
    `payment`, `pathPaymentStrictReceive`, `pathPaymentStrictSend`, or
    `mergeAccount`, then the SDK will load the destination account from Aurora and check if
    `config.memo_required` is set to `1` as defined in [SEP0029](https://github.com/diamcircle/diamcircle-protocol/blob/master/ecosystem/sep-0029.md).

    For performance reasons, you may choose to skip the check by setting the `SkipMemoRequiredCheck` to `true`:

	```
		client.SubmitTransactionWithOptions(tx, auroraclient.SubmitTxOpts{SkipMemoRequiredCheck: true})
	```

	Additionally, the check will be skipped automatically if the transaction includes a memo.

## Changed

-  Change `client.SubmitTransaction` to always check if memo is required.
	If you want to skip the check, call `client.SubmitTransactionWithOptions` instead.

## [v2.1.0](https://github.com/diamcircle/go/releases/tag/auroraclient-v2.1.0) - 2020-02-24

### Added

- Add `client.StrictReceivePaths` and  `client.StrictSendPaths` ([#2237](https://github.com/diamcircle/go/pull/2237)).

`client.StrictReceivePaths`:

```go
	client := auroraclient.DefaultPublicNetClient
	// Find paths for XLM->NGN
	pr := auroraclient.PathsRequest{
		DestinationAccount:     "GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU",
		DestinationAmount:      "100",
		DestinationAssetCode:   "NGN",
		DestinationAssetIssuer: "GDZST3XVCDTUJ76ZAV2HA72KYQODXXZ5PTMAPZGDHZ6CS7RO7MGG3DBM",
		DestinationAssetType:   auroraclient.AssetType4,
		SourceAccount:          "GDZST3XVCDTUJ76ZAV2HA72KYQODXXZ5PTMAPZGDHZ6CS7RO7MGG3DBM",
	}
	paths, err := client.StrictReceivePaths(pr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(paths)
```

`client.StrictSendPaths`:

```go
	client := auroraclient.DefaultPublicNetClient
	// Find paths for USD->EUR
	pr := auroraclient.StrictSendPathsRequest{
		SourceAmount:      "20",
		SourceAssetCode:   "USD",
		SourceAssetIssuer: "GDUKMGUGDZQK6YHYA5Z6AY2G4XDSZPSZ3SW5UN3ARVMO6QSRDWP5YLEX",
		SourceAssetType:   auroraclient.AssetType4,
		DestinationAssets: "EURT:GAP5LETOV6YIE62YAM56STDANPRDO7ZFDBGSNHJQIYGGKSMOZAHOOS2S",
	}
	paths, err := client.StrictSendPaths(pr)
```

- Add `client.OfferDetails` ([#2303](https://github.com/diamcircle/go/pull/2303)).

```go
	client := auroraclient.DefaultPublicNetClient
	offer, err := client.OfferDetails("2")
	if err != nil {
		fmt.Println(err)
		return
	}	
	fmt.Print(offer)
```

- Add support to `client.Offers` for the filters: `Seller`, `Selling` and `Buying` ([#2230](https://github.com/diamcircle/go/pull/2230)).
```go
	offerRequest = auroraclient.OfferRequest{
		Seller:  "GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU",
		Selling: "COP:GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU",
		Buying:  "EUR:GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU",
		Order:   auroraclient.OrderDesc,
	}	
	offers, err = client.Offers(offerRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(offers)
```
- Add `client.Accounts` ([#2229](https://github.com/diamcircle/go/pull/2229)).

This feature allows account retrieval filtering by signer or by a trustline to an asset.

```go
	client := auroraclient.DefaultPublicNetClient
	accountID := "GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU"
	accountsRequest := auroraclient.AccountsRequest{Signer: accountID}	
	account, err := client.Accounts(accountsRequest)
	if err != nil {
		fmt.Println(err)
		return
	}	
	fmt.Print(account)
```

- Add `IsNotFoundError` ([#2197](https://github.com/diamcircle/go/pull/2197)).

### Deprecated

- Make `hProtocol.FeeStats` backwards compatible with Aurora `0.24.1` and `1.0` deprecating usage of `*_accepted_fee` ([#2290](https://github.com/diamcircle/go/pull/2290)).

All the `_accepted_fee` fields were removed in Aurora 1.0, however we extended this version of the SDK to backfill the `FeeStat` struct using data from `MaxFee`. This is a temporary workaround and it will be removed in auroraclient 3.0. Please start using data from `FeeStat.MaxFee` instead.


## [v2.0.0](https://github.com/diamcircle/go/releases/tag/auroraclient-v2.0.0) - 2020-01-13

- Add custom `UnmarshalJSON()` implementations to Aurora protocol structs so `int64` fields can be parsed as JSON numbers or JSON strings
- Remove deprecated `fee_paid field` from Transaction response
- Dropped support for Go 1.10, 1.11.

## [v1.4.0](https://github.com/diamcircle/go/releases/tag/auroraclient-v1.4.0) - 2019-08-09

- Add support for querying operation endpoint with `join` parameter [#1521](https://github.com/diamcircle/go/issues/1521).
- Add support for querying previous and next trade aggregations with `Client.NextTradeAggregationsPage` and `Client.PrevTradeAggregationsPage` methods.


## [v1.3.0](https://github.com/diamcircle/go/releases/tag/auroraclient-v1.3.0) - 2019-07-08

- Transaction information returned by methods now contain new fields: `FeeCharged` and `MaxFee`. `FeePaid` is deprecated and will be removed in later versions.
- Improved unit test for `Client.FetchTimebounds` method.
- Added `Client.HomeDomainForAccount` helper method for retrieving the home domain of an account.

## [v1.2.0](https://github.com/diamcircle/go/releases/tag/auroraclient-v1.2.0) - 2019-05-16

- Added support for returning the previous and next set of pages for a aurora response; issue [#985](https://github.com/diamcircle/go/issues/985).
- Fixed bug reported in [#1254](https://github.com/diamcircle/go/issues/1254)  that causes a panic when using auroraclient in goroutines.


## [v1.1.0](https://github.com/diamcircle/go/releases/tag/auroraclient-v1.1.0) - 2019-05-02

### Added

- `Client.Root()` method for querying the root endpoint of a aurora server.
- Support for returning concrete effect types[#1217](https://github.com/diamcircle/go/pull/1217)
- Fix when no HTTP client is provided

### Changes

- `Client.Fund()` now returns `TransactionSuccess` instead of a http response pointer.

- Querying the effects endpoint now supports returning the concrete effect type for each effect. This is also supported in streaming mode. See the [docs](https://godoc.org/github.com/diamcircle/go/clients/auroraclient#Client.Effects) for examples.

## [v1.0.0](https://github.com/diamcircle/go/releases/tag/auroraclient-v1.0) - 2019-04-26

 * Initial release
