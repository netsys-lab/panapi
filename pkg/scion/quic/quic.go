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
	l quic.Listener
}

func (l *listener) Accept() (taps.Connection, error) {
	if l.l == nil {
		return nil, errors.New("not a listener")
	}
	session, err := l.l.Accept(context.Background())
	if err != nil {
		return nil, err
	}
	stream, err := session.AcceptStream(context.Background())
	return taps.Connection(stream), err
}

func (l *listener) Close() error {
	return l.l.Close()
}

type Protocol struct {
	QuicConfig *quic.Config
	TLSConfig  *tls.Config
	Selector   taps.Selector
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
		err = q.Selector.SetPreferences(p.ConnectionPreferences)
		if err != nil {
			return nil, err
		}
	}
	l, err := pan.ListenQUIC(
		context.Background(),
		netaddr.IPPortFrom(addr.IP, addr.Port),
		nil,
		q.TLSConfig,
		q.QuicConfig,
	)
	return &listener{l: l}, err
}

func (q *Protocol) Initiate(p *taps.Preconnection) (taps.Connection, error) {
	addr, err := pan.ResolveUDPAddr(p.RemoteEndpoint.Address)
	if err != nil {
		return nil, err
	}
	if q.Selector != nil {
		err = q.Selector.SetPreferences(p.ConnectionPreferences)
		if err != nil {
			return nil, err
		}
	}
	session, err := pan.DialQUIC(
		context.Background(),
		netaddr.IPPort{},
		addr,
		nil,
		q.Selector,
		"",
		q.TLSConfig,
		q.QuicConfig,
	)
	if err != nil {
		return nil, err
	}

	stream, err := session.OpenStream() //Sync(context.Background())
	return taps.Connection(stream), err

}
