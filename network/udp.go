package network

import (
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
