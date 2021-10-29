package main

import (
	//"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/netsys-lab/panapi"
	"github.com/netsys-lab/panapi/network"
	//"github.com/scionproto/scion/go/lib/snet"
)

type SpateClientSpawner struct {
	net            string
	transport      string
	server_address string
	packet_size    int
	single_path    bool
	interactive    bool
	bandwidth      int64
	parallel       int
}

// e.g. NewSpateClientSpawner("16-ffaa:0:1001,[172.31.0.23]:1337")
func NewSpateClientSpawner() SpateClientSpawner {
	return SpateClientSpawner{
		packet_size: 1208,
		single_path: false,
		interactive: false,
		bandwidth:   0,
		parallel:    8,
	}
}

func (s SpateClientSpawner) ServerAddress(server_address string) SpateClientSpawner {
	s.server_address = server_address
	return s
}

func (s SpateClientSpawner) Network(network string) SpateClientSpawner {
	s.net = network
	return s
}
func (s SpateClientSpawner) Transport(transport string) SpateClientSpawner {
	s.transport = transport
	return s
}

func (s SpateClientSpawner) PacketSize(packet_size int) SpateClientSpawner {
	s.packet_size = packet_size
	return s
}

func (s SpateClientSpawner) Parallel(parallel int) SpateClientSpawner {
	s.parallel = parallel
	return s
}

func (s SpateClientSpawner) SinglePath(single_path bool) SpateClientSpawner {
	s.single_path = single_path
	return s
}

func (s SpateClientSpawner) Interactive(interactive bool) SpateClientSpawner {
	s.interactive = interactive
	return s
}

func (s SpateClientSpawner) Bandwidth(bandwidth int64) SpateClientSpawner {
	s.bandwidth = bandwidth
	return s
}

func (s SpateClientSpawner) Spawn() error {
	//Info("Resolving address %s...", s.server_address)
	/*serverAddr, err := appnet.ResolveUDPAddr(s.server_address)
	if err != nil {
		Error("Resolution of UDP address (%s) failed: %v", s.server_address, err)
		return err
	}

	paths := []snet.Path{}
	if s.interactive {
	selection:
		for {
			chosen_path, err := appnet.ChoosePathInteractive(serverAddr.IA)
			if err != nil || chosen_path == nil {
				Error("Error while selecting chosen path! Please check the given paths.")
				os.Exit(1)
			}
			paths = append(paths, chosen_path)
			if s.single_path {
				break selection
			}
			fmt.Print("Would you like to choose an additional path? (y/N): ")
			var input string
		prompt:
			for {
				fmt.Scanln(&input)
				switch input {
				case "y", "Y":
					break prompt
				case "n", "N", "":
					break selection
				default:
					fmt.Print("Choose yes or no (y/N): ")
					continue
				}
			}
		}

	} else {
		Info("Searching paths to remote...")
		paths, err = appnet.QueryPaths(serverAddr.IA)
		if err != nil {
			Warn("Could not query for available paths: %v", err)
			Error("Could not find valid paths!")
			os.Exit(1)
		}
		if paths == nil {
			Warn("Detected test on direct connection. Multipath via SCION is not available...")
			paths = []snet.Path{nil}
		}
		if s.single_path {
			// Use first available path
			Info("Using single path")
			paths = paths[:1]
		}
	}*/

	paths := []int{1}
	Info("Choosing the following paths: %v", paths)

	rs := panapi.NewRemoteEndpoint()
	rs.WithNetwork(s.net)
	rs.WithTransport(s.transport)
	rs.WithAddress(s.server_address)

	preconn, err := panapi.NewPreconnection(rs)
	if err != nil {
		return err
	}

	Info("Establishing connections with server...")

	bytes_sent := 0
	packets_sent := 0

	complete := make(chan struct{}, len(paths))

	counter := make(chan int, 1024)
	stop := make(chan struct{}, 1)
	cancel := make(chan os.Signal, 1)
	signal.Notify(cancel, os.Interrupt)

	Info("Starting sending data for measurements...")
	start := time.Now()
	var wg sync.WaitGroup

	for _, path := range paths {
		// Set selected singular path
		Info("Creating %d new connections on path %v...", s.parallel, path)
		//appnet.SetPath(serverAddr, path)

		// Let's take care of appnet here and ensure that we dare not create multiple connections in parallel
		var connections = make([]network.Connection, s.parallel)

		for i := 0; i < s.parallel; i++ {
			//conn, err := appnet.DialAddrUDP(serverAddr)
			conn, err := preconn.Initiate()
			// Checking on err != nil will not work here as non-critical errors are returned
			if conn != nil {
				go awaitCompletion(conn, complete)
				connections[i] = conn
			} else {
				Warn("Connection on path failed: %v", err)
			}
		}

		go operatorThread(connections, complete, counter, stop, &wg, s)
	}

	closed_conn := 0
	total_conn := len(paths) * s.parallel
runner:
	for {
		select {
		case bytes := <-counter:
			bytes_sent += bytes
			packets_sent += 1
		case <-cancel:
			Info("Received interrupt signal, stopping flooding of available paths...")
			break runner
		case <-complete:
			closed_conn += 1
			if closed_conn >= total_conn {
				Info("Measurements finished on server!")
				break runner
			}
		}
	}

	elapsed := time.Since(start)
	actual_bandwidth := float64(bytes_sent) / elapsed.Seconds() * 8.0 / 1024.0 / 1024.0

	stop <- struct{}{}
	//wg.Wait()

	heading := color.New(color.Bold, color.Underline).Sprint("Summary")
	deco := color.New(color.Bold).Sprint("=====")
	lower := color.New(color.Bold).Sprint("===================")
	Info("         %s %s %s", deco, heading, deco)
	Info("         Sent data: %v KiB", bytes_sent/1024.0)
	Info("      Sent packets: %v packets", packets_sent)
	Info("       Packet size: %v B", s.packet_size)
	Info("          Duration: %s", elapsed)
	Info("  Target bandwidth: %v Mib/s", float64(s.bandwidth)/1024.0/1024.0)
	Info("  Actual bandwidth: %v Mib/s", actual_bandwidth)
	Info("         %s", lower)
	Info(">>> Please check the server measurements for the throughput achieved through")
	Info(">>> the network!")

	return nil
}

