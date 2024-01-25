# Ticker API Documentation

## Market (Ticker) Data
Provides trade data about each trade pair within the last 7-day period. Asset pairs that did not have any activity in the last 7 days are omitted from the response.

Assets from different issuers but with the same code are aggregated, so trades between, for instance:
- `native` and `BTC:GDT3ZKQZXXHDPJUKNHUMANMNIT4JWSUYXUGN7EQZDVXBO7NPNFVFPBAK`
- `native` and `BTC:GATEMHCCKCY67ZUCKTROYN24ZYT5GK4EQZ65JJLDHKHRUZI3EUEKMTCH`

are aggregated in the `XLM_BTC` pair.

### Trade Pairs

Trade pairs are ordered `<Counter>_<Base>`.

Example:

The pair `XLM_ZZZ` has the `XLM` as the counter currency and `ZZZ` as the base. For that pair if the API returns a `close` value of `2`, then the last trade for the pair was `2 XLM` exchanged for `1 ZZZ`.

### Response Fields

* `generated_at`: UNIX timestamp of when data was generated
* `generated_at_rfc3339 `: RFC 3339 formatted string of when data was generated
* `name`: name of the trade pair
* `base_volume`: accumulated amount of base traded in the last 24h
* `counter_volume`: accumulated amount of counter traded in the last 24h
* `trade_count`: number of trades in the last 24h
* `open`: open price in the last 24h period
* `low`: lowest price in the last 24h
* `high`: highest price in the last 24h
* `change`: price difference between open and low in the last 24h
* `base_volume_7d`: accumulated amount of base traded in the last 7 days
* `counter_volume_7d`: accumulated amount of counter traded in the last 7 days
* `trade_count_7d`: number of trades in the last 7 days
* `open_7d`: open price in the last 7-day period
* `low_7d`: lowest price in the last 7 days
* `high_7d`: highest price in the last 7 days
* `change_7d`: price difference between open and low in the last 7 days
* `price`: (DEPRECATED) price of the most recent trade in this market
* `close`: price of the most recent trade in this market
* `close_time`: ledger close time of the most recent trade in this market
* `bid_count`: number of open bids on order book
* `bid_volume`: volume of open bids on order book
* `bid_max`: maximum open bid price on order book
* `ask_count`: number of open asks on order book
* `ask_volume`: volume of open asks on order book
* `ask_min`: minimum asked price on order book
* `spread`: spread between bid_max an ask_min
* `spread_mid_point`: spread mid point

### Example
#### Endpoint
GET `https://ticker.diamcircle.org/markets.json`
#### Response (application/json)

