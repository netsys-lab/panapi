package scion

import (
	"context"
	//"crypto/rand"
	//"crypto/rsa"
	"crypto/tls"
	//"crypto/x509"
	//"encoding/pem"
	"errors"
	"fmt"
	"log"
	//"math/big"
	"net"

	"github.com/lucas-clemente/quic-go"
	"github.com/netsec-ethz/scion-apps/pkg/pan"
	"github.com/netsec-ethz/scion-apps/pkg/quicutil"
	"github.com/netsys-lab/panapi/network"
)

type UDPDialer struct {
	raddr pan.UDPAddr
}

func NewUDPDialer(address string) (*UDPDialer, error) {
	addr, err := pan.ResolveUDPAddr(address)
	return &UDPDialer{addr}, err
}

func (d *UDPDialer) Dial() (network.Connection, error) {
	conn, err := pan.DialUDP(context.Background(), nil, d.raddr, nil, nil)
	if err != nil {
		return nil, err
	}
	return network.NewUDP(conn, nil, conn.LocalAddr(), conn.RemoteAddr()), err
}

type UDPListener struct {
	laddr net.UDPAddr
}

func NewUDPListener(address string) (*UDPListener, error) {
	addr, err := pan.ResolveUDPAddr(address)
	if err != nil {
		return nil, err
	}
	return &UDPListener{net.UDPAddr{Port: addr.Port}}, nil
}

func (l *UDPListener) Listen() (network.Connection, error) {
	pconn, err := pan.ListenUDP(context.Background(), &l.laddr, nil)
	if err != nil {
		return &network.UDP{}, err
	}
	return network.NewUDP(nil, pconn, pconn.LocalAddr(), nil), nil
}

func (l *UDPListener) Stop() error {
	return nil
}

type QUICDialer struct {
	raddr pan.UDPAddr
}

func NewQUICDialer(address string) (*QUICDialer, error) {
	addr, err := pan.ResolveUDPAddr(address)
	return &QUICDialer{addr}, err
}

func (d *QUICDialer) Dial() (network.Connection, error) {
	tlsConf := &tls.Config{
		//Certificates: quicutil.MustGenerateSelfSignedCert(),
		InsecureSkipVerify: true,
		NextProtos:         []string{"panapi-quic-test"},
	}
	conn, err := pan.DialQUIC(context.Background(), nil, d.raddr, nil, nil, "", tlsConf, nil)
	if err != nil {
		return nil, err
	}
	stream, err := conn.OpenStream() //Sync(context.Background())
	if err != nil {
		return nil, err
	}
	return network.NewQUIC(conn, stream, conn.LocalAddr(), conn.RemoteAddr()), nil
}

type QUICListener struct {
	listener quic.Listener
}

func NewQUICListener(address string) (*QUICListener, error) {
	addr, err := pan.ResolveUDPAddr(address)
	if err != nil {
		return nil, err
	}
	tlsConf := &tls.Config{
		Certificates: quicutil.MustGenerateSelfSignedCert(),
		//InsecureSkipVerify: true,
		NextProtos: []string{"panapi-quic-test"},
	}
	//tlsConf, err := generateTLSConfig()
	if err != nil {
		return nil, err
	}
	listener, err := pan.ListenQUIC(context.Background(), &net.UDPAddr{Port: addr.Port}, nil, tlsConf, nil)
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
	log.Println("accepted stream")
	return network.NewQUIC(conn, stream, conn.LocalAddr(), conn.RemoteAddr()), nil
}

func (l *QUICListener) Stop() error {
	return nil
}

type scion struct{}

func (scion *scion) NewListener(e *network.Endpoint) (network.Listener, error) {
	switch e.Transport {
	case network.TRANSPORT_UDP:
		return NewUDPListener(e.LocalAddress)
	case network.TRANSPORT_QUIC:
		return NewQUICListener(e.LocalAddress)
	default:
		return nil, errors.New(fmt.Sprintf("Transport %s not implemented for SCION", e.Transport))
	}
}

func (scion *scion) NewDialer(e *network.Endpoint) (network.Dialer, error) {
	switch e.Transport {
	case network.TRANSPORT_UDP:
		return NewUDPDialer(e.RemoteAddress)
	case network.TRANSPORT_QUIC:
		return NewQUICDialer(e.RemoteAddress)
	default:
		return nil, errors.New(fmt.Sprintf("Transport %s not implemented for SCION", e.Transport))
	}
}

func Network() network.Network {
	return &scion{}
}

/*func generateTLSConfig() (*tls.Config, error) {
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
        }*/