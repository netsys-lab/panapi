package main

import (
	"flag"
	"fmt"
	"github.com/netsys-lab/panapi"
	"github.com/netsys-lab/panapi/network"
	"log"
)

func fcheck(err error) {
	if err != nil {
		log.Fatalf("Error! %s\n", err)
	}
}

func check(err error) {
	if err != nil {
		log.Printf("Error! %s\n", err)
	}
}

func main() {
	var net, address, transport string

	flag.StringVar(&net, "n", network.NETWORK_IP, "network type")
	flag.StringVar(&address, "a", "[127.0.0.1]:1337", "network address and port")
	flag.StringVar(&transport, "t", network.TRANSPORT_TCP, "transport protocol")

	flag.Parse()

	//taps.GetFlags(&network, &address, &transport)

	RemoteSpecifier := panapi.NewRemoteEndpoint()

	RemoteSpecifier.WithNetwork(net)
	RemoteSpecifier.WithAddress(address)
	RemoteSpecifier.WithTransport(transport)

	// RemoteSpecifier.WithNetwork(taps.NETWORK_IP)
	// RemoteSpecifier.WithAddress("[127.0.0.1]:1337")
	// RemoteSpecifier.WithNetwork(taps.NETWORK_SCION)
	// RemoteSpecifier.WithAddress("19-ffaa:1:e9e,[127.0.0.1]:1337")
	// RemoteSpecifier.WithTransport(taps.TRANSPORT_UDP)
	// RemoteSpecifier.WithTransport(taps.TRANSPORT_TCP)
	// RemoteSpecifier.WithTransport(taps.TRANSPORT_QUIC)

	Preconnection, err := panapi.NewPreconnection(RemoteSpecifier)
	fcheck(err)

	Connection, err := Preconnection.Initiate()
	fcheck(err)

	err = Connection.Send(network.DummyMessage("Hi from client!\n"))
	check(err)

	Message, err := Connection.Receive()
	check(err)
	fmt.Printf("Message: %v\n", Message)

	err = Connection.Close()
	check(err)
}
