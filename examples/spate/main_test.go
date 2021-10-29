package main

import (
	"testing"
)

// achieves ~22 GB/s on Ryzen 5 3600
func BenchmarkXorshift(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var frand = NewFastRand(134217728)
		// get 1 GibiByte of random data
		frand.Get()
	}
}
