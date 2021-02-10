package ip

import (
	"code.ovgu.de/hausheer/taps-api/connection"
	"code.ovgu.de/hausheer/taps-api/network"
	"errors"
	"fmt"
	"net"
	"sync"
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
	return connection.NewUDP(conn, conn.LocalAddr(), conn.RemoteAddr()), err
}

type UDPListener struct {
	laddr *net.UDPAddr
	mutex sync.Mutex
}

func NewUDPListener(address string) (*UDPListener, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}
	return &UDPListener{addr, sync.Mutex{}}, nil
}

func (l *UDPListener) Listen() (network.Connection, error) {
	l.mutex.Lock()
	//defer l.mutex.Unlock()
	//fmt.Printf("after %v", l.mutex)
	conn, err := net.ListenUDP("udp", l.laddr)
	if err != nil {
		return nil, err
	}
	return connection.NewUDP(conn, conn.LocalAddr(), conn.RemoteAddr()), err
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
	return connection.NewTCP(conn, conn.LocalAddr(), conn.RemoteAddr()), err

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
		return nil, err
	}
	return connection.NewTCP(conn, conn.LocalAddr(), conn.RemoteAddr()), err
}

type ip struct{}

func (ip *ip) NewListener(e *network.Endpoint) (network.Listener, error) {
	switch e.Transport {
	case "UDP":
		return NewUDPListener(e.LocalAddress)
	case "TCP":
		return NewTCPListener(e.LocalAddress)
	default:
		return nil, errors.New(fmt.Sprintf("Transport %s not implemented for IP", e.Transport))
	}
}

func (ip *ip) NewDialer(e *network.Endpoint) (network.Dialer, error) {
	switch e.Transport {
	case "UDP":
		return NewUDPDialer(e.RemoteAddress)
	case "TCP":
		return NewTCPDialer(e.RemoteAddress)
	default:
		return nil, errors.New(fmt.Sprintf("Transport %s not implemented for IP", e.Transport))
	}
}

func Network() network.Network {
	return &ip{}
}
