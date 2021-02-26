package main

import (
	"fmt"
	"log"

	"code.ovgu.de/hausheer/taps-api/taps"
)

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
	// RemoteSpecifier.WithTransport(taps.TRANSPORT_QUIC)
	// RemoteSpecifier.WithTransport(taps.TRANSPORT_TCP)

	Preconnection, err := taps.NewPreconnection(RemoteSpecifier)
	if err != nil {
		log.Fatalf("Error! %s\n", err)
	}

	Connection := Preconnection.Initiate()

	err = Connection.Send(taps.Message("Hai!\n"))
	if err != nil {
		fmt.Printf("Error! %s\n", err)
	}

	Message, err := Connection.Receive()
	if err != nil {
		fmt.Printf("Error! %s\n", err)
	}
	fmt.Printf("Message: %v\n", Message)

	Connection.Close()
}
