# PANAPI - Path-aware Networking API

[![Go](https://github.com/netsys-lab/panapi/actions/workflows/go.yml/badge.svg)](https://github.com/netsys-lab/panapi/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/netsys-lab/panapi)](https://goreportcard.com/report/github.com/netsys-lab/panapi) 
[![Go Reference](https://pkg.go.dev/badge/github.com/netsys-lab/panapi.svg)](https://pkg.go.dev/github.com/netsys-lab/panapi)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

PANAPI is an early [research](https://dl.acm.org/doi/10.1145/3472727.3472808) implementation of a next-generation networking [API to the transport layer](https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html). The latter is currently under development in the IETF [TAPS working group](https://datatracker.ietf.org/wg/taps/about/). PANAPI is an [open-source project](https://www.ngi.eu/funded_solution/ngi-pointer-project-33/) that adds support for the [SCION network architecture](https://scion-architecture.net/) to a general purpose TAPS-like networking API. It is funded by the [EU NGI Pointer initiative](https://pointer.ngi.eu/).

For more details, please check out the following links: 
* [PANAPI Paper](https://dl.acm.org/doi/pdf/10.1145/3472727.3472808)
* [PANAPI presentation to the IETF TAPS Meeting](https://datatracker.ietf.org/meeting/113/materials/slides-113-taps-panapi-implementation-00) (March 23, 2022)
* [Basic Path Selection presentation](assets/presentation.pdf) (May 5, 2022)
  * [Demo video](https://www.youtube.com/watch?v=2_I7xbsk89I) demonstrating Basic Path Selection
* [Advanced Path Selection](assets/Presentation_Milestone3.pdf) (Oct 7, 2022)
* [Evaluation of Path Selection](assets/Presentation_Milestone4.pdf) (Oct 7, 2022)

## `import "panapi"` - The PANAPI Library

* [x] Simple working code example, see [examples/concurrent/README.md](examples/concurrent/README.md)
* [ ] Add more code examples

## `cmd/daemon` - The PANAPI Daemon

* [x] applications selecting SCION as transport benefit from daemon running in the backend
* [x] Graceful fallback to default behavior when daemon not running
* [x] Expose Quic performance monitoring via RPC to Lua script executed by Daemon
* [ ] Create dedicated daemon README
* [x] Lua scripting examples 
  * [x] [cmd/daemon/simple.lua](cmd/daemon/simple.lua)
  * [x] [cmd/daemon/pathselection.lua](cmd/daemon/pathselection.lua)
  * [x] [cmd/daemon/selector_with_stats.lua](cmd/daemon/selector_with_stats.lua)

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
  - [x] live access to connection preferences like `CapacityProfile`
- [x] Central path selection Daemon

### Path quality

- [x] Passive throughput monitoring
- [x] Exposed to Lua script

### Convenience features
- [ ] Different log levels

### Other
- [ ] Full test coverage
- [ ] Code Documentation
- [ ] Move scripting selector to `/pkg` such that it could be used without the rest of PANAPI

## Ported Applications
- [ ] `spate` traffic generator
- [x] `concurrent` code example client/server timestamp echoing
- [ ] `http`
  - [ ] server
  - [ ] client

## Affiliations

[![OVGU](assets/ovgu-small.png)](https://netsys.ovgu.de)

[![NGI Pointer](assets/NGI-Pointer-logo-small.png)](https://pointer.ngi.eu)

[![SCION](assets/scion-small.png)](https://scion-architecture.net)
