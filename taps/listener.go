package taps

import "fmt"

// Listener passively waits for Connections from RemoteEndpoints.
//
// See https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-7.2
type Listener struct {
	Events  chan Event
	preConn Preconnection
}

func newListener(preConn Preconnection) *Listener {
	//TODO make deep copy of preConn
	l := Listener{
		make(chan Event),
		preConn,
	}

	if len(preConn.LocalEndpoints) == 0 {
		l.Events <- EstablishmentErrorEvent{Error: fmt.Errorf("no local endpoint for listening specified")}
		l.Stop()
	}
	return &l
}

// Stop sends the StoppedEvent and closes the Events channel
func (l *Listener) Stop() {
	l.Events <- StoppedEvent{}
	close(l.Events)
}

// SetNewConnectionLimit sets a cap on the number of inbound
// Connections that will be delivered.
func (l *Listener) SetNewConnectionLimit(value uint) {
	// TODO
}
