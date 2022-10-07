// Copyright 2021,2022 Thorben Krüger (thorben.krueger@ovgu.de)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
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
		remote, local, t, n, p                       string
		proto                                        taps.Protocol
		server, client, daemontracer, daemonselector bool
		bytes                                        int64
	)

	flag.StringVar(&remote, "remote", "", `[Client] Remote (i.e. the server's) Address
        (e.g. 17-ffaa:1:1,[127.0.0.1]:1337 or 192.0.2.1:1337, depending on chosen network type)`)
	flag.StringVar(&local, "local", "", `[Server] Local Address to listen on
        (e.g. 17-ffaa:1:1,[127.0.0.1]:1337 or 0.0.0.0:1337, depending on chosen network type)`)
	flag.StringVar(&n, "net", "IP", "network type (IP|SCION)")
	flag.StringVar(&t, "transport", "QUIC", "transport protocol (TCP|QUIC)")
	flag.StringVar(&p, "profile", "CapacitySeeking", "SCION capacity profile (Default|CapacitySeeking|Scavenger|LowLatency)")
	flag.BoolVar(&daemontracer, "daemontracer", false, "use PANAPI daemon tracer")
	flag.BoolVar(&daemonselector, "daemonselector", false, "use PANAPI daemon selector")
	flag.Int64Var(&bytes, "bytes", 1000*1000*10, "amount of bytes to transfer during experiment")

	flag.Parse()

	//log.SetFlags(log.Lshortfile)

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

				if !daemontracer || err != nil {
					config.Tracer = nil
				}
				if !daemonselector || err != nil {
					selector = &taps.DefaultSelector{}
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

	profile := taps.Default
	switch p {
	case "CapacitySeeking":
		profile = taps.CapacitySeeking
	case "Scavenger":
		profile = taps.Scavenger
	case "LowLatency":
		profile = taps.LowLatencyNonInteractive
	case "Default", "":
	default:
		log.Fatalf("Unknown capacity profile \"%s\"", p)
	}

	if len(local) > 0 {
		log.Println(runServer(local, proto))
	} else {
		err := runClient(bytes, remote, proto, profile)
		if err != nil {
			log.Println(err)
		}
	}
}

func myCopy(w io.Writer, r io.Reader, c chan int) (total int64, err error) {
	buf := make([]byte, 1024*32)
	for {
		var nr int
		nr, err = r.Read(buf)
		nw, erw := w.Write(buf[:nr])
		total += int64(nw)
		c <- nr
		if erw == nil && nw != nr {
			err = fmt.Errorf("short write")
			break
		}
		if err == io.EOF {
			err = nil
			break
		}
		if erw != nil {
			err = fmt.Errorf("Write error: %s", erw)
			break
		}
		if err != nil {
			err = fmt.Errorf("Read err: %s", err)
			break
		}
	}
	return
}

func report(c chan int, verbose bool) {
	total := 0
	subtotal := 0
	ticker := time.Tick(time.Second)
	begin := time.Now()
	for {
		select {
		case bytes := <-c:
			subtotal += bytes
		case <-ticker:
			total += subtotal
			dur := time.Since(begin)
			fmt.Printf("%d,%d,%d\n", int(dur.Seconds()), subtotal, total)
			if verbose {
				log.Printf("%.3f kb/s", float64(subtotal)/1000)
			}
			subtotal = 0

		}
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

	c := make(chan int)
	running := false

	for {
		Connection, err := Listener.Accept()
		if !running {
			go report(c, true)
			running = true
		}
		log.Printf("Got Connection from: %s", Connection.Preconnection().RemoteEndpoint.Address)
		if err != nil {
			log.Println(err)
		} else {
			go myCopy(ioutil.Discard, Connection, c)
		}
	}
}

func runClient(bytes int64, remote string, proto taps.Protocol, profile taps.CapacityProfile) error {
	RemoteSpecifier := taps.RemoteEndpoint{}
	RemoteSpecifier.Address = remote
	RemoteSpecifier.Protocol = proto

	Preconnection := taps.Preconnection{
		RemoteEndpoint: &RemoteSpecifier,
		ConnectionPreferences: &taps.ConnectionPreferences{
			ConnCapacityProfile: profile,
		},
	}

	Connection, err := Preconnection.Initiate()
	if err != nil {
		return fmt.Errorf("Initate error: %s", err)
	}
	defer Connection.Close()

	c := make(chan int)
	go report(c, false)

	reader := io.LimitReader(rand.Reader, bytes)
	begin := time.Now()
	n, err := myCopy(Connection, reader, c)
	if err == nil {
		log.Printf("%s: Average: %.3f kb/s", profile, float64(n)/(1000*time.Since(begin).Seconds()))
	}
	return err
}
