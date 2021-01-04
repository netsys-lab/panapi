package taps

import (
	"net"

	quic "github.com/lucas-clemente/quic-go"
)

//

func NewLocalEndpoint() *LocalEndpoint {
	return &LocalEndpoint{Endpoint: Endpoint{serviceType: SERV_NONE}}
}

func NewRemoteEndpoint() *RemoteEndpoint {
	return &RemoteEndpoint{Endpoint: Endpoint{serviceType: SERV_NONE}}
}

func NewTransportProperties() *TransportProperties {
	return &TransportProperties{}
}
func NewSecurityParameters() *SecurityParameters {
	return &SecurityParameters{}
}

func NewPreconnection(endPo interface{}, transProp *TransportProperties, secParam *SecurityParameters) (*Preconnection, error) {
	switch endPo.(type) {
	case *LocalEndpoint:
		return &Preconnection{endPo.(*LocalEndpoint), nil, transProp, secParam}, nil
	case *RemoteEndpoint:
		return &Preconnection{nil, endPo.(*RemoteEndpoint), transProp, secParam}, nil
	default:
		return nil, &tapsError{Op: "NewPreconnection", Err: errInvalidEndpointType}
	}
}

func NewListener(lis interface{}, preconn *Preconnection) (*Listener, error) {
	servType, err := preconn.getServiceType()
	if err != nil {
		return nil, &tapsError{Op: "NewListener", Err: err}
	}
	ret := &Listener{}
	state := (lis != nil)
	connChan := make(chan Connection)
	switch servType {
	case SERV_TCP:
		ret = &Listener{lis.(net.Listener), nil, preconn, connChan, state}
	case SERV_QUIC:
		ret = &Listener{nil, lis.(quic.Listener), preconn, connChan, state}
	}
	return ret, nil
}

func NewConnection(conn interface{}, preconn *Preconnection) (*Connection, error) {
	servType, err := preconn.getServiceType()
	if err != nil {
		return nil, &tapsError{Op: "NewConnection", Err: err}
	}
	ret := &Connection{}
	state := (conn != nil)
	switch servType {
	case SERV_TCP:
		ret = &Connection{conn.(net.Conn), nil, preconn, state, nil}
	case SERV_QUIC:
		ret = &Connection{nil, conn.(quic.Session), preconn, state, nil}
	}
	return ret, nil
}

func NewMessage(msg string, ctx string) *Message {
	return &Message{msg, ctx}
}
