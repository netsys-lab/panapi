// Copyright 2021,2022 Thorben Krüger (thorben.krueger@ovgu.de)
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
	"log"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/netsys-lab/panapi/pkg/convenience"
	"github.com/netsys-lab/panapi/taps"

	iquic "github.com/netsys-lab/panapi/pkg/inet/quic"
	tcp "github.com/netsys-lab/panapi/pkg/inet/tcp"
	squic "github.com/netsys-lab/panapi/pkg/scion/quic"
)

func main() {
	var (
		remote, local, t, n, p string
		proto                  taps.Protocol
		server, client         bool
	)

	flag.StringVar(&remote, "remote", "", `[Client] Remote (i.e. the server's) Address
        (e.g. 17-ffaa:1:1,[127.0.0.1]:1337 or 192.0.2.1:1337, depending on chosen network type)`)
	flag.StringVar(&local, "local", "", `[Server] Local Address to listen on
        (e.g. 17-ffaa:1:1,[127.0.0.1]:1337 or 0.0.0.0:1337, depending on chosen network type)`)
	flag.StringVar(&n, "net", "IP", "network type (IP|SCION)")
	flag.StringVar(&t, "transport", "QUIC", "transport protocol (TCP|QUIC)")
	flag.StringVar(&p, "profile", "", "profile to select ("+taps.Scavenger.String()+"|"+taps.LowLatencyInteractive.String()+")")
	flag.Parse()

	log.SetFlags(log.Ltime)

	if len(local) > 0 {
		server = true
	}
	if len(remote) > 0 {
		client = true
	}
	if server == client {
		log.Fatalln("Either specify -port for server or -remote for client")
	}

	if t == "TCP" {
		if n == "SCION" {
			log.Fatalln("Transport TCP is not supported for Network Type SCION")
		}
		proto = &tcp.Protocol{}
	} else if t == "QUIC" {
		tlsConf := convenience.DummyTLSConfig()
		if n == "IP" {
			proto = &iquic.Protocol{
				TLSConfig: &tlsConf,
			}
		} else if n == "SCION" {
			var (
				config   = &quic.Config{}
				selector taps.Selector
				err      error
			)
			if client {
				selector, config.Tracer, err = convenience.RPCClientHelper()
				if err != nil {
					log.Println(err)
				}
			}
			proto = &squic.Protocol{
				squic.Config{
					TLS:      &tlsConf,
					Selector: selector,
					Quic:     config,
				},
			}
		} else {
			log.Fatalln("Either specify -n IP or -n SCION")
		}
	} else {
		log.Fatalln("Either specify -t TCP or -t QUIC")
	}

	if len(local) > 0 {
		log.Println(runServer(local, proto))
	} else {
		if p == taps.Scavenger.String() {
			log.Println(runClient(remote, proto, taps.Scavenger))
		} else if p == taps.LowLatencyInteractive.String() {
			log.Println(runClient(remote, proto, taps.LowLatencyInteractive))
		} else {
			log.Println(runClient(remote, proto, taps.Default))
		}
	}
}

func serverWorker(conn taps.Connection) {
	defer conn.Close()
	var (
		r    = bufio.NewReader(conn)
		err  error
		ping string
	)
	for {
		ping, err = r.ReadString('\n')
		if err != nil {
			break
		}
		_, err = conn.Write([]byte(ping))
		if err != nil {
			break
		}
	}
	log.Println(err)
}

func clientWorker(pconn taps.Preconnection) {
	conn, err := pconn.Initiate()
	if err == nil {
		defer conn.Close()
		var (
			ticker    = time.Tick(time.Second)
			r         = bufio.NewReader(conn)
			prefs     = pconn.ConnectionPreferences
			now, then time.Time
			response  string
			i         uint
			sw        bool
		)
		if prefs.ConnCapacityProfile == taps.Default {
			// Switch Profiles occasionally
			sw = true
		}
		for {
			i += 1
			if i%10 == 0 && sw {
				if i%20 == 0 {
					prefs.ConnCapacityProfile = taps.Scavenger
				} else {
					prefs.ConnCapacityProfile = taps.LowLatencyNonInteractive
				}
				pconn.SetPreferences(prefs)
			}

			_, err = conn.Write([]byte(time.Now().Format(time.RFC3339Nano) + "\n"))
			if err != nil {
				break
			}
			response, err = r.ReadString('\n')
			if err != nil {
				break
			}
			now = time.Now()
			then, err = time.Parse(time.RFC3339Nano+"\n", response)
			if err != nil {
				break
			}
			log.Printf(
				"read %d bytes from %s (Profile: %s): %s",
				len(response),
				pconn.RemoteEndpoint.Address,
				pconn.ConnectionPreferences.ConnCapacityProfile,
				now.Sub(then),
			)
			<-ticker
		}
	}
	log.Println(err)
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
		} else {
			log.Println("Got connection")
			go serverWorker(Connection)
		}
	}
}

func runClient(remote string, proto taps.Protocol, profile taps.CapacityProfile) error {
	RemoteSpecifier := taps.RemoteEndpoint{}
	RemoteSpecifier.Address = remote
	RemoteSpecifier.Protocol = proto

	Preconnection := taps.Preconnection{
		RemoteEndpoint: &RemoteSpecifier,
		ConnectionPreferences: &taps.ConnectionPreferences{
			ConnCapacityProfile: profile,
		},
	}

	clientWorker(Preconnection)

	return nil

}
