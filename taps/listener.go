package taps

// Listener passively waits for Connections from RemoteEndpoints.
//
// See https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-7.2
type Listener struct {
	events  chan Event
	preConn Preconnection
}

func newListener(preConn Preconnection) *Listener {
	//TODO make deep copy of preConn
	l := Listener{
		make(chan Event),
		preConn,
	}

	if len(preConn.LocalEndpoints) == 0 {
		l.events <- ErrorEvent{Error: NewEstablishmentError("can't create listener without at least 1 local endpoint")}
		l.Stop()
	}
	return &l
}

func (l *Listener) Events() <-chan Event {
	return l.events
}

// Stop sends the StoppedEvent and closes l's event channel
func (l *Listener) Stop() {
	l.events <- StoppedEvent{}
	close(l.events)
}

// SetNewConnectionLimit sets a cap on the number of inbound
// Connections that will be delivered.
func (l *Listener) SetNewConnectionLimit(value uint) {
	// TODO
}
