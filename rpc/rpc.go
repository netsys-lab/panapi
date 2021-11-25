package rpc

import (
	"net/rpc"

	"github.com/lucas-clemente/quic-go/logging"
)

func NewServer(selector ServerSelector, tracer logging.Tracer) (*rpc.Server, error) {
	err := rpc.Register(NewSelectorServer(selector))
	if err != nil {
		return nil, err
	}
	err = rpc.Register(NewTracerServer(tracer))
	if err != nil {
		return nil, err
	}
	return rpc.DefaultServer, nil
}
