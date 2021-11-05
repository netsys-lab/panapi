package rpc

import (
	//"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/netsec-ethz/scion-apps/pkg/pan"
)

type Path struct {
	Source      pan.IA
	Destination pan.IA
	Metadata    *pan.PathMetadata
	Fingerprint pan.PathFingerprint
	Expiry      time.Time
}

func NewPathFromPanPath(p *pan.Path) *Path {
	return &Path{
		Source:      p.Source,
		Destination: p.Destination,
		//Metadata:    p.Metadata.Copy(),
		Fingerprint: p.Fingerprint,
		Expiry:      p.Expiry,
	}
}

func PanPathFromPath(p *Path) *pan.Path {
	return &pan.Path{
		Source:      p.Source,
		Destination: p.Destination,
		//Metadata:    p.Metadata.Copy(),
		Fingerprint: p.Fingerprint,
		Expiry:      p.Expiry,
	}
}

type SetPathArg struct {
	Remote pan.UDPAddr
	Paths  []*Path
}

type OnPathDownArg struct {
	Fingerprint   pan.PathFingerprint
	PathInterface pan.PathInterface
}

type SelectorServer struct {
	selector pan.Selector
}

type Response struct {
	Msg string
}

func NewSelectorServer(selector pan.Selector) (*SelectorServer, error) {
	s := &SelectorServer{selector}
	err := rpc.Register(s)
	if err != nil {
		return nil, err
	}
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":4711")
	if err != nil {
		return nil, err
	}
	http.Serve(l, nil)
	return s, nil
}

func (s *SelectorServer) SetPaths(arg *SetPathArg, resp *Response) error {
	log.Println("SetPaths invoked")
	paths := make([]*pan.Path, len(arg.Paths))
	for i, p := range arg.Paths {
		paths[i] = PanPathFromPath(p)
		//log.Printf("%s", paths[i].Source)
	}
	s.selector.SetPaths(arg.Remote, paths)
	msg := "SetPaths done"
	resp.Msg = msg
	log.Println(msg)
	return nil
}

func (s *SelectorServer) Path(arg *SetPathArg, path *Path) error {
	log.Println("Path invoked")

	p := s.selector.Path()
	{
		path.Source = p.Source
		path.Destination = p.Destination
		path.Fingerprint = p.Fingerprint
		path.Expiry = p.Expiry
	}
	log.Printf("Path done: %+v", *path)
	return nil
}

func (s *SelectorServer) OnPathDown(arg *OnPathDownArg, _ *interface{}) error {
	log.Println("OnPathDown called")
	s.selector.OnPathDown(arg.Fingerprint, arg.PathInterface)
	return nil
}

func (s *SelectorServer) Close(_, _ *interface{}) error {
	log.Println("Close called")
	return s.selector.Close()
}

type SelectorClient struct {
	//TODO, remove?
	server     *SelectorServer
	client     *rpc.Client
	references []*pan.Path
}

func NewSelectorClient() (*SelectorClient, error) {
	client, err := rpc.DialHTTP("tcp", "localhost:4711")
	return &SelectorClient{&SelectorServer{}, client, []*pan.Path{}}, err
}

func (s *SelectorClient) SetPaths(remote pan.UDPAddr, paths []*pan.Path) {
	log.Println("SetPaths called")
	s.references = paths
	ps := make([]*Path, len(paths))
	for i, p := range paths {
		ps[i] = NewPathFromPanPath(p)
		//log.Printf("%s", ps[i].Source)
	}
	//paths = []*pan.Path{}
	resp := Response{}
	s.client.Call("SelectorServer.SetPaths", &SetPathArg{remote, ps}, &resp)
	log.Printf("SetPaths returned: %s", resp)
}

func (s *SelectorClient) Path() *pan.Path {
	log.Println("Path called")
	path := pan.Path{}
	s.client.Call("SelectorServer.Path", &SetPathArg{}, &path)
	log.Printf("%+v", path)
	return s.references[0] //path
}
func (s *SelectorClient) OnPathDown(fp pan.PathFingerprint, pi pan.PathInterface) {
	log.Println("OnPathDown called")
	s.client.Call("SelectorServer.OnPathDown", &OnPathDownArg{fp, pi}, nil)
}
func (s *SelectorClient) Close() error {
	log.Println("Close called")
	return s.client.Call("SelectorServer.Close", nil, nil)
}
