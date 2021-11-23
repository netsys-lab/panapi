# PANAPI -- Work In Progress
[![Go](https://github.com/netsys-lab/panapi/actions/workflows/go.yml/badge.svg)](https://github.com/netsys-lab/panapi/actions/workflows/go.yml)

PANAPI is an early [research](https://netsys-lab.ovgu.de) implementation of a next-generation networking [API to the transport layer](https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html). The latter is currently under development in the IETF [TAPS working group](https://datatracker.ietf.org/wg/taps/about/). PANAPI is an [EU-funded](https://pointer.ngi.eu/) open-source project, that [adds support](https://dl.acm.org/doi/10.1145/3472727.3472808) for the [SCION network architecture](https://scion-architecture.net/) to a general purpose TAPS-like networking API.

## `import "panapi"` - The PANAPI Library

* [ ] Add code example here

## `cmd/daemon` - The PANAPI Daemon

* [ ] Provide further details here

* SCION-enabled applications benefit from running daemon
* Graceful fallback to default behavior when daemon not running

## Protocol support

- [x] TCP/IP support
- [x] QUIC/IP support
- [ ] UDP/IP support _currently broken_
- [x] QUIC/SCION support
- [ ] UDP/SCION support _currently broken_

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

