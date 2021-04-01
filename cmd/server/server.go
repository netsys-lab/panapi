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
	taps.GetFlags(&network, &address, &transport)

	LocalSpecifier := taps.NewLocalEndpoint()
	LocalSpecifier.WithNetwork(network)
	LocalSpecifier.WithAddress(address)
	LocalSpecifier.WithTransport(transport)

	// LocalSpecifier.WithNetwork(taps.NETWORK_IP)
	// LocalSpecifier.WithAddress(":1337")
	// LocalSpecifier.WithNetwork(taps.NETWORK_SCION)
	// LocalSpecifier.WithAddress("19-ffaa:1:e9e,[127.0.0.1]:1337")
	// LocalSpecifier.WithTransport(taps.TRANSPORT_UDP)
	// LocalSpecifier.WithTransport(taps.TRANSPORT_TCP)
	// LocalSpecifier.WithTransport(taps.TRANSPORT_QUIC)

	Preconnection, err := taps.NewPreconnection(LocalSpecifier)
	fcheck(err)

	Listener := Preconnection.Listen()
	Connection := <-Listener.ConnectionReceived
	fcheck(Connection.GetError())

	err = Listener.Stop()
	check(err)

	Message, err := Connection.Receive()
	check(err)
	fmt.Printf("Message: %v\n", Message.String())

	err = Connection.Send(taps.Message("Hi from server!\n"))
	check(err)

	err = Connection.Close()
	check(err)
}
