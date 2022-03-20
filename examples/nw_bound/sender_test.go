package main

import (
	"github.com/docker/go-units"
	"log"
	"testing"
)

func Benchmark(b *testing.B) {
	size, _ := units.FromHumanSize("50 MiB")
	log.Println(runSender("SCION", "QUIC", "19-ffaa:1:e94,[78.46.228.180]:1337", size))
}
