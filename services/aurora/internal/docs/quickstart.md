---
title: aurora Quickstart
replacement: https://developers.diamcircle.org/docs/run-api-server/quickstart/
---
## aurora Quickstart
This document describes how to quickly set up a **test** diamcircle Core + aurora node, that you can play around with to get a feel for how a diamcircle node operates. **This configuration is not secure!** It is **not** intended as a guide for production administration.

For detailed information about running aurora and diamcircle Core safely in production see the [aurora Administration Guide](admin.md) and the [diamcircle Core Administration Guide](https://www.diamcircle.org/developers/diamcircle-core/software/admin.html).

If you're ready to roll up your sleeves and dig into the code, check out the [Developer Guide](developing.md).

### Install and run the Quickstart Docker Image
The fastest way to get up and running is using the [diamcircle Quickstart Docker Image](https://github.com/diamcircle/docker-diamcircle-core-aurora). This is a Docker container that provides both `diamcircle-core` and `aurora`, pre-configured for testing.

1. Install [Docker](https://www.docker.com/get-started).
2. Verify your Docker installation works: `docker run hello-world`
3. Create a local directory that the container can use to record state. This is helpful because it can take a few minutes to sync a new `diamcircle-core` with enough data for testing, and because it allows you to inspect and modify the configuration if needed. Here, we create a directory called `diamcircle` to use as the persistent volume:
`cd $HOME; mkdir diamcircle`
4. Download and run the diamcircle Quickstart container, replacing `USER` with your username:

```bash
docker run --rm -it -p "8000:8000" -p "11626:11626" -p "11625:11625" -p"8002:5432" -v $HOME/diamcircle:/opt/diamcircle --name diamcircle diamcircle/quickstart --testnet
```

You can check out diamcircle Core status by browsing to http://localhost:11626.

You can check out your aurora instance by browsing to http://localhost:8000.

You can tail logs within the container to see what's going on behind the scenes:
```bash
docker exec -it diamcircle /bin/bash
supervisorctl tail -f diamcircle-core
supervisorctl tail -f aurora stderr
```

On a modern laptop this test setup takes about 15 minutes to synchronise with the last couple of days of testnet ledgers. At that point aurora will be available for querying. 

See the [Quickstart Docker Image](https://github.com/diamcircle/docker-diamcircle-core-aurora) documentation for more details, and alternative ways to run the container. 

You can test your aurora instance with a query like: http://localhost:8000/transactions?cursor=&limit=10&order=asc. Use the [diamcircle Laboratory](https://www.diamcircle.org/laboratory/) to craft other queries to try out,
and read about the available endpoints and see examples in the [aurora API reference](https://www.diamcircle.org/developers/aurora/reference/).

