package network

type Message interface {
	String() string
}

type Connection interface {
	Send(Message) error
	Receive() (Message, error)
}

type Dialer interface {
	Dial() (Connection, error)
}

type Listener interface {
	Listen() (Connection, error)
}

type Preconnection interface {
	Listen() chan Connection
	Initiate() Connection
}

type Network interface {
	NewListener(*Endpoint) (Listener, error)
	NewDialer(*Endpoint) (Dialer, error)
}
