# Bitcoin Node
A Bitcoin node implementation written from scratch in Go.

The project currently implements:

* Inbound and outbound TCP connections.
* Bitcoin version / verack handshake.
* ping / pong handling.
* getaddr / addr handling.
* DNS seed bootstrapping.
* Persistent peer address storage.
* Explicit outbound peer connections.
* Configurable listening port.
* Configurable maximum number of outbound peers.

## Architecture

                                   +----------------+
                                   |    main.go     |
                                   |                |
                                   | flags/config   |
                                   +-------+--------+
                                           |
                                           v
                                   +----------------+
                                   |      Node      |
                                   +-------+--------+
                                           |
                         +-----------------+-----------------+
                         |                 |                 |
                         v                 v                 v
                  +-------------+   +-------------+   +-------------+
                  |  Bootstrap  |   |  Listener   |   |   Outbound  |
                  +------+------+   +------+------+   +------+------+
                         |                 |                 |
              +----------+----------+      |                 |
              |                     |      |                 |
              v                     v      v                 v
       +-------------+       +-------------+          +-------------+
       | addresses   |       |  DNS Seeds  |          | Candidates  |
       |    .json    |       +------+------+          +------+------+
       +------+------+              |                        |
              |                     |                        |
              +----------+----------+                        |
                         |                                   |
                         v                                   |
                  +-------------+                            |
                  |   Address   |<---------------------------+
                  |   Manager   |
                  +------+------+ 
                         |
                  +------+------+
                  |             |
                  v             v
                Load()        Save()
                  |
                  |
                  +---------------------------------------------+
                  |                                             |
                  |                                             |
                  v                                             v
           +-------------+                               +-------------+
           |   Inbound   |                               |  Outbound   |
           |   Accept    |                               |    Dial     |
           +------+------+                               +------+------+
                  |                                             |
                  |                                             |
                  v                                             v
             +---------+                                   +---------+
             |  Peer   |                                   |  Peer   |
             |  IN     |                                   |  OUT    |
             +----+----+                                   +----+----+
                  |                                             |
                  +----------------------+----------------------+
                                         |
                                         v
                                  +--------------+
                                  |    run()     |
                                  +------+-------+
                                         |
                                         v
                                  +--------------+
                                  |  Handshake   |
                                  |              |
                                  | version      |
                                  | verack       |
                                  +------+-------+
                                         |
                                         v
                                  +--------------+
                                  | Message Loop |
                                  +------+-------+
                                         |
                  +----------------------+----------------------+
                  |                      |                      |
                  v                      v                      v
                ping                    pong                   addr
                  |                      |                      |
                  v                      v                      v
                pong             validate nonce         AddressManager
                                                               |
                                                               |
                                                               v
                                                        known addresses
                                                               |
                                                               v
                                                        persistAddresses()
                                                               |
                                                               v
                                                         addresses.json


## Usage examples

### Show help
```
go run . -help
```


### Run with default configuration
```
go run .
```

Uses the default listening port and the addresses available in addresses.json. If the file does not exist, the node resolves the DNS seeds.

### Specify the listening port
```
go run . -port 8334
```

The node listens for inbound connections on port 8334.

### Connect to a specific peer
```
go run . -connect 1.2.3.4:8333
```

The peer is added as an explicitOutboundPeer.

### Connect to multiple specific peers
```
go run . \
  -connect 1.2.3.4:8333 \
  -connect 5.6.7.8:8333 \
  -connect 9.10.11.12:8333
```

Explicit peers have priority over discovered peers.

### Limit the number of outbound peers
```
go run . -max-outbound-peers 4
```

The node will use at most 4 outbound peers.

### Combine options
```
go run . \
  -port 8334 \
  -max-outbound-peers 8 \
  -connect 1.2.3.4:8333 \
  -connect 5.6.7.8:8333
```

In this case:

The node listens on port 8334.

It allows up to 8 outbound peers.

It first attempts to connect to 1.2.3.4:8333.

It then attempts to connect to 5.6.7.8:8333.

The remaining slots are filled with discovered peers.


### Testing locally:
In two terminals execute:
```
go run . -port 8334
```

```
go run . -connect 127.0.0.1:8334 -max-outbound-peers 1
```

## Roadmap
Possible next steps include:

* Concurrent outbound connection attempts.
* Maintaining the configured number of outbound peers.
* Reconnecting when a peer disconnects.
* Better peer address selection.
* Address freshness and timestamps.
* Peer scoring / connection management.
* More Bitcoin P2P messages.
* Block header synchronization.
* Block synchronization.
* Transaction relay.
* Peer banning and misbehavior handling.
* Graceful node shutdown.
* More persistent peer metadata.