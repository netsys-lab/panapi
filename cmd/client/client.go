package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"code.ovgu.de/hausheer/taps-api/taps"
)

func main() {
	cli := *taps.NewRemoteEndpoint()
	cli.WithInterface("any")
	cli.WithService("tcp")
	cli.WithHostname("127.0.0.1")
	cli.WithPort("5555")

	transProp := *taps.NewTransportProperties()
	secParam := *taps.NewSecurityParameters()

	preconn := *taps.NewPreconnection(cli, transProp, secParam)
	conn := preconn.Initiate()

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

	go func() {
		for {
			msg := conn.Receive()
			fmt.Print(msg.Data)
		}
	}()

loop:
	for {
		select {
		case msg := <-sender:
			conn.Send(*taps.NewMessage(msg, ""))
		case <-quitter:
			break loop
		}
	}

	conn.Close()
}
