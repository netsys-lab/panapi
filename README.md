# PANAPI - Path-aware Networking API
[![Go](https://github.com/netsys-lab/panapi/actions/workflows/go.yml/badge.svg)](https://github.com/netsys-lab/panapi/actions/workflows/go.yml)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

PANAPI aims to implement [path-awareness support](https://dl.acm.org/doi/10.1145/3472727.3472808) to [TAPS](https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html), a next-generation networking API to the transport layer. TAPS is currently under development in the [IETF TAPS working group](https://datatracker.ietf.org/wg/taps/about/). Path-awareness is a novel networking concept supported by the [SCION network architecture](https://scion-architecture.net/), enabling multi-path communication at the inter-domain level. PANAPI is an open-source project funded by [the EU NGI Pointer initiative](https://pointer.ngi.eu/).

## Development branch

_In this development branch of the repository, the depreceated old API from the
[main branch](https://github.com/netsys-lab/panapi) is being replaced
by a more faithful implementation of the [current TAPS API specification](https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html).
This work is currently still in an early stage._ **See [here](doc/Implementation.md) for a detailed discussion of our particular design choices.** _A first impression of the Go package structure can be found here: [![Go Reference](https://pkg.go.dev/badge/github.com/netsys-lab/panapi.svg)](https://pkg.go.dev/github.com/netsys-lab/panapi@v0.3.0-alpha7/taps)_



## Affiliations

[![OVGU](assets/ovgu-small.png)](https://netsys.ovgu.de)

[![NGI Pointer](assets/NGI-Pointer-logo-small.png)](https://pointer.ngi.eu)

[![SCION](assets/scion-small.png)](https://scion-architecture.net)
