package taps

import "io"

type Protocol interface {
	Dial(Endpoint) (io.ReadWriteCloser, error)
	Accept(Endpoint) (io.ReadWriteCloser, error)
	//Rendezvous(Endpoint, Endpoint) (io.ReadWriteCloser, error)
}
