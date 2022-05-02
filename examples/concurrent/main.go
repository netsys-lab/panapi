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
	"bufio"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/netsys-lab/panapi/pkg/convenience"
	iquic "github.com/netsys-lab/panapi/pkg/inet/quic"
	"github.com/netsys-lab/panapi/pkg/inet/tcp"
	squic "github.com/netsys-lab/panapi/pkg/scion/quic"
	"github.com/netsys-lab/panapi/rpc"
	"github.com/netsys-lab/panapi/taps"
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
		remote, local, t string
		proto            taps.Protocol
		server, client   bool
	)

	flag.StringVar(&remote, "remote", "", "[Client] Remote (i.e. the server's) Address (e.g. 17-ffaa:1:1,[127.0.0.1]:1337 or 192.0.2.1:1337, depending on chosen network type)")
	flag.StringVar(&local, "local", "", "[Server] Local Address to listen on, (e.g. 17-ffaa:1:1,[127.0.0.1]:1337 or 0.0.0.0:1337, depending on chosen network type)")
	//flag.StringVar(&n, "net", network.NETWORK_IP, "network type")
	flag.StringVar(&t, "transport", "tcp", "transport protocol (tcp|quic|squic")
	//flag.StringVar(&script, "script", "", "[Client] Lua script for path selection")
	//flag.UintVar(&port, "port", 0, "[Server] local port to listen on")
	flag.Parse()

	log.SetFlags(log.Lshortfile)

	if len(local) > 0 {
		server = true
	}

	if len(remote) > 0 {
		client = true
	}

	if server == client {
		fcheck(fmt.Errorf("Either specify -port for server or -remote for client"))
	}

	if t == "tcp" {
		proto = &tcp.Protocol{}
	} else if t == "quic" || t == "squic" {
		tlsConf := convenience.GenerateTLSConfig()
		tlsConf.NextProtos = []string{"concurrent-quic-test"}
		tlsConf.InsecureSkipVerify = true
		if t == "quic" {
			proto = &iquic.Protocol{
				TLSConfig: &tlsConf,
			}
		} else {
			var (
				config   *quic.Config
				selector taps.Selector
			)
			if client {
				c, err := convenience.NewRPCClient()
				if err != nil {
					log.Fatalln(err)
				}
				selector = rpc.NewSelectorClient(c)
				config = &quic.Config{Tracer: rpc.NewTracerClient(c)}
			}
			proto = &squic.Protocol{
				TLSConfig:  &tlsConf,
				Selector:   selector,
				QuicConfig: config,
			}
		}

	} else {
		fcheck(fmt.Errorf("Either specify -t tcp, -t quic or -t squic"))
	}

	if len(local) > 0 {
		check(runServer(local, proto))
	} else {
		check(runClient(remote, proto))
	}
}

func worker(conn taps.Connection) {
	defer conn.Close()
	ticker := time.Tick(time.Second)

	r := bufio.NewReader(conn)
	for {
		request := (<-ticker).String() + "\n"
		_, err := conn.Write([]byte(request))
		if !check(err) {
			break
		}
		response, err := r.ReadString('\n')
		if !check(err) {
			break
		}
		log.Printf("Message: %s", response)
	}

}

func runServer(local string, proto taps.Protocol) error {
	LocalSpecifier := taps.LocalEndpoint{}
	LocalSpecifier.Address = local
	LocalSpecifier.Protocol = proto

	Preconnection := taps.Preconnection{
		LocalEndpoint: &LocalSpecifier,
	}

	Listener, err := Preconnection.Listen()
	if err != nil {
		return err
	}

	for {
		Connection, err := Listener.Accept()
		if err != nil {
			log.Println(err)
		}
		go worker(Connection)
	}

	return nil

}

func runClient(remote string, proto taps.Protocol) error {
	RemoteSpecifier := taps.RemoteEndpoint{}
	RemoteSpecifier.Address = remote
	RemoteSpecifier.Protocol = proto

	Preconnection := taps.Preconnection{
		RemoteEndpoint: &RemoteSpecifier,
		ConnectionPreferences: &taps.ConnectionPreferences{
			ConnCapacityProfile: taps.CapacitySeeking,
		},
	}

	Connection, err := Preconnection.Initiate()
	if err != nil {
		return err
	}

	worker(Connection)

	return nil

}
