# @block-api/block-node

```text
█▄▄ █░░ █▀█ █▀▀ █▄▀ ▄▄ █▄░█ █▀█ █▀▄ █▀▀
█▄█ █▄▄ █▄█ █▄▄ █░█ ░░ █░▀█ █▄█ █▄▀ ██▄
```

<p align="center" width="100%">
<img src="./docs/images/gopher-golang.png" alt="golang gopher"/>
</p>

> :warning: **This project is in development**: Do not use in production environment

- [Overview](#overview)
- [Configuration](#configuration)
- [How it works](#how-it-works)
  - [P2P](/docs/p2p.md)
  - [Redis/Nats](/docs/redis_nats.md)
  - [Database](/docs/database.md)
  - [File Storage](/docs/file_storage.md)
- [TBD](#tbd)

## Overview

**block-node** is open source, framework written in Go language.

Main purpose of this project is to provide communication layer for application over protocols/services listed below:

- [ ] P2P
- [ ] TCP
- [ ] Redis
- [ ] NATS

This should allow you to build decentralized applications (Web 3.0) as well as microservices (Web 2.0).

In the future it will provide out of the box support for couple of databases as well:

- [ ] LevelDB
- [ ] MongoDB
- [ ] CouchDB

![image](./docs/images/block_node_web20.png)

![image](./docs/images/block_node_web30.png)

## Configuration

In root directory you can find `config.example.yml` file which includes available options to configure.

## TBD

This section presents features to be discussed if they should be implemented in the future:

- Generation of ETH wallet for node, pub/priv keys, for identification and to sign `TransportPocket`
- WebSockets
- HTTP (eg. for REST API)
- Command like "git" to commit files to node - ak'a file storage
  - compressing files before sending to node
  - options if distributed and kept on node itself or on separate "storage" type node, or on cloud storage like s3/azure blob

## Contact

If you have any questions or ideas feel free to reach us out on [twitter](https://twitter.com/blockapi_dev).
