package taps

import (
	"net"
	"reflect"

	quic "github.com/lucas-clemente/quic-go"
	"github.com/scionproto/scion/go/lib/snet"
)

//

func NewLocalEndpoint() *LocalEndpoint {
	return &LocalEndpoint{
		Endpoint: Endpoint{
			serviceType: SERV_NONE}}
}

func NewRemoteEndpoint() *RemoteEndpoint {
	return &RemoteEndpoint{
		Endpoint: Endpoint{
			serviceType: SERV_NONE}}
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
		return &Preconnection{
				locEnd:    endPo.(*LocalEndpoint),
				remEnd:    nil,
				transProp: transProp,
				secParam:  secParam},
			nil
	case *RemoteEndpoint:
		return &Preconnection{
				locEnd:    nil,
				remEnd:    endPo.(*RemoteEndpoint),
				transProp: transProp,
				secParam:  secParam},
			nil
	default:
		return nil, &tapsError{
			Op:   "NewPreconnection",
			Desc: reflect.TypeOf(endPo).String(),
			Err:  errInvalidEndpointType}
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
		ret = &Listener{
			nlis:    lis.(net.Listener),
			qlis:    nil,
			preconn: preconn,
			ConnRec: connChan,
			active:  state}
	case SERV_QUIC:
		ret = &Listener{
			nlis:    nil,
			qlis:    lis.(quic.Listener),
			preconn: preconn,
			ConnRec: connChan,
			active:  state}
	case SERV_SCION:
		ret = &Listener{
			nlis:    nil,
			qlis:    nil,
			preconn: preconn,
			ConnRec: connChan,
			active:  state}
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
		ret = &Connection{
			nconn:   conn.(net.Conn),
			qconn:   nil,
			sconn:   nil,
			preconn: preconn,
			active:  state,
			Err:     nil,
			saddr:   nil}
	case SERV_QUIC:
		ret = &Connection{
			nconn:   nil,
			qconn:   conn.(quic.Session),
			sconn:   nil,
			preconn: preconn,
			active:  state,
			Err:     nil,
			saddr:   nil}
	case SERV_SCION:
		ret = &Connection{
			nconn:   nil,
			qconn:   nil,
			sconn:   conn.(*snet.Conn),
			preconn: preconn,
			active:  state,
			Err:     nil,
			saddr:   nil}
	}
	return ret, nil
}

func NewMessage(msg string, ctx string) *Message {
	return &Message{
		Data:    msg,
		Context: ctx}
}
