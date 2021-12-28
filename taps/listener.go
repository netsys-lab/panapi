package taps

// Listener passively waits for Connections from RemoteEndpoints.
//
// See https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-7.2
type Listener struct {
	Preconnection Preconnection
}

//
func newListener(Preconnection Preconnection) (*Listener, error) {
	if len(Preconnection.LocalEndpoints) == 0 {
		return nil, NewEstablishmentError("can't create listener without at least 1 local endpoint")
	}
	l := Listener{
		Preconnection,
	}
	return &l, nil
}

// Accept blocks until it receives the next incoming connection and returns it
func (l *Listener) Accept() (Connection, error) {
	return Connection{}, NotYetImplementendError
}

// Stop stops the Listener l
func (l *Listener) Stop() error {
	return NotYetImplementendError
}

/*
// Stopped blocks until the Listener l is stopped.
func (l *Listener) Stopped() {

}
*/

// SetNewConnectionLimit sets a cap on the number of inbound
// Connections that will be delivered.
func (l *Listener) SetNewConnectionLimit(value uint) {
	// TODO
}
