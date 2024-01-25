# Ticker
This project aims to provide an easy-to-deploy diamcircle ticker.

## Quick Start
This project provides a docker setup that makes it easy to get a Ticker up and running (you can
check an architecture overview [here](docs/Architecture.md)). In order to get up and running,
follow these steps:

1. Install [Docker](https://hub.docker.com/editions/community/docker-ce-desktop-mac)
2. Clone the [monorepo](https://go)
3. Build the Ticker's docker image. At the repo's root, run `$ docker build -t ticker -f services/ticker/docker/Dockerfile-dev .`
4. Run the Ticker: `$ docker run --rm -it -p "8000:8000" ticker` (you'll be asked to enter a
   PostgreSQL password)
5. After the initial setup (after the `supervisord started` message), you should be able to visit
   the two available endpoints: http://localhost:8000/markets.json and
   http://localhost:8000/assets.json

### Persisting the data
The quickstart guide creates an ephemeral database that will be deleted once the Docker image stops
running. If you wish to have a persisting ticker, you'll have to mount a volume inside of it. If
you want to do this, replace step `4` with the following steps:

1. Create a folder for the persisting data `$ mkdir /path/to/data/folder`
2. Run the ticker with the mounted folder: `$ docker run --rm -it -p "8000:8000" -v
   "/path/to/data/folder:/opt/diamcircle/postgresql" ticker` (you'll also be asked to enter a
   PostgreSQL password in the first time you run, but shouldn't happen the next time you run this
   command).
3. Voilà! After the initial setup / population is done, you should be able to visit you should be
   able to visit the two available endpoints: http://localhost:8000/markets.json and
   http://localhost:8000/assets.json

## Using the CLI
You can also test the Ticker locally, without the Docker setup. For that, you'll need a PostgreSQL
instance running. In order to build the Ticker project, follow these steps:
1. See the details in [README.md](../../../../README.md#dependencies) for installing dependencies.
2. Run `$ go run main.go --help` to see the list of available commands.