```json
{
    "generated_at": 1556828634778,
    "generated_at_rfc3339": "2019-05-02T17:23:54-03:00",
    "pairs": [
        {
            "name": "ABDT_DOP",
            "base_volume": 27933.1306978,
            "counter_volume": 703779.0492835,
            "trade_count": 73,
            "open": 0.03987601218950153,
            "low": 0.038593480963638155,
            "high": 0.03989875591737053,
            "change": -0.0011995715564988227,
            "base_volume_7d": 199598.77306629982,
            "counter_volume_7d": 5004887.537185903,
            "trade_count_7d": 488,
            "open_7d": 0.03988668687332845,
            "low_7d": 0.038593480963638155,
            "high_7d": 0.04145936964569084,
            "change_7d": -0.0012102462403257436,
            "price": 0.038676440633002705,
            "close": 0.038676440633002705,
            "close_time": "2019-05-02T12:23:57-03:00",
            "bid_count": 200,
            "bid_volume": 229694.35403809912,
            "bid_max": 25.8555333333,
            "ask_count": 36,
            "ask_volume": 149041.62309569685,
            "ask_min": 25.902828723,
            "spread": 0.0018258774053509135,
            "spread_mid_point": 25.856446272002675
        },
        {
            "name": "BTC_CNY",
            "base_volume": 0.0737282,
            "counter_volume": 2686.9835871000005,
            "trade_count": 49,
            "open": 0.0000276,
            "low": 0.0000269,
            "high": 0.0000278,
            "change": -3.9999999999999956e-7,
            "base_volume_7d": 0.37105660000000024,
            "counter_volume_7d": 13616.162691900003,
            "trade_count_7d": 285,
            "open_7d": 0.0000264,
            "low_7d": 0.0000263,
            "high_7d": 0.000028,
            "change_7d": 7.999999999999991e-7,
            "price": 0.0000272,
            "close": 0.0000272,
            "close_time": "2019-05-02T12:48:41-03:00",
            "bid_count": 27,
            "bid_volume": 22126.4118872,
            "bid_max": 36630.03663003663,
            "ask_count": 21,
            "ask_volume": 4438.404611090742,
            "ask_min": 36900.36900369004,
            "spread": 0.007326007326007345,
            "spread_mid_point": 36630.04029304029
        }
    ]
}
```
## Asset (Currency) Data
Lists all the valid assets within the Diamcircle network. The provided fields are based on the [Currency Documentation of SEP-0001](https://github.com/diamcircle/diamcircle-protocol/blob/master/ecosystem/sep-0001.md#currency-documentation) and the [Asset fields from Aurora](https://developers.diamcircle.org/api/resources/assets/).
### Response Fields

* `generated_at`: UNIX timestamp of when data was generated
* `generated_at_rfc3339 `: RFC 3339 formatted string of when data was generated
* `code`: code of the asset
* `issuer`: token issuer Diamcircle public key
* `type`: type of the asset (e.g. `native` or `credit_alphanum4`)
* `num_accounts`: the number of accounts that: 1) trust this asset and 2) where if the asset has the auth_required flag then the account is authorized to hold the asset.
* `auth_required`: an anchor must approve anyone who wants to hold its asset
* `auth_revocable`: an anchor can set the authorize flag of an existing trustline to freeze the assets held by an asset holder
* `amount`: number of units of credit issued
* `asset_controlled_by_domain`: whether the reported issuer domain controls the asset
* `is_asset_anchored`: whether the asset is anchored to another one
* `anchor_asset`: anchor asset associated (e.g. USDT is anchored to USD)
* `anchor_asset_type`: type of anchor asset
* `display_decimals`: preference for number of decimals to show when a client displays currency balance
* `name`: name of the token
* `desc`: description of the token
* `conditions`: conditions on token
* `fixed_number`: if the number of tokens issued will never change
* `max_number`: max number of tokens, if there will never be more than max_number tokens
* `is_unlimited`: if the number of tokens is dilutable at the issuer's discretion
* `redemption_instructions`: if anchored token, these are instructions to redeem the underlying asset from tokens.
* `collateral_addresses`: if this is an anchored crypto token, list of one or more public addresses that hold the assets for which you are issuing tokens.
* `collateral_address_signatures`: occasional collateral address signatures
* `countries`: countries in which the asset is available
* `status`: status of token
* `last_valid`: last the time the asset info was validated

### Example
#### Endpoint
GET `https://ticker.diamcircle.org/assets.json`

#### Response (application/json)

```json
{
	"generated_at": 1556828621410,
	"generated_at_rfc3339": "2019-05-02T17:23:41-03:00",
    "assets": [
        {
            "code": "AngelXYZ",
            "issuer": "GANZBUS4726LBT2CBJ3BGF3TP4NJT5MHZMI423NHMXFRWGO2DCBQEXYZ",
            "type": "credit_alphanum12",
            "num_accounts": 282,
            "auth_required": false,
            "auth_revocable": false,
            "amount": 4999999999.999953,
            "asset_controlled_by_domain": true,
            "anchor_asset": "",
            "anchor_asset_type": "",
            "display_decimals": 0,
            "name": "",
            "desc": "",
            "conditions": "",
            "is_asset_anchored": false,
            "fixed_number": 0,
            "max_number": 0,
            "is_unlimited": false,
            "redemption_instructions": "",
            "collateral_addresses": [],
            "collateral_address_signatures": [],
            "countries": "",
            "status": "",
            "last_valid": 1555509989002
        },
        {
            "code": "PUSH",
            "issuer": "GBB5TTFQE5KT3TEBCR7Z3FZR3R3WTVD654XL2KHKVONRIOBEI5UGOFQQ",
            "type": "credit_alphanum4",
            "num_accounts": 15,
            "auth_required": false,
            "auth_revocable": false,
            "amount": 1000000000,
            "asset_controlled_by_domain": true,
            "anchor_asset": "",
            "anchor_asset_type": "",
            "display_decimals": 2,
            "name": "Push",
            "desc": "1 PUSH token entitles you to access the push API.",
            "conditions": "Token used to access the PUSH api to send a push request to the diamcircle network.",
            "is_asset_anchored": false,
            "fixed_number": 0,
            "max_number": 0,
            "is_unlimited": false,
            "redemption_instructions": "",
            "collateral_addresses": [],
            "collateral_address_signatures": [],
            "countries": "",
            "status": "",
            "last_valid": 1555509990457
        }
    ]
}

```

