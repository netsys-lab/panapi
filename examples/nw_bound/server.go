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

func runServer(net, transport, listenAddr string, size int64) error {
	localEndpoint := panapi.NewLocalEndpoint()
	localEndpoint.WithNetwork(net)
	localEndpoint.WithAddress(listenAddr)
	localEndpoint.WithTransport(transport)

	var buf bytes.Buffer
	if _, err := io.CopyN(&buf, rand.Reader, size); err != nil {
		return err
	}

	data := buf.Bytes()

	log.Printf("serving data, hashsum: %x", md5.Sum(data))

	pcon, err := panapi.NewPreconnection(localEndpoint, nil)
	if err != nil {
		return err
	}

	listen := pcon.Listen()

	for {
		con := <-listen.ConnectionReceived
		go handleCon(&data, con)
	}
}

func handleCon(data *[]byte, con network.Connection) {
	log.Printf("handling con of %s", con.RemoteAddr())

	o := network.NewLineMessage()
	if err := con.Receive(o); err != nil {
		log.Printf("error receiving opening msg of %s: %s", con.RemoteAddr(), err)
	}

	msg := network.NewFixedMessageByte(*data)
	if err := con.Send(msg); err != nil {
		log.Printf("error transmitting pkt to %s: %s\n", con.RemoteAddr(), err)
	} else {
		log.Printf("successfully transmitted filedata to %s\n", con.RemoteAddr())
	}
}
