package taps

import (
	"fmt"
	"io"
	"net"
)

//

func (preconn *Preconnection) tpcListen() *Listener {
	var lis *Listener
	nlis, err := net.Listen("tcp", preconn.locEnd.ipv4address+":"+preconn.locEnd.port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
	}
	lis = NewListener(nlis, preconn)
	go func() {
		conn, err := lis.nlis.Accept()
		if err != nil && lis.isOpen() {
			fmt.Println("Error accepting:", err.Error())
		}
		lis.ConnRec <- *NewConnection(conn, preconn)
	}()
	return lis
}

func (preconn *Preconnection) tcpInitiate() *Connection {
	conn, err := net.Dial("tcp", preconn.remEnd.hostName+":"+preconn.remEnd.port)
	if err != nil {
		fmt.Println(err.Error())
	}
	return NewConnection(conn, preconn)
}

func (lis *Listener) tcpStop() {
	lis.nlis.Close()
}

func (conn *Connection) tcpReceive() *Message {
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
	return &Message{string(buf), "context"}
}

func (conn *Connection) tcpSend(msg *Message) {
	if conn.isOpen() {
		_, err := conn.nconn.Write([]byte(msg.Data))
		if err != nil {
			fmt.Println("Error sending:", err.Error())
			conn.Close()
		}
	}
}

func (conn *Connection) tcpClose() {
	conn.nconn.Close()
}
