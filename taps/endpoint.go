package taps

import "net"

type Endpoint struct {
	Address   net.Addr
	Transport string
	Network   string
	Interface string
}

type LocalEndpoint struct{ Endpoint }

type RemoteEndpoint struct{ Endpoint }

func NewRemoteEndpoint() *RemoteEndpoint {
	return &RemoteEndpoint{}
}

func NewLocalEndpoint() *LocalEndpoint {
	return &LocalEndpoint{}
}

func (e *Endpoint) WithNetwork(network string) {
	e.Network = network
}

func (e *Endpoint) WithTransport(transport string) {
	e.Transport = transport
}

func (e *Endpoint) WithIPv4Address(addr net.IPAddr) {
	e.Address = &addr
}

func (e *Endpoint) WithIPv6Address(addr net.IPAddr) {
	e.Address = &addr
}

func (e *Endpoint) WithPort(port int) {

}

func (e *Endpoint) WithHostname(name string) {

}

func (e *Endpoint) WithInterface(intf string) {
	e.Interface = intf
}

func (e *Endpoint) WithProtocol(proto Protocol) {

}