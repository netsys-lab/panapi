package measured_appnet

import (
	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/snet"
	"github.com/scionproto/scion/go/lib/sock/reliable"
)

func NewNetwork(ia addr.IA, dispatcher reliable.Dispatcher,
	revHandler snet.RevocationHandler) *snet.SCIONNetwork {

	return &snet.SCIONNetwork{
		LocalIA: ia,
		Dispatcher: &TimedPacketDispatcherService{
			Dispatcher: dispatcher,
			SCMPHandler: &snet.DefaultSCMPHandler{
				RevocationHandler: revHandler,
			},
		},
	}
}
