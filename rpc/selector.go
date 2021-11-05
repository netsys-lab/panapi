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

type Msg struct {
	Remote        *pan.UDPAddr
	Fingerprint   *pan.PathFingerprint
	PathInterface *pan.PathInterface
	Paths         []*Path
}

type SelectorServer struct {
	selector pan.Selector
}

func NewSelectorServer(selector pan.Selector) (*rpc.Server, error) {
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
	p := s.selector.Path()
	//fmt.Printf("%+v", resp)
	resp.Fingerprint = &p.Fingerprint
	//log.Printf("Path done")
	return nil
}

func (s *SelectorServer) OnPathDown(args, resp *Msg) error {
	log.Println("OnPathDown called")
	s.selector.OnPathDown(*args.Fingerprint, *args.PathInterface)
	return nil
}

func (s *SelectorServer) Close(args, resp *Msg) error {
	log.Println("Close called")
	return s.selector.Close()
}

type SelectorClient struct {
	client *rpc.Client
	paths  map[pan.PathFingerprint]*pan.Path
}

func NewSelectorClient(conn io.ReadWriteCloser) *SelectorClient {
	client := rpc.NewClient(conn)
	log.Printf("RPC connection etablished")
	return &SelectorClient{client, map[pan.PathFingerprint]*pan.Path{}}
}

func (s *SelectorClient) SetPaths(remote pan.UDPAddr, paths []*pan.Path) {
	log.Println("SetPaths called")
	ps := make([]*Path, len(paths))
	for i, p := range paths {
		s.paths[p.Fingerprint] = p
		ps[i] = NewPathFrom(p)
	}
	err := s.client.Call("SelectorServer.SetPaths", &Msg{Remote: &remote, Paths: ps}, &Msg{})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("SetPaths returned")
}

func (s *SelectorClient) Path() *pan.Path {
	//log.Println("Path called")
	msg := Msg{}
	err := s.client.Call("SelectorServer.Path", &Msg{}, &msg)
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
	err := s.client.Call("SelectorServer.OnPathDown", &Msg{Fingerprint: &fp, PathInterface: &pi}, nil)
	if err != nil {
		log.Fatalln(err)
	}

}
func (s *SelectorClient) Close() error {
	log.Println("Close called")
	err := s.client.Call("SelectorServer.Close", &Msg{}, &Msg{})
	if err != nil {
		log.Println(err)
		log.Println(s.client.Close())
		return err
	}
	return s.client.Close()
}
