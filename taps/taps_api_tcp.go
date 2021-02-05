package taps

import (
	"io"
	"net"
)

//

func (preconn *Preconnection) tpcListen() (*Listener, error) {
	nlis, err := net.Listen("tcp", "["+preconn.locEnd.address+"]:"+preconn.locEnd.port)
	if err != nil {
		return nil, &tapsError{
			Op:   "tpcListen",
			Endp: preconn.locEnd,
			Err:  err}
	}
	lis, err := NewListener(nlis, preconn)
	if err != nil {
		return nil, &tapsError{
			Op:   "tpcListen",
			Endp: preconn,
			Err:  err}
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
	nconn, err := net.Dial("tcp", "["+preconn.remEnd.address+"]:"+preconn.remEnd.port)
	if err != nil {
		return nil, &tapsError{
			Op:   "tcpInitiate",
			Endp: preconn,
			Err:  err}
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
			conn.Close()
			return nil, &tapsError{
				Op:   "tcpReceive",
				Endp: conn.preconn,
				Err:  err}
		}
		return &Message{string(buf), "context"}, nil
	}
	return nil, &tapsError{Op: "tcpReceive", Err: errReadOnClosedConnection}
}

func (conn *Connection) tcpSend(msg *Message) error {
	if conn.isOpen() {
		_, err := conn.nconn.Write([]byte(msg.Data))
		if err != nil {
			conn.Close()
			return &tapsError{
				Op:   "tcpSend",
				Endp: conn.preconn,
				Err:  err}
		}
		return nil
	}
	return &tapsError{
		Op:   "tcpSend",
		Endp: conn.preconn,
		Err:  errWriteOnClosedConnection}
}

func (conn *Connection) tcpClose() error {
	return conn.nconn.Close()
}
