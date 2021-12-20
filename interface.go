package panapi

import (
	"crypto"
	"crypto/tls"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
)

type Preference uint8

const (
	// (Implementation detail: Indicate that recommended default
	// value for Property should be used)
	unset Preference = iota

	// No preference
	Ignore

	// Select only protocols/paths providing the property, fail
	// otherwise
	Require

	// Prefer protocols/paths providing the property, proceed
	// otherwise
	Prefer

	// Prefer protocols/paths not providing the property, proceed
	// otherwise
	Avoid

	// Select only protocols/paths not providing the property,
	// fail otherwise
	Prohibit
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

type MultipathPreference uint8

const (
	// (Implementation detail: need to use different defaults
	// depending on endpoint)
	dynamic MultipathPreference = iota

	// The connection will not use multiple paths once
	// established, even if the chosen transport supports using
	// multiple paths.
	Disabled

	// The connection will negotiate the use of multiple paths if
	// the chosen transport supports this.
	Active

	// The connection will support the use of multiple paths if
	// the Remote Endpoint requests it.
	Passive
)

func (p MultipathPreference) String() string {
	return [...]string{
		"(unset)",
		"Disabled",
		"Active",
		"Passive",
	}[p-dynamic]
}

type MultipathPolicy uint8

const (
	// The connection ought only to attempt to migrate between
	// different paths when the original path is lost or becomes
	// unusable.
	Handover MultipathPolicy = iota

	// The connection ought only to attempt to minimize the
	// latency for interactive traffic patterns by transmitting
	// data across multiple paths when this is beneficial. The
	// goal of minimizing the latency will be balanced against the
	// cost of each of these paths. Depending on the cost of the
	// lower-latency path, the scheduling might choose to use a
	// higher-latency path. Traffic can be scheduled such that
	// data may be transmitted on multiple paths in parallel to
	// achieve a lower latency.
	Interactive

	// The connection ought to attempt to use multiple paths in
	// parallel to maximize available capacity and possibly
	// overcome the capacity limitations of the individual paths.
	Aggregate
)

type Directionality uint8

const (
	// The connection must support sending and receiving data
	Bidirectional Directionality = iota

	// The connection must support sending data, and the application cannot use the connection to receive any data
	UnidirectionalSend

	// The connection must support receiving data, and the application cannot use the connection to send any data
	UnidirectionalReceive
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

func set(st interface{}, key string, value interface{}) error {
	s := reflect.ValueOf(st).Elem()
	reg := regexp.MustCompile("[^a-z]+")
	stripKey := reg.ReplaceAllString(strings.ToLower(key), "")
	f := s.FieldByNameFunc(func(k string) bool {
		if reg.ReplaceAllString(strings.ToLower(k), "") == stripKey {
			return true
		} else {
			return false
		}
	})
	if !f.IsValid() {
		return fmt.Errorf("Type %T has no Field %s (%s)", st, key, stripKey)
	}
	p := reflect.ValueOf(value)
	if p.IsValid() && p.Type().AssignableTo(f.Type()) {
		f.Set(p)
	} else {
		return fmt.Errorf("Can not assign value of Type %T to Field %s of Type %T (expect %s)", value, key, st, f.Type())
	}
	return nil

}

// Set stores value for property, which is stripped of case and
// non-alphabetic characters before being matched against the (equally
// stripped) exported Field names of tp. The type of value must be
// assignable to type of the targeted property Field, otherwise an
// error is returned.
//
// For the sake of respecting the TAPS (draft) spec as closely as
// possible, this function allows you to instead say:
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

// Require has the effect of selecting only protocols/paths providing the property and failing otherwise.
//
// It is equivalent to calling tp.Set(property, Require) - the caveats of func tp.Set apply in full
func (tp *TransportProperties) Require(property string) error {
	return tp.Set(property, Require)
}

// Prefer has the effect of prefering protocols/paths providing the property and proceeding otherwise
//
// It is equivalent to calling tp.Set(property, Prefer) - the caveats of func tp.Set apply in full
func (tp *TransportProperties) Prefer(property string) error {
	return tp.Set(property, Prefer)
}

// Ignore has the effect of expressing no preference for the given property
//
// It is equivalent to calling tp.Set(property, Ignore) - the caveats of func tp.Set apply in full
func (tp *TransportProperties) Ignore(property string) error {
	return tp.Set(property, Ignore)

}

// Avoid has the effect of avoiding protocols/paths with the property and proceeding otherwise
//
// It is equivalent to calling tp.Set(property, Avoid) - the caveats of func tp.Set apply in full
func (tp *TransportProperties) Avoid(property string) error {
	return tp.Set(property, Avoid)
}

// Prohibit has the effect of failing if the property can not be avoided
//
// It is equivalent to calling tp.Set(property, Prohibit) - the caveats of func tp.Set apply in full
func (tp *TransportProperties) Prohibit(property string) error {
	return tp.Set(property, Prohibit)
}

// SecurityParameters is a structure used to configure security for a Preconnection
type SecurityParameters struct {
	// Local identity and private keys: Used to perform private key operations and prove one's identity to the Remote Endpoint.
	Identity string
	KeyPair  KeyPair

	// Supported algorithms: Used to restrict what parameters are
	// used by underlying transport security protocols. When not
	// specified, these algorithms should use known and safe
	// defaults for the system. Parameters include: ciphersuites,
	// supported groups, and signature algorithms. These
	// parameters take a collection of supported algorithms as
	// parameter.
	SupportedGroup        tls.CurveID
	CipherSuite           *tls.CipherSuite
	SignatureAlgorithm    tls.SignatureScheme
	MaxCachedSessions     uint
	CachedSessionLifetime time.Duration

	// Unsupported
	PSK crypto.PrivateKey
}

// NewSecurityParameters TODO
func NewSecurityParameters() *SecurityParameters {
	return &SecurityParameters{}
}

// NewDisabledSecurityParameters is intended for compatibility with
// endpoints that do not support transport security protocols (such as
// TCP without support for TLS)
func NewDisabledSecurityParameters() *SecurityParameters {
	return &SecurityParameters{}
}

// NewOpportunisticSecurityParameters() is not yet implemented
func NewOpportunisticSecurityParameters() *SecurityParameters {
	return &SecurityParameters{}
}

// SetTrustVerificationCallback is not yet implemented
func (sp SecurityParameters) SetTrustVerificationCallback() {
}

// SetIdentityChallengeCallback is not yet implemented
func (sp SecurityParameters) SetIdentityChallengeCallback() {
}

// Set stores value for parameter, which is stripped of case and
// non-alphabetic characters before being matched against the (equally
// stripped) exported Field names of sp. The type of value must be
// assignable to type of the targeted parameter Field, otherwise an
// error is returned.
//
// For the sake of respecting the TAPS (draft) spec as closely as
// possible, this function allows you to instead say:
//  err := sp.Set("supported-group", tls.CurveP521)
//  if err != nil {
//    ... // handle runtime error
//  }
//
// In idiomatic Go, you would (and should) instead say:
//  sp.SupportedGroup = tls.CurveP521
//
// Deprecated: Use func sp.Set only if you must. Direct access of the
// SecurityParameters struct Fields is usually preferred. This
// function is implemented using reflection and dynamic string
// matching, which is inherently inefficient and prone to bugs
// triggered at runtime.
func (sp *SecurityParameters) Set(parameter string, value interface{}) error {
	return set(sp, parameter, value)
}

type CapacityProfile uint8

const (
	// The application provides no information about its expected
	// capacity profile.
	Default CapacityProfile = iota

	// The application is not interactive. It expects to send
	// and/or receive data without any urgency. This can, for
	// example, be used to select protocol stacks with scavenger
	// transmission control and/or to assign the traffic to a
	// lower-effort service.
	Scavenger

	// The application is interactive, and prefers loss to
	// latency. Response time should be optimized at the expense
	// of delay variation and efficient use of the available
	// capacity when sending on this connection. This can be used
	// by the system to disable the coalescing of multiple small
	// Messages into larger packets (Nagle's algorithm); to prefer
	// immediate acknowledgment from the peer endpoint when
	// supported by the underlying transport; and so on.
	LowLatencyInteractive

	// The application prefers loss to latency, but is not
	// interactive. Response time should be optimized at the
	// expense of delay variation and efficient use of the
	// available capacity when sending on this connection.
	LowLatencyNonInteractive

	// The application expects to send/receive data at a constant
	// rate after Connection establishment. Delay and delay
	// variation should be minimized at the expense of efficient
	// use of the available capacity. This implies that the
	// Connection might fail if the Path is unable to maintain the
	// desired rate.
	ConstantRateStreaming

	// The application expects to send/receive data at the maximum
	// rate allowed by its congestion controller over a relatively
	// long period of time.
	CapacitySeeking
)

func (p CapacityProfile) String() string {
	return [...]string{
		"Default",
		"Scavenger",
		"Low Latency/Interactive",
		"Low Latency/Non-Interactive",
		"Constant-Rate Streaming",
		"Capacity-Seeking",
	}[p-Default]
}

type StreamScheduler uint8

const (
	SCTP_SS_FCFS   StreamScheduler = iota // First-Come, First-Served Scheduler
	SCTP_SS_RR                            // Round-Robin Scheduler
	SCTP_SS_RR_PKT                        // Round-Robin Scheduler per Packet
	SCTP_SS_PRIO                          // Priority-Based Scheduler
	SCTP_SS_FC                            // Fair Capacity Scheduler
	SCTP_SS_WFQ                           // Weighted Fair Queueing Scheduler
)

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
	// property. (See https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-8.1.2)
	ConnPrio uint

	// ConnTimeout specifies how long to wait before deciding that
	// an active Connection has failed when trying to reliably
	// deliver data to the Remote Endpoint. Adjusting this
	// Property will only take effect when the underlying stack
	// supports reliability. A value of 0 means that no timeout is
	// scheduled. (See https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-8.1.3)
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

	zeroRTTMsgMaxLen, singularTransmissionMsgMaxLen, sendMsgMaxLen, recMsgMaxLen uint
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
	return cp.recMsgMaxLen
}

// KeyPair clearly associates a Private and Public Key into a pair
type KeyPair struct {
	PrivateKey crypto.PrivateKey
	PublicKey  crypto.PublicKey
}

//
type Connection struct {
	Ready      chan bool
	SoftError  chan error
	PathChange chan bool
}
