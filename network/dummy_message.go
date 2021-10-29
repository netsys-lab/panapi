package network

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
