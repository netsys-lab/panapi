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
package panapi

import (
	"github.com/netsys-lab/panapi/network"
	"github.com/netsys-lab/panapi/network/ip"
	"github.com/netsys-lab/panapi/network/scion"
)

type Listener struct {
	listener           network.Listener
	ConnectionReceived chan network.Connection
}

func (l *Listener) Stop() error {
	return l.listener.Stop()
}

type Preconnection struct {
	endpoint *network.Endpoint
	listener network.Listener
	dialer   network.Dialer
}

func (p *Preconnection) Listen() Listener {
	c := make(chan network.Connection)
	go func(p *Preconnection, c chan network.Connection) {
		for {
			conn, err := p.listener.Listen()
			if err != nil {
				conn.SetError(err)
			} else {
				c <- conn
			}
		}
	}(p, c)
	return Listener{ConnectionReceived: c, listener: p.listener}
}

func (p *Preconnection) Initiate() (network.Connection, error) {
	return p.dialer.Dial()
}

func NewPreconnection(e *network.Endpoint, tp *network.TransportProperties) (Preconnection, error) {
	var (
		l      network.Listener
		dialer network.Dialer
		net    network.Network
		p      Preconnection
		err    error
	)

	switch e.Transport {
	case network.TRANSPORT_UDP:
	case network.TRANSPORT_TCP:
	case network.TRANSPORT_QUIC:
	default:
		return p, network.AddrTypeError
	}

	switch e.Network {
	case network.NETWORK_IP, network.NETWORK_IPV4, network.NETWORK_IPV6:
		net = ip.Network(tp)
	case network.NETWORK_SCION:
		net = scion.Network(tp)
	default:
		return p, network.NetTypeError
	}

	if e.Local {
		l, err = net.NewListener(e)
	} else {
		dialer, err = net.NewDialer(e)
	}
	if err != nil {
		return p, err
	}

	p = Preconnection{
		endpoint: e,
		listener: l,
		dialer:   dialer,
	}
	return p, err
}
