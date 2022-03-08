# nw_bound

This example implements a client/server for simulating network bound IO.
The server serves randomized data to any client.

## Build

```bash
user@host:~/panapi/nw_bound$ go build . 
```

## Arguments

```bash
user@host:~/panapi/nw_bound$ ./nw_bound --help
```

## Example: Running a server with SCION/QUIC

```bash
./nw_bound \
  -mode server \
  -net SCION \
  -transport QUIC \
  -listenAddr "19-ffaa:1:1303,[127.0.0.1]:1337" \
  -size 5MiB
```

Replace the value of _listenAddr_ with the IA identifier, IP address and port 
the server should listen to.

## Example: Running a client with SCION/QUIC

```bash
./nw_bound \
  -mode client \
  -net SCION \
  -transport QUIC \
  -remoteAddr "19-ffaa:1:1303,[127.0.0.1]:1337" \
  -size 5MiB
```

Replace the value of _remoteAddr_ with the IA identifier, IP address and port 
of your network_bound server.