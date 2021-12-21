package taps

import "fmt"

// Listener passively waits for Connections from RemoteEndpoints.
//
// Deprecated : It is at this point unclear, whether representing the
// notion of "Event"s from the TAPS draft as Go channels is a good
// idea. Feedback welcome.
//
// See https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-7.2
type Listener struct {
	// ConnectionReceived is a channel that receives a new
	// Connection when a Remote Endpoint has established a
	// transport-layer connection to this Listener (for
	// Connection-oriented transport protocols), or when the first
	// Message has been received from the Remote Endpoint (for
	// Connectionless protocols), causing a new Connection to be
	// created.
	ConnectionReceived chan Connection

	// Stopped is closed when the Listener has stopped listening
	Stopped chan struct{}

	// Error is a channel that receives an error, either when the
	// Properties and Security Parameters of the Preconnection
	// cannot be fulfilled for listening or cannot be reconciled
	// with the Local Endpoint (and/or Remote Endpoint, if
	// specified), when the Local Endpoint (or Remote Endpoint, if
	// specified) cannot be resolved, or when the application is
	// prohibited from listening by policy.
	Error chan error

	preConn Preconnection
}

func newListener(preConn Preconnection) *Listener {
	//TODO make deep copy of preConn
	l := Listener{
		make(chan Connection),
		make(chan struct{}),
		make(chan error),
		preConn,
	}

	if len(preConn.LocalEndpoints) == 0 {
		l.Error <- fmt.Errorf("no local endpoint for listening specified")
		l.Stop()
	}
	return &l
}

// Stop closes the channels Stopped, Error and ConnectionReceived channels of Listener l in order.
func (l *Listener) Stop() {
	close(l.Stopped)
	close(l.Error)
	close(l.ConnectionReceived)
}

// SetNewConnectionLimit sets a cap on the number of inbound
// Connections that will be delivered.
func (l *Listener) SetNewConnectionLimit(value uint) {
	// TODO
}
