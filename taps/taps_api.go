package taps

import (
	"crypto/rsa"
	"net"

	quic "github.com/lucas-clemente/quic-go"
	"github.com/scionproto/scion/go/lib/snet"
)

//

const (
	SERV_NONE  = 1
	SERV_TCP   = 2
	SERV_QUIC  = 3
	SERV_SCION = 4

	KEYPAIR   = "keypair"
	NAGLE_ON  = "nagle_on"
	NAGLE_OFF = "nagle_off"
)

var (
	SERV_NAMES = []string{"none (invalid)", "none", "tcp", "quic", "scion"}
)

//

type Endpoint struct {
	interfaceName string
	serviceType   int
	port          string
	ipv4address   string
	// ipv6address   string
	// scionAddress  string
}

type LocalEndpoint struct {
	Endpoint
}

type RemoteEndpoint struct {
	Endpoint
	hostName string
}

type TransportProperties struct {
	nagle bool
}

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
	sconn   *snet.Conn
	preconn *Preconnection
	active  bool
	Err     error
	saddr   net.Addr
}

type Message struct {
	Data    string
	Context string
}
