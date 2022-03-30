package main

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"github.com/docker/go-units"
	"github.com/lucas-clemente/quic-go"
	"github.com/netsys-lab/panapi"
	"github.com/netsys-lab/panapi/network"
	"io"
	"log"
	"time"
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
	sum := md5.Sum(msg.Bytes())

	start := time.Now()
	pcon, err := panapi.NewPreconnection(remoteEndpoint, transportProps)
	if err != nil {
		return err
	}

	con, err := pcon.Initiate()
	if err != nil {
		return err
	}
	defer con.Close()

	err = con.Send(msg)
	if err != nil {
		qerr, ok := err.(*quic.ApplicationError)
		if !ok || qerr.ErrorCode != 0 {
			log.Printf("error sending message: %s", err)
			return err
		}
	}

	log.Printf("successfully sent message with hahsum: %x", sum)

	dur := time.Since(start)
	throughput := units.BytesSize(float64(size) / dur.Seconds())
	log.Printf("SUMMARY: send %d bytes in %s - throughput: %s", size, dur, throughput)
	return nil
}

func getRandomMsg(size int64) (*network.FixedMessage, error) {
	var buf bytes.Buffer
	if _, err := io.CopyN(&buf, rand.Reader, size); err != nil {
		return nil, err
	}
	return network.NewFixedMessageByte(buf.Bytes()), nil
}
