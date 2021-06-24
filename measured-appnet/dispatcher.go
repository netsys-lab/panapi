package measured_appnet

import (
	"context"
	"net"

	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/snet"
	"github.com/scionproto/scion/go/lib/sock/reliable"
)

type TimedPacketDispatcherService struct {
	Dispatcher reliable.Dispatcher
	SCMPHandler snet.SCMPHandler
}

func (s *TimedPacketDispatcherService) Register(ctx context.Context, ia addr.IA,
	registration *net.UDPAddr, svc addr.HostSVC) (snet.PacketConn, uint16, error) {

	rconn, port, err := s.Dispatcher.Register(ctx, ia, registration, svc)
	if err != nil {
		return nil, 0, err
	}
	return &TimedSCIONPacketConn{
		conn:        rconn,
		scmpHandler: s.SCMPHandler,
	}, port, nil
}
