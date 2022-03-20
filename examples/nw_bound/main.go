package main

import (
	"flag"
	"github.com/docker/go-units"
	"github.com/netsys-lab/panapi/network"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {

	var (
		net, transport, mode, listenAddr, remoteAddr, sizeHuman string
	)

	// common flags
	flag.StringVar(&net, "net", network.NETWORK_IP, "network type")
	flag.StringVar(&transport, "transport", network.TRANSPORT_QUIC, "transport protocol")
	flag.StringVar(&mode, "mode", "server", "mode, either 'receiver' or 'sender'")
	flag.StringVar(&sizeHuman, "size", "1MiB", "amount of (random) data the sender generates and uploads / the receiver expects")

	// receiver-only flags
	flag.StringVar(&listenAddr, "listenAddr", "", "[Receiver] Local Address to listen on, (e.g. 17-ffaa:1:1,[127.0.0.1]:1337 or 0.0.0.0:1337, depending on chosen network type)")

	// sender-only flags
	flag.StringVar(&remoteAddr, "remoteAddr", "", "[Sender] Server's Address (e.g. 17-ffaa:1:1,[127.0.0.1]:1337 or 192.0.2.1:1337, depending on chosen network type)")

	flag.Parse()
	size, err := units.FromHumanSize(sizeHuman)
	if err != nil {
		log.Fatalf("could not parse size %s", sizeHuman)
	}

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	switch mode {
	case "receiver":
		err := runReceiver(net, transport, listenAddr, size)
		if err != nil {
			log.Fatalf("Error running server: %s", err)
		}
	case "sender":
		err := runSender(net, transport, remoteAddr, size)
		if err != nil {
			log.Fatalf("Error running client: %s", err)
		}
	default:
		log.Fatalf("unknown mode, must be either 'receiver' or 'sender'")
	}

}
