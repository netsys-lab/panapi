package taps

import (
	"time"
)

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
	// change.
	//
	// The policy for using multiple paths is specified using the
	// separate MultipathPolicy ConnectionProperty. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2.14)
	Multipath MultipathPreference

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

// NewTransportProperties creates TransportProperties with the
// recommended defaults from
// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2
func NewTransportProperties() *TransportProperties {
	return &TransportProperties{
		Reliability:           Require,
		PreserveMsgBoundaries: Ignore,
		PerMsgReliability:     Ignore,
		PreserveOrder:         Require,
		ZeroRTTMsg:            Ignore,
		Multistreaming:        Prefer,
		FullChecksumSend:      Require,
		FullChecksumRecv:      Require,
		CongestionControl:     Require,
		KeepAlive:             Ignore,
		Interface:             map[string]Preference{},
		PvD:                   map[string]Preference{},
		// Needs to be resolved at runtime: Avoid for Listeners and Rendezvous Connections, else Prefer
		UseTemporaryLocalAddress: unset,
		// Needs to be resolved at runtime: Disabled for Initiated and Rendezvous Connections, else Passive
		Multipath:            dynamic,
		AdvertisesAltAddr:    false,
		Direction:            Bidirectional,
		SoftErrorNotify:      Ignore,
		ActiveReadBeforeSend: Ignore,
	}
}

// Set stores value for property, which is stripped of case and
// non-alphabetic characters before being matched against the (equally
// stripped) exported Field names of tp. The type of value must be
// assignable to type of the targeted property Field, otherwise an
// error is returned.
//
// For the sake of respecting the TAPS (draft) spec as closely as
// possible, this function allows you to say:
//  err := tp.Set("preserve-msg-boundaries", Require)
//  if err != nil {
//    ... // handle runtime error
//  }
//
// In idiomatic Go, you would (and should) instead say:
//  tp.PreserveMsgBoundaries = Require
//
// Deprecated: Use func tp.Set only if you must. Direct access of the
// TransportProperties struct Fields is usually preferred. This
// function is implemented using reflection and dynamic string
// matching, which is inherently inefficient and prone to bugs
// triggered at runtime.
func (tp *TransportProperties) Set(property string, value interface{}) error {
	return set(tp, property, value)
}

// Require has the effect of selecting only protocols/paths providing
// the property and failing otherwise.
//
// Deprecated: This is equivalent to calling tp.Set(property, Require) - the caveats
// of func tp.Set apply in full
func (tp *TransportProperties) Require(property string) error {
	return tp.Set(property, Require)
}

// Prefer has the effect of prefering protocols/paths providing the
// property and proceeding otherwise
//
// Deprecated: This is equivalent to calling tp.Set(property, Prefer) - the caveats
// of func tp.Set apply in full
func (tp *TransportProperties) Prefer(property string) error {
	return tp.Set(property, Prefer)
}

// Ignore has the effect of expressing no preference for the given property
//
// Deprecated: This is equivalent to calling tp.Set(property, Ignore) - the caveats
// of func tp.Set apply in full
func (tp *TransportProperties) Ignore(property string) error {
	return tp.Set(property, Ignore)

}

// Avoid has the effect of avoiding protocols/paths with the property
// and proceeding otherwise
//
// Deprecated: This is equivalent to calling tp.Set(property, Avoid) - the caveats
// of func tp.Set apply in full
func (tp *TransportProperties) Avoid(property string) error {
	return tp.Set(property, Avoid)
}

// Prohibit has the effect of failing if the property can not be
// avoided
//
// Deprecated: This is equivalent to calling tp.Set(property, Prohibit) - the
// caveats of func tp.Set apply in full
func (tp *TransportProperties) Prohibit(property string) error {
	return tp.Set(property, Prohibit)
}

