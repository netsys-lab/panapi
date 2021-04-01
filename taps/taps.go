package taps

import (
	"flag"
	"strings"

	"code.ovgu.de/hausheer/taps-api/errs"
	"code.ovgu.de/hausheer/taps-api/ip"
	"code.ovgu.de/hausheer/taps-api/network"
	"code.ovgu.de/hausheer/taps-api/scion"
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

func GetFlags(network, address, transport *string) {
	address_p := flag.String("a", "[127.0.0.1]:1337", "ip or scion address and port")
	network_p := flag.String("n", NETWORK_IP, "network type: ip or scion")
	transport_p := flag.String("t", TRANSPORT_TCP, "transport protocol: udp, tcp, quic")

	flag.Parse()

	*address = *address_p
	*network = strings.ToUpper(*network_p)
	*transport = strings.ToUpper(*transport_p)

	// fmt.Println(*network, " | ", *transport, " | ", *address)
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
	default:
		return p, errs.TransportType
	}

	switch e.Network {
	case NETWORK_IP, NETWORK_IPV4, NETWORK_IPV6:
		network = ip.Network()
	case NETWORK_SCION:
		network = scion.Network()
	default:
		return p, errs.NetworkType
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
