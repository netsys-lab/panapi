package debug

import (
	"github.com/netsec-ethz/scion-apps/pkg/pan"
	"log"
	"time"
)

type DebugSelector struct {
	delay time.Duration
	s     pan.Selector
}

func NewDebugSelector(delay time.Duration, selector pan.Selector) (pan.Selector, error) {
	log.SetPrefix("Debug Selector")
	if selector == nil {
		selector = pan.NewDefaultSelector()
	}
	return &DebugSelector{delay, selector}, nil
}

func (s *DebugSelector) Initialize(local, remote pan.UDPAddr, paths []*pan.Path) {
	log.Println("Enter Initialize")
	time.Sleep(s.delay)
	s.s.Initialize(local, remote, paths)
	log.Println("Return Initialize")
}

func (s *DebugSelector) Path() *pan.Path {
	log.Println("Enter Path")
	time.Sleep(s.delay)
	res := s.s.Path()
	log.Println("Return Path")
	return res
}

func (s *DebugSelector) PathDown(fp pan.PathFingerprint, pi pan.PathInterface) {
	log.Println("Enter PathDown")
	time.Sleep(s.delay)
	s.s.PathDown(fp, pi)
	log.Println("Return PathDown")

}

func (s *DebugSelector) Refresh(paths []*pan.Path) {
	log.Println("Enter Refresh")
	time.Sleep(s.delay)
	s.s.Refresh(paths)
	log.Println("Return Refresh")
}

func (s *DebugSelector) Close() error {
	log.Println("Enter Close")
	time.Sleep(s.delay)
	err := s.s.Close()
	log.Println("Return Close")
	return err
}
