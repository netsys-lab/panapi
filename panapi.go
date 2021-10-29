package panapi

import (
	"github.com/netsys-lab/panapi/network"
	"github.com/netsys-lab/panapi/network/ip"
	"github.com/netsys-lab/panapi/network/scion"
)

type Preconnection struct {
	endpoint *network.Endpoint
	listener network.Listener
	dialer   network.Dialer
}

type Listener struct {
	listener           network.Listener
	ConnectionReceived chan network.Connection
}

func NewRemoteEndpoint() *network.Endpoint {
	return &network.Endpoint{Local: false}
}

func NewLocalEndpoint() *network.Endpoint {
	return &network.Endpoint{Local: true}
}

func (p *Preconnection) Listen() Listener {
	c := make(chan network.Connection)
	go func(p *Preconnection, c chan network.Connection) {
		conn, err := p.listener.Listen()
		if err != nil {
			conn.SetError(err)
		}
		c <- conn
	}(p, c)
	return Listener{ConnectionReceived: c, listener: p.listener}
}

func (p *Preconnection) Initiate() (network.Connection, error) {
	return p.dialer.Dial()
}

func (l *Listener) Stop() error {
	return l.listener.Stop()
}

func NewPreconnection(e *network.Endpoint) (Preconnection, error) {
	var (
		l      network.Listener
		dialer network.Dialer
		net    network.Network
		p      Preconnection
		err    error
	)
	switch e.Transport {
	case network.TRANSPORT_UDP:
	case network.TRANSPORT_TCP:
	case network.TRANSPORT_QUIC:
	default:
		return p, network.AddrTypeError
	}

	switch e.Network {
	case network.NETWORK_IP, network.NETWORK_IPV4, network.NETWORK_IPV6:
		net = ip.Network()
	case network.NETWORK_SCION:
		net = scion.Network()
	default:
		return p, network.NetTypeError
	}

	if e.Local {
		l, err = net.NewListener(e)
	} else {
		dialer, err = net.NewDialer(e)
	}
	if err != nil {
		return p, err
	}

	p = Preconnection{
		endpoint: e,
		listener: l,
		dialer:   dialer,
	}
	return p, err
}
