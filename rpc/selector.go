package rpc

import (
	//"fmt"
	"fmt"
	"io"
	"log"
	"net"
	"net/rpc"
	"time"

	"github.com/netsec-ethz/scion-apps/pkg/pan"
)

var DefaultDaemonAddress = &net.UnixAddr{
	Name: "/tmp/panapid.sock",
	Net:  "unix",
}

type Path struct {
	Source      pan.IA
	Destination pan.IA
	Metadata    *pan.PathMetadata
	Fingerprint pan.PathFingerprint
	//ForwardingPath pan.ForwardingPath
	Expiry time.Time
}

func (p *Path) PanPath() *pan.Path {
	return &pan.Path{
		Source:      p.Source,
		Destination: p.Destination,
		//Metadata:    p.Metadata.Copy(), // do we need this?
		Metadata:    p.Metadata,
		Fingerprint: p.Fingerprint,
		Expiry:      p.Expiry,
	}
}

func NewPathFrom(p *pan.Path) *Path {
	return &Path{
		Source:      p.Source,
		Destination: p.Destination,
		//Metadata:    p.Metadata.Copy(),
		Metadata:    p.Metadata,
		Fingerprint: p.Fingerprint,
		//ForwardingPath: p.ForwardingPath,
		Expiry: p.Expiry,
	}
}

type ServerSelector interface {
	Path(pan.UDPAddr) *pan.Path
	SetPaths(pan.UDPAddr, []*pan.Path)
	OnPathDown(pan.UDPAddr, pan.PathFingerprint, pan.PathInterface)
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

func (s serverSelector) SetPaths(raddr pan.UDPAddr, paths []*pan.Path) {
	s.getSelector(raddr).SetPaths(raddr, paths)
}

func (s serverSelector) OnPathDown(raddr pan.UDPAddr, fp pan.PathFingerprint, pi pan.PathInterface) {
	s.getSelector(raddr).OnPathDown(fp, pi)
}

func (s serverSelector) Close(raddr pan.UDPAddr) error {
	err := s.getSelector(raddr).Close()
	delete(s.selectors, raddr.String())
	return err
}

type Msg struct {
	Remote        *pan.UDPAddr
	Fingerprint   *pan.PathFingerprint
	PathInterface *pan.PathInterface
	Paths         []*Path
}

type SelectorServer struct {
	selector ServerSelector
}

func NewSelectorServer(selector ServerSelector) (*rpc.Server, error) {
	err := rpc.Register(&SelectorServer{selector})
	if err != nil {
		return nil, err
	}
	return rpc.DefaultServer, nil
}

func (s *SelectorServer) SetPaths(args, resp *Msg) error {
	fmt.Println("SetPaths invoked")
	paths := make([]*pan.Path, len(args.Paths))
	for i, p := range args.Paths {
		paths[i] = p.PanPath()
		//log.Printf("%s", paths[i].Source)
	}
	s.selector.SetPaths(*args.Remote, paths)
	msg := "SetPaths done"
	fmt.Println(msg)
	return nil
}

func (s *SelectorServer) Path(args, resp *Msg) error {
	//log.Println("Path invoked")
	//log.Printf("%+v", args)
	p := s.selector.Path(*args.Remote)
	//fmt.Printf("%+v", resp)
	resp.Fingerprint = &p.Fingerprint
	//log.Printf("Path done")
	return nil
}

func (s *SelectorServer) OnPathDown(args, resp *Msg) error {
	log.Println("OnPathDown called")
	s.selector.OnPathDown(*args.Remote, *args.Fingerprint, *args.PathInterface)
	return nil
}

func (s *SelectorServer) Close(args, resp *Msg) error {
	log.Println("Close called")
	return s.selector.Close(*args.Remote)
}

type SelectorClient struct {
	client *rpc.Client
	paths  map[pan.PathFingerprint]*pan.Path
	raddr  *pan.UDPAddr
}

func NewSelectorClient(conn io.ReadWriteCloser) *SelectorClient {
	client := rpc.NewClient(conn)
	log.Printf("RPC connection etablished")
	return &SelectorClient{client, map[pan.PathFingerprint]*pan.Path{}, nil}
}

func (s *SelectorClient) SetPaths(remote pan.UDPAddr, paths []*pan.Path) {
	log.Println("SetPaths called")
	if s.raddr != nil {
		if s.raddr.Equal(remote) {
			log.Fatalf("%s != %s", *s.raddr, remote)
		} else {
			log.Println("Setpaths update apparently")
		}
	} else {
		s.raddr = &remote
	}
	ps := make([]*Path, len(paths))
	for i, p := range paths {
		s.paths[p.Fingerprint] = p
		ps[i] = NewPathFrom(p)
	}
	err := s.client.Call("SelectorServer.SetPaths", &Msg{
		Remote: s.raddr,
		Paths:  ps,
	}, &Msg{})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("SetPaths returned")
}

func (s *SelectorClient) Path() *pan.Path {
	//log.Println("Path called")
	msg := Msg{}
	err := s.client.Call("SelectorServer.Path", &Msg{
		Remote: s.raddr,
	}, &msg)
	if err != nil {
		log.Fatalln(err)
	}
	//log.Printf("Returning %+v", s.paths[*msg.Fingerprint])
	return s.paths[*msg.Fingerprint]
	//return s.references[0] //path
}
func (s *SelectorClient) OnPathDown(fp pan.PathFingerprint, pi pan.PathInterface) {
	log.Println("OnPathDown called")
	s.paths[fp] = nil // remove from local table
	err := s.client.Call("SelectorServer.OnPathDown", &Msg{
		Remote:        s.raddr,
		Fingerprint:   &fp,
		PathInterface: &pi,
	}, &Msg{})
	if err != nil {
		log.Fatalln(err)
	}

}
func (s *SelectorClient) Close() error {
	log.Println("Close called")
	err := s.client.Call("SelectorServer.Close", &Msg{Remote: s.raddr}, &Msg{})
	if err != nil {
		log.Println(err)
		log.Println(s.client.Close())
		return err
	}
	return s.client.Close()
}
