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
	"time"

	"github.com/lucas-clemente/quic-go"
)

var timeout = 1 * time.Second
var timeStep = 10 * time.Millisecond

type QUIC struct {
	conn   quic.Session
	stream quic.Stream
	laddr  net.Addr
	raddr  net.Addr
	last   time.Time
	err    error
}

func NewQUIC(conn quic.Session, stream quic.Stream, laddr, raddr net.Addr) Connection {
	return &QUIC{conn, stream, laddr, raddr, time.Now().Add(-1 * timeout), nil}
}

func (c *QUIC) Send(message Message) error {
	_, err := io.Copy(c.stream, message)
	c.last = time.Now()
	return err
}

func (c *QUIC) Receive(message Message) error {
	_, err := io.Copy(message, c.stream)
	if err == EOM {
		// End Of Message, not an error we need to propagate beyond this point
		err = nil
	}
	return err
}

func (c *QUIC) Read(p []byte) (n int, err error) {
	return c.stream.Read(p)
}

func (c *QUIC) Write(p []byte) (n int, err error) {
	return c.stream.Write(p)
}

func (c *QUIC) Close() error {
	for time.Now().Sub(c.last) < timeout {
		time.Sleep(timeStep)
	}
	return c.conn.CloseWithError(0, "connection closed.")
}

func (c *QUIC) LocalAddr() net.Addr {
	return c.laddr
}

func (c *QUIC) RemoteAddr() net.Addr {
	return c.raddr
}

func (c *QUIC) SetError(err error) {
	c.err = err
}

func (c *QUIC) GetError() error {
	return c.err
}
