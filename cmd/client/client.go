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

func check(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func main() {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()

	servF := flag.String("serv", "tcp", "tcp or quic")
	ipF := flag.String("ip", "127.0.0.1", "ip address")
	portF := flag.String("port", "1111", "port")
	interF := flag.String("inter", "any", "interface name")
	flag.Parse()

	cli := taps.NewRemoteEndpoint()
	cli.WithInterface(*interF)
	cli.WithService(*servF)
	cli.WithIPv4Address(*ipF)
	cli.WithPort(*portF)

	transProp := taps.NewTransportProperties()

	privatKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	secParam := taps.NewSecurityParameters()
	secParam.Set("keypair", privatKey, &privatKey.PublicKey)

	preconn, err := taps.NewPreconnection(cli, transProp, secParam)

	conn, err := preconn.Initiate()

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
		case <-quitter:
			break loop
		}
	}

	conn.Close()
}
