# Run keystored in development

Generate the certificate and the private key pair for localhost
if you haven't done so:

```sh
cd go/exp/services/keystore
./tls/regen.sh
```
Simply choose all the default options. This will create three files:
tls/server.crt, tls/server.key, and tls/server.csr.

We will only be using `server.crt` and `server.key`.

## Install the `keystored` command:

```sh
cd go/exp/services/keystore
go install ./cmd/keystored
```

## Set up `keystore` Postgres database locally:

```sh
createdb keystore
keystored migrate up
```

You can undo the last migration by running
```sh
keystored migrate down
```

You can redo the last migration by running
```sh
keystored migrate redo
```

You can check whether there is any unapplied migrations by running
```sh
keystored migrate status
```

## Run `keystored` in development with authentication disabled:

You might want to set the `KEYSTORE_LISTENER_PORT` environment variable
for the keystored listener. Otherwise, the default value is port 8000.

```sh
keystored -tls-cert=tls/server.crt -tls-key=tls/server.key -auth=false serve
```

Before you have a valid endpoint that can handle your auth token and return a
user id in plaintext, you might want to disable authentication for testing.

## Run `keystored` in production:

There are five environment variables used for starting keystored:
`KEYSTORE_DATABASE_URL`, `DB_MAX_IDLE_CONNS`, `DB_MAX_OPEN_CONNS`,
`KEYSTORE_AUTHFORWARDING_URL`, and `KEYSTORE_LISTENER_PORT`.
* `KEYSTORE_DATABASE_URL` is required.
* `KEYSTORE_AUTHFORWARDING_URL` is required if authentication is turned on.
* `DB_MAX_IDLE_CONNS` and `DB_MAX_OPEN_CONNS` are default to 5.
* `KEYSTORE_LISTENER_PORT` is default to 8000.

```sh
keystored -tls-cert=PATH_TO_TLS_CERT -tls-key=PATH_TO_TLS_KEY serve
```

To disable authentication, you can simply add the `-auth=false` flag.

## Build docker image:

To build docker image:
```sh
cd go/services/keystore
make docker-build
```

to use custom tag:
```sh
cd go/services/keystore
TAG=my-registry.example.com/keystore:dev make docker-build
```

to push image built using the command above:
```sh
cd go/services/keystore
TAG=my-registry.example.com/keystore:dev make docker-push
```

## Logging

You can put the log messages in a designated file with the `-log-file` flag as well as determine
the log severity level with the `-log-level` flag:

```sh
keystored -log-file=PATH_TO_YOUR_LOG_FILE -log-level=[debug|info|warn|error] serve
```
