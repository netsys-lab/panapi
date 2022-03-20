package main

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"github.com/netsys-lab/panapi"
	"github.com/netsys-lab/panapi/network"
	"io"
	"log"
)

func runSender(net, transport, remote string, size int64) error {
	remoteEndpoint := panapi.NewRemoteEndpoint()
	remoteEndpoint.WithNetwork(net)
	remoteEndpoint.WithAddress(remote)
	remoteEndpoint.WithTransport(transport)

	transportProps := network.NewTransportProperties()

	msg, err := getRandomMsg(size)
	if err != nil {
		return err
	}
	sum := md5.Sum([]byte(msg.String()))

	pcon, err := panapi.NewPreconnection(remoteEndpoint, transportProps)
	if err != nil {
		return err
	}

	con, err := pcon.Initiate()
	if err != nil {
		return err
	}

	if err := con.Send(msg); err != nil {
		log.Printf("error sending message: %s", err)
		return err
	} else {
		log.Printf("successfully sent message, hahsum: %x", sum)
		return con.Close()
	}
}

func getRandomMsg(size int64) (*network.FixedMessage, error) {
	var buf bytes.Buffer
	if _, err := io.CopyN(&buf, rand.Reader, size); err != nil {
		return nil, err
	}
	return network.NewFixedMessageByte(buf.Bytes()), nil
}
