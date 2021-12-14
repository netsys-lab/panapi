// Copyright 2021 Thorben Krüger (thorben.krueger@ovgu.de)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package rpc

import (
	"errors"
	"log"
	"net"

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
	Initialize(pan.UDPAddr, pan.UDPAddr, []*pan.Path) error
	Path(pan.UDPAddr, pan.UDPAddr) (*pan.Path, error)
	PathDown(pan.UDPAddr, pan.UDPAddr, pan.PathFingerprint, pan.PathInterface) error
	Refresh(pan.UDPAddr, pan.UDPAddr, []*pan.Path) error
	Close(pan.UDPAddr, pan.UDPAddr) error
}

type serverSelector struct {
	fn        func(pan.UDPAddr, pan.UDPAddr) pan.Selector
	selectors map[string]pan.Selector
}

func NewServerSelectorFunc(fn func(pan.UDPAddr, pan.UDPAddr) pan.Selector) ServerSelector {
	return serverSelector{fn, map[string]pan.Selector{}}
}

func (s *serverSelector) getSelector(local, remote pan.UDPAddr) pan.Selector {
	addr := local.String() + remote.String()
	selector, ok := s.selectors[addr]
	if !ok {
		selector = s.fn(local, remote)
		s.selectors[addr] = selector
	}
	return selector
}

func (s serverSelector) Path(local, remote pan.UDPAddr) (*pan.Path, error) {
	return s.getSelector(local, remote).Path(), nil
}

func (s serverSelector) Initialize(local, remote pan.UDPAddr, paths []*pan.Path) error {
	s.getSelector(local, remote).Initialize(local, remote, paths)
	return nil
}

func (s serverSelector) PathDown(local, remote pan.UDPAddr, fp pan.PathFingerprint, pi pan.PathInterface) error {
	s.getSelector(local, remote).PathDown(fp, pi)
	return nil
}

func (s serverSelector) Refresh(local, remote pan.UDPAddr, paths []*pan.Path) error {
	s.getSelector(local, remote).Refresh(paths)
	return nil
}

func (s serverSelector) Close(local, remote pan.UDPAddr) error {
	err := s.getSelector(local, remote).Close()
	delete(s.selectors, local.String()+remote.String())
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
	log.Println("Initialize invoked")
	paths := make([]*pan.Path, len(args.Paths))
	for i, p := range args.Paths {
		paths[i] = p.PanPath()
		//log.Printf("%s", paths[i].Source)
	}
	if args.Local == nil || args.Remote == nil {
		return ErrDeref
	}
	return s.selector.Initialize(*args.Local, *args.Remote, paths)
}

func (s *SelectorServer) Path(args, resp *SelectorMsg) error {
	if args.Remote == nil {
		return ErrDeref
	}
	p, err := s.selector.Path(*args.Local, *args.Remote)
	if p != nil {
		resp.Fingerprint = &p.Fingerprint
	}
	return err
}

func (s *SelectorServer) PathDown(args, resp *SelectorMsg) error {
	log.Println("PathDown called")
	if args.Remote == nil || args.Fingerprint == nil || args.PathInterface == nil {
		return ErrDeref
	}
	return s.selector.PathDown(*args.Local, *args.Remote, *args.Fingerprint, *args.PathInterface)
}

func (s *SelectorServer) Refresh(args, resp *SelectorMsg) error {
	log.Println("Refresh invoked")
	paths := make([]*pan.Path, len(args.Paths))
	for i, p := range args.Paths {
		paths[i] = p.PanPath()
		//log.Printf("%s", paths[i].Source)
	}
	if args.Remote == nil {
		return ErrDeref
	}
	return s.selector.Refresh(*args.Local, *args.Remote, paths)
}

func (s *SelectorServer) Close(args, resp *SelectorMsg) error {
	log.Println("Close called")
	if args.Remote == nil {
		return ErrDeref
	}
	return s.selector.Close(*args.Local, *args.Remote)
}

type SelectorClient struct {
	client *Client
	paths  map[pan.PathFingerprint]*pan.Path
	local  *pan.UDPAddr
	remote *pan.UDPAddr
	l      *log.Logger
}

func NewSelectorClient(client *Client) *SelectorClient {
	client.l.Printf("RPC connection etablished")
	return &SelectorClient{client, map[pan.PathFingerprint]*pan.Path{}, nil, nil, client.l}
}

func (s *SelectorClient) Initialize(local, remote pan.UDPAddr, paths []*pan.Path) {
	s.l.Println("Initialize called")
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
		s.l.Fatalln(err)
	}
	s.l.Printf("Initialize returned")
}

func (s *SelectorClient) Path() *pan.Path {
	//s.l.Println("Path called")
	msg := SelectorMsg{}
	err := s.client.Call("SelectorServer.Path", &SelectorMsg{
		Local:  s.local,
		Remote: s.remote,
	}, &msg)
	if err != nil {
		s.l.Fatalln(err)
	}
	if msg.Fingerprint != nil {
		return s.paths[*msg.Fingerprint]
	}
	return nil
}

func (s *SelectorClient) PathDown(fp pan.PathFingerprint, pi pan.PathInterface) {
	s.l.Println("PathDown called")
	s.paths[fp] = nil // remove from local table
	err := s.client.Call("SelectorServer.PathDown", &SelectorMsg{
		Local:         s.local,
		Remote:        s.remote,
		Fingerprint:   &fp,
		PathInterface: &pi,
	}, &SelectorMsg{})
	if err != nil {
		s.l.Fatalln(err)
	}

}

func (s *SelectorClient) Refresh(paths []*pan.Path) {
	s.l.Println("Refresh called")
	ps := make([]*Path, len(paths))
	for i, p := range paths {
		s.paths[p.Fingerprint] = p
		ps[i] = NewPathFrom(p)
	}
	err := s.client.Call("SelectorServer.Refresh", &SelectorMsg{
		Local:  s.local,
		Remote: s.remote,
		Paths:  ps,
	}, &SelectorMsg{})
	if err != nil {
		s.l.Fatalln(err)
	}
	s.l.Printf("Refresh returned")
}

func (s *SelectorClient) Close() error {
	s.l.Println("Close called")
	err := s.client.Call("SelectorServer.Close", &SelectorMsg{Local: s.local, Remote: s.remote}, &SelectorMsg{})
	if err != nil {
		s.l.Println(err)
		s.l.Println(s.client.client.Close())
		return err
	}
	return s.client.client.Close()
}
