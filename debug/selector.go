package debug

import (
	"log"
	"time"

	"github.com/netsec-ethz/scion-apps/pkg/pan"
)

type DebugSelector struct {
	delay time.Duration
	s     pan.Selector
}

func NewDebugSelector(delay time.Duration, selector pan.Selector) (pan.Selector, error) {
	log.SetPrefix("Debug Selector")
	if selector == nil {
		selector = &pan.DefaultSelector{}
	}
	return &DebugSelector{delay, selector}, nil
}

func (s *DebugSelector) SetPaths(remote pan.UDPAddr, paths []*pan.Path) {
	log.Println("Enter SetPaths")
	time.Sleep(s.delay)
	ps := make([]*pan.Path, len(paths))
	for i, p := range paths {
		ps[i] = &pan.Path{
			Source:         p.Source,
			Destination:    p.Destination,
			ForwardingPath: p.ForwardingPath,
			//Metadata:    p.Metadata.Copy(),
			//Fingerprint: p.Fingerprint,
			//Expiry:      p.Expiry,
		}
	}
	s.s.SetPaths(remote, ps)
	log.Println("Return SetPaths")
}

func (s *DebugSelector) Path() *pan.Path {
	log.Println("Enter Path")
	time.Sleep(s.delay)
	res := s.s.Path()
	log.Println("Return Path")
	return res
}

func (s *DebugSelector) OnPathDown(fp pan.PathFingerprint, pi pan.PathInterface) {
	log.Println("Enter OnPathDown")
	time.Sleep(s.delay)
	s.s.OnPathDown(fp, pi)
	log.Println("Return OnPathDown")

}

func (s *DebugSelector) Close() error {
	log.Println("Enter Close")
	time.Sleep(s.delay)
	err := s.s.Close()
	log.Println("Return Close")
	return err
}
