package taps

// Listener passively waits for Connections from RemoteEndpoints.
//
// See https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-7.2
type Listener struct {
	preConn Preconnection
}

//
func newListener(preConn Preconnection) (*Listener, error) {
	if len(preConn.LocalEndpoints) == 0 {
		return nil, NewEstablishmentError("can't create listener without at least 1 local endpoint")
	}
	l := Listener{
		preConn,
	}
	return &l, nil
}

func (l *Listener) Accept() (Connection, error) {
	return Connection{}, NotYetImplementendError
}

// Stop sends the StoppedEvent and closes l's event channel
func (l *Listener) Stop() error {
	return NotYetImplementendError
}

// Stopped blocks until the Listener l is stopped.
func (l *Listener) Stopped() {

}

// SetNewConnectionLimit sets a cap on the number of inbound
// Connections that will be delivered.
func (l *Listener) SetNewConnectionLimit(value uint) {
	// TODO
}
