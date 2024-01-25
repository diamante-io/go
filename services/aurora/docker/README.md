# Overview

Files related to docker and docker-compose
* `Dockerfile` and `Makefile` - used to build the official, package-based docker image for diamcircle-aurora
* `Dockerfile.dev` - used with docker-compose

# Running Diamcircle with Docker Compose

## Dependencies

The only dependency you will need to install is [Docker](https://www.docker.com/products/docker-desktop).

## Start script

[start.sh](./start.sh) will setup the env file and run docker-compose to start the Diamcircle docker containers. Feel free to use this script, otherwise continue with the next two steps.

The script takes one optional parameter which configures the Diamcircle network used by the docker containers. If no parameter is supplied, the containers will run on the Diamcircle test network.

`./start.sh pubnet` will run the containers on the Diamcircle public network.

`./start.sh standalone` will run the containers on a private standalone Diamcircle network.

## Run docker-compose

Run the following command to start all the Diamcircle docker containers:

```
docker-compose up -d --build
```

Aurora will be exposed on port 8000. Diamcircle Core will be exposed on port 11626. The Diamcircle Core postgres instance will be exposed on port 5641.
The Aurora postgres instance will be exposed on port 5432.

## Swapping in a local service

If you're developing a service locally you may want to run that service locally while also being able to interact with the other Diamcircle components running in Docker. You can do that by stopping the container corresponding to the service you're developing.

For example, to run Aurora locally from source, you would perform the following steps:

```
# stop aurora in docker-compose
docker-compose stop aurora
```

Now you can run aurora locally in vscode using the following configuration:
```
    {
        "name": "Launch",
        "type": "go",
        "request": "launch",
        "mode": "debug",
        "remotePath": "",
        "port": 2345,
        "host": "127.0.0.1",
        "program": "${workspaceRoot}/services/aurora/main.go",
        "env": {
            "DATABASE_URL": "postgres://postgres@localhost:5432/aurora?sslmode=disable",
            "DIAMNET_CORE_DATABASE_URL": "postgres://postgres:mysecretpassword@localhost:5641/diamcircle?sslmode=disable",
            "NETWORK_PASSPHRASE": "Test SDF Network ; September 2015",
            "DIAMNET_CORE_URL": "http://localhost:11626",
            "INGEST": "true",
        },
        "args": []
    }
```

Similarly, to run Diamcircle core locally from source and have it interact with Aurora in docker, all you need to do is run `docker-compose stop core` before running Diamcircle core from source.

## Connecting to the Diamcircle Public Network

By default, the Docker Compose file configures Diamcircle Core to connect to the Diamcircle test network. If you would like to run the docker containers on the
Diamcircle public network, run `docker-compose -f docker-compose.yml -f docker-compose.pubnet.yml up -d --build`. 

To run the containers on a private stand-alone network, run `docker-compose -f docker-compose.yml -f docker-compose.standalone.yml up -d --build`.
When you run Diamcircle Core on a private stand-alone network, an account will be created which will hold 100 billion Lumens.
The seed for the account will be emitted in the Diamcircle Core logs:

```
2020-04-22T18:39:19.248 GD5KD [Ledger INFO] Root account seed: SC5O7VZUXDJ6JBDSZ74DSERXL7W3Y5LTOAMRF7RQRL3TAGAPS7LUVG3L
```

When running Aurora on a private stand-alone network, Aurora will not start ingesting until Diamcircle Core creates its first history archive snapshot. Diamcircle Core creates snapshots every 64 ledgers, which means ingestion will be delayed until ledger 64.

When you switch between different networks you will need to clear the Diamcircle Core and Diamcircle Aurora databases. You can wipe out the databases by running `docker-compose down --remove-orphans -v`.

## Using a specific version of Diamcircle Core

By default the Docker Compose file is configured to use version 18 of Protocol and Diamcircle Core. You want the Core version to be at same level as the version aurora repo expects for ingestion. You can specify optional environment variables from the command shell for stating version overrides for either the docker-compose or start.sh invocations. 

PROTOCOL_VERSION=18                              // the Diamcircle Protocol version number
CORE_IMAGE=diamcircle/diamcircle-core:18               // the docker hub image:tag 
DIAMNET_CORE_VERSION=18.1.1-779.ef0f44b44.focal  // the apt deb package version from apt.diamcircle.org

Example:

Runs Diamcircle Protocol and Core version 18, for any mode of testnet,standalone,pubnet
```PROTOCOL_VERSION=18 CORE_IMAGE=diamcircle/diamcircle-core:18 DIAMNET_CORE_VERSION=18.1.1-779.ef0f44b44.focal ./start.sh [standalone|pubnet]```
