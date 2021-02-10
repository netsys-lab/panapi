package scion

import (
	"code.ovgu.de/hausheer/taps-api/connection"
	"code.ovgu.de/hausheer/taps-api/network"
	"errors"
	"fmt"
	"github.com/netsec-ethz/scion-apps/pkg/appnet"
	"net"
)

type UDPDialer struct {
	addr string
}

func NewUDPDialer(address string) (*UDPDialer, error) {
	return &UDPDialer{address}, nil
}

func (d *UDPDialer) Dial() (network.Connection, error) {
	conn, err := appnet.Dial(d.addr)
	if err != nil {
		return nil, err
	}
	return connection.NewUDP(conn, conn.LocalAddr(), conn.RemoteAddr()), err
}

type UDPListener struct {
	addr *net.UDPAddr
}

func NewUDPListener(address string) (*UDPListener, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}
	return &UDPListener{addr}, nil
}

func (l *UDPListener) Listen() (network.Connection, error) {
	conn, err := appnet.Listen(l.addr)
	if err != nil {
		return nil, err
	}
	return connection.NewUDP(conn, conn.LocalAddr(), conn.RemoteAddr()), err
}

type scion struct{}

func (scion *scion) NewListener(e *network.Endpoint) (network.Listener, error) {
	switch e.Transport {
	case "UDP":
		return NewUDPListener(e.LocalAddress)
	default:
		return nil, errors.New(fmt.Sprintf("Transport %s not implemented for SCION", e.Transport))
	}
}

func (scion *scion) NewDialer(e *network.Endpoint) (network.Dialer, error) {
	switch e.Transport {
	case "UDP":
		return NewUDPDialer(e.RemoteAddress)
	default:
		return nil, errors.New(fmt.Sprintf("Transport %s not implemented for SCION", e.Transport))
	}
}

func Network() network.Network {
	return &scion{}
}
