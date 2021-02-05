package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"os"
	"strings"

	"code.ovgu.de/hausheer/taps-api/taps"
)

func check(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func main() {
	var err error

	servF, addrF, portF, interF := taps.Init()

	cli := taps.NewRemoteEndpoint()
	err = cli.WithInterface(*interF)
	check(err)
	err = cli.WithService(*servF)
	check(err)
	err = cli.WithAddress(*addrF)
	// err = cli.WithHostname("localhost")
	check(err)
	err = cli.WithPort(*portF)
	check(err)

	transProp := taps.NewTransportProperties()

	privatKey, err := rsa.GenerateKey(rand.Reader, 1024)
	check(err)
	secParam := taps.NewSecurityParameters()
	secParam.Set("keypair", privatKey, &privatKey.PublicKey)
	check(err)

	preconn, err := taps.NewPreconnection(cli, transProp, secParam)
	check(err)
	conn, err := preconn.Initiate()
	check(err)

	quitter := make(chan bool)
	sender := make(chan string)

	go func() {
		var b []byte = make([]byte, 1)
		for {
			os.Stdin.Read(b)
			str := string(b)
			// str, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			if strings.Contains(str, ".") {
				quitter <- true
			} else {
				sender <- str
			}
		}
	}()

	go func() {
		for {
			msg, err := conn.Receive()
			check(err)
			fmt.Print(msg.Data)
		}
	}()

loop:
	for {
		select {
		case msg := <-sender:
			conn.Send(taps.NewMessage(msg, ""))
			check(err)
		case <-quitter:
			fmt.Println()
			break loop
		}
	}

	conn.Close()
	check(err)
}
