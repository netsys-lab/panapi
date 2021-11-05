package main

import (
	"log"

	"github.com/netsec-ethz/scion-apps/pkg/pan"
	"github.com/netsys-lab/panapi/rpc"
)

func main() {
	s := pan.DefaultSelector{}
	_, err := rpc.NewSelectorServer(&s)
	if err != nil {
		log.Fatalln(err)
	}
}