type ConnectionProperties struct {
	// RecvChecksumLen specifies the minimum number of bytes in a
	// received message that need to be covered by a checksum. A
	// special value of 0 means that a received packet does not
	// need to have a non-zero checksum field. A receiving
	// endpoint will not forward messages that have less coverage
	// to the application. The application is responsible for
	// handling any corruption within the non-protected part of
	// the message [RFC8085]. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-8.1.1)
	RecvChecksumLen uint

	// ConnPrio is a non-negative integer representing the
	// relative inverse priority (i.e., a lower value reflects a
	// higher priority) of this Connection relative to other
	// Connections in the same Connection Group. It has no effect
	// on Connections not part of a Connection Group. This
	// property is not entangled when Connections are cloned,
	// i.e., changing the Priority on one Connection in a
	// Connection Group does not change it on the other
	// Connections in the same Connection Group. No guarantees of
	// a specific behavior regarding Connection Priority are
	// given; a Transport Services system may ignore this
	// property. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-8.1.2)
	ConnPrio uint

	// ConnTimeout specifies how long to wait before deciding that
	// an active Connection has failed when trying to reliably
	// deliver data to the Remote Endpoint. Adjusting this
	// Property will only take effect when the underlying stack
	// supports reliability. A value of 0 means that no timeout is
	// scheduled. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-8.1.3)
	ConnTimeout time.Duration

	// KeepAliveTimeout specifies the maximum length of time an
	// idle connection (one for which no transport packets have
	// been sent) should wait before the Local Endpoint sends a
	// keep-alive packet to the Remote Endpoint. Adjusting this
	// Property will only take effect when the underlying stack
	// supports sending keep-alive packets. Guidance on setting
	// this value for datagram transports is provided in
	// [RFC8085]. A value greater than ConnTimeout or the special
	// value 0 will disable the sending of keep-alive
	// packets. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-8.1.4)
	KeepAliveTimeout time.Duration

	// ConnScheduler specifies which scheduler should be used
	// among Connections within a Connection Group. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-8.1.5)
	ConnScheduler StreamScheduler

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
	// Multipath is not set to Disabled in TransportProperty. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-8.1.7)
	MultipathPolicy MultipathPolicy

	// [Max|Min][Send|Recv]Rate specifies an upper-bound rate that
	// a transfer is not expected to exceed (even if flow control
	// and congestion control allow higher rates), and/or a
	// lower-bound rate below which the application does not deem
	// it will be useful. These are specified in bits per
	// second. The special value of 0 (alternatively: a Max Rate
	// set lower than the corresponding Min Rate) indicates that
	// no bound is specified.  (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-8.1.8)
	MinSendRate, MinRecvRate, MaxSendRate, MaxRecvRate uint

	// GroupConnLimit controls the number of Connections that can
	// be accepted from a peer as new members of the Connection's
	// group. Similar to SetNewConnectionLimit(), this limits the
	// number of ConnectionReceived Events that will occur, but
	// constrained to the group of the Connection associated with
	// this property. For a multi-streaming transport, this limits
	// the number of allowed streams. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-8.1.9)
	GroupConnLimit uint

	// IsolateSession, when set, will initiate new Connections
	// using as little cached information (such as session tickets
	// or cookies) as possible from previous connections that are
	// not in the same Connection Group. Any state generated by
	// this Connection will only be shared with Connections in the
	// same Connection Group. Cloned Connections will use saved
	// state from within the Connection Group. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-8.1.10)
	IsolateSession bool

	zeroRTTMsgMaxLen, singularTransmissionMsgMaxLen, sendMsgMaxLen, recvMsgMaxLen uint
}

// ZeroRTTMsgMaxLen returns the maximum Message size that can be sent
// before or during Connection establishment in bytes.
//
// (See https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-8.1.11.1)
func (cp *ConnectionProperties) ZeroRTTMsgMaxLen() uint {
	return cp.zeroRTTMsgMaxLen
}

// SingularTransmissionMsgMaxLen, if applicable, returns the maximum
// Message size that can be sent without incurring network-layer
// fragmentation at the sender, in bytes.
//
// (See https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-8.1.11.2)
func (cp *ConnectionProperties) SingularTransmissionMsgMaxLen() uint {
	return cp.singularTransmissionMsgMaxLen
}

// SendMsgMaxLen returns the maximum Message size that an application
// can send, in bytes.
//
// (See https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-8.1.11.3)
func (cp *ConnectionProperties) SendMsgMaxLen() uint {
	return cp.sendMsgMaxLen
}

// RecvMsgMaxLen returns the maximum Message size that an application
// can receive, in bytes.
//
// (See https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-8.1.11.4)
func (cp *ConnectionProperties) RecvMsgMaxLen() uint {
	return cp.recvMsgMaxLen
}

// Get returns value associated with Field property of cp. If property
// is not a Field of cp, an error is returned
//
// For the sake of respecting the TAPS (draft) spec, this function
// allows you to say:
//  value, err := cp.Get("multipath-policy")
//  if err != nil {
//    ... // handle runtime error
//  }
//  policy, ok := value.(MultipathPolicy)
//  if !ok {
//    ... // handle failed type assertion
//  }
//
// In idiomatic Go, you would (and should) instead say:
//  policy := cp.MultipathPolicy
//
// Deprecated: Use func cp.Get only if you must. Direct access
// of the underlying ConnectionProperties struct Fields is usually
// preferred. This function is implemented using reflection and
// dynamic string matching, which is inherently inefficient and prone
// to bugs triggered at runtime.
func (cp *ConnectionProperties) Get(property string) (value interface{}, err error) {
	v, err := get(cp, property)
	return v.Interface(), err
}

// Connection
type Connection struct {
	Ready      chan bool
	SoftError  chan error
	PathChange chan bool
	cp         *ConnectionProperties
}

// GetProperties can be called at any time by the application to query ConnectionProperties
//
// Deprecated: Per https://go.dev/doc/effective_go#Getters, it is not
// idiomatic Go to put "Get" into a getter's name
func (c *Connection) GetProperties() *ConnectionProperties {
	return c.Properties()
}

// Properties can be called at any time by the application to query ConnectionProperties
func (c *Connection) Properties() *ConnectionProperties {
	return c.cp
}

// SetProperty stores value for property, which is stripped of case
// and non-alphabetic characters before being matched against the
// (equally stripped) exported Field names of c. The type of value
// must be assignable to type of the targeted property Field,
// otherwise an error is returned.
//
// For the sake of respecting the TAPS (draft) spec as closely as
// possible, this function allows you to say:
//  err := c.SetProperty("connPrio", 100)
//  if err != nil {
//    ... // handle runtime error
//  }
//
// In idiomatic Go, you would (and should) instead say:
//  c.Properties().ConnPrio = 100
//
// Deprecated: Use func c.SetProperty only if you must. Direct access
// of the underlying ConnectionProperties struct Fields is usually
// preferred. This function is implemented using reflection and
// dynamic string matching, which is inherently inefficient and prone
// to bugs triggered at runtime.
func (c *Connection) SetProperty(property string, value interface{}) error {
	return set(c.cp, property, value)
}
