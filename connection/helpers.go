package connection

import (
	"net"
	"time"

	"code.ovgu.de/hausheer/taps-api/glob"
	"code.ovgu.de/hausheer/taps-api/network"
	"github.com/lucas-clemente/quic-go"
)

// conn includes the lowest common denominator of member functions of
// net.UDPConn and snet.Conn. This way, both the ip and the scion package
// can make use of the UDP helper
type conn interface {
	Write([]byte) (int, error)
	WriteTo([]byte, net.Addr) (int, error)
	ReadFrom([]byte) (int, net.Addr, error)
	Close() error
}

type UDP struct {
	conn  conn
	laddr net.Addr
	raddr net.Addr
	write bool
	err   error
}

func NewUDP(conn conn, laddr, raddr net.Addr, write bool) network.Connection {
	return &UDP{conn, laddr, raddr, write, nil}
}

func (c *UDP) Send(message network.Message) error {
	var err error
	if c.write {
		_, err = c.conn.Write([]byte(message.String()))
	} else {
		_, err = c.conn.WriteTo([]byte(message.String()), c.raddr)
	}
	return err
}

func (c *UDP) Receive() (network.Message, error) {
	var (
		m   glob.Message
		n   int
		err error
	)
	buffer := make([]byte, 1024)
	n, c.raddr, err = c.conn.ReadFrom(buffer)
	m = glob.Message(string(buffer[:n]))
	return &m, err
}

func (c *UDP) Close() error {
	return c.conn.Close()
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

func NewTCP(conn *net.TCPConn, laddr, raddr net.Addr) network.Connection {
	return &TCP{conn, laddr, raddr, nil}
}

func (c TCP) Send(message network.Message) error {
	_, err := c.conn.Write([]byte(message.String()))
	return err
}

func (c TCP) Receive() (network.Message, error) {
	var (
		m   glob.Message
		n   int
		err error
	)
	buffer := make([]byte, 1024)
	n, err = c.conn.Read(buffer)
	m = glob.Message(string(buffer[:n]))
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
	err    error
}

func NewQUIC(conn quic.Session, stream quic.Stream, laddr, raddr net.Addr) network.Connection {
	return &QUIC{conn, stream, laddr, raddr, nil}
}

func (c QUIC) Send(message network.Message) error {
	_, err := c.stream.Write([]byte(message.String()))
	return err
}

func (c QUIC) Receive() (network.Message, error) {
	var (
		m   glob.Message
		n   int
		err error
	)
	buffer := make([]byte, 1024)
	n, err = c.stream.Read(buffer)
	if err != nil {
		return nil, err
	}
	m = glob.Message(string(buffer[:n]))
	return &m, err
}

func (c *QUIC) Close() error {
	// todo
	time.Sleep(1000 * time.Microsecond)
	return c.conn.CloseWithError(0, "closed by server.")
}

func (c *QUIC) SetError(err error) {
	c.err = err
}

func (c *QUIC) GetError() error {
	return c.err
}
