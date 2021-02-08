package taps

import (
	"crypto/rsa"
	"flag"
	"net"
	"os/exec"
	"reflect"
	"strconv"

	"github.com/netsec-ethz/scion-apps/pkg/appnet"
)

//
// Setup
//

func Init() (*string, *string, *string, *string) {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()

	servF := flag.String("serv", "tcp", "tcp or quic or scion")
	addrF := flag.String("addr", "127.0.0.1", "address")
	portF := flag.String("port", "1111", "port")
	interF := flag.String("inter", "any", "interface name")
	flag.Parse()

	return servF, addrF, portF, interF
}

//
// Endpoint
//

func (endPo *Endpoint) WithInterface(interfaceName string) error {
	endPo.interfaceName = interfaceName
	var err error = nil
	if interfaceName != "any" {
		_, err = net.InterfaceByName(interfaceName)
	}
	return err
}

// todo
func (endPo *Endpoint) WithNetType(serviceType string) error {
	switch serviceType {
	case "quic":
		endPo.listener = net.Listener
		// endPo.netType = SERV_QUIC
	case "scion":
		endPo.listener = nil

	default:
		return &tapsError{
			Op:   "WithService",
			Endp: endPo,
			Desc: serviceType,
			Err:  errUnknownServiceType}
	}
	return nil
}

func (endPo *Endpoint) WithPort(port string) error {
	_, err := strconv.ParseUint(port, 10, 0)
	if err != nil {
		return &tapsError{
			Op:   "WithPort",
			Endp: endPo,
			Desc: port,
			Err:  err}
	}
	if len(port) > 5 {
		return &tapsError{
			Op:   "WithPort",
			Endp: endPo,
			Desc: port,
			Err:  errInvalidPort}
	}
	endPo.port = port
	return nil
}

func (endPo *Endpoint) WithIPv4Address(addr string) error {
	if net.ParseIP(addr) == nil {
		return &tapsError{
			Op:   "WithIPv4Address",
			Endp: endPo,
			Desc: addr,
			Err:  errInvalidIPAddress}
	}
	endPo.address = addr
	return nil
}

func (endPo *Endpoint) WithIPv6Address(addr string) error {
	if net.ParseIP(addr) == nil {
		return &tapsError{
			Op:   "WithIPv6Address",
			Endp: endPo,
			Desc: addr,
			Err:  errInvalidIPAddress}
	}
	endPo.address = addr
	return nil
}

func (endPo *Endpoint) WithScionAddress(addr string) error {
	_, err := appnet.ResolveUDPAddrAt(addr, appnet.DefaultResolver())
	if err != nil {
		return &tapsError{
			Op:   "WithScionAddress",
			Endp: endPo,
			Desc: addr,
			Err:  err}
	}
	endPo.address = addr
	return nil
}

func (endPo *Endpoint) WithHostname(hostname string) error {
	endPo.hostname = hostname
	ips, err := net.LookupIP(hostname)
	if err != nil {
		return &tapsError{
			Op:   "WithHostname",
			Endp: endPo,
			Err:  err}
	}
	endPo.address = ips[0].String()
	return nil
}

func (endPo *Endpoint) WithAddress(addr string) error {
	for i := 0; i < len(addr); i++ {
		if addr[i] == ',' {
			return endPo.WithScionAddress(addr)
		}
	}
	for i := 0; i < len(addr); i++ {
		switch addr[i] {
		case '.':
			return endPo.WithIPv4Address(addr)
		case ':':
			return endPo.WithIPv6Address(addr)
		}
	}
	return &tapsError{
		Op:   "WithAddress",
		Desc: addr,
		Err:  errInvalidAddressType}
}

//
// TransportProperties
//

func (tranProp *TransportProperties) Require(name string) error {
	switch name {
	case NAGLE_ON:
		tranProp.nagle = true
	case NAGLE_OFF:
		tranProp.nagle = false
	default:
		return &tapsError{
			Op:   "Set",
			Desc: name,
			Err:  errUnknownRequireName}
	}
	return nil
}

//
// SecurityParameters
//

func (secParam *SecurityParameters) Set(name string, args ...interface{}) error {
	switch name {
	case KEYPAIR:
		for i, arg := range args {
			switch arg.(type) {
			case *rsa.PrivateKey:
				secParam.privateKey = arg.(*rsa.PrivateKey)
			case *rsa.PublicKey:
				secParam.publicKey = arg.(*rsa.PublicKey)
			default:
				return &tapsError{
					Op:     "Set",
					ArgNum: i + 1,
					Desc:   reflect.TypeOf(arg).String(),
					Err:    errInvalidArgument}
			}
		}
	default:
		return &tapsError{
			Op:   "Set",
			Desc: name,
			Err:  errUnknownSetName}
	}
	return nil
}

