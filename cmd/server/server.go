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

	ser := taps.NewLocalEndpoint()
	err = ser.WithInterface(*interF)
	check(err)
	err = ser.WithService(*servF)
	check(err)
	err = ser.WithAddress(*addrF)
	check(err)
	err = ser.WithPort(*portF)
	check(err)

	transProp := taps.NewTransportProperties()
	// transProp.Require(taps.NAGLE_ON)

	privatKey, err := rsa.GenerateKey(rand.Reader, 1024)
	check(err)
	secParam := taps.NewSecurityParameters()
	// err = secParam.Set("keypair", 1, &privatKey.PublicKey)
	err = secParam.Set("keypair", privatKey, &privatKey.PublicKey)
	check(err)

	preconn, err := taps.NewPreconnection(ser, transProp, secParam)
	check(err)
	lis, err := preconn.Listen()
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

	var conn taps.Connection

loop:
	for {
		select {
		case conn = <-lis.ConnRec:
			check(conn.Err)
			go func() {
				for {
					msg, err := conn.Receive()
					check(err)
					fmt.Print(msg.Data)
				}
			}()
		case msg := <-sender:
			err = conn.Send(taps.NewMessage(msg, ""))
			check(err)
		case <-quitter:
			fmt.Println()
			lis.Stop()
			break loop
		}
	}

	conn.Close()
	check(err)
}
