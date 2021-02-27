package network

type Message interface {
	String() string
}

type Connection interface {
	Send(Message) error
	Receive() (Message, error)
	Close() error
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
