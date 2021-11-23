package scion

import (
	"context"
	"net"

	"github.com/netsec-ethz/scion-apps/pkg/pan"
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
