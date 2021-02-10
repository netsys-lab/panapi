package network

type Endpoint struct {
	Local         bool
	LocalAddress  string
	RemoteAddress string
	Transport     string
	Network       string
}

func NewRemoteEndpoint() *Endpoint {
	return &Endpoint{Local: false}
}

func NewLocalEndpoint() *Endpoint {
	return &Endpoint{Local: true}
}

func (e *Endpoint) WithNetwork(network string) {
	e.Network = network
}

func (e *Endpoint) WithTransport(transport string) {
	e.Transport = transport
}

func (e *Endpoint) WithAddress(addr string) {
	if e.Local {
		e.LocalAddress = addr
	} else {
		e.RemoteAddress = addr
	}
}
