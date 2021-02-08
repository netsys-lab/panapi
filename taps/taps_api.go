package taps

import (
	"crypto/rsa"
	"net"

	quic "github.com/lucas-clemente/quic-go"
	"github.com/scionproto/scion/go/lib/snet"
)

//

const (
	TRANS_INVALID = iota
	TRANS_NONE
	TRANS_UDP
	TRANS_TCP
	TRANS_QUIC

	NET_INVALID
	NET_NONE
	NET_SCION
	NET_IP

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
	listener      net.Listener
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
