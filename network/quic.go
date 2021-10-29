package network

import (
	"github.com/lucas-clemente/quic-go"
	"net"
	"time"
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
	_, err := c.stream.Write([]byte(message.String()))
	c.last = time.Now()
	return err
}

func (c *QUIC) Receive() (Message, error) {
	var (
		m   DummyMessage
		n   int
		err error
	)
	buffer := make([]byte, 1024)
	n, err = c.stream.Read(buffer)
	m = DummyMessage(string(buffer[:n]))
	return &m, err
}

func (c *QUIC) Close() error {
	for time.Now().Sub(c.last) < timeout {
		time.Sleep(timeStep)
	}
	return c.conn.CloseWithError(0, "connection closed.")
}

func (c *QUIC) SetError(err error) {
	c.err = err
}

func (c *QUIC) GetError() error {
	return c.err
}