package quic

import (
	"context"
	"crypto/tls"
	"errors"

	"github.com/lucas-clemente/quic-go"
	"github.com/netsys-lab/panapi/taps"
)

type listener struct {
	p *taps.Preconnection
	l quic.Listener
}

type Connection struct {
	quic.Stream
	p *taps.Preconnection
}

func (c *Connection) Preconnection() *taps.Preconnection {
	return c.p
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
	return &Connection{stream, l.p}, err
}

func (l *listener) Close() error {
	return l.l.Close()
}

type Protocol struct {
	QuicConfig *quic.Config
	TLSConfig  *tls.Config
}

func (q *Protocol) Selector() taps.Selector {
	return nil
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
		Multipath:         taps.Disabled,
	}, err

}

func (q *Protocol) NewListener(p *taps.Preconnection) (taps.Listener, error) {
	_, err := q.Satisfy(p)
	if err != nil {
		return nil, err
	}

	//tlsConf := convenience.GenerateTLSConfig()
	// TODO what does this do?
	//tlsConf.NextProtos = []string{"panapi-quic-test"}
	l, err := quic.ListenAddr(
		p.LocalEndpoint.Address,
		q.TLSConfig,
		q.QuicConfig,
	)
	return &listener{l: l}, err
}

func (q *Protocol) Initiate(p *taps.Preconnection) (taps.Connection, error) {
	/*tlsConf := &tls.Config{
		//Certificates: quicutil.MustGenerateSelfSignedCert(),
		InsecureSkipVerify: true,
		NextProtos:         []string{"panapi-quic-test"},
	}*/
	session, err := quic.DialAddr(p.RemoteEndpoint.Address, q.TLSConfig, q.QuicConfig)
	if err != nil {
		return nil, err
	}

	stream, err := session.OpenStream() //Sync(context.Background())
	return &Connection{stream, p}, err

}
