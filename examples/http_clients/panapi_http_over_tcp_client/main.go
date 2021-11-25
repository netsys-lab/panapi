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
	RemoteSpecifier.WithAddress("127.0.0.1:8080")
	RemoteSpecifier.WithTransport("TCP")
	tps := network.NewTransportProperties()

	Preconnection, err := panapi.NewPreconnection(RemoteSpecifier, tps)
	fcheck(err)

	Connection, err := Preconnection.Initiate()
	fcheck(err)

	defer Connection.Close()

	toSend := network.NewFixedMessageString("GET / HTTP/1.0\r\n\r\n")
	toSend.Header.Add("content-type", "text")
	err = toSend.AddMIMEHeaderToMesaage()
	fcheck(err)
	Connection.Send(toSend)

	response := network.NewFixedMessage(1024)
	err = Connection.Receive(response)
	fcheck(err)

	log.Printf("Message: %s", response)

	return nil

}

func fcheck(err error) {
	if err != nil {
		log.Fatalf("Error! %s\n", err)
	}
}
