package taps

import "fmt"

//
// Endpoint
//

func (endPo *Endpoint) WithInterface(interfaceName string) {
	endPo.interfaceName = interfaceName
	// i, _ := net.Interfaces()
	// _, e := net.InterfaceByName("")
	// fmt.Print(i)
}

func (endPo *Endpoint) WithPort(port string) {
	endPo.port = port
}

func (endPo *Endpoint) WithIPv4Address(addr string) {
	endPo.ipv4address = addr
}

func (endPo *Endpoint) WithService(serviceType string) {
	switch serviceType {
	case "tcp":
		endPo.serviceType = SERV_TCP
	case "quic":
		endPo.serviceType = SERV_QUIC
	}
}

func (remEndPo *RemoteEndpoint) WithHostname(hostName string) {
	remEndPo.hostName = hostName
}

//
// TransportProperties
//

func (tranProp *TransportProperties) Require(method string) {
	// process data
}

//
// SecurityParameters
//

func (secParam *SecurityParameters) Set(id string, list ...int) {
	// process data
}

//
// Preconnection
//

func (preconn *Preconnection) Listen() *Listener {
	var lis *Listener
	switch preconn.getServiceType() {
	case SERV_TCP:
		lis = preconn.tpcListen()
	case SERV_QUIC:
		lis = preconn.quicListen()
	}
	return lis
}

func (preconn *Preconnection) Initiate() *Connection {
	var ret *Connection
	switch preconn.getServiceType() {
	case SERV_TCP:
		ret = preconn.tcpInitiate()
	case SERV_QUIC:
		ret = preconn.quicInitiate()
	}
	return ret
}

//
// Listener
//

func (lis *Listener) Stop() {
	if lis.isOpen() {
		lis.active = false
		lis.ConnRec = nil
		switch lis.preconn.getServiceType() {
		case SERV_TCP:
			lis.tcpStop()
		case SERV_QUIC:
			lis.quicStop()
		}
	}
}

//
// Connection
//

func (conn *Connection) Clone() *Connection {
	return &Connection{conn.nconn, conn.qconn, conn.preconn, conn.active}
}

func (conn *Connection) Receive() *Message {
	var ret *Message
	switch conn.preconn.getServiceType() {
	case SERV_TCP:
		ret = conn.tcpReceive()
	case SERV_QUIC:
		ret = conn.quicReceive()
	}
	return ret
}

func (conn *Connection) Send(msg *Message) {
	switch conn.preconn.getServiceType() {
	case SERV_TCP:
		conn.tcpSend(msg)
	case SERV_QUIC:
		conn.quicSend(msg)
	}
}

func (conn *Connection) Close() {
	if conn.isOpen() {
		conn.active = false
		switch conn.preconn.getServiceType() {
		case SERV_TCP:
			conn.tcpClose()
		case SERV_QUIC:
			conn.quicClose()
		}
	}
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

func (preconn *Preconnection) getServiceType() int {
	if preconn.locEnd == nil {
		if preconn.remEnd != nil {
			return preconn.remEnd.serviceType
		}
	}
	if preconn.remEnd == nil {
		if preconn.locEnd != nil {
			return preconn.locEnd.serviceType
		}
	}
	fmt.Println("Error no Service Type")
	return -1
}
