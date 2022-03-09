package main

import (
	"crypto/md5"
	"github.com/netsys-lab/panapi"
	"github.com/netsys-lab/panapi/network"
	"log"
	"time"
)

func runClient(net, transport, remote string, size int64) error {
	remoteEndpoint := panapi.NewRemoteEndpoint()
	remoteEndpoint.WithNetwork(net)
	remoteEndpoint.WithAddress(remote)
	remoteEndpoint.WithTransport(transport)

	transportProps := network.NewTransportProperties()

	pcon, err := panapi.NewPreconnection(remoteEndpoint, transportProps)
	if err != nil {
		return err
	}

	con, err := pcon.Initiate()
	if err != nil {
		return err
	}

	o, _ := network.NewLineMessageString("hello")
	err = con.Send(o)
	if err != nil {
		return err
	}

	time.Sleep(1 * time.Second)

	msg := network.NewFixedMessage(size)
	if err := con.Receive(msg); err != nil {
		log.Printf("error receiving packet: %s", err)
		return err
	} else {
		log.Printf("successfully received message, hahsum: %x", md5.Sum([]byte(msg.String())))
		return con.Close()
	}
}
