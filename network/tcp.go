package network

import (
	"net"
)

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

func (c *TCP) Read(p []byte) (int, error) {
	return c.conn.Read(p)
}

func (c *TCP) Write(p []byte) (int, error) {
	return c.conn.Write(p)
}

func (c *TCP) Close() error {
	return c.conn.Close()
}

func (c *TCP) LocalAddr() net.Addr {
	return c.laddr
}

func (c *TCP) RemoteAddr() net.Addr {
	return c.raddr
}

func (c *TCP) SetError(err error) {
	c.err = err
}

func (c *TCP) GetError() error {
	return c.err
}
