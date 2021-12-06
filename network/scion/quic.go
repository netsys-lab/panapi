package scion

import (
	"context"
	"crypto/tls"

	"log"
	"net"

	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/logging"
	"inet.af/netaddr"

	"github.com/netsec-ethz/scion-apps/pkg/pan"
	"github.com/netsec-ethz/scion-apps/pkg/quicutil"

	"github.com/netsys-lab/panapi/network"
	"github.com/netsys-lab/panapi/rpc"
)

type QUICDialer struct {
	raddr    pan.UDPAddr
	selector pan.Selector
	tp       *network.TransportProperties
	client   *rpc.Client
}

func (d *QUICDialer) getSelector() (selector pan.Selector) {
	return
}

func NewQUICDialer(address string, tp *network.TransportProperties) (*QUICDialer, error) {
	var (
		selector pan.Selector
		client   *rpc.Client
	)

	if tp != nil {
		conn, err := net.DialUnix("unix", nil, rpc.DefaultDaemonAddress)
		if err != nil {
			log.Printf("Could not connect to PANAPI Deamon: %s", err)
			log.Println("Using default selector")
			selector = &pan.DefaultSelector{}
		} else {
			log.Println("using daemon selector")
			client, err = rpc.NewClient(conn)
			if err != nil {
				return nil, err
			}
			selector = rpc.NewSelectorClient(client)

		}
	} else {
		log.Println("no transport properties given")
	}

	addr, err := pan.ResolveUDPAddr(address)
	return &QUICDialer{addr, selector, tp, client}, err
}

func (d *QUICDialer) Dial() (network.Connection, error) {
	tlsConf := &tls.Config{
		//Certificates: quicutil.MustGenerateSelfSignedCert(),
		InsecureSkipVerify: true,
		NextProtos:         []string{"panapi-quic-test"},
	}
	var tracer logging.Tracer
	if d.client == nil {
		log.Printf("Could not connect to PANAPI Deamon")
		log.Println("Not using tracer")
	} else {
		tracer = rpc.NewTracerClient(d.client)
	}

	conn, err := pan.DialQUIC(context.Background(), netaddr.IPPort{}, d.raddr, nil, d.selector, "", tlsConf, &quic.Config{Tracer: tracer})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	stream, err := conn.OpenStream() //Sync(context.Background())
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return network.NewQUIC(conn, stream, conn.LocalAddr(), conn.RemoteAddr()), nil
}

type QUICListener struct {
	listener quic.Listener
	client   *rpc.Client
}

func NewQUICListener(address string, tp *network.TransportProperties) (*QUICListener, error) {
	addr, err := pan.ResolveUDPAddr(address)
	if err != nil {
		return nil, err
	}
	tlsConf := &tls.Config{
		Certificates: quicutil.MustGenerateSelfSignedCert(),
		//InsecureSkipVerify: true,
		NextProtos: []string{"panapi-quic-test"},
	}
	//tlsConf, err := generateTLSConfig()
	var (
		tracer logging.Tracer
		client *rpc.Client
	)
	conn, err := net.DialUnix("unix", nil, rpc.DefaultDaemonAddress)
	if err != nil {
		log.Printf("Could not connect to PANAPI Deamon: %s", err)
		log.Println("Not using tracer")
	} else {
		client, err = rpc.NewClient(conn)
		if err != nil {
			return nil, err
		} else {
			tracer = rpc.NewTracerClient(client)
		}
	}

	listener, err := pan.ListenQUIC(context.Background(), netaddr.IPPortFrom(addr.IP, addr.Port), nil, tlsConf, &quic.Config{
		Tracer: tracer,
	})
	if err != nil {
		return nil, err
	}
	return &QUICListener{listener, client}, nil
}

func (l *QUICListener) Listen() (network.Connection, error) {
	conn, err := l.listener.Accept(context.Background())
	if err != nil {
		return &network.QUIC{}, err
	}
	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		return &network.QUIC{}, err
	}
	log.Println("accepted stream")
	return network.NewQUIC(conn, stream, conn.LocalAddr(), conn.RemoteAddr()), nil
}

func (l *QUICListener) Stop() error {
	log.Println("Stop called")
	return nil
}
