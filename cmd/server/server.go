package main

import (
	"fmt"
	"log"

	"code.ovgu.de/hausheer/taps-api/taps"
)

func main() {
	var network, address, transport string
	taps.Init(&network, &address, &transport)

	LocalSpecifier := taps.NewLocalEndpoint()
	LocalSpecifier.WithNetwork(network)
	LocalSpecifier.WithAddress(address)
	LocalSpecifier.WithTransport(transport)

	// LocalSpecifier.WithNetwork(taps.NETWORK_IP)
	// LocalSpecifier.WithAddress(":1337")
	// LocalSpecifier.WithNetwork(taps.NETWORK_SCION)
	// LocalSpecifier.WithAddress("19-ffaa:1:e9e,[127.0.0.1]:1337")
	// LocalSpecifier.WithTransport(taps.TRANSPORT_UDP)
	// LocalSpecifier.WithTransport(taps.TRANSPORT_QUIC)
	// LocalSpecifier.WithTransport(taps.TRANSPORT_TCP)

	Preconnection, err := taps.NewPreconnection(LocalSpecifier)
	if err != nil {
		log.Fatalf("Error! %s\n", err)
	}

	Listener := Preconnection.Listen()
	Connection := <-Listener.ConnectionReceived
	Listener.Stop()

	Message, err := Connection.Receive()
	if err != nil {
		fmt.Printf("Error! %s\n", err)
	}
	fmt.Printf("Message: %v\n", Message.String())

	err = Connection.Send(taps.Message("Got your message!\n"))
	if err != nil {
		fmt.Printf("Error! %s\n", err)
	}

	Connection.Close()
}
