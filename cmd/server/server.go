package main

import (
	"crypto/rand"
	"crypto/rsa"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"code.ovgu.de/hausheer/taps-api/taps"
)

// (http anfrage händisch)
// flag: service port ip4 -> done
// jedes symbol übertragen -> done, non portable
// setnodelay tcp -> is default
// tls nach set -> done
// error handling
// scion bsp

func check(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func main() {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()

	var err error

	servF := flag.String("serv", "tcp", "tcp or quic")
	ipF := flag.String("ip", "127.0.0.1", "ip address")
	portF := flag.String("port", "1111", "port")
	interF := flag.String("inter", "any", "interface name")
	flag.Parse()

	serv := taps.NewLocalEndpoint()
	err = serv.WithInterface(*interF)
	check(err)
	err = serv.WithService(*servF)
	check(err)
	err = serv.WithIPv4Address(*ipF)
	check(err)
	err = serv.WithPort(*portF)
	check(err)

	transProp := taps.NewTransportProperties()

	privatKey, err := rsa.GenerateKey(rand.Reader, 1024)
	check(err)

	secParam := taps.NewSecurityParameters()
	err = secParam.Set("keypair", privatKey, &privatKey.PublicKey)
	check(err)

	preconn, err := taps.NewPreconnection(serv, transProp, secParam)
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
			conn.Send(taps.NewMessage(msg, ""))
		case <-quitter:
			lis.Stop()
			break loop
		}
	}

	conn.Close()
}