func operatorThread(conns []network.Connection, complete chan struct{}, counter chan int, stop chan struct{}, finalize *sync.WaitGroup, spawner SpateClientSpawner) {
	target_duration := time.Duration((float64(spawner.packet_size*8) / float64(spawner.bandwidth) * float64(spawner.parallel)) * float64(time.Second))
	data_cs := make([]chan *[]byte, spawner.parallel)
	reply := make(chan struct{}, spawner.parallel)
	rand := NewFastRand(uint64(spawner.packet_size))

	for i, _ := range data_cs {
		data_cs[i] = make(chan *[]byte, 1024)
		if spawner.bandwidth > 0 {
			go workerThread(conns[i], counter, stop, data_cs[i], reply, false)
		} else {
			go workerThread(conns[i], counter, stop, data_cs[i], reply, true)
		}
	}

	sum_error := time.Duration(0)

	for {
		start := time.Now()
		val := rand.Get()
		for _, c := range data_cs {
			c <- val
		}
		if spawner.bandwidth > 0 {
			packets := 0
			for _ = range reply {
				packets += 1
				if packets >= spawner.parallel {
					break
				}
			}
			end := time.Now()
			duration := end.Sub(start)
			sum_error += target_duration - duration
			// fmt.Println("Sending took ", duration)
			// fmt.Println("Supposed to take ", target_duration)
			if sum_error > 0 {
				sum_error = target_duration - duration
				// fmt.Println("Waiting ", target_duration - duration)
				time.Sleep(target_duration - duration)
			}
		}
	}
}

func workerThread(conn network.Connection, counter chan int, stop chan struct{}, data chan *[]byte, reply chan struct{}, skip bool) {
worker:
	for {
		select {
		case <-stop:
			break worker
		case dat := <-data:
			sent_bytes, err := conn.Write(*dat)
			if err != nil {
				Error("Sending data failed: %v", err)
				break
			}

			counter <- sent_bytes
			if skip {
				continue worker
			}
			reply <- struct{}{}
		}
	}
	//close(control_points)
}

func awaitCompletion(conn network.Connection, complete chan struct{}) {
	buf := make([]byte, 4)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			Error("Waiting for completion of measurement failed: %v", err)
			complete <- struct{}{}
			break
		}
		if string(buf) == "stop" {
			complete <- struct{}{}
			break
		}
	}
}
