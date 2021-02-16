package connection

import (
	"net"

	"code.ovgu.de/hausheer/taps-api/network"
	"github.com/lucas-clemente/quic-go"
)

//conn includes the lowest common denominator of member functions of
//net.UDPConn and snet.Conn. This way, both the ip and the scion package
//can make use of the UDP helper
type conn interface {
	Write([]byte) (int, error)
	WriteTo([]byte, net.Addr) (int, error)
	ReadFrom([]byte) (int, net.Addr, error)
}

type message string

func (m message) String() string {
	return string(m)
}

type UDP struct {
	conn  conn
	laddr net.Addr
	raddr net.Addr
	write bool
}

func NewUDP(conn conn, laddr, raddr net.Addr, write bool) network.Connection {
	return &UDP{conn, laddr, raddr, write}
}

func (c *UDP) Send(message network.Message) error {
	var err error
	if c.write {
		// fmt.Println("udp Write start : raddr =", c.raddr)
		_, err = c.conn.Write([]byte(message.String()))
		// fmt.Println("udp Write end")
	} else {
		// fmt.Println("udp WriteTo start : raddr=", c.raddr)
		_, err = c.conn.WriteTo([]byte(message.String()), c.raddr)
		// fmt.Println("udp WriteTo end")
	}
	return err
}

func (c *UDP) Receive() (network.Message, error) {
	var (
		m   message
		n   int
		err error
	)
	buffer := make([]byte, 1024)
	// fmt.Println("udp ReadFrom start")
	n, c.raddr, err = c.conn.ReadFrom(buffer)
	// fmt.Println("udp ReadFrom end : raddr = ", c.raddr.String())
	m = message(string(buffer[:n]))
	return &m, err
}

type TCP struct {
	conn  *net.TCPConn
	laddr net.Addr
	raddr net.Addr
}

func NewTCP(conn *net.TCPConn, laddr, raddr net.Addr) network.Connection {
	return &TCP{conn, laddr, raddr}
}

func (c TCP) Send(message network.Message) error {
	_, err := c.conn.Write([]byte(message.String()))
	return err
}

func (c TCP) Receive() (network.Message, error) {
	var (
		m   message
		n   int
		err error
	)
	buffer := make([]byte, 1024)
	n, err = c.conn.Read(buffer)
	m = message(string(buffer[:n]))
	return &m, err
}

type QUIC struct {
	conn   quic.Session
	stream quic.Stream
	laddr  net.Addr
	raddr  net.Addr
}

func NewQUIC(conn quic.Session, stream quic.Stream, laddr, raddr net.Addr) network.Connection {
	return &QUIC{conn, stream, laddr, raddr}
}

func (c QUIC) Send(message network.Message) error {
	// fmt.Println("quic Write start : local:", c.conn.LocalAddr(), "remote:", c.conn.RemoteAddr())
	c.stream.Write([]byte(message.String()))
	// fmt.Println("quic Write end")
	return nil
}

func (c QUIC) Receive() (network.Message, error) {
	var (
		m   message
		n   int
		err error
	)
	buffer := make([]byte, 1024)
	// fmt.Println("quic Read start : local:", c.conn.LocalAddr(), "remote:", c.conn.RemoteAddr())
	n, err = c.stream.Read(buffer)
	// fmt.Println("quic Read end")
	if err != nil {
		return nil, err
	}
	m = message(string(buffer[:n]))
	return &m, err
}
