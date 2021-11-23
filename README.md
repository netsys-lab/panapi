# PANAPI -- Work In Progress
[![Go](https://github.com/netsys-lab/panapi/actions/workflows/go.yml/badge.svg)](https://github.com/netsys-lab/panapi/actions/workflows/go.yml)

## `cmd/daemon` - The PANAPI Daemon

* SCION-enabled applications benefit from running daemon
* Graceful fallback to default behavior when daemon not running

## Protocol support

- [x] TCP/IP support
- [x] QUIC/IP support
- [ ] UDP/IP support
- [x] QUIC/SCION support
- [ ] UDP/SCION support

## Features

### Path selection

- [x] Scriptable path selector, implementing `pan.Selector`
  - [x] working Lua Data model
  - [x] working path ranking
- [x] Central path selection Daemon

### Path quality

- [ ] Passive throughput monitoring

### Convenience features
- [ ] Different log levels


## Ported Applications
- [x] `spate` traffic generator
- [x] `concurrent` code example client/server timestamp echoing
- [ ] `http`
  - [ ] server
  - [ ] client
