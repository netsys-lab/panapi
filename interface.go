package panapi

import (
	"crypto"
	"crypto/tls"
	"fmt"
	"reflect"
	"time"
)

type Preference uint8

const (
	unset    Preference = iota // (Implementation detail: Indicate that recommended default value for Property should be used)
	Ignore                     // No preference
	Require                    // Select only protocols/paths providing the property, fail otherwise
	Prefer                     // Prefer protocols/paths providing the property, proceed otherwise
	Avoid                      // Prefer protocols/paths not providing the property, proceed otherwise
	Prohibit                   // Select only protocols/paths not providing the property, fail otherwise
)

func (p Preference) String() string {
	return [...]string{
		"(unset)",
		"Ignore",
		"Require",
		"Prefer",
		"Avoid",
		"Prohibit",
	}[p-unset]
}

type MultiPathPreference uint8

const (
	dynamic  MultiPathPreference = iota // (Implementation detail: need to use different defaults depending on endpoint)
	Disabled                            // The connection will not use multiple paths once established, even if the chosen transport supports using multiple paths.
	Active                              // The connection will negotiate the use of multiple paths if the chosen transport supports this.
	Passive                             // The connection will support the use of multiple paths if the Remote Endpoint requests it.
)

func (p MultiPathPreference) String() string {
	return [...]string{
		"(unset)",
		"Disabled",
		"Active",
		"Passive",
	}[p-dynamic]
}

type Directionality uint8

const (
	Bidirectional         Directionality = iota // The connection must support sending and receiving data
	UnidirectionalSend                          // The connection must support sending data, and the application cannot use the connection to receive any data
	UnidirectionalReceive                       // The connection must support receiving data, and the application cannot use the connection to send any data

)

func (d Directionality) String() string {
	return [...]string{
		"Bidirectional",
		"Unidirectional Send",
		"Unidirectional Receive",
	}[d-Bidirectional]
}

