// Copyright 2021 Thorben Krüger (thorben.krueger@ovgu.de)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"flag"
	"fmt"
	"github.com/netsys-lab/panapi"
	"github.com/netsys-lab/panapi/network"
	"log"
	"time"
)

func fcheck(err error) {
	if err != nil {
		log.Fatalf("Error! %s\n", err)
	}
}

func check(err error) bool {
	if err != nil {
		log.Printf("Error! %s\n", err)
		return false
	}
	return true
}

func main() {
	var (
		n, remote, local, t, script string
		//		port         uint
	)

	flag.StringVar(&remote, "remote", "", "[Client] Remote (i.e. the server's) Address (e.g. 17-ffaa:1:1,[127.0.0.1]:1337 or 192.0.2.1:1337, depending on chosen network type)")
	flag.StringVar(&local, "local", "", "[Server] Local Address to listen on, (e.g. 17-ffaa:1:1,[127.0.0.1]:1337 or 0.0.0.0:1337, depending on chosen network type)")
	flag.StringVar(&n, "net", network.NETWORK_IP, "network type")
	flag.StringVar(&t, "transport", network.TRANSPORT_QUIC, "transport protocol")
	flag.StringVar(&script, "script", "", "[Client] Lua script for path selection")
	//flag.UintVar(&port, "port", 0, "[Server] local port to listen on")
	flag.Parse()

	log.SetFlags(log.Lshortfile)

	if (len(local) > 0) == (len(remote) > 0) {
		check(fmt.Errorf("Either specify -port for server or -remote for client"))
	}

	if len(local) > 0 {
		check(runServer(n, t, local))
	} else {
		check(runClient(n, t, remote, script))
	}
}

func worker(conn network.Connection) {
	defer conn.Close()
	ticker := time.Tick(time.Second)

	for check(conn.GetError()) {
		request, err := network.NewLineMessageString((<-ticker).String())
		if !check(err) {
			break
		}
		if !check(conn.Send(request)) {
			break
		}
		response := network.NewLineMessage()
		err = conn.Receive(response)
		if !check(err) {
			break
		}
		log.Printf("Message: %s", response)
	}

}

func runServer(net, t, local string) error {
	LocalSpecifier := panapi.NewLocalEndpoint()
	LocalSpecifier.WithNetwork(net)
	LocalSpecifier.WithAddress(local)
	LocalSpecifier.WithTransport(t)

	Preconnection, err := panapi.NewPreconnection(LocalSpecifier, nil)
	if err != nil {
		return err
	}

	Listener := Preconnection.Listen()

	for {
		Connection := <-Listener.ConnectionReceived
		go worker(Connection)
	}

	return nil

}

func runClient(net, t, remote, script string) error {
	RemoteSpecifier := panapi.NewRemoteEndpoint()
	RemoteSpecifier.WithNetwork(net)
	RemoteSpecifier.WithAddress(remote)
	RemoteSpecifier.WithTransport(t)

	TransportProperties := network.NewTransportProperties()
	if script != "" {
		TransportProperties.Set("lua-script", script)
	}

	Preconnection, err := panapi.NewPreconnection(RemoteSpecifier, TransportProperties)
	if err != nil {
		return err
	}

	Connection, err := Preconnection.Initiate()
	if err != nil {
		return err
	}
	worker(Connection)

	return nil

}
