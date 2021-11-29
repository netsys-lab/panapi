package rpc

import (
	"fmt"
	"io"
	"log"
	"net/rpc"
	"os"
	"time"

	"github.com/lucas-clemente/quic-go/logging"
)

type IDMsg struct {
	Value *int
}

type IDServer struct {
	id int
}

func (s *IDServer) GetID(arg, resp *IDMsg) error {
	resp.Value = &s.id
	s.id += 1
	return nil
}

func NewServer(selector ServerSelector, tracer logging.Tracer) (*rpc.Server, error) {
	/*err := rpc.Register(IDServer{42})
	if err != nil {
		panic(err)
		return nil, err
	}*/
	err := rpc.Register(NewSelectorServer(selector))
	if err != nil {
		return nil, err
	}
	err = rpc.Register(NewTracerServer(tracer))
	if err != nil {
		return nil, err
	}
	err = rpc.Register(NewConnectionTracerServer(tracer, log.Default()))
	if err != nil {
		return nil, err
	}
	return rpc.DefaultServer, nil
}

type Client struct {
	client *rpc.Client
	l      *log.Logger
	id     int
}

func NewClient(conn io.ReadWriteCloser) (*Client, error) {
	client := rpc.NewClient(conn)
	fname := fmt.Sprintf("/tmp/%s-quic-rpc-client.log", time.Now().Format("2006-01-02-15-04"))
	log.Println("quic rpc client file opened as", fname)
	f, err := os.Create(fname)
	if err != nil {
		return nil, err
	}
	//FIXME
	go func(f *os.File) {
		for f.Sync() == nil {
			time.Sleep(time.Second)
		}
	}(f)
	//f := os.Stderr

	n := 42
	var id = &IDMsg{Value: &n}
	/*err = client.Call("IDServer.GetID", &IDMsg{}, &id)
	if err != nil {
		return nil, err
	} else {*/

	log.Printf("RPC connection etablished with ID %d", *id.Value)
	//}

	return &Client{client, log.New(f, "rpc-client", log.Lshortfile|log.Ltime), *id.Value}, nil
}

func (c *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	return c.client.Call(serviceMethod, args, reply)
}

func (c *Client) Close() error {
	return c.client.Close()
}