type TransportProperties struct {
	// Reliability pecifies whether the application needs to use a
	// transport protocol that ensures that all data is received
	// at the Remote Endpoint without corruption. When reliable
	// data transfer is enabled, this also entails being notified
	// when a Connection is closed or aborted. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2.1)
	Reliability Preference

	// PreserveMsgBoundaries specifies whether the application
	// needs or prefers to use a transport protocol that preserves
	// message boundaries. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2.2)
	PreserveMsgBoundaries Preference

	// PerMsgReliability specifies whether an application
	// considers it useful to specify different reliability
	// requirements for individual Messages in a Connection. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2.3)
	PerMsgReliability Preference

	// PreserveOrder specifies whether the application wishes to
	// use a transport protocol that can ensure that data is
	// received by the application on the other end in the same
	// order as it was sent. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2.4)
	PreserveOrder Preference

	// ZeroRTTMsg specifies whether an application would like to
	// supply a Message to the transport protocol before
	// Connection establishment that will then be reliably
	// transferred to the other side before or during Connection
	// establishment. This Message can potentially be received
	// multiple times (i.e., multiple copies of the message data
	// may be passed to the Remote Endpoint). (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2.5)
	ZeroRTTMsg Preference

	// Multistreaming specifies that the application would prefer
	// multiple Connections within a Connection Group to be
	// provided by streams of a single underlying transport
	// connection where possible.  (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2.6)
	Multistreaming Preference

	// FullChecksumSend specifies the application's need for
	// protection against corruption for all data transmitted on
	// this Connection. Disabling this property could enable later
	// control of the sender checksum coverage. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2.7)
	FullChecksumSend Preference

	// FullChecksumRecv specifies the application's need for
	// protection against corruption for all data received on this
	// Connection. Disabling this property could enable later
	// control of the required minimum receiver checksum
	// coverage. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2.8)
	FullChecksumRecv Preference

	// CongestionControl specifies whether the application would like
	// the Connection to be congestion controlled or not. Note
	// that if a Connection is not congestion controlled, an
	// application using such a Connection SHOULD itself perform
	// congestion control in accordance with [RFC2914] or use a
	// circuit breaker in accordance with [RFC8084], whichever is
	// appropriate. Also note that reliability is usually combined
	// with congestion control in protocol implementations,
	// rendering "reliable but not congestion controlled" a
	// request that is unlikely to succeed. If the Connection is
	// congestion controlled, performing additional congestion
	// control in the application can have negative performance
	// implications. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2.9)
	CongestionControl Preference

	// KeepAlive specifies whether the application would like the
	// Connection to send keep-alive packets or not. Note that if
	// a Connection determines that keep-alive packets are being
	// sent, the applicaton should itself avoid generating
	// additional keep alive messages. Note that when supported,
	// the system will use the default period for generation of
	// the keep alive-packets. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2.10)
	KeepAlive Preference

	// Interface allows the application to select any specific
	// network interfaces or categories of interfaces it wants to
	// Require, Prohibit, Prefer, or Avoid. Note that marking a
	// specific interface as Require strictly limits path
	// selection to that single interface, and often leads to less
	// flexible and resilient connection establishment. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2.11)
	//
	// Implementation Detail: In contrast to other Selection
	// Properties, this property maps interface identifier strings
	// to Preferences. In future, common interface types might
	// exist as constants.
	Interface map[string]Preference

	// PvD allows the application to control path selection by
	// selecting which specific Provisioning Domain (PvD) or
	// categories of PVDs it wants to Require, Prohibit, Prefer,
	// or Avoid. Provisioning Domains define consistent sets of
	// network properties that may be more specific than network
	// interfaces [RFC7556]. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2.12)
	//
	// Implementation Detail: In contrast to other Selection
	// Properties, this property maps PvD identifier strings to
	// Preferences. In future, common PvD types and categories
	// might exist as constants.
	PvD map[string]Preference

	// Multipath specifies whether and how applications want to
	// take advantage of transferring data across multiple paths
	// between the same end hosts. Using multiple paths allows
	// connections to migrate between interfaces or aggregate
	// bandwidth as availability and performance properties
	// change. Possible values are:
	//
	// Disabled: The connection will not use multiple paths once
	//    established, even if the chosen transport supports using
	//    multiple paths.
	//
	// Active: The connection will negotiate the use of multiple
	//    paths if the chosen transport supports this.
	//
	// Passive: The connection will support the use of multiple
	//    paths if the Remote Endpoint requests it.
	//
	//
	// The policy for using multiple paths is specified using the
	// separate multipath-policy property. (TODO) (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2.14)
	Multipath MultiPathPreference

	// UseTemporaryLocalAddress allows the application to express
	// a preference for the use of temporary local addresses,
	// sometimes called "privacy" addresses [RFC4941]. Temporary
	// addresses are generally used to prevent linking connections
	// over time when a stable address, sometimes called
	// "permanent" address, is not needed. There are some caveats
	// to note when specifying this property. First, if an
	// application Requires the use of temporary addresses, the
	// resulting Connection cannot use IPv4, because temporary
	// addresses do not exist in IPv4. Second, temporary local
	// addresses might involve trading off privacy for
	// performance. For instance, temporary addresses can
	// interfere with resumption mechanisms that some protocols
	// rely on to reduce initial latency. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2.13)
	UseTemporaryLocalAddress Preference

	// AdvertisesAltAddr specifies whether alternative addresses,
	// e.g., of other interfaces, should be advertised to the peer
	// endpoint by the protocol stack. Advertising these addresses
	// enables the peer-endpoint to establish additional
	// connectivity, e.g., for connection migration or using
	// multiple paths. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2.15)
	AdvertisesAltAddr bool

	// Direction specifies whether an application wants to use the connection for sending and/or receiving data.
	//
	// Since unidirectional communication can be supported by
	// transports offering bidirectional communication, specifying
	// unidirectional communication may cause a transport stack
	// that supports bidirectional communication to be
	// selected. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2.16)
	Direction Directionality

	// SoftErrorNotify specifies whether an application considers
	// it useful to be informed when an ICMP error message arrives
	// that does not force termination of a connection. When set
	// to true, received ICMP errors are available as
	// SoftErrors. Note that even if a protocol supporting this
	// property is selected, not all ICMP errors will necessarily
	// be delivered, so applications cannot rely upon receiving
	// them [RFC8085]. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2.17)
	SoftErrorNotify Preference

	// ActiveReadBeforeSend specifies whether an application wants
	// to diverge from the most common communication pattern - the
	// client actively opening a connection, then sending data to
	// the server, with the server listening (passive open),
	// reading and then answering. ActiveReadBeforeSend departs
	// from this pattern, either by actively opening with
	// Initiate(), immediately followed by reading, or passively
	// opening with Listen(), immediately followed by
	// writing. This property is ignored when establishing
	// connections using Rendezvous(). Requiring this property
	// limits the choice of mappings to underlying protocols,
	// which can reduce efficiency. For example, it prevents the
	// Transport Services system from mapping Connections to SCTP
	// streams, where the first transmitted data takes the role of
	// an active open signal. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2.18)
	ActiveReadBeforeSend Preference
}

