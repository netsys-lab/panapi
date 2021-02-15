package main

import (
	"log"

	"code.ovgu.de/hausheer/taps-api/taps"
)

func main() {
	LocalSpecifier := taps.NewLocalEndpoint()
	LocalSpecifier.WithNetwork(taps.NETWORK_IP)
	LocalSpecifier.WithTransport(taps.TRANSPORT_QUIC)
	LocalSpecifier.WithAddress(":1337")

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
