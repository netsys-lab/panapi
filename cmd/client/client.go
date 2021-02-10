package main

import (
	"code.ovgu.de/hausheer/taps-api/taps"
	"fmt"
)

func main() {
	RemoteSpecifier := taps.NewRemoteEndpoint()
	RemoteSpecifier.WithNetwork(taps.NETWORK_IP)
	RemoteSpecifier.WithTransport(taps.TRANSPORT_TCP)
	RemoteSpecifier.WithAddress("localhost:1337")

	Preconnection, err := taps.NewPreconnection(RemoteSpecifier)
	if err != nil {
		fmt.Printf("Error! %s\n", err)
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
}
