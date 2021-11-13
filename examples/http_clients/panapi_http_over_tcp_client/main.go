package main

import (
	"log"

	"github.com/netsys-lab/panapi"
	"github.com/netsys-lab/panapi/network"
)

func main() {
	runClient()
}

func runClient() error {
	RemoteSpecifier := panapi.NewRemoteEndpoint()
	RemoteSpecifier.WithNetwork("IP")
	RemoteSpecifier.WithAddress("www.google.com:80")
	RemoteSpecifier.WithTransport("TCP")

	Preconnection, err := panapi.NewPreconnection(RemoteSpecifier)
	fcheck(err)

	Connection, err := Preconnection.Initiate()
	fcheck(err)

	defer Connection.Close()

	Connection.Send(network.DummyMessage("GET / HTTP/1.0\r\n\r\n"))

	m, err := Connection.Receive()
	fcheck(err)

	log.Printf("Message: %s\n", m)

	return nil

}

func fcheck(err error) {
	if err != nil {
		log.Fatalf("Error! %s\n", err)
	}
}
