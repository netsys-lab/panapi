package main

import (
	"flag"
	"fmt"
	"github.com/netsys-lab/panapi"
	"github.com/netsys-lab/panapi/pkg/network"
	"log"
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
	var net, address, transport string

	flag.StringVar(&net, "n", network.NETWORK_IP, "network type")
	flag.StringVar(&address, "a", "[127.0.0.1]:1337", "network address and port")
	flag.StringVar(&transport, "t", network.TRANSPORT_TCP, "transport protocol")

	flag.Parse()

	LocalSpecifier := panapi.NewLocalEndpoint()
	LocalSpecifier.WithNetwork(net)
	LocalSpecifier.WithAddress(address)
	LocalSpecifier.WithTransport(transport)

	// LocalSpecifier.WithNetwork(taps.NETWORK_IP)
	// LocalSpecifier.WithAddress(":1337")
	// LocalSpecifier.WithNetwork(taps.NETWORK_SCION)
	// LocalSpecifier.WithAddress("19-ffaa:1:e9e,[127.0.0.1]:1337")
	// LocalSpecifier.WithTransport(taps.TRANSPORT_UDP)
	// LocalSpecifier.WithTransport(taps.TRANSPORT_TCP)
	// LocalSpecifier.WithTransport(taps.TRANSPORT_QUIC)

	Preconnection, err := panapi.NewPreconnection(LocalSpecifier)
	fcheck(err)

	Listener := Preconnection.Listen()
	Connection := <-Listener.ConnectionReceived
	fcheck(Connection.GetError())

	err = Listener.Stop()
	check(err)

	Message, err := Connection.Receive()
	check(err)
	fmt.Printf("Message: %v\n", Message.String())

	err = Connection.Send(network.DummyMessage("Hi from server!\n"))
	check(err)

	err = Connection.Close()
	check(err)
}
