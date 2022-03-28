package taps

// Preconnection is a passive data structure that merely maintains the
// state that describes the properties of a Connection that might
// exist in the future.
type Preconnection struct {
	LocalEndpoints      []*LocalEndpoint
	RemoteEndpoints     []*RemoteEndpoint
	TransportProperties *TransportProperties
	SecurityParameters  *SecurityParameters
	Framers             []*Framer
}

// NewPreconnection returns a struct representing a potential
// Connection.
//
// At least one Local Endpoint MUST be specified if the Preconnection
// is used to Listen() for incoming Connections, but the list of Local
// Endpoints MAY be empty if the Preconnection is used to Initiate()
// connections.
//
//  myPreconnection := NewPreconnection(
//    myEndpoints,
//    []*RemoteEndpoint{},   // leave empty, we just want to Listen()
//    myTransportProperties,
//    nil,                   // no SecurityParameters for this NewPreconnection
//  )
//
// Note that it would be idiomatic Go to do the above with a struct
// literal, leaving all unset Fields at their appropriate zero value
//  myPreconnection := Preconnection{
//    LocalEndpoints: myEndpoints,
//    TransportProperties: myTransportProperties,
//  }
//
// See https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6
func NewPreconnection(
	LocalEndpoints []*LocalEndpoint,
	RemoteEndpoints []*RemoteEndpoint,
	TransportProperties *TransportProperties,
	SecurityParameters *SecurityParameters,
) *Preconnection {
	return &Preconnection{
		LocalEndpoints:      LocalEndpoints,
		RemoteEndpoints:     RemoteEndpoints,
		TransportProperties: TransportProperties,
		SecurityParameters:  SecurityParameters,
	}
}

// Copy returns a new Preconnection struct with its values deeply
// copied from p
func (p *Preconnection) Copy() *Preconnection {
	locals := make([]*LocalEndpoint, len(p.LocalEndpoints))
	for i := 0; i < len(p.LocalEndpoints); i += 1 {
		locals[i] = &LocalEndpoint{*p.LocalEndpoints[i].Copy()}
	}
	remotes := make([]*RemoteEndpoint, len(p.RemoteEndpoints))
	for i := 0; i < len(p.RemoteEndpoints); i += 1 {
		remotes[i] = &RemoteEndpoint{*p.RemoteEndpoints[i].Copy()}
	}

	return &Preconnection{
		LocalEndpoints:      locals,
		RemoteEndpoints:     remotes,
		TransportProperties: p.TransportProperties.Copy(),
		SecurityParameters:  p.SecurityParameters.Copy(),
	}
}

// Resolve called on a Preconnection p can be used by the application
// to force early binding when required, for example with some Network
// Address Translator (NAT) traversal protocols.
//
// See
// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.1
// and
// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-7.3
func (p *Preconnection) Resolve() (les []*LocalEndpoint, res []*RemoteEndpoint, err error) {
	err = NotYetImplementendError
	return
}

// Listen returns a Listener object. Once Listen() has been called,
// any changes to the Preconnection do not have any effect on the
// Listener. The Preconnection can be disposed of or reused, e.g., to
// create another Listener.
//
// See https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-7.2
func (p *Preconnection) Listen() (*Listener, error) {
	return newListener(*p.Copy())
}

// Rendezvous listens on the Local Endpoint candidates for an incoming
// Connection from the Remote Endpoint candidates, while also
// simultaneously trying to establish a Connection from the Local
// Endpoint candidates to the Remote Endpoint candidates.
// See https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-7.3
func (p *Preconnection) Rendezvous() (Connection, error) {
	// TODO
	return Connection{}, NotYetImplementendError
}

// AddRemote can add RemoteEndpoints obtained via p.Resolve() to the
// Preconnection p.
//
// Deprecated: The spec is unclear why p.Resolve() should not modify p
// directly. It is also unclear how calling AddRemote modifies the
// existing set of RemoteEndpoints configured in p. Are they
// overwritten or merely appended to? Feedback welcome.
func (p *Preconnection) AddRemote([]*RemoteEndpoint) {
	// TODO
}

func (p *Preconnection) Initiate() (Connection, error) {
	return Connection{}, NotYetImplementendError
}
