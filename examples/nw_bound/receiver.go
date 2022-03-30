package main

import (
	"crypto/md5"
	"github.com/netsys-lab/panapi"
	"github.com/netsys-lab/panapi/network"
	"log"
)

func runReceiver(net, transport, listenAddr string, size int64) error {
	localEndpoint := panapi.NewLocalEndpoint()
	localEndpoint.WithNetwork(net)
	localEndpoint.WithAddress(listenAddr)
	localEndpoint.WithTransport(transport)

	pcon, err := panapi.NewPreconnection(localEndpoint, nil)
	if err != nil {
		return err
	}

	listen := pcon.Listen()

	for {
		con := <-listen.ConnectionReceived
		go handleCon(size, con)
	}
}

func handleCon(size int64, con network.Connection) {
	log.Printf("handling con of %s", con.RemoteAddr())

	defer con.Close()
	msg := network.NewFixedMessage(size)

	if err := con.Receive(msg); err != nil {
		log.Printf("error receiving msg from %s: %s\n", con.RemoteAddr(), err)
		return
	}
	log.Printf("successfully received msg (hashsum: %x) from %s\n", md5.Sum(msg.Bytes()), con.RemoteAddr())
}
