# webauth

This is a [SEP-10] Web Authentication implementation based on SEP-10 v3.2.1
that requires a user to prove they possess a signing key(s) that meets the high
threshold for an account, i.e. they have the ability to perform any high
threshold operation on the given account. If an account does not exist it may
be optionally verified using the account's master key.

SEP-10 defines an endpoint for authenticating a user in possession of a diamcircle
account using their diamcircle account as credentials. This implementation is a
standalone microservice that implements the minimum requirements as defined by
the SEP-10 protocol and will be adapted as the protocol evolves.

This implementation is not polished and is still experimental.
Running this implementation in production is not recommended.

## Usage

```
$ webauth --help
SEP-10 Web Authentication Server

Usage:
  webauth [command] [flags]
  webauth [command]

Available Commands:
  genjwk      Generate a JSON Web Key (ECDSA/ES256) for JWT issuing
  serve       Run the SEP-10 Web Authentication server

Use "webauth [command] --help" for more information about a command.
```

## Usage: Serve

```
$ webauth serve --help
Run the SEP-10 Web Authentication server

Usage:
  webauth serve [flags]

Flags:
      --allow-accounts-that-do-not-exist   Allow accounts that do not exist (ALLOW_ACCOUNTS_THAT_DO_NOT_EXIST)
      --auth-home-domain string            Home domain(s) of the service(s) requiring SEP-10 authentication comma separated (first domain is the default domain) (AUTH_HOME_DOMAIN)
      --challenge-expires-in int           The time period in seconds after which the challenge transaction expires (CHALLENGE_EXPIRES_IN) (default 300)
      --domain string                      Domain that this service is hosted at (DOMAIN)
      --aurora-url string                 aurora URL used for looking up account details (aurora_URL) (default "https://aurora-testnet.diamcircle.org/")
      --jwk string                         JSON Web Key (JWK) used for signing JWTs (if the key is an asymmetric key that has separate public and private key, the JWK must contain the private key) (JWK)
      --jwt-expires-in int                 The time period in seconds after which the JWT expires (JWT_EXPIRES_IN) (default 300)
      --jwt-issuer string                  The issuer to set in the JWT iss claim (JWT_ISSUER)
      --network-passphrase string          Network passphrase of the diamcircle network transactions should be signed for (NETWORK_PASSPHRASE) (default "Test SDF Network ; September 2015")
      --port int                           Port to listen and serve on (PORT) (default 8000)
      --signing-key string                 diamcircle signing key(s) used for signing transactions comma separated (first key is used for signing, others used for verifying challenges) (SIGNING_KEY)
```

[SEP-10]: https://github.com/diamcircle/diamcircle-protocol/blob/28c636b4ef5074ca0c3d46bbe9bf0f3f38095233/ecosystem/sep-0010.md
