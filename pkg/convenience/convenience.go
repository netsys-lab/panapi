package convenience

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"net"

	"github.com/netsys-lab/panapi/rpc"
)

func GenerateTLSConfig() tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
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
	return tls.Config{
		Certificates: []tls.Certificate{tlsCert},
	}
}

func NewRPCClient() (*rpc.Client, error) {
	conn, err := net.Dial(rpc.DefaultDaemonAddress.Net, rpc.DefaultDaemonAddress.Name)
	if err != nil {
		return nil, err
	}
	return rpc.NewClient(conn)
}
