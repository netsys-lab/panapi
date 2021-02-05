package taps

import (
	"crypto/rsa"
	"net"

	quic "github.com/lucas-clemente/quic-go"
	"github.com/scionproto/scion/go/lib/snet"
)

//

const (
	// serviceType
	SERV_INVALID int = 0
	SERV_NONE    int = 1
	SERV_TCP     int = 2
	SERV_QUIC    int = 3
	SERV_SCION   int = 4

	KEYPAIR   string = "keypair"
	NAGLE_ON  string = "nagle_on"
	NAGLE_OFF string = "nagle_off"
)

var (
	SERV_NAMES = []string{"invalid", "none", "tcp", "quic", "scion"}
)

//

type Endpoint struct {
	interfaceName string
	serviceType   int
	port          string
	address       string
	hostname      string
	// ipv4address   string
	// ipv6address   string
	// scionAddress  string
}

type LocalEndpoint struct {
	Endpoint
}

type RemoteEndpoint struct {
	Endpoint
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