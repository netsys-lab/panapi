package network

import (
	"io"
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
	_, err := io.Copy(c.conn, message)
	return err
}

func (c *TCP) Receive(m Message) error {
	_, err := io.Copy(m, c.conn)
	return err
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
