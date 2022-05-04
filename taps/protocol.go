package taps

type Listener interface {
	Accept() (Connection, error)
	Close() error
	//Addr() net.Addr
}

type Protocol interface {
	Satisfy(*Preconnection) (*TransportProperties, error)
	NewListener(*Preconnection) (Listener, error)
	Initiate(*Preconnection) (Connection, error)
	Selector() Selector
}
