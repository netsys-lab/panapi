package rpc

import (
	//"fmt"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/rpc"

	"github.com/netsec-ethz/scion-apps/pkg/pan"
)

var (
	DefaultDaemonAddress = &net.UnixAddr{
		Name: "/tmp/panapid.sock",
		Net:  "unix",
	}
	ErrDeref = errors.New("Can not dereference Nil value")
)

type ServerSelector interface {
	Initialize(pan.UDPAddr, pan.UDPAddr, []*pan.Path)
	Path(pan.UDPAddr) *pan.Path
	PathDown(pan.UDPAddr, pan.PathFingerprint, pan.PathInterface)
	Refresh(pan.UDPAddr, []*pan.Path)
	Close(pan.UDPAddr) error
}

type serverSelector struct {
	fn        func(pan.UDPAddr) pan.Selector
	selectors map[string]pan.Selector
}

func NewServerSelectorFunc(fn func(pan.UDPAddr) pan.Selector) ServerSelector {
	return serverSelector{fn, map[string]pan.Selector{}}
}

func (s *serverSelector) getSelector(raddr pan.UDPAddr) pan.Selector {
	selector, ok := s.selectors[raddr.String()]
	if !ok {
		selector = s.fn(raddr)
		s.selectors[raddr.String()] = selector
	}
	return selector
}

func (s serverSelector) Path(raddr pan.UDPAddr) *pan.Path {
	return s.getSelector(raddr).Path()
}

func (s serverSelector) Initialize(local, remote pan.UDPAddr, paths []*pan.Path) {
	s.getSelector(remote).Initialize(local, remote, paths)
}

func (s serverSelector) PathDown(raddr pan.UDPAddr, fp pan.PathFingerprint, pi pan.PathInterface) {
	s.getSelector(raddr).PathDown(fp, pi)
}

func (s serverSelector) Refresh(remote pan.UDPAddr, paths []*pan.Path) {
	s.getSelector(remote).Refresh(paths)
}

func (s serverSelector) Close(raddr pan.UDPAddr) error {
	err := s.getSelector(raddr).Close()
	delete(s.selectors, raddr.String())
	return err
}

type SelectorMsg struct {
	Local         *pan.UDPAddr
	Remote        *pan.UDPAddr
	Fingerprint   *pan.PathFingerprint
	PathInterface *pan.PathInterface
	Paths         []*Path
}

type SelectorServer struct {
	selector ServerSelector
}

/*func NewSelectorServer(selector ServerSelector) (*rpc.Server, error) {
	err := rpc.Register(&SelectorServer{selector})
	if err != nil {
		return nil, err
	}
	return rpc.DefaultServer, nil
}

func (s *SelectorServer) RegisterAt(r *rpc.Server) error {
	return r.Register(s)
        }*/

func NewSelectorServer(selector ServerSelector) *SelectorServer {
	return &SelectorServer{selector}
}

func (s *SelectorServer) Initialize(args, resp *SelectorMsg) error {
	fmt.Println("Initialize invoked")
	paths := make([]*pan.Path, len(args.Paths))
	for i, p := range args.Paths {
		paths[i] = p.PanPath()
		//log.Printf("%s", paths[i].Source)
	}
	if args.Local == nil || args.Remote == nil {
		return ErrDeref
	}
	s.selector.Initialize(*args.Local, *args.Remote, paths)
	msg := "Initialize done"
	fmt.Println(msg)
	return nil
}

func (s *SelectorServer) Path(args, resp *SelectorMsg) error {
	if args.Remote == nil {
		return ErrDeref
	}
	p := s.selector.Path(*args.Remote)
	if p != nil {
		resp.Fingerprint = &p.Fingerprint
	}
	return nil
}

func (s *SelectorServer) PathDown(args, resp *SelectorMsg) error {
	log.Println("PathDown called")
	if args.Remote == nil || args.Fingerprint == nil || args.PathInterface == nil {
		return ErrDeref
	}
	s.selector.PathDown(*args.Remote, *args.Fingerprint, *args.PathInterface)
	return nil
}

func (s *SelectorServer) Refresh(args, resp *SelectorMsg) error {
	fmt.Println("Refresh invoked")
	paths := make([]*pan.Path, len(args.Paths))
	for i, p := range args.Paths {
		paths[i] = p.PanPath()
		//log.Printf("%s", paths[i].Source)
	}
	if args.Remote == nil {
		return ErrDeref
	}
	s.selector.Refresh(*args.Remote, paths)
	msg := "Refresh done"
	fmt.Println(msg)
	return nil
}

func (s *SelectorServer) Close(args, resp *SelectorMsg) error {
	log.Println("Close called")
	if args.Remote == nil {
		return ErrDeref
	}
	return s.selector.Close(*args.Remote)
}

type SelectorClient struct {
	client *rpc.Client
	paths  map[pan.PathFingerprint]*pan.Path
	local  *pan.UDPAddr
	remote *pan.UDPAddr
}

func NewSelectorClient(conn io.ReadWriteCloser) *SelectorClient {
	client := rpc.NewClient(conn)
	log.Printf("RPC connection etablished")
	return &SelectorClient{client, map[pan.PathFingerprint]*pan.Path{}, nil, nil}
}

func (s *SelectorClient) Initialize(local, remote pan.UDPAddr, paths []*pan.Path) {
	log.Println("Initialize called")
	s.remote = &remote
	s.local = &local
	ps := make([]*Path, len(paths))
	for i, p := range paths {
		s.paths[p.Fingerprint] = p
		ps[i] = NewPathFrom(p)
	}
	err := s.client.Call("SelectorServer.Initialize", &SelectorMsg{
		Local:  s.local,
		Remote: s.remote,
		Paths:  ps,
	}, &SelectorMsg{})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Initialize returned")
}

func (s *SelectorClient) Path() *pan.Path {
	//log.Println("Path called")
	msg := SelectorMsg{}
	err := s.client.Call("SelectorServer.Path", &SelectorMsg{
		Remote: s.remote,
	}, &msg)
	if err != nil {
		log.Fatalln(err)
	}
	if msg.Fingerprint != nil {
		return s.paths[*msg.Fingerprint]
	}
	return nil
}

func (s *SelectorClient) PathDown(fp pan.PathFingerprint, pi pan.PathInterface) {
	log.Println("PathDown called")
	s.paths[fp] = nil // remove from local table
	err := s.client.Call("SelectorServer.PathDown", &SelectorMsg{
		Remote:        s.remote,
		Fingerprint:   &fp,
		PathInterface: &pi,
	}, &SelectorMsg{})
	if err != nil {
		log.Fatalln(err)
	}

}

func (s *SelectorClient) Refresh(paths []*pan.Path) {
	log.Println("Refresh called")
	ps := make([]*Path, len(paths))
	for i, p := range paths {
		s.paths[p.Fingerprint] = p
		ps[i] = NewPathFrom(p)
	}
	err := s.client.Call("SelectorServer.Refresh", &SelectorMsg{
		Remote: s.remote,
		Paths:  ps,
	}, &SelectorMsg{})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Refresh returned")
}

func (s *SelectorClient) Close() error {
	log.Println("Close called")
	err := s.client.Call("SelectorServer.Close", &SelectorMsg{Remote: s.remote}, &SelectorMsg{})
	if err != nil {
		log.Println(err)
		log.Println(s.client.Close())
		return err
	}
	return s.client.Close()
}