## GraphQL interface
Asset, issuer, markets and ticker data can be queried through a GraphQL interface, which is also provided by the Ticker.

To explore the GraphQL queries, you can access the GraphiQL URL: https://ticker.diamcircle.org/graphiql

## Orderbook
Apart from the orderbook data provided by `markets.json`, orderbook data can be retrieved directly from Aurora. In order to retrieve `ask` and `bid` data, you have to provide the following parameters from the asset pairs:

- `selling_asset_type`: type of selling asset (e.g. `native`, `credit_alphanum4`)
- `selling_asset_code`: code of the selling asset. Omit if `selling_asset_type` = `native`
- `selling_asset_issuer`: selling asset's issuer ID. Omit if `selling_asset_type` = `native`
- `buying_asset_type`: type of buying asset (e.g. `native` or `credit_alphanum4`)
- `buying_asset_code`: code of the buying asset. Omit if `buying_asset_type` = `native`
- `buying_asset_issuer`: buying asset's issuer ID. Omit if `buying_asset_type` = `native`

The `type`, `code` and `issuer` parameters for any given asset can be found in the Ticker's `assets.json` endpoint described in the previous section.


Full documentation on Aurora's Orderbook endpoint can be found [here](https://developers.diamcircle.org/api/aggregations/order-books/).

### Example
#### Endpoint
GET `https://aurora.diamcircle.org/order_book?selling_asset_type=native&buying_asset_type=credit_alphanum4&buying_asset_code=BTC&buying_asset_issuer=GATEMHCCKCY67ZUCKTROYN24ZYT5GK4EQZ65JJLDHKHRUZI3EUEKMTCH`

#### Response (application/json)
```json
{
    "bids": [
        {
            "price_r": {
                "n": 223,
                "d": 10000000
            },
            "price": "0.0000223",
            "amount": "0.0006261"
        },
        {
            "price_r": {
                "n": 16469,
                "d": 739692077
            },
            "price": "0.0000223",
            "amount": "0.0037850"
        },
        {
            "price_r": {
                "n": 16469,
                "d": 741750702
            },
            "price": "0.0000222",
            "amount": "0.0037745"
        },
        {
            "price_r": {
                "n": 111,
                "d": 5000000
            },
            "price": "0.0000222",
            "amount": "0.0040000"
        }
    ],
    "asks": [
        {
            "price_r": {
                "n": 7,
                "d": 312500
            },
            "price": "0.0000224",
            "amount": "150.8482143"
        },
        {
            "price_r": {
                "n": 9,
                "d": 400000
            },
            "price": "0.0000225",
            "amount": "348.4311112"
        },
        {
            "price_r": {
                "n": 113,
                "d": 5000000
            },
            "price": "0.0000226",
            "amount": "335.6238939"
        }
    ],
    "base": {
        "asset_type": "native"
    },
    "counter": {
        "asset_type": "credit_alphanum4",
        "asset_code": "BTC",
        "asset_issuer": "GATEMHCCKCY67ZUCKTROYN24ZYT5GK4EQZ65JJLDHKHRUZI3EUEKMTCH"
    }
}
```
