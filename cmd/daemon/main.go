package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"runtime/pprof"

	"github.com/netsec-ethz/scion-apps/pkg/pan"
	"github.com/netsys-lab/panapi/network/scion"
	"github.com/netsys-lab/panapi/rpc"
)

func main() {
	var (
		script   string
		cpulog   string
		selector rpc.ServerSelector
		err      error
	)

	flag.StringVar(&script, "script", "", "Lua script for path selection")
	flag.StringVar(&cpulog, "cpulog", "", "Write profiling information to file")
	flag.Parse()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, os.Kill, os.Interrupt)

	l, err := net.ListenUnix("unix", rpc.DefaultDaemonAddress)
	if err != nil {
		log.Fatalf("Could not start daemon: %s", err)
	}
	log.Println("Starting daemon")

	if cpulog != "" {
		f, err := os.Create(cpulog)
		if err != nil {
			log.Fatal("cpuprofile:", err)
		}
		if err = pprof.StartCPUProfile(f); err != nil {
			log.Fatal("cpuprofile:", err)
		}
	}

	selector, err = scion.NewLuaSelector(script)
	if err != nil {
		log.Printf("Could not load path-selection script: %s", err)
		log.Println("Falling back to default selector")
		selector = rpc.NewServerSelectorFunc(func(pan.UDPAddr) pan.Selector {
			return &pan.DefaultSelector{}
		})
	}

	//serverselector := rpc.NewServerSelectorFunc(func(raddr,
	server, err := rpc.NewSelectorServer(selector)
	if err != nil {
		log.Fatalln(err)
	}
	go func() {
		log.Println("Started listening for rpc calls")
		server.Accept(l)
	}()
	sig := <-c
	log.Printf("Got signal [%s]: running defered cleanup and exiting.", sig)
	err = l.Close()
	if err != nil {
		log.Println(err)
	}
	//should be NOP if profiler is not running
	pprof.StopCPUProfile()
	os.Exit(0)
}
