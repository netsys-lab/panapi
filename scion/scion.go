package scion

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"

	"code.ovgu.de/hausheer/taps-api/connection"
	"code.ovgu.de/hausheer/taps-api/glob"
	"code.ovgu.de/hausheer/taps-api/network"
	"github.com/lucas-clemente/quic-go"
	"github.com/netsec-ethz/scion-apps/pkg/appnet"
	"github.com/netsec-ethz/scion-apps/pkg/appnet/appquic"
	"github.com/scionproto/scion/go/lib/snet"
)

type UDPDialer struct {
	raddr string
}

func NewUDPDialer(address string) (*UDPDialer, error) {
	return &UDPDialer{address}, nil
}

func (d *UDPDialer) Dial() (network.Connection, error) {
	conn, err := appnet.Dial(d.raddr)
	if err != nil {
		return nil, err
	}
	return connection.NewUDP(conn, conn.LocalAddr(), conn.RemoteAddr(), false), err
}

type UDPListener struct {
	laddr *snet.UDPAddr
}

func NewUDPListener(address string) (*UDPListener, error) {
	addr, err := appnet.ResolveUDPAddr(address)
	if err != nil {
		return nil, err
	}
	return &UDPListener{addr}, nil
}

func (l *UDPListener) Listen() (network.Connection, error) {
	conn, err := appnet.ListenPort(uint16(l.laddr.Host.Port))
	if err != nil {
		return &connection.UDP{}, err
	}
	return connection.NewUDP(conn, conn.LocalAddr(), conn.RemoteAddr(), false), nil
}

func (l *UDPListener) Stop() error {
	return nil
}

type QUICDialer struct {
	raddr string
}

func NewQUICDialer(address string) (*QUICDialer, error) {
	return &QUICDialer{address}, nil
}

func (d *QUICDialer) Dial() (network.Connection, error) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"taps-quic-test"},
	}
	conn, err := appquic.Dial(d.raddr, tlsConf, nil)
	if err != nil {
		return nil, err
	}
	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		return nil, err
	}
	return connection.NewQUIC(conn, stream, conn.LocalAddr(), conn.RemoteAddr()), nil
}

type QUICListener struct {
	listener quic.Listener
}

func NewQUICListener(address string) (*QUICListener, error) {
	addr, err := appnet.ResolveUDPAddr(address)
	if err != nil {
		return nil, err
	}
	tlsConf, err := generateTLSConfig()
	if err != nil {
		return nil, err
	}
	listener, err := appquic.ListenPort(uint16(addr.Host.Port), tlsConf, nil)
	if err != nil {
		return nil, err
	}
	return &QUICListener{listener}, nil
}

func (l *QUICListener) Listen() (network.Connection, error) {
	conn, err := l.listener.Accept(context.Background())
	if err != nil {
		return &connection.QUIC{}, err
	}
	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		return &connection.QUIC{}, err
	}
	return connection.NewQUIC(conn, stream, conn.LocalAddr(), conn.RemoteAddr()), nil
}

func (l *QUICListener) Stop() error {
	return nil
}

type scion struct{}

func (scion *scion) NewListener(e *network.Endpoint) (network.Listener, error) {
	switch e.Transport {
	case glob.TRANSPORT_UDP:
		return NewUDPListener(e.LocalAddress)
	case glob.TRANSPORT_QUIC:
		return NewQUICListener(e.LocalAddress)
	default:
		return nil, errors.New(fmt.Sprintf("Transport %s not implemented for SCION", e.Transport))
	}
}

func (scion *scion) NewDialer(e *network.Endpoint) (network.Dialer, error) {
	switch e.Transport {
	case glob.TRANSPORT_UDP:
		return NewUDPDialer(e.RemoteAddress)
	case glob.TRANSPORT_QUIC:
		return NewQUICDialer(e.RemoteAddress)
	default:
		return nil, errors.New(fmt.Sprintf("Transport %s not implemented for SCION", e.Transport))
	}
}

func Network() network.Network {
	return &scion{}
}

func generateTLSConfig() (*tls.Config, error) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return nil, err
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"taps-quic-test"},
	}, nil
}
