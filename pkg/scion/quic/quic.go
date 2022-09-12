package quic

import (
	"context"
	"crypto/tls"
	"errors"

	"github.com/lucas-clemente/quic-go"
	"github.com/netsec-ethz/scion-apps/pkg/pan"
	"github.com/netsys-lab/panapi/taps"
	"inet.af/netaddr"
)

type listener struct {
	p *taps.Preconnection
	l quic.Listener
}

type Connection struct {
	quic.Stream
	p *taps.Preconnection
	quic.Session
}

func (c *Connection) Preconnection() *taps.Preconnection {
	return c.p
}

func (c *Connection) Close() error {
	c.Stream.Close()
	return c.Session.CloseWithError(0, "closed")
}

func (l *listener) Accept() (taps.Connection, error) {
	if l.l == nil {
		return nil, errors.New("not a listener")
	}
	session, err := l.l.Accept(context.Background())
	if err != nil {
		return nil, err
	}
	ep := taps.Endpoint{Address: session.RemoteAddr().String()}
	l.p.RemoteEndpoint = &taps.RemoteEndpoint{Endpoint: ep}
	stream, err := session.AcceptStream(context.Background())
	return &Connection{stream, l.p, session}, err
}

func (l *listener) Close() error {
	return l.l.Close()
}

type Config struct {
	Quic     *quic.Config
	TLS      *tls.Config
	Selector taps.Selector
}

type Protocol struct {
	Config Config
}

func (q *Protocol) Selector() taps.Selector {
	return q.Config.Selector
}

func (q *Protocol) Satisfy(p *taps.Preconnection) (*taps.TransportProperties, error) {
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
		Multipath:         taps.Passive,
	}, err

}

func (q *Protocol) NewListener(p *taps.Preconnection) (taps.Listener, error) {
	_, err := q.Satisfy(p)
	if err != nil {
		return nil, err
	}
	addr, err := pan.ResolveUDPAddr(p.LocalEndpoint.Address)
	if err != nil {
		return nil, err
	}
	if p.ConnectionPreferences != nil {
		err = q.Config.Selector.SetPreferences(p.ConnectionPreferences)
		if err != nil {
			return nil, err
		}
	}
	l, err := pan.ListenQUIC(
		context.Background(),
		netaddr.IPPortFrom(addr.IP, addr.Port),
		nil,
		q.Config.TLS,
		q.Config.Quic,
	)
	return &listener{p: p, l: l}, err
}

func (q *Protocol) Initiate(p *taps.Preconnection) (taps.Connection, error) {
	addr, err := pan.ResolveUDPAddr(p.RemoteEndpoint.Address)
	if err != nil {
		return nil, err
	}
	if q.Config.Selector != nil {
		err = q.Config.Selector.SetPreferences(p.ConnectionPreferences)
		if err != nil {
			return nil, err
		}
	}
	session, err := pan.DialQUIC(
		context.Background(),
		netaddr.IPPort{},
		addr,
		nil,
		q.Config.Selector,
		"",
		q.Config.TLS,
		q.Config.Quic,
	)
	if err != nil {
		return nil, err
	}

	stream, err := session.OpenStream() //Sync(context.Background())
	return &Connection{stream, p, session}, err

}
