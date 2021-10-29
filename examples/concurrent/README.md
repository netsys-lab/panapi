# Example Implementation for concurrent client/server

Building the `concurrent` binary:

```bash
go build .
```

## Example: Running in TCP/IP mode

Run a server with "regular IP" and "TCP" like so:

```bash
./concurrent -net IP -transport TCP -local :1337
```

From a second terminal, you can connect to the server like so:

```bash
./concurrent -net IP -transport TCP -remote 127.0.0.1:1337
```

After a second (or so), you should start seeing Timestamps being sent back and forth:
```
main.go:64: Message: 2021-10-29 12:44:42.98597813 +0200 CEST m=+1.006368091
main.go:64: Message: 2021-10-29 12:44:43.986855783 +0200 CEST m=+2.007244410
...
```

## Example: Running in QUIC over SCION

Server:

``` bash
./concurrent -net SCION -transport QUIC -local "19-ffaa:1:eb6,[127.0.0.1]:1337"
```
(Replace `19-ffaa:1:eb6` with the local ISD-AS)

Client:

```bash
./concurrent -net SCION -transport QUIC -remote "19-ffaa:1:eb6,[127.0.0.1]:1337"
```
(Replace `19-ffaa:1:eb6` with the remote ISD-AS. No need to call the client from a different AS though, testing locally works fine, but you still need to specify the full SCION address for now)


## Example: Running over UDP

_currently broken_