// NewTransportProperties creates TransportProperties with the recommended defaults from https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2
func NewTransportProperties() *TransportProperties {
	return &TransportProperties{
		Reliability:              Require,
		PreserveMsgBoundaries:    Ignore,
		PerMsgReliability:        Ignore,
		PreserveOrder:            Require,
		ZeroRTTMsg:               Ignore,
		Multistreaming:           Prefer,
		FullChecksumSend:         Require,
		FullChecksumRecv:         Require,
		CongestionControl:        Require,
		KeepAlive:                Ignore,
		Interface:                map[string]Preference{},
		PvD:                      map[string]Preference{},
		UseTemporaryLocalAddress: unset,   // Needs to be resolved at runtime: Avoid for Listeners and Rendezvous Connections, else Prefer
		Multipath:                dynamic, // Needs to be resolved at runtime: Disabled for Initiated and Rendezvous Connections, else Passive
		AdvertisesAltAddr:        false,
		Direction:                Bidirectional,
		SoftErrorNotify:          Ignore,
		ActiveReadBeforeSend:     Ignore,
	}
}

// SetProperty is not yet implemented
func (tp *TransportProperties) SetProperty(property, value string) {

}

// SetPreference stores preference for property.
//
// CAUTION: Use SetPreference only if you must. Direct access of
// struct Fields is usually preferred. This function is implemented
// using reflection and probably a bit slow. It inherently also could
// have runtime bugs.
func (tp *TransportProperties) SetPreference(property string, preference Preference) error {
	s := reflect.ValueOf(tp).Elem()
	f := s.FieldByName(property)
	if !f.IsValid() {
		return fmt.Errorf("Preference %s not found", property)
	}
	p := reflect.ValueOf(preference)
	if p.Type().AssignableTo(f.Type()) {
		f.Set(p)
	} else {
		return fmt.Errorf("Can not assign a Preference to Property %s", property)
	}
	return nil
}

/* GetPreference returns the stored Preference for property, or Ignore if preference is not found

   Deprecated: GetX is not idiomatic Go, use X instead, i.e., "Preference(property)"

func (tp *TransportProperties) GetPreference(property string) Preference {
	return tp.Preference(property)
}

// Preference returns the stored Preference for property, or Ignore if preference is not found
func (tp *TransportProperties) Preference(property string) Preference {
	p, ok := tp.preferences[property]
	if !ok {
		return Ignore
	}
	return p
        }*/

// Require has the effect of selecting only protocols/paths providing the property and failing otherwise
func (tp *TransportProperties) Require(property string) {
	tp.SetPreference(property, Require)

}

// Prefer has the effect of prefering protocols/paths providing the property and proceeding otherwise
func (tp *TransportProperties) Prefer(property string) {
	tp.SetPreference(property, Prefer)
}

// Ignore has the effect of expressing no preference for the given property
func (tp *TransportProperties) Ignore(property string) {
	tp.SetPreference(property, Ignore)

}

// Avoid has the effect of avoiding protocols/paths with the property and proceeding otherwise
func (tp *TransportProperties) Avoid(property string) {
	tp.SetPreference(property, Avoid)
}

// Prohibit has the effect of failing if the property can not be avoided
func (tp *TransportProperties) Prohibit(property string) {
	tp.SetPreference(property, Prohibit)
}

// SecurityParameters is a structure used to configure security for a Preconnection
type SecurityParameters struct {
	Identity              string //? []byte?
	KeyPair               KeyPair
	SupportedGroup        tls.CurveID
	CipherSuite           tls.CipherSuite
	SignatureAlgorithm    tls.SignatureScheme
	MaxCachedSessions     uint
	CachedSessionLifetime time.Duration
	//PSK not supported
}

// SetTrustVerificationCallback is not yet implemented
func (sp SecurityParameters) SetTrustVerificationCallback() {
}

// SetIdentityChallengeCallback is not yet implemented
func (sp SecurityParameters) SetIdentityChallengeCallback() {
}

type CapacityProfile uint8

const (
	Default CapacityProfile = iota
	Scavenger
	LowLatencyInteractive
	LowLatencyNonInteractive
	ConstantRateStreaming
	CapacitySeeking
)

var profileNames = [...]string{
	"Default",
	"Scavenger",
	"Low Latency/Interactive",
	"Low Latency/Non-Interactive",
	"Constant-Rate Streaming",
	"Capacity-Seeking",
}

func (p CapacityProfile) String() string {
	return profileNames[p-Default]
}

type ConnectionProperties struct {
	RecvChecksumLen uint
	ConnPrio        uint
	ConnTimeout     time.Duration
	// ConnScheduler not yet implemented
	ConnCapacityProfile CapacityProfile
	/*
	   MultipathPolicy
	   MinSendRate
	   MinRecvRate
	   MaxSendRate
	   MaxRecvRate
	   GroupConnLimit
	   IsolateSession
	   not yet implemented
	*/
}

// KeyPair clearly associates a Private and Public Key into a pair
type KeyPair struct {
	PrivateKey crypto.PrivateKey
	PublicKey  crypto.PublicKey
}

type Connection struct {
	Ready      chan bool
	SoftError  chan error
	PathChange chan bool
}
