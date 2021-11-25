package rpc

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/rpc"
	"os"
	"time"

	"github.com/lucas-clemente/quic-go/logging"
)

type TracerClient struct {
	client *rpc.Client
	f      *os.File
	l      *log.Logger
}

func NewTracerClient(conn io.ReadWriteCloser) logging.Tracer {
	client := rpc.NewClient(conn)
	log.Printf("RPC connection etablished")
	fname := fmt.Sprintf("/tmp/%s-quic-client.log", time.Now().Format("2006-01-02-15-04"))
	log.Println("quic tracer file opened as", fname)
	f, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	go func(f *os.File) {
		for f.Sync() == nil {
			time.Sleep(time.Second)
		}
	}(f)

	return &TracerClient{client, f, log.New(f, "tracer", log.Lshortfile|log.Ltime)}
}

func (c TracerClient) TracerForConnection(ctx context.Context, p logging.Perspective, odcid logging.ConnectionID) logging.ConnectionTracer {

	c.l.Printf("TracerForConnection %+v %+v %+v", ctx, p, odcid)
	err := c.client.Call(
		"TracerServer.TracerForConnection",
		&TracerMsg{
			//Context:      ctx
			Perspective:  &p,
			ConnectionID: &odcid,
		},
		&TracerMsg{},
	)
	if err != nil {
		c.l.Println(err)
	}
	return nil
}

func (c TracerClient) SentPacket(addr net.Addr, hdr *logging.Header, n logging.ByteCount, fs []logging.Frame) {
	c.l.Printf("SentPacket %+v %+v %+v %+v", addr, hdr, n, fs)
	c.client.Call(
		"TracerServer.SentPacket",
		&TracerMsg{
			Addr:      addr,
			Header:    hdr,
			ByteCount: &n,
			Frames:    fs,
		},
		&TracerMsg{},
	)
}

func (c TracerClient) DroppedPacket(addr net.Addr, tp logging.PacketType, n logging.ByteCount, r logging.PacketDropReason) {
	c.l.Printf("DroppedPacket %+v %+v %+v %+v", addr, tp, n, r)
	c.client.Call(
		"TracerServer.DroppedPacket",
		&TracerMsg{
			Addr:       addr,
			PacketType: &tp,
			ByteCount:  &n,
			DropReason: &r,
		},
		&TracerMsg{},
	)
}

type TracerMsg struct {
	//Context      context.Context
	Perspective  *logging.Perspective
	ConnectionID *logging.ConnectionID
	Addr         net.Addr
	Header       *logging.Header
	ByteCount    *logging.ByteCount
	Frames       []logging.Frame
	PacketType   *logging.PacketType
	DropReason   *logging.PacketDropReason
}

type TracerServer struct {
	tracer logging.Tracer
	f      *os.File
	l      *log.Logger
}

func NewTracerServer(tracer logging.Tracer) *TracerServer {
	fname := fmt.Sprintf("/tmp/%s-quic-server.log", time.Now().Format("2006-01-02-15-04"))
	log.Println("quic tracer file opened as", fname)
	f, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	go func(f *os.File) {
		for f.Sync() == nil {
			time.Sleep(time.Second)
		}
	}(f)

	return &TracerServer{tracer, f, log.New(f, "tracer", log.Lshortfile|log.Ltime)}
}

func (s *TracerServer) TracerForConnection(args, resp *TracerMsg) error {
	if args.Perspective != nil && args.ConnectionID != nil {
		ctx := context.TODO()
		s.l.Printf("TracerForConnection %+v %+v %+v", ctx, *args.Perspective, *args.ConnectionID)
		s.tracer.TracerForConnection(ctx, *args.Perspective, *args.ConnectionID)
	} else {
		return ErrDeref
	}
	return nil
}

func (s *TracerServer) SentPacket(args, resp *TracerMsg) error {
	if args.Addr != nil && args.ByteCount != nil {
		s.l.Printf("SentPacket %+v %+v %+v %+v", args.Addr, args.Header, *args.ByteCount, args.Frames)
		s.tracer.SentPacket(args.Addr, args.Header, *args.ByteCount, args.Frames)
	} else {
		return ErrDeref
	}
	return nil
}

func (s *TracerServer) DroppedPacket(args, resp *TracerMsg) error {
	if args.Addr != nil && args.PacketType != nil && args.ByteCount != nil && args.DropReason != nil {
		s.l.Printf("DroppedPacket %+v %+v %+v %+v", args.Addr, *args.PacketType, *args.ByteCount, *args.DropReason)
		s.tracer.DroppedPacket(args.Addr, *args.PacketType, *args.ByteCount, *args.DropReason)
	} else {
		return ErrDeref
	}
	return nil
}
