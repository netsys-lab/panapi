package main

import (
	"flag"
	"github.com/docker/go-units"
	"github.com/netsys-lab/panapi/network"
	"log"
)

func main() {

	var (
		net, transport, mode, listenAddr, remoteAddr, script, sizeHuman string
	)

	// common flags
	flag.StringVar(&net, "net", network.NETWORK_IP, "network type")
	flag.StringVar(&transport, "transport", network.TRANSPORT_QUIC, "transport protocol")
	flag.StringVar(&mode, "mode", "server", "mode, server or client")
	flag.StringVar(&sizeHuman, "size", "1MiB", "amount of (random) data the server generates and serves / the clients expects")

	// server-only flags
	flag.StringVar(&listenAddr, "listenAddr", "", "[Server] Local Address to listen on, (e.g. 17-ffaa:1:1,[127.0.0.1]:1337 or 0.0.0.0:1337, depending on chosen network type)")

	// client-only flags
	flag.StringVar(&remoteAddr, "remoteAddr", "", "[Client] Server's Address (e.g. 17-ffaa:1:1,[127.0.0.1]:1337 or 192.0.2.1:1337, depending on chosen network type)")
	flag.StringVar(&script, "script", "", "[Client] Lua script for path selection")

	flag.Parse()
	size, err := units.FromHumanSize(sizeHuman)
	if err != nil {
		log.Fatalf("could not parse size %s", sizeHuman)
	}

	switch mode {
	case "server":
		err := runServer(net, transport, listenAddr, size)
		if err != nil {
			log.Fatalf("Error running server: %s", err)
		}
	case "client":
		err := runClient(net, transport, remoteAddr, script, size)
		if err != nil {
			log.Fatalf("Error running client: %s", err)
		}
	default:
		log.Fatalf("unknown mode, must be either 'server' or 'client'")
	}
}
