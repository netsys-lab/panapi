package taps

import (
	"github.com/netsys-lab/panapi/internal/enum"
)

type SelectionProperties struct {
	// Reliability pecifies whether the application needs to use a
	// transport protocol that ensures that all data is received
	// at the Remote Endpoint without corruption. When reliable
	// data transfer is enabled, this also entails being notified
	// when a Connection is closed or aborted. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2.1)
	Reliability enum.Preference

	// PreserveMsgBoundaries specifies whether the application
	// needs or prefers to use a transport protocol that preserves
	// message boundaries. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2.2)
	PreserveMsgBoundaries enum.Preference

	// PerMsgReliability specifies whether an application
	// considers it useful to specify different reliability
	// requirements for individual Messages in a Connection. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2.3)
	PerMsgReliability enum.Preference

	// PreserveOrder specifies whether the application wishes to
	// use a transport protocol that can ensure that data is
	// received by the application on the other end in the same
	// order as it was sent. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2.4)
	PreserveOrder enum.Preference

	// ZeroRTTMsg specifies whether an application would like to
	// supply a Message to the transport protocol before
	// Connection establishment that will then be reliably
	// transferred to the other side before or during Connection
	// establishment. This Message can potentially be received
	// multiple times (i.e., multiple copies of the message data
	// may be passed to the Remote Endpoint). (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2.5)
	ZeroRTTMsg enum.Preference

	// Multistreaming specifies that the application would prefer
	// multiple Connections within a Connection Group to be
	// provided by streams of a single underlying transport
	// connection where possible.  (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2.6)
	Multistreaming enum.Preference

	// FullChecksumSend specifies the application's need for
	// protection against corruption for all data transmitted on
	// this Connection. Disabling this property could enable later
	// control of the sender checksum coverage. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2.7)
	FullChecksumSend enum.Preference

	// FullChecksumRecv specifies the application's need for
	// protection against corruption for all data received on this
	// Connection. Disabling this property could enable later
	// control of the required minimum receiver checksum
	// coverage. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2.8)
	FullChecksumRecv enum.Preference

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
	CongestionControl enum.Preference

	// KeepAlive specifies whether the application would like the
	// Connection to send keep-alive packets or not. Note that if
	// a Connection determines that keep-alive packets are being
	// sent, the applicaton should itself avoid generating
	// additional keep alive messages. Note that when supported,
	// the system will use the default period for generation of
	// the keep alive-packets. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2.10)
	KeepAlive enum.Preference

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
	Interface map[string]enum.Preference

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
	PvD map[string]enum.Preference

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
	Multipath enum.MultipathPreference

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
	UseTemporaryLocalAddress enum.Preference

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
	Direction enum.Directionality

	// SoftErrorNotify specifies whether an application considers
	// it useful to be informed when an ICMP error message arrives
	// that does not force termination of a connection. When set
	// to true, received ICMP errors are available as
	// SoftErrors. Note that even if a protocol supporting this
	// property is selected, not all ICMP errors will necessarily
	// be delivered, so applications cannot rely upon receiving
	// them [RFC8085]. (See
	// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2.17)
	SoftErrorNotify enum.Preference

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
	ActiveReadBeforeSend enum.Preference
}

// NewSelectionProperties creates SelectionProperties with the
// recommended defaults from
// https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-6.2
func NewSelectionProperties() *SelectionProperties {
	return &SelectionProperties{
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
		Interface:             map[string]enum.Preference{},
		PvD:                   map[string]enum.Preference{},
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

// Copy returns a new SelectionProperties struct with its values deeply copied from sp
func (sp *SelectionProperties) Copy() *SelectionProperties {
	var (
		newInterface = make(map[string]enum.Preference)
		newPvD       = make(map[string]enum.Preference)
	)
	for key, value := range sp.Interface {
		newInterface[key] = value
	}
	for key, value := range sp.PvD {
		newPvD[key] = value
	}
	return &SelectionProperties{
		Reliability:              sp.Reliability,
		PreserveMsgBoundaries:    sp.PreserveMsgBoundaries,
		PerMsgReliability:        sp.PerMsgReliability,
		PreserveOrder:            sp.PreserveOrder,
		ZeroRTTMsg:               sp.ZeroRTTMsg,
		Multistreaming:           sp.Multistreaming,
		FullChecksumSend:         sp.FullChecksumSend,
		FullChecksumRecv:         sp.FullChecksumRecv,
		CongestionControl:        sp.CongestionControl,
		KeepAlive:                sp.KeepAlive,
		Interface:                newInterface,
		PvD:                      newPvD,
		Multipath:                sp.Multipath,
		UseTemporaryLocalAddress: sp.UseTemporaryLocalAddress,
		AdvertisesAltAddr:        sp.AdvertisesAltAddr,
		Direction:                sp.Direction,
		SoftErrorNotify:          sp.SoftErrorNotify,
		ActiveReadBeforeSend:     sp.ActiveReadBeforeSend,
	}
}
