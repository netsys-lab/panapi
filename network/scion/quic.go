package scion

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	//"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/logging"
	"github.com/lucas-clemente/quic-go/qlog"

	"github.com/netsec-ethz/scion-apps/pkg/pan"

	"github.com/netsec-ethz/scion-apps/pkg/quicutil"
	"github.com/netsys-lab/panapi/network"
	"github.com/netsys-lab/panapi/rpc"
)

type QUICDialer struct {
	raddr    pan.UDPAddr
	selector pan.Selector
	tp       *network.TransportProperties
}

func getSelector(tp *network.TransportProperties) (selector pan.Selector, err error) {
	if tp != nil {
		/*if script, ok := tp.Properties["lua-script"]; ok {
			selector, err = NewLuaSelector(script)
			if err != nil {
				return
			}
		} else {*/
		//log.Println("no selector script found in transport properties")
		log.Println("using daemon selector")
		var conn *net.UnixConn
		conn, err = net.DialUnix("unix", nil, rpc.DefaultDaemonAddress)
		if err != nil {
			return
		}
		selector = rpc.NewSelectorClient(conn)
		return
		//		}
	} else {
		log.Println("no transport properties given")
	}
	return
}

func NewQUICDialer(address string, tp *network.TransportProperties) (*QUICDialer, error) {
	var (
		selector pan.Selector
		err      error
		addr     pan.UDPAddr
	)
	selector, err = getSelector(tp)
	if err != nil {
		return nil, err
	}
	addr, err = pan.ResolveUDPAddr(address)
	return &QUICDialer{addr, selector, tp}, err
}

func (d *QUICDialer) Dial() (network.Connection, error) {
	tlsConf := &tls.Config{
		//Certificates: quicutil.MustGenerateSelfSignedCert(),
		InsecureSkipVerify: true,
		NextProtos:         []string{"panapi-quic-test"},
	}
	//log.Printf("%+v", d.selector)
	conn, err := pan.DialQUIC(context.Background(), nil, d.raddr, nil, d.selector, "", tlsConf, &quic.Config{
		Tracer: qlog.NewTracer(func(p logging.Perspective, connectionID []byte) io.WriteCloser {
			fname := fmt.Sprintf("/tmp/%s-%x-quic-dialer-%d.log", time.Now().Format("2006-01-02-15-04"), connectionID, p)
			log.Println("quic tracer file opened as", fname)
			f, err := os.Create(fname)
			if err != nil {
				panic(err)
			}
			return f
		}),
	})

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
}

func NewQUICListener(address string, tp *network.TransportProperties) (*QUICListener, error) {
	/*var (
		selector pan.Selector
		err      error
	)
	selector, err = getSelector(tp)
	if err != nil {
		return nil, err
	}*/

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
	if err != nil {
		return nil, err
	}
	listener, err := pan.ListenQUIC(context.Background(), &net.UDPAddr{Port: addr.Port}, nil, tlsConf, &quic.Config{
		Tracer: qlog.NewTracer(func(p logging.Perspective, connectionID []byte) io.WriteCloser {
			fname := fmt.Sprintf("/tmp/%s-%x-quic-listener-%d.log", time.Now().Format("2006-01-02-15-04"), connectionID, p)
			log.Println("quic tracer file opened as", fname)
			f, err := os.Create(fname)
			if err != nil {
				panic(err)
			}
			return f
		}),
	})
	if err != nil {
		return nil, err
	}
	return &QUICListener{listener}, nil
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
	return nil
}
