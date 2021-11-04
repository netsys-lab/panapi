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
	_, err = c.Write([]byte(message.String()))
	return err
}

func (c *UDP) Receive() (Message, error) {
	var (
		m   DummyMessage
		n   int
		err error
	)
	buffer := make([]byte, 16*1024)
	n, err = c.Read(buffer)
	m = DummyMessage(string(buffer[:n]))
	return &m, err
}

var AddrFromRead net.Addr

// var Conn net.PacketConn
// var UdpErr error

func (c *UDP) Read(p []byte) (n int, err error) {

	// laddr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 1337}
	// if Conn == nil {
	// 	Conn, UdpErr = net.ListenPacket("udp", "127.0.0.1:1333")
	// }

	// if UdpErr != nil {
	// 	return 0, UdpErr
	// }

	// // conn.ReadFrom(p)

	// n, from, err := Conn.ReadFrom(p)
	// Conn.Close()
	// if err != nil {
	// 	return 0, err
	// }
	// data := p[:n]
	// fmt.Printf("Received %s: %s\n", from, data)
	// AddrFromRead = from

	// return n, nil

	if c.pconn != nil {
		n, c.raddr, err = c.pconn.ReadFrom(p)
		AddrFromRead = c.raddr
	} else if c.conn != nil {
		n, err = c.conn.Read(p)
	}
	return
}

func (c *UDP) Write(p []byte) (n int, err error) {
	if c.conn != nil {
		n, err = c.conn.Write(p)
	} else if c.pconn != nil && AddrFromRead != nil {
		// newCon, newConErr := net.Dial("udp", AddrFromRead.String())
		// if newConErr != nil {
		// 	return 0, newConErr
		// }
		n, err = c.pconn.WriteTo(p, AddrFromRead)
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
