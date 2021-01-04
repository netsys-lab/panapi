package taps

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math/big"

	quic "github.com/lucas-clemente/quic-go"
)

//

func (preconn *Preconnection) generateTLSConfig() *tls.Config {
	key := preconn.secParam.privateKey
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"taps-quic-test"},
	}
}

func (preconn *Preconnection) quicListen() (*Listener, error) {
	var lis *Listener
	qlis, err := quic.ListenAddr(preconn.locEnd.ipv4address+":"+preconn.locEnd.port, preconn.generateTLSConfig(), nil)
	if err != nil {
		return nil, &tapsError{Op: "quicListen", Err: err}
	}
	lis, err = NewListener(qlis, preconn)
	go func() {
		sess, err := lis.qlis.Accept(context.Background())
		if err != nil {
			lis.ConnRec <- Connection{Err: err}
			return
		}
		conn, err := NewConnection(sess, preconn)
		if err != nil {
			lis.ConnRec <- Connection{Err: err}
			return
		}
		lis.ConnRec <- *conn
	}()
	return lis, err
}

func (preconn *Preconnection) quicInitiate() (*Connection, error) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"taps-quic-test"},
	}
	session, err := quic.DialAddr(preconn.remEnd.ipv4address+":"+preconn.remEnd.port, tlsConf, nil)
	if err != nil {
		return nil, &tapsError{Op: "quicInitiate", Err: err}
	}
	return NewConnection(session, preconn)
}

func (lis *Listener) quicStop() error {
	// closes current connections
	// lis.qlis.Close()
	return nil
}

func (conn *Connection) quicReceive() (*Message, error) {
	bufSize := 1024
	buf := make([]byte, bufSize)
	if conn.isOpen() {
		stream, err := conn.qconn.AcceptUniStream(context.Background())
		if err != nil {
			err2 := conn.Close()
			if err2 != nil {
				return nil, &tapsError{Op: "quicReceive", Err: err2}
			}
			return nil, &tapsError{Op: "quicReceive", Err: err}
		}
		_, err = stream.Read(buf)
		if err != nil && conn.isOpen() {
			err2 := conn.Close()
			if err2 != nil {
				return nil, &tapsError{Op: "quicReceive", Err: err2}
			}
			return nil, &tapsError{Op: "quicReceive", Err: err}
		}
		return &Message{string(buf), "context"}, nil
	}
	return nil, &tapsError{Op: "quicReceive", Err: errReadOnClosedConnection}
}

func (conn *Connection) quicSend(msg *Message) error {
	if conn.isOpen() {
		stream, err := conn.qconn.OpenUniStreamSync(context.Background())
		if err != nil {
			err2 := conn.Close()
			if err2 != nil {
				return &tapsError{Op: "quicSend", Err: err2}
			}
			return &tapsError{Op: "quicSend", Err: err}
		}
		stream.Write([]byte(msg.Data))
		return nil
	}
	return &tapsError{Op: "quicSend", Err: errWriteOnClosedConnection}
}

func (conn *Connection) quicClose() error {
	return conn.qconn.CloseWithError(1, "closed.")
}
