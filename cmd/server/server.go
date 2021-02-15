package main

import (
	"log"

	"code.ovgu.de/hausheer/taps-api/taps"
)

func main() {
	LocalSpecifier := taps.NewLocalEndpoint()
	// LocalSpecifier.WithNetwork(taps.NETWORK_IP)
	// LocalSpecifier.WithAddress(":1234")
	LocalSpecifier.WithNetwork(taps.NETWORK_SCION)
	LocalSpecifier.WithAddress("19-ffaa:1:e9e,[127.0.0.1]:1234")
	// LocalSpecifier.WithTransport(taps.TRANSPORT_UDP)
	LocalSpecifier.WithTransport(taps.TRANSPORT_QUIC)

	Preconnection, err := taps.NewPreconnection(LocalSpecifier)
	if err != nil {
		log.Fatalf("Error! %s\n", err)
	}

	Listener := Preconnection.Listen()

	Connection := <-Listener

	Message, err := Connection.Receive()
	if err != nil {
		log.Printf("Error! %s\n", err)
	}
	log.Printf("Message: %v\n", Message)

	err = Connection.Send(taps.Message("Got your message!\n"))
	if err != nil {
		log.Printf("Error! %s\n", err)
	}
}
