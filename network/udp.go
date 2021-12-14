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
package network

import (
	"io"
	"net"
)

type UDP struct {
	conn  net.Conn
	pconn net.PacketConn
	laddr net.Addr
	raddr net.Addr
	//write bool
	err error
}

func NewUDP(conn net.Conn, pconn net.PacketConn, laddr, raddr net.Addr) Connection {
	return &UDP{conn, pconn, laddr, raddr, nil}
}

func (c *UDP) Send(message Message) error {
	_, err := io.Copy(c, message)
	return err
}

func (c *UDP) Receive(m Message) error {
	_, err := io.Copy(m, c)
	return err
}

func (c *UDP) Read(p []byte) (n int, err error) {
	if c.pconn != nil {
		n, c.raddr, err = c.pconn.ReadFrom(p)
	} else if c.conn != nil {
		n, err = c.conn.Read(p)
	}
	return
}

func (c *UDP) Write(p []byte) (n int, err error) {
	if c.conn != nil {
		n, err = c.conn.Write(p)
	} else if c.pconn != nil {
		n, err = c.pconn.WriteTo(p, c.raddr)
	}
	return
}

func (c *UDP) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	if c.pconn != nil {
		return c.pconn.Close()
	}
	return nil
}

func (c *UDP) SetError(err error) {
	c.err = err
}

func (c *UDP) GetError() error {
	return c.err
}

func (c *UDP) LocalAddr() net.Addr {
	return c.laddr
}

func (c *UDP) RemoteAddr() net.Addr {
	return c.raddr
}
