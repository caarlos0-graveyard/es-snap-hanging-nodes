# es-snap-hanging-nodes

Long story short, sometimes, a snapshot would hang, neither finishing nor
failing, and canceling it also has no effects.

Often, its a couple of bad nodes for whatever reason.

This little tool pokes `/_snapshot/_status` and prints the list of nodes that
have snapshot shards stuck in `INIT`. Often you'll want to restart those bad
nodes.

## Install

Grab it from the latest release, or, on macs:

```shell
brew install caarlos0/tap/es-snap-hanging-nodes
```
