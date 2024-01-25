---
title: Overview
---

aurora is an API server for the diamcircle ecosystem.  It acts as the interface between [diamcircle-core](https://github.com/diamcircle/diamcircle-core) and applications that want to access the diamcircle network. It allows you to submit transactions to the network, check the status of accounts, subscribe to event streams, etc. See [an overview of the diamcircle ecosystem](https://www.diamcircle.org/developers/guides/) for details of where aurora fits in.

aurora provides a RESTful API to allow client applications to interact with the diamcircle network. You can communicate with aurora using cURL or just your web browser. However, if you're building a client application, you'll likely want to use a diamcircle SDK in the language of your client.
SDF provides a [JavaScript SDK](https://www.diamcircle.org/developers/js-diamcircle-sdk/reference/index.html) for clients to use to interact with aurora.

SDF runs a instance of aurora that is connected to the test net: [https://aurora-testnet.diamcircle.org/](https://aurora-testnet.diamcircle.org/) and one that is connected to the public diamcircle network:
[https://aurora.diamcircle.org/](https://aurora.diamcircle.org/).

## Libraries

SDF maintained libraries:<br />
- [JavaScript](https://github.com/diamcircle/js-diamcircle-sdk)
- [Go](https://go/tree/master/clients/auroraclient)
- [Java](https://github.com/diamcircle/java-diamcircle-sdk)

Community maintained libraries for interacting with aurora in other languages:<br>
- [Python](https://github.com/diamcircleCN/py-diamcircle-base)
- [C# .NET Core 2.x](https://github.com/elucidsoft/dotnetcore-diamcircle-sdk)
- [Ruby](https://github.com/astroband/ruby-diamcircle-sdk)
- [iOS and macOS](https://github.com/Soneso/diamcircle-ios-mac-sdk)
- [Scala SDK](https://github.com/synesso/scala-diamcircle-sdk)
- [C++ SDK](https://github.com/bnogalm/diamcircleQtSDK)
