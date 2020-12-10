package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"code.ovgu.de/hausheer/taps-api/taps"
)

func main() {
	serv := taps.NewLocalEndpoint()
	serv.WithInterface("any")
	serv.WithService("quic")
	// serv.WithService("tcp")
	serv.WithIPv4Address("127.0.0.1")
	serv.WithPort("5555")

	transProp := taps.NewTransportProperties()
	secParam := taps.NewSecurityParameters()
	preconn := taps.NewPreconnection(serv, transProp, secParam)
	lis := preconn.Listen()

	quitter := make(chan bool)
	sender := make(chan string)

	go func() {
		for {
			str, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			if strings.Contains(str, "#") {
				quitter <- true
			} else {
				sender <- str
			}
		}
	}()

	var conn taps.Connection

loop:
	for {
		select {
		case conn = <-lis.ConnRec:
			go func() {
				for {
					msg := conn.Receive()
					fmt.Print(msg.Data)
				}
			}()
		case msg := <-sender:
			conn.Send(taps.NewMessage(msg, ""))
		case <-quitter:
			lis.Stop()
			break loop
		}
	}

	conn.Close()
}
