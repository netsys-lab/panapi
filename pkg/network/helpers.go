package network

import (
	"net"
	"time"

	"github.com/lucas-clemente/quic-go"
)

var timeout = 1 * time.Second
var timeStep = 10 * time.Millisecond

/**
 * conn includes the lowest common denominator of member
 * functions of net.UDPConn and snet.Conn. This way, both the
 * ip and the scion package can make use of the UDP helper.

type conn interface {
	Write([]byte) (int, error)
	WriteTo([]byte, net.Addr) (int, error)
	ReadFrom([]byte) (int, net.Addr, error)
	Close() error
        }*/

// TODO, placeholder stub implementation for message
type DummyMessage string

func (m DummyMessage) String() string {
	return string(m)
}

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
	var err error
	if c.conn != nil {
		_, err = c.conn.Write([]byte(message.String()))
	} else if c.pconn != nil {
		_, err = c.pconn.WriteTo([]byte(message.String()), c.raddr)
	}
	return err
}

func (c *UDP) Receive() (Message, error) {
	var (
		m   DummyMessage
		n   int
		err error
	)
	buffer := make([]byte, 16*1024)
	if c.pconn != nil {
		n, c.raddr, err = c.pconn.ReadFrom(buffer)
	} else if c.conn != nil {
		n, err = c.conn.Read(buffer)
	}
	m = DummyMessage(string(buffer[:n]))
	return &m, err
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

type TCP struct {
	conn  *net.TCPConn
	laddr net.Addr
	raddr net.Addr
	err   error
}

func NewTCP(conn *net.TCPConn, laddr, raddr net.Addr) Connection {
	return &TCP{conn, laddr, raddr, nil}
}

func (c *TCP) Send(message Message) error {
	_, err := c.conn.Write([]byte(message.String()))
	return err
}

func (c *TCP) Receive() (Message, error) {
	var (
		m   DummyMessage
		n   int
		err error
	)
	buffer := make([]byte, 1024)
	n, err = c.conn.Read(buffer)
	m = DummyMessage(string(buffer[:n]))
	return &m, err
}

func (c *TCP) Close() error {
	return c.conn.Close()
}

func (c *TCP) SetError(err error) {
	c.err = err
}

func (c *TCP) GetError() error {
	return c.err
}

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
