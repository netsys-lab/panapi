package panapi

import (
	"github.com/netsys-lab/panapi/network"
)

const (
	DaemonSocketPath = "/tmp/panapi.sock"
)

func NewRemoteEndpoint() *network.Endpoint {
	return &network.Endpoint{Local: false}
}

func NewLocalEndpoint() *network.Endpoint {
	return &network.Endpoint{Local: true}
}

// HACK, let's see if this works
type TransportProperties struct {
	network.TransportProperties
}

func NewTransportProperties() *TransportProperties {
	return &TransportProperties{
		*network.NewTransportProperties(),
	}
}
