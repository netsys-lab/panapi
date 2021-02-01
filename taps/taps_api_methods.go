package taps

import (
	"crypto/rsa"
	"errors"
	"net"
	"strconv"
)

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

func (endPo *Endpoint) WithPort(port string) error {
	_, err := strconv.Atoi(port)
	if err != nil {
		return &tapsError{Op: "WithPort", Port: port, Err: err}
	}
	if len(port) > 5 {
		return &tapsError{Op: "WithPort", Port: port, Err: errInvalidPort}
	}
	endPo.port = port
	return nil
}

func (endPo *Endpoint) WithIPv4Address(addr string) error {
	// if net.ParseIP(addr) == nil {
	// 	return &tapsError{Op: "WithIPv4Address", Ipv4address: addr, Err: errInvalidIPAddress}
	// }
	endPo.ipv4address = addr
	return nil
}

func (endPo *Endpoint) WithAddress(addr, addrType string) error {
	switch addrType {
	case "IPv4":
	case "ipv4":

	case "IPv6":
	case "ipv6":

	case "scion":

	default:
		return &tapsError{Op: "WithAddress", Ipv4address: addr, Err: nil}
	}
	endPo.ipv4address = addr
	return nil
}

func (endPo *Endpoint) WithService(serviceType string) error {
	switch serviceType {
	case "tcp":
		endPo.serviceType = SERV_TCP
	case "quic":
		endPo.serviceType = SERV_QUIC
	case "scion":
		endPo.serviceType = SERV_SCION
	default:
		return &tapsError{Op: "WithService", ServiceTypeInvalid: serviceType, Err: errUnknownServiceType}
	}
	return nil
}

func (remEndPo *RemoteEndpoint) WithHostname(hostName string) error {
	remEndPo.hostName = hostName
	return nil
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
		return &tapsError{Op: "Set", SetName: name, Err: errUnknownSetName}
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
				return &tapsError{Op: "Set", ArgNum: i, Err: errInvalidArgument}
			}
		}
	default:
		return &tapsError{Op: "Set", SetName: name, Err: errUnknownSetName}
	}
	return nil
}

//
// Preconnection
//

func (preconn *Preconnection) Listen() (*Listener, error) {
	servType, err := preconn.getServiceType()
	if err != nil {
		return nil, &tapsError{Op: "Listen", Err: err}
	}
	var lis *Listener
	switch servType {
	case SERV_TCP:
		lis, err = preconn.tpcListen()
	case SERV_QUIC:
		lis, err = preconn.quicListen()
	case SERV_SCION:
		lis, err = preconn.scionListen()
	}
	return lis, err
}

func (preconn *Preconnection) Initiate() (*Connection, error) {
	servType, err := preconn.getServiceType()
	if err != nil {
		return nil, &tapsError{Op: "Initiate", Err: err}
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
			return &tapsError{Op: "Stop", Err: err}
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
	return &Connection{conn.nconn, conn.qconn, conn.sconn, conn.preconn, conn.active, nil, conn.saddr}
}

func (conn *Connection) Receive() (*Message, error) {
	servType, err := conn.preconn.getServiceType()
	if err != nil {
		return nil, &tapsError{Op: "Receive", Err: err}
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
		return &tapsError{Op: "Send", Err: err}
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
			return &tapsError{Op: "Close", Err: err}
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
	return 0, errors.New("no service type")
}
