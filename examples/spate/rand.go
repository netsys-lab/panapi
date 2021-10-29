package main

import (
	"encoding/binary"
	"math/rand"
)

// Xorshift
type FastRand struct {
	value uint64
	size  uint64
	buf   []byte
}

func NewFastRand(size uint64) FastRand {
	return FastRand{value: rand.Uint64(), buf: make([]byte, size), size: size}
}

func (r FastRand) Get() *[]byte {
	for k := uint64(0); k < r.size/uint64(8); k++ {
		r.value ^= r.value << 13
		r.value ^= r.value >> 7
		r.value ^= r.value << 17
		binary.BigEndian.PutUint64(r.buf[k*8:], r.value)
	}
	r.value ^= r.value << 13
	r.value ^= r.value >> 7
	r.value ^= r.value << 17
	binary.BigEndian.PutUint64(r.buf[r.size-8:], r.value)
	return &r.buf
}
