# PANAPI -- Work In Progress

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
- [ ] Central path selection Daemon

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
