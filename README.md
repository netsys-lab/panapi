# PANAPI -- Work In Progress

[![Go](https://github.com/netsys-lab/panapi/actions/workflows/go.yml/badge.svg)](https://github.com/netsys-lab/panapi/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/netsys-lab/panapi)](https://goreportcard.com/report/github.com/netsys-lab/panapi) 
[![Go Reference](https://pkg.go.dev/badge/github.com/netsys-lab/panapi.svg)](https://pkg.go.dev/github.com/netsys-lab/panapi)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/netsys-lab/panapi/LICENSE)

PANAPI is an early [research](https://netsys.ovgu.de) implementation of a next-generation networking [API to the transport layer](https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html). The latter is currently under development in the IETF [TAPS working group](https://datatracker.ietf.org/wg/taps/about/). PANAPI is an [EU-funded](https://pointer.ngi.eu/) open-source project, that [adds support](https://dl.acm.org/doi/10.1145/3472727.3472808) for the [SCION network architecture](https://scion-architecture.net/) to a general purpose TAPS-like networking API.

## `import "panapi"` - The PANAPI Library

* [x] Simple working code example, see [examples/concurrent/README.md](examples/concurrent/README.md)
* [ ] Add more code examples

## `cmd/daemon` - The PANAPI Daemon

* [x] applications selecting SCION as transport benefit from daemon running in the backend
* [x] Graceful fallback to default behavior when daemon not running
* [ ] Create dedicated daemon README
* [x] Lua scripting examples 
  * [x] [cmd/daemon/simple.lua](cmd/daemon/simple.lua)


## Protocol support

- [x] TCP/IP support
- [x] QUIC/IP support
- [ ] UDP/IP support (_currently broken_)
- [x] QUIC/SCION support
- [ ] UDP/SCION support (_currently broken_)

## Features

### Path selection

- [x] Scriptable path selector, implementing `pan.Selector`
  - [x] working Lua Data model
  - [ ] working path ranking (_currently broken, script needs to be ported to new API_)
- [x] Central path selection Daemon

### Path quality

- [ ] Passive throughput monitoring

### Convenience features
- [ ] Different log levels

### Other
- [ ] Full test coverage
- [ ] Code Documentation

## Ported Applications
- [x] `spate` traffic generator
- [x] `concurrent` code example client/server timestamp echoing
- [ ] `http`
  - [ ] server
  - [ ] client

