package ip

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
	"net"

	"github.com/lucas-clemente/quic-go"
	"github.com/netsys-lab/panapi/pkg/network"
)

type UDPDialer struct {
	raddr *net.UDPAddr
}

func NewUDPDialer(address string) (*UDPDialer, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}
	return &UDPDialer{addr}, nil
}

func (d *UDPDialer) Dial() (network.Connection, error) {
	conn, err := net.DialUDP("udp", nil, d.raddr)
	if err != nil {
		return nil, err
	}
	return network.NewUDP(conn, conn.LocalAddr(), conn.RemoteAddr(), true), nil
}

type UDPListener struct {
	laddr *net.UDPAddr
}

func NewUDPListener(address string) (*UDPListener, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}
	return &UDPListener{addr}, nil
}

func (l *UDPListener) Listen() (network.Connection, error) {
	conn, err := net.ListenUDP("udp", l.laddr)
	if err != nil {
		return &network.UDP{}, err
	}
	return network.NewUDP(conn, conn.LocalAddr(), conn.RemoteAddr(), false), nil
}

func (l *UDPListener) Stop() error {
	return nil
}

type TCPDialer struct {
	raddr *net.TCPAddr
}

func NewTCPDialer(address string) (*TCPDialer, error) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}
	return &TCPDialer{addr}, nil
}

func (d *TCPDialer) Dial() (network.Connection, error) {
	conn, err := net.DialTCP("tcp", nil, d.raddr)
	if err != nil {
		return nil, err
	}
	return network.NewTCP(conn, conn.LocalAddr(), conn.RemoteAddr()), nil
}

type TCPListener struct {
	listener *net.TCPListener
}

func NewTCPListener(address string) (*TCPListener, error) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &TCPListener{listener}, nil
}

func (l *TCPListener) Listen() (network.Connection, error) {
	conn, err := l.listener.AcceptTCP()
	if err != nil {
		return &network.TCP{}, err
	}
	return network.NewTCP(conn, conn.LocalAddr(), conn.RemoteAddr()), nil
}

func (l *TCPListener) Stop() error {
	return l.listener.Close()
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
	conn, err := quic.DialAddr(d.raddr, tlsConf, nil)
	if err != nil {
		return nil, err
	}
	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		return nil, err
	}
	return network.NewQUIC(conn, stream, conn.LocalAddr(), conn.RemoteAddr()), nil
}

type QUICListener struct {
	listener quic.Listener
}

func NewQUICListener(address string) (*QUICListener, error) {
	listener, err := quic.ListenAddr(address, generateTLSConfig(), nil)
	if err != nil {
		return nil, err
	}
	return &QUICListener{listener}, nil
}

func (l *QUICListener) Listen() (network.Connection, error) {
	conn, err := l.listener.Accept(context.Background())
	if err != nil {
		return &network.QUIC{}, err
	}
	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		return &network.QUIC{}, err
	}
	return network.NewQUIC(conn, stream, conn.LocalAddr(), conn.RemoteAddr()), err
}

func (l *QUICListener) Stop() error {
	return nil
}

type ip struct{}

func (ip *ip) NewListener(e *network.Endpoint) (network.Listener, error) {
	switch e.Transport {
	case network.TRANSPORT_UDP:
		return NewUDPListener(e.LocalAddress)
	case network.TRANSPORT_TCP:
		return NewTCPListener(e.LocalAddress)
	case network.TRANSPORT_QUIC:
		return NewQUICListener(e.LocalAddress)
	default:
		return nil, errors.New(fmt.Sprintf("Transport %s not implemented for IP", e.Transport))
	}
}

func (ip *ip) NewDialer(e *network.Endpoint) (network.Dialer, error) {
	switch e.Transport {
	case network.TRANSPORT_UDP:
		return NewUDPDialer(e.RemoteAddress)
	case network.TRANSPORT_TCP:
		return NewTCPDialer(e.RemoteAddress)
	case network.TRANSPORT_QUIC:
		return NewQUICDialer(e.RemoteAddress)
	default:
		return nil, errors.New(fmt.Sprintf("Transport %s not implemented for IP", e.Transport))
	}
}

func Network() network.Network {
	return &ip{}
}

func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"taps-quic-test"},
	}
}
