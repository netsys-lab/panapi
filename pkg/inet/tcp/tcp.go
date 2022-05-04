package tcp

import (
	"errors"
	"net"

	"github.com/netsys-lab/panapi/taps"
)

type listener struct {
	l net.Listener
}

func (l *listener) Accept() (taps.Connection, error) {
	if l.l == nil {
		return nil, errors.New("not a listener")
	}
	conn, err := l.l.Accept()
	return taps.Connection(conn), err
}

func (l *listener) Close() error {
	return l.l.Close()
}

type Protocol struct{}

func (t *Protocol) Satisfy(p *taps.Preconnection) (*taps.TransportProperties, error) {
	sp := p.TransportPreferences
	var err error
	if sp.Reliability == taps.Prohibit ||
		sp.PreserveOrder == taps.Prohibit ||
		sp.CongestionControl == taps.Prohibit {
		err = errors.New("Can't satisfy all constraints")
	}
	return &taps.TransportProperties{
		Reliability:       true,
		PreserveOrder:     true,
		CongestionControl: true,
		Multipath:         taps.Disabled,
	}, err

}

func (t *Protocol) NewListener(p *taps.Preconnection) (taps.Listener, error) {
	_, err := t.Satisfy(p)
	if err != nil {
		return nil, err
	}
	addr := p.LocalEndpoint.Address
	l, err := net.Listen("tcp", addr)
	return &listener{l: l}, err

}

func (t *Protocol) Initiate(p *taps.Preconnection) (taps.Connection, error) {
	_, err := t.Satisfy(p)
	if err != nil {
		return nil, err
	}
	addr := p.RemoteEndpoint.Address
	conn, err := net.Dial("tcp", addr)
	return taps.Connection(conn), err
}
