package taps

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"

	quic "github.com/lucas-clemente/quic-go"
)

//

func (preconn *Preconnection) quicListen() *Listener {
	var lis *Listener
	qlis, err := quic.ListenAddr(preconn.locEnd.ipv4address+":"+preconn.locEnd.port, generateTLSConfig(), nil)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
	}
	lis = NewListener(qlis, preconn)
	go func() {
		sess, err := lis.qlis.Accept(context.Background())
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
		}
		lis.ConnRec <- *NewConnection(sess, preconn)
	}()
	return lis
}

func (preconn *Preconnection) quicInitiate() *Connection {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"taps-quic-test"},
	}
	session, err := quic.DialAddr(preconn.remEnd.ipv4address+":"+preconn.remEnd.port, tlsConf, nil)
	if err != nil {
		fmt.Println("Error quic init:", err.Error())
	}
	return NewConnection(session, preconn)
}

func (lis *Listener) quicStop() {
	// lis.qlis.Close()
}

func (conn *Connection) quicReceive() *Message {
	bufSize := 1024
	buf := make([]byte, bufSize)
	if conn.isOpen() {
		stream, err := conn.qconn.AcceptUniStream(context.Background())
		if err != nil {
			if conn.isOpen() {
				fmt.Println("Error quic acc str:", err.Error())
			}
			conn.Close()
			return &Message{}
		}
		n, err := stream.Read(buf)
		if err != nil && conn.isOpen() {
			fmt.Println("Error quic rec:", err.Error())
			conn.Close()
			return &Message{}
		}
		if n > bufSize {
			fmt.Println("Read buffer overflow:", err.Error())
		}
	}
	return &Message{string(buf), "context"}
}

func (conn *Connection) quicSend(msg *Message) {
	if conn.isOpen() {
		stream, err := conn.qconn.OpenUniStreamSync(context.Background())
		if err != nil {
			fmt.Println("Error quic send:", err.Error())
		}
		stream.Write([]byte(msg.Data))
	}
}

func (conn *Connection) quicClose() {
	conn.qconn.CloseWithError(0x1234, "closed.")
}

//

func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
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
