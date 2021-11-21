package scion

import (
	//"crypto/rand"
	//"crypto/rsa"
	//"crypto/x509"
	//"encoding/pem"
	"errors"
	"fmt"
	//"math/big"

	"github.com/netsys-lab/panapi/network"
)

type scion struct {
	tp *network.TransportProperties
}

func (scion *scion) NewListener(e *network.Endpoint) (network.Listener, error) {
	switch e.Transport {
	case network.TRANSPORT_UDP:
		return NewUDPListener(e.LocalAddress)
	case network.TRANSPORT_QUIC:
		return NewQUICListener(e.LocalAddress, scion.tp)
	default:
		return nil, errors.New(fmt.Sprintf("Transport %s not implemented for SCION", e.Transport))
	}
}

func (scion *scion) NewDialer(e *network.Endpoint) (network.Dialer, error) {
	switch e.Transport {
	case network.TRANSPORT_UDP:
		return NewUDPDialer(e.RemoteAddress)
	case network.TRANSPORT_QUIC:
		return NewQUICDialer(e.RemoteAddress, scion.tp)
	default:
		return nil, errors.New(fmt.Sprintf("Transport %s not implemented for SCION", e.Transport))
	}
}

func Network(tp *network.TransportProperties) network.Network {
	return &scion{tp}
}

/*func generateTLSConfig() (*tls.Config, error) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return nil, err
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"taps-quic-test"},
	}, nil
        }*/
