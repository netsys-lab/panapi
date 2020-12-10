package taps

import (
	"fmt"
	"net"

	quic "github.com/lucas-clemente/quic-go"
)

//

func NewLocalEndpoint() *LocalEndpoint {
	return &LocalEndpoint{}
}

func NewRemoteEndpoint() *RemoteEndpoint {
	return &RemoteEndpoint{}
}

func NewTransportProperties() *TransportProperties {
	return &TransportProperties{}
}
func NewSecurityParameters() *SecurityParameters {
	return &SecurityParameters{}
}

func NewPreconnection(endPo interface{}, transProp *TransportProperties, secParam *SecurityParameters) *Preconnection {
	var ret *Preconnection
	switch endPo.(type) {
	case *LocalEndpoint:
		ret = &Preconnection{endPo.(*LocalEndpoint), nil, transProp, secParam}
	case *RemoteEndpoint:
		ret = &Preconnection{nil, endPo.(*RemoteEndpoint), transProp, secParam}
	default:
		fmt.Println("Error no endpoint type.")
		ret = &Preconnection{}
	}
	return ret
}

func NewListener(lis interface{}, preconn *Preconnection) *Listener {
	var ret *Listener
	state := (lis != nil)
	connChan := make(chan Connection)
	switch preconn.getServiceType() {
	case SERV_TCP:
		ret = &Listener{lis.(net.Listener), nil, preconn, connChan, state}
	case SERV_QUIC:
		ret = &Listener{nil, lis.(quic.Listener), preconn, connChan, state}
	default:
		ret = &Listener{}
	}
	return ret
}

func NewConnection(conn interface{}, preconn *Preconnection) *Connection {
	state := (conn != nil)
	var ret *Connection
	switch preconn.getServiceType() {
	case SERV_TCP:
		ret = &Connection{conn.(net.Conn), nil, preconn, state}
	case SERV_QUIC:
		ret = &Connection{nil, conn.(quic.Session), preconn, state}
	default:
		ret = &Connection{}
	}
	return ret
}

func NewMessage(msg string, ctx string) *Message {
	return &Message{msg, ctx}
}
