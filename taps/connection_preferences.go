package taps

import (
	"time"
)

type ConnectionPreferences struct {
	ConnTimeout time.Duration
	// ConnCapacityProfile specifies the desired network treatment
	// for traffic sent by the application and the tradeoffs the
	// application is prepared to make in path and protocol
	// selection to receive that desired treatment. When the
	// capacity profile is set to a value other than Default, the
	// Transport Services system SHOULD select paths and configure
	// protocols to optimize the tradeoff between delay, delay
	// variation, and efficient use of the available capacity
	// based on the capacity profile specified. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-8.1.6)
	ConnCapacityProfile CapacityProfile

	// MultipathPolicy specifies the local policy for transferring
	// data across multiple paths between the same end hosts if
	// Multipath is not set to Disabled in TransportPreference. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-8.1.7)
	MultipathPolicy MultipathPolicy

	IsolateSession bool
}

// Copy returns a new ConnectionPreferences struct with its values deeply copied from cp
func (cp *ConnectionPreferences) Copy() *ConnectionPreferences {
	return &ConnectionPreferences{
		ConnTimeout:         cp.ConnTimeout,
		ConnCapacityProfile: cp.ConnCapacityProfile,
		MultipathPolicy:     cp.MultipathPolicy,
		IsolateSession:      cp.IsolateSession,
	}
}
