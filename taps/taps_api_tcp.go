package taps

import (
	"io"
	"net"
)

//

func (preconn *Preconnection) tpcListen() (*Listener, error) {
	nlis, err := net.Listen("tcp", "["+preconn.locEnd.ipv4address+"]:"+preconn.locEnd.port)
	if err != nil {
		return nil, &tapsError{Op: "tpcListen", Err: err}
	}
	lis, err := NewListener(nlis, preconn)
	if err != nil {
		return nil, &tapsError{Op: "tpcListen", Err: err}
	}
	go func() {
		nconn, err := lis.nlis.Accept()
		if err != nil && lis.isOpen() {
			lis.ConnRec <- Connection{Err: err}
			return
		}
		nconn.(*net.TCPConn).SetNoDelay(!preconn.transProp.nagle)
		conn, err := NewConnection(nconn, preconn)
		if err != nil {
			lis.ConnRec <- Connection{Err: err}
			return
		}
		lis.ConnRec <- *conn
	}()
	return lis, nil
}

func (preconn *Preconnection) tcpInitiate() (*Connection, error) {
	nconn, err := net.Dial("tcp", "["+preconn.remEnd.ipv4address+"]:"+preconn.remEnd.port)
	if err != nil {
		return nil, &tapsError{Op: "tcpInitiate", Err: err}
	}
	nconn.(*net.TCPConn).SetNoDelay(!preconn.transProp.nagle)
	return NewConnection(nconn, preconn)
}

func (lis *Listener) tcpStop() error {
	return lis.nlis.Close()
}

func (conn *Connection) tcpReceive() (*Message, error) {
	bufSize := 1024
	buf := make([]byte, bufSize)
	if conn.isOpen() {
		_, err := conn.nconn.Read(buf)
		if err != nil && err != io.EOF && conn.isOpen() {
			// err2 :=
			conn.Close()
			// if err2 != nil {
			// 	return nil, &tapsError{Op: "tcpReceive", Err: err2}
			// }
			return nil, &tapsError{Op: "tcpReceive", Err: err}
		}
		return &Message{string(buf), "context"}, nil
	}
	return nil, &tapsError{Op: "tcpReceive", Err: errReadOnClosedConnection}
}

func (conn *Connection) tcpSend(msg *Message) error {
	if conn.isOpen() {
		_, err := conn.nconn.Write([]byte(msg.Data))
		if err != nil {
			err2 := conn.Close()
			if err2 != nil {
				return &tapsError{Op: "tcpSend", Err: err2}
			}
			return &tapsError{Op: "tcpSend", Err: err}
		}
		return nil
	}
	return &tapsError{Op: "tcpSend", Err: errWriteOnClosedConnection}
}

func (conn *Connection) tcpClose() error {
	return conn.nconn.Close()
}
