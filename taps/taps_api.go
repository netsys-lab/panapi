package taps

import (
	"crypto/rsa"
	"net"

	quic "github.com/lucas-clemente/quic-go"
)

//

const (
	SERV_NONE = 1
	SERV_TCP  = 2
	SERV_QUIC = 3

	KEYPAIR = "keypair"
)

var (
	SERV_NAMES = []string{"none (invalid)", "none", "tcp", "quic"}
)

//

type Endpoint struct {
	interfaceName string
	serviceType   int
	port          string
	ipv4address   string
}

type LocalEndpoint struct {
	Endpoint
}

type RemoteEndpoint struct {
	Endpoint
	hostName string
}

type TransportProperties struct{}

type SecurityParameters struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

type Preconnection struct {
	locEnd    *LocalEndpoint
	remEnd    *RemoteEndpoint
	transProp *TransportProperties
	secParam  *SecurityParameters
}

type Listener struct {
	nlis    net.Listener
	qlis    quic.Listener
	preconn *Preconnection
	ConnRec chan Connection
	active  bool
}

type Connection struct {
	nconn   net.Conn
	qconn   quic.Session
	preconn *Preconnection
	active  bool
	Err     error
}

type Message struct {
	Data    string
	Context string
}

// type interf interface {
// 	(endPo *Endpoint) WithInterface(interfaceName string)
// 	(endPo *Endpoint) WithPort(port string)
// 	(endPo *Endpoint) WithIPv4Address(addr string)
// 	(endPo *Endpoint) WithService(serviceType string)
// 	(remEndPo *RemoteEndpoint) WithHostname(hostName string)
// 	(tranProp *TransportProperties) Require(method string)
// 	(secParam *SecurityParameters) Set(id string, list ...int)
// 	(preconn *Preconnection) Listen() *Listener
// 	(preconn *Preconnection) Initiate() *Connection
// 	(lis *Listener) Stop()
// 	(conn *Connection) Clone() *Connection
// 	(conn *Connection) Receive() *Message
// 	(conn *Connection) Send(msg *Message)
// 	(conn *Connection) Close()
// }
