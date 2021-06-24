// Copyright 2018 ETH Zurich
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
	"math"
	"os"
	"time"

	"github.com/netsys-lab/scion-measurement-lib/measured-appnet"
)

func main() {
	var err error
	// get local and remote addresses from program arguments:
	port := flag.Uint("port", 0, "[Server] local port to listen on")
	remoteAddr := flag.String("remote", "", "[Client] Remote (i.e. the server's) SCION Address (e.g. 17-ffaa:1:1,[127.0.0.1]:12345)")
	config := flag.String("config", "", "[Client] Config file containing several remotes to ping. One per line. (e.g. 17-ffaa:1:1,[127.0.0.1]:12345)")
	flag.Parse()

	switch {
	case *port > 0:
		err = runServer(uint16(*port))
		check(err)
	case len(*remoteAddr) > 0:
		err = runClient(*remoteAddr)
		check(err)
	case len(*config) > 0:
		lines, err := readConfig(*config)
		if err != nil {
			log.Fatal(err)
			return
		}
		resp := make(chan error)
		for _, line := range lines {
			go runClientAsync(line, resp)
		}
		check(<-resp)
	default:
		check(fmt.Errorf("Either specify -port for server or -remote, -config for client"))
	}
}

func runServer(port uint16) error {
	conn, err := measured_appnet.ListenPort(port)
	if err != nil {
		return err
	}
	defer conn.Close()

	buffer := make([]byte, 4*1024)
	for {
		_, from, err := conn.ReadFrom(buffer)
		if err != nil {
			return err
		}

		log.Printf("Received ping from %s\n", from)

		if err != nil {
			return err
		}

		_, err = conn.WriteTo([]byte("ack"), from)

		if err != nil {
			return err
		}
	}
}

func runClientAsync(a string, c chan error) {
	c <- runClient(a)
}

func runClient(address string) error {
	conn, err := measured_appnet.Dial(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	for {
		_, err := conn.Write([]byte("ping"))
		if err != nil {
			return err
		}
		send := time.Now()

		buffer := make([]byte, 4*1024)
		n, from, err := conn.ReadFrom(buffer)
		if err != nil {
			fmt.Print(err)
		}

		log.Printf("%v bytes from %s: time=%v\n", n, from, math.Abs(send.Sub(time.Now()).Seconds()))
	}
}

// Check just ensures the error is nil, or complains and quits
func check(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, "Fatal error. Exiting.", "err", e)
		os.Exit(1)
	}
}

func readConfig(f string) ([]string, error) {
	file, err := os.Open(f)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return lines, nil
}