//
// Preconnection
//

func (preconn *Preconnection) Listen() (*Listener, error) {
	servType, err := preconn.getServiceType()
	if err != nil {
		return nil, &tapsError{
			Op:   "Listen",
			Endp: preconn,
			Err:  err}
	}
	var lis *Listener
	switch transType {
	case SERV_TCP:
		lis, err = preconn.tpcListen()
	case SERV_QUIC:
		lis, err = preconn.quicListen()
	case NET_SCION:
		switch netType {
		case TRANS_UDP:
			lis, err = preconn.scionListen()
		case TRANS_QUIC:
			return nil, &tapsError{
				Op:   "Listen",
				Endp: preconn,
				Err:  errNotImplemented}
		}
	}
	return lis, err
}

func (preconn *Preconnection) Initiate() (*Connection, error) {
	servType, err := preconn.getServiceType()
	if err != nil {
		return nil, &tapsError{
			Op:   "Initiate",
			Endp: preconn,
			Err:  err}
	}
	var conn *Connection
	switch servType {
	case SERV_TCP:
		conn, err = preconn.tcpInitiate()
	case SERV_QUIC:
		conn, err = preconn.quicInitiate()
	case SERV_SCION:
		conn, err = preconn.scionInitiate()
	}
	return conn, err
}

//
// Listener
//

func (lis *Listener) Stop() error {
	var err error = nil
	if lis.isOpen() {
		lis.active = false
		lis.ConnRec = nil
		servType, err := lis.preconn.getServiceType()
		if err != nil {
			return &tapsError{
				Op:   "Stop",
				Endp: lis.preconn,
				Err:  err}
		}
		switch servType {
		case SERV_TCP:
			err = lis.tcpStop()
		case SERV_QUIC:
			err = lis.quicStop()
		case SERV_SCION:
			err = lis.scionStop()
		}
	}
	return err
}

//
// Connection
//

func (conn *Connection) Clone() *Connection {
	return &Connection{
		conn.nconn,
		conn.qconn,
		conn.sconn,
		conn.preconn,
		conn.active,
		nil,
		conn.saddr}
}

func (conn *Connection) Receive() (*Message, error) {
	servType, err := conn.preconn.getServiceType()
	if err != nil {
		return nil, &tapsError{
			Op:   "Receive",
			Endp: conn.preconn,
			Err:  err}
	}
	var ret *Message
	switch servType {
	case SERV_TCP:
		ret, err = conn.tcpReceive()
	case SERV_QUIC:
		ret, err = conn.quicReceive()
	case SERV_SCION:
		ret, err = conn.scionReceive()
	}
	return ret, err
}

func (conn *Connection) Send(msg *Message) error {
	servType, err := conn.preconn.getServiceType()
	if err != nil {
		return &tapsError{
			Op:   "Send",
			Endp: conn.preconn,
			Err:  err}
	}
	switch servType {
	case SERV_TCP:
		err = conn.tcpSend(msg)
	case SERV_QUIC:
		err = conn.quicSend(msg)
	case SERV_SCION:
		err = conn.scionSend(msg)
	}
	return err
}

func (conn *Connection) Close() error {
	var err error = nil
	if conn.isOpen() {
		conn.active = false
		servType, err := conn.preconn.getServiceType()
		if err != nil {
			return &tapsError{
				Op:   "Close",
				Endp: conn.preconn,
				Err:  err}
		}
		switch servType {
		case SERV_TCP:
			err = conn.tcpClose()
		case SERV_QUIC:
			err = conn.quicClose()
		case SERV_SCION:
			err = conn.scionClose()
		}
	}
	return err
}

//

//
// unexported funcs
//

//

func (conn *Connection) isOpen() bool {
	return conn.active
}

func (lis *Listener) isOpen() bool {
	return lis.active
}

func (preconn *Preconnection) getServiceType() (int, error) {
	if preconn.locEnd == nil {
		if preconn.remEnd != nil {
			return preconn.remEnd.serviceType, nil
		}
	}
	if preconn.remEnd == nil {
		if preconn.locEnd != nil {
			return preconn.locEnd.serviceType, nil
		}
	}
	return SERV_INVALID, &tapsError{
		Op:   "getServiceType",
		Endp: preconn,
		Err:  errNoServiceType}
}
