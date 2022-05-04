package taps

import "errors"

// Preconnection is a passive data structure that merely maintains the
// state that describes the properties of a Connection that might
// exist in the future.
type Preconnection struct {
	LocalEndpoint         *LocalEndpoint
	RemoteEndpoint        *RemoteEndpoint
	TransportPreferences  TransportPreferences
	SecurityParameters    SecurityParameters
	ConnectionPreferences *ConnectionPreferences
}

/*// NewPreconnection returns a struct representing a potential
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
*/

// Copy returns a new Preconnection struct with its values deeply
// copied from p
func (p *Preconnection) Copy() *Preconnection {
	var (
		local  *LocalEndpoint
		remote *RemoteEndpoint
		cp     *ConnectionPreferences
	)
	if p.LocalEndpoint != nil {
		local = &LocalEndpoint{*p.LocalEndpoint.Copy()}
	}
	if p.RemoteEndpoint != nil {
		remote = &RemoteEndpoint{*p.RemoteEndpoint.Copy()}
	}
	if p.ConnectionPreferences != nil {
		cp = p.ConnectionPreferences.Copy()
	}

	return &Preconnection{
		LocalEndpoint:         local,
		RemoteEndpoint:        remote,
		TransportPreferences:  *p.TransportPreferences.Copy(),
		SecurityParameters:    *p.SecurityParameters.Copy(),
		ConnectionPreferences: cp,
	}
}

// Listen returns a Listener object. Once Listen() has been called,
// any changes to the Preconnection do not have any effect on the
// Listener. The Preconnection can be disposed of or reused, e.g., to
// create another Listener.
//
// See https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-7.2
func (p *Preconnection) Listen() (Listener, error) {
	if p.LocalEndpoint == nil {
		return nil, NewEstablishmentError("can't create listener without a local endpoint")
	}
	if p.LocalEndpoint.Protocol == nil {
		return nil, NewEstablishmentError("no protocol specified")
	}
	return p.LocalEndpoint.Protocol.NewListener(p.Copy())
}

/*// Rendezvous listens on the Local Endpoint candidates for an incoming
// Connection from the Remote Endpoint candidates, while also
// simultaneously trying to establish a Connection from the Local
// Endpoint candidates to the Remote Endpoint candidates.
// See https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-7.3
func (p *Preconnection) Rendezvous() (Connection, error) {
	// TODO
	return Connection{}, NotYetImplementendError
        }*/

/*// AddRemote can add RemoteEndpoints obtained via p.Resolve() to the
// Preconnection p.
//
// Deprecated: The spec is unclear why p.Resolve() should not modify p
// directly. It is also unclear how calling AddRemote modifies the
// existing set of RemoteEndpoints configured in p. Are they
// overwritten or merely appended to? Feedback welcome.
func (p *Preconnection) AddRemote([]*RemoteEndpoint) {
	// TODO
        }*/

func (p *Preconnection) Initiate() (Connection, error) {
	if p.RemoteEndpoint == nil {
		return nil, NewEstablishmentError("can't initiate without a remote endpoint")
	}
	if p.RemoteEndpoint.Protocol == nil {
		return nil, NewEstablishmentError("no protocol specified")
	}
	return p.RemoteEndpoint.Protocol.Initiate(p.Copy())
}

func (p *Preconnection) SetPreferences(cps *ConnectionPreferences) error {
	var proto Protocol
	if p.LocalEndpoint == nil {
		if p.RemoteEndpoint == nil {
			return errors.New("No endpoint specified")
		} else {
			proto = p.RemoteEndpoint.Protocol
		}
	} else {
		proto = p.LocalEndpoint.Protocol
	}
	if proto == nil {
		return errors.New("No protocol specified")
	}
	s := proto.Selector()
	if s != nil {
		return s.SetPreferences(cps)
	} else {
		return errors.New("Can not set preferences on empty selector")
	}
}
