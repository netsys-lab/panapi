package main

import "fmt"

type Tuple struct {
}

type Endpoint struct {
	interfaceName string
	serviceType   string
}

type LocalEndpoint struct {
	Endpoint
}

type RemoteEndpoint struct {
	Endpoint
}

type TransportProperties struct {
}

type SecurityParameters struct {
}

type Preconnection struct {
	locEnd    LocalEndpoint
	transProp TransportProperties
	secParam  SecurityParameters
}

type Listener struct {
	port string
}

type Connection struct {
	ip   string
	port string
}

type Message struct {
	data    string
	context string
}

func (endPo Endpoint) WithInterface(interfaceName string) {
	endPo.interfaceName = interfaceName
}

func (endPo Endpoint) WithService(serviceType string) {
	endPo.serviceType = serviceType
}

func (tranProp TransportProperties) Require(method string) {
	// process data
}

func (secParam SecurityParameters) Set(id string, list ...int) {
	// process data
}

func (preconn Preconnection) Listen() Listener {
	// start thread that
	return Listener{}
}

func (preconn Preconnection) Initiate() Connection {
	// start thread that
	return Connection{}
}

func (lis Listener) ConnectionReceived() Connection {
	// return conn
	return Connection{}
}

func (lis Listener) Stop() {
	// start thread that
}

func (conn Connection) Clone() Connection {
	// clone conn and ret
	return Connection{}
}

func (conn Connection) Receive() Message {
	// wait for msg
	return Message{"data", "context"}
}

func (conn Connection) Send(msg Message) {
	// send(msg)
}

func (conn Connection) Close() {
	// close conn
}

func main() {
	locEnd := LocalEndpoint{}
	locEnd.WithInterface("any")
	locEnd.WithService("https")

	tranProp := TransportProperties{}
	tranProp.Require("preserve-msg-boundaries")

	secParam := SecurityParameters{}
	secParam.Set("identity", 100)
	secParam.Set("keypair", 1234, 12345678)

	preconn := Preconnection{locEnd, tranProp, secParam}

	listener := preconn.Listen()

	conn := listener.ConnectionReceived()

	reqMsg := conn.Receive()
	fmt.Println(reqMsg)
	respMsg := Message{"testResp", "testContext"}
	conn.Send(respMsg)
	conn.Close()
	listener.Stop()
}
