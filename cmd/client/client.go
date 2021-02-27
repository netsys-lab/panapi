package main

import (
	"fmt"
	"log"

	"code.ovgu.de/hausheer/taps-api/taps"
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
	var network, address, transport string
	taps.Init(&network, &address, &transport)

	RemoteSpecifier := taps.NewRemoteEndpoint()

	RemoteSpecifier.WithNetwork(network)
	RemoteSpecifier.WithAddress(address)
	RemoteSpecifier.WithTransport(transport)

	// RemoteSpecifier.WithNetwork(taps.NETWORK_IP)
	// RemoteSpecifier.WithAddress("[127.0.0.1]:1337")
	// RemoteSpecifier.WithNetwork(taps.NETWORK_SCION)
	// RemoteSpecifier.WithAddress("19-ffaa:1:e9e,[127.0.0.1]:1337")
	// RemoteSpecifier.WithTransport(taps.TRANSPORT_UDP)
	// RemoteSpecifier.WithTransport(taps.TRANSPORT_TCP)
	// RemoteSpecifier.WithTransport(taps.TRANSPORT_QUIC)

	Preconnection, err := taps.NewPreconnection(RemoteSpecifier)
	fcheck(err)

	Connection, err := Preconnection.Initiate()
	fcheck(err)

	err = Connection.Send(taps.Message("Hai!\n"))
	check(err)

	Message, err := Connection.Receive()
	check(err)
	fmt.Printf("Message: %v\n", Message)

	err = Connection.Close()
	check(err)
}
