// Copyright 2021 Thorben Krüger (thorben.krueger@ovgu.de)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package scion

import (
	"context"

	"github.com/netsec-ethz/scion-apps/pkg/pan"
	"github.com/netsys-lab/panapi/network"
	"inet.af/netaddr"
)

type UDPDialer struct {
	raddr pan.UDPAddr
}

func NewUDPDialer(address string) (*UDPDialer, error) {
	addr, err := pan.ResolveUDPAddr(address)
	return &UDPDialer{addr}, err
}

func (d *UDPDialer) Dial() (network.Connection, error) {
	conn, err := pan.DialUDP(context.Background(), netaddr.IPPort{}, d.raddr, nil, nil)
	if err != nil {
		return nil, err
	}
	return network.NewUDP(conn, nil, conn.LocalAddr(), conn.RemoteAddr()), err
}

type UDPListener struct {
	laddr netaddr.IPPort
}

func NewUDPListener(address string) (*UDPListener, error) {
	addr, err := pan.ResolveUDPAddr(address)
	if err != nil {
		return nil, err
	}
	return &UDPListener{netaddr.IPPortFrom(addr.IP, addr.Port)}, nil
}

func (l *UDPListener) Listen() (network.Connection, error) {
	pconn, err := pan.ListenUDP(context.Background(), l.laddr, nil)
	if err != nil {
		return &network.UDP{}, err
	}
	return network.NewUDP(nil, pconn, pconn.LocalAddr(), nil), nil
}

func (l *UDPListener) Stop() error {
	return nil
}
