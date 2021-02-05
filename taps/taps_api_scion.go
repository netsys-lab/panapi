package taps

import (
	"net"
	"strconv"

	"github.com/netsec-ethz/scion-apps/pkg/appnet"
)

//

func (preconn *Preconnection) scionListen() (*Listener, error) {
	lis, err := NewListener(nil, preconn)
	port, err := strconv.ParseUint(preconn.locEnd.port, 10, 16)
	sconn, err := appnet.ListenPort(uint16(port))
	if err != nil && lis.isOpen() {
		return nil, &tapsError{
			Op:   "scionListen",
			Endp: preconn,
			Err:  err}
	}
	conn, err := NewConnection(sconn, preconn)
	if err != nil {
		return nil, &tapsError{
			Op:   "scionListen",
			Endp: preconn,
			Err:  err}
	}
	go func() {
		lis.ConnRec <- *conn
	}()
	return lis, err
}

func (preconn *Preconnection) scionInitiate() (*Connection, error) {
	sconn, err := appnet.Dial(preconn.remEnd.address)
	if err != nil {
		return nil, &tapsError{
			Op:   "scionInitiate",
			Endp: preconn,
			Err:  err}
	}
	return NewConnection(sconn, preconn)
}

func (lis *Listener) scionStop() error {
	return nil
}

func (conn *Connection) scionReceive() (*Message, error) {
	buffer := make([]byte, 1024)
	n, from, err := conn.sconn.ReadFrom(buffer)
	if err != nil {
		return nil, &tapsError{
			Op:   "scionReceive",
			Endp: conn.preconn,
			Err:  err}
	}
	if conn.saddr == nil {
		conn.saddr = from
	}
	data := buffer[:n]
	return &Message{string(data), "context"}, nil
}

func (conn *Connection) scionSend(msg *Message) error {
	var addr net.Addr
	var err error
	if conn.preconn.locEnd == nil {
		if conn.preconn.remEnd != nil {
			addr = conn.sconn.RemoteAddr()
		}
	}
	if conn.preconn.remEnd == nil {
		if conn.preconn.locEnd != nil {
			if conn.saddr != nil {
				addr = conn.saddr
			} else {
				return &tapsError{
					Op:   "scionSend",
					Endp: conn.preconn,
					Err:  errNoClientAddr}
			}
		}
	}
	_, err = conn.sconn.WriteTo([]byte(msg.Data), addr)
	if err != nil {
		conn.Close()
		return &tapsError{
			Op:   "scionSend",
			Endp: conn.preconn,
			Err:  err}
	}
	return nil
}

func (conn *Connection) scionClose() error {
	err := conn.sconn.Close()
	if err != nil {
		return &tapsError{
			Op:   "scionClose",
			Endp: conn.preconn,
			Err:  err}
	}
	return nil
}
