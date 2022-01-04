# PANAPI - Work in Progress
[![Go](https://github.com/netsys-lab/panapi/actions/workflows/go.yml/badge.svg)](https://github.com/netsys-lab/panapi/actions/workflows/go.yml)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

PANAPI is an early [research](https://dl.acm.org/doi/10.1145/3472727.3472808) implementation of a next-generation networking [API to the transport layer](https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html). The latter is currently under development in the IETF [TAPS working group](https://datatracker.ietf.org/wg/taps/about/). PANAPI is an [EU-funded](https://pointer.ngi.eu/) open-source project, that adds support for the [SCION network architecture](https://scion-architecture.net/) to a general purpose TAPS-like networking API.

## Development branch

_In this development branch of the repository, the depreceated old API from the
[main branch](https://github.com/netsys-lab/panapi) is being replaced
by a more faithful implementation of the [current TAPS API specification](https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html).
This work is currently still in mockup stage. See [here](doc/Implementation.md) for a detailled discussion of our particular design choices. A first impression of the Go package structure can be found here:__

[![Go Reference](https://pkg.go.dev/badge/github.com/netsys-lab/panapi.svg)](https://pkg.go.dev/github.com/netsys-lab/panapi@v0.3.0-alpha7/taps)



## Affiliations

[![OVGU](assets/ovgu-small.png)](https://netsys.ovgu.de)

[![NGI Pointer](assets/NGI-Pointer-logo-small.png)](https://pointer.ngi.eu)

[![SCION](assets/scion-small.png)](https://scion-architecture.net)
