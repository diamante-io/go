# federation server


Go implementation of [Federation](https://developers.diamcircle.org/docs/glossary/federation/) protocol server. This federation server is designed to be dropped in to your existing infrastructure. It can be configured to pull the data it needs out of your existing DB.

## Downloading the server

[Prebuilt binaries](https://github.com/diamcircle/go/releases) of the federation server are available on the [releases page](https://github.com/diamcircle/go/releases).

## Config

By default this server uses a config file named `federation.cfg` in the current working directory. This configuration file should be [TOML](https://github.com/toml-lang/toml) and the following fields are supported:

* `port` - server listening port
* `database`
  * `type` - database type (sqlite3, postgres)
  * `dsn` - The DSN (data source name) used to connect to the database connection.  This value should be appropriate for the database type chosen.
    * for `postgres`: `postgres://user:password@host/dbname?sslmode=sslmode` ([more info](https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters))
* `queries`
  * `federation` - Implementation dependent query to fetch federation results, should return either 1 or 3 columns. These columns should be labeled `id`,`memo`,`memo_type`. Memo and memo_type are optional - see [Federation](https://developers.diamcircle.org/docs/glossary/federation/) docs for more detail).  When executed, this query will be provided with two input parameters, the first will be the name portion of a diamcircle address and the second will be the domain portion of a diamcircle address.  For example, a request for `scott*diamcircle.org` would trigger a query with two input parameters, `scott` and `diamcircle.org` respectively. 
  * `reverse-federation` - A SQL query to fetch reverse federation results that should return two columns, labeled `name` and `domain`.   When executed, this query will be provided with one input parameter, a [diamcircle account ID](https://developers.diamcircle.org/docs/glossary/accounts/#account-id) used to lookup the name and domain mapping.

    If reverse-lookup isn't supported (e.g. you have a single Diamcircle account for all users), leave this entry out.

* `tls` (only when running HTTPS server)
  * `certificate-file` - a file containing a certificate
  * `private-key-file` - a file containing a matching private key

`memo_type` should be one of the following:
* `id` - then `memo` field should contain unsigned 64-bit integer, please note that this value will be converted to integer so the field should be an integer or a string representing an integer,
* `text` - then `memo` field should contain string, up to 28 characters.
* `hash` - then `memo` field should contain string that is 32bytes base64 encoded.

## Example `federation.cfg`
In this section you can find config examples for the two main ways of setting up a federation server.

### #1: Every user has their own Diamcircle account

In the case that every user owns a Diamcircle account, you don't need `memo`. You can simply return `id` based on the username. Your `queries` section could look like this:

```toml
port = 8000

[database]
type = "postgres"
dsn = "postgres://user:password@host/dbname?sslmode=sslmode"

[queries]
federation = "SELECT account_id as id FROM Users WHERE username = ? AND domain = ?"
reverse-federation = "SELECT username as name, domain FROM Users WHERE account_id = ?"
```


### #2: Single Diamcircle account for all incoming transactions

If you have a single Diamcircle account for all incoming transactions you need to use `memo` to check which internal account should receive the payment.

Let's say that your Diamcircle account ID is: `GD6WU64OEP5C4LRBH6NK3MHYIA2ADN6K6II6EXPNVUR3ERBXT4AN4ACD` and every user has an `id` and `username` in your database. Then your `queries` section could look like this:

```toml
port = 8000

[database]
type = "postgres"
dsn = "postgres://user:password@host/dbname?sslmode=sslmode"

[queries]
federation = "SELECT username as memo, 'text' as memo_type, 'GD6WU64OEP5C4LRBH6NK3MHYIA2ADN6K6II6EXPNVUR3ERBXT4AN4ACD' as id FROM Users WHERE username = ? AND domain = ?"
# No entry for `reverse-federation` since a reverse-lookup isn't possible
```

## Providing federation for a single domain

In the event that your organization only wants to offer federation for a single domain, a little bit of trickery can be used to configure your queries to satisfy this use case.  For example, let's say you own `acme.org` and want to provide only results for that domain.  The following example config illustrates:

```toml
federation = "SELECT username as memo, 'text' as memo_type, 'GD6WU64OEP5C4LRBH6NK3MHYIA2ADN6K6II6EXPNVUR3ERBXT4AN4ACD' as id FROM Users WHERE username = ? AND ? = 'acme.org'"
reverse-federation = "SELECT username as name, 'acme.org' as domain FROM Users WHERE account_id = ?"
```

Notice that SQL fragment `? = 'acme.org"` on the `federation` query:  It ensures the incoming query is for the correct domain.  Additionally, the `reverse-federation` query always returns `acme.org` for the domain.

## Postgresql sample

Bundled with the source code of this project is a sample configuration file and a shell script that can be used to populate a sample database.  These two items can be used to play around with the service.  See (./federation.cfg) and (./build_sample.sh).

## Usage

```
./federation [-c=CONFIGPATH]
```

## Building

This service can built from source, provided you have installed the [go tools](https://golang.org/doc/install), by issuing the following command in a terminal:

Given you have a running golang installation, you can build the server with:

```
go get -u github.com/diamcircle/go/services/federation
```

After successful completion, you should find `bin/federation` is present in your configured GOPATH.

## Running tests

```
go test
```
