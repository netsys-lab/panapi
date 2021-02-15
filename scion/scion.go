package scion

import (
	"errors"
	"fmt"
	"net"

	"code.ovgu.de/hausheer/taps-api/connection"
	"code.ovgu.de/hausheer/taps-api/network"
	"github.com/netsec-ethz/scion-apps/pkg/appnet"
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
	return connection.NewUDP(conn, conn.LocalAddr(), conn.RemoteAddr()), err
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
	conn, err := appnet.Listen(l.laddr)
	if err != nil {
		return nil, err
	}
	return connection.NewUDP(conn, conn.LocalAddr(), conn.RemoteAddr()), err
}

// type QUICDialer struct {
// 	raddr *snet.UDPAddr
// }

// func NewQUICDialer(address string) (*QUICDialer, error) {
// 	addr, err := appnet.ResolveUDPAddr(address)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &QUICDialer{addr}, nil
// }

// func (d *QUICDialer) Dial() (network.Connection, error) {
// 	conn, err := appquic.DialAddr(d.raddr, string(d.raddr.Host.IP), nil, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return connection.NewQUIC(conn, nil, conn.LocalAddr(), conn.RemoteAddr()), err
// }

// type QUICListener struct {
// 	listener quic.Listener
// }

// func NewQUICListener(address string) (*QUICListener, error) {
// 	addr, err := appnet.ResolveUDPAddr(address)
// 	if err != nil {
// 		return nil, err
// 	}
// 	listener, err := appquic.ListenPort(uint16(addr.Host.Port), nil, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &QUICListener{listener}, nil
// }

// func (l *QUICListener) Listen() (network.Connection, error) {
// 	conn, err := l.listener.Accept(context.Background())
// 	if err != nil {
// 		return nil, err
// 	}
// 	return connection.NewQUIC(conn, nil, conn.LocalAddr(), conn.RemoteAddr()), err
// }

type scion struct{}

func (scion *scion) NewListener(e *network.Endpoint) (network.Listener, error) {
	switch e.Transport {
	// case taps.TRANSPORT_UDP:
	case "UDP":
		return NewUDPListener(e.LocalAddress)
	// case taps.TRANSPORT_QUIC:
	// case "QUIC":
	// 	return NewQUICListener(e.LocalAddress)
	default:
		return nil, errors.New(fmt.Sprintf("Transport %s not implemented for SCION", e.Transport))
	}
}

func (scion *scion) NewDialer(e *network.Endpoint) (network.Dialer, error) {
	switch e.Transport {
	// case taps.TRANSPORT_UDP:
	case "UDP":
		return NewUDPDialer(e.RemoteAddress)
	// case taps.TRANSPORT_QUIC:
	// case "QUIC":
	// 	return NewQUICDialer(e.LocalAddress)
	default:
		return nil, errors.New(fmt.Sprintf("Transport %s not implemented for SCION", e.Transport))
	}
}

func Network() network.Network {
	return &scion{}
}
