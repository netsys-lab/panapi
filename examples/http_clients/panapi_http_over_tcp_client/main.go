package main

import (
	"log"
	"net/textproto"

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

	toSend := network.NewFixedMessageString("This is my message body which is made of text")
	header := make(textproto.MIMEHeader)
	header.Add("content-type", "text")
	toSend.SetHeader(&header)
	toSend.ToHTTPMessage()
	fcheck(err)
	Connection.Send(toSend)

	response := network.NewFixedMessage(1024)
	err = Connection.Receive(response)
	fcheck(err)
	responseHeader, err := response.GetHeader()
	fcheck(err)

	log.Printf("Message: %s", response)
	log.Printf("Message: %s", responseHeader)

	return nil

}

func fcheck(err error) {
	if err != nil {
		log.Fatalf("Error! %s\n", err)
	}
}
