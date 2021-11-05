package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	//"runtime/pprof"

	"github.com/netsec-ethz/scion-apps/pkg/pan"
	"github.com/netsys-lab/panapi/rpc"
)

func main() {
	/*f, err := os.Create("cpuprofile")
	if err != nil {
		log.Fatal("could not create cpuprofile:", err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile:", err)
	}*/
	l, err := net.ListenUnix("unix", rpc.DefaultDaemonAddress)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Started listening for rpc calls")
	server, err := rpc.NewSelectorServer(&pan.DefaultSelector{})
	if err != nil {
		log.Fatalln(err)
	}
	go func() {
		server.Accept(l)
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, os.Kill, os.Interrupt)
	sig := <-c
	log.Printf("Got signal [%s]: running defered cleanup and exiting.", sig)
	log.Println(l.Close())
	/*pprof.StopCPUProfile()*/
	os.Exit(0)
}
