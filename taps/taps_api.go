package taps

import (
	"code.ovgu.de/hausheer/taps-api/ip"
	"code.ovgu.de/hausheer/taps-api/network"
	"code.ovgu.de/hausheer/taps-api/scion"
	"errors"
)

const (
	NETWORK_IP    = "IP"
	NETWORK_IPV4  = "IPv4"
	NETWORK_IPV6  = "IPv6"
	NETWORK_SCION = "SCION"

	TRANSPORT_UDP  = "UDP"
	TRANSPORT_TCP  = "TCP"
	TRANSPORT_QUIC = "QUIC"
)

var (
	errNetworkType    = errors.New("invalid network type")
	errTransportType  = errors.New("invalid address type")
	errNotImplemented = errors.New("not (yet) supported")
)

type Message string

func (m Message) String() string {
	return string(m)
}

func NewRemoteEndpoint() *network.Endpoint {
	return &network.Endpoint{Local: false}
}

func NewLocalEndpoint() *network.Endpoint {
	return &network.Endpoint{Local: true}
}

type Preconnection struct {
	endpoint *network.Endpoint
	listener network.Listener
	dialer   network.Dialer
}

// TODO, detect closed connections, do failure recovery
func (p Preconnection) Listen() chan network.Connection {
	c := make(chan network.Connection)
	go func(c chan network.Connection, listener network.Listener) {
		for {
			conn, err := listener.Listen()
			if err != nil {
				// FIXME TODO
				panic(err)
			}
			c <- conn
		}
	}(c, p.listener)
	return c
}

func (p Preconnection) Initiate() network.Connection {
	conn, err := p.dialer.Dial()
	if err != nil {
		panic(err)
	}
	return conn
}

func NewPreconnection(e *network.Endpoint) (Preconnection, error) {
	var (
		listener network.Listener
		dialer   network.Dialer
		network  network.Network
		p        Preconnection
		err      error
	)
	switch e.Transport {
	case TRANSPORT_UDP:
	case TRANSPORT_TCP:
	case TRANSPORT_QUIC:
		return p, errNotImplemented
	default:
		return p, errTransportType
	}

	switch e.Network {
	case NETWORK_IP:
		fallthrough
	case NETWORK_IPV4:
		fallthrough
	case NETWORK_IPV6:
		network = ip.Network()
	case NETWORK_SCION:
		network = scion.Network()
	default:
		return p, errNetworkType
	}

	if e.Local {
		listener, err = network.NewListener(e)
	} else {
		dialer, err = network.NewDialer(e)
	}

	if err != nil {
		return p, err
	}

	p = Preconnection{
		endpoint: e,
		listener: listener,
		dialer:   dialer,
	}
	return p, err
}
