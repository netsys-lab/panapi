package panapi

import (
	"github.com/netsys-lab/panapi/network"
)

func NewRemoteEndpoint() *network.Endpoint {
	return &network.Endpoint{Local: false}
}

func NewLocalEndpoint() *network.Endpoint {
	return &network.Endpoint{Local: true}
}

// HACK, let's see if this works
type TransportProperties network.TransportProperties

func NewTransportProperties() *TransportProperties {
	return &TransportProperties{}
}
