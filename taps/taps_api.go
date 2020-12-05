package taps

import (
	"fmt"
	"io"
	"net"
	"reflect"
)

//

//
// structs
//

//

type Endpoint struct {
	interfaceName string
	serviceType   string
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

type SecurityParameters struct{}

type Preconnection struct {
	locEnd    LocalEndpoint
	remEnd    RemoteEndpoint
	transProp TransportProperties
	secParam  SecurityParameters
}

type Listener struct {
	nlis    net.Listener
	ConnRec chan Connection
	active  bool
}

type Connection struct {
	nconn  net.Conn
	active bool
}

type Message struct {
	Data    string
	Context string
}

//

//
// public constructors
//

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

func NewPreconnection(endPo interface{}, transProp TransportProperties, secParam SecurityParameters) *Preconnection {
	switch reflect.TypeOf(endPo).Name() {
	case "LocalEndpoint":
		return &Preconnection{endPo.(LocalEndpoint), *NewRemoteEndpoint(), transProp, secParam}
	case "RemoteEndpoint":
		return &Preconnection{*NewLocalEndpoint(), endPo.(RemoteEndpoint), transProp, secParam}
	default:
		return &Preconnection{}
	}
}

func NewListener(nlis net.Listener, ConnRec chan Connection) *Listener {
	if nlis == nil {
		return &Listener{nlis, ConnRec, false}
	} else {
		return &Listener{nlis, ConnRec, true}
	}
}

func NewConnection(nconn net.Conn) *Connection {
	if nconn == nil {
		return &Connection{nconn, false}
	} else {
		return &Connection{nconn, true}
	}
}

func NewMessage(msg string, ctx string) *Message {
	return &Message{msg, ctx}
}

//

//
// methods
//

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
	endPo.serviceType = serviceType
}

func (remEndPo *RemoteEndpoint) WithHostname(hostName string) {
	remEndPo.hostName = hostName
}

func (tranProp *TransportProperties) Require(method string) {
	// process data
}

func (secParam *SecurityParameters) Set(id string, list ...int) {
	// process data
}

func (preconn *Preconnection) Listen() Listener {
	nlis, err := net.Listen(preconn.locEnd.serviceType, preconn.locEnd.ipv4address+":"+preconn.locEnd.port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
	}
	ConnRec := make(chan Connection)
	lis := *NewListener(nlis, ConnRec)
	go func() {
		conn, err := lis.nlis.Accept()
		if err != nil && lis.isOpen() {
			fmt.Println("Error accepting:", err.Error())
		}
		lis.ConnRec <- *NewConnection(conn)
	}()
	return lis
}

func (preconn *Preconnection) Initiate() Connection {
	conn, err := net.Dial(preconn.remEnd.serviceType, preconn.remEnd.hostName+":"+preconn.remEnd.port)
	if err != nil {
		fmt.Println(err)
	}
	return *NewConnection(conn)
}

func (lis *Listener) Stop() {
	if lis.isOpen() {
		lis.active = false
		lis.ConnRec = nil
		lis.nlis.Close()
	}
}

func (conn *Connection) Clone() Connection {
	return *NewConnection(conn.nconn)
}

func (conn *Connection) Receive() Message {
	bufSize := 1024
	buf := make([]byte, bufSize)
	if conn.isOpen() {
		n, err := conn.nconn.Read(buf)
		if err != nil && err != io.EOF && conn.isOpen() {
			fmt.Println("Error reading:", err.Error())
			conn.Close()
		}
		if n > bufSize {
			fmt.Println("Read buffer overflow:", err.Error())
		}
	}
	return Message{string(buf), "context"}
}

func (conn *Connection) Send(msg Message) {
	if conn.isOpen() {
		_, err := conn.nconn.Write([]byte(msg.Data))
		if err != nil {
			fmt.Println("Error sending:", err.Error())
			conn.Close()
		}
	}
}

func (conn *Connection) Close() {
	if conn.isOpen() {
		conn.active = false
		conn.nconn.Close()
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
