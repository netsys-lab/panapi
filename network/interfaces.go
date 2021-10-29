package network

import (
	"io"
	"net"
)

type Message interface {
	String() string
}

type Connection interface {
	io.ReadWriteCloser
	Send(Message) error
	Receive() (Message, error)
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	SetError(error)
	GetError() error
}

type Dialer interface {
	Dial() (Connection, error)
}

type Listener interface {
	Listen() (Connection, error)
	Stop() error
}

type Preconnection interface {
	Listen() (Listener, error)
	Initiate() (Connection, error)
}

type Network interface {
	NewListener(*Endpoint) (Listener, error)
	NewDialer(*Endpoint) (Dialer, error)
}
