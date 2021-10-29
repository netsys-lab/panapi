package main

import (
	"flag"
	"fmt"
	"github.com/netsys-lab/panapi"
	"github.com/netsys-lab/panapi/network"
	"log"
	"time"
)

func fcheck(err error) {
	if err != nil {
		log.Fatalf("Error! %s\n", err)
	}
}

func check(err error) bool {
	if err != nil {
		log.Printf("Error! %s\n", err)
		return false
	}
	return true
}

func main() {
	var (
		n, remote, local, t string
		//		port         uint
	)

	flag.StringVar(&remote, "remote", "", "[Client] Remote (i.e. the server's) Address (e.g. 17-ffaa:1:1,[127.0.0.1]:1337 or 192.0.2.1:1337, depending on chosen network type)")
	flag.StringVar(&local, "local", "", "[Server] Local Address to listen on, (e.g. 17-ffaa:1:1,[127.0.0.1]:1337 or 0.0.0.0:1337, depending on chosen network type)")
	flag.StringVar(&n, "net", network.NETWORK_IP, "network type")
	flag.StringVar(&t, "transport", network.TRANSPORT_QUIC, "transport protocol")
	//flag.UintVar(&port, "port", 0, "[Server] local port to listen on")
	flag.Parse()

	log.SetFlags(log.Lshortfile)

	if (len(local) > 0) == (len(remote) > 0) {
		check(fmt.Errorf("Either specify -port for server or -remote for client"))
	}

	if len(local) > 0 {
		check(runServer(n, t, local))
	} else {
		check(runClient(n, t, remote))
	}
}

func worker(conn network.Connection) {
	defer conn.Close()
	ticker := time.Tick(time.Second)

	for check(conn.GetError()) {
		if !check(conn.Send(network.DummyMessage((<-ticker).String()))) {
			break
		}
		m, err := conn.Receive()
		if !check(err) {
			break
		}
		log.Printf("Message: %s\n", m)
	}

}

func runServer(net, t, local string) error {
	LocalSpecifier := panapi.NewLocalEndpoint()
	LocalSpecifier.WithNetwork(net)
	LocalSpecifier.WithAddress(local)
	LocalSpecifier.WithTransport(t)

	Preconnection, err := panapi.NewPreconnection(LocalSpecifier)
	log.Printf("%v, %v", Preconnection, err)
	if err != nil {
		return err
	}

	Listener := Preconnection.Listen()

	for {
		Connection := <-Listener.ConnectionReceived
		go worker(Connection)
	}

	return nil

}

func runClient(net, t, remote string) error {
	RemoteSpecifier := panapi.NewRemoteEndpoint()
	RemoteSpecifier.WithNetwork(net)
	RemoteSpecifier.WithAddress(remote)
	RemoteSpecifier.WithTransport(t)

	Preconnection, err := panapi.NewPreconnection(RemoteSpecifier)
	if err != nil {
		return err
	}

	Connection, err := Preconnection.Initiate()
	if err != nil {
		return err
	}
	worker(Connection)

	return nil

}
