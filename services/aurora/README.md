# aurora
[![Build Status](https://circleci.com/gh/diamcircle/go.svg?style=shield)](https://circleci.com/gh/diamcircle/go)

aurora is the client facing API server for the [diamcircle ecosystem](https://developers.diamcircle.org/docs/start/introduction/).  It acts as the interface between [diamcircle Core](https://developers.diamcircle.org/docs/run-core-node/) and applications that want to access the diamcircle network. It allows you to submit transactions to the network, check the status of accounts, subscribe to event streams and more.

## Try it out
See aurora in action by running your own diamcircle node as part of the diamcircle [testnet](https://developers.diamcircle.org/docs/glossary/testnet/). With our Docker quick-start image, you can be running your own fully functional node in around 20 minutes. See the [Quickstart Guide](internal/docs/quickstart.md) to get up and running.

## Run a production server
If you're an administrator planning to run a production instance of aurora as part of the public diamcircle network, check out the detailed [Administration Guide](internal/docs/admin.md). It covers installation, monitoring, error scenarios and more.

## Contributing
As an open source project, development of aurora is public, and you can help! We welcome new issue reports, documentation and bug fixes, and contributions that further the project roadmap. The [Development Guide](internal/docs/developing.md) will show you how to build aurora, see what's going on behind the scenes, and set up an effective develop-test-push cycle so that you can get your work incorporated quickly.
